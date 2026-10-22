# Шпаргалка по Go (Cheat Sheet)

> [!TIP]
> 📚 **Навигация:** [📖 Главное оглавление и путь изучения](README.md) — эта страница является **краткой справкой** по синтаксису и командам. Подробные объяснения, «ловушки» и внутреннее устройство — в томах [1.1](01-go-core-basics.md), [1.2](01-go-core-data-structures.md), [2](02-go-core-oop-methods.md), [3](03-go-core-concurrency.md), [4](04-go-core-backend-web.md), [5](05-go-core-testing-arch.md), [6](06-go-core-algorithms-interview.md). Пошаговый план курса — в [course-codebase-guide.md](course-codebase-guide.md), а все 32 домашних задания — в [homework-tasks.md](homework-tasks.md).

---

## 1. Команды тулчейна

| Команда | Что делает |
| :--- | :--- |
| `go version` | Показать версию Go |
| `go mod init example.com/app` | Создать модуль (файл `go.mod`) |
| `go mod tidy` | Добавить недостающие и убрать лишние зависимости |
| `go get github.com/user/pkg@latest` | Добавить/обновить зависимость |
| `go run .` / `go run main.go` | Скомпилировать во временную папку и запустить |
| `go build -o app ./cmd/app` | Собрать бинарный файл |
| `go fmt ./...` (`gofmt -w .`) | Отформатировать код |
| `go vet ./...` | Статический анализ (подозрительные конструкции) |
| `go test ./...` | Запустить все тесты модуля |
| `go test -v -run TestName -count 1 ./pkg` | Один тест, подробный вывод, без кэша |
| `go test -race ./...` | Тесты с детектором гонок данных |
| `go test -cover` / `-coverprofile=c.out` | Покрытие кода |
| `go test -bench=. -benchmem -run '^$'` | Бенчмарки с учётом памяти |
| `go doc fmt.Printf` | Документация из терминала |
| `go generate ./...` | Выполнить директивы `//go:generate` |
| `go tool pprof`, `go tool trace` | Профилирование и трассировка |
| `govulncheck ./...` | Сканирование уязвимостей в коде и зависимостях (Vulnerability Scanner) |

---

## 2. Типы и значения по умолчанию (Zero Value)

| Тип | Zero Value | Примечание |
| :--- | :--- | :--- |
| `int`, `int8…int64`, `uint…`, `float32/64` | `0` | `int` — 32 или 64 бита в зависимости от платформы |
| `bool` | `false` | |
| `string` | `""` | неизменяемая последовательность байт |
| `byte` (= `uint8`), `rune` (= `int32`) | `0` | `rune` — кодовая точка Unicode |
| указатель `*T`, `func`, `chan`, `map`, `[]T`, `interface` | `nil` | у `nil`-слайса и `nil`-мапы есть свои особенности (см. §9) |
| массив `[N]T` | все элементы — Zero Value | размер — часть типа |
| `struct` | все поля — Zero Value | |

---

## 3. Объявления и управление потоком

```go
var x int = 10          // явно
var y = "text"          // тип выводится
z := 3.14               // короткое объявление (только внутри функций)
const Pi = 3.14159      // константа
const (                 // перечисление через iota
    Red = iota          // 0
    Green               // 1
    Blue                // 2
)

if v, ok := m["key"]; ok { /* ... */ } else { /* ... */ }

for i := 0; i < 10; i++ { }      // классический цикл
for cond { }                     // как while
for { break }                    // бесконечный
for i, v := range slice { }      // индекс и значение (значение — копия!)
for i := range 5 { }             // Go 1.22+: 0..4
for k, v := range someMap { }    // порядок случайный

switch x {                       // break не нужен
case 1, 2:
    fmt.Println("один или два")
case 3:
    fmt.Println("три")
    fallthrough                  // явный «провал» в следующую ветку
default:
    fmt.Println("другое")
}

switch {                         // без выражения: вместо цепочки if/else
case x < 0:
case x == 0:
default:
}

switch v := i.(type) {           // type switch для интерфейсов
case int:
case string:
}
```

