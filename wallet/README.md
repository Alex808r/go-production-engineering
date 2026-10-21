# High-Load Wallet REST API Service

Высоконагруженный сервис управления балансом кошельков на **Go**, **PostgreSQL** и **Docker Compose**, оптимизированный для работы в условиях жесткой конкуренции (**1000+ RPS по одному кошельку**) без потерь данных и без 50x ошибок.

---

## 🚀 Архитектура и решение проблемы конкурентности (1000 RPS на 1 кошелек)

При одновременной отправке 1000 запросов на изменение баланса **одного и того же кошелька** наивные решения (напр. `SELECT FOR UPDATE` в каждой транзакции) приводят к:
1. **Истощению пула соединений БД** (`connection pool exhaustion`) — 1000 горутин захватывают соединения и встают в очередь ожидания блокировки строки в PostgreSQL.
2. **Взрыву задержек и таймаутам** (`context deadline exceeded`) $\rightarrow$ каскадные ошибки **HTTP 500/504**.
3. **Взаимным блокировкам** (`deadlocks` `40P01` или `serialization failures` `40001`).

### Примененный многоуровневый подход:

1. **Шардированный In-Memory замок (`KeyLocker` на 4096 шардов):**
   * Перед тем как запросить соединение из пула PostgreSQL, запрос сериализуется в памяти рантайма Go по хешу UUID кошелька (`Striped Mutex`).
   * Запросы к **одному** кошельку выполняются последовательно, сменяя друг друга каждые ~0.3–0.5 мс.
   * Запросы к **разным** кошелькам выполняются параллельно без блокировок.
   * За счет этого для 1000 одновременных запросов по одному кошельку используется **ровно одно соединение с БД одновременно**! Остальные 999 запросов эффективно запаркованы планировщиком Go в userspace, не потребляя соединений БД и не создавая очередей в PostgreSQL.
   * Суммарно 1000 запросов отрабатывают за **~380 мс** (более **2600 RPS**).

2. **Атомарный SQL-запрос с проверкой баланса:**
   ```sql
   UPDATE wallets
   SET balance = balance + $1, updated_at = NOW()
   WHERE id = $2 AND (balance + $1 >= 0)
   RETURNING id, balance, created_at, updated_at;
   ```
   * Если средств недостаточно (`WITHDRAW` при нехватке денег), строка не обновляется вовсе.
   * Нет состояния гонки между проверкой баланса и списанием (Check-Then-Act).
   * Целостность гарантируется на уровне базы данных констрейнтом `CHECK (balance >= 0)`.

3. **Retry-политика с экспоненциальным бэкоффом:**
   * При возникновении транзиентных ошибок БД (`40001 serialization_failure`, `40P01 deadlock_detected`, временный обрыв сети) операция автоматически повторяется с экспоненциальной задержкой и джиттером (до 5 попыток).

4. **Тюнинг пула соединений (`pgxpool`):**
   * `MaxConns`: 100 соединений.
   * `MinConns`: 20 горячих соединений.
   * Асинхронные healthcheck-пинги и тайм-ауты.

---

## 📡 REST API Спецификация

### 1. Изменение баланса кошелька (DEPOSIT / WITHDRAW)
* **URL:** `POST /api/v1/wallet`
* **Content-Type:** `application/json`
* **Поддержка полей:** Поддерживает как `valletId` (из формулировки ТЗ), так и стандартный `walletId`.
* **Тело запроса:**
  ```json
  {
    "valletId": "11111111-1111-1111-1111-111111111111",
    "operationType": "DEPOSIT",
    "amount": 1000
  }
  ```
* **Успешный ответ (200 OK):**
  ```json
  {
    "walletId": "11111111-1111-1111-1111-111111111111",
    "balance": 11000
  }
  ```
* **Коды ответов:**
  * `200 OK` — операция успешно выполнена.
  * `400 Bad Request` — некорректный JSON, невалидный UUID, недопустимый тип операции, сумма $\le 0$, либо **недостаточно средств на балансе**.
  * `404 Not Found` — кошелек с данным UUID не найден.
  * `500 Internal Server Error` — исключена в штатных сценариях и под нагрузкой.

