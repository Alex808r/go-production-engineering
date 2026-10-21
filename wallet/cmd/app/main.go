package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	// --------------------------------------------------------------------------
	// Внутренние пакеты сервиса wallet:
	// Префикс импорта "github.com/Alex808r/wallet" совпадает с именем модуля в go.mod.
	// Компилятор Go разрешает эти импорты ЛОКАЛЬНО в каталоге проекта без сетевых обращений.
	// Использование полного URL-пути — каноничный стандарт (Best Practice) экосистемы Go.
	// --------------------------------------------------------------------------
	"github.com/Alex808r/wallet/config"
	httpdelivery "github.com/Alex808r/wallet/internal/delivery/http"
	"github.com/Alex808r/wallet/internal/repository/postgres"
	"github.com/Alex808r/wallet/internal/service"
	internalsync "github.com/Alex808r/wallet/internal/sync"
)

func main() {
	log.Println("Starting Wallet Service...")

	// --------------------------------------------------------------------------
	// 1. Загрузка конфигурации:
	// Считывает параметры из файла config.env (согласно ТЗ) и системных переменных окружения.
	// --------------------------------------------------------------------------
	cfg, err := config.Load("config.env")
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}
	log.Printf("Configuration loaded (Port: %s, DB: %s@%s:%s/%s)",
		cfg.AppPort, cfg.DBUser, cfg.DBHost, cfg.DBPort, cfg.DBName)

	// --------------------------------------------------------------------------
	// 2. Подключение к базе данных PostgreSQL:
	// Инициализирует пул соединений pgxpool с таймаутом на случай задержек сети.
	// --------------------------------------------------------------------------
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	repo, err := postgres.NewWalletRepo(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer repo.Close()
	log.Println("Connected to PostgreSQL successfully.")

	// --------------------------------------------------------------------------
	// 3. Автоматические миграции базы данных:
	// Проверяет наличие таблиц wallets и wallet_transactions, создает их при необходимости.
	// Это обеспечивает полную автономность сервиса при первом старте.
	// --------------------------------------------------------------------------
	if err := repo.AutoMigrate(context.Background()); err != nil {
		log.Fatalf("failed to run database migrations: %v", err)
	}
	log.Println("Database schema checked and migrations applied.")

	// --------------------------------------------------------------------------
	// 4. Сборка слоев приложения (Clean Architecture / Dependency Injection):
	// - KeyLocker(4096): шардированный локер на 4096 полос для защиты от 1000 RPS.
	// - WalletService: бизнес-логика изменения баланса и координации с локером.
	// - WalletHandler: HTTP-обработчики (REST API).
	// - Router: маршрутизация Go 1.22+ с middleware логирования и паник-рекавери.
	// --------------------------------------------------------------------------
	locker := internalsync.NewKeyLocker(4096)
	svc := service.NewWalletService(repo, locker)
	handler := httpdelivery.NewWalletHandler(svc)
	router := httpdelivery.NewRouter(handler)

	// --------------------------------------------------------------------------
	// 5. Конфигурация HTTP-сервера:
	// Задаем явные таймауты чтения/записи, чтобы защитить сервер от зависших Slowloris-атак
	// и утечек ресурсов под нагрузкой.
	// --------------------------------------------------------------------------
	server := &http.Server{
		Addr:         ":" + cfg.AppPort,
		Handler:      router,
		ReadTimeout:  cfg.ServerReadTimeout,
		WriteTimeout: cfg.ServerWriteTimeout,
		IdleTimeout:  60 * time.Second,
	}

	// --------------------------------------------------------------------------
	// 6. Асинхронный запуск HTTP-сервера:
	// Запускаем сервер в фоновой горутине, чтобы основной поток мог слушать сигналы ОС.
	// --------------------------------------------------------------------------
	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("Wallet HTTP server listening on port %s", cfg.AppPort)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		}
	}()

	// --------------------------------------------------------------------------
	// 7. Элегантное завершение (Graceful Shutdown):
	// Слушаем системные сигналы SIGINT (Ctrl+C) и SIGTERM (Docker stop / Kubernetes).
	// Даем активным транзакциям завершиться без потери клиентских денег (ServerShutdownTimeout),
	// после чего закрываем пул соединений с базой.
	// --------------------------------------------------------------------------
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Fatalf("server error: %v", err)
	case sig := <-shutdown:
		log.Printf("Received signal %v. Initiating graceful shutdown...", sig)

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.ServerShutdownTimeout)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("Graceful shutdown failed: %v. Forcing close...", err)
			_ = server.Close()
		}
		log.Println("Server stopped cleanly. Goodbye!")
	}
}
