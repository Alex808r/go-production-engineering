package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-core-4/18-microservices/homework-18/pkg/api"
	"go-core-4/18-microservices/homework-18/pkg/storage"
)

func setupTestAPI() (*api.API, *storage.SliceStore) {
	store := storage.NewSliceStore("http://localhost:8080")
	srv := api.New(store)
	return srv, store
}

func TestAPI_CreateAndGetLink(t *testing.T) {
	srv, _ := setupTestAPI()

	// 1. Создание короткой ссылки
	payload := `{"url": "https://golang.org/doc/"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/links", bytes.NewBufferString(payload))
	rec := httptest.NewRecorder()
	srv.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d: %s", rec.Code, rec.Body.String())
	}

	var created storage.Link
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	if created.Original != "https://golang.org/doc/" {
		t.Errorf("expected original 'https://golang.org/doc/', got %q", created.Original)
	}
	if created.ID == "" || created.ShortURL == "" {
		t.Errorf("expected non-empty ID and ShortURL: %+v", created)
	}

	// 2. Получение оригинальной ссылки по ID
	reqGet := httptest.NewRequest(http.MethodGet, "/api/v1/links/"+created.ID, nil)
	recGet := httptest.NewRecorder()
	srv.Router().ServeHTTP(recGet, reqGet)

	if recGet.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recGet.Code)
	}

	var fetched storage.Link
	_ = json.Unmarshal(recGet.Body.Bytes(), &fetched)
	if fetched.Original != created.Original {
		t.Errorf("expected %q, got %q", created.Original, fetched.Original)
	}

	// 3. Получение оригинальной ссылки через POST /api/v1/links/resolve
	resolveBody := `{"short_url": "` + created.ShortURL + `"}`
	reqResolve := httptest.NewRequest(http.MethodPost, "/api/v1/links/resolve", bytes.NewBufferString(resolveBody))
	recResolve := httptest.NewRecorder()
	srv.Router().ServeHTTP(recResolve, reqResolve)

	if recResolve.Code != http.StatusOK {
		t.Fatalf("expected status 200 from resolve, got %d", recResolve.Code)
	}

	var resolved storage.Link
	_ = json.Unmarshal(recResolve.Body.Bytes(), &resolved)
	if resolved.Original != created.Original {
		t.Errorf("expected resolved %q, got %q", created.Original, resolved.Original)
	}

	// 4. Редирект HTTP 302 Found
	reqRedir := httptest.NewRequest(http.MethodGet, "/r/"+created.ID, nil)
	recRedir := httptest.NewRecorder()
	srv.Router().ServeHTTP(recRedir, reqRedir)

	if recRedir.Code != http.StatusFound {
		t.Fatalf("expected status 302, got %d", recRedir.Code)
	}
	if loc := recRedir.Header().Get("Location"); loc != created.Original {
		t.Errorf("expected redirect location %q, got %q", created.Original, loc)
	}
}

func TestAPI_ErrorCases(t *testing.T) {
	srv, _ := setupTestAPI()

	// Пустой URL
	reqEmpty := httptest.NewRequest(http.MethodPost, "/api/v1/links", bytes.NewBufferString(`{"url": ""}`))
	recEmpty := httptest.NewRecorder()
	srv.Router().ServeHTTP(recEmpty, reqEmpty)
	if recEmpty.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for empty url, got %d", recEmpty.Code)
	}

	// Некорректный JSON
	reqBadJSON := httptest.NewRequest(http.MethodPost, "/api/v1/links", bytes.NewBufferString(`{not-json}`))
	recBadJSON := httptest.NewRecorder()
	srv.Router().ServeHTTP(recBadJSON, reqBadJSON)
	if recBadJSON.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for bad JSON, got %d", recBadJSON.Code)
	}

	// Несуществующая ссылка
	reqNotFound := httptest.NewRequest(http.MethodGet, "/api/v1/links/unknown_id", nil)
	recNotFound := httptest.NewRecorder()
	srv.Router().ServeHTTP(recNotFound, reqNotFound)
	if recNotFound.Code != http.StatusNotFound {
		t.Errorf("expected status 404 for unknown link, got %d", recNotFound.Code)
	}
}

func TestAPI_HealthzAndMetrics(t *testing.T) {
	srv, _ := setupTestAPI()

	// Healthz
	reqH := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	recH := httptest.NewRecorder()
	srv.Router().ServeHTTP(recH, reqH)
	if recH.Code != http.StatusOK {
		t.Errorf("expected status 200 for healthz, got %d", recH.Code)
	}

	// Metrics
	reqM := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	recM := httptest.NewRecorder()
	srv.Router().ServeHTTP(recM, reqM)
	if recM.Code != http.StatusOK {
		t.Errorf("expected status 200 for metrics, got %d", recM.Code)
	}
	if !strings.Contains(recM.Body.String(), "http_requests_total") {
		t.Errorf("expected metrics to contain http_requests_total, got: %s", recM.Body.String())
	}
}

func BenchmarkAPI_ResolveLink(b *testing.B) {
	srv, st := setupTestAPI()
	link, _ := st.Save(context.Background(), "https://golang.org/doc")
	req := httptest.NewRequest(http.MethodGet, "/r/"+link.ID, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		srv.Router().ServeHTTP(rec, req)
	}
}

func BenchmarkAPI_CreateLink(b *testing.B) {
	srv, _ := setupTestAPI()
	body := `{"url":"https://example.com/benchmark"}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/links", bytes.NewBufferString(body))
		rec := httptest.NewRecorder()
		srv.Router().ServeHTTP(rec, req)
	}
}

