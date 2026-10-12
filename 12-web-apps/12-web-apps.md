# 12-web-apps — Веб-разработка на Go: net/http, маршрутизация и HTML-шаблоны

В данной директории изучаются основы создания веб-сервисов и сайтов на Go с использованием пакета `net/http`, стороннего маршрутизатора `gorilla/mux`, шаблонизатора `html/template` и пакета тестирования `httptest`.

---

## 📋 Структура каталога

```text
12-web-apps/
├── 1-http/           # Базовый HTTP:
│   ├── server/       # Простейший HTTP-сервер (http.HandleFunc, http.ListenAndServe)
│   └── client/       # HTTP-клиент (http.Get, http.Client, http.Transport, Context)
├── 2-custom_server/  # Тонкая настройка http.Server (ReadTimeout, WriteTimeout, IdleTimeout, net.Listener)
├── 3-gorilla_mux/    # Маршрутизатор Gorilla Mux: path variables (/{name}), фильтрация методов (.Methods)
├── 4-template/       # Безопасный HTML-рендеринг через html/template (XSS auto-escaping) и цепочка Middleware (mux.Use)
├── 5-testing/        # Модульное тестирование хэндлеров без сетевых сокетов через net/http/httptest
├── hw/               # Анализ занимаемой памяти (Memory Layout, alignment, padding, shallow vs deep size)
└── homework-12/      # Практическое ДЗ: Веб-интерфейс GoSearch (HTML/JSON, gorilla/mux, httptest)
    ├── cmd/webapp/   # Точка входа HTTP-сервера GoSearch
    └── pkg/webapp/   # Хэндлеры, шаблоны, in-memory DB и модульные тесты с бенчмарками
```

---

## 🔍 Ключевые концепции

1. **Интерфейс `http.Handler`:**
   ```go
   type Handler interface {
       ServeHTTP(ResponseWriter, *Request)
   }
   ```
   Любая функция с сигнатурой `func(http.ResponseWriter, *http.Request)` адаптируется к интерфейсу через адаптер `http.HandlerFunc`.

2. **Безопасная конфигурация сервера (`2-custom_server`):**
   - Никогда не используйте `http.ListenAndServe(":8080", nil)` в продакшене без настройки таймаутов, так как медленные клиенты (Slowloris атака) исчерпают пул открытых сокетов и горутин.
   - Всегда задавайте `ReadTimeout`, `ReadHeaderTimeout`, `WriteTimeout` и `IdleTimeout`.

3. **Паттерн Middleware (`4-template`):**
   - Сигнатура middleware в Go: `func(http.Handler) http.Handler`.
   - Цепочка обязанностей («луковая» архитектура): код до `next.ServeHTTP(w, r)` выполняется на входящем запросе, код после — на исходящем ответе.

4. **Шаблонизация с защитой от XSS (`4-template`):**
   - Пакет `html/template` контекстно экранирует HTML, атрибуты, URL и JS в пользовательских данных, исключая уязвимости межсайтового скриптинга (в отличие от `text/template`).

5. **Тестирование через `net/http/httptest` (`5-testing`):**
   - `httptest.NewRequest(...)` создает виртуальный HTTP-запрос.
   - `httptest.NewRecorder()` перехватывает код ответа, заголовки и тело в оперативную память без открытия реальных TCP-портов.

---

## 🚀 Запуск и тестирование

```bash
# Запуск модульных тестов и бенчмарков модуля 12
go test -v -race -bench=. ./12-web-apps/...

# Запуск примера с Gorilla Mux
go run ./12-web-apps/3-gorilla_mux/gorilla.go

# Запуск сервера с Middleware и шаблонизацией
go run ./12-web-apps/4-template/template.go

# Запуск веб-приложения GoSearch (ДЗ-12)
go run ./12-web-apps/homework-12/cmd/webapp/main.go
```

---

## 🔗 Связанные материалы
- **Теория:** [Том 4: Backend и Web (Раздел 12: net/http, роутинг, middleware)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/04-go-core-backend-web.md)
- **Практическое ДЗ:** [homework-12](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/12-web-apps/homework-12) (Веб-интерфейс поисковика GoSearch с HTML-шаблонами)

