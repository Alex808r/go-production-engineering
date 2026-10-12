// Package api реализует HTTP REST API микросервиса Lynks Shortener.
package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"

	"go-core-4/lynks/shortener/pkg/cacheclient"
	"go-core-4/lynks/shortener/pkg/metrics"
	"go-core-4/lynks/shortener/pkg/storage"
	"go-core-4/lynks/shortener/pkg/urls"
)

// ShortenRequest описывает тело запроса сокращения ссылки.
type ShortenRequest struct {
	Destination string `json:"destination"`
}

// ShortenResponse описывает тело ответа создания ссылки.
type ShortenResponse struct {
	ShortURL    string `json:"shortUrl"`
	Destination string `json:"destination"`
}

// API инкапсулирует зависимости микросервиса shortener.
type API struct {
	router      *mux.Router
	storage     storage.Storage
	cacheClient cacheclient.Client
	baseURL     string
	logger      zerolog.Logger
	metrics     *metrics.Metrics
}

// New создает экземпляр API shortener.
func New(storage storage.Storage, cacheClient cacheclient.Client, baseURL string, logger zerolog.Logger) *API {
	baseURL = strings.TrimRight(baseURL, "/")

	a := &API{
		router:      mux.NewRouter(),
		storage:     storage,
		cacheClient: cacheClient,
		baseURL:     baseURL,
		logger:      logger,
		metrics:     metrics.New("lynks_shortener"),
	}

	// Middleware
	a.router.Use(a.observabilityMiddleware)

	// Маршруты по ТЗ (final.md):
	// 1. Создание короткой ссылки: POST / и POST /api/v1/links
	a.router.HandleFunc("/", a.handleShorten).Methods(http.MethodPost)
	a.router.HandleFunc("/api/v1/links", a.handleShorten).Methods(http.MethodPost)

	// 2. Системные эндпоинты
	a.router.HandleFunc("/healthz", a.handleHealthz).Methods(http.MethodGet)
	a.router.Handle("/metrics", promhttp.Handler()).Methods(http.MethodGet)

	// 3. Редирект по короткой ссылке: GET /{short} (регистрируется последним)
	a.router.HandleFunc("/{short:[a-zA-Z0-9_-]+}", a.handleRedirect).Methods(http.MethodGet)

	return a
}

// Router возвращает настроенный HTTP-роутер.
func (a *API) Router() http.Handler {
	return a.router
}

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

		// Метрики Prometheus
		statusStr := strconv.Itoa(wrapper.statusCode)
		a.metrics.RequestsTotal.WithLabelValues(r.Method, r.URL.Path, statusStr).Inc()
		a.metrics.RequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration.Seconds())

		// Zerolog логирование
		a.logger.Info().
			Str("service", "shortener").
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status", wrapper.statusCode).
			Dur("duration", duration).
			Str("remote_addr", r.RemoteAddr).
			Msg("HTTP request processed")
	})
}

// handleShorten: принимает destination, генерирует shortUrl, сохраняет в PostgreSQL и обновляет кэш.
func (a *API) handleShorten(w http.ResponseWriter, r *http.Request) {
	var req ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.logger.Error().Err(err).Msg("failed to decode JSON in handleShorten")
		a.respondError(w, http.StatusBadRequest, "некорректный JSON в теле запроса")
		return
	}

	cleanURL := strings.TrimSpace(req.Destination)
	if cleanURL == "" {
		a.respondError(w, http.StatusBadRequest, "поле destination не может быть пустым")
		return
	}

	parsed, err := url.ParseRequestURI(cleanURL)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		a.respondError(w, http.StatusBadRequest, "некорректный URL (поддерживаются только http и https)")
		return
	}

	// Проверяем, есть ли уже такая ссылка в БД
	existingShort, err := a.storage.FindByDestination(r.Context(), cleanURL)
	if err == nil && existingShort != "" {
		fullShortURL := fmt.Sprintf("%s/%s", a.baseURL, existingShort)
		a.respondJSON(w, http.StatusOK, ShortenResponse{
			ShortURL:    fullShortURL,
			Destination: cleanURL,
		})
		return
	}

	// Генерируем уникальный короткий идентификатор через urls.Shorten (например, 5 символов как qz6d7)
	shortID := urls.Shorten(cleanURL)
	fullShortURL := fmt.Sprintf("%s/%s", a.baseURL, shortID)

	// Сохранение в PostgreSQL
	if err := a.storage.Save(r.Context(), shortID, cleanURL); err != nil {
		a.logger.Error().Err(err).Str("short", shortID).Msg("failed to save in PostgreSQL")
		a.respondError(w, http.StatusInternalServerError, "ошибка сохранения ссылки в базе данных: "+err.Error())
		return
	}

	// Асинхронное прогревание кэша в memcache
	if a.cacheClient != nil {
		go func(sID, dest string) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			if err := a.cacheClient.Set(ctx, sID, dest); err != nil {
				a.logger.Warn().Err(err).Str("short", sID).Msg("не удалось отправить пару в memcache")
			}
		}(shortID, cleanURL)
	}

	a.respondJSON(w, http.StatusCreated, ShortenResponse{
		ShortURL:    fullShortURL,
		Destination: cleanURL,
	})
}

// handleRedirect: перенаправляет пользователя по короткой ссылке (HTTP 302 Found)
// Реализует паттерн Cache-Aside: Сначала кэш Memcache (Redis), затем PostgreSQL при промахе.
func (a *API) handleRedirect(w http.ResponseWriter, r *http.Request) {
	short := mux.Vars(r)["short"]
	short = strings.TrimSpace(short)

	// 1. Проверяем кэш микросервиса Memcache
	if a.cacheClient != nil {
		dest, err := a.cacheClient.Get(r.Context(), short)
		if err == nil && dest != "" {
			a.metrics.CacheHitsTotal.Inc()
			a.logger.Info().Str("short", short).Str("destination", dest).Msg("Cache Hit! Redirecting from cache")
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}
	}

	// 2. Cache Miss: запрашиваем постоянное хранилище PostgreSQL
	a.metrics.CacheMissTotal.Inc()
	a.logger.Info().Str("short", short).Msg("Cache Miss! Querying PostgreSQL")

	dest, err := a.storage.FindByShort(r.Context(), short)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		a.logger.Error().Err(err).Str("short", short).Msg("PostgreSQL query failed")
		a.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// 3. Сохраняем в кэш для последующих быстрых запросов (Cache Warming)
	if a.cacheClient != nil {
		go func(sID, destination string) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			_ = a.cacheClient.Set(ctx, sID, destination)
		}(short, dest)
	}

	// 4. HTTP Redirect на оригинальную страницу
	http.Redirect(w, r, dest, http.StatusFound)
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
