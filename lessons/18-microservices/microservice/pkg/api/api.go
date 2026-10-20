// Package api реализует HTTP API демонстрационного микросервиса для Kubernetes.
// Включает возврат hostname пода для демонстрации балансировки нагрузки в k8s.
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/mux"
)

// API представляет HTTP-интерфейс микросервиса.
type API struct {
	Router *mux.Router

	sync.Mutex
	db []Review
}

// New инициализирует API и регистрирует маршруты.
func New() *API {
	api := API{
		Router: mux.NewRouter(),
	}

	api.Router.HandleFunc("/", api.handler).Methods(http.MethodGet)
	api.Router.HandleFunc("/healthz", api.healthz).Methods(http.MethodGet)
	api.Router.HandleFunc("/api/v1/reviews", api.reviews).Methods(http.MethodGet)

	return &api
}

// Review представляет модель отзыва к фильму.
type Review struct {
	MovieID int    `json:"movie_id"`
	Text    string `json:"text"`
}

// handler возвращает имя хоста (в Kubernetes — имя пода), на котором запущен контейнер.
// Это позволяет наглядно увидеть работу Service Load Balancer при нескольких репликах.
func (api *API) handler(w http.ResponseWriter, r *http.Request) {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(fmt.Sprintf("Microapp на машине %v\n", hostname)))
}

// healthz — проба жизнеспособности (Liveness Probe) для Kubernetes.
func (api *API) healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

// reviews возвращает список сохраненных отзывов.
func (api *API) reviews(w http.ResponseWriter, r *http.Request) {
	api.Lock()
	defer api.Unlock()

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(api.db)
}
