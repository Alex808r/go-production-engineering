package srp_test

import (
	"strings"
	"testing"

	"go-core-4/lessons/17-system-design/homework-17/solid/srp"
)

func TestMovieService_SRP(t *testing.T) {
	svc := srp.NewMovieService()

	// Валидный фильм
	err := svc.AddMovie(srp.Movie{ID: 1, Title: "Inception", ReleaseYear: 2010, BoxOffice: 836000000})
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}

	// Невалидный фильм
	if err := svc.AddMovie(srp.Movie{Title: ""}); err == nil {
		t.Error("expected error for empty title")
	}

	if err := svc.AddMovie(srp.Movie{Title: "Old", ReleaseYear: 1700}); err == nil {
		t.Error("expected error for year < 1800")
	}

	movies := svc.AllMovies()
	if len(movies) != 1 {
		t.Fatalf("expected 1 movie, got %d", len(movies))
	}

	// Тестирование форматирования отчетов (отдельная ответственность)
	formatter := srp.ReportFormatter{}
	txt := formatter.ToPlainText(movies)
	if !strings.Contains(txt, "Inception") {
		t.Errorf("expected plain text report to contain 'Inception', got: %s", txt)
	}

	csv := formatter.ToCSV(movies)
	if !strings.Contains(csv, "1,Inception,2010,836000000") {
		t.Errorf("expected CSV report to contain CSV line, got: %s", csv)
	}
}
