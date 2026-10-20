package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/rs/zerolog"

	"go-core-4/lessons/20-NoSQL/homework-20/shortener/pkg/cacheclient"
	"go-core-4/lessons/20-NoSQL/homework-20/shortener/pkg/storage"
)

// mockCacheClient имитирует сервис memcache для тестов.
type mockCacheClient struct {
	mu    sync.RWMutex
	items map[string]string
}

func newMockCacheClient() *mockCacheClient {
	return &mockCacheClient{items: make(map[string]string)}
}

func (m *mockCacheClient) Get(_ context.Context, shortURL string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	val, ok := m.items[shortURL]
	if !ok {
		return "", cacheclient.ErrCacheMiss
	}
	return val, nil
}

func (m *mockCacheClient) Set(_ context.Context, shortURL, destination string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.items[shortURL] = destination
	return nil
}

func TestShortenerAPI_CreateAndRedirect_CacheAside(t *testing.T) {
	store := storage.NewMemoryStorage()
	cache := newMockCacheClient()
	logger := zerolog.New(os.Stdout)
	appAPI := New(store, cache, "https://lynks.org", logger)

	// 1. POST / (Создание короткой ссылки)
	reqBody := `{"destination":"https://github.com/thinknetica/go_course_4"}`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	appAPI.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp ShortenResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !strings.HasPrefix(resp.ShortURL, "https://lynks.org/") {
		t.Errorf("unexpected short URL prefix: %s", resp.ShortURL)
	}

	shortID := strings.TrimPrefix(resp.ShortURL, "https://lynks.org/")

	// 2. GET /{short} (Редирект: Cache Miss -> Postgres -> Cache Warming -> Redirect)
	req2 := httptest.NewRequest(http.MethodGet, "/"+shortID, nil)
	rec2 := httptest.NewRecorder()
	appAPI.Router().ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusFound {
		t.Fatalf("expected 302 Found, got %d", rec2.Code)
	}
	if rec2.Header().Get("Location") != "https://github.com/thinknetica/go_course_4" {
		t.Errorf("unexpected redirect location: %s", rec2.Header().Get("Location"))
	}

	// 3. GET /{short} повторно (Cache Hit)
	// Добавим напрямую в mock кэш
	_ = cache.Set(context.Background(), shortID, "https://github.com/thinknetica/go_course_4")

	req3 := httptest.NewRequest(http.MethodGet, "/"+shortID, nil)
	rec3 := httptest.NewRecorder()
	appAPI.Router().ServeHTTP(rec3, req3)

	if rec3.Code != http.StatusFound {
		t.Fatalf("expected 302 Found on Cache Hit, got %d", rec3.Code)
	}

	// 4. GET несуществующий ключ -> 404
	req4 := httptest.NewRequest(http.MethodGet, "/unknown99", nil)
	rec4 := httptest.NewRecorder()
	appAPI.Router().ServeHTTP(rec4, req4)
	if rec4.Code != http.StatusNotFound {
		t.Errorf("expected 404 Not Found, got %d", rec4.Code)
	}

	// 5. Metrics
	req5 := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec5 := httptest.NewRecorder()
	appAPI.Router().ServeHTTP(rec5, req5)
	if rec5.Code != http.StatusOK || !strings.Contains(rec5.Body.String(), "http_requests_total") {
		t.Errorf("expected prometheus metrics")
	}
}
