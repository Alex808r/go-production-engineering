package kvdb

import (
	"fmt"
	"sync"
	"testing"
)

func TestDB_CRUD(t *testing.T) {
	db := New()

	if db.Len() != 0 {
		t.Fatalf("ожидался пустой размер, получено: %d", db.Len())
	}

	// 1. Тест SET и GET
	db.SET("key1", "val1")
	db.SET("key2", "val2")

	if got := db.GET("key1"); got != "val1" {
		t.Errorf("GET(key1) = %q; ожидалось %q", got, "val1")
	}

	// 2. Тест Get с флагом существования
	val, ok := db.Get("key2")
	if !ok || val != "val2" {
		t.Errorf("Get(key2) = (%q, %v); ожидалось (%q, true)", val, ok, "val2")
	}

	// 3. Тест отсутствующего ключа
	if got := db.GET("nonexistent"); got != "" {
		t.Errorf("GET(nonexistent) = %q; ожидалась пустая строка", got)
	}

	val, ok = db.Get("nonexistent")
	if ok || val != "" {
		t.Errorf("Get(nonexistent) = (%q, %v); ожидалось (\"\", false)", val, ok)
	}

	// 4. Тест Len
	if db.Len() != 2 {
		t.Errorf("Len() = %d; ожидалось 2", db.Len())
	}

	// 5. Тест DELETE
	db.DELETE("key1")
	if db.Len() != 1 {
		t.Errorf("Len() после DELETE = %d; ожидалось 1", db.Len())
	}

	if got := db.GET("key1"); got != "" {
		t.Errorf("GET(key1) после удаления = %q; ожидалась пустая строка", got)
	}
}

func TestDB_Concurrency(t *testing.T) {
	db := New()
	const numGoroutines = 50
	const numOps = 100

	var wg sync.WaitGroup
	wg.Add(numGoroutines * 2)

	// Писатели
	for i := 0; i < numGoroutines; i++ {
		go func(writerID int) {
			defer wg.Done()
			for j := 0; j < numOps; j++ {
				key := fmt.Sprintf("key_%d_%d", writerID, j)
				val := fmt.Sprintf("val_%d_%d", writerID, j)
				db.SET(key, val)
			}
		}(i)
	}

	// Читатели
	for i := 0; i < numGoroutines; i++ {
		go func(readerID int) {
			defer wg.Done()
			for j := 0; j < numOps; j++ {
				key := fmt.Sprintf("key_%d_%d", readerID, j)
				_ = db.GET(key)
				_, _ = db.Get(key)
				_ = db.Len()
			}
		}(i)
	}

	wg.Wait()

	if db.Len() != numGoroutines*numOps {
		t.Errorf("ожидалось %d записей, получено: %d", numGoroutines*numOps, db.Len())
	}
}

func BenchmarkDB_SET(b *testing.B) {
	db := New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		db.SET("benchmark_key", "benchmark_val")
	}
}

func BenchmarkDB_GET(b *testing.B) {
	db := New()
	db.SET("benchmark_key", "benchmark_val")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = db.GET("benchmark_key")
	}
}

func BenchmarkDB_ParallelGET(b *testing.B) {
	db := New()
	db.SET("benchmark_key", "benchmark_val")
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = db.GET("benchmark_key")
		}
	})
}
