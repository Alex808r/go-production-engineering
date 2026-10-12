# 09-interfaces — Интерфейсы, полиморфизм, утверждения типов и дженерики

В данной директории детально исследуются интерфейсы в Go: утиная типизация, внутренняя структура `iface` и `eface`, коварство `nil`-интерфейсов, утверждения типов (`Type Assertions`), переключатели типов (`Type Switches`) и параметрический полиморфизм (Generics, Go 1.18+).

---

## 📋 Структура каталога

```text
09-interfaces/
├── 1-polymorph/    # Полиморфизм: обработка различных структур через общий интерфейс
├── 2-quiz/         # Классическая ловушка: "nil interface" vs "interface containing nil pointer"
├── 3-assertion/    # Утверждение типов (val, ok := i.(MyType)) и наборы методов (Method Sets)
├── 4-type_switch/  # Переключатель типов по switch v := i.(type)
├── 5-important/    # Каталог ключевых интерфейсов стандартной библиотеки Go
├── 6-generics/     # Обобщенное программирование (type parameters, ~, cmp.Ordered, slices)
├── 7-DIP/          # Архитектурное проектирование через интерфейсы и инверсию зависимостей
└── homework-09/    # Домашнее задание 09 (полиморфизм Ager, type switch, io.Writer)
```

---

## 🔍 Ключевые концепции

1. **Внутреннее устройство интерфейса:**
   - Интерфейс состоит из двух машинных слов: указателя на таблицу типов (`itab`/`_type`) и указателя на сами данные (`data`).
   - Значение интерфейса равно `nil` **только тогда**, когда оба указателя (`type` и `data`) равны `nil`.
   - Если переменной интерфейса присвоить типизированный указатель `(*MyStruct)(nil)`, интерфейс **НЕ равен `nil`**! Это частый источник трудноуловимых багов на собеседованиях (`2-quiz`).

2. **Безопасное приведение типов (`3-assertion`):**
   - Всегда используйте безопасную форму с двумя значениями: `s, ok := iface.(string)`. Форма с одним значением `s := iface.(string)` вызывает панику при несоответствии типов.
   - Учитывайте наборы методов (Method Sets): метод с `pointer receiver` не входит в набор методов значения!

3. **Дженерики (`6-generics`):**
   - Параметризация функций и структур типами без потери производительности и без накладных расходов `any`:
     ```go
     func min[T ordered](in []T) (T, bool)
     ```

---

## 🚀 Запуск и проверка

```bash
# Запуск примеров с полиморфизмом
go run ./09-interfaces/1-polymorph/polymorph.go

# Запуск разбора задачи с quiz (nil interface vs typed nil)
go run ./09-interfaces/2-quiz/quiz.go

# Запуск примеров с утверждением типов (type assertion)
go run ./09-interfaces/3-assertion/assertion.go

# Запуск примеров с type switch
go run ./09-interfaces/4-type_switch/switch.go

# Запуск примеров с дженериками
go run ./09-interfaces/6-generics/generics.go

# Запуск тестов модуля и домашнего задания
go test -v -cover ./09-interfaces/...
```

---

## 🔗 Связанные материалы
- **Теория:** [Том 2: ООП и методы (Разделы 9, 15: Полиморфизм, дженерики, empty interface)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/02-go-core-oop-methods.md)
- **Практическое ДЗ:** [homework-09](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/09-interfaces/homework-09) (Разработка интерфейсов и реализаций для пакета пользователей)

