package storage

import (
	"errors"
	"sync"
	"testing"
)

func TestStorage_CRUD(t *testing.T) {
	db := New()

	doc1 := Document{
		Title: "Go Concurrency",
		URL:   "https://go.dev/concurrency",
		Body:  "Goroutines and channels in Go",
	}

	created := db.Add(doc1)
	if created.ID <= 0 {
		t.Fatalf("Add() didn't assign ID: %d", created.ID)
	}

	// Read
	got, err := db.Document(created.ID)
	if err != nil {
		t.Fatalf("Document(%d) unexpected error: %v", created.ID, err)
	}
	if got.Title != doc1.Title {
		t.Errorf("got title %s, want %s", got.Title, doc1.Title)
	}

	// Update
	created.Title = "Go Concurrency Advanced"
	err = db.Update(created)
	if err != nil {
		t.Fatalf("Update() failed: %v", err)
	}
	updated, _ := db.Document(created.ID)
	if updated.Title != "Go Concurrency Advanced" {
		t.Errorf("got %s, want Go Concurrency Advanced", updated.Title)
	}

	// Update non-existing
	err = db.Update(Document{ID: 9999, Title: "Does not exist"})
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("Update(9999) error = %v, want ErrNotFound", err)
	}

	// Search
	results := db.Search("concurrency")
	if len(results) == 0 {
		t.Error("Search('concurrency') returned 0 results, want >= 1")
	}

	// Delete
	err = db.Delete(created.ID)
	if err != nil {
		t.Fatalf("Delete(%d) failed: %v", created.ID, err)
	}

	// Delete again should return ErrNotFound
	err = db.Delete(created.ID)
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("Delete(%d) second time error = %v, want ErrNotFound", created.ID, err)
	}

	// Read deleted
	_, err = db.Document(created.ID)
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("Document(%d) after delete error = %v, want ErrNotFound", created.ID, err)
	}
}

func TestStorage_DocumentsAndSearch(t *testing.T) {
	db := New()

	db.Add(Document{ID: 10, Title: "Go Routines", URL: "https://go.dev/routines", Body: "Channel basics"})
	db.Add(Document{ID: 5, Title: "Effective Go", URL: "https://go.dev/effective", Body: "Best practices"})

	all := db.Documents()
	if len(all) != 2 {
		t.Fatalf("len(all) = %d, want 2", len(all))
	}
	if all[0].ID != 5 || all[1].ID != 10 {
		t.Errorf("all is not sorted: %v", all)
	}

	// Search with empty query
	if res := db.Search(""); len(res) != 0 {
		t.Errorf("Search('') = %v, want nil/empty", res)
	}

	// Search with non-existent word
	if res := db.Search("python"); len(res) != 0 {
		t.Errorf("Search('python') = %v, want empty", res)
	}
}

func TestStorage_Concurrency(t *testing.T) {
	db := New()
	var wg sync.WaitGroup

	for i := 1; i <= 50; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			doc := db.Add(Document{
				Title: "Test Concurrent Doc",
				URL:   "https://go.dev",
				Body:  "Concurrent testing in Go",
			})
			_ = db.Search("concurrent")
			_ = db.Documents()
			doc.Title = "Updated title"
			_ = db.Update(doc)
			_ = db.Delete(doc.ID)
		}(i)
	}

	wg.Wait()
}

// BenchmarkStorage_Search измеряет скорость поиска по инвертированному индексу в storage.DB.
func BenchmarkStorage_Search(b *testing.B) {
	db := New()
	for i := 1; i <= 100; i++ {
		db.Add(Document{
			ID:    i,
			Title: "Concurrency and channels in Go programming",
			URL:   "https://go.dev/doc",
			Body:  "Detailed guide on goroutines and mutexes",
		})
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = db.Search("concurrency")
	}
}

// BenchmarkStorage_Add измеряет скорость добавления и индексации документов.
func BenchmarkStorage_Add(b *testing.B) {
	doc := Document{
		Title: "Concurrency and channels in Go programming",
		URL:   "https://go.dev/doc",
		Body:  "Detailed guide on goroutines and mutexes",
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		db := New()
		db.Add(doc)
	}
}

