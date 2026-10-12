// Package controller содержит первичные адаптеры (Primary Adapters / Driving Adapters)
// для взаимодействия с внешним миром через HTTP.
package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"go-core-4/17-system-design/homework-17/cleanarch/usecase"
)

// MovieHTTPController преобразует HTTP-запросы в вызовы Input Port (usecase.MovieUseCase).
type MovieHTTPController struct {
	useCase usecase.MovieUseCase
}

func NewMovieHTTPController(uc usecase.MovieUseCase) *MovieHTTPController {
	return &MovieHTTPController{useCase: uc}
}

// RegisterRoutes регистрирует маршруты в gorilla/mux роутере.
func (c *MovieHTTPController) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/v1/movies", c.HandleListMovies).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/movies/{id:[0-9]+}", c.HandleGetMovie).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/movies", c.HandleCreateMovie).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/movies/batch", c.HandleCreateMoviesBatch).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/movies/{id:[0-9]+}", c.HandleDeleteMovie).Methods(http.MethodDelete)
}

func (c *MovieHTTPController) HandleListMovies(w http.ResponseWriter, r *http.Request) {
	studioID := 0
	if param := r.URL.Query().Get("studio_id"); param != "" {
		if id, err := strconv.Atoi(param); err == nil && id > 0 {
			studioID = id
		}
	}

	movies, err := c.useCase.ListMovies(r.Context(), studioID)
	if err != nil {
		c.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	c.respondJSON(w, http.StatusOK, movies)
}

func (c *MovieHTTPController) HandleGetMovie(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		c.respondError(w, http.StatusBadRequest, "неверный ID фильма")
		return
	}

	movie, err := c.useCase.GetMovieByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.respondError(w, http.StatusNotFound, "фильм не найден")
			return
		}
		c.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	c.respondJSON(w, http.StatusOK, movie)
}

func (c *MovieHTTPController) HandleCreateMovie(w http.ResponseWriter, r *http.Request) {
	var dto usecase.CreateMovieDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		c.respondError(w, http.StatusBadRequest, "некорректный JSON")
		return
	}

	created, err := c.useCase.CreateMovie(r.Context(), dto)
	if err != nil {
		c.respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	c.respondJSON(w, http.StatusCreated, created)
}

func (c *MovieHTTPController) HandleCreateMoviesBatch(w http.ResponseWriter, r *http.Request) {
	var dtos []usecase.CreateMovieDTO
	if err := json.NewDecoder(r.Body).Decode(&dtos); err != nil {
		c.respondError(w, http.StatusBadRequest, "некорректный JSON: ожидался массив")
		return
	}

	if err := c.useCase.CreateMoviesBatch(r.Context(), dtos); err != nil {
		c.respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	c.respondJSON(w, http.StatusCreated, map[string]interface{}{
		"status": "created",
		"count":  len(dtos),
	})
}

func (c *MovieHTTPController) HandleDeleteMovie(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		c.respondError(w, http.StatusBadRequest, "неверный ID фильма")
		return
	}

	if err := c.useCase.DeleteMovie(r.Context(), id); err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.respondError(w, http.StatusNotFound, "фильм не найден")
			return
		}
		c.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (c *MovieHTTPController) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func (c *MovieHTTPController) respondError(w http.ResponseWriter, status int, msg string) {
	c.respondJSON(w, status, map[string]string{
		"error": strings.TrimSpace(msg),
	})
}
