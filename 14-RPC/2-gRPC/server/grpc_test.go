package main

import (
	"context"
	"io"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	pb "go-core-4/14-RPC/2-gRPC/books_proto"
)

const bufSize = 1024 * 1024

func setupGRPCTestServer(t *testing.T) (pb.BookinistClient, func()) {
	lis := bufconn.Listen(bufSize)
	s := grpc.NewServer()

	srv := &Bookinist{
		Data: []*pb.Book{
			{Id: 1, Title: "1984", AuthorName: "George Orwell"},
			{Id: 2, Title: "The Go Programming Language", AuthorName: "Alan Donovan"},
		},
	}
	pb.RegisterBookinistServer(s, srv)

	go func() {
		if err := s.Serve(lis); err != nil && err != grpc.ErrServerStopped {
			t.Errorf("Serve error: %v", err)
		}
	}()

	conn, err := grpc.DialContext(context.Background(), "bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}

	client := pb.NewBookinistClient(conn)

	cleanup := func() {
		_ = conn.Close()
		s.GracefulStop()
		_ = lis.Close()
	}

	return client, cleanup
}

func TestBookinist_Books_Streaming(t *testing.T) {
	client, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	stream, err := client.Books(context.Background(), &pb.Empty{})
	if err != nil {
		t.Fatalf("Books streaming call failed: %v", err)
	}

	var count int
	for {
		book, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("stream.Recv error: %v", err)
		}
		count++
		if book.Id <= 0 || book.Title == "" {
			t.Errorf("invalid book received: %+v", book)
		}
	}

	if count != 2 {
		t.Errorf("expected 2 books, received: %d", count)
	}
}

func TestBookinist_AddBook(t *testing.T) {
	client, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	newBook := &pb.Book{
		Id:         3,
		Title:      "The Hobbit",
		AuthorName: "J.R.R. Tolkien",
	}

	_, err := client.AddBook(context.Background(), newBook)
	if err != nil {
		t.Fatalf("AddBook failed: %v", err)
	}

	// Проверяем, что книга появилась в потоке
	stream, err := client.Books(context.Background(), &pb.Empty{})
	if err != nil {
		t.Fatalf("Books call failed: %v", err)
	}

	var total int
	for {
		_, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("stream.Recv error: %v", err)
		}
		total++
	}

	if total != 3 {
		t.Errorf("expected 3 books after AddBook, got %d", total)
	}
}

func BenchmarkBookinist_AddBook(b *testing.B) {
	lis := bufconn.Listen(bufSize)
	s := grpc.NewServer()
	srv := &Bookinist{}
	pb.RegisterBookinistServer(s, srv)
	go func() { _ = s.Serve(lis) }()

	conn, _ := grpc.DialContext(context.Background(), "bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	defer func() {
		_ = conn.Close()
		s.Stop()
		_ = lis.Close()
	}()

	client := pb.NewBookinistClient(conn)
	book := &pb.Book{Id: 10, Title: "Benchmark Book", AuthorName: "Author"}
	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = client.AddBook(ctx, book)
	}
}
