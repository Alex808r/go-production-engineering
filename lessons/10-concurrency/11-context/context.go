// Package main демонстрирует работу со стандартным пакетом context в Go:
// управление временем жизни горутин, отмену по таймауту (WithTimeout) и передачу контекстных метаданных (WithValue).
package main

import (
	"context"
	"fmt"
	"time"
)

// type contextKey используется для предотвращения коллизий ключей в context.WithValue
type contextKey string

const requestIDKey contextKey = "requestID"

// worker выполняет обработку поступающих данных с поддержкой отмены через контекст.
func worker(ctx context.Context, in <-chan string, out chan<- string) {
	for {
		select {
		case <-ctx.Done():
			// Ветка отмены: срабатывает при вызове cancel(), истечении таймаута или дедлайна.
			// ctx.Err() содержит причину завершения: context.Canceled или context.DeadlineExceeded
			fmt.Printf("Воркер получил сигнал завершения: %v\n", ctx.Err())
			if reqID := ctx.Value(requestIDKey); reqID != nil {
				fmt.Printf("Контекстный RequestID: %v\n", reqID)
			}
			close(out)
			return

		case val, ok := <-in:
			if !ok {
				// Входной канал закрыт отправителем
				close(out)
				return
			}
			fmt.Println("Воркер обрабатывает:", val)
			out <- fmt.Sprintf("Обработано: %s", val)
		}
	}
}

func main() {
	// ========================================================================
	// ТЕОРИЯ: ПАКЕТ CONTEXT В GO
	// ========================================================================
	// context.Context решает 3 ключевые задачи в конкурентных системах:
	//  1. Сигнал отмены (Cancellation): досрочная остановка долгих операций при разрыве соединения.
	//  2. Дедлайны и таймауты (Timeouts): ограничение максимального времени ожидания ответа БД/API.
	//  3. Метаданные запроса (Values): сквозная передача trace ID, токенов авторизации через слои.
	//
	// ЗОЛОТОЕ ПРАВИЛО:
	// Всегда вызывайте `defer cancel()` сразу после создания контекста с таймаутом/отменой,
	// чтобы остановить внутренний таймер и не допустить утечки памяти (context leak)!

	data := make(chan string)
	result := make(chan string)

	// Создаем контекст с таймаутом 600 миллисекунд
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Millisecond)
	defer cancel()

	// Оборачиваем в контекст с метаданными (WithValue)
	ctx = context.WithValue(ctx, requestIDKey, "req-xyz-98765")

	fmt.Println("--- Запуск горутины с контекстом (таймаут 600 мс) ---")
	go worker(ctx, data, result)

	// Горутина-отправитель
	go func() {
		defer close(data)
		for i := 1; i <= 5; i++ {
			select {
			case <-ctx.Done():
				fmt.Println("Отправитель прерван контекстом!")
				return
			case data <- fmt.Sprintf("сообщение #%d", i):
				time.Sleep(200 * time.Millisecond)
			}
		}
	}()

	// Вычитываем результаты, пока канал `result` не закроется воркером при отмене контекста
	for val := range result {
		fmt.Println("-> Получен результат:", val)
	}

	fmt.Println("--- Завершение работы программы ---")
}

