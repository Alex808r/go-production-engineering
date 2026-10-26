# 21-interview — Разбор каверзных задач и паттернов для технических собеседований

В данной директории собраны фундаментальные вопросы-ловушки, архитектурные паттерны многопоточности и практические задачи, регулярно встречающиеся на технических интервью по Go (уровни Middle / Senior).

---

## 📋 Структура каталога

```text
21-interview/
├── 01-basics.go       # Базовые ловушки: интерфейсы nil!=nil, defer, срезы, замыкания, UTF-8
├── 01-basics_test.go  # Тесты и бенчмарки преаллокации срезов
├── 02-channels.go     # Каналы и конкурентность: неблокирующий select, Producer-Consumer, Fan-In, WorkerPool
├── 02-channels_test.go # Тесты без дедлоков и бенчмарки многопоточных паттернов
└── 21-interview.md    # Данный справочник
```

---

## 🔍 Подробный разбор тем и ловушек

### 1. Ловушка типизированного `nil != nil` в интерфейсах ([`01-basics.go`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/21-interview/01-basics.go))
Интерфейс в рантайме Go представляется структурой `iface` (для непустых интерфейсов с методами) или `eface` (для `any` / `interface{}`), состоящей из двух машинных слов:
1. `tab` (или `_type`): указатель на тип данных и его таблицу методов.
2. `data`: указатель на конкретное значение в памяти.

> [!IMPORTANT]
> Переменная интерфейсного типа равна `nil` **ТОЛЬКО** тогда, когда оба поля (`type` и `data`) равны `nil`!

```go
type CustomError struct{ Msg string }
func (e *CustomError) Error() string { return e.Msg }

func Fail() error {
    var err *CustomError = nil
    return err // ВНИМАНИЕ: возвращает iface{type: *CustomError, data: nil}
}

func main() {
    err := Fail()
    if err != nil {
        // Условие ИСТИННО! Интерфейс НЕ равен nil, хотя err внутри nil.
        // Вызов err.Error() приведет к панике, если метод не защищен от nil-receiver!
    }
}
```
**Как исправить:** возвращать нетипизированный `nil`:
```go
if err == nil {
    return nil // Корректно: iface{type: nil, data: nil}
}
```

---

### 2. Механика `defer` и именованные возвращаемые значения ([`01-basics.go`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/21-interview/01-basics.go))
1. **Вычисление аргументов:** аргументы отложенного вызова вычисляются **в момент объявления `defer`**, а не в момент фактического выполнения.
2. **Стек LIFO:** отложенные вызовы выполняются в порядке, обратном их объявлению.
3. **Модификация результатов:** анонимная функция в `defer` выполняется **после** инструкции `return`, но **до** выхода из вызывающей функции. Она имеет доступ к именованным возвращаемым значениям и может их перезаписать:
```go
func DeferWithNamedReturn() (res int) {
    defer func() {
        res *= 2 // изменяет результат на 20
    }()
    return 10
}
```

---

### 3. Мутации и реалокация срезов ([`01-basics.go`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/21-interview/01-basics.go))
Срез в Go — это заголовок (`SliceHeader`): указатель `Data *T`, длина `Len int` и вместимость `Cap int`.
- При взятии подслайса `s2 := s1[1:3]` создается новый заголовок, но указатель `Data` ссылается на **тот же базовый массив**.
- Мутация `s2[1] = 999` изменяет элемент исходного массива, поэтому значение меняется и в `s1[2]`.
- Если вызов `append` превышает текущую емкость `Cap`, рантайм Go аллоцирует **новый массив в куче**, копирует данные и возвращает срез с указателем на новый массив. Начиная с этого момента `s1` и `s2` изолированы.

#### Бенчмарк преаллокации (`make([]T, 0, cap)` vs `append` без капа):
```text
BenchmarkSliceAppend_NoPrealloc-8       81253     12489 ns/op   (множественные реаллокации в куче)
BenchmarkSliceAppend_WithPrealloc-8    363015      3300 ns/op   (~4x быстрее, ровно 1 аллокация!)
```

---

### 4. Строки, UTF-8 и руны ([`01-basics.go`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/21-interview/01-basics.go))
- `string` — неизменяемый срез байт (`[]byte`).
- Символы кириллицы в UTF-8 кодируются **2 байтами**. Функция `len("Привет")` возвращает `12`, а не `6`!
- Цикл `for idx, r := range "Привет"` декодирует UTF-8 поток в руны (`rune` = `int32`), при этом `idx` смещается на 2 байта на каждый символ.

