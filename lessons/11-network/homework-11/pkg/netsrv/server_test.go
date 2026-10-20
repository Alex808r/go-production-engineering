package netsrv

import (
	"bufio"
	"net"
	"strings"
	"testing"
	"time"

	"go-core-4/lessons/07-testing/homework-07/pkg/crawler"
)

// mockSearcher реализует интерфейс Searcher для тестов.
type mockSearcher struct {
	docs []crawler.Document
}

func (m *mockSearcher) Search(query string) []crawler.Document {
	var result []crawler.Document
	q := strings.ToLower(query)
	for _, doc := range m.docs {
		if strings.Contains(strings.ToLower(doc.Title), q) ||
			strings.Contains(strings.ToLower(doc.URL), q) {
			result = append(result, doc)
		}
	}
	return result
}

func TestServer_HandleConn_WithNetPipe(t *testing.T) {
	mock := &mockSearcher{
		docs: []crawler.Document{
			{ID: 1, Title: "Go Programming Language", URL: "https://go.dev"},
			{ID: 2, Title: "A Tour of Go", URL: "https://tour.golang.org"},
			{ID: 3, Title: "PostgreSQL Database", URL: "https://postgresql.org"},
		},
	}

	srv := New(":8000", mock)
	srv.SetTimeout(0) // Отключаем таймауты для in-memory pipe

	// Создаем виртуальное двустороннее соединение в памяти через net.Pipe()
	serverConn, clientConn := net.Pipe()
	defer clientConn.Close()

	// 1. Запускаем обработчик сервера в отдельной горутине
	go srv.HandleConn(serverConn)

	// 2. Согласно паттерну тестирования из лекции (слайд 22):
	// запускаем отдельный поток для асинхронного чтения из канала net.Pipe,
	// чтобы синхронные вызовы Write с обеих сторон не блокировали друг друга.
	reader := bufio.NewReader(clientConn)
	linesCh := make(chan string, 100)

	go func() {
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				close(linesCh)
				return
			}
			linesCh <- line
		}
	}()

	waitForResponse := func(expected string) bool {
		timer := time.After(2 * time.Second)
		for {
			select {
			case line, ok := <-linesCh:
				if !ok {
					return false
				}
				if strings.Contains(line, expected) {
					return true
				}
			case <-timer:
				return false
			}
		}
	}

	// Проверяем приветственное сообщение
	if !waitForResponse("GoSearch") {
		t.Fatal("не дождались приветственного баннера от сервера")
	}

	// Отправляем поисковый запрос "Go"
	if _, err := clientConn.Write([]byte("Go\n")); err != nil {
		t.Fatalf("ошибка отправки запроса: %v", err)
	}
	if !waitForResponse("Найдено документов: 2") {
		t.Fatal("не найдены результаты поиска по запросу 'Go'")
	}

	// Отправляем пустой запрос
	if _, err := clientConn.Write([]byte("\n")); err != nil {
		t.Fatalf("ошибка отправки пустого запроса: %v", err)
	}

	// Отправляем запрос, по которому ничего нет
	if _, err := clientConn.Write([]byte("python\n")); err != nil {
		t.Fatalf("ошибка отправки: %v", err)
	}
	if !waitForResponse("ничего не найдено") {
		t.Fatal("не получено сообщение о пустом результате поиска")
	}

	// Отправляем команду выхода
	if _, err := clientConn.Write([]byte("exit\n")); err != nil {
		t.Fatalf("ошибка отправки exit: %v", err)
	}
	if !waitForResponse("Сеанс завершен") {
		t.Fatal("не получено сообщение о завершении сеанса")
	}
}

func TestServer_RealTCPListenerAndStop(t *testing.T) {
	mock := &mockSearcher{
		docs: []crawler.Document{
			{ID: 10, Title: "Тестовый документ Go", URL: "http://example.com"},
		},
	}

	// Запускаем сервер на случайном свободном порту 127.0.0.1:0
	srv := New("127.0.0.1:0", mock)
	srv.SetTimeout(1 * time.Second)

	done := make(chan struct{})
	go func() {
		defer close(done)
		_ = srv.Start()
	}()

	// Даем серверу время привязать сокет
	time.Sleep(20 * time.Millisecond)

	addr := srv.Addr()
	if addr == nil {
		t.Fatalf("сервер не вернул локальный адрес")
	}

	// Подключаемся реальным клиентом через net.Dial
	conn, err := net.Dial("tcp", addr.String())
	if err != nil {
		t.Fatalf("не удалось подключиться к серверу: %v", err)
	}

	// Читаем заголовок
	reader := bufio.NewReader(conn)
	banner, err := reader.ReadString('\n')
	if err != nil {
		t.Fatalf("ошибка чтения приветствия: %v", err)
	}
	if !strings.Contains(banner, "GoSearch") {
		t.Errorf("неожиданный баннер: %s", banner)
	}

	// Отправляем команду выхода
	if _, err := conn.Write([]byte("exit\n")); err != nil {
		t.Fatalf("ошибка отправки: %v", err)
	}

	// Закрываем клиент
	_ = conn.Close()

	// Останавливаем сервер
	if err := srv.Stop(); err != nil {
		t.Errorf("ошибка остановки сервера: %v", err)
	}

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("таймаут ожидания остановки сервера")
	}
}

func BenchmarkServer_HandleConn_WithNetPipe(b *testing.B) {
	mock := &mockSearcher{
		docs: []crawler.Document{
			{ID: 1, Title: "Go Programming Language", URL: "https://go.dev"},
			{ID: 2, Title: "A Tour of Go", URL: "https://tour.golang.org"},
		},
	}
	srv := New(":8000", mock)
	srv.SetTimeout(0)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		serverConn, clientConn := net.Pipe()

		go srv.HandleConn(serverConn)

		reader := bufio.NewReader(clientConn)

		// Читаем приветственный баннер
		_, _ = reader.ReadString('>')

		// Отправляем поисковый запрос
		_, _ = clientConn.Write([]byte("Go\n"))
		_, _ = reader.ReadString('>')

		// Завершаем сессию
		_, _ = clientConn.Write([]byte("exit\n"))
		_, _ = reader.ReadString('\n')

		_ = clientConn.Close()
	}
}

