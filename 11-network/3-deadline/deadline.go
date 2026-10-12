// Package main демонстрирует управление сетевыми таймаутами (Deadlines) в Go
// для защиты сервера от зависших соединений, утечек файловых дескрипторов и Slowloris атак.
package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

const timeoutAddr = ":12345"
const idleTimeout = 5 * time.Second

// handler обрабатывает подключение с контролем времени бездействия (Idle Timeout).
func handler(conn net.Conn) {
	defer conn.Close()

	remoteAddr := conn.RemoteAddr().String()
	fmt.Printf("[%s] Соединение установлено. Таймаут бездействия: %v\n", remoteAddr, idleTimeout)
	defer fmt.Printf("[%s] Соединение закрыто\n", remoteAddr)

	// ========================================================================
	// ТЕОРИЯ: ДЕДЛАЙНЫ В NET.CONN
	// ========================================================================
	// Интерфейс net.Conn предоставляет три метода для установки абсолютного времени завершения:
	//   1. conn.SetDeadline(t time.Time) — таймаут на чтение и запись.
	//   2. conn.SetReadDeadline(t time.Time) — таймаут только на чтение.
	//   3. conn.SetWriteDeadline(t time.Time) — таймаут только на запись.
	//
	// Если передать нулевое время (time.Time{}), дедлайн отключается.

	r := bufio.NewReader(conn)

	for {
		// Устанавливаем абсолютное время дедлайна на 5 секунд вперед от ТЕКУЩЕГО момента.
		// Если за 5 секунд клиент не пришлет данных, вызов ReadBytes вернет ошибку таймаута!
		err := conn.SetReadDeadline(time.Now().Add(idleTimeout))
		if err != nil {
			log.Printf("[%s] Ошибка установки дедлайна: %v", remoteAddr, err)
			return
		}

		line, err := r.ReadBytes('\n')
		if err != nil {
			// Проверяем, произошел ли разрыв соединения по таймауту
			var netErr net.Error
			if errors.As(err, &netErr) && netErr.Timeout() {
				fmt.Printf("[%s] Соединение разорвано: истек таймаут бездействия (%v)!\n", remoteAddr, idleTimeout)
				return
			}
			if errors.Is(err, os.ErrDeadlineExceeded) {
				fmt.Printf("[%s] Таймаут дедлайна истек!\n", remoteAddr)
				return
			}

			// Обычное отключение клиента (клиент сам закрыл соединение)
			return
		}

		fmt.Printf("[%s] Принято: %s", remoteAddr, string(line))

		// Отправляем ответ клиенту
		if _, err = conn.Write(line); err != nil {
			return
		}
	}
}

func main() {
	listener, err := net.Listen("tcp4", timeoutAddr)
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
	defer listener.Close()

	fmt.Printf("Сервер с таймаутами успешно запущен на порту %s\n", timeoutAddr)
	fmt.Printf("Если после подключения не отправлять данные %v, соединение закроется автоматически.\n", idleTimeout)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Ошибка Accept: %v", err)
			continue
		}
		go handler(conn)
	}
}

