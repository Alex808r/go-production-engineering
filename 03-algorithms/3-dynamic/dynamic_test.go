package main

import (
	"math/big"
	"testing"
)

func TestFibonacciMethods(t *testing.T) {
	expected := []int{1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144, 233, 377, 610}

	// 1. Проверка fibo (рекурсия с мемоизацией)
	cache := make(map[int]int)
	for i, exp := range expected {
		if got := fibo(i, cache); got != exp {
			t.Errorf("fibo(%d) = %d; want %d", i, got, exp)
		}
	}

	// 2. Проверка fiboIterative
	for i, exp := range expected {
		if got := fiboIterative(i); got != exp {
			t.Errorf("fiboIterative(%d) = %d; want %d", i, got, exp)
		}
	}

	// 3. Проверка fiboClosure
	gen := fiboClosure()
	for i, exp := range expected {
		if got := gen(); got != exp {
			t.Errorf("fiboClosure() step %d = %d; want %d", i, got, exp)
		}
	}

	// 4. Проверка fiboSeq
	seqResults := make([]int, 0, len(expected))
	for v := range fiboSeq() {
		seqResults = append(seqResults, v)
		if len(seqResults) == len(expected) {
			break
		}
	}
	for i, exp := range expected {
		if seqResults[i] != exp {
			t.Errorf("fiboSeq()[%d] = %d; want %d", i, seqResults[i], exp)
		}
	}

	// 5. Проверка fiboBig
	for i, exp := range expected {
		if got := fiboBig(i); got.Cmp(big.NewInt(int64(exp))) != 0 {
			t.Errorf("fiboBig(%d) = %s; want %d", i, got.String(), exp)
		}
	}
}

// Бенчмарк вычисления F(50) через рекурсию с мемоизацией
func BenchmarkFiboMemo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cache := make(map[int]int, 64)
		_ = fibo(50, cache)
	}
}

// Бенчмарк вычисления F(50) через итеративный цикл (Bottom-Up)
func BenchmarkFiboIterative(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = fiboIterative(50)
	}
}

// Бенчмарк генератора через замыкание (50 шагов)
func BenchmarkFiboClosure(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gen := fiboClosure()
		for j := 0; j < 50; j++ {
			_ = gen()
		}
	}
}

// Бенчмарк вычисления F(50) через math/big
func BenchmarkFiboBig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = fiboBig(50)
	}
}
