// Package entity содержит сущности ядра предметной области (Enterprise Business Rules).
// По Чистой архитектуре Роберта Мартина (Clean Architecture), этот слой является самым внутренним:
// он не зависит ни от каких внешних библиотек, фреймворков, баз данных или протоколов передачи данных.
package entity

import (
	"errors"
	"strings"
	"time"
)

// Доменные ошибки уровня сущности.
var (
	ErrEmptyTitle    = errors.New("название фильма не может быть пустым")
	ErrInvalidYear   = errors.New("год выпуска фильма не может быть меньше 1800 или в далеком будущем")
	ErrInvalidRating = errors.New("недопустимый возрастной рейтинг: разрешены только PG-10, PG-13, PG-18")
	ErrNegativeBox   = errors.New("кассовые сборы не могут быть отрицательными")
)

// AllowedRatings — допустимые возрастные рейтинги по правилам кинопроката.
var AllowedRatings = map[string]bool{
	"PG-10": true,
	"PG-13": true,
	"PG-18": true,
}

// Movie представляет доменную сущность "Фильм".
type Movie struct {
	ID          int64
	Title       string
	ReleaseYear int
	BoxOffice   int64
	Rating      string
	StudioID    int
	StudioName  string
}

// NewMovie — фабричный метод создания и валидации фильма по инвариантам предметной области.
func NewMovie(id int64, title string, year int, boxOffice int64, rating string, studioID int) (Movie, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return Movie{}, ErrEmptyTitle
	}

	currentYear := time.Now().Year()
	if year < 1800 || year > currentYear+10 {
		return Movie{}, ErrInvalidYear
	}

	if boxOffice < 0 {
		return Movie{}, ErrNegativeBox
	}

	if !AllowedRatings[rating] {
		return Movie{}, ErrInvalidRating
	}

	return Movie{
		ID:          id,
		Title:       title,
		ReleaseYear: year,
		BoxOffice:   boxOffice,
		Rating:      rating,
		StudioID:    studioID,
	}, nil
}
