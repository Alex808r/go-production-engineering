// Package handler демонстрирует создание и тестирование HTTP-обработчиков (HTTP Handlers) в Go.
package handler

import (
	"encoding/json"
	"net/http"
)

// API представляет HTTP-сервис и хранит необходимые зависимости (строка подключения к БД и др.).
type API struct {
	dbConn string
}

// product описывает товар в формате JSON для ответа клиенту.
type product struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
}

// ProductsHandler обрабатывает GET-запросы на получение списка товаров.
// Соответствует стандартной сигнатуре `http.HandlerFunc(w http.ResponseWriter, r *http.Request)`.
func (api API) ProductsHandler(w http.ResponseWriter, r *http.Request) {
	data := []product{
		{
			Name:  "Книга",
			Price: 1000,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

