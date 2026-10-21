package repository

import (
	"context"

	"github.com/Alex808r/wallet/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// WalletRepository defines data access methods for wallets.
type WalletRepository interface {
	// GetByID retrieves a wallet by its ID. Returns domain.ErrWalletNotFound if not found.
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Wallet, error)

	// UpdateBalance atomically adjusts the wallet balance (delta > 0 for DEPOSIT, delta < 0 for WITHDRAW).
	// Returns updated wallet, domain.ErrWalletNotFound, or domain.ErrInsufficientFunds.
	UpdateBalance(ctx context.Context, id uuid.UUID, delta decimal.Decimal, opType domain.OperationType) (*domain.Wallet, error)

	// Create creates a new wallet with initial balance.
	Create(ctx context.Context, id uuid.UUID, initialBalance decimal.Decimal) (*domain.Wallet, error)
}
