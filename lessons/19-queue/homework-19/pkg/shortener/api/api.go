// Package api реализует REST API микросервиса сокращения гиперссылок
// с публикацией событий создания ссылок в брокер сообщений (Kafka / Очередь).
package api

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gorilla/mux"

	"go-core-4/lessons/19-queue/homework-19/pkg/event"
	"go-core-4/lessons/19-queue/homework-19/pkg/shortener/storage"
)

// API инкапсулирует маршрутизатор, хранилище ссылок и продюсер очереди событий.
type API struct {
	router       *mux.Router
	store        storage.Interface
	producer     event.Producer
	requestCount uint64
}

// New создает экземпляр API с зарегистрированными маршрутами, хранилищем и брокером событий.
func New(store storage.Interface, producer event.Producer) *API {
	a := &API{
		router:   mux.NewRouter(),
		store:    store,
		producer: producer,
	}

	a.router.Use(a.loggingMiddleware)
	a.router.Use(a.metricsMiddleware)

	// Основные маршруты по ТЗ:
	// 1. Добавление ссылки: POST /api/v1/links
	a.router.HandleFunc("/api/v1/links", a.handleCreateLink).Methods(http.MethodPost)
	// 2. Получение ссылки по ID: GET /api/v1/links/{id}
	a.router.HandleFunc("/api/v1/links/{id}", a.handleGetLink).Methods(http.MethodGet)
	// 3. Получение оригинала по короткой ссылке в JSON: POST /api/v1/links/resolve
	a.router.HandleFunc("/api/v1/links/resolve", a.handleResolveLink).Methods(http.MethodPost)
	// 4. Редирект на оригинальный URL: GET /r/{id}
	a.router.HandleFunc("/r/{id}", a.handleRedirect).Methods(http.MethodGet)
	// 5. Просмотр всех ссылок в памяти: GET /api/v1/links
	a.router.HandleFunc("/api/v1/links", a.handleListLinks).Methods(http.MethodGet)

	// Системные эндпоинты (Liveness, Readiness, Prometheus metrics):
	a.router.HandleFunc("/healthz", a.handleHealthz).Methods(http.MethodGet)
	a.router.HandleFunc("/ready", a.handleReady).Methods(http.MethodGet)
	a.router.HandleFunc("/metrics", a.handleMetrics).Methods(http.MethodGet)

	return a
}

// Router возвращает настроенный HTTP-маршрутизатор.
func (a *API) Router() http.Handler {
	return a.router
}

// --- Middleware ---

func (a *API) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("[SHORTENER-HTTP] %s %s %s - %v", r.Method, r.URL.Path, r.RemoteAddr, time.Since(start))
	})
}

func (a *API) metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddUint64(&a.requestCount, 1)
		next.ServeHTTP(w, r)
	})
}

// --- Обработчики методов ---

// CreateLinkRequest представляет тело запроса создания короткой ссылки.
type CreateLinkRequest struct {
	URL string `json:"url"`
}

// ResolveLinkRequest представляет тело запроса разрешения короткой ссылки.
type ResolveLinkRequest struct {
	ShortURL string `json:"short_url"`
}

// handleCreateLink: сохраняет ссылку в хранилище и публикует событие в брокер сообщений.
func (a *API) handleCreateLink(w http.ResponseWriter, r *http.Request) {
	var req CreateLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.respondError(w, http.StatusBadRequest, "некорректный JSON в теле запроса")
		return
	}

	link, err := a.store.Save(r.Context(), req.URL)
	if err != nil {
		if errors.Is(err, storage.ErrEmptyURL) || errors.Is(err, storage.ErrInvalidURL) {
			a.respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		a.respondError(w, http.StatusInternalServerError, "ошибка сохранения ссылки: "+err.Error())
		return
	}

	// Публикация события в очередь сообщений (Kafka / Channel)
	if a.producer != nil {
		ev := event.LinkCreatedEvent{
			EventID:   generateEventID(),
			Type:      "link.created",
			LinkID:    link.ID,
			Original:  link.Original,
			ShortURL:  link.ShortURL,
			CreatedAt: link.CreatedAt,
		}

		// Асинхронно или с таймаутом отправляем событие, не блокируя клиентский HTTP-ответ надолго
		go func(eventPayload event.LinkCreatedEvent) {
			pubCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if pubErr := a.producer.Publish(pubCtx, eventPayload.LinkID, eventPayload); pubErr != nil {
				log.Printf("[PRODUCER-ERROR] Не удалось отправить событие создания ссылки ID=%s: %v", eventPayload.LinkID, pubErr)
			} else {
				log.Printf("[PRODUCER] Событие создания ссылки успешно опубликовано: link_id=%s, url=%s", eventPayload.LinkID, eventPayload.Original)
			}
		}(ev)
	}

	a.respondJSON(w, http.StatusCreated, link)
}

// handleGetLink: возвращает ссылку по её ID.
func (a *API) handleGetLink(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	link, err := a.store.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			a.respondError(w, http.StatusNotFound, "ссылка не найдена")
			return
		}
		a.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	a.respondJSON(w, http.StatusOK, link)
}

// handleResolveLink: возвращает оригинальную ссылку по значению short_url.
func (a *API) handleResolveLink(w http.ResponseWriter, r *http.Request) {
	var req ResolveLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.respondError(w, http.StatusBadRequest, "некорректный JSON")
		return
	}

	link, err := a.store.Get(r.Context(), req.ShortURL)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			a.respondError(w, http.StatusNotFound, "короткая ссылка не найдена")
			return
		}
		a.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	a.respondJSON(w, http.StatusOK, link)
}

// handleRedirect: выполняет редирект (HTTP 302 Found).
func (a *API) handleRedirect(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	link, err := a.store.Get(r.Context(), id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, link.Original, http.StatusFound)
}

// handleListLinks: возвращает все ссылки из памяти.
func (a *API) handleListLinks(w http.ResponseWriter, r *http.Request) {
	links, err := a.store.All(r.Context())
	if err != nil {
		a.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	a.respondJSON(w, http.StatusOK, links)
}

// --- Микросервисные эндпоинты ---

func (a *API) handleHealthz(w http.ResponseWriter, _ *http.Request) {
	a.respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (a *API) handleReady(w http.ResponseWriter, _ *http.Request) {
	a.respondJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}

func (a *API) handleMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	all, _ := a.store.All(r.Context())
	fmt.Fprintf(w, "# HELP http_requests_total Total number of HTTP requests\n")
	fmt.Fprintf(w, "# TYPE http_requests_total counter\n")
	fmt.Fprintf(w, "http_requests_total %d\n\n", atomic.LoadUint64(&a.requestCount))
	fmt.Fprintf(w, "# HELP shortener_links_stored_total Number of links stored in memory\n")
	fmt.Fprintf(w, "# TYPE shortener_links_stored_total gauge\n")
	fmt.Fprintf(w, "shortener_links_stored_total %d\n", len(all))
}

// --- Вспомогательные методы ---

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

func generateEventID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