---

### 2. Получение баланса кошелька
* **URL:** `GET /api/v1/wallets/{WALLET_UUID}` (также поддерживается `/api/v1/wallet/{WALLET_UUID}`)
* **Успешный ответ (200 OK):**
  ```json
  {
    "walletId": "11111111-1111-1111-1111-111111111111",
    "balance": 11000
  }
  ```
* **Коды ответов:**
  * `200 OK` — баланс получен.
  * `400 Bad Request` — неверный формат UUID.
  * `404 Not Found` — кошелек не найден.

---

### 3. Вспомогательные эндпоинты
* `POST /api/v1/wallets` — создание нового кошелька (для тестирования).
  ```json
  {
    "walletId": "11111111-1111-1111-1111-111111111111",
    "initialBalance": 10000
  }
  ```
* `GET /health` — проверка работоспособности контейнера (Healthcheck).

---

## ⚙️ Конфигурация (`config.env`)

Все параметры загружаются из файла `config.env`:

```env
# Приложение
APP_PORT=8080
APP_ENV=production

# База данных PostgreSQL
DB_HOST=postgres
DB_PORT=5432
DB_USER=wallet_user
DB_PASSWORD=wallet_password
DB_NAME=wallet_db
DB_SSLMODE=disable
DB_MAX_CONNS=100
DB_MIN_CONNS=20
DB_MAX_IDLE_TIME=5m
DB_MAX_CONN_LIFETIME=1h

# Таймауты HTTP сервера
SERVER_READ_TIMEOUT=15s
SERVER_WRITE_TIMEOUT=30s
SERVER_SHUTDOWN_TIMEOUT=10s
```

---

## 🐳 Запуск через Docker Compose

Вся система (PostgreSQL 16 и приложение на Go) поднимается одной командой:

```bash
docker compose up --build -d
```

Проверить статус работы контейнеров:
```bash
docker compose ps
```

Посмотреть логи приложения:
```bash
docker compose logs -f app
```

Остановить сервисы:
```bash
docker compose down -v
```

---

## 🧪 Запуск тестов

Проект покрыт полным набором модульных и интеграционных стресс-тестов:

1. **Все тесты с детектором гонок памяти (`-race`):**
   ```bash
   go test -v -race ./...
   ```

2. **Нагрузочный стресс-тест на 1000 конкурентных RPS по одному кошельку:**
   ```bash
   go test -v -run TestConcurrency ./tests
   ```

### Результат прогона нагрузочного теста (1000 конкурентных запросов):
```text
=== RUN   TestConcurrency_1000_RPS_SingleWallet
    concurrency_test.go:162: 1000 requests completed in 380.075584ms (throughput: 2631.06 req/sec)
    concurrency_test.go:164: Results: 200 OK = 1000, 50x Errors = 0, Other Errors = 0
    concurrency_test.go:199: SUCCESS: Final balance verified exactly: 102500 (zero lost updates, zero 50x errors!)
--- PASS: TestConcurrency_1000_RPS_SingleWallet (0.65s)
=== RUN   TestConcurrency_WithdrawExceedingBalance
--- PASS: TestConcurrency_WithdrawExceedingBalance (0.02s)
PASS
```

---

## 📋 Примеры curl-запросов

```bash
# 1. Пополнение кошелька (DEPOSIT) на 1000
curl -X POST http://localhost:8080/api/v1/wallet \
  -H "Content-Type: application/json" \
  -d '{
    "valletId": "11111111-1111-1111-1111-111111111111",
    "operationType": "DEPOSIT",
    "amount": 1000
  }'

# 2. Списание с кошелька (WITHDRAW) на 500
curl -X POST http://localhost:8080/api/v1/wallet \
  -H "Content-Type: application/json" \
  -d '{
    "valletId": "11111111-1111-1111-1111-111111111111",
    "operationType": "WITHDRAW",
    "amount": 500
  }'

# 3. Получение текущего баланса
curl http://localhost:8080/api/v1/wallets/11111111-1111-1111-1111-111111111111
```
