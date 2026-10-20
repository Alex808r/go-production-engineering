package interview

import (
	"fmt"
	"sync"
)

func Select() {
	ch := make(chan string)
	select {
	case ch <- "string":
		fmt.Println("получилось записать в канал")
	default:
		fmt.Println("не получилось записать в канал")
	}
}

// ProduceConsume демонстрирует паттерн "поставщик-потребитель".
// Важно: горутина-поставщик обязана закрыть канал (defer close(ch)),
// иначе цикл `for val := range ch` у потребителя заблокируется навсегда (deadlock).
func ProduceConsume() {
	ch := make(chan int)

	go func() {
		defer close(ch)
		for i := 0; i < 10; i++ {
			ch <- i
		}
	}()

	for val := range ch {
		fmt.Println(val)
	}
}

func FanIn[T any](channels ...<-chan T) <-chan T {
	ch := make(chan T)
	var wg sync.WaitGroup
	wg.Add(len(channels))

	for _, c := range channels {
		go func(in <-chan T) {
			defer wg.Done()
			for i := range in {
				ch <- i
			}
		}(c)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	return ch
}

// WorkerPool реализует классический пул горутин-обработчиков:
// Запускается numWorkers воркеров, которые параллельно читают задачи из входного канала jobs,
// применяют к ним функцию process и отправляют результат в результирующий канал.
// Канал results закрывается только после того, как все воркеры завершат обработку (wg.Wait).
func WorkerPool[T any, R any](numWorkers int, jobs <-chan T, process func(T) R) <-chan R {
	results := make(chan R)
	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for w := 0; w < numWorkers; w++ {
		go func() {
			defer wg.Done()
			for job := range jobs {
				results <- process(job)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}

