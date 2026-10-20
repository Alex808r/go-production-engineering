package dip_test

import (
	"testing"

	"go-core-4/lessons/17-system-design/SOLID/5-DIP"
)

type mockStorage struct {
	items []dip.Book
}

func (m *mockStorage) Books() []dip.Book {
	return m.items
}

func TestGoodServer_DIP(t *testing.T) {
	mock := &mockStorage{
		items: []dip.Book{
			{ID: 1, Title: "The Pragmatic Programmer"},
			{ID: 2, Title: "Clean Code"},
		},
	}

	server := dip.NewGoodServer(mock)
	server.BusinessLogic()

	if count := server.BooksCount(); count != 2 {
		t.Fatalf("expected 2 books, got %d", count)
	}
}
