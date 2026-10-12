package ocp_test

import (
	"testing"

	"go-core-4/17-system-design/homework-17/solid/ocp"
)

func sampleMovies() []ocp.Movie {
	return []ocp.Movie{
		{ID: 1, Title: "The Godfather", ReleaseYear: 1972, Rating: "PG-18", BoxOffice: 250000000},
		{ID: 2, Title: "Inception", ReleaseYear: 2010, Rating: "PG-13", BoxOffice: 836000000},
		{ID: 3, Title: "Toy Story", ReleaseYear: 1995, Rating: "PG-10", BoxOffice: 373000000},
		{ID: 4, Title: "The Dark Knight", ReleaseYear: 2008, Rating: "PG-13", BoxOffice: 1006000000},
	}
}

func TestOCP_Filters(t *testing.T) {
	movies := sampleMovies()

	// 1. Фильтр по году (>= 2000)
	modern := ocp.FilterMovies(movies, ocp.YearFilter{MinYear: 2000})
	if len(modern) != 2 {
		t.Fatalf("expected 2 modern movies, got %d", len(modern))
	}

	// 2. Фильтр по рейтингу (PG-13)
	pg13 := ocp.FilterMovies(movies, ocp.RatingFilter{Rating: "PG-13"})
	if len(pg13) != 2 {
		t.Fatalf("expected 2 PG-13 movies, got %d", len(pg13))
	}

	// 3. Композиция (OCP): Фильтр по году >= 2000 И сборам >= 1,000,000,000 без модификации базового кода
	combo := ocp.FilterMovies(movies, ocp.AndSpecification{
		Specs: []ocp.FilterSpecification{
			ocp.YearFilter{MinYear: 2000},
			ocp.BoxOfficeFilter{MinBoxOffice: 1000000000},
		},
	})
	if len(combo) != 1 || combo[0].Title != "The Dark Knight" {
		t.Fatalf("expected 1 movie 'The Dark Knight', got %+v", combo)
	}
}
