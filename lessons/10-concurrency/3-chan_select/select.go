// Package main демонстрирует работу оператора select в Go:
// мультиплексирование каналов, псевдослучайный выбор, отключение веток через nil-канал и таймауты.
package main

import (
	"fmt"
	"math/rand"
	"time"
)

// generator генерирует сообщения с задержкой и закрывает канал.
func generator(ch chan<- string, num int) {
	for i := 1; i <= 3; i++ {
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		ch <- fmt.Sprintf("Канал #%d: сообщение %d", num, i)
	}
	close(ch)
}

func main() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go generator(ch1, 1)
	go generator(ch2, 2)

	// ========================================================================
	// ТЕОРИЯ: ПРАВИЛЬНАЯ РАБОТА С SELECT И ЗАКРЫТЫМИ КАНАЛАМИ
	// ========================================================================
	// 1. Если несколько case готовы одновременно, select выбирает один из них СЛУЧАЙНО.
	// 2. Если канал закрылся, чтение из него будет ВСЕГДА возвращать zero-value мгновенно!
	//    Поэтому нужно проверять `ok`: если !ok, мы обнуляем переменную канала (ch = nil).
	//    Так как чтение из nil-канала блокируется навсегда, этот case больше НИКОГДА не сработает!
	// 3. АНТИПАТТЕРН: В бесконечном цикле `default: continue` превращает select в 100% busy-wait,
	//    сжигая ядро процессора! Без default горутина мирно спит в ожидании данных.

	fmt.Println("--- Мультиплексирование двух каналов через select ---")

	for ch1 != nil || ch2 != nil {
		select {
		case val, ok := <-ch1:
			if !ok {
				fmt.Println("Канал #1 закрыт -> отключаем ветку (ch1 = nil)")
				ch1 = nil // Отключаем ветку в select!
				continue
			}
			fmt.Println("<-", val)

		case val, ok := <-ch2:
			if !ok {
				fmt.Println("Канал #2 закрыт -> отключаем ветку (ch2 = nil)")
				ch2 = nil // Отключаем ветку в select!
				continue
			}
			fmt.Println("<-", val)

		case <-time.After(500 * time.Millisecond):
			// Защита от зависания: таймаут на операцию
			fmt.Println("Таймаут: за 500 мс сообщений не поступило!")
		}
	}

	fmt.Println("--- Все каналы прочитаны и отключены ---")

	// 4. НЕБЛОКИРУЮЩЕЕ ЧТЕНИЕ / ЗАПИСЬ ЧЕРЕЗ DEFAULT
	// Ветка default нужна ТОЛЬКО тогда, когда мы хотим попытаться прочитать/записать
	// БЕЗ блокировки, и если канал не готов — сразу выполнить альтернативное действие:
	nonBlockingCh := make(chan int, 1)
	select {
	case val := <-nonBlockingCh:
		fmt.Println("Прочитано без блокировки:", val)
	default:
		fmt.Println("Канал пуст: выполнена ветка default без блокировки!")
	}
}

