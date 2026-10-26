# 17-system-design — Архитектурные паттерны, SOLID и Чистая Архитектура (Clean Architecture)

В данном модуле изучаются фундаментальные принципы проектирования надежного, тестируемого и масштабируемого серверного ПО на Go: принципы SOLID в контексте языка Go, методология «Чистой архитектуры» Роберта Мартина (Clean Architecture / Hexagonal Architecture), изоляция бизнес-правил через порты и адаптеры, а также анализ архитектурных компромиссов между модульным монолитом и микросервисами.

---

## 📋 Структура каталога

```text
17-system-design/
├── SOLID/                       # Примеры реализации и рефакторинга принципов SOLID
│   ├── 1-SRP/                   # Single Responsibility: srp.go, srp_refactored.go, srp_test.go
│   ├── 2-OCP/                   # Open-Closed: ocp.go, ocp_refactored.go, ocp_test.go
│   ├── 3-LSP/                   # Liskov Substitution: lsp.go, lsp_test.go
│   ├── 4-ISP/                   # Interface Segregation: isp.go, isp_test.go
│   └── 5-DIP/                   # Dependency Inversion: dip.go, dip_test.go
├── app/                         # Учебный пример сервиса каталога книг с разделением слоев
│   ├── cmd/
│   │   ├── appserver/           # Точка входа HTTP-сервера
│   │   └── appclient/           # Клиентское приложение
│   └── internal/                # Внутренние пакеты сервиса:
│       ├── api/                 # HTTP API (api.go, api_test.go) на базе gorilla/mux
│       ├── db/                  # Адаптер подключения к СУБД (pgxpool)
│       ├── models/              # Доменные модели (Book)
│       └── server/              # Точка сборки зависимостей (Composition Root)
└── homework-17/                  # Домашнее задание №17: Clean Architecture и SOLID
    ├── README.md                # Конспект книги Р. Мартина, анализ Монолит vs Микросервисы
    ├── cmd/app/main.go          # Composition Root сервиса онлайн-кинотеатра
    ├── cleanarch/               # Слои Чистой архитектуры:
    │   ├── entity/movie.go      # Ядро: корпоративные бизнес-правила (Entity)
    │   ├── usecase/             # Сценарии использования: Input/Output порты, DTO
    │   ├── adapter/             # Адаптеры доставки (HTTP) и хранилища (In-Memory)
    │   └── cleanarch_test.go    # E2E тесты и бенчмарки прохождения через все слои
    └── solid/                   # Полные практические примеры SRP, OCP, LSP, ISP, DIP
```

---

## 🔍 Ключевые концепции

1. **Принципы SOLID в контексте Go:**
   - **SRP (Single Responsibility Principle):** Пакет, структура или функция должны отвечать перед одним и только одним актором. Код, меняющийся по разным причинам, разносится по разным пакетам.
   - **OCP (Open-Closed Principle):** Сущности открыты для расширения, но закрыты для модификации. В Go достигается через полиморфизм интерфейсов (паттерны Specification, Strategy).
   - **LSP (Liskov Substitution Principle):** Поведение интерфейса гарантирует, что любая конкретная реализация не нарушит инвариантов вызывающего кода.
   - **ISP (Interface Segregation Principle):** «Чем больше интерфейс, тем слабее абстракция». Клиенты не должны зависеть от методов, которые не используют. Объявляйте узкие интерфейсы на стороне потребителя (1–2 метода).
   - **DIP (Dependency Inversion Principle):** Модули высокого уровня не зависят от модулей низкого уровня; оба зависят от абстракций. Зависимости внедряются через параметры конструкторов (`Dependency Injection`).

2. **Clean Architecture (Правило зависимостей — The Dependency Rule):**
   - Зависимости исходного кода направлены **только внутрь**:
     - Внешний слой: транспорт (HTTP/gRPC) и хранилище (PostgreSQL/Redis) — это детали реализации.
     - Промежуточный слой: Адаптеры (Контроллеры, Репозитории, Шлюзы).
     - Внутренний слой: Сценарии использования (Use Cases) и Сущности ядра (Entities).
   - Бизнес-правила ничего не знают о фреймворках, СУБД или библиотеках сериализации.

3. **Монолит vs Микросервисы:**
   - Рекомендация разработки: практически всегда начинать проект с **модульного монолита** с чистой архитектурой.
   - Чистая архитектура внутри монолита обеспечивает строгие границы между доменами, позволяя при необходимости безболезненно вынести домен в отдельный микросервис.

---

## 🚀 Запуск тестов и бенчмарков

```bash
# Запуск всех тестов модуля 17 с проверкой состояний гонки
go test -v -race ./lessons/17-system-design/...

# Запуск тестов с бенчмарками
go test -v -race -bench=. ./lessons/17-system-design/...
```

---

## 🔗 Файлы и связанные материалы
- [homework-17/README.md](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/17-system-design/homework-17/README.md) — Домашнее задание №17: Конспект и эталонная архитектура
- [homework-17/cleanarch/entity/movie.go](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/17-system-design/homework-17/cleanarch/entity/movie.go) — Доменные сущности (Entities)
- [homework-17/cleanarch/usecase/movie_usecase.go](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/17-system-design/homework-17/cleanarch/usecase/movie_usecase.go) — Сценарии использования и порты (Use Cases)
- [homework-17/cleanarch/cleanarch_test.go](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/17-system-design/homework-17/cleanarch/cleanarch_test.go) — Тесты и бенчмарки чистой архитектуры
- [app/internal/api/api.go](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/17-system-design/app/internal/api/api.go) — Пример HTTP API сервиса
- **Теория:** [Том 5: Тестирование и архитектура (Разделы 21–22: SOLID, Clean Architecture)](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/go-documentation/05-go-core-testing-arch.md)
- **Гайд по системному дизайну:** [go-system-design-guide.md](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/go-documentation/go-system-design-guide.md)
