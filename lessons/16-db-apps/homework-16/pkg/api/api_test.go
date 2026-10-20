package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-core-4/lessons/16-db-apps/homework-16/pkg/api"
	"go-core-4/lessons/16-db-apps/homework-16/pkg/db"
	"go-core-4/lessons/16-db-apps/homework-16/pkg/db/memsql"
)

func setupTestServer() (*api.API, *memsql.DB) {
	store := memsql.New()
	server := api.New(store)
	return server, store
}

func TestHandleMovies_Empty(t *testing.T) {
	server, _ := setupTestServer()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/movies", nil)
	rec := httptest.NewRecorder()

	server.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var movies []db.Movie
	if err := json.Unmarshal(rec.Body.Bytes(), &movies); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(movies) != 0 {
		t.Errorf("expected 0 movies, got %d", len(movies))
	}
}

func TestHandleMovies_WithFilterByStudio(t *testing.T) {
	server, store := setupTestServer()
	ctx := context.Background()

	_ = store.AddMovies(ctx, []db.Movie{
		{Title: "Inception", ReleaseYear: 2010, StudioID: 1},
		{Title: "Interstellar", ReleaseYear: 2014, StudioID: 2},
		{Title: "Dunkirk", ReleaseYear: 2017, StudioID: 1},
	})

	// 1. Без фильтра (studio_id не указан) -> все 3 фильма
	reqAll := httptest.NewRequest(http.MethodGet, "/api/v1/movies", nil)
	recAll := httptest.NewRecorder()
	server.Router().ServeHTTP(recAll, reqAll)

	var all []db.Movie
	_ = json.Unmarshal(recAll.Body.Bytes(), &all)
	if len(all) != 3 {
		t.Fatalf("expected 3 movies without filter, got %d", len(all))
	}

	// 2. С фильтром studio_id=1 -> ровно 2 фильма
	reqFilter := httptest.NewRequest(http.MethodGet, "/api/v1/movies?studio_id=1", nil)
	recFilter := httptest.NewRecorder()
	server.Router().ServeHTTP(recFilter, reqFilter)

	var filtered []db.Movie
	_ = json.Unmarshal(recFilter.Body.Bytes(), &filtered)
	if len(filtered) != 2 {
		t.Fatalf("expected 2 movies for studio 1, got %d", len(filtered))
	}

	// 3. С фильтром studio_id=0 -> все 3 фильма (согласно ТЗ)
	reqZero := httptest.NewRequest(http.MethodGet, "/api/v1/movies?studio_id=0", nil)
	recZero := httptest.NewRecorder()
	server.Router().ServeHTTP(recZero, reqZero)

	var zero []db.Movie
	_ = json.Unmarshal(recZero.Body.Bytes(), &zero)
	if len(zero) != 3 {
		t.Fatalf("expected 3 movies with studio_id=0, got %d", len(zero))
	}
}

func TestHandleMovieByID(t *testing.T) {
	server, store := setupTestServer()
	ctx := context.Background()

	_ = store.AddMovies(ctx, []db.Movie{
		{Title: "Inception", ReleaseYear: 2010, Rating: "PG-13", StudioID: 1},
	})
	movies, _ := store.Movies(ctx, 0)
	id := movies[0].ID

	// Существующий фильм
	req := httptest.NewRequest(http.MethodGet, "/api/v1/movies/1", nil)
	rec := httptest.NewRecorder()
	server.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var m db.Movie
	_ = json.Unmarshal(rec.Body.Bytes(), &m)
	if m.ID != id || m.Title != "Inception" {
		t.Errorf("unexpected movie: %+v", m)
	}

	// Несуществующий фильм
	reqNotFound := httptest.NewRequest(http.MethodGet, "/api/v1/movies/999", nil)
	recNotFound := httptest.NewRecorder()
	server.Router().ServeHTTP(recNotFound, reqNotFound)

	if recNotFound.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", recNotFound.Code)
	}
}

func TestHandleCreateMovies_BatchAndSingle(t *testing.T) {
	server, store := setupTestServer()
	ctx := context.Background()

	// 1. Добавление массива фильмов
	payload := `[
		{"title": "The Matrix", "release_year": 1999, "box_office": 467200000, "rating": "PG-18", "studio_id": 1},
		{"title": "The Matrix Reloaded", "release_year": 2003, "box_office": 741800000, "rating": "PG-18", "studio_id": 1}
	]`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/movies", bytes.NewBufferString(payload))
	rec := httptest.NewRecorder()
	server.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d: %s", rec.Code, rec.Body.String())
	}

	movies, _ := store.Movies(ctx, 0)
	if len(movies) != 2 {
		t.Fatalf("expected 2 movies in store, got %d", len(movies))
	}

	// 2. Добавление одного фильма
	singlePayload := `{"title": "The Dark Knight", "release_year": 2008, "box_office": 1006000000, "rating": "PG-13", "studio_id": 1}`
	reqSingle := httptest.NewRequest(http.MethodPost, "/api/v1/movies", bytes.NewBufferString(singlePayload))
	recSingle := httptest.NewRecorder()
	server.Router().ServeHTTP(recSingle, reqSingle)

	if recSingle.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", recSingle.Code)
	}

	movies, _ = store.Movies(ctx, 0)
	if len(movies) != 3 {
		t.Fatalf("expected 3 movies in store, got %d", len(movies))
	}
}

func TestHandleUpdateMovie(t *testing.T) {
	server, store := setupTestServer()
	ctx := context.Background()

	_ = store.AddMovies(ctx, []db.Movie{
		{Title: "Original Title", ReleaseYear: 2020, Rating: "PG-13", StudioID: 1},
	})
	movies, _ := store.Movies(ctx, 0)
	id := movies[0].ID

	updatePayload := `{"title": "Updated Title", "release_year": 2021, "box_office": 500, "rating": "PG-18", "studio_id": 2}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/movies/1", bytes.NewBufferString(updatePayload))
	rec := httptest.NewRecorder()
	server.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
	}

	updated, _ := store.MovieByID(ctx, id)
	if updated.Title != "Updated Title" || updated.ReleaseYear != 2021 {
		t.Errorf("movie was not updated correctly: %+v", updated)
	}
}

func TestHandleDeleteMovie(t *testing.T) {
	server, store := setupTestServer()
	ctx := context.Background()

	_ = store.AddMovies(ctx, []db.Movie{
		{Title: "To Be Deleted", ReleaseYear: 2020, Rating: "PG-13", StudioID: 1},
	})

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/movies/1", nil)
	rec := httptest.NewRecorder()
	server.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", rec.Code)
	}

	remaining, _ := store.Movies(ctx, 0)
	if len(remaining) != 0 {
		t.Errorf("expected 0 movies, got %d", len(remaining))
	}
}

func BenchmarkAPI_GetMovies(b *testing.B) {
	server, store := setupTestServer()
	_ = store.AddMovies(context.Background(), []db.Movie{
		{Title: "Inception", ReleaseYear: 2010, StudioID: 1},
		{Title: "Interstellar", ReleaseYear: 2014, StudioID: 2},
	})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/movies", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		server.Router().ServeHTTP(rec, req)
	}
}

func BenchmarkAPI_CreateMovie(b *testing.B) {
	server, _ := setupTestServer()
	payload := `{"title": "Bench Film", "release_year": 2024, "box_office": 1000000, "rating": "PG-13", "studio_id": 1}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/movies", bytes.NewBufferString(payload))
		rec := httptest.NewRecorder()
		server.Router().ServeHTTP(rec, req)
	}
}

