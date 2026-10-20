# Домашнее задание №16: Приложения с базами данных (Пакет для БД фильмов с паттерном Repository)

Решение домашнего задания по курсу «Современная разработка на Go» (Занятие №16: Приложения баз данных).  
Реализовано в отдельной директории `homework-16` на основе схемы БД информационного сайта о фильмах из Занятия №15.

---

## 📑 Содержание

1. [Постановка задачи](#постановка-задачи)
2. [Задача №1: Интерфейс взаимодействия с БД](#задача-1-интерфейс-взаимодействия-с-бд)
3. [Задача №2: Пакет для PostgreSQL (pgsql)](#задача-2-пакет-для-postgresql-pgsql)
4. [Задача №3: Тестирование разработанных функций](#задача-3-тестирование-разработанных-функций)
5. [In-Memory хранилище (memsql)](#in-memory-хранилище-memsql)
6. [REST API и веб-сервер](#rest-api-и-веб-сервер)
7. [Инструкция по запуску и примеры запросов](#инструкция-по-запуску-и-примеры-запросов)

---

## Постановка задачи

Разработать пакет для работы с БД информационного сайта о фильмах.  
Схему БД взять из предыдущего задания (Занятие №15: таблицы `movies`, `studios`, `actors`, `directors`).

- **Задача №1**: Разработать интерфейс для взаимодействия с БД, позволяющий:
  - добавлять массив фильмов;
  - удалять фильм;
  - обновлять фильм;
  - получать массив фильмов (с возможностью указать ID студии для фильтра; если ID не указан или равен 0 — возвращать все фильмы).
- **Задача №2**: Разработать пакет для БД Postgres, выполняющий контракт интерфейса из Задачи №1.
- **Задача №3**: Написать тесты для разработанных функций.

---

## Задача №1: Интерфейс взаимодействия с БД

Файл [`pkg/db/db.go`](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/16-db-apps/homework-16/pkg/db/db.go) определяет доменные модели и контракт хранилища:

```go
type Movie struct {
    ID          int64  `json:"id"`
    Title       string `json:"title"`
    ReleaseYear int    `json:"release_year"`
    BoxOffice   int64  `json:"box_office"`
    Rating      string `json:"rating"` // 'PG-10', 'PG-13', 'PG-18'
    StudioID    int    `json:"studio_id"`
    StudioName  string `json:"studio_name,omitempty"` // заполняется при JOIN со студией
}

type Interface interface {
    // AddMovies добавляет массив фильмов в хранилище.
    AddMovies(ctx context.Context, movies []Movie) error
    // DeleteMovie удаляет фильм по его ID.
    DeleteMovie(ctx context.Context, id int64) error
    // UpdateMovie обновляет данные фильма.
    UpdateMovie(ctx context.Context, m Movie) error
    // Movies возвращает массив фильмов.
    // Если studioID == 0, возвращаются все фильмы.
    // Если studioID > 0, возвращаются фильмы только указанной студии.
    Movies(ctx context.Context, studioID int) ([]Movie, error)
    // MovieByID возвращает фильм по его ID.
    MovieByID(ctx context.Context, id int64) (Movie, error)
}
```

Файл [`pkg/db/errors.go`](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/16-db-apps/homework-16/pkg/db/errors.go) объявляет доменные ошибки (`db.ErrNotFound`), изолирующие бизнес-логику от драйверов СУБД.

---

## Задача №2: Пакет для PostgreSQL (pgsql)

Пакет [`pkg/db/pgsql/pgsql.go`](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/16-db-apps/homework-16/pkg/db/pgsql/pgsql.go) реализует интерфейс `db.Interface` для PostgreSQL:

- **Пул соединений**: на базе `github.com/jackc/pgx/v5/pgxpool` с поддержкой `Ping()` и таймаутов контекста.
- **Пакетная вставка (`AddMovies`)**: объединение вставки массива фильмов в одну транзакцию с использованием `pgx.Batch` (согласно слайду 21 презентации) для минимизации сетевых задержек.
- **Фильтрация по студии (`Movies`)**: SQL-запрос с `LEFT JOIN studios` и предикатом `WHERE ($1 = 0 OR m.studio_id = $1)` — если `studioID == 0`, возвращаются фильмы всех студий.
- **Обновление и удаление (`UpdateMovie`, `DeleteMovie`)**: контроль `RowsAffected() == 0` с возвратом `db.ErrNotFound`.

---

## Задача №3: Тестирование разработанных функций

Реализованы тесты:

1. **Модульные тесты хранилища** ([`pkg/db/memsql/memsql_test.go`](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/16-db-apps/homework-16/pkg/db/memsql/memsql_test.go)):
   - `TestAddMovies_And_Movies_All`: добавление массива фильмов и выборка всех (`studioID == 0`).
   - `TestMovies_FilterByStudio`: проверка фильтрации по конкретной студии (`studioID > 0`) и по несуществующей студии.
   - `TestUpdateMovie`: успешное обновление фильма и валидация `db.ErrNotFound` при несуществующем ID.
   - `TestDeleteMovie`: удаление фильма, проверка повторного удаления и проверка отсутствия в списке.
   - `TestConcurrency`: параллельное чтение и запись 40 горутинами с флагом `-race`.

2. **Тесты соответствия контракту pgsql** ([`pkg/db/pgsql/pgsql_test.go`](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/16-db-apps/homework-16/pkg/db/pgsql/pgsql_test.go)):
   - Статическая проверка реализации интерфейса: `var _ db.Interface = (*pgsql.DB)(nil)`.
   - Обработка ошибок при недоступности СУБД.
   - Корректная обработка пустого батча.

3. **Тесты REST API** ([`pkg/api/api_test.go`](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/16-db-apps/homework-16/pkg/api/api_test.go)):
   - `TestHandleMovies_Empty` и `TestHandleMovies_WithFilterByStudio` (включая query param `?studio_id=X` и `?studio_id=0`).
   - `TestHandleMovieByID`, `TestHandleCreateMovies_BatchAndSingle`, `TestHandleUpdateMovie`, `TestHandleDeleteMovie`.

---

## In-Memory хранилище (memsql)

Пакет [`pkg/db/memsql/memsql.go`](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/16-db-apps/homework-16/pkg/db/memsql/memsql.go) реализует `db.Interface` в оперативной памяти с потокобезопасностью на базе `sync.RWMutex`. Используется для мгновенных модульных тестов и демонстрационного запуска сервера без установленного PostgreSQL.

---

## REST API и веб-сервер

- Пакет [`pkg/api/api.go`](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/16-db-apps/homework-16/pkg/api/api.go) предоставляет HTTP-эндпоинты (на базе `gorilla/mux`):
  - `GET /api/v1/movies?studio_id=X` — получение фильмов (все или фильтр по студии).
  - `GET /api/v1/movies/{id}` — получение фильма по ID.
  - `POST /api/v1/movies` — добавление фильма или массива фильмов.
  - `PUT /api/v1/movies/{id}` — обновление фильма.
  - `DELETE /api/v1/movies/{id}` — удаление фильма.
- Приложение [`cmd/server/main.go`](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/16-db-apps/homework-16/cmd/server/main.go) запускает HTTP-сервер с таймаутами и Graceful Shutdown.

---

## Инструкция по запуску и примеры запросов

### Запуск тестов:

```bash
cd go_course_4
go test -v -race -cover ./lessons/16-db-apps/homework-16/...
```

### Запуск сервера:

```bash
go run ./lessons/16-db-apps/homework-16/cmd/server/main.go
```

### Примеры HTTP-запросов:

```bash
# 1. Получить все фильмы
curl http://localhost:8080/api/v1/movies | jq

# 2. Получить фильмы с фильтром по студии (studio_id = 1)
curl http://localhost:8080/api/v1/movies?studio_id=1 | jq

# 3. Добавить массив фильмов
curl -X POST http://localhost:8080/api/v1/movies \
  -H "Content-Type: application/json" \
  -d '[
    {"title": "The Matrix", "release_year": 1999, "box_office": 467200000, "rating": "PG-18", "studio_id": 1},
    {"title": "The Matrix Reloaded", "release_year": 2003, "box_office": 741800000, "rating": "PG-18", "studio_id": 1}
  ]' | jq

# 4. Обновить фильм
curl -X PUT http://localhost:8080/api/v1/movies/1 \
  -H "Content-Type: application/json" \
  -d '{"title": "Inception: Extended Cut", "release_year": 2010, "box_office": 850000000, "rating": "PG-13", "studio_id": 1}' | jq

# 5. Удалить фильм
curl -X DELETE http://localhost:8080/api/v1/movies/1
```
