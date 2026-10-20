package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"go-core-4/lessons/13-api/homework-13/pkg/storage"
)

// API представляет REST API поисковика GoSearch.
type API struct {
	router *mux.Router
	db     storage.Interface
}

// New инициализирует API, регистрирует middleware и маршруты.
func New(db storage.Interface) *API {
	api := &API{
		router: mux.NewRouter(),
		db:     db,
	}
	api.endpoints()
	return api
}

// Router возвращает настроенный HTTP-маршрутизатор.
func (api *API) Router() *mux.Router {
	return api.router
}

// endpoints конфигурирует цепочки middleware и обработчики маршрутов.
func (api *API) endpoints() {
	// Подключение middleware
	api.router.Use(headersMiddleware)
	api.router.Use(requestIDMiddleware)
	api.router.Use(loggingMiddleware)

	// Роуты API v1
	api.router.HandleFunc("/api/v1/search", api.handleSearch).Methods(http.MethodGet)
	api.router.HandleFunc("/api/v1/docs", api.handleDocs).Methods(http.MethodGet)
	api.router.HandleFunc("/api/v1/docs", api.handleCreateDoc).Methods(http.MethodPost)
	api.router.HandleFunc("/api/v1/docs/{id:[0-9]+}", api.handleDocByID).Methods(http.MethodGet)
	api.router.HandleFunc("/api/v1/docs/{id:[0-9]+}", api.handleUpdateDoc).Methods(http.MethodPut)
	api.router.HandleFunc("/api/v1/docs/{id:[0-9]+}", api.handleDeleteDoc).Methods(http.MethodDelete)
}

// handleSearch обрабатывает поиск документов по текстовому запросу: GET /api/v1/search?q=...
func (api *API) handleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		query = r.URL.Query().Get("query")
	}

	if query == "" {
		http.Error(w, `{"error":"query parameter 'q' is required"}`, http.StatusBadRequest)
		return
	}

	results := api.db.Search(query)
	if results == nil {
		results = []storage.Document{}
	}

	api.respondJSON(w, http.StatusOK, results)
}

// handleDocs возвращает полный список документов: GET /api/v1/docs
func (api *API) handleDocs(w http.ResponseWriter, r *http.Request) {
	docs := api.db.Documents()
	if docs == nil {
		docs = []storage.Document{}
	}
	api.respondJSON(w, http.StatusOK, docs)
}

// handleDocByID возвращает документ по его идентификатору: GET /api/v1/docs/{id}
func (api *API) handleDocByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, `{"error":"invalid document id"}`, http.StatusBadRequest)
		return
	}

	doc, err := api.db.Document(id)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			http.Error(w, `{"error":"document not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	api.respondJSON(w, http.StatusOK, doc)
}

// handleCreateDoc создает новый документ: POST /api/v1/docs
func (api *API) handleCreateDoc(w http.ResponseWriter, r *http.Request) {
	var doc storage.Document
	if err := json.NewDecoder(r.Body).Decode(&doc); err != nil {
		http.Error(w, `{"error":"invalid json body"}`, http.StatusBadRequest)
		return
	}

	if doc.Title == "" || doc.URL == "" {
		http.Error(w, `{"error":"fields 'title' and 'url' cannot be empty"}`, http.StatusBadRequest)
		return
	}

	created := api.db.Add(doc)
	api.respondJSON(w, http.StatusCreated, created)
}

// handleUpdateDoc обновляет документ по ID: PUT /api/v1/docs/{id}
func (api *API) handleUpdateDoc(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, `{"error":"invalid document id"}`, http.StatusBadRequest)
		return
	}

	var doc storage.Document
	if err := json.NewDecoder(r.Body).Decode(&doc); err != nil {
		http.Error(w, `{"error":"invalid json body"}`, http.StatusBadRequest)
		return
	}

	doc.ID = id
	if err := api.db.Update(doc); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			http.Error(w, `{"error":"document not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	api.respondJSON(w, http.StatusOK, doc)
}

// handleDeleteDoc удаляет документ по ID: DELETE /api/v1/docs/{id}
func (api *API) handleDeleteDoc(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, `{"error":"invalid document id"}`, http.StatusBadRequest)
		return
	}

	if err := api.db.Delete(id); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			http.Error(w, `{"error":"document not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	api.respondJSON(w, http.StatusOK, map[string]string{"message": "document deleted successfully"})
}

// respondJSON сериализует payload в JSON и отправляет с указанным HTTP-статусом.
func (api *API) respondJSON(w http.ResponseWriter, status int, data any) {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, `{"error":"failed to serialize response"}`, http.StatusInternalServerError)
	}
}
