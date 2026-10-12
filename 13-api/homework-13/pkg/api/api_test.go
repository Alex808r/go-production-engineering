package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"go-core-4/13-api/homework-13/pkg/storage"
)

var (
	testAPI *API
	testDB  *storage.DB
)

func TestMain(m *testing.M) {
	testDB = storage.New()

	// Наполняем тестовыми данными
	testDB.Add(storage.Document{
		ID:    1,
		URL:   "https://go.dev/doc/tutorial",
		Title: "Go Tutorial and Basics",
		Body:  "Getting started with Go programming language",
	})
	testDB.Add(storage.Document{
		ID:    2,
		URL:   "https://go.dev/doc/effective_go",
		Title: "Effective Go Concurrency",
		Body:  "Tips on goroutines and channel design",
	})

	testAPI = New(testDB)
	os.Exit(m.Run())
}

func TestAPI_Search(t *testing.T) {
	tests := []struct {
		name       string
		queryURL   string
		wantCode   int
		wantMinLen int
	}{
		{
			name:       "Search with match",
			queryURL:   "/api/v1/search?q=concurrency",
			wantCode:   http.StatusOK,
			wantMinLen: 1,
		},
		{
			name:       "Search with alternative query param",
			queryURL:   "/api/v1/search?query=tutorial",
			wantCode:   http.StatusOK,
			wantMinLen: 1,
		},
		{
			name:       "Search with no matches",
			queryURL:   "/api/v1/search?q=nonexistentkeyword",
			wantCode:   http.StatusOK,
			wantMinLen: 0,
		},
		{
			name:       "Search with empty query",
			queryURL:   "/api/v1/search",
			wantCode:   http.StatusBadRequest,
			wantMinLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.queryURL, nil)
			rr := httptest.NewRecorder()

			testAPI.Router().ServeHTTP(rr, req)

			if rr.Code != tt.wantCode {
				t.Fatalf("Code = %d, want %d", rr.Code, tt.wantCode)
			}

			if tt.wantCode == http.StatusOK {
				var docs []storage.Document
				if err := json.Unmarshal(rr.Body.Bytes(), &docs); err != nil {
					t.Fatalf("Unmarshal error: %v", err)
				}
				if len(docs) < tt.wantMinLen {
					t.Errorf("len(docs) = %d, want >= %d", len(docs), tt.wantMinLen)
				}
			}
		})
	}
}

func TestAPI_Docs(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/docs", nil)
	rr := httptest.NewRecorder()

	testAPI.Router().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Code = %d, want %d", rr.Code, http.StatusOK)
	}

	var docs []storage.Document
	if err := json.Unmarshal(rr.Body.Bytes(), &docs); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	if len(docs) == 0 {
		t.Error("expected non-empty list of documents")
	}
}

func TestAPI_CreateDoc(t *testing.T) {
	t.Run("Valid Document Creation", func(t *testing.T) {
		newDoc := storage.Document{
			URL:   "https://go.dev/ref/spec",
			Title: "The Go Programming Language Specification",
			Body:  "Detailed specification of the Go language.",
		}
		payload, _ := json.Marshal(newDoc)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/docs", bytes.NewBuffer(payload))
		rr := httptest.NewRecorder()

		testAPI.Router().ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Fatalf("Code = %d, want %d. Body: %s", rr.Code, http.StatusCreated, rr.Body.String())
		}

		var created storage.Document
		if err := json.Unmarshal(rr.Body.Bytes(), &created); err != nil {
			t.Fatalf("Unmarshal error: %v", err)
		}
		if created.ID <= 0 {
			t.Errorf("Expected assigned ID > 0, got %d", created.ID)
		}
		if created.Title != newDoc.Title {
			t.Errorf("Title = %s, want %s", created.Title, newDoc.Title)
		}
	})

	t.Run("Invalid JSON Payload", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/docs", bytes.NewBufferString("{invalid-json"))
		rr := httptest.NewRecorder()

		testAPI.Router().ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Code = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("Empty Required Fields", func(t *testing.T) {
		emptyDoc := storage.Document{Body: "No title or url"}
		payload, _ := json.Marshal(emptyDoc)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/docs", bytes.NewBuffer(payload))
		rr := httptest.NewRecorder()

		testAPI.Router().ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Code = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})
}

func TestAPI_DocByID(t *testing.T) {
	t.Run("Existing Document", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/docs/2", nil)
		rr := httptest.NewRecorder()

		testAPI.Router().ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("Code = %d, want %d", rr.Code, http.StatusOK)
		}

		var doc storage.Document
		if err := json.Unmarshal(rr.Body.Bytes(), &doc); err != nil {
			t.Fatalf("Unmarshal error: %v", err)
		}
		if doc.ID != 2 {
			t.Errorf("doc.ID = %d, want 2", doc.ID)
		}
	})

	t.Run("Non-existing Document", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/docs/99999", nil)
		rr := httptest.NewRecorder()

		testAPI.Router().ServeHTTP(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("Code = %d, want %d", rr.Code, http.StatusNotFound)
		}
	})
}

