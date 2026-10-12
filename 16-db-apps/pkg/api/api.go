// Package api предоставляет HTTP API для работы с репозиторием книг.
package api

import (
	"encoding/json"
	"net/http"

	"go-core-4/16-db-apps/pkg/db"
)

// API инкапсулирует абстрактное хранилище книг.
type API struct {
	db db.Interface
}

// New создает экземпляр API с переданным интерфейсом БД.
func New(db db.Interface) *API {
	return &API{db: db}
}

// Handler возвращает HTTP-обработчик для маршрута получения книг.
func (api *API) Handler() http.HandlerFunc {
	return api.handler
}

// ServeHTTP реализует интерфейс http.Handler.
func (api *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	api.handler(w, r)
}

func (api *API) handler(w http.ResponseWriter, r *http.Request) {
	books, err := api.db.Books(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if books == nil {
		books = []db.Book{}
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(books)
}
