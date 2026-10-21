# 🏭 Руководство по промышленной разработке на Go (Go Production Engineering Guide)

> [!IMPORTANT]
> **От учебного кода — к боевому продакшену.**
> На собеседованиях и в академических курсах проверяют знание структур данных, алгоритмов и рантайма Go (слайсы, планировщик GMP, каналы). Но в **ежедневной коммерческой разработке** 80% времени уходит на надежность, инфраструктурную обвязку, архитектурные контракты и наблюдаемость сервисов.
>
> Этот гайд содержит проверенные практикой шаблоны и рецепты production-grade кода, обязательные для любого Enterprise/Fintech Go-бэкенда.

---

## 📑 Оглавление

1. [12-Factor Конфигурация и Fail-Fast старт сервиса](#1-12-factor-конфигурация-и-fail-fast-старт-сервиса)
2. [Graceful Shutdown и управление жизненным циклом (Lifecycle)](#2-graceful-shutdown-и-управление-жизненным-циклом-lifecycle)
3. [Базы данных в реальном мире: sqlc, миграции и транзакции](#3-базы-данных-в-реальном-мире-sqlc-миграции-и-транзакции)
4. [Observability: Метрики (Prometheus), Трейсинг (OpenTelemetry) и Логи (slog)](#4-observability-метрики-prometheus-трейсинг-opentelemetry-и-логи-slog)
5. [Отказоустойчивый HTTP-клиент: Таймауты, Backoff-Retry и Circuit Breaker](#5-отказоустойчивый-http-клиент-таймауты-backoff-retry-и-circuit-breaker)
6. [Валидация входящих данных (DTO Validation) и ошибки RFC 7807](#6-валидация-входящих-данных-dto-validation-и-ошибки-rfc-7807)
7. [Аутентификация: JWT (Access/Refresh), Context Injection и Pipeline Middleware](#7-аутентификация-jwt-accessrefresh-context-injection-и-pipeline-middleware)
8. [Production Tooling: golangci-lint, Makefile и Multi-Stage Dockerfile](#8-production-tooling-golangci-lint-makefile-и-multi-stage-dockerfile)
9. [Боевой CI/CD пайплайн на GitHub Actions (Lint, Test, Vulncheck, Docker Buildx)](#9-боевой-cicd-пайплайн-на-github-actions)
10. [Продвинутая работа с Apache Kafka: Transactional Outbox, Consumer Groups и DLQ](#10-продвинутая-работа-с-apache-kafka-transactional-outbox-consumer-groups-и-dlq)
11. [gRPC Interceptors: Цепочки Middleware для Enterprise-микросервисов](#11-grpc-interceptors-цепочки-middleware-для-enterprise-микросервисов)
12. [Безопасность приложений на Go (AppSec): Защита от OOM, SSRF, Timing Attacks](#12-безопасность-приложений-на-go-appsec-защита-от-oom-ssrf-timing-attacks)

---

## 1. 12-Factor Конфигурация и Fail-Fast старт сервиса

### Проблема «наивного» кода
В учебных проектах адреса БД и порты часто объявляют константами в коде или читают разрозненными вызовами `os.Getenv("PORT")`. В production:
1. Значения приходят из переменных окружения (Kubernetes ConfigMaps/Secrets, Docker Compose, `.env`).
2. Опечатка в названии переменной или отсутствие критического параметра (например, `JWT_SECRET`) приводит к падению сервиса спустя часы работы, когда пользователь нажимает на определенную ручку.

### Паттерн Fail-Fast и библиотека `cleanenv`
Сервис должен **упасть в первую же секунду запуска** с понятным описанием проблемы, если обязательные параметры отсутствуют или невалидны.

```go
package config

import (
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App  AppConfig  `yaml:"app"`
	HTTP HTTPConfig `yaml:"http"`
	DB   DBConfig   `yaml:"db"`
}

type AppConfig struct {
	Env     string `yaml:"env" env:"APP_ENV" env-default:"development"` // development, staging, production
	Name    string `yaml:"name" env:"APP_NAME" env-default:"order-service"`
	Version string `yaml:"version" env:"APP_VERSION" env-default:"1.0.0"`
}

type HTTPConfig struct {
	Port         int           `yaml:"port" env:"HTTP_PORT" env-default:"8080"`
	ReadTimeout  time.Duration `yaml:"read_timeout" env:"HTTP_READ_TIMEOUT" env-default:"5s"`
	WriteTimeout time.Duration `yaml:"write_timeout" env:"HTTP_WRITE_TIMEOUT" env-default:"10s"`
	IdleTimeout  time.Duration `yaml:"idle_timeout" env:"HTTP_IDLE_TIMEOUT" env-default:"60s"`
}

type DBConfig struct {
	URL          string        `yaml:"url" env:"DATABASE_URL" env-required:"true"`
	MaxConns     int32         `yaml:"max_conns" env:"DB_MAX_CONNS" env-default:"25"`
	MinConns     int32         `yaml:"min_conns" env:"DB_MIN_CONNS" env-default:"5"`
	MaxIdleTime  time.Duration `yaml:"max_idle_time" env:"DB_MAX_IDLE_TIME" env-default:"15m"`
	ConnTimeout  time.Duration `yaml:"conn_timeout" env:"DB_CONN_TIMEOUT" env-default:"5s"`
}

// MustLoad загружает конфигурацию или паникует/завершает процесс при ошибке
func MustLoad(configPath string) *Config {
	var cfg Config

	// 1. Если передан путь к yaml/json файлу (локальная разработка)
	if configPath != "" {
		if _, err := os.Stat(configPath); err == nil {
			if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
				panic(fmt.Sprintf("failed to read config file %s: %v", configPath, err))
			}
		}
	}

	// 2. Читаем переменные окружения, переопределяющие файл
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic(fmt.Sprintf("failed to parse environment variables: %v", err))
	}

	return &cfg
}
```

> [!TIP]
> Функция называется `MustLoad`, подчеркивая идиому Go: функции с префиксом `Must` вызываются только на старте приложения и паникуют при невозможности инициализации.

---

## 2. Graceful Shutdown и управление жизненным циклом (Lifecycle)

### Зачем нужен Graceful Shutdown?
В Kubernetes поды постоянно завершаются (rolling updates, pod eviction, autoscaling).
Когда Kubelet решает завершить контейнер:
1. Он отправляет процессу сигнал **`SIGTERM`**.
2. Ждет время `terminationGracePeriodSeconds` (по умолчанию 30 секунд).
3. Если процесс не завершился — принудительно убивает его сигналом **`SIGKILL` (`kill -9`)**.

Если ваш сервис не обрабатывает `SIGTERM`, все текущие HTTP-запросы пользователей обрываются, сокеты закрываются с ошибкой `connection reset by peer`, а клиенты видят `502 Bad Gateway`.

### Каскадный порядок остановки сервиса:
```mermaid
graph TD
    Signal["Получен сигнал SIGTERM / SIGINT"] --> Step1["1. Остановка приема входящих запросов<br/>(httpServer.Shutdown)"]
    Step1 --> Step2["2. Завершение обработки текущих запросов<br/>и фоновых воркеров (wg.Wait)"]
    Step2 --> Step3["3. Остановка консьюмеров очередей<br/>(Kafka / RabbitMQ close)"]
    Step3 --> Step4["4. Закрытие соединений с БД и кэшем<br/>(dbPool.Close)"]
    Step4 --> Exit["Успешный выход с кодом 0"]
```

### Эталонный шаблон `main.go`:

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	// 1. Инициализация структурированного логера
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	// 2. Создаем контекст, который отменяется при поступлении сигналов SIGINT или SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 3. Подключение к БД (псевдокод)
	// dbPool := initDatabase(...)
	// defer dbPool.Close()

	// 4. Настройка HTTP-сервера
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("POST /api/v1/orders", func(w http.ResponseWriter, r *http.Request) {
		// Имитация долгой полезной нагрузки
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"status":"created"}`))
	})

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 5. Запуск фоновых воркеров с WaitGroup
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		backgroundWorker(ctx, logger)
	}()

	// 6. Запуск HTTP-сервера в отдельной горутине
	go func() {
		logger.Info("HTTP server starting", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("HTTP server failed", "error", err)
			os.Exit(1)
		}
	}()

	// 7. Ожидание сигнала завершения
	<-ctx.Done()
	logger.Info("Shutdown signal received, initiating graceful shutdown...")

	// 8. Контекст с таймаутом на безопасную остановку (например, 15 секунд)
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	// Останавливаем HTTP-сервер: он закрывает слушающий сокет (новые запросы не приходят),
	// но дожидается завершения уже обрабатывающихся запросов
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("HTTP server forced shutdown", "error", err)
	} else {
		logger.Info("HTTP server stopped gracefully")
	}

	// Ждем остановки фоновых воркеров
	workerChan := make(chan struct{})
	go func() {
		wg.Wait()
		close(workerChan)
	}()

	select {
	case <-workerChan:
		logger.Info("All background workers finished")
	case <-shutdownCtx.Done():
		logger.Warn("Background workers shutdown timed out")
	}

	logger.Info("Service exited cleanly")
}

func backgroundWorker(ctx context.Context, logger *slog.Logger) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("Background worker received cancellation, cleaning up...")
			return
		case <-ticker.C:
			// фоновая работа
		}
	}
}
```

---

## 3. Базы данных в реальном мире: sqlc, миграции и транзакции

### Проблема: `database/sql` vs Тяжелые ORM (GORM)
- **`database/sql` с ручным сканированием:** при изменении таблицы в БД разработчик забывает обновить поле в `rows.Scan(&a, &b, ...)`, что приводит к падениям в рантайме.
- **ORM (GORM / Ent):** прячут реальный SQL. В простых запросах удобны, но на высоконагруженных проектах приводят к неоптимальным JOIN, скрытым аллокациям рефлексии и проблеме $N+1$.

### Стандарт индустрии: `sqlc` (Compile-time Type-Safe SQL)
Инструмент **`sqlc`** компилирует чистый SQL в строгий, идиоматичный Go-код без использования рефлексии.

#### 1. Конфиг `sqlc.yaml`:
```yaml
version: "2"
sql:
  - engine: "postgresql"
    schema: "migrations/schema.sql"
    queries: "queries.sql"
    gen:
      go:
        package: "db"
        out: "internal/db"
        sql_package: "pgx/v5"
        emit_json_tags: true
        emit_empty_slices: true
```

#### 2. SQL-запросы (`queries.sql`):
```sql
-- name: GetUserByID :one
SELECT id, email, full_name, created_at
FROM users
WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT id, email, full_name, created_at
FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CreateUser :one
INSERT INTO users (email, full_name)
VALUES ($1, $2)
RETURNING id, email, full_name, created_at;
```

`sqlc generate` создаст Go-структуру `User`, методы `GetUserByID(ctx, id)` и `CreateUser(ctx, params)` с полной проверкой типов аргументов и возвращаемых значений на этапе сборки.

---

### Миграции базы данных (`goose`)
Схема базы данных должна версионироваться вместе с кодом. Самый популярный инструмент в Go-сообществе — **`goose`** (или `golang-migrate`).

Файлы миграций хранятся в `migrations/`:
```sql
-- 20261021100000_create_users_table.sql
-- +goose Up
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    full_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);

-- +goose Down
DROP TABLE IF EXISTS users;
```

---

### Управление транзакциями в Clean Architecture
Вечный вопрос: *как сделать транзакцию, затрагивающую несколько репозиториев (например, списать деньги в `AccountRepo` и создать заказ в `OrderRepo`), не протаскивая `*sql.Tx` в сервисный слой?*

**Идиоматичный паттерн: Репозиторий с замыканием (`WithTx`):**

```go
type Transactor interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type PostgresTransactor struct {
	pool *pgxpool.Pool
}

type txKey struct{}

func (t *PostgresTransactor) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := t.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx) // Игнорируется, если выполнен Commit()

	// Вкладываем транзакцию в контекст
	txCtx := context.WithValue(ctx, txKey{}, tx)

	if err := fn(txCtx); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

// В репозитории извлекаем транзакцию из контекста или используем обычный пул:
func (r *OrderRepository) getExecutor(ctx context.Context) DBExecutor {
	if tx, ok := ctx.Value(txKey{}).(pgx.Tx); ok {
		return tx
	}
	return r.pool
}
```

---

## 4. Observability: Метрики (Prometheus), Трейсинг (OpenTelemetry) и Логи (slog)

### Три столпа наблюдаемости (Observability)
```mermaid
graph LR
    Metrics["Метрики (Prometheus)<br/>Симптомы: что и где болит? (RPS, 5xx, p99 Latency)"] --- Traces["Трейсы (OpenTelemetry)<br/>Локализация: в каком микросервисе задержка?"]
    Traces --- Logs["Логи (slog)<br/>Причина: почему произошла ошибка?"]
```

### 1. Метрики Prometheus (`client_golang`)
Для продакшена критически важны 4 золотых сигнала (Google SRE): **Latency, Traffic, Errors, Saturation**.

```go
package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Количество запросов с разбивкой по методу, пути и коду ответа
	HttpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests processed",
		},
		[]string{"method", "path", "status"},
	)

	// Гистограмма времени выполнения запросов (для вычисления p50, p95, p99)
	HttpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latencies in seconds",
			Buckets: []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5},
		},
		[]string{"method", "path"},
	)
)

// Middleware для автоматического сбора метрик
func PrometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(ww, r)

		duration := time.Since(start).Seconds()
		path := r.Pattern // Go 1.22+ сохраняет паттерн маршрута, предотвращая кардинальность меток

		HttpRequestsTotal.WithLabelValues(r.Method, path, strconv.Itoa(ww.statusCode)).Inc()
		HttpRequestDuration.WithLabelValues(r.Method, path).Observe(duration)
	})
}

type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *statusResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}
```

### 2. Корреляция Трейсов (OpenTelemetry) и Логов (`slog`)
Когда в логах появляется ошибка, разработчик должен иметь возможность скопировать `trace_id` и открыть полный путь запроса в Jaeger/Grafana Tempo.

```go
// slogHandler обогащает каждую запись лога текущим trace_id из контекста
type TraceHandler struct {
	slog.Handler
}

func (h *TraceHandler) Handle(ctx context.Context, r slog.Record) error {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		r.AddAttrs(
			slog.String("trace_id", span.SpanContext().TraceID().String()),
			slog.String("span_id", span.SpanContext().SpanID().String()),
		)
	}
	return h.Handler.Handle(ctx, r)
}
```

---

## 5. Отказоустойчивый HTTP-клиент: Таймауты, Backoff-Retry и Circuit Breaker

### Почему `http.DefaultClient` запрещен в проде?
В стандартном `http.DefaultClient` таймаут равен **0 (бесконечность)**. Если внешний сервис зависнет, горутина вызова заблокируется навсегда. При потоке в 100 запросов в секунду за минуту утекут тысячи горутин и сокетов, что приведет к OOM (Out Of Memory).

### Настройка эталонного `http.Client`:
```go
func NewResilientHTTPClient() *http.Client {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,  // Таймаут TCP handshake
			KeepAlive: 30 * time.Second, // Поддержание TCP соединения
		}).DialContext,
		MaxIdleConns:          100,              // Всего свободных соединений в пуле
		MaxIdleConnsPerHost:   20,               // Свободных соединений на каждый хост
		IdleConnTimeout:       90 * time.Second, // Время жизни простаивающего соединения
		TLSHandshakeTimeout:   5 * time.Second,  // Таймаут TLS-рукопожатия
		ExpectContinueTimeout: 1 * time.Second,
	}

	return &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second, // Общий таймаут всего запроса (включая чтение тела)
	}
}
```

### Retry с Exponential Backoff и Full Jitter
Если внешний сервис временно перегружен, нельзя повторять запрос сразу (это добьет его). Нужно увеличивать интервал экспоненциально и добавлять случайный разброс (jitter), чтобы клиенты не стреляли синхронно (Thundering Herd Problem).

```mermaid
graph LR
    Att1["Попытка 1: Ошибка 503"] --> Delay1["Пауза ~100ms ± random"]
    Delay1 --> Att2["Попытка 2: Ошибка 503"]
    Att2 --> Delay2["Пауза ~400ms ± random"]
    Delay2 --> Att3["Попытка 3: Успех 200"]
```

```go
// DoWithRetry выполняет операцию с экспоненциальной задержкой и джиттером
func DoWithRetry(ctx context.Context, maxRetries int, baseDelay time.Duration, op func() error) error {
	var err error
	for attempt := 0; attempt < maxRetries; attempt++ {
		err = op()
		if err == nil {
			return nil
		}

		// Экспоненциальный рост: baseDelay * 2^attempt
		backoff := float64(baseDelay * (1 << attempt))
		// Добавляем Full Jitter: случайное число от 0 до backoff
		jitter := rand.Float64() * backoff
		sleepDuration := time.Duration(jitter)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(sleepDuration):
		}
	}
	return fmt.Errorf("all %d attempts failed: %w", maxRetries, err)
}
```

---

## 6. Валидация входящих данных (DTO Validation) и ошибки RFC 7807

### Валидация через `go-playground/validator/v10`

```go
package dto

type CreateOrderRequest struct {
	CustomerID string  `json:"customer_id" validate:"required,uuid4"`
	UserEmail  string  `json:"email" validate:"required,email"`
	Amount     float64 `json:"amount" validate:"required,gt=0"`
	Currency   string  `json:"currency" validate:"required,oneof=USD EUR RUB KZT"`
	ItemsCount int     `json:"items_count" validate:"required,min=1,max=100"`
}
```

### Формат ошибок: RFC 7807 (Problem Details)
Клиенты API (фронтенд, мобильные приложения) ожидают структурированный ответ при ошибках валидации (HTTP 422 / 400).

```json
{
  "type": "https://api.example.com/errors/validation-failed",
  "title": "Validation Failed",
  "status": 422,
  "detail": "One or more fields failed validation checks",
  "invalid_params": [
    {
      "field": "email",
      "reason": "must be a valid email address"
    },
    {
      "field": "amount",
      "reason": "must be greater than 0"
    }
  ]
}
```

```go
type ProblemDetails struct {
	Type          string          `json:"type"`
	Title         string          `json:"title"`
	Status        int             `json:"status"`
	Detail        string          `json:"detail"`
	InvalidParams []InvalidParam  `json:"invalid_params,omitempty"`
}

type InvalidParam struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

func FormatValidationError(err error) ProblemDetails {
	var invalidParams []InvalidParam
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range validationErrors {
			invalidParams = append(invalidParams, InvalidParam{
				Field:  fe.Field(),
				Reason: fmt.Sprintf("failed check on tag '%s'", fe.Tag()),
			})
		}
	}

	return ProblemDetails{
		Type:          "https://example.com/errors/validation-failed",
		Title:         "Validation Failed",
		Status:        http.StatusUnprocessableEntity,
		Detail:        "Invalid input parameters",
		InvalidParams: invalidParams,
	}
}
```

---

## 7. Аутентификация: JWT (Access/Refresh), Context Injection и Pipeline Middleware

### Конвейер Middleware (Строгий порядок выполнения)
Порядок объявления middleware в коде критичен:
```mermaid
graph TD
    Req["Входящий запрос"] --> M1["1. Recovery (защита от паник)"]
    M1 --> M2["2. RequestID & TraceID (генерация сквозного ID)"]
    M2 --> M3["3. Structured Logger (логирование старта/завершения)"]
    M3 --> M4["4. Metrics (Prometheus счетчики)"]
    M4 --> M5["5. CORS & Security Headers"]
    M5 --> M6["6. Rate Limiter (ограничение спама)"]
    M6 --> M7["7. Auth / JWT (проверка токена и проброс пользователя)"]
    M7 --> Handler["Целевой хендлер (Бизнес-логика)"]
```

### Безопасное внедрение пользователя в `context.Context`

> [!WARNING]
> Никогда не используйте встроенные типы (`string`, `int`) в качестве ключей `context.WithValue`! Это приводит к коллизиям между сторонними пакетами. Всегда объявляйте неэкспортируемый приватный тип `type ctxKey struct{}`.

```go
package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type userContextKey struct{}

type UserClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// AuthMiddleware проверяет валидность Bearer JWT и инжектит пользователя в контекст
func AuthMiddleware(secretKey []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			claims := &UserClaims{}

			token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("unexpected signing method")
				}
				return secretKey, nil
			})

			if err != nil || !token.Valid {
				http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
				return
			}

			// Инжектим данные в контекст
			ctx := context.WithValue(r.Context(), userContextKey{}, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// UserFromContext безопасно извлекает утверждения пользователя из контекста
func UserFromContext(ctx context.Context) (*UserClaims, bool) {
	claims, ok := ctx.Value(userContextKey{}).(*UserClaims)
	return claims, ok
}
```

---

## 8. Production Tooling: golangci-lint, Makefile и Multi-Stage Dockerfile

### 1. Конфигурация `.golangci.yml`
В промышленной разработке линтеры запускаются на каждый Pull Request. Базовый набор защищает от утечек памяти, пропущенных ошибок и небезопасного кода:

```yaml
run:
  timeout: 5m
  tests: true

linters:
  enable:
    - errcheck      # Проверяет, что ошибки не игнорируются
    - govet         # Официальный анализатор Go
    - staticcheck   # Мощнейший набор статических проверок
    - bodyclose     # Ловит утечки незакрытых http.Response.Body!
    - noctx         # Требует передачи context в HTTP-запросы
    - gocritic      # Советы по стилю и производительности
    - revive        # Быстрая замена устаревшего golint
    - gosec         # Поиск уязвимостей (SQL injection, слабые хэши)
```

---

### 2. Универсальный `Makefile` для Go-проекта:

```makefile
.PHONY: run build test lint migrate-up migrate-down sqlc docker-build

APP_NAME = order-service
BIN_DIR = bin

run:
	go run cmd/app/main.go

build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -o $(BIN_DIR)/$(APP_NAME) cmd/app/main.go

test:
	go test -race -v -coverprofile=coverage.out ./...

lint:
	golangci-lint run ./...

sqlc:
	sqlc generate

migrate-up:
	goose -dir migrations postgres "$$DATABASE_URL" up

migrate-down:
	goose -dir migrations postgres "$$DATABASE_URL" down

docker-build:
	docker build -t $(APP_NAME):latest .
```

---

### 3. Production Multi-Stage `Dockerfile` (Размер образа: ~15 МБ)

```dockerfile
# ---------------------------------------------------------------------
# ЭТАП 1: Сборка бинарника (Build Stage)
# ---------------------------------------------------------------------
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Устанавливаем сертификаты и зависимости для сборки
RUN apk add --no-cache git ca-certificates tzdata

# Сначала копируем только файлы зависимостей (для эффективного кэширования слоев Docker)
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Копируем исходный код
COPY . .

# Собираем минималистичный статический бинарник:
# - CGO_ENABLED=0 (без зависимостей от системного libc)
# - -ldflags="-s -w" (удаление отладочной таблицы символов DWARF для снижения размера на 40%)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w" \
    -o /bin/app ./cmd/app/main.go

# ---------------------------------------------------------------------
# ЭТАП 2: Финальный минимальный контейнер (Runtime Stage)
# ---------------------------------------------------------------------
FROM alpine:3.21

# Создаем непривилегированного пользователя appuser (защита от container escape)
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Копируем корневые сертификаты для исходящих HTTPS-запросов и таймзоны
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /bin/app /bin/app

USER appuser

EXPOSE 8080

ENTRYPOINT ["/bin/app"]
```

> [!TIP]
> **Почему этот Dockerfile эталонный:**
> 1. **Безопасность:** Запуск от `appuser` (не root).
> 2. **Размер:** Всего ~18 МБ вместо ~1 ГБ стандартного Go-образа.
> 3. **Быстрый CI:** Слой `go mod download` кэшируется и не перезакачивается, пока не изменились зависимости.
> 4. **Корректный TLS:** В образ скопированы `ca-certificates.crt`, поэтому сервис сможет ходить в сторонние HTTPS API без ошибок `x509: certificate signed by unknown authority`.

---

## 9. Боевой CI/CD пайплайн на GitHub Actions

### Архитектура современного Go CI/CD
В современных высоконагруженных командах CI/CD пайплайн выполняет четыре критических барьера качества (Quality Gates) перед мержем в основную ветку:
1. **Fast-Fail Линтинг (`golangci-lint`):** проверка кодстайла, утечек памяти, необработанных ошибок (`errcheck`), невыровненных структур (`fieldalignment`) и подозрительных конструкций (`govet`, `staticcheck`).
2. **Тестирование с детектором гонок (`-race`) и покрытием:** обязательный запуск unit- и integration-тестов под флагом `-race`. В проде флаг `-race` не используется из-за замедления в 2–5 раз и роста потребления памяти, но в CI он предотвращает 99% плавающих дедлоков и дата-рейсов.
3. **Сканирование уязвимостей (`govulncheck`):** официальный анализатор Go Vulnerability Database от команды Go. В отличие от общих CVE-сканеров, `govulncheck` строит граф вызовов и алертит **только тогда, когда ваш код реально вызывает уязвимую функцию из библиотеки**, полностью устраняя ложные срабатывания (False Positives).
4. **Сборка и публикация Multi-Arch Docker-образа:** сборка через Docker Buildx с кэшированием слоев в GitHub Actions Cache (`type=gha`), что сокращает время сборки с 3 минут до 15 секунд.

### Эталонный рабочий файл `.github/workflows/ci.yml`:
```yaml
name: CI/CD Pipeline

on:
  push:
    branches: [ master, main ]
  pull_request:
    branches: [ master, main ]

# Автоматически отменяем предыдущие незавершенные запуски для того же PR:
concurrency:
  group: ${{ github.workflow }}-${{ github.ref }}
  cancel-in-progress: true

env:
  GO_VERSION: '1.23'
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}

jobs:
  # ---------------------------------------------------------------------------
  # 1. ЛИНТИНГ И СТАТИЧЕСКИЙ АНАЛИЗ
  # ---------------------------------------------------------------------------
  lint:
    name: Golangci-lint
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
          cache: false # Кэшированием управляет golangci-lint action

      - name: Run golangci-lint
        uses: golangci/golangci-lint-action@v6
        with:
          version: v1.61.0
          args: --timeout=5m --verbose

  # ---------------------------------------------------------------------------
  # 2. ТЕСТИРОВАНИЕ С ДЕТЕКТОРОМ ГОНОК И ПОКРЫТИЕМ
  # ---------------------------------------------------------------------------
  test:
    name: Tests & Race Detector
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
          cache: true

      - name: Download Go modules
        run: go mod download && go mod verify

      - name: Run tests with -race
        run: |
          go test -v -race -covermode=atomic -coverprofile=coverage.out ./...

      - name: Check minimum coverage (threshold: 70%)
        run: |
          TOTAL_COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print substr($3, 1, length($3)-1)}')
          echo "Total coverage: ${TOTAL_COVERAGE}%"
          if (( $(echo "$TOTAL_COVERAGE < 70.0" | bc -l) )); then
            echo "❌ Coverage is below required threshold of 70%!"
            exit 1
          fi

  # ---------------------------------------------------------------------------
  # 3. АУДИТ УЯЗВИМОСТЕЙ (GOVULNCHECK)
  # ---------------------------------------------------------------------------
  security:
    name: Security Vulnerability Scan
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
          cache: true

      - name: Run govulncheck
        uses: golang/govulncheck-action@v1
        with:
          go-version-input: ${{ env.GO_VERSION }}
          go-package: ./...

  # ---------------------------------------------------------------------------
  # 4. СБОРКА И КЭШИРОВАНИЕ DOCKER-ОБРАЗА ЧЕРЕЗ BUILDX
  # ---------------------------------------------------------------------------
  build-docker:
    name: Build & Push Docker Image
    needs: [lint, test, security]
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Log in to GitHub Container Registry
        if: github.event_name != 'pull_request'
        uses: docker/login-action@v3
        with:
          registry: ${{ env.REGISTRY }}
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: Extract Docker metadata
        id: meta
        uses: docker/metadata-action@v5
        with:
          images: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}
          tags: |
            type=ref,event=branch
            type=ref,event=pr
            type=semver,pattern={{version}}
            type=sha,format=short

      - name: Build and Push with GHA Caching
        uses: docker/build-push-action@v6
        with:
          context: .
          push: ${{ github.event_name != 'pull_request' }}
          tags: ${{ steps.meta.outputs.tags }}
          labels: ${{ steps.meta.outputs.labels }}
          cache-from: type=gha
          cache-to: type=gha,mode=max
```

---

## 10. Продвинутая работа с Apache Kafka: Transactional Outbox, Consumer Groups и DLQ

### Проблема двойной записи (Dual-Write Problem)
Типичная архитектурная катастрофа в микросервисах:
1. Сервис оформляет платеж: сохраняет заказ в PostgreSQL со статусом `PAID`.
2. Сервис отправляет событие `OrderPaidEvent` в Kafka для складского сервиса.
3. **Что если сеть до Kafka моргнула или брокер упал?**  
   Заказ оплачен в базе, деньги списаны, но склад не получил уведомление — заказ завис навсегда.
4. **Что если сначала писать в Kafka, а потом в Postgres?**  
   Если база вернет ошибку (deadlock, уникальный ключ), сообщение в Kafka уже ушло — склад отправит товар, за который не списаны деньги.

---

### Паттерн Transactional Outbox (Золотой стандарт надежности)
Решение: **сохранять доменную сущность и исходящее событие в одной локальной ACID-транзакции PostgreSQL**.

```
[HTTP Request: Оплата]
         │
         ▼
 ┌───────────────────────────────────────────────┐
 │ PostgreSQL Local ACID Transaction             │
 │  1. INSERT INTO orders (...) VALUES (...)     │
 │  2. INSERT INTO outbox_events (topic, payload)│
 └───────────────────────────────────────────────┘
         │ (Гарантия атомарности 100%)
         ▼
 ┌───────────────────────────────────────────────┐
 │ Background Outbox Worker (Go Goroutine)       │
 │  SELECT FOR UPDATE SKIP LOCKED ...            │
 │  kafkaProducer.Publish(event)                 │
 │  UPDATE outbox_events SET status = 'PUBLISHED'│
 └───────────────────────────────────────────────┘
         │
         ▼
     [Kafka Broker]
```

#### DDL таблицы Outbox:
```sql
CREATE TABLE outbox_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    aggregate_type VARCHAR(64) NOT NULL,
    aggregate_id VARCHAR(64) NOT NULL,
    event_type VARCHAR(64) NOT NULL,
    topic VARCHAR(128) NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    status VARCHAR(32) NOT NULL DEFAULT 'PENDING',
    retry_count INT NOT NULL DEFAULT 0,
    last_error TEXT
);

CREATE INDEX idx_outbox_pending ON outbox_events (created_at) WHERE status = 'PENDING';
```

#### Реализация воркера Outbox на Go:
```go
package outbox

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/segmentio/kafka-go"
)

type Event struct {
	ID      string          `db:"id"`
	Topic   string          `db:"topic"`
	Payload json.RawMessage `db:"payload"`
}

type OutboxProcessor struct {
	pool   *pgxpool.Pool
	writer *kafka.Writer
	logger *slog.Logger
}

func (p *OutboxProcessor) Start(ctx context.Context) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := p.processBatch(ctx); err != nil {
				p.logger.Error("failed to process outbox batch", "error", err)
			}
		}
	}
}

func (p *OutboxProcessor) processBatch(ctx context.Context) error {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Блокируем до 100 записей через SKIP LOCKED, чтобы параллельные поды не конфликтовали:
	rows, err := tx.Query(ctx, `
		SELECT id, topic, payload 
		FROM outbox_events 
		WHERE status = 'PENDING' 
		ORDER BY created_at ASC 
		LIMIT 100 
		FOR UPDATE SKIP LOCKED
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var e Event
		if err := rows.Scan(&e.ID, &e.Topic, &e.Payload); err != nil {
			return err
		}
		events = append(events, e)
	}

	if len(events) == 0 {
		return nil
	}

	// Отправляем сообщения в Kafka батчем:
	kafkaMessages := make([]kafka.Message, len(events))
	eventIDs := make([]string, len(events))
	for i, e := range events {
		kafkaMessages[i] = kafka.Message{
			Topic: e.Topic,
			Key:   []byte(e.ID),
			Value: e.Payload,
			Time:  time.Now(),
		}
		eventIDs[i] = e.ID
	}

	if err := p.writer.WriteMessages(ctx, kafkaMessages...); err != nil {
		return fmt.Errorf("kafka write failed: %w", err)
	}

	// Помечаем отправленные события:
	_, err = tx.Exec(ctx, `
		UPDATE outbox_events 
		SET status = 'PUBLISHED' 
		WHERE id = ANY($1)
	`, eventIDs)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
```

---

### Consumer Groups, ручной коммит и Graceful Rebalance
При автокоммите (`commit.interval`) оффсет сдвигается в фоне по таймеру. Если под умер во время обработки сложного сообщения, следующее сообщение уже скоммичено $\to$ **потеря данных**!

#### ✅ Эталонный цикл консьюмера с ручным коммитом:
```go
package consumer

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/segmentio/kafka-go"
)

type OrderConsumer struct {
	reader *kafka.Reader
	logger *slog.Logger
}

func (c *OrderConsumer) Run(ctx context.Context) error {
	for {
		// 1. FetchMessage забирает сообщение, но НЕ коммитит оффсет:
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return nil // Штатный shutdown
			}
			c.logger.Error("failed to fetch message", "error", err)
			time.Sleep(1 * time.Second)
			continue
		}

		// 2. Бизнес-обработка сообщения:
		if err := c.handleMessage(ctx, msg); err != nil {
			c.logger.Error("business handling failed", "key", string(msg.Key), "error", err)
			// При критической ошибке НЕ коммитим — сообщение повторится или уйдет в DLQ
			continue
		}

		// 3. Ручной коммит оффсета ТОЛЬКО после успешной бизнес-обработки:
		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			c.logger.Error("failed to commit offset", "error", err)
		}
	}
}

func (c *OrderConsumer) handleMessage(ctx context.Context, msg kafka.Message) error {
	// Десериализация и вызов UseCase
	return nil
}
```

---

### Dead Letter Queue (DLQ) и защита от Poison Pill
Если в топик попало поврежденное сообщение («ядовитая таблетка»), консьюмер будет падать на нем бесконечно, останавливая вычитку всей партиции.

#### Паттерн DLQ с заголовком попыток (`x-retry-count`):
```go
func (c *OrderConsumer) HandleWithDLQ(ctx context.Context, msg kafka.Message, dlqWriter *kafka.Writer) {
	err := c.processOrder(ctx, msg.Value)
	if err == nil {
		_ = c.reader.CommitMessages(ctx, msg)
		return
	}

	// Считываем текущее число ретраев из заголовков Kafka:
	retries := 0
	for _, h := range msg.Headers {
		if h.Key == "x-retry-count" && len(h.Value) > 0 {
			retries = int(h.Value[0])
		}
	}

	if retries < 3 {
		// Отправляем в retry-топик с инкрементом счетчика:
		msg.Headers = append(msg.Headers, kafka.Header{Key: "x-retry-count", Value: []byte{byte(retries + 1)}})
		_ = dlqWriter.WriteMessages(ctx, kafka.Message{
			Topic:   "orders-retry",
			Key:     msg.Key,
			Value:   msg.Value,
			Headers: msg.Headers,
		})
	} else {
		// Превышен лимит ретраев -> отправляем в DLQ топик для ручного разбора инженерами:
		c.logger.Warn("message moved to DLQ", "key", string(msg.Key), "error", err)
		_ = dlqWriter.WriteMessages(ctx, kafka.Message{
			Topic:   "orders-dlq",
			Key:     msg.Key,
			Value:   msg.Value,
			Headers: append(msg.Headers, kafka.Header{Key: "x-dead-reason", Value: []byte(err.Error())}),
		})
	}

	// Коммитим оффсет основного топика, чтобы очередь не зависала:
	_ = c.reader.CommitMessages(ctx, msg)
}
```

---

## 11. gRPC Interceptors: Цепочки Middleware для Enterprise-микросервисов

### Зачем нужны интерцепторы
В gRPC интерцепторы выполняют роль HTTP Middleware: они перехватывают каждый входящий вызов до того, как управление попадет в хэндлер.  
В микросервисной архитектуре сервер **обязан защищаться цепочкой интерцепторов**:
1. **Recovery:** перехватывать паники внутри методов и возвращать статус `codes.Internal`, не роняя pod.
2. **Logging:** логировать имя RPC метода, параметры вызова, время выполнения и код завершения через `slog`.
3. **Metrics:** собирать Prometheus-гистограммы латентности и счетчики кодов ошибок.
4. **Tracing:** извлекать контекст W3C Trace Context (`traceparent`) из gRPC-метаданных.

---

### Полная реализация цепочки интерцепторов на Go:

```go
package grpcserver

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"runtime/debug"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	grpcDurationHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "grpc_request_duration_seconds",
		Help:    "gRPC request duration in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "code"})
)

// ---------------------------------------------------------------------------
// 1. RECOVERY INTERCEPTOR: Защита от падений при panic
// ---------------------------------------------------------------------------
func RecoveryInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		defer func() {
			if r := recover(); r != nil {
				stack := string(debug.Stack())
				logger.ErrorContext(ctx, "panic recovered in gRPC handler",
					"method", info.FullMethod,
					"panic", fmt.Sprintf("%v", r),
					"stack", stack,
				)
				err = status.Errorf(codes.Internal, "internal server error: panic recovered")
			}
		}()
		return handler(ctx, req)
	}
}

// ---------------------------------------------------------------------------
// 2. LOGGING & METRICS INTERCEPTOR: Наблюдаемость каждого RPC
// ---------------------------------------------------------------------------
func LoggingMetricsInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		start := time.Now()

		resp, err := handler(ctx, req)

		duration := time.Since(start)
		st, _ := status.FromError(err)
		codeStr := st.Code().String()

		// Записываем метрику в Prometheus:
		grpcDurationHistogram.WithLabelValues(info.FullMethod, codeStr).Observe(duration.Seconds())

		// Структурированный лог вызова:
		level := slog.LevelInfo
		if st.Code() != codes.OK {
			level = slog.LevelWarn
		}
		if st.Code() == codes.Internal || st.Code() == codes.Unknown {
			level = slog.LevelError
		}

		logger.Log(ctx, level, "grpc request processed",
			"method", info.FullMethod,
			"code", codeStr,
			"duration_ms", duration.Milliseconds(),
		)

		return resp, err
	}
}

// ---------------------------------------------------------------------------
// 3. TRACE CONTEXT INTERCEPTOR: Проброс W3C TraceID из метаданных
// ---------------------------------------------------------------------------
func TraceContextInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			// Читаем заголовок X-Request-ID или W3C traceparent:
			if requestIDs := md.Get("x-request-id"); len(requestIDs) > 0 {
				ctx = context.WithValue(ctx, "request_id", requestIDs[0])
			}
		}
		return handler(ctx, req)
	}
}

