# Том 5: Тестирование, Бенчмарки, Архитектура сервисов и Микросервисы

> [!TIP]
> 📚 **Навигация:** [⬅️ Назад: Том 4](04-go-core-backend-web.md) | [📖 Главное оглавление](README.md) | [Вперед: Том 6 (Алгоритмы и Собеседования) ➡️](06-go-core-algorithms-interview.md)

> [!IMPORTANT]
> 💻 **Практика, примеры кода и домашние задания к этому тому:**
> - 📂 **Код лекций:** [07-testing](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/07-testing) (Unit-тесты, TestMain, моки), [08-prof_debug](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/08-prof_debug) (pprof, bench, trace), [14-RPC](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/14-RPC) (RPC, gRPC), [17-system-design](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/17-system-design) (Clean Architecture), [18-microservices](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/18-microservices) (Docker, URL Shortener), [19-queue](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/19-queue) (Kafka), [20-NoSQL](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/20-NoSQL) (Redis, проект Lynks).
> - 🎯 **Задачник:** [homework-tasks.md (Занятие 7, ДЗ 8, 14, 17–20)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/homework-tasks.md#занятие-7-тестирование-в-go--unit-тесты-табличные-тесты-бенчмарки-и-имитация-зависимостей).
> - 🛠️ **Решения домашних заданий:**
>   - [homework-07](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/07-testing/homework-07) — Табличные unit-тесты и бенчмарки алгоритмов сортировки.
>   - [homework-08](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/08-prof_debug/homework-08) — Профилирование pprof, выявление узких мест и оптимизация CPU/RAM.
>   - [homework-14](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/14-RPC/homework-14) — Реализация службы RPC-сообщений для GoSearch.
>   - [homework-17](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/17-system-design/homework-17) — Чистая слоистая архитектура (Clean Architecture, SOLID).
>   - [homework-18](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/18-microservices/homework-18) — Микросервис URL Shortener в Docker-контейнере.
>   - [homework-19](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/19-queue/homework-19) — Асинхронная очередь событий (Kafka / Event Sourcing).
>   - [homework-20](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/20-NoSQL/homework-20) — Полнофункциональный сервис «Lynks» (Redis + PostgreSQL).
> - 🏛️ **Архитектурный гайд:** [go-system-design-guide.md](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/go-system-design-guide.md).
> - 🗺️ **Сквозной путеводитель:** [course-codebase-guide.md (Этап 3 & 4)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/course-codebase-guide.md#этап-3-middle-go-developer--production-инженерия-базы-данных-и-микросервисы).

---

## 17. Тестирование, Бенчмарки и Профилирование (Testing & pprof)

> 🔗 **Практика и примеры кода:** [07-testing](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/07-testing), [08-prof_debug](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/08-prof_debug) | 🛠️ **Домашние задания:** [homework-07 (TDD)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/07-testing/homework-07), [homework-08 (pprof)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/08-prof_debug/homework-08) | 🎯 **Задачи:** [homework-tasks.md (Занятие 7: Тестирование)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/homework-tasks.md#занятие-7-тестирование-в-go--unit-тесты-табличные-тесты-бенчмарки-и-имитация-зависимостей), [homework-tasks.md (ДЗ 8: Профилирование)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/homework-tasks.md#домашнее-задание-8-профилирование-отладка-и-трассировка-в-go)

### 17.1. Философия и фундаментальные цели тестирования

Тестирование ПО — это процесс испытания программного кода на специально подготовленных данных для определения соответствия реального и ожидаемого поведения программы.

#### Зачем нужны тесты:
1. **Проверка корректности кода:** Убедиться, что логика функции на подготовленных тестовых примерах (фикстурах) или случайных данных возвращает ожидаемый математический и бизнесовый результат.
2. **Развитие, масштабирование и рефакторинг системы:** Чем больше становится проект, тем выше когнитивная нагрузка на разработчиков и тем сложнее отслеживать неявные межмодульные связи. Без тестов любое изменение превращается в рулетку — неизвестно, где и что сломается. Тесты служат «страховочной сеткой» (safety net), защищая от **регрессий** (ситуаций, когда правка в одном месте ломает уже отлаженный функционал в другом).
3. **Живая документация:** Качественный тест лучше любой документации демонстрирует, как использовать функцию, какие аргументы передавать и какие ошибки она способна вернуть.

#### Главное правило: Если тесты сложные — то это плохо!
- **Тесты должны быть простыми и линейными.** Тестовый код вызывает тестируемый метод и прямо сравнивает результат с эталоном.
- Если в самом тесте появляется ветвистая логика, сложные вложенные циклы, математические вычисления или сложный стейт — в коде теста неизбежно появляются **собственные баги**.
- Тесты не должны требовать написания тестов для самих себя.

#### Когда пишутся тесты:
- **До написания кода (TDD — Test-Driven Development):** Сначала пишется падающий тест (RED), затем минимальный код для его прохождения (GREEN), после чего код улучшается (REFACTOR). Отлично подходит, когда алгоритм решения и требования заранее строго определены.
- **Вместе с написанием кода:** Разработчик пишет метод и параллельно создает тестовые сценарии, тестируя контракт API на лету.
- **После написания кода:** Классический подход при проверке гипотез и быстром прототипировании.

#### Парадокс покрытия кода: почему 0% и 100% — это плохо

```mermaid
graph LR
    Zero["0% Покрытие
    (Катастрофа: слепая система, страх любых изменений)"] --> Optimum["70% — 85% Покрытие (ОПТИМУМ)
    (Проверена вся критическая бизнес-логика и пограничные случаи)"]
    Optimum --> Hundred["100% Покрытие
    (Антипаттерн: хрупкость тестов, замедление разработки)"]
```

1. **Почему 0% — это плохо:**
   - Полная неизвестность. Код невозможно безопасно рефакторить, а любая доработка требует ручной сквозной проверки.
2. **Почему 100% — это плохо:**
   - **Хрупкость тестов (Test Brittleness):** Чтобы добиться 100%, приходится тестировать мельчайшие детали внутренней реализации (тривиальные геттеры, однострочные конструкторы `New()`, редкие системные ошибки). При малейшем рефакторинге падают десятки тестов, хотя бизнес-логика работает корректно. Разработка замедляется в разы.
   - **Ложное чувство защищенности:** 100% покрытие операторов (statement coverage) означает лишь то, что каждая строчка кода была исполнена хотя бы один раз. Это **не гарантирует** отсутствие ошибок при неожиданных комбинациях входных данных, гонках данных (data race) или граничных условиях.
   - **Экономическая неэффективность:** Затраты времени на покрытие последних 10–15% тривиального кода превышают пользу от них.
3. **Золотое правило:** Покрытию подлежит **только нетривиальный код** и бизнес-логика. Если некоторый внутренний код вызывается функцией, для которой уже написан подробный тест — писать для него отдельный тест избыточно.

---

### 17.2. Пирамида тестирования и виды тестов в Go

```mermaid
graph TD
    subgraph "Пирамида тестирования (Testing Pyramid)"
        E2E["1. Функциональные / E2E / Приёмочные (End-to-End)
        (Мало, медленные, дорогие, тестируют всю систему целиком)"]
        Integration["2. Интеграционные тесты (Integration Tests)
        (Взаимодействие компонентов: СУБД, кэш, сеть, файловая система)"]
        Unit["3. Модульные / Юнит-тесты (Unit Tests)
        (Подавляющее большинство, быстрые, дешевые, изоляция зависимостей)"]
    end
    Unit --> Integration
    Integration --> E2E
```

#### Закономерности пирамиды:
1. **Чем ниже уровень пирамиды:**
   - Тем проще и быстрее пишутся и исполняются тесты.
   - Тем больше их количественно в кодовой базе (юнит-тесты составляют 70–80%+).
   - Тем выше степень автоматизации (запускаются на каждый коммит и при сборке CI/CD).
2. **Разделение зон ответственности:**
   - **Разработчики:** создают автоматизированные модульные и интеграционные тесты.
   - **Тестировщики (QA) / Автоматизаторы:** создают сквозные функциональные (E2E), приёмочные, нагрузочные и ручные сценарии.

#### «Белый ящик» (White-box) vs «Чёрный ящик» (Black-box) в Go:
- **White-box (`package mypkg`):** Тестовый файл находится в том же пакете. Тесты имеют доступ к приватным (неэкспортируемым) функциям, структурам и полям. Используется для глубокого тестирования внутренних алгоритмов.
- **Black-box (`package mypkg_test`):** Тестовый файл объявляет отдельный пакет с суффиксом `_test`. Тест импортирует тестируемый пакет извне и проверяет **исключительно публичный API-контракт**, как внешний потребитель библиотеки.

> [!WARNING]
> **Миф о 100% покрытии кода (Code Coverage):**
> Покрытие 0% недопустимо, но и стремление к 100% покрытию любой ценой — опасный антипаттерн. 100% покрытие приводит к тестированию тривиальных геттеров, однострочных конструкторов и созданию хрупких тестов, затрудняющих рефакторинг. Распространённый ориентир в индустрии — **70–85% покрытия нетривиальной бизнес-логики** (это эвристика, а не правило: важнее, чтобы были покрыты критичные ветки и граничные случаи). Если вспомогательный код вызывается протестированной функцией, писать отдельный тест для него не требуется.

---

### 17.3. Команды запуска, флаги `go test` и сброс кэша

Тесты в Go размещаются в файлах с суффиксом `*_test.go` рядом с тестируемым кодом. Каждая тестовая функция обязана иметь сигнатуру:
```go
func TestFuncName(t *testing.T) { ... }
```

Минимальный пример (файл `calc_test.go` рядом с `calc.go`, в котором объявлена функция `Add`):

```go
package calc

import "testing"

func TestAdd(t *testing.T) {
    got := Add(2, 3)
    want := 5
    if got != want {
        t.Errorf("Add(2, 3) = %d, want %d", got, want) // Errorf помечает тест упавшим и продолжает выполнение
    }
}
```

Запуск: `go test` (текущий пакет) или `go test ./...` (все пакеты модуля). Важно различать: `t.Errorf` / `t.Error` — пометить тест как упавший и **продолжить**; `t.Fatalf` / `t.Fatal` — пометить упавшим и **немедленно прекратить** этот тест (используйте, когда продолжать бессмысленно, например, после ошибки создания объекта, который используется дальше).

#### Основные флаги CLI команды `go test`:
- `go test ./...` — запустить все тесты в текущем модуле и всех его подпакетах.
- `-v` (Verbose) — детальный вывод статуса (`=== RUN`, `--- PASS`, `--- FAIL`) и логов (`t.Logf`) для каждого теста.
- `-run <regexp>` — запуск только тех тестов, имя которых соответствует регулярному выражению:
  ```bash
  go test -v -run '^TestSort' ./pkg/sorts/...
  ```
- `-cover` — вывод процента покрытия операторов кода тестами.
- `-race` — включить детектор гонок данных (см. Том 3, раздел 15.4); `-short` — режим «быстрых» тестов (см. `testing.Short()`); `-failfast` — остановиться на первом упавшем тесте; `-timeout 30s` — общий таймаут (по умолчанию 10 минут); `-shuffle=on` — случайный порядок тестов (помогает найти зависимости между тестами).
- `-count 1` — **принудительный сброс кэша тестов:**
  > [!IMPORTANT]
  > **Кэширование тестов в Go:** Если исходный код пакета и флаги запуска не менялись с момента прошлого прогона, Go выдает результат из кэша (`ok  mypkg (cached)`), экономя время сборки. Однако при тестировании внешних ресурсов, сетевых запросов, времени или плавающих багов (flaky tests) кэш скрывает реальный результат. Флаг `-count 1` гарантирует реальное выполнение каждого теста. (Кэш работает в режиме списка пакетов — `go test ./...` или `go test pkg`; при запуске просто `go test` без аргументов в каталоге пакета кэш не используется.)

#### Анализ и визуализация покрытия кода (Coverage Report):
```bash
# 1. Сбор профиля покрытия в файл:
go test -coverprofile=coverage.out ./...

# 2. Текстовая сводка по функциям:
go tool cover -func=coverage.out

# 3. Интерактивная HTML-визуализация в браузере (подсветка зеленым покрытых строк, красным — непокрытых):
go tool cover -html=coverage.out
```

---

### 17.4. Управление жизненным циклом тестов пакета: `TestMain(m *testing.M)`

Функция `TestMain` позволяет выполнить глобальную инициализацию (**Setup**) до запуска тестов пакета и гарантированную очистку ресурсов (**Teardown**) после их завершения (например, поднять тестовую БД, применить миграции, закрыть пул соединений):

```go
package repository_test

import (
    "context"
    "log"
    "os"
    "testing"

    "github.com/jackc/pgx/v5/pgxpool"
)

var testPool *pgxpool.Pool

func TestMain(m *testing.M) {
    ctx := context.Background()

    // 1. SETUP: Инициализация тестовой инфраструктуры
    var err error
    testPool, err = pgxpool.New(ctx, "postgres://postgres:pass@localhost:5432/testdb?sslmode=disable")
    if err != nil {
        log.Fatalf("Не удалось подключиться к тестовой БД: %v", err)
    }

    // 2. RUN: Запуск всех тестов в текущем пакете
    exitCode := m.Run()

    // 3. TEARDOWN: Освобождение глобальных ресурсов
    testPool.Close()

    // 4. Завершение с кодом возврата тестов
    os.Exit(exitCode)
}
```

> Начиная с Go 1.15 явный `os.Exit(m.Run())` не обязателен — если `TestMain` просто вернётся, тестовый фреймворк сам завершит процесс с кодом результата `m.Run()`. Но `os.Exit` **не выполняет `defer`**, поэтому, если в `TestMain` используется `defer`, выносите логику в отдельную функцию `run() int` и вызывайте `os.Exit(run())`. Строку подключения к БД обычно берут из переменной окружения, а не зашивают в код.

---

### 17.5. Табличные тесты (Table-Driven Tests) и подтесты `t.Run`

Табличные тесты — общепринятый идиоматичный паттерн в Go. Набор сценариев описывается в виде анонимной структуры, а каждый сценарий запускается в изолированном подтесте с помощью метода `t.Run()`:

```go
func TestReverse(t *testing.T) {
    tests := []struct {
        name string
        in   string
        want string
    }{
        {name: "латиница", in: "ABCdef", want: "fedCBA"},
        {name: "пустая строка", in: "", want: ""},
        {name: "один символ", in: "Z", want: "Z"},
        {name: "палиндром", in: "radar", want: "radar"},
        {name: "кириллица UTF-8", in: "Привет", want: "тевирП"},
    }

    for _, tt := range tests {
        // Подтест отображается отдельной строкой в отчете: TestReverse/латиница
        t.Run(tt.name, func(t *testing.T) {
            got := Reverse(tt.in)
            if got != tt.want {
                t.Errorf("Reverse(%q) = %q, want %q", tt.in, got, tt.want)
            }
        })
    }
}
```

> [!CAUTION]
> **Ловушка замыкания переменной цикла при `t.Parallel()`:**
> До Go 1.22 при вызове `t.Parallel()` внутри цикла `for _, tt := range tests` все горутины подтестов захватывали указатель на одну и ту же переменную `tt`!
> **Правило:** в версиях до Go 1.22 всегда объявляйте локальную копию `tt := tt` внутри тела цикла перед вызовом `t.Run()`. В Go 1.22+ область видимости переменных цикла создается на каждой итерации автоматически.

> 💡 Для сравнения сложных структур (срезы, мапы, вложенные структуры) вместо ручного `if got != want` (который для срезов и мап вообще не скомпилируется) используйте `slices.Equal`, `reflect.DeepEqual` или, лучше, `cmp.Diff` из `github.com/google/go-cmp/cmp` — он печатает читаемое различие двух значений.

---

### 17.6. Имитация зависимостей (Mocks, Stubs) и рефакторинг: Выделение функции

Под **зависимостью** в тестировании понимается любой внешний объект, недоступный или недетерминированный в тестовой среде: СУБД, сеть, стороннее API, файловая система, системные часы.

#### 1. Мокирование через интерфейсы (Interface-based Mocking)
Для изоляции кода зависимость обязана скрываться за интерфейсом. В продуктовом коде передается реальный объект, а в тестах — легковесный Mock/Stub:

```go
// Контракт зависимости
type WeatherClient interface {
    GetTemperature(city string) (float64, error)
}

// Бизнес-логика зависит только от абстракции
type TravelAdvisor struct {
    client WeatherClient
}

func (ta *TravelAdvisor) ShouldTakeJacket(city string) (bool, error) {
    temp, err := ta.client.GetTemperature(city)
    if err != nil {
        return false, err
    }
    return temp < 15.0, nil
}

// В тесте: In-Memory Stub без сетевых вызовов
type mockWeatherClient struct {
    temp float64
    err  error
}

func (m *mockWeatherClient) GetTemperature(city string) (float64, error) {
    return m.temp, m.err
}

func TestShouldTakeJacket(t *testing.T) {
    stub := &mockWeatherClient{temp: 10.5}
    advisor := &TravelAdvisor{client: stub}

    needJacket, err := advisor.ShouldTakeJacket("Moscow")
    if err != nil || !needJacket {
        t.Errorf("expected true, got %v (err: %v)", needJacket, err)
    }
}
```

> 💡 Писать заглушки вручную удобно для небольших интерфейсов. Для больших интерфейсов используют генераторы моков: `go.uber.org/mock` (`mockgen`) или `vektra/mockery`, а также библиотеку утверждений `stretchr/testify` (`assert`, `require`).

#### 2. Рефакторинг «Выделение функции» (Extract Function)
Если зависимость невозможно или нецелесообразно замокать (например, легаси-код или жестко связанная глобальная библиотека), применяется стандартный прием рефакторинга: **выделение чистой функции**:
- Из монолитной функции, совмещающей запрос к сети и расчеты, выделяется чистая детерминированная функция с чистой бизнес-логикой.
- Для чистой функции пишутся 100% надежные юнит-тесты.
- Функция-обертка с побочными эффектами ввода-вывода проверяется редкими интеграционными тестами.

```go
// БЫЛО: Невозможно протестировать расчет без доступа к cbr-xml-daily.ru
func RublesToUSD(rubles float64) (float64, error) {
    resp, err := http.Get("https://www.cbr-xml-daily.ru/daily_json.js")
    // ... парсинг JSON курса ...
    return rubles / rate, nil
}

// СТАЛО: Разделение на I/O и чистые вычисления (Extract Function)
func CalcRublesToUSD(rubles, usdRate float64) (float64, error) {
    if usdRate <= 0 {
        return 0, errors.New("некорректный курс валюты")
    }
    return rubles / usdRate, nil // Чистая функция: тестируется мгновенно и без сети!
}
```

---

### 17.7. Интеграционные тесты и идемпотентность работы с СУБД

Интеграционные тесты проверяют корректность совместной работы компонентов и внешних сервисов (СУБД, шины сообщений, внешние REST API).

#### Принципы надежного тестирования с СУБД:
1. **Идемпотентность тестов:** Повторный или многократный запуск одного и того же теста обязан давать идентичный результат.
2. **Очистка данных:** Тест обязан убирать за собой созданные записи.
   - **Стратегия транзакционного отката (Transaction Rollback):** Тест начинает транзакцию `tx, _ := pool.Begin(ctx)`, выполняет репозиторные операции и в конце делает `defer tx.Rollback(ctx)`. База остается нетронутой!
   - **Стратегия Testcontainers (Золотой стандарт):** Поднятие изолированного эфемерного Docker-контейнера через `testcontainers-go`.

#### 🐳 Практический пример: Интеграционный тест репозитория с `testcontainers-go`:

```go
package repository_test

import (
    "context"
    "testing"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/modules/postgres"
    "github.com/testcontainers/testcontainers-go/wait"
)

// setupTestDB поднимает настоящий PostgreSQL в Docker и возвращает готовый пул соединений:
func setupTestDB(t *testing.T) *pgxpool.Pool {
    t.Helper()
    ctx := context.Background()

    // 1. Запуск изолированного контейнера на свободном случайном порту
    pgContainer, err := postgres.Run(ctx,
        "postgres:16-alpine",
        postgres.WithDatabase("testdb"),
        postgres.WithUsername("postgres"),
        postgres.WithPassword("secret"),
        testcontainers.WithWaitStrategy(
            wait.ForLog("database system is ready to accept connections").
                WithOccurrence(2).
                WithStartupTimeout(15*time.Second)),
    )
    if err != nil {
        t.Fatalf("не удалось запустить контейнер postgres: %v", err)
    }

    // 2. Автоматическое удаление контейнера по завершении теста (t.Cleanup):
    t.Cleanup(func() {
        _ = pgContainer.Terminate(ctx)
    })

    // 3. Получаем динамический Connection String (localhost:случайный_порт)
    connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
    if err != nil {
        t.Fatalf("ошибка получения connection string: %v", err)
    }

    pool, err := pgxpool.New(ctx, connStr)
    if err != nil {
        t.Fatalf("ошибка подключения к пулу: %v", err)
    }
    t.Cleanup(pool.Close)

    // Здесь можно накатить миграции (например, через goose или golang-migrate)
    return pool
}
```

> 💎 **Мостик для Ruby-разработчика:**  
> В Ruby on Rails для изоляции тестов традиционно полагаются на гем `database_cleaner` или механизм `use_transactional_tests = true`, откатывающий транзакции в заранее поднятом локальном PostgreSQL. В мире микросервисов на Go стандартом де-факто стал **Testcontainers**: каждый пакет тестов поднимает абсолютно чистый, изолированный контейнер в Docker, что обеспечивает высокую воспроизводимость на любом CI/CD без предварительной настройки серверов (нужен лишь запущенный Docker).

3. **Разделение тестов по скорости:**
   ```go
   func TestDatabaseHeavyQuery(t *testing.T) {
       if testing.Short() {
           t.Skip("Пропуск тяжелого интеграционного теста в режиме -short")
       }
       // ... работа с реальной БД ...
   }
   ```
   Запуск быстрых тестов: `go test -short ./...`

---

### 17.8. Тестирование HTTP-хендлеров (`net/http/httptest`)

Пакет `net/http/httptest` из стандартной библиотеки Go предоставляет инструменты для двух типов тестирования веб-слоя:
1. **Модульное (изолированное) тестирование** без открытия реальных TCP-сокетов через `httptest.NewRecorder()`.
2. **Интеграционное / сквозное (E2E) тестирование** с подъемом локального HTTP-сервера на случайном порту через `httptest.NewServer()`.

---

#### 1. Модульное тестирование без сетевого сокета (`httptest.NewRecorder`):
`httptest.ResponseRecorder` реализует интерфейс `http.ResponseWriter`, перехватывая статус-код, заголовки и тело ответа в память без сетевых задержек:

```go
package handler_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestCreateProductHandler(t *testing.T) {
    // 1. Подготовка входного JSON тела
    reqBody := `{"name": "MacBook Pro", "price": 1999}`
    req := httptest.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewBufferString(reqBody))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer test-token")

    // 2. Создаем регистратор ответа (заменяет реальный TCP socket)
    rec := httptest.NewRecorder()

    // 3. Вызываем HTTP-обработчик напрямую
    CreateProductHandler(rec, req)

    // 4. Проверка HTTP статус-кода (например, 201 Created)
    if rec.Code != http.StatusCreated {
        t.Fatalf("expected status %d, got %d. Body: %s", http.StatusCreated, rec.Code, rec.Body.String())
    }

    // 5. Проверка заголовков ответа
    if contentType := rec.Header().Get("Content-Type"); contentType != "application/json" {
        t.Errorf("expected Content-Type application/json, got %s", contentType)
    }

    // 6. Валидация структуры ответа
    var resp ProductResponse
    if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
        t.Fatalf("failed to decode response JSON: %v", err)
    }
    if resp.ID == 0 || resp.Name != "MacBook Pro" {
        t.Errorf("unexpected response data: %+v", resp)
    }
}
```

---

> 💡 `rec.Code`, `rec.Header()` и `rec.Body` доступны напрямую; чтобы посмотреть ответ «как его увидит клиент» (с корректными заголовками и `Content-Length`), вызовите `res := rec.Result()` и не забудьте закрыть `res.Body`.

#### 2. Табличные тесты для хендлеров (Table-Driven Handler Tests):
Идиоматичный подход для проверки различных сценариев: успешный ответ, невалидный JSON, отсутствие авторизации, бизнес-ошибки:

```go
func TestUserHandler_TableDriven(t *testing.T) {
    tests := []struct {
        name           string
        method         string
        target         string
        body           string
        authHeader     string
        expectedStatus int
        expectedBody   string
    }{
        {
            name:           "успешный запрос",
            method:         http.MethodPost,
            target:         "/users",
            body:           `{"email": "user@example.com"}`,
            authHeader:     "Bearer valid",
            expectedStatus: http.StatusCreated,
            expectedBody:   `"email":"user@example.com"`,
        },
        {
            name:           "отсутствие авторизации -> 401",
            method:         http.MethodPost,
            target:         "/users",
            body:           `{"email": "user@example.com"}`,
            authHeader:     "",
            expectedStatus: http.StatusUnauthorized,
            expectedBody:   "unauthorized",
        },
        {
            name:           "невалидный JSON -> 400",
            method:         http.MethodPost,
            target:         "/users",
            body:           `{broken json}`,
            authHeader:     "Bearer valid",
            expectedStatus: http.StatusBadRequest,
            expectedBody:   "invalid json payload",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req := httptest.NewRequest(tt.method, tt.target, bytes.NewBufferString(tt.body))
            if tt.authHeader != "" {
                req.Header.Set("Authorization", tt.authHeader)
            }
            rec := httptest.NewRecorder()

            UserHandler(rec, req)

            if rec.Code != tt.expectedStatus {
                t.Errorf("got status %d, want %d", rec.Code, tt.expectedStatus)
            }
            if !strings.Contains(rec.Body.String(), tt.expectedBody) {
                t.Errorf("body %q does not contain %q", rec.Body.String(), tt.expectedBody)
            }
        })
    }
}
```

---

#### 3. Тестирование Path-параметров и контекста:
- **Go 1.22+ Standard Mux (`r.SetPathValue`)**:
  Если хэндлер извлекает параметры пути через стандартный `req.PathValue("id")`:
  ```go
  req := httptest.NewRequest(http.MethodGet, "/users/42", nil)
  req.SetPathValue("id", "42") // Установка значения пути напрямую в тесте
  ```
- **Сторонние роутеры (`chi`, `gorilla/mux`)**:
  Для проверки с роутером запрос прогоняется через инстанс роутера `router.ServeHTTP(rec, req)`.
- **Передача контекста (`req.WithContext`)**:
  Если в хендлере используется пользователь или TraceID из контекста:
  ```go
  ctx := context.WithValue(req.Context(), userCtxKey, &User{ID: 10})
  req = req.WithContext(ctx)
  ```

---

#### 4. Интеграционное E2E тестирование (`httptest.NewServer`):
Поднимает реальный локальный веб-сервер на свободном случайном порту `127.0.0.1`. Идеально подходит для проверки полного цикла: роутер + middleware (CORS, gzip, rate limiter) + клиент `http.Client`:

```go
func TestServerE2E_FullPipeline(t *testing.T) {
    // 1. Инициализируем полноценный роутер со всеми middleware
    router := setupRouter()

    // 2. Поднимаем локальный тестовый HTTP-сервер
    ts := httptest.NewServer(router)
    defer ts.Close() // Автоматически освобождает порт после завершения теста

    // 3. Выполняем реальный сетевой запрос через стандартный клиент
    client := ts.Client() // Клиент, преднастроенный на работу с тестовым сервером
    resp, err := client.Get(ts.URL + "/api/v1/health")
    if err != nil {
        t.Fatalf("HTTP GET failed: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        t.Errorf("got status %d, want 200", resp.StatusCode)
    }
}
```

---

### 17.9. Test-Driven Development (TDD): Red, Green, Refactor

TDD — методология разработки, в которой тесты пишутся **до** написания исполняемого кода.

```mermaid
graph LR
    Red["1. RED
    (Пишем падающий тест на новое требование)"] --> Green["2. GREEN
    (Пишем минимальный код, чтобы тест прошел)"]
    Green --> Refactor["3. REFACTOR
    (Улучшаем архитектуру без изменения поведения)"]
    Refactor --> Red
```

- **Когда TDD эффективен:** Четко формализованные требования, алгоритмические задачи, математические расчеты, парсеры данных.
- **Когда применять с осторожностью:** Исследовательские задачи (R&D), прототипирование меняющегося UI, начальная разработка не устоявшейся архитектуры.

---

### 17.10. Бенчмаркинг производительности (`testing.B`)

Бенчмарк в Go измеряет среднее время выполнения функции и потребление ресурсов (памяти и аллокаций) при многократном запуске.

#### Анатомия бенчмарка:

> [!TIP]
> **Go 1.24+:** вместо `for i := 0; i < b.N; i++` рекомендуется цикл `for b.Loop() { ... }`. Он сам подбирает число итераций, автоматически исключает код подготовки до цикла из замера (не нужен `b.ResetTimer()`) и не даёт компилятору выбросить «бесполезный» вызов. Пример: `for b.Loop() { sort.Ints(buf) }`.

```go
func BenchmarkSortInts(b *testing.B) {
    // 1. Подготовка тестовых данных ДО замера
    data := generateData(1000)
    buf := make([]int, len(data))

    b.ReportAllocs() // Включает подсчет B/op и allocs/op
    b.ResetTimer()   // Сбрасывает таймер, исключая время подготовки данных!

    // 2. Цикл от 0 до b.N:
    // Число b.N рантайм Go подбирает динамически (1, 100, 10000...),
    // пока суммарное время работы бенчмарка не достигнет ~1 секунды (-benchtime)
    for i := 0; i < b.N; i++ {
        copy(buf, data) // Восстанавливаем несортированный срез
        sort.Ints(buf)
    }
}
```

#### Табличные бенчмарки (`b.Run`):
```go
func BenchmarkSort(b *testing.B) {
    sizes := []int{10, 100, 1000, 10000}
    for _, size := range sizes {
        b.Run(fmt.Sprintf("N=%d", size), func(b *testing.B) {
            master := generateData(size)
            buf := make([]int, size)
            b.ResetTimer()

            for i := 0; i < b.N; i++ {
                copy(buf, master)
                sort.Ints(buf)
            }
        })
    }
}
```

#### Запуск и расшифровка метрик бенчмарков:
```bash
go test -bench=. -benchmem -run '^$' ./...
```
*(Флаг `-run '^$'` гарантирует, что ни один юнит-тест не запустится, исключая лишний шум).*

**Пример вывода консоли** (число после дефиса в имени — значение `GOMAXPROCS`):
```text
BenchmarkSortInts/N=1000-8         12744      9645 ns/op         0 B/op        0 allocs/op
BenchmarkSortFloat64s/N=1000-8      7758     16249 ns/op         0 B/op        0 allocs/op
```

- **`12744`:** Количество итераций `b.N`, выполненных за время замера.
- **`9645 ns/op`:** Наносекунд на одну операцию (чем меньше, тем быстрее).
- **`0 B/op`:** Байт памяти, выделенных в куче (heap) за одну операцию.
- **`0 allocs/op`:** Количество вызовов аллокатора памяти в куче за операцию.

> [!TIP]
> **Почему в этом примере `sort.Ints` быстрее `sort.Float64s`:**
> 1. Сравнение целых чисел — простая инструкция, а сравнение `float64` должно корректно обрабатывать значения `NaN` (Not-a-Number), для которых обычные законы порядка не работают (`NaN != NaN`), — это дополнительные проверки.
> 2. Обе функции имеют `0 allocs/op`, так как сортировка выполняется на месте (in-place).
> 3. Конкретные числа зависят от версии Go, процессора и данных (в Go 1.22+ обе функции реализованы через обобщённую `slices.Sort`) — **всегда измеряйте на своём окружении**, а не переносите выводы из чужих бенчмарков.

---

### 17.11. Fuzz-тестирование (Go 1.18+)

Фаззинг автоматически генерирует псевдослучайные мутации входных данных для выявления паник, переполнений буфера и скрытых логических ошибок:

```go
func FuzzParseJSON(f *testing.F) {
    // Начальный корпус данных (seed corpus):
    f.Add([]byte(`{"name": "Alex", "age": 25}`))
    f.Add([]byte(`{}`))
    f.Add([]byte(`{"invalid": true}`))

    // Фазз-таргет:
    f.Fuzz(func(t *testing.T, data []byte) {
        var user User
        // Функция НЕ должна паниковать ни на каком входе:
        _ = json.Unmarshal(data, &user)
    })
}
```
Запуск: `go test -fuzz=FuzzParseJSON -fuzztime=30s .` (флаг `-fuzz` работает только для **одного** пакета и одного фазз-теста, поэтому `./...` использовать нельзя).

> [!NOTE]
> - Найденные «сломавшие» вход данные сохраняются в `testdata/fuzz/<ИмяТеста>/` и затем автоматически проверяются при каждом обычном `go test` — так фаззинг превращается в регрессионный тест.
> - Пример выше проверяет лишь отсутствие паники. Более сильные фазз-тесты проверяют **свойства** (инвариант): например, «`Decode(Encode(x)) == x`» или «`Reverse(Reverse(s)) == s`».

---

### 17.12. Профилирование бенчмарков (`-cpuprofile`, `-memprofile`)

Профиль — это набор статистических данных многократных замеров потребления ресурсов (CPU, куча, блокировки) во время выполнения программы. Профилирование позволяет обнаружить узкие места производительности (Hotspots) на уровне конкретных функций и строк кода.

#### 1. Снятие профиля при запуске тестов и бенчмарков:
Профилирование во время обычных тестов не дает объективной картины из-за малого времени работы. Снимать профиль имеет смысл **только при выполнении бенчмарков**:

```bash
# Профиль процессора (CPU):
go test -run=^$ -bench=. -benchmem -cpuprofile=cpu.out

# Профиль распределения памяти в куче (Memory/Heap):
go test -run=^$ -bench=. -benchmem -memprofile=mem.out

# Профиль блокировок горутин (Block profile) и борьбы за мьютексы (Mutex profile):
go test -run=^$ -bench=. -blockprofile=block.out -mutexprofile=mutex.out
```

#### 2. Анализ профиля через интерактивный инструмент `go tool pprof`:
```bash
# Интерактивный CLI-режим:
go tool pprof cpu.out

# Вывод топ-10 самых ресурсоемких функций в терминал:
go tool pprof -text cpu.out

# Запуск веб-интерфейса в браузере с графом вызовов (Call Graph):
go tool pprof -http=:8080 cpu.out

# Экспорт графа вызовов в файл SVG или PDF:
go tool pprof -svg cpu.out > cpu.svg
go tool pprof -pdf cpu.out > cpu.pdf
```

#### 3. Интерпретация метрик и графа вызовов:
- **`flat`**: Время или объем ресурсов, потраченные непосредственно внутри данной функции.
- **`cum` (cumulative)**: Суммарное время функции, включая все дочерние вызовы.
- **Размер прямоугольника**: Чем больше размер узла, тем больше ресурсов (`flat`) потребляет именно эта функция.
- **Цвет узла**: От белого к ярко-красному. Красный прямоугольник — критическая «горячая точка» (hotspot).
- **Толщина стрелки**: Показывает относительный объем работы или памяти, передаваемый по данному пути выполнения.
- В веб-интерфейсе (`-http=:8080`) в меню **View** доступен удобный режим **Flame Graph** («пламенный график»): ширина блока пропорциональна времени, а вложенность отражает стек вызовов.

---

### 17.13. Профилирование живых приложений в runtime (`net/http/pprof`)

Для долгоживущих сетевых сервисов и микросервисов профиль снимается непосредственно в процессе работы под боевой или тестовой нагрузкой.

#### 1. Подключение профилировщика:
Используется безымянный импорт (side-effect import) пакета `net/http/pprof`, который автоматически регистрирует стандартные HTTP-хэндлеры в `http.DefaultServeMux`:

```go
package main

import (
    "log"
    "net/http"
    _ "net/http/pprof" // Регистрирует маршруты /debug/pprof/*
)

func main() {
    // В production рекомендуется выносить pprof на отдельный внутренний порт,
    // закрытый от публичного доступа через фаервол / ingress:
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()

    // Основной сервер приложения:
    log.Println(http.ListenAndServe(":8080", appHandler))
}
```

> [!CAUTION]
> Пакет `net/http/pprof` регистрирует обработчики в **`http.DefaultServeMux`**. Если основной публичный сервер вашего приложения тоже использует `DefaultServeMux` (например, `http.ListenAndServe(":8080", nil)`), профилировщик станет доступен из интернета — это утечка внутренних данных и вектор DoS. Поэтому в примере выше основной сервер использует собственный обработчик (`appHandler`), а pprof слушает `localhost:6060`.

#### 2. Список стандартных эндпоинтов `/debug/pprof/`:
- `http://localhost:6060/debug/pprof/` — веб-страница со списком всех профилей.
- `/debug/pprof/profile?seconds=30` — сбор 30-секундного семплированного профиля CPU.
- `/debug/pprof/heap` — снимок распределения памяти в куче.
- `/debug/pprof/goroutine` — снимок стеков всех активных горутин (поиск утечек горутин / Goroutine Leaks).
- `/debug/pprof/block` — задержки на синхронизации (каналы, `sync.Mutex`, `sync.WaitGroup`).
- `/debug/pprof/mutex` — трассировка борьбы за мьютексы (Lock Contention).

> Профили `block` и `mutex` по умолчанию **выключены** и будут пустыми, пока вы не включите их в коде: `runtime.SetBlockProfileRate(1)` и `runtime.SetMutexProfileFraction(1)` (на боевых сервисах используйте значения побольше, чтобы снизить накладные расходы).
- `/debug/pprof/threadcreate` — активность создания потоков операционной системы (M).

#### 3. Снятие и анализ профилей приложения:
```bash
# Сбор CPU профиля (30 секунд) с открытием интерактивного веб-графа:
go tool pprof -http=:8081 'http://localhost:6060/debug/pprof/profile?seconds=30'   # URL в кавычках, иначе zsh воспримет ? как шаблон

# Анализ памяти в куче (inuse_space — занято прямо сейчас):
go tool pprof -inuse_space -http=:8081 http://localhost:6060/debug/pprof/heap

# Анализ общего объема аллокаций (alloc_space — сколько памяти было выделено за все время):
go tool pprof -alloc_space -http=:8081 http://localhost:6060/debug/pprof/heap

# Сравнение профилей до и после оптимизации (дифф):
go tool pprof -base profile_before.pb.gz profile_after.pb.gz
```

> [!TIP]
> **`inuse_space` vs `alloc_space`**:
> - Используйте `inuse_space` для поиска **утечек памяти (Memory Leaks)** — объектов, которые не собираются GC и продолжают расти.
> - Используйте `alloc_space` для поиска **паразитных временных аллокаций**, создающих высокое давление на сборщик мусора (GC Pressure), даже если память вовремя освобождается.

---

### 17.14. Отладка программного обеспечения и отладчик Delve (`dlv`)

#### 1. Инженерия и культура отладки:
- **Отладка (Debugging)** — процесс локализации и устранения дефектов в уже написанном исходном коде.
- **Отладка vs Тесты**: Модульные тесты изолируют отдельные функции и пакеты. Отладка применяется ко всему работающему комплексу приложения, когда на искусственных тестах программа работает штатно, но сбоит на реальных данных или в краевых сценариях.
- **Отладчик — крайнее средство**: Профессиональный инженер сначала анализирует логи (`slog`), метрики и трассировки. Если логи не позволяют установить причину аномалии, привлекается интерактивный отладчик.
- **Отладка не предназначена для поиска проблем с производительностью**: Производительность исследуется исключительно с помощью профилировщиков (`pprof`) и трассировщика (`trace`).

#### 2. Специализированный отладчик Go — Delve (`dlv`):
Delve спроектирован с глубоким пониманием рантайма Go: он умеет работать с горутинами, понимать структуру интерфейсов `iface`, слайсов и мап, а также корректно отслеживать стек вызовов при динамическом расширении стека горутин.

```bash
# Установка Delve:
go install github.com/go-delve/delve/cmd/dlv@latest

# Запуск отладки пакета main (Delve сам компилирует программу с отключёнными оптимизациями):
dlv debug ./cmd/api

# Если вы отлаживаете заранее собранный бинарник, собирайте его так: go build -gcflags="all=-N -l"
# (-N отключает оптимизации, -l — инлайнинг), иначе переменные могут быть «оптимизированы»

# Запуск отладки конкретного тестового пакета:
dlv test ./pkg/service

# Подключение к уже запущенному PID процесса:
dlv attach 12345
```

#### 3. Основные команды CLI Delve:
| Команда CLI | Сокращение | Назначение |
| :--- | :---: | :--- |
| `break <target>` | `b` | Установить точку останова (по функции: `b main.main`, по строке: `b worker.go:42`). |
| `condition <id> <expr>` | `cond` | Установить условие для брейкпоинта (например: `cond 1 i == 10 && user.ID != 0`). |
| `continue` | `c` | Продолжить выполнение программы до следующего брейкпоинта или завершения. |
| `next` | `n` | Выполнить следующую строчку кода без захода внутрь функций (Step Over). |
| `step` | `s` | Шагнуть внутрь вызываемой функции (Step Into). |
| `stepout` | `so` | Довыполнить текущую функцию и остановиться на строке возврата (Step Out). |
| `print <expr>` | `p` | Распечатать значение переменной или выражения (например: `p nums[i]`, `p len(users)`). |
| `locals` | — | Вывести все локальные переменные текущей функции. |
| `args` | — | Вывести значения входных аргументов текущей функции. |
| `vars <regex>` | — | Поиск и вывод глобальных переменных пакета. |
| `set <var> = <val>` | — | Изменить значение переменной в памяти во время паузы приложения. |
| `stack` | — | Показать трассировку вызовов (Call Stack) текущей горутины. |
| `frame <n>` | — | Переключиться на `n`-й кадр стека вызовов для просмотра его локальных переменных. |
| `goroutines` | — | Показать список всех горутин в рантайме с их статусом. |
| `goroutine <id>` | — | Переключить контекст отладчика на конкретную горутину. |

#### 4. Графическая отладка в VS Code (`.vscode/launch.json`):
Для комфортной работы в VS Code настраивается конфигурация отладки:

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Service (Delve)",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/cmd/api"
        },
        {
            "name": "Debug Test Under Cursor",
            "type": "go",
            "request": "launch",
            "mode": "test",
            "program": "${fileDirname}",
            "args": ["-test.run", "^TestTargetFunction$"]
        }
    ]
}
```

---

### 17.15. Трассировка выполнения (`runtime/trace` & `go tool trace`)

Трассировка выполнения (Execution Tracing) фиксирует хронологический таймлайн событий рантайма Go с наносекундной точностью. В отличие от статистического семплирования `pprof`, трассировщик протоколирует каждое изменение состояния горутин и планировщика.

#### 1. Какие события фиксирует трассировщик:
- Выполнение вычислений на логических процессорах `P` и системных потоках `M`.
- Переключение контекста между горутинами (`G` status: `Runnable`, `Running`, `Waiting`).
- Полный жизненный цикл сборщика мусора (GC): запуск маркинга, фазы STW (Sweep Termination, Mark Termination), помощь аллоцирующих горутин (`GC Assist`).
- Блокировки и простои: сетевой ввод/вывод (`netpoll`), системные вызовы (`syscall`), захват каналов и мьютексов.
- Давление на сборщик мусора (GC Pressure) при частых аллокациях памяти.

#### 2. Включение трассировки в приложении:
```go
package main

