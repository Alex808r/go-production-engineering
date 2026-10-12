# Домашнее задание №12: Веб-приложения на Go

Веб-приложение для просмотра содержимого поискового индекса и хранилища проиндексированных документов поисковика GoSearch.

## Архитектура решения

Решение полностью реализовано в рамках Занятия №12 («Веб-приложения») и состоит из:

1. **Пакет `pkg/webapp` (Задача №1):**
   - **Маршрутизация**: маршрутизатор `github.com/gorilla/mux` со строгой фильтрацией методов (`http.MethodGet`).
   - **Шаблонизация**: безопасные шаблоны стандартной библиотеки `html/template` с автоматическим экранированием от XSS-атак.
   - **Эндпоинты**:
     - `GET /docs` — просмотр всех документов в хранилище (таблица с ID, Title, кликабельными URL).
     - `GET /index` — просмотр инвертированного поискового индекса (термины, количество вхождений, список ID документов).
   - **Поддержка форматов (Content Negotiation)**:
     - По умолчанию в браузере отдается стилизованный HTML (`Content-Type: text/html`).
     - При передаче заголовка `Accept: application/json` или query-параметра `?format=json` отдается чистый JSON (`Content-Type: application/json`).
   - **Потокобезопасность**: внутреннее хранилище `SearchDB` защищено `sync.RWMutex` для конкурентных чтений и записей.

2. **Модульные тесты `pkg/webapp/webapp_test.go` (Задача №2):**
   - Реализованы с использованием `net/http/httptest` (`httptest.NewRequest`, `httptest.NewRecorder`).
   - Выполняются без поднятия реальных сетевых сокетов (zero port collision).
   - Проверяют:
     - Рендеринг HTML для `/docs` и `/index`.
     - Корректность JSON-выгрузки.
     - Защиту методов: возврат `405 Method Not Allowed` для `POST`, `PUT`, `DELETE`.
     - Обработку пустых данных и ошибок `nil`-провайдера (`500 Internal Server Error`).
     - Несуществующие маршруты (`404 Not Found`).

3. **Серверное приложение `cmd/webapp/main.go`:**
   - Инициализирует базу с демонстрационным набором документов.
   - Настраивает таймауты `http.Server` (`ReadTimeout`, `WriteTimeout`, `IdleTimeout`).
   - Обеспечивает корректное завершение работы (**Graceful Shutdown**) по сигналам `SIGINT`/`SIGTERM`.

---

## Запуск приложения

```bash
go run ./12-web-apps/homework-12/cmd/webapp/main.go
```

Приложение будет доступно по адресам:
- Документы (HTML): `http://localhost:8080/docs`
- Индекс (HTML): `http://localhost:8080/index`
- Документы (JSON): `http://localhost:8080/docs?format=json`
- Индекс (JSON): `http://localhost:8080/index?format=json`

---

## Запуск тестов

```bash
go test -v -race -cover ./12-web-apps/homework-12/...
```
