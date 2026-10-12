package index

import (
	"reflect"
	"sort"
	"testing"

	"go-core-4/GoSearch/pkg/crawler"
)

func TestService_AddAndSearch(t *testing.T) {
	docs := []crawler.Document{
		{
			ID:    0,
			URL:   "https://go.dev",
			Title: "The Go Programming Language",
		},
		{
			ID:    1,
			URL:   "https://golang.org/doc",
			Title: "Documentation - The Go Programming Language",
		},
		{
			ID:    2,
			URL:   "https://yandex.ru",
			Title: "Яндекс - быстрый поиск в интернете",
		},
	}

	idx := New()
	idx.Add(docs)

	tests := []struct {
		name    string
		word    string
		wantIDs []int
	}{
		{
			name:    "search 'go' (lowercase)",
			word:    "go",
			wantIDs: []int{0, 1},
		},
		{
			name:    "search 'Go' (case-insensitive)",
			word:    "Go",
			wantIDs: []int{0, 1},
		},
		{
			name:    "search 'programming'",
			word:    "programming",
			wantIDs: []int{0, 1},
		},
		{
			name:    "search 'яндекс'",
			word:    "яндекс",
			wantIDs: []int{2},
		},
		{
			name:    "search non-existent word",
			word:    "python",
			wantIDs: nil,
		},
		{
			name:    "search empty word",
			word:    "",
			wantIDs: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := idx.Search(tt.word)
			if !reflect.DeepEqual(got, tt.wantIDs) {
				t.Errorf("Search(%q) = %v, want %v", tt.word, got, tt.wantIDs)
			}
		})
	}
}

func TestService_Deduplication(t *testing.T) {
	docs := []crawler.Document{
		{
			ID:    10,
			URL:   "https://example.com",
			Title: "Go Go Go! Modern Go Development",
		},
	}

	idx := New()
	idx.Add(docs)

	got := idx.Search("go")
	want := []int{10}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected deduplicated IDs for word 'go': got %v, want %v", got, want)
	}
}

func TestBinarySearch(t *testing.T) {
	docs := []crawler.Document{
		{ID: 1, Title: "Doc 1"},
		{ID: 3, Title: "Doc 3"},
		{ID: 7, Title: "Doc 7"},
		{ID: 15, Title: "Doc 15"},
		{ID: 20, Title: "Doc 20"},
	}

	// Убеждаемся, что документы отсортированы
	sort.Slice(docs, func(i, j int) bool {
		return docs[i].ID < docs[j].ID
	})

	tests := []struct {
		name      string
		targetID  int
		wantIndex int
	}{
		{
			name:      "first element",
			targetID:  1,
			wantIndex: 0,
		},
		{
			name:      "middle element",
			targetID:  7,
			wantIndex: 2,
		},
		{
			name:      "last element",
			targetID:  20,
			wantIndex: 4,
		},
		{
			name:      "not found (smaller than min)",
			targetID:  0,
			wantIndex: -1,
		},
		{
			name:      "not found (between elements)",
			targetID:  5,
			wantIndex: -1,
		},
		{
			name:      "not found (greater than max)",
			targetID:  100,
			wantIndex: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BinarySearch(docs, tt.targetID)
			if got != tt.wantIndex {
				t.Errorf("BinarySearch(targetID=%d) = %d, want %d", tt.targetID, got, tt.wantIndex)
			}
		})
	}
}

func BenchmarkIndex_Add(b *testing.B) {
	docs := []crawler.Document{
		{ID: 1, Title: "Go Programming Language and Tools"},
		{ID: 2, Title: "Concurrency in Go: Goroutines and Channels"},
		{ID: 3, Title: "Microservices with Go and gRPC"},
		{ID: 4, Title: "Database optimization with PostgreSQL and Redis"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		idx := New()
		idx.Add(docs)
	}
}

func BenchmarkIndex_Search(b *testing.B) {
	docs := []crawler.Document{
		{ID: 1, Title: "Go Programming Language and Tools"},
		{ID: 2, Title: "Concurrency in Go: Goroutines and Channels"},
		{ID: 3, Title: "Microservices with Go and gRPC"},
	}
	idx := New()
	idx.Add(docs)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = idx.Search("Go")
	}
}

func BenchmarkBinarySearch(b *testing.B) {
	docs := make([]crawler.Document, 1000)
	for i := range docs {
		docs[i] = crawler.Document{ID: i, Title: "Doc Title"}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = BinarySearch(docs, 777)
	}
}

