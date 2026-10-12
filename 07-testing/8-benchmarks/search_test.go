package bsearch

import (
	"math/rand"
	"sort"
	"testing"
	"time"
)

// TestSearches проверяет корректность работы обоих алгоритмов поиска на эталонных примерах.
func TestSearches(t *testing.T) {
	tests := []struct {
		name string
		data []int
		item int
		want int
	}{
		{
			name: "Найден в середине нечетного среза",
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
			name: "Элемент отсутствует в срезе",
			data: []int{1, 2, 4, 5, 6, 10},
			item: 16,
			want: -1,
		},
		{
			name: "Первые позиции среза",
			data: []int{1, 2, 4, 5, 6, 10},
			item: 2,
			want: 1,
		},
		{
			name: "Большой диапазон чисел",
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

// sampleData генерирует отсортированный срез из 1_000_000 случайных чисел для бенчмарков.
func sampleData() []int {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	data := make([]int, 1_000_000)
	for i := 0; i < 1_000_000; i++ {
		data[i] = rng.Intn(1000)
	}

	sort.Ints(data)
	return data
}

// Глобальная переменная для предотвращения оптимизации компилятора (Compiler Optimization).
var benchmarkResult int

// BenchmarkBinary замеряет скорость бинарного поиска O(log N).
//
// АНАТОМИЯ БЕНЧМАРКА В GO:
// 1. Имя функции должно начинаться с `Benchmark` и принимать `(b *testing.B)`.
// 2. b.ResetTimer(): сбрасывает счетчик времени бенчмарка после тяжелой фазы подготовки данных
//    (sampleData), чтобы генерация и сортировка 1 млн чисел не учитывались в метриках поиска.
// 3. b.ReportAllocs(): включает автоматический учет выделения памяти (B/op и allocs/op).
// 4. Цикл `for i := 0; i < b.N; i++`: фреймворк Go автоматически подбирает значение b.N
//    так, чтобы тест длился заданное время (по умолчанию 1 секунду) для репрезентативной статистики.
func BenchmarkBinary(b *testing.B) {
	data := sampleData()
	b.ReportAllocs()
	b.ResetTimer()

	var res int
	for i := 0; i < b.N; i++ {
		res = Binary(data, 500)
	}
	benchmarkResult = res
}

// BenchmarkSimple замеряет скорость линейного поиска O(N).
// Показывает разницу в производительности: линейный поиск будет медленнее бинарного в тысячи раз.
func BenchmarkSimple(b *testing.B) {
	data := sampleData()
	b.ReportAllocs()
	b.ResetTimer()

	var res int
	for i := 0; i < b.N; i++ {
		res = Simple(data, 500)
	}
	benchmarkResult = res
}

