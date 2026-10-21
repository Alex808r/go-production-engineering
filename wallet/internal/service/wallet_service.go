package service

import (
	"context"

	"github.com/Alex808r/wallet/internal/domain"
	"github.com/Alex808r/wallet/internal/repository"
	internalsync "github.com/Alex808r/wallet/internal/sync"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// WalletService defines high-level business operations for wallets.
type WalletService interface {
	ProcessOperation(ctx context.Context, req domain.WalletOperationRequest) (*domain.WalletResponse, error)
	GetBalance(ctx context.Context, walletID uuid.UUID) (*domain.WalletResponse, error)
	CreateWallet(ctx context.Context, walletID uuid.UUID, initialBalance decimal.Decimal) (*domain.WalletResponse, error)
}

type walletService struct {
	repo   repository.WalletRepository
	locker *internalsync.KeyLocker
}

// NewWalletService creates a new WalletService instance with keyed lock protection.
func NewWalletService(repo repository.WalletRepository, locker *internalsync.KeyLocker) WalletService {
	if locker == nil {
		locker = internalsync.NewKeyLocker()
	}
	return &walletService{
		repo:   repo,
		locker: locker,
	}
}

// ProcessOperation handles DEPOSIT and WITHDRAW operations with strict serialization on the same wallet.
func (s *walletService) ProcessOperation(ctx context.Context, req domain.WalletOperationRequest) (*domain.WalletResponse, error) {
	if req.WalletID == uuid.Nil {
		return nil, domain.ErrInvalidWalletID
	}
	if req.Amount.LessThanOrEqual(decimal.Zero) {
		return nil, domain.ErrInvalidAmount
	}

	var delta decimal.Decimal
	switch req.OperationType {
	case domain.OperationDeposit:
		delta = req.Amount
	case domain.OperationWithdraw:
		delta = req.Amount.Neg()
	default:
		return nil, domain.ErrInvalidOperation
	}

	// Lock wallet in-memory to prevent database connection starvation
	// and lock contention in Postgres under 1000 RPS.
	s.locker.Lock(req.WalletID)
	defer s.locker.Unlock(req.WalletID)

	wallet, err := s.repo.UpdateBalance(ctx, req.WalletID, delta, req.OperationType)
	if err != nil {
		return nil, err
	}

	return &domain.WalletResponse{
		WalletID: wallet.ID,
		Balance:  wallet.Balance,
	}, nil
}

// GetBalance retrieves current wallet balance.
func (s *walletService) GetBalance(ctx context.Context, walletID uuid.UUID) (*domain.WalletResponse, error) {
	if walletID == uuid.Nil {
		return nil, domain.ErrInvalidWalletID
	}

	wallet, err := s.repo.GetByID(ctx, walletID)
	if err != nil {
		return nil, err
	}

	return &domain.WalletResponse{
		WalletID: wallet.ID,
		Balance:  wallet.Balance,
	}, nil
}

// CreateWallet initializes a new wallet with given balance.
func (s *walletService) CreateWallet(ctx context.Context, walletID uuid.UUID, initialBalance decimal.Decimal) (*domain.WalletResponse, error) {
	if walletID == uuid.Nil {
		return nil, domain.ErrInvalidWalletID
	}
	if initialBalance.LessThan(decimal.Zero) {
		return nil, domain.ErrInvalidAmount
	}

	s.locker.Lock(walletID)
	defer s.locker.Unlock(walletID)

	wallet, err := s.repo.Create(ctx, walletID, initialBalance)
	if err != nil {
		return nil, err
	}

	return &domain.WalletResponse{
		WalletID: wallet.ID,
		Balance:  wallet.Balance,
	}, nil
}
