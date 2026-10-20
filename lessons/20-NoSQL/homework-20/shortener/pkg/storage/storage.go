// Package storage определяет интерфейс постоянного хранения ссылок в PostgreSQL.
package storage

import (
	"context"
	"errors"
)

var (
	// ErrNotFound возвращается, если запись не найдена в БД.
	ErrNotFound = errors.New("ссылка не найдена в базе данных")
	// ErrEmptyField возвращается при пустых значениях аргументов.
	ErrEmptyField = errors.New("ссылка или идентификатор не могут быть пустыми")
)

// Storage описывает контракт взаимодействия с СУБД PostgreSQL.
type Storage interface {
	// Save сохраняет пару short_url и destination в таблицу links.
	Save(ctx context.Context, shortURL, destination string) error
	// FindByShort находит длинную ссылку по короткой с использованием индекса.
	FindByShort(ctx context.Context, shortURL string) (string, error)
	// FindByDestination возвращает существующий short_url для указанного destination.
	FindByDestination(ctx context.Context, destination string) (string, error)
	// Close освобождает соединения с базой данных.
	Close() error
}
