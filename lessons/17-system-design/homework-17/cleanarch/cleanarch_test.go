package cleanarch_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"

	"go-core-4/lessons/17-system-design/homework-17/cleanarch/adapter/controller"
	"go-core-4/lessons/17-system-design/homework-17/cleanarch/adapter/repository"
	"go-core-4/lessons/17-system-design/homework-17/cleanarch/usecase"
)

func setupApp() (*mux.Router, usecase.MovieUseCase) {
	// Composition Root: сборка зависимостей
	repo := repository.NewMemoryMovieRepository()
	uc := usecase.NewMovieUseCase(repo)
	ctrl := controller.NewMovieHTTPController(uc)

	r := mux.NewRouter()
	ctrl.RegisterRoutes(r)
	return r, uc
}

func TestCleanArchitecture_FullFlow(t *testing.T) {
	router, uc := setupApp()
	ctx := context.Background()

	// 1. Создание через usecase
	created, err := uc.CreateMovie(ctx, usecase.CreateMovieDTO{
		Title:       "Interstellar",
		ReleaseYear: 2014,
		BoxOffice:   701000000,
		Rating:      "PG-13",
		StudioID:    2,
	})
	if err != nil {
		t.Fatalf("failed to create movie: %v", err)
	}

	// 2. Получение через HTTP GET
	req := httptest.NewRequest(http.MethodGet, "/api/v1/movies/1", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}

	var resDTO usecase.MovieResponseDTO
	if err := json.Unmarshal(rec.Body.Bytes(), &resDTO); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if resDTO.Title != "Interstellar" || resDTO.ID != created.ID {
		t.Errorf("unexpected movie response: %+v", resDTO)
	}

	// 3. Создание пакета фильмов через HTTP POST /api/v1/movies/batch
	batchPayload := `[
		{"title":"Inception","release_year":2010,"box_office":836000000,"rating":"PG-13","studio_id":1},
		{"title":"Dunkirk","release_year":2017,"box_office":526000000,"rating":"PG-13","studio_id":1}
	]`
	reqBatch := httptest.NewRequest(http.MethodPost, "/api/v1/movies/batch", bytes.NewBufferString(batchPayload))
	recBatch := httptest.NewRecorder()
	router.ServeHTTP(recBatch, reqBatch)

	if recBatch.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created for batch, got %d", recBatch.Code)
	}

	// 4. Фильтрация по studio_id = 1
	reqFilter := httptest.NewRequest(http.MethodGet, "/api/v1/movies?studio_id=1", nil)
	recFilter := httptest.NewRecorder()
	router.ServeHTTP(recFilter, reqFilter)

	var filtered []usecase.MovieResponseDTO
	_ = json.Unmarshal(recFilter.Body.Bytes(), &filtered)
	if len(filtered) != 2 {
		t.Fatalf("expected 2 movies for studio 1, got %d", len(filtered))
	}

	// 5. Удаление фильма
	reqDel := httptest.NewRequest(http.MethodDelete, "/api/v1/movies/1", nil)
	recDel := httptest.NewRecorder()
	router.ServeHTTP(recDel, reqDel)

	if recDel.Code != http.StatusNoContent {
		t.Fatalf("expected 204 No Content, got %d", recDel.Code)
	}
}

func BenchmarkCleanArchitecture_UseCase(b *testing.B) {
	_, uc := setupApp()
	ctx := context.Background()
	dto := usecase.CreateMovieDTO{
		Title:       "Benchmark Movie",
		ReleaseYear: 2024,
		BoxOffice:   1000000,
		Rating:      "PG-13",
		StudioID:    1,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = uc.CreateMovie(ctx, dto)
	}
}

func BenchmarkCleanArchitecture_HTTP(b *testing.B) {
	router, uc := setupApp()
	ctx := context.Background()
	_, _ = uc.CreateMovie(ctx, usecase.CreateMovieDTO{
		Title:       "Test",
		ReleaseYear: 2020,
		Rating:      "PG-13",
		StudioID:    1,
	})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/movies/1", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
	}
}

