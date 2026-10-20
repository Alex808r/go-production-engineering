// Package pgsql реализует хранилище фильмов для PostgreSQL на базе pgxpool.
// Соответствует Задаче №2 домашнего задания №16.
package pgsql

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"go-core-4/lessons/16-db-apps/homework-16/pkg/db"
)

// DB — реализация хранилища фильмов на PostgreSQL через pgxpool.Pool.
type DB struct {
	pool *pgxpool.Pool
}

// New создает новое подключение к PostgreSQL через пул соединений.
// connString — строка подключения в формате: "postgres://user:pass@host/dbname".
func New(ctx context.Context, connString string) (*DB, error) {
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

// NewWithPool создает экземпляр DB с переданным пулом соединений (удобно для тестов).
func NewWithPool(pool *pgxpool.Pool) *DB {
	return &DB{pool: pool}
}

// Close закрывает пул соединений с БД.
func (d *DB) Close() {
	if d.pool != nil {
		d.pool.Close()
	}
}

// Movies возвращает массив фильмов.
// Если studioID == 0, возвращаются все фильмы.
// Если studioID > 0, выполняется фильтрация по ID студии.
func (d *DB) Movies(ctx context.Context, studioID int) ([]db.Movie, error) {
	query := `
		SELECT 
			m.id, 
			m.title, 
			m.release_year, 
			m.box_office, 
			m.rating, 
			m.studio_id,
			COALESCE(s.name, '') AS studio_name
		FROM movies m
		LEFT JOIN studios s ON m.studio_id = s.id
		WHERE ($1 = 0 OR m.studio_id = $1)
		ORDER BY m.id;
	`

	rows, err := d.pool.Query(ctx, query, studioID)
	if err != nil {
		return nil, fmt.Errorf("query movies: %w", err)
	}
	defer rows.Close()

	var movies []db.Movie
	for rows.Next() {
		var m db.Movie
		err := rows.Scan(
			&m.ID,
			&m.Title,
			&m.ReleaseYear,
			&m.BoxOffice,
			&m.Rating,
			&m.StudioID,
			&m.StudioName,
		)
		if err != nil {
			return nil, fmt.Errorf("scan movie: %w", err)
		}
		movies = append(movies, m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return movies, nil
}

// MovieByID возвращает фильм по его ID. Возвращает db.ErrNotFound, если запись отсутствует.
func (d *DB) MovieByID(ctx context.Context, id int64) (db.Movie, error) {
	query := `
		SELECT 
			m.id, 
			m.title, 
			m.release_year, 
			m.box_office, 
			m.rating, 
			m.studio_id,
			COALESCE(s.name, '') AS studio_name
		FROM movies m
		LEFT JOIN studios s ON m.studio_id = s.id
		WHERE m.id = $1;
	`

	var m db.Movie
	err := d.pool.QueryRow(ctx, query, id).Scan(
		&m.ID,
		&m.Title,
		&m.ReleaseYear,
		&m.BoxOffice,
		&m.Rating,
		&m.StudioID,
		&m.StudioName,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Movie{}, db.ErrNotFound
		}
		return db.Movie{}, fmt.Errorf("query movie by id: %w", err)
	}

	return m, nil
}

// AddMovies добавляет массив фильмов в рамках одной транзакции с использованием pgx.Batch
// (как показано на занятии №16 для эффективной пакетной вставки).
func (d *DB) AddMovies(ctx context.Context, movies []db.Movie) error {
	if len(movies) == 0 {
		return nil
	}

	tx, err := d.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	batch := new(pgx.Batch)
	for _, m := range movies {
		batch.Queue(`
			INSERT INTO movies (title, release_year, box_office, rating, studio_id)
			VALUES ($1, $2, $3, $4, $5);
		`, m.Title, m.ReleaseYear, m.BoxOffice, m.Rating, m.StudioID)
	}

	br := tx.SendBatch(ctx, batch)
	for range movies {
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

// UpdateMovie обновляет данные фильма. Если запись с таким ID не существует, возвращает db.ErrNotFound.
func (d *DB) UpdateMovie(ctx context.Context, m db.Movie) error {
	cmdTag, err := d.pool.Exec(ctx, `
		UPDATE movies
		SET title = $1, release_year = $2, box_office = $3, rating = $4, studio_id = $5
		WHERE id = $6;
	`, m.Title, m.ReleaseYear, m.BoxOffice, m.Rating, m.StudioID, m.ID)
	if err != nil {
		return fmt.Errorf("update movie: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return db.ErrNotFound
	}

	return nil
}

// DeleteMovie удаляет фильм по ID. Если запись с таким ID не существует, возвращает db.ErrNotFound.
func (d *DB) DeleteMovie(ctx context.Context, id int64) error {
	cmdTag, err := d.pool.Exec(ctx, `
		DELETE FROM movies
		WHERE id = $1;
	`, id)
	if err != nil {
		return fmt.Errorf("delete movie: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return db.ErrNotFound
	}

	return nil
}
