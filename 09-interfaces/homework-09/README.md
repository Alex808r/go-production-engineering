# Домашнее задание №9: Интерфейсы в Go

Данная директория содержит решение домашнего задания к занятию №9 курса **«Современная разработка на Go»** (Thinknetica).

---

## 📑 Содержание

1. [Задача №1: Полиморфизм через общий контракт интерфейса](#задача-1-полиморфизм-через-общий-контракт-интерфейса)
   - Контракт `Ager` и утиная типизация для `Employee` и `Customer`.
   - Функция поиска максимального возраста `MaxAge(people ...Ager) int`.
2. [Задача №2: Обобщенная обработка без методов через `any` и `Type Switch`](#задача-2-обобщенная-обработка-без-методов-через-any-и-type-switch)
   - Структуры без методов `EmployeeSimple` и `CustomerSimple`.
   - Извлечение данных и возврат самого объекта через `OldestPerson(people ...any) any`.
3. [Задача №3: Фильтрация типов и потоковый вывод в `io.Writer`](#задача-3-фильтрация-типов-и-потоковый-вывод-в-iowriter)
   - Функция `WriteOnlyStrings(w io.Writer, args ...any) (int, error)`.
   - Проверка типов с помощью `switch v := arg.(type)`.
4. [Инструкции по запуску тестов](#инструкции-по-запуску-тестов)

---

## Задача №1: Полиморфизм через общий контракт интерфейса

Реализовано в [`pkg/users/users.go`](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/09-interfaces/homework-09/pkg/users/users.go).

### Описание решения:
1. Объявлен интерфейс:
   ```go
   type Ager interface {
       Age() int
   }
   ```
2. Типы `Employee` и `Customer` имеют инкапсулированное поле `age int` и реализуют метод `Age() int`:
   ```go
   func (e Employee) Age() int { return e.age }
   func (c Customer) Age() int { return c.age }
   ```
3. Функция `MaxAge` принимает вариативный список `...Ager`, благодаря утиной типизации может принимать экземпляры обоих типов в произвольном порядке и возвращает максимальный возраст:
   ```go
   func MaxAge(people ...Ager) int {
       var maxAge int
       for _, p := range people {
           if p != nil && p.Age() > maxAge {
               maxAge = p.Age()
           }
       }
       return maxAge
   }
   ```

---

## Задача №2: Обобщенная обработка без методов через `any` и `Type Switch`

### Описание решения:
1. Типы `EmployeeSimple` и `CustomerSimple` не имеют методов и представляют собой простые структуры:
   ```go
   type EmployeeSimple struct { Name string; Age int }
   type CustomerSimple struct { Name string; Age int }
   ```
2. Функция `OldestPerson` принимает `...any` (пустой интерфейс), использует конструкцию **`switch v := p.(type)`** для безопасной проверки типов, находит человека с максимальным возрастом и возвращает сам объект:
   ```go
   func OldestPerson(people ...any) any {
       var oldest any
       maxAge := -1
       for _, p := range people {
           switch v := p.(type) {
           case EmployeeSimple:
               if v.Age > maxAge {
                   maxAge = v.Age
                   oldest = v
               }
           case CustomerSimple:
               if v.Age > maxAge {
                   maxAge = v.Age
                   oldest = v
               }
           }
       }
       return oldest
   }
   ```

---

## Задача №3: Фильтрация типов и потоковый вывод в `io.Writer`

### Описание решения:
Функция принимает интерфейс `io.Writer` и список аргументов произвольных типов `...any`. С помощью `switch v := arg.(type)` отфильтровываются только строковые аргументы (`string` и `*string`), которые затем записываются через `io.WriteString`:

```go
func WriteOnlyStrings(w io.Writer, args ...any) (int, error) {
    if w == nil {
        return 0, fmt.Errorf("writer cannot be nil")
    }
    var totalWritten int
    for _, arg := range args {
        switch v := arg.(type) {
        case string:
            n, err := io.WriteString(w, v)
            totalWritten += n
            if err != nil {
                return totalWritten, err
            }
        case *string:
            if v != nil {
                n, err := io.WriteString(w, *v)
                totalWritten += n
                if err != nil {
                    return totalWritten, err
                }
            }
        }
    }
    return totalWritten, nil
}
```

---

## 🚀 Инструкции по запуску тестов

Запуск тестов пакета с проверкой покрытия:
```bash
cd go_course_4
go test -v -cover ./09-interfaces/homework-09/pkg/users/...
```
