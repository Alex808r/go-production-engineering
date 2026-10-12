package index

import (
	"reflect"
	"testing"

	"go-core-4/06-oop/homework-06/pkg/crawler"
)

func TestService_AddAndSearch(t *testing.T) {
	idx := New()

	docs := []crawler.Document{
		{ID: 0, Title: "The Go Programming Language"},
		{ID: 1, Title: "Writing Web Applications in Go"},
		{ID: 2, Title: "Python for Data Science"},
	}

	idx.Add(docs)

	tests := []struct {
		word     string
		expected []int
	}{
		{word: "go", expected: []int{0, 1}},
		{word: "programming", expected: []int{0}},
		{word: "applications", expected: []int{1}},
		{word: "science", expected: []int{2}},
		{word: "rust", expected: nil},
	}

	for _, tt := range tests {
		t.Run(tt.word, func(t *testing.T) {
			got := idx.Search(tt.word)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("idx.Search(%q) = %v, want %v", tt.word, got, tt.expected)
			}
		})
	}
}

func TestBinarySearch(t *testing.T) {
	docs := []crawler.Document{
		{ID: 10, Title: "Doc 10"},
		{ID: 20, Title: "Doc 20"},
		{ID: 30, Title: "Doc 30"},
		{ID: 40, Title: "Doc 40"},
	}

	if pos := BinarySearch(docs, 20); pos != 1 {
		t.Errorf("BinarySearch(20) = %d, want 1", pos)
	}

	if pos := BinarySearch(docs, 50); pos != -1 {
		t.Errorf("BinarySearch(50) = %d, want -1", pos)
	}
}
