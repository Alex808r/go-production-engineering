package webapp

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// mockDataProvider реализует интерфейс DataProvider для изолированных модульных тестов.
type mockDataProvider struct {
	docs []Document
	idx  map[string][]int
}

func (m *mockDataProvider) Documents() []Document {
	return m.docs
}

func (m *mockDataProvider) Index() map[string][]int {
	return m.idx
}

func setupTestApp() (*Service, *mockDataProvider) {
	mock := &mockDataProvider{
		docs: []Document{
			{ID: 1, URL: "https://go.dev/doc/", Title: "Go Documentation"},
			{ID: 2, URL: "https://golang.org/pkg/", Title: "Standard Library Packages"},
		},
		idx: map[string][]int{
			"documentation": {1},
			"golang":        {1, 2},
			"packages":      {2},
		},
	}
	app := New(mock)
	return app, mock
}

// TestHandleDocs_HTML проверяет рендеринг страницы со списком документов в формате HTML.
func TestHandleDocs_HTML(t *testing.T) {
	app, _ := setupTestApp()

	req := httptest.NewRequest(http.MethodGet, "/docs", nil)
	rec := httptest.NewRecorder()

	app.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("ожидался HTTP 200, получен: %d", rec.Code)
	}

	ct := rec.Header().Get("Content-Type")
	if !strings.Contains(ct, "text/html") {
		t.Errorf("ожидался заголовок text/html, получен: %s", ct)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "Go Documentation") {
		t.Errorf("тело ответа не содержит 'Go Documentation'")
	}
	if !strings.Contains(body, "https://go.dev/doc/") {
		t.Errorf("тело ответа не содержит URL 'https://go.dev/doc/'")
	}
	if !strings.Contains(body, "Standard Library Packages") {
		t.Errorf("тело ответа не содержит 'Standard Library Packages'")
	}
}

// TestHandleDocs_JSON проверяет отдачу документов в формате JSON (через query param и Accept заголовок).
func TestHandleDocs_JSON(t *testing.T) {
	app, _ := setupTestApp()

	tests := []struct {
		name    string
		url     string
		headers map[string]string
	}{
		{
			name: "query param ?format=json",
			url:  "/docs?format=json",
		},
		{
			name:    "Accept header application/json",
			url:     "/docs",
			headers: map[string]string{"Accept": "application/json"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}
			rec := httptest.NewRecorder()

			app.Router().ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("ожидался HTTP 200, получен: %d", rec.Code)
			}

			var docs []Document
			if err := json.Unmarshal(rec.Body.Bytes(), &docs); err != nil {
				t.Fatalf("ошибка парсинга JSON: %v", err)
			}

			if len(docs) != 2 {
				t.Fatalf("ожидалось 2 документа, получено: %d", len(docs))
			}
			if docs[0].ID != 1 || docs[0].Title != "Go Documentation" {
				t.Errorf("неожиданные данные первого документа: %+v", docs[0])
			}
		})
	}
}

// TestHandleIndex_HTML проверяет рендеринг содержимого поискового индекса в формате HTML.
func TestHandleIndex_HTML(t *testing.T) {
	app, _ := setupTestApp()

	req := httptest.NewRequest(http.MethodGet, "/index", nil)
	rec := httptest.NewRecorder()

	app.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("ожидался HTTP 200, получен: %d", rec.Code)
	}

	ct := rec.Header().Get("Content-Type")
	if !strings.Contains(ct, "text/html") {
		t.Errorf("ожидался заголовок text/html, получен: %s", ct)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "golang") {
		t.Errorf("тело ответа не содержит токен 'golang'")
	}
	if !strings.Contains(body, "documentation") {
		t.Errorf("тело ответа не содержит токен 'documentation'")
	}
	if !strings.Contains(body, "packages") {
		t.Errorf("тело ответа не содержит токен 'packages'")
	}
}

