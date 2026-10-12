package storage

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"
)

// SliceStore реализует Interface с хранением ссылок в срезе в оперативной памяти.
type SliceStore struct {
	mu      sync.RWMutex
	items   []Link
	baseURL string
}

// NewSliceStore создает новое хранилище на базе среза.
func NewSliceStore(baseURL string) *SliceStore {
	baseURL = strings.TrimRight(baseURL, "/")
	return &SliceStore{
		items:   make([]Link, 0),
		baseURL: baseURL,
	}
}

// Save валидирует URL, проверяет отсутствие дубликатов и сохраняет ссылку в срезе.
func (s *SliceStore) Save(ctx context.Context, originalURL string) (Link, error) {
	cleanURL := strings.TrimSpace(originalURL)
	if cleanURL == "" {
		return Link{}, ErrEmptyURL
	}

	parsed, err := url.ParseRequestURI(cleanURL)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		return Link{}, fmt.Errorf("%w: %s (поддерживаются только схемы http/https)", ErrInvalidURL, originalURL)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Проверяем, не сохраняли ли уже такую ссылку
	for _, l := range s.items {
		if l.Original == cleanURL {
			return l, nil
		}
	}

	id := s.generateID()
	shortURL := fmt.Sprintf("%s/r/%s", s.baseURL, id)

	link := Link{
		ID:        id,
		Original:  cleanURL,
		ShortURL:  shortURL,
		CreatedAt: time.Now().UTC(),
	}

	s.items = append(s.items, link)
	return link, nil
}

// Get возвращает ссылку по её ID или по полному short_url.
func (s *SliceStore) Get(ctx context.Context, idOrShort string) (Link, error) {
	target := strings.TrimSpace(idOrShort)
	if target == "" {
		return Link{}, ErrNotFound
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, l := range s.items {
		if l.ID == target || l.ShortURL == target {
			return l, nil
		}
	}

	return Link{}, ErrNotFound
}

// All возвращает копию среза всех ссылок.
func (s *SliceStore) All(ctx context.Context) ([]Link, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make([]Link, len(s.items))
	copy(res, s.items)
	return res, nil
}

// generateID генерирует уникальный случайный 8-значный hex-идентификатор.
func (s *SliceStore) generateID() string {
	b := make([]byte, 4)
	for {
		if _, err := rand.Read(b); err != nil {
			return fmt.Sprintf("%x", time.Now().UnixNano())[:8]
		}
		candidate := hex.EncodeToString(b)
		collision := false
		for _, item := range s.items {
			if item.ID == candidate {
				collision = true
				break
			}
		}
		if !collision {
			return candidate
		}
	}
}