import (
    "log"
    "os"
    "runtime/trace"
)

func main() {
    // Создаем файл журнала трассировки
    f, err := os.Create("trace.out")
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()

    // Включаем трассировщик
    if err := trace.Start(f); err != nil {
        log.Fatal(err)
    }
    defer trace.Stop() // Гарантирует сброс буферов на диск

    // ... выполнение бизнес-логики приложения ...
}
```

Для разметки собственных участков кода в трассе используют «задачи» и «регионы»: `ctx, task := trace.NewTask(ctx, "handleOrder"); defer task.End()` и `defer trace.StartRegion(ctx, "db-query").End()` — они отображаются в `go tool trace` отдельными полосами.

Также трассировку можно записать непосредственно во время прогона тестов:
```bash
go test -trace=trace.out ./pkg/service
```

#### 3. Визуализация и анализ через `go tool trace`:
```bash
go tool trace trace.out
```
Команда запускает локальный веб-сервер и открывает интерфейс Chrome Trace Viewer:
- **View trace**: интерактивная временная шкала (Timeline). Навигация клавишами `W` (приблизить), `S` (отдалить), `A` (влево), `D` (вправо).
- **Goroutine analysis**: таблица всех горутин с группировкой по типам и статистикой времени выполнения, ожидания на каналах и блокировок планировщика.
- **Network / Sync / Syscall blocking profile**: профили задержек на системных вызовах и синхронизации.
- **Scheduler latency profile**: гистограмма задержек перехода горутины из очереди готовности к фактическому исполнению на ядре CPU.

---

### 17.16. Сводная таблица методов пакета `testing`

| Метод `testing.T` | Прерывает тест? | Описание |
| :--- | :---: | :--- |
| `t.Log(args...)`, `t.Logf(format, args...)` | ❌ Нет | Выводит информационное сообщение (отображается при `-v` или при падении теста). |
| `t.Fail()` | ❌ Нет | Помечает тест как проваленный, но **продолжает** выполнение текущей функции. |
| `t.FailNow()` | 🚨 **Да** | Помечает тест как проваленный и **немедленно прерывает** выполнение текущей горутины (`runtime.Goexit`). |
| `t.Error(...)`, `t.Errorf(...)` | ❌ Нет | Эквивалентно вызову `t.Log()` + `t.Fail()`. Тест идет дальше. |
| `t.Fatal(...)`, `t.Fatalf(...)` | 🚨 **Да** | Эквивалентно вызову `t.Log()` + `t.FailNow()`. Тест немедленно завершается. |
| `t.Skip(...)`, `t.Skipf(...)`, `t.SkipNow()` | ⏩ **Да** | Пропускает выполнение теста (помечается как `SKIP`). |
| `t.Helper()` | ❌ Нет | Помечает функцию как вспомогательную; в трассировке ошибки будет показана строка вызывающего кода, а не файла хелпера. |
| `t.Parallel()` | ❌ Нет | Разрешает параллельный запуск данного теста с другими параллельными тестами в рамках `GOMAXPROCS`. |
| `t.Cleanup(fn)` | ❌ Нет | Регистрирует функцию отложенной очистки, вызываемую после завершения теста и всех подтестов. |
| `t.TempDir()` | ❌ Нет | Создает уникальную временную директорию на диске, которая автоматически удаляется после завершения теста. |
| `t.Setenv(key, val)` | ❌ Нет | Задаёт переменную окружения на время теста и восстанавливает старое значение после (нельзя использовать вместе с `t.Parallel()`). |
| `t.Context()` | ❌ Нет | (Go 1.24+) Контекст, который автоматически отменяется прямо перед выполнением функций `Cleanup` — удобно для остановки фоновых горутин теста. |

---

### 17.17. Детерминированное тестирование конкурентности: `testing/synctest` (Go 1.25+)

Исторически тестирование кода с таймаутами (`time.Sleep`, `time.NewTicker`, `context.WithTimeout`) было источником нестабильных тестов (Flaky Tests): тесты либо искусственно замедляли CI/CD пайплайны на секунды, либо падали из-за случайных задержек планировщика на перегруженном раннере.

Пакет **`testing/synctest`** (экспериментальный в Go 1.24 под `GOEXPERIMENT=synctest`, стабильный начиная с **Go 1.25**) запускает тест в изолированном «пузыре» (_Bubble_) с синтетическим виртуальным временем:

```go
package mytest

