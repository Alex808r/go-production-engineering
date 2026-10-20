package dip

import (
	"reflect"
	"testing"
)

// mockStorage — мок-реализация интерфейса Storage для юнит-тестирования бизнес-логики.
// Не требует базы данных, сетевых подключений или файловой системы.
type mockStorage struct {
	books []Book
}

func (m *mockStorage) Books() []Book {
	return m.books
}

func TestGoodServer_BusinessLogic(t *testing.T) {
	expectedBooks := []Book{
		{ID: 1, Title: "The Go Programming Language"},
		{ID: 2, Title: "Concurrency in Go"},
	}

	// Внедряем мок в GoodServer через интерфейс Storage
	mock := &mockStorage{books: expectedBooks}
	server := NewGoodServer(mock)

	got := server.BusinessLogic()
	if !reflect.DeepEqual(got, expectedBooks) {
		t.Errorf("BusinessLogic() = %+v, want %+v", got, expectedBooks)
	}
}

func TestDB_Books(t *testing.T) {
	dbBooks := []Book{
		{ID: 10, Title: "Learning Go"},
	}
	db := NewDB(dbBooks...)

	if got := db.Books(); !reflect.DeepEqual(got, dbBooks) {
		t.Errorf("db.Books() = %+v, want %+v", got, dbBooks)
	}
}
