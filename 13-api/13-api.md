# 13-api — Проектирование REST API на Go

В данной директории изучаются принципы построения надежных RESTful API на языке Go: обработка JSON-запросов и ответов, корректная работа с HTTP статус-кодами, валидация входных данных, промежуточные обработчики (Middleware), аутентификация (Sessions и JWT) и модульное тестирование эндпоинтов.

---

## 📋 Структура каталога

```text
13-api/
├── 1-api/               # Базовый REST API: CRUD-хэндлеры, сериализация JSON, потоковый декодер, коды ошибок
│   ├── server.go        # Точка входа HTTP-сервера
│   ├── api.go           # Хэндлеры books, newBook, deleteBook
│   ├── middleware.go    # Логирование и установка JSON-заголовков
│   └── api_test.go      # Тестирование хэндлеров через httptest с ответами на вопросы о stateful тестах
├── 2-api/               # Промышленная архитектура REST API: аутентификация и контексты
│   ├── cmd/server/      # Исполняемый файл сервера API
│   └── pkg/api/         # Пакет API: Cookie-сессии (gorilla/sessions), JWT-токены, requestID middleware
└── homework-13/         # Практическое ДЗ: REST API поисковика GoSearch и исследование памяти
    ├── cmd/server/      # Сервер GoSearch API с Graceful Shutdown
    ├── cmd/task3/       # Задача №3: Разбор модели памяти, выравнивания и size.Of vs unsafe.Sizeof
    ├── pkg/storage/     # Потокобезопасное хранилище (storage.DB) с инвертированным индексом
    └── pkg/api/         # REST API хэндлеры (Search, CRUD) и middleware (Request-ID, Logging)
```

---

## 🔍 Ключевые концепции

1. **Идиоматичная работа с JSON:**
   - Чтение тела запроса через потоковый `json.NewDecoder(r.Body).Decode(&req)` без загрузки всего тела в промежуточный срез байт.
   - Ответ клиенту: установка заголовка `w.Header().Set("Content-Type", "application/json; charset=utf-8")`, отправка статус-кода `w.WriteHeader(http.StatusOK)` и потоковое кодирование `json.NewEncoder(w).Encode(resp)`.

2. **Статусы ответов REST API:**
   - `200 OK`: Успешное выполнение GET, PUT, DELETE.
   - `201 Created`: Успешное создание ресурса через POST.
   - `400 Bad Request`: Ошибка валидации входного JSON или параметров запроса.
   - `401 Unauthorized`: Отсутствует или недействителен токен аутентификации.
   - `403 Forbidden`: Недостаточно прав для доступа к ресурсу.
   - `404 Not Found`: Запрашиваемый ресурс не найден.
   - `405 Method Not Allowed`: Метод HTTP не поддерживается данным эндпоинтом.
   - `500 Internal Server Error`: Внутренняя ошибка сервера.

3. **Middleware (Промежуточное ПО):**
   - Паттерн "матрешка" (цепочка обязанностей), оборачивающий `http.Handler`:
     ```go
     func loggingMiddleware(next http.Handler) http.Handler {
         return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
             start := time.Now()
             next.ServeHTTP(w, r)
             log.Printf("%s %s took %v", r.Method, r.RequestURI, time.Since(start))
         })
     }
     ```
   - Применяется для сквозного логирования, сбора метрик, проброса Request ID в `r.Context()`, аутентификации и восстановления после паник.

4. **Аутентификация: Sessions vs JWT:**
   - **Cookie Sessions (`gorilla/sessions`):** состояние сессии хранится в зашифрованной/подписанной HMAC Cookie на клиенте.
   - **JWT (`dgrijalva/jwt-go`):** полностью stateless токен (Header.Payload.Signature), передаваемый в заголовке `Authorization: Bearer <token>`.

5. **Модель памяти (Задача №3):**
   - `unsafe.Sizeof`: размер структуры без данных, на которые ссылаются указатели (shallow size).
   - `size.Of`: рекурсивный размер объекта со всеми аллокациями в куче (deep size).

---

## 🚀 Запуск и тестирование

```bash
# Тестирование всех пакетов REST API и бенчмарки
go test -v -race -bench=. ./13-api/...

# Запуск базового сервера 1-api
go run ./13-api/1-api/server.go ./13-api/1-api/api.go ./13-api/1-api/middleware.go

# Запуск продвинутого сервера 2-api (с JWT и сессиями)
go run ./13-api/2-api/cmd/server/server.go

# Запуск REST API сервера GoSearch (ДЗ-13)
go run ./13-api/homework-13/cmd/server/main.go

# Запуск демонстрации модели памяти (Задача №3)
go run ./13-api/homework-13/cmd/task3/main.go
```

---

## 🔗 Связанные материалы
- **Теория:** [Том 4: Backend и Web (Раздел 12: REST архитектура, Middleware)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/04-go-core-backend-web.md)
- **Практическое ДЗ:** [homework-13](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/13-api/homework-13) (Полноценный REST API сервис поисковика GoSearch с CRUD-операциями)

