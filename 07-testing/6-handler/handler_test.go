package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

var api API

// TestAPI_ProductsHandler демонстрирует идиоматичное тестирование HTTP-обработчиков
// с использованием пакета стандартной библиотеки net/http/httptest.
//
// ПРЕИМУЩЕСТВА httptest:
// 1. Без открытия сетевых сокетов: тест выполняется прямо в памяти без вызова listen()
//    на портах операционной системы (быстро, надежно, исключает конфликты портов в CI).
// 2. httptest.NewRequest: формирует синтетический *http.Request с нужным методом, URL,
//    заголовками и телом.
// 3. httptest.NewRecorder: создает реализацию http.ResponseWriter, которая записывает
//    HTTP-статус, заголовки и тело ответа во внутренний буфер bytes.Buffer.
// 4. Прямой вызов: обработчик вызывается как обычная функция: `handler(rr, req)`.
func TestAPI_ProductsHandler(t *testing.T) {
	// 1. Создание виртуального HTTP-запроса (GET /api/products)
	req := httptest.NewRequest(http.MethodGet, "/api/products", nil)

	// 2. Создание виртуального ResponseWriter для перехвата ответа
	rr := httptest.NewRecorder()

	// 3. Прямой вызов обработчика
	api.ProductsHandler(rr, req)

	// 4. Проверка HTTP статус-кода
	if rr.Code != http.StatusOK {
		t.Errorf("код ответа неверен: получили %d, ожидалось %d", rr.Code, http.StatusOK)
	}

	// 5. Проверка заголовка Content-Type
	if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Content-Type: получили %q, ожидалось %q", contentType, "application/json")
	}

	// 6. Чтение и десериализация тела ответа
	resp := rr.Result()
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("ошибка чтения тела ответа: %v", err)
	}

	var data []product
	err = json.Unmarshal(body, &data)
	if err != nil {
		t.Fatalf("ошибка парсинга JSON: %v", err)
	}

	if len(data) == 0 || data[0].Name != "Книга" {
		t.Errorf("неверные данные ответа: %v", data)
	}

	t.Logf("Ответ сервера успешно проверен:\n%+v", data)
}

