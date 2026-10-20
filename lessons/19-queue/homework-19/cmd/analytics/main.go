// Микросервис аналитики ссылок (Consumer из Kafka).
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

	"go-core-4/lessons/19-queue/homework-19/pkg/analytics/api"
	"go-core-4/lessons/19-queue/homework-19/pkg/analytics/config"
	"go-core-4/lessons/19-queue/homework-19/pkg/analytics/consumer"
	"go-core-4/lessons/19-queue/homework-19/pkg/analytics/service"
	"go-core-4/lessons/19-queue/homework-19/pkg/event"
)

func main() {
	// 1. Загрузка параметров окружения
	cfg := config.Load()

	// 2. Инициализация бизнес-сервиса аналитики
	svc := service.New()

	// 3. Обработка сигналов завершения (Graceful Shutdown)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 4. Подключение Kafka Consumer
	kafkaConsumer, err := event.NewKafkaConsumer(cfg.KafkaBrokers, cfg.KafkaTopic, cfg.KafkaGroupID)
	if err != nil {
		log.Printf("⚠️ Предупреждение: Не удалось инициализировать Kafka Consumer (%v). Вычитка очереди отключена.", err)
	} else {
		defer kafkaConsumer.Close()
		log.Printf("🔌 Kafka Consumer подключен к брокерам %v (топик: %s, group: %s)",
			cfg.KafkaBrokers, cfg.KafkaTopic, cfg.KafkaGroupID)

		// Запуск фонового воркера вычитки очереди
		worker := consumer.NewWorker(kafkaConsumer, svc)
		go func() {
			if err := worker.Start(ctx); err != nil {
				log.Printf("ошибка воркера консьюмера: %v", err)
			}
		}()
	}

	// 5. Создание HTTP API
	appAPI := api.New(svc)

	server := &http.Server{
		Addr:         cfg.Port,
		Handler:      appAPI.Router(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		fmt.Printf("📊 [ANALYTICS SERVICE] Запущен на порту %s\n", cfg.Port)
		fmt.Println("   GET    /api/v1/stats — агрегированная статистика по ссылкам (JSON)")
		fmt.Println("   GET    /healthz      — Liveness probe")
		fmt.Println("   GET    /ready        — Readiness probe")
		fmt.Println("   GET    /metrics      — Prometheus метрики")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ошибка HTTP-сервера аналитики: %v", err)
		}
	}()

	<-ctx.Done()
	fmt.Println("\n🛑 [ANALYTICS SERVICE] Остановка микросервиса...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("ошибка остановки HTTP-сервера: %v", err)
	}

	fmt.Println("✅ [ANALYTICS SERVICE] Остановлен корректно.")
}
