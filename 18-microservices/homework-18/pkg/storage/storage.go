// Package storage определяет контракт и типы хранилища коротких ссылок.
package storage

import (
	"context"
	"errors"
	"time"
)

var (
	// ErrNotFound возвращается, если ссылка не найдена в хранилище.
	ErrNotFound = errors.New("ссылка не найдена")
	// ErrEmptyURL возвращается при попытке сохранить пустой URL.
	ErrEmptyURL = errors.New("URL не может быть пустым")
	// ErrInvalidURL возвращается при некорректном формате URL.
	ErrInvalidURL = errors.New("некорректный формат URL")
)

// Link представляет сущность короткой гиперссылки.
type Link struct {
	ID        string    `json:"id"`
	Original  string    `json:"original"`
	ShortURL  string    `json:"short_url"`
	CreatedAt time.Time `json:"created_at"`
}

// Interface определяет методы работы с хранилищем ссылок микросервиса.
type Interface interface {
	// Save сохраняет оригинальный URL и возвращает сформированную ссылку.
	Save(ctx context.Context, originalURL string) (Link, error)
	// Get возвращает оригинальную ссылку по её идентификатору или полной короткой ссылке.
	Get(ctx context.Context, idOrShort string) (Link, error)
	// All возвращает срез всех сохранённых ссылок.
	All(ctx context.Context) ([]Link, error)
}
