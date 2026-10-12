# Домашнее задание №17: Архитектура Go-приложения (Clean Architecture и SOLID)

Решение домашнего задания по курсу «Современная разработка на Go» (Занятие №17: Архитектура Go-приложения).  
Реализовано в отдельной изолированной директории `homework-17` без изменения или затирания предыдущих решений.

---

## 📑 Содержание

1. [Постановка задачи](#постановка-задачи)
2. [Конспект и анализ книги Р. Мартина «Чистая архитектура»](#конспект-и-анализ-книги-р-мартина-чистая-архитектура)
   - [Что такое архитектура ПО](#что-такое-архитектура-по)
   - [Парадигмы программирования](#парадигмы-программирования)
   - [Принципы SOLID в контексте Go](#принципы-solid-в-контексте-go)
   - [Принципы организации компонентов](#принципы-организации-компонентов)
   - [Правило зависимостей (The Dependency Rule)](#правило-зависимостей-the-dependency-rule)
   - [Кричащая архитектура (Screaming Architecture)](#кричащая-архитектура-screaming-architecture)
   - [База данных и Веб — это детали реализации](#база-данных-и-веб--это-детали-реализации)
3. [Практическая реализация SOLID в Go (`solid/`)](#практическая-реализация-solid-в-go-solid)
   - [SRP: Single Responsibility Principle](#srp-single-responsibility-principle)
   - [OCP: Open-Closed Principle](#ocp-open-closed-principle)
   - [LSP: Liskov Substitution Principle](#lsp-liskov-substitution-principle)
   - [ISP: Interface Segregation Principle](#isp-interface-segregation-principle)
   - [DIP: Dependency Inversion Principle](#dip-dependency-inversion-principle)
4. [Эталонная реализация Clean Architecture в Go (`cleanarch/`)](#эталонная-реализация-clean-architecture-в-go-cleanarch)
   - [Структура слоев](#структура-слоев)
   - [Сущности (Entities)](#сущности-entities)
   - [Сценарии использования (Use Cases & Ports)](#сценарии-использования-use-cases--ports)
   - [Адаптеры (Adapters: Repository & HTTP Controller)](#адаптеры-adapters-repository--http-controller)
   - [Точка сборки (Composition Root)](#точка-сборки-composition-root)
5. [Анализ: Монолит vs Микросервисы](#анализ-монолит-vs-микросервисы)
6. [Инструкция по запуску и тесты](#инструкция-по-запуску-и-тесты)

---

## Постановка задачи

- **Задание**: Изучить книгу Роберта Мартина (Дядюшки Боба) «Чистая архитектура. Руководство для профессионала по структуре и проектированию ПО» (*Clean Architecture: A Craftsman's Guide to Software Structure and Design*), материалы курса и статьи Мартина Фаулера.
- **Цель**: Применить принципы Clean Architecture и SOLID при проектировании архитектуры Go-приложений, исключить антипаттерны (обезличенные слои MVC, чрезмерную зависимость от фреймворков и СУБД) и реализовать эталонную архитектуру с полным покрытием тестами.

---

## Конспект и анализ книги Р. Мартина «Чистая архитектура»

### Что такое архитектура ПО
Архитектура программной системы — это форма, которую придают системе её создатели. Главная цель архитектуры — **минимизировать человеческие трудозатраты на создание и сопровождение системы**.
Хороший архитектор откладывает принятие решений, касающихся деталей (выбор СУБД, веб-фреймворка, протоколов связи), оставляя максимальное количество опций открытыми как можно дольше.

### Парадигмы программирования
Роберт Мартин показывает, что все парадигмы программирования не дают новых возможностей, а **накладывают ограничения**:
1. **Структурное программирование**: дисциплинирует передачу управления (устранение `goto` в пользу циклов и ветвлений).
2. **Объектно-ориентированное программирование**: дисциплинирует косвенную передачу управления с помощью полиморфизма (инверсия зависимостей).
3. **Функциональное программирование**: дисциплинирует присваивание (неизменяемость состояния, отсутствие side effects).

### Принципы SOLID в контексте Go
- **S (SRP)**: Модуль должен отвечать перед одним и только одним актором (источником требований). В Go это означает: код, который меняется вместе по одной причине, группируется в один пакет; код, меняющийся по разным причинам — разносится.
- **O (OCP)**: Программные сущности открыты для расширения, но закрыты для модификации. В Go достигается через интерфейсы и комбинацию строительных блоков (паттерны Specification, Strategy, Decorator).
- **L (LSP)**: Если для каждого объекта $o_1$ типа $S$ существует объект $o_2$ типа $T$ такой, что для всех программ поведение не изменяется при замене $o_1$ на $o_2$, то $S$ — подтип $T$. В Go достигается утиной типизацией (Duck Typing) интерфейсов при сохранении инвариантов и контрактов.
- **I (ISP)**: Клиенты не должны зависеть от интерфейсов, которые они не используют. Правило Go: *"The bigger the interface, the weaker the abstraction"*. Создавать узкие интерфейсы (1-2 метода, как `io.Reader`), объявляя их на стороне потребителя.
- **D (DIP)**: Модули высокого уровня не зависят от модулей низкого уровня. Оба зависят от абстракций. В Go: зависимости передаются через параметры конструкторов (`NewService(repo Repository)`).

### Принципы организации компонентов
- **REP** (Reuse/Release Equivalence Principle) — единица повторного использования эквивалентна единице выпуска.
- **CCP** (Common Closure Principle) — в один компонент собираются классы/пакеты, меняющиеся по одной причине.
- **CRP** (Common Reuse Principle) — не заставляйте клиентов зависеть от того, что им не нужно.
- **ADP** (Acyclic Dependencies Principle) — граф зависимостей пакетов должен быть ациклическим (в Go циклические импорты запрещены на уровне компилятора!).
- **SDP** (Stable Dependencies Principle) — модуль должен зависеть только от более стабильных модулей.
- **SAP** (Stable Abstractions Principle) — стабильный модуль должен быть абстрактным (состоять из интерфейсов).

### Правило зависимостей (The Dependency Rule)
Исходный код может иметь зависимости, направленные **только внутрь**, в сторону высокоуровневых политик:
```
  [ Внешний мир: Web, DB, UI, Очереди ]
                  │
                  ▼  (зависит от)
     [ Адаптеры: Controllers, Gateways ]
                  │
                  ▼  (зависит от)
       [ Сценарии использования: Use Cases ]
                  │
                  ▼  (зависит от)
          [ Сущности ядра: Entities ]
```
Внутренние слои **ничего не знают** о внешних. Если сценарию использования нужно сохранить сущность в БД, он объявляет интерфейс репозитория (Output Port) внутри слоя Use Cases, а конкретный адаптер БД реализует этот интерфейс снаружи.

### Кричащая архитектура (Screaming Architecture)
Структура пакетов должна «кричать» о предметной области системы (кинопрокат, медицина, бухгалтерия), а не о библиотеках (`controllers`, `models`, `views`). Использование плоских технических слоев MVC — антипаттерн для крупных сервисов.

### База данных и Веб — это детали реализации
- **База данных** — это механизм персистентности данных (деталь). Бизнес-правилам все равно, где лежат данные: в PostgreSQL, MongoDB, памяти или файле.
- **Веб / HTTP / gRPC** — это транспорт (механизм ввода-вывода). Сценарий не должен знать про `http.Request` или `http.ResponseWriter`.

---

## Практическая реализация SOLID в Go (`solid/`)

Все 5 принципов реализованы и покрыты тестами в каталоге [`homework-17/solid`](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/17-system-design/homework-17/solid):

### SRP: Single Responsibility Principle
- Код: [`solid/srp/srp.go`](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/17-system-design/homework-17/solid/srp/srp.go)
- Бизнес-логика каталога фильмов (`MovieService`) отделена от форматирования отчетов (`ReportFormatter`).
- Изменение формата вывода отчета (например, добавление CSV) не затрагивает логику управления фильмами.

### OCP: Open-Closed Principle
- Код: [`solid/ocp/ocp.go`](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/17-system-design/homework-17/solid/ocp/ocp.go)
- Функция `FilterMovies(movies, spec)` закрыта для модификации, но открыта для добавления новых фильтров (по году, по рейтингу, по кассовым сборам, логическое «И» через `AndSpecification`).

### LSP: Liskov Substitution Principle
- Код: [`solid/lsp/lsp.go`](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/17-system-design/homework-17/solid/lsp/lsp.go)
- Сервис `NotificationService` ожидает контракт `Notifier`. Реализации `EmailNotifier` и `SMSNotifier` строго соблюдают инварианты поведения без побочных эффектов.

### ISP: Interface Segregation Principle
- Код: [`solid/isp/isp.go`](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/17-system-design/homework-17/solid/isp/isp.go)
- Вместо гигантского «божественного» интерфейса объявлены узкие интерфейсы: `Reader`, `Writer`, `Deleter`. Сервис чтения `ReadOnlyService` зависит только от `Reader`, не требуя заглушек для методов записи и удаления.

### DIP: Dependency Inversion Principle
- Код: [`solid/dip/dip.go`](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/17-system-design/homework-17/solid/dip/dip.go)
- Сервис бронирования билетов `TicketBookingService` зависит от абстракций `PaymentGateway` и `TicketRepository`. Реализации моков и продакшен-адаптеров передаются через Dependency Injection.

---

## Эталонная реализация Clean Architecture в Go (`cleanarch/`)

Каталог [`homework-17/cleanarch`](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/17-system-design/homework-17/cleanarch) содержит чистую структуру:

```
homework-17/
├── cleanarch/
│   ├── entity/
│   │   └── movie.go                 # Enterprise Business Rules (0 внешних импортов)
│   ├── usecase/
│   │   └── movie_usecase.go         # Application Business Rules, Input/Output Ports, DTO
│   ├── adapter/
│   │   ├── repository/
│   │   │   └── memory_repo.go       # Secondary Adapter (Gateway к хранилищу)
│   │   └── controller/
│   │       └── http_controller.go   # Primary Adapter (HTTP Delivery)
│   └── cleanarch_test.go            # E2E тесты взаимодействия слоев
├── solid/                           # Демонстрации SOLID на Go
└── cmd/app/main.go                  # Composition Root
```

### 1. Сущности (`cleanarch/entity/movie.go`):
- Содержит `Movie` и фабричный метод `NewMovie(...)` с проверкой бизнес-правил (год $\ge 1800$, валидный рейтинг PG-10/13/18, неотрицательные сборы).
- **0 зависимостей** от БД или внешних пакетов.

### 2. Сценарии использования (`cleanarch/usecase/movie_usecase.go`):
- Объявляет входной порт `MovieUseCase` (бизнес-операции: создание, пакетное создание, получение, фильтрация по студии, удаление).
- Объявляет выходной порт `MovieRepository` (контракт хранения).
- DTO (`CreateMovieDTO`, `MovieResponseDTO`) защищают внутренние структуры от внешних изменений.

### 3. Адаптеры (`cleanarch/adapter/`):
- `repository.MemoryMovieRepository`: реализует `usecase.MovieRepository`. Легко заменяется на `PostgreSQLMovieRepository` без изменения Use Cases.
- `controller.MovieHTTPController`: транслирует HTTP-запросы в вызовы `usecase.MovieUseCase`, форматирует JSON-ответы.

### 4. Точка сборки (`cmd/app/main.go`):
- Собирает граф зависимостей приложения: `Repo -> UseCase -> Controller -> Server`.

---

## Анализ: Монолит vs Микросервисы

| Критерий | Модульный Монолит | Микросервисы |
|---|---|---|
| **Развертывание** | Один процесс / контейнер | Множество независимых служб, оркестрация (k8s) |
| **Сетевые накладные расходы** | Нулевые (вызовы функций Go в памяти) | Высокие (сеть, TCP, сериализация JSON/Protobuf) |
| **Транзакции и ACID** | Полная поддержка СУБД в рамках единой БД | Отсутствие распределенных ACID; саги (Saga Pattern) |
| **Сложность инфраструктуры** | Низкая | Высокая (Service Discovery, Tracing, API Gateway) |
| **Масштабирование команды** | До определенного размера (merge conflicts) | Высокая независимость команд по фичам |
| **Рекомендация лекции 17** | **Начинать проект практически всегда следует с монолита** с чистой архитектурой. Переходить к микросервисам только при веских организационных причинах. |

---

## Инструкция по запуску и тесты

### Запуск всех тестов:

```bash
cd go_course_4
go test -v -race ./17-system-design/homework-17/...
```

### Запуск сервиса:

```bash
go run ./17-system-design/homework-17/cmd/app/main.go
```

### Проверка эндпоинтов:

```bash
# Получить все фильмы
curl http://localhost:8080/api/v1/movies | jq

# Фильтрация по студии
curl http://localhost:8080/api/v1/movies?studio_id=1 | jq

# Добавление фильма
curl -X POST http://localhost:8080/api/v1/movies \
  -H "Content-Type: application/json" \
  -d '{"title":"Tenet","release_year":2020,"box_office":365000000,"rating":"PG-13","studio_id":1}' | jq
```
