# 22-modern-go — Современные возможности Go 1.21–1.24 (State-of-the-Art Go)

Данный модуль расширяет базовую программу курса актуальными передовыми практиками, внедренными в версиях **Go 1.21, 1.22, 1.23 и 1.24**, которые активно применяются в современном продакшене и проверяются на собеседованиях уровня Middle+ / Senior.

---

## 📋 Структура каталога

```text
22-modern-go/
├── 01-slog/               # Структурированное логирование log/slog (Go 1.21+)
│   ├── slog_demo.go       # JSON/Text хэндлеры, маскирование LogValuer, контекст
│   └── slog_demo_test.go  # Тесты безопасности маскирования и бенчмарки
├── 02-fuzzing/            # Встроенный фаззинг testing.F (Go 1.18+)
│   ├── parser.go          # Парсер URL-параметров с защитой от паник
│   └── parser_test.go     # Нативный фаззинг-тест FuzzParseQuery
├── 03-std-generics/       # Стандартные generic-пакеты slices, maps, cmp (Go 1.21+)
│   ├── std_generics.go    # BinarySearchFunc, SortFunc, Compact, Clone, cmp.Or
│   └── std_generics_test.go # Сравнительные бенчмарки со старым sort.Slice
├── 04-enhanced-routing/   # Продвинутая маршрутизация net/http.ServeMux (Go 1.22+)
│   ├── router.go          # Паттерны методов "POST /items", PathValue("{id}"), wildcards
│   └── router_test.go     # Тесты автоматического 405 Method Not Allowed
├── 05-iterators/          # Range-over-func iter.Seq и iter.Seq2 (Go 1.23+)
│   ├── iterators.go       # Симметричный обход BST дерева без аллокаций, Filter
│   └── iterators_test.go  # Тесты early-break и бенчмарк iter.Seq vs Slice
├── 06-graceful-shutdown/  # Паттерн безопасного завершения серверов
│   ├── server.go          # signal.NotifyContext, srv.Shutdown, завершение in-flight
│   └── server_test.go     # Тест отсутствия обрывов активных запросов
├── 07-secrets-and-traps/  # Скрытые механизмы и ловушки Go (Uber Style & 100 Mistakes)
│   ├── traps.go           # Struct padding, slice retention, Functional Options, timer leaks
│   └── traps_test.go      # Тесты выравнивания памяти (33% savings), бенчмарки
└── 22-modern-go.md        # Данное руководство
```

---

## 🔍 Подробный разбор концепций

### 1. `log/slog` — Официальный структурированный логгер ([`01-slog`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/22-modern-go/01-slog/slog_demo.go))
До Go 1.21 в стандартной библиотеке существовал только простейший `log.Println`. Команды были вынуждены подключать тяжелые внешние библиотеки (`uber-go/zap`, `rs/zerolog`).
- **`slog`** объединил производительность, гибкость и стандартный API.
- Поддерживает структурированные атрибуты: `slog.Int`, `slog.String`, `slog.Group`.
- **Маскирование чувствительных данных (LogValuer):**
  Любой тип, реализующий `slog.LogValuer`, может самостоятельно определить, как он отображается в логах. Это полностью исключает утечку паролей и токенов:
  ```go
  type Password string
  func (p Password) LogValue() slog.Value {
      return slog.StringValue("******")
  }
  ```

---

### 2. Фаззинг-тестирование `testing.F` ([`02-fuzzing`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/22-modern-go/02-fuzzing/parser.go))
В отличие от классических тестов, где разработчик сам придумывает сценарии, фаззер непрерывно генерирует случайные мутации входных данных (Seed Corpus), стремясь вызвать панику (`panic`), деление на ноль или выход за пределы массива:
```go
func FuzzParseQuery(f *testing.F) {
    f.Add("key=val&a=b") // Сид-корпус
    f.Fuzz(func(t *testing.T, input string) {
        _, _ = ParseQuery(input) // Главный инвариант: никаких паник!
    })
}
```
Команда запуска:
```bash
go test -fuzz=FuzzParseQuery -fuzztime=5s ./lessons/22-modern-go/02-fuzzing/...
```

---

### 3. Пакеты `slices`, `maps`, `cmp` ([`03-std-generics`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/22-modern-go/03-std-generics/std_generics.go))
С выходом Go 1.21 стандартная библиотека получила мощный набор типизированных функций:
- `slices.Sort(s)` и `slices.SortFunc`: работают в **~2 раза быстрее** старого `sort.Slice` за счет отсутствия рефлексии и интерфейсных оберток.
- `slices.BinarySearch` / `slices.BinarySearchFunc`: типизированный бинарный поиск $O(\log N)$.
- `slices.Compact`: удаление дубликатов на месте.
- `maps.Clone`, `maps.Equal`, `maps.DeleteFunc`: потокобезопасная и чистая работа с ассоциативными массивами.
- `cmp.Or(a, b, c)` (Go 1.22+): возвращает первое ненулевое значение. Идеально для конфигураций:
  ```go
  port := cmp.Or(os.Getenv("PORT"), cfg.Port, "8080")
  ```

