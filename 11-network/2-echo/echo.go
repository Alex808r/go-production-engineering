// Package main реализует конкурентный многопоточный TCP эхо-сервер на Go.
//
// ТЕОРИЯ: МОДЕЛЬ ОБРАБОТКИ СЕТЕВЫХ СОЕДИНЕНИЙ В GO
// В языках вроде C или Java для масштабирования сетевых серверов часто используют
// сложные асинхронные неблокирующие циклы событий (Epoll / Event Loop / Reactor).
// В Go сетевой планировщик рантайма (Network Poller) прозрачно интегрирован с GMP:
// мы пишем простой синхронный блокирующий код в горутине, а рантайм Go под капотом
// использует epoll (Linux) или kqueue (macOS) без каких-либо накладных расходов!
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

const echoAddr = ":12345"

// handler обрабатывает сессию одного подключившегося клиента.
// Запускается в отдельной изолированной горутине.
func handler(conn net.Conn) {
	// Гарантируем закрытие сетевого сокета при завершении сессии клиента
	defer conn.Close()

	remoteAddr := conn.RemoteAddr().String()
	fmt.Printf("[%s] Клиент подключился\n", remoteAddr)
	defer fmt.Printf("[%s] Соединение закрыто\n", remoteAddr)

	// bufio.Reader буферизует чтение из сетевого сокета, снижая количество системных вызовов
	r := bufio.NewReader(conn)

	for {
		// ReadBytes читает данные до символа разделителя строки ('\n')
		line, err := r.ReadBytes('\n')
		if err != nil {
			// Ошибка чтения (клиент отключился, закрыл сокет или нажал Ctrl+C)
			return
		}

		fmt.Printf("[%s] Получено эхо-сообщение: %s", remoteAddr, string(line))

		// Отправляем эхо обратно клиенту
		_, err = conn.Write(line)
		if err != nil {
			log.Printf("[%s] Ошибка отправки данных клиенту: %v", remoteAddr, err)
			return
		}
	}
}

func main() {
	// 1. Открываем слушающий сокет
	listener, err := net.Listen("tcp4", echoAddr)
	if err != nil {
		log.Fatalf("Не удалось запустить эхо-сервер на %s: %v", echoAddr, err)
	}
	defer listener.Close()

	fmt.Printf("TCP Эхо-сервер успешно запущен на порту %s\n", echoAddr)
	fmt.Println("Для подключения используйте: telnet localhost 12345 или nc localhost 12345")

	// 2. Бесконечный цикл приема входящих клиентских соединений
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Ошибка при вызове Accept: %v", err)
			continue
		}

		// Запускаем горутину на каждое клиентское подключение:
		// сервер может одновременно обслуживать десятки и сотни тысяч клиентов!
		go handler(conn)
	}
}

