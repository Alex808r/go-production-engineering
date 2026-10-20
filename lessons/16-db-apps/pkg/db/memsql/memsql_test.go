package memsql_test

import (
	"context"
	"sync"
	"testing"

	"go-core-4/lessons/16-db-apps/pkg/db"
	"go-core-4/lessons/16-db-apps/pkg/db/memsql"
)

func TestInterfaceCompliance(t *testing.T) {
	var _ db.Interface = (*memsql.DB)(nil)
}

func TestMemsql_AddAndGetBooks(t *testing.T) {
	ctx := context.Background()
	store := memsql.New()

	books, err := store.Books(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(books) != 0 {
		t.Fatalf("expected 0 books, got %d", len(books))
	}

	testData := []db.Book{
		{Title: "The Hobbit", Year: 1937, Public: false, PublisherID: 1, Publisher: "Allen & Unwin"},
		{Title: "1984", Year: 1949, Public: true, PublisherID: 2, Publisher: "Secker & Warburg"},
	}

	err = store.AddBooks(ctx, testData)
	if err != nil {
		t.Fatalf("AddBooks failed: %v", err)
	}

	books, err = store.Books(ctx)
	if err != nil {
		t.Fatalf("Books failed: %v", err)
	}
	if len(books) != 2 {
		t.Fatalf("expected 2 books, got %d", len(books))
	}

	if books[0].Title != "The Hobbit" || books[1].Title != "1984" {
		t.Errorf("unexpected books content: %+v", books)
	}
}

func TestMemsql_Concurrency(t *testing.T) {
	ctx := context.Background()
	store := memsql.New()

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			_ = store.AddBooks(ctx, []db.Book{{Title: "Book", Year: 2000 + idx}})
		}(i)
	}

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = store.Books(ctx)
		}()
	}

	wg.Wait()
}

func BenchmarkMemsql_AddBooks(b *testing.B) {
	ctx := context.Background()
	store := memsql.New()
	item := []db.Book{{Title: "Benchmark Book", Year: 2024}}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = store.AddBooks(ctx, item)
	}
}

func BenchmarkMemsql_Books(b *testing.B) {
	ctx := context.Background()
	store := memsql.New()
	_ = store.AddBooks(ctx, []db.Book{
		{Title: "Book 1", Year: 2001},
		{Title: "Book 2", Year: 2002},
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = store.Books(ctx)
	}
}
