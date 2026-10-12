package books_test

import (
	"fmt"
	"net"
	"net/rpc"
	"testing"

	"go-core-4/14-RPC/1-Go-RPC/pkg/books"
)

// mockServer реализует RPC-сервер для тестирования контрактов пакета books.
type mockServer struct {
	items []books.Book
}

func (s *mockServer) Books(req books.Request, resp *[]books.Book) error {
	*resp = s.items
	return nil
}

func (s *mockServer) Book(req books.Request, resp *books.Book) error {
	for _, b := range s.items {
		if b.ID == req.ID {
			*resp = b
			return nil
		}
	}
	return fmt.Errorf("книга с ID %d не найдена", req.ID)
}

func TestBooksRPC_OverNetPipe(t *testing.T) {
	svc := &mockServer{
		items: []books.Book{
			{ID: 1, Title: "The Go Programming Language", Author: "Alan Donovan"},
			{ID: 2, Title: "Design Patterns", Author: "GoF"},
		},
	}

	server := rpc.NewServer()
	if err := server.RegisterName("BooksService", svc); err != nil {
		t.Fatalf("RegisterName error: %v", err)
	}

	serverConn, clientConn := net.Pipe()
	defer serverConn.Close()
	defer clientConn.Close()

	go server.ServeConn(serverConn)

	client := rpc.NewClient(clientConn)
	defer client.Close()

	// 1. Получение всех книг
	var all []books.Book
	if err := client.Call("BooksService.Books", books.Request{}, &all); err != nil {
		t.Fatalf("Books call failed: %v", err)
	}
	if len(all) != 2 {
		t.Errorf("len(all) = %d, want 2", len(all))
	}

	// 2. Получение книги по ID
	var single books.Book
	if err := client.Call("BooksService.Book", books.Request{ID: 1}, &single); err != nil {
		t.Fatalf("Book call failed: %v", err)
	}
	if single.Title != "The Go Programming Language" {
		t.Errorf("unexpected book title: %s", single.Title)
	}

	// 3. Несуществующая книга
	var notFound books.Book
	if err := client.Call("BooksService.Book", books.Request{ID: 999}, &notFound); err == nil {
		t.Errorf("expected error for non-existing book ID, got nil")
	}
}

func BenchmarkBooksRPC_Call(b *testing.B) {
	svc := &mockServer{
		items: []books.Book{
			{ID: 1, Title: "The Go Programming Language", Author: "Alan Donovan"},
		},
	}
	server := rpc.NewServer()
	_ = server.RegisterName("BooksService", svc)

	serverConn, clientConn := net.Pipe()
	defer serverConn.Close()
	defer clientConn.Close()

	go server.ServeConn(serverConn)

	client := rpc.NewClient(clientConn)
	defer client.Close()

	req := books.Request{ID: 1}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var single books.Book
		_ = client.Call("BooksService.Book", req, &single)
	}
}
