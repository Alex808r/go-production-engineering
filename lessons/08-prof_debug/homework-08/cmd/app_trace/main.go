package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"runtime/trace"
	"sync"
	"time"
)

func main() {
	traceFilePath := "trace.out"
	f, err := os.Create(traceFilePath)
	if err != nil {
		log.Fatalf("Ошибка создания trace.out: %v", err)
	}
	defer f.Close()

	// 1. Запуск трассировки execution tracer
	if err := trace.Start(f); err != nil {
		log.Fatalf("Ошибка старта трассировки: %v", err)
	}
	defer trace.Stop()

	fmt.Printf("=== Трассировка выполнения начата, запись в %s ===\n", traceFilePath)

	ctx, task := trace.NewTask(context.Background(), "homework08_trace_task")
	defer task.End()

	// Фаза 1: Многопоточные CPU-вычисления (распределение по процессорам P)
	trace.WithRegion(ctx, "CPU_Computation_Phase", func() {
		runParallelComputation(runtime.GOMAXPROCS(0) * 2)
	})

	// Фаза 2: Память и принудительное давление на сборщик мусора GC
	trace.WithRegion(ctx, "Memory_GC_Pressure_Phase", func() {
		runMemoryChurn()
	})

	// Фаза 3: Взаимодействие через каналы (блокировки и передача сообщений)
	trace.WithRegion(ctx, "Channel_Synchronization_Phase", func() {
		runChannelPipeline(5000)
	})

	fmt.Println("=== Трассировка успешно завершена. Для просмотра выполните: go tool trace trace.out ===")
}

func runParallelComputation(workers int) {
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for i := 0; i < 200_000; i++ {
				_ = math.Sqrt(float64(i)) * math.Sin(float64(i+id))
			}
		}(w)
	}
	wg.Wait()
}

func runMemoryChurn() {
	for cycle := 0; cycle < 5; cycle++ {
		// Создаем аллокации
		data := make([][]byte, 1000)
		for i := range data {
			data[i] = make([]byte, 10*1024) // 10 KB
		}
		// Запускаем сборку мусора для отображения в трейсе
		runtime.GC()
		time.Sleep(10 * time.Millisecond)
	}
}

func runChannelPipeline(items int) {
	ch := make(chan int, 10)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < items; i++ {
			ch <- i
		}
		close(ch)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for val := range ch {
			_ = val * 2
		}
	}()

	wg.Wait()
}