import (
    "testing"
    "testing/synctest"
    "time"
)

func TestTimeoutWorker_Deterministic(t *testing.T) {
    // synctest.Test запускает функцию в пузыре виртуального времени
    // (в Go 1.24 экспериментальный вариант назывался synctest.Run — в актуальных версиях его нет):
    synctest.Test(t, func(t *testing.T) {
        done := make(chan bool)

        go func() {
            // Внутри пузыря задержка 10 минут симулируется виртуальным временем:
            time.Sleep(10 * time.Minute)
            done <- true
        }()

        // Когда все горутины пузыря заблокированы, виртуальные часы мгновенно
        // «перематываются» к ближайшему таймеру:
        <-done
    }) // Тест завершается за миллисекунды вместо 10 минут реального ожидания
}
```

- **Преимущества:** тесты с таймерами выполняются мгновенно и детерминированно — без случайных сбоев из-за нагрузки на CI.
- Функция `synctest.Wait()` блокируется до тех пор, пока все горутины пузыря не окажутся в состоянии «надёжной блокировки» (ожидание канала, `time.Sleep`, `sync.WaitGroup` внутри пузыря) — это позволяет проверить состояние программы ровно в момент, когда она «затихла».
- Ограничения: все горутины пузыря должны завершиться до окончания теста; блокировки на реальном I/O (сеть, файлы) не считаются «надёжными» и мешают перемотке времени.

---

## 18. Архитектура сервисов и Структурированное логирование

> 🔗 **Практика и примеры кода:** [17-system-design](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/17-system-design) | 🛠️ **Домашнее задание:** [homework-17 (Clean Arch)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/17-system-design/homework-17) | 🎯 **Задачи:** [homework-tasks.md (ДЗ 17: Архитектура)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/homework-tasks.md#домашнее-задание-17-архитектура-go-приложения-clean-architecture-и-solid) | 🏛️ **Гайд:** [go-system-design-guide.md](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/go-system-design-guide.md)

### 18.1. Clean Architecture и Standard Go Project Layout

В реальном промышленном Go-бэкенде принята трехслойная архитектура с инверсией зависимостей (Dependency Inversion):

```mermaid
graph TD
    subgraph "Слои приложения (Clean Architecture)"
        Handler["1. Handler / Transport Layer (HTTP / gRPC)"] -->|Вызывает интерфейс| Service["2. Service / Usecase Layer (Бизнес-логика)"]
        Service -->|Вызывает интерфейс| Repo["3. Repository / Infrastructure Layer (PostgreSQL / Redis)"]
    end
