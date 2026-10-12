// Package sorts предоставляет утилиты генерации тестовых наборов данных
// для тестирования и бенчмаркинга стандартных алгоритмов сортировки Go (sort.Ints, sort.Float64s, sort.Strings).
package sorts

import (
	"math/rand"
	"time"
)

// GenerateInts генерирует срез случайных целых чисел размера n в диапазоне [0, 1_000_000).
func GenerateInts(n int) []int {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	res := make([]int, n)
	for i := 0; i < n; i++ {
		res[i] = rng.Intn(1_000_000)
	}
	return res
}

// GenerateFloat64s генерирует срез случайных чисел с плавающей точкой размера n.
func GenerateFloat64s(n int) []float64 {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	res := make([]float64, n)
	for i := 0; i < n; i++ {
		res[i] = rng.Float64() * 1_000_000
	}
	return res
}
