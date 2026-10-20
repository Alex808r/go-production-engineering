package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"go-core-4/lessons/19-queue/homework-19/pkg/event"
	"go-core-4/lessons/19-queue/homework-19/pkg/shortener/storage"
)

// mockProducer перехватывает опубликованные события для проверки в тестах.
type mockProducer struct {
	mu     sync.Mutex
	events []event.LinkCreatedEvent
	notify chan struct{}
}

func newMockProducer() *mockProducer {
	return &mockProducer{
		notify: make(chan struct{}, 10),
	}
}

func (m *mockProducer) Publish(_ context.Context, _ string, ev event.LinkCreatedEvent) error {
	m.mu.Lock()
	m.events = append(m.events, ev)
	m.mu.Unlock()
	m.notify <- struct{}{}
	return nil
}

func (m *mockProducer) Close() error {
	return nil
}

func TestAPI_CreateLink_PublishesEvent(t *testing.T) {
	store := storage.NewSliceStore("http://localhost:8080")
	producer := newMockProducer()
	appAPI := New(store, producer)

	// 1. Создание ссылки через POST /api/v1/links
	reqBody := `{"url":"https://go.dev/play/"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/links", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	appAPI.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected HTTP 201 Created, got %d: %s", rec.Code, rec.Body.String())
	}

	var created storage.Link
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if created.ID == "" || created.Original != "https://go.dev/play/" {
		t.Errorf("unexpected created link: %+v", created)
	}

	// Проверяем, что событие прилетело в producer
	select {
	case <-producer.notify:
	case <-time.After(1 * time.Second):
		t.Fatal("timed out waiting for event publication")
	}

	producer.mu.Lock()
	defer producer.mu.Unlock()

	if len(producer.events) != 1 {
		t.Fatalf("expected 1 published event, got %d", len(producer.events))
	}

	ev := producer.events[0]
	if ev.LinkID != created.ID || ev.Original != "https://go.dev/play/" || ev.ShortURL != created.ShortURL {
		t.Errorf("published event mismatch: %+v", ev)
	}
}

func TestAPI_GetAndResolve(t *testing.T) {
	store := storage.NewSliceStore("http://localhost:8080")
	producer := newMockProducer()
	appAPI := New(store, producer)

	ctx := context.Background()
	link, err := store.Save(ctx, "https://kubernetes.io")
	if err != nil {
		t.Fatalf("unexpected save error: %v", err)
	}

	// GET /api/v1/links/{id}
	req := httptest.NewRequest(http.MethodGet, "/api/v1/links/"+link.ID, nil)
	rec := httptest.NewRecorder()
	appAPI.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	// POST /api/v1/links/resolve
	resolveBody := `{"short_url":"` + link.ShortURL + `"}`
	req2 := httptest.NewRequest(http.MethodPost, "/api/v1/links/resolve", bytes.NewBufferString(resolveBody))
	rec2 := httptest.NewRecorder()
	appAPI.Router().ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec2.Code)
	}

	// GET /r/{id} redirect
	req3 := httptest.NewRequest(http.MethodGet, "/r/"+link.ID, nil)
	rec3 := httptest.NewRecorder()
	appAPI.Router().ServeHTTP(rec3, req3)

	if rec3.Code != http.StatusFound {
		t.Fatalf("expected 302, got %d", rec3.Code)
	}
	if rec3.Header().Get("Location") != "https://kubernetes.io" {
		t.Errorf("unexpected redirect location: %s", rec3.Header().Get("Location"))
	}
}

func TestAPI_ErrorsAndHealth(t *testing.T) {
	store := storage.NewSliceStore("http://localhost:8080")
	appAPI := New(store, nil)

	// Невалидный JSON
	req := httptest.NewRequest(http.MethodPost, "/api/v1/links", bytes.NewBufferString("invalid json"))
	rec := httptest.NewRecorder()
	appAPI.Router().ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}

	// Несуществующий ID
	req2 := httptest.NewRequest(http.MethodGet, "/api/v1/links/unknown", nil)
	rec2 := httptest.NewRecorder()
	appAPI.Router().ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec2.Code)
	}

	// Healthz
	req3 := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec3 := httptest.NewRecorder()
	appAPI.Router().ServeHTTP(rec3, req3)
	if rec3.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec3.Code)
	}

	// Metrics
	req4 := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec4 := httptest.NewRecorder()
	appAPI.Router().ServeHTTP(rec4, req4)
	if rec4.Code != http.StatusOK || !strings.Contains(rec4.Body.String(), "http_requests_total") {
		t.Errorf("expected valid metrics, got %s", rec4.Body.String())
	}
}
