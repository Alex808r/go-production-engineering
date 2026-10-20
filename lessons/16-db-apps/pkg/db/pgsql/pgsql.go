// Package pgsql реализует репозиторий книг для PostgreSQL на базе pgxpool.
package pgsql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"go-core-4/lessons/16-db-apps/pkg/db"
)

// DB — реализация интерфейса db.Interface для PostgreSQL.
type DB struct {
	pool *pgxpool.Pool
}

// New создает экземпляр DB с переданным пулом соединений.
func New(pool *pgxpool.Pool) *DB {
	return &DB{pool: pool}
}

// Connect создает и проверяет новый пул соединений pgxpool.Pool по строке подключения.
func Connect(ctx context.Context, connString string) (*DB, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.New: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping failed: %w", err)
	}
	return &DB{pool: pool}, nil
}

// Close закрывает пул соединений.
func (d *DB) Close() {
	if d.pool != nil {
		d.pool.Close()
	}
}

// Books возвращает список книг из БД PostgreSQL.
func (d *DB) Books(ctx context.Context) ([]db.Book, error) {
	if d.pool == nil {
		return nil, fmt.Errorf("connection pool is not initialized")
	}

	rows, err := d.pool.Query(ctx, `
		SELECT b.id, b.title, b.year, b.public_domain, b.publisher_id, COALESCE(p.name, '')
		FROM books b
		LEFT JOIN publishers p ON b.publisher_id = p.id
		ORDER BY b.id;
	`)
	if err != nil {
		return nil, fmt.Errorf("query books: %w", err)
	}
	defer rows.Close()

	var books []db.Book
	for rows.Next() {
		var b db.Book
		err := rows.Scan(
			&b.ID,
			&b.Title,
			&b.Year,
			&b.Public,
			&b.PublisherID,
			&b.Publisher,
		)
		if err != nil {
			return nil, fmt.Errorf("scan book: %w", err)
		}
		books = append(books, b)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return books, nil
}

// AddBooks пакетно добавляет книги в PostgreSQL одной транзакцией.
func (d *DB) AddBooks(ctx context.Context, books []db.Book) error {
	if len(books) == 0 {
		return nil
	}
	if d.pool == nil {
		return fmt.Errorf("connection pool is not initialized")
	}

	tx, err := d.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	batch := new(pgx.Batch)
	for _, b := range books {
		batch.Queue(`
			INSERT INTO books (title, year, public_domain, publisher_id)
			VALUES ($1, $2, $3, $4);
		`, b.Title, b.Year, b.Public, b.PublisherID)
	}

	br := tx.SendBatch(ctx, batch)
	for range books {
		if _, err := br.Exec(); err != nil {
			_ = br.Close()
			return fmt.Errorf("batch exec: %w", err)
		}
	}

	if err := br.Close(); err != nil {
		return fmt.Errorf("batch close: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
