// Package enhancedrouting демонстрирует возможности улучшенного роутера net/http.ServeMux
// из Go 1.22+ со встроенной поддержкой методов, параметров путей {id} и wildcards {path...}.
package enhancedrouting

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

// Item описывает сущность для демонстрации REST API.
type Item struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Service представляет хранилище элементов и обработчики HTTP.
type Service struct {
	mu    sync.RWMutex
	items map[string]Item
}

// NewService инициализирует сервис с тестовыми данными.
func NewService() *Service {
	return &Service{
		items: map[string]Item{
			"1": {ID: "1", Name: "Laptop"},
			"2": {ID: "2", Name: "Smartphone"},
		},
	}
}

// NewRouter настраивает и возвращает стандартный http.ServeMux с паттернами Go 1.22+.
func (s *Service) NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	// 1. Метод + точный путь
	mux.HandleFunc("GET /items", s.handleListItems)
	mux.HandleFunc("POST /items", s.handleCreateItem)

	// 2. Метод + именованный параметр пути {id}
	mux.HandleFunc("GET /items/{id}", s.handleGetItem)
	mux.HandleFunc("DELETE /items/{id}", s.handleDeleteItem)

	// 3. Catch-all / Wildcard {path...}
	mux.HandleFunc("GET /files/{path...}", s.handleServeFile)

	return mux
}

func (s *Service) handleListItems(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var list []Item
	for _, it := range s.items {
		list = append(list, it)
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(list)
}

func (s *Service) handleGetItem(w http.ResponseWriter, r *http.Request) {
	// В Go 1.22+ параметры извлекаются напрямую через r.PathValue(paramName)!
	id := r.PathValue("id")

	s.mu.RLock()
	item, ok := s.items[id]
	s.mu.RUnlock()

	if !ok {
		http.Error(w, "item not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(item)
}

func (s *Service) handleCreateItem(w http.ResponseWriter, r *http.Request) {
	var item Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil || item.ID == "" {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	s.items[item.ID] = item
	s.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(item)
}

func (s *Service) handleDeleteItem(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	s.mu.Lock()
	delete(s.items, id)
	s.mu.Unlock()

	w.WriteHeader(http.StatusNoContent)
}

func (s *Service) handleServeFile(w http.ResponseWriter, r *http.Request) {
	// Wildcard параметр захватывает весь оставшийся путь (например, "docs/api/v1.pdf")
	filePath := r.PathValue("path")
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, "Serving file: %s", filePath)
}
