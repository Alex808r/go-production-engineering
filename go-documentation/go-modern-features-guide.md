# Modern Go Features Guide (Go 1.21–1.24)

> **Назначение документа:** Комплексное руководство по ключевым нововведениям, архитектурным паттернам и библиотекам современных релизов Go (1.21, 1.22, 1.23 и 1.24).  
> **Практические примеры с тестами и бенчмарками:** см. модуль [`22-modern-go`](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/22-modern-go/22-modern-go.md).

---

## 🧭 Хронология ключевых изменений языка

| Версия Go | Главные фичи | Практическая польза |
| :--- | :--- | :--- |
| **Go 1.21** | `log/slog`, пакеты `slices`, `maps`, `cmp`, встроенные `min`, `max`, `clear` | Официальный структурированный логгер, generic-утилиты без рефлексии, удаление бойлерплейта |
| **Go 1.22** | Улучшенный роутер `net/http.ServeMux`, исправление `for` loop variable capture | Маршрутизация REST из коробки без сторонних библиотек (`chi`, `gorilla/mux`), устранение утечек переменных цикла |
| **Go 1.23** | Range-over-func (`iter.Seq`, `iter.Seq2`), пакет `unique` (интернирование строк) | Кастомные итераторы в цикле `for range` для деревьев и генераторов без аллокаций памяти |
| **Go 1.24** | Пакет `weak` (слабые указатели), Swiss Tables для runtime map, директива `tool` в `go.mod` | Эффективные кэши без утечек памяти, ускорение мап на ~30%, управление CLI-инструментами |

---

## 1. Структурированное логирование: `log/slog` (Go 1.21+)

### Проблема прошлого:
До Go 1.21 стандартная библиотека предоставляла только неструктурированный логгер `log.Println("user", u.ID, "logged in")`. В проде использовались `uber-go/zap` или `rs/zerolog`. Возникала фрагментация: библиотеки не могли договориться об общем интерфейсе логгера.

### Современное решение:
Пакет `log/slog` предлагает унифицированный API со сменными обработчиками (`slog.Handler`):
- `slog.NewJSONHandler` — для машинной обработки в ELK / Loki / CloudWatch;
- `slog.NewTextHandler` — для локальной разработки.

```go
logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelInfo,
}))

logger.InfoContext(ctx, "order created",
    slog.String("order_id", "ord_123"),
    slog.Float64("amount", 99.50),
    slog.Group("customer",
        slog.Int("id", 42),
        slog.String("tier", "gold"),
    ),
)
```

### Защита от утечки секретов (`slog.LogValuer`):
Любой тип может реализовать метод `LogValue() slog.Value`:
```go
type AuthToken string

func (t AuthToken) LogValue() slog.Value {
    return slog.StringValue("[REDACTED]")
}
```

---

## 2. Стандартные generic-пакеты: `slices`, `maps`, `cmp` (Go 1.21+)

### Сравнение с устаревшими подходами:

| Задача | До Go 1.21 | Начиная с Go 1.21 |
| :--- | :--- | :--- |
| **Сортировка** | `sort.Slice(s, func(i, j int) bool { return s[i] < s[j] })` (рефлексия) | `slices.Sort(s)` / `slices.SortFunc` (zero-alloc, ~2x быстрее) |
| **Бинарный поиск** | `sort.Search(len(s), func(i int) bool { return s[i] >= target })` | `idx, found := slices.BinarySearch(s, target)` |
| **Клонирование** | `clone := make([]T, len(s)); copy(clone, s)` | `clone := slices.Clone(s)` |
| **Очистка мапы** | `for k := range m { delete(m, k) }` | `clear(m)` (Go 1.21 builtin) |
| **Дефолтное значение** | `if v == "" { v = fallback }` | `v = cmp.Or(v, fallback)` |

---

## 3. Встроенная маршрутизация: `net/http.ServeMux` (Go 1.22+)

Начиная с Go 1.22, стандартный роутер полностью заменяет сторонние легковесные библиотеки:

```go
mux := http.NewServeMux()

// 1. Фильтрация по методу
mux.HandleFunc("POST /api/v1/orders", handleCreateOrder)

// 2. Именованный параметр пути {id}
mux.HandleFunc("GET /api/v1/orders/{id}", func(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id") // извлечение значения
    fmt.Fprintf(w, "Order: %s", id)
})

// 3. Catch-all параметр {path...}
mux.HandleFunc("GET /static/{path...}", handleStaticFiles)
```

> [!TIP]
> Если клиент отправляет `PUT /api/v1/orders` (для которого зарегистрирован только `POST`), Go 1.22+ автоматически отвечает статусом **`405 Method Not Allowed`** и заголовком `Allow: POST`.

---

## 4. Range-over-func: Итераторы `iter.Seq` (Go 1.23+)

Go 1.23 расширил цикл `for ... := range` для пользовательских функций-итераторов.

### Типы итераторов:
1. `iter.Seq[V any]` — итерация по одиночным значениям (`for v := range seq`);
2. `iter.Seq2[K, V any]` — итерация парами (`for k, v := range seq2`).

### Реализация обхода дерева без срезов и каналов:
```go
func (n *Node) InOrder() iter.Seq[int] {
    return func(yield func(int) bool) {
        var walk func(*Node) bool
        walk = func(curr *Node) bool {
            if curr == nil {
                return true
            }
            if !walk(curr.Left) { return false }
            if !yield(curr.Val)  { return false } // прерывание при break
            return walk(curr.Right)
        }
        walk(n)
    }
}

// Использование:
for val := range root.InOrder() {
    if val == 42 {
        break // немедленно прекратит рекурсию
    }
}
```

---

## 5. Встроенный фаззинг: `testing.F`

Фаззинг встроен непосредственно в утилиту `go test`:
```go
func FuzzDecode(f *testing.F) {
    // Начальный корпус эталонных входных данных
    f.Add([]byte("{\"key\": \"value\"}"))

    f.Fuzz(func(t *testing.T, data []byte) {
        var obj MyStruct
        // Инвариант: парсер не должен падать с panic
        _ = json.Unmarshal(data, &obj)
    })
}
```
Запуск:
```bash
go test -fuzz=FuzzDecode -fuzztime=30s
```

---

## 6. Паттерн Graceful Shutdown

Эталонный подход к остановке HTTP-серверов в микросервисной инфраструктуре:

```go
func RunServer(addr string) error {
    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer stop()

    server := &http.Server{Addr: addr, Handler: appMux}

    go func() {
        if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
            log.Fatalf("server error: %v", err)
        }
    }()

    <-ctx.Done() // ожидаем SIGINT / SIGTERM

    shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    return server.Shutdown(shutdownCtx) // ждем завершения активных соединений
}
```

---

## 🔗 Связанные материалы
- [Практический модуль 22-modern-go](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/22-modern-go/22-modern-go.md)
- [Том 1: Основы и синтаксис Go](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/01-go-core-basics.md)
- [Том 3: Конкурентность и рантайм](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/03-go-core-concurrency.md)
- [Том 5: Тестирование и архитектура](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/05-go-core-testing-arch.md)
