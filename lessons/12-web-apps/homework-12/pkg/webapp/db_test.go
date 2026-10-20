package webapp

import (
	"fmt"
	"testing"
)

func TestSearchDB(t *testing.T) {
	db := NewSearchDB()

	docs := []Document{
		{ID: 1, URL: "https://golang.org", Title: "The Go Programming Language", Body: "Go is an open source language"},
		{ID: 2, URL: "https://pkg.go.dev", Title: "Go Packages and Documentation", Body: "Find packages and libraries"},
	}

	db.Add(docs)

	storedDocs := db.Documents()
	if len(storedDocs) != 2 {
		t.Fatalf("ожидалось 2 документа, получено %d", len(storedDocs))
	}

	idx := db.Index()
	if len(idx) == 0 {
		t.Fatalf("индекс пуст после добавления документов")
	}

	res := db.Search("programming")
	if len(res) == 0 || res[0] != 1 {
		t.Errorf("ожидался документ ID 1 для слова 'programming', получено: %v", res)
	}

	resPkg := db.Search("packages")
	if len(resPkg) == 0 || resPkg[0] != 2 {
		t.Errorf("ожидался документ ID 2 для слова 'packages', получено: %v", resPkg)
	}

	emptyRes := db.Search("nonexistent")
	if len(emptyRes) != 0 {
		t.Errorf("ожидался пустой результат для несуществующего слова, получено: %v", emptyRes)
	}

	cleanRes := db.Search("")
	if cleanRes != nil {
		t.Errorf("ожидался nil для пустой строки, получено: %v", cleanRes)
	}
}

// BenchmarkSearchDB_Search измеряет скорость поиска по инвертированному индексу в SearchDB.
func BenchmarkSearchDB_Search(b *testing.B) {
	db := NewSearchDB()
	var docs []Document
	for i := 1; i <= 100; i++ {
		docs = append(docs, Document{
			ID:    i,
			URL:   fmt.Sprintf("https://example.com/doc/%d", i),
			Title: fmt.Sprintf("Golang Documentation Article %d", i),
			Body:  "Concurrent programming with goroutines and channels in Go.",
		})
	}
	db.Add(docs)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = db.Search("goroutines")
	}
}

// BenchmarkSearchDB_Add измеряет скорость добавления и индексации документов.
func BenchmarkSearchDB_Add(b *testing.B) {
	sampleDocs := []Document{
		{ID: 1, URL: "https://golang.org", Title: "The Go Programming Language", Body: "Open source programming"},
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		db := NewSearchDB()
		db.Add(sampleDocs)
	}
}