---

### 5. Сводная матрица состояний каналов

Знание этой таблицы проверяется практически на каждом собеседовании:

| Состояние канала | Чтение (`<-ch`) | Запись (`ch <- val`) | Закрытие (`close(ch)`) |
| :--- | :--- | :--- | :--- |
| **`nil`** (`var ch chan T`) | Вечная блокировка (deadlock) | Вечная блокировка (deadlock) | **`panic: close of nil channel`** |
| **Открыт, пуст** | Блокировка (до появления данных) | Запись в буфер / передача | Закрывается успешно |
| **Открыт, заполнен** | Чтение элемента | Блокировка (до освобождения буфера) | Закрывается успешно |
| **Закрыт** (`closed`) | Считывает остатки буфера, затем zero-value (`val, ok := <-ch`, `ok==false`) | **`panic: send on closed channel`** | **`panic: close of closed channel`** |

---

### 6. Паттерны многопоточности ([`02-channels.go`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/21-interview/02-channels.go))

#### А. Неблокирующий select (`Select`)
Использование ветки `default` предотвращает блокировку текущей горутины, если канал не готов к отправке или приему:
```go
select {
case ch <- "event":
    // успешно отправлено
default:
    // канал переполнен или читателя нет, действие пропускается
}
```

#### Б. Поставщик-Потребитель (`ProduceConsume`)
Поставщик генерирует данные и обязан закрыть канал через `defer close(ch)`. Потребитель читает через `for val := range ch`. Если поставщик забудет закрыть канал, цикл `range` потребителя навсегда заблокируется в ожидании новых элементов, приводя к утечке горутины (**goroutine leak**).

#### В. Fan-In: Мультиплексирование потоков (`FanIn`)
Объединение $N$ независимых каналов в один результирующий поток. Канал закрывается строго после того, как все поставщики завершат работу (`sync.WaitGroup`).

#### Г. Пул воркеров (`WorkerPool`)
Пул из $K$ горутин, параллельно обрабатывающих задачи из канала `jobs` и пишущих результаты в `results`.

---

### 7. Внутреннее устройство `map`: коллизии и инкрементальная эвакуация
- **Коллизии:** 1-го рода (хэш-функции $h(K_1) == h(K_2)$) и 2-го рода (слота $h \pmod{2^B}$). Внутри бакета (до 8 пар) ячейки отбираются по `tophash [8]uint8`. При заполнении бакета рантайм связывает `overflow`-бакеты (**Separate Chaining** в Go $\le$ 1.23) или использует **Swiss Tables** (**Open Addressing** в Go 1.24+). Защита от DoS обеспечивается случайным зерном `hash0`.
- **Эвакуация (Evacuation):** Плавный перенос ключей из `oldbuckets` в `buckets`.
  - *Триггеры:* **Load Factor > 6.5** (удвоение бакетов $2^{B+1}$) или **избыток `overflow`-бакетов** при низком LF (дефрагментация $2^B$).
  - *Инкрементальность:* Переносится ровно по **2 бакета** на операцию записи/удаления для устранения пауз (STW).
  - *Чтение и адресация:* Чтение `m[k]` не эвакуирует бакеты (lock-free). Из-за смещения физических адресов в RAM элементы мапы **неадресуемы** (`&m[k]` запрещён).  
  *(Подробный разбор: см. [Том 1.2, разд. 6.7](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/go-documentation/01-go-core-data-structures.md#67-внутреннее-устройство-map-в-go-hmap-bmap-коллизии-и-swiss-tables) и [Том 6, вопросы 26–27](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/go-documentation/06-go-core-algorithms-interview.md#361-топ-27-ключевых-вопросов-и-ответов-технического-собеседования))*.

---

## 🚀 Запуск и валидация

```bash
# Запуск всех тестов с проверкой race-condition и бенчмарками
go test -v -race -bench=. ./lessons/21-interview/...
```

---

## 🔗 Связанные материалы
- **Теория:** [Том 6: Алгоритмы и собеседования (Разделы 32–36)](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/go-documentation/06-go-core-algorithms-interview.md)
- **Лайвкодинг:** [go-livecoding-guide.md](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/go-documentation/go-livecoding-guide.md)
- **Шпаргалка:** [go-core-cheatsheet.md](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/go-documentation/go-core-cheatsheet.md)
- **Базы данных на собеседовании:** [go-database-interview-guide.md](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/go-documentation/go-database-interview-guide.md)
