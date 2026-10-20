package enhancedrouting

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEnhancedRouting_PathValue(t *testing.T) {
	svc := NewService()
	router := svc.NewRouter()

	// 1. GET /items/1 -> 200 OK
	req := httptest.NewRequest(http.MethodGet, "/items/1", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}

	var it Item
	_ = json.NewDecoder(rec.Body).Decode(&it)
	if it.ID != "1" || it.Name != "Laptop" {
		t.Errorf("item mismatch: %+v", it)
	}

	// 2. GET /items/999 -> 404 Not Found
	req = httptest.NewRequest(http.MethodGet, "/items/999", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestEnhancedRouting_MethodMatching(t *testing.T) {
	svc := NewService()
	router := svc.NewRouter()

	// 1. POST /items -> 201 Created
	body := `{"id": "3", "name": "Tablet"}`
	req := httptest.NewRequest(http.MethodPost, "/items", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d", rec.Code)
	}

	// 2. В Go 1.22+ отправка незарегистрированного метода автоматически возвращает 405 Method Not Allowed!
	req = httptest.NewRequest(http.MethodPut, "/items", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405 Method Not Allowed, got %d", rec.Code)
	}
}

func TestEnhancedRouting_Wildcard(t *testing.T) {
	svc := NewService()
	router := svc.NewRouter()

	req := httptest.NewRequest(http.MethodGet, "/files/images/2026/diagram.png", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}

	if !strings.Contains(rec.Body.String(), "images/2026/diagram.png") {
		t.Errorf("wildcard path mismatch: %s", rec.Body.String())
	}
}

func BenchmarkEnhancedRouting_ServeHTTP(b *testing.B) {
	svc := NewService()
	router := svc.NewRouter()
	req := httptest.NewRequest(http.MethodGet, "/items/1", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
	}
}
