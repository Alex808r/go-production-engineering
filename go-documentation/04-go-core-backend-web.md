# Том 4: Backend, Web-разработка, Базы данных, Сетевые сокеты и Время

> [!TIP]
> 📚 **Навигация:** [⬅️ Назад: Том 3](03-go-core-concurrency.md) | [📖 Главное оглавление](README.md) | [Вперед: Том 5 (Тестирование и Архитектура) ➡️](05-go-core-testing-arch.md)

> [!IMPORTANT]
> 💻 **Практика, примеры кода и домашние задания к этому тому:**
> - 📂 **Код лекций:** [11-network](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/11-network) (TCP/UDP, сокеты, дедлайны), [12-web-apps](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/12-web-apps) (HTTP-сервер, роутинг, middleware), [13-api](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/13-api) (REST API, CRUD, JSON), [15-sql](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/15-sql) (DDL/DML, индексы), [16-db-apps](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/16-db-apps) (PostgreSQL `pgxpool`, транзакции, репозиторий).
> - 🎯 **Задачник:** [homework-tasks.md (Уроки 11–12, ДЗ 11–13, 15–16)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/homework-tasks.md#урок-11-web-http-rest-api-и-middleware).
> - 🛠️ **Решения домашних заданий:**
>   - [homework-11](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/11-network/homework-11) — Сетевая служба TCP для поисковика GoSearch.
>   - [homework-12](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/12-web-apps/homework-12) — Веб-служба GoSearch на `net/http`.
>   - [homework-13](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/13-api/homework-13) — REST API с CRUD-эндпоинтами для поисковых документов.
>   - [homework-15](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/15-sql/homework-15) — Схема базы данных PostgreSQL для онлайн-кинотеатра.
>   - [homework-16](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/16-db-apps/homework-16) — Пакет репозитория PostgreSQL с пулом `pgxpool`.
> - 🗄️ **Справочники:** [go-database-interview-guide.md](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/go-database-interview-guide.md) (Полный гайд по БД), [sql-livecoding-guide.md](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/sql-livecoding-guide.md) (28 SQL-задач).
> - 🗺️ **Сквозной путеводитель:** [course-codebase-guide.md (Этап 3: Middle+ Developer)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/course-codebase-guide.md#этап-3-middle-go-developer--production-инженерия-базы-данных-и-микросервисы).

---

## 12. Backend: Web-разработка и HTTP (`net/http`)

> 🔗 **Практика и примеры кода:** [12-web-apps](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/12-web-apps), [13-api](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/13-api) | 🛠️ **Домашние задания:** [homework-12 (Web)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/12-web-apps/homework-12), [homework-13 (REST API)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/13-api/homework-13) | 🎯 **Задачи:** [homework-tasks.md (Урок 11: Web REST API)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/homework-tasks.md#урок-11-web-http-rest-api-и-middleware), [homework-tasks.md (ДЗ 12–13)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/homework-tasks.md#домашнее-задание-12-веб-приложения-на-go-веб-служба-gosearch)

### 12.1. Базовый HTTP-сервер и `http.ServeMux`

#### 💡 На пальцах: Как устроен HTTP-сервер в Go

В отличие от многих других платформ (Node.js с одним потоком событий или Python WSGI с пулом тяжёлых процессов), в Go веб-сервер устроен максимально элегантно:
- 🚀 **1 запрос = 1 горутина:**  
  Когда клиент подключается по TCP, метод `server.ListenAndServe()` моментально запускает отдельную горутину: `go c.serve(connCtx)`. Если к серверу одновременно придут 50 000 пользователей, Go создаст 50 000 ультра-легких горутин (~2 КБ памяти каждая), которые будут обслуживаться параллельно на всех доступных ядрах CPU!
- 📤 **`w http.ResponseWriter` — труба для отправки ответа:**  
  Интерфейс для отправки данных клиенту. Через него задаются HTTP-заголовки, статус-код ответа и пишется само тело ответа.
- 📥 **`r *http.Request` — входящий конверт запроса:**  
  Структура, содержащая всё, что прислал клиент: HTTP-метод (`GET`, `POST`), путь URL, query-параметры (`r.URL.Query()`), заголовки и поток чтения тела (`r.Body`).

> [!CAUTION]
> **Золотое правило `ResponseWriter`:**  
> Метод `w.Header().Set(...)` должен вызываться **строго ДО** первого вызова `w.WriteHeader(...)` или `w.Write(...)`! В момент первого вызова `Write()` рантайм Go немедленно отправляет HTTP-заголовки в сокет со статусом `200 OK`. Любые попытки изменить заголовки после этого будут проигнорированы.

```go
package main

import (
    "encoding/json"
    "log"
    "net/http"
    "time"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
    // w (ResponseWriter) — куда пишем ответ
    // r (*Request) — откуда читаем данные запроса (headers, query, body)
    w.Header().Set("Content-Type", "application/json") // заголовки — строго ДО WriteHeader/Write
    w.WriteHeader(http.StatusOK)                         // 200 OK (если не вызвать, отправится автоматически при первом Write)
    // Формируем JSON через кодировщик, а не через Sprintf: так спецсимволы в name будут корректно экранированы
    _ = json.NewEncoder(w).Encode(map[string]string{"message": "Hello " + r.URL.Query().Get("name")})
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("GET /api/v1/hello", helloHandler) // В Go 1.22+ методы встроены в маршрут!

    server := &http.Server{
        Addr:              ":8080",
        Handler:           mux,
        ReadHeaderTimeout: 2 * time.Second,   // лимит на чтение заголовков (защита от Slowloris-атак)
        ReadTimeout:       5 * time.Second,   // лимит на чтение всего запроса
        WriteTimeout:      10 * time.Second,  // лимит на отправку ответа
        IdleTimeout:       120 * time.Second, // сколько держать Keep-Alive соединение без запросов
    }

    // ListenAndServe блокируется и возвращает ошибку (например, порт занят) — её нельзя игнорировать:
    log.Fatal(server.ListenAndServe())
}
```

> [!WARNING]
> Сервер без таймаутов (`http.ListenAndServe(":8080", mux)` в «голом» виде) уязвим: медленные или зависшие клиенты удерживают соединения и горутины бесконечно. Всегда создавайте `http.Server` с таймаутами.

Запуск и проверка:

```bash
go run .                                          # в первом терминале
curl 'localhost:8080/api/v1/hello?name=Go'        # во втором терминале → {"message":"Hello Go"}
```

---

### 12.2. Экосистема роутеров и фреймворков в Go (`chi`, `gin`, `echo`, `fiber`)

В Go отношение к фреймворкам отличается от других языков (где правят Django, Spring Boot или Ruby on Rails). Стандартная библиотека `net/http` самодостаточна, поэтому сообщество чаще использует термин **роутер (router)**, а не полноценный фреймворк.

| Роутер / Фреймворк          | Особенности & Применимость                                                                                                                                         | Совместимость с `net/http` |
| :-------------------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------- | :------------------------- |
| **`net/http` (Go 1.22+)**   | Стандартный `http.ServeMux` теперь поддерживает HTTP-методы (`"GET /users/{id}"`) и path-параметры. Внешние библиотеки для простых API больше не требуются.        | 100% стандарт              |
| **`chi` (go-chi/chi)**      | 🌟 **Один из самых популярных выборов:** лёгкий, 100% совместим со стандартным `http.Handler`, идеальная поддержка цепочек middleware.                                  | 100% совместим             |
| **`gin` (gin-gonic/gin)**   | Самый популярный по звездам на GitHub фреймворк. Имеет собственный `gin.Context`, быструю валидацию JSON, встроенный биндинг.                                      | Собственный контекст       |
| **`echo` (labstack/echo)**  | Быстрый, лаконичный микрофреймворк с удобным рендерингом и встроенным middleware.                                                                                  | Собственный контекст       |
| **`fiber` (gofiber/fiber)** | Построен поверх `fasthttp` (не использует `net/http`), дизайн вдохновлен Express.js (Node.js). Рекордный RPS в бенчмарках, но требует осторожности с пулом памяти. | Не совместим с `net/http`  |

---

### 12.3. Middleware (Промежуточное ПО)

#### 💡 На пальцах: Middleware как слои луковицы (или матрёшка)

Middleware в Go — это слои вокруг хендлера:
1. Запрос летит **снаружи внутрь к ядру**:  
   `Клиент ➔ [Recovery] ➔ [Logging] ➔ [Auth] ➔ Сам Хендлер`
2. Хендлер выполняет бизнес-логику и генерирует ответ.
3. Ответ возвращается **изнутри наружу к клиенту**:  
   `Сам Хендлер ➔ [Auth] ➔ [Logging] (замер времени) ➔ [Recovery] ➔ Клиент`

```mermaid
flowchart LR
    Client["Клиент"] --> Rec["1. Recovery Middleware"]
    Rec --> Log["2. Logging Middleware"]
    Log --> Auth["3. Auth Middleware"]
    Auth --> Handler["4. Бизнес-хендлер"]

    Handler -.-> Auth
    Auth -.-> Log
    Log -.-> Rec
    Rec -.-> Client
```

Сигнатура Middleware в Go предельно проста: функция принимает следующий обработчик `http.Handler` и возвращает новый обёрнутый `http.Handler`:

$$\text{Middleware Signature: } \text{func}(\text{http.Handler}) \rightarrow \text{http.Handler}$$

```go
// Recovery Middleware (Защита сервера от паники в хендлере):
func RecoveryMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                if err == http.ErrAbortHandler {
                    panic(err) // специальная «тихая» паника net/http — пробрасываем дальше
                }
                log.Printf("panic: %v\n%s", err, debug.Stack()) // ВСЕГДА логируем причину и стек
                http.Error(w, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
            }
        }()
        next.ServeHTTP(w, r)
    })
}
```

> Сам `http.Server` тоже перехватывает панику каждого обработчика (иначе процесс бы падал), но в этом случае он лишь пишет стек в лог и обрывает соединение — клиент не получит корректный ответ. Middleware выше позволяет вернуть внятный `500`. Порядок вызова middleware важен: цепочка `Recovery(Logging(Auth(handler)))` выполняется «снаружи внутрь», поэтому `Recovery` обычно ставят самым внешним.

---

### 12.4. Graceful Shutdown (Элегантное завершение сервера)

При получении сигнала завершения (`SIGINT`/`SIGTERM` от Docker или Kubernetes) сервер **не должен мгновенно умирать**, обрывая клиентские запросы:

```go
ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
defer stop()

go func() {
    if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
        log.Fatalf("listen: %s\n", err)
    }
}()

<-ctx.Done() // Ждем сигнал от ОС (Ctrl+C или k8s SIGTERM)

// Даем серверу 10 секунд на плавное завершение активных запросов:
shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

if err := server.Shutdown(shutdownCtx); err != nil {
    log.Fatal("Server forced to shutdown:", err)
}
fmt.Println("Server gracefully stopped")
```

> [!NOTE]
> - `Shutdown` перестаёт принимать новые соединения и ждёт завершения активных запросов, но **не** ждёт соединения, которые были «захвачены» (hijacked), например WebSocket. Для них используйте `server.RegisterOnShutdown(...)` и закрывайте вручную.
> - После `Shutdown` вызов `ListenAndServe` возвращает `http.ErrServerClosed` — это штатное завершение, а не ошибка.
> - В Kubernetes после `SIGTERM` под ещё некоторое время получает трафик от балансировщика; таймаут остановки (`terminationGracePeriodSeconds`) должен быть больше вашего `10*time.Second`.
> - Здесь `stop()` из `signal.NotifyContext` восстанавливает поведение по умолчанию: повторный `Ctrl+C` завершит программу немедленно.

---

### 12.5. Архитектура REST API и слои приложения (Standard Go Project Layout)

В реальных production-проектах на Go код не пишут в одном файле или пакете. Общепринятым стандартом организации микросервисов является **Standard Go Project Layout** с четким разделением ответственности (Separation of Concerns) и чистой слоистой архитектурой:

```text
my-rest-service/
├── cmd/
│   └── api/
│       └── main.go           # Точка входа (Composition Root): инициализация конфигов, пула БД, DI и запуск
├── internal/                 # Приватный код проекта (защищен компилятором Go от внешнего импорта)
│   ├── database/             # Подключение к СУБД, конфигурация пула соединений
│   ├── handlers/             # HTTP-слой (Transport): парсинг запросов, вызов бизнес-логики/репозитория, статус-коды
│   ├── models/               # Доменные структуры данных и DTO (Data Transfer Objects)
│   └── repository/           # Паттерн Repository: изоляция чистых SQL-запросов от бизнес-логики
├── sql/                      # Скрипты инициализации БД и миграции (schema.sql, init.sql)
├── docker-compose.yml        # Локальное окружение (PostgreSQL, Redis)
└── go.mod
```

> [!IMPORTANT]
> **Магия пакета `internal/` в компиляторе Go:**  
> Имя директории `internal` — зарезервированное ключевое соглашение тулчейна Go. Любой пакет внутри `internal/` **не может быть импортирован** из внешних репозиториев или других модулей. Попытка сделать `import "github.com/user/project/internal/..."` в чужом коде вызовет ошибку компиляции:  
> `use of internal package ... not allowed`. Это запрещает импортировать такие пакеты любому коду за пределами родительского дерева каталога `internal` (в типичном проекте — за пределами вашего модуля).

#### Внедрение зависимостей (Dependency Injection) через структуру Handler

Чтобы обработчики были легко тестируемыми (Unit-тесты без реальной БД) и не использовали глобальные переменные, зависимости передаются явно через структуры и конструкторы:

```go
package handlers

import (
    "my-rest-service/internal/repository"
)

type Handler struct {
    store *repository.TaskStore // Зависимость на хранилище данных (или интерфейс)
}

func NewHandler(store *repository.TaskStore) *Handler {
    return &Handler{store: store}
}
```

---

### 12.6. Идиоматичный REST API без фреймворков (CRUD, JSON Helpers, ловушки)

Хотя в индустрии популярны фреймворки вроде Gin или Chi, глубокое понимание стандартной библиотеки `net/http` является ключевым требованием на собеседованиях.

#### Вспомогательные функции (JSON Helpers)

Для соблюдения принципа DRY (Don't Repeat Yourself) и единообразных ответов создаются хелперы сериализации:

```go
package handlers

import (
    "encoding/json"
    "net/http"
)

// respondWithJSON сериализует payload в JSON и устанавливает правильный Content-Type
func respondWithJSON(w http.ResponseWriter, statusCode int, payload any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    if payload != nil {
        _ = json.NewEncoder(w).Encode(payload)
    }
}

// respondWithError возвращает унифицированный JSON с описанием ошибки
func respondWithError(w http.ResponseWriter, statusCode int, message string) {
    respondWithJSON(w, statusCode, map[string]string{"error": message})
}
```

#### Потоковый декодинг: `json.NewDecoder` vs `json.Unmarshal`

```go
// 💡 Идиоматичный прием: чтение тела запроса на лету без загрузки всего слайса байт в память:
var input models.CreateTaskInput
if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
    respondWithError(w, http.StatusBadRequest, "Некорректный формат JSON")
    return // ⚠️ КРИТИЧЕСКИ ВАЖНО: выход из функции!
}
```

> [!WARNING]
> **Главная ловушка новичков в HTTP-хендлерах:**  
> Функции `http.Error(...)`, `respondWithError(...)` или `w.WriteHeader(...)` **НЕ прерывают выполнение функции хендлера!**  
> Если забыть написать ключевое слово `return` после отправки ошибки, функция продолжит исполняться дальше, попытается выполнить запрос к базе данных и повторно запишет заголовки в сокет:  
> в лог попадёт предупреждение `http: superfluous response.WriteHeader call`, а клиент получит смешанный/искажённый ответ (код статуса уже отправлен и изменить его нельзя).

---

### 12.7. Паттерн DTO и указатели для частичного обновления (Partial Update)

При реализации **частичного** обновления (по семантике REST это метод `PATCH`; `PUT` — полная замена ресурса) возникает проблема нулевых значений (Zero Values) языка Go:

```go
// ❌ ПРОБЛЕМА: Невозможно отличить "поле не передали" от "поле передали со значением по умолчанию"
type BadUpdateInput struct {
    Title     string `json:"title"`     // Если поле пропущено в JSON, здесь будет ""
    Completed bool   `json:"completed"` // Если поле пропущено в JSON, здесь будет false!
}
// Если клиент прислал `{"title": "New"}`, мы случайно сбросим `completed` в false!
```

#### Решение через указатели (`*string`, `*bool`):

```go
type UpdateTaskInput struct {
    Title       *string `json:"title"`       // nil = не передано, non-nil = новое значение
    Description *string `json:"description"` // nil = не передано, non-nil = новое значение
    Completed   *bool   `json:"completed"`   // nil = не передано, non-nil = true/false
}
```

В коде обработчика или репозитория проверка выполняется тривиально:

```go
if input.Title != nil {
    // Клиент явно передал заголовок (даже если передал пустую строку ""):
    task.Title = *input.Title
}
if input.Completed != nil {
    // Клиент явно изменил статус задачи:
    task.Completed = *input.Completed
}
```

> ⚠️ Ограничение: при декодировании JSON поле `"description": null` и отсутствующее поле дают одинаковый результат (`nil`). Если вам нужно отличать «очистить значение» от «не менять», используют собственный тип-обёртку (`Optional[T]` с флагом `Set` и кастомным `UnmarshalJSON`) или `map[string]json.RawMessage`.

---

### 12.8. Эволюция маршрутизации: ручной диспатчинг (Pre-1.22) vs Native Routing (Go 1.22+)

#### 1. Исторический подход (Go < 1.22):

Стандартный `http.ServeMux` умел матчить только префикс URL. Поэтому маршрутизацию по методам (`GET`, `POST`, `DELETE`) и разбор параметров пути приходилось реализовывать вручную:

```go
// Ручной роутинг в хендлере /tasks/
func (h *Handler) TaskByIDHandler(w http.ResponseWriter, r *http.Request) {
    // Извлечение параметра ID из пути вручную:
    // r.URL.Path = "/tasks/42" -> idStr = "42"
    pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/tasks/"), "/")
    id, err := strconv.Atoi(pathParts[0])
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Некорректный ID задачи")
        return
    }

    // Ручной диспатчинг HTTP-метода через switch:
    switch r.Method {
    case http.MethodGet:
        h.getTaskByID(w, r, id)
    case http.MethodPut:
        h.updateTask(w, r, id)
    case http.MethodDelete:
        h.deleteTask(w, r, id)
    default:
        respondWithError(w, http.StatusMethodNotAllowed, "Метод не поддерживается")
    }
}
```

#### 2. Современный подход (Go 1.22+):

В Go 1.22 роутер `http.ServeMux` получил первоклассную поддержку HTTP-методов и шаблонов путей (`Path Values`), сделав сторонние библиотеки для простых API ненужными:

```go
mux := http.NewServeMux()

// Маршруты регистрируются с явным указанием метода и плейсхолдера {id}:
mux.HandleFunc("GET /tasks", h.GetAllTasks)
mux.HandleFunc("POST /tasks", h.CreateTask)
mux.HandleFunc("GET /tasks/{id}", h.GetTaskByID)
mux.HandleFunc("PUT /tasks/{id}", h.UpdateTask)
mux.HandleFunc("DELETE /tasks/{id}", h.DeleteTask)

// Внутри хендлера параметр извлекается нативно:
func (h *Handler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
    idStr := r.PathValue("id") // Встроенный метод Go 1.22+
    id, err := strconv.Atoi(idStr)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Невалидный ID")
        return
    }
    // ...
}
```

> [!NOTE]
> Дополнительные возможности шаблонов Go 1.22:
> - `{path...}` — «хвост» пути (`"GET /files/{path...}"`), `{$}` — точное совпадение корня (`"GET /{$}"`).
> - Шаблон `"GET /tasks"` также обрабатывает `HEAD`; если путь существует, но метод не подходит, мультиплексор сам вернёт `405 Method Not Allowed`.
> - При конфликте шаблонов побеждает более специфичный; неразрешимые конфликты вызывают `panic` при регистрации (то есть ловятся сразу на старте).

---

### 12.9. Docker Compose для PostgreSQL в локальной разработке

Для комфортной локальной разработки без необходимости установки СУБД на хост-машину используется контейнеризация через Docker Compose.

```yaml
# Поле `version:` в современных версиях Docker Compose не нужно (считается устаревшим).
# Учётные данные ниже — только для ЛОКАЛЬНОЙ разработки. В проде секреты берут из менеджера секретов.
services:
  postgres:
    image: postgres:15-alpine
    container_name: tasks_postgres
    restart: unless-stopped
    environment:
      POSTGRES_DB: tasks_db
      POSTGRES_USER: task_user
      POSTGRES_PASSWORD: task_password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      # Скрипты из этой папки автоматически выполняются при ПЕРВОМ старте контейнера:
      - ./sql/init.sql:/docker-entrypoint-initdb.d/init.sql
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U task_user -d tasks_db"]
      interval: 5s
      timeout: 5s
      retries: 5

volumes:
  postgres_data:
```

> [!TIP]
> **Полезные команды управления контейнером:**
>
> - `docker compose up -d` — запуск базы данных в фоновом режиме.
> - `docker compose ps` — проверка статуса (включая `healthy`).
> - `docker exec -it tasks_postgres psql -U task_user -d tasks_db` — интерактивный вход в консоль СУБД.
> - `docker compose down -v` — остановка и удаление тома с данными (полная очистка базы).

### 12.10. Подводные камни `http.Client` и `http.Transport` в Highload-продакшене

> [!WARNING]
> Стандартный `http.DefaultClient` **не следует использовать в микросервисном продакшене**.  
> По умолчанию у него `Timeout = 0` (таймаут отсутствует). Если внешний сервис или шлюз зависнет, горутина будет ждать ответа бесконечно. Тысячи таких запросов приводят к исчерпанию файловых дескрипторов ОС (`too many open files`) и краху приложения по памяти.

#### 1. Три смертельные ловушки `net/http` клиента:

1. **`DefaultTransport.MaxIdleConnsPerHost = 2`:**
   По умолчанию пул открытых Keep-Alive соединений к **одному хосту** равен всего **2**! Если ваш микросервис отправляет 500 RPS к соседнему сервису (например, платежному шлюзу или auth-провайдеру), 498 соединений будут ежесекундно открываться и тут же закрываться (`TCP FIN`), забивая операционную систему тысячами сокетов в состоянии `TIME_WAIT` и истощая диапазон эфемерных портов.

2. **Непрочитанное до конца тело ответа (`resp.Body`):**
   Даже если вызвать `defer resp.Body.Close()`, соединение **не вернется в пул Keep-Alive**, если из него не было вычитано все тело до байта `io.EOF`. При разрыве соединение закрывается на уровне TCP, требуя заново выполнять TCP Handshake и TLS-рукопожатие на следующий запрос (+50-100 мс задержки).

3. **Чтение тела без ограничения размера (`io.ReadAll`):**
   Использование `io.ReadAll(resp.Body)` без обертки в `io.LimitReader` позволяет злоумышленнику или сбойному сервису вернуть ответ размером в несколько гигабайт (Decompression Bomb или огромный JSON), из-за чего процесс будет убит по нехватке памяти (`OOM`) в вашем контейнере.

#### 2. Эталонная настройка `http.Client` для Highload-продакшена:

```go
package httpClient

import (
    "context"
    "fmt"
    "io"
    "net"
    "net/http"
    "time"
)

func NewProductionHTTPClient() *http.Client {
    // 1. Тонкая настройка сетевого транспорта (пул TCP-соединений)
    transport := &http.Transport{
        Proxy: http.ProxyFromEnvironment,
        DialContext: (&net.Dialer{
            Timeout:   5 * time.Second,  // Таймаут установления TCP-соединения
            KeepAlive: 30 * time.Second, // Период отправки TCP keep-alive зондов
        }).DialContext,
        MaxIdleConns:          500,               // Общий лимит свободных соединений в пуле
        MaxIdleConnsPerHost:   100,               // ⚠️ Лимит свободных соединений к КАЖДОМУ хосту (дефолт 2!)
        MaxConnsPerHost:       200,               // Максимум одновременных соединений к хосту
        IdleConnTimeout:       90 * time.Second,  // Время жизни неиспользуемого коннекта в пуле
        TLSHandshakeTimeout:   5 * time.Second,   // Таймаут на криптографическое рукопожатие TLS
        ResponseHeaderTimeout: 10 * time.Second,  // Таймаут ожидания ПЕРВЫХ байт заголовков ответа
        ExpectContinueTimeout: 1 * time.Second,
    }

    // 2. Клиент с глобальным жестким таймаутом на ВЕСЬ запрос (включая редиректы и чтение)
    return &http.Client{
        Transport: transport,
        Timeout:   15 * time.Second,
    }
}

// DoRequest демонстрирует безопасное выполнение запроса с дренажом тела
func DoRequest(ctx context.Context, client *http.Client, url string) ([]byte, error) {
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
    if err != nil {
        return nil, fmt.Errorf("создание запроса: %w", err)
    }

    resp, err := client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("выполнение HTTP запроса: %w", err)
    }

    // ⚠️ ГАРАНТИЯ ПЕРЕИСПОЛЬЗОВАНИЯ СОЕДИНЕНИЯ (Drain + Close):
    defer func() {
        // Вычитываем остатки тела в черную дыру, чтобы сокет остался чистым для пула
        // (с ограничением: не читаем бесконечно длинное тело ради переиспользования соединения)
        _, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 1<<20))
        _ = resp.Body.Close()
    }()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("неожиданный статус код: %d", resp.StatusCode)
    }

    // Безопасное чтение: ограничиваем максимальный размер тела 10 МБ от OOM
    limitedReader := io.LimitReader(resp.Body, 10<<20)
    data, err := io.ReadAll(limitedReader)
    if err != nil {
        return nil, fmt.Errorf("чтение ответа: %w", err)
    }

    return data, nil
}
```

> 💡 `http.Client` потокобезопасен и предназначен для **повторного использования**: создайте один клиент при старте приложения и передавайте его в зависимости. Создание нового клиента (и транспорта) на каждый запрос лишает вас пула соединений и приводит к утечке сокетов.

---

### 12.11. Маршрутизатор Gorilla Mux (`github.com/gorilla/mux`)

До появления Go 1.22 стандартный `http.ServeMux` не поддерживал переменные пути и фильтрацию по HTTP-методам. Одним из самых популярных и зрелых решений в экосистеме Go стал пакет `github.com/gorilla/mux`.

> [!NOTE]
> Проект Gorilla был архивирован в конце 2022 года и затем возобновлён силами сообщества. Для **новых** проектов с простой маршрутизацией сегодня, как правило, достаточно стандартного `http.ServeMux` (Go 1.22+) либо `chi`; `gorilla/mux` полезно знать, так как он широко встречается в существующем коде.

#### 1. Ключевые возможности `gorilla/mux`:
- **Переменные пути (Path Variables)**: сопоставление сегментов URL и извлечение через `mux.Vars(r)`.
- **Регулярные выражения в путях**: валидация параметров прямо на уровне роутера (например, только цифры для ID).
- **Ограничение методов**: сопоставление по HTTP-глаголам (`GET`, `POST`, `PUT`, `DELETE`).
- **Middleware**: встроенный стек промежуточного ПО через `r.Use(middleware)`.
- **Субмаршрутизация (Subrouting)**: группировка путей с общими префиксами (например, `/api/v1`).

```go
package main

import (
    "fmt"
    "net/http"

    "github.com/gorilla/mux"
)

func main() {
    r := mux.NewRouter()

    // Регистрация глобального middleware
    r.Use(loggingMiddleware)

    // Маршрут с переменной пути и ограничением по HTTP-методу GET
    r.HandleFunc("/users/{id:[0-9]+}", getUserHandler).Methods(http.MethodGet)

    // Подмаршрутизатор (Subrouter) для версионирования API
    api := r.PathPrefix("/api/v1").Subrouter()
    api.HandleFunc("/items", listItemsHandler).Methods(http.MethodGet)
    api.HandleFunc("/items", createItemHandler).Methods(http.MethodPost)

    http.ListenAndServe(":8080", r)
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID := vars["id"] // извлечение значения {id}
    fmt.Fprintf(w, "Просмотр пользователя с ID: %s", userID)
}

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Printf("[%s] %s\n", r.Method, r.URL.Path)
        next.ServeHTTP(w, r)
    })
}
```

---

### 12.12. Традиционные веб-приложения и SSR: Шаблонизация (`html/template`)

В традиционных веб-приложениях HTML генерируется сервером (Server-Side Rendering — SSR) на основе данных и отправляется клиенту. Для этого в стандартной библиотеке Go предназначен пакет **`html/template`**.

> [!IMPORTANT]
> **Фундаментальное отличие `html/template` от `text/template`:**  
> Пакет `html/template` реализует **контекстно-зависимое автоэкранирование (Context-aware Auto-escaping)**. Он анализирует контекст вставки переменной (внутри HTML-тега, внутри атрибута `href=""`, внутри CSS-стиля или внутри блока `<script>`) и автоматически экранирует опасные спецсимволы, гарантируя защиту от **XSS-атак (Cross-Site Scripting)**. Пакет `text/template` такого экранирования не производит и предназначен только для генерации простого текста/конфигураций.

#### 1. Основные конструкции шаблонов:
- `{{ .Field }}` — подстановка значения поля структуры. Точка `.` обозначает текущий контекст (pipeline).
- `{{ range .Slice }}...{{ else }}...{{ end }}` — итерация по срезу. Внутри блока `.` указывает на текущий элемент.
- `{{ if .Condition }}...{{ else }}...{{ end }}` — ветвление по условию.
- `{{ template "name" . }}` — включение другого именованного шаблона.

#### 2. Пример безопасного рендеринга страницы:

```go
package main

import (
    "html/template"
    "net/http"
)

type DocumentItem struct {
    ID    int
    Title string
    URL   string
}

type PageData struct {
    Heading string
    Items   []DocumentItem
}

const pageTmpl = `<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>{{ .Heading }}</title>
</head>
<body>
    <h1>{{ .Heading }}</h1>
    <ul>
    {{ range .Items }}
        <li><strong>#{{ .ID }}</strong>: <a href="{{ .URL }}">{{ .Title }}</a></li>
    {{ else }}
        <li>Документов пока нет.</li>
    {{ end }}
    </ul>
</body>
</html>`

// Шаблон парсится ОДИН раз при старте программы, а не на каждый запрос:
var tmpl = template.Must(template.New("page").Parse(pageTmpl))

func renderHandler(w http.ResponseWriter, r *http.Request) {

    data := PageData{
        Heading: "Каталог ресурсов",
        Items: []DocumentItem{
            {ID: 1, Title: "Официальный сайт Go", URL: "https://go.dev"},
            {ID: 2, Title: "<script>alert('xss')</script>", URL: "https://example.com"}, // Автоматически экранируется!
        },
    }

    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    if err := tmpl.Execute(w, data); err != nil {
        // Если часть HTML уже отправлена, статус изменить нельзя — в проде рендерите в bytes.Buffer
        // и только при успехе отправляйте клиенту
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
```

#### 3. Пользовательские функции (`template.FuncMap`):
Если требуется трансформация данных в шаблоне (форматирование дат, перевод в верхний регистр), регистрируют словарь функций:

```go
funcMap := template.FuncMap{
    "toUpper": strings.ToUpper,
    "formatDate": func(t time.Time) string {
        return t.Format("02.01.2006")
    },
}
tmpl := template.Must(template.New("report").Funcs(funcMap).Parse(rawTmpl))
```

---

### 12.13. Файловый сервер, раздача статики и Single Page Applications (SPA)

Современные веб-приложения часто состоят из Single Page Application (React, Vue, Angular) на фронтенде и Go REST API на бэкенде. Для раздачи скомпилированных ассетов (HTML, JS, CSS, изображений) стандартная библиотека предоставляет `http.FileServer`.

#### 1. Базовая раздача статических файлов (`http.FileServer`):

```go
// Раздаем содержимое директории ./public по адресу /static/
fs := http.FileServer(http.Dir("./public"))
http.Handle("/static/", http.StripPrefix("/static/", fs))
```
> [!NOTE]
> Функция `http.StripPrefix` обязательна, если URL-префикс (`/static/`) не совпадает с путем внутри директории (`./public`), иначе файловый сервер будет искать файл в `./public/static/file.css`.

#### 2. Встраивание статики в бинарный файл через `embed.FS` (Go 1.16+):
В Go можно упаковать весь фронтенд прямо внутрь скомпилированного бинарника, исключая проблемы с отсутствующими файлами на сервере:

```go
package main

import (
    "embed"
    "io/fs"
    "log"
    "net/http"
)

//go:embed static
var staticFS embed.FS

func main() {
    // Встроенная ФС содержит каталог static/ целиком, поэтому «входим» в него через fs.Sub:
    // иначе файлы были бы доступны по URL /static/index.html, а не /index.html
    sub, err := fs.Sub(staticFS, "static")
    if err != nil {
        log.Fatal(err)
    }
    // http.FS адаптирует io/fs.FS к http.FileSystem
    http.Handle("/", http.FileServer(http.FS(sub)))
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

#### 3. SPA Handler (Fallback на `index.html`):
Для SPA с клиентским роутингом (HTML5 History API) любой неизвестный путь (например, `/dashboard/settings`) должен отдавать `index.html`, чтобы клиентский роутер приложения сам отобразил нужный экран:

```go
type spaHandler struct {
    staticPath string
    indexPath  string
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // filepath.Clean("/"+...) убирает ".." и не даёт выйти за пределы staticPath (Path Traversal)
    path := filepath.Join(h.staticPath, filepath.Clean("/"+r.URL.Path))
    // Проверяем, существует ли физический файл на диске
    fi, err := os.Stat(path)
    if os.IsNotExist(err) || fi.IsDir() {
        // Если файла нет — отдаем главный index.html
        http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
        return
    }
    // Если файл существует — отдаем его обычным способом
    http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}
```

#### 4. Защита от Path Traversal и Directory Escape: `os.Root` (Go 1.24+)

Исторически раздача пользовательских файлов по пути из URL таила риск уязвимостей **Path Traversal** (например, запрос `GET /files/../../../../etc/passwd` или обход через симлинки).
В **Go 1.24** представлен тип `os.Root` (`os.OpenRoot(dir)`), который изолирует операции с файловой системой строго внутри указанной корневой директории:

```go
// Go 1.24+: Доступ ограничен пределами каталога (Sandboxed Directory): пути проверяются на уровне ОС-вызовов (openat)
root, err := os.OpenRoot("/var/www/uploads")
if err != nil {
    log.Fatal(err)
}
defer root.Close()

// Любая попытка выйти за пределы root (через ../ или symlink) вернёт ошибку (path escapes from parent)
file, err := root.Open(userProvidedPath)
```

---

### 12.14. Запуск защищенного веб-сервера по TLS/HTTPS (`ListenAndServeTLS`)

В промышленной эксплуатации веб-трафик обязан быть зашифрован с использованием протокола TLS (HTTPS).

#### 1. Запуск через `http.ListenAndServeTLS`:

```go
// Для запуска требуются путь к сертификату (cert.pem) и приватному ключу (key.pem)
err := http.ListenAndServeTLS(":443", "cert.pem", "key.pem", handler)
```

#### 2. Тонкая настройка TLS и автоматический редирект HTTP $\to$ HTTPS:

```go
package main

import (
    "crypto/tls"
    "net/http"
    "time"
)

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Безопасное соединение по HTTPS!"))
    })

    // 1. Безопасная конфигурация TLS (отключаем устаревшие TLS 1.0 / 1.1)
    tlsConfig := &tls.Config{
        MinVersion:       tls.VersionTLS12,
        CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
    }
    // Примечание: поле PreferServerCipherSuites устарело и игнорируется начиная с Go 1.18

    httpsServer := &http.Server{
        Addr:         ":443",
        Handler:      mux,
        TLSConfig:    tlsConfig,
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  120 * time.Second,
    }

    // 2. В отдельной горутине запускаем HTTP-сервер на 80 порту для редиректа на HTTPS
    go func() {
        redirectServer := &http.Server{
            Addr: ":80",
            Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                target := "https://" + r.Host + r.URL.RequestURI()
                http.Redirect(w, r, target, http.StatusMovedPermanently) // 301 Redirect
            }),
            ReadTimeout: 3 * time.Second,
        }
        _ = redirectServer.ListenAndServe()
    }()

    // Запускаем основной HTTPS сервер
    _ = httpsServer.ListenAndServeTLS("cert.pem", "key.pem")
}
```

> 💡 Для локальной разработки сертификат можно создать через `mkcert` или `openssl`; в продакшене чаще всего TLS завершают на балансировщике/Ingress. Если TLS терминируется в самом Go-приложении, автоматическое получение сертификатов Let's Encrypt даёт пакет `golang.org/x/crypto/acme/autocert`.

---

### 12.15. Архитектура пакета API и ООП-модель сервиса

В современной Go-разработке программный интерфейс приложения (**API**) рекомендуется оформлять в виде **самостоятельного изолированного пакета** (например, `pkg/api` или `internal/api`).

#### 1. Public API vs Private API
- **Public API (Внешний API):** Предназначен для взаимодействия со сторонними системами и внешними разработчиками. Требует строгой обратной совместимости, версионирования (`/api/v1/`, `/api/v2/`), публичной документации (OpenAPI/Swagger) и строгой аутентификации (API-ключи, OAuth2, JWT).
- **Private API (Внутренний API):** Обслуживает собственные клиентские приложения (Single Page Applications на React/Vue, мобильные клиенты iOS/Android, сервисы фронтенда). Может содержать специализированные endpoints под экраны UI (BFF — Backend for Frontend) и сессионную аутентификацию.
- **API Gateway (Шлюз API):** В микросервисной архитектуре единая точка входа, которая принимает клиентские запросы, выполняет сквозную аутентификацию, rate limiting, логирование и маршрутизирует трафик к целевым микросервисам.

#### 2. Структура API как объект с внедрением зависимостей (Dependency Injection)

Вместо глобальных переменных и анонимных функций пакет API проектируется в ООП-стиле:
- Создается структура `API`, содержащая маршрутизатор и ссылки на зависимости (база данных, сервис кэширования, логгер, конфигурация).
- Обработчики эндпоинтов реализуются как **методы структуры `API`** (`func (api *API) handleSearch(...)`), получая прямой доступ ко всем сервисам без глобального состояния.
- Конструктор `New(...)` создает экземпляр API и регистрирует маршруты и middleware через метод `Endpoints()`:

```go
package api

