# Go Production Engineering & Knowledge Base | Инженерное портфолио и база знаний Go

> **Highload (1000+ RPS) • Распределенные системы • 7 томов Go Core • Linux & System Design • 85+ задач с собеседований • Modern Go 1.25+**

[![CI](https://github.com/Alex808r/go-production-engineering/actions/workflows/ci.yml/badge.svg)](https://github.com/Alex808r/go-production-engineering/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/Alex808r/go-production-engineering)](https://goreportcard.com/report/github.com/Alex808r/go-production-engineering)
[![Go Version](https://img.shields.io/badge/Go-1.22%20%7C%201.23%20%7C%201.24%20%7C%201.25%2B-00ADD8?logo=go)](https://go.dev/)
[![Highload](https://img.shields.io/badge/Highload-1000%2B%20RPS-success?logo=speedtest)](wallet/)
[![Agentic AI](https://img.shields.io/badge/Agentic%20AI-Ready%20(AGENTS.md)-blueviolet)](AGENTS.md)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> [!IMPORTANT]
> 🎓 **О происхождении репозитория и авторском вкладе:**  
> Репозиторий изначально брал за основу программу и базовые учебные заготовки курса **«Современная разработка на Go»** (онлайн-школа [Thinknetica](https://thinknetica.com/), оригинальный репозиторий курса: [thinknetica/go_course_4](https://github.com/thinknetica/go_course_4), автор курса — Дмитрий Титов).  
> 
> Однако в процессе обучения и глубокого погружения проект перерос в **масштабную самостоятельную работу и полноценную базу инженерных знаний**:
> - Все решения практических задач и домашних заданий (`homework-02` — `homework-20`), а также сквозные проекты **GoSearch**, **lynks** и высоконагруженный финансовый сервис **wallet** полностью спроектированы, написаны, отлажены и покрыты тестами **самостоятельно**.
> - Исходные лекционные примеры были существенно переработаны, дополнены подробными комментариями к механизмам рантайма и стандартам production-кода.
> - Создана фундаментальная авторская энциклопедия **[go-documentation](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/go-documentation)** (7 томов глубокой теории Go Core, руководства по Modern Go 1.21–1.26+, Docker & Compose, ядру Linux, лайвкодингу, базам данных, System Design, архитектурным паттернам Enterprise и разборам рантайм-ловушек), выходящая далеко за рамки исходной программы курса.
> - **Использование материалов книги Джона Боднера:** При написании и расширении энциклопедии [go-documentation](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/go-documentation) использованы теоретические материалы и практические примеры из книги **Джона Боднера «Learning Go: An Idiomatic Approach to Real-World Go Programming» (2nd Edition, O'Reilly, 2024)**, а канонические реализации и архитектурные паттерны взяты напрямую из официального репозитория книги: **[learning-go-book-2e (GitHub)](https://github.com/learning-go-book-2e)** с сохранением атрибуции первоисточника.
> - Теоретическая база также систематизирует и цитирует эталонные практики других мировых источников: классики **Алана Донована и Брайана Кернигана (*The Go Programming Language*)**, стайлгайдов **Uber Go** и **Google**, а также разборов ошибок **Тейвы Харшаньи (*100 Go Mistakes*)**.

Комплексный инженерный репозиторий: витрина практических навыков разработки высоконагруженных систем на Go, библиотека продакшен-микросервисов и фундаментальная база знаний для подготовки к техническим собеседованиям (Middle / Senior Go Developer в BigTech).

Репозиторий включает:
- **22 тематических модуля** с примерами кода от базового синтаксиса до распределенных систем в каталоге [lessons](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons) (`00-docs` — `22-modern-go`);
- **19 практических домашних заданий** с полным тестовым покрытием (`homework-02` — `homework-20`), реализованных самостоятельно с нуля;
- **3 сквозных продакшн-проекта**:
  - [GoSearch](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/GoSearch) — многопоточный поисковый движок (краулер, инвертированный индекс, бинарный поиск, TCP и REST API, веб-интерфейс);
  - [lynks](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lynks) — распределенный микросервис сокращения ссылок (Clean Architecture, PostgreSQL, Redis/Memcached, Kafka/Event-Driven, Prometheus метрики, Docker Compose);
  - [wallet](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/wallet) — высоконагруженный финансовый REST-сервис управления балансом (1000+ RPS на кошелек, шардированный `KeyLocker`, пул соединений `pgxpool`, строгая изоляция транзакций, идемпотентность, Docker Compose, e2e стресс-тесты под нагрузкой);
- **Фундаментальную базу знаний** в каталоге [go-documentation](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/go-documentation): 7 подробных томов теории (включая материалы и практические реализации из книги Джона Боднера *«Learning Go 2nd Edition»* и её репозитория [learning-go-book-2e](https://github.com/learning-go-book-2e)), шпаргалки, настольный справочник по стандартным пакетам и коллекциям, гайды по Docker & Compose, ядру Linux и системному программированию, базам данных, System Design, Enterprise-стандартам (Kafka Outbox/DLQ, gRPC, CI/CD, AppSec), лайвкодингу (52 задачи по Go, 33 задачи по SQL, 9 квизов рантайма) и 14 ловушкам рантайма;
- **Промышленный регламент работы с AI-ассистентами** ([AGENTS.md](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/AGENTS.md) и [CLAUDE.md](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/CLAUDE.md)) для стандартизации Agentic Pair Programming.

---

## 📚 Структура репозитория

Каждый тематический каталог модуля внутри директории [`lessons/`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons) (`02-syntax`, `03-algorithms` и т.д.) является **самодостаточным**: он содержит лекционные материалы, конспект темы и вложенное практическое домашнее задание (`homework-XX`).

```text
go-course-4/
├── AGENTS.md                   # Промышленный регламент для AI-ассистентов (Agentic Engineering)
├── CLAUDE.md                   # Точка входа для Claude Code (ссылка на AGENTS.md)
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
│   ├── 05-go-core-testing-arch.md          # Том 5: Тестирование, Clean Arch, gRPC streaming, Kafka/NATS, K8s, SLA/SLO
│   ├── 06-go-core-algorithms-interview.md  # Том 6: Алгоритмы, структуры данных, unsafe/reflect, Swiss Tables, STAR, ATS
│   ├── course-codebase-guide.md            # Путеводитель по кодовой базе репозитория
│   ├── docker-and-compose-starter-guide.md # Практическое руководство по Docker и Docker Compose
│   ├── go-core-cheatsheet.md               # Быстрая шпаргалка по Go
│   ├── go-database-interview-guide.md      # Вопросы и практика по SQL и NoSQL (блокировки, транзакции, MVCC)
│   ├── go-design-patterns-guide.md         # Паттерны проектирования в Go (GoF, Concurrency, Highload, антипаттерны)
│   ├── go-livecoding-guide.md              # Задачи для лайвкодинга с разбором (52 задачи)
│   ├── go-modern-features-guide.md         # Руководство по фичам Go 1.21–1.26+ (slog, iter, synctest, Swiss Tables)
│   ├── go-os-and-linux-interview-guide.md  # Гайд по ОС, Ядру Linux, Сетям и Системному программированию
│   ├── go-production-backend-guide.md      # Промышленный backend: CI/CD, Kafka Outbox/DLQ, gRPC, AppSec
│   ├── go-secrets-and-traps-guide.md       # 14 ловушек и секретов Go (память, утечки горутин, typed nil)
│   ├── go-std-collections-guide.md         # Справочник по стандартным коллекциям и пакетам Go
│   ├── go-system-design-guide.md           # Проектирование распределенных систем
│   ├── homework-tasks.md                   # Формулировки домашних заданий
│   ├── sql-livecoding-guide.md             # Практика SQL запросов (33 задачи от базовых до LeetCode Hard)
│   └── mermaid-diagrams/                  # Стандарты визуализации архитектуры Mermaid
└── README.md                   # Главный путеводитель (этот файл)
```

---

## 🚀 Сквозные продакшн-проекты

В репозитории реализованы 3 полноценных продакшн-сервиса, демонстрирующих реальные сценарии проектирования:

| Проект | Архитектура и стек | Ключевые инженерные решения |
|---|---|---|
| **[GoSearch](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/GoSearch)** | Стандартная библиотека Go (`net/http`, `net`, `sync`), Gorilla Mux, HTML templates | Параллельный веб-краулер, инвертированный индекс документов, быстрый бинарный поиск, строковое ранжирование, файловое хранение (JSON), TCP-сервер и веб-интерфейс. |
| **[lynks](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lynks)** | Clean Architecture, PostgreSQL (`pgx`), Redis / Memcached, Apache Kafka, Prometheus, Docker Compose | Распределенный сервис сокращения URL: Cache-Aside кэширование, асинхронный стриминг событий переходов в Kafka, консьюмер аналитики, метрики попаданий в кэш для Prometheus, Docker Compose окружение. |
| **[wallet](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/wallet)** | Clean Architecture, PostgreSQL (`pgxpool`), шардированный `KeyLocker`, Docker Compose, e2e тесты | Высоконагруженный финансовый REST API кошелька: поддержка **1000+ RPS на один кошелек**, сериализация операций через шардированный in-memory `KeyLocker` (4096 шардов), атомарный SQL UPDATE, защита пула соединений PostgreSQL от истощения, автоматический retry с бэкоффом, e2e тесты на состояние гонки (`-race`). |

---

## 🎯 Инженерное портфолио: Матрица компетенций и подтверждения в коде

Данный репозиторий спроектирован как **демонстрация реальной инженерной квалификации (Middle+ / Senior Go Developer)** для технических интервьюеров, тимлидов и архитекторов:

| Инженерная компетенция | Продакшн-реализация в репозитории | Ключевые архитектурные решения и пруфы в кодовой базе |
|---|---|---|
| **Highload & Concurrency (1000+ RPS)** | [wallet/](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/wallet) | Шардированный in-memory мьютекс `KeyLocker` (4096 шардов), предотвращение гонок на балансе кошелька при 1000 RPS, атомарный SQL UPDATE, защита пула соединений PostgreSQL от истощения, e2e стресс-тесты параллелизма под `-race`.<br>🔗 **Код:** [`wallet/internal/sync`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/wallet/internal/sync), [`wallet/tests`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/wallet/tests), [`wallet/Makefile`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/wallet/Makefile) |
| **Clean Architecture & Enterprise** | [wallet/internal](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/wallet/internal), [lynks/](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lynks) | Строгая 4-слойная Clean Architecture (`Delivery` → `Service` → `Repository` → `Domain`), соблюдение принципа инверсии зависимостей (DIP), доменные модели без инфраструктурных утечек.<br>🔗 **Код:** [`wallet/internal`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/wallet/internal), [`lessons/17-system-design/homework-17`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/17-system-design/homework-17) |
| **Распределенные системы & Kafka** | [lynks/](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lynks), [lessons/19-queue](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/19-queue) | Асинхронная шина событий Apache Kafka (Producer/Consumer группы, ребалансировка, оффсеты, Consumer Lag), паттерн Transactional Outbox (решение проблемы Dual Write), двухуровневый Cache-Aside с Redis.<br>🔗 **Код:** [`lynks/`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lynks), [`lessons/19-queue/homework-19`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/19-queue/homework-19), [Том 5: разд. 21.3](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/go-documentation/05-go-core-testing-arch.md#213-transactional-outbox-pattern-решение-проблемы-dual-write) |
| **Поисковые алгоритмы & Low-level Go** | [GoSearch/](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/GoSearch) | Параллельный краулер веб-страниц, инвертированный индекс документов (`map[string][]int`), строковый бинарный поиск, выравнивание структур (Memory padding), zero-copy преобразования через `unsafe`.<br>🔗 **Код:** [`GoSearch/pkg`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/GoSearch/pkg), [`lessons/22-modern-go/07-secrets-and-traps`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/22-modern-go/07-secrets-and-traps) |
| **Observability, Метрики & SRE** | [lessons/08-prof_debug](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/08-prof_debug), [lynks](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lynks) | Мониторинг надежности (SLA, SLO, SLI, перцентили $p50..p99.9$, гистограммы Prometheus `_bucket` и `histogram_quantile`), OpenTelemetry трассировка, структурированные логи `log/slog`, `pprof` CPU/MEM/Block/Goroutine.<br>🔗 **Код:** [`lessons/08-prof_debug`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/08-prof_debug), [`lessons/20-NoSQL/homework-20`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/20-NoSQL/homework-20), [Том 5: разд. 22.1](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/go-documentation/05-go-core-testing-arch.md#221-метрики-prometheus) |
| **Linux Internals, Docker & K8s** | [go-os-and-linux-guide](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/go-documentation/go-os-and-linux-interview-guide.md) | Multi-stage Dockerfile (Scratch образ до 15 МБ), Kubernetes Probes (Startup, Liveness, Readiness), защита от CFS CPU Throttling (`automaxprocs`) и OOM Killer (`GOMEMLIMIT`), устройство виртуальной памяти Linux, TLB, Hugepages, сокеты, сисколы.<br>🔗 **Код:** [`wallet/Dockerfile`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/wallet/Dockerfile), [`go-os-and-linux-interview-guide.md`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/go-documentation/go-os-and-linux-interview-guide.md) |
| **Modern Go (Go 1.22–1.26+)** | [lessons/22-modern-go](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/22-modern-go) | Swiss Tables (Go 1.24+), `testing/synctest` детерминированное тестирование конкурентности (Go 1.25+), Range-over-func (`iter.Seq`), `slices`/`maps`, пакеты `unique` и `weak`.<br>🔗 **Код:** [`lessons/22-modern-go`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/22-modern-go), [go-modern-features-guide.md](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/go-documentation/go-modern-features-guide.md) |
| **Agentic AI Engineering (Тренды 2026)** | [AGENTS.md](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/AGENTS.md), [CLAUDE.md](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/CLAUDE.md) | Промышленный регламент парного программирования с AI-ассистентами: жесткий Git Safety (запрет коммитов без согласования), Code Evidence, мультимодульная верификация.<br>🔗 **Код:** [`AGENTS.md`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/AGENTS.md), [`CLAUDE.md`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/CLAUDE.md) |

---

## ⚡ Быстрый старт для технического ревьюера (Reviewer Quickstart)

Для быстрой проверки работоспособности кодовой базы и запуска продакшн-сервисов:

```bash
# 1. Запуск всех unit-тестов проекта (100% PASS)
go test ./...

# 2. Запуск тестов финансового микросервиса Wallet с детектором гонок данных (-race)
cd wallet && make test-race

# 3. Стресс-тест параллелизма Wallet под нагрузкой 1000+ RPS на кошелек
cd wallet && make test-concurrency

# 4. Поднятие микросервисов в Docker Compose (PostgreSQL + автомиграции + REST API)
cd wallet && make docker-up

# 5. Просмотр логов сервиса в реальном времени
cd wallet && make docker-logs
```

---

## 🗺️ Карта модулей курса

| № | Тема | Каталог лекций | Домашнее задание | Ключевые концепции |
|---|---|---|---|---|
| **00** | Стандарты и Go Modules | [00-docs](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/00-docs) | — | Go modules, структура `cmd/` и `pkg/`, code review |
| **01** | Основы языка | [01-intro](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/01-intro) | — | Компиляция, пакеты, переменные, `fmt`, точка входа |
| **02** | Синтаксис и структуры памяти | [02-syntax](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/02-syntax) | [homework-02](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/02-syntax/homework-02) | Слайсы, мапы, указатели, структуры, срезы, cap/len |
| **03** | Алгоритмы | [03-algorithms](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/03-algorithms) | [homework-03](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/03-algorithms/homework-03) | Бинарный поиск, рекурсия, динамическое программирование, графы |
| **04** | Структуры данных | [04-datastructs](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/04-datastructs) | [homework-04](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/04-datastructs/homework-04) | Двусвязный список, BST (дерево поиска), массивы |
| **05** | Потоковый ввод-вывод | [05-io](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/05-io) | [homework-05](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/05-io/homework-05) | `io.Reader`, `io.Writer`, `bufio`, файлы, JSON сериализация |
| **06** | ООП и композиция | [06-oop](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/06-oop) | [homework-06](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/06-oop/homework-06) | Методы, получатели (значение/указатель), встраивание, DIP |
| **07** | Тестирование | [07-testing](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/07-testing) | [homework-07](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/07-testing/homework-07) | Unit-тесты, `TestMain`, Table-driven tests, TDD, моки |
| **08** | Профилирование и отладка | [08-prof_debug](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/08-prof_debug) | [homework-08](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/08-prof_debug/homework-08) | `pprof` (CPU/MEM), бенчмарки, `b.ResetTimer`, `go tool trace` |
| **09** | Интерфейсы и дженерики | [09-interfaces](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/09-interfaces) | [homework-09](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/09-interfaces/homework-09) | Полиморфизм, утиная типизация, Type Assertions, Generics |
| **10** | Конкурентность | [10-concurrency](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/10-concurrency) | [homework-10](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/10-concurrency/homework-10) | Горутины, каналы, `select`, `sync.WaitGroup`, `sync.Mutex`, `atomic`, `context`, Fan-Out / Fan-In |
| **11** | Сетевое взаимодействие | [11-network](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/11-network) | [homework-11](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/11-network/homework-11) | TCP сокеты, `net.Conn`, протоколы, таймауты и дедлайны |
| **12** | Веб-приложения | [12-web-apps](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/12-web-apps) | [homework-12](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/12-web-apps/homework-12) | `net/http`, маршрутизация, Gorilla Mux, HTML шаблоны (`html/template`) |
| **13** | REST API | [13-api](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/13-api) | [homework-13](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/13-api/homework-13) | RESTful API, JSON кодирование/декодирование, Middleware |
| **14** | RPC и gRPC | [14-RPC](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/14-RPC) | [homework-14](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/14-RPC/homework-14) | Go RPC (`net/rpc`), Protobuf v3, gRPC клиенты и серверы, streaming |
| **15** | Реляционные базы данных | [15-sql](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/15-sql) | [homework-15](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/15-sql/homework-15) | Реляционная модель, DDL, DML, индексы, внешние ключи, JOIN, транзакции |
| **16** | Go и СУБД | [16-db-apps](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/16-db-apps) | [homework-16](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/16-db-apps/homework-16) | `database/sql`, PostgreSQL драйвер `pgx/v4` и `pgx/v5`, пул коннектов, паттерн Репозиторий |
| **17** | Системный дизайн и архитектура | [17-system-design](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/17-system-design) | [homework-17](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/17-system-design/homework-17) | Clean Architecture, SOLID, Hexagonal / Onion, инверсия зависимостей |
| **18** | Микросервисы | [18-microservices](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/18-microservices) | [homework-18](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/18-microservices/homework-18) | 12 Factor App, Docker, Docker Compose, Healthcheck, конфигурация окружения |
| **19** | Очереди сообщений | [19-queue](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/19-queue) | [homework-19](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/19-queue/homework-19) | Apache Kafka, Event-Driven Architecture, Producer / Consumer группы, асинхронная аналитика |
| **20** | NoSQL и кэширование | [20-NoSQL](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/20-NoSQL) | [homework-20](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/20-NoSQL/homework-20) | Redis, In-Memory Cache (Cache-Aside), MongoDB, Prometheus метрики |
| **21** | Подготовка к собеседованиям | [21-interview](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/21-interview) | — | Разбор типовых ошибок, неблокирующий select, Fan-In, слайсы, UTF-8 руны |
| **22** | Современный Go и Best Practices | [22-modern-go](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/22-modern-go) | — | `log/slog`, `iter.Seq`, ServeMux, `slices`/`maps`, выравнивание памяти (padding), защита от утечек |

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

### 5. Запуск тестов финансового сервиса Wallet (самостоятельный модуль)
```bash
# Из корня переходим в каталог wallet:
cd wallet && make test-race
# или прямой запуск через go test:
cd wallet && go test -v -race ./...
```

### 6. Проверка уязвимостей сторонних зависимостей (govulncheck)
```bash
go run golang.org/x/vuln/cmd/govulncheck@latest ./...
```

---

## 📖 Рекомендуемый план самостоятельного обучения

1. **Базовый уровень (Junior Go Developer):**
   - Прочитайте [Том 1.1: Основы и синтаксис](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/go-documentation/01-go-core-basics.md) и [Том 1.2: Коллекции и память](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/go-documentation/01-go-core-data-structures.md).
   - Освойте каталоги [01-intro](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/01-intro), [02-syntax](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/02-syntax), [03-algorithms](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/03-algorithms), [04-datastructs](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/04-datastructs), [05-io](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/05-io).
   - Выполните задания [homework-02](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/02-syntax/homework-02) — [homework-05](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/05-io/homework-05).
   - Пользуйтесь настольным справочником [go-std-collections-guide.md](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/go-documentation/go-std-collections-guide.md).

2. **Продвинутый уровень (Middle Go Developer):**
   - Изучите [Том 2: ООП и методы](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/go-documentation/02-go-core-oop-methods.md) и [Том 3: Конкурентность](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/go-documentation/03-go-core-concurrency.md).
   - Освойте профилирование [08-prof_debug](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/08-prof_debug), интерфейсы [09-interfaces](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/09-interfaces) и конкурентность [10-concurrency](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/10-concurrency).
   - Изучите сетевое программирование [11-network](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/11-network), REST API [13-api](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/13-api), базы данных [15-sql](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/15-sql), [16-db-apps](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/16-db-apps).
   - Освойте контейнеризацию в [docker-and-compose-starter-guide.md](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/go-documentation/docker-and-compose-starter-guide.md).

3. **Senior / Архитектура / Собеседования:**
   - Освойте [Том 5: Тестирование, архитектура и микросервисы](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/go-documentation/05-go-core-testing-arch.md).
   - Разберите Clean Architecture в [homework-17](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/17-system-design/homework-17).
   - Изучите распределенные очереди Kafka в [homework-19](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/19-queue/homework-19) и кэширование с метриками в [homework-20](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/20-NoSQL/homework-20).
   - Изучите высоконагруженные сценарии в [wallet](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/wallet): шардированные мьютексы `KeyLocker`, изоляция транзакций в `pgxpool`, предотвращение deadlocks и гонок при 1000+ RPS на кошелек.
   - Освойте промышленные паттерны в [go-production-backend-guide.md](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/go-documentation/go-production-backend-guide.md) (Kafka Transactional Outbox, DLQ, gRPC Interceptors, AppSec).
   - Освойте системное программирование и ядро Linux в [go-os-and-linux-interview-guide.md](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/go-documentation/go-os-and-linux-interview-guide.md) (виртуальная память, TLB, Hugepages, CFS Scheduler, сокеты, top 146%).
   - Отработайте 33 практические задачи по базам данных в [sql-livecoding-guide.md](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/go-documentation/sql-livecoding-guide.md) и проверьте себя на 14 коварных ловушках рантайма в [go-secrets-and-traps-guide.md](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/go-documentation/go-secrets-and-traps-guide.md).
   - Отработайте лайвкодинг по гайду [go-livecoding-guide.md (52 задачи с разбором)](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/go-documentation/go-livecoding-guide.md), решите 9 квизов-ловушек рантайма и изучите стратегию прохождения воронки в [Том 6: разд. 36.6](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/go-documentation/06-go-core-algorithms-interview.md#366-карьерная-стратегия-воронка-headhunter-atsai-фильтры-резюме-самопрезентация-мотивационные-вопросы-и-торг-за-оффер).
   - Применяйте современные инженерные стандарты парного AI-программирования из [AGENTS.md](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/AGENTS.md).

---

## 📚 Рекомендованная литература и первоисточники

При проектировании архитектуры учебных сервисов и составлении базы знаний использованы авторитетные профильные издания и официальные стандарты сообщества Go:

1. **Джон Боднер — «Learning Go: An Idiomatic Approach to Real-World Go Programming» (2nd Edition, O'Reilly, 2024):**  
   Эталонный учебник по современному идиоматичному Go (дженерики, итераторы, паттерны `errors.Is`, `http.ResponseController`, безопасный FFI/Cgo и культура архитектуры). Практические примеры и архитектурные решения в [go-documentation](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/go-documentation) взяты и адаптированы напрямую из официального репозитория книги: [github.com/learning-go-book-2e](https://github.com/learning-go-book-2e).
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