// TestHandleIndex_JSON проверяет отдачу содержимого индекса в формате JSON.
func TestHandleIndex_JSON(t *testing.T) {
	app, _ := setupTestApp()

	req := httptest.NewRequest(http.MethodGet, "/index?format=json", nil)
	rec := httptest.NewRecorder()

	app.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("ожидался HTTP 200, получен: %d", rec.Code)
	}

	var idxMap map[string][]int
	if err := json.Unmarshal(rec.Body.Bytes(), &idxMap); err != nil {
		t.Fatalf("ошибка парсинга JSON: %v", err)
	}

	if len(idxMap) != 3 {
		t.Fatalf("ожидалось 3 термина в индексе, получено: %d", len(idxMap))
	}
	if len(idxMap["golang"]) != 2 {
		t.Errorf("для токена 'golang' ожидалось 2 документа, получено: %v", idxMap["golang"])
	}
}

// TestMethodNotAllowed проверяет, что неподдерживаемые HTTP-методы возвращают 405.
func TestMethodNotAllowed(t *testing.T) {
	app, _ := setupTestApp()

	endpoints := []string{"/docs", "/index"}
	for _, ep := range endpoints {
		t.Run("POST "+ep, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, ep, nil)
			rec := httptest.NewRecorder()

			app.Router().ServeHTTP(rec, req)

			if rec.Code != http.StatusMethodNotAllowed {
				t.Errorf("ожидался статус 405 MethodNotAllowed, получен: %d", rec.Code)
			}
		})
	}
}

// TestEmptyData проверяет корректную обработку пустых коллекций.
func TestEmptyData(t *testing.T) {
	emptyMock := &mockDataProvider{
		docs: []Document{},
		idx:  map[string][]int{},
	}
	app := New(emptyMock)

	// Проверяем /docs с пустой базой
	req := httptest.NewRequest(http.MethodGet, "/docs", nil)
	rec := httptest.NewRecorder()
	app.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("ожидался 200 OK при пустых docs, получен: %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Хранилище документов пусто") {
		t.Errorf("ожидалось сообщение о пустом хранилище")
	}

	// Проверяем /index с пустым индексом
	req = httptest.NewRequest(http.MethodGet, "/index", nil)
	rec = httptest.NewRecorder()
	app.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("ожидался 200 OK при пустом индексе, получен: %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Поисковый индекс пуст") {
		t.Errorf("ожидалось сообщение о пустом индексе")
	}
}

// TestNilProvider проверяет обработку ошибки, когда провайдер равен nil.
func TestNilProvider(t *testing.T) {
	app := New(nil)

	endpoints := []string{"/docs", "/index"}
	for _, ep := range endpoints {
		req := httptest.NewRequest(http.MethodGet, ep, nil)
		rec := httptest.NewRecorder()

		app.Router().ServeHTTP(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("для %s ожидался статус 500, получен: %d", ep, rec.Code)
		}
	}
}

// TestNotFoundRoute проверяет несуществующий маршрут.
func TestNotFoundRoute(t *testing.T) {
	app, _ := setupTestApp()

	req := httptest.NewRequest(http.MethodGet, "/unknown-page", nil)
	rec := httptest.NewRecorder()

	app.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("ожидался 404 Not Found, получен: %d", rec.Code)
	}
}

// BenchmarkHandleDocs_HTML измеряет скорость рендеринга страницы документов в HTML.
func BenchmarkHandleDocs_HTML(b *testing.B) {
	app, _ := setupTestApp()
	req := httptest.NewRequest(http.MethodGet, "/docs", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		app.Router().ServeHTTP(rec, req)
	}
}

// BenchmarkHandleDocs_JSON измеряет скорость сериализации документов в JSON.
func BenchmarkHandleDocs_JSON(b *testing.B) {
	app, _ := setupTestApp()
	req := httptest.NewRequest(http.MethodGet, "/docs?format=json", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		app.Router().ServeHTTP(rec, req)
	}
}

// BenchmarkHandleIndex_HTML измеряет скорость рендеринга поискового индекса в HTML.
func BenchmarkHandleIndex_HTML(b *testing.B) {
	app, _ := setupTestApp()
	req := httptest.NewRequest(http.MethodGet, "/index", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		app.Router().ServeHTTP(rec, req)
	}
}

// BenchmarkHandleIndex_JSON измеряет скорость сериализации поискового индекса в JSON.
func BenchmarkHandleIndex_JSON(b *testing.B) {
	app, _ := setupTestApp()
	req := httptest.NewRequest(http.MethodGet, "/index?format=json", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		app.Router().ServeHTTP(rec, req)
	}
}

