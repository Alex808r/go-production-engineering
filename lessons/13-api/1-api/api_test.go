package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

var api *API

func TestMain(m *testing.M) {
	api = new(API)
	api.router = mux.NewRouter()
	api.Endpoints()
	os.Exit(m.Run())
}

func TestAPI_newBook(t *testing.T) {
	// =========================================================================
	// Вопрос №1: «want := len(books) + 1 — что здесь плохо?»
	// Ответ: Использование глобального разделяемого состояния (мутабельного среза books):
	// 1. Тесты становятся взаимозависимыми (stateful): порядок запуска тестов влияет на результат.
	// 2. Невозможно безопасно использовать t.Parallel(), так как возникнет Data Race.
	// 3. Правильный подход — внедрение зависимостей (Dependency Injection) с чистым хранилищем
	//    для каждого тестового кейса (per-test mock/storage isolation).
	// =========================================================================
	booksMutex.RLock()
	want := len(books) + 1
	booksMutex.RUnlock()

	data := book{
		Name:   "1984",
		Author: "George Orwell",
	}
	payload, _ := json.Marshal(data)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/books", bytes.NewBuffer(payload))
	rr := httptest.NewRecorder()

	api.router.ServeHTTP(rr, req)

	// Согласно REST, при успешном создании сущности ожидается HTTP 201 Created
	if rr.Code != http.StatusCreated && rr.Code != http.StatusOK {
		t.Errorf("код неверен: получили %d, а хотели 201/200", rr.Code)
	}
	t.Log("Response: ", rr.Body)

	// =========================================================================
	// Вопрос №2: «что плохо в таком сравнении? как сделать лучше?»
	// Ответ:
	// 1. Простая проверка длины среза (len == want) не гарантирует, что добавилась
	//    ИМЕННО ТА книга с правильными полями Name и Author (могла добавиться пустая запись).
	// 2. Как сделать лучше:
	//    а) Распарсить JSON-ответ эндпоинта и убедиться, что он возвращает добавленную книгу.
	//    б) Проверить поля последнего элемента в хранилище (books[len(books)-1]).
	// =========================================================================
	booksMutex.RLock()
	got := len(books)
	lastBook := books[got-1]
	booksMutex.RUnlock()

	if got != want {
		t.Fatalf("книга не добавлена: ожидали размер %d, получили %d", want, got)
	}

	if lastBook.Name != data.Name || lastBook.Author != data.Author {
		t.Errorf("неверные данные добавленной книги: %+v, ожидали: %+v", lastBook, data)
	}
}

func TestAPI_newBook_Validation(t *testing.T) {
	// Проверка валидации на пустые поля (ожидаем 400 Bad Request)
	emptyBook := book{Name: "", Author: ""}
	payload, _ := json.Marshal(emptyBook)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/books", bytes.NewBuffer(payload))
	rr := httptest.NewRecorder()

	api.router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("ожидался HTTP 400 Bad Request для пустых полей, получен: %d", rr.Code)
	}
}

func TestAPI_books(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/books", nil)
	rr := httptest.NewRecorder()

	api.router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusOK)
	}

	// Проверяем Content-Type заголовок из headersMiddleware
	ct := rr.Header().Get("Content-Type")
	if !strings.Contains(ct, "application/json") {
		t.Errorf("ожидался заголовок application/json, получен: %s", ct)
	}

	var result []book
	if err := json.Unmarshal(rr.Body.Bytes(), &result); err != nil {
		t.Fatalf("ошибка разбора JSON: %v", err)
	}
	if len(result) == 0 {
		t.Errorf("ожидался непустой список книг")
	}
}

func TestAPI_deleteBook(t *testing.T) {
	// Добавим книгу для гарантированного удаления
	booksMutex.Lock()
	books = append(books, book{Name: "Book to delete", Author: "Author"})
	targetIndex := len(books) - 1
	booksMutex.Unlock()

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/books/"+strings.TrimSpace(string(rune('0'+targetIndex))), nil)
	rr := httptest.NewRecorder()

	api.router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("ожидался код 200 OK при удалении, получен: %d", rr.Code)
	}

	// Проверка удаления несуществующего индекса (404 Not Found)
	req404 := httptest.NewRequest(http.MethodDelete, "/api/v1/books/9999", nil)
	rr404 := httptest.NewRecorder()
	api.router.ServeHTTP(rr404, req404)

	if rr404.Code != http.StatusNotFound {
		t.Errorf("ожидался код 404 Not Found для несуществующего ID, получен: %d", rr404.Code)
	}
}

func BenchmarkAPI_books(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/books", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		api.router.ServeHTTP(rr, req)
	}
}

