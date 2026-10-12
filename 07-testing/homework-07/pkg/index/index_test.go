package index

import (
	"reflect"
	"testing"

	"go-core-4/07-testing/homework-07/pkg/crawler"
)

func TestService_AddAndSearch(t *testing.T) {
	idx := New()

	docs := []crawler.Document{
		{ID: 0, Title: "The Go Programming Language", URL: "https://golang.org/doc"},
		{ID: 1, Title: "Writing Web Applications in Go", URL: "https://golang.org/doc/articles/wiki"},
		{ID: 2, Title: "Python for Data Science and Machine Learning", URL: "https://python.org"},
	}

	idx.Add(docs)

	tests := []struct {
		name     string
		word     string
		expected []int
	}{
		{name: "поиск слова 'go' в нескольких документах", word: "go", expected: []int{0, 1}},
		{name: "поиск слова 'programming' (регистронезависимый)", word: "Programming", expected: []int{0}},
		{name: "поиск токена из URL", word: "wiki", expected: []int{1}},
		{name: "поиск слова со знаками препинания", word: "science!", expected: []int{2}},
		{name: "поиск отсутствующего слова", word: "rust", expected: nil},
		{name: "поиск пустой строки", word: "", expected: nil},
		{name: "поиск только знаков препинания", word: "---", expected: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := idx.Search(tt.word)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("idx.Search(%q) = %v, want %v", tt.word, got, tt.expected)
			}
		})
	}
}

func TestBinarySearch_Table(t *testing.T) {
	docs := []crawler.Document{
		{ID: 10, Title: "Doc 10"},
		{ID: 20, Title: "Doc 20"},
		{ID: 30, Title: "Doc 30"},
		{ID: 40, Title: "Doc 40"},
		{ID: 50, Title: "Doc 50"},
	}

	tests := []struct {
		name     string
		docs     []crawler.Document
		targetID int
		wantPos  int
	}{
		{
			name:     "элемент в начале среза",
			docs:     docs,
			targetID: 10,
			wantPos:  0,
		},
		{
			name:     "элемент в середине среза",
			docs:     docs,
			targetID: 30,
			wantPos:  2,
		},
		{
			name:     "элемент в конце среза",
			docs:     docs,
			targetID: 50,
			wantPos:  4,
		},
		{
			name:     "элемент отсутствует (меньше минимального)",
			docs:     docs,
			targetID: 5,
			wantPos:  -1,
		},
		{
			name:     "элемент отсутствует (больше максимального)",
			docs:     docs,
			targetID: 100,
			wantPos:  -1,
		},
		{
			name:     "элемент отсутствует (между элементами)",
			docs:     docs,
			targetID: 25,
			wantPos:  -1,
		},
		{
			name:     "пустой срез документов",
			docs:     []crawler.Document{},
			targetID: 10,
			wantPos:  -1,
		},
		{
			name:     "срез из одного элемента (найден)",
			docs:     []crawler.Document{{ID: 42, Title: "Answer"}},
			targetID: 42,
			wantPos:  0,
		},
		{
			name:     "срез из одного элемента (не найден)",
			docs:     []crawler.Document{{ID: 42, Title: "Answer"}},
			targetID: 99,
			wantPos:  -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BinarySearch(tt.docs, tt.targetID)
			if got != tt.wantPos {
				t.Errorf("BinarySearch(..., %d) = %d, want %d", tt.targetID, got, tt.wantPos)
			}
		})
	}
}
