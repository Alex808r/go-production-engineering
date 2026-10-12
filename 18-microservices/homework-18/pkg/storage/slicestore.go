// Package storage реализует хранилище ссылок на базе слайса в памяти.
// Требование ТЗ Урока №18: "Хранилище ссылок должно быть реализовано в виде слайса в памяти."
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

// SliceStore — потокобезопасная реализация хранилища на базе среза в памяти.
type SliceStore struct {
	mu      sync.RWMutex
	links   []Link // Срез в памяти
	baseURL string
}

// NewSliceStore создает новое хранилище ссылок.
// baseURL используется для формирования абсолютной короткой ссылки (например, "http://localhost:8080").
func NewSliceStore(baseURL string) *SliceStore {
	baseURL = strings.TrimRight(baseURL, "/")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	return &SliceStore{
		links:   make([]Link, 0),
		baseURL: baseURL,
	}
}

// Save сохраняет оригинальную ссылку в слайс и возвращает структуру Link.
func (s *SliceStore) Save(_ context.Context, originalURL string) (Link, error) {
	originalURL = strings.TrimSpace(originalURL)
	if originalURL == "" {
		return Link{}, ErrEmptyURL
	}

	// Базовая валидация URL
	parsed, err := url.ParseRequestURI(originalURL)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return Link{}, fmt.Errorf("%w: требуется схема http или https", ErrInvalidURL)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Проверяем, может ссылка уже сохранена в срезе (идемпотентность)
	for _, l := range s.links {
		if l.Original == originalURL {
			return l, nil
		}
	}

	// Генерируем уникальный короткий идентификатор (6 символов hex)
	id := s.generateID()
	shortURL := fmt.Sprintf("%s/r/%s", s.baseURL, id)

	link := Link{
		ID:        id,
		Original:  originalURL,
		ShortURL:  shortURL,
		CreatedAt: time.Now(),
	}

	// Добавление в слайс в памяти
	s.links = append(s.links, link)

	return link, nil
}

// Get ищет ссылку по ID или по полному короткому URL.
func (s *SliceStore) Get(_ context.Context, idOrShort string) (Link, error) {
	idOrShort = strings.TrimSpace(idOrShort)
	if idOrShort == "" {
		return Link{}, ErrNotFound
	}

	// Извлекаем чистый ID, если передали полный URL (например, "http://localhost:8080/r/abcd12")
	id := idOrShort
	if idx := strings.LastIndex(idOrShort, "/"); idx != -1 {
		id = idOrShort[idx+1:]
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	// Линейный поиск по слайсу в памяти
	for _, l := range s.links {
		if l.ID == id || l.ShortURL == idOrShort {
			return l, nil
		}
	}

	return Link{}, ErrNotFound
}

// All возвращает копию всех ссылок из слайса.
func (s *SliceStore) All(_ context.Context) ([]Link, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make([]Link, len(s.links))
	copy(res, s.links)
	return res, nil
}

// generateID генерирует уникальный случайный идентификатор для короткой ссылки.
func (s *SliceStore) generateID() string {
	for {
		bytes := make([]byte, 4)
		_, _ = rand.Read(bytes)
		candidate := hex.EncodeToString(bytes)

		// Проверяем уникальность среди существующих в срезе
		exists := false
		for _, l := range s.links {
			if l.ID == candidate {
				exists = true
				break
			}
		}
		if !exists {
			return candidate
		}
	}
}