---

### 4. Продвинутый роутер `net/http.ServeMux` ([`04-enhanced-routing`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/22-modern-go/04-enhanced-routing/router.go))
В Go 1.22 стандартный мультиплексор научился тому, ради чего раньше ставили `gorilla/mux` или `chi`:
1. **Метод в строке паттерна:** `"POST /items"`, `"DELETE /items/{id}"`.
2. **Параметры пути:** `"GET /items/{id}"` извлекаются через `id := r.PathValue("id")`.
3. **Wildcard пути:** `"GET /files/{path...}"` захватывает остаток пути.
4. **Автоматический статус 405:** если путь совпал, но метод не совпал, Go автоматически возвращает `405 Method Not Allowed`.

---

### 5. Range-over-func итераторы `iter.Seq` ([`05-iterators`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/22-modern-go/05-iterators/iterators.go))
В Go 1.23 появилась возможность использовать `for ... := range` для обхода **любых кастомных структур** (деревьев, графов, связных списков, генераторов чисел).
- Сигнатура итератора одного значения: `type Seq[V any] func(yield func(V) bool)`.
- Сигнатура для пар (ключ-значение): `type Seq2[K, V any] func(yield func(K, V) bool)`.
- **Преимущества перед старыми подходами:**
  - Не требует выделять промежуточный срез в памяти (`tree.Slice()`);
  - Не требует создавать горутины и каналы (`chan T`), которые были в 100 раз медленнее;
  - Поддерживает оператор `break`: когда тело цикла делает `break`, функция `yield` возвращает `false`, и обход мгновенно останавливается.

```go
// Обход бинарного дерева без выделения памяти в куче
for val := range root.All() {
    if val == 40 {
        break // Мгновенный выход
    }
}
```

---

### 6. Паттерн Graceful Shutdown ([`06-graceful-shutdown`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/22-modern-go/06-graceful-shutdown/server.go))
В продакшене (особенно в Kubernetes) контейнеры регулярно перезапускаются при деплоях. Если завершить процесс через `os.Exit(0)`, активные HTTP-запросы пользователей будут оборваны с ошибкой 502/Connection Reset.

**Эталонная схема Graceful Shutdown:**
1. Перехват сигналов через `signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)`.
2. Ожидание отмены контекста.
3. Вызов `srv.Shutdown(shutdownCtx)` с таймаутом (например, 5 секунд).
4. Сервер прекращает принимать новые запросы, но **дожидается завершения всех текущих**.
5. Ожидание завершения фоновых воркеров через `sync.WaitGroup`.

---

### 7. Секреты рантайма и ловушки памяти ([`07-secrets-and-traps`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/22-modern-go/07-secrets-and-traps/traps.go))
Критически важные знания из **Uber Go Style Guide** и книги **100 Go Mistakes**:
- **Выравнивание структур (Data Alignment):** Размещение полей по убыванию размера сокращает размер структуры на 33.3% за счет устранения пустого байтового отступа (padding).
- **Утечки памяти при срезах:** Подсрез `huge[:2]` удерживает весь гигантский базовый массив в RAM. Использование `slices.Clone()` изолирует память и позволяет GC удалить исходный массив.
- **Очистка map:** `delete()` не уменьшает аллокации бакетов `hmap`. Для освобождения памяти мапу необходимо пересоздавать.
- **Таймеры в циклах:** `time.After` внутри `select` в бесконечном цикле создает утечку памяти до истечения таймаута; решение — `time.NewTimer` с явным `Stop()`.
- **Kubernetes CFS Throttling:** Go не знает про лимиты cgroups из коробки — решение `go.uber.org/automaxprocs`.
- **Compile-time проверки интерфейсов:** `var _ Runner = (*Service)(nil)` защищает от рассинхронизации сигнатур еще на этапе компиляции.
- Полное руководство: см. [`go-documentation/go-secrets-and-traps-guide.md`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/go-documentation/go-secrets-and-traps-guide.md).

---

## 📊 Результаты бенчмарков

```text
BenchmarkSlog_JSONHandler-8       133734     8818 ns/op
BenchmarkSlog_TextHandler-8       131674     9141 ns/op
BenchmarkSlices_SortFunc-8        917683     1245 ns/op  (~2x быстрее sort.Slice!)
BenchmarkClassic_SortSlice-8      451970     2636 ns/op
BenchmarkTree_IterSeq-8          1542219      768 ns/op  (быстрее среза, 0 аллокаций)
BenchmarkTree_FullSlice-8        1386534      865 ns/op
BenchmarkGoodStruct_Allocation-8 369654117    3.15 ns/op  (16 байт vs 24 байта, 33% экономия RAM)
BenchmarkBadStruct_Allocation-8  378972750    3.30 ns/op
```

---

## 🚀 Запуск тестов модуля

```bash
# Запуск всех тестов и бенчмарков модуля
go test -v -race -bench=. ./lessons/22-modern-go/...

# Запуск фаззинг-теста на 3 секунды
go test -v -fuzz=FuzzParseQuery -fuzztime=3s ./lessons/22-modern-go/02-fuzzing/...
```
