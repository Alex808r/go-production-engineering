// Package store определяет интерфейс и реализации кэша гиперссылок.
package store

import (
	"context"
	"errors"
)

var (
	// ErrNotFound возвращается, если ключ отсутствует в кэше.
	ErrNotFound = errors.New("ссылка не найдена в кэше")
	// ErrEmptyKey возвращается при попытке передать пустой ключ короткой ссылки.
	ErrEmptyKey = errors.New("ключ короткой ссылки не может быть пустым")
)

// Store определяет контракт работы с кэшем пар коротких и длинных ссылок.
type Store interface {
	// Set сохраняет пару (shortURL -> destination) с установленным TTL (не более 24 часов).
	Set(ctx context.Context, shortURL, destination string) error
	// Get возвращает оригинальную ссылку destination по ключу короткой ссылки.
	Get(ctx context.Context, shortURL string) (string, error)
	// Delete удаляет пару из кэша.
	Delete(ctx context.Context, shortURL string) error
	// Close освобождает ресурсы подключения.
	Close() error
}