```

#### 💡 На пальцах: Чистая архитектура как средневековая крепость

Чтобы раз и навсегда понять идею Clean Architecture (Порты и Адаптеры), представьте неприступную крепость:
- 👑 **Ядро крепости (Domain / Entities):**  
  Самый защищённый центр. Здесь живут чистые доменные структуры (`User`, `Account`, `Order`) и фундаментальные правила бизнеса. Ядро **абсолютно ничего не знает о внешнем мире**: в нём нет ни строчки кода про SQL, HTTP, JSON, PostgreSQL или Docker.
- 🏰 **Внутренний двор (Service / Use Cases):**  
  Бизнес-сценарии («зарегистрировать пользователя», «перевести деньги с баланса»). Сервис общается с внешним миром **только через интерфейсы**: он сам объявляет контракт `type UserRepository interface { Get(id) (*User, error) }`. Кто именно будет сохранять данные — реальный Postgres, кэш Redis или тестовый mock-объект в памяти — сервису не важно!
- 🚪 **Внешние ворота и стены (Adapters / Transport & Infrastructure):**  
  - *Ворота на вход (Primary / Driving Adapters):* HTTP-хендлеры, gRPC-серверы, CLI-команды. Они принимают сырой запрос (JSON, protobuf), валидируют его и вызывают чистый метод сервиса.
  - *Ворота на выход (Secondary / Driven Adapters):* драйвер БД `pgxpool`, кэш Redis, брокер сообщений Kafka, клиент отправки SMS.
- 🎯 **Главная победа чистой архитектуры:**  
  Зависимости направлены **всегда снаружи внутрь к ядру**. Если завтра вы решите сменить HTTP на gRPC, а PostgreSQL на ClickHouse — **код чистой бизнес-логики в `service` и сущности в `domain` не изменятся ни на один символ!**

> [!NOTE]
> **Осторожно с термином «стандарт»:** репозиторий `golang-standards/project-layout` **не является официальным стандартом** Go (в нём нет одобрения команды Go), а популярная в нём папка `pkg/` вызывает споры. Официальные рекомендации по организации модулей: <https://go.dev/doc/modules/layout>. Для небольших проектов достаточно плоской структуры; каталоги `cmd/` и `internal/` добавляют, когда в проекте появляется несколько бинарников или потребность скрыть код от внешних импортов. Ниже приведена типичная для сервисов схема — как удобный ориентир, а не обязательное правило.

#### Структура директорий (по мотивам `golang-standards/project-layout`):

- `cmd/api/main.go` — точка входа: инициализация конфигов, пулов соединений и DI (внедрение зависимостей).
- `internal/` — приватный код проекта, который компилятор Go запрещает импортировать из других модулей:
  - `internal/domain/` — доменные структуры (`User`, `Movie`) и кастомные ошибки (`ErrNotFound`).
  - `internal/handler/` — HTTP/gRPC контроллеры (сериализация JSON, статус-коды).
  - `internal/service/` — бизнес-логика (валидация, расчет баланса, отправка писем).
  - `internal/repository/` — работа с БД (`pgxpool`, транзакции).
- `pkg/` — публичные переиспользуемые библиотеки (клиенты API, утилиты).
- `migrations/` — файлы SQL-миграций для `goose` или `golang-migrate`.

---

### 18.2. Структурированное логирование: `log/slog` (в Go 1.21+)

Начиная с **Go 1.21**, в стандартную библиотеку добавлен пакет структурированного логирования `log/slog`, которого во многих проектах достаточно вместо внешних библиотек (`zap`, `logrus`; для экстремально нагруженных сервисов `zap` и `zerolog` по-прежнему могут быть быстрее):

```go
import (
    "log/slog"
    "os"
    "time"
)

func main() {
    // 1. JSON-форматтер для сбора логов в ELK / Grafana Loki:
    logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelInfo, // Debug, Info, Warn, Error
    }))
    slog.SetDefault(logger)

    // 2. Логирование со структурированными полями:
    slog.Info("пользователь успешно авторизован",
        slog.Int64("user_id", 42),
        slog.String("ip", "192.168.1.1"),
        slog.Duration("latency", 15*time.Millisecond),
    )

    // Краткая форма: чередующиеся пары «ключ, значение» (следите за парностью — `go vet` проверяет вызовы slog):
    slog.Info("заказ создан", "order_id", 1001, "total", 99.9)

    // Логгер с общими полями для всех записей (например, идентификатор запроса):
    reqLogger := slog.With("request_id", "abc-123")
    reqLogger.Warn("медленный запрос")

    // Версии с контекстом (InfoContext/ErrorContext) позволяют обработчику достать trace_id из ctx
}
```

> Для локальной разработки вместо JSON удобнее читаемый `slog.NewTextHandler(os.Stdout, nil)`. Не логируйте пароли, токены и персональные данные; для скрытия значений тип может реализовать интерфейс `slog.LogValuer`.

---

### 18.3. Принципы декомпозиции пакетов: Screaming Architecture, Сцепление и Связность

Пакет в Go — это не просто директория с кодом, а **неделимый строительный блок повторного использования**.

#### 1. Screaming Architecture (Кричащая архитектура)

Структура пакетов должна с первого взгляда «кричать» о предметной области и назначении приложения, а не о технических библиотеках или паттернах:

- ❌ **Антипаттерн (техническая группировка):** пакеты `models/`, `controllers/`, `utils/`, `common/`, `helpers/`. В таких пакетах скапливается несвязанный разнородный код, провоцирующий циклические импорты.
- ✅ **Идиоматичный подход (доменная декомпозиция):** пакеты `crawler/`, `index/`, `billing/`, `auth/`, `order/`.

#### 2. Слабое сцепление (Low Coupling) vs Сильная связность (High Cohesion)

- **Сильная связность (High Cohesion):** Все типы, функции и константы внутри одного пакета тесно связаны и решают **одну законченную задачу**. Если функционал не относится к этой задаче, его выносят в отдельный пакет.
- **Слабое сцепление (Low Coupling):** Пакеты минимально зависят от внутренних деталей друг друга. Изменение внутренней реализации одного пакета не требует переписывания других пакетов. Взаимодействие на границах строится через **минималистичные интерфейсы**.

---

### 18.4. Архитектурный шаблон Ядро и плагины (Core & Plugins / Microkernel)

Шаблон декомпозиции, при котором центральное приложение (ядро) изолировано от деталей внешних подсистем:

```mermaid
graph TD
    subgraph "Ядро системы (Core / Server)"
        Server["Объект Server / App Core"]
    end

    subgraph "Плагины через Интерфейсы"
        DB_I["Интерфейс DB / Storage"]
        Log_I["Интерфейс Logger"]
        Crawler_I["Интерфейс Crawler"]
        API_I["Интерфейс Handler / Router"]
    end

    Server --> DB_I
    Server --> Log_I
    Server --> Crawler_I
    Server --> API_I

    PG["Плагин: PostgreSQL (pgxpool)"] -.->|Реализует| DB_I
    Redis["Плагин: Redis Cache"] -.->|Реализует| DB_I
    Spider["Плагин: Web Spider"] -.->|Реализует| Crawler_I
    Membot["Плагин: Mock Membot"] -.->|Реализует| Crawler_I
    Slog["Плагин: log/slog"] -.->|Реализует| Log_I
```

- **Преимущества:**
  - Ядро системы (`Server`) не знает ничего о конкретных библиотеках (`pgx`, `net/http`, сторонних API).
  - Любая подсистема (например, замена поискового робота `spider` на `membot` или замена PostgreSQL на In-Memory) происходит прозрачно без изменения логики ядра.
  - Позволяет получить ортогональные подсистемы: проект легко тестировать (подмена моками за 1 строку) и дешево поддерживать.

---

### 18.5. Циклические зависимости (`import cycle not allowed`) и стратегии их устранения

Компилятор Go **принципиально запрещает** циклические зависимости между пакетами ($A \to B \to C \to A$). При наличии цикла компиляция завершается фатальной ошибкой: `import cycle not allowed`.

> [!NOTE]
> **Почему Go запрещает циклические импорты?**
>
> 1. Это гарантирует построение графа зависимостей в виде **ориентированного ациклического графа (DAG)**.
> 2. Позволяет собирать пакеты в правильном порядке и параллельно (сначала листья графа, затем зависящие от них) и ускоряет компиляцию.
> 3. Защищает архитектуру от превращения в «большой комок грязи» (Big Ball of Mud).

#### Пример проблемы:

```go
package a
import "project/pkg/b"
type ServiceA struct { B *b.ServiceB }

package b
import "project/pkg/a" // ❌ Ошибка: import cycle not allowed!
type ServiceB struct { A *a.ServiceA }
```

#### 3 паттерна разрешения циклических импортов:

```mermaid
flowchart TD
    subgraph "1. Вынесение общего контракта (Extract Domain / Interface Package)"
        D["pkg/domain или pkg/types"]
        PA1["pkg/a"] --> D
        PB1["pkg/b"] --> D
    end

    subgraph "2. Инверсия зависимостей (Dependency Inversion)"
        PA2["pkg/a (объявляет интерфейс BConsumer)"]
        PB2["pkg/b (реализует BConsumer)"] --> PA2
    end
```

1. **Вынесение общего интерфейса или структур в независимый пакет (Extract Interface / Types):**  
   Создается общий пакет `domain` или `types`, не зависящий ни от `a`, ни от `b`. И `a`, и `b` импортируют только его.
2. **Инверсия зависимостей (Dependency Inversion):**  
   Пакет `a` объявляет **узкий интерфейс** того, что ему нужно от `b`. Пакет `b` реализует этот интерфейс и внедряется в `a` на уровне `main.go`. Пакет `a` перестает импортировать `b`.
3. **Объединение пакетов (Merge Packages):**  
   Если два пакета настолько тесно связаны, что постоянно требуют доступа к типам друг друга — они логически представляют собой **один пакет** и должны быть объединены в одной директории.

---

## 19. Управление зависимостями: Go Modules и Vendoring

### 19.1. Ключевые команды `go mod`

- `go mod init <module_name>` — инициализация нового Go-модуля.
- `go mod tidy` — автоматическое удаление неиспользуемых и скачивание недостающих зависимостей.
- `go mod vendor` — создание локальной копии всех сторонних библиотек в папке `vendor/` (для offline-сборки в закрытом контуре/Enterprise).
- `go mod verify` — проверка контрольных сумм скачанных пакетов на соответствие `go.sum`.
- `go get github.com/user/pkg@v1.2.3` — добавить/обновить зависимость до конкретной версии (`@latest`, `@none` для удаления); `go get -u ./...` — обновить зависимости.
- `go list -m all` — список всех модулей сборки; `go mod why <pkg>` — почему модуль нужен проекту; `go mod graph` — граф зависимостей.
- `go work init` / `go.work` — режим рабочего пространства: разработка нескольких связанных модулей локально без `replace`.

### 19.2. Назначение файлов `go.mod` и `go.sum`

1. **`go.mod`:** Декларирует имя модуля, минимальную версию Go (включающую языковую семантику, например `go 1.22`) и список прямых/транзитивных (`// indirect`) зависимостей с их версиями.
   - Директива `replace`: позволяет временно подменить внешний модуль на локальную папку при разработке (`replace github.com/user/pkg => ../pkg`).
   - Директива `go 1.22` — минимально необходимая версия Go (с Go 1.21 более старый компилятор откажется собирать модуль), а `toolchain go1.22.3` задаёт конкретный рекомендуемый тулчейн: команда `go` при необходимости сама скачает нужную версию.
   - Версии выбираются алгоритмом **Minimal Version Selection (MVS)**: берётся минимальная версия, удовлетворяющая требованиям всех зависимостей, что делает сборки воспроизводимыми без lock-файла (роль «lock-файла» частично играет `go.sum`).
2. **`go.sum`:** Криптографические хэши (SHA-256) содержимого модулей для защиты от подмены кода в репозиториях (Supply Chain Attacks).

---

### 19.3. Встраивание файлов в бинарник: `//go:embed` (Go 1.16+)

Директива `//go:embed` позволяет включать статические файлы (HTML-шаблоны, SQL-миграции, конфиги, сертификаты) **прямо внутрь скомпилированного бинарника**. При деплое не нужно копировать файлы рядом с исполняемым файлом:

