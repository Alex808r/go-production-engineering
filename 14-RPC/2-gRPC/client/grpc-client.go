// Package main демонстрирует работу клиента gRPC в Go.
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "go-core-4/14-RPC/2-gRPC/books_proto"
)

func main() {
	// 1. Установка соединения с gRPC-сервером.
	// gRPC работает поверх HTTP/2. Для локальной разработки без TLS используем insecure.NewCredentials().
	conn, err := grpc.Dial("localhost:12345",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Не удалось подключиться к gRPC-серверу (проверьте, запущен ли grpc-server): %v", err)
	}
	defer conn.Close()

	// 2. Создаем клиентский стаб (Client Stub) из сгенерированного пакета pb
	client := pb.NewBookinistClient(conn)

	// Контекст для контроля длительности вызовов
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 3. Чтение потока книг от сервера
	err = printAllBooksOnServer(ctx, client)
	if err != nil {
		log.Fatalf("Ошибка при чтении книг: %v", err)
	}

	// 4. Унарный вызов: добавление новой книги на сервер
	fmt.Println("\nДобавляем новую книгу через Unary RPC (AddBook)...")
	_, err = client.AddBook(ctx, &pb.Book{
		Id:         3,
		Title:      "The Lord Of The Rings",
		AuthorName: "J.R.R. Tolkien",
	})
	if err != nil {
		log.Fatalf("Ошибка AddBook: %v", err)
	}

	// 5. Повторное чтение потока для проверки добавления
	err = printAllBooksOnServer(ctx, client)
	if err != nil {
		log.Fatalf("Ошибка при чтении книг: %v", err)
	}
}

// printAllBooksOnServer запрашивает поток книг через Server-side Streaming RPC.
func printAllBooksOnServer(ctx context.Context, client pb.BookinistClient) error {
	fmt.Println("\nЗапрашиваю книги с gRPC-сервера (Server Streaming)...")

	// Вызов метода Books возвращает клиентский стрим pb.Bookinist_BooksClient
	stream, err := client.Books(ctx, &pb.Empty{})
	if err != nil {
		return err
	}

	// Цикл получения сообщений из потока
	for {
		select {
		case <-ctx.Done():
			// Контекст отменен или истек таймаут
			return ctx.Err()
		default:
			// stream.Recv() блокируется до получения очередного сообщения
			book, err := stream.Recv()
			if err == io.EOF {
				// Сервер завершил отправку потока (EOF)
				fmt.Println("Поток успешно завершен (получен EOF).")
				return nil
			}
			if err != nil {
				return fmt.Errorf("ошибка чтения из стрима: %w", err)
			}
			fmt.Printf("  📚 Получена книга: [ID: %d] %q (Автор: %s)\n",
				book.GetId(), book.GetTitle(), book.GetAuthorName())
		}
	}
}

