# Курс «Современная разработка на Go» (go-core-4)

[![CI](https://github.com/Alex808r/go_course_4/actions/workflows/ci.yml/badge.svg)](https://github.com/Alex808r/go_course_4/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/Alex808r/go_course_4)](https://goreportcard.com/report/github.com/Alex808r/go_course_4)
[![Go Version](https://img.shields.io/badge/Go-1.22%20%7C%201.23%20%7C%201.24%2B-00ADD8?logo=go)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> [!IMPORTANT]
> 🎓 **О происхождении репозитория и авторском вкладе:**  
> Репозиторий изначально брал за основу программу и базовые учебные заготовки курса **«Современная разработка на Go»** (Thinknetica, автор курса Дмитрий Титов).  
> 
> Однако в процессе обучения и глубокого погружения проект перерос в **масштабную самостоятельную работу и полноценную базу инженерных знаний**:
> - Все решения практических задач и домашних заданий (`homework-02` — `homework-20`), а также сквозные проекты **GoSearch** и **lynks** полностью спроектированы, написаны, отлажены и покрыты тестами **самостоятельно**.
> - Исходные лекционные примеры были существенно переработаны, дополнены подробными комментариями к механизмам рантайма и стандартам production-кода.
> - Создана фундаментальная авторская энциклопедия **[go-documentation](go-documentation)** (6 томов глубокой теории Go Core, руководства по Modern Go 1.21–1.25, лайвкодингу, базам данных, System Design и архитектуре), выходящая далеко за рамки исходной программы курса.

Комплексный репозиторий для изучения языка Go с нуля до уверенного Production-уровня и подготовки к техническим собеседованиям (Middle/Senior Go Developer).

Репозиторий включает:
- **22 тематических модуля** с примерами кода от базового синтаксиса до распределенных систем (`00-docs` — `22-modern-go`);
- **19 практических домашних заданий** с полным тестовым покрытием (`homework-02` — `homework-20`), реализованных самостоятельно с нуля;
- **2 сквозных продакшн-проекта**:
  - [GoSearch](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/GoSearch) — многопоточный поисковый движок (краулер, инвертированный индекс, бинарный поиск, TCP и REST API, веб-интерфейс);
  - [lynks](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lynks) — распределенный микросервис сокращения ссылок (Clean Architecture, PostgreSQL, Redis/Memcached, Kafka/Event-Driven, Prometheus метрики, Docker Compose);
- **Фундаментальную базу знаний** в каталоге [go-documentation](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation): 6 подробных томов теории, шпаргалки, настольный справочник по стандартным пакетам, гайды по базам данных, System Design и лайвкодингу.

---

## 📚 Структура репозитория

Каждый тематический каталог модуля (`02-syntax`, `03-algorithms` и т.д.) является **самодостаточным**: он содержит лекционные материалы, конспект темы и вложенное практическое домашнее задание (`homework-XX`).

```text
go-course-4/
├── 00-docs ... 21-interview    # Тематические модули (теория, примеры лекций + вложенные homework-XX)
│   ├── 02-syntax/
│   │   ├── 02-syntax.md        # Конспект темы
│   │   ├── ...                 # Примеры с лекции
│   │   └── homework-02/        # Практическое ДЗ к теме
│   └── ...
├── GoSearch/                   # Сквозной проект: поисковая система
├── lynks/                      # Сквозной проект: микросервисный URL Shortener
├── go-documentation/           # Энциклопедия Go: 6 томов теории и гайды
│   ├── 01-go-core-basics.md                # Том 1: Базовый синтаксис, типы, память
│   ├── 02-go-core-oop-methods.md           # Том 2: ООП, интерфейсы, io
│   ├── 03-go-core-concurrency.md           # Том 3: Конкурентность, каналы, sync
│   ├── 04-go-core-backend-web.md           # Том 4: Web, net/http, REST, БД
│   ├── 05-go-core-testing-arch.md          # Том 5: Тестирование, pprof, Clean Arch, gRPC, очереди
│   ├── 06-go-core-algorithms-interview.md  # Том 6: Алгоритмы, структуры данных, интервью
│   ├── course-codebase-guide.md            # Путеводитель по кодовой базе репозитория
│   ├── go-core-cheatsheet.md               # Быстрая шпаргалка по Go
│   ├── go-livecoding-guide.md              # Задачи для лайвкодинга с разбором
│   ├── go-database-interview-guide.md      # Вопросы и практика по SQL и NoSQL
│   ├── go-system-design-guide.md           # Проектирование распределенных систем
│   ├── go-modern-features-guide.md         # Руководство по фичам Go 1.21–1.24 (slog, iter, routing)
│   ├── go-secrets-and-traps-guide.md       # Секреты и ловушки Go (Uber Style & 100 Mistakes)
│   ├── homework-tasks.md                   # Формулировки домашних заданий
│   └── sql-livecoding-guide.md             # Практика SQL запросов
└── README.md                   # Главный путеводитель (этот файл)
```

---

## 🗺️ Карта модулей курса

| № | Тема | Каталог лекций | Домашнее задание | Ключевые концепции |
|---|---|---|---|---|
| **00** | Стандарты и Go Modules | [00-docs](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/00-docs) | — | Go modules, структура `cmd/` и `pkg/`, code review |
| **01** | Основы языка | [01-intro](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/01-intro) | — | Компиляция, пакеты, переменные, `fmt`, точка входа |
| **02** | Синтаксис и структуры памяти | [02-syntax](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/02-syntax) | [homework-02](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/02-syntax/homework-02) | Слайсы, мапы, указатели, структуры, срезы, cap/len |
| **03** | Алгоритмы | [03-algorithms](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/03-algorithms) | [homework-03](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/03-algorithms/homework-03) | Бинарный поиск, рекурсия, динамическое программирование, графы |
| **04** | Структуры данных | [04-datastructs](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/04-datastructs) | [homework-04](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/04-datastructs/homework-04) | Двусвязный список, BST (дерево поиска), массивы |
| **05** | Потоковый ввод-вывод | [05-io](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/05-io) | [homework-05](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/05-io/homework-05) | `io.Reader`, `io.Writer`, `bufio`, файлы, JSON сериализация |
| **06** | ООП и композиция | [06-oop](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/06-oop) | [homework-06](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/06-oop/homework-06) | Методы, получатели (значение/указатель), встраивание, DIP |
| **07** | Тестирование | [07-testing](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/07-testing) | [homework-07](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/07-testing/homework-07) | Unit-тесты, `TestMain`, Table-driven tests, TDD, моки |
| **08** | Профилирование и отладка | [08-prof_debug](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/08-prof_debug) | [homework-08](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/08-prof_debug/homework-08) | `pprof` (CPU/MEM), бенчмарки, `b.ResetTimer`, `go tool trace` |
| **09** | Интерфейсы и дженерики | [09-interfaces](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/09-interfaces) | [homework-09](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/09-interfaces/homework-09) | Полиморфизм, утиная типизация, Type Assertions, Generics |
| **10** | Конкурентность | [10-concurrency](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/10-concurrency) | [homework-10](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/10-concurrency/homework-10) | Горутины, каналы, `select`, `sync.WaitGroup`, `sync.Mutex`, `atomic`, `context`, Fan-Out / Fan-In |
| **11** | Сетевое взаимодействие | [11-network](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/11-network) | [homework-11](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/11-network/homework-11) | TCP сокеты, `net.Conn`, протоколы, таймауты и дедлайны |
| **12** | Веб-приложения | [12-web-apps](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/12-web-apps) | [homework-12](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/12-web-apps/homework-12) | `net/http`, маршрутизация, Gorilla Mux, HTML шаблоны (`html/template`) |
| **13** | REST API | [13-api](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/13-api) | [homework-13](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/13-api/homework-13) | RESTful API, JSON кодирование/декодирование, Middleware |
| **14** | RPC и gRPC | [14-RPC](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/14-RPC) | [homework-14](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/14-RPC/homework-14) | Go RPC (`net/rpc`), Protobuf v3, gRPC клиенты и серверы, streaming |
| **15** | Реляционные базы данных | [15-sql](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/15-sql) | [homework-15](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/15-sql/homework-15) | Реляционная модель, DDL, DML, индексы, внешние ключи, JOIN, транзакции |
| **16** | Go и СУБД | [16-db-apps](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/16-db-apps) | [homework-16](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/16-db-apps/homework-16) | `database/sql`, PostgreSQL драйвер `pgx/v4` и `pgx/v5`, пул коннектов, паттерн Репозиторий |
| **17** | Системный дизайн и архитектура | [17-system-design](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/17-system-design) | [homework-17](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/17-system-design/homework-17) | Clean Architecture, SOLID, Hexagonal / Onion, инверсия зависимостей |
| **18** | Микросервисы | [18-microservices](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/18-microservices) | [homework-18](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/18-microservices/homework-18) | 12 Factor App, Docker, Docker Compose, Healthcheck, конфигурация окружения |
| **19** | Очереди сообщений | [19-queue](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/19-queue) | [homework-19](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/19-queue/homework-19) | Apache Kafka, Event-Driven Architecture, Producer / Consumer группы, асинхронная аналитика |
| **20** | NoSQL и кэширование | [20-NoSQL](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/20-NoSQL) | [homework-20](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/20-NoSQL/homework-20) | Redis, In-Memory Cache (Cache-Aside), MongoDB, Prometheus метрики |
| **21** | Подготовка к собеседованиям | [21-interview](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/21-interview) | — | Разбор типовых ошибок, неблокирующий select, Fan-In, слайсы, UTF-8 руны |
| **22** | Современный Go и Best Practices | [22-modern-go](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/22-modern-go) | — | `log/slog`, `iter.Seq`, ServeMux, `slices`/`maps`, выравнивание памяти (padding), защита от утечек |

---

## 🛠️ Запуск тестов и проверка кода

В проекте настроена строгая проверка корректности кода.

### 1. Запуск всех тестов репозитория
```bash
go test ./...
```

### 2. Запуск статического анализатора (vet)
```bash
go vet ./...
```

### 3. Запуск тестов конкретного домашнего задания
```bash
go test -v ./17-system-design/homework-17/...
go test -v ./20-NoSQL/homework-20/...
```

### 4. Запуск с детекцией гонок (Race Detector)
```bash
go test -race ./10-concurrency/homework-10/...
go test -race ./21-interview/...
```

---

## 📖 Рекомендуемый план самостоятельного обучения

1. **Базовый уровень (Junior Go Developer):**
   - Прочитайте [Том 1: Основы и синтаксис](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/01-go-core-basics.md).
   - Освойте каталоги [01-intro](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/01-intro), [02-syntax](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/02-syntax), [03-algorithms](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/03-algorithms), [04-datastructs](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/04-datastructs), [05-io](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/05-io).
   - Выполните задания [homework-02](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/02-syntax/homework-02) — [homework-05](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/05-io/homework-05).

2. **Продвинутый уровень (Middle Go Developer):**
   - Изучите [Том 2: ООП и методы](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/02-go-core-oop-methods.md) и [Том 3: Конкурентность](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/03-go-core-concurrency.md).
   - Освойте профилирование [08-prof_debug](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/08-prof_debug), интерфейсы [09-interfaces](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/09-interfaces) и конкурентность [10-concurrency](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/10-concurrency).
   - Изучите сетевое программирование [11-network](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/11-network), REST API [13-api](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/13-api), базы данных [15-sql](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/15-sql), [16-db-apps](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/16-db-apps).

3. **Senior / Архитектура / Собеседования:**
   - Освойте [Том 5: Тестирование, архитектура и микросервисы](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/05-go-core-testing-arch.md).
   - Разберите Clean Architecture в [homework-17](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/17-system-design/homework-17).
   - Изучите распределенные очереди Kafka в [homework-19](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/19-queue/homework-19) и кэширование с метриками в [homework-20](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/20-NoSQL/homework-20).
   - Отработайте лайвкодинг по гайдам [go-livecoding-guide.md](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/go-livecoding-guide.md) и [21-interview](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/21-interview).
