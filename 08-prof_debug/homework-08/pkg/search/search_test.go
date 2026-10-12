package search

import (
	"math/rand"
	"sort"
	"testing"
)

func TestSearch(t *testing.T) {
	data := []int{1, 3, 5, 7, 9, 11, 15, 20, 25, 30}

	tests := []struct {
		item int
		want int
	}{
		{item: 1, want: 0},
		{item: 7, want: 3},
		{item: 30, want: 9},
		{item: 4, want: -1},
		{item: 100, want: -1},
	}

	for _, tt := range tests {
		if got := Binary(data, tt.item); got != tt.want {
			t.Errorf("Binary(data, %d) = %d, want %d", tt.item, got, tt.want)
		}
		if got := Linear(data, tt.item); got != tt.want {
			t.Errorf("Linear(data, %d) = %d, want %d", tt.item, got, tt.want)
		}
	}
}

func prepareSortedData(n int) []int {
	r := rand.New(rand.NewSource(42)) // Детерминированный генератор
	data := make([]int, n)
	for i := 0; i < n; i++ {
		data[i] = r.Intn(n * 2)
	}
	sort.Ints(data)
	return data
}

// BenchmarkBinarySearch замеряет скорость O(log N) поиска
func BenchmarkBinarySearch_100k(b *testing.B) {
	data := prepareSortedData(100_000)
	r := rand.New(rand.NewSource(42))
	queries := make([]int, 1000)
	for i := range queries {
		queries[i] = r.Intn(200_000)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		target := queries[i%len(queries)]
		_ = Binary(data, target)
	}
}

// BenchmarkLinearSearch замеряет скорость O(N) поиска
func BenchmarkLinearSearch_100k(b *testing.B) {
	data := prepareSortedData(100_000)
	r := rand.New(rand.NewSource(42))
	queries := make([]int, 1000)
	for i := range queries {
		queries[i] = r.Intn(200_000)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		target := queries[i%len(queries)]
		_ = Linear(data, target)
	}
}
