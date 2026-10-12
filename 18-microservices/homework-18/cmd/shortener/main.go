// Микросервис сокращения ссылок (Homework-18).
// Архитектура типичного микросервиса (слайд 9 лекции №18):
// - Конфигурация через переменные окружения и флаги.
// - In-memory хранилище на базе среза (по условию ТЗ).
// - REST API + Prometheus метрики + Healthcheck.
// - Graceful Shutdown для безопасного завершения в Kubernetes / Docker.
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

	"go-core-4/18-microservices/homework-18/pkg/api"
	"go-core-4/18-microservices/homework-18/pkg/config"
	"go-core-4/18-microservices/homework-18/pkg/storage"
)

func main() {
	// 1. Загрузка конфигурации (переменные окружения)
	cfg := config.Load()

	// 2. Инициализация хранилища (срез в памяти по требованию ТЗ)
	store := storage.NewSliceStore(cfg.BaseURL)

	// Предзаполнение демонстрационными ссылками
	seed(store)

	// 3. Создание API
	appAPI := api.New(store)

	// 4. Настройка HTTP-сервера с промышленными таймаутами
	server := &http.Server{
		Addr:         cfg.Port,
		Handler:      appAPI.Router(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// 5. Обработка сигналов завершения (Graceful Shutdown)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		fmt.Printf("🚀 URL Shortener Microservice запущен на порту %s (BaseURL: %s)\n", cfg.Port, cfg.BaseURL)
		fmt.Println("   POST   /api/v1/links         — добавить ссылку и получить короткую")
		fmt.Println("   GET    /api/v1/links/{id}    — получить оригинальную ссылку по ID")
		fmt.Println("   POST   /api/v1/links/resolve — получить оригинал по короткой ссылке")
		fmt.Println("   GET    /r/{id}               — редирект (HTTP 302 Found)")
		fmt.Println("   GET    /api/v1/links         — список всех ссылок в памяти")
		fmt.Println("   GET    /healthz              — проверка здоровья (Liveness Probe)")
		fmt.Println("   GET    /metrics              — метрики Prometheus")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ошибка HTTP-сервера: %v", err)
		}
	}()

	<-ctx.Done()
	fmt.Println("\n🛑 Получен сигнал завершения. Остановка микросервиса...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("ошибка принудительной остановки: %v", err)
	}

	fmt.Println("✅ Микросервис остановлен корректно.")
}

func seed(store storage.Interface) {
	ctx := context.Background()
	_, _ = store.Save(ctx, "https://go.dev")
	_, _ = store.Save(ctx, "https://golang.org/doc/")
	_, _ = store.Save(ctx, "https://thinknetica.com")
}
