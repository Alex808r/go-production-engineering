package main

import (
	"log"
	"net"
	"strings"
	"time"
)

// Сетевой адрес.
// Служба будет слушать входящие запросы на всех сетевых интерфейсах (0.0.0.0) на порту 12345.
// Например, подключиться можно по 127.0.0.1:12345 или localhost:12345.
const addr = "0.0.0.0:12345"

// Протокол сетевой службы: IPv4 TCP.
const proto = "tcp4"

func main() {
	// 1. net.Listen открывает сокет и начинает слушать входящие TCP-подключения.
	listener, err := net.Listen(proto, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	// 2. Бесконечный цикл обработки: после закрытия одного клиентского соединения
	// сервер должен оставаться активным и ожидать следующих клиентов.
	for {
		// listener.Accept() блокирует выполнение до тех пор, пока не подключится новый клиент.
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		// В реальных серверах обработчик запускают в отдельной горутине: "go handleConn(conn)",
		// чтобы медленный клиент не блокировал остальных.
		handleConn(conn)
	}
}

// handleConn читает запрос клиента посимвольно/побайтово и отправляет ответ.
// net.Conn реализует стандартные интерфейсы io.Reader, io.Writer и io.Closer.
func handleConn(conn net.Conn) {
	// Буфер для накопления прочитанного сообщения.
	var req []byte

	// Искусственно маленький буфер на 2 байта: демонстрирует потоковую природу ввода-вывода (streaming I/O).
	// Метод Read читает данные порциями по мере их поступления из сети.
	var buf = make([]byte, 2)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Println(err)
			return
		}

		end := false
		// Важно: итерируемся только по фактически прочитанным n байтам (buf[:n]),
		// чтобы не захватить устаревшие байты от предыдущей итерации!
		for _, b := range buf[:n] {
			if b == '\n' {
				end = true
				break
			}
			req = append(req, b)
		}
		if end {
			break
		}
	}

	// Удаление символов конца строки (\n и \r для совместимости с Windows CRLF).
	msg := strings.TrimSuffix(string(req), "\n")
	msg = strings.TrimSuffix(msg, "\r")

	// Простой текстовый протокол: если клиент прислал "time", отвечаем текущим временем.
	if msg == "time" {
		n, err := conn.Write([]byte(time.Now().String() + "\n"))
		if err != nil {
			log.Printf("ошибка отправки: %v", err)
			return
		}
		log.Printf("клиенту отправлено %d байт", n)
	}

	// Обязательно закрываем соединение после обслуживания запроса, чтобы освободить сокет/дескриптор.
	conn.Close()
}