```go
import "embed"

// 1. Встраивание одного файла как строки:
//go:embed config/defaults.yaml
var defaultConfig string

// 2. Встраивание одного файла как []byte:
//go:embed static/logo.png
var logoPNG []byte

// 3. Встраивание целой директории как виртуальной файловой системы (embed.FS):
//go:embed templates
var templateFS embed.FS // пути внутри FS начинаются с "templates/..."; чтобы «зайти» внутрь каталога, используйте fs.Sub(templateFS, "templates")

// 4. Использование embed.FS в HTTP-сервере для отдачи статики:
func main() {
    mux := http.NewServeMux()
    mux.Handle("/static/", http.FileServer(http.FS(templateFS)))
    http.ListenAndServe(":8080", mux)
}

// 5. Встраивание SQL-миграций для goose:
//go:embed migrations/*.sql
var migrationsFS embed.FS
```

> [!WARNING]
> **Правила директивы `//go:embed`:**
>
> - Комментарий `//go:embed` должен стоять **непосредственно перед** переменной (без пустых строк).
> - Нельзя встраивать файлы из родительских директорий (`../secret.key` — ошибка компиляции).
> - Переменная должна быть типа `string`, `[]byte` или `embed.FS` и объявлена на **уровне пакета** (внутри функции нельзя); пакет `embed` нужно импортировать (для `string`/`[]byte` достаточно `import _ "embed"`).
> - Файлы, имена которых начинаются с `.` или `_`, при встраивании каталога по умолчанию пропускаются; чтобы включить их, используйте префикс `all:` (`//go:embed all:static`).

---

### 19.4. Кодогенерация: `go generate`

Директива `//go:generate` запускает произвольные команды во время разработки (но **НЕ** во время `go build`). Используется для генерации кода из шаблонов:

```go
// В файле с исходным кодом:
//go:generate stringer -type=Weekday
//go:generate mockgen -source=repository.go -destination=mock_repository.go
//go:generate protoc --go_out=. --go-grpc_out=. api/proto/movie.proto
//go:generate sqlc generate
```

```bash
# Запуск всех //go:generate директив в проекте:
go generate ./...
```

> 💡 Для воспроизводимости инструменты генерации фиксируют версией: в Go 1.24+ их можно объявить директивой `tool` в `go.mod` (`go get -tool github.com/vektra/mockery/v2`) и запускать как `go tool mockery` — тогда у всех разработчиков и в CI используется одна и та же версия. Сгенерированные файлы принято начинать с комментария `// Code generated ... DO NOT EDIT.` — линтеры и `gofmt`-инструменты такие файлы пропускают.

#### Популярные генераторы кода:

| Инструмент                  | Назначение                                            |
| :-------------------------- | :---------------------------------------------------- |
| **`stringer`**              | Генерация метода `String()` для `iota`-констант       |
| **`mockgen`** (`uber/mock`) | Генерация моков интерфейсов для тестов                |
| **`sqlc`**                  | Генерация типобезопасного Go-кода из `.sql` файлов    |
| **`protoc`**                | Генерация Go-структур и gRPC-сервисов из `.proto`     |
| **`oapi-codegen`**          | Генерация HTTP-клиентов и серверов из OpenAPI/Swagger |

---

## 20. Удаленный вызов процедур (RPC, `net/rpc`) и gRPC

