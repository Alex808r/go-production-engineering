# Домашнее задание №19: Очереди сообщений и Асинхронная аналитика (Kafka / Event-Driven)

Решение домашнего задания по курсу «Современная разработка на Go» (**Занятие №19: Брокеры сообщений**).  
Реализовано в отдельной изолированной директории [`homework-19/`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/19-queue/homework-19) без модификации и затирания решений предыдущих уроков.

---

## 📑 Содержание

1. [Постановка задачи](#постановка-задачи)
2. [Архитектура системы и компонентов](#архитектура-системы-и-компонентов)
3. [Событийная модель (Event Sourcing)](#событийная-модель-event-sourcing)
4. [Реализация брокера сообщений (Apache Kafka & In-Memory Channel)](#реализация-брокера-сообщений-apache-kafka--in-memory-channel)
5. [Микросервис URL Shortener (Producer)](#микросервис-url-shortener-producer)
6. [Микросервис Analytics (Consumer)](#микросервис-analytics-consumer)
7. [Контейнеризация и Docker Compose](#контейнеризация-и-docker-compose)
8. [Инструкция по запуску и примеры запросов](#инструкция-по-запуску-и-примеры-запросов)
9. [Результаты тестирования](#результаты-тестирования)

---

## Постановка задачи

**Задача №1 (Слайд 19 лекции №19):**
1. Добавить микросервис аналитики в систему сокращения гиперссылок.
2. Новый микросервис должен содержать статистику по количеству ссылок: всего штук, средняя длина, статистика по доменам и т.д.
3. При сохранении новой гиперссылки существующий микросервис должен отправлять информацию о ней в очередь (**Apache Kafka** или RabbitMQ).
4. Микросервис аналитики должен подписаться на эту очередь и обновлять статистику при получении каждого нового сообщения.

---

## Архитектура системы и компонентов

```
                                  +-------------------------------------+
                                  |     Apache Kafka Cluster (Topic)    |
                                  |            "links-events"           |
                                  +-------------------------------------+
                                            ^                 |
                                (Publish)   |                 | (Subscribe)
                             LinkCreatedEvent                 | Consumer Group
                                            |                 v
+-----------------------------+             |   +-----------------------------+
|    Shortener Microservice   |-------------+   |    Analytics Microservice   |
|         (Port :8080)        |                 |         (Port :8081)        |
+-----------------------------+                 +-----------------------------+
| • POST /api/v1/links        |                 | • GET  /api/v1/stats        |
| • GET  /api/v1/links/{id}   |                 | • GET  /healthz             |
| • POST /api/v1/links/resolve|                 | • GET  /ready               |
| • GET  /r/{id} (Redirect)   |                 | • GET  /metrics (Prometheus)|
| • In-Memory Slice Storage   |                 | • Thread-safe Metrics Store |
+-----------------------------+                 +-----------------------------+
```

### Структура проекта `homework-19`:

```
homework-19/
├── Dockerfile                   # Multi-stage сборка для shortener и analytics
├── docker-compose.yml           # Zookeeper, Kafka, Shortener, Analytics
├── README.md                    # Документация и техническое описание
├── e2e_test.go                  # Сквозной E2E-тест асинхронного пайплайна
├── cmd/
│   ├── shortener/
│   │   └── main.go              # Точка входа URL Shortener (Producer)
│   └── analytics/
│       └── main.go              # Точка входа сервиса аналитики (Consumer)
└── pkg/
    ├── event/
    │   ├── event.go             # Интерфейсы Producer, Consumer, Handler и модель события
    │   ├── kafka.go             # Реализация Kafka Producer/Consumer на github.com/segmentio/kafka-go
    │   ├── channel.go           # In-Memory брокер на каналах Go для тестов без брокера
    │   └── broker_test.go       # Тесты брокера и валидации
    ├── shortener/
    │   ├── api/                 # HTTP API с публикацией событий в очередь
    │   ├── config/              # Конфигурация (PORT, BASE_URL, KAFKA_BROKERS)
    │   └── storage/             # Потокобезопасное хранилище на срезе (SliceStore)
    └── analytics/
        ├── api/                 # REST API аналитики (GET /api/v1/stats, метрики)
        ├── config/              # Конфигурация аналитики (PORT, KAFKA_GROUP_ID)
        ├── consumer/            # Фоновый воркер подписки и обработки сообщений
        └── service/             # Бизнес-логика расчета средней длины, доменов и счетчиков
```

---

## Событийная модель (Event Sourcing)

Для обеспечения слабой связности (Loose Coupling) сервисы общаются асинхронными событиями.

Структура события [`pkg/event/event.go`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/19-queue/homework-19/pkg/event/event.go):

```go
type LinkCreatedEvent struct {
    EventID   string    `json:"event_id"`   // Уникальный ID события (идемпотентность)
    Type      string    `json:"type"`       // "link.created"
    LinkID    string    `json:"link_id"`    // Идентификатор короткой ссылки
    Original  string    `json:"original"`   // Исходный длинный URL
    ShortURL  string    `json:"short_url"`  // Сформированный короткий URL
    CreatedAt time.Time `json:"created_at"` // Временная метка генерации
}
```

---

## Реализация брокера сообщений (Apache Kafka & In-Memory Channel)

1. **Гарантии доставки (At-Least-Once Delivery)**:
   - В соответствии со слайдом 14 лекции (*«Выбор варианта подтверждения сообщения критически важен!»*), консьюмер применяет стратегию:  
     **Чтение (`FetchMessage`) ➔ Обработка бизнес-логикой (`RecordLink`) ➔ Подтверждение смещения (`CommitMessages`)**.
   - Это гарантирует, что если сервис аналитики аварийно завершится во время обработки, непрочитанные сообщения будут обработаны повторно после перезапуска.

2. **Защита от Poison Pill**:
   - При получении битого JSON-сообщения консьюмер логирует ошибку и коммитит оффсет, чтобы очередь не зависала в бесконечном цикле.

3. **In-Memory Channel Broker**:
   - Для автономного запуска модульных тестов без необходимости поднимать локальный кластер Kafka реализован [`ChannelBroker`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/19-queue/homework-19/pkg/event/channel.go). Все тесты выполняются за сотые доли секунды.

---

## Микросервис URL Shortener (Producer)

- Эндпоинт создания ссылки `POST /api/v1/links`:
  1. Валидирует URL и сохраняет его в потокобезопасный срез памяти `SliceStore` (`sync.RWMutex`).
  2. Генерирует уникальный ID и формирует `ShortURL`.
  3. Асинхронно отправляет событие `LinkCreatedEvent` в топик `links-events` Kafka.
  4. Возвращает клиенту HTTP 201 Created.

---

## Микросервис Analytics (Consumer)

Сервис агрегирует расширенную статистику:
- `total_links`: Общее количество сокращенных ссылок.
- `total_original_length` & `avg_original_length`: Суммарная и средняя длина исходных URL.
- `total_short_length` & `avg_short_length`: Суммарная и средняя длина коротких URL.
- `min_original_length` & `max_original_length`: Длина самого короткого и длинного URL.
- `domain_counts`: Частотный словарь по доменам (например, `github.com`: 4, `go.dev`: 2).
- `last_event_at`, `updated_at`: Временные метки последних событий.

Пример ответа `GET /api/v1/stats`:

```json
{
  "total_links": 3,
  "total_original_length": 115,
  "avg_original_length": 38.33,
  "total_short_length": 87,
  "avg_short_length": 29.0,
  "min_original_length": 14,
  "max_original_length": 58,
  "domain_counts": {
    "github.com": 1,
    "go.dev": 1,
    "microservices.io": 1
  },
  "last_link_id": "9411840c",
  "last_original": "https://microservices.io/patterns/data/event-sourcing.html",
  "last_event_at": "2026-09-26T10:59:14.412Z",
  "updated_at": "2026-09-26T10:59:14.415Z"
}
```

---

## Контейнеризация и Docker Compose

Файл [`docker-compose.yml`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/19-queue/homework-19/docker-compose.yml) поднимает полный стек:
- `zookeeper` (порт 22181)
- `kafka` (порты 9092, 29092)
- `shortener` (порт 8080)
- `analytics` (порт 8081)

Запуск:
```bash
cd homework-19
docker-compose up --build -d
```

---

## Инструкция по запуску и примеры запросов

### 1. Добавление новой ссылки в Shortener (порт 8080):
```bash
curl -X POST http://localhost:8080/api/v1/links \
  -H "Content-Type: application/json" \
  -d '{"url": "https://pkg.go.dev/github.com/segmentio/kafka-go"}'
```

Ответ:
```json
{
  "id": "e2f0a14b",
  "original": "https://pkg.go.dev/github.com/segmentio/kafka-go",
  "short_url": "http://localhost:8080/r/e2f0a14b",
  "created_at": "2026-09-26T11:00:00Z"
}
```

### 2. Получение агрегированной аналитики (порт 8081):
```bash
curl http://localhost:8081/api/v1/stats
```

### 3. Проверка здоровья сервисов:
```bash
curl http://localhost:8080/healthz
curl http://localhost:8081/healthz
```

---

## Результаты тестирования

Все пакеты протестированы с детектором гонок (`-race`) и проверкой покрытия:

```bash
$ cd go-production-engineering && go test -v -race ./lessons/19-queue/homework-19/...
=== RUN   TestE2E_ShortenerToAnalytics_AsyncFlow
--- PASS: TestE2E_ShortenerToAnalytics_AsyncFlow (0.02s)
PASS
ok  	go-core-4/homework-19	1.288s
=== RUN   TestAnalyticsAPI_GetStats
--- PASS: TestAnalyticsAPI_GetStats (0.00s)
=== RUN   TestAnalyticsAPI_HealthAndMetrics
--- PASS: TestAnalyticsAPI_HealthAndMetrics (0.00s)
PASS
ok  	go-core-4/homework-19/pkg/analytics/api	0.512s
=== RUN   TestWorker_Processing
--- PASS: TestWorker_Processing (0.02s)
PASS
ok  	go-core-4/homework-19/pkg/analytics/consumer	0.485s
=== RUN   TestService_RecordLink_CalculatesStats
--- PASS: TestService_RecordLink_CalculatesStats (0.00s)
=== RUN   TestService_Concurrency
--- PASS: TestService_Concurrency (0.00s)
PASS
ok  	go-core-4/homework-19/pkg/analytics/service	0.450s
=== RUN   TestChannelBroker_PublishAndSubscribe
--- PASS: TestChannelBroker_PublishAndSubscribe (0.00s)
=== RUN   TestChannelBroker_Closed
--- PASS: TestChannelBroker_Closed (0.00s)
=== RUN   TestKafkaConfig_Validation
--- PASS: TestKafkaConfig_Validation (0.00s)
PASS
ok  	go-core-4/homework-19/pkg/event	0.520s
=== RUN   TestAPI_CreateLink_PublishesEvent
--- PASS: TestAPI_CreateLink_PublishesEvent (0.00s)
=== RUN   TestAPI_GetAndResolve
--- PASS: TestAPI_GetAndResolve (0.00s)
=== RUN   TestAPI_ErrorsAndHealth
--- PASS: TestAPI_ErrorsAndHealth (0.00s)
PASS
ok  	go-core-4/homework-19/pkg/shortener/api	0.490s
=== RUN   TestSliceStore_SaveAndGet
--- PASS: TestSliceStore_SaveAndGet (0.00s)
=== RUN   TestSliceStore_Validation
--- PASS: TestSliceStore_Validation (0.00s)
=== RUN   TestSliceStore_Concurrency
--- PASS: TestSliceStore_Concurrency (0.00s)
PASS
ok  	go-core-4/homework-19/pkg/shortener/storage	0.510s
```
