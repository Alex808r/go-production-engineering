package api_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-core-4/lessons/16-db-apps/pkg/api"
	"go-core-4/lessons/16-db-apps/pkg/db"
	"go-core-4/lessons/16-db-apps/pkg/db/memsql"
)

func TestAPI_BooksHandler_Empty(t *testing.T) {
	store := memsql.New()
	handler := api.New(store)

	req := httptest.NewRequest(http.MethodGet, "/books", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var books []db.Book
	if err := json.Unmarshal(rr.Body.Bytes(), &books); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}
	if len(books) != 0 {
		t.Fatalf("expected 0 books, got %d", len(books))
	}
}

func TestAPI_BooksHandler_WithData(t *testing.T) {
	store := memsql.New()
	_ = store.AddBooks(context.Background(), []db.Book{
		{Title: "Clean Architecture", Year: 2017},
		{Title: "The Go Programming Language", Year: 2015},
	})
	handler := api.New(store)

	req := httptest.NewRequest(http.MethodGet, "/books", nil)
	rr := httptest.NewRecorder()

	handler.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var books []db.Book
	if err := json.Unmarshal(rr.Body.Bytes(), &books); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}
	if len(books) != 2 {
		t.Fatalf("expected 2 books, got %d", len(books))
	}
	if books[0].Title != "Clean Architecture" {
		t.Errorf("expected 'Clean Architecture', got %s", books[0].Title)
	}
}

func BenchmarkAPI_ServeHTTP(b *testing.B) {
	store := memsql.New()
	_ = store.AddBooks(context.Background(), []db.Book{
		{Title: "Benchmark Book", Year: 2024},
	})
	handler := api.New(store)
	req := httptest.NewRequest(http.MethodGet, "/books", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
	}
}
