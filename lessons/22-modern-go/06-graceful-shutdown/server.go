// Package gracefulshutdown демонстрирует эталонный паттерн корректной остановки
// HTTP-сервера и фоновых горутин с использованием signal.NotifyContext и srv.Shutdown.
package gracefulshutdown

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

// App инкапсулирует HTTP-сервер и фоновые процессы.
type App struct {
	server       *http.Server
	inFlightReqs atomic.Int64
	bgWorkersWg  sync.WaitGroup
	shutdownDone chan struct{}
}

// NewApp инициализирует приложение.
func NewApp(addr string) *App {
	app := &App{
		shutdownDone: make(chan struct{}),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /slow", app.handleSlowTask)
	mux.HandleFunc("GET /healthz", app.handleHealthz)

	app.server = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return app
}

func (a *App) handleSlowTask(w http.ResponseWriter, r *http.Request) {
	a.inFlightReqs.Add(1)
	defer a.inFlightReqs.Add(-1)

	// Имитация долгой работы (обработка платежа или тяжелого запроса)
	select {
	case <-time.After(50 * time.Millisecond):
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, "task completed successfully")
	case <-r.Context().Done():
		// Клиент оборвал соединение
		http.Error(w, "request cancelled", http.StatusRequestTimeout)
	}
}

func (a *App) handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, "OK")
}

// Start запускает сервер и ожидает сигнала завершения из контекста ctx.
func (a *App) Start(ctx context.Context, shutdownTimeout time.Duration) error {
	serverErrChan := make(chan error, 1)

	// 1. Запуск слушателя в отдельной горутине
	go func() {
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrChan <- err
		}
		close(serverErrChan)
	}()

	// 2. Блокируемся до получения сигнала отмены (SIGINT / SIGTERM)
	select {
	case <-ctx.Done():
		// Получен сигнал на остановку
	case err := <-serverErrChan:
		return fmt.Errorf("ошибка запуска HTTP-сервера: %w", err)
	}

	// 3. Создаем таймаут для завершения активных соединений (Graceful Shutdown)
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// Shutdown прекращает прием новых соединений и ждет завершения активных
	if err := a.server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("принудительное завершение сервера по таймауту: %w", err)
	}

	// 4. Ожидаем завершения фоновых воркеров
	a.bgWorkersWg.Wait()
	close(a.shutdownDone)

	return nil
}

// ShutdownDone возвращает канал, сигнализирующий о полном завершении приложения.
func (a *App) ShutdownDone() <-chan struct{} {
	return a.shutdownDone
}
