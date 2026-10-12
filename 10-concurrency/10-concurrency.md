# 10-concurrency — Конкурентное программирование в Go: горутины, каналы и паттерны

В данной директории всесторонне рассматривается главная визитная карточка языка Go — модель конкурентности на основе легковесных потоков (горутин), каналов (CSP) и примитивов синхронизации пакета `sync`.

---

## 📋 Структура каталога

```text
10-concurrency/
├── 1-goroutine/           # Запуск горутин, легковесность (~2 КБ стек), планировщик Go (GMP)
├── 2-chan/                # Каналы (chan): однонаправленные, чтение, запись, закрытие (close)
├── 3-chan_select/         # Оператор select, неблокирующие операции (default), таймауты
├── 4-waitgroup/           # Синхронизация завершения горутин через sync.WaitGroup
├── 5-shared_mem/          # Разделяемая память: sync.Mutex и sync.RWMutex
├── 6-common_errors/       # Типовые ошибки: гонки данных (data race), утечки горутин, дедлоки
├── 7-atomic/              # Атомарные операции без блокировок (пакет sync/atomic)
├── 8-pattern-fan-out/     # Паттерн Fan-Out (пул воркеров для параллельной обработки)
├── 9-pattern-fan-in/      # Паттерн Fan-In (объединение данных из нескольких каналов в один)
├── 10-pattern-pipeline/   # Паттерн Pipeline (конвейерная потоковая обработка)
├── 11-context/            # Пакет context: отмена операций (WithCancel), дедлайны (WithTimeout)
└── homework-10/           # Домашнее задание 10: игра Ping-Pong между горутинами
```

---

## 🔍 Ключевые концепции

1. **Планировщик GMP:**
   - **G (Goroutine):** легковесный поток со стеком от 2 КБ.
   - **M (Machine):** поток операционной системы (OS thread).
   - **P (Processor):** логический процессор, контекст выполнения (`GOMAXPROCS`).

2. **Каналы (Channels):**
   - «Не общайтесь, делясь памятью; делитесь памятью, общаясь» (Rob Pike).
   - Чтение из закрытого канала возвращает zero-value и `false` (вторым значением).
   - Запись в закрытый канал или повторное закрытие приводит к **панике**.
   - Чтение или запись в `nil`-канал блокирует горутину **навсегда**.

3. **Контекст (`context.Context`):**
   - Распространение сигналов отмены и таймаутов вниз по дереву вызовов.
   - Освобождение сетевых и вычислительных ресурсов при отключении клиента или истечении дедлайна.

4. **Детектор гонок (Race Detector):**
   - Проверка наличия состояния гонки (data race): `go test -race ./...`.

---

## 🚀 Запуск и проверка

```bash
# Запуск демонстраций уроков
go run ./10-concurrency/1-goroutine/goroutine.go
go run ./10-concurrency/2-chan/chan.go
go run ./10-concurrency/3-chan_select/select.go
go run ./10-concurrency/4-waitgroup/sync.go
go run -race ./10-concurrency/5-shared_mem/shared.go
go run ./10-concurrency/6-common_errors/errors.go
go run -race ./10-concurrency/7-atomic/atomic.go
go run -race ./10-concurrency/8-pattern-fan-out/fan-out.go
go run -race ./10-concurrency/9-pattern-fan-in/fan-in.go
go run -race ./10-concurrency/10-pattern-pipeline/pipeline.go
go run -race ./10-concurrency/11-context/context.go

# Проверка конкурентного кода на гонки данных и запуск тестов
go test -v -race -cover ./10-concurrency/...
```

---

## 🔗 Связанные материалы
- **Теория:** [Том 3: Конкурентное программирование (Разделы 11, 14–16: Каналы, sync, context)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/03-go-core-concurrency.md)
- **Практическое ДЗ:** [homework-10](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/10-concurrency/homework-10) (Синхронизация двух игроков Ping-Pong через каналы)

