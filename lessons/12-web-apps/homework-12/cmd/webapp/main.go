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

	"go-core-4/lessons/12-web-apps/homework-12/pkg/webapp"
)

func main() {
	fmt.Println("=== Запуск веб-приложения GoSearch (webapp) ===")

	// 1. Инициализируем хранилище и наполняем демонстрационным корпусом документов
	db := webapp.NewSearchDB()
	sampleDocs := []webapp.Document{
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
			URL:   "https://pkg.go.dev/html/template",
			Title: "Package html/template - Data-driven templates for HTML",
			Body:  "Package html/template implements data-driven templates for generating HTML code safely.",
		},
	}

	db.Add(sampleDocs)
	fmt.Printf("Успешно проиндексировано документов: %d\n", len(db.Documents()))
	fmt.Printf("Количество терминов в поисковом индексе: %d\n", len(db.Index()))

	// 2. Создаем и конфигурируем веб-приложение
	app := webapp.New(db)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      app.Router(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// 3. Graceful Shutdown по системным сигналам
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nПолучен сигнал завершения. Остановка веб-сервера...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Ошибка при остановке сервера: %v", err)
		}
	}()

	fmt.Println("Веб-сервер запущен и доступен:")
	fmt.Println("  🌐 Документы:           http://localhost:8080/docs")
	fmt.Println("  🔍 Поисковый индекс:    http://localhost:8080/index")
	fmt.Println("  📄 Документы (JSON):    http://localhost:8080/docs?format=json")
	fmt.Println("  📊 Индекс (JSON):       http://localhost:8080/index?format=json")

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Ошибка веб-сервера: %v", err)
	}

	fmt.Println("Сервер успешно остановлен.")
}
