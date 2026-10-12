// Package srp демонстрирует принцип единственной ответственности (Single Responsibility Principle) в Go.
// Формулировка Р. Мартина: "У модуля должна быть одна и только одна причина для изменения"
// (модуль отвечает перед одним и только одним актором).
package srp

import (
	"fmt"
	"strings"
	"time"
)

// Movie представляет доменную модель фильма.
type Movie struct {
	ID          int64
	Title       string
	ReleaseYear int
	BoxOffice   int64
}

// MovieService отвечает исключительно за бизнес-логику управления фильмами.
// Единственный актор: менеджер каталога фильмов.
type MovieService struct {
	movies []Movie
}

// NewMovieService создает сервис фильмов.
func NewMovieService() *MovieService {
	return &MovieService{
		movies: make([]Movie, 0),
	}
}

// AddMovie добавляет фильм после валидации бизнес-правил.
func (s *MovieService) AddMovie(m Movie) error {
	if strings.TrimSpace(m.Title) == "" {
		return fmt.Errorf("название фильма не может быть пустым")
	}
	if m.ReleaseYear < 1800 || m.ReleaseYear > time.Now().Year()+5 {
		return fmt.Errorf("недопустимый год выхода: %d", m.ReleaseYear)
	}
	s.movies = append(s.movies, m)
	return nil
}

// AllMovies возвращает список фильмов.
func (s *MovieService) AllMovies() []Movie {
	res := make([]Movie, len(s.movies))
	copy(res, s.movies)
	return res
}

// ReportFormatter отвечает исключительно за форматирование отчетов.
// Единственный актор: аналитик / отдел отчетности.
// Вынесение форматирования из MovieService предотвращает смешение причин для изменения.
type ReportFormatter struct{}

// ToPlainText формирует текстовое представление отчета.
func (ReportFormatter) ToPlainText(movies []Movie) string {
	var b strings.Builder
	b.WriteString("Отчет по фильмам:\n")
	for _, m := range movies {
		b.WriteString(fmt.Sprintf("- %s (%d): $%d\n", m.Title, m.ReleaseYear, m.BoxOffice))
	}
	return b.String()
}

// ToCSV формирует CSV-представление отчета.
func (ReportFormatter) ToCSV(movies []Movie) string {
	var b strings.Builder
	b.WriteString("id,title,year,box_office\n")
	for _, m := range movies {
		b.WriteString(fmt.Sprintf("%d,%s,%d,%d\n", m.ID, m.Title, m.ReleaseYear, m.BoxOffice))
	}
	return b.String()
}
