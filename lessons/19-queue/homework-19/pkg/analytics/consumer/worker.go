// Package consumer реализует фоновый воркер вычитки событий из очереди сообщений.
package consumer

import (
	"context"
	"errors"
	"log"
	"sync/atomic"

	"go-core-4/lessons/19-queue/homework-19/pkg/analytics/service"
	"go-core-4/lessons/19-queue/homework-19/pkg/event"
)

// Worker управляет подпиской на очередь и передачей событий в сервис аналитики.
type Worker struct {
	consumer       event.Consumer
	analyticsSvc   *service.Service
	processedCount uint64
}

// NewWorker создает новый воркер обработки очереди событий.
func NewWorker(consumer event.Consumer, analyticsSvc *service.Service) *Worker {
	return &Worker{
		consumer:     consumer,
		analyticsSvc: analyticsSvc,
	}
}

// Start запускает вычитку сообщений. Метод блокирует выполнение до отмены ctx.
func (w *Worker) Start(ctx context.Context) error {
	log.Println("[ANALYTICS-CONSUMER] Запуск подписки на события создания ссылок...")

	err := w.consumer.Subscribe(ctx, func(_ context.Context, ev event.LinkCreatedEvent) error {
		w.analyticsSvc.RecordLink(ev)
		atomic.AddUint64(&w.processedCount, 1)
		log.Printf("[ANALYTICS-CONSUMER] Обработано событие link_id=%s, original=%s (всего обработано: %d)",
			ev.LinkID, ev.Original, atomic.LoadUint64(&w.processedCount))
		return nil
	})

	if err != nil && !errors.Is(err, context.Canceled) {
		return err
	}

	log.Println("[ANALYTICS-CONSUMER] Подписка остановлена корректно.")
	return nil
}

// ProcessedCount возвращает количество успешно обработанных событий.
func (w *Worker) ProcessedCount() uint64 {
	return atomic.LoadUint64(&w.processedCount)
}
