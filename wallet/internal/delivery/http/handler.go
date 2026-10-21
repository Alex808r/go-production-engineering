package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/Alex808r/wallet/internal/domain"
	"github.com/Alex808r/wallet/internal/service"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// WalletHandler обрабатывает входящие HTTP-запросы REST API.
type WalletHandler struct {
	service service.WalletService
}

// NewWalletHandler создает новый экземпляр HTTP-обработчика.
func NewWalletHandler(svc service.WalletService) *WalletHandler {
	return &WalletHandler{service: svc}
}

// HandleOperation обрабатывает POST /api/v1/wallet (пополнение или списание).
// -----------------------------------------------------------------------------
// Безопасность и обработка:
// 1. Ограничение размера тела запроса (1 МБ) через io.LimitReader для защиты от DoS-атак переполнением памяти.
// 2. Десериализация через кастомный UnmarshalJSON, поддерживающий valletId и walletId.
// 3. Маппинг доменных ошибок на стандартизированные HTTP-коды:
//    - 400 Bad Request: недостаток средств, невалидный UUID, неверный тип операции, сумма <= 0.
//    - 404 Not Found: кошелек не найден.
//    - 200 OK: операция успешно зафиксирована.
// -----------------------------------------------------------------------------
func (h *WalletHandler) HandleOperation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Читаем тело с жестким лимитом 1 МБ
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		writeError(w, http.StatusBadRequest, "failed to read request body")
		return
	}
	defer r.Body.Close()

	// Парсинг и базовая валидация JSON
	var req domain.WalletOperationRequest
	if err := json.Unmarshal(body, &req); err != nil {
		if errors.Is(err, domain.ErrInvalidWalletID) ||
			errors.Is(err, domain.ErrInvalidOperation) ||
			errors.Is(err, domain.ErrInvalidAmount) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusBadRequest, "invalid json payload: "+err.Error())
		return
	}

	// Выполнение операции на уровне сервиса
	resp, err := h.service.ProcessOperation(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrWalletNotFound):
			writeError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, domain.ErrInsufficientFunds),
			errors.Is(err, domain.ErrInvalidOperation),
			errors.Is(err, domain.ErrInvalidAmount),
			errors.Is(err, domain.ErrInvalidWalletID):
			writeError(w, http.StatusBadRequest, err.Error())
		default:
			// Внутренняя непредвиденная ошибка
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// HandleGetBalance обрабатывает GET /api/v1/wallets/{WALLET_UUID}.
// -----------------------------------------------------------------------------
// Извлекает UUID кошелька из пути запроса, валидирует формат и возвращает текущий баланс.
// -----------------------------------------------------------------------------
func (h *WalletHandler) HandleGetBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Извлечение параметра пути через встроенный функционал Go 1.22+ PathValue
	idStr := r.PathValue("id")
	if idStr == "" {
		// Резервный разбор пути для обратной совместимости
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) > 0 {
			idStr = parts[len(parts)-1]
		}
	}

	walletID, err := uuid.Parse(strings.TrimSpace(idStr))
	if err != nil || walletID == uuid.Nil {
		writeError(w, http.StatusBadRequest, "invalid wallet id format")
		return
	}

	resp, err := h.service.GetBalance(r.Context(), walletID)
	if err != nil {
		if errors.Is(err, domain.ErrWalletNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// HandleCreateWallet обрабатывает POST /api/v1/wallets.
// Вспомогательный эндпоинт для начального создания кошельков в автотестах.
func (h *WalletHandler) HandleCreateWallet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	type createReq struct {
		WalletID       *uuid.UUID       `json:"walletId"`
		InitialBalance *decimal.Decimal `json:"initialBalance"`
	}

	var req createReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil && err != io.EOF {
		writeError(w, http.StatusBadRequest, "invalid json: "+err.Error())
		return
	}

	targetID := uuid.New()
	if req.WalletID != nil && *req.WalletID != uuid.Nil {
		targetID = *req.WalletID
	}

	balance := decimal.Zero
	if req.InitialBalance != nil {
		balance = *req.InitialBalance
	}

	resp, err := h.service.CreateWallet(r.Context(), targetID, balance)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidAmount) || errors.Is(err, domain.ErrInvalidWalletID) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

// HandleHealth возвращает статус доступности сервиса для Docker / Kubernetes healthcheck.
func (h *WalletHandler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// writeJSON сериализует ответ в JSON с установкой заголовка Content-Type.
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// writeError формирует унифицированный JSON с сообщением об ошибке.
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
