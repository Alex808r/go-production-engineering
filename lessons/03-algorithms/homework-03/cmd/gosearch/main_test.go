package main

import (
	"reflect"
	"sort"
	"testing"

	"go-core-4/lessons/03-algorithms/homework-03/pkg/crawler"
	"go-core-4/lessons/03-algorithms/homework-03/pkg/index"
)

func TestSearchDocs(t *testing.T) {
	docs := []crawler.Document{
		{
			ID:    0,
			URL:   "https://go.dev/doc/",
			Title: "Go Documentation and Guides",
		},
		{
			ID:    1,
			URL:   "https://go.dev/blog/",
			Title: "The Go Blog articles",
		},
		{
			ID:    2,
			URL:   "https://golang.org/pkg/",
			Title: "Standard library packages documentation",
		},
		{
			ID:    3,
			URL:   "https://example.com",
			Title: "Random web page with no matching words",
		},
	}

	// Сортировка по ID
	sort.Slice(docs, func(i, j int) bool {
		return docs[i].ID < docs[j].ID
	})

	idx := index.New()
	idx.Add(docs)

	tests := []struct {
		name      string
		query     string
		wantCount int
		wantIDs   []int
	}{
		{
			name:      "search 'documentation' (found in 2 docs)",
			query:     "documentation",
			wantCount: 2,
			wantIDs:   []int{0, 2},
		},
		{
			name:      "search 'blog' (found in 1 doc)",
			query:     "blog",
			wantCount: 1,
			wantIDs:   []int{1},
		},
		{
			name:      "search 'go' (found in doc 0 and doc 1)",
			query:     "go",
			wantCount: 2,
			wantIDs:   []int{0, 1},
		},
		{
			name:      "search word not in index",
			query:     "python",
			wantCount: 0,
			wantIDs:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := SearchDocs(docs, idx, tt.query)
			if len(results) != tt.wantCount {
				t.Fatalf("SearchDocs(%q) count = %d, want %d", tt.query, len(results), tt.wantCount)
			}
			var gotIDs []int
			for _, r := range results {
				gotIDs = append(gotIDs, r.ID)
			}
			if !reflect.DeepEqual(gotIDs, tt.wantIDs) {
				t.Errorf("SearchDocs(%q) IDs = %v, want %v", tt.query, gotIDs, tt.wantIDs)
			}
		})
	}
}
