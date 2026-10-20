package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-core-4/lessons/13-api/homework-13/pkg/api"
	"go-core-4/lessons/13-api/homework-13/pkg/storage"
)

func main() {
	fmt.Println("=== Запуск REST API сервера GoSearch (Занятие №13) ===")

	// 1. Инициализация хранилища данных
	db := storage.New()

	// Наполнение демонстрационным корпусом документов
	sampleDocs := []storage.Document{
		{
			ID:    1,
			URL:   "https://go.dev/doc/tutorial/getting-started",
			Title: "Tutorial: Get started with Go",
			Body:  "A tutorial that introduces the basics of the Go language.",
		},
		{
			ID:    2,
			URL:   "https://go.dev/doc/effective_go",
			Title: "Effective Go - Best Practices",
			Body:  "Tips for writing clear, performant, and idiomatic Go code.",
		},
		{
			ID:    3,
			URL:   "https://pkg.go.dev/net/http",
			Title: "Package net/http - HTTP client and server implementations",
			Body:  "Package http provides HTTP client and server implementations in Go.",
		},
		{
			ID:    4,
			URL:   "https://pkg.go.dev/github.com/gorilla/mux",
			Title: "Gorilla Mux - Powerful URL router and dispatcher",
			Body:  "A powerful HTTP router and URL matcher for building Go web servers.",
		},
		{
			ID:    5,
			URL:   "https://swagger.io/specification/",
			Title: "OpenAPI Specification (Swagger)",
			Body:  "The OpenAPI Specification is a broadly adopted industry standard for REST APIs.",
		},
	}

	for _, d := range sampleDocs {
		db.Add(d)
	}

	fmt.Printf("Загружено документов в индекс: %d\n", len(db.Documents()))

	// 2. Инициализация API-слоя
	apiService := api.New(db)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      apiService.Router(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// 3. Graceful Shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nПолучен сигнал завершения. Остановка REST API сервера...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Ошибка при остановке сервера: %v", err)
		}
	}()

	fmt.Println("REST API сервер успешно запущен на порту :8080")
	fmt.Println("Доступные эндпоинты:")
	fmt.Println("  🔍 Поиск документов:      GET    http://localhost:8080/api/v1/search?q=concurrency")
	fmt.Println("  📄 Список документов:     GET    http://localhost:8080/api/v1/docs")
	fmt.Println("  ➕ Создание документа:    POST   http://localhost:8080/api/v1/docs")
	fmt.Println("  🔎 Документ по ID:        GET    http://localhost:8080/api/v1/docs/1")
	fmt.Println("  ✏️ Обновление по ID:      PUT    http://localhost:8080/api/v1/docs/1")
	fmt.Println("  🗑️ Удаление по ID:        DELETE http://localhost:8080/api/v1/docs/1")

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Ошибка сервера: %v", err)
	}

	fmt.Println("Сервер успешно остановлен.")
}
