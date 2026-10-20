// Package main_test демонстрирует тестирование сетевого кода в Go с помощью net.Pipe().
package main

import (
	"bufio"
	"net"
	"strings"
	"sync"
	"testing"
)

// Test_handleConn тестирует обработчик сетевого соединения без поднятия реального TCP-порта.
//
// ТЕОРИЯ: ВСТРОЕННАЯ ФУНКЦИЯ net.Pipe()
// net.Pipe() возвращает пару связанных сетевых соединений (srv, cl), реализующих интерфейс net.Conn:
//  - Соединение полностью синхронное и работает в оперативной памяти.
//  - Запись в один конец блокируется до тех пор, пока второй конец не прочитает данные.
//  - Идеально подходит для модульного тестирования: тесты работают мгновенно, не требуют
//    прав на открытие портов, никогда не конфликтуют с занятыми портами ("address already in use")!
func Test_handleConn(t *testing.T) {
	t.Run("успешный запрос времени (time)", func(t *testing.T) {
		srv, cl := net.Pipe()
		defer cl.Close()

		var wg sync.WaitGroup
		wg.Add(1)

		// ВАЖНО: Обработчик сервера ОБЯЗАТЕЛЬНО запускается в отдельной горутине,
		// так как net.Pipe() синхронен, и одновременная запись/чтение в одном потоке
		// привела бы к взаимной блокировке (deadlock)!
		go func() {
			defer wg.Done()
			handleConn(srv)
		}()

		// Клиент пишет команду "time\n"
		_, err := cl.Write([]byte("time\n"))
		if err != nil {
			t.Fatalf("ошибка записи в трубу: %v", err)
		}

		// Клиент читает ответ
		reader := bufio.NewReader(cl)
		response, err := reader.ReadString('\n')
		if err != nil {
			t.Fatalf("ошибка чтения ответа: %v", err)
		}

		// Проверяем, что получен непустой ответ (дата в формате RFC3339)
		if len(strings.TrimSpace(response)) < 10 {
			t.Errorf("слишком короткий ответ времени: %q", response)
		}

		wg.Wait()
	})

	t.Run("неизвестная команда", func(t *testing.T) {
		srv, cl := net.Pipe()
		defer cl.Close()

		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer wg.Done()
			handleConn(srv)
		}()

		_, err := cl.Write([]byte("unknown\n"))
		if err != nil {
			t.Fatalf("ошибка записи: %v", err)
		}

		reader := bufio.NewReader(cl)
		response, err := reader.ReadString('\n')
		if err != nil {
			t.Fatalf("ошибка чтения: %v", err)
		}

		if !strings.Contains(response, "Неизвестная команда") {
			t.Errorf("ожидалось сообщение о неизвестной команде, получено: %q", response)
		}

		wg.Wait()
	})
}

