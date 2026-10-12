package storage

import (
	"context"
	"strings"
	"sync"
)

// MemoryStorage реализует Storage в оперативной памяти с потокобезопасностью.
type MemoryStorage struct {
	mu            sync.RWMutex
	shortToDest   map[string]string
	destToShort   map[string]string
}

// NewMemoryStorage создает экземпляр MemoryStorage.
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		shortToDest: make(map[string]string),
		destToShort: make(map[string]string),
	}
}

// Save сохраняет соответствие короткой ссылки и назначения.
func (m *MemoryStorage) Save(_ context.Context, shortURL, destination string) error {
	shortURL = strings.TrimSpace(shortURL)
	destination = strings.TrimSpace(destination)
	if shortURL == "" || destination == "" {
		return ErrEmptyField
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.shortToDest[shortURL] = destination
	m.destToShort[destination] = shortURL
	return nil
}

// FindByShort находит destination по shortURL.
func (m *MemoryStorage) FindByShort(_ context.Context, shortURL string) (string, error) {
	shortURL = strings.TrimSpace(shortURL)
	if shortURL == "" {
		return "", ErrEmptyField
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	dest, ok := m.shortToDest[shortURL]
	if !ok {
		return "", ErrNotFound
	}
	return dest, nil
}

// FindByDestination находит shortURL по destination.
func (m *MemoryStorage) FindByDestination(_ context.Context, destination string) (string, error) {
	destination = strings.TrimSpace(destination)
	if destination == "" {
		return "", ErrEmptyField
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	short, ok := m.destToShort[destination]
	if !ok {
		return "", ErrNotFound
	}
	return short, nil
}

// Close очищает мапы.
func (m *MemoryStorage) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.shortToDest = make(map[string]string)
	m.destToShort = make(map[string]string)
	return nil
}
