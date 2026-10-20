package gracefulshutdown

import (
	"context"
	"io"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestGracefulShutdown_InFlightRequestCompletes(t *testing.T) {
	// Находим свободный случайный порт
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen on random port: %v", err)
	}
	addr := l.Addr().String()
	_ = l.Close() // освобождаем порт для сервера

	app := NewApp(addr)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serverDone := make(chan error, 1)
	go func() {
		serverDone <- app.Start(ctx, 2*time.Second)
	}()

	// Ждем, пока сервер поднимется
	baseURL := "http://" + addr
	for i := 0; i < 20; i++ {
		resp, err := http.Get(baseURL + "/healthz")
		if err == nil && resp.StatusCode == http.StatusOK {
			_ = resp.Body.Close()
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	// Запускаем долгий запрос в отдельной горутине
	reqFinished := make(chan struct{})
	go func() {
		defer close(reqFinished)
		resp, err := http.Get(baseURL + "/slow")
		if err != nil {
			t.Errorf("slow request failed: %v", err)
			return
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusOK || string(body) != "task completed successfully" {
			t.Errorf("unexpected response: status=%d, body=%s", resp.StatusCode, string(body))
		}
	}()

	// Даем запросу начаться
	time.Sleep(15 * time.Millisecond)

	// Инициируем остановку сервера (имитация SIGINT / SIGTERM)
	cancel()

	// Убеждаемся, что долгий запрос успел успешно завершиться!
	select {
	case <-reqFinished:
		// Успешно! In-flight запрос не был оборван
	case <-time.After(1 * time.Second):
		t.Fatalf("slow request timed out")
	}

	// Убеждаемся, что Start завершился без ошибок
	select {
	case err := <-serverDone:
		if err != nil {
			t.Fatalf("unexpected server start/shutdown error: %v", err)
		}
	case <-time.After(1 * time.Second):
		t.Fatalf("server shutdown did not complete in time")
	}
}
