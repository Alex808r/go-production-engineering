# Домашнее задание №14: Удалённый вызов процедур (RPC-служба сообщений)

Решение домашнего задания по курсу «Современная разработка на Go» (Занятие №14: Удалённый вызов процедур (RPC)).  
Реализовано в отдельной директории `homework-14` без затирания или изменения предыдущих решений.

---

## 📑 Содержание задач

1. [Задача №1: Приложение-сервер RPC (Send и Messages)](#задача-1-приложение-сервер-rpc)
2. [Задача №2: Приложение-клиент RPC](#задача-2-приложение-клиент-rpc)
3. [Тестирование: модульные тесты и сетевое тестирование в памяти (net.Pipe)](#тестирование)
4. [Инструкция по запуску](#инструкция-по-запуску)

---

## Задача №1: Приложение-сервер RPC

Сервер реализован на базе стандартного пакета Go `net/rpc` в соответствии со всеми требованиями:
- **Структура сообщения ([`pkg/messages/messages.go`](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/14-RPC/homework-14/pkg/messages/messages.go)):**
  ```go
  type Message struct {
      ID   int       `json:"id"`
      Time time.Time `json:"time"`
      Text string    `json:"text"`
  }
  ```
- **Сервис RPC ([`pkg/rpcservice/service.go`](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/14-RPC/homework-14/pkg/rpcservice/service.go)):**
  - Экспортируемый тип `Service` с потокобезопасным хранилищем (`sync.RWMutex`).
  - Метод `Send(req []messages.Message, resp *int) error` — принимает срез сообщений, безопасно добавляет их в накопительный массив и возвращает число добавленных сообщений.
  - Метод `Messages(req struct{}, resp *[]messages.Message) error` — возвращает клиенту копию среза всех накопленных сообщений.
- **Серверное приложение ([`cmd/server/main.go`](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/14-RPC/homework-14/cmd/server/main.go)):**
  - Регистрирует сервис через `rpc.Register(svc)`.
  - Открывает сетевой TCP-порт `:8080` (`net.Listen`).
  - Каждое входящее клиентское соединение обрабатывается конкурентно в отдельной горутине через `go rpc.ServeConn(conn)`.
  - Поддерживает Graceful Shutdown по системным сигналам `SIGINT` и `SIGTERM`.

---

## Задача №2: Приложение-клиент RPC

Клиент ([`cmd/client/main.go`](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/14-RPC/homework-14/cmd/client/main.go)):
- Устанавливает TCP-соединение с сервером по адресу `localhost:8080` через `rpc.Dial`.
- Формирует пакет сообщений со структурой `{ID, Time, Text}`.
- Вызывает удаленный метод `client.Call("Service.Send", newBatch, &ackCount)`.
- Вызывает метод `client.Call("Service.Messages", struct{}{}, &allMessages)` для получения истории сообщений.
- Форматирует и выводит результат в виде таблицы в консоль.

---

## Тестирование

Тесты реализованы в [`pkg/rpcservice/service_test.go`](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/14-RPC/homework-14/pkg/rpcservice/service_test.go):
1. **Прямое тестирование методов:** проверка логики добавления и извлечения сообщений.
2. **Сетевое тестирование в памяти через `net.Pipe()`:**
   - Позволяет тестировать полный цикл RPC (сериализацию `encoding/gob`, передачу по сетевому интерфейсу, демаршалинг и вызов на сервере) **без открытия реальных TCP-портов операционной системы**.
   - Тестирование синхронного (`client.Call`) и асинхронного (`client.Go`) взаимодействия.
3. **Конкурентные тесты под нагрузкой:** 20 одновременных горутин-клиентов с проверкой на race conditions (`go test -race`).

---

## Инструкция по запуску

### Запуск тестов:

```bash
cd go_course_4
go test -v -race -cover ./lessons/14-RPC/homework-14/...
```

### Запуск сервера и клиента:

1. В первом терминале запустите RPC-сервер:
   ```bash
   go run ./lessons/14-RPC/homework-14/cmd/server/main.go
   ```

2. Во втором терминале запустите RPC-клиент:
   ```bash
   go run ./lessons/14-RPC/homework-14/cmd/client/main.go
   ```
