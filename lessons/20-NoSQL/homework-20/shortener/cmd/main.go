// Микросервис Lynks Shortener (PostgreSQL + Memcache клиент).
// Итоговый проект курса (Занятие №20: NoSQL).
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"

	"go-core-4/lessons/20-NoSQL/homework-20/shortener/pkg/api"
	"go-core-4/lessons/20-NoSQL/homework-20/shortener/pkg/cacheclient"
	"go-core-4/lessons/20-NoSQL/homework-20/shortener/pkg/config"
	"go-core-4/lessons/20-NoSQL/homework-20/shortener/pkg/storage"
)

func main() {
	// 1. Zerolog
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	logger.Info().Msg("Запуск микросервиса Lynks Shortener...")

	// 2. Чтение конфигурации
	cfg := config.Load()

	// 3. Подключение к PostgreSQL (или In-Memory при отсутствии СУБД)
	var mainStorage storage.Storage
	pgStorage, err := storage.NewPostgresStorage(cfg.PgDSN)
	if err != nil {
		logger.Warn().Err(err).Msg("Не удалось подключиться к PostgreSQL. Используется резервное In-Memory хранилище.")
		mainStorage = storage.NewMemoryStorage()
	} else {
		mainStorage = pgStorage
		logger.Info().Str("dsn", cfg.PgDSN).Msg("Успешное подключение к PostgreSQL")
	}
	defer mainStorage.Close()

	// 4. Подключение HTTP-клиента к микросервису Memcache
	cacheClient := cacheclient.NewHTTPClient(cfg.MemcacheURL, 2*time.Second)
	logger.Info().Str("memcache_url", cfg.MemcacheURL).Msg("Клиент кэширования инициализирован")

	// 5. Инициализация API
	appAPI := api.New(mainStorage, cacheClient, cfg.BaseURL, logger)

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
		fmt.Printf("🚀 [LYNKS SHORTENER] Запущен на порту %s (BaseURL: %s)\n", cfg.Port, cfg.BaseURL)
		fmt.Println("   POST /                      — сократить гиперссылку")
		fmt.Println("   POST /api/v1/links          — сократить гиперссылку (REST-алиас)")
		fmt.Println("   GET  /{short}               — редирект на оригинал (Cache-Aside)")
		fmt.Println("   GET  /healthz               — проверка здоровья")
		fmt.Println("   GET  /metrics               — Prometheus метрики")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("Ошибка запуска HTTP-сервера shortener")
		}
	}()

	<-ctx.Done()
	logger.Info().Msg("Получен сигнал завершения. Остановка Shortener...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error().Err(err).Msg("Ошибка при остановке сервера shortener")
	}

	logger.Info().Msg("Микросервис Shortener успешно остановлен.")
}
