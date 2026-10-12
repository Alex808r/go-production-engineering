package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-core-4/17-system-design/app/internal/api"
	"go-core-4/17-system-design/app/internal/models"
)

func TestAPI_BooksEndpoints(t *testing.T) {
	appAPI := api.New(nil)

	// 1. GET /api/v1/books
	req := httptest.NewRequest(http.MethodGet, "/api/v1/books", nil)
	rec := httptest.NewRecorder()
	appAPI.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var books []models.Book
	if err := json.Unmarshal(rec.Body.Bytes(), &books); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}
	if len(books) == 0 {
		t.Fatalf("expected at least 1 initial book")
	}

	// 2. POST /api/v1/books
	newBook := models.Book{Name: "1984", Author: "George Orwell"}
	payload, _ := json.Marshal(newBook)
	postReq := httptest.NewRequest(http.MethodPost, "/api/v1/books", bytes.NewReader(payload))
	postRec := httptest.NewRecorder()
	appAPI.Router().ServeHTTP(postRec, postReq)

	if postRec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", postRec.Code)
	}

	// 3. DELETE /api/v1/books/{id}
	delReq := httptest.NewRequest(http.MethodDelete, "/api/v1/books/1984", nil)
	delRec := httptest.NewRecorder()
	appAPI.Router().ServeHTTP(delRec, delReq)

	if delRec.Code != http.StatusOK {
		t.Fatalf("expected status 200 on delete, got %d", delRec.Code)
	}

	// 4. DELETE non-existent book
	fakeDelReq := httptest.NewRequest(http.MethodDelete, "/api/v1/books/non_existent_book", nil)
	fakeDelRec := httptest.NewRecorder()
	appAPI.Router().ServeHTTP(fakeDelRec, fakeDelReq)

	if fakeDelRec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", fakeDelRec.Code)
	}
}
