// Package api реализует HTTP REST API микросервиса хранения коротких ссылок.
// Соответствует Задаче №1 домашнего задания №18.
package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gorilla/mux"

	"go-core-4/lessons/18-microservices/homework-18/pkg/storage"
)

// API инкапсулирует маршрутизатор, хранилище и метрики микросервиса.
type API struct {
	router       *mux.Router
	store        storage.Interface
	requestCount uint64
}

// New создает новый экземпляр API с зарегистрированными маршрутами и middleware.
func New(store storage.Interface) *API {
	a := &API{
		router: mux.NewRouter(),
		store:  store,
	}

	// Middleware
	a.router.Use(a.loggingMiddleware)
	a.router.Use(a.metricsMiddleware)

	// Основные методы API по ТЗ Задачи №1:
	// 1. Добавление ссылки: POST /api/v1/links
	a.router.HandleFunc("/api/v1/links", a.handleCreateLink).Methods(http.MethodPost)
	// 2. Получение оригинальной ссылки по ID: GET /api/v1/links/{id}
	a.router.HandleFunc("/api/v1/links/{id}", a.handleGetLink).Methods(http.MethodGet)
	// 3. Получение оригинальной ссылки по полной короткой ссылке в теле: POST /api/v1/links/resolve
	a.router.HandleFunc("/api/v1/links/resolve", a.handleResolveLink).Methods(http.MethodPost)
	// 4. Редирект по короткой ссылке: GET /r/{id}
	a.router.HandleFunc("/r/{id}", a.handleRedirect).Methods(http.MethodGet)
	// 5. Список всех ссылок: GET /api/v1/links
	a.router.HandleFunc("/api/v1/links", a.handleListLinks).Methods(http.MethodGet)

	// Системные микросервисные эндпоинты (Kubernetes / Prometheus / Cloud-Ready):
	a.router.HandleFunc("/healthz", a.handleHealthz).Methods(http.MethodGet)
	a.router.HandleFunc("/ready", a.handleReady).Methods(http.MethodGet)
	a.router.HandleFunc("/metrics", a.handleMetrics).Methods(http.MethodGet)

	return a
}

// Router возвращает настроенный маршрутизатор.
func (a *API) Router() http.Handler {
	return a.router
}

// --- Middleware ---

func (a *API) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("[%s] %s %s - %v", r.Method, r.URL.Path, r.RemoteAddr, time.Since(start))
	})
}

func (a *API) metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddUint64(&a.requestCount, 1)
		next.ServeHTTP(w, r)
	})
}

// --- Обработчики методов ТЗ ---

// CreateLinkRequest — тело запроса создания короткой ссылки.
type CreateLinkRequest struct {
	URL string `json:"url"`
}

// ResolveLinkRequest — тело запроса поиска по короткой ссылке.
type ResolveLinkRequest struct {
	ShortURL string `json:"short_url"`
}

// 1. handleCreateLink: принимает ссылку, добавляет её в хранилище и возвращает короткую ссылку.
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

	a.respondJSON(w, http.StatusCreated, link)
}

// 2. handleGetLink: принимает ID ссылки в URL и возвращает оригинальную ссылку.
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

// handleResolveLink: принимает полную короткую ссылку в теле POST-запроса и возвращает оригинальную.
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

// handleRedirect: перенаправляет пользователя на оригинальный URL (HTTP 302 Found).
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

// --- Вспомогательные функции ответов ---

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
