package storage

import (
	"context"
	"fmt"
	"sync"
	"testing"
)

func TestSliceStore_SaveAndGet(t *testing.T) {
	store := NewSliceStore("http://localhost:8080")
	ctx := context.Background()

	link, err := store.Save(ctx, "https://golang.org/doc/")
	if err != nil {
		t.Fatalf("unexpected error saving link: %v", err)
	}

	if link.ID == "" {
		t.Error("expected non-empty ID")
	}
	if link.ShortURL != "http://localhost:8080/r/"+link.ID {
		t.Errorf("unexpected short url: %s", link.ShortURL)
	}

	// Повторное сохранение не дублирует
	sameLink, err := store.Save(ctx, "https://golang.org/doc/")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sameLink.ID != link.ID {
		t.Errorf("expected duplicate to return same link id %s, got %s", link.ID, sameLink.ID)
	}

	// Поиск по ID
	byID, err := store.Get(ctx, link.ID)
	if err != nil {
		t.Fatalf("unexpected error getting by ID: %v", err)
	}
	if byID.Original != "https://golang.org/doc/" {
		t.Errorf("unexpected original URL: %s", byID.Original)
	}

	// Поиск по short URL
	byShort, err := store.Get(ctx, link.ShortURL)
	if err != nil {
		t.Fatalf("unexpected error getting by short url: %v", err)
	}
	if byShort.ID != link.ID {
		t.Errorf("expected same ID")
	}

	// All
	all, err := store.All(ctx)
	if err != nil {
		t.Fatalf("unexpected error getting all: %v", err)
	}
	if len(all) != 1 {
		t.Errorf("expected 1 link, got %d", len(all))
	}
}

func TestSliceStore_Validation(t *testing.T) {
	store := NewSliceStore("http://localhost:8080")
	ctx := context.Background()

	_, err := store.Save(ctx, "")
	if err != ErrEmptyURL {
		t.Errorf("expected ErrEmptyURL, got %v", err)
	}

	_, err = store.Save(ctx, "not-a-valid-url")
	if err == nil {
		t.Error("expected error for invalid URL, got nil")
	}

	_, err = store.Get(ctx, "nonexistent")
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestSliceStore_Concurrency(t *testing.T) {
	store := NewSliceStore("http://localhost:8080")
	ctx := context.Background()

	const count = 50
	var wg sync.WaitGroup
	wg.Add(count)

	for i := 0; i < count; i++ {
		go func(idx int) {
			defer wg.Done()
			url := fmt.Sprintf("https://example.com/page/%d", idx)
			link, err := store.Save(ctx, url)
			if err != nil {
				t.Errorf("error saving %s: %v", url, err)
				return
			}
			_, _ = store.Get(ctx, link.ID)
		}(i)
	}

	wg.Wait()

	all, err := store.All(ctx)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(all) != count {
		t.Errorf("expected %d links, got %d", count, len(all))
	}
}
