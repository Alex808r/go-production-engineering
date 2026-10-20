// Package api реализует REST API для информационного сайта о фильмах.
// Зависит только от интерфейса db.Interface — реализацию хранилища
// можно подменять (in-memory для тестов, PostgreSQL для продакшена).
package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"go-core-4/lessons/16-db-apps/homework-16/pkg/db"
)

// API инкапсулирует HTTP-маршрутизатор и хранилище данных.
type API struct {
	router *mux.Router
	db     db.Interface
}

// New создает новый экземпляр API с зарегистрированными маршрутами и middleware.
func New(store db.Interface) *API {
	api := &API{
		router: mux.NewRouter(),
		db:     store,
	}

	// Middleware.
	api.router.Use(api.loggingMiddleware)
	api.router.Use(api.contentTypeMiddleware)

	// Маршруты CRUD для фильмов.
	api.router.HandleFunc("/api/v1/movies", api.handleMovies).Methods(http.MethodGet)
	api.router.HandleFunc("/api/v1/movies/{id:[0-9]+}", api.handleMovieByID).Methods(http.MethodGet)
	api.router.HandleFunc("/api/v1/movies", api.handleCreateMovies).Methods(http.MethodPost)
	api.router.HandleFunc("/api/v1/movies/{id:[0-9]+}", api.handleUpdateMovie).Methods(http.MethodPut)
	api.router.HandleFunc("/api/v1/movies/{id:[0-9]+}", api.handleDeleteMovie).Methods(http.MethodDelete)

	return api
}

// Router возвращает маршрутизатор для передачи в http.Server.
func (api *API) Router() http.Handler {
	return api.router
}

// --- Middleware ---

// loggingMiddleware логирует метод, URL и продолжительность обработки каждого запроса.
func (api *API) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("[%s] %s — %v", r.Method, r.URL.String(), time.Since(start))
	})
}

// contentTypeMiddleware устанавливает заголовок Content-Type: application/json.
func (api *API) contentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

// --- Обработчики маршрутов ---

// handleMovies обрабатывает GET /api/v1/movies?studio_id=X.
// При указании studio_id фильтрует по студии, при 0 или отсутствии параметра возвращает все фильмы.
func (api *API) handleMovies(w http.ResponseWriter, r *http.Request) {
	studioID := 0
	if val := r.URL.Query().Get("studio_id"); val != "" {
		parsed, err := strconv.Atoi(val)
		if err == nil && parsed > 0 {
			studioID = parsed
		}
	}

	movies, err := api.db.Movies(r.Context(), studioID)
	if err != nil {
		api.respondError(w, http.StatusInternalServerError, "failed to get movies: "+err.Error())
		return
	}

	if movies == nil {
		movies = []db.Movie{}
	}

	api.respondJSON(w, http.StatusOK, movies)
}

// handleMovieByID обрабатывает GET /api/v1/movies/{id}.
func (api *API) handleMovieByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		api.respondError(w, http.StatusBadRequest, "invalid movie id")
		return
	}

	movie, err := api.db.MovieByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			api.respondError(w, http.StatusNotFound, "movie not found")
			return
		}
		api.respondError(w, http.StatusInternalServerError, "failed to get movie: "+err.Error())
		return
	}

	api.respondJSON(w, http.StatusOK, movie)
}

// handleCreateMovies обрабатывает POST /api/v1/movies.
// Поддерживает как добавление массива фильмов [{"title": ...}], так и одного фильма {"title": ...}.
func (api *API) handleCreateMovies(w http.ResponseWriter, r *http.Request) {
	var body json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		api.respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	var movies []db.Movie
	// Сначала пробуем распарсить как массив фильмов
	if err := json.Unmarshal(body, &movies); err != nil {
		// Если не массив, пробуем распарсить как один объект фильма
		var single db.Movie
		if errSingle := json.Unmarshal(body, &single); errSingle != nil {
			api.respondError(w, http.StatusBadRequest, "invalid JSON: expected movie object or array of movies")
			return
		}
		movies = []db.Movie{single}
	}

	if len(movies) == 0 {
		api.respondError(w, http.StatusBadRequest, "movies array cannot be empty")
		return
	}

	for _, m := range movies {
		if m.Title == "" {
			api.respondError(w, http.StatusBadRequest, "movie title cannot be empty")
			return
		}
		if m.ReleaseYear < 1800 {
			api.respondError(w, http.StatusBadRequest, "release_year cannot be less than 1800")
			return
		}
	}

	if err := api.db.AddMovies(r.Context(), movies); err != nil {
		api.respondError(w, http.StatusInternalServerError, "failed to add movies: "+err.Error())
		return
	}

	api.respondJSON(w, http.StatusCreated, map[string]interface{}{
		"status": "created",
		"count":  len(movies),
	})
}

// handleUpdateMovie обрабатывает PUT /api/v1/movies/{id}.
func (api *API) handleUpdateMovie(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		api.respondError(w, http.StatusBadRequest, "invalid movie id")
		return
	}

	var m db.Movie
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		api.respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	m.ID = id

	if m.Title == "" {
		api.respondError(w, http.StatusBadRequest, "movie title cannot be empty")
		return
	}
	if m.ReleaseYear < 1800 {
		api.respondError(w, http.StatusBadRequest, "release_year cannot be less than 1800")
		return
	}

	if err := api.db.UpdateMovie(r.Context(), m); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			api.respondError(w, http.StatusNotFound, "movie not found")
			return
		}
		api.respondError(w, http.StatusInternalServerError, "failed to update movie: "+err.Error())
		return
	}

	api.respondJSON(w, http.StatusOK, m)
}

// handleDeleteMovie обрабатывает DELETE /api/v1/movies/{id}.
func (api *API) handleDeleteMovie(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		api.respondError(w, http.StatusBadRequest, "invalid movie id")
		return
	}

	if err := api.db.DeleteMovie(r.Context(), id); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			api.respondError(w, http.StatusNotFound, "movie not found")
			return
		}
		api.respondError(w, http.StatusInternalServerError, "failed to delete movie: "+err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// --- Хелперы ответов ---

func (api *API) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("error encoding response: %v", err)
	}
}

func (api *API) respondError(w http.ResponseWriter, status int, msg string) {
	api.respondJSON(w, status, map[string]string{"error": msg})
}
