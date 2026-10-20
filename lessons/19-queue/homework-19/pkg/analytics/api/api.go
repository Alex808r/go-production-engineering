// Package api реализует HTTP REST API микросервиса аналитики ссылок.
package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gorilla/mux"

	"go-core-4/lessons/19-queue/homework-19/pkg/analytics/service"
)

// API инкапсулирует HTTP-роутер и сервис аналитики.
type API struct {
	router       *mux.Router
	svc          *service.Service
	requestCount uint64
}

// New создает экземпляр API аналитики.
func New(svc *service.Service) *API {
	a := &API{
		router: mux.NewRouter(),
		svc:    svc,
	}

	a.router.Use(a.loggingMiddleware)
	a.router.Use(a.metricsMiddleware)

	// Основной аналитический эндпоинт по ТЗ (Задача №1, п.2):
	a.router.HandleFunc("/api/v1/stats", a.handleGetStats).Methods(http.MethodGet)
	a.router.HandleFunc("/stats", a.handleGetStats).Methods(http.MethodGet) // Алиас для быстрого доступа

	// Системные микросервисные эндпоинты:
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
		log.Printf("[ANALYTICS-HTTP] %s %s %s - %v", r.Method, r.URL.Path, r.RemoteAddr, time.Since(start))
	})
}

func (a *API) metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddUint64(&a.requestCount, 1)
		next.ServeHTTP(w, r)
	})
}

// --- Обработчики ---

// handleGetStats: возвращает агрегированную статистику ссылок в формате JSON.
func (a *API) handleGetStats(w http.ResponseWriter, _ *http.Request) {
	stats := a.svc.GetStats()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(stats)
}

func (a *API) handleHealthz(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (a *API) handleReady(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ready"})
}

func (a *API) handleMetrics(w http.ResponseWriter, _ *http.Request) {
	stats := a.svc.GetStats()
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")

	fmt.Fprintf(w, "# HELP analytics_http_requests_total Total number of HTTP requests to Analytics API\n")
	fmt.Fprintf(w, "# TYPE analytics_http_requests_total counter\n")
	fmt.Fprintf(w, "analytics_http_requests_total %d\n\n", atomic.LoadUint64(&a.requestCount))

	fmt.Fprintf(w, "# HELP analytics_links_total Total number of tracked links\n")
	fmt.Fprintf(w, "# TYPE analytics_links_total gauge\n")
	fmt.Fprintf(w, "analytics_links_total %d\n\n", stats.TotalLinks)

	fmt.Fprintf(w, "# HELP analytics_avg_original_url_length Average original URL length in characters\n")
	fmt.Fprintf(w, "# TYPE analytics_avg_original_url_length gauge\n")
	fmt.Fprintf(w, "analytics_avg_original_url_length %.2f\n\n", stats.AvgOriginalLength)

	fmt.Fprintf(w, "# HELP analytics_avg_short_url_length Average short URL length in characters\n")
	fmt.Fprintf(w, "# TYPE analytics_avg_short_url_length gauge\n")
	fmt.Fprintf(w, "analytics_avg_short_url_length %.2f\n", stats.AvgShortLength)
}
