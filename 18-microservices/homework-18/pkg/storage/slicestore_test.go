package storage_test

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"

	"go-core-4/18-microservices/homework-18/pkg/storage"
)

func TestSliceStore_SaveAndGet(t *testing.T) {
	ctx := context.Background()
	store := storage.NewSliceStore("http://localhost:8080")

	// 1. Сохранение валидного URL
	original := "https://thinknetica.com/courses/golang"
	link, err := store.Save(ctx, original)
	if err != nil {
		t.Fatalf("unexpected error on save: %v", err)
	}

	if link.Original != original {
		t.Errorf("expected original %q, got %q", original, link.Original)
	}
	if link.ID == "" {
		t.Error("expected non-empty ID")
	}
	if link.ShortURL != "http://localhost:8080/r/"+link.ID {
		t.Errorf("unexpected ShortURL: %s", link.ShortURL)
	}

	// 2. Получение по ID
	foundByID, err := store.Get(ctx, link.ID)
	if err != nil {
		t.Fatalf("failed to get by ID: %v", err)
	}
	if foundByID.Original != original {
		t.Errorf("expected %q, got %q", original, foundByID.Original)
	}

	// 3. Получение по полному короткому URL
	foundByShort, err := store.Get(ctx, link.ShortURL)
	if err != nil {
		t.Fatalf("failed to get by short URL: %v", err)
	}
	if foundByShort.ID != link.ID {
		t.Errorf("expected link ID %s, got %s", link.ID, foundByShort.ID)
	}

	// 4. Повторное сохранение того же URL (идемпотентность)
	dup, err := store.Save(ctx, original)
	if err != nil {
		t.Fatalf("duplicate save error: %v", err)
	}
	if dup.ID != link.ID {
		t.Errorf("expected same ID for duplicate original URL, got %s vs %s", link.ID, dup.ID)
	}

	all, _ := store.All(ctx)
	if len(all) != 1 {
		t.Errorf("expected 1 unique link in slice, got %d", len(all))
	}
}

func TestSliceStore_ValidationAndNotFound(t *testing.T) {
	ctx := context.Background()
	store := storage.NewSliceStore("http://localhost:8080")

	// Пустой URL
	if _, err := store.Save(ctx, ""); !errors.Is(err, storage.ErrEmptyURL) {
		t.Errorf("expected ErrEmptyURL, got: %v", err)
	}

	// Невалидный URL
	if _, err := store.Save(ctx, "ftp://invalid-scheme.com"); !errors.Is(err, storage.ErrInvalidURL) {
		t.Errorf("expected ErrInvalidURL, got: %v", err)
	}

	// Несуществующая ссылка
	if _, err := store.Get(ctx, "non-existent-id"); !errors.Is(err, storage.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got: %v", err)
	}
}

func TestSliceStore_Concurrency(t *testing.T) {
	ctx := context.Background()
	store := storage.NewSliceStore("http://localhost:8080")

	var wg sync.WaitGroup
	// Параллельная запись
	for i := 0; i < 30; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			_, _ = store.Save(ctx, fmt.Sprintf("https://example.com/page/%d", idx))
		}(i)
	}

	// Параллельное чтение
	for i := 0; i < 30; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = store.All(ctx)
		}()
	}

	wg.Wait()
}

func BenchmarkSliceStore_Save(b *testing.B) {
	ctx := context.Background()
	store := storage.NewSliceStore("http://localhost:8080")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = store.Save(ctx, fmt.Sprintf("https://example.com/%d", i))
	}
}

func BenchmarkSliceStore_Get(b *testing.B) {
	ctx := context.Background()
	store := storage.NewSliceStore("http://localhost:8080")
	link, _ := store.Save(ctx, "https://example.com/target")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = store.Get(ctx, link.ID)
	}
}

