package postgres

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/Alex808r/wallet/config"
	"github.com/Alex808r/wallet/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

// maxRetries задает максимальное число повторных попыток при транзиентных сбоях БД.
const maxRetries = 5

// WalletRepo реализует интерфейс repository.WalletRepository на базе PostgreSQL.
type WalletRepo struct {
	pool *pgxpool.Pool // Пул постоянных соединений к базе данных
}

// NewWalletRepo инициализирует и проверяет пул соединений pgxpool.
func NewWalletRepo(ctx context.Context, cfg *config.Config) (*WalletRepo, error) {
	// Парсинг строки подключения DSN (Data Source Name)
	poolConfig, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to parse database DSN: %w", err)
	}

	// --------------------------------------------------------------------------
	// Тюнинг пула соединений pgxpool:
	// MaxConns: максимальное число одновременных открытых соединений (100).
	// MinConns: пул "горячих" готовых соединений (20), исключающий задержки на handshake.
	// MaxConnIdleTime: время закрытия неиспользуемых соединений.
	// MaxConnLifetime: принудительная ротация соединений для сброса серверных ресурсов.
	// --------------------------------------------------------------------------
	poolConfig.MaxConns = cfg.DBMaxConns
	poolConfig.MinConns = cfg.DBMinConns
	poolConfig.MaxConnIdleTime = cfg.DBMaxIdleTime
	poolConfig.MaxConnLifetime = cfg.DBMaxConnLifetime

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Проверяем физическую доступность базы данных перед стартом сервиса
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &WalletRepo{pool: pool}, nil
}

// Close безопасно закрывает все соединения в пуле pgxpool.
func (r *WalletRepo) Close() {
	if r.pool != nil {
		r.pool.Close()
	}
}

// Pool возвращает экземпляр пула для тестов и проверок состояния.
func (r *WalletRepo) Pool() *pgxpool.Pool {
	return r.pool
}

// AutoMigrate автоматически применяет DDL-схему при первом запуске сервиса.
func (r *WalletRepo) AutoMigrate(ctx context.Context) error {
	query := `
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

	-- Таблица кошельков
	-- CHECK (balance >= 0) — ключевая защита на уровне ядра СУБД от отрицательного баланса
	CREATE TABLE IF NOT EXISTS wallets (
		id UUID PRIMARY KEY,
		balance NUMERIC(18, 2) NOT NULL DEFAULT 0.00 CHECK (balance >= 0),
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	CREATE INDEX IF NOT EXISTS idx_wallets_id ON wallets(id);

	-- Таблица аудита всех операций (для финансовой прозрачности и истории)
	CREATE TABLE IF NOT EXISTS wallet_transactions (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		wallet_id UUID NOT NULL REFERENCES wallets(id) ON DELETE CASCADE,
		operation_type VARCHAR(20) NOT NULL,
		amount NUMERIC(18, 2) NOT NULL,
		balance_after NUMERIC(18, 2) NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	CREATE INDEX IF NOT EXISTS idx_wallet_transactions_wallet_id ON wallet_transactions(wallet_id);
	`
	_, err := r.pool.Exec(ctx, query)
	return err
}

// GetByID возвращает кошелек по первичному ключу UUID.
func (r *WalletRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Wallet, error) {
	query := `
		SELECT id, balance, created_at, updated_at
		FROM wallets
		WHERE id = $1
	`

	var w domain.Wallet
	err := r.pool.QueryRow(ctx, query, id).Scan(&w.ID, &w.Balance, &w.CreatedAt, &w.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrWalletNotFound
		}
		return nil, fmt.Errorf("failed to query wallet: %w", err)
	}

	return &w, nil
}

// Create создает новый кошелек с заданным начальным балансом.
func (r *WalletRepo) Create(ctx context.Context, id uuid.UUID, initialBalance decimal.Decimal) (*domain.Wallet, error) {
	if initialBalance.LessThan(decimal.Zero) {
		return nil, domain.ErrInvalidAmount
	}

	query := `
		INSERT INTO wallets (id, balance, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		ON CONFLICT (id) DO UPDATE
		SET updated_at = NOW()
		RETURNING id, balance, created_at, updated_at
	`

	var w domain.Wallet
	err := r.pool.QueryRow(ctx, query, id, initialBalance).Scan(&w.ID, &w.Balance, &w.CreatedAt, &w.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create wallet: %w", err)
	}

	return &w, nil
}

