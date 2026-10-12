// Package main демонстрирует трассировку планировщика горутин и загрузки ядер процессора
// с помощью стандартного инструмента runtime/trace и утилиты go tool trace.
package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"runtime/trace"
	"time"
)

// КАК ЗАПУСТИТЬ И АНАЛИЗИРОВАТЬ ТРАССИРОВКУ (EXECUTION TRACER):
//
// 1. Скомпилировать и запустить программу:
//    go run trace.go
//    (Создастся бинарный файл трассировки `trace.out`)
//
// 2. Открыть визуализатор трассировки в браузере:
//    go tool trace trace.out
//
// ЧТО ПОКАЗЫВАЕТ GO TOOL TRACE:
// - View trace: интерактивный хронологический таймлайн (навигация клавишами W, A, S, D).
// - Процессоры P (Logical Processors): загрузка виртуальных ядер, на которых выполняются горутины.
// - Потоки ОС (OS Threads / M): переключение контекста между потоками операционной системы.
// - Горутины (Goroutines / G): состояния горутин (Running, Runnable, Blocked, Waiting on GC).
// - Heap & GC: динамика выделения кучи и фазы сборщика мусора (GC sweep, Mark assist, STW).
// - Network & Sync: задержки сетевого ввода-вывода и блокировки каналов/мьютексов.
func main() {
	// Создаем файл для сохранения бинарных событий трассировки
	f, err := os.Create("trace.out")
	if err != nil {
		log.Fatalf("ошибка создания trace.out: %v", err)
	}
	defer f.Close()

	// Запуск трассировщика runtime/trace
	if err := trace.Start(f); err != nil {
		log.Fatalf("ошибка запуска трассировки: %v", err)
	}
	defer trace.Stop()

	fmt.Println("Запись трассировки CPU начата (1 секунда)...")

	// Порождаем конкурентные горутины, соревнующиеся за кванты времени CPU
	go compute1()
	go compute1()
	go compute2()
	go compute2()

	// Даем горутинам поработать 1 секунду под наблюдением трассировщика
	time.Sleep(1 * time.Second)

	fmt.Println("Трассировка записана в trace.out. Для просмотра выполните:")
	fmt.Println("  go tool trace trace.out")
}

func compute1() {
	for {
		_ = math.Sin(1) * math.Cos(10)
	}
}

func compute2() {
	for {
		_ = math.Tan(1)
	}
}