> 🔗 **Практика и примеры кода:** [14-RPC](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/14-RPC) | 🛠️ **Домашнее задание:** [homework-14 (RPC)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/14-RPC/homework-14) | 🎯 **Задачи:** [homework-tasks.md (ДЗ 14: RPC-служба)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/homework-tasks.md#домашнее-задание-14-удалённый-вызов-процедур-rpc-служба-сообщений)

### 20.1. Концепция RPC (Remote Procedure Call) и сравнение с REST API

**Удалённый вызов процедур (RPC)** — это технология, позволяющая программе вызывать функции или методы в другом адресном пространстве (в другом процессе на той же машине или на удалённом сервере по сети) так, будто это обычный локальный вызов функции.

#### 1. Принцип работы и архитектурные компоненты (Stubs):
- **Клиентская заглушка (Client Stub):** локальный прокси-объект, который упаковывает (маршалирует/сериализует) параметры вызова в бинарное или текстовое сообщение и передает его по сетевому протоколу.
- **Сетевой транспорт:** RPC является технологией, а не жестким протоколом. Транспортом может выступать **TCP, HTTP, HTTP/2 или UDP**.
- **Серверная заглушка (Server Stub / Skeleton):** принимает сетевой запрос, распаковывает (демаршалирует) параметры, вызывает реальный метод на стороне сервера, упаковывает полученный результат и отсылает обратно клиенту.
- **Клиент:** принимает ответ, распаковывает результат и возвращает его в вызывающий код. По умолчанию вызов является **синхронным (блокирующим)**, однако поддерживается и асинхронный режим.

#### 2. Сравнительный анализ: RPC vs REST API

| Характеристика | HTTP REST API | RPC (`net/rpc`, gRPC) |
| :--- | :--- | :--- |
| **Основная парадигма** | **Ориентация на ресурсы (Существительные):** моделирование сущностей данных (`/books/123`, `/users/45`). | **Ориентация на действия (Глаголы):** система команд и вызовов бизнес-процедур (`SendMessages`, `CalculateTotal`, `TransferFunds`). |
| **Сетевые методы** | Стандартные глаголы HTTP (`GET`, `POST`, `PUT`, `PATCH`, `DELETE`). | Имена функций предметной области (`Service.Method`). |
| **Формат данных** | Обычно текстовый JSON (реже XML). | Бинарный компактный формат (`encoding/gob` в Go, `Protobuf` в gRPC). |
| **Контракт и типизация** | Слабая или опциональная (OpenAPI/JSON Schema). | **Строгая статическая типизация** типов параметров и результатов на уровне компилятора. |
| **Когда выбирать** | Публичные Web API, мобильные клиенты, CRUD-сервисы, где важна совместимость с браузерами. | Высоконагруженное межсервисное взаимодействие (microservices), сложные команды, которые неестественно ложатся на REST-глаголы. |

---

### 20.2. Стандартный пакет Go `net/rpc`

В стандартную библиотеку Go входит встроенная легковесная реализация RPC поверх TCP или HTTP: [`net/rpc`](https://pkg.go.dev/net/rpc). По умолчанию в качестве формата сериализации используется бинарный кодировщик **`encoding/gob`**.

> [!WARNING]
> Пакет `net/rpc` **заморожен**: он не развивается (команда Go рекомендует gRPC и другие решения), не имеет встроенных аутентификации и шифрования и поддерживает только клиентов на Go. Его изучают для понимания принципов RPC; в новых проектах вместо него берут gRPC.

#### 6 строгих требований к сигнатуре методов в `net/rpc`:
Для того чтобы метод пользовательского типа данных мог быть зарегистрирован и вызван через `net/rpc`, он **обязан** удовлетворять следующим правилам (проверяются через `reflect` при регистрации):
1. Тип данных сервера должен быть **экспортируемым** (`type Service struct`).
2. Метод должен быть **экспортируемым** (`func (s *Service) Method(...) error`).
3. Метод имеет **ровно два аргумента**, оба экспортируемого или встроенного типа.
4. Первый аргумент метода — **входящие параметры** (значение или указатель: `req RequestType`).
5. Второй аргумент метода — **указатель на результат** (`resp *ResponseType`).
6. Метод возвращает **исключительно `error`** (`error`, и ничего больше).

```go
// Сигнатура эталонного RPC-метода:
func (t *T) MethodName(argType T1, replyType *T2) error
```

#### Реализация RPC-сервера:

```go
package main

import (
    "fmt"
    "log"
    "net"
    "net/rpc"
    "sync"
)

type Book struct {
    ID     int
    Title  string
    Author string
}

type LibraryService struct {
    mu    sync.RWMutex
    books []Book
}

// GetBook возвращает книгу по запросу
func (s *LibraryService) GetBook(id int, resp *Book) error {
    s.mu.RLock()
    defer s.mu.RUnlock()
    for _, b := range s.books {
        if b.ID == id {
            *resp = b
            return nil
        }
    }
    return fmt.Errorf("книга с id=%d не найдена", id) // ошибку получит клиент в результате Call
}

func main() {
    svc := &LibraryService{}
    // 1. Регистрация сервиса в RPC
    if err := rpc.Register(svc); err != nil {
        log.Fatal(err)
    }

    // 2. Открытие TCP-слушателя
    listener, err := net.Listen("tcp", ":8080")
    if err != nil {
        log.Fatal(err)
    }
    defer listener.Close()

    // 3. Конкурентная обработка входящих соединений
    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Printf("Accept error: %v", err)
            continue
        }
        go rpc.ServeConn(conn)
    }
}
```

#### Реализация RPC-клиента (синхронный и асинхронный вызов):

```go
package main

import (
    "fmt"
    "log"
    "net/rpc"
)

func main() {
    // 1. Подключение к серверу
    client, err := rpc.Dial("tcp", "localhost:8080")
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // 2. Синхронный (блокирующий) вызов: client.Call
    var foundBook Book
    err = client.Call("LibraryService.GetBook", 1, &foundBook)
    if err != nil {
        log.Fatal("RPC error:", err)
    }
    fmt.Printf("Получена книга: %+v\n", foundBook)

    // 3. Асинхронный (неблокирующий) вызов: client.Go
    asyncCall := client.Go("LibraryService.GetBook", 2, &foundBook, nil)
    // Ожидание завершения через канал Done:
    replyCall := <-asyncCall.Done
    if replyCall.Error != nil {
        log.Fatal("Async error:", replyCall.Error)
    }
}
```

---

### 20.3. Тестирование RPC-сервисов в памяти (`net.Pipe()`)

Тестирование сетевых служб `net/rpc` можно проводить изолированно **без открытия реальных сетевых сокетов операционной системы** с помощью **`net.Pipe()`**:

```go
func TestLibraryService_RPC(t *testing.T) {
    server := rpc.NewServer()
    svc := &LibraryService{}
    _ = server.Register(svc)

    // Создаем синхронную двунаправленную виртуальную сетевую пару в памяти
    serverConn, clientConn := net.Pipe()
    defer serverConn.Close()
    defer clientConn.Close()

    go server.ServeConn(serverConn)

    client := rpc.NewClient(clientConn)
    defer client.Close()

    var res Book
    err := client.Call("LibraryService.GetBook", 1, &res)
    if err != nil {
        t.Fatalf("RPC call failed: %v", err)
    }
}
```

---

### 20.4. Почему gRPC быстрее REST JSON

| Критерий          | REST + JSON                        | gRPC + Protocol Buffers                                                 |
| :---------------- | :--------------------------------- | :---------------------------------------------------------------------- |
| **Формат данных** | Текстовый JSON (медленный парсинг) | Компактный бинарный формат Protobuf (быстрая сериализация, меньше байт) |
| **Транспорт**     | HTTP/1.1 (запросы на одном соединении идут по очереди) | **HTTP/2** (мультиплексирование множества запросов в одном соединении, сжатие заголовков HPACK) |
| **Контракт API**  | OpenAPI/Swagger (часто не строгий) | Строгий `.proto` файл со статической компиляцией                        |
| **Стриминг**      | Ограничен (WebSockets/SSE)         | Двунаправленный потоковый обмен (Client/Server/Bidirectional Streaming) |

---

### 20.5. Protocol Buffers (`.proto`) и кодогенерация

```protobuf
syntax = "proto3";

package movie;
option go_package = "github.com/company/project/gen/movie/v1;moviev1";

service MovieService {
    rpc GetMovie (GetMovieRequest) returns (GetMovieResponse);
}

message GetMovieRequest {
    string id = 1; // 1 — уникальный тег поля в бинарном Protobuf
}

message GetMovieResponse {
    string id = 1;
    string title = 2;
    double rating = 3;
}
```

Генерация Go-кода:

```bash
# Один раз установить плагины генерации:
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Затем (компилятор protoc устанавливается отдельно):
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       api/proto/movie.proto
```

> 💡 Удобная современная альтернатива ручному `protoc` — инструмент **`buf`** (линтинг `.proto`, проверка обратной совместимости и генерация одной командой `buf generate`). Правило совместимости Protobuf: **номера полей (`= 1`, `= 2`) нельзя менять и переиспользовать** после публикации API; удалённые поля помечают `reserved`.

---

### 20.6. Реализация gRPC сервера и клиента на Go

```go
// 1. Сервер:
type MovieServer struct {
    moviev1.UnimplementedMovieServiceServer
}

func (s *MovieServer) GetMovie(ctx context.Context, req *moviev1.GetMovieRequest) (*moviev1.GetMovieResponse, error) {
    if req.GetId() == "" {
        return nil, status.Error(codes.InvalidArgument, "id фильма обязателен")
    }
    return &moviev1.GetMovieResponse{
        Id:     req.GetId(),
        Title:  "Интерстеллар",
        Rating: 8.6,
    }, nil
}

func main() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatal(err)
    }
    grpcServer := grpc.NewServer()
    moviev1.RegisterMovieServiceServer(grpcServer, &MovieServer{})
    // Для плавной остановки при сигнале используйте grpcServer.GracefulStop()
    log.Fatal(grpcServer.Serve(lis))
}
```

```go
// 2. Клиент:
import (
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

func main() {
    // NewClient (grpc-go 1.63+) заменяет устаревший grpc.Dial. insecure — только для локальной разработки: в проде используйте TLS!
    conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    client := moviev1.NewMovieServiceClient(conn)

    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second) // дедлайн на вызов — обязателен
    defer cancel()

    resp, err := client.GetMovie(ctx, &moviev1.GetMovieRequest{Id: "mov_42"})
    if err != nil {
        log.Fatal(err) // ошибки gRPC содержат код: status.Code(err) == codes.InvalidArgument и т.д.
    }
    log.Println(resp.GetTitle())
}
```

---

### 20.7. Enterprise gRPC: Интерцепторы, Метаданные и grpc-gateway

В высоконагруженных микросервисных системах «голый» gRPC не используется — его оборачивают в цепочки связующего ПО (Middleware / Interceptors) для аутентификации, логирования, трейсинга и метрик.

```mermaid
graph LR
    Client["HTTP/JSON клиент (браузер, curl)"] -->|HTTP/1.1 + JSON| Gateway["grpc-gateway (HTTP Reverse Proxy)"]
    Gateway -->|gRPC / HTTP2| Interceptor["Unary Interceptor Chain"]
    Interceptor -->|Context + Tracing| Handler["gRPC Service Handler"]
```

#### 1. Цепочка Интерцепторов (Unary Server Interceptors)

Интерцептор перехватывает каждый входящий RPC-вызов аналогично HTTP-middleware:

```go
package main

import (
    "context"
    "log/slog"
    "time"

    "google.golang.org/grpc"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/metadata"
    "google.golang.org/grpc/status"
)

// LoggingAndRecoveryInterceptor объединяет логирование, тайминг и защиту от паник
func LoggingAndRecoveryInterceptor() grpc.UnaryServerInterceptor {
    return func(
        ctx context.Context,
        req any,
        info *grpc.UnaryServerInfo,
        handler grpc.UnaryHandler,
    ) (resp any, err error) {
        start := time.Now()

        // 1. Извлечение метаданных (Tracing Headers, Auth Tokens):
        if md, ok := metadata.FromIncomingContext(ctx); ok {
            if traceIDs := md.Get("x-trace-id"); len(traceIDs) > 0 {
                slog.Info("incoming RPC", "method", info.FullMethod, "trace_id", traceIDs[0])
            }
        }

        // 2. Защита от паник (Panic Recovery):
        defer func() {
            if r := recover(); r != nil {
                slog.Error("panic recovered in gRPC handler", "error", r, "method", info.FullMethod)
                err = status.Errorf(codes.Internal, "internal server error")
            }
        }()

        // 3. Вызов реального обработчика RPC:
        resp, err = handler(ctx, req)

        // 4. Логирование длительности и статуса:
        duration := time.Since(start)
        slog.Info("completed RPC",
            "method", info.FullMethod,
            "duration_ms", duration.Milliseconds(),
            "code", status.Code(err).String(),
        )

        return resp, err
    }
}

// Регистрация цепочки интерцепторов на gRPC-сервере:
func newEnterpriseGRPCServer() *grpc.Server {
    return grpc.NewServer(
        grpc.ChainUnaryInterceptor(
            LoggingAndRecoveryInterceptor(),
            // Дополнительные интерцепторы: auth, opentelemetry, prometheus
        ),
    )
}
```

#### 2. Передача метаданных и контекстных дедлайнов (Metadata & Deadlines)

> [!NOTE]
> В gRPC контекст — это не просто отмена вызова внутри процесса. Дедлайны (`context.WithTimeout`) автоматически транслируются по сети через HTTP/2 заголовки (`grpc-timeout`). Если вызывающий сервис установил дедлайн 500 мс, принимающий микросервис увидит оставшееся время в своем `ctx.Deadline()`.

```go
// Отправка метаданных клиентом:
func callWithMetadata(ctx context.Context, client moviev1.MovieServiceClient) (*moviev1.GetMovieResponse, error) {
    // Добавляем Outgoing-метаданные (например, Bearer токен):
    md := metadata.Pairs(
        "authorization", "Bearer eyJhbGciOi...",
        "x-trace-id", "req-abc-12345",
    )
    ctx = metadata.NewOutgoingContext(ctx, md)

    // Устанавливаем строгий дедлайн:
    ctx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
    defer cancel()

    return client.GetMovie(ctx, &moviev1.GetMovieRequest{Id: "mov_42"})
}
```

#### 3. Dual Protocol: `grpc-gateway` (HTTP/JSON + gRPC на одном сервисе)

`grpc-gateway` — это плагин компилятора `protoc`, который читает спецификацию `.proto` с аннотациями `google.api.http` и генерирует обратный прокси-сервер. Клиент может отправить обычный HTTP `GET /v1/movies/mov_42`, а шлюз преобразует его в бинарный gRPC-запрос:

```protobuf
// movie.proto с аннотацией HTTP
syntax = "proto3";
package movie.v1;
import "google/api/annotations.proto";

service MovieService {
  rpc GetMovie (GetMovieRequest) returns (GetMovieResponse) {
    option (google.api.http) = {
      get: "/v1/movies/{id}"
    };
  }
}
```

```go
// Инициализация HTTP reverse-proxy шлюза:
import (
    "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

func runHTTPGateway(grpcEndpoint string, httpPort string) error {
    ctx := context.Background()
    mux := runtime.NewServeMux()
    opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())} // insecure — только для локальной среды

    // Регистрируем сгенерированный обработчик шлюза:
    err := moviev1.RegisterMovieServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
    if err != nil {
        return err
    }

    return http.ListenAndServe(httpPort, mux)
}
```

---

## 21. Очереди сообщений и Кэширование (RabbitMQ / Kafka & Redis)

> 🔗 **Практика и примеры кода:** [19-queue](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/19-queue), [20-NoSQL](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/20-NoSQL) | 🛠️ **Домашние задания:** [homework-19 (Kafka)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/19-queue/homework-19), [homework-20 (Lynks)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/20-NoSQL/homework-20) | 🎯 **Задачи:** [homework-tasks.md (ДЗ 19–20)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/homework-tasks.md#домашнее-задание-19-очереди-сообщений-и-асинхронная-аналитика-kafka--event-sourcing)

### 21.1. Асинхронные очереди (Producer-Consumer Pattern)

Очереди сообщений изолируют сервисы от пиковых нагрузок (Traffic Spikes) и обеспечивают надежную доставку:

```mermaid
graph LR
    API["API Gateway (Producer)"] -->|Event: VideoUploaded| Queue["Message Broker (RabbitMQ / Kafka)"]
    Queue -->|Consume| Worker1["Transcoder Worker 1"]
    Queue -->|Consume| Worker2["Transcoder Worker 2"]
```

- **RabbitMQ:** Очередь на базе AMQP брокера (Push-модель, routing keys, Dead Letter Queue).
- **Apache Kafka:** Распределенный лог коммитов (Pull-модель, партиции, Consumer Groups, гарантированный порядок внутри партиции).

---

### 21.2. In-Memory Кэширование с Redis (`go-redis`)

```go
import (
    "context"
    "errors"
    "log/slog"
    "time"

    "github.com/redis/go-redis/v9"
)

func getCachedMovie(ctx context.Context, rdb *redis.Client, movieID string) (string, error) {
    key := "movie:" + movieID

    // 1. Cache-Aside Pattern: проверяем кэш
    val, err := rdb.Get(ctx, key).Result()
    if err == nil {
        return val, nil // Cache Hit!
    }
    if !errors.Is(err, redis.Nil) {
        // Любая ошибка, кроме redis.Nil (ключа нет) — проблема с Redis: логируем и идём в БД (кэш не критичен)
        slog.Warn("redis недоступен", "error", err)
    }

    // 2. Cache Miss: читаем из PostgreSQL
    movieData := fetchFromPostgres(movieID)

    // 3. Сохраняем в кэш с TTL (Time-To-Live) на 10 минут:
    _ = rdb.Set(ctx, key, movieData, 10*time.Minute).Err()

    return movieData, nil
}
```

> ⚠️ Типичные проблемы кэширования (частые вопросы на собеседованиях):
> - **Cache Stampede** («давка»): когда популярный ключ истёк, сотни запросов одновременно идут в БД. Решения: `singleflight` (`golang.org/x/sync/singleflight`) для объединения одинаковых запросов, случайный разброс TTL (jitter), фоновое обновление.
> - **Инвалидация:** при изменении данных нужно удалять/обновлять ключ (`rdb.Del`); без этого пользователи видят устаревшие данные до истечения TTL.
> - **Сериализация:** в Redis лежат строки/байты — структуры кодируют в JSON или Protobuf.

---

### 21.3. `Transactional Outbox Pattern` (Решение проблемы Dual Write)

> [!IMPORTANT]
> **Проблема двух записей (Dual Write Problem):**  
> Представьте, что при оформлении заказа вам нужно:
>
> 1. Сохранить заказ в PostgreSQL (`db.Exec`).
> 2. Отправить событие `OrderCreated` в Kafka (`kafka.Produce`).
>
> Если упадет сеть между 1 и 2 шагом — заказ в БД сохранится, а сообщение в Kafka **не уйдет** (рассинхрон данных!). Распределенные транзакции (2PC) медленные и ненадежные.

#### Решение через паттерн Transactional Outbox:

Мы сохраняем бизнес-сущность и само событие **в одной локальной ACID-транзакции PostgreSQL**:

```mermaid
graph LR
    API["API Gateway"] -->|1. Tx: Save Order + Insert Outbox| DB[("PostgreSQL\norders & outbox tables")]
    Worker["Background Relay Worker"] -->|2. Poll / CDC (Debezium)| DB
    Worker -->|3. Publish Event| Kafka["Kafka / RabbitMQ"]
    Worker -->|4. Mark as Processed| DB
```

```go
// Выполняем строго в одной транзакции (проверки ошибок после каждого Exec опущены для краткости —
// в реальном коде после каждого вызова нужно делать `if err != nil { return err }`):
tx, err := dbPool.Begin(ctx)
if err != nil {
    return err
}
defer tx.Rollback(ctx)

// 1. Сохраняем заказ в таблицу orders:
_, err = tx.Exec(ctx, "INSERT INTO orders (id, user_id, amount) VALUES ($1, $2, $3)", orderID, userID, 1500)

// 2. Сохраняем событие в таблицу outbox:
_, err = tx.Exec(ctx, `INSERT INTO outbox (id, event_type, payload, status)
                       VALUES ($1, 'OrderCreated', $2, 'PENDING')`, eventID, payloadJSON)

// 3. Фиксация транзакции (или оба сохранятся, или ни один):
return tx.Commit(ctx)
```

_Фоновый воркер (Relay)_ периодически считывает `PENDING` записи из `outbox`, отправляет их в Kafka и после подтверждения меняет статус на `PROCESSED`. Чтобы несколько экземпляров воркера не отправляли одни и те же записи, выборку делают с блокировкой: `SELECT ... FROM outbox WHERE status = 'PENDING' ORDER BY created_at LIMIT 100 FOR UPDATE SKIP LOCKED`.

> ⚠️ Если воркер упадёт после отправки в Kafka, но до смены статуса, событие уйдёт повторно. Поэтому Outbox даёт гарантию **at-least-once** (не «ровно один раз»), а потребители событий обязаны быть **идемпотентными** (например, хранить идентификаторы уже обработанных событий).

---

### 21.4. Распределенные блокировки (`Distributed Locks`)

В микросервисной архитектуре, когда запущено 10 реплик одного сервиса, обычный `sync.Mutex` защищает память только _одного инстанса_. Чтобы операцию гарантированно выполнял только один сервис (например, генерацию ежемесячного отчета), используют **распределенный лок через Redis**:

```go
// Захват блокировки на 10 секунд через атомарную команду SET NX EX.
// uniqueToken — случайное уникальное значение (например, UUID) владельца лока:
func acquireLock(ctx context.Context, rdb *redis.Client, lockKey, uniqueToken string) bool {
    success, err := rdb.SetNX(ctx, lockKey, uniqueToken, 10*time.Second).Result()
    return err == nil && success
}

// Освобождение блокировки (через Lua-скрипт, чтобы случайно не снять чужой лок):
const unlockScript = `
if redis.call("get", KEYS[1]) == ARGV[1] then
    return redis.call("del", KEYS[1])
else
    return 0
end`

func releaseLock(ctx context.Context, rdb *redis.Client, lockKey, uniqueToken string) {
    rdb.Eval(ctx, unlockScript, []string{lockKey}, uniqueToken)
}
```

> [!WARNING]
> Такая блокировка **не абсолютна**: если работа заняла дольше TTL (10 секунд), лок истечёт, и второй экземпляр войдёт в критическую секцию одновременно с первым (а при отказе и переключении Redis-реплики лок может «потеряться»). Поэтому для критичных данных добавляют **fencing token** (монотонно растущий номер, который проверяет само хранилище) и/или продлевают лок в фоне, а если нужна строгая гарантия — используют блокировки на уровне СУБД (`pg_advisory_lock`, `SELECT ... FOR UPDATE`) или консенсус-системы (etcd, ZooKeeper). Для простых задач вроде «раз в сутки запустить отчёт» Redis-лок подходит.

---

### 21.5. Практика работы с Apache Kafka в Go (`segmentio/kafka-go`)

Apache Kafka представляет собой распределенный горизонтально масштабируемый журнал коммитов (Distributed Commit Log). В отличие от классических брокеров (RabbitMQ), где сообщения удаляются сразу после подтверждения (ACK), Kafka сохраняет все события на диск в течение заданного времени хранения (Retention Period).

#### 1. Сравнение RabbitMQ и Apache Kafka:

| Характеристика | RabbitMQ | Apache Kafka |
| :--- | :--- | :--- |
| **Модель доставки** | **Push-модель** (брокер отправляет потребителю) | **Pull-модель** (потребитель сам вычитывает пакеты) |
| **Хранение сообщений** | В очереди (в памяти и, для durable-очередей, на диске), удаляются после ACK | На диске в append-only логе, хранятся в течение срока хранения (retention) |
| **Маршрутизация** | Гибкая (Exchange: direct, topic, fanout, headers) | Простая (по топикам и партициям через ключ) |
| **Повторное чтение** | Нет (сообщения удаляются) | **Да** (перемещение Offset назад во времени) |
| **Пропускная способность** (порядок величин, зависит от железа и настроек) | ~10k–50k сообщений/сек | **100k–1M+ сообщений/сек** (Sequential I/O, Zero-Copy) |
| **Сценарии** | Сложная маршрутизация задач, RPC | Event Sourcing, аудит, потоковая аналитика, CDC |

#### 2. Реализация Kafka Producer (`segmentio/kafka-go`):

```go
import (
    "context"
    "encoding/json"
    "time"
    "github.com/segmentio/kafka-go"
)

type EventProducer struct {
    writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) *EventProducer {
    return &EventProducer{
        writer: &kafka.Writer{
            Addr:                   kafka.TCP(brokers...),
            Topic:                  topic,
            // Hash: партиция выбирается по ключу сообщения, поэтому сообщения с одним ключом
            // попадают в одну партицию и сохраняют порядок (LeastBytes/RoundRobin порядок по ключу НЕ гарантируют!)
            Balancer:               &kafka.Hash{},
            AllowAutoTopicCreation: true, // удобно локально; в проде топики обычно создают заранее
            WriteTimeout:           10 * time.Second,
            RequiredAcks:           kafka.RequireAll, // ждём подтверждения от всех in-sync реплик (RequireOne — быстрее, но менее надёжно)
        },
    }
}

func (p *EventProducer) Send(ctx context.Context, key string, data any) error {
    payload, err := json.Marshal(data)
    if err != nil {
        return err
    }
    return p.writer.WriteMessages(ctx, kafka.Message{
        Key:   []byte(key), // Ключ определяет партицию (при Balancer Hash) и, значит, порядок сообщений с этим ключом
        Value: payload,
        Time:  time.Now(),
    })
}
```

#### 3. Реализация Kafka Consumer и выбор стратегии фиксации (Commit):

> [!IMPORTANT]
> **Выбор момента подтверждения сообщения (Commit):**
> 1. **At-Most-Once (Подтверждение ➔ Обработка):** оффсет коммитится сразу после получения сообщения. Если сервис упадет во время обработки — сообщение потеряно навсегда.
> 2. **At-Least-Once (Обработка ➔ Подтверждение):** оффсет фиксируется строго **после** успешной обработки в бизнес-логике. При сбое сообщение вычитается повторно. Требует **идемпотентности** на стороне потребителя!

```go
type EventConsumer struct {
    reader *kafka.Reader
}

func NewConsumer(brokers []string, topic, groupID string) *EventConsumer {
    return &EventConsumer{
        reader: kafka.NewReader(kafka.ReaderConfig{
            Brokers:        brokers,
            Topic:          topic,
            GroupID:        groupID, // Consumer Group для параллельного чтения партиций
            MinBytes:       10e1,    // 100B
            MaxBytes:       10e6,    // 10MB
            CommitInterval: 0,       // 0 = коммиты выполняются синхронно в момент вызова CommitMessages (мы сами решаем, когда)
        }),
    }
}

func (c *EventConsumer) Start(ctx context.Context, handler func([]byte) error) error {
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
        }

        // 1. Вычитываем сообщение:
        msg, err := c.reader.FetchMessage(ctx)
        if err != nil {
            return err
        }

        // 2. Обрабатываем в бизнес-логике:
        if err := handler(msg.Value); err != nil {
            // ⚠️ Ошибка обработки. Просто `continue` здесь НЕ приведёт к повторной доставке:
            // читатель уже продвинулся к следующему сообщению, а последующий CommitMessages
            // зафиксирует смещение уже ПОСЛЕ этого сообщения, и оно будет потеряно.
            // Нужно либо повторить обработку в цикле с backoff, либо отправить сообщение в DLQ (см. 21.6),
            // либо остановить консьюмер (return err) и перечитать сообщение после перезапуска.
            continue
        }

        // 3. Фиксируем оффсет только после успешной обработки (At-Least-Once):
        _ = c.reader.CommitMessages(ctx, msg)
    }
}
```

---

### 21.6. Отказоустойчивость Apache Kafka: Rebalance, Poison Pill, DLQ и Debezium CDC

#### 1. Протоколы ребалансировки (Rebalance Protocols): Eager vs Cooperative Sticky

Когда в Consumer Group добавляется новый консьюмер или падает существующий, Kafka перераспределяет партиции между участниками:

```mermaid
graph TD
    subgraph Eager ["Eager Protocol (Stop-The-World)"]
        E1["Все консьюмеры отдают партиции"] --> E2["Полная пауза обработки (Lag Spike)"]
        E2 --> E3["Переназначение партиций с нуля"]
    end
    subgraph Cooperative ["Cooperative Sticky Protocol"]
        C1["Забираются только перемещаемые партиции"] --> C2["Остальные консьюмеры продолжают непрерывно читать"]
    end
```

- **Eager Rebalance (`Range`, `RoundRobin`):** Консьюмеры бросают все партиции и ждут полного перераспределения. Происходит «Stop-The-World», резко растут задержки и лаг.
- **Cooperative Sticky Rebalance (`CooperativeStickyAssignor`):** Партиции отзываются инкрементально в две фазы. Консьюмеры, чьи партиции не затронуты миграцией, **продолжают чтение без остановки**. (Поддержка зависит от клиентской библиотеки: в Java-клиенте и в `franz-go` она есть, а `segmentio/kafka-go` использует Range/RoundRobin-стратегии — проверяйте документацию своей библиотеки.)

#### 2. Проблема «Ядовитой таблетки» (Poison Pill) и Dead Letter Queue (DLQ)

> [!WARNING]
> Если сообщение в топике повреждено (битый JSON, невалидная схема Protobuf), наивный консьюмер падает в панику или бесконечно повторяет чтение. Оффсет не сдвигается, и вся очередь для данной партиции намертво блокируется (Head-of-Line Blocking).

**Паттерн Dead Letter Queue (Очередь мертвых писем):**
При превышении лимита повторных попыток (Max Retries) консьюмер перекладывает невалидное сообщение в специальный топик `orders.dlq` с сохранением исходных заголовков ошибки, после чего коммитит оффсет в основном топике:

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    "github.com/segmentio/kafka-go"
)

type ConsumerWithDLQ struct {
    reader    *kafka.Reader
    dlqWriter *kafka.Writer
}

func (c *ConsumerWithDLQ) ProcessMessages(ctx context.Context) {
    for {
        msg, err := c.reader.FetchMessage(ctx)
        if err != nil {
            return
        }

        var order map[string]any
        if err := json.Unmarshal(msg.Value, &order); err != nil {
            // Poison Pill обнаружен: отправляем в топик orders.dlq
            c.sendToDLQ(ctx, msg, fmt.Sprintf("unmarshal error: %v", err))
            // (в реальном коде: если запись в DLQ не удалась, оффсет НЕ коммитим, иначе сообщение потеряется)
            _ = c.reader.CommitMessages(ctx, msg) // сдвигаем оффсет, разблокируя партицию
            continue
        }

        // Успешная бизнес-обработка:
        _ = c.reader.CommitMessages(ctx, msg)
    }
}

func (c *ConsumerWithDLQ) sendToDLQ(ctx context.Context, origMsg kafka.Message, reason string) {
    dlqMsg := kafka.Message{
        Key:   origMsg.Key,
        Value: origMsg.Value,
        Headers: []kafka.Header{
            {Key: "x-original-topic", Value: []byte(origMsg.Topic)},
            {Key: "x-error-reason", Value: []byte(reason)},
            {Key: "x-failed-at", Value: []byte(time.Now().UTC().Format(time.RFC3339))},
        },
    }
    _ = c.dlqWriter.WriteMessages(ctx, dlqMsg)
}
```

#### 3. Change Data Capture (CDC) и Debezium

Вместо периодического опроса БД (Polling) или двухфазных распределенных транзакций (2PC) современная архитектура использует **CDC (Change Data Capture)**:

```mermaid
graph LR
    App["Go Application"] -->|1. Локальная транзакция| DB[("PostgreSQL\n(Таблица outbox)")]
    DB -->|2. Чтение WAL (Write-Ahead Log)| Debezium["Debezium CDC Connector"]
    Debezium -->|3. Streaming без задержек| Kafka["Kafka Topic: events.order"]
```

1. Go-сервис пишет в основную бизнес-таблицу и в таблицу `outbox` в рамках **одной локальной SQL-транзакции** (ACID).
2. **Debezium** слушает системный лог репликации PostgreSQL (`pgoutput` / WAL).
3. При коммите транзакции Debezium с минимальной задержкой парсит изменения и отправляет сообщение в Kafka.
4. **Результат:** гарантия доставки **at-least-once** (потребители должны быть идемпотентными), отсутствие нагрузки на БД опросными `SELECT`-запросами и отсутствие рассинхронизации Dual-Write.

> ⚠️ Слот репликации PostgreSQL удерживает WAL, пока Debezium не прочитает данные: если коннектор остановлен надолго, WAL может заполнить диск сервера БД — за слотами нужно следить.

---

## 22. Наблюдаемость систем (Observability: Prometheus & OpenTelemetry)

### 22.1. Метрики Prometheus

В продакшене здоровье и производительность Go-сервисов оценивают по 3 ключевым типам метрик пакета `github.com/prometheus/client_golang/prometheus`:

1. **Counter (Счетчик):** монотонно возрастающее значение (общее количество HTTP-запросов, число ошибок). Сбрасывается только при перезапуске сервиса.
2. **Gauge (Датчик):** значение, которое может как увеличиваться, так и уменьшаться (число активных WebSocket-соединений, размер очереди в памяти, занятые потоки пула БД).
3. **Histogram (Гистограмма):** распределение значений по корзинам (buckets), используется для замера **Latency (времени ответа)**: $p50, p95, p99$.

```go
import (
    "net/http"
    "strconv"
    "time"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    httpRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Общее число обработанных HTTP запросов",
        },
        []string{"method", "path", "status"},
    )

    httpRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "Длительность обработки запросов в секундах",
            Buckets: prometheus.DefBuckets, // [0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10]
        },
        []string{"path"},
    )
)

