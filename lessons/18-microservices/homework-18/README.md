# Домашнее задание №18: Разработка микросервисов (URL Shortener Service)

Решение домашнего задания по курсу «Современная разработка на Go» (Занятие №18: Разработка микросервисов).  
Реализовано в отдельной директории `homework-18` со всеми требованиями к микросервисной архитектуре и контейнеризации.

---

## 📑 Содержание

1. [Постановка задачи](#постановка-задачи)
2. [Архитектура микросервиса](#архитектура-микросервиса)
3. [Задача №1: HTTP REST API и хранилище в памяти](#задача-1-http-rest-api-и-хранилище-в-памяти)
4. [Задача №2: Контейнеризация (Dockerfile и docker-compose)](#задача-2-контейнеризация-dockerfile-и-docker-compose)
5. [Мониторинг, логирование и Cloud-Ready эндпоинты](#мониторинг-логирование-и-cloud-ready-эндпоинты)
6. [Инструкция по запуску и примеры запросов](#инструкция-по-запуску-и-примеры-запросов)
7. [Результаты тестирования](#результаты-тестирования)

---

## Постановка задачи

- **Задача №1**: Разработать микросервис хранения коротких гиперссылок с HTTP-API:
  1. *Добавление ссылки*: принимает оригинальную ссылку, добавляет её в хранилище и возвращает короткую ссылку.
  2. *Получение оригинальной ссылки*: принимает короткую ссылку и возвращает оригинальную.
  - **Критическое требование**: Хранилище ссылок должно быть реализовано в виде **слайса в памяти** (`[]Link`). Постоянное хранилище не требуется.
- **Задача №2**: Запустить разработанный микросервис в контейнере:
  1. Написать `Dockerfile`.
  2. Запустить у себя на компьютере микросервис и подключиться к приложению браузером или HTTP-клиентом (curl, Postman).

---

## Архитектура микросервиса

```
homework-18/
├── Dockerfile                   # Multi-stage production сборка (alpine + non-root user)
├── docker-compose.yml           # Декларативный запуск сервиса с Healthcheck
├── README.md                    # Документация и примеры работы
├── cmd/
│   └── shortener/
│       └── main.go              # Точка входа: инициализация, Graceful Shutdown
└── pkg/
    ├── api/
    │   ├── api.go               # HTTP API, роутер, middleware, healthz, metrics
    │   └── api_test.go          # Модульные тесты API (100% pass)
    ├── config/
    │   └── config.go            # Конфигурация через ENV (PORT, BASE_URL)
    └── storage/
        ├── storage.go           # Контракт хранилища и модель Link
        ├── slicestore.go        # Реализация хранилища на базе среза в памяти
        └── slicestore_test.go   # Модульные тесты хранилища (с проверкой конкурентности)
```

---

## Задача №1: HTTP REST API и хранилище в памяти

### 1. Хранилище на слайсе ([`pkg/storage/slicestore.go`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/18-microservices/homework-18/pkg/storage/slicestore.go)):
- В соответствии с ТЗ хранилище использует срез в памяти `[]Link`.
- Потокобезопасность обеспечена `sync.RWMutex` (безопасное параллельное чтение и запись под `-race`).
- Генерация уникальных компактных идентификаторов (URL-safe hex).
- Поддерживается идемпотентность (повторное сохранение того же URL не плодит дубликаты).

### 2. HTTP API ([`pkg/api/api.go`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/18-microservices/homework-18/pkg/api/api.go)):
- `POST /api/v1/links` — создание короткой ссылки (`{"url": "https://go.dev"}`).
- `GET /api/v1/links/{id}` — получение объекта ссылки и оригинала по идентификатору.
- `POST /api/v1/links/resolve` — получение оригинального URL по полной короткой ссылке.
- `GET /r/{id}` — HTTP 302 Redirect на оригинальную страницу.
- `GET /api/v1/links` — список всех сохранённых ссылок.

---

## Задача №2: Контейнеризация (Dockerfile и docker-compose)

Файл [`Dockerfile`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/18-microservices/homework-18/Dockerfile) построен по лучшим индустриальным практикам:
1. **Multi-stage build**:
   - `builder` на базе `golang:1.22-alpine`: компиляция статического бинарника с флагами `CGO_ENABLED=0` и `-ldflags="-s -w"`.
   - `runtime` на базе ультра-легковесного `alpine:3.19` (итоговый образ < 25 МБ).
2. **Безопасность**:
   - Запуск от имени непривилегированного пользователя `appuser` (UID 1000).
   - Включены корневые сертификаты `ca-certificates` и зона времени `tzdata`.
3. **Оркестрация**:
   - [`docker-compose.yml`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/18-microservices/homework-18/docker-compose.yml) с проверкой здоровья контейнера (`healthcheck` через `/healthz`).

---

## Мониторинг, логирование и Cloud-Ready эндпоинты

- **Структурированное логирование**: HTTP Middleware логирует метод, путь, IP клиента и время ответа.
- **Healthcheck & Kubernetes Probes**:
  - `GET /healthz` — проверка живучести (Liveness Probe).
  - `GET /ready` — проверка готовности принимать трафик (Readiness Probe).
- **Метрики Prometheus**:
  - `GET /metrics` — экспорт счетчика запросов `http_requests_total` и количества ссылок `shortener_links_stored_total`.
- **Graceful Shutdown**: корректная обработка `SIGINT`/`SIGTERM` с таймаутом завершения соединений 5 секунд.

---

## Инструкция по запуску и примеры запросов

### 1. Локальный запуск:

```bash
cd go_course_4
go run ./lessons/18-microservices/homework-18/cmd/shortener/main.go
```

### 2. Сборка и запуск в Docker:

```bash
# Сборка образа
docker build -t url-shortener:latest -f homework-18/Dockerfile .

# Запуск контейнера
docker run --rm -p 8080:8080 url-shortener:latest

# Либо через Docker Compose
docker compose -f homework-18/docker-compose.yml up -d
```

### 3. Примеры HTTP-запросов:

```bash
# 1. Сократить ссылку
curl -X POST http://localhost:8080/api/v1/links \
  -H "Content-Type: application/json" \
  -d '{"url": "https://thinknetica.com/courses/golang"}' | jq

# Пример ответа:
# {
#   "id": "e3a891f0",
#   "original": "https://thinknetica.com/courses/golang",
#   "short_url": "http://localhost:8080/r/e3a891f0",
#   "created_at": "2026-09-25T19:35:00Z"
# }

# 2. Получить ссылку по ID
curl http://localhost:8080/api/v1/links/e3a891f0 | jq

# 3. Разрешить ссылку через POST /resolve
curl -X POST http://localhost:8080/api/v1/links/resolve \
  -H "Content-Type: application/json" \
  -d '{"short_url": "http://localhost:8080/r/e3a891f0"}' | jq

# 4. Проверить редирект
curl -I http://localhost:8080/r/e3a891f0

# 5. Проверить метрики Prometheus
curl http://localhost:8080/metrics
```

---

## Результаты тестирования

```bash
$ cd go_course_4 && go test -v -race -cover ./lessons/18-microservices/homework-18/...
=== RUN   TestAPI_CreateAndGetLink
--- PASS: TestAPI_CreateAndGetLink (0.00s)
=== RUN   TestAPI_ErrorCases
--- PASS: TestAPI_ErrorCases (0.00s)
=== RUN   TestAPI_HealthzAndMetrics
--- PASS: TestAPI_HealthzAndMetrics (0.00s)
PASS
coverage: 89.1% of statements
ok  	go-core-4/homework-18/pkg/api	1.555s
=== RUN   TestSliceStore_SaveAndGet
--- PASS: TestSliceStore_SaveAndGet (0.00s)
=== RUN   TestSliceStore_ValidationAndNotFound
--- PASS: TestSliceStore_ValidationAndNotFound (0.00s)
=== RUN   TestSliceStore_Concurrency
--- PASS: TestSliceStore_Concurrency (0.00s)
PASS
coverage: 94.9% of statements
ok  	go-core-4/homework-18/pkg/storage	1.797s
```
