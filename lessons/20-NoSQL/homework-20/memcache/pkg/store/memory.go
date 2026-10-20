package store

import (
	"context"
	"strings"
	"sync"
	"time"
)

type cacheItem struct {
	value     string
	expiresAt time.Time
}

// MemoryStore реализует Store в памяти с поддержкой TTL.
type MemoryStore struct {
	mu    sync.RWMutex
	items map[string]cacheItem
	ttl   time.Duration
}

// NewMemoryStore создает экземпляр MemoryStore.
func NewMemoryStore(ttl time.Duration) *MemoryStore {
	if ttl <= 0 || ttl > 24*time.Hour {
		ttl = 24 * time.Hour
	}
	return &MemoryStore{
		items: make(map[string]cacheItem),
		ttl:   ttl,
	}
}

// Set сохраняет пару в локальной хэш-таблице с временем истечения.
func (m *MemoryStore) Set(_ context.Context, shortURL, destination string) error {
	key := strings.TrimSpace(shortURL)
	if key == "" {
		return ErrEmptyKey
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.items[key] = cacheItem{
		value:     destination,
		expiresAt: time.Now().Add(m.ttl),
	}
	return nil
}

// Get возвращает значение или ErrNotFound, если запись отсутствует или истек TTL.
func (m *MemoryStore) Get(_ context.Context, shortURL string) (string, error) {
	key := strings.TrimSpace(shortURL)
	if key == "" {
		return "", ErrEmptyKey
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	item, exists := m.items[key]
	if !exists {
		return "", ErrNotFound
	}

	if time.Now().After(item.expiresAt) {
		return "", ErrNotFound
	}

	return item.value, nil
}

// Delete удаляет ключ.
func (m *MemoryStore) Delete(_ context.Context, shortURL string) error {
	key := strings.TrimSpace(shortURL)
	if key == "" {
		return ErrEmptyKey
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.items, key)
	return nil
}

// Close очищает данные.
func (m *MemoryStore) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.items = make(map[string]cacheItem)
	return nil
}
