package bsearch

import (
	"math/rand"
	"sort"
	"testing"
	"time"
)

// TestSearches проверяет корректность работы алгоритмов на наборе тестовых сценариев.
func TestSearches(t *testing.T) {
	tests := []struct {
		name string
		data []int
		item int
		want int
	}{
		{
			name: "Найден в нечетном срезе",
			data: []int{1, 2, 4, 6, 10},
			item: 6,
			want: 3,
		},
		{
			name: "Найден в четном срезе",
			data: []int{1, 2, 4, 5, 6, 10},
			item: 6,
			want: 4,
		},
		{
			name: "Элемент отсутствует",
			data: []int{1, 2, 4, 5, 6, 10},
			item: 16,
			want: -1,
		},
		{
			name: "Первый элемент среза",
			data: []int{1, 2, 4, 5, 6, 10},
			item: 2,
			want: 1,
		},
		{
			name: "Крайний элемент среза",
			data: []int{1, 2, 4, 5, 6, 10, 100},
			item: 4,
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Binary(tt.data, tt.item); got != tt.want {
				t.Errorf("Binary() = %v, want %v", got, tt.want)
			}
			if got := Simple(tt.data, tt.item); got != tt.want {
				t.Errorf("Simple() = %v, want %v", got, tt.want)
			}
		})
	}
}

// sampleData генерирует миллион случайных чисел для профилирования бенчмарков.
func sampleData() []int {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	data := make([]int, 1_000_000)
	for i := 0; i < 1_000_000; i++ {
		data[i] = rng.Intn(1000)
	}

	sort.Ints(data)
	return data
}

// Глобальная переменная для предотвращения оптимизации компилятором.
var benchSink int

// BenchmarkBinary замеряет бинарный поиск и используется для снятия CPU/Heap профилей.
//
// КОМАНДЫ ДЛЯ СБОРА И АНАЛИЗА ПРОФИЛЕЙ PPROF:
//
// 1. Сбор профиля процессора (CPU) и памяти (Heap):
//    go test -bench=BenchmarkBinary -cpuprofile=cpu.out -memprofile=mem.out
//
// 2. Интерактивный анализ в веб-интерфейсе браузера:
//    go tool pprof -http=:8080 cpu.out
//    (Откроется веб-UI со вкладками: Top, Graph, Flame Graph, Peek, Source, Disassemble)
//
// 3. Анализ в терминале (CLI):
//    go tool pprof cpu.out
//    В интерактивном шелле pprof доступны команды:
//    - `top` (вывести самые тяжелые функции)
//    - `list Binary` (построчный анализ исходного кода с затратами процессорного времени)
//    - `web` (сгенерировать SVG-граф вызовов)
func BenchmarkBinary(b *testing.B) {
	data := sampleData()
	b.ResetTimer()

	var res int
	for i := 0; i < b.N; i++ {
		res = Binary(data, 500)
	}
	benchSink = res
}

// BenchmarkSimple замеряет линейный поиск и демонстрирует высокую нагрузку на CPU в профилях pprof.
func BenchmarkSimple(b *testing.B) {
	data := sampleData()
	b.ResetTimer()

	var res int
	for i := 0; i < b.N; i++ {
		res = Simple(data, 500)
	}
	benchSink = res
}

