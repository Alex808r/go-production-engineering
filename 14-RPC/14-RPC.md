# 14-RPC — Удаленный вызов процедур: net/rpc и gRPC (Protocol Buffers)

В данной директории сравниваются два ключевых подхода к реализации RPC (Remote Procedure Call) в экосистеме Go:
1. Встроенный пакет стандартной библиотеки `net/rpc` на базе бинарного кодека `encoding/gob`.
2. Индустриальный стандарт межсервисного взаимодействия `gRPC` на базе HTTP/2 и бинарной сериализации Protocol Buffers v3.

---

## 📋 Структура каталога

```text
14-RPC/
├── 1-Go-RPC/            # Встроенный стандартный механизм net/rpc
│   ├── cmd/client/      # RPC-клиент (rpc.Dial, client.Call, обработка ошибок)
│   ├── cmd/server/      # RPC-сервер (rpc.Register, net.Listen, rpc.ServeConn)
│   └── pkg/books/       # Контракт данных (Book, Request) и тесты над net.Pipe
├── 2-gRPC/              # Высокопроизводительный gRPC сервис
│   ├── books_proto/     # Protobuf-контракт (books.proto) и сгенерированный код (.pb.go)
│   ├── client/          # gRPC-клиент (Server Streaming через stream.Recv)
│   └── server/          # gRPC-сервер (Streaming + Unary) и in-memory тесты (bufconn)
└── homework-14/         # Практическое ДЗ: RPC-служба накопления сообщений
    ├── cmd/server/      # Сервер сообщений на net/rpc с Graceful Shutdown
    ├── cmd/client/      # Клиент сообщений (отправка батчей Send и вычитка Messages)
    ├── pkg/messages/    # Модель сообщения (Message{ID, Time, Text})
    └── pkg/rpcservice/  # Потокобезопасный сервис (sync.RWMutex) и тесты над net.Pipe
```

---

## 🔍 Ключевые концепции

1. **Стандартный `net/rpc`:**
   - Простота: не требует внешних компиляторов и сторонних зависимостей.
   - Требования к сигнатуре экспортируемого метода:
     ```go
     func (t *T) MethodName(req T1, resp *T2) error
     ```
   - **Важность разыменования указателя:** передача аргументов в Go происходит по значению. Поэтому `resp = &data` не изменит память вызывающей стороны. Данные необходимо записывать через разыменование: `*resp = data`.
   - Сериализация: пакет `encoding/gob` (все сериализуемые поля структур обязаны быть экспортируемыми).

2. **gRPC и Protocol Buffers (`2-gRPC`):**
   - **Кроссплатформенность:** контракт описывается в `.proto` файле и компилируется утилитой `protoc` в типизированный код для Go, Java, Python, TypeScript и др.
   - **Эффективность:** бинарный кодек Protobuf значительно превосходит JSON по компактности и скорости сериализации.
   - **Транспорт HTTP/2:** мультиплексирование сотен запросов в рамках одного TCP-сокета, бинарное фреймирование, сжатие заголовков (HPACK) и потоковая передача (Streaming).
   - **Прямая совместимость:** встраивание `pb.Unimplemented<Service>Server` защищает от поломки сборки при добавлении новых методов в интерфейс сервиса.

3. **Изолированное тестирование без конфликтов портов:**
   - Для `net/rpc`: тестирование полного сетевого цикла через виртуальный двунаправленный канал `net.Pipe()`.
   - Для `gRPC`: тестирование через in-memory буферизированный слушатель `google.golang.org/grpc/test/bufconn`.

---

## 🚀 Запуск и проверка

```bash
# Запуск всех модульных тестов и бенчмарков модуля 14
go test -v -race -bench=. ./14-RPC/...

# Запуск сервера 1-Go-RPC
go run ./14-RPC/1-Go-RPC/cmd/server/rpc-server.go

# Запуск клиента 1-Go-RPC (в отдельном терминале)
go run ./14-RPC/1-Go-RPC/cmd/client/rpc-client.go

# Запуск gRPC-сервера
go run ./14-RPC/2-gRPC/server/grpc-server.go

# Запуск gRPC-клиента (в отдельном терминале)
go run ./14-RPC/2-gRPC/client/grpc-client.go

# Запуск RPC-сервера сообщений (ДЗ-14)
go run ./14-RPC/homework-14/cmd/server/main.go

# Запуск RPC-клиента сообщений (ДЗ-14)
go run ./14-RPC/homework-14/cmd/client/main.go
```

---

## 🔗 Связанные материалы
- **Теория:** [Том 5: Тестирование и архитектура (Раздел 23: RPC, Protobuf, gRPC)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/05-go-core-testing-arch.md)
- **Практическое ДЗ:** [homework-14](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/14-RPC/homework-14) (Реализация микросервиса сообщений на Go-RPC)

