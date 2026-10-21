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

type WalletHandler struct {
	service service.WalletService
}

func NewWalletHandler(svc service.WalletService) *WalletHandler {
	return &WalletHandler{service: svc}
}

// HandleOperation processes POST /api/v1/wallet.
func (h *WalletHandler) HandleOperation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20)) // 1MB max
	if err != nil {
		writeError(w, http.StatusBadRequest, "failed to read request body")
		return
	}
	defer r.Body.Close()

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
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// HandleGetBalance processes GET /api/v1/wallets/{WALLET_UUID}.
func (h *WalletHandler) HandleGetBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Extract wallet ID from URL path
	idStr := r.PathValue("id")
	if idStr == "" {
		// Fallback for paths like /api/v1/wallets/{id}
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

// HandleCreateWallet processes POST /api/v1/wallets (convenience endpoint for creating wallets).
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

// HandleHealth returns system health status.
func (h *WalletHandler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
