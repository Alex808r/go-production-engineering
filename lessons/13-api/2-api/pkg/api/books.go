// Package api содержит обработчики для работы с ресурсом книг.
package api

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

// booksMutex защищает срез books от конкурентных чтений и записей.
var booksMutex sync.RWMutex

// books обрабатывает запрос GET /api/v1/books и возвращает список книг в формате JSON.
func (api *API) books(w http.ResponseWriter, r *http.Request) {
	booksMutex.RLock()
	defer booksMutex.RUnlock()

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err := json.NewEncoder(w).Encode(books)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// newBook обрабатывает запрос POST /api/v1/newBook для добавления книги.
// Демонстрирует использование контекста (context.Context):
// 1. context.WithTimeout — ограничение времени выполнения операции.
// 2. context.WithValue — проброс метаданных (например, requestID) сквозь слои приложения.
func (api *API) newBook(w http.ResponseWriter, r *http.Request) {
	// Создаем контекст с таймаутом на основе контекста входящего HTTP-запроса r.Context()
	timeoutCtx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	// Обязательно вызываем cancel() через defer для своевременного освобождения таймера
	defer cancel()

	// Производный контекст со значением (Request ID)
	ctx := context.WithValue(timeoutCtx, "requestID", rand.Intn(1_000_000_000))
	_ = ctx

	var b book
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	booksMutex.Lock()
	books = append(books, b)
	booksMutex.Unlock()

	w.WriteHeader(http.StatusOK)
}

// book описывает модель книги.
type book struct {
	Name   string `json:"name"`
	Author string `json:"author"`
}

var books = []book{
	{
		Name:   "The Lord Of The Rings",
		Author: "J.R.R. Tolkien",
	},
}