// Эндпоинт метрик для сбора Prometheus-сервером:
// http.Handle("/metrics", promhttp.Handler())

// Использование в middleware:
func MetricsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK} // обёртка для перехвата статус-кода
        next.ServeHTTP(rec, r)

        // ВАЖНО: в label кладём ШАБЛОН маршрута ("/users/{id}"), а не реальный путь ("/users/42") — см. предупреждение ниже
        route := r.Pattern // Go 1.23+: шаблон, по которому был найден маршрут
        httpRequestsTotal.WithLabelValues(r.Method, route, strconv.Itoa(rec.status)).Inc()
        httpRequestDuration.WithLabelValues(route).Observe(time.Since(start).Seconds())
    })
}

type statusRecorder struct {
    http.ResponseWriter
    status int
}

func (s *statusRecorder) WriteHeader(code int) {
    s.status = code
    s.ResponseWriter.WriteHeader(code)
}
```

> [!WARNING]
> **Кардинальность меток (label cardinality):** каждая уникальная комбинация значений меток создаёт отдельный временной ряд в памяти Prometheus. Если поместить в метку `user_id` или реальный URL с идентификаторами, число рядов вырастет до миллионов, и Prometheus «умрёт» от нехватки памяти. Метки — только для значений из небольшого фиксированного набора (метод, шаблон маршрута, код статуса).

Перцентили ($p95$, $p99$) считаются на стороне Prometheus запросом PromQL: `histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (le, path))`.

---

### 22.2. Распределенная трассировка (OpenTelemetry / Jaeger)

В микросервисной архитектуре один пользовательский клик порождает цепочку вызовов через 10 микросервисов. OpenTelemetry передает `TraceID` и `SpanID` через заголовки (`traceparent`):

```mermaid
graph LR
    Client["Client Request"] -->|TraceID: abc, Span 1| GW["API Gateway"]
    GW -->|TraceID: abc, Span 2 (traceparent)| Auth["Auth Service"]
    GW -->|TraceID: abc, Span 3 (traceparent)| Billing["Billing Service"]
    Billing -->|TraceID: abc, Span 4| DB[(PostgreSQL)]
```

#### 1. Стандарт W3C TraceContext (`traceparent` header)

Межсервисный контекст передается по сети в едином формате W3C:
`traceparent: 00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01`
- `00` — версия спецификации.
- `4bf92f3577b34da6a3ce929d0e0e4736` — глобальный **Trace ID** (один на всю цепочку запроса).
- `00f067aa0ba902b7` — **Parent Span ID** (идентификатор текущего шага).
- `01` — флаги трассировки (`01` = сэмплировать и записать в хранилище).

#### 2. Создание спанов и проброс контекста в Go (`go.opentelemetry.io/otel`)

```go
package main

import (
    "context"
    "net/http"

    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/propagation"
    "go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("order-service")

// ⚠️ При старте приложения ОБЯЗАТЕЛЬНО настройте глобальные провайдер и пропагатор, иначе спаны и
// Inject/Extract будут «пустышками» (по умолчанию пропагатор — no-op):
//   otel.SetTracerProvider(tp)                                   // tp — TracerProvider с экспортёром (например, OTLP → Jaeger/Tempo)
//   otel.SetTextMapPropagator(propagation.TraceContext{})        // формат W3C traceparent

func ProcessOrder(ctx context.Context, orderID string) error {
    // 1. Создаем дочерний спан с атрибутами для Jaeger / Grafana Tempo:
    ctx, span := tracer.Start(ctx, "ProcessOrder",
        trace.WithAttributes(attribute.String("order.id", orderID)),
    )
    defer span.End() // Обязательно завершаем спан для фиксации тайминга

    // 2. Вызов внешнего сервиса оплаты с пробросом W3C заголовков:
    req, _ := http.NewRequestWithContext(ctx, "POST", "http://payment-service/pay", nil)
    
    // Внедряем traceparent заголовок в исходящий HTTP-запрос:
    otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

    // 3. Выполняем сетевой запрос...
    return nil
}

// На практике вместо ручных Inject/Extract используют готовые обёртки go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp
// (otelhttp.NewHandler для сервера и otelhttp.NewTransport для http.Client) и otelgrpc для gRPC.

// Middleware на принимающей стороне (Payment Service):
func TraceMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Извлекаем TraceID и ParentSpanID из заголовков входящего запроса:
        ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))
        
        ctx, span := tracer.Start(ctx, r.Method+" "+r.URL.Path)
        defer span.End()

        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

---

## 23. Cloud-Ready приложения и 12-Factor App в Go

### 23.1. Конфигурация через переменные окружения

По методологии **12-Factor App** конфигурация сервиса должна передаваться через **Environment Variables** (переменные среды), а не жестко зашиваться в код или файлы:

```go
import "github.com/caarlos0/env/v10"

type Config struct {
    Port        int    `env:"PORT" envDefault:"8080"`
    DatabaseURL string `env:"DATABASE_URL,required"`
    LogLevel    string `env:"LOG_LEVEL" envDefault:"info"`
}

func LoadConfig() (*Config, error) {
    var cfg Config
    if err := env.Parse(&cfg); err != nil {
        return nil, err
    }
    return &cfg, nil
}
```

---

### 23.2. Health checks: Liveness & Readiness Probes для Kubernetes

В контейнерах Kubernetes разделяют 2 типа проверок здоровья:

- **`/healthz/liveness` (Жив ли процесс):** проверяет, что процесс не завис в дедлоке. Если возвращает ошибку (не 200 OK) — Kubernetes **перезапускает (restart)** контейнер.
- **`/healthz/readiness` (Готов ли принимать трафик):** проверяет доступность зависимостей (успешный `Ping()` в PostgreSQL и Redis, прогрев кэша). Если не готов — Kubernetes временно **снимает с пода входящий трафик**, не убивая его.

> [!WARNING]
> Liveness-проба **не должна** проверять внешние зависимости (БД, Redis): если БД временно недоступна, Kubernetes начнёт перезапускать все поды сразу — и усугубит аварию. Зависимости проверяет только readiness.

Подключение проб в манифесте Deployment:

```yaml
livenessProbe:
  httpGet: { path: /healthz/liveness, port: 8080 }
  periodSeconds: 10
readinessProbe:
  httpGet: { path: /healthz/readiness, port: 8080 }
  periodSeconds: 5
# startupProbe нужен медленно стартующим сервисам: пока он не пройдёт, liveness не запускается
```

```go
func RegisterHealthRoutes(mux *http.ServeMux, db *pgxpool.Pool) {
    // 1. Liveness:
    mux.HandleFunc("GET /healthz/liveness", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        _, _ = w.Write([]byte(`{"status":"alive"}`))
    })

    // 2. Readiness:
    mux.HandleFunc("GET /healthz/readiness", func(w http.ResponseWriter, r *http.Request) {
        if err := db.Ping(r.Context()); err != nil {
            http.Error(w, `{"status":"database unreachable"}`, http.StatusServiceUnavailable)
            return
        }
        w.WriteHeader(http.StatusOK)
        _, _ = w.Write([]byte(`{"status":"ready"}`))
    })
}
```

---

### 23.3. Docker: Multi-Stage Build для Go-приложений

> [!IMPORTANT]
> **Вопрос на собеседовании:** «Как уменьшить Docker-образ Go-приложения с 1 ГБ до 10 МБ?»  
> **Ответ:** Multi-stage build + статическая линковка + образ `scratch` (или `distroless`).

```dockerfile
# ========== Этап 1: Сборка (Build Stage) ==========
# Версия образа должна соответствовать директиве go в go.mod (актуальную смотрите на hub.docker.com/_/golang)
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Копируем и скачиваем зависимости отдельно для кэширования Docker-слоя:
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код и собираем статический бинарник:
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /app/server ./cmd/api

# ========== Этап 2: Минимальный образ для запуска ==========
FROM scratch

# Копируем только скомпилированный бинарник (без Go SDK, без ОС):
COPY --from=builder /app/server /server

# Копируем CA-сертификаты (если нужны HTTPS-запросы наружу):
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Не запускаем процесс от root: в scratch нет /etc/passwd, поэтому используем числовой UID
USER 65532:65532

EXPOSE 8080
ENTRYPOINT ["/server"]
```

#### Ключевые флаги сборки:

| Флаг                      | Назначение                                                                          |
| :------------------------ | :---------------------------------------------------------------------------------- |
| `CGO_ENABLED=0`           | Отключает привязку к C-библиотекам (позволяет запускать в `scratch` / `distroless`) |
| `-ldflags="-s -w"`        | Убирает таблицу символов и DWARF-отладку (заметно уменьшает бинарник, обычно на 25–35%) |
| `-trimpath`               | Убирает из бинарника локальные пути к исходникам (воспроизводимая сборка, меньше утечки информации) |
| `GOOS=linux GOARCH=amd64` | Кросс-компиляция под целевую ОС                                                     |

#### Выбор базового образа:

| Образ                      | Размер     | Когда использовать                                                     |
| :------------------------- | :--------- | :--------------------------------------------------------------------- |
| `scratch`                  | **0 байт** | Максимальная минимальность (нет shell, нет пакетов)                    |
| `gcr.io/distroless/static` | **~2 МБ**  | Рекомендация Google: без shell, но есть CA-сертификаты и timezone data |
| `alpine`                   | **~7 МБ**  | Когда нужен shell для отладки (`sh`, `wget`)                           |

---

## 24. Production-экосистема Go: Топ сторонних библиотек и инструментов

В реальной Enterprise-разработке (BigTech и FinTech) стандартную библиотеку дополняет устоявшийся набор проверенных библиотек («Золотой стандарт»):

### 24.1. Сводная таблица золотого стандарта Go-библиотек