func TestAPI_UpdateDoc(t *testing.T) {
	t.Run("Successful Update", func(t *testing.T) {
		updated := storage.Document{
			URL:   "https://go.dev/doc/effective_go_v2",
			Title: "Effective Go 2.0",
			Body:  "Updated concurrency tips",
		}
		payload, _ := json.Marshal(updated)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/docs/2", bytes.NewBuffer(payload))
		rr := httptest.NewRecorder()

		testAPI.Router().ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("Code = %d, want %d", rr.Code, http.StatusOK)
		}

		var res storage.Document
		if err := json.Unmarshal(rr.Body.Bytes(), &res); err != nil {
			t.Fatalf("Unmarshal error: %v", err)
		}
		if res.Title != "Effective Go 2.0" {
			t.Errorf("Title = %s, want Effective Go 2.0", res.Title)
		}
	})

	t.Run("Update Non-existing Document", func(t *testing.T) {
		doc := storage.Document{Title: "Missing", URL: "http://missing"}
		payload, _ := json.Marshal(doc)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/docs/99999", bytes.NewBuffer(payload))
		rr := httptest.NewRecorder()

		testAPI.Router().ServeHTTP(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("Code = %d, want %d", rr.Code, http.StatusNotFound)
		}
	})

	t.Run("Update with Invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/v1/docs/2", bytes.NewBufferString("{bad-json"))
		rr := httptest.NewRecorder()

		testAPI.Router().ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Code = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})
}

func TestAPI_DeleteDoc(t *testing.T) {
	// Создаем отдельный документ для удаления
	target := testDB.Add(storage.Document{
		URL:   "https://go.dev/delete-me",
		Title: "Document to Delete",
		Body:  "Will be deleted",
	})

	t.Run("Delete Existing Document", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/docs/"+strconvItoa(target.ID), nil)
		rr := httptest.NewRecorder()

		testAPI.Router().ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("Code = %d, want %d", rr.Code, http.StatusOK)
		}
	})

	t.Run("Delete Already Deleted Document", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/docs/"+strconvItoa(target.ID), nil)
		rr := httptest.NewRecorder()

		testAPI.Router().ServeHTTP(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("Code = %d, want %d", rr.Code, http.StatusNotFound)
		}
	})
}

func TestAPI_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/search", nil)
	rr := httptest.NewRecorder()

	testAPI.Router().ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Code = %d, want %d", rr.Code, http.StatusMethodNotAllowed)
	}
}

func strconvItoa(n int) string {
	var b [20]byte
	i := len(b) - 1
	for n >= 10 {
		b[i] = byte('0' + n%10)
		n /= 10
		i--
	}
	b[i] = byte('0' + n)
	return string(b[i:])
}

// BenchmarkAPI_Search измеряет скорость сквозного поиска через HTTP API.
func BenchmarkAPI_Search(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/search?q=tutorial", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		testAPI.Router().ServeHTTP(rr, req)
	}
}

// BenchmarkAPI_Docs измеряет скорость получения всех документов через HTTP API.
func BenchmarkAPI_Docs(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/docs", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		testAPI.Router().ServeHTTP(rr, req)
	}
}

// BenchmarkAPI_DocByID измеряет скорость получения документа по ID через HTTP API.
func BenchmarkAPI_DocByID(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/docs/1", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		testAPI.Router().ServeHTTP(rr, req)
	}
}

