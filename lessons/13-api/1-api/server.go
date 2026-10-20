// Package main демонстрирует основы проектирования базового REST API на Go
// с использованием маршрутизатора Gorilla Mux (пакет github.com/gorilla/mux).
package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

// server инкапсулирует зависимости веб-сервера: маршрутизатор и API-контроллер.
type server struct {
	api    *API
	router *mux.Router
}

func main() {
	// Создаем экземпляр сервера и инициализируем маршрутизатор
	srv := new(server)
	srv.router = mux.NewRouter()

	// Внедряем маршрутизатор в структуру API (Dependency Injection)
	srv.api = &API{router: srv.router}

	// Регистрируем middleware и маршруты (Endpoints)
	srv.api.Endpoints()

	fmt.Println("REST API сервер (1-api) запущен на http://localhost:8081")
	fmt.Println("Доступные эндпоинты:")
	fmt.Println("  GET    http://localhost:8081/api/v1/books      - получение списка книг")
	fmt.Println("  POST   http://localhost:8081/api/v1/books      - добавление книги")
	fmt.Println("  DELETE http://localhost:8081/api/v1/books/{id} - удаление книги по индексу")

	// Запуск HTTP-сервера
	log.Fatal(http.ListenAndServe(":8081", srv.router))
}

// book представляет модель сущности книги в REST API.
type book struct {
	Name   string `json:"name"`
	Author string `json:"author"`
}

// Мьютекс для защиты среза books от состояний гонки (Data Race)
// при конкурентных HTTP-запросах.
var booksMutex sync.RWMutex

// books — коллекция книг в оперативной памяти (демонстрационная база данных).
var books = []book{
	{
		Name:   "The Lord Of The Rings",
		Author: "J.R.R. Tolkien",
	},
}

