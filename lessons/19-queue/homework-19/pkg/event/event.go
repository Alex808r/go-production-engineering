// Package event определяет структуры событий и абстракции брокера сообщений
// для асинхронного взаимодействия микросервисов (согласно лекции №19).
package event

import (
	"context"
	"time"
)

// LinkCreatedEvent описывает событие успешного создания новой короткой гиперссылки.
type LinkCreatedEvent struct {
	EventID   string    `json:"event_id"`
	Type      string    `json:"type"` // "link.created"
	LinkID    string    `json:"link_id"`
	Original  string    `json:"original"`
	ShortURL  string    `json:"short_url"`
	CreatedAt time.Time `json:"created_at"`
}

// Handler определяет функцию обратного вызова для обработки события консьюмером.
type Handler func(ctx context.Context, event LinkCreatedEvent) error

// Producer описывает контракт поставщика событий в очередь (Kafka / RabbitMQ / Channel).
type Producer interface {
	// Publish отправляет событие с указанным ключом партиционирования.
	Publish(ctx context.Context, key string, event LinkCreatedEvent) error
	// Close освобождает ресурсы подключения к брокеру.
	Close() error
}

// Consumer описывает контракт потребителя сообщений из очереди.
type Consumer interface {
	// Subscribe запускает блокирующее чтение сообщений из очереди до отмены контекста.
	Subscribe(ctx context.Context, handler Handler) error
	// Close останавливает консьюмер и закрывает соединение.
	Close() error
}
