// Package main демонстрирует базовые паттерны работы с распределенным брокером сообщений Apache Kafka.
// Используется библиотека github.com/segmentio/kafka-go — идиоматичный чистый Go-клиент без зависимости от librdkafka на C.
package main

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

// Client объединяет публикацию (Producer/Writer) и вычитку (Consumer/Reader) сообщений Kafka.
type Client struct {
	Reader *kafka.Reader
	Writer *kafka.Writer
}

// New создает и инициализирует клиента Kafka.
// brokers — список адресов брокеров кластера (bootstrap servers), e.g. ["localhost:29092"].
// topic — имя целевого топика Kafka.
// groupID — идентификатор группы потребителей (Consumer Group).
func New(brokers []string, topic string, groupID string) (*Client, error) {
	if len(brokers) == 0 || brokers[0] == "" || topic == "" || groupID == "" {
		return nil, errors.New("не указаны обязательные параметры подключения к Kafka")
	}

	c := Client{}

	// kafka.Reader настраивает высокоуровневого потребителя (Consumer Group).
	// Преимущества Consumer Group:
	// 1. Автоматический ребалансинг партиций топика между несколькими репликами приложения.
	// 2. Отслеживание смещения (Offset) прочитанных сообщений на стороне брокера.
	c.Reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e1, // 100 байт — минимальный размер порции данных от брокера
		MaxBytes: 10e6, // 10 МБ — максимальный размер пакета
	})

	// kafka.Writer настраивает производителя (Producer).
	// Обеспечивает буферизацию, объединение в батчи и балансировку сообщений по партициям.
	c.Writer = &kafka.Writer{
		Addr:                   kafka.TCP(brokers[0]),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{}, // Алгоритм выбора партиции с наименьшим объемом данных
		AllowAutoTopicCreation: true,                // Автоматическое создание топика при первой отправке
	}

	return &c, nil
}

// Close освобождает открытые TCP-соединения Reader и Writer.
func (c *Client) Close() error {
	var errReader, errWriter error
	if c.Reader != nil {
		errReader = c.Reader.Close()
	}
	if c.Writer != nil {
		errWriter = c.Writer.Close()
	}
	if errReader != nil {
		return errReader
	}
	return errWriter
}

func main() {
	// Инициализация клиента Kafka.
	kfk, err := New(
		[]string{"localhost:29092"},
		"test-topic",
		"test-consumer-group",
	)
	if err != nil {
		log.Fatal(err)
	}
	defer kfk.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Запуск публикации сообщений в отдельной горутине
	go kfk.producer(ctx)

	// Чтение сообщений в основном цикле
	kfk.consumer(ctx)
}

// producer отправляет сообщения с временной меткой в топик с интервалом в 1 секунду.
func (c *Client) producer(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case t := <-ticker.C:
			msg := kafka.Message{
				Key:   []byte("timestamp-key"),
				Value: []byte(t.Format(time.RFC3339Nano)),
			}
			err := c.Writer.WriteMessages(ctx, msg)
			if err != nil {
				log.Printf("[PRODUCER] Ошибка отправки сообщения: %v\n", err)
			} else {
				log.Printf("[PRODUCER] Отправлено: %s\n", string(msg.Value))
			}
		}
	}
}

// consumer реализует семантику доставки At-Least-Once (как минимум один раз):
// 1. FetchMessage: считывание очередного сообщения из топика (без авто-коммита).
// 2. Обработка сообщения бизнес-логикой.
// 3. CommitMessages: подтверждение смещения (commit offset) ТОЛЬКО после успешной обработки.
func (c *Client) consumer(ctx context.Context) {
	for {
		// FetchMessage блокирует поток до появления нового сообщения или отмены контекста
		msg, err := c.Reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			log.Printf("[CONSUMER] Ошибка получения сообщения: %v\n", err)
			continue
		}

		log.Printf("[CONSUMER] Получено: %s (Partition: %d, Offset: %d)\n",
			string(msg.Value), msg.Partition, msg.Offset)

		// Подтверждаем обработку сообщения в Kafka
		err = c.Reader.CommitMessages(ctx, msg)
		if err != nil {
			log.Printf("[CONSUMER] Ошибка фиксации смещения (Commit): %v\n", err)
		}
	}
}
