package memsql_test

import (
	"context"
	"errors"
	"sync"
	"testing"

	"go-core-4/lessons/16-db-apps/homework-16/pkg/db"
	"go-core-4/lessons/16-db-apps/homework-16/pkg/db/memsql"
)

func sampleMovies() []db.Movie {
	return []db.Movie{
		{Title: "Inception", ReleaseYear: 2010, BoxOffice: 836800000, Rating: "PG-13", StudioID: 1, StudioName: "Warner Bros."},
		{Title: "Interstellar", ReleaseYear: 2014, BoxOffice: 701729206, Rating: "PG-13", StudioID: 2, StudioName: "Paramount Pictures"},
		{Title: "Dunkirk", ReleaseYear: 2017, BoxOffice: 526940665, Rating: "PG-13", StudioID: 1, StudioName: "Warner Bros."},
		{Title: "Pulp Fiction", ReleaseYear: 1994, BoxOffice: 213928762, Rating: "PG-18", StudioID: 3, StudioName: "Miramax"},
	}
}

// Проверка добавления массива фильмов и получения всех фильмов (studioID == 0)
func TestAddMovies_And_Movies_All(t *testing.T) {
	ctx := context.Background()
	store := memsql.New()

	movies := sampleMovies()
	err := store.AddMovies(ctx, movies)
	if err != nil {
		t.Fatalf("AddMovies failed: %v", err)
	}

	all, err := store.Movies(ctx, 0)
	if err != nil {
		t.Fatalf("Movies(0) failed: %v", err)
	}

	if len(all) != len(movies) {
		t.Fatalf("expected %d movies, got %d", len(movies), len(all))
	}

	// Проверяем, что ID сгенерированы
	for i, m := range all {
		if m.ID == 0 {
			t.Errorf("movie[%d] ID is 0", i)
		}
		if m.Title != movies[i].Title {
			t.Errorf("movie[%d] expected title %q, got %q", i, movies[i].Title, m.Title)
		}
	}
}

// Проверка фильтрации по ID студии (studioID > 0)
func TestMovies_FilterByStudio(t *testing.T) {
	ctx := context.Background()
	store := memsql.New()

	movies := sampleMovies()
	if err := store.AddMovies(ctx, movies); err != nil {
		t.Fatalf("AddMovies failed: %v", err)
	}

	// Фильтр по StudioID = 1 (Warner Bros.) -> должно быть 2 фильма
	wbMovies, err := store.Movies(ctx, 1)
	if err != nil {
		t.Fatalf("Movies(1) failed: %v", err)
	}
	if len(wbMovies) != 2 {
		t.Fatalf("expected 2 movies for studio 1, got %d", len(wbMovies))
	}
	for _, m := range wbMovies {
		if m.StudioID != 1 {
			t.Errorf("expected studio_id 1, got %d", m.StudioID)
		}
	}

	// Фильтр по несуществующей студии StudioID = 999 -> должно быть 0 фильмов
	nonExisting, err := store.Movies(ctx, 999)
	if err != nil {
		t.Fatalf("Movies(999) failed: %v", err)
	}
	if len(nonExisting) != 0 {
		t.Fatalf("expected 0 movies for non-existing studio, got %d", len(nonExisting))
	}
}

// Проверка обновления фильма
func TestUpdateMovie(t *testing.T) {
	ctx := context.Background()
	store := memsql.New()

	if err := store.AddMovies(ctx, sampleMovies()); err != nil {
		t.Fatalf("AddMovies failed: %v", err)
	}

	all, _ := store.Movies(ctx, 0)
	target := all[0]
	target.Title = "Inception: Extended Edition"
	target.BoxOffice = 900000000

	if err := store.UpdateMovie(ctx, target); err != nil {
		t.Fatalf("UpdateMovie failed: %v", err)
	}

	updated, err := store.MovieByID(ctx, target.ID)
	if err != nil {
		t.Fatalf("MovieByID failed: %v", err)
	}
	if updated.Title != "Inception: Extended Edition" {
		t.Errorf("expected title 'Inception: Extended Edition', got %q", updated.Title)
	}
	if updated.BoxOffice != 900000000 {
		t.Errorf("expected box office 900000000, got %d", updated.BoxOffice)
	}

	// Попытка обновить несуществующий фильм
	fakeMovie := db.Movie{ID: 99999, Title: "Ghost"}
	err = store.UpdateMovie(ctx, fakeMovie)
	if !errors.Is(err, db.ErrNotFound) {
		t.Errorf("expected ErrNotFound for non-existing movie, got: %v", err)
	}
}

// Проверка удаления фильма
func TestDeleteMovie(t *testing.T) {
	ctx := context.Background()
	store := memsql.New()

	if err := store.AddMovies(ctx, sampleMovies()); err != nil {
		t.Fatalf("AddMovies failed: %v", err)
	}

	all, _ := store.Movies(ctx, 0)
	toDelete := all[0]

	if err := store.DeleteMovie(ctx, toDelete.ID); err != nil {
		t.Fatalf("DeleteMovie failed: %v", err)
	}

	// Повторная попытка удаления должна вернуть ErrNotFound
	err := store.DeleteMovie(ctx, toDelete.ID)
	if !errors.Is(err, db.ErrNotFound) {
		t.Errorf("expected ErrNotFound on second delete, got: %v", err)
	}

	// Поиск удаленного фильма должен вернуть ErrNotFound
	_, err = store.MovieByID(ctx, toDelete.ID)
	if !errors.Is(err, db.ErrNotFound) {
		t.Errorf("expected ErrNotFound when getting deleted movie, got: %v", err)
	}

	// Всего фильмов должно остаться на 1 меньше
	remaining, _ := store.Movies(ctx, 0)
	if len(remaining) != len(all)-1 {
		t.Errorf("expected %d movies remaining, got %d", len(all)-1, len(remaining))
	}
}

// Проверка конкурентного доступа
func TestConcurrency(t *testing.T) {
	ctx := context.Background()
	store := memsql.New()

	var wg sync.WaitGroup
	// Параллельная запись
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			_ = store.AddMovies(ctx, []db.Movie{
				{Title: "Film", ReleaseYear: 2000 + idx, StudioID: (idx % 3) + 1},
			})
		}(i)
	}

	// Параллельное чтение
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			_, _ = store.Movies(ctx, idx%4)
		}(i)
	}

	wg.Wait()
}

func BenchmarkMemsql_AddMovies(b *testing.B) {
	ctx := context.Background()
	store := memsql.New()
	sample := []db.Movie{
		{Title: "Bench Movie", ReleaseYear: 2024, BoxOffice: 1000000, Rating: "PG-13", StudioID: 1},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = store.AddMovies(ctx, sample)
	}
}

func BenchmarkMemsql_Movies_All(b *testing.B) {
	ctx := context.Background()
	store := memsql.New()
	_ = store.AddMovies(ctx, sampleMovies())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = store.Movies(ctx, 0)
	}
}

func BenchmarkMemsql_Movies_FilterByStudio(b *testing.B) {
	ctx := context.Background()
	store := memsql.New()
	_ = store.AddMovies(ctx, sampleMovies())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = store.Movies(ctx, 1)
	}
}
