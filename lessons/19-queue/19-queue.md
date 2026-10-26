# 19-queue — Очереди сообщений: Apache Kafka и Event-Driven архитектура

В данном модуле изучается событийно-ориентированная архитектура (Event-Driven Architecture) и работа с распределенным брокером сообщений Apache Kafka на языке Go: публикация событий (Producer / Writer), вычитка и обработка сообщений (Consumer / Reader), группы потребителей (Consumer Groups), семантики доставки сообщений, а также сквозная асинхронная аналитика микросервисов.

---

## 📋 Структура каталога

```text
19-queue/
├── 1-kafka/                     # Базовый Go-клиент Apache Kafka
│   ├── kafka.go                 # Client, New, Producer (Writer), Consumer (Reader), Commit
│   └── kafka_test.go            # Модульные тесты валидации подключения и конфигурации
├── docker-compose.yml           # Развертывание кластера Kafka и Zookeeper
└── homework-19/                 # Домашнее задание №19: Асинхронный сервис аналитики ссылок
    ├── Dockerfile               # Multi-stage Dockerfile для сервисов shortener и analytics
    ├── docker-compose.yml       # Стек Zookeeper + Kafka + Shortener + Analytics
    ├── README.md                # Полная документация, Event Sourcing и примеры curl
    ├── e2e_test.go              # Сквозной E2E-тест асинхронной доставки событий
    ├── cmd/
    │   ├── shortener/main.go    # Точка входа сервиса коротких ссылок (Producer)
    │   └── analytics/main.go    # Точка входа сервиса аналитики (Consumer)
    └── pkg/
        ├── event/               # Абстракция брокера: Event, Producer, Consumer, ChannelBroker
        ├── shortener/           # HTTP API, Storage и Config микросервиса shortener
        └── analytics/           # HTTP API, Consumer Worker, Stats Service и Config аналитики
```

---

## 🔍 Ключевые концепции

1. **Зачем нужны брокеры сообщений (Event-Driven Architecture):**
   - **Асинхронность:** тяжелые операции (отправка email, агрегация аналитики, индексация) выносятся из синхронного жизненного цикла HTTP-запроса.
   - **Сглаживание пиковых нагрузок (Backpressure / Load Leveling):** потребители считывают события из очереди с комфортной для себя скоростью, не перегружая базу данных.
   - **Слабая связность (Decoupling):** сервис `shortener` ничего не знает о существовании сервиса `analytics`; он лишь публикует факт `LinkCreatedEvent` в топик.

2. **Модель данных и компонентов Kafka:**
   - **Topic (Топик):** логическая категория или поток сообщений.
   - **Partition (Партиция):** строго упорядоченный неизменяемый лог записей (append-only commit log). Партиции являются базовой единицей параллелизма.
   - **Consumer Group (Группа потребителей):** механизм автоматического распределения партиций топика между несколькими репликами одного сервиса. Каждая партиция читается ровно одним консьюмером из группы.
   - **Offset (Смещение):** монотонно возрастающий числовой идентификатор позиции сообщения внутри партиции.

3. **Семантика доставки сообщений:**
   - **At-Most-Once (не более одного раза):** оффсет коммитится *до* обработки. При сбое сообщение теряется.
   - **At-Least-Once (как минимум один раз — РЕКОМЕНДУЕМЫЙ):** чтение (`FetchMessage`) ➔ обработка бизнес-логикой ➔ фиксация смещения (`CommitMessages`). Сообщения не теряются; бизнес-логика должна быть идемпотентной.
   - **Exactly-Once (ровно один раз):** требует двухфазных коммитов (транзакционный Producer/Consumer).

4. **Тестирование очередей без внешних зависимостей:**
   - Использование паттерна Ports & Adapters: интерфейсы `event.Producer` и `event.Consumer` позволяют в юнит- и E2E-тестах использовать легковесный in-memory брокер на каналах Go (`ChannelBroker`), прогоняя тесты за миллисекунды без запуска Kafka.

---

## 🚀 Запуск тестов и локального кластера

```bash
# Запуск всех тестов модуля 19 с детектором гонок
go test -v -race ./lessons/19-queue/...

# Запуск тестов с замером бенчмарков
go test -v -race -bench=. ./lessons/19-queue/...

# Запуск кластера Kafka в Docker Compose
docker compose -f ./lessons/19-queue/docker-compose.yml up -d
```

---

## 🔗 Файлы и связанные материалы
- [1-kafka/kafka.go](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/19-queue/1-kafka/kafka.go) — Пример работы с `kafka-go` (Writer, Reader, CommitMessages)
- [docker-compose.yml](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/19-queue/docker-compose.yml) — Docker Compose для локального Kafka
- [homework-19/README.md](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/19-queue/homework-19/README.md) — Домашнее задание №19: Асинхронная аналитика
- [homework-19/e2e_test.go](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/19-queue/homework-19/e2e_test.go) — Сквозной E2E-тест пайплайна Shortener -> Kafka -> Analytics
- [homework-19/pkg/event/event.go](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/19-queue/homework-19/pkg/event/event.go) — Модели событий и контракты брокера
- **Теория:** [Том 5: Тестирование и архитектура (Раздел 24: Очереди сообщений, брокеры, Kafka)](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/go-documentation/05-go-core-testing-arch.md)
