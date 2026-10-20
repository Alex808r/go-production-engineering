// Package memsql реализует in-memory хранилище фильмов для тестирования.
// Не зависит от внешних СУБД и позволяет запускать тесты мгновенно.
package memsql

import (
	"context"
	"sort"
	"sync"

	"go-core-4/lessons/16-db-apps/homework-16/pkg/db"
)

// DB — потокобезопасное in-memory хранилище фильмов.
type DB struct {
	mu     sync.RWMutex
	movies map[int64]db.Movie
	seq    int64 // автоинкрементный счетчик ID
}

// New создает новое пустое in-memory хранилище.
func New() *DB {
	return &DB{
		movies: make(map[int64]db.Movie),
	}
}

// Movies возвращает массив фильмов.
// Если studioID == 0, возвращаются все фильмы.
// Если studioID > 0, выполняется фильтрация по ID студии.
func (d *DB) Movies(_ context.Context, studioID int) ([]db.Movie, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	result := make([]db.Movie, 0, len(d.movies))
	for _, m := range d.movies {
		if studioID == 0 || m.StudioID == studioID {
			result = append(result, m)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})

	return result, nil
}

// MovieByID возвращает фильм по ID. Возвращает db.ErrNotFound, если фильм не найден.
func (d *DB) MovieByID(_ context.Context, id int64) (db.Movie, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	m, ok := d.movies[id]
	if !ok {
		return db.Movie{}, db.ErrNotFound
	}
	return m, nil
}

// AddMovies добавляет массив фильмов в хранилище.
// Каждому фильму при необходимости генерируется уникальный ID.
func (d *DB) AddMovies(_ context.Context, movies []db.Movie) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	for _, m := range movies {
		d.seq++
		if m.ID == 0 {
			m.ID = d.seq
		} else if m.ID > d.seq {
			d.seq = m.ID
		}
		d.movies[m.ID] = m
	}
	return nil
}

// UpdateMovie обновляет существующий фильм. Возвращает db.ErrNotFound, если фильм не найден.
func (d *DB) UpdateMovie(_ context.Context, m db.Movie) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, ok := d.movies[m.ID]; !ok {
		return db.ErrNotFound
	}
	d.movies[m.ID] = m
	return nil
}

// DeleteMovie удаляет фильм по ID. Возвращает db.ErrNotFound, если фильм не найден.
func (d *DB) DeleteMovie(_ context.Context, id int64) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, ok := d.movies[id]; !ok {
		return db.ErrNotFound
	}
	delete(d.movies, id)
	return nil
}
