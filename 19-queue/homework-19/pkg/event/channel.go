package event

import (
	"context"
	"errors"
	"sync"
)

var (
	// ErrBrokerClosed возвращается при попытке отправки в закрытый брокер.
	ErrBrokerClosed = errors.New("брокер очереди закрыт")
)

// ChannelBroker реализует легковесный брокер сообщений на каналах Go.
// Идеален для модульных тестов и автономного тестирования взаимодействия сервисов без Kafka.
type ChannelBroker struct {
	mu     sync.RWMutex
	ch     chan LinkCreatedEvent
	closed bool
}

// NewChannelBroker создает новый брокер на каналах Go с буфером указанного размера.
func NewChannelBroker(bufferSize int) *ChannelBroker {
	if bufferSize <= 0 {
		bufferSize = 100
	}
	return &ChannelBroker{
		ch: make(chan LinkCreatedEvent, bufferSize),
	}
}

// Publish отправляет событие в канал брокера.
func (b *ChannelBroker) Publish(ctx context.Context, _ string, event LinkCreatedEvent) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.closed {
		return ErrBrokerClosed
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case b.ch <- event:
		return nil
	}
}

// Subscribe непрерывно читает события из канала и передает их в handler.
func (b *ChannelBroker) Subscribe(ctx context.Context, handler Handler) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case ev, ok := <-b.ch:
			if !ok {
				return nil
			}
			if err := handler(ctx, ev); err != nil {
				// Логирование или обработка ошибки
				continue
			}
		}
	}
}

// Close закрывает канал брокера.
func (b *ChannelBroker) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if !b.closed {
		b.closed = true
		close(b.ch)
	}
	return nil
}