// ---------------------------------------------------------------------------
// СБОРКА И ЗАПУСК СЕРВЕРА С ЦЕПОЧКОЙ INTERCEPTORS:
// ---------------------------------------------------------------------------
func StartGRPCServer(port string, logger *slog.Logger) (*grpc.Server, error) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return nil, err
	}

	// Порядок чейнинга критичен:
	// Сначала Recovery, затем Tracing, затем Logging/Metrics, затем сам Handler:
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			RecoveryInterceptor(logger),
			TraceContextInterceptor(),
			LoggingMetricsInterceptor(logger),
		),
	)

	logger.Info("gRPC server listening", "port", port)
	go func() {
		if err := server.Serve(lis); err != nil {
			logger.Error("gRPC server stopped", "error", err)
		}
	}()

	return server, nil
}
```

---

## 12. Безопасность приложений на Go (AppSec): Защита от OOM, SSRF, Timing Attacks

### 1. Защита от отказа в обслуживании (DoS & Memory Exhaustion)

#### ❌ Ловушка `io.ReadAll(r.Body)`:
```go
// ❌ СМЕРТЕЛЬНО ДЛЯ PRODUCTION:
body, err := io.ReadAll(r.Body) // Если злоумышленник отправит 2 ГБ поток, сервер получит OOMKilled!
```

#### ✅ Эталонная защита через `http.MaxBytesReader`:
```go
func HandleUpload(w http.ResponseWriter, r *http.Request) {
    const maxUploadSize = 10 << 20 // 10 Мегабайт

    // MaxBytesReader оборачивает r.Body и возвращает ошибку при превышении лимита:
    r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

    body, err := io.ReadAll(r.Body)
    if err != nil {
        var maxBytesErr *http.MaxBytesError
        if errors.As(err, &maxBytesErr) {
            http.Error(w, "Request body is too large (max 10MB)", http.StatusRequestEntityTooLarge)
            return
        }
        http.Error(w, "Failed to read request body", http.StatusBadRequest)
        return
    }
    // Обработка body
}
```

#### Обязательная защита от Slowloris атак на `http.Server`:
Злоумышленник открывает тысячи TCP-соединений и шлет заголовки по 1 байту раз в минуту, парализуя сервер.
```go
server := &http.Server{
    Addr:              ":8080",
    Handler:           router,
    ReadHeaderTimeout: 2 * time.Second,  // 🛡️ Защита от Slowloris: заголовки должны прийти за 2 секунды!
    ReadTimeout:       10 * time.Second, // Время на чтение тела запроса
    WriteTimeout:      15 * time.Second, // Время на формирование ответа клиенту
    IdleTimeout:       60 * time.Second, // Таймаут keep-alive соединений
    MaxHeaderBytes:    1 << 19,          // 512 КБ максимум на заголовки
}
```

---

### 2. Криптография и защита от Timing Attacks

#### ❌ Опасность `math/rand`:
Пакет `math/rand` генерирует **детерминированные псевдослучайные числа**. Если сгенерировать токен сброса пароля или сессионный ID через `math/rand`, злоумышленник, зная алгоритм и пару сгенерированных чисел, может математически восстановить все будущие токены.

#### ✅ Безопасная генерация токенов через `crypto/rand`:
```go
import (
    "crypto/rand"
    "encoding/hex"
)