// UpdateBalance выполняет атомарное изменение баланса с retry-политикой при временных сбоях.
// -----------------------------------------------------------------------------
// Retry-механизм:
// Если база данных вернула транзиентную ошибку (deadlock 40P01 или serialization failure 40001),
// запрос не завершается ошибкой 500, а повторяется до 5 раз с экспоненциальной задержкой.
// -----------------------------------------------------------------------------
func (r *WalletRepo) UpdateBalance(ctx context.Context, id uuid.UUID, delta decimal.Decimal, opType domain.OperationType) (*domain.Wallet, error) {
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		wallet, err := r.executeUpdate(ctx, id, delta, opType)
		if err == nil {
			return wallet, nil // Успех
		}

		// Доменные бизнес-ошибки повторять бессмысленно — сразу возвращаем клиенту
		if errors.Is(err, domain.ErrWalletNotFound) || errors.Is(err, domain.ErrInsufficientFunds) {
			return nil, err
		}

		// Проверяем, является ли ошибка временной ошибкой базы данных
		if isRetryable(err) {
			lastErr = err
			// Exponential Backoff с добавлением случайного джиттера (0..3 мс)
			// для предотвращения эффекта "толпы" (thundering herd)
			sleepTime := time.Duration(1<<attempt)*time.Millisecond + time.Duration(rand.Intn(3))*time.Millisecond
			select {
			case <-time.After(sleepTime):
				continue
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}

		// Неизвестная невосстановимая ошибка
		return nil, err
	}

	return nil, fmt.Errorf("update balance failed after %d retries: %w", maxRetries, lastErr)
}

// executeUpdate выполняет единую ACID-транзакцию изменения баланса и записи в аудит-лог.
func (r *WalletRepo) executeUpdate(ctx context.Context, id uuid.UUID, delta decimal.Decimal, opType domain.OperationType) (*domain.Wallet, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// --------------------------------------------------------------------------
	// АТОМАРНЫЙ SQL-ЗАПРОС:
	// 1. SET balance = balance + $1 — изменяет баланс (+ для DEPOSIT, - для WITHDRAW).
	// 2. WHERE id = $2 AND (balance + $1 >= 0) — проверяет, что баланс не станет отрицательным!
	// 3. RETURNING — возвращает обновленные значения за 1 сетевой раундтрип.
	// --------------------------------------------------------------------------
	updateQuery := `
		UPDATE wallets
		SET balance = balance + $1, updated_at = NOW()
		WHERE id = $2 AND (balance + $1 >= 0)
		RETURNING id, balance, created_at, updated_at
	`

	var w domain.Wallet
	err = tx.QueryRow(ctx, updateQuery, delta, id).Scan(&w.ID, &w.Balance, &w.CreatedAt, &w.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Если запрос ничего не обновил (pgx.ErrNoRows), выясняем точную причину:
			// Было ли это отсутствие кошелька или нехватка баланса?
			var existingID uuid.UUID
			checkErr := tx.QueryRow(ctx, `SELECT id FROM wallets WHERE id = $1`, id).Scan(&existingID)
			if errors.Is(checkErr, pgx.ErrNoRows) {
				return nil, domain.ErrWalletNotFound // Кошелька нет -> 404 Not Found
			}
			return nil, domain.ErrInsufficientFunds // Кошелек есть, но условие (balance + delta >= 0) ложно -> 400 Bad Request
		}
		return nil, err
	}

	// Запись факта операции в аудит-лог внутри той же транзакции
	insertAuditQuery := `
		INSERT INTO wallet_transactions (wallet_id, operation_type, amount, balance_after, created_at)
		VALUES ($1, $2, $3, $4, NOW())
	`
	_, _ = tx.Exec(ctx, insertAuditQuery, id, string(opType), delta.Abs(), w.Balance)

	// Фиксация транзакции
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &w, nil
}

// isRetryable проверяет коды ошибок PostgreSQL:
// 40001: serialization_failure (ошибка сериализации в конкурентной среде)
// 40P01: deadlock_detected (взаимная блокировка)
// 55P03: lock_not_available (таймаут захвата блокировки)
func isRetryable(err error) bool {
	if err == nil {
		return false
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "40001" || pgErr.Code == "40P01" || pgErr.Code == "55P03"
	}
	return false
}
