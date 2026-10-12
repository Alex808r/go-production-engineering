package main

import (
	"reflect"
	"testing"

	"go-core-4/02-syntax/homework-02/pkg/crawler"
)

func TestFilterDocs(t *testing.T) {
	docs := []crawler.Document{
		{URL: "https://go.dev/doc/", Title: "Go Documentation and Guides"},
		{URL: "https://go.dev/blog/", Title: "The Go Blog articles"},
		{URL: "https://golang.org/pkg/", Title: "Standard library packages documentation"},
		{URL: "https://example.com", Title: "Random page"},
	}

	tests := []struct {
		name      string
		query     string
		wantCount int
		wantURLs  []string
	}{
		{
			name:      "match in title (case insensitive)",
			query:     "documentation",
			wantCount: 2,
			wantURLs:  []string{"https://go.dev/doc/", "https://golang.org/pkg/"},
		},
		{
			name:      "match in URL",
			query:     "blog",
			wantCount: 1,
			wantURLs:  []string{"https://go.dev/blog/"},
		},
		{
			name:      "match multiple",
			query:     "go",
			wantCount: 3,
			wantURLs:  []string{"https://go.dev/doc/", "https://go.dev/blog/", "https://golang.org/pkg/"},
		},
		{
			name:      "no match",
			query:     "python",
			wantCount: 0,
			wantURLs:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FilterDocs(docs, tt.query)
			if len(got) != tt.wantCount {
				t.Fatalf("FilterDocs() returned %d items, want %d", len(got), tt.wantCount)
			}
			var gotURLs []string
			for _, d := range got {
				gotURLs = append(gotURLs, d.URL)
			}
			if !reflect.DeepEqual(gotURLs, tt.wantURLs) {
				t.Errorf("FilterDocs() URLs = %v, want %v", gotURLs, tt.wantURLs)
			}
		})
	}
}
