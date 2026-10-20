// Package main демонстрирует профилирование работающего HTTP-сервиса в реальном времени.
package main

import (
	"fmt"
	"math"
	"net/http"
	// Побочный импорт net/http/pprof регистрирует эндпоинты профилирования в http.DefaultServeMux:
	// - /debug/pprof/           (индексная страница)
	// - /debug/pprof/profile    (CPU-профиль, по умолчанию за 30 сек: ?seconds=30)
	// - /debug/pprof/heap       (профиль памяти/кучи)
	// - /debug/pprof/goroutine  (состояние и стектрейсы всех горутин)
	// - /debug/pprof/block      (задержки блокировок на каналах/мьютексах)
	// - /debug/pprof/mutex      (конкуренция за мьютексы)
	_ "net/http/pprof"
)

func main() {
	// Фоновые вычислительные горутины, создающие интенсивную нагрузку на CPU
	go compute1(1_000_000, 1_000_000)
	go compute2(2_000_000, 2_000_000)

	// Фоновые горутины, аллоцирующие память в куче
	go mem1()
	go mem2()

	// Используем порт :8080 (порт :80 на macOS/Linux требует прав администратора sudo)
	const addr = ":8080"
	fmt.Printf("HTTP-сервер с pprof запущен на http://localhost%s/debug/pprof/\n", addr)
	fmt.Println("Для снятия профиля выполните:")
	fmt.Println("  go tool pprof http://localhost:8080/debug/pprof/profile?seconds=10")
	fmt.Println("  go tool pprof -http=:8081 http://localhost:8080/debug/pprof/heap")

	// nil означает использование http.DefaultServeMux, куда pprof зарегистрировал свои роуты
	_ = http.ListenAndServe(addr, nil)
}

// compute1 генерирует CPU-нагрузку через тригонометрические вычисления math.Sin.
func compute1(n, m int) {
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			_ = math.Sin(float64(i)) * math.Sin(float64(j))
		}
	}
}

// compute2 генерирует CPU-нагрузку через math.Cos. В профиле CPU обе функции будут видны отдельно.
func compute2(n, m int) {
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			_ = math.Cos(float64(i)) * math.Cos(float64(j))
		}
	}
}

// mem1 аллоцирует срез из 10 млн элементов int32 (~40 МБ в куче).
func mem1() {
	var arr []int32
	for i := 0; i < 10_000_000; i++ {
		arr = append(arr, int32(i))
	}
	_ = arr
	select {} // удерживает горутину и массив в памяти
}

// mem2 аллоцирует срез из 20 млн элементов int64 (~160 МБ в куче).
// В отчете heap-профиля эта функция займет наибольший объем inuse_space.
func mem2() {
	var arr []int64
	for i := 20_000_000; i > 0; i-- {
		arr = append(arr, int64(i))
	}
	_ = arr
	select {} // удерживает горутину и массив в памяти
}