---

## 4. Функции

```go
func add(a, b int) int { return a + b }

func divide(a, b float64) (float64, error) {      // несколько результатов
    if b == 0 {
        return 0, errors.New("деление на ноль")
    }
    return a / b, nil
}

func sum(nums ...int) int { /* nums — это []int */ }  // вариативная функция; вызов: sum(s...)

counter := func() func() int {                    // замыкание
    n := 0
    return func() int { n++; return n }
}()

defer file.Close()   // выполнится при выходе из функции; порядок LIFO;
                     // аргументы вычисляются сразу, в момент defer
```

> **Всё передаётся по значению.** Срез, мапа, канал — это небольшие «дескрипторы», которые копируются, но указывают на общие данные (см. Том 1, раздел 2.3).

---

## 5. Слайсы, массивы, мапы, строки

```go
a := [3]int{1, 2, 3}              // массив (значение, размер — часть типа)
s := []int{1, 2, 3}               // слайс: {указатель, len, cap}
s = append(s, 4, 5)               // ВСЕГДА присваивайте результат append
s2 := make([]int, 0, 10)          // len=0, cap=10
sub := s[1:3]                     // общий базовый массив! (len=2)
sub3 := s[1:3:3]                  // ограничиваем cap, чтобы append не портил s
c := slices.Clone(s)              // независимая копия
copy(dst, src)                    // копирует min(len(dst), len(src)) элементов
s = slices.Delete(s, i, i+1)      // удалить элемент i (Go 1.21+)
slices.Sort(s); slices.Contains(s, 3); slices.Index(s, 3)
clear(s)                          // обнулить элементы (Go 1.21+)

m := map[string]int{"a": 1}       // или make(map[string]int)
v, ok := m["x"]                   // ok=false, если ключа нет (v — Zero Value)
delete(m, "a")                    // безопасно и для отсутствующего ключа
keys := slices.Sorted(maps.Keys(m)) // отсортированные ключи (Go 1.23+)

str := "Привет"
len(str)                          // 12 — байты, не символы!
utf8.RuneCountInString(str)       // 6 — руны
for i, r := range str { }         // i — смещение в байтах, r — руна
var sb strings.Builder            // эффективная склейка строк
sb.WriteString("go"); sb.String()
```

---

## 6. Структуры, методы, интерфейсы, дженерики

```go
type User struct {
    Name string `json:"name"`         // тег для encoding/json
    Age  int    `json:"age,omitempty"`
}

func (u User) Greet() string { return "Hi, " + u.Name }   // value receiver: работает с копией
func (u *User) Birthday()     { u.Age++ }                  // pointer receiver: меняет оригинал

u := User{Name: "Alex"}            // именованная инициализация — предпочтительно
p := &User{Name: "Bob"}            // указатель на структуру; p.Name — автоматическое разыменование

type Employee struct {             // встраивание (композиция, НЕ наследование)
    User
    Salary int
}

type Stringer interface{ String() string }   // интерфейс = набор методов; реализуется неявно

var s fmt.Stringer = user           // присваивание интерфейсу
str, ok := s.(fmt.Stringer)         // type assertion (используйте форму с ok)

type Number interface{ ~int | ~float64 }   // ограничение типа (constraint)
func Sum[T Number](xs []T) T {             // дженерик-функция
    var t T
    for _, x := range xs { t += x }
    return t
}
type Stack[T any] struct{ items []T }      // дженерик-тип
```

**Правила ресиверов:** если хотя бы одному методу нужен `*T` — используйте `*T` у всех методов типа. Экспортируется (виден из других пакетов) только то, что начинается с **заглавной буквы**.

---

## 7. Ошибки, panic, recover

