// Package main демонстрирует работу клиента standard Go RPC (пакет net/rpc).
package main

import (
	"fmt"
	"log"
	"net/rpc"

	"go-core-4/lessons/14-RPC/1-Go-RPC/pkg/books"
)

// book2 — альтернативная локальная структура на стороне клиента.
// Так как net/rpc использует encoding/gob, имена и типы полей должны совпадать,
// но сам тип структуры может быть другим.
type book2 struct {
	ID     int
	Title  string
	Author string
}

func main() {
	// 1. Установка TCP-соединения с RPC-сервером
	client, err := rpc.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("Не удалось подключиться к RPC-серверу (убедитесь, что запущен rpc-server): %v", err)
	}
	defer client.Close()

	// 2. Вызов удаленного метода "Server.Books" с типом результата из пакета "books"
	// client.Call("Сервис.Метод", запрос, &ответ)
	var req = &books.Request{}
	var data []books.Book
	err = client.Call("Server.Books", req, &data)
	if err != nil {
		log.Fatalf("Ошибка вызова Server.Books: %v", err)
	}
	fmt.Printf("1. Список книг (тип books.Book): %+v\n\n", data)

	// 3. Вызов удаленного метода "Server.Book" с локальным типом данных book2
	req = &books.Request{ID: 1}
	var item book2
	err = client.Call("Server.Book", req, &item)
	if err != nil {
		log.Fatalf("Ошибка вызова Server.Book: %v", err)
	}
	fmt.Printf("2. Книга по ID=1 (тип book2): %+v\n\n", item)

	// 4. Демонстрация обработки ошибок RPC:
	// а) Несуществующее имя метода
	fmt.Println("3. Проверка обработки ошибок:")
	err = client.Call("Server.WrongName", req, &item)
	if err != nil {
		fmt.Printf("   Ожидаемая ошибка (несуществующий метод): %v\n", err)
	}

	// б) Неверные типы аргументов (строка вместо указателя на структуру)
	err = client.Call("Server.Book", 100, "ABC")
	if err != nil {
		fmt.Printf("   Ожидаемая ошибка (несовпадение типов): %v\n", err)
	}
}

