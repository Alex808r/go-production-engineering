// Package repository содержит вторичные адаптеры (Secondary Adapters / Gateways)
// для взаимодействия со слоем хранения данных.
package repository

import (
	"context"
	"sort"
	"sync"

	"go-core-4/lessons/17-system-design/homework-17/cleanarch/entity"
	"go-core-4/lessons/17-system-design/homework-17/cleanarch/usecase"
)

// MemoryMovieRepository — потокобезопасная реализация usecase.MovieRepository в памяти.
// Позволяет запускать тесты без поднятия внешних СУБД.
type MemoryMovieRepository struct {
	mu     sync.RWMutex
	movies map[int64]entity.Movie
	seq    int64
}

func NewMemoryMovieRepository() *MemoryMovieRepository {
	return &MemoryMovieRepository{
		movies: make(map[int64]entity.Movie),
	}
}

func (r *MemoryMovieRepository) Save(_ context.Context, m entity.Movie) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.seq++
	m.ID = r.seq
	r.movies[m.ID] = m
	return m.ID, nil
}

func (r *MemoryMovieRepository) SaveBatch(_ context.Context, movies []entity.Movie) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, m := range movies {
		r.seq++
		m.ID = r.seq
		r.movies[m.ID] = m
	}
	return nil
}

func (r *MemoryMovieRepository) FindByID(_ context.Context, id int64) (entity.Movie, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	m, ok := r.movies[id]
	if !ok {
		return entity.Movie{}, usecase.ErrNotFound
	}
	return m, nil
}

func (r *MemoryMovieRepository) FindAll(_ context.Context, studioID int) ([]entity.Movie, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []entity.Movie
	for _, m := range r.movies {
		if studioID == 0 || m.StudioID == studioID {
			result = append(result, m)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})

	return result, nil
}

func (r *MemoryMovieRepository) Update(_ context.Context, m entity.Movie) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.movies[m.ID]; !ok {
		return usecase.ErrNotFound
	}
	r.movies[m.ID] = m
	return nil
}

func (r *MemoryMovieRepository) Delete(_ context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.movies[id]; !ok {
		return usecase.ErrNotFound
	}
	delete(r.movies, id)
	return nil
}
