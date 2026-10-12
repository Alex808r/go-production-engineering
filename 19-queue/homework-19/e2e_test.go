package homework19_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	analyticsAPI "go-core-4/19-queue/homework-19/pkg/analytics/api"
	analyticsConsumer "go-core-4/19-queue/homework-19/pkg/analytics/consumer"
	analyticsService "go-core-4/19-queue/homework-19/pkg/analytics/service"
	"go-core-4/19-queue/homework-19/pkg/event"
	shortenerAPI "go-core-4/19-queue/homework-19/pkg/shortener/api"
	shortenerStorage "go-core-4/19-queue/homework-19/pkg/shortener/storage"
)

// TestE2E_ShortenerToAnalytics_AsyncFlow проверяет полный сквозной сценарий взаимодействия
// микросервиса сокращения ссылок и микросервиса аналитики через брокер сообщений:
// 1. Клиент делает POST /api/v1/links в Shortener.
// 2. Shortener сохраняет ссылку в срез памяти и публикует LinkCreatedEvent в очередь.
// 3. Analytics Consumer вычитывает событие из очереди и обновляет агрегированную статистику.
// 4. Клиент делает GET /api/v1/stats в Analytics и получает актуальные метрики (кол-во, средняя длина, домены).
func TestE2E_ShortenerToAnalytics_AsyncFlow(t *testing.T) {
	// Создаем общий брокер событий (ChannelBroker)
	broker := event.NewChannelBroker(50)
	defer broker.Close()

	// 1. Инициализация Shortener микросервиса
	store := shortenerStorage.NewSliceStore("http://localhost:8080")
	shortener := shortenerAPI.New(store, broker)

	// 2. Инициализация Analytics микросервиса
	analyticsSvc := analyticsService.New()
	worker := analyticsConsumer.NewWorker(broker, analyticsSvc)
	analytics := analyticsAPI.New(analyticsSvc)

	// Запускаем фоновый воркер консьюмера аналитики
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		_ = worker.Start(ctx)
	}()

	testLinks := []string{
		"https://go.dev/doc/tutorial/getting-started",
		"https://github.com/segmentio/kafka-go",
		"https://microservices.io/patterns/data/event-sourcing.html",
	}

	for _, u := range testLinks {
		body, _ := json.Marshal(map[string]string{"url": u})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/links", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		shortener.Router().ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("failed to create link for %s: code %d, body: %s", u, rec.Code, rec.Body.String())
		}
	}

	// Ожидаем асинхронную обработку всех 3 событий консьюмером
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if worker.ProcessedCount() == uint64(len(testLinks)) {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}

	if worker.ProcessedCount() != uint64(len(testLinks)) {
		t.Fatalf("expected 3 processed events in analytics, got %d", worker.ProcessedCount())
	}

	// 3. Запрос статистики через HTTP API микросервиса аналитики: GET /api/v1/stats
	reqStats := httptest.NewRequest(http.MethodGet, "/api/v1/stats", nil)
	recStats := httptest.NewRecorder()
	analytics.Router().ServeHTTP(recStats, reqStats)

	if recStats.Code != http.StatusOK {
		t.Fatalf("expected 200 OK from analytics, got %d", recStats.Code)
	}

	var stats analyticsService.Stats
	if err := json.Unmarshal(recStats.Body.Bytes(), &stats); err != nil {
		t.Fatalf("failed to decode stats: %v", err)
	}

	// Проверяем точность вычисленных метрик
	if stats.TotalLinks != 3 {
		t.Errorf("expected 3 total links, got %d", stats.TotalLinks)
	}

	expectedLen := int64(len(testLinks[0]) + len(testLinks[1]) + len(testLinks[2]))
	if stats.TotalOriginalLength != expectedLen {
		t.Errorf("expected total length %d, got %d", expectedLen, stats.TotalOriginalLength)
	}

	if stats.DomainCounts["go.dev"] != 1 || stats.DomainCounts["github.com"] != 1 || stats.DomainCounts["microservices.io"] != 1 {
		t.Errorf("unexpected domain distribution: %+v", stats.DomainCounts)
	}

	if stats.AvgOriginalLength <= 0 || stats.AvgShortLength <= 0 {
		t.Errorf("average lengths should be positive: orig=%f, short=%f", stats.AvgOriginalLength, stats.AvgShortLength)
	}
}