```go
var ErrNotFound = errors.New("not found")                 // sentinel-ошибка

return fmt.Errorf("find user %d: %w", id, ErrNotFound)     // %w — обернуть (сохранить цепочку)

if errors.Is(err, ErrNotFound) { }                         // есть ли ошибка в цепочке
var ve *ValidationError
if errors.As(err, &ve) { /* ve.Field */ }                  // достать ошибку нужного типа
err = errors.Join(err1, err2)                              // объединить ошибки (Go 1.20+)

defer func() {                                             // recover — ТОЛЬКО внутри defer
    if r := recover(); r != nil { log.Println("panic:", r) }
}()
```

> `panic` — для невосстановимых программных ошибок, а не для обычных ошибок. `os.Exit`/`log.Fatal` не выполняют `defer`.

---

## 8. Конкурентность

```go
go doWork()                                   // запустить горутину
var wg sync.WaitGroup
wg.Add(1)                                     // Add ДО запуска горутины
go func() { defer wg.Done(); work() }()
wg.Wait()                                     // Go 1.25+: wg.Go(func() { work() })

ch := make(chan int)                          // небуферизованный: отправитель ждёт получателя
bch := make(chan string, 10)                  // буферизованный
ch <- 1; v := <-ch; v, ok := <-ch             // ok=false: канал закрыт и пуст
close(ch)                                     // закрывает ТОЛЬКО отправляющая сторона
for v := range ch { }                         // читает, пока канал не закрыт

select {                                      // ожидание нескольких каналов
case v := <-ch1:
case ch2 <- 42:
case <-ctx.Done():
case <-time.After(time.Second):
default:                                      // без default select блокируется
}

var mu sync.Mutex
mu.Lock(); defer mu.Unlock()                  // защита общей памяти
var n atomic.Int64; n.Add(1); n.Load()        // атомарные операции
var once sync.Once; once.Do(initFn)           // выполнить один раз

ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()                                // ВСЕГДА вызывайте cancel
```

---

## 9. Поведение `nil` и закрытых каналов

| Операция | `nil`-слайс | `nil`-мапа | `nil`-канал | Закрытый канал |
| :--- | :--- | :--- | :--- | :--- |
| чтение | паника при `s[i]` | вернёт Zero Value | блокируется навсегда | остаток буфера, затем Zero Value и `ok=false` |
| запись | паника при `s[i]=x` | **паника** | блокируется навсегда | **паника** |
| `len`, `range` | `0`, 0 итераций | `0`, 0 итераций | — | — |
| `append` / `delete` / `close` | `append` работает | `delete` безопасен | `close` — паника | повторный `close` — паника |

---

## 10. Тестирование

```go
func TestAdd(t *testing.T) {                       // файл *_test.go
    tests := []struct {
        name       string
        a, b, want int
    }{
        {"положительные", 2, 3, 5},
        {"с нулём", 0, 4, 4},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {        // подтесты
            if got := Add(tt.a, tt.b); got != tt.want {
                t.Errorf("Add(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
            }
        })
    }
}

func BenchmarkAdd(b *testing.B) {
    for b.Loop() { Add(1, 2) }                     // Go 1.24+; раньше: for i := 0; i < b.N; i++
}
```

`t.Errorf` — отметить ошибку и продолжить; `t.Fatalf` — отметить и остановить тест; `t.Helper()`, `t.Cleanup(fn)`, `t.TempDir()`, `t.Parallel()`.

---

## 11. Частые операции стандартной библиотеки

```go
fmt.Printf("%v %+v %#v %T %q %d %5.2f %x %p\n", v, v, v, v, s, n, f, n, &n)
fmt.Sprintf("id=%d", 42);  fmt.Errorf("...: %w", err)

strings.Split(s, ","); strings.Join(parts, "-"); strings.Contains(s, "x")
strings.TrimSpace(s);  strings.ToUpper(s);  strings.HasPrefix(s, "go")
strconv.Atoi("42");  strconv.Itoa(42);  strconv.ParseFloat("3.14", 64)

b, _ := json.Marshal(v);  json.Unmarshal(b, &v)
json.NewEncoder(w).Encode(v);  json.NewDecoder(r).Decode(&v)

resp, err := http.Get(url); defer resp.Body.Close()      // в реальном коде — свой http.Client с таймаутом
http.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")                               // Go 1.22+
})

os.ReadFile("f.txt");  os.WriteFile("f.txt", data, 0o644)
bufio.NewScanner(os.Stdin);  io.Copy(dst, src);  io.ReadAll(r)
time.Now();  time.Since(start);  t.Format("2006-01-02 15:04:05")   // референсное время Go!
```

