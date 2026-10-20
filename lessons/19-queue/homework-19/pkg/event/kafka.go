package event

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

var (
	// ErrNoBrokers возвращается, если не указаны адреса брокеров Kafka.
	ErrNoBrokers = errors.New("не указаны брокеры Kafka")
	// ErrEmptyTopic возвращается, если имя топика пустое.
	ErrEmptyTopic = errors.New("не указан топик Kafka")
)

// KafkaProducer реализует Producer для Apache Kafka на базе github.com/segmentio/kafka-go.
type KafkaProducer struct {
	writer *kafka.Writer
	topic  string
}

// NewKafkaProducer создает нового Kafka-продюсера.
func NewKafkaProducer(brokers []string, topic string) (*KafkaProducer, error) {
	if len(brokers) == 0 || brokers[0] == "" {
		return nil, ErrNoBrokers
	}
	if topic == "" {
		return nil, ErrEmptyTopic
	}

	w := &kafka.Writer{
		Addr:                   kafka.TCP(brokers...),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
		WriteTimeout:           10 * time.Second,
		ReadTimeout:            10 * time.Second,
	}

	return &KafkaProducer{
		writer: w,
		topic:  topic,
	}, nil
}

// Publish сериализует событие в JSON и отправляет его в партиционированный топик Kafka.
func (p *KafkaProducer) Publish(ctx context.Context, key string, event LinkCreatedEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("ошибка сериализации события: %w", err)
	}

	msg := kafka.Message{
		Key:   []byte(key),
		Value: data,
		Time:  time.Now(),
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("ошибка отправки сообщения в топик %s: %w", p.topic, err)
	}

	return nil
}

// Close закрывает сетевые соединения писателя Kafka.
func (p *KafkaProducer) Close() error {
	if p.writer != nil {
		return p.writer.Close()
	}
	return nil
}

// KafkaConsumer реализует Consumer для Apache Kafka на базе consumer-группы.
type KafkaConsumer struct {
	reader *kafka.Reader
	topic  string
}

// NewKafkaConsumer создает нового потребителя сообщений Kafka.
func NewKafkaConsumer(brokers []string, topic, groupID string) (*KafkaConsumer, error) {
	if len(brokers) == 0 || brokers[0] == "" {
		return nil, ErrNoBrokers
	}
	if topic == "" {
		return nil, ErrEmptyTopic
	}
	if groupID == "" {
		return nil, errors.New("не указан group_id для Kafka Consumer")
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        groupID,
		MinBytes:       10e1, // 100B
		MaxBytes:       10e6, // 10MB
		CommitInterval: 0,    // Ручной коммит после успешной обработки (At-Least-Once Delivery)
		StartOffset:    kafka.FirstOffset,
	})

	return &KafkaConsumer{
		reader: r,
		topic:  topic,
	}, nil
}

// Subscribe запускает непрерывный цикл вычитки сообщений, обработки и фиксации оффсетов.
// Реализует паттерн «обработка -> подтверждение прочтения» (At-Least-Once Delivery).
func (c *KafkaConsumer) Subscribe(ctx context.Context, handler Handler) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return nil
			}
			return fmt.Errorf("ошибка чтения сообщения из Kafka: %w", err)
		}

		var ev LinkCreatedEvent
		if err := json.Unmarshal(msg.Value, &ev); err != nil {
			log.Printf("[KAFKA] Ошибка десериализации сообщения offset=%d: %v. Пропуск.", msg.Offset, err)
			// Коммитим некорректное сообщение (Poison Pill), чтобы не зависнуть в бесконечном цикле
			_ = c.reader.CommitMessages(ctx, msg)
			continue
		}

		// Вызов обработчика бизнес-логики
		if err := handler(ctx, ev); err != nil {
			log.Printf("[KAFKA] Ошибка обработки события link_id=%s: %v", ev.LinkID, err)
			// Не коммитим оффсет при временных сбоях для ретрая или обработка через DLQ
			continue
		}

		// Фиксация смещения после успешной обработки
		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			log.Printf("[KAFKA] Ошибка фиксации оффсета offset=%d: %v", msg.Offset, err)
		}
	}
}

// Close останавливает читателя Kafka.
func (c *KafkaConsumer) Close() error {
	if c.reader != nil {
		return c.reader.Close()
	}
	return nil
}
