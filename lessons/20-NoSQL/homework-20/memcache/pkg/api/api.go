// Package api реализует HTTP REST API микросервиса кэширования memcache.
package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"

	"go-core-4/lessons/20-NoSQL/homework-20/memcache/pkg/metrics"
	"go-core-4/lessons/20-NoSQL/homework-20/memcache/pkg/store"
)

// CachePair описывает JSON-тело запроса/ответа кэширования пары ссылок.
type CachePair struct {
	ShortURL    string `json:"shortUrl"`
	Destination string `json:"destination"`
}

// API инкапсулирует маршрутизатор, хранилище и логирование сервиса кэша.
type API struct {
	router  *mux.Router
	store   store.Store
	logger  zerolog.Logger
	metrics *metrics.Metrics
}

// New создает экземпляр API кэширования.
func New(store store.Store, logger zerolog.Logger) *API {
	a := &API{
		router:  mux.NewRouter(),
		store:   store,
		logger:  logger,
		metrics: metrics.New("lynks_memcache"),
	}

	// Middleware логирования Zerolog и сбора Prometheus метрик
	a.router.Use(a.observabilityMiddleware)

	// REST API маршруты (коллекция-ресурс):
	// 1. Сохранение пары гиперссылок в кэш: POST /api/v1/cache
	a.router.HandleFunc("/api/v1/cache", a.handleSet).Methods(http.MethodPost)
	// 2. Получение длинной ссылки по короткой из кэша: GET /api/v1/cache/{short}
	a.router.HandleFunc("/api/v1/cache/{short}", a.handleGet).Methods(http.MethodGet)

	// Системные эндпоинты
	a.router.HandleFunc("/healthz", a.handleHealthz).Methods(http.MethodGet)
	a.router.Handle("/metrics", promhttp.Handler()).Methods(http.MethodGet)

	return a
}

// Router возвращает маршрутизатор HTTP.
func (a *API) Router() http.Handler {
	return a.router
}

// responseWriterWrapper перехватывает HTTP статус ответа для метрик и логов.
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriterWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (a *API) observabilityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapper := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapper, r)

		duration := time.Since(start)

		// Фиксация метрик Prometheus
		statusStr := strconv.Itoa(wrapper.statusCode)
		a.metrics.RequestsTotal.WithLabelValues(r.Method, r.URL.Path, statusStr).Inc()
		a.metrics.RequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration.Seconds())

		// Структурированное логирование Zerolog (в stdout по ТЗ)
		a.logger.Info().
			Str("service", "memcache").
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status", wrapper.statusCode).
			Dur("duration", duration).
			Str("remote_addr", r.RemoteAddr).
			Msg("HTTP request processed")
	})
}

// handleSet: сохраняет пару ссылок в кэш.
func (a *API) handleSet(w http.ResponseWriter, r *http.Request) {
	var pair CachePair
	if err := json.NewDecoder(r.Body).Decode(&pair); err != nil {
		a.logger.Error().Err(err).Msg("failed to decode JSON in handleSet")
		a.respondError(w, http.StatusBadRequest, "некорректный JSON в теле запроса")
		return
	}

	pair.ShortURL = strings.TrimSpace(pair.ShortURL)
	pair.Destination = strings.TrimSpace(pair.Destination)
	if pair.ShortURL == "" || pair.Destination == "" {
		a.respondError(w, http.StatusBadRequest, "поля shortUrl и destination обязательны")
		return
	}

	if err := a.store.Set(r.Context(), pair.ShortURL, pair.Destination); err != nil {
		a.logger.Error().Err(err).Str("short", pair.ShortURL).Msg("failed to write into cache store")
		a.respondError(w, http.StatusInternalServerError, "ошибка сохранения в кэш: "+err.Error())
		return
	}

	a.respondJSON(w, http.StatusCreated, pair)
}

// handleGet: возвращает пару по ключу короткой ссылки.
func (a *API) handleGet(w http.ResponseWriter, r *http.Request) {
	short := mux.Vars(r)["short"]
	short = strings.TrimSpace(short)

	dest, err := a.store.Get(r.Context(), short)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			a.respondError(w, http.StatusNotFound, "ссылка не найдена в кэше")
			return
		}
		a.logger.Error().Err(err).Str("short", short).Msg("failed to get from cache store")
		a.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	a.respondJSON(w, http.StatusOK, CachePair{
		ShortURL:    short,
		Destination: dest,
	})
}

func (a *API) handleHealthz(w http.ResponseWriter, _ *http.Request) {
	a.respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (a *API) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func (a *API) respondError(w http.ResponseWriter, status int, msg string) {
	a.respondJSON(w, status, map[string]string{
		"error": strings.TrimSpace(msg),
	})
}
