package store

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestMemoryStore_BasicOperations(t *testing.T) {
	s := NewMemoryStore(10 * time.Minute)
	defer s.Close()

	ctx := context.Background()

	// 1. Успешная запись и чтение
	err := s.Set(ctx, "abc12", "https://golang.org")
	if err != nil {
		t.Fatalf("unexpected Set error: %v", err)
	}

	val, err := s.Get(ctx, "abc12")
	if err != nil {
		t.Fatalf("unexpected Get error: %v", err)
	}
	if val != "https://golang.org" {
		t.Errorf("expected https://golang.org, got %s", val)
	}

	// 2. Несуществующий ключ
	_, err = s.Get(ctx, "nonexistent")
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}

	// 3. Пустой ключ
	err = s.Set(ctx, "", "https://go.dev")
	if err != ErrEmptyKey {
		t.Errorf("expected ErrEmptyKey, got %v", err)
	}

	// 4. Удаление
	_ = s.Delete(ctx, "abc12")
	_, err = s.Get(ctx, "abc12")
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound after Delete, got %v", err)
	}
}

func TestMemoryStore_TTLExpiration(t *testing.T) {
	// TTL 50 миллисекунд
	s := NewMemoryStore(50 * time.Millisecond)
	ctx := context.Background()

	_ = s.Set(ctx, "short_ttl", "https://example.com")

	val, err := s.Get(ctx, "short_ttl")
	if err != nil || val != "https://example.com" {
		t.Fatalf("key should be present immediately")
	}

	time.Sleep(70 * time.Millisecond)

	_, err = s.Get(ctx, "short_ttl")
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound after TTL expiration, got %v", err)
	}
}

func TestMemoryStore_Concurrency(t *testing.T) {
	s := NewMemoryStore(time.Hour)
	ctx := context.Background()

	const workers = 50
	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func(idx int) {
			defer wg.Done()
			key := "key"
			_ = s.Set(ctx, key, "val")
			_, _ = s.Get(ctx, key)
		}(i)
	}

	wg.Wait()
}

func BenchmarkMemoryStore_Set(b *testing.B) {
	s := NewMemoryStore(time.Hour)
	defer s.Close()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.Set(ctx, "bench_key", "https://golang.org")
	}
}

func BenchmarkMemoryStore_Get(b *testing.B) {
	s := NewMemoryStore(time.Hour)
	defer s.Close()
	ctx := context.Background()
	_ = s.Set(ctx, "bench_key", "https://golang.org")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = s.Get(ctx, "bench_key")
	}
}

func BenchmarkMemoryStore_ParallelGet(b *testing.B) {
	s := NewMemoryStore(time.Hour)
	defer s.Close()
	ctx := context.Background()
	_ = s.Set(ctx, "bench_key", "https://golang.org")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = s.Get(ctx, "bench_key")
		}
	})
}