---

## 12. Топ-ловушек для начинающих (контрольный список)

1. `:=` внутри вложенного блока **создаёт новую переменную** и затеняет внешнюю (типичный баг с `err`).
2. `append` **возвращает** новый слайс — результат надо присваивать; общий базовый массив у под-срезов может неожиданно меняться.
3. В `for ... range` значение — **копия** элемента: `for _, v := range s { v.X = 1 }` оригинал не меняет.
4. Запись в `nil`-мапу — паника; мапу создают через `make` или литерал.
5. Мапа и `sync.WaitGroup`/`sync.Mutex` **нельзя** безопасно использовать без синхронизации / копировать после первого использования.
6. `defer` внутри цикла откладывает вызов до конца **функции**, а не итерации.
7. `err != nil` может быть `true` для «пустого» указателя, завёрнутого в интерфейс `error` (ловушка «typed nil»).
8. `len(str)` — байты, а не символы; для текста с кириллицей/эмодзи используйте руны.
9. Горутина, заблокированная на канале, «утекает» — всегда продумывайте, кто и когда закрывает канал/отменяет контекст.
10. Деньги — не `float64`: используйте целые копейки (`int64`) или decimal-тип.
11. Не забывайте `defer cancel()` после `context.WithTimeout/WithCancel` и `defer resp.Body.Close()` после HTTP-запроса.
12. `http.Server` и `http.Client` без таймаутов опасны в продакшене.

---

## 13. Словарь терминов

| Термин | Что означает в Go |
| :--- | :--- |
| **Модуль (module)** | Единица версионирования кода: папка с файлом `go.mod` |
| **Пакет (package)** | Набор `.go`-файлов одной папки с общим именем пакета; единица повторного использования |
| **Экспортируемый идентификатор** | Имя с заглавной буквы — доступно из других пакетов (аналог `public`) |
| **Ресивер (receiver)** | Параметр метода перед его именем: `func (u *User) Name()` |
| **Zero Value** | Значение по умолчанию, которое получает переменная без явной инициализации |
| **Слайс (slice)** | Динамическое «окно» над массивом: `{указатель, длина len, ёмкость cap}` |
| **Руна (rune)** | Кодовая точка Unicode (`int32`), «символ» в широком смысле |
| **Горутина (goroutine)** | Легковесный поток выполнения, управляемый рантаймом Go (не поток ОС) |
| **Канал (channel)** | Типобезопасная очередь для обмена данными между горутинами |
| **Интерфейс (interface)** | Набор сигнатур методов; тип реализует его **неявно**, если имеет эти методы |
| **Встраивание (embedding)** | Включение одного типа в другой с «пробросом» полей и методов; **не** наследование |
| **Type assertion / type switch** | Проверка и извлечение конкретного типа из значения интерфейса |
| **Sentinel-ошибка** | Заранее объявленное значение ошибки (`var ErrX = errors.New(...)`), проверяемое через `errors.Is` |
| **Escape-анализ** | Решение компилятора, разместить переменную на стеке или в куче |
| **GC (сборщик мусора)** | Фоновый механизм освобождения памяти, на которую больше нет ссылок |
| **Data race (гонка данных)** | Одновременный доступ к одной памяти из разных горутин, где хотя бы одно обращение — запись, без синхронизации |
| **Идемпотентность** | Повторное выполнение операции даёт тот же результат, что и однократное |
