package storage

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// PostgresStorage реализует Storage для СУБД PostgreSQL.
type PostgresStorage struct {
	pool *pgxpool.Pool
}

// NewPostgresStorage подключается к пулу соединений PostgreSQL и накатывает схему с индексами.
func NewPostgresStorage(dsn string) (*PostgresStorage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга DSN: %w", err)
	}

	config.MaxConns = 25
	config.MinConns = 2

	pool, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к пулу PostgreSQL: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ошибка ping PostgreSQL: %w", err)
	}

	s := &PostgresStorage{pool: pool}
	if err := s.initSchema(ctx); err != nil {
		return nil, fmt.Errorf("ошибка инициализации схемы PostgreSQL: %w", err)
	}

	return s, nil
}

// initSchema гарантирует создание таблицы links и индекса по short_url согласно ТЗ.
func (s *PostgresStorage) initSchema(ctx context.Context) error {
	schema := `
	CREATE TABLE IF NOT EXISTS links (
		id SERIAL PRIMARY KEY,
		short_url VARCHAR(64) UNIQUE NOT NULL,
		destination TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_links_short_url ON links (short_url);
	CREATE INDEX IF NOT EXISTS idx_links_destination ON links (destination);
	`
	_, err := s.pool.Exec(ctx, schema)
	return err
}

// Save сохраняет пару short_url и destination в PostgreSQL.
func (s *PostgresStorage) Save(ctx context.Context, shortURL, destination string) error {
	shortURL = strings.TrimSpace(shortURL)
	destination = strings.TrimSpace(destination)
	if shortURL == "" || destination == "" {
		return ErrEmptyField
	}

	query := `INSERT INTO links (short_url, destination) VALUES ($1, $2) ON CONFLICT (short_url) DO NOTHING`
	_, err := s.pool.Exec(ctx, query, shortURL, destination)
	if err != nil {
		return fmt.Errorf("ошибка вставки в PostgreSQL: %w", err)
	}
	return nil
}

// FindByShort находит длинную ссылку по короткой с использованием индекса idx_links_short_url.
func (s *PostgresStorage) FindByShort(ctx context.Context, shortURL string) (string, error) {
	shortURL = strings.TrimSpace(shortURL)
	if shortURL == "" {
		return "", ErrEmptyField
	}

	query := `SELECT destination FROM links WHERE short_url = $1 LIMIT 1`
	var destination string
	err := s.pool.QueryRow(ctx, query, shortURL).Scan(&destination)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrNotFound
		}
		return "", fmt.Errorf("ошибка поиска ссылки в PostgreSQL: %w", err)
	}

	return destination, nil
}

// FindByDestination возвращает существующий short_url для указанного destination.
func (s *PostgresStorage) FindByDestination(ctx context.Context, destination string) (string, error) {
	destination = strings.TrimSpace(destination)
	if destination == "" {
		return "", ErrEmptyField
	}

	query := `SELECT short_url FROM links WHERE destination = $1 LIMIT 1`
	var shortURL string
	err := s.pool.QueryRow(ctx, query, destination).Scan(&shortURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrNotFound
		}
		return "", fmt.Errorf("ошибка поиска дубликата в PostgreSQL: %w", err)
	}

	return shortURL, nil
}

// Close закрывает пул соединений с базой данных.
func (s *PostgresStorage) Close() error {
	if s.pool != nil {
		s.pool.Close()
	}
	return nil
}
