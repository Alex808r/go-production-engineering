# 11-network — Сетевое программирование на Go: сокеты, TCP и дедлайны

В данной директории изучаются низкоуровневые сетевые возможности стандартного пакета `net`: разработка TCP-серверов и клиентов, обработка параллельных клиентских подключений, управление сетевыми дедлайнами и тестирование сетевых служб.

---

## 📋 Структура каталога

```text
11-network/
├── 1-daytime/       # Простая сетевая служба времени (Daytime protocol, RFC 867)
│   ├── client/      # Клиент: подключение через net.Dial
│   └── server/      # Сервер: прослушивание порта через net.Listen("tcp4", ":1313")
├── 2-echo/          # Полноценный TCP эхо-сервер (обработка каждого подключения в горутине)
├── 3-deadline/      # Управление таймаутами сокета: conn.SetDeadline, SetReadDeadline (idle timeout)
├── 4-testing/       # Тестирование сетевых компонентов в памяти с помощью net.Pipe()
└── homework-11/     # Домашнее задание 11: сетевая служба GoSearch (netsrv на :8000) и netclient
```

---

## 🔍 Ключевые концепции

1. **Архитектура конкурентного TCP-сервера (`net`):**
   ```go
   listener, err := net.Listen("tcp4", ":8000")
   for {
       conn, err := listener.Accept()
       if err != nil { continue }
       go handleConnection(conn) // Обработка каждого клиента в отдельной горутине
   }
   ```

2. **Сетевые таймауты (`3-deadline`):**
   - Сокеты без дедлайнов могут "висеть" неделями при обрыве сети или зависании удаленного узла.
   - Вызовы `conn.SetReadDeadline(time.Now().Add(timeout))` перед чтением гарантируют закрытие сокета при отсутствии активности (Idle Timeout).

3. **Интерфейс `net.Conn`:**
   - Реализует `io.Reader`, `io.Writer`, `io.Closer`, что позволяет использовать `bufio.Scanner`, `io.Copy` и стандартные потоки непосредственно на сокетах.

4. **Тестирование сетевого кода через `net.Pipe()` (`4-testing`):**
   - Виртуальный синхронный дуплексный канал в оперативной памяти без открытия сетевых портов ОС и без конфликтов `address already in use`.

---

## 🚀 Запуск и проверка

```bash
# 1. Запуск Daytime сервера и клиента
go run ./lessons/11-network/1-daytime/server/daytime-server.go
go run ./lessons/11-network/1-daytime/client/daytime-client.go

# 2. Запуск TCP эхо-сервера
go run ./lessons/11-network/2-echo/echo.go

# 3. Запуск сервера с таймаутами
go run ./lessons/11-network/3-deadline/deadline.go

# 4. Запуск тестов сетевого модуля и ДЗ
go test -v -race -cover ./lessons/11-network/...
```

---

## 🔗 Связанные материалы
- **Теория:** [Том 4: Backend и Web (Раздел 29: Сетевое программирование)](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/go-documentation/04-go-core-backend-web.md)
- **Практическое ДЗ:** [homework-11](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/11-network/homework-11) (Сетевая TCP-служба поиска для GoSearch с клиентом netclient)

