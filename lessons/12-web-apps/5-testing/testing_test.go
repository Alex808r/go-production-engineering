package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

var testMux *mux.Router

// TestMain запускается перед выполнением тестов пакета.
// Позволяет единоразово инициализировать общие ресурсы (например, тестовый роутер).
func TestMain(m *testing.M) {
	testMux = newRouter()
	m.Run()
}

// Test_mainHandler демонстрирует использование пакета net/http/httptest
// для быстрого и детерминированного тестирования HTTP-обработчиков.
//
// Преимущества httptest:
// 1. Вызовы происходят в оперативной памяти (ServeHTTP), без открытия реальных сетевых TCP-портов.
// 2. Нет конфликтов занятых портов при параллельном запуске в CI/CD (Zero port collision).
// 3. Тесты выполняются за микросекунды.
func Test_mainHandler(t *testing.T) {
	// -------------------------------------------------------------------------
	// Сценарий 1: Проверка запрета метода (POST вместо GET -> 405 Method Not Allowed)
	// -------------------------------------------------------------------------
	data := []int{1, 2, 3}
	payload, _ := json.Marshal(data)

	// httptest.NewRequest создает *http.Request, оптимизированный для вызова в тестах
	req := httptest.NewRequest(http.MethodPost, "/Alex", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	// httptest.NewRecorder() возвращает *httptest.ResponseRecorder,
	// который реализует интерфейс http.ResponseWriter и накапливает в памяти
	// код ответа (Code), заголовки (HeaderMap) и буфер тела ответа (Body).
	rr := httptest.NewRecorder()

	// Напрямую вызываем метод ServeHTTP роутера
	testMux.ServeHTTP(rr, req)

	// Проверяем, что Gorilla Mux отклонил POST-запрос с кодом 405
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("ожидался HTTP %d (Method Not Allowed), получен: %d", http.StatusMethodNotAllowed, rr.Code)
	}

	// -------------------------------------------------------------------------
	// Сценарий 2: Успешный GET-запрос (200 OK + проверка тела ответа)
	// -------------------------------------------------------------------------
	req = httptest.NewRequest(http.MethodGet, "/Alex", nil)
	rr = httptest.NewRecorder()

	testMux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("ожидался HTTP %d (OK), получен: %d", http.StatusOK, rr.Code)
	}

	body := rr.Body.String()
	expected := "Hi, Alex"
	if !strings.Contains(body, expected) {
		t.Errorf("тело ответа должно содержать %q, получено:\n%s", expected, body)
	}

	// Проверяем Content-Type заголовок
	ct := rr.Header().Get("Content-Type")
	if !strings.Contains(ct, "text/html") {
		t.Errorf("ожидался заголовок text/html, получен: %s", ct)
	}
}

// Benchmark_mainHandler измеряет производительность сквозной обработки запроса через роутер и хэндлер.
func Benchmark_mainHandler(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "/Gopher", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		testMux.ServeHTTP(rr, req)
	}
}

