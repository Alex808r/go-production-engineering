# Домашнее задание №20: Итоговый проект «Lynks» (NoSQL, Redis, PostgreSQL)

Решение итогового проекта курса «Современная разработка на Go» (**Занятие №20: Нереляционные базы данных** / [`final.md`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lynks/final.md)).  
Реализовано в отдельной изолированной директории [`homework-20/`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/20-NoSQL/homework-20) без затирания решений предыдущих уроков.

---

## 📑 Содержание

1. [Постановка задачи и спецификация](#постановка-задачи-и-спецификация)
2. [Архитектура распределенной системы Lynks](#архитектура-распределенной-системы-lynks)
3. [Схема данных PostgreSQL и индексация](#схема-данных-postgresql-и-индексация)
4. [Кэширование в Redis (Cache-Aside Pattern)](#кэширование-в-redis-cache-aside-pattern)
5. [Наблюдаемость: Prometheus и Zerolog](#наблюдаемость-prometheus-и-zerolog)
6. [Контейнеризация и Docker Compose](#контейнеризация-и-docker-compose)
7. [Примеры HTTP-запросов (curl)](#примеры-http-запросов-curl)
8. [Результаты тестирования](#результаты-тестирования)

---

## Постановка задачи и спецификация

Сервис **Lynks** предоставляет пользователям возможности:
1. Принимать на вход длинную ссылку (`destination`), формировать на её основе короткую на базе собственного домена (`https://lynks.org/<id>`), сохранять в БД пару ссылок и возвращать пользователю короткую ссылку.
2. При переходе по короткой ссылке выполнять HTTP Redirect (`HTTP 302 Found`) на оригинальную страницу.
3. Ускорять чтение за счет микросервиса кэширования **Memcache** на базе in-memory СУБД **Redis** с временем жизни ключей (TTL) не более 24 часов.

---

## Архитектура распределенной системы Lynks

```
                                      +-------------------------+
                                      |     Redis Container     |
                                      |       (Port 6379)       |
                                      +-------------------------+
                                                   ^
                                                   | (TTL <= 24h)
                                                   v
+-------------------------+            +-------------------------+
|     PostgreSQL (DB)     |            |    Memcache Service     |
|   Таблица 'links'       |            |       (Port 8082)       |
+-------------------------+            +-------------------------+
             ^                                      ^
             | (Cache Miss)                         | (Cache Hit/Set)
             v                                      v
+----------------------------------------------------------------+
|                   Shortener Microservice                       |
|                        (Port 8080)                             |
|                                                                |
| • POST / или /api/v1/links (Создание короткой ссылки)         |
| • GET  /{short}            (Редирект: Cache-Aside алгоритм)    |
| • GET  /healthz            (Liveness Probe)                    |
| • GET  /metrics            (Prometheus CounterVec/Histogram)   |
| • Zerolog                  (Структурированное логирование)     |
+----------------------------------------------------------------+
```

### Алгоритм чтения (Cache-Aside Pattern):
1. Пользователь обращается по короткой ссылке: `GET /{short}`.
2. `Shortener` запрашивает микросервис `Memcache`: `GET /api/v1/cache/{short}`.
3. **Cache Hit**: если пара есть в Redis, `Shortener` мгновенно возвращает `HTTP 302 Found` (без нагрузки на PostgreSQL).
4. **Cache Miss**: если пары нет в кэше, `Shortener` делает `SELECT` из PostgreSQL по индексированному полю `short_url`.
5. При обнаружении ссылки `Shortener` асинхронно прогревает кэш (`POST /api/v1/cache`) и перенаправляет пользователя на оригинальный URL.

---

## Схема данных PostgreSQL и индексация

Файл [`schema.sql`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/20-NoSQL/homework-20/schema.sql):

```sql
CREATE TABLE IF NOT EXISTS links (
    id SERIAL PRIMARY KEY,
    short_url VARCHAR(64) UNIQUE NOT NULL,
    destination TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Индекс для эффективного поиска длинной гиперссылки по короткой (по ТЗ):
CREATE INDEX IF NOT EXISTS idx_links_short_url ON links (short_url);

-- Индекс для предотвращения повторной генерации дубликатов:
CREATE INDEX IF NOT EXISTS idx_links_destination ON links (destination);
```

---

## Кэширование в Redis (Cache-Aside Pattern)

- Драйвер: `github.com/go-redis/redis/v8`.
- Префикс ключей: `lynks:<short_url>`.
- Срок экспирации: строго **24 часа** (согласно ТЗ: не более 24 часов).
- Для автономного прохождения unit-тестов реализован потокобезопасный `MemoryStore` с поддержкой TTL.

---

## Наблюдаемость: Prometheus и Zerolog

1. **Prometheus Metrics**:
   - `http_requests_total` (`CounterVec` с метками `method`, `path`, `status`).
   - `http_request_duration_seconds` (`HistogramVec` с метками `method`, `path`).
   - `cache_hits_total` и `cache_misses_total` в сервисе shortener.
   - Эндпоинты экспорта метрик: `GET /metrics`.
2. **Структурированное логирование (Zerolog)**:
   - Пакет `github.com/rs/zerolog` пишет в `stdout` структурированные JSON/Console логи каждого входящего HTTP-запроса и возникающие ошибки.

---

## Контейнеризация и Docker Compose

Файл [`docker-compose.yml`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/20-NoSQL/homework-20/docker-compose.yml) оркестрирует 4 сервиса:
- `postgres` (образ `postgres:15-alpine`) с автонакатом `schema.sql`.
- `redis` (образ `redis:7-alpine`).
- `memcache` (порт 8082, образ на базе [`Dockerfile.memcache`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/20-NoSQL/homework-20/Dockerfile.memcache)).
- `shortener` (порт 8080, образ на базе [`Dockerfile.shortener`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/20-NoSQL/homework-20/Dockerfile.shortener)).

Команда запуска:
```bash
cd homework-20
docker-compose up --build -d
```

---

## Примеры HTTP-запросов (curl)

### 1. Сокращение гиперссылки:
```bash
curl --location 'http://localhost:8080' \
  --header 'Content-Type: application/json' \
  --data '{
    "destination": "https://github.com/Alex808r/go-production-engineering"
  }'
```

Ответ:
```json
{
  "shortUrl": "https://lynks.org/qz6d7",
  "destination": "https://github.com/Alex808r/go-production-engineering"
}
```

### 2. Переход по короткой ссылке (Redirect):
```bash
curl -I http://localhost:8080/qz6d7
```

Ответ:
```http
HTTP/1.1 302 Found
Location: https://github.com/Alex808r/go-production-engineering
```

### 3. Проверка наличия в кэше Memcache (порт 8082):
```bash
curl http://localhost:8082/api/v1/cache/qz6d7
```

Ответ:
```json
{
  "shortUrl": "qz6d7",
  "destination": "https://github.com/Alex808r/go-production-engineering"
}
```

---

## Результаты тестирования

```bash
$ cd go-production-engineering && go test -v -race ./lessons/20-NoSQL/homework-20/...
=== RUN   TestE2E_Lynks_FullFlow
--- PASS: TestE2E_Lynks_FullFlow (0.02s)
PASS
ok  	go-core-4/homework-20	1.698s
=== RUN   TestMemcacheAPI_SetAndGet
--- PASS: TestMemcacheAPI_SetAndGet (0.00s)
PASS
ok  	go-core-4/homework-20/memcache/pkg/api	1.850s
=== RUN   TestMemoryStore_BasicOperations
--- PASS: TestMemoryStore_BasicOperations (0.00s)
=== RUN   TestMemoryStore_TTLExpiration
--- PASS: TestMemoryStore_TTLExpiration (0.07s)
=== RUN   TestMemoryStore_Concurrency
--- PASS: TestMemoryStore_Concurrency (0.00s)
PASS
ok  	go-core-4/homework-20/memcache/pkg/store	1.710s
=== RUN   TestShortenerAPI_CreateAndRedirect_CacheAside
--- PASS: TestShortenerAPI_CreateAndRedirect_CacheAside (0.00s)
PASS
ok  	go-core-4/homework-20/shortener/pkg/api	1.800s
=== RUN   TestMemoryStorage_CRUD
--- PASS: TestMemoryStorage_CRUD (0.00s)
=== RUN   TestMemoryStorage_Concurrency
--- PASS: TestMemoryStorage_Concurrency (0.00s)
PASS
ok  	go-core-4/homework-20/shortener/pkg/storage	1.527s
```
