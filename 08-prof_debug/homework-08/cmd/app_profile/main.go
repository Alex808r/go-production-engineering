package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	_ "net/http/pprof"
	"sync"
	"time"
)

// Программа демонстрирует профилирование живого сервиса через net/http/pprof:
// 1. CPU-профиль: http://localhost:6060/debug/pprof/profile?seconds=5
// 2. Heap-профиль (память): http://localhost:6060/debug/pprof/heap
// 3. Goroutine-профиль: http://localhost:6060/debug/pprof/goroutine
// 4. Block/Mutex профилирование
func main() {
	addr := ":6060"
	fmt.Printf("=== Сервер профилирования pprof запущен на http://localhost%s/debug/pprof/ ===\n", addr)

	// Фоновые нагрузочные воркеры для генерации реального профиля:
	go runCPUWorkers()
	go runMemoryLeakSimulation()
	go runGoroutineWorkers()

	// Запуск HTTP-сервера со стандартными обработчиками pprof
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Ошибка запуска HTTP сервера: %v", err)
	}
}

// runCPUWorkers генерирует интенсивную нагрузку на CPU (тригонометрические вычисления)
func runCPUWorkers() {
	for {
		var wg sync.WaitGroup
		for w := 0; w < 4; w++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()
				computeHeavyTrig(workerID, 500_000)
			}(w)
		}
		wg.Wait()
		time.Sleep(100 * time.Millisecond)
	}
}

func computeHeavyTrig(id, iterations int) float64 {
	var acc float64
	for i := 0; i < iterations; i++ {
		acc += math.Sin(float64(i)) * math.Cos(float64(i+id))
	}
	return acc
}

// runMemoryLeakSimulation непрерывно аллоцирует структуры в куче без преаллокации,
// создавая давление на память и сборщик мусора (GC pressure).
func runMemoryLeakSimulation() {
	var retainedChunks [][]byte

	for {
		// Создаем временные аллокации (GC pressure)
		for i := 0; i < 50; i++ {
			temp := make([]byte, 100*1024) // 100 KB
			for j := range temp {
				temp[j] = byte(i)
			}
			_ = temp
		}

		// Удерживаем часть памяти в куче (демонстрация inuse_space)
		retained := make([]byte, 512*1024) // 512 KB
		retainedChunks = append(retainedChunks, retained)
		if len(retainedChunks) > 40 {
			// Ограничиваем максимальный размер удержания ~20 MB
			retainedChunks = retainedChunks[10:]
		}

		time.Sleep(200 * time.Millisecond)
	}
}

// runGoroutineWorkers порождает пулы горутин для анализа через /debug/pprof/goroutine
func runGoroutineWorkers() {
	for {
		for i := 0; i < 15; i++ {
			go func(id int) {
				time.Sleep(500 * time.Millisecond)
			}(i)
		}
		time.Sleep(1 * time.Second)
	}
}