import (
    "net/http"
    "github.com/gorilla/mux"
)

// Storage описывает контракт хранилища данных
type Storage interface {
    Find(id int) (*Document, error)
    Save(doc Document) error
    Delete(id int) error
}

// API инкапсулирует транспортный HTTP-интерфейс и его зависимости
type API struct {
    router  *mux.Router
    storage Storage
}

// New инициализирует API с внедрением зависимостей
func New(storage Storage) *API {
    api := &API{
        router:  mux.NewRouter(),
        storage: storage,
    }
    api.endpoints()
    return api
}

// Router возвращает готовый маршрутизатор для запуска в http.Server
func (api *API) Router() *mux.Router {
    return api.router
}

// endpoints регистрирует промежуточное ПО и маршруты
func (api *API) endpoints() {
    // Подключение сквозных middleware
    api.router.Use(requestIDMiddleware)
    api.router.Use(loggingMiddleware)

    // Регистрация бизнес-маршрутов
    api.router.HandleFunc("/api/v1/docs", api.handleDocs).Methods(http.MethodGet)
    api.router.HandleFunc("/api/v1/docs", api.handleCreateDoc).Methods(http.MethodPost)
    api.router.HandleFunc("/api/v1/docs/{id:[0-9]+}", api.handleDocByID).Methods(http.MethodGet)
    api.router.HandleFunc("/api/v1/docs/{id:[0-9]+}", api.handleUpdateDoc).Methods(http.MethodPut)
    api.router.HandleFunc("/api/v1/docs/{id:[0-9]+}", api.handleDeleteDoc).Methods(http.MethodDelete)
}
```

---

### 12.16. Контекст в HTTP-обработчиках: отмена цепочки вызовов и проброс метаданных

Контекст (`context.Context`) в веб-разработке на Go решает **две фундаментальные задачи**:
1. **Каскадная отмена операций (Cancellation & Timeouts):** При разрыве клиентом TCP-соединения (закрытие вкладки браузера, нажатие Cancel) или превышении таймаута запрос должен немедленно отменяться на всех уровнях: обработчик $\to$ бизнес-логика $\to$ выполнение SQL-запроса в СУБД $\to$ HTTP-вызовы сторонних API.
2. **Передача сквозных метаданных запроса (Request-Scoped Values):** Идентификатор трассировки (`Request ID`), аутентифицированный пользователь (`UserID`, `Role`), IP-клиента.

#### 1. Каскадная отмена и таймауты (`context.WithTimeout`)

Каждый входящий `*http.Request` уже содержит базовый контекст `r.Context()`, который автоматически отменяется рантаймом Go при закрытии клиентского сокета. Обработчик может задать жесткий дедлайн на выполнение операции:

```go
func (api *API) handleHeavyTask(w http.ResponseWriter, r *http.Request) {
    // Создаем производный контекст с жестким таймаутом в 3 секунды
    ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
    defer cancel() // Обязательное освобождение ресурсов таймера!

    // Передаем контекст глубже в базу данных или внешний сервис
    result, err := api.db.QueryWithContext(ctx, "SELECT ...")
    if err != nil {
        if errors.Is(ctx.Err(), context.DeadlineExceeded) {
            http.Error(w, "Gateway Timeout", http.StatusGatewayTimeout) // 504
            return
        }
        if errors.Is(ctx.Err(), context.Canceled) {
            // Клиент оборвал соединение — прерываем обработку без записи в сокет
            return
        }
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    _ = json.NewEncoder(w).Encode(result)
}
```

#### 2. Передача значений через контекст и Middleware (`r.WithContext`)

Для передачи значений создается производный контекст через `context.WithValue`, а затем запрос обогащается новым контекстом с помощью **`r.WithContext(ctx)`**:

```go
type contextKey string

const requestIDKey contextKey = "request_id"

// requestIDMiddleware генерирует уникальный ID запроса и внедряет его в контекст
func requestIDMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        reqID := r.Header.Get("X-Request-ID")
        if reqID == "" {
            reqID = fmt.Sprintf("req-%d", time.Now().UnixNano()) // упрощённо; в проде используйте UUID или crypto/rand — UnixNano не уникален при конкурентных запросах
        }

        // 1. Обогащаем контекст запроса значением
        ctx := context.WithValue(r.Context(), requestIDKey, reqID)

        // 2. Устанавливаем заголовок ответа для клиента
        w.Header().Set("X-Request-ID", reqID)

        // 3. Вызываем следующий обработчик с обновленным запросом r.WithContext(ctx)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// Извлечение значения из контекста в нижележащих слоях
func GetRequestID(ctx context.Context) string {
    if id, ok := ctx.Value(requestIDKey).(string); ok {
        return id
    }
    return ""
}
```

> [!CAUTION]
> Никогда не используйте встроенные типы (например, обычный `string`) в качестве ключей `context.WithValue`! Всегда объявляйте **пользовательский неэкспортируемый тип** (`type contextKey string`), чтобы исключить коллизии ключей между сторонними библиотеками.

---

### 12.17. Аутентификация и авторизация в API: Сессии vs JWT

| Критерий | Сессии на сервере (`gorilla/sessions`) | JWT (JSON Web Token, `golang-jwt/jwt`) |
| :--- | :--- | :--- |
| **Хранение состояния** | **Stateful**: данные сессии хранятся на сервере (память, Redis, БД), клиенту выдается Cookie с ID сессии. | **Stateless**: все данные (claims) упакованы в сам токен, сервер не хранит состояние сессии. |
| **Масштабирование** | Требуется централизованное хранилище (Redis/Memcached) при работе за балансировщиком. | Идеально для микросервисов: любой сервис может проверить подпись открытым ключом без обращений к БД. |
| **Отзыв прав (Logout)** | Мгновенный: удаление ключа из Redis немедленно аннулирует доступ. | Затруднен: токен валиден до истечения `exp`. Требуются черные списки (blacklist) в Redis или короткий TTL (5–15 мин) + Refresh-токены. |
| **Типичные клиенты** | Традиционные веб-приложения и браузеры (HttpOnly Cookies). | Мобильные клиенты, SPA (React, Vue, Flutter), межсервисное взаимодействие (M2M). |

#### 1. Анатомия JWT-токена

JWT представляет собой строку из трех частей, разделенных точками (`.`):
$$\text{JWT} = \text{Base64Url}(\text{Header}) \,.\, \text{Base64Url}(\text{Payload}) \,.\, \text{Base64Url}(\text{Signature})$$

1. **Header (Заголовок):** Содержит алгоритм подписи (`alg`, например `HS256`, `RS256`) и тип токена (`typ: "JWT"`).
2. **Payload (Полезная нагрузка):** JSON-объект с утверждениями (**Claims**):
   - Стандартные клеймы: `sub` (идентификатор субъекта/пользователя), `exp` (unix timestamp истечения), `iat` (время выпуска).
   - Пользовательские клеймы: `role`, `email`, `permissions`.
3. **Signature (Цифровая подпись):** Вычисляется от склеенных заголовка и полезной нагрузки с секретным ключом сервера:
   $$\text{Signature} = \text{HMAC-SHA256}(\text{Header} + "." + \text{Payload}, \text{SecretKey})$$

> [!WARNING]
> Данные в Payload **не зашифрованы**, а лишь закодированы в Base64URL! Любой клиент может декодировать токен и прочитать данные. **Никогда не храните пароли, секреты или персональные данные в Payload JWT!** Цифровая подпись гарантирует исключительно **целостность и неизменность** данных, а не их секретность.

#### 2. Реализация на Go

Полная пошаговая реализация — выпуск и проверка токена, middleware, пара Access + Refresh токенов с Redis и вход через OAuth2 — приведена ниже в разделе [12.20. Аутентификация: JWT и OAuth2 в Go](#1220-аутентификация-jwt-и-oauth2-в-go). Используйте актуальную библиотеку `github.com/golang-jwt/jwt/v5`; старый пакет `dgrijalva/jwt-go` больше не поддерживается и содержит известные уязвимости.

---

### 12.18. Документирование и контрактное тестирование API (OpenAPI / Swagger и Postman)

В промышленной разработке API должно сопровождаться машиночитаемой и интерактивной документацией.

#### 1. Стандарт OpenAPI (Swagger)
- **OpenAPI Specification (OAS 3.0 / 3.1):** Международный стандарт описания REST API в формате JSON или YAML. Описывает эндпоинты, входные параметры, структуры JSON (Schemas), статус-коды ответов и методы авторизации.
- **Swagger UI:** Веб-интерфейс, визуализирующий OpenAPI-спецификацию и позволяющий отправлять запросы прямо из браузера ("Try it out").
- В Go документацию генерируют:
  - Либо методом **Code-First** через комментарии к обработчикам с помощью утилиты `swag` (`swaggo/swag`).
  - Либо методом **Spec-First (Contract-First)**, описывая `openapi.yaml`, из которого с помощью `oapi-codegen` генерируются серверные интерфейсы и клиентские SDK.

#### 2. Интеграционное тестирование через Postman
- **Postman Collections:** Наборы сохраненных HTTP-запросов с переменными окружения (`{{baseUrl}}`, `{{token}}`).
- **Автоматические assertions:** На вкладке *Tests* в Postman пишутся JavaScript-проверки:
  ```javascript
  pm.test("Status code is 200", function () {
      pm.response.to.have.status(200);
  });
  pm.test("Response contains docs", function () {
      var jsonData = pm.response.json();
      pm.expect(jsonData.length).to.be.above(0);
  });
  ```
- **CI/CD запуск:** Консольный раннер **Newman** (`newman run collection.json -e env.json`) позволяет запускать коллекции Postman в пайплайнах GitHub Actions / GitLab CI для регрессионного end-to-end тестирования сервиса.

---

### 12.19. Веб-безопасность (Web Security): CORS, Rate Limiting и защита Cookies/CSRF

При разработке публичных Web API безопасность — ключевой критерий приёмки кода на Code Review и технических собеседованиях.

#### 1. CORS (Cross-Origin Resource Sharing): Preflight и Middleware

Когда веб-фронтенд (например, `https://myapp.com`) обращается к Go API на другом домене (`https://api.myapp.com`), браузер блокирует ответ согласно политике Same-Origin Policy (SOP).

##### Предварительный запрос (Preflight Request):
Для запросов с методами `PUT`, `DELETE` или заголовками `Authorization` / `Content-Type: application/json` браузер сначала отправляет легкий запрос `OPTIONS`:

```go
func CORSMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Разрешенные домены (в проде выносить в конфиг, избегать "*")
        w.Header().Set("Access-Control-Allow-Origin", "https://myapp.com")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        w.Header().Set("Access-Control-Allow-Credentials", "true")

        // Если это preflight-запрос OPTIONS, прерываем цепочку и сразу отдаем 204 No Content:
        if r.Method == http.MethodOptions {
            w.Header().Set("Access-Control-Max-Age", "600") // браузер кэширует результат preflight на 10 минут
            w.WriteHeader(http.StatusNoContent)
            return
        }

        next.ServeHTTP(w, r)
    })
}
```

> ⚠️ Если нужно разрешить **несколько** доменов, нельзя просто отправлять список в `Access-Control-Allow-Origin` — заголовок принимает только одно значение. Сверьте `r.Header.Get("Origin")` со списком разрешённых, верните его же и добавьте заголовок `Vary: Origin`. Комбинация `Access-Control-Allow-Origin: *` с `Allow-Credentials: true` браузеры отвергают. CORS — механизм защиты **браузера**; он не защищает API от запросов через `curl` или другие серверы.

---

#### 2. Ограничение частоты запросов (Rate Limiting) по алгоритму Token Bucket

Для защиты API от перегрузки, DoS-атак и перебора паролей стандартным в Go является пакет `golang.org/x/time/rate`, реализующий алгоритм **Token Bucket (Корзина токенов)**:

```go
package main

import (
    "net"
    "net/http"
    "sync"

    "golang.org/x/time/rate"
)

// In-Memory Rate Limiter с привязкой по IP-адресу клиента:
type IPRateLimiter struct {
    mu      sync.Mutex
    clients map[string]*rate.Limiter
    r       rate.Limit // количество токенов в секунду
    b       int        // максимальный всплеск (burst)
}

func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
    return &IPRateLimiter{
        clients: make(map[string]*rate.Limiter),
        r:       r,
        b:       b,
    }
}

func (lim *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
    lim.mu.Lock()
    defer lim.mu.Unlock()

    l, exists := lim.clients[ip]
    if !exists {
        // Например, NewIPRateLimiter(5, 10): 5 запросов в секунду с всплеском до 10
        l = rate.NewLimiter(lim.r, lim.b)
        lim.clients[ip] = l
    }
    return l
}

// Middleware:
func (lim *IPRateLimiter) LimitMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // r.RemoteAddr имеет вид "ip:port", а порт у каждого соединения свой —
        // берём только IP, иначе лимит будет считаться отдельно для каждого соединения!
        ip, _, err := net.SplitHostPort(r.RemoteAddr)
        if err != nil {
            ip = r.RemoteAddr
        }
        if !lim.GetLimiter(ip).Allow() {
            http.Error(w, `{"error": "Too Many Requests"}`, http.StatusTooManyRequests)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

> [!WARNING]
> Учебная реализация имеет два ограничения: (1) мапа `clients` растёт бесконечно — в реальном коде нужна периодическая очистка неактивных IP (или LRU/TTL-кэш); (2) за прокси/балансировщиком `RemoteAddr` — это адрес самого прокси. Заголовку `X-Forwarded-For` можно доверять **только** если он выставлен вашим собственным доверенным прокси (клиент может подделать его). Для нескольких инстансов сервиса лимиты хранят в общем хранилище (Redis).

---

#### 3. Безопасность Cookie и защита от XSS / CSRF

При аутентификации через сессионные Cookie необходимо строго выставлять флаги защиты:

```go
http.SetCookie(w, &http.Cookie{
    Name:     "session_token",
    Value:    token,
    Path:     "/",
    Expires:  time.Now().Add(24 * time.Hour),
    HttpOnly: true,                  // 🛡️ Защита от XSS: JavaScript (document.cookie) НЕ имеет доступа к Cookie!
    Secure:   true,                  // 🛡️ Передача Cookie разрешена ТОЛЬКО по зашифрованному HTTPS-каналу
    SameSite: http.SameSiteStrictMode, // 🛡️ Защита от CSRF: Cookie не прикрепляется при переходах со сторонних сайтов
})
```

> [!NOTE]
> `SameSite` — важный, но не единственный слой защиты от CSRF. Для изменяющих запросов (`POST`/`PUT`/`DELETE`) дополнительно используют CSRF-токены (например, пакет `gorilla/csrf`) либо проверку заголовков `Origin`/`Sec-Fetch-Site`. Начиная с Go 1.25 в `net/http` есть готовый `http.CrossOriginProtection`, реализующий такую проверку. Режим `Strict` может ломать сценарии «перешёл по ссылке из письма и остался залогиненным» — часто выбирают `Lax`.

---

### 12.20. Аутентификация: JWT и OAuth2 в Go

> [!TIP]
> **Зачем это знать?** Вопросы об аутентификации задают на почти каждом собеседовании на позицию Backend Go Developer. Это практический стандарт индустрии.

---

#### 12.20.1. Что такое JWT? (JSON Web Token)

**JWT** — это компактный токен, который сервер выдаёт клиенту после входа. Клиент хранит его и отправляет с каждым запросом, чтобы доказать «кто я».

##### Структура JWT — три части через точку:

```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9   <- Header  (алгоритм + тип)
.
eyJ1c2VySWQiOjQyLCJyb2xlIjoiYWRtaW4ifQ  <- Payload (данные: userID, role, exp)
.
SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQ  <- Signature (подпись = гарантия подлинности)
```

> [!NOTE]
> **Ключевое правило:** Payload **НЕ зашифрован** — его можно декодировать из base64 без ключа. Подпись лишь **гарантирует**, что токен не был изменён. **Никогда не кладите пароли в JWT!**

##### Как это работает — простая схема:

```mermaid
sequenceDiagram
    participant Client as Клиент (браузер/мобилка)
    participant Server as Go Сервер

    Client->>Server: POST /login {email, password}
    Server->>Server: Проверяет пароль в БД
    Server-->>Client: 200 OK + JWT токен (подписан секретным ключом)

    Note over Client: Сохраняет токен

    Client->>Server: GET /profile<br/>Authorization: Bearer <JWT>
    Server->>Server: Проверяет подпись токена
    Server-->>Client: 200 OK + данные профиля
```

---

#### 12.20.2. Реализация JWT в Go (`golang-jwt/jwt`)

```bash
go get github.com/golang-jwt/jwt/v5
```

##### Шаг 1: Генерация токена при входе

```go
package auth

import (
    "os"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

// Секретный ключ — хранится в переменной окружения, НИКОГДА не в коде!
// (для HS256 используйте длинный случайный ключ, не менее 32 байт; в проде проверьте, что он не пустой)
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// Claims — данные, которые мы кладём внутрь токена.
// jwt.RegisteredClaims содержит стандартные поля: exp (срок), iss (издатель) и т.д.
type Claims struct {
    UserID int64  `json:"userId"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}

// GenerateToken создаёт подписанный JWT для пользователя.
func GenerateToken(userID int64, role string) (string, error) {
    claims := Claims{
        UserID: userID,
        Role:   role,
        RegisteredClaims: jwt.RegisteredClaims{
            // Токен действует 15 минут (Access Token должен быть короткоживущим!)
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "my-go-service",
        },
    }

    // Создаём токен с алгоритмом HMAC-SHA256 и подписываем секретным ключом:
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret) // Возвращает строку "header.payload.signature"
}
```

##### Шаг 2: Проверка токена при каждом запросе

```go
// (дополнительно понадобится импорт "fmt")
// ParseToken проверяет подпись и срок действия токена.
// Возвращает данные из токена (Claims) или ошибку, если токен недействителен.
func ParseToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(
        tokenString,
        &Claims{},
        // Эта функция вызывается библиотекой для получения ключа проверки.
        // Важно: проверяем алгоритм! Защита от атаки "alg:none".
        func(token *jwt.Token) (any, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("неожиданный алгоритм подписи: %v", token.Header["alg"])
            }
            return jwtSecret, nil
        },
        // Дополнительная защита в v5: разрешаем только явно указанные алгоритмы
        // и требуем наличия срока действия (exp) в токене:
        jwt.WithValidMethods([]string{"HS256"}),
        jwt.WithExpirationRequired(),
    )

    if err != nil {
        return nil, fmt.Errorf("невалидный токен: %w", err)
    }

    claims, ok := token.Claims.(*Claims)
    if !ok || !token.Valid {
        return nil, fmt.Errorf("не удалось прочитать данные токена")
    }

    return claims, nil
}
```

##### Шаг 3: HTTP Middleware для защиты маршрутов

```go
// (дополнительно понадобятся импорты "context", "net/http", "strings")
// AuthMiddleware — обёртка над хендлером. Пропускает запрос только с валидным JWT.
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 1. Читаем заголовок: "Authorization: Bearer eyJhbGci..."
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "отсутствует заголовок Authorization", http.StatusUnauthorized)
            return
        }

        // 2. Проверяем формат "Bearer <token>"
        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) != 2 || parts[0] != "Bearer" {
            http.Error(w, "неверный формат токена", http.StatusUnauthorized)
            return
        }

        // 3. Парсим и проверяем токен
        claims, err := ParseToken(parts[1])
        if err != nil {
            // Детали ошибки клиенту не раскрываем — логируем на сервере
            http.Error(w, "недействительный токен", http.StatusUnauthorized)
            return
        }

        // 4. Кладём данные пользователя в контекст запроса,
        //    чтобы хендлеры могли их получить без повторного парсинга.
        ctx := context.WithValue(r.Context(), contextKeyUserID, claims.UserID)
        ctx = context.WithValue(ctx, contextKeyRole, claims.Role)

        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// Типобезопасные ключи контекста (не используем строки напрямую!)
type contextKey string
const (
    contextKeyUserID contextKey = "userID"
    contextKeyRole   contextKey = "role"
)

// Хелперы для извлечения данных из контекста в хендлерах:
func UserIDFromCtx(ctx context.Context) int64 {
    v, _ := ctx.Value(contextKeyUserID).(int64)
    return v
}

// Пример использования middleware при регистрации маршрутов:
// mux.Handle("GET /profile", AuthMiddleware(http.HandlerFunc(profileHandler)))
```

---

#### 12.20.3. Паттерн: Access Token + Refresh Token

**Проблема:** Если Access Token живёт долго (например, 7 дней) и его украдут — злоумышленник имеет полный доступ неделю.

**Решение:** Два токена:

```mermaid
graph LR
    AT["Access Token\n(короткий: 15 мин)\nПередаётся в каждом запросе"]
    RT["Refresh Token\n(длинный: 30 дней)\nХранится безопасно, используется РЕДКО"]

    AT -->|"Истёк? Используй"| RT
    RT -->|"Выдаёт новый"| AT
```

```go
import (
    "context"
    "crypto/rand"
    "encoding/hex"
    "encoding/json"
    "errors"
    "fmt"
    "net/http"
    "strconv"
    "time"

    "github.com/redis/go-redis/v9"
)

// TokenPair — пара токенов, которую сервер отправляет клиенту при входе.
type TokenPair struct {
    AccessToken  string `json:"accessToken"`  // Короткоживущий, идёт в Authorization header
    RefreshToken string `json:"refreshToken"` // Долгоживущий, хранится в httpOnly cookie
}

// IssueTokenPair выдаёт пару токенов и сохраняет refresh token в Redis.
func IssueTokenPair(ctx context.Context, rdb *redis.Client, userID int64, role string) (TokenPair, error) {
    // 1. Создаём Access Token (JWT, 15 минут)
    accessToken, err := GenerateToken(userID, role)
    if err != nil {
        return TokenPair{}, err
    }

    // 2. Создаём Refresh Token — просто случайные байты (не JWT!)
    //    Он не несёт данных — только служит «билетом» для обновления.
    rawBytes := make([]byte, 32)
    if _, err := rand.Read(rawBytes); err != nil { // crypto/rand, НЕ math/rand!
        return TokenPair{}, fmt.Errorf("генерация refresh token: %w", err)
    }
    refreshToken := hex.EncodeToString(rawBytes) // "a3f8c1d2e4..."

    // 3. Сохраняем в Redis: ключ = refreshToken, значение = userID, TTL = 30 дней
    //    При logout — просто удаляем ключ из Redis, и токен станет недействительным.
    key := "refresh:" + refreshToken
    err = rdb.Set(ctx, key, userID, 30*24*time.Hour).Err()
    if err != nil {
        return TokenPair{}, fmt.Errorf("не удалось сохранить refresh token: %w", err)
    }

    return TokenPair{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

// RefreshHandler — хендлер обновления токенов.
// Клиент присылает refresh token -> получает новую пару токенов.
func RefreshHandler(rdb *redis.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // 1. Читаем refresh token из тела запроса или httpOnly cookie
        var req struct {
            RefreshToken string `json:"refreshToken"`
        }
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "неверный формат запроса", http.StatusBadRequest)
            return
        }

        // 2-3. Атомарно читаем И удаляем токен из Redis (GETDEL, Redis 6.2+).
        //      Rotation: каждый refresh token одноразовый; атомарность защищает от гонки,
        //      когда два запроса одновременно предъявляют один и тот же токен.
        key := "refresh:" + req.RefreshToken
        userIDStr, err := rdb.GetDel(r.Context(), key).Result()
        if errors.Is(err, redis.Nil) {
            http.Error(w, "refresh token не найден или истёк", http.StatusUnauthorized)
            return
        }
        if err != nil { // например, Redis недоступен — это НЕ то же самое, что «токен не найден»
            http.Error(w, "ошибка сервера", http.StatusInternalServerError)
            return
        }

        userID, err := strconv.ParseInt(userIDStr, 10, 64)
        if err != nil {
            http.Error(w, "ошибка сервера", http.StatusInternalServerError)
            return
        }

        // 4. Выдаём новую пару токенов
        //    (роль здесь для краткости зашита в код; в реальном проекте её читают из БД по userID)
        pair, err := IssueTokenPair(r.Context(), rdb, userID, "user")
        if err != nil {
            http.Error(w, "ошибка сервера", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(pair)
    }
}
```

> [!WARNING]
> **Частые ошибки с JWT:**
>
> 1. **Длинный Access Token** — типичная ошибка: токен на 7 дней. Типичные значения: **5–15 минут** для Access, **дни–недели** для Refresh (и Refresh в БД лучше хранить в виде хэша, а не «как есть»).
> 2. **Не проверять `exp`** — библиотека `golang-jwt/jwt` делает это автоматически при `ParseWithClaims`. Не отключайте!
> 3. **Атака `alg:none`** — злоумышленник подделывает токен с `"alg":"none"`. Защита: явно проверяйте метод (как в `ParseToken` выше).
> 4. **JWT в localStorage** — уязвимо к XSS. Refresh Token храните в **httpOnly cookie** (недоступна JavaScript; при таком хранении нужна защита от CSRF — см. 12.19).
> 5. **Нет logout** — JWT нельзя «отозвать» без хранилища. Паттерн с Redis Refresh Tokens решает это: `DEL refresh:<token>`.

---

#### 12.20.4. OAuth2 — вход через Google / GitHub

**OAuth2** — это стандарт, позволяющий вашему приложению использовать аккаунт другого сервиса (Google, GitHub, VK) для входа, **не зная пароль пользователя**.

##### Как это работает — Authorization Code Flow:

```mermaid
sequenceDiagram
    participant User as Пользователь
    participant App as Ваш Go Сервис
    participant Google as Google

    User->>App: Нажимает "Войти через Google"
    App-->>User: Перенаправляет на accounts.google.com<br/>(с client_id и redirect_uri)
    User->>Google: Вводит логин и даёт разрешение
    Google-->>User: Перенаправляет на ваш /oauth/callback?code=XXXX
    User->>App: GET /oauth/callback?code=XXXX
    App->>Google: POST /token (code + client_secret)
    Google-->>App: access_token для Google API
    App->>Google: GET /userinfo (кто это?)
    Google-->>App: {email: "user@gmail.com", name: "Ivan"}
    App-->>User: Выдаёт свой JWT токен (авторизован!)
```

##### Реализация на Go (`golang.org/x/oauth2`):

```bash
go get golang.org/x/oauth2
go get golang.org/x/oauth2/google
```

```go
package main

import (
    "crypto/rand"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "net/http"
    "os"

    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
)

// Конфигурация OAuth2 для Google.
// client_id и client_secret берём из Google Cloud Console (console.cloud.google.com).
var googleOAuthConfig = &oauth2.Config{
    ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
    ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
    // RedirectURL — куда Google перенаправит пользователя после разрешения.
    // Должен быть зарегистрирован в Google Cloud Console!
    RedirectURL: "http://localhost:8080/oauth/callback",
    // Scopes — какие данные просим у Google:
    Scopes: []string{
        "https://www.googleapis.com/auth/userinfo.email",
        "https://www.googleapis.com/auth/userinfo.profile",
    },
    Endpoint: google.Endpoint,
}

// GoogleUserInfo — данные, которые Google вернёт о пользователе.
type GoogleUserInfo struct {
    ID      string `json:"id"`
    Email   string `json:"email"`
    Name    string `json:"name"`
    Picture string `json:"picture"`
}

// LoginHandler — перенаправляет пользователя на страницу входа Google.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
    // state — случайная строка для защиты от CSRF атак: генерируем для КАЖДОГО входа
    // и запоминаем в короткоживущей httpOnly-cookie (в реальном проекте — можно и в серверной сессии).
    b := make([]byte, 16)
    if _, err := rand.Read(b); err != nil {
        http.Error(w, "ошибка сервера", http.StatusInternalServerError)
        return
    }
    state := base64.URLEncoding.EncodeToString(b)
    http.SetCookie(w, &http.Cookie{
        Name: "oauth_state", Value: state, Path: "/", MaxAge: 600,
        HttpOnly: true, Secure: false, // Secure: true — обязательно в проде (HTTPS)
        SameSite: http.SameSiteLaxMode,
    })

    // Получаем URL страницы авторизации Google:
    authURL := googleOAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)

    // Перенаправляем пользователя (он увидит окно "Войти через Google"):
    http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// CallbackHandler — Google перенаправляет сюда пользователя с одноразовым кодом.
func CallbackHandler(w http.ResponseWriter, r *http.Request) {
    // 1. Проверяем state (защита от CSRF): значение из URL должно совпасть со значением из cookie
    stateCookie, err := r.Cookie("oauth_state")
    if err != nil || stateCookie.Value == "" || r.FormValue("state") != stateCookie.Value {
        http.Error(w, "неверный state", http.StatusBadRequest)
        return
    }

    // 2. Обмениваем одноразовый code на access_token Google.
    //    Это происходит на сервере — клиент не видит client_secret!
    code := r.FormValue("code")
    googleToken, err := googleOAuthConfig.Exchange(r.Context(), code)
    if err != nil {
        http.Error(w, "ошибка обмена кода на токен", http.StatusInternalServerError) // детали err — только в лог
        return
    }

    // 3. Используем Google токен для получения данных пользователя.
    client := googleOAuthConfig.Client(r.Context(), googleToken)
    resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
    if err != nil {
        http.Error(w, "ошибка получения данных пользователя", http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    var userInfo GoogleUserInfo
    if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
        http.Error(w, "ошибка парсинга данных", http.StatusInternalServerError)
        return
    }

    // 4. Ищем или создаём пользователя в нашей БД по email.
    //    user, err := userRepo.FindOrCreate(ctx, userInfo.Email, userInfo.Name)

    // 5. Выдаём наш собственный JWT токен (не Google токен!).
    //    Дальше клиент работает с нашими токенами — Google больше не нужен.
    myJWT, err := GenerateToken(/* user.ID */ 42, "user")
    if err != nil {
        http.Error(w, "ошибка генерации токена", http.StatusInternalServerError)
        return
    }

    // 6. Возвращаем токен клиенту (или редиректим с токеном):
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "accessToken": myJWT,
        "email":       userInfo.Email,
        "name":        userInfo.Name,
    })
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("GET /login/google", LoginHandler)      // Шаг 1: начало входа
    mux.HandleFunc("GET /oauth/callback", CallbackHandler) // Шаг 2: Google возвращает код

    // Защищённые маршруты — требуют JWT:
    mux.Handle("GET /profile", AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        userID := UserIDFromCtx(r.Context())
        fmt.Fprintf(w, "Привет, пользователь %d!", userID)
    })))

    http.ListenAndServe(":8080", mux)
}
```

> [!NOTE]
> Для «публичных» клиентов (мобильные/SPA-приложения) и в целом как современная практика к Authorization Code Flow добавляют **PKCE**: в `golang.org/x/oauth2` для этого есть `oauth2.GenerateVerifier()`, `oauth2.S256ChallengeOption(verifier)` и `oauth2.VerifierOption(verifier)`. Для проверки Google ID-токена (OpenID Connect) используют библиотеку `github.com/coreos/go-oidc`.

> [!TIP]
> **Сводная таблица: что выбрать?**
>
> | Сценарий | Решение |
> | :--- | :--- |
> | Свой логин/пароль в вашей БД | JWT (Access + Refresh) |
> | «Войти через Google/GitHub» | OAuth2 Authorization Code Flow |
> | Межсервисная аутентификация (service-to-service) | OAuth2 Client Credentials с короткоживущими токенами или mTLS (не выдавайте бессрочные/долгоживущие JWT) |
> | Мобильное приложение | JWT в памяти (не в localStorage!) + Refresh в secure storage |

---

---

## 13. Работа с Базами Данных (PostgreSQL & `pgxpool`)

> 🔗 **Практика и примеры кода:** [15-sql](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/15-sql), [16-db-apps](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/16-db-apps) | 🛠️ **Домашние задания:** [homework-15 (Схема SQL)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/15-sql/homework-15), [homework-16 (Пакет БД фильмов)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/16-db-apps/homework-16) | 🎯 **Задачи:** [homework-tasks.md (Урок 12: PostgreSQL)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/homework-tasks.md#урок-12-postgresql-pgxpool-acid-транзакции-и-миграции), [homework-tasks.md (ДЗ 15–16)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/homework-tasks.md#домашнее-задание-15-реляционные-базы-данных-схема-бд-онлайн-кинотеатра)

> [!TIP]
> 🚀 **Полный справочник по БД для подготовки к собеседованиям:**  
> Для глубокого погружения в устройство СУБД, бинарный протокол pgx, MVCC, аномалии транзакций, распределенные блокировки Redis, паттерн Outbox и топ-30 каверзных вопросов уровня Junior-Senior читайте отдельный расширенный гайд:  
> 👉 **[go-database-interview-guide.md (Базы Данных для Go-разработчика: от архитектуры пулов до Senior-вопросов)](go-database-interview-guide.md)**

### 13.1. Пул соединений (`pgxpool.Pool`)

> [!IMPORTANT]
> **Почему нельзя открывать одно соединение на запрос?**  
> Создание TCP + TLS-хендшейка с базой данных — крайне тяжелая операция. Пул соединений (`pgxpool.Pool`) держит готовые открытые коннекты и выдаёт их горутинам на время запроса.

```go
import "github.com/jackc/pgx/v5/pgxpool"

// Строку подключения (DSN) в реальном проекте берут из переменной окружения; sslmode=disable — только для локальной разработки
dbPool, err := pgxpool.New(ctx, "postgres://user:pass@localhost:5432/dbname?sslmode=disable")
if err != nil {
    log.Fatal(err)
}
defer dbPool.Close()

// Проверка доступности БД (health check):
if err := dbPool.Ping(ctx); err != nil {
    log.Fatal("База недоступна:", err)
}
```

Размером пула управляют через конфигурацию: `cfg, _ := pgxpool.ParseConfig(dsn); cfg.MaxConns = 20; pool, _ := pgxpool.NewWithConfig(ctx, cfg)`. Число соединений подбирают с учётом лимита `max_connections` самой СУБД и числа экземпляров вашего сервиса (по умолчанию `MaxConns` = число ядер, но не меньше 4).

---

### 13.2. Запросы, Параметризация и обработка `ErrNoRows`

```go
// 1. Параметризованный запрос ($1, $2 предотвращают SQL-инъекции):
query := `SELECT id, title, rating FROM movies WHERE id = $1`

var m Movie
err := dbPool.QueryRow(ctx, query, 42).Scan(&m.ID, &m.Title, &m.Rating)

if errors.Is(err, pgx.ErrNoRows) {
    // Запись не найдена — это штатная ситуация, мапим в доменную ошибку ErrNotFound!
    return nil, ErrNotFound
}
if err != nil {
    return nil, fmt.Errorf("query movie failed: %w", err)
}
```

Для выборки **нескольких** строк используйте `Query` и обязательно закрывайте `rows` и проверяйте `rows.Err()`:

```go
rows, err := dbPool.Query(ctx, `SELECT id, title, rating FROM movies WHERE rating > $1`, 8.0)
if err != nil {
    return nil, err
}
// pgx.CollectRows читает все строки, закрывает rows и возвращает ошибку итерации
movies, err := pgx.CollectRows(rows, pgx.RowToStructByName[Movie])
```

---

### 13.3. Транзакции (`ACID Transactions` в Go)

Каноничный паттерн выполнения транзакции:

```go
tx, err := dbPool.Begin(ctx)
if err != nil {
    return err
}
// defer tx.Rollback безопасен: если дойдет до tx.Commit(), Rollback вернёт pgx.ErrTxClosed и ничего не сделает
defer tx.Rollback(ctx)

// 1. Списание баланса
_, err = tx.Exec(ctx, "UPDATE accounts SET balance = balance - $1 WHERE id = $2", 100, fromID)
if err != nil {
    return err
}

// 2. Пополнение баланса
_, err = tx.Exec(ctx, "UPDATE accounts SET balance = balance + $1 WHERE id = $2", 100, toID)
if err != nil {
    return err
}

// 3. Фиксация транзакции:
return tx.Commit(ctx)
```

#### Уровни изоляции транзакций и Аномалии (Вопрос с собеседований):

> [!NOTE]
> Таблица показывает требования стандарта SQL. **В PostgreSQL** уровень `Read Uncommitted` работает как `Read Committed` (грязное чтение невозможно ни на каком уровне), а `Repeatable Read` — это snapshot isolation, при котором фантомов нет, но возможна аномалия *write skew*; полную защиту даёт только `Serializable`, но приложение обязано быть готово **повторить транзакцию** при ошибке сериализации (SQLSTATE `40001`).

| Уровень изоляции                         | Грязное чтение (_Dirty Read_) | Неповторяемое чтение (_Non-Repeatable Read_) |      Фантомное чтение (_Phantom Read_)      |
| :--------------------------------------- | :---------------------------: | :------------------------------------------: | :-----------------------------------------: |
| **Read Uncommitted**                     |          ❌ Возможен          |                 ❌ Возможен                  |                 ❌ Возможен                 |
| **Read Committed** (_дефолт в Postgres_) |          ✅ Защищен           |                 ❌ Возможен                  |                 ❌ Возможен                 |
| **Repeatable Read**                      |          ✅ Защищен           |                  ✅ Защищен                  | ❌ Возможен (_в PG защищен благодаря MVCC_) |
| **Serializable**                         |          ✅ Защищен           |                  ✅ Защищен                  |                 ✅ Защищен                  |

#### Пессимистические и Оптимистические блокировки:

- **Пессимистическая блокировка (`SELECT ... FOR UPDATE`):**
  Блокирует строки на уровне БД до завершения транзакции (`Commit`/`Rollback`). Защищает от состояния гонки (_Lost Update_) при одновременных списаниях:
  ```sql
  SELECT balance FROM accounts WHERE id = $1 FOR UPDATE;
  ```
- **Оптимистическая блокировка (через `version` / `updated_at`):**
  Без блокировки строк. При обновлении проверяется, не изменилась ли версия:
  ```sql
  UPDATE accounts SET balance = 500, version = version + 1 WHERE id = 1 AND version = 3;
  -- Если возвращено 0 затронутых строк (RowsAffected == 0), значит была конкурентная модификация -> повторяем попытку (Retry).
  ```

---

### 13.4. Продвинутый SQL для собеседований: Индексы, `EXPLAIN ANALYZE`, CTE и Оконные функции

#### 1. Индексы в PostgreSQL (B-Tree, HASH, GIN, Partial, Covering)

- **B-Tree (дефолт):** Сбалансированное дерево поиска ($O(\log N)$). Идеально для операций `=`, `<`, `>`, `BETWEEN`, `IN`, `ORDER BY`.
- **HASH:** Поиск $O(1)$ строго по равенству (`=`). Хранит 4-байтный хеш, поэтому компактнее B-Tree для длинных строк (UUID, токены, URL). Не умеет диапазоны, сортировку, составные ключи и Index Only Scan.
- **GIN (Generalized Inverted Index):** Для поиска внутри составных структур: `JSONB`, массивов `ARRAY`, полнотекстового поиска (FTS).
- **Составной индекс (Composite Index) и правило левого префикса:**
  - Индекс `CREATE INDEX idx_user_status ON users(status, created_at);`
  - Будет работать для: `WHERE status = 'active'` и `WHERE status = 'active' AND created_at > NOW()`.
  - **Как правило, не будет эффективно использоваться** для: `WHERE created_at > NOW()` (пропущен ведущий столбец `status`; в новых версиях PostgreSQL для отдельных случаев есть «skip scan», но рассчитывать на него не стоит).
- **Частичный индекс (Partial Index):** Экономит память и диск, индексируя только нужные строки:
  ```sql
  CREATE INDEX idx_active_orders ON orders(user_id) WHERE status = 'active';
  ```
- **Покрывающий индекс (Covering Index с `INCLUDE`):** Позволяет выполнить **Index Only Scan** без обращения к основной таблице (Heap):
  ```sql
  CREATE INDEX idx_users_email_name ON users(email) INCLUDE (full_name);
  ```

#### 2. Анализ производительности запросов: `EXPLAIN ANALYZE`

- `EXPLAIN` — показывает предположительный план выполнения (Cost-based Optimizer).
- `EXPLAIN ANALYZE` — **реально выполняет запрос** и возвращает точное время выполнения (`actual time`), количество строк и задействованные узлы. ⚠️ Для `INSERT`/`UPDATE`/`DELETE` это значит, что данные **изменятся**: выполняйте такие проверки внутри `BEGIN; ... ROLLBACK;`.
- **Ключевые типы сканирования:**
  - `Seq Scan` (Sequential Scan) — полное сканирование всей таблицы (плохо на больших объемах данных, нужен индекс).
  - `Index Scan` — поиск по индексу + чтение строки из таблицы (Heap).
  - `Index Only Scan` — идеальный случай: все нужные колонки взяты прямо из самого индекса (эффективен, если карта видимости таблицы актуальна — за это отвечает `VACUUM`).
  - `Bitmap Index Scan / Bitmap Heap Scan` — поиск множества страниц по индексу с последующим их чтением пакетом.

#### 3. CTE (Common Table Expressions) и Рекурсивные запросы

Позволяют разбивать сложные многоуровневые запросы на читаемые логические блоки (а с `WITH RECURSIVE` — обходить деревья и иерархии):

```sql
WITH top_customers AS (
    SELECT user_id, SUM(amount) AS total_spent
    FROM orders
    WHERE status = 'completed'
    GROUP BY user_id
    HAVING SUM(amount) > 10000
)
SELECT u.id, u.email, tc.total_spent
FROM top_customers tc
JOIN users u ON u.id = tc.user_id
ORDER BY tc.total_spent DESC;
```

Пример рекурсивного CTE — обход иерархии сотрудников (каждый сотрудник хранит `manager_id`):

```sql
WITH RECURSIVE subordinates AS (
    SELECT id, name, manager_id, 1 AS depth FROM employees WHERE id = 1   -- якорь: начальник
    UNION ALL
    SELECT e.id, e.name, e.manager_id, s.depth + 1                        -- рекурсивный шаг
    FROM employees e
    JOIN subordinates s ON e.manager_id = s.id
)
SELECT * FROM subordinates;
```

#### 4. Оконные функции (Window Functions)

Выполняют вычисления над набором строк, связанных с текущей строкой, **без схлопывания (GROUP BY) общего числа строк**:

- **`ROW_NUMBER() OVER (PARTITION BY category_id ORDER BY price DESC)`**: Нумерация строк внутри каждой категории.
- **`DENSE_RANK() / RANK()`**: Ранжирование с учетом одинаковых значений.
- **`LAG() / LEAD()`**: Доступ к предыдущей / следующей строке (удобно для подсчета разницы по времени между событиями).

```sql
-- Пример: Получить Топ-3 самых дорогих товаров в каждой категории:
WITH ranked_products AS (
    SELECT
        id,
        name,
        category_id,
        price,
        ROW_NUMBER() OVER (PARTITION BY category_id ORDER BY price DESC) as rank
    FROM products
)
SELECT id, name, category_id, price
FROM ranked_products
WHERE rank <= 3;
```

---

### 13.5. Миграции Баз Данных: `golang-migrate` vs `goose` и генераторы кода

> [!IMPORTANT]
> **Зачем нужны миграции?**  
> Ручное создание таблиц через GUI (DBeaver, DataGrip) приводит к неконсистентности между окружениями (local, dev, prod). Миграции версионируют схему БД в Git и автоматически применяются при старте сервиса.

#### Сравнение популярных инструментов:

- **`golang-migrate`**: Требует 2 отдельных файла на каждую миграцию: `000001_create_users.up.sql` (накатить) и `000001_create_users.down.sql` (откатить).
- **`goose`**: Объединяет секции `-- +goose Up` и `-- +goose Down` в **один `.sql` файл**, а также поддерживает миграции на чистом Go-коде (удобно для сложного переноса данных).

#### Подходы к написанию SQL в Go:

1. **Raw SQL (`pgx`)**: Максимальный контроль над планами запросов (EXPLAIN ANALYZE) и наивысшая скорость.
2. **Query Builder (`Squirrel`)**: Типобезопасная сборка динамических запросов с защитой от SQL-инъекций (`PlaceholderFormat(sq.Dollar)`).
3. **`sqlc` (Type-safe SQL)**: Вы пишете чистые SQL-запросы в файлах `.sql`, а компилятор `sqlc` автоматически генерирует строго типизированные Go-структуры и методы (Zero Boilerplate).
4. **ORM (`GORM`)**: Удобно для простых CRUD, но скрывает реальные SQL-запросы, медленнее работает и трудно оптимизируется на высоких нагрузках (на собеседованиях часто ценят опыт чистого SQL выше ORM).

---


---

## 28. Сетевое Программирование: TCP и UDP сокеты (Пакет `net`)

> 🔗 **Практика и примеры кода:** [11-network](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/11-network) | 🛠️ **Домашнее задание:** [homework-11](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/11-network/homework-11) (Сетевая служба GoSearch) | 🎯 **Задачи:** [homework-tasks.md (ДЗ 11)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/homework-tasks.md#домашнее-задание-11-сетевое-программирование-сетевая-служба-gosearch)



Вся веб-разработка (`net/http`, `gRPC`, драйверы баз данных) под капотом построена на сокетах стандартного пакета `net`.

### 28.0. Сетевые модели (OSI vs TCP/IP), Сокеты и Принципы транспорта

> [!NOTE]
> 📚 **Рекомендуемая литература:**  
> Книга: *«Network Programming with Go»*, Jan Newmarch.

#### 1. Модель OSI vs Стек протоколов TCP/IP

В теории сети описываются 7-уровневой эталонной моделью **OSI**, однако реальный интернет и операционные системы работают по 4-уровневой практической модели **TCP/IP**:

| Уровень OSI (Теория) | Стек TCP/IP (Практика) | Что делает на практике | Протоколы |
| :--- | :--- | :--- | :--- |
| **7. Приложений**<br>**6. Представления**<br>**5. Сеансовый** | **Уровень приложений (Application)** | Бизнес-логика, форматы данных, шифрование | HTTP, HTTPS, gRPC, DNS, SMTP, FTP |
| **4. Транспортный** | **Транспортный уровень (Transport)** | Доставка данных между процессами на узлах | **TCP** (надежный поток), **UDP** (дейтаграммы) |
| **3. Сетевой** | **Межсетевой уровень (Internet)** | Маршрутизация пакетов между сетями | IPv4, IPv6, ICMP, IPsec |
| **2. Канальный**<br>**1. Физический** | **Уровень сетевого доступа (Link)** | Физическая адресация (MAC), передача сигналов по кабелю/радио | Ethernet, Wi-Fi (802.11), Оптика |

> [!TIP]
> Для прикладного разработчика на Go ключевое значение имеют **Транспортный уровень (`net.Conn`, TCP, UDP)** и **Уровень приложений (`net/http`, gRPC)**.

#### 2. Сетевой сокет, адресация и порты

* **Сетевой сокет (Socket)** — это программная абстракция конечной точки сетевого соединения, однозначно идентифицируемая парой:  
  $$\text{Сокет} = \text{IP-адрес} : \text{Порт} \quad (\text{например, } 192.168.1.10:8000)$$
  Само TCP-соединение однозначно определяется парой сокетов (адрес и порт клиента + адрес и порт сервера), поэтому к одному порту сервера могут одновременно подключаться тысячи клиентов.
* **Номера портов (0–65535):**
  * **0–1023 (Зарезервированные / Well-Known):** требуют привилегий суперпользователя (`root`/администратора): 80 (HTTP), 443 (HTTPS), 22 (SSH), 13 (Daytime RFC 867).
  * **1024–49151 (Зарегистрированные порты):** используются пользовательскими сервисами и СУБД: 5432 (PostgreSQL), 6379 (Redis), 8000/8080 (Web apps).
  * **49152–65535 (Динамические / Ephemeral):** выделяются ОС для исходящих клиентских соединений (диапазон по рекомендации IANA; в Linux по умолчанию 32768–60999).
* **Адрес привязки (Binding IP):**
  * `0.0.0.0` (или пустая строка `":8000"` в Go) — слушать входящие подключения **на всех сетевых интерфейсах** машины (Loopback, LAN, WAN).
  * `127.0.0.1` (`localhost`) — служба доступна **только локально** на текущей машине.

---

### 28.1. TCP-сервер и TCP-клиент на сокетах

**TCP (Transmission Control Protocol)** — надежный, дуплексный, ориентированный на соединение протокол:
1. **3-Way Handshake:** Установка соединения (`SYN` $\rightarrow$ `SYN-ACK` $\rightarrow$ `ACK`).
2. **Гарантия доставки и порядка:** Пакеты нумеруются, потерянные сегменты запрашиваются повторно, дубликаты отсекаются.
3. **Потоковая передача (Byte Stream):** В TCP **нет концепции сообщений**! Данные передаются как непрерывный поток байт. Если клиент отправил 2 сообщения подряд, сервер может прочитать их как один кусок или как пять мелких фрагментов. Для разделения потока на сообщения требуется протокол прикладного уровня (например, символ переноса строки `\n` или префиксная длина сообщения).

#### 1. Ключевые интерфейсы пакета `net`:

```go
// net.Listener — интерфейс слушающего сервера:
type Listener interface {
    Accept() (Conn, error) // Ожидает и возвращает новое входящее соединение
    Close() error          // Закрывает слушающий сокет
    Addr() Addr            // Возвращает сетевой адрес слушателя
}

// net.Conn — абстракция сетевого подключения (реализует io.Reader, io.Writer, io.Closer):
type Conn interface {
    Read(b []byte) (n int, err error)
    Write(b []byte) (n int, err error)
    Close() error
    LocalAddr() Addr
    RemoteAddr() Addr
    SetDeadline(t time.Time) error
    SetReadDeadline(t time.Time) error
    SetWriteDeadline(t time.Time) error
}
```

#### 2. Многопоточный TCP-сервер (Echo Server):

```go
package main

import (
    "bufio"
    "errors"
    "fmt"
    "log"
    "net"
)

func handleConnection(conn net.Conn) {
    defer conn.Close()
    scanner := bufio.NewScanner(conn)

    for scanner.Scan() {
        text := scanner.Text()
        fmt.Fprintf(conn, "ECHO: %s\n", text) // отправляем ответ клиенту
    }
}

func main() {
    listener, err := net.Listen("tcp", ":9000")
    if err != nil {
        log.Fatal(err)
    }
    defer listener.Close()
    fmt.Println("TCP сервер слушает порт 9000...")

    for {
        conn, err := listener.Accept() // блокируется, пока не подключится клиент
        if err != nil {
            if errors.Is(err, net.ErrClosed) {
                return // слушающий сокет закрыт (например, при остановке сервера)
            }
            log.Println("accept:", err) // временные ошибки (например, too many open files) — логируем и пробуем снова
            continue
        }
        go handleConnection(conn) // обрабатываем каждого клиента в своей горутине
    }
}
```

#### 3. TCP-клиент:

```go
// Всегда задавайте таймаут установления соединения (net.Dial без таймаута может ждать очень долго):
conn, err := net.DialTimeout("tcp", "localhost:9000", 5*time.Second)
if err != nil {
    log.Fatal(err)
}
defer conn.Close()

fmt.Fprintf(conn, "Привет от клиента!\n")
response, _ := bufio.NewReader(conn).ReadString('\n')
fmt.Print("Ответ сервера: ", response)
```

---

### 28.2. UDP: Дейтаграммы и отличие от TCP

- **UDP (User Datagram Protocol):** Протокол без установки соединения (Connectionless). Пакеты отправляются без гарантии доставки и порядка, но без задержек на рукопожатия и повторную пересылку.
- **Где применяется на проде:** DNS-запросы (53 порт), видеостриминг, голосовая связь (VoIP), онлайн-игры, метрики StatsD.
- **Границы сообщений сохраняются:** один `WriteTo` = одна дейтаграмма = одно `ReadFrom` (в отличие от потока TCP). Дейтаграмму стараются держать в пределах ~1400 байт, чтобы она не фрагментировалась на уровне IP; буфер чтения меньше размера дейтаграммы «обрежет» её остаток.

```go
// UDP сервер:
packetConn, err := net.ListenPacket("udp", ":9001")
if err != nil {
    log.Fatal(err)
}
defer packetConn.Close()

buf := make([]byte, 1024)
n, clientAddr, err := packetConn.ReadFrom(buf) // чтение дейтаграммы
packetConn.WriteTo([]byte("ACK"), clientAddr) // ответ клиенту
```

---

### 28.3. Таймауты на сокетах (`SetDeadline`)

> [!WARNING]
> Если не выставлять таймауты на `net.Conn`, «зависший» или медленный клиент может навечно заблокировать горутину на чтении (`Read`), вызвав **утечку горутин (Goroutine Leak)**.

```go
// Устанавливаем общий дедлайн (и на чтение, и на запись) на 5 секунд:
conn.SetDeadline(time.Now().Add(5 * time.Second))

// Или только на чтение входящих данных:
conn.SetReadDeadline(time.Now().Add(3 * time.Second))
```

> ⚠️ Дедлайн — это **абсолютный момент времени**, а не «таймаут бездействия»: после его наступления все последующие `Read`/`Write` будут сразу возвращать ошибку (`os.ErrDeadlineExceeded`). В цикле чтения дедлайн нужно **обновлять перед каждой операцией**, если вы хотите ограничить время ожидания именно каждого сообщения.

---

### 28.4. WebSockets: Полнодуплексный обмен в реальном времени

В отличие от стандартного HTTP (модель «запрос-ответ»), **WebSocket** обеспечивает постоянное двунаправленное соединение поверх одного TCP-сокета:

1. **HTTP Upgrade Handshake:** Клиент отправляет заголовок `Upgrade: websocket`. Сервер отвечает `101 Switching Protocols`.
2. **Паттерн Read/Write Pump:** Для каждого соединения на сервере выделяются 2 независимые горутины:
   - **ReadPump:** читает входящие сообщения сокета и передает в центральный хаб.
   - **WritePump:** забирает сообщения из очереди исходящих и отправляет в сокет.
3. Популярная в Go библиотека: `github.com/coder/websocket` (ранее `nhooyr.io/websocket`) или `github.com/gorilla/websocket`.

#### 💻 Практический пример: Чат-сервер на WebSocket (`gorilla/websocket`)

```go
package main

import (
    "log"
    "net/http"
    "time"

    "github.com/gorilla/websocket"
)

const (
    writeWait  = 10 * time.Second    // Таймаут на запись одного сообщения
    pongWait   = 60 * time.Second    // Таймаут ожидания Pong от клиента
    pingPeriod = (pongWait * 9) / 10 // Период отправки Ping (< pongWait)
    maxMsgSize = 512                  // Максимальный размер входящего сообщения (байт)
)

// upgrader конвертирует HTTP соединение в WebSocket.
// ⚠️ CheckOrigin: true разрешает подключения с ЛЮБОГО Origin — это упрощение для примера и уязвимость
// (Cross-Site WebSocket Hijacking). В продакшене сверяйте r.Header.Get("Origin") со списком разрешённых доменов.
var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin:     func(r *http.Request) bool { return true },
}

// Client — представляет одно WebSocket соединение на сервере.
type Client struct {
    hub  *Hub
    conn *websocket.Conn
    send chan []byte // Буферизованный канал исходящих сообщений
}

// Hub — центральный диспетчер: хранит список клиентов и рассылает сообщения.
type Hub struct {
    clients    map[*Client]bool // Активные подключения
    broadcast  chan []byte      // Входящие сообщения для рассылки всем
    register   chan *Client     // Регистрация нового клиента
    unregister chan *Client     // Отключение клиента
}

func newHub() *Hub {
    return &Hub{
        broadcast:  make(chan []byte),
        register:   make(chan *Client),
        unregister: make(chan *Client),
        clients:    make(map[*Client]bool),
    }
}

// Run — основной цикл хаба (запускается в отдельной горутине).
func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.clients[client] = true
        case client := <-h.unregister:
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                close(client.send) // Сигнал WritePump завершить работу
            }
        case message := <-h.broadcast:
            for client := range h.clients {
                select {
                case client.send <- message:
                default:
                    // Буфер переполнен — клиент слишком медленный, отключаем
                    close(client.send)
                    delete(h.clients, client)
                }
            }
        }
    }
}

// ReadPump — читает сообщения из WebSocket соединения и публикует в хаб.
// Запускается в отдельной горутине для каждого клиента.
func (c *Client) ReadPump() {
    defer func() {
        c.hub.unregister <- c
        c.conn.Close()
    }()
    c.conn.SetReadLimit(maxMsgSize)
    c.conn.SetReadDeadline(time.Now().Add(pongWait))
    // Обновляем дедлайн при каждом получении Pong от клиента:
    c.conn.SetPongHandler(func(string) error {
        c.conn.SetReadDeadline(time.Now().Add(pongWait))
        return nil
    })
    for {
        _, message, err := c.conn.ReadMessage()
        if err != nil {
            // websocket.IsUnexpectedCloseError — проверяет штатное закрытие
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("ws read error: %v", err)
            }
            break
        }
        c.hub.broadcast <- message
    }
}

// WritePump — читает из канала c.send и пишет в WebSocket соединение.
// Посылает Ping-сообщения для проверки живости соединения.
func (c *Client) WritePump() {
    ticker := time.NewTicker(pingPeriod)
    defer func() {
        ticker.Stop()
        c.conn.Close()
    }()
    for {
        select {
        case message, ok := <-c.send:
            c.conn.SetWriteDeadline(time.Now().Add(writeWait))
            if !ok {
                // Хаб закрыл канал — отправляем Close Frame клиенту
                c.conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }
            w, err := c.conn.NextWriter(websocket.TextMessage)
            if err != nil {
                return
            }
            w.Write(message)
            // Отправляем накопленные в очереди сообщения пакетом (batching):
            n := len(c.send)
            for i := 0; i < n; i++ {
                w.Write([]byte{'\n'})
                w.Write(<-c.send)
            }
            if err := w.Close(); err != nil {
                return
            }
        case <-ticker.C:
            // Периодический Ping для поддержания соединения живым:
            c.conn.SetWriteDeadline(time.Now().Add(writeWait))
            if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return // Клиент недоступен — завершаем WritePump
            }
        }
    }
}

// serveWs — HTTP хендлер, апгрейдящий соединение до WebSocket.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("upgrade error:", err)
        return
    }
    client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
    client.hub.register <- client

    // Запускаем Pump-горутины для этого клиента:
    go client.WritePump()
    go client.ReadPump()
}

func main() {
    hub := newHub()
    go hub.Run() // Хаб работает в фоне весь жизненный цикл приложения

    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        serveWs(hub, w, r)
    })
    log.Println("WebSocket сервер запущен на :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

> [!NOTE]
> **Разница между `coder/websocket` и `gorilla/websocket`:**
> - `gorilla/websocket` — классический выбор, очень зрелый (проект Gorilla был архивирован в 2022 году и затем возобновлён сообществом; API стабилен). Соединение допускает **одного** конкурентного писателя и одного читателя — поэтому в примере запись идёт только из `WritePump`.
> - `coder/websocket` (форк Nhooyr) — современная альтернатива: context-aware, поддерживает WASM, лучше интегрируется с `net/http` стандарта Go 1.22+.

---



### 28.5. Безопасность сетевых соединений: TLS и mTLS в Go

Безопасность транспортного уровня в Go реализуется пакетом `crypto/tls`.

#### Что такое mTLS (Mutual TLS / Взаимный TLS)?

- При стандартном HTTPS клиент проверяет сертификат сервера.
- При **mTLS** и клиент, и сервер **взаимно проверяют сертификаты друг друга** с помощью доверенного корневого центра (Root CA). Это золотой стандарт безопасности в межсервисном взаимодействии (Zero-Trust Service Mesh).

#### Настройка mTLS-сервера в Go:

```go
package main

import (
    "crypto/tls"
    "crypto/x509"
    "net/http"
    "os"
)

func createMTLSServer() (*http.Server, error) {
    // 1. Загружаем сертификат доверенного CA для проверки клиентов:
    caCert, err := os.ReadFile("ca.crt")
    if err != nil {
        return nil, err
    }
    caCertPool := x509.NewCertPool()
    caCertPool.AppendCertsFromPEM(caCert)

    // 2. Конфигурируем TLS с обязательной проверкой сертификата клиента:
    tlsConfig := &tls.Config{
        ClientCAs:  caCertPool,
        ClientAuth: tls.RequireAndVerifyClientCert, // ⚠️ Строгая проверка клиента!
        MinVersion: tls.VersionTLS13,               // Только современный TLS 1.3
    }

    return &http.Server{
        Addr:      ":8443",
        TLSConfig: tlsConfig,
    }, nil
}

// Запуск: собственный сертификат сервера всё равно нужен:
//   srv.ListenAndServeTLS("server.crt", "server.key")
```

Клиентская сторона mTLS: клиент предъявляет свой сертификат и проверяет сертификат сервера по тому же корневому CA:

```go
cert, _ := tls.LoadX509KeyPair("client.crt", "client.key")
tlsCfg := &tls.Config{
    Certificates: []tls.Certificate{cert}, // сертификат клиента
    RootCAs:      caCertPool,              // доверенный CA, подписавший сертификат сервера
    MinVersion:   tls.VersionTLS13,
}
client := &http.Client{Transport: &http.Transport{TLSClientConfig: tlsCfg}}
```

---

### 28.6. Тестирование сетевых служб в памяти: `net.Pipe()`

Стандартный подход к тестированию сетевых сервисов через поднятие реального TCP-порта (`127.0.0.1:8000`) порождает ряд проблем:
- **Конфликт портов (Flaky tests):** ошибка `bind: address already in use` при параллельном запуске тестов (`go test -p`).
- **Сетевые ограничения:** фаерволы CI/CD и требования прав администратора.
- **Накладные расходы:** задержки на системные вызовы ядра ОС и сетевой стек.

Стандартная библиотека Go решает это функцией **`net.Pipe()`**:

```go
// Создает синхронное дуплексное соединение в оперативной памяти:
serverConn, clientConn := net.Pipe()
```

Оба объекта реализуют полноценный интерфейс **`net.Conn`**. То, что пишется в `clientConn`, считывается из `serverConn`, и наоборот — без участия сетевой карты и ОС!

> ⚠️ `net.Pipe` **синхронный и без буфера**: `Write` блокируется, пока другая сторона не прочитает данные. Поэтому обработчик сервера нужно запускать в отдельной горутине (как в примере ниже), иначе тест зависнет (deadlock).

#### 💡 Идиоматичный паттерн тестирования сетевого обработчика:

```go
func TestEchoHandler(t *testing.T) {
    serverConn, clientConn := net.Pipe()
    defer clientConn.Close()

    // 1. Запускаем обработчик сервера в отдельной горутине:
    go func() {
        handleConnection(serverConn) // функция сервера закрывает serverConn при выходе
    }()

    // 2. В основном потоке теста эмулируем действия клиента:
    message := "Привет, Go!\n"
    if _, err := clientConn.Write([]byte(message)); err != nil {
        t.Fatalf("ошибка отправки: %v", err)
    }

    // 3. Читаем ответ сервера:
    reader := bufio.NewReader(clientConn)
    response, err := reader.ReadString('\n')
    if err != nil {
        t.Fatalf("ошибка чтения ответа: %v", err)
    }

    expected := "ECHO: Привет, Go!\n"
    if response != expected {
        t.Errorf("Ожидалось %q, получено %q", expected, response)
    }
}
```

---

## 29. Работа со временем: пакет `time`

### 29.1. Форматирование и парсинг дат (Референсное время Go)

> [!WARNING]
> **Ловушка для новичков:** В Go используется уникальный подход к форматированию дат — **референсное время** `Mon Jan 2 15:04:05 MST 2006` (мнемоника: `01/02 03:04:05PM '06 -0700`). Это НЕ шаблон `YYYY-MM-DD`, как в других языках!

```go
import "time"

now := time.Now()

// Форматирование:
fmt.Println(now.Format("2006-01-02 15:04:05")) // "2026-09-08 22:10:48"
fmt.Println(now.Format(time.RFC3339))            // "2026-09-08T22:10:48+03:00"
fmt.Println(now.Format("02.01.2006"))            // "08.09.2026" (европейский формат)

// Парсинг строки в time.Time (без указания зоны результат — в UTC):
t, err := time.Parse("2006-01-02", "2024-12-25")
if err != nil {
    log.Fatal(err)
}

// Парсинг с учетом таймзоны:
loc, _ := time.LoadLocation("Europe/Moscow")
t, err = time.ParseInLocation("2006-01-02 15:04", "2024-12-25 18:00", loc)
```

> [!TIP]
> - `time.LoadLocation` читает базу часовых поясов ОС. В минимальных Docker-образах (`scratch`, `alpine` без tzdata) она отсутствует — тогда вернётся ошибка. Решение: добавьте `import _ "time/tzdata"` (встроит базу в бинарник, +~450 КБ) или установите пакет `tzdata` в образ.
> - Храните и передавайте время в **UTC** (`time.Now().UTC()`), переводя в локальную зону только для показа пользователю.
> - `time.Duration` — это `int64` наносекунд: `time.Sleep(5)` спит 5 **наносекунд**, а не 5 секунд. Пишите `5 * time.Second`. Если количество секунд лежит в переменной: `time.Duration(n) * time.Second`.

### 29.2. Измерение времени и таймеры

```go
// 1. Замер длительности операции:
start := time.Now()
doExpensiveWork()
elapsed := time.Since(start) // эквивалент time.Now().Sub(start)
fmt.Printf("Операция заняла: %s\n", elapsed) // "Операция заняла: 1.234s"

// 2. time.Ticker — периодическое срабатывание (замена while + sleep):
ticker := time.NewTicker(5 * time.Second)
defer ticker.Stop() // ОБЯЗАТЕЛЬНО для go.mod < 1.23 (иначе утечка ресурсов), рекомендуется всегда

for range ticker.C {
    fmt.Println("Тик каждые 5 секунд")
}

// 3. time.Timer — одноразовое срабатывание через интервал:
timer := time.NewTimer(10 * time.Second)
<-timer.C // блокируется до срабатывания
```

> [!NOTE]
> **Новинка Go 1.23: Сборка неактивных таймеров и тикеров сборщиком мусора**  
> До Go 1.23 забытый `time.NewTimer` или `time.NewTicker` без вызова `.Stop()` удерживался в памяти внутренним рантаймом до истечения времени.  
> Начиная с **Go 1.23**, сборщик мусора (GC) умеет автоматически собирать недостижимые таймеры и тикеры, даже если `.Stop()` не был вызван (новое поведение включено, если в `go.mod` указана версия `go 1.23` или выше). Однако явный вызов `defer ticker.Stop()` остается строгой рекомендацией для детерминированного освобождения ресурсов.

### 29.3. Ловушки при сравнении `time.Time`

> **Вопрос на собеседовании:** Почему `t1 == t2` может вернуть `false`, даже если время одинаковое?

```go
t1 := time.Now()
t2 := t1.Round(0) // убирает монотонный таймер

fmt.Println(t1 == t2)       // false (!) — разные внутренние представления
fmt.Println(t1.Equal(t2))   // true  ✅ — сравнивает именно время

// Правило: ВСЕГДА используйте .Equal() для сравнения time.Time
// Оператор == сравнивает и монотонный таймер, и wall clock, и таймзону
// Для упорядочивания (например, в slices.SortFunc) есть t1.Compare(t2) (Go 1.20+): -1, 0, +1
// Для «нулевого» времени используйте t.IsZero(), а не сравнение с time.Time{}
```

---

### 29.4. Каверзный вопрос собеседования: Всегда ли второй `time.Now()` больше первого? (Wall Clock vs Monotonic Clock)

> **Вопрос на собеседовании (Senior Go, 400k+):**  
> *«Мы взяли `t1 := time.Now()`, выполнили некоторую работу и взяли `t2 := time.Now()`. Всегда ли `t2` будет позже `t1`?»*

```go
t1 := time.Now()
doWork()
t2 := time.Now()

// 1. Строковый вывод (настенные часы):
fmt.Println(t1.Format(time.RFC3339))
fmt.Println(t2.Format(time.RFC3339))

// 2. Вычисление интервала (монотонные часы):
diff := t2.Sub(t1)
```

#### Эталонный ответ:

1. **Нет, если выводить настенное время (`wall clock`):**
   - Системные настенные часы операционной системы отображают текущее календарное время.
   - Они **могут быть переведены назад**:
     - Фоновой синхронизацией по протоколу NTP (Network Time Protocol) при коррекции расхождения часов сервера;
     - При ручном переводе системного времени системным администратором;
     - При наступлении високосной секунды (Leap Second).
   - В этом случае вывод `t2.Format(...)` или строковое представление `t2.String()` может показать момент времени **раньше**, чем `t1`!

2. **Да, если вычислять разницу (`t2.Sub(t1)` или `time.Since(t1)`):**
   - Начиная с Go 1.9, структура `time.Time` внутри содержит два независимых компонента:
     - `wall uint64` — настенное время (секунды и наносекунды от 1 января 1 года н.э.);
     - `ext int64` — **монотонные часы (Monotonic Clock)**, отсчитывающие тики процессора с момента старта ОС.
   - Монотонные часы **идут строго вперед** и физически не подвержены переводам системного времени операционной системы.
   - При вызове `t2.Sub(t1)` или `time.Since(t1)` рантайм Go проверяет наличие монотонной составляющей у обоих объектов и вычисляет разницу **по монотонным часам**. Поэтому `diff` гарантированно будет положительным и отразит реальную длительность выполнения `doWork()`.

3. **Когда монотонный таймер теряется:**
   - При сериализации времени (в JSON, protobuf, SQL-базу данных).
   - При явном вызове `t.Round(0)` (это идиоматичный способ отсечь монотонную составляющую, если нужно сравнить только календарное время).
   - Если один из объектов `time.Time` потерял монотонную составляющую, `t2.Sub(t1)` откатывается к вычитанию настенного времени и может вернуть отрицательное число при скачке системных часов!

> 💡 **Фраза для собеседования:** *«`time.Now()` возвращает и настенное время, и монотонные часы. Настенные часы могут быть переведены назад демоном NTP, поэтому строковый вывод `t2` может оказаться раньше `t1`. Но операция `t2.Sub(t1)` использует монотонные часы, которые идут только вперед, поэтому вычисленный интервал всегда будет неотрицательным и корректным».*

---

> [!TIP]
> 📚 **Навигация:** [⬅️ Назад: Том 3](03-go-core-concurrency.md) | [📖 Главное оглавление](README.md) | [Вперед: Том 5 (Тестирование и Архитектура) ➡️](05-go-core-testing-arch.md)

