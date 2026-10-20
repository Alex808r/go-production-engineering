package consumer

import (
	"context"
	"testing"
	"time"

	"go-core-4/lessons/19-queue/homework-19/pkg/analytics/service"
	"go-core-4/lessons/19-queue/homework-19/pkg/event"
)

func TestWorker_Processing(t *testing.T) {
	broker := event.NewChannelBroker(10)
	defer broker.Close()

	svc := service.New()
	w := NewWorker(broker, svc)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- w.Start(ctx)
	}()

	// Отправляем тестовые события
	ev1 := event.LinkCreatedEvent{
		EventID:   "evt-101",
		LinkID:    "id-101",
		Original:  "https://pkg.go.dev/github.com/segmentio/kafka-go",
		ShortURL:  "http://localhost:8080/r/id-101",
		CreatedAt: time.Now(),
	}
	_ = broker.Publish(context.Background(), ev1.LinkID, ev1)

	// Ждем обработки
	for i := 0; i < 20; i++ {
		if w.ProcessedCount() >= 1 {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}

	if w.ProcessedCount() != 1 {
		t.Fatalf("expected 1 processed event, got %d", w.ProcessedCount())
	}

	stats := svc.GetStats()
	if stats.TotalLinks != 1 || stats.DomainCounts["pkg.go.dev"] != 1 {
		t.Errorf("stats mismatch: %+v", stats)
	}

	cancel()
	select {
	case err := <-errCh:
		if err != nil {
			t.Errorf("expected clean exit, got error: %v", err)
		}
	case <-time.After(1 * time.Second):
		t.Error("timeout waiting for worker shutdown")
	}
}
