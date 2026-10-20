package event

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestChannelBroker_PublishAndSubscribe(t *testing.T) {
	broker := NewChannelBroker(10)
	defer broker.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		received []LinkCreatedEvent
		mu       sync.Mutex
		wg       sync.WaitGroup
	)

	wg.Add(3)

	go func() {
		_ = broker.Subscribe(ctx, func(_ context.Context, ev LinkCreatedEvent) error {
			mu.Lock()
			received = append(received, ev)
			mu.Unlock()
			wg.Done()
			return nil
		})
	}()

	now := time.Now()
	events := []LinkCreatedEvent{
		{
			EventID:   "evt-1",
			Type:      "link.created",
			LinkID:    "abc12345",
			Original:  "https://go.dev",
			ShortURL:  "http://localhost:8080/r/abc12345",
			CreatedAt: now,
		},
		{
			EventID:   "evt-2",
			Type:      "link.created",
			LinkID:    "xyz67890",
			Original:  "https://github.com/golang/go",
			ShortURL:  "http://localhost:8080/r/xyz67890",
			CreatedAt: now,
		},
		{
			EventID:   "evt-3",
			Type:      "link.created",
			LinkID:    "def11223",
			Original:  "https://thinknetica.com",
			ShortURL:  "http://localhost:8080/r/def11223",
			CreatedAt: now,
		},
	}

	for _, ev := range events {
		err := broker.Publish(ctx, ev.LinkID, ev)
		if err != nil {
			t.Fatalf("unexpected error publishing event: %v", err)
		}
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for events in subscriber")
	}

	mu.Lock()
	defer mu.Unlock()
	if len(received) != 3 {
		t.Fatalf("expected 3 received events, got %d", len(received))
	}
	if received[0].LinkID != "abc12345" || received[1].Original != "https://github.com/golang/go" {
		t.Errorf("unexpected event content: %+v", received)
	}
}

func TestChannelBroker_Closed(t *testing.T) {
	broker := NewChannelBroker(1)
	if err := broker.Close(); err != nil {
		t.Fatalf("failed to close broker: %v", err)
	}

	err := broker.Publish(context.Background(), "test", LinkCreatedEvent{EventID: "1"})
	if err != ErrBrokerClosed {
		t.Fatalf("expected ErrBrokerClosed, got %v", err)
	}
}

func TestKafkaConfig_Validation(t *testing.T) {
	_, err := NewKafkaProducer(nil, "topic")
	if err != ErrNoBrokers {
		t.Errorf("expected ErrNoBrokers, got %v", err)
	}

	_, err = NewKafkaProducer([]string{"localhost:9092"}, "")
	if err != ErrEmptyTopic {
		t.Errorf("expected ErrEmptyTopic, got %v", err)
	}

	_, err = NewKafkaConsumer(nil, "topic", "group")
	if err != ErrNoBrokers {
		t.Errorf("expected ErrNoBrokers, got %v", err)
	}

	_, err = NewKafkaConsumer([]string{"localhost:9092"}, "topic", "")
	if err == nil {
		t.Error("expected error for empty groupID, got nil")
	}
}

func BenchmarkChannelBroker_Publish(b *testing.B) {
	broker := NewChannelBroker(b.N + 10)
	defer broker.Close()
	ctx := context.Background()
	ev := LinkCreatedEvent{
		EventID:  "bench-id",
		LinkID:   "abc",
		Original: "https://example.com/bench",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = broker.Publish(ctx, ev.LinkID, ev)
	}
}

