package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"

	"go-core-4/20-NoSQL/homework-20/memcache/pkg/store"
)

func TestMemcacheAPI_SetAndGet(t *testing.T) {
	memStore := store.NewMemoryStore(time.Hour)
	logger := zerolog.New(os.Stdout)
	appAPI := New(memStore, logger)

	// 1. POST /api/v1/cache
	reqBody := `{"shortUrl":"qz6d7","destination":"https://github.com/thinknetica/go_course_4"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/cache", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	appAPI.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d: %s", rec.Code, rec.Body.String())
	}

	var pair CachePair
	if err := json.Unmarshal(rec.Body.Bytes(), &pair); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if pair.ShortURL != "qz6d7" || pair.Destination != "https://github.com/thinknetica/go_course_4" {
		t.Errorf("unexpected saved pair: %+v", pair)
	}

	// 2. GET /api/v1/cache/qz6d7
	req2 := httptest.NewRequest(http.MethodGet, "/api/v1/cache/qz6d7", nil)
	rec2 := httptest.NewRecorder()
	appAPI.Router().ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", rec2.Code, rec2.Body.String())
	}

	var fetched CachePair
	if err := json.Unmarshal(rec2.Body.Bytes(), &fetched); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if fetched.Destination != "https://github.com/thinknetica/go_course_4" {
		t.Errorf("expected destination match, got %s", fetched.Destination)
	}

	// 3. GET несуществующий ключ -> 404
	req3 := httptest.NewRequest(http.MethodGet, "/api/v1/cache/unknown", nil)
	rec3 := httptest.NewRecorder()
	appAPI.Router().ServeHTTP(rec3, req3)

	if rec3.Code != http.StatusNotFound {
		t.Errorf("expected 404 Not Found, got %d", rec3.Code)
	}

	// 4. Metrics & Healthz
	req4 := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec4 := httptest.NewRecorder()
	appAPI.Router().ServeHTTP(rec4, req4)
	if rec4.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec4.Code)
	}

	req5 := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec5 := httptest.NewRecorder()
	appAPI.Router().ServeHTTP(rec5, req5)
	if rec5.Code != http.StatusOK || !strings.Contains(rec5.Body.String(), "http_requests_total") {
		t.Errorf("expected prometheus metrics in response")
	}
}