func GenerateSecureToken(lengthBytes int) (string, error) {
    bytes := make([]byte, lengthBytes)
    if _, err := rand.Read(bytes); err != nil {
        return "", fmt.Errorf("failed to generate secure random bytes: %w", err)
    }
    return hex.EncodeToString(bytes), nil
}
```

#### Защита от атак по времени (Timing Attack) при сравнении хешей:
При стандартном сравнении строк `token == expectedToken` процессор прерывает сравнение **на первом несовпавшем байте**. Замеряя наносекунды отклика сервера, злоумышленник может посимвольно угадать API-ключ или HMAC-подпись!

```go
import "crypto/subtle"

// ❌ Уязвимо к Timing Attack:
// if clientToken == secretApiKey { ... }

// ✅ Безопасно: константное время сравнения независимо от совпадения символов:
func SecureCompare(a, b string) bool {
    return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}
```

---

### 3. Защита от Server-Side Request Forgery (SSRF)

#### Угроза:
Если ваше приложение умеет загружать вебхуки или аватарки по URL, переданному пользователем (`POST /fetch-avatar?url=...`), атакующий передаст адрес внутренней инфраструктуры:
* `http://169.254.169.254/latest/meta-data/` — сервис метаданных AWS/GCP, отдающий IAM-ключи и секреты инфраструктуры.
* `http://10.0.0.5:5432` — сканирование внутренних баз данных компании.

