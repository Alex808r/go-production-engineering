// Микросервис сокращения ссылок (Producer в Kafka).
// Домашнее задание №19: Занятие «Очереди сообщений».
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

	"go-core-4/lessons/19-queue/homework-19/pkg/event"
	"go-core-4/lessons/19-queue/homework-19/pkg/shortener/api"
	"go-core-4/lessons/19-queue/homework-19/pkg/shortener/config"
	"go-core-4/lessons/19-queue/homework-19/pkg/shortener/storage"
)

func main() {
	// 1. Загрузка параметров окружения
	cfg := config.Load()

	// 2. Инициализация In-Memory хранилища ссылок (по ТЗ курса)
	store := storage.NewSliceStore(cfg.BaseURL)
	seed(store)

	// 3. Инициализация продюсера Kafka
	producer, err := event.NewKafkaProducer(cfg.KafkaBrokers, cfg.KafkaTopic)
	if err != nil {
		log.Printf("⚠️ Предупреждение: Не удалось инициализировать Kafka Producer (%v). События не будут отправляться в Kafka.", err)
	} else {
		defer producer.Close()
		log.Printf("🔌 Kafka Producer подключен к брокерам %v (топик: %s)", cfg.KafkaBrokers, cfg.KafkaTopic)
	}

	// 4. Создание HTTP API
	appAPI := api.New(store, producer)

	// 5. Настройка HTTP-сервера с таймаутами
	server := &http.Server{
		Addr:         cfg.Port,
		Handler:      appAPI.Router(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// 6. Graceful Shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		fmt.Printf("🚀 [SHORTENER SERVICE] Запущен на порту %s\n", cfg.Port)
		fmt.Printf("   URL: %s\n", cfg.BaseURL)
		fmt.Println("   POST   /api/v1/links         — создание ссылки (+ публикация события в Kafka)")
		fmt.Println("   GET    /api/v1/links/{id}    — получение ссылки по ID")
		fmt.Println("   POST   /api/v1/links/resolve — поиск по короткой ссылке")
		fmt.Println("   GET    /r/{id}               — редирект (HTTP 302 Found)")
		fmt.Println("   GET    /api/v1/links         — просмотр всех ссылок")
		fmt.Println("   GET    /healthz              — Liveness probe")
		fmt.Println("   GET    /metrics              — Prometheus метрики")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ошибка HTTP-сервера: %v", err)
		}
	}()

	<-ctx.Done()
	fmt.Println("\n🛑 [SHORTENER SERVICE] Остановка микросервиса...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("ошибка принудительной остановки: %v", err)
	}

	fmt.Println("✅ [SHORTENER SERVICE] Остановлен корректно.")
}

func seed(store storage.Interface) {
	ctx := context.Background()
	_, _ = store.Save(ctx, "https://go.dev")
	_, _ = store.Save(ctx, "https://golang.org/doc/")
	_, _ = store.Save(ctx, "https://github.com/segmentio/kafka-go")
}
