package service

import (
	"context"

	"github.com/Alex808r/wallet/internal/domain"
	"github.com/Alex808r/wallet/internal/repository"
	internalsync "github.com/Alex808r/wallet/internal/sync"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// WalletService описывает контракт уровня бизнес-логики (Use-Cases) для кошельков.
type WalletService interface {
	ProcessOperation(ctx context.Context, req domain.WalletOperationRequest) (*domain.WalletResponse, error)
	GetBalance(ctx context.Context, walletID uuid.UUID) (*domain.WalletResponse, error)
	CreateWallet(ctx context.Context, walletID uuid.UUID, initialBalance decimal.Decimal) (*domain.WalletResponse, error)
}

// walletService реализует интерфейс WalletService с поддержкой конкурентных блокировок.
type walletService struct {
	repo   repository.WalletRepository // Слой доступа к базе данных
	locker *internalsync.KeyLocker     // Шардированный in-memory замок
}

// NewWalletService создает экземпляр сервиса с внедренными зависимостями.
func NewWalletService(repo repository.WalletRepository, locker *internalsync.KeyLocker) WalletService {
	if locker == nil {
		locker = internalsync.NewKeyLocker()
	}
	return &walletService{
		repo:   repo,
		locker: locker,
	}
}

// ProcessOperation выполняет пополнение (DEPOSIT) или списание (WITHDRAW) средств.
// -----------------------------------------------------------------------------
// Ключевая логика надежности:
// 1. Вычисляет математическую дельту (положительная для DEPOSIT, отрицательная для WITHDRAW).
// 2. Захватывает мьютекс кошелька s.locker.Lock(req.WalletID) ДО обращения к базе данных.
//    Благодаря этому запросы к одному кошельку не конкурируют за соединения из pgxpool,
//    а выполняются друг за другом за микросекунды.
// 3. Освобождает мьютекс через defer s.locker.Unlock(req.WalletID) при любом исходе.
// -----------------------------------------------------------------------------
func (s *walletService) ProcessOperation(ctx context.Context, req domain.WalletOperationRequest) (*domain.WalletResponse, error) {
	// Валидация входных данных
	if req.WalletID == uuid.Nil {
		return nil, domain.ErrInvalidWalletID
	}
	if req.Amount.LessThanOrEqual(decimal.Zero) {
		return nil, domain.ErrInvalidAmount
	}

	// Определение дельты изменения баланса
	var delta decimal.Decimal
	switch req.OperationType {
	case domain.OperationDeposit:
		delta = req.Amount
	case domain.OperationWithdraw:
		delta = req.Amount.Neg() // Отрицательное число для списания
	default:
		return nil, domain.ErrInvalidOperation
	}

	// Захват блокировки конкретного кошелька в оперативной памяти
	s.locker.Lock(req.WalletID)
	defer s.locker.Unlock(req.WalletID)

	// Выполнение атомарного обновления в PostgreSQL
	wallet, err := s.repo.UpdateBalance(ctx, req.WalletID, delta, req.OperationType)
	if err != nil {
		return nil, err
	}

	return &domain.WalletResponse{
		WalletID: wallet.ID,
		Balance:  wallet.Balance,
	}, nil
}

// GetBalance возвращает актуальный баланс кошелька по его идентификатору.
// Чтение не блокирует другие операции благодаря модели MVCC в PostgreSQL.
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

// CreateWallet инициализирует новый кошелек с начальным балансом (хелпер для тестов/онбординга).
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