| Категория                 | Стандартная библиотека | Production-стандарт (Библиотеки)                               | Зачем используется                                          |
| :------------------------ | :--------------------- | :------------------------------------------------------------- | :---------------------------------------------------------- |
| **HTTP Роутинг**          | `net/http`             | **`go-chi/chi/v5`**, **`gin-gonic/gin`**                       | Роутинг с middleware и path-параметрами                     |
| **PostgreSQL**            | `database/sql`         | **`jackc/pgx/v5`** (`pgxpool`)                                 | Высокопроизводительный пул и бинарный протокол PG           |
| **SQL & Query Builder**   | Raw SQL строки         | **`sqlc`**, **`Masterminds/squirrel`**                         | Генерация типобезопасного кода из SQL / сборщик запросов    |
| **Миграции БД**           | Отсутствует            | **`pressly/goose/v3`**, **`golang-migrate`**                   | Версионирование схемы БД в `.sql` и Go коде                 |
| **Кэширование**           | Отсутствует            | **`redis/go-redis/v9`**                                        | Клиент Redis с поддержкой Cluster и Sentinel                |
| **Очереди сообщений**     | Каналы (in-memory)     | **`segmentio/kafka-go`**, **`rabbitmq/amqp091-go`**            | Брокеры сообщений для микросервисов                         |
| **RPC / Межсервис**       | `net/rpc`              | **`google.golang.org/grpc`**, **Protobuf**                     | Бинарное сетевое взаимодействие по HTTP/2                   |
| **Конфигурация**          | `os.Getenv`            | **`caarlos0/env/v10`**, **`spf13/viper`**                      | Декларативный маппинг переменных окружения в структуры      |
| **Валидация данных**      | Ручные `if/else`       | **`go-playground/validator/v10`**                              | Валидация структур через теги (`validate:"required,email"`) |
| **Тестирование & Моки**   | `testing`              | **`stretchr/testify`** (`assert`, `require`, `mock`)           | Удобные утверждения тестов и генерация моков                |
| **Метрики & Трассировка** | Отсутствует            | **`prometheus/client_golang`**, **`go.opentelemetry.io/otel`** | Мониторинг Prometheus, Grafana, Jaeger                      |
| **Линтер**                | `go vet`               | **`golangci-lint`**                                            | Мета-линтер (объединяет 50+ анализаторов кода)              |

---

### 24.2. Линтеры и Статический анализ (`golangci-lint`)

В CI/CD пайплайнах запуск `golangci-lint` обязателен. Он перехватывает ошибки, которые пропускает компилятор:

```bash
# Установка (macOS): brew install golangci-lint   (другие способы — на golangci-lint.run)
# Запуск:
golangci-lint run ./...
```

#### Ключевые линтеры:

- `govet` — встроенный анализатор Go на базовые ошибки.
- `errcheck` — проверяет, что ни одна возвращаемая ошибка не проигнорирована (`_ = func()`).
- `staticcheck` — глубокий статический анализ (устаревший код, неиспользуемые переменные).
- `gosec` — аудит безопасности (SQL-инъекции, хардкод паролей, небезопасные числа).
- `bodyclose` — проверяет, закрыто ли тело HTTP-ответа (`defer resp.Body.Close()`).

---

### 24.3. Валидация входных данных (`go-playground/validator`)

Автоматическая валидация JSON-запросов через теги структур:

```go
import "github.com/go-playground/validator/v10"

type CreateUserRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Age      int    `json:"age" validate:"required,gte=18,lte=120"`
    Role     string `json:"role" validate:"oneof=admin moderator viewer"`
}

var validate = validator.New()

func validateRequest(req CreateUserRequest) error {
    return validate.Struct(req) // вернет детальную ошибку, если email невалиден или age < 18
}
```

---

### 24.4. Зеркалирование и тестирование трафика: GoReplay (утилита `gor`)

> ⚠️ Зеркалирование боевого трафика означает копирование **реальных пользовательских данных** (токены, персональные данные) в другую среду. Маскируйте чувствительные поля и убедитесь, что это допустимо с точки зрения политики безопасности и законодательства о персональных данных.

> [!IMPORTANT]
> **Что такое `gor` (GoReplay):**  
> `GoReplay` (команда в терминале: `gor`) — это популярный open-source инструмент, написанный на Go (`github.com/buger/goreplay`), предназначенный для перехвата сетевого HTTP-трафика и его повторного воспроизведения (_Traffic Shadowing / Mirroring_).  
> Он работает на уровне сырых сокетов (pcap / raw sockets), перехватывая реальные входящие HTTP-запросы пользователей на продакшене без модификации кода приложения и без добавления задержек (latency).

#### 💡 Что позволяет делать команда `gor`:

1. **Теневое тестирование (Shadow Testing / Traffic Mirroring):**
   Дублирует реальный боевой трафик с продакшена на тестовый контур (Staging) в реальном времени. Это позволяет проверить новую версию микросервиса на реальном пользовательском потоке запросов до релиза в прод:

   ```bash
   # Захват входящего HTTP-трафика с порта 8080 и репликация на staging:
   gor --input-raw :8080 --output-http "http://staging.example.com:8080"
   ```

2. **Запись и воспроизведение трафика (Record & Replay):**
   Сохраняет боевые запросы в файл и воспроизводит их в любой момент (например, для локального воспроизведения бага или регрессионного тестирования):

   ```bash
   # Запись боевого трафика в файл:
   gor --input-raw :8080 --output-file requests.gor

   # Воспроизведение записанного трафика на локальный сервис:
   gor --input-file requests.gor --output-http "http://localhost:8080"
   ```

3. **Стресс-тестирование с ускорением нагрузки (Load Testing):**
   Позволяет воспроизвести реальный пик нагрузки с ускорением в 2x (200%), 5x или 10x для проверки устойчивости:

   ```bash
   # Воспроизведение с удвоенной скоростью нагрузки:
   gor --input-file "requests.gor|200%" --output-http "http://staging.example.com:8080"
   ```

4. **Фильтрация и семплирование трафика:**
   Позволяет реплицировать только процент запросов (например, 10%) или фильтровать по HTTP-методу/URL:
   ```bash
   # Репликация только 10% трафика и только POST-запросов:
   gor --input-raw :8080 --output-http "http://staging:8080|10%" --http-allow-method POST
   ```

---

### 24.5. Золотой стандарт конфигурации `.golangci.yml`

В реальных Enterprise-проектах стандартного набора линтеров недостаточно. Ниже приведена конфигурация с разделением на критические категории.

> [!NOTE]
> Пример написан для формата конфигурации **golangci-lint v1**. В версии **v2** (2025) формат изменился: файл начинается с `version: "2"`, вместо `disable-all: true` используется `linters: default: none`, линтер `typecheck` убран (проверка типов выполняется всегда), а `gosimple`/`stylecheck` объединены в `staticcheck`. Сверяйтесь с документацией установленной у вас версии (`golangci-lint --version`, команда `golangci-lint migrate` помогает преобразовать конфиг).

```yaml
run:
  timeout: 5m
  tests: true

linters:
  disable-all: true
  enable:
    # Базовые проверки и баги
    - errcheck      # Проверка необработанных ошибок
    - govet         # Официальный анализатор Go
    - staticcheck   # Статический анализ от автора Go-инструментов
    - typecheck     # Проверка типов
    - unused        # Поиск неиспользуемых констант, переменных, функций

    # Безопасность и утечки памяти
    - gosec         # Аудит уязвимостей (SQLi, слабые крипто-алгоритмы)
    - bodyclose     # Контроль закрытия resp.Body в HTTP-клиентах
    - sqlclosecheck # Контроль закрытия sql.Rows и sql.Stmt
    - noctx         # Проверка отправки HTTP-запросов без context.Context

    # Сложность и архитектурный стиль
    - gocognit      # Когнитивная сложность функций (порог: <= 20)
    - revive        # Быстрая и гибкая замена golint
    - misspell      # Исправление опечаток в комментариях и именах
    - unparam       # Поиск неиспользуемых параметров функций
    - prealloc      # Рекомендации по pre-allocate слайсов (make([]T, 0, n))

issues:
  exclude-use-default: false
  max-issues-per-linter: 0
  max-same-issues: 0
```

---

### 24.6. Продакшен `Makefile` для автоматизации разработки

В профессиональных репозиториях сборка, линтинг, миграции и запуск тестов с проверкой гонок стандартизированы через `Makefile`:

```makefile
.PHONY: build test lint cover run clean docker-build

APP_NAME = server
BIN_DIR = ./bin
MAIN_SRC = ./cmd/api/main.go

# Сборка статического бинарника без отладочных символов (-s -w уменьшает вес на ~30%)
build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -o $(BIN_DIR)/$(APP_NAME) $(MAIN_SRC)

# Запуск тестов со включенным детектором гонок памяти
test:
	go test -v -race -timeout 60s ./...

# Запуск тестов с генерацией отчета покрытия кода (Code Coverage)
cover:
	go test -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html

# Запуск мета-линтера (результаты анализа кэшируются автоматически)
lint:
	golangci-lint run ./...

# Локальный запуск приложения
run:
	go run $(MAIN_SRC)

# Сборка минимального контейнера
docker-build:
	docker build -t $(APP_NAME):latest .

clean:
	rm -rf $(BIN_DIR) coverage.out coverage.html
```

---

### 24.7. Паттерн Feature Flags (Управление фичами на лету)

> [!NOTE]
> Feature Flags позволяют включать или выключать новую функциональность в продакшене без пересборки бинарника и без деплоя, а также проводить канареечные релизы (Canary) и A/B тестирование. В примере ниже показан только потокобезопасный «переключатель» в памяти: источником значения обычно служит внешняя система (админ-эндпоинт, конфиг-сервис, Unleash, LaunchDarkly, стандарт OpenFeature), которая вызывает `SetNewCheckout` при изменении.

```go
package main

import (
    "context"
    "sync/atomic"
)

// FeatureManager потокобезопасно хранит состояние флагов
type FeatureManager struct {
    newCheckoutFlow atomic.Bool
}

func (fm *FeatureManager) IsNewCheckoutEnabled() bool {
    return fm.newCheckoutFlow.Load()
}

func (fm *FeatureManager) SetNewCheckout(enabled bool) {
    fm.newCheckoutFlow.Store(enabled)
}

// Использование в обработчике заказа:
func HandleCheckout(ctx context.Context, fm *FeatureManager) {
    if fm.IsNewCheckoutEnabled() {
        // Запуск нового алгоритма оплаты
        return
    }
    // Старый надежный Fallback-путь
}
```

---

## 25. Паттерны Отказоустойчивости в Микросервисах (Resilience & Stability)

### 25.1. `Circuit Breaker` (Предохранитель)

> [!NOTE]
> **Зачем это нужно на проде?**  
> Если внешний платежный шлюз или смежный микросервис «завис» или отвечает 500 ошибками, отправка тысяч новых запросов окончательно добьет сервис и исчерпает стек горутин вашего собственного приложения.  
> **Circuit Breaker** «размыкает цепь» при всплеске ошибок, мгновенно возвращая ошибку клиенту без выполнения тяжелых сетевых вызовов (Fast Fail).

```mermaid
stateDiagram-v2
    [*] --> Closed
    Closed --> Open : Процент ошибок > 50%
    Open --> HalfOpen : Таймаут ожидания (10 сек)
    HalfOpen --> Closed : Тестовые запросы успешны
    HalfOpen --> Open : Тестовый запрос упал
```

- **Closed (Замкнут):** Нормальная работа, запросы идут к сервису.
- **Open (Разомкнут):** Сервис недоступен. Запросы **не отправляются**, клиенту сразу отдается ошибка или Fallback-ответ (кэш).
- **Half-Open (Полуразомкнут):** Пропускает 1–2 тестовых запроса для проверки, ожил ли зависимый сервис.

Популярная библиотека в Go: `github.com/sony/gobreaker`. Пример использования (версия v2 с дженериками):

```go
import "github.com/sony/gobreaker/v2"

var cb = gobreaker.NewCircuitBreaker[[]byte](gobreaker.Settings{
    Name:    "payments",
    Timeout: 10 * time.Second, // сколько цепь остаётся разомкнутой, прежде чем перейти в Half-Open
    ReadyToTrip: func(c gobreaker.Counts) bool { // условие размыкания
        return c.Requests >= 10 && float64(c.TotalFailures)/float64(c.Requests) > 0.5
    },
})

func callPayments(ctx context.Context) ([]byte, error) {
    // Execute сам считает успехи/ошибки и возвращает gobreaker.ErrOpenState, когда цепь разомкнута
    return cb.Execute(func() ([]byte, error) {
        return doHTTPCall(ctx) // реальный сетевой вызов с таймаутом
    })
}
```

---

### 25.2. `Rate Limiting` и `Throttling` (Ограничение частоты запросов)

Защищает сервис от DDoS-атак, перегрузок и контролирует лимиты бесплатных/платных API тарифов.

1. **Token Bucket (Корзина токенов — `golang.org/x/time/rate`):**
   - В «корзину» с фиксированной скоростью капают токены.
   - Каждый запрос забирает 1 токен.
   - Позволяет кратковременные всплески трафика (_Burst_), пока есть накопленные токены.
2. **Leaky Bucket (Протекающее ведро — `go.uber.org/ratelimit`):**
   - Запросы обрабатываются со строго постоянной скоростью, сглаживая любые всплески в ровный поток.

```go
import "golang.org/x/time/rate"

// Лимит: 100 запросов в секунду, с возможностью всплеска (burst) до 200:
var limiter = rate.NewLimiter(rate.Limit(100), 200)

func handleRequest(w http.ResponseWriter, r *http.Request) {
    if !limiter.Allow() {
        w.Header().Set("Retry-After", "1") // подсказка клиенту, через сколько секунд повторить
        http.Error(w, "Too Many Requests (429)", http.StatusTooManyRequests)
        return
    }
    // Обработка запроса...
}
```

`limiter.Allow()` отвечает мгновенно (пропустить или отказать); если запрос нужно не отклонить, а **подождать** своей очереди, используйте `limiter.Wait(ctx)`. Такой лимитер хранит состояние в памяти одного экземпляра сервиса — для лимита на весь кластер нужен общий счётчик (например, в Redis).

---

### 25.3. `Retry + Exponential Backoff + Jitter`

> [!WARNING]
> **Почему нельзя просто повторять упавший запрос в цикле `for`?**  
> Если база данных кратковременно перегрузилась, 10 000 клиентов, одновременно и непрерывно повторяющих запрос, устроят **«Шторм повторных попыток» (Retry Storm)** и никогда не дадут базе восстановиться.

#### Правильный алгоритм:

1. **Exponential Backoff:** Время между попытками растет экспоненциально: $100\text{ms} \rightarrow 200\text{ms} \rightarrow 400\text{ms} \rightarrow 800\text{ms}$.
2. **Jitter (Случайный шум):** Добавление случайной дельты к задержке ($+\text{rand}(0, 50)\text{ms}$), чтобы размазать нагрузку от параллельных клиентов по времени.

```go
import (
    "context"
    "math/rand"
    "time"
)

func retryWithBackoff(ctx context.Context, maxAttempts int, op func() error) error {
    var err error
    delay := 100 * time.Millisecond

    for attempt := 1; attempt <= maxAttempts; attempt++ {
        if err = op(); err == nil {
            return nil // Успех!
        }

        if attempt == maxAttempts {
            break
        }

        // Вычисляем задержку с джиттером (Jitter):
        jitter := time.Duration(rand.Int63n(int64(delay / 2)))
        sleepTime := delay + jitter

        select {
        case <-time.After(sleepTime):
            delay *= 2 // экспоненциальный рост
        case <-ctx.Done():
            return ctx.Err()
        }
    }
    return err
}
```

> [!WARNING]
> - Повторять можно только **идемпотентные** операции (чтение, `PUT`, операции с ключом идемпотентности): повтор `POST /payments` без такого ключа может списать деньги дважды.
> - Повторяют только **временные** ошибки (таймаут, `503`, `429`), но не `400`/`404`/ошибки валидации.
> - Ограничивайте общее число попыток и суммарное время (через `ctx`); на практике часто используют готовые библиотеки (`cenkalti/backoff`).

---

> [!TIP]
> 📚 **Навигация:** [⬅️ Назад: Том 4](04-go-core-backend-web.md) | [📖 Главное оглавление](README.md) | [Вперед: Том 6 (Алгоритмы и Собеседования) ➡️](06-go-core-algorithms-interview.md)
