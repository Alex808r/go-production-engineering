// Package db определяет интерфейс хранилища данных и модели информационного сайта о фильмах.
// Схема БД основана на задании №15 (таблицы movies, studios).
// Паттерн Repository позволяет легко подменять реализацию
// (in-memory для тестов, PostgreSQL для продакшена).
package db

import "context"

// Movie описывает сущность "Фильм" из схемы БД информационного сайта о фильмах.
type Movie struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	ReleaseYear int    `json:"release_year"`
	BoxOffice   int64  `json:"box_office"`
	Rating      string `json:"rating"` // 'PG-10', 'PG-13', 'PG-18'
	StudioID    int    `json:"studio_id"`
	StudioName  string `json:"studio_name,omitempty"` // заполняется при JOIN со студией
}

// Studio описывает сущность "Киностудия".
type Studio struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Interface определяет контракт взаимодействия с БД информационного сайта о фильмах.
// Соответствует Задаче №1 домашнего задания №16:
// - добавлять массив фильмов (AddMovies);
// - удалять фильм (DeleteMovie);
// - обновлять фильм (UpdateMovie);
// - получать массив фильмов (Movies) с возможностью фильтрации по ID студии.
//   Если ID студии не указан (== 0), возвращаются все фильмы.
type Interface interface {
	// AddMovies добавляет массив фильмов в хранилище.
	AddMovies(ctx context.Context, movies []Movie) error
	// DeleteMovie удаляет фильм по его ID.
	DeleteMovie(ctx context.Context, id int64) error
	// UpdateMovie обновляет данные фильма.
	UpdateMovie(ctx context.Context, m Movie) error
	// Movies возвращает массив фильмов.
	// Если studioID == 0, возвращаются все фильмы.
	// Если studioID > 0, возвращаются фильмы только указанной студии.
	Movies(ctx context.Context, studioID int) ([]Movie, error)
	// MovieByID возвращает фильм по его ID. Возвращает ErrNotFound, если фильм не найден.
	MovieByID(ctx context.Context, id int64) (Movie, error)
}

