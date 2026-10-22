# Курс «Современная разработка на Go» (go-core-4)

[![CI](https://github.com/Alex808r/go_course_4/actions/workflows/ci.yml/badge.svg)](https://github.com/Alex808r/go_course_4/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/Alex808r/go_course_4)](https://goreportcard.com/report/github.com/Alex808r/go_course_4)
[![Go Version](https://img.shields.io/badge/Go-1.22%20%7C%201.23%20%7C%201.24%2B-00ADD8?logo=go)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> [!IMPORTANT]
> 🎓 **О происхождении репозитория и авторском вкладе:**  
> Репозиторий изначально брал за основу программу и базовые учебные заготовки курса **«Современная разработка на Go»** (онлайн-школа [Thinknetica](https://thinknetica.com/), оригинальный репозиторий курса: [thinknetica/go_course_4](https://github.com/thinknetica/go_course_4), автор курса — Дмитрий Титов).  
> 
> Однако в процессе обучения и глубокого погружения проект перерос в **масштабную самостоятельную работу и полноценную базу инженерных знаний**:
> - Все решения практических задач и домашних заданий (`homework-02` — `homework-20`), а также сквозные проекты **GoSearch**, **lynks** и высоконагруженный финансовый сервис **wallet** полностью спроектированы, написаны, отлажены и покрыты тестами **самостоятельно**.
> - Исходные лекционные примеры были существенно переработаны, дополнены подробными комментариями к механизмам рантайма и стандартам production-кода.
> - Создана фундаментальная авторская энциклопедия **[go-documentation](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation)** (6 томов глубокой теории Go Core, руководства по Modern Go 1.21–1.24+, Docker & Compose, лайвкодингу, базам данных, System Design, архитектурным паттернам Enterprise и разборам рантайм-ловушек), выходящая далеко за рамки исходной программы курса.
> - **Использование материалов книги Джона Боднера:** При написании и расширении энциклопедии [go-documentation](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation) использованы теоретические материалы и практические примеры из книги **Джона Боднера «Learning Go: An Idiomatic Approach to Real-World Go Programming» (2nd Edition, O'Reilly, 2024)**, а канонические реализации и архитектурные паттерны взяты напрямую из официального репозитория книги: **[learning-go-book-2e (GitHub)](https://github.com/learning-go-book-2e)** с сохранением атрибуции первоисточника.
> - Теоретическая база также систематизирует и цитирует эталонные практики других мировых источников: классики **Алана Донована и Брайана Кернигана (*The Go Programming Language*)**, стайлгайдов **Uber Go** и **Google**, а также разборов ошибок **Тейвы Харшаньи (*100 Go Mistakes*)**.

Комплексный репозиторий для изучения языка Go с нуля до уверенного Production-уровня и подготовки к техническим собеседованиям (Middle/Senior Go Developer).

Репозиторий включает:
- **22 тематических модуля** с примерами кода от базового синтаксиса до распределенных систем в каталоге [lessons](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons) (`00-docs` — `22-modern-go`);
- **19 практических домашних заданий** с полным тестовым покрытием (`homework-02` — `homework-20`), реализованных самостоятельно с нуля;
- **3 сквозных продакшн-проекта**:
  - [GoSearch](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/GoSearch) — многопоточный поисковый движок (краулер, инвертированный индекс, бинарный поиск, TCP и REST API, веб-интерфейс);
  - [lynks](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lynks) — распределенный микросервис сокращения ссылок (Clean Architecture, PostgreSQL, Redis/Memcached, Kafka/Event-Driven, Prometheus метрики, Docker Compose);
  - [wallet](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/wallet) — высоконагруженный финансовый REST-сервис управления балансом (1000+ RPS на кошелек, шардированный `KeyLocker`, пул соединений `pgxpool`, строгая изоляция транзакций, идемпотентность, Docker Compose, e2e стресс-тесты под нагрузкой);
- **Фундаментальную базу знаний** в каталоге [go-documentation](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation): 6 подробных томов теории (включая материалы и практические реализации из книги Джона Боднера *«Learning Go 2nd Edition»* и её репозитория [learning-go-book-2e](https://github.com/learning-go-book-2e)), шпаргалки, настольный справочник по стандартным пакетам и коллекциям, гайды по Docker & Compose, базам данных, System Design, Enterprise-стандартам (Kafka Outbox/DLQ, gRPC, CI/CD, AppSec), лайвкодингу (29 задач по SQL и алгоритмы Go) и 14 ловушкам рантайма.

---

## 📚 Структура репозитория

Каждый тематический каталог модуля внутри директории [`lessons/`](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons) (`02-syntax`, `03-algorithms` и т.д.) является **самодостаточным**: он содержит лекционные материалы, конспект темы и вложенное практическое домашнее задание (`homework-XX`).

```text
go-course-4/
├── lessons/                    # Тематические модули курса (00-docs ... 22-modern-go)
│   ├── 00-docs/                # Стандарты и Go Modules
│   ├── 01-intro/               # Первая программа и CLI-приложения
│   ├── 02-syntax/              # Синтаксис, память, срезы, структуры
│   │   ├── 02-syntax.md        # Конспект темы
│   │   ├── ...                 # Примеры с лекции
│   │   └── homework-02/        # Практическое ДЗ к теме
│   └── ...
├── GoSearch/                   # Сквозной проект: поисковая система
├── lynks/                      # Сквозной проект: микросервисный URL Shortener
├── wallet/                     # Сквозной проект: High-Load Wallet REST API (1000+ RPS)
├── go-documentation/           # Энциклопедия Go: 7 томов теории и прикладные гайды
│   ├── 01-go-core-basics.md                # Том 1.1: Базовый синтаксис, control flow, функции, числа, константы
│   ├── 01-go-core-data-structures.md       # Том 1.2: Срезы, мапы, строки, указатели, структуры, память
│   ├── 02-go-core-oop-methods.md           # Том 2: ООП, интерфейсы, io, Google Wire, wildcard errors, errors.Join
│   ├── 03-go-core-concurrency.md           # Том 3: Конкурентность, каналы, sync, GMP, PressureGauge, cgo.Handle
│   ├── 04-go-core-backend-web.md           # Том 4: Web, net/http (Go 1.22+), ResponseController streaming, REST, БД
│   ├── 05-go-core-testing-arch.md          # Том 5: Тестирование (testdata, cache), internal-правила, Clean Arch, gRPC, Kafka
│   ├── 06-go-core-algorithms-interview.md  # Том 6: Алгоритмы, структуры данных, unsafe/reflect, Swiss Tables, STAR
│   ├── course-codebase-guide.md            # Путеводитель по кодовой базе репозитория
│   ├── docker-and-compose-starter-guide.md # Практическое руководство по Docker и Docker Compose
│   ├── go-core-cheatsheet.md               # Быстрая шпаргалка по Go
│   ├── go-database-interview-guide.md      # Вопросы и практика по SQL и NoSQL
│   ├── go-livecoding-guide.md              # Задачи для лайвкодинга с разбором (47 задач)
│   ├── go-modern-features-guide.md         # Руководство по фичам Go 1.21–1.24 (slog, iter, routing)
│   ├── go-production-backend-guide.md      # Промышленный backend: CI/CD, Kafka Outbox/DLQ, gRPC, AppSec
│   ├── go-secrets-and-traps-guide.md       # 14 ловушек и секретов Go (память, утечки горутин, typed nil)
│   ├── go-std-collections-guide.md         # Справочник по стандартным коллекциям и пакетам Go
│   ├── go-system-design-guide.md           # Проектирование распределенных систем
│   ├── homework-tasks.md                   # Формулировки домашних заданий
│   ├── sql-livecoding-guide.md             # Практика SQL запросов (29 задач от базовых до LeetCode Hard)
│   └── mermaid-diagrams/                  # Стандарты визуализации архитектуры Mermaid
└── README.md                   # Главный путеводитель (этот файл)
```

---

## 🚀 Сквозные продакшн-проекты

В репозитории реализованы 3 полноценных продакшн-сервиса, демонстрирующих реальные сценарии проектирования:

| Проект | Архитектура и стек | Ключевые инженерные решения |
|---|---|---|
| **[GoSearch](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/GoSearch)** | Стандартная библиотека Go (`net/http`, `net`, `sync`), Gorilla Mux, HTML templates | Параллельный веб-краулер, инвертированный индекс документов, быстрый бинарный поиск, строковое ранжирование, файловое хранение (JSON), TCP-сервер и веб-интерфейс. |
| **[lynks](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lynks)** | Clean Architecture, PostgreSQL (`pgx`), Redis / Memcached, Apache Kafka, Prometheus, Docker Compose | Распределенный сервис сокращения URL: Cache-Aside кэширование, асинхронный стриминг событий переходов в Kafka, консьюмер аналитики, метрики попаданий в кэш для Prometheus, Docker Compose окружение. |
| **[wallet](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/wallet)** | Clean Architecture, PostgreSQL (`pgxpool`), шардированный `KeyLocker`, Docker Compose, e2e тесты | Высоконагруженный финансовый REST API кошелька: поддержка **1000+ RPS на один кошелек**, сериализация операций через шардированный in-memory `KeyLocker` (4096 шардов), атомарный SQL UPDATE, защита пула соединений PostgreSQL от истощения, автоматический retry с бэкоффом, e2e тесты на состояние гонки (`-race`). |

---

## 🗺️ Карта модулей курса

| № | Тема | Каталог лекций | Домашнее задание | Ключевые концепции |
|---|---|---|---|---|
| **00** | Стандарты и Go Modules | [00-docs](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/00-docs) | — | Go modules, структура `cmd/` и `pkg/`, code review |
| **01** | Основы языка | [01-intro](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/01-intro) | — | Компиляция, пакеты, переменные, `fmt`, точка входа |
| **02** | Синтаксис и структуры памяти | [02-syntax](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/02-syntax) | [homework-02](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/02-syntax/homework-02) | Слайсы, мапы, указатели, структуры, срезы, cap/len |
| **03** | Алгоритмы | [03-algorithms](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/03-algorithms) | [homework-03](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/03-algorithms/homework-03) | Бинарный поиск, рекурсия, динамическое программирование, графы |
| **04** | Структуры данных | [04-datastructs](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/04-datastructs) | [homework-04](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/04-datastructs/homework-04) | Двусвязный список, BST (дерево поиска), массивы |
| **05** | Потоковый ввод-вывод | [05-io](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/05-io) | [homework-05](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/05-io/homework-05) | `io.Reader`, `io.Writer`, `bufio`, файлы, JSON сериализация |
| **06** | ООП и композиция | [06-oop](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/06-oop) | [homework-06](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/06-oop/homework-06) | Методы, получатели (значение/указатель), встраивание, DIP |
| **07** | Тестирование | [07-testing](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/07-testing) | [homework-07](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/07-testing/homework-07) | Unit-тесты, `TestMain`, Table-driven tests, TDD, моки |
| **08** | Профилирование и отладка | [08-prof_debug](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/08-prof_debug) | [homework-08](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/08-prof_debug/homework-08) | `pprof` (CPU/MEM), бенчмарки, `b.ResetTimer`, `go tool trace` |
| **09** | Интерфейсы и дженерики | [09-interfaces](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/09-interfaces) | [homework-09](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/09-interfaces/homework-09) | Полиморфизм, утиная типизация, Type Assertions, Generics |
| **10** | Конкурентность | [10-concurrency](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/10-concurrency) | [homework-10](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/10-concurrency/homework-10) | Горутины, каналы, `select`, `sync.WaitGroup`, `sync.Mutex`, `atomic`, `context`, Fan-Out / Fan-In |
| **11** | Сетевое взаимодействие | [11-network](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/11-network) | [homework-11](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/11-network/homework-11) | TCP сокеты, `net.Conn`, протоколы, таймауты и дедлайны |
| **12** | Веб-приложения | [12-web-apps](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/12-web-apps) | [homework-12](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/12-web-apps/homework-12) | `net/http`, маршрутизация, Gorilla Mux, HTML шаблоны (`html/template`) |
| **13** | REST API | [13-api](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/13-api) | [homework-13](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/13-api/homework-13) | RESTful API, JSON кодирование/декодирование, Middleware |
| **14** | RPC и gRPC | [14-RPC](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/14-RPC) | [homework-14](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/14-RPC/homework-14) | Go RPC (`net/rpc`), Protobuf v3, gRPC клиенты и серверы, streaming |
| **15** | Реляционные базы данных | [15-sql](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/15-sql) | [homework-15](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/15-sql/homework-15) | Реляционная модель, DDL, DML, индексы, внешние ключи, JOIN, транзакции |
| **16** | Go и СУБД | [16-db-apps](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/16-db-apps) | [homework-16](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/16-db-apps/homework-16) | `database/sql`, PostgreSQL драйвер `pgx/v4` и `pgx/v5`, пул коннектов, паттерн Репозиторий |
| **17** | Системный дизайн и архитектура | [17-system-design](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/17-system-design) | [homework-17](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/17-system-design/homework-17) | Clean Architecture, SOLID, Hexagonal / Onion, инверсия зависимостей |
| **18** | Микросервисы | [18-microservices](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/18-microservices) | [homework-18](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/18-microservices/homework-18) | 12 Factor App, Docker, Docker Compose, Healthcheck, конфигурация окружения |
| **19** | Очереди сообщений | [19-queue](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/19-queue) | [homework-19](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/19-queue/homework-19) | Apache Kafka, Event-Driven Architecture, Producer / Consumer группы, асинхронная аналитика |
| **20** | NoSQL и кэширование | [20-NoSQL](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/20-NoSQL) | [homework-20](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/20-NoSQL/homework-20) | Redis, In-Memory Cache (Cache-Aside), MongoDB, Prometheus метрики |
| **21** | Подготовка к собеседованиям | [21-interview](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/21-interview) | — | Разбор типовых ошибок, неблокирующий select, Fan-In, слайсы, UTF-8 руны |
| **22** | Современный Go и Best Practices | [22-modern-go](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/22-modern-go) | — | `log/slog`, `iter.Seq`, ServeMux, `slices`/`maps`, выравнивание памяти (padding), защита от утечек |

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
go test -v ./lessons/17-system-design/homework-17/...
go test -v ./lessons/20-NoSQL/homework-20/...
```

### 4. Запуск с детекцией гонок (Race Detector)
```bash
go test -race ./lessons/10-concurrency/homework-10/...
go test -race ./lessons/21-interview/...
```

### 5. Запуск тестов сервиса Wallet (включая стресс-тесты конкурентности)
```bash
go test -v -race ./wallet/...
```

### 6. Проверка уязвимостей сторонних зависимостей (govulncheck)
```bash
govulncheck ./...
```

---

## 📖 Рекомендуемый план самостоятельного обучения

1. **Базовый уровень (Junior Go Developer):**
   - Прочитайте [Том 1.1: Основы и синтаксис](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/01-go-core-basics.md) и [Том 1.2: Коллекции и память](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/01-go-core-data-structures.md).
   - Освойте каталоги [01-intro](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/01-intro), [02-syntax](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/02-syntax), [03-algorithms](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/03-algorithms), [04-datastructs](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/04-datastructs), [05-io](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/05-io).
   - Выполните задания [homework-02](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/02-syntax/homework-02) — [homework-05](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/05-io/homework-05).
   - Пользуйтесь настольным справочником [go-std-collections-guide.md](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/go-std-collections-guide.md).

2. **Продвинутый уровень (Middle Go Developer):**
   - Изучите [Том 2: ООП и методы](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/02-go-core-oop-methods.md) и [Том 3: Конкурентность](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/03-go-core-concurrency.md).
   - Освойте профилирование [08-prof_debug](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/08-prof_debug), интерфейсы [09-interfaces](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/09-interfaces) и конкурентность [10-concurrency](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/10-concurrency).
   - Изучите сетевое программирование [11-network](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/11-network), REST API [13-api](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/13-api), базы данных [15-sql](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/15-sql), [16-db-apps](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/16-db-apps).
   - Освойте контейнеризацию в [docker-and-compose-starter-guide.md](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/docker-and-compose-starter-guide.md).

3. **Senior / Архитектура / Собеседования:**
   - Освойте [Том 5: Тестирование, архитектура и микросервисы](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/05-go-core-testing-arch.md).
   - Разберите Clean Architecture в [homework-17](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/17-system-design/homework-17).
   - Изучите распределенные очереди Kafka в [homework-19](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/19-queue/homework-19) и кэширование с метриками в [homework-20](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/20-NoSQL/homework-20).
   - Изучите высоконагруженные сценарии в [wallet](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/wallet): шардированные мьютексы `KeyLocker`, изоляция транзакций в `pgxpool`, предотвращение deadlocks и гонок при 1000+ RPS на кошелек.
   - Освойте промышленные паттерны в [go-production-backend-guide.md](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/go-production-backend-guide.md) (Kafka Transactional Outbox, DLQ, gRPC Interceptors, AppSec).
   - Отработайте 29 практических задач по базам данных в [sql-livecoding-guide.md](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/sql-livecoding-guide.md) и проверьте себя на 14 коварных ловушках рантайма в [go-secrets-and-traps-guide.md](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/go-secrets-and-traps-guide.md).
   - Отработайте лайвкодинг по гайдам [go-livecoding-guide.md](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/go-livecoding-guide.md) и [21-interview](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/21-interview).

---

## 📚 Рекомендованная литература и первоисточники

При проектировании архитектуры учебных сервисов и составлении базы знаний использованы авторитетные профильные издания и официальные стандарты сообщества Go:

1. **Джон Боднер — «Learning Go: An Idiomatic Approach to Real-World Go Programming» (2nd Edition, O'Reilly, 2024):**  
   Эталонный учебник по современному идиоматичному Go (дженерики, итераторы, паттерны `errors.Is`, `http.ResponseController`, безопасный FFI/Cgo и культура архитектуры). Практические примеры и архитектурные решения в [go-documentation](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation) взяты и адаптированы напрямую из официального репозитория книги: [github.com/learning-go-book-2e](https://github.com/learning-go-book-2e).
2. **Алан Донован, Брайан Керниган — «Язык программирования Go» (The Go Programming Language, Addison-Wesley):**  
   Фундаментальный труд от соавторов Unix и языка C об основах дизайна языка, модели памяти и философии Go.
3. **Тейва Харшаньи — «100 ошибок Go и как их избежать» (100 Go Mistakes and How to Avoid Them, Manning):**  
   Практическое руководство по предотвращению скрытых ловушек: struct padding, утечки памяти срезов/map, гонки данных, CFS throttling в Kubernetes.
4. **Стандарты и стайлгайды ведущих инженерных команд:**  
   - [Effective Go (Официальное руководство команды Go)](https://go.dev/doc/effective_go)  
   - [Uber Go Style Guide (Стандарт промышленного кода)](https://github.com/uber-go/guide/blob/master/style.md)  
   - [Google Go Style Decisions](https://google.github.io/styleguide/go/decisions)
5. **Системный дизайн и архитектура распределенных систем:**  
   - Роберт Мартин — «Чистая архитектура» (Clean Architecture)  
   - Алекс Сюй — «System Design Interview – An Insider's Guide» (Volumes 1 & 2)  
   - Донн Мартин — [The System Design Primer](https://github.com/donnemartin/system-design-primer)
