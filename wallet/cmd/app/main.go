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

	"github.com/Alex808r/wallet/config"
	httpdelivery "github.com/Alex808r/wallet/internal/delivery/http"
	"github.com/Alex808r/wallet/internal/repository/postgres"
	"github.com/Alex808r/wallet/internal/service"
	internalsync "github.com/Alex808r/wallet/internal/sync"
)

func main() {
	log.Println("Starting Wallet Service...")

	// 1. Load configuration from config.env / environment
	cfg, err := config.Load("config.env")
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}
	log.Printf("Configuration loaded (Port: %s, DB: %s@%s:%s/%s)",
		cfg.AppPort, cfg.DBUser, cfg.DBHost, cfg.DBPort, cfg.DBName)

	// 2. Connect to PostgreSQL
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	repo, err := postgres.NewWalletRepo(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer repo.Close()
	log.Println("Connected to PostgreSQL successfully.")

	// 3. Run database migrations
	if err := repo.AutoMigrate(context.Background()); err != nil {
		log.Fatalf("failed to run database migrations: %v", err)
	}
	log.Println("Database schema checked and migrations applied.")

	// 4. Initialize layers
	locker := internalsync.NewKeyLocker(4096)
	svc := service.NewWalletService(repo, locker)
	handler := httpdelivery.NewWalletHandler(svc)
	router := httpdelivery.NewRouter(handler)

	// 5. Configure HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.AppPort,
		Handler:      router,
		ReadTimeout:  cfg.ServerReadTimeout,
		WriteTimeout: cfg.ServerWriteTimeout,
		IdleTimeout:  60 * time.Second,
	}

	// 6. Start server in background
	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("Wallet HTTP server listening on port %s", cfg.AppPort)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		}
	}()

	// 7. Graceful Shutdown on SIGINT / SIGTERM
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
