// Package main реализует высокопроизводительный gRPC-сервер
// с поддержкой унарных запросов и серверного стриминга (Server-side Streaming).
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"

	pb "go-core-4/lessons/14-RPC/2-gRPC/books_proto"
)

// Bookinist реализует интерфейс gRPC-сервера, сгенерированный компилятором protoc.
type Bookinist struct {
	mu   sync.RWMutex
	Data []*pb.Book

	// Встраивание (композиция) нереализованного сервера из сгенерированного пакета.
	// Это требование gRPC для обеспечения прямой совместимости (forward compatibility):
	// при добавлении новых RPC-методов в .proto файл сервер продолжит компилироваться,
	// возвращая статус codes.Unimplemented для ещё не реализованных методов.
	pb.UnimplementedBookinistServer
}

// Books реализует серверный стриминг: клиент отправляет один запрос (Empty),
// а сервер отправляет книги одну за другой через stream.Send(book).
func (b *Bookinist) Books(_ *pb.Empty, stream pb.Bookinist_BooksServer) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, book := range b.Data {
		// Отправка очередного элемента в стрим
		if err := stream.Send(book); err != nil {
			return err
		}
	}
	// Возврат nil сигнализирует gRPC-фреймворку об успешном завершении потока (EOF для клиента)
	return nil
}

// AddBook реализует унарный RPC-вызов (Unary RPC): принимает книгу и сохраняет её в памяти.
func (b *Bookinist) AddBook(_ context.Context, book *pb.Book) (*pb.Empty, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.Data = append(b.Data, book)
	return new(pb.Empty), nil
}

func main() {
	srv := Bookinist{}
	srv.Data = append(srv.Data,
		&pb.Book{Id: 1, Title: "1984", AuthorName: "George Orwell"},
		&pb.Book{Id: 2, Title: "The Go Programming Language", AuthorName: "Alan Donovan"},
	)

	// Создаем сетевой сокет TCP на порту :12345
	lis, err := net.Listen("tcp", ":12345")
	if err != nil {
		log.Fatalf("Ошибка привязки порта: %v", err)
	}

	// Инициализируем gRPC сервер
	grpcServer := grpc.NewServer()

	// Регистрируем реализацию сервиса в gRPC сервере
	pb.RegisterBookinistServer(grpcServer, &srv)

	fmt.Println("gRPC-сервер запущен на порту :12345 (HTTP/2 + Protobuf)")

	// Запуск цикла обработки входящих gRPC-запросов (блокирующий вызов)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Ошибка работы gRPC-сервера: %v", err)
	}
}

