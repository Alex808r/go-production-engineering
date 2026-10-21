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

const maxRetries = 5

// WalletRepo implements repository.WalletRepository using PostgreSQL.
type WalletRepo struct {
	pool *pgxpool.Pool
}

// NewWalletRepo initializes a connection pool to PostgreSQL.
func NewWalletRepo(ctx context.Context, cfg *config.Config) (*WalletRepo, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to parse database DSN: %w", err)
	}

	poolConfig.MaxConns = cfg.DBMaxConns
	poolConfig.MinConns = cfg.DBMinConns
	poolConfig.MaxConnIdleTime = cfg.DBMaxIdleTime
	poolConfig.MaxConnLifetime = cfg.DBMaxConnLifetime

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Verify connection
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &WalletRepo{pool: pool}, nil
}

// Close closes the database connection pool.
func (r *WalletRepo) Close() {
	if r.pool != nil {
		r.pool.Close()
	}
}

// Pool returns the underlying pgxpool.Pool (for testing/healthchecks).
func (r *WalletRepo) Pool() *pgxpool.Pool {
	return r.pool
}

// AutoMigrate runs initial table creation and indexes.
func (r *WalletRepo) AutoMigrate(ctx context.Context) error {
	query := `
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

	CREATE TABLE IF NOT EXISTS wallets (
		id UUID PRIMARY KEY,
		balance NUMERIC(18, 2) NOT NULL DEFAULT 0.00 CHECK (balance >= 0),
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	CREATE INDEX IF NOT EXISTS idx_wallets_id ON wallets(id);

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

// GetByID retrieves a wallet by its UUID.
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

// Create creates a new wallet with initial balance.
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

// UpdateBalance atomically applies delta to wallet balance with retry logic for transient errors.
func (r *WalletRepo) UpdateBalance(ctx context.Context, id uuid.UUID, delta decimal.Decimal, opType domain.OperationType) (*domain.Wallet, error) {
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		wallet, err := r.executeUpdate(ctx, id, delta, opType)
		if err == nil {
			return wallet, nil
		}

		// Non-retryable domain errors
		if errors.Is(err, domain.ErrWalletNotFound) || errors.Is(err, domain.ErrInsufficientFunds) {
			return nil, err
		}

		// Check if error is transient and retryable
		if isRetryable(err) {
			lastErr = err
			// Exponential backoff with jitter
			sleepTime := time.Duration(1<<attempt)*time.Millisecond + time.Duration(rand.Intn(3))*time.Millisecond
			select {
			case <-time.After(sleepTime):
				continue
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}

		return nil, err
	}

	return nil, fmt.Errorf("update balance failed after %d retries: %w", maxRetries, lastErr)
}

func (r *WalletRepo) executeUpdate(ctx context.Context, id uuid.UUID, delta decimal.Decimal, opType domain.OperationType) (*domain.Wallet, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Atomic update: only updates if (balance + delta >= 0)
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
			// Determine if wallet is missing or balance was insufficient
			var existingID uuid.UUID
			checkErr := tx.QueryRow(ctx, `SELECT id FROM wallets WHERE id = $1`, id).Scan(&existingID)
			if errors.Is(checkErr, pgx.ErrNoRows) {
				return nil, domain.ErrWalletNotFound
			}
			return nil, domain.ErrInsufficientFunds
		}
		return nil, err
	}

	// Insert audit record
	insertAuditQuery := `
		INSERT INTO wallet_transactions (wallet_id, operation_type, amount, balance_after, created_at)
		VALUES ($1, $2, $3, $4, NOW())
	`
	_, _ = tx.Exec(ctx, insertAuditQuery, id, string(opType), delta.Abs(), w.Balance)

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &w, nil
}

func isRetryable(err error) bool {
	if err == nil {
		return false
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		// 40001: serialization_failure, 40P01: deadlock_detected, 55P03: lock_not_available
		return pgErr.Code == "40001" || pgErr.Code == "40P01" || pgErr.Code == "55P03"
	}
	return false
}
