package storage

import (
	"context"
	"sync"
	"testing"
)

func TestMemoryStorage_CRUD(t *testing.T) {
	s := NewMemoryStorage()
	defer s.Close()
	ctx := context.Background()

	// 1. Сохранение
	err := s.Save(ctx, "qz6d7", "https://github.com/thinknetica/go_course_4")
	if err != nil {
		t.Fatalf("unexpected Save error: %v", err)
	}

	// 2. Поиск по короткой ссылке
	dest, err := s.FindByShort(ctx, "qz6d7")
	if err != nil {
		t.Fatalf("unexpected FindByShort error: %v", err)
	}
	if dest != "https://github.com/thinknetica/go_course_4" {
		t.Errorf("destination mismatch: %s", dest)
	}

	// 3. Поиск по длинной ссылке
	short, err := s.FindByDestination(ctx, "https://github.com/thinknetica/go_course_4")
	if err != nil {
		t.Fatalf("unexpected FindByDestination error: %v", err)
	}
	if short != "qz6d7" {
		t.Errorf("short mismatch: %s", short)
	}

	// 4. Ошибки при пустых аргументах
	if err := s.Save(ctx, "", "https://go.dev"); err != ErrEmptyField {
		t.Errorf("expected ErrEmptyField, got %v", err)
	}

	if _, err := s.FindByShort(ctx, "nonexistent"); err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestMemoryStorage_Concurrency(t *testing.T) {
	s := NewMemoryStorage()
	ctx := context.Background()

	const count = 50
	var wg sync.WaitGroup
	wg.Add(count)

	for i := 0; i < count; i++ {
		go func(idx int) {
			defer wg.Done()
			short := "short"
			dest := "https://example.com"
			_ = s.Save(ctx, short, dest)
			_, _ = s.FindByShort(ctx, short)
		}(i)
	}

	wg.Wait()
}

func BenchmarkMemoryStorage_Save(b *testing.B) {
	s := NewMemoryStorage()
	defer s.Close()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.Save(ctx, "bench_short", "https://golang.org")
	}
}

func BenchmarkMemoryStorage_FindByShort(b *testing.B) {
	s := NewMemoryStorage()
	defer s.Close()
	ctx := context.Background()
	_ = s.Save(ctx, "bench_short", "https://golang.org")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = s.FindByShort(ctx, "bench_short")
	}
}

func BenchmarkMemoryStorage_ParallelFindByShort(b *testing.B) {
	s := NewMemoryStorage()
	defer s.Close()
	ctx := context.Background()
	_ = s.Save(ctx, "bench_short", "https://golang.org")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = s.FindByShort(ctx, "bench_short")
		}
	})
}
