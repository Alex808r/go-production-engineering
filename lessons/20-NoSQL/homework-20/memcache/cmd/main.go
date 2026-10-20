// Микросервис кэширования Lynks Memcache (Redis / In-Memory).
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

	"go-core-4/lessons/20-NoSQL/homework-20/memcache/pkg/api"
	"go-core-4/lessons/20-NoSQL/homework-20/memcache/pkg/config"
	"go-core-4/lessons/20-NoSQL/homework-20/memcache/pkg/store"
)

func main() {
	// 1. Инициализация структурированного логирования Zerolog (stdout)
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	logger.Info().Msg("Запуск микросервиса Lynks Memcache...")

	// 2. Чтение конфигурации
	cfg := config.Load()

	// 3. Подключение к Redis (или fallback на MemoryStore при отсутствии связи)
	var cacheStore store.Store
	redisStore, err := store.NewRedisStore(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB, cfg.CacheTTL)
	if err != nil {
		logger.Warn().Err(err).Msg("Не удалось подключиться к Redis. Используется резервное In-Memory хранилище кэша.")
		cacheStore = store.NewMemoryStore(cfg.CacheTTL)
	} else {
		cacheStore = redisStore
		logger.Info().Str("addr", cfg.RedisAddr).Dur("ttl", cfg.CacheTTL).Msg("Успешное подключение к Redis")
	}
	defer cacheStore.Close()

	// 4. Инициализация API
	appAPI := api.New(cacheStore, logger)

	server := &http.Server{
		Addr:         cfg.Port,
		Handler:      appAPI.Router(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// 5. Graceful Shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		fmt.Printf("⚡ [LYNKS MEMCACHE] Запущен на порту %s\n", cfg.Port)
		fmt.Println("   POST /api/v1/cache         — сохранить пару в кэш (TTL: 24h)")
		fmt.Println("   GET  /api/v1/cache/{short} — получить пару из кэша")
		fmt.Println("   GET  /healthz              — проверка здоровья")
		fmt.Println("   GET  /metrics              — Prometheus метрики")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("Ошибка запуска HTTP-сервера кэширования")
		}
	}()

	<-ctx.Done()
	logger.Info().Msg("Получен сигнал завершения. Остановка Memcache...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error().Err(err).Msg("Ошибка остановки сервера кэша")
	}

	logger.Info().Msg("Микросервис Memcache успешно остановлен.")
}