#### ✅ Безопасный HTTP-клиент с фильтрацией приватных IP:
```go
package security

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"syscall"
	"time"
)

// NewSafeHTTPClient создает клиент, запрещающий подключения к loopback, private и link-local адресам:
func NewSafeHTTPClient() *http.Client {
	dialer := &net.Dialer{
		Timeout: 5 * time.Second,
		Control: func(network, address string, c syscall.RawConn) error {
			host, _, err := net.SplitHostPort(address)
			if err != nil {
				return err
			}

			ip := net.ParseIP(host)
			if ip == nil {
				return fmt.Errorf("invalid host IP: %s", host)
			}

			// Блокируем приватные и служебные диапазоны:
			if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
				return fmt.Errorf("access to internal / private network is forbidden (SSRF protection): %s", ip)
			}

			return nil
		},
	}

	transport := &http.Transport{
		DialContext:         dialer.DialContext,
		TLSHandshakeTimeout: 5 * time.Second,
	}

	return &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}
}
```

---

### 4. Защита JWT: предотвращение атаки «alg: none»
Если библиотека парсит JWT без строгой проверки метода подписи, атакующий может подменить заголовок токена на `{"alg": "none"}` (без подписи) или подсунуть публичный RSA-ключ в HMAC-валидатор.

```go
import (
    "fmt"
    "github.com/golang-jwt/jwt/v5"
)

func ValidateToken(tokenString string, secretKey []byte) (*jwt.Token, error) {
    return jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
        // 🛡️ Строго проверяем, что алгоритм подписи именно HMAC-SHA256:
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return secretKey, nil
    })
}
```
