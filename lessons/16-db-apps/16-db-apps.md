# 16-db-apps — Работа с базами данных из Go: database/sql, pgx и паттерн Репозиторий

В данном модуле рассматриваются практические способы взаимодействия Go-приложений с реляционными СУБД (PostgreSQL и MySQL): стандартный интерфейс `database/sql`, специализированный драйвер `pgx` (с пулом `pgxpool`), пакетная обработка запросов (`pgx.Batch`), управление транзакциями и реализация паттерна «Репозиторий» (Repository Pattern) для чистой архитектуры сервисов.

---

## 📋 Структура каталога

```text
16-db-apps/
├── 1-database-sql/
│   └── database-sql.go   # Стандартный пакет database/sql: sql.Open, пул соединений, Query, Exec, Tx
├── 2-pgx/
│   ├── pgx.go            # Высокопроизводительный драйвер pgx (pgxpool.Pool, контекст, pgx.Batch)
│   └── pgx_test.go       # Тест подключения и выборки с поддержкой graceful skip
├── pkg/                  # Пакеты учебной архитектуры репозиториев (Книги)
│   ├── api/              # HTTP API (api.go, api_test.go), зависящий от db.Interface
│   └── db/               # Интерфейс db.Interface, Book и реализации:
│       ├── memsql/       # In-memory потокобезопасная БД (memsql.go, memsql_test.go)
│       └── pgsql/        # Репозиторий на базе pgxpool (pgsql.go, pgsql_test.go)
├── schema.sql            # Схема учебной базы данных книг (DDL)
└── homework-16/          # Домашнее задание №16: Сервис фильмов и кинотеатра
    ├── README.md         # Подробное ТЗ, архитектурный разбор и примеры curl
    ├── cmd/server/       # Точка входа HTTP-сервера с Graceful Shutdown
    └── pkg/
        ├── api/          # REST API фильмов (CRUD, фильтрация по студиям, gorilla/mux)
        └── db/           # Интерфейс db.Interface (Movie), memsql и pgsql реализации
```

---

## 🔍 Ключевые концепции

1. **Пул соединений (`sql.DB` и `pgxpool.Pool`):**
   - Вызовы `sql.Open` или `pgxpool.New` не открывают сетевое соединение мгновенно, а создают пул потокобезопасных соединений.
   - Метод `Ping()` (или `PingContext()`) выполняет реальную проверку сетевой доступности СУБД.
   - Обязательно настраивать параметры пула: `SetMaxOpenConns`, `SetMaxIdleConns`, `SetConnMaxLifetime`.

2. **`database/sql` vs специализированный драйвер `pgx`:**
   - `database/sql` универсален для любых SQL-СУБД через единый интерфейс Go.
   - `pgx` использует нативный бинарный протокол PostgreSQL, поддерживает специфичные типы Postgres (массивы, JSONB, UUID), пакетные запросы (`pgx.Batch`) и имеет встроенный эффективный пул `pgxpool.Pool`.

3. **Защита от SQL Injection:**
   - **Никогда** не конкатенируйте строки в SQL-запросах (`fmt.Sprintf("SELECT ... WHERE id = %s", id)`).
   - Всегда используйте параметризованные плейсхолдеры (`$1, $2` в PostgreSQL или `?` в MySQL/SQLite).

4. **Освобождение ресурсов и безопасность транзакций:**
   - Всегда закрывайте `rows`: `defer rows.Close()` — это возвращает соединение в пул.
   - Проверяйте ошибки после цикла обхода строк: `if err := rows.Err(); err != nil`.
   - В транзакциях используйте паттерн `defer tx.Rollback()` сразу после `db.Begin()`: если произойдет ошибка или паника, откат выполнится автоматически; при успешном `tx.Commit()` повторный `Rollback` безопасен (идемпотентен).

5. **Паттерн «Репозиторий» (Repository Pattern):**
   - Бизнес-логика и обработчики HTTP API зависят исключительно от интерфейса `db.Interface`.
   - Это позволяет запускать быстрые юнит-тесты с `memsql` в памяти без внешней зависимости от СУБД, а в проде переключаться на `pgsql` простой сменой инициализации в `main.go`.

---

## 🚀 Запуск тестов и бенчмарков

```bash
# Запуск всех тестов модуля 16 с детектором гонок
go test -v -race ./lessons/16-db-apps/...

# Запуск с замером производительности (бенчмарки)
go test -v -bench=. ./lessons/16-db-apps/...
```

---

## 🔗 Файлы и связанные материалы
- [1-database-sql/database-sql.go](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/16-db-apps/1-database-sql/database-sql.go) — Пример работы со стандартным пакетом `database/sql`
- [2-pgx/pgx.go](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/16-db-apps/2-pgx/pgx.go) — Пример работы с `pgxpool` и `pgx.Batch`
- [schema.sql](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/16-db-apps/schema.sql) — Схема таблиц БД книг
- [homework-16/README.md](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/16-db-apps/homework-16/README.md) — Домашнее задание №16: Пакет для БД фильмов
- **Теория:** [Том 4: Backend и Web (Раздел 13: Go и СУБД)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/04-go-core-backend-web.md)
- **Гайд по базам данных:** [go-database-interview-guide.md](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/go-database-interview-guide.md)
