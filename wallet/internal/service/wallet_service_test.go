package service

import (
	"context"
	"testing"
	"time"

	"github.com/Alex808r/wallet/internal/domain"
	internalsync "github.com/Alex808r/wallet/internal/sync"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type mockWalletRepo struct {
	wallets map[uuid.UUID]*domain.Wallet
}

func newMockRepo() *mockWalletRepo {
	return &mockWalletRepo{
		wallets: make(map[uuid.UUID]*domain.Wallet),
	}
}

func (m *mockWalletRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Wallet, error) {
	w, ok := m.wallets[id]
	if !ok {
		return nil, domain.ErrWalletNotFound
	}
	return &domain.Wallet{ID: w.ID, Balance: w.Balance}, nil
}

func (m *mockWalletRepo) UpdateBalance(ctx context.Context, id uuid.UUID, delta decimal.Decimal, opType domain.OperationType) (*domain.Wallet, error) {
	w, ok := m.wallets[id]
	if !ok {
		return nil, domain.ErrWalletNotFound
	}

	newBalance := w.Balance.Add(delta)
	if newBalance.LessThan(decimal.Zero) {
		return nil, domain.ErrInsufficientFunds
	}

	w.Balance = newBalance
	w.UpdatedAt = time.Now()
	return &domain.Wallet{ID: w.ID, Balance: w.Balance}, nil
}

func (m *mockWalletRepo) Create(ctx context.Context, id uuid.UUID, initialBalance decimal.Decimal) (*domain.Wallet, error) {
	w := &domain.Wallet{
		ID:        id,
		Balance:   initialBalance,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	m.wallets[id] = w
	return w, nil
}

func TestWalletService_ProcessOperation(t *testing.T) {
	repo := newMockRepo()
	locker := internalsync.NewKeyLocker()
	svc := NewWalletService(repo, locker)
	ctx := context.Background()

	walletID := uuid.New()
	_, _ = repo.Create(ctx, walletID, decimal.NewFromInt(1000))

	// 1. Successful Deposit
	resp, err := svc.ProcessOperation(ctx, domain.WalletOperationRequest{
		WalletID:      walletID,
		OperationType: domain.OperationDeposit,
		Amount:        decimal.NewFromInt(500),
	})
	if err != nil {
		t.Fatalf("expected no error on deposit, got: %v", err)
	}
	if !resp.Balance.Equal(decimal.NewFromInt(1500)) {
		t.Errorf("expected balance 1500, got: %v", resp.Balance)
	}

	// 2. Successful Withdraw
	resp, err = svc.ProcessOperation(ctx, domain.WalletOperationRequest{
		WalletID:      walletID,
		OperationType: domain.OperationWithdraw,
		Amount:        decimal.NewFromInt(300),
	})
	if err != nil {
		t.Fatalf("expected no error on withdraw, got: %v", err)
	}
	if !resp.Balance.Equal(decimal.NewFromInt(1200)) {
		t.Errorf("expected balance 1200, got: %v", resp.Balance)
	}

	// 3. Insufficient Funds Withdraw
	_, err = svc.ProcessOperation(ctx, domain.WalletOperationRequest{
		WalletID:      walletID,
		OperationType: domain.OperationWithdraw,
		Amount:        decimal.NewFromInt(5000),
	})
	if err != domain.ErrInsufficientFunds {
		t.Errorf("expected ErrInsufficientFunds, got: %v", err)
	}

	// 4. Non-existent wallet
	_, err = svc.ProcessOperation(ctx, domain.WalletOperationRequest{
		WalletID:      uuid.New(),
		OperationType: domain.OperationDeposit,
		Amount:        decimal.NewFromInt(100),
	})
	if err != domain.ErrWalletNotFound {
		t.Errorf("expected ErrWalletNotFound, got: %v", err)
	}
}

func TestWalletService_GetBalance(t *testing.T) {
	repo := newMockRepo()
	svc := NewWalletService(repo, nil)
	ctx := context.Background()

	walletID := uuid.New()
	_, _ = repo.Create(ctx, walletID, decimal.NewFromInt(2500))

	resp, err := svc.GetBalance(ctx, walletID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.Balance.Equal(decimal.NewFromInt(2500)) {
		t.Errorf("expected balance 2500, got: %v", resp.Balance)
	}

	// Not found
	_, err = svc.GetBalance(ctx, uuid.New())
	if err != domain.ErrWalletNotFound {
		t.Errorf("expected ErrWalletNotFound, got: %v", err)
	}
}
