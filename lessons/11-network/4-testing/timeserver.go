// Package main реализует простую сетевую службу времени по протоколу TCP.
// Предназначена для демонстрации тестирования сетевых компонентов с помощью net.Pipe().
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

// listenAddr — адрес прослушивания сетевой службы.
const listenAddr = "0.0.0.0:12345"

// handleConn обрабатывает соединение одного клиента.
// Принимает интерфейс net.Conn, благодаря чему легко тестируется
// как с реальными сокетами, так и с виртуальной трубой net.Pipe().
func handleConn(conn net.Conn) {
	defer conn.Close()

	// Чтение строкового запроса от клиента
	reader := bufio.NewReader(conn)
	b, err := reader.ReadBytes('\n')
	if err != nil {
		return
	}

	// Удаление символов конца строки (\n и \r)
	msg := strings.TrimSpace(string(b))

	// Если клиент прислал команду "time" — отправляем текущее время
	if msg == "time" {
		currentTime := time.Now().Format(time.RFC3339) + "\n"
		_, _ = conn.Write([]byte(currentTime))
	} else {
		_, _ = conn.Write([]byte("Неизвестная команда. Доступно: time\n"))
	}
}

func main() {
	listener, err := net.Listen("tcp4", listenAddr)
	if err != nil {
		log.Fatalf("Ошибка запуска timeserver на %s: %v", listenAddr, err)
	}
	defer listener.Close()

	fmt.Printf("Сервер времени слушает на %s\n", listenAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Ошибка Accept: %v", err)
			continue
		}
		// Запуск обработчика в отдельной горутине
		go handleConn(conn)
	}
}

