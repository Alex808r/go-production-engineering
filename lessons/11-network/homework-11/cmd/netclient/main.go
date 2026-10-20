package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	addr := flag.String("addr", "localhost:8000", "Сетевой адрес службы netsrv")
	flag.Parse()

	fmt.Printf("Подключение к поисковой службе GoSearch на %s...\n", *addr)
	conn, err := net.Dial("tcp", *addr)
	if err != nil {
		log.Fatalf("Ошибка подключения: %v\nУбедитесь, что сервер netsrv запущен на %s", err, *addr)
	}
	defer conn.Close()

	fmt.Println("Соединение успешно установлено! Для завершения работы введите 'exit'.")

	// Горутина асинхронно выводит поступающие от сервера ответы в консоль
	done := make(chan struct{})
	go func() {
		defer close(done)
		_, _ = io.Copy(os.Stdout, conn)
	}()

	// Основной цикл интерактивного диалога (запрос -> отправка -> ожидание)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		_, err := fmt.Fprintf(conn, "%s\n", text)
		if err != nil {
			fmt.Printf("Ошибка передачи данных на сервер: %v\n", err)
			break
		}

		trimmed := strings.TrimSpace(text)
		if strings.EqualFold(trimmed, "exit") || strings.EqualFold(trimmed, "quit") {
			break
		}
	}

	<-done
	fmt.Println("\nСеанс завершён. Соединение закрыто.")
}
