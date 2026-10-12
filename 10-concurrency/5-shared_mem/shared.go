// Package main демонстрирует работу с разделяемой памятью (Shared Memory)
// и защиту критических секций от состояния гонки (Data Race) с помощью sync.Mutex.
package main

import (
	"fmt"
	"sync"
)

// SafeCounter инкапсулирует разделяемое состояние и мьютекс для защиты данных.
// В Go принято упаковывать мьютекс вместе с защищаемыми им полями в структуру!
type SafeCounter struct {
	mu sync.Mutex // защищает поле count
	// Примечание: если операций чтения намного больше, чем записи,
	// вместо sync.Mutex используют sync.RWMutex (RLock/RUnlock для чтения).
	count int
}

// Inc безопасно увеличивает счетчик на 1 внутри критической секции.
func (c *SafeCounter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock() // гарантирует освобождение блокировки при выходе
	c.count++
}

// Value безопасно возвращает текущее значение счетчика.
func (c *SafeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

func main() {
	// ========================================================================
	// ТЕОРИЯ: КРИТИЧЕСКАЯ СЕКЦИЯ И ГОНКИ ДАННЫХ (DATA RACE)
	// ========================================================================
	// Операция `count++` в Go НЕ является атомарной!
	// На уровне процессора она раскладывается минимум на 3 инструкции:
	//  1. Чтение count из памяти в регистр CPU.
	//  2. Инкремент регистра.
	//  3. Запись регистра обратно в память.
	//
	// Если две горутины делают это одновременно без мьютекса, возникает Data Race
	// («потерянное обновление» / lost update), и результат будет непредсказуемым!
	//
	// Проверить наличие гонок данных всегда можно флагом компилятора:
	//   go run -race ./10-concurrency/5-shared_mem/shared.go

	var counter SafeCounter
	var wg sync.WaitGroup

	const goroutines = 2
	const iterations = 10_000

	wg.Add(goroutines)

	// Запускаем 2 конкурентные горутины
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				counter.Inc()
			}
		}()
	}

	wg.Wait()

	// Ожидаемый результат: ровно 20 000
	fmt.Printf("Итоговое значение счётчика: %d (ожидалось: %d)\n", counter.Value(), goroutines*iterations)
}

