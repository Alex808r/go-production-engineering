// Package main содержит реализацию REST API хэндлеров.
package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// API предоставляет интерфейс программного взаимодействия (REST API).
type API struct {
	router *mux.Router
}

// Endpoints регистрирует глобальные middleware и конечные точки API.
func (api *API) Endpoints() {
	// Подключение цепочки middleware
	api.router.Use(logMiddleware)
	api.router.Use(headersMiddleware)

	// Регистрация маршрутов с привязкой к HTTP-методам
	api.router.HandleFunc("/api/v1/books", api.books).Methods(http.MethodGet)
	api.router.HandleFunc("/api/v1/books", api.newBook).Methods(http.MethodPost)
	api.router.HandleFunc("/api/v1/books/{id:[0-9]+}", api.deleteBook).Methods(http.MethodDelete)
}

// books обрабатывает GET-запрос на получение списка всех книг: GET /api/v1/books.
func (api *API) books(w http.ResponseWriter, r *http.Request) {
	booksMutex.RLock()
	defer booksMutex.RUnlock()

	// json.NewEncoder(w).Encode(books) производит потоковую сериализацию напрямую
	// в http.ResponseWriter без создания промежуточного []byte буфера в памяти.
	err := json.NewEncoder(w).Encode(books)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// newBook обрабатывает POST-запрос на создание новой книги: POST /api/v1/books.
func (api *API) newBook(w http.ResponseWriter, r *http.Request) {
	var b book

	// json.NewDecoder(r.Body).Decode(&b) десериализует JSON прямо из входящего потока запроса.
	// Если JSON невалиден, возвращаем клиенту 400 Bad Request.
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		http.Error(w, "некорректный JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Валидация обязательных полей модели
	if b.Name == "" || b.Author == "" {
		http.Error(w, "поля 'name' и 'author' обязательны для заполнения", http.StatusBadRequest)
		return
	}

	booksMutex.Lock()
	books = append(books, b)
	booksMutex.Unlock()

	// По стандартам REST при создании сущности возвращается статус 201 Created
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(b)
}

// deleteBook обрабатывает DELETE-запрос на удаление книги по индексу: DELETE /api/v1/books/{id}.
func (api *API) deleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "некорректный ID книги", http.StatusBadRequest)
		return
	}

	booksMutex.Lock()
	defer booksMutex.Unlock()

	// Проверяем, существует ли книга с таким индексом
	if id < 0 || id >= len(books) {
		http.Error(w, "книга с таким ID не найдена", http.StatusNotFound)
		return
	}

	// Удаление элемента из среза с сохранением порядка
	books = append(books[:id], books[id+1:]...)

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"deleted"}`))
}

