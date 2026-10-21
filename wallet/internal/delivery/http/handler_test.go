package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Alex808r/wallet/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type mockService struct {
	wallets map[uuid.UUID]decimal.Decimal
}

func newMockService() *mockService {
	return &mockService{
		wallets: make(map[uuid.UUID]decimal.Decimal),
	}
}

func (m *mockService) ProcessOperation(ctx context.Context, req domain.WalletOperationRequest) (*domain.WalletResponse, error) {
	bal, ok := m.wallets[req.WalletID]
	if !ok {
		return nil, domain.ErrWalletNotFound
	}

	if req.OperationType == domain.OperationWithdraw {
		if bal.LessThan(req.Amount) {
			return nil, domain.ErrInsufficientFunds
		}
		bal = bal.Sub(req.Amount)
	} else if req.OperationType == domain.OperationDeposit {
		bal = bal.Add(req.Amount)
	} else {
		return nil, domain.ErrInvalidOperation
	}

	m.wallets[req.WalletID] = bal
	return &domain.WalletResponse{
		WalletID: req.WalletID,
		Balance:  bal,
	}, nil
}

func (m *mockService) GetBalance(ctx context.Context, walletID uuid.UUID) (*domain.WalletResponse, error) {
	bal, ok := m.wallets[walletID]
	if !ok {
		return nil, domain.ErrWalletNotFound
	}
	return &domain.WalletResponse{
		WalletID: walletID,
		Balance:  bal,
	}, nil
}

func (m *mockService) CreateWallet(ctx context.Context, walletID uuid.UUID, initialBalance decimal.Decimal) (*domain.WalletResponse, error) {
	m.wallets[walletID] = initialBalance
	return &domain.WalletResponse{
		WalletID: walletID,
		Balance:  initialBalance,
	}, nil
}

func TestHandler_PostOperation(t *testing.T) {
	svc := newMockService()
	h := NewWalletHandler(svc)
	router := NewRouter(h)

	walletID := uuid.New()
	svc.wallets[walletID] = decimal.NewFromInt(1000)

	// 1. Successful Deposit using "valletId" (from spec typo)
	body1 := []byte(`{
		"valletId": "` + walletID.String() + `",
		"operationType": "DEPOSIT",
		"amount": 500
	}`)
	req1 := httptest.NewRequest(http.MethodPost, "/api/v1/wallet", bytes.NewReader(body1))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)

	if w1.Code != http.StatusOK {
		t.Fatalf("expected status 200, got: %d (%s)", w1.Code, w1.Body.String())
	}

	var resp1 domain.WalletResponse
	if err := json.Unmarshal(w1.Body.Bytes(), &resp1); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if !resp1.Balance.Equal(decimal.NewFromInt(1500)) {
		t.Errorf("expected balance 1500, got: %v", resp1.Balance)
	}

	// 2. Successful Withdraw using "walletId"
	body2 := []byte(`{
		"walletId": "` + walletID.String() + `",
		"operationType": "WITHDRAW",
		"amount": 200
	}`)
	req2 := httptest.NewRequest(http.MethodPost, "/api/v1/wallet", bytes.NewReader(body2))
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Fatalf("expected status 200, got: %d (%s)", w2.Code, w2.Body.String())
	}

	// 3. Insufficient Funds (returns 400 Bad Request)
	body3 := []byte(`{
		"valletId": "` + walletID.String() + `",
		"operationType": "WITHDRAW",
		"amount": 999999
	}`)
	req3 := httptest.NewRequest(http.MethodPost, "/api/v1/wallet", bytes.NewReader(body3))
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)

	if w3.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for insufficient funds, got: %d (%s)", w3.Code, w3.Body.String())
	}

	// 4. Non-existent wallet (returns 404 Not Found)
	unknownID := uuid.New()
	body4 := []byte(`{
		"valletId": "` + unknownID.String() + `",
		"operationType": "DEPOSIT",
		"amount": 100
	}`)
	req4 := httptest.NewRequest(http.MethodPost, "/api/v1/wallet", bytes.NewReader(body4))
	w4 := httptest.NewRecorder()
	router.ServeHTTP(w4, req4)

	if w4.Code != http.StatusNotFound {
		t.Errorf("expected status 404 for unknown wallet, got: %d (%s)", w4.Code, w4.Body.String())
	}

	// 5. Invalid operation type
	body5 := []byte(`{
		"valletId": "` + walletID.String() + `",
		"operationType": "UNKNOWN",
		"amount": 100
	}`)
	req5 := httptest.NewRequest(http.MethodPost, "/api/v1/wallet", bytes.NewReader(body5))
	w5 := httptest.NewRecorder()
	router.ServeHTTP(w5, req5)

	if w5.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for unknown op, got: %d (%s)", w5.Code, w5.Body.String())
	}

	// 6. Invalid amount (negative)
	body6 := []byte(`{
		"valletId": "` + walletID.String() + `",
		"operationType": "DEPOSIT",
		"amount": -50
	}`)
	req6 := httptest.NewRequest(http.MethodPost, "/api/v1/wallet", bytes.NewReader(body6))
	w6 := httptest.NewRecorder()
	router.ServeHTTP(w6, req6)

	if w6.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for negative amount, got: %d (%s)", w6.Code, w6.Body.String())
	}
}

func TestHandler_GetBalance(t *testing.T) {
	svc := newMockService()
	h := NewWalletHandler(svc)
	router := NewRouter(h)

	walletID := uuid.New()
	svc.wallets[walletID] = decimal.NewFromInt(3000)

	// 1. Success
	req := httptest.NewRequest(http.MethodGet, "/api/v1/wallets/"+walletID.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got: %d (%s)", w.Code, w.Body.String())
	}

	var resp domain.WalletResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if !resp.Balance.Equal(decimal.NewFromInt(3000)) {
		t.Errorf("expected balance 3000, got: %v", resp.Balance)
	}

	// 2. Not found
	req2 := httptest.NewRequest(http.MethodGet, "/api/v1/wallets/"+uuid.New().String(), nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusNotFound {
		t.Errorf("expected 404, got: %d", w2.Code)
	}

	// 3. Invalid UUID format
	req3 := httptest.NewRequest(http.MethodGet, "/api/v1/wallets/not-a-valid-uuid", nil)
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)

	if w3.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got: %d", w3.Code)
	}
}

func TestHandler_Health(t *testing.T) {
	svc := newMockService()
	h := NewWalletHandler(svc)
	router := NewRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got: %d", w.Code)
	}
}
