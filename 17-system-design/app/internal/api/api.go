package api

import (
	"encoding/json"
	"go-core-4/17-system-design/app/internal/db"
	"go-core-4/17-system-design/app/internal/models"
	"net/http"

	"github.com/gorilla/mux"
)

// API предоставляет интерфейс программного взаимодействия.
type API struct {
	router *mux.Router
	db     *db.DB
}

func New(db *db.DB) *API {
	api := API{
		router: mux.NewRouter(),
		db:     db,
	}
	api.Endpoints()
	return &api
}

// Router возвращает настроенный HTTP-маршрутизатор.
func (api *API) Router() *mux.Router {
	return api.router
}

func (api *API) Run(addr string) error {
	return http.ListenAndServe(addr, api.router)
}

// Endpoints регистрирует конечные точки API.
func (api *API) Endpoints() {
	api.router.HandleFunc("/api/v1/books", api.books).Methods(http.MethodGet)
	api.router.HandleFunc("/api/v1/books", api.newBook).Methods(http.MethodPost)
	api.router.HandleFunc("/api/v1/books/{id}", api.deleteBook).Methods(http.MethodDelete)
}

func (api *API) books(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err := json.NewEncoder(w).Encode(models.Books)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (api *API) newBook(w http.ResponseWriter, r *http.Request) {
	var b models.Book
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	models.Books = append(models.Books, b)
	w.WriteHeader(http.StatusCreated)
}

func (api *API) deleteBook(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["id"]
	for i, b := range models.Books {
		if b.Name == name {
			models.Books = append(models.Books[:i], models.Books[i+1:]...)
			w.WriteHeader(http.StatusOK)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}
