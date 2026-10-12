// Package usecase содержит сценарии использования приложения (Application Business Rules).
// Реализует паттерн Input/Output Ports из Чистой Архитектуры:
// - Входной порт (Input Port): интерфейс, через который внешние адаптеры (контроллеры) вызывают сценарии.
// - Выходной порт (Output Port): интерфейс вторичного порта (репозиторий, логер, уведомитель),
//   который реализуется во внешнем слое адаптеров (DIP).
package usecase

import (
	"context"
	"errors"
	"fmt"

	"go-core-4/17-system-design/homework-17/cleanarch/entity"
)

var (
	ErrNotFound = errors.New("фильм не найден")
)

// DTO для передачи данных между слоями без утечки сущностей и деталей БД.
type CreateMovieDTO struct {
	Title       string `json:"title"`
	ReleaseYear int    `json:"release_year"`
	BoxOffice   int64  `json:"box_office"`
	Rating      string `json:"rating"`
	StudioID    int    `json:"studio_id"`
}

type MovieResponseDTO struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	ReleaseYear int    `json:"release_year"`
	BoxOffice   int64  `json:"box_office"`
	Rating      string `json:"rating"`
	StudioID    int    `json:"studio_id"`
	StudioName  string `json:"studio_name,omitempty"`
}

// Output Port (Вторичный порт): контракт хранилища данных.
type MovieRepository interface {
	Save(ctx context.Context, m entity.Movie) (int64, error)
	SaveBatch(ctx context.Context, movies []entity.Movie) error
	FindByID(ctx context.Context, id int64) (entity.Movie, error)
	FindAll(ctx context.Context, studioID int) ([]entity.Movie, error)
	Update(ctx context.Context, m entity.Movie) error
	Delete(ctx context.Context, id int64) error
}

// Input Port (Первичный порт): контракт доступных бизнес-сценариев.
type MovieUseCase interface {
	CreateMovie(ctx context.Context, dto CreateMovieDTO) (MovieResponseDTO, error)
	CreateMoviesBatch(ctx context.Context, dtos []CreateMovieDTO) error
	GetMovieByID(ctx context.Context, id int64) (MovieResponseDTO, error)
	ListMovies(ctx context.Context, studioID int) ([]MovieResponseDTO, error)
	DeleteMovie(ctx context.Context, id int64) error
}

// interactor — реализация бизнес-логики сценариев использования.
type interactor struct {
	repo MovieRepository
}

// NewMovieUseCase создает интерактор бизнес-сценариев.
func NewMovieUseCase(r MovieRepository) MovieUseCase {
	return &interactor{repo: r}
}

func (uc *interactor) CreateMovie(ctx context.Context, dto CreateMovieDTO) (MovieResponseDTO, error) {
	m, err := entity.NewMovie(0, dto.Title, dto.ReleaseYear, dto.BoxOffice, dto.Rating, dto.StudioID)
	if err != nil {
		return MovieResponseDTO{}, fmt.Errorf("ошибка валидации: %w", err)
	}

	id, err := uc.repo.Save(ctx, m)
	if err != nil {
		return MovieResponseDTO{}, fmt.Errorf("ошибка сохранения: %w", err)
	}
	m.ID = id

	return toDTO(m), nil
}

func (uc *interactor) CreateMoviesBatch(ctx context.Context, dtos []CreateMovieDTO) error {
	var movies []entity.Movie
	for _, dto := range dtos {
		m, err := entity.NewMovie(0, dto.Title, dto.ReleaseYear, dto.BoxOffice, dto.Rating, dto.StudioID)
		if err != nil {
			return fmt.Errorf("ошибка валидации фильма %q: %w", dto.Title, err)
		}
		movies = append(movies, m)
	}
	return uc.repo.SaveBatch(ctx, movies)
}

func (uc *interactor) GetMovieByID(ctx context.Context, id int64) (MovieResponseDTO, error) {
	m, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return MovieResponseDTO{}, err
	}
	return toDTO(m), nil
}

func (uc *interactor) ListMovies(ctx context.Context, studioID int) ([]MovieResponseDTO, error) {
	movies, err := uc.repo.FindAll(ctx, studioID)
	if err != nil {
		return nil, err
	}
	res := make([]MovieResponseDTO, 0, len(movies))
	for _, m := range movies {
		res = append(res, toDTO(m))
	}
	return res, nil
}

func (uc *interactor) DeleteMovie(ctx context.Context, id int64) error {
	return uc.repo.Delete(ctx, id)
}

func toDTO(m entity.Movie) MovieResponseDTO {
	return MovieResponseDTO{
		ID:          m.ID,
		Title:       m.Title,
		ReleaseYear: m.ReleaseYear,
		BoxOffice:   m.BoxOffice,
		Rating:      m.Rating,
		StudioID:    m.StudioID,
		StudioName:  m.StudioName,
	}
}
