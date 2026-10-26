# Домашнее задание №8: Профилирование, Отладка и Трассировка в Go

Данная директория содержит полное решение домашнего задания к занятию №8 курса **«Современная разработка на Go»** (Thinknetica).

---

## 📑 Содержание

1. [Задача №1: Отладка приложения с помощью Delve (`dlv`)](#задача-1-отладка-приложения-с-помощью-delve-dlv)
   - Воспроизведение и анализ бага в задаче Two Sum.
   - Сценарий интерактивной отладки в Delve (точки останова, условные брейкпоинты, стек, инспекция переменных).
   - Настройка отладки в VS Code (`launch.json`).
   - Исправление ошибки: наивный $O(N^2)$ и оптимальный $O(N)$ алгоритм на хеш-таблице.
   - Сравнительные бенчмарки решений (ускорение в 14 раз).
2. [Задача №2: Профилирование процессора и памяти (`pprof`)](#задача-2-профилирование-процессора-и-памяти-pprof)
   - Снятие профилей бенчмарков (`-cpuprofile`, `-memprofile`).
   - Анализ CPU-профиля алгоритмов поиска (`Binary` vs `Linear`).
   - Профилирование живого приложения через HTTP (`net/http/pprof`).
   - Интерпретация графа вызовов (цвета, размеры блоков, толщина стрелок, `cum` vs `flat`).
   - Анализ профиля кучи: `inuse_space` vs `alloc_space` и локализация утечек.
3. [Задача №3: Трассировка выполнения приложения (`runtime/trace` & `go tool trace`)](#задача-3-трассировка-выполнения-приложения-runtimetrace--go-tool-trace)
   - Инструментирование кода с помощью `runtime/trace`.
   - Запуск визуализатора `go tool trace trace.out`.
   - Анализ распределения горутин по логическим процессорам `P` и потокам ОС `M`.
   - Мониторинг фаз сборщика мусора (Concurrent Mark, STW, GC Assist) и выявление GC Pressure.
4. [Инструкции по запуску тестов, профилирования и трассировки](#инструкции-по-запуску-тестов-профилирования-и-трассировки)

---

## Задача №1: Отладка приложения с помощью Delve (`dlv`)

Реализовано в пакете [`pkg/twosum`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/08-prof_debug/homework-08/pkg/twosum).

### 1. Анализ проблемы в Two Sum
В исходном примере из урока функция `TwoSum` имеет коварный дефект:
```go
func TwoSumBuggy(nums []int, target int) [2]int {
    for i := 0; i < len(nums); i++ {
        for j := 0; j < len(nums); j++ {
            if nums[i]+nums[j] == target {
                return [2]int{i, j}
            }
        }
    }
    return [2]int{-1, -1}
}
```
- **Симптом**: На тестовом примере `nums = [2, 7, 11, 15], target = 13` алгоритм возвращает правильный ответ `[0, 2]`.
- **Ошибка**: На данных `nums = [2, 10, 6, 12, 14], target = 20` алгоритм находит `nums[1] + nums[1] == 10 + 10 == 20` и возвращает `[1, 1]`. Однако по условию задачи **запрещено использовать один и тот же элемент дважды**. Правильный ответ: `6 + 14 = 20` (индексы `[2, 4]`).

### 2. Сценарий отладочной сессии в Delve CLI
Для локализации дефекта используется специализированный отладчик Go — **Delve**:

```bash
# Запуск отладочной сессии тестов или main:
$ dlv test ./lessons/08-prof_debug/homework-08/pkg/twosum
Type 'help' for list of commands.
(dlv) 
```

**Ключевые команды отладчика:**
1. **Установка точки останова с условием (Conditional Breakpoint)**:
   Останавливаемся только тогда, когда алгоритм находит пару, равную `target`, но при этом `i == j`:
   ```text
   (dlv) break twosum.go:10
   Breakpoint 1 set at 0x... for go-core-4/homework-08/pkg/twosum.TwoSumBuggy() twosum.go:10
   (dlv) condition 1 i == j && nums[i] + nums[j] == target
   ```
2. **Продолжение выполнения до срабатывания условия**:
   ```text
   (dlv) continue
   > go-core-4/homework-08/pkg/twosum.TwoSumBuggy() twosum.go:10 (hits goroutine(7):1 total:1)
        5: func TwoSumBuggy(nums []int, target int) [2]int {
        6: 	for i := 0; i < len(nums); i++ {
        7: 		for j := 0; j < len(nums); j++ {
        8: 			if nums[i]+nums[j] == target {
   =>   9: 				return [2]int{i, j}
       10: 			}
       11: 		}
   ```
3. **Инспекция локальных переменных (`locals`) и выражений (`print`)**:
   ```text
   (dlv) locals
   i = 1
   j = 1
   (dlv) print nums[i]
   10
   (dlv) print nums[j]
   10
   (dlv) print target
   20
   ```
   *Вывод отладчика*: Срез содержит элемент `10` на индексе 1, и вложенный цикл взял его самого с собой, нарушив контракт задачи.
4. **Просмотр стека вызовов (`stack`) и фреймов (`frame`)**:
   ```text
   (dlv) stack
   0  0x0000000100f8e910 in go-core-4/homework-08/pkg/twosum.TwoSumBuggy
      at twosum.go:9
   1  0x0000000100f8ed48 in go-core-4/homework-08/pkg/twosum.TestTwoSumBuggy_FailsOnDuplicateElement
      at twosum_test.go:18
   2  0x0000000100f28a1c in testing.tRunner
      at src/testing/testing.go:1689
   ```
5. **Просмотр горутин (`goroutines`)**:
   ```text
   (dlv) goroutines
     Goroutine 1 - User: ...
   * Goroutine 7 - User: go-core-4/homework-08/pkg/twosum.TwoSumBuggy twosum.go:9 (thread 12345)
   ```

### 3. Интеграция с VS Code (`launch.json`)
Для графической отладки в VS Code создается конфигурация:
```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Delve: Debug TwoSum Test",
            "type": "go",
            "request": "launch",
            "mode": "test",
            "program": "${workspaceFolder}/homework-08/pkg/twosum",
            "args": ["-test.run", "TestTwoSumBuggy"]
        }
    ]
}
```
В графическом интерфейсе доступны панели:
- **Variables**: отображение локальных переменных `nums`, `target`, `i`, `j`.
- **Watch**: контрольные выражения (например, `i != j`, `nums[i] + nums[j]`).
- **Call Stack**: навигация по кадрам стека вызовов.
- **Breakpoints**: точки останова с поддержкой Condition (выражение) и Hit Count.

### 4. Исправление дефекта и алгоритмическая оптимизация
В [`pkg/twosum/twosum.go`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/08-prof_debug/homework-08/pkg/twosum/twosum.go) реализованы два корректных решения:

1. **Исправленный Brute-Force ($O(N^2)$ время, $O(1)$ память)**:
   Начинаем внутренний цикл с `j = i + 1`:
   ```go
   func TwoSumBruteForce(nums []int, target int) [2]int {
       for i := 0; i < len(nums); i++ {
           for j := i + 1; j < len(nums); j++ {
               if nums[i]+nums[j] == target {
                   return [2]int{i, j}
               }
           }
       }
       return [2]int{-1, -1}
   }
   ```
2. **Оптимальный алгоритм на хеш-таблице ($O(N)$ время, $O(N)$ память)**:
   Однопроходный поиск через сохранение дополнения `target - num` в `map[int]int`:
   ```go
   func TwoSumOptimal(nums []int, target int) [2]int {
       seen := make(map[int]int, len(nums))
       for i, num := range nums {
           diff := target - num
           if idx, ok := seen[diff]; ok {
               return [2]int{idx, i}
           }
           seen[num] = i
       }
       return [2]int{-1, -1}
   }
   ```

### 5. Результаты бенчмарков ($N=1\,000$ элементов):
```text
BenchmarkTwoSumBruteForce_1000-8         5846       205892 ns/op
BenchmarkTwoSumOptimal_1000-8           81206        14813 ns/op
```
> 🚀 **Вывод**: Оптимальный алгоритм $O(N)$ работает в **14 раз быстрее** квадратичного перебора на 1 000 элементов (и разрыв растет квадратично с увеличением $N$).

---

## Задача №2: Профилирование процессора и памяти (`pprof`)

### 1. Профилирование бенчмарков
Для пакета [`pkg/search`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/08-prof_debug/homework-08/pkg/search) запущен бенчмарк со сбором профилей:
```bash
go test -bench=BenchmarkBinarySearch_100k -benchmem -cpuprofile=cpu.out -memprofile=mem.out ./lessons/08-prof_debug/homework-08/pkg/search
```
Результат замера:
```text
BenchmarkBinarySearch_100k-8   24036921   48.54 ns/op   0 B/op   0 allocs/op
```

### 2. Анализ профиля процессора (`cpu.out`)
Анализ профиля с помощью `go tool pprof -text cpu.out`:
```text
Type: cpu
Duration: 2.22s, Total samples = 1.79s (80.65%)
Showing nodes accounting for 1.79s, 100% of 1.79s total
      flat  flat%   sum%        cum   cum%
     1.55s 86.59% 86.59%      1.60s 89.39%  go-core-4/homework-08/pkg/search.Binary
     0.15s  8.38% 94.97%      1.76s 98.32%  go-core-4/homework-08/pkg/search.BenchmarkBinarySearch_100k
     0.05s  2.79% 97.77%      0.05s  2.79%  runtime.asyncPreempt
```

**Интерпретация метрик:**
- **`flat`**: Время, проведенное непосредственно в теле функции (без учета вызываемых ею других функций). `Binary` занял **1.55s (86.59%)** чистого времени CPU.
- **`cum` (cumulative)**: Суммарное время функции вместе со всеми дочерними вызовами. Бенчмарк `BenchmarkBinarySearch_100k` имеет `cum = 1.76s (98.32%)`.
- **Аллокации**: `0 B/op, 0 allocs/op` — бинарный поиск не производит аллокаций в куче, все переменные (`low`, `high`, `mid`) аллоцируются на стеке.

### 3. Профилирование живого сервиса (`cmd/app_profile`)
В [`cmd/app_profile/main.go`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/08-prof_debug/homework-08/cmd/app_profile/main.go) реализован HTTP-сервер с подключением `_ "net/http/pprof"`:

```bash
# Сбор CPU профиля за 10 секунд:
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=10

# Инспекция кучи (память):
go tool pprof http://localhost:6060/debug/pprof/heap
```

#### Правила чтения визуального графа pprof:
1. **Размер прямоугольника**: Чем больше размер узла, тем больше ресурсов (`flat`) потребляет именно эта функция.
2. **Цвет узла**: От белого к ярко-красному. Красный узел указывает на критическую «горячую точку» (hotspot).
3. **Толщина стрелки**: Показывает относительный объем работы или объем памяти, переданный по данному пути вызова.
4. **Типы профилей кучи**:
   - `inuse_space`: объем памяти, занятый живыми объектами прямо сейчас (инструмент поиска утечек памяти).
   - `alloc_space`: суммарный объем памяти, выделенный за все время работы (инструмент поиска аллокаций, создающих паразитное давление на сборщик мусора GC).

---

## Задача №3: Трассировка выполнения приложения (`runtime/trace` & `go tool trace`)

Реализовано в [`cmd/app_trace/main.go`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/08-prof_debug/homework-08/cmd/app_trace/main.go).

### 1. Сбор трассировки
```go
f, _ := os.Create("trace.out")
trace.Start(f)
defer trace.Stop()
```
Генерируются три различные фазы работы системы:
- `CPU_Computation_Phase`: параллельные вычисления с распределением по ядрам процессора (`GOMAXPROCS`).
- `Memory_GC_Pressure_Phase`: частые аллокации и принудительные вызовы `runtime.GC()`.
- `Channel_Synchronization_Phase`: передача сообщений по буферизованному каналу.

### 2. Анализ трассировки через `go tool trace trace.out`
Инструмент запускает встроенный веб-сервер (на движке Chrome Trace Viewer):

```bash
go tool trace trace.out
```

**Что видно на временной шкале (Timeline):**
1. **Секция `STATS`**:
   - **`Goroutines`**: График количества одновременно существующих горутин (всплеск при запуске пула воркеров и возврат к baseline).
   - **`Heap`**: Пилообразный график потребления кучи: рост при выделении срезов и резкий спад после срабатывания GC.
   - **`Threads`**: Количество реальных потоков операционной системы (M), удерживаемых рантаймом Go.
2. **Секция `PROCS` (логические процессоры GMP: `Proc 0`, `Proc 1`, ...)**:
   - Отображает, какая именно горутина (`G1`, `G7`, `G8`...) выполнялась на каждом виртуальном процессоре в каждый момент времени.
   - Наглядно видно параллельное выполнение горутин на всех доступных ядрах без простоев.
3. **События GC (Garbage Collection)**:
   - Синие и красные полосы: фазы `GC (mark)`, `GC (dedicated)`, `GC (idle)` и `GC (assist)`.
   - Если горутина тратит слишком много времени в `GC (assist)`, это означает, что она выделяет память быстрее, чем сборщик успевает её очищать.
4. **Блокировки и планировщик**:
   - Отслеживание задержек планировщика (время между переходом горутины в состояние `Runnable` и фактическим стартом на `P`).

---

## 🚀 Инструкции по запуску тестов, профилирования и трассировки

### 1. Запуск всех тестов и бенчмарков TwoSum:
```bash
cd go-production-engineering
go test -v -bench=. ./lessons/08-prof_debug/homework-08/pkg/twosum/...
```

### 2. Запуск бенчмарков поиска и генерация профилей pprof:
```bash
cd go-production-engineering
go test -bench=BenchmarkBinarySearch_100k -benchmem -cpuprofile=homework-08/cpu.out -memprofile=homework-08/mem.out ./lessons/08-prof_debug/homework-08/pkg/search/...
go tool pprof -text ./lessons/08-prof_debug/homework-08/cpu.out
```

### 3. Запуск сервиса профилирования:
```bash
cd go-production-engineering
go run ./lessons/08-prof_debug/homework-08/cmd/app_profile/main.go
# В отдельном терминале:
go tool pprof -text http://localhost:6060/debug/pprof/profile?seconds=5
```

### 4. Генерация и просмотр трассировки:
```bash
cd go-production-engineering
go run ./lessons/08-prof_debug/homework-08/cmd/app_trace/main.go
go tool trace ./lessons/08-prof_debug/homework-08/trace.out
```
