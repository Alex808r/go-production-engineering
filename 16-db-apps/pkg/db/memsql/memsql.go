// Package memsql реализует in-memory хранилище книг для быстрого тестирования.
// Позволяет изолировать бизнес-логику и HTTP API от внешних СУБД.
package memsql

import (
	"context"
	"sync"

	"go-core-4/16-db-apps/pkg/db"
)

// DB — потокобезопасное хранилище книг в оперативной памяти.
type DB struct {
	mu    sync.RWMutex
	books []db.Book
}

// New создает новый экземпляр in-memory БД книг.
func New() *DB {
	return &DB{
		books: make([]db.Book, 0),
	}
}

// Books возвращает копию списка книг из памяти.
func (d *DB) Books(_ context.Context) ([]db.Book, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	res := make([]db.Book, len(d.books))
	copy(res, d.books)
	return res, nil
}

// AddBooks добавляет книги в хранилище в памяти с автогенерацией ID при необходимости.
func (d *DB) AddBooks(_ context.Context, books []db.Book) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	for _, b := range books {
		if b.ID == 0 {
			b.ID = len(d.books) + 1
		}
		d.books = append(d.books, b)
	}
	return nil
}
