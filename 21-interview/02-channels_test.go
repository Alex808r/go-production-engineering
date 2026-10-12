package interview

import (
	"testing"
)

func TestSelect(t *testing.T) {
	Select()
}

func TestProduceConsume(t *testing.T) {
	ProduceConsume()
}

func TestFanIn(t *testing.T) {
	ch1 := make(chan int)
	ch2 := make(chan int)
	go func() {
		defer close(ch1)
		for i := 0; i < 10; i++ {
			ch1 <- i
		}
	}()
	go func() {
		defer close(ch2)
		for i := 11; i < 20; i++ {
			ch2 <- i
		}
	}()

	ch := FanIn(ch1, ch2)
	count := 0
	for val := range ch {
		count++
		t.Log(val)
	}
	if count != 19 {
		t.Fatalf("expected 19 values from fan-in, got %d", count)
	}
}

func TestWorkerPool(t *testing.T) {
	const totalJobs = 50
	jobs := make(chan int, totalJobs)
	for i := 1; i <= totalJobs; i++ {
		jobs <- i
	}
	close(jobs)

	// Запуск пула из 5 воркеров
	results := WorkerPool(5, jobs, func(x int) int {
		return x * 10
	})

	var sum int
	for res := range results {
		sum += res
	}

	// Сумма от 1 до 50 = (1+50)*50/2 = 1275, умноженная на 10 = 12750
	expected := 12750
	if sum != expected {
		t.Fatalf("expected sum %d, got %d", expected, sum)
	}
}

func BenchmarkFanIn(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ch1 := make(chan int, 100)
		ch2 := make(chan int, 100)
		for j := 0; j < 100; j++ {
			ch1 <- j
			ch2 <- j
		}
		close(ch1)
		close(ch2)

		out := FanIn(ch1, ch2)
		for range out {
		}
	}
}

func BenchmarkWorkerPool(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		jobs := make(chan int, 100)
		for j := 0; j < 100; j++ {
			jobs <- j
		}
		close(jobs)

		results := WorkerPool(4, jobs, func(x int) int {
			return x * 2
		})
		for range results {
		}
	}
}

