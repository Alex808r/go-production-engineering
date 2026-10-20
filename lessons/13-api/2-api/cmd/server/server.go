// Package main запускает production-grade REST API сервер
// с поддержкой сессий, JWT-аутентификации и контекстов запросов.
package main

import (
	"fmt"
	"log"
	"net/http"

	"go-core-4/lessons/13-api/2-api/pkg/api"
)

type server struct {
	api *api.API
}

func main() {
	srv := new(server)

	// Инициализируем сервис API со всеми маршрутами и middleware
	srv.api = api.New()

	fmt.Println("REST API сервер (2-api) запущен на порту :8082")
	fmt.Println("Доступные эндпоинты:")
	fmt.Println("  POST http://localhost:8082/api/v1/authJWT     - получение JWT токена")
	fmt.Println("  POST http://localhost:8082/api/v1/authSession - аутентификация через Cookie сессию")
	fmt.Println("  GET  http://localhost:8082/api/v1/books       - получение книг (требует JWT)")
	fmt.Println("  POST http://localhost:8082/api/v1/newBook     - создание книги (требует JWT)")

	// Запуск HTTP-сервера на порту 8082
	log.Fatal(http.ListenAndServe(":8082", srv.api.Router()))
}

