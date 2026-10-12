// Package main реализует сервер удалённого вызова процедур (RPC)
// с использованием стандартного пакета net/rpc языка Go.
package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"sync"

	"go-core-4/14-RPC/1-Go-RPC/pkg/books"
)

// Server — экспортируемый тип данных RPC-сервера.
//
// Требования net/rpc к методам удалённого вызова:
// 1. Метод должен принадлежать экспортируемому типу.
// 2. Метод должен быть экспортируемым (начинаться с заглавной буквы).
// 3. Метод имеет ровно два аргумента: (req T1, resp *T2), оба типа должны быть экспортируемыми.
// 4. Второй аргумент (resp) ОБЯЗАТЕЛЬНО должен быть указателем.
// 5. Метод должен возвращать ровно одно значение типа error:
//    func (t *T) MethodName(req T1, resp *T2) error
type Server int

// Мьютекс для безопасного одновременного чтения/записи списка книг
var booksMu sync.RWMutex

// Books возвращает список всех книг клиенту.
func (s *Server) Books(req books.Request, resp *[]books.Book) error {
	booksMu.RLock()
	defer booksMu.RUnlock()

	// ВАЖНО: Почему присваивание resp = &allBooks НЕ РАБОТАЕТ?
	// В Go все аргументы передаются ПО ЗНАЧЕНИЮ (by value).
	// Параметр 'resp' является локальной копией указателя. Если мы переприсвоим
	// локальную переменную resp = &allBooks, мы лишь изменим адрес, на который
	// указывает локальная переменная внутри этой функции. Вызывающая сторона (net/rpc)
	// об этом изменении не узнает и вернет клиенту пустой результат!
	//
	// Единственный правильный способ вернуть данные — разыменовать указатель:
	// *resp = allBooks (запись данных по переданному адресу памяти).
	res := make([]books.Book, len(allBooks))
	copy(res, allBooks)
	*resp = res
	return nil
}

// Book ищет и возвращает книгу по её идентификатору ID.
func (s *Server) Book(req books.Request, resp *books.Book) error {
	booksMu.RLock()
	defer booksMu.RUnlock()

	for _, b := range allBooks {
		if b.ID == req.ID {
			*resp = b
			return nil
		}
	}
	return fmt.Errorf("книга с ID %d не найдена", req.ID)
}

func main() {
	srv := new(Server)

	// rpc.Register публикует методы объекта srv для удалённого вызова.
	// По умолчанию имя сервиса совпадает с именем типа: "Server".
	err := rpc.Register(srv)
	if err != nil {
		log.Fatalf("Ошибка регистрации RPC-сервиса: %v", err)
	}

	// Открываем TCP-слушатель на порту 8080
	listener, err := net.Listen("tcp4", ":8080")
	if err != nil {
		log.Fatalf("Ошибка запуска TCP-слушателя: %v", err)
	}
	defer listener.Close()

	fmt.Println("RPC-сервер (net/rpc) запущен на порту :8080")
	fmt.Println("Ожидание клиентских вызовов...")

	// Цикл обработки входящих подключений:
	// для каждого TCP-клиента запускаем go rpc.ServeConn(conn),
	// что обеспечивает конкурентное обслуживание сотен клиентов параллельно.
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Ошибка Accept: %v", err)
			continue
		}

		go rpc.ServeConn(conn)
	}
}

var allBooks = []books.Book{
	{
		ID:     1,
		Title:  "The Lord Of The Rings",
		Author: "J.R.R. Tolkien",
	},
	{
		ID:     2,
		Title:  "The Chronicles of Amber",
		Author: "Roger Zelazny",
	},
}

