# Практический задачник и разбор домашних заданий по Go (Goflex Project)

Данный документ содержит все домашние задания из 3-часового курса Антона Назарова («Быстрый старт в Go с нуля») по разработке backend-сервиса онлайн-кинотеатра **Goflex**.

> [!NOTE]
> **Как работать с этим задачником.** Задания расположены по нарастанию сложности (Уроки 1–12 → Занятия 2–7 → Домашние задания 8–20). Перед чтением решения попробуйте выполнить задание сами. Код решений лежит в каталогах `homework-NN/` рядом с этой документацией (все перечисленные пакеты компилируются и проходят тесты командой `go test ./homework-NN/...` из корня учебного модуля). Для запуска потребуется установленный Go (см. Том 1, раздел 0.8) и — для Уроков 12, 16, 19, 20 — Docker (см. Том 4, раздел 12.9).

Для каждой задачи приведены:

1. 🎯 **Условие задачи** (что требуется разработать).
2. 💡 **Ключевые концепции и архитектурные ловушки**.
3. 💻 **Эталонное решение (Clean Code / Production-ready)**.
4. 🔍 **Пошаговый разбор под капотом (Memory Model & Runtime)**.

> [!TIP]
> 📚 **Связанные материалы и справочники:**
>
> - 📖 **[README.md (Главное оглавление базы знаний по Go)](README.md)** — глубокий разбор языка, рантайма, структур данных и конкурентности в 6 томах; краткая справка — [go-core-cheatsheet.md](go-core-cheatsheet.md).
> - 🗄️ **[go-database-interview-guide.md (Гайд по Базам Данных)](go-database-interview-guide.md)** — PostgreSQL, `pgxpool`, транзакции ACID, Redis, миграции и 30 вопросов.
> - 💻 **[go-livecoding-guide.md (Live-Coding на Go)](go-livecoding-guide.md)** — 47 задач по многопоточности, структурам данных и алгоритмам в Go (Junior, Middle, Senior).
> - 🗄️ **[sql-livecoding-guide.md (Live SQL Задачник)](sql-livecoding-guide.md)** — 28 практических задач на SQL-запросы, JOINs, EXPLAIN ANALYZE и оконные функции.
> - 🏛️ **[go-system-design-guide.md (System Design для Go)](go-system-design-guide.md)** — 4-шаговый фреймворк, расчеты и 4 кейса архитектуры.
> - 🗺️ **[course-codebase-guide.md (Путеводитель по кодовой базе и примерам курса)](course-codebase-guide.md)** — сквозная карта курса Thinknetica: 22 лекционных модуля (00–21), проекты GoSearch и Lynks, навигатор по домашним заданиям (homework-02–20).

---

## 📑 Оглавление задач

1. [Урок 1: Первая программа и базовый ввод-вывод (`fmt`)](#урок-1-первая-программа-и-базовый-ввод-вывод-fmt)
2. [Урок 2: Переменные, константы и калькулятор](#урок-2-переменные-константы-и-калькулятор)
3. [Урок 3: Приведение типов и расчет рейтинга](#урок-3-приведение-типов-и-расчет-рейтинга)
4. [Урок 4: Условия (if с областью видимости и switch)](#урок-4-условия-if-с-областью-видимости-и-switch)
5. [Урок 5: Циклы for, break и continue](#урок-5-циклы-for-break-и-continue)
6. [Урок 6: Слайсы (Slices) и управление каталогом](#урок-6-слайсы-slices-и-управление-каталогом)
7. [Урок 7: Словари (Maps) и быстрый поиск тайтлов](#урок-7-словари-maps-и-быстрый-поиск-тайтлов)
8. [Урок 8: Структуры, указатели и методы (Бизнес-модель Goflex)](#урок-8-структуры-указатели-и-методы-бизнес-модель-goflex)
9. [Урок 9: Интерфейсы, полиморфизм и Type Switch](#урок-9-интерфейсы-полиморфизм-и-type-switch)
10. [Урок 10: Конкурентность (Горутины, Каналы, sync, Worker Pool, Context)](#урок-10-конкурентность-горутины-каналы-sync-worker-pool-context)
11. [Урок 11: Web HTTP REST API и Middleware](#урок-11-web-http-rest-api-и-middleware)
12. [Урок 12: PostgreSQL, pgxpool, ACID-транзакции и Миграции](#урок-12-postgresql-pgxpool-acid-транзакции-и-миграции)
13. [Упражнение (A Tour of Go): Квадратный корень методом Ньютона (Циклы и функции)](#упражнение-a-tour-of-go-квадратный-корень-методом-ньютона-циклы-и-функции)
14. [Занятие 2: Поисковый робот GoSearch (Пакет crawler и консольный поиск)](#занятие-2-поисковый-робот-gosearch-пакет-crawler-и-консольный-поиск)
15. [Занятие 3: Быстрый поисковый индекс (Инвертированный индекс и бинарный поиск)](#занятие-3-быстрый-поисковый-индекс-инвертированный-индекс-и-бинарный-поиск)
16. [Занятие 4: Структуры данных — Циклический двусвязный список со сторожевым элементом (Pop и Reverse)](#занятие-4-структуры-данных--циклический-двусвязный-список-со-сторожевым-элементом-pop-и-reverse)
17. [Занятие 5: Ввод-вывод в Go — Персистентность поисковых данных (io.Reader, io.Writer, JSON/GOB кеширование)](#занятие-5-ввод-вывод-в-go--персистентность-поисковых-данных-ioreader-iowriter-jsongob-кеширование)
18. [Занятие 6: ООП в Go — Структуры, методы и идиоматичный рефакторинг](#занятие-6-ооп-в-go--структуры-методы-и-идиоматичный-рефакторинг)
19. [Занятие 7: Тестирование в Go — Unit-тесты, Табличные тесты, Бенчмарки и Имитация зависимостей](#занятие-7-тестирование-в-go--unit-тесты-табличные-тесты-бенчмарки-и-имитация-зависимостей)
20. [Домашнее задание 8: Профилирование, Отладка и Трассировка в Go](#домашнее-задание-8-профилирование-отладка-и-трассировка-в-go)
21. [Домашнее задание 9: Интерфейсы в Go](#домашнее-задание-9-интерфейсы-в-go)
22. [Домашнее задание 10: Конкурентное программирование (Игра в Пинг-Понг)](#домашнее-задание-10-конкурентное-программирование-игра-в-пинг-понг)
23. [Домашнее задание 11: Сетевое программирование (Сетевая служба GoSearch)](#домашнее-задание-11-сетевое-программирование-сетевая-служба-gosearch)
24. [Домашнее задание 12: Веб-приложения на Go (Веб-служба GoSearch)](#домашнее-задание-12-веб-приложения-на-go-веб-служба-gosearch)
25. [Домашнее задание 13: Разработка REST API (Поисковик GoSearch API и Модель памяти)](#домашнее-задание-13-разработка-rest-api-поисковик-gosearch-api-и-модель-памяти)
26. [Домашнее задание 14: Удалённый вызов процедур (RPC-служба сообщений)](#домашнее-задание-14-удалённый-вызов-процедур-rpc-служба-сообщений)
27. [Домашнее задание 15: Реляционные базы данных (Схема БД онлайн-кинотеатра)](#домашнее-задание-15-реляционные-базы-данных-схема-бд-онлайн-кинотеатра)
28. [Домашнее задание 16: Приложения с базами данных (Пакет для БД фильмов с паттерном Repository)](#домашнее-задание-16-приложения-с-базами-данных-пакет-для-бд-фильмов-с-паттерном-repository)
29. [Домашнее задание 17: Архитектура Go-приложения (Clean Architecture и SOLID)](#домашнее-задание-17-архитектура-go-приложения-clean-architecture-и-solid)
30. [Домашнее задание 18: Разработка микросервисов (URL Shortener Service)](#домашнее-задание-18-разработка-микросервисов-url-shortener-service)
31. [Домашнее задание 19: Очереди сообщений и Асинхронная аналитика (Kafka / Event Sourcing)](#домашнее-задание-19-очереди-сообщений-и-асинхронная-аналитика-kafka--event-sourcing)
32. [Домашнее задание 20: Итоговый проект «Lynks» (NoSQL, Redis, PostgreSQL)](#домашнее-задание-20-итоговый-проект-lynks-nosql-redis-postgresql)

---

## Урок 1: Первая программа и базовый ввод-вывод (`fmt`)

> 📖 **Теория:** [Том 1: Разделы 0.1, 0.7 (Старт с нуля и fmt)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/01-go-core-basics.md#01-структура-программы-package-main-func-main-и-первая-программа)  
> 💻 **Код лекции:** [01-intro/hello](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/01-intro/hello)  
> 🗺️ **Путеводитель:** [course-codebase-guide.md (Урок 01)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/course-codebase-guide.md#урок-01-первая-программа-и-структура-приложения)

### 🎯 Условие:

Создайте программу `main.go`, которая выводит в консоль визитную карточку разработчика сервиса Goflex:

1. Имя разработчика и текущую позицию через `fmt.Println`.
2. Опыт в месяцах, целевой стек и ожидаемый оклад с использованием форматированного вывода `fmt.Printf` (спецификаторы `%s`, `%d`, `%f`, `%T`).

### 💻 Эталонное решение:

```go
package main

import "fmt"

func main() {
    // 1. Вывод простых строк с переносом:
    fmt.Println("=== Сервис Goflex: Карточка Инженера ===")
    fmt.Println("Разработчик: Антон")
    fmt.Println("Специализация: Golang Backend Developer")

    // 2. Форматированный вывод данных:
    experienceMonths := 6
    targetSalary := 250000.50
    stack := "Go, PostgreSQL, Docker, Redis"

    fmt.Printf("Опыт обучения: %d мес.\n", experienceMonths)
    fmt.Printf("Ожидаемый оклад: %.2f руб.\n", targetSalary)
    fmt.Printf("Стек: %s (Тип переменной: %T)\n", stack, stack)
}
```

### 🔍 Разбор:

- Пакет `package main` сообщает компилятору, что создается исполняемый бинарник, а не переиспользуемая библиотека.
- Функция `fmt.Println` автоматически разделяет аргументы пробелами и ставит символ `\n` в конце.
- `fmt.Printf` позволяет форматировать данные: `%.2f` ограничивает число с плавающей точкой 2 знаками после запятой, `%T` выводит внутренний тип переменной.

---

## Урок 2: Переменные, константы и калькулятор

> 📖 **Теория:** [Том 1: Разделы 0.2, 0.3, 3.1 (Переменные, Zero Value, константы)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/01-go-core-basics.md#02-переменные-в-go-явное-var-vs-короткое-объявление)  
> 💻 **Код лекции:** [02-syntax/1-basic](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/02-syntax/1-basic)  
> 🗺️ **Путеводитель:** [course-codebase-guide.md (Урок 02)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/course-codebase-guide.md#урок-02-фундаментальный-синтаксис-и-модель-памяти)

### 🎯 Условие:

Напишите калькулятор стоимости подписки для пользователя Goflex:

- Объявите константу базовой стоимости подписки на 1 месяц.
- Объявите переменные для количества месяцев и процента скидки.
- Рассчитайте итоговую стоимость и проверьте дефолтные нулевые значения (`zero value`) для неинициализированных переменных.

### 💻 Эталонное решение:

```go
package main

import "fmt"

const BaseMonthlyPrice = 499 // нетипизированная константа

func main() {
    var (
        months          int     = 6
        discountPercent float64 = 0.15 // 15% скидка
        promoCodeApplied bool          // zero value: false
    )

    // Расчет без скидки:
    totalBeforeDiscount := float64(months * BaseMonthlyPrice)
    // Расчет со скидкой:
    totalAfterDiscount := totalBeforeDiscount * (1.0 - discountPercent)

    fmt.Printf("Количество месяцев: %d\n", months)
    fmt.Printf("Базовая цена: %d руб/мес\n", BaseMonthlyPrice)
    fmt.Printf("Применен промокод: %t (дефолтное значение)\n", promoCodeApplied)
    fmt.Printf("Итого к оплате со скидкой 15%%: %.2f руб.\n", totalAfterDiscount)
}
```

### 🔍 Разбор:

- В Go константы (`const`) вычисляются во время компиляции. Если тип не указан явно (`untyped constant`), она может автоматически использоваться в выражениях совместимых типов без явного каста.
- Все неинициализированные переменные получают гарантированный `Zero Value` (`0`, `0.0`, `""`, `false`, `nil`), предотвращая чтение мусора из неинициализированной памяти (в отличие от C/C++).
- ⚠️ В учебном примере цена и скидка хранятся в `float64` ради простоты. В реальных финансовых расчётах так делать нельзя (ошибки округления): храните деньги в целых копейках (`int64`) или в десятичном типе — см. Том 1, раздел 0.3.

---

## Урок 3: Приведение типов и расчет рейтинга

> 📖 **Теория:** [Том 1: Раздел 0.4 (Строгая типизация и конверсия типов)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/01-go-core-basics.md#04-строгая-типизация-и-явное-приведение-типов-type-conversion)  
> 💻 **Код лекции:** [02-syntax/1-basic](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/02-syntax/1-basic)  
> 🗺️ **Путеводитель:** [course-codebase-guide.md (Урок 02)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/course-codebase-guide.md#урок-02-фундаментальный-синтаксис-и-модель-памяти)

### 🎯 Условие:

В кинотеатре Goflex есть суммарное количество оценок фильма (целое число) и сумма всех баллов (целое число).
Рассчитайте точный средний рейтинг (число с плавающей точкой) и выведите округленный рейтинг до 1 знака.

### 💻 Эталонное решение:

```go
package main

import "fmt"

func main() {
    var totalVotes int = 143
    var totalScore int = 1247

    // ❌ ОШИБКА: totalScore / totalVotes даст целочисленное деление 8 (потеря дробной части)
    // ✅ ПРАВИЛЬНО: явное приведение хотя бы одного операнда к float64
    averageRating := float64(totalScore) / float64(totalVotes)

    fmt.Printf("Всего голосов: %d, Сумма баллов: %d\n", totalVotes, totalScore)
    fmt.Printf("Точный рейтинг: %f\n", averageRating)
    fmt.Printf("Отображаемый рейтинг на карточке фильма: %.1f ⭐\n", averageRating)
}
```

### 🔍 Разбор:

- В Go **нет неявного (автоматического) приведения типов**. Выражение `int / int` всегда выполняет целочисленное деление с отсечением остатка. Для вещественного результата требуется явный вызов `float64(...)`.

---

## Урок 4: Условия (if с областью видимости и switch)

> 📖 **Теория:** [Том 1: Разделы 1.2, 1.3 (if short statement, switch)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/01-go-core-basics.md#12-конструкция-if-с-коротким-объявлением-if-with-short-statement)  
> 💻 **Код лекции:** [02-syntax/1-basic](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/02-syntax/1-basic)  
> 🗺️ **Путеводитель:** [course-codebase-guide.md (Урок 02)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/course-codebase-guide.md#урок-02-фундаментальный-синтаксис-и-модель-памяти)

### 🎯 Условие:

Реализуйте проверку доступа к фильму по возрасту:

- Функция `checkAgeLimit(rating string) int` возвращает минимальный возраст для рейтинга: `"G"` (0+), `"PG-13"` (13+), `"R-18"` (18+) — используйте `switch`.
- В `main` сравните возраст пользователя с ограничением через `if` с коротким объявлением переменной (short statement).

### 💻 Эталонное решение:

```go
package main

import "fmt"

func checkAgeLimit(rating string) int {
    switch rating {
    case "G":
        return 0
    case "PG-13":
        return 13
    case "R-18":
        return 18
    default:
        return 99 // неизвестный рейтинг — максимальное ограничение
    }
}

func main() {
    userAge := 15
    targetMovieRating := "PG-13"

    // Идиоматичный if с объявлением переменной requiredAge в локальной области видимости:
    if requiredAge := checkAgeLimit(targetMovieRating); userAge >= requiredAge {
        fmt.Printf("✅ Доступ разрешен: пользователю %d лет, требуемый возраст %d+\n", userAge, requiredAge)
    } else {
        fmt.Printf("⛔ Доступ запрещен: требуется возраст %d+, текущий %d\n", requiredAge, userAge)
    }

    // Переменная requiredAge здесь недоступна — область видимости изолирована!
}
```

### 🔍 Разбор:

- Переменная `requiredAge` создается перед точкой с запятой в `if` и уничтожается после выхода из блока `if/else`. Это избавляет внешнюю функцию от лишних временных переменных.
- В блоке `switch` в Go оператор `break` подставляется автоматически в конце каждой ветки `case`.

---

## Урок 5: Циклы for, break и continue

> 📖 **Теория:** [Том 1: Разделы 1.1, 1.4 (Циклы for, метки break/continue)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/01-go-core-basics.md#11-формы-цикла-for-в-go)  
> 💻 **Код лекции:** [02-syntax/1-basic](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/02-syntax/1-basic)  
> 🗺️ **Путеводитель:** [course-codebase-guide.md (Урок 02)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/course-codebase-guide.md#урок-02-фундаментальный-синтаксис-и-модель-памяти)

### 🎯 Условие:

В плеере Goflex загружается сезон сериала из 8 серий:

- Пропустите серию 4 (спешл/рекап, не входящий в основной сюжет) с помощью `continue`.
- Остановите просмотр, если пользователь нажал паузу/выход на 6 серии (`break`).

### 💻 Эталонное решение:

```go
package main

import "fmt"

func main() {
    totalEpisodes := 8
    stopAtEpisode := 6

    fmt.Println("▶ Начат автопросмотр сезона:")

    for ep := 1; ep <= totalEpisodes; ep++ {
        if ep == 4 {
            fmt.Printf("  [Серия %d]: Пропуск (бонусный рекап)...\n", ep)
            continue // переход к следующей итерации
        }

        if ep == stopAtEpisode {
            fmt.Printf("  [Серия %d]: Пользователь закрыл плеер. Остановка сезона.\n", ep)
            break // досрочный выход из цикла
        }

        fmt.Printf("  [Серия %d]: Воспроизведение...\n", ep)
    }

    fmt.Println("⏹ Просмотр завершен.")
}
```

---

## Урок 6: Слайсы (Slices) и управление каталогом

> 📖 **Теория:** [Том 1: Раздел 5 (Массивы и слайсы, append, cap)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/01-go-core-basics.md#5-массивы-и-слайсы-arrays--slices)  
> 💻 **Код лекции:** [02-syntax/3-array_slice_map](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/02-syntax/3-array_slice_map)  
> 🛠️ **Домашнее задание:** [homework-02 (GoSearch crawler)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/02-syntax/homework-02)  
> 🗺️ **Путеводитель:** [course-codebase-guide.md (Урок 02)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/course-codebase-guide.md#урок-02-фундаментальный-синтаксис-и-модель-памяти)

### 🎯 Условие:

Создайте каталог фильмов пользователя:

1. Создайте слайс избранного емкостью 2 с помощью `make([]string, 0, 2)`.
2. Добавьте 4 фильма через `append` и проследите за изменением `len` и `cap`.
3. Создайте независимую копию первых 2 фильмов (Top-2) без риска утечки памяти под базовый массив.

### 💻 Эталонное решение:

```go
package main

import (
    "fmt"
    "slices"
)

func main() {
    // 1. Создаем слайс с заранее аллоцированной емкостью:
    favorites := make([]string, 0, 2)
    fmt.Printf("Init: len=%d, cap=%d\n", len(favorites), cap(favorites))

    // 2. Добавляем элементы:
    movies := []string{"Матрица", "Интерстеллар", "Начало", "Дюна"}
    for _, m := range movies {
        favorites = append(favorites, m)
        fmt.Printf("Добавлен '%s': len=%d, cap=%d\n", m, len(favorites), cap(favorites))
    }

    // 3. Безопасное создание Top-2 через slices.Clone (Go 1.21+):
    // Это исключает удержание всего 4-элементного массива в памяти, если favorites будет удален.
    top2 := slices.Clone(favorites[:2])

    fmt.Println("Весь список:", favorites)
    fmt.Printf("Top-2 (независимая копия): %v (len=%d, cap=%d)\n", top2, len(top2), cap(top2))
}
```

### 🔍 Разбор:

- При превышении `cap` рантайм выделяет новый массив большего размера в куче и копирует старые элементы.
- Срез `favorites[:2]` ссылается на тот же базовый массив. Использование `slices.Clone` или `copy()` создает полностью независимый массив и предотвращает утечки памяти.

---

## Урок 7: Словари (Maps) и быстрый поиск тайтлов

> 📖 **Теория:** [Том 1: Раздел 6 (Карты maps, устройство и доступ)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/01-go-core-basics.md#6-карты-maps)  
> 💻 **Код лекции:** [02-syntax/3-array_slice_map](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/02-syntax/3-array_slice_map)  
> 🛠️ **Домашнее задание:** [homework-03 (Инвертированный индекс)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/03-algorithms/homework-03)  
> 🗺️ **Путеводитель:** [course-codebase-guide.md (Урок 03)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/course-codebase-guide.md#урок-03-алгоритмы-и-вычислительная-сложность)

### 🎯 Условие:

Создайте каталог цен аренды фильмов `map[string]int`:

1. Проверьте наличие фильма через паттерн `comma ok`.
2. Удалите фильм из каталога через `delete()`.
3. Реализуйте стабильный детерминированный вывод каталога (отсортированный по алфавиту).

### 💻 Эталонное решение:

```go
package main

import (
    "fmt"
    "slices"
)

func main() {
    // 1. Инициализация мапы с хинтом размера:
    prices := make(map[string]int, 4)
    prices["Матрица"] = 299
    prices["Интерстеллар"] = 399
    prices["Оппенгеймер"] = 499

    // 2. Проверка наличия через comma ok:
    title := "Аватар"
    if price, ok := prices[title]; ok {
        fmt.Printf("Фильм '%s' доступен за %d руб.\n", title, price)
    } else {
        fmt.Printf("Фильм '%s' не найден в каталоге!\n", title)
    }

    // 3. Удаление:
    delete(prices, "Оппенгеймер")

    // 4. Сортированная итерация (так как for range по map выдает случайный порядок):
    keys := make([]string, 0, len(prices))
    for k := range prices {
        keys = append(keys, k)
    }
    slices.Sort(keys) // Go 1.21+

    fmt.Println("\nТекущий каталог (по алфавиту):")
    for _, k := range keys {
        fmt.Printf("  • %s: %d руб.\n", k, prices[k])
    }
}
```

---

## Урок 8: Структуры, указатели и методы (Бизнес-модель Goflex)

> 📖 **Теория:** [Том 1: Разделы 7, 8 (Указатели и структуры)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/01-go-core-basics.md#7-указатели-в-деталях-pointers), [Том 2: Раздел 9 (Методы)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/02-go-core-oop-methods.md#9-методы-и-интерфейсы-methods--interfaces)  
> 💻 **Код лекции:** [02-syntax/2-structures](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/02-syntax/2-structures), [02-syntax/4-pointers](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/02-syntax/4-pointers), [06-oop](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/06-oop)  
> 🛠️ **Домашнее задание:** [homework-06 (ООП движок)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/06-oop/homework-06)  
> 🗺️ **Путеводитель:** [course-codebase-guide.md (Урок 06)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/course-codebase-guide.md#урок-06-ооп-в-go-методы-композиция-и-dip)

### 🎯 Условие:

Смоделируйте сущности онлайн-кинотеатра:

- Структура `User` (`ID`, `Email`, `Balance`, `IsPremium`).
- Структура `Device` (`DeviceID`, `Model`, `LastActive`).
- Конструктор `NewUser(id int64, email string, initialBalance float64) *User`.
- Методы с ресивером-указателем: `Deposit(amount float64) error`, `BuySubscription(price float64) error` и `RegisterDevice(deviceID, model string)`.
- (В учебном решении баланс хранится в `float64` для простоты; в реальных системах для денег используют целые копейки или decimal-тип, см. Том 1, раздел 0.3.)

### 💻 Эталонное решение:

```go
package main

import (
    "errors"
    "fmt"
    "time"
)

var (
    ErrInsufficientFunds = errors.New("недостаточно средств на балансе")
    ErrInvalidAmount     = errors.New("сумма должна быть больше нуля")
)

type Device struct {
    DeviceID   string
    Model      string
    LastActive time.Time
}

type User struct {
    ID        int64
    Email     string
    Balance   float64
    IsPremium bool
    Devices   []Device
}

// Конструктор:
func NewUser(id int64, email string, initialBalance float64) *User {
    return &User{
        ID:        id,
        Email:     email,
        Balance:   initialBalance,
        IsPremium: false,
        Devices:   make([]Device, 0, 2),
    }
}

// Pointer Receiver: мутирует исходную структуру
func (u *User) Deposit(amount float64) error {
    if amount <= 0 {
        return ErrInvalidAmount
    }
    u.Balance += amount
    return nil
}

func (u *User) BuySubscription(price float64) error {
    if u.Balance < price {
        return ErrInsufficientFunds
    }
    u.Balance -= price
    u.IsPremium = true
    return nil
}

func (u *User) RegisterDevice(deviceID, model string) {
    u.Devices = append(u.Devices, Device{
        DeviceID:   deviceID,
        Model:      model,
        LastActive: time.Now(),
    })
}

func main() {
    user := NewUser(1, "alex@goflex.io", 200)
    user.RegisterDevice("dev-99", "Apple TV 4K")

    fmt.Printf("Создан пользователь: %+v\n", user)

    // Попытка купить подписку за 499:
    err := user.BuySubscription(499)
    fmt.Printf("Покупка 1 (баланс %.2f): %v\n", user.Balance, err)

    // Пополняем баланс:
    _ = user.Deposit(500)
    _ = user.BuySubscription(499)
    fmt.Printf("Покупка 2 (после пополнения): Премиум = %t, Баланс = %.2f\n", user.IsPremium, user.Balance)
}
```

### 🔍 Разбор:

- Если метод должен изменить поля структуры, **обязательно используется ресивер-указатель `(u *User)`**. Если указать `(u User)` (value receiver), Go молча скопирует структуру на стек, изменения применятся к копии и бесследно пропадут.

---

## Урок 9: Интерфейсы, полиморфизм и Type Switch

> 📖 **Теория:** [Том 2: Раздел 9 (Интерфейсы, полиморфизм, duck typing)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/02-go-core-oop-methods.md#9-методы-и-интерфейсы-methods--interfaces)  
> 💻 **Код лекции:** [06-oop](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/06-oop), [09-interfaces](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/09-interfaces)  
> 🛠️ **Домашние задания:** [homework-06](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/06-oop/homework-06), [homework-09](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/09-interfaces/homework-09)  
> 🗺️ **Путеводитель:** [course-codebase-guide.md (Урок 09)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/course-codebase-guide.md#урок-09-интерфейсы-и-полиморфизм)

### 🎯 Условие:

Создайте систему отправки уведомлений пользователям Goflex:

1. Интерфейс `Notifier` с методом `Notify(userID int64, message string) error`.
2. Реализации: `EmailNotifier` и `PushNotifier`.
3. Универсальную функцию рассылки `SendAlert(n Notifier, userID int64, msg string)`.
4. Демонстрацию работы `Type Switch` для извлечения дополнительной информации о канале доставки.

### 💻 Эталонное решение:

```go
package main

import (
    "fmt"
)

// 1. Интерфейс
type Notifier interface {
    Notify(userID int64, message string) error
}

// 2. Реализация Email:
type EmailNotifier struct {
    SMTPServer string
}

func (e EmailNotifier) Notify(userID int64, message string) error {
    fmt.Printf("[EMAIL через %s] Для User #%d: %s\n", e.SMTPServer, userID, message)
    return nil
}

// 3. Реализация Push:
type PushNotifier struct {
    FirebaseAppID string
}

func (p PushNotifier) Notify(userID int64, message string) error {
    fmt.Printf("[PUSH via Firebase %s] Для User #%d: %s\n", p.FirebaseAppID, userID, message)
    return nil
}

// 4. Универсальная функция с Type Switch:
func SendAlert(n Notifier, userID int64, msg string) {
    // Полиморфный вызов:
    _ = n.Notify(userID, msg)

    // Type Switch для проверки типа:
    switch v := n.(type) {
    case EmailNotifier:
        fmt.Printf("  ↳ Инфо: Отправлено по почтовому шлюзу %s\n", v.SMTPServer)
    case PushNotifier:
        fmt.Printf("  ↳ Инфо: Отправлено на мобильное устройство через Firebase (%s)\n", v.FirebaseAppID)
    default:
        fmt.Println("  ↳ Инфо: Неизвестный провайдер доставки")
    }
}

func main() {
    emailService := EmailNotifier{SMTPServer: "smtp.goflex.internal:587"}
    pushService := PushNotifier{FirebaseAppID: "goflex-ios-prod"}

    notifiers := []Notifier{emailService, pushService}

    for _, s := range notifiers {
        SendAlert(s, 42, "Вышел новый эпизод вашего любимого сериала!")
    }
}
```

---

## Урок 10: Конкурентность (Горутины, Каналы, sync, Worker Pool, Context)

> 📖 **Теория:** [Том 3: Разделы 11, 15 (Горутины, каналы, sync, Worker Pool)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/03-go-core-concurrency.md#11-конкурентность-горутины-и-каналы-concurrency)  
> 💻 **Код лекции:** [10-concurrency](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/10-concurrency), [21-interview](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/21-interview)  
> 🛠️ **Домашнее задание:** [homework-10 (Ping-Pong)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/10-concurrency/homework-10)  
> 💡 **Live-Coding:** [go-livecoding-guide.md](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/go-livecoding-guide.md)  
> 🗺️ **Путеводитель:** [course-codebase-guide.md (Урок 10)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/course-codebase-guide.md#урок-10-конкурентность-горутины-и-каналы)

### 🎯 Комплексное задание со всеми 5 подзадачами из видео:

1. **Подзадача 1 (`sync.WaitGroup`):** Параллельная проверка доступности эпизодов в CDN.
2. **Подзадача 2 (Каналы):** Сбор статусов транскодирования видео через канал без использования `time.Sleep`.
3. **Подзадача 3 (Буферизованный канал):** Очередь событий плеера (`Play`, `Pause`, `Stop`) с закрытием и чтением через `for range`.
4. **Подзадача 4 (`select` & Timeout):** Загрузка предпочтений пользователя с таймаутом через `time.After`.
5. **Подзадача 5 (Worker Pool + `context.Context` + `sync.Mutex`):** Пул из 3 воркеров, обрабатывающих очередь аналитики, с подсчетом обработанных задач и отменой по таймауту.

### 💻 Эталонное решение:

```go
package main

import (
    "context"
    "fmt"
    "sync"
    "time"
)

// === 1. Проверка доступности эпизодов с WaitGroup ===
func checkEpisodeAvailability(wg *sync.WaitGroup, episodeName string, delay time.Duration) {
    defer wg.Done()
    time.Sleep(delay) // имитация сетевого запроса к CDN
    fmt.Printf("  ✔ Эпизод '%s' проверен и доступен\n", episodeName)
}

// === 2. Транскодирование видео через канал ===
func transcodeVideo(filename string, out chan<- string) {
    time.Sleep(50 * time.Millisecond) // имитация тяжелой работы
    out <- fmt.Sprintf("Видео '%s' успешно перекодировано в 1080p H.264", filename)
}

// === 4. Загрузка рекомендаций с таймаутом ===
func loadUserPreferences(out chan<- string) {
    time.Sleep(80 * time.Millisecond) // имитация обращения к ML-сервису
    out <- "Рекомендации: [Sci-Fi, Cyberpunk, Noir]"
}

// === 5. Worker Pool для аналитики ===
type AnalyticsEvent struct {
    UserID int
    Action string
}

func analyticsWorker(ctx context.Context, id int, jobs <-chan AnalyticsEvent, mu *sync.Mutex, counter *int, wg *sync.WaitGroup) {
    defer wg.Done()
    for {
        select {
        case <-ctx.Done():
            fmt.Printf("  [Воркер %d]: Остановлен по контексту (%v)\n", id, ctx.Err())
            return
        case job, ok := <-jobs:
            if !ok {
                return // канал закрыт, задач больше нет
            }
            // Обрабатываем событие:
            time.Sleep(10 * time.Millisecond)
            mu.Lock()
            *counter++
            mu.Unlock()
            fmt.Printf("  [Воркер %d]: Записано событие '%s' юзера %d\n", id, job.Action, job.UserID)
        }
    }
}

func main() {
    fmt.Println("=== 1. Проверка CDN эпизодов (sync.WaitGroup) ===")
    var wg sync.WaitGroup
    episodes := []string{"s01e01", "s01e02", "s01e03"}
    for _, ep := range episodes {
        wg.Add(1)
        go checkEpisodeAvailability(&wg, ep, 30*time.Millisecond)
    }
    wg.Wait()
    fmt.Println("Все эпизоды готовы к раздаче!")

    fmt.Println("\n=== 2. Сбор статусов транскодирования через канал ===")
    transcodeChan := make(chan string, len(episodes))
    for _, ep := range episodes {
        go transcodeVideo(ep+".mp4", transcodeChan)
    }
    for i := 0; i < len(episodes); i++ {
        result := <-transcodeChan
        fmt.Println("  ↳", result)
    }

    fmt.Println("\n=== 3. Очередь событий плеера (Buffered Channel) ===")
    playerEvents := make(chan string, 4)
    playerEvents <- "EVENT_PLAY"
    playerEvents <- "EVENT_SEEK_10m"
    playerEvents <- "EVENT_PAUSE"
    playerEvents <- "EVENT_STOP"
    close(playerEvents) // обязательно закрываем перед for range!

    for event := range playerEvents {
        fmt.Printf("  Обработано событие плеера: %s\n", event)
    }

    fmt.Println("\n=== 4. Загрузка рекомендаций с select и timeout ===")
    prefChan := make(chan string, 1)
    go loadUserPreferences(prefChan)

    select {
    case data := <-prefChan:
        fmt.Println("  Получены данные:", data)
    case <-time.After(50 * time.Millisecond): // таймаут меньше времени ответа
        fmt.Println("  ⚠ Таймаут! Сервер рекомендаций отвечает слишком долго. Отдаем дефолтный топ.")
    }

    fmt.Println("\n=== 5. Аналитический Worker Pool с Context и Mutex ===")
    ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
    defer cancel()

    jobs := make(chan AnalyticsEvent, 10)
    var poolWg sync.WaitGroup
    var mu sync.Mutex
    processedCount := 0

    // Запускаем 3 воркера:
    for w := 1; w <= 3; w++ {
        poolWg.Add(1)
        go analyticsWorker(ctx, w, jobs, &mu, &processedCount, &poolWg)
    }

    // Отправляем события в очередь:
    events := []AnalyticsEvent{
        {UserID: 1, Action: "START_PLAY"},
        {UserID: 2, Action: "LIKE_MOVIE"},
        {UserID: 3, Action: "ADD_TO_FAVORITES"},
        {UserID: 1, Action: "STOP_PLAY"},
    }
    for _, ev := range events {
        jobs <- ev
    }
    close(jobs)

    poolWg.Wait()
    fmt.Printf("Всего событий обработано безопасно под мьютексом: %d\n", processedCount)
}
```

---

## Урок 11: Web HTTP REST API и Middleware

> 📖 **Теория:** [Том 4: Раздел 12 (HTTP сервер, маршрутизация, middleware)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/04-go-core-backend-web.md#12-backend-web-разработка-и-http-nethttp)  
> 💻 **Код лекции:** [12-web-apps](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/12-web-apps), [13-api](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/13-api)  
> 🛠️ **Домашние задания:** [homework-12 (Web-сервер)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/12-web-apps/homework-12), [homework-13 (REST API)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/13-api/homework-13)  
> 🗺️ **Путеводитель:** [course-codebase-guide.md (Уроки 12–13)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/course-codebase-guide.md#урок-12-веб-приложения-и-http-сервер)

### 🎯 Условие:

Создайте сервер REST API онлайн-кинотеатра Goflex на стандартной библиотеке `net/http` (Go 1.22+):

1. Маршрут `GET /api/v1/movies/{id}` — возвращает JSON с фильмом.
2. Маршрут `POST /api/v1/movies` — создает новый фильм из JSON тела запроса.
3. Middleware логирования времени выполнения запросов.
4. Middleware защиты от паники (Panic Recovery Middleware).
5. Graceful Shutdown по сигналу операционной системы (`SIGINT`/`SIGTERM`).

### 💻 Эталонное решение:

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)

type Movie struct {
    ID     string  `json:"id"`
    Title  string  `json:"title"`
    Rating float64 `json:"rating"`
}

// 1. Middleware логирования:
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("[%s] %s %s - %v", r.Method, r.URL.Path, r.RemoteAddr, time.Since(start))
    })
}

// 2. Recovery Middleware:
func RecoveryMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("🔥 ПАНИКА В ХЕНДЛЕРЕ: %v", err)
                http.Error(w, `{"error": "Внутренняя ошибка сервера"}`, http.StatusInternalServerError)
            }
        }()
        next.ServeHTTP(w, r)
    })
}

func main() {
    mux := http.NewServeMux()

    // Go 1.22+ встроенная поддержка метода и path-параметров {id}:
    mux.HandleFunc("GET /api/v1/movies/{id}", func(w http.ResponseWriter, r *http.Request) {
        movieID := r.PathValue("id")
        movie := Movie{ID: movieID, Title: "Начало (Inception)", Rating: 8.8}

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        _ = json.NewEncoder(w).Encode(movie)
    })

    mux.HandleFunc("POST /api/v1/movies", func(w http.ResponseWriter, r *http.Request) {
        var m Movie
        if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
            http.Error(w, `{"error": "Невалидный JSON"}`, http.StatusBadRequest)
            return
        }
        m.ID = "generated-101"
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)
        _ = json.NewEncoder(w).Encode(m)
    })

    // Сборка цепочки Middleware:
    handler := RecoveryMiddleware(LoggingMiddleware(mux))

    server := &http.Server{
        Addr:         ":8080",
        Handler:      handler,
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 10 * time.Second,
    }

    // Graceful Shutdown:
    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer stop()

    go func() {
        fmt.Println("🚀 Goflex API запущен на http://localhost:8080")
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Ошибка запуска: %v", err)
        }
    }()

    <-ctx.Done() // Ожидание сигнала остановки
    fmt.Println("\n🛑 Получен сигнал завершения. Плавная остановка сервера...")

    shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := server.Shutdown(shutdownCtx); err != nil {
        log.Fatalf("Принудительная остановка: %v", err)
    }
    fmt.Println("✅ Сервер успешно остановлен без потери клиентских запросов.")
}
```

---

## Урок 12: PostgreSQL, pgxpool, ACID-транзакции и Миграции

> 📖 **Теория:** [Том 4: Раздел 13 (PostgreSQL, pgxpool, транзакции)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/04-go-core-backend-web.md#13-работа-с-базами-данных-postgresql--pgxpool)  
> 🗄️ **Гайд по БД:** [go-database-interview-guide.md](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/go-database-interview-guide.md)  
> 💻 **Код лекции:** [15-sql](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/15-sql), [16-db-apps](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/16-db-apps)  
> 🛠️ **Домашние задания:** [homework-15](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/15-sql/homework-15), [homework-16](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/16-db-apps/homework-16)  
> 🗺️ **Путеводитель:** [course-codebase-guide.md (Уроки 15–16)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/course-codebase-guide.md#урок-15-реляционные-базы-данных-и-sql)

### 🎯 Условие:

Реализуйте слой работы с базой данных для оплаты подписки в Goflex:

1. Подключение через `pgxpool.Pool` с методом `Ping()`.
2. Маппинг ошибки `pgx.ErrNoRows` в доменную ошибку `ErrUserNotFound`.
3. ACID-транзакция перевода денег за подписку (списание с баланса пользователя и зачисление на счет кинотеатра) с каноничным `defer tx.Rollback(ctx)`.
4. Пример файла миграции для `goose`.

### 💻 Эталонное решение:

#### 1. SQL-миграция (`00001_create_users.sql` для goose):

```sql
-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    balance NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
    is_premium BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE system_accounts (
    id VARCHAR(50) PRIMARY KEY,
    balance NUMERIC(15, 2) NOT NULL DEFAULT 0.00
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS system_accounts;
-- +goose StatementEnd
```

#### 2. Go-код репозитория (`repository.go`):

```go
package main

import (
    "context"
    "errors"
    "fmt"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

var (
    ErrUserNotFound     = errors.New("пользователь не найден")
    ErrInsufficientCash = errors.New("недостаточно средств для оплаты")
)

type PostgresRepository struct {
    pool *pgxpool.Pool
}

func NewPostgresRepository(ctx context.Context, connString string) (*PostgresRepository, error) {
    pool, err := pgxpool.New(ctx, connString)
    if err != nil {
        return nil, fmt.Errorf("ошибка создания пула: %w", err)
    }

    if err := pool.Ping(ctx); err != nil {
        pool.Close() // не оставляем «висящий» пул при неудачной проверке
        return nil, fmt.Errorf("база данных недоступна: %w", err)
    }

    return &PostgresRepository{pool: pool}, nil
}

// Покупка подписки через ACID-транзакцию.
// Учебное упрощение: сумма и баланс — float64. В реальных проектах деньги хранят в целых минимальных единицах (int64)
// или используют pgtype.Numeric / shopspring decimal, чтобы избежать ошибок округления.
func (r *PostgresRepository) SubscribeUserTx(ctx context.Context, userID int64, price float64) error {
    tx, err := r.pool.Begin(ctx)
    if err != nil {
        return fmt.Errorf("begin tx failed: %w", err)
    }
    // Каноничный defer Rollback: если транзакция не закоммичена, она гарантированно откатится
    defer tx.Rollback(ctx)

    // 1. Проверяем баланс и списываем средства (FOR UPDATE блокирует строку от Race Condition):
    var currentBalance float64
    err = tx.QueryRow(ctx, "SELECT balance FROM users WHERE id = $1 FOR UPDATE", userID).Scan(&currentBalance)
    if errors.Is(err, pgx.ErrNoRows) {
        return ErrUserNotFound
    }
    if err != nil {
        return fmt.Errorf("select balance failed: %w", err)
    }

    if currentBalance < price {
        return ErrInsufficientCash
    }

    // 2. Списываем средства и активируем премиум:
    _, err = tx.Exec(ctx, "UPDATE users SET balance = balance - $1, is_premium = true WHERE id = $2", price, userID)
    if err != nil {
        return fmt.Errorf("update user failed: %w", err)
    }

    // 3. Зачисляем выручку на системный счет кинотеатра:
    _, err = tx.Exec(ctx, "UPDATE system_accounts SET balance = balance + $1 WHERE id = 'goflex_revenue'", price)
    if err != nil {
        return fmt.Errorf("update system account failed: %w", err)
    }

    // 4. Коммитим изменения:
    if err := tx.Commit(ctx); err != nil {
        return fmt.Errorf("commit failed: %w", err)
    }

    return nil
}
```

### 🔍 Разбор:

1. `defer tx.Rollback(ctx)`: Если функция завершается раньше (по ошибке или панике), транзакция автоматически откатывается. Если был вызван `tx.Commit(ctx)`, последующий `tx.Rollback()` безопасен: он возвращает `pgx.ErrTxClosed`, а мы эту ошибку в `defer` сознательно не проверяем.
2. `FOR UPDATE` в PostgreSQL: предотвращает «двойное списание» (Double-Spending Attack), если пользователь одновременно отправит два параллельных запроса на оплату подписки.
3. Маппинг `pgx.ErrNoRows -> ErrUserNotFound`: слой бизнес-логики не должен зависеть от драйвера БД, он оперирует чистыми доменными ошибками (Domain Errors / Sentinel Errors).

---

## Упражнение (A Tour of Go): Квадратный корень методом Ньютона (Циклы и функции)

> 📖 **Теория:** [Том 1: Разделы 1.1, 2.1 (Функции и циклы)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/01-go-core-basics.md#11-формы-цикла-for-в-go)  
> 💻 **Код лекции:** [02-syntax/1-basic](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/02-syntax/1-basic)

### 🎯 Условие:

Реализовать функцию `Sqrt(x float64) float64` для вычисления квадратного корня методом касательных Ньютона:

1. **Базовый алгоритм:** Начиная с начального приближения $z = 1.0$, итеративно корректировать значение по формуле:
   $$z = z - \frac{z^2 - x}{2z}$$
2. **Шаг 1 (Фиксированный цикл):** Выполнить вычисление ровно 10 раз, выводя значение $z$ и дельту с `math.Sqrt(x)` на каждом шаге.
3. **Шаг 2 (Адаптивный цикл):** Изменить условие цикла так, чтобы он останавливался автоматически, когда разница между двумя последовательными шагами становится меньше заданной машинной погрешности ($\Delta \le 10^{-15}$).
4. **Шаг 3 (Граничные условия):** Обработать отрицательные числа (`math.NaN()`) и ноль (`x == 0`).
5. **Шаг 4 (Исследование):** Сравнить начальные приближения $z = 1.0$, $z = x$ и $z = x/2$ по количеству итераций до полной сходимости.

### 💻 Эталонное решение:

```go
package main

import (
	"fmt"
	"math"
)

// Sqrt10 — вычисление за фиксированные 10 итераций с логированием шагов
func Sqrt10(x float64) float64 {
	if x < 0 {
		return math.NaN()
	}
	if x == 0 {
		return 0
	}

	z := 1.0
	fmt.Printf("--- Вычисление Sqrt10(x = %.1f) ---\n", x)
	for i := 1; i <= 10; i++ {
		z -= (z*z - x) / (2 * z)
		delta := math.Abs(z - math.Sqrt(x))
		fmt.Printf("Шаг %2d: z = %.15f | разница с math.Sqrt: %e\n", i, z, delta)
	}
	return z
}

// SqrtAdaptive — адаптивное вычисление до стабилизации значения (точность машинного нуля)
func SqrtAdaptive(x float64, initialGuess float64) (float64, int) {
	if x < 0 {
		return math.NaN(), 0
	}
	if x == 0 {
		return 0, 0
	}

	z := initialGuess
	const epsilon = 1e-15 // абсолютный допуск (машинный эпсилон float64 — около 2.2e-16 в относительном смысле)
	iterations := 0

	for {
		iterations++
		prev := z
		z -= (z*z - x) / (2 * z)

		// Условие выхода: изменение меньше машинной погрешности
		if math.Abs(z-prev) <= epsilon {
			break
		}
	}
	return z, iterations
}

func main() {
	// 1. Демонстрация фиксированных 10 итераций (x = 2):
	Sqrt10(2)

	// 2. Сравнение сходимости при разных начальных приближениях z:
	testValues := []float64{2, 3, 25, 1234567}
	fmt.Println("\n=== Сравнение стратегий начального приближения z ===")

	for _, x := range testValues {
		res1, iter1 := SqrtAdaptive(x, 1.0)
		resX, iterX := SqrtAdaptive(x, x)
		resHalf, iterHalf := SqrtAdaptive(x, x/2.0)
		expected := math.Sqrt(x)

		fmt.Printf("\nЧисло x = %-9.1f (math.Sqrt = %.15f):\n", x, expected)
		fmt.Printf("  • z0 = 1.0 : результат = %.15f (итераций: %d, совпало: %t)\n", res1, iter1, res1 == expected)
		fmt.Printf("  • z0 = x   : результат = %.15f (итераций: %d, совпало: %t)\n", resX, iterX, resX == expected)
		fmt.Printf("  • z0 = x/2 : результат = %.15f (итераций: %d, совпало: %t)\n", resHalf, iterHalf, resHalf == expected)
	}
}
```

### 🔍 Разбор:

1. **Математическая суть (Метод Ньютона-Рафсона):**
   - Мы ищем корень уравнения $f(z) = z^2 - x = 0$.
   - Метод Ньютона строит касательную в точке $z_n$ и находит точку её пересечения с осью $OX$:
     $$z_{n+1} = z_n - \frac{f(z_n)}{f'(z_n)} = z_n - \frac{z_n^2 - x}{2z_n}$$
2. **Квадратичная скорость сходимости:**
   - Каждая итерация удваивает количество точных значащих цифр после запятой.
   - Для $x = 2$ уже на **5-й итерации** результат идентичен аппаратному `math.Sqrt(2)` вплоть до последнего бита мантиссы (все 53 бита IEEE 754 совпадают).
3. **Выбор начального приближения ($z_0$):**
   - Из-за квадратичной сходимости выбор $z_0$ почти не влияет на число итераций. Замеры функции `SqrtAdaptive` (число итераций включает последний шаг, подтверждающий сходимость):

     | $x$ | $z_0 = 1.0$ | $z_0 = x$ | $z_0 = x/2$ |
     | :--- | :---: | :---: | :---: |
     | 2 | 6 | 6 | 6 |
     | 3 | 6 | 6 | 5 |
     | 25 | 8 | 8 | 7 |
     | 1 234 567 | 16 | 16 | 15 |

   - Разница составляет 0–1 итерацию: основную часть шагов занимает не «подъём из единицы», а выход на нужный порядок величины (число шагов растёт примерно как $\log_2$ от порядка $x$). Более близкое к ответу начальное приближение (например, $z_0 = x/2$ для малых $x$) даёт лишь небольшой выигрыш.
4. **Подводные камни IEEE 754:**
   - ⚠️ **Не полагайтесь на точное сравнение `z == prev`.** Из-за округления чисел с плавающей точкой итерации в общем случае могут не остановиться на одном значении, а колебаться между двумя соседними числами последнего разряда. Общее правило работы с `float`: сравнивать через допуск, например `math.Abs(z - prev) <= epsilon`. В примере это работает, так как метод Ньютона для `sqrt` сходится точно (на тестовых значениях разность в итоге равна `0`), но надёжная реализация использует **относительный допуск** (`math.Abs(z-prev) <= epsilon*z`) **и ограничение числа итераций**, чтобы цикл гарантированно завершался при любых входных данных: абсолютный `epsilon = 1e-15` меньше расстояния между соседними числами `float64` для больших значений (для $z \approx 10^3$ оно около $2 \cdot 10^{-13}$).
   - Защита от деления на ноль: при $x = 0$ значение $2z = 0$, поэтому необходима проверка `if x == 0 { return 0 }`.

---

## Занятие 2: Поисковый робот GoSearch (Пакет crawler и консольный поиск)

> 📖 **Теория:** [Том 1: Раздел 5 (Массивы и слайсы)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/01-go-core-basics.md#5-массивы-и-слайсы-arrays--slices)  
> 💻 **Код лекции:** [02-syntax](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/02-syntax)  
> 📁 **Каталог решения:** [homework-02](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/02-syntax/homework-02)  
> 💻 **Исполняемый пакет:** [homework-02/cmd/gosearch/main.go](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/02-syntax/homework-02/cmd/gosearch/main.go)  
> 📦 **Пакеты:** [crawler](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/02-syntax/homework-02/pkg/crawler), [spider](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/02-syntax/homework-02/pkg/crawler/spider), [membot](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/02-syntax/homework-02/pkg/crawler/membot)  
> 🧪 **Тесты:** [homework-02/cmd/gosearch/main_test.go](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/02-syntax/homework-02/cmd/gosearch/main_test.go)  
> 🗺️ **Путеводитель:** [course-codebase-guide.md (Урок 02)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/course-codebase-guide.md#урок-02-фундаментальный-синтаксис-и-модель-памяти)

### 🎯 Условие задачи:

1. Создать приложение со стандартной структурой каталогов `cmd/` и `pkg/`. Скопировать пакет `crawler`.
2. Использовать пакет `crawler` для сканирования сайтов `https://go.dev` и `https://golang.org` (глубина 1). Объединить результаты сканирования.
3. Добавить обработку флага `-s <слово>`: при его указании напечатать все найденные страницы, содержащие данное слово в Title или URL. При отсутствии флага — вывести полный список найденных ссылок.

### 💡 Архитектурные особенности:

- Для автономной работы и изолированного тестирования без реального выхода в сеть реализован эмулятор сканера `membot`.
- Поиск на первой итерации выполняется линейным сканированием среза (`FilterDocs`) со сложностью $O(N)$.

---

## Занятие 3: Быстрый поисковый индекс (Инвертированный индекс и бинарный поиск)

> 📖 **Теория:** [Том 1: Раздел 6 (Карты maps)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/01-go-core-basics.md#6-карты-maps), [Том 6: Разделы 26–27 (Сложность O, поиск)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/06-go-core-algorithms-interview.md#26-базовые-алгоритмы-и-o-нотация-в-go)  
> 💻 **Код лекции:** [03-algorithms](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/03-algorithms)  
> 📁 **Каталог решения:** [homework-03](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/03-algorithms/homework-03)  
> 💻 **Исполняемый пакет:** [homework-03/cmd/gosearch/main.go](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/03-algorithms/homework-03/cmd/gosearch/main.go)  
> 📦 **Пакеты:** [pkg/index](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/03-algorithms/homework-03/pkg/index), [pkg/crawler](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/03-algorithms/homework-03/pkg/crawler)  
> 🧪 **Тесты:** [index_test.go](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/03-algorithms/homework-03/pkg/index/index_test.go), [main_test.go](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/03-algorithms/homework-03/cmd/gosearch/main_test.go)  
> 🗺️ **Путеводитель:** [course-codebase-guide.md (Урок 03)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/course-codebase-guide.md#урок-03-алгоритмы-и-вычислительная-сложность)

### 🎯 Условие задачи:

1. **Задача №1 (Обратный индекс):**  
   Создать обратный поисковый индекс в пакете `pkg/index`. Ключом индекса является нормализованное слово из заголовка страницы (`Title`), значением — срез идентификаторов документов (`[]int`). Для каждого документа добавить уникальный номер `ID int`. Массив документов отсортировать по `ID` с помощью стандартной библиотеки (`sort.Slice`).
2. **Задача №2 (Поисковая выдача через индекс):**  
   Переделать поисковую выдачу: при передаче флага `-s` выполнять поиск номеров документов через индекс (`idx.Search(word)`), а не линейным перебором массива.
3. **Задача №3 (Бинарный поиск по массиву документов):**  
   Для каждого найденного в индексе `docID` выполнять бинарный поиск (`sort.Search`) по отсортированному по `ID` массиву документов за время $O(\log N)$.

### 🔍 Разбор сложности алгоритма:

- **Построение индекса:** $O(D \times W)$, где $D$ — количество документов, $W$ — среднее число слов в заголовке.
- **Сортировка документов по ID:** $O(D \log D)$ через `sort.Slice` (алгоритм pdqsort).
- **Поиск в индексе:** $O(1)$ в среднем (доступ по ключу в хэш-таблице `map[string][]int`).
- **Извлечение каждого найденного документа:** $O(\log D)$ с помощью бинарного поиска `index.BinarySearch(docs, id)`.

---

## Занятие 4: Структуры данных — Циклический двусвязный список со сторожевым элементом (Pop и Reverse)

> 📖 **Теория:** [Том 6: Раздел 27 (Классические структуры данных)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/06-go-core-algorithms-interview.md#27-классические-структуры-данных-data-structures-в-go)  
> 💻 **Код лекции:** [04-datastructs](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/04-datastructs)  
> 📁 **Каталог решения:** [homework-04](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/04-datastructs/homework-04) и [04-datastructs/1-list](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/04-datastructs/1-list)  
> 🧪 **Тесты:** [list_test.go](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/04-datastructs/homework-04/pkg/list/list_test.go)  
> 🗺️ **Путеводитель:** [course-codebase-guide.md (Урок 04)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/course-codebase-guide.md#урок-04-классические-структуры-данных)

### 🎯 Условие задачи:

Доделать реализацию пакета для двусвязного списка `list`:

1. **`func (l *List) Pop() *List`**: удаляет первый значащий элемент списка (голову) и возвращает указатель на список `*List` для поддержки цепочечных вызовов (`l.Pop().Pop()`). Если список уже пуст, паники возникать не должно.
2. **`func (l *List) Reverse() *List`**: разворачивает список на месте (_in-place_) за время $O(N)$ и память $O(1)$, меняя направление указателей `next` и `prev` у всех элементов, и возвращает `*List`.
3. Обеспечить успешное прохождение всех модульных тестов `go test -v ./lessons/04-datastructs/1-list/...`.

---

### 💡 Архитектурные ловушки и концепции:

1. **Сторожевой элемент (Sentinel Node / Dummy Root):**
   - В отличие от простых списков, где `head == nil` обозначает пустоту, здесь список закольцован через фиктивный корневой элемент `root`:
     - В пустом списке: `l.root.next = l.root` и `l.root.prev = l.root`.
     - Первый элемент списка находится в `l.root.next`.
     - Последний элемент списка находится в `l.root.prev`.
   - **Преимущество паттерна Sentinel:** устраняет необходимость ветвлений `if head == nil` и краевых проверок при вставке и удалении первого/последнего узла.
2. **Цепочечные вызовы методов (Fluent Interface / Method Chaining):**
   - Сигнатура `Pop() *List` возвращает `*List`, а не удаленное значение, что позволяет писать выразительный код в тестах: `l.Pop().Pop().String()`.
3. **Корректность разворота (In-Place Reverse):**
   - Разворот циклического списка со сторожевым элементом требует аккуратного перенаправления связей:
     - Для каждого значащего узла от `l.root.next` до `l.root` меняем местами ссылки `next` и `prev`.
     - После завершения цикла обновляем указатели самого `l.root`: `l.root.next` должен указывать на бывший хвост, а `l.root.prev` — на бывшую голову.
4. **Защита от потери указателя `root.prev` при вставке:**
   - При добавлении самого первого элемента в пустой список (`Push`), необходимо не забыть обновить `l.root.prev`, иначе обратный указатель останется указывать на сам `root`.

---

### 💻 Эталонное решение:

Файл `04-datastructs/1-list/list.go`:

```go
// Реализация двусвязного циклического списка со сторожевым элементом.
package list

import (
	"fmt"
)

// List - двусвязный список со сторожевым элементом root.
type List struct {
	root *Elem
}

// Elem - элемент списка.
type Elem struct {
	Val        interface{}
	next, prev *Elem
}

// New создаёт список и возвращает указатель на него.
func New() *List {
	var l List
	l.root = &Elem{}
	l.root.next = l.root
	l.root.prev = l.root
	return &l
}

// Push вставляет элемент в начало списка.
func (l *List) Push(e Elem) *Elem {
	e.prev = l.root
	e.next = l.root.next
	l.root.next = &e
	if e.next != l.root {
		e.next.prev = &e
	} else {
		l.root.prev = &e // первый добавленный элемент также является хвостом
	}
	return &e
}

// String реализует интерфейс fmt.Stringer, представляя список в виде строки значений через пробел.
func (l *List) String() string {
	el := l.root.next
	var s string
	for el != l.root {
		s += fmt.Sprintf("%v ", el.Val)
		el = el.next
	}
	if len(s) > 0 {
		s = s[:len(s)-1]
	}
	return s
}

// Pop удаляет первый элемент списка. Поддерживает цепочечные вызовы.
// Сложность: Time O(1), Space O(1)
func (l *List) Pop() *List {
	if l == nil || l.root == nil || l.root.next == l.root {
		return l // список пуст
	}

	first := l.root.next
	l.root.next = first.next
	first.next.prev = l.root

	// Если список опустел, восстанавливаем указатель prev на корень
	if l.root.next == l.root {
		l.root.prev = l.root
	}

	return l
}

// Reverse разворачивает список на месте (in-place).
// Сложность: Time O(N), Space O(1)
func (l *List) Reverse() *List {
	if l == nil || l.root == nil || l.root.next == l.root {
		return l
	}

	curr := l.root.next
	var prev *Elem = l.root

	for curr != l.root {
		next := curr.next
		curr.next = prev
		curr.prev = next
		prev = curr
		curr = next
	}

	oldHead := l.root.next
	l.root.next = prev
	l.root.prev = oldHead

	return l
}
```

---

### 🔍 Разбор работы под капотом и валидация тестов:

1. **Схема связей до и после `Reverse()`:**
   - Исходный список после `Push(3)`, `Push(2)`, `Push(1)`:
     ```text
     root -> [1] <-> [2] <-> [3] -> root
     root <- [1]            [3] <- root
     ```
   - Разворот пошагово перенаправляет ссылки `next`:
     - Узел `[1]`: `next` становится `root`, `prev` становится `[2]`
     - Узел `[2]`: `next` становится `[1]`, `prev` становится `[3]`
     - Узел `[3]`: `next` становится `[2]`, `prev` становится `root`
   - Итог: `root.next` указывает на `[3]`, `root.prev` указывает на `[1]`:
     ```text
     root -> [3] <-> [2] <-> [1] -> root
     ```
2. **Прохождение тестов:**
   ```bash
   $ go test -v ./lessons/04-datastructs/1-list/...
   === RUN   TestList_Push
   --- PASS: TestList_Push (0.00s)
   === RUN   TestList_Pop
   --- PASS: TestList_Pop (0.00s)
   === RUN   TestList_Reverse
   --- PASS: TestList_Reverse (0.00s)
   PASS
   ok      go-core-4/04-datastructs/1-list 0.441s
   ```

---

## Занятие 5: Ввод-вывод в Go — Персистентность поисковых данных (io.Reader, io.Writer, JSON/GOB кеширование)

> 📖 **Теория:** [Том 2: Раздел 9.6 (Потоковый ввод-вывод io.Reader / io.Writer)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/02-go-core-oop-methods.md#96-потоковый-ввод-вывод-интерфейсы-ioreader-и-iowriter)  
> 💻 **Код лекции:** [05-io](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/05-io)  
> 📁 **Каталог решения:** [homework-05](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/05-io/homework-05)  
> 📦 **Пакеты:** [homework-05/pkg/storage](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/05-io/homework-05/pkg/storage)  
> 🧪 **Тесты:** `go test -v ./lessons/05-io/homework-05/...`  
> 🗺️ **Путеводитель:** [course-codebase-guide.md (Урок 05)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/course-codebase-guide.md#урок-05-потоковый-ввод-вывод-файлы-и-json)

В данном домашнем задании (каталог `go_course_4/05-io/homework-05/`) реализована персистентность поисковых данных робота **GoSearch**: сохранение отсканированных веб-страниц в файл на диске и их мгновенная загрузка при повторном запуске без необходимости повторного сетевого сканирования.

### 🎯 Условия задач занятия №5:

1. **Задача №1:** Разработать механизм записи полученных от сканера данных в файл. Механизм должен использовать интерфейс `io.Writer` в качестве аргумента. Формат представления данных для хранения выбирается на своё усмотрение (JSON / GOB).
2. **Задача №2:** Загрузка данных из файла при старте поисковика. Если файл с данными существует, то сканировать сайты не нужно — данные загружаются из файла. В процессе обязательно должен использоваться интерфейс `io.Reader`.

### 💡 Ключевые концепции и архитектурные решения:

1. **Инверсия зависимостей через `io.Writer` и `io.Reader`:**
   Функции сохранения и загрузки не привязываются к конкретному файлу `*os.File`. Принимая `io.Writer` и `io.Reader`, функции `storage.Store(w, docs)` и `storage.Load(r)` могут работать с любым источником: файлом, сетевым сокетом, HTTP-телом или `bytes.Buffer` в модульных тестах.
2. **Формат хранения:**
   - Основной формат: **JSON** (`json.NewEncoder(w)` / `json.NewDecoder(r)`). Человекочитаемый, структурированный, поддерживающий отступы `SetIndent("", "  ")`.
   - Альтернативный бинарный формат: **GOB** (`encoding/gob`), нативный бинарный формат Go-to-Go для максимальной скорости и компактности.
3. **Логика старта поисковика:**
   - Проверка наличия кэш-файла `docs.json` через `os.Stat()`.
   - Если файл существует: открытие через `os.Open()` $\rightarrow$ передача дескриптора как `io.Reader` в `storage.Load()` $\rightarrow$ пропуск сетевого сканирования.
   - Если файла нет (или передан флаг `-rescan`): запуск краулера $\rightarrow$ сортировка по ID $\rightarrow$ создание файла через `os.Create()` $\rightarrow$ передача как `io.Writer` в `storage.Store()`.

### 💻 Реализация пакета `pkg/storage` (`homework-05/pkg/storage/storage.go`):

```go
package storage

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"

	"go-core-4/05-io/homework-05/pkg/crawler"
)

// Store выполняет сериализацию документов в формате JSON и записывает их в io.Writer.
// Задача №1: Механизм записи данных сканера с использованием io.Writer.
func Store(w io.Writer, docs []crawler.Document) error {
	if w == nil {
		return fmt.Errorf("storage.Store: writer не может быть nil")
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(docs); err != nil {
		return fmt.Errorf("storage.Store: ошибка кодирования JSON: %w", err)
	}
	return nil
}

// Load считывает и десериализует срез документов из io.Reader в формате JSON.
// Задача №2: Механизм загрузки данных из файла с использованием io.Reader.
func Load(r io.Reader) ([]crawler.Document, error) {
	if r == nil {
		return nil, fmt.Errorf("storage.Load: reader не может быть nil")
	}

	var docs []crawler.Document
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&docs); err != nil {
		return nil, fmt.Errorf("storage.Load: ошибка декодирования JSON: %w", err)
	}
	return docs, nil
}

// StoreGOB сохраняет документы в бинарном формате GOB (нативный формат Go-to-Go).
func StoreGOB(w io.Writer, docs []crawler.Document) error {
	if w == nil {
		return fmt.Errorf("storage.StoreGOB: writer не может быть nil")
	}
	if err := gob.NewEncoder(w).Encode(docs); err != nil {
		return fmt.Errorf("storage.StoreGOB: ошибка кодирования GOB: %w", err)
	}
	return nil
}

// LoadGOB загружает документы из бинарного формата GOB.
func LoadGOB(r io.Reader) ([]crawler.Document, error) {
	if r == nil {
		return nil, fmt.Errorf("storage.LoadGOB: reader не может быть nil")
	}
	var docs []crawler.Document
	if err := gob.NewDecoder(r).Decode(&docs); err != nil {
		return nil, fmt.Errorf("storage.LoadGOB: ошибка декодирования GOB: %w", err)
	}
	return docs, nil
}
```

### 🔍 Прохождение тестов:

```bash
$ cd go_course_4/05-io/homework-05 && go test -v ./...
=== RUN   TestFileExists
--- PASS: TestFileExists (0.00s)
=== RUN   TestStorageWorkflowWithTempFile
--- PASS: TestStorageWorkflowWithTempFile (0.00s)
=== RUN   TestExecuteSearch
--- PASS: TestExecuteSearch (0.00s)
=== RUN   TestScanSitesMembot
--- PASS: TestScanSitesMembot (0.00s)
=== RUN   TestBytesBufferReader
--- PASS: TestBytesBufferReader (0.00s)
PASS
ok      go-core-4/homework-05/cmd/gosearch      0.258s
ok      go-core-4/homework-05/pkg/index         0.216s
ok      go-core-4/homework-05/pkg/storage       0.512s
```

---

## Занятие 6: ООП в Go — Структуры, методы и идиоматичный рефакторинг

> 📖 **Теория:** [Том 2: Раздел 9 (Методы, получатели, встраивание, DIP)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/02-go-core-oop-methods.md#9-методы-и-интерфейсы-methods--interfaces)  
> 💻 **Код лекции:** [06-oop](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/06-oop)  
> 📁 **Каталог решения:** [homework-06](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/06-oop/homework-06)  
> 📦 **Пакеты:** [engine](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/06-oop/homework-06/pkg/engine), [geom](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/06-oop/homework-06/pkg/geom)  
> 🧪 **Тесты:** `go test -v ./lessons/06-oop/homework-06/...`  
> 🗺️ **Путеводитель:** [course-codebase-guide.md (Урок 06)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/course-codebase-guide.md#урок-06-ооп-в-go-методы-композиция-и-dip)

В рамках данного занятия (каталоги `go_course_4/06-oop/5-hw/` и `go_course_4/06-oop/homework-06/`) решаются две ключевые задачи:

1. **Задача №1 (Рефакторинг пакета геометрии):** Исправление стилистических, математических и архитектурных ошибок в исходном коде `hw.go` с выписыванием всех исправленных недостатков.
2. **Задача №2 (ООП-архитектура поискового движка GoSearch):** Превращение процедурного кода поисковика в полноценный объект `engine.Engine`, инкапсулирующий состояние (список документов) и зависимости (краулер, инвертированный индекс), с инициализацией зависимостей на верхнем уровне (`main`) и передачей через конструктор `New()`.

---

### 🎯 Разбор задачи №1: Исходный неидиоматичный код и список дефектов

#### Исходный код (`06-oop/5-hw/hw.go`):

```go
package hw

import (
	"fmt"
	"math"
)

// По условиям задачи, координаты не могут быть меньше 0.

type Geom struct {
	X1, Y1, X2, Y2 float64
}

func (geom Geom) CalculateDistance() (distance float64) {

	if geom.X1 < 0 || geom.X2 < 0 || geom.Y1 < 0 || geom.Y2 < 0 {
		fmt.Println("Координаты не могут быть меньше нуля")
		return -1
	} else {
		distance = math.Sqrt(math.Pow(geom.X2-geom.X1, 2) + math.Pow(geom.Y2-geom.Y1, 2))
	}

	// возврат расстояния между точками
	return distance
}
```

#### 🔍 Детальный перечень исправленных дефектов (согласно *Effective Go* и *Go Code Review Comments*):

| № | Дефект / Антипаттерн | Было | Стало | Почему это антипаттерн в Go |
|---|----------------------|------|-------|-----------------------------|
| 1 | **Имя получателя метода (Receiver Name)** | `func (geom Geom)` | `func (g Geom)` / `func (p Point)` | *Go Code Review Comments (Receiver Names)*: имя ресивера должно быть кратким (1-2 буквы), отражать тип. Запрещено дублировать полное имя типа и использовать `this`/`self`. |
| 2 | **Имя метода (Method Naming)** | `CalculateDistance()` | `Distance()` | *Effective Go (Getters)*: геттеры и вычисляемые характеристики объекта именуются лаконично (`Distance()`, `Area()`). Избыточные префиксы `Calculate...`/`Compute...` засоряют API. |
| 3 | **Побочные эффекты (Side Effects)** | `fmt.Println(...)` | Консольный вывод удален | Библиотечные пакеты никогда не должны самовольно писать в `os.Stdout`, нарушая изоляцию вызывающего приложения. |
| 4 | **Магическое значение вместо ошибки** | `if ... < 0 { return -1 }` | Ограничение удалено | В декартовой геометрии координаты $X, Y \in (-\infty, +\infty)$ законны, а евклидово расстояние всегда $\ge 0$, поэтому возврат `-1` нарушает контракт. ⚠️ Если по условиям конкретной задачи отрицательные координаты действительно недопустимы (как сказано в комментарии к исходному коду), проверку не удаляют, а выносят в валидацию/конструктор и возвращают `error`, а не `-1`. |
| 5 | **Неэффективность вычислений** | `math.Pow(dx, 2)` | `math.Hypot(dx, dy)` | `math.Pow` — универсальная функция возведения в вещественную степень и для квадрата заметно медленнее обычного умножения `dx*dx`. `math.Hypot` корректно обрабатывает крайние случаи (программно избегает переполнения и потери точности на очень больших/малых значениях), поэтому подходит для «общего» расстояния; если важна максимальная скорость и значения ограничены, достаточно `math.Sqrt(dx*dx + dy*dy)`. |
| 6 | **Избыточный `else` после `return`** | `if ... return } else { ... }` | Плоский код (Guard Clause) | Идиома Go: при наличии терминального `return` блок `else` опускается, основной поток выполнения пишется по левому краю без отступов. |
| 7 | **Именованный результат и лишняя переменная** | `(distance float64)` + `distance = ...` | `return math.Hypot(dx, dy)` | Избыточный шум в короткой функции. Прямой возврат значения чище и исключает случайные ошибки инициализации. |
| 8 | **Смешение абстракций (DDD / ООП)** | `Geom{X1, Y1, X2, Y2}` | `Point{X, Y}` + `Distance(Point)` | `Geom` смешивал понятия точки и отрезка. Базовая сущность геометрии — `Point`. Метод `p1.Distance(p2)` естественно выражает расстояние от точки до цели. |

---

### 💻 Эталонное решение задачи №1 (`go_course_4/06-oop/homework-06/pkg/geom/geom.go`):

```go
package geom

import (
	"math"
)

// Point представляет точку на двумерной декартовой плоскости.
type Point struct {
	X float64
	Y float64
}

// Distance вычисляет евклидово расстояние между текущей точкой p и целевой точкой target.
func (p Point) Distance(target Point) float64 {
	return math.Hypot(target.X-p.X, target.Y-p.Y)
}

// Geom представляет пару точек для обратной совместимости со старым API.
type Geom struct {
	X1, Y1, X2, Y2 float64
}

// Distance возвращает расстояние между точками (X1, Y1) и (X2, Y2).
func (g Geom) Distance() float64 {
	p1 := Point{X: g.X1, Y: g.Y1}
	p2 := Point{X: g.X2, Y: g.Y2}
	return p1.Distance(p2)
}

// CalculateDistance сохраняет обратную совместимость с тестами курса.
func (g Geom) CalculateDistance() float64 {
	return g.Distance()
}
```

---

### 🎯 Разбор задачи №2: ООП-архитектура поискового движка GoSearch

Согласно лекции №6:
1. **Когда нужен ООП в Go:** Если пакет описывает сущность, хранящую и изменяющую свое состояние между вызовами из других пакетов (кэш, корзина, поисковый движок), он оформляется в виде структуры с методами.
2. **Внедрение зависимостей (Dependency Injection):** Зависимости (краулер, обратный индекс, хранилище) инициализируются на самом верхнем уровне системы (`main.go`) и передаются в конструктор `New()` подсистемы.
3. **Инкапсуляция:** Документы хранятся во внутреннем неэкспортируемом поле `docs []crawler.Document`. Метод `Documents()` возвращает изолированную копию среза, предотвращая внешнюю непреднамеренную мутацию.

#### Архитектура объекта `Engine` (`pkg/engine/engine.go`):

```go
package engine

import (
	"fmt"
	"io"
	"sort"

	"go-core-4/06-oop/homework-06/pkg/crawler"
	"go-core-4/06-oop/homework-06/pkg/index"
	"go-core-4/06-oop/homework-06/pkg/storage"
)

// Engine инкапсулирует состояние поискового сервиса и его зависимости.
type Engine struct {
	crawler crawler.Interface
	index   *index.Service
	docs    []crawler.Document
}

// New — конструктор поискового движка (DI верхнего уровня).
func New(cr crawler.Interface, idx *index.Service) *Engine {
	if idx == nil {
		idx = index.New()
	}
	return &Engine{
		crawler: cr,
		index:   idx,
		docs:    make([]crawler.Document, 0),
	}
}

// Scan обходит сайты и сохраняет состояние внутри объекта.
func (e *Engine) Scan(urls []string, depth int) error {
	if e.crawler == nil {
		return fmt.Errorf("engine: crawler не инициализирован")
	}

	var allDocs []crawler.Document
	for _, u := range urls {
		scanned, err := e.crawler.Scan(u, depth)
		if err != nil {
			continue
		}
		allDocs = append(allDocs, scanned...)
	}

	for i := range allDocs {
		allDocs[i].ID = i
	}
	sort.Slice(allDocs, func(i, j int) bool {
		return allDocs[i].ID < allDocs[j].ID
	})

	e.docs = allDocs
	e.index.Add(e.docs)
	return nil
}

// Search ищет документы по слову через инвертированный индекс и бинарный поиск.
func (e *Engine) Search(query string) []crawler.Document {
	ids := e.index.Search(query)
	if len(ids) == 0 {
		return nil
	}

	results := make([]crawler.Document, 0, len(ids))
	for _, id := range ids {
		pos := index.BinarySearch(e.docs, id)
		if pos != -1 {
			results = append(results, e.docs[pos])
		}
	}
	return results
}

// Save сохраняет документы в io.Writer.
func (e *Engine) Save(w io.Writer) error {
	return storage.Store(w, e.docs)
}

// Load загружает документы из io.Reader и восстанавливает индекс.
func (e *Engine) Load(r io.Reader) error {
	docs, err := storage.Load(r)
	if err != nil {
		return err
	}
	e.SetDocuments(docs)
	return nil
}
```

---

### 🔍 Прохождение тестов и бенчмарков:

```bash
$ cd go_course_4 && go test -v -bench=. ./lessons/06-oop/homework-06/...
=== RUN   TestFileExists
--- PASS: TestFileExists (0.00s)
=== RUN   TestEngineIntegration
--- PASS: TestEngineIntegration (0.00s)
=== RUN   TestExecuteSearch
--- PASS: TestExecuteSearch (0.00s)
PASS
ok  	go-core-4/homework-06/cmd/gosearch	2.204s
=== RUN   TestEngine_ScanAndSearch
--- PASS: TestEngine_ScanAndSearch (0.00s)
=== RUN   TestEngine_SaveAndLoad
--- PASS: TestEngine_SaveAndLoad (0.00s)
=== RUN   TestEngine_Encapsulation
--- PASS: TestEngine_Encapsulation (0.00s)
PASS
ok  	go-core-4/homework-06/pkg/engine	0.191s
=== RUN   TestPoint_Distance
--- PASS: TestPoint_Distance (0.00s)
=== RUN   TestGeom_BackwardCompatibility
--- PASS: TestGeom_BackwardCompatibility (0.00s)
BenchmarkPoint_Distance-8   	427682140	         2.709 ns/op
PASS
ok  	go-core-4/homework-06/pkg/geom	1.613s
=== RUN   TestService_AddAndSearch
--- PASS: TestService_AddAndSearch (0.00s)
=== RUN   TestBinarySearch
--- PASS: TestBinarySearch (0.00s)
PASS
ok  	go-core-4/homework-06/pkg/index	0.191s
=== RUN   TestJSONStorageRoundtrip
--- PASS: TestJSONStorageRoundtrip (0.00s)
=== RUN   TestGOBStorageRoundtrip
--- PASS: TestGOBStorageRoundtrip (0.00s)
PASS
ok  	go-core-4/homework-06/pkg/storage	0.173s
```

---

## Занятие 7: Тестирование в Go — Unit-тесты, Табличные тесты, Бенчмарки и Имитация зависимостей

> 📖 **Теория:** [Том 5: Разделы 17.1–17.5 (Unit-тесты, Table-Driven, бенчмарки)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/05-go-core-testing-arch.md#17-тестирование-бенчмарки-и-профилирование-testing--pprof)  
> 💻 **Код лекции:** [07-testing](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/07-testing)  
> 📁 **Каталог решения:** [homework-07](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/07-testing/homework-07)  
> 📦 **Пакеты:** [sorts](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/07-testing/homework-07/pkg/sorts), [storage](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/07-testing/homework-07/pkg/storage)  
> 🧪 **Тесты:** `go test -v -bench=. ./lessons/07-testing/homework-07/...`  
> 🗺️ **Путеводитель:** [course-codebase-guide.md (Урок 07)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/course-codebase-guide.md#урок-07-промышленное-тестирование-и-tdd)

### 🎯 Условие задачи:

В рамках занятия №7 курса **«Современная разработка на Go»** (Thinknetica) требуется:
1. **Задача №1:**
   - Написать простой юнит-тест для функции standard library `sort.Ints()`.
   - Написать идиоматичный табличный тест (Table-Driven Test) для функции `sort.Strings()`, проверяющий пограничные и нетривиальные случаи (пустой срез, `nil`, 1 элемент, дубликаты, регистрозависимость, Unicode кириллица, общие префиксы).
2. **Задача №2:**
   - Написать бенчмарки (`testing.B`) для функций `sort.Ints()` и `sort.Float64s()`.
   - Сравнить их производительность на различных объемах входных данных ($N = 10, 100, 1\,000, 10\,000$).
   - Выполнить замер потребления памяти (`-benchmem`: `B/op` и `allocs/op`) и дать техническое обоснование разницы в скорости.
3. **Задача №3 (Архитектура поисковика GoSearch):**
   - Написать набор тестов для проекта GoSearch, защищающий кодовую базу от регрессий при дальнейшей разработке.
   - Изолировать модульные тесты поискового движка `Engine` от реальной внешней сети, используя интерфейс `crawler.Interface` и тестовую имитацию (Mock/Stub `membot`).
   - Покрыть тестами инвертированный индекс (`pkg/index`), алгоритм бинарного поиска (`BinarySearch`), а также персистентное сохранение и загрузку документов (`pkg/storage`).

---

### 💡 Ключевые концепции и архитектурные ловушки:

1. **Мутация входных данных в табличных тестах:**
   - Функции сортировки `sort.Ints` и `sort.Strings` мутируют переданный срез **на месте (in-place)**.
   - Если передать срез напрямую из тестовой структуры `tt.in`, первая же итерация теста навсегда изменит данные в памяти процесса. При повторном запуске теста (или при запуске с флагом `-count 2`) тест будет проверять уже отсортированный массив!
   - *Решение:* Всегда создавать изолированную копию среза через `copy()` перед передачей в мутирующую функцию:
     ```go
     input := make([]string, len(tt.in))
     copy(input, tt.in)
     sort.Strings(input)
     ```

2. **Ловушка `StopTimer()` / `StartTimer()` в быстрых бенчмарках:**
   - Попытка остановить и перезапустить таймер внутри цикла `for i := 0; i < b.N; i++` для подготовки несортированного среза приводит к катастрофическому замедлению теста: каждый вызов `b.StopTimer()` и `b.StartTimer()` совершает системные вызовы замера времени и атомарные операции.
   - *Решение:* Выполнять быструю операцию `copy(buf, master)` прямо в цикле замера. Поскольку копирование среза из 8-байтных чисел для `[]int` и `[]float64` выполняется за идентичное число тактов CPU, базовый оверхед одинаков для обеих функций, а бенчмарк выполняется за доли секунды без системных прерываний.

3. **Почему `sort.Ints` в целом быстрее `sort.Float64s` (и почему на малых $N$ это не всегда заметно):**
   - ℹ️ В таблице замеров ниже при $N = 10$ `Float64s` оказался даже быстрее, а при $N \ge 100$ `Ints` выигрывает: на очень малых объёмах различия тонут в шуме измерений, поэтому выводы делают по большим $N$ и по нескольким прогонам (`-count 5`, сравнение через `benchstat`).
   - **Специфика IEEE 754 и `NaN`:** Значение `NaN` (Not a Number) не равно ничему, включая себя (`NaN == NaN` возвращает `false`), а выражения `NaN < x` и `NaN > x` всегда ложны. Чтобы алгоритм сортировки (pdqsort) не нарушал транзитивность и не зацикливался, `sort.Float64s` обязан проверять наличие `NaN` и выносить их в конец среза.
   - **Инструкции CPU (ALU vs FPU):** Сравнение целых чисел `int` (`CMP` + условные инструкции) выполняется конвейером ALU за 1 такт с практически 100% точностью branch predictor. Сравнение чисел с плавающей точкой в регистрах FPU/SIMD имеет большую латентность инструкции и сложнее оптимизируется процессором.
   - **Аллокации:** Обе функции показывают `0 B/op` и `0 allocs/op`, так как работают in-place.

4. **Тестирование зависимостей через интерфейсы (Interface-based Mocking):**
   - Модульные тесты не должны зависеть от внешних систем (сеть, интернет, файловая система). Реальный краулер `spider` зависит от доступности `go.dev` и скорости соединения.
   - Внедрение зависимости `crawler.Interface` в объект `engine.New(cr, idx)` позволяет в тестах передать in-memory Mock `membot.New()`, который мгновенно возвращает фиксированный набор страниц без сетевых задержек и флакания.

5. **Кэширование тестов в Go (`-count 1`):**
   - По умолчанию Go кэширует успешные результаты `(cached)`. Если код не менялся, повторный запуск `go test` не выполняет тест физически.
   - Для гарантированного прогона всегда используется флаг `-count 1`: `go test -v -count 1 ./...`.

---

### 💻 Эталонное решение (Clean Code / Production-ready):

#### 1. Модульные и табличные тесты сортировок ([`pkg/sorts/sorts_test.go`](../lessons/07-testing/homework-07/pkg/sorts/sorts_test.go)):

```go
package sorts

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
)

// Задача №1 (Часть 1): Простой тест для sort.Ints()
func TestSortInts(t *testing.T) {
	data := []int{42, -5, 0, 17, -100, 8, 3}
	want := []int{-100, -5, 0, 3, 8, 17, 42}

	sort.Ints(data)

	if !sort.IntsAreSorted(data) {
		t.Errorf("sort.IntsAreSorted(data) = false, want true")
	}
	if !reflect.DeepEqual(data, want) {
		t.Errorf("sort.Ints() got %v, want %v", data, want)
	}
}

// Задача №1 (Часть 2): Табличный тест для sort.Strings()
func TestSortStrings(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want []string
	}{
		{name: "пустой срез", in: []string{}, want: []string{}},
		{name: "nil срез", in: nil, want: nil},
		{name: "один элемент", in: []string{"golang"}, want: []string{"golang"}},
		{name: "уже отсортированный срез", in: []string{"apple", "banana", "cherry", "date"}, want: []string{"apple", "banana", "cherry", "date"}},
		{name: "обратный порядок", in: []string{"zebra", "tiger", "lion", "cat"}, want: []string{"cat", "lion", "tiger", "zebra"}},
		{name: "дубликаты и повторяющиеся строки", in: []string{"go", "rust", "go", "python", "rust", "c"}, want: []string{"c", "go", "go", "python", "rust", "rust"}},
		{name: "регистрозависимость (ASCII порядок)", in: []string{"banana", "Apple", "apple", "Banana"}, want: []string{"Apple", "Banana", "apple", "banana"}},
		{name: "кириллица (Unicode UTF-8)", in: []string{"яблоко", "арбуз", "банан", "груша"}, want: []string{"арбуз", "банан", "груша", "яблоко"}},
		{name: "общие префиксы разной длины", in: []string{"testing", "test", "tester", "testcase"}, want: []string{"test", "testcase", "tester", "testing"}},
		{name: "спецсимволы, цифры и пробелы", in: []string{" 10", "!important", "01", " 02", "@user"}, want: []string{" 02", " 10", "!important", "01", "@user"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var input []string
			if tt.in != nil {
				input = make([]string, len(tt.in))
				copy(input, tt.in)
			}

			sort.Strings(input)

			if !sort.StringsAreSorted(input) {
				t.Errorf("StringsAreSorted() = false for input %v", input)
			}
			if !reflect.DeepEqual(input, tt.want) {
				t.Errorf("sort.Strings() = %v, want %v", input, tt.want)
			}
		})
	}
}

// Задача №2: Бенчмарки sort.Ints() и sort.Float64s()
func BenchmarkSortInts(b *testing.B) {
	sizes := []int{10, 100, 1_000, 10_000}
	for _, size := range sizes {
		b.Run(fmt.Sprintf("N=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			master := GenerateInts(size)
			buf := make([]int, size)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				copy(buf, master)
				sort.Ints(buf)
			}
		})
	}
}

func BenchmarkSortFloat64s(b *testing.B) {
	sizes := []int{10, 100, 1_000, 10_000}
	for _, size := range sizes {
		b.Run(fmt.Sprintf("N=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			master := GenerateFloat64s(size)
			buf := make([]float64, size)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				copy(buf, master)
				sort.Float64s(buf)
			}
		})
	}
}
```

#### 2. Мокирование зависимостей краулера ([`pkg/crawler/membot/membot.go`](../lessons/07-testing/homework-07/pkg/crawler/membot/membot.go)):

```go
package membot

import "go-core-4/07-testing/homework-07/pkg/crawler"

// Service реализует crawler.Interface полностью в памяти без обращения к сети.
type Service struct{}

func New() *Service {
	return &Service{}
}

func (s *Service) Scan(url string, depth int) ([]crawler.Document, error) {
	return []crawler.Document{
		{ID: 0, URL: "https://go.dev", Title: "The Go Programming Language documents and tutorials"},
		{ID: 1, URL: "https://golang.org", Title: "The Go Programming Language official documents"},
		{ID: 2, URL: "https://yandex.ru", Title: "Яндекс - быстрый поиск в интернете"},
		{ID: 3, URL: "https://google.com", Title: "Google Search Engine"},
	}, nil
}
```

#### 3. Тестирование поискового движка ([`pkg/engine/engine_test.go`](../lessons/07-testing/homework-07/pkg/engine/engine_test.go)):

```go
package engine

import (
	"bytes"
	"testing"

	"go-core-4/07-testing/homework-07/pkg/crawler/membot"
	"go-core-4/07-testing/homework-07/pkg/index"
)

func TestEngine_ScanAndSearch(t *testing.T) {
	// Изолированное тестирование через Mock-краулер membot
	bot := membot.New()
	idx := index.New()
	eng := New(bot, idx)

	if err := eng.Scan([]string{"https://go.dev"}, 1); err != nil {
		t.Fatalf("eng.Scan() failed: %v", err)
	}

	results := eng.Search("golang")
	if len(results) == 0 {
		t.Errorf("eng.Search(\"golang\") returned 0 results, want > 0")
	}

	if empty := eng.Search("unknownword"); len(empty) != 0 {
		t.Errorf("expected 0 results, got %v", empty)
	}
}

func TestEngine_SaveAndLoad(t *testing.T) {
	bot := membot.New()
	eng1 := New(bot, nil)
	_ = eng1.Scan([]string{""}, 1)

	buf := new(bytes.Buffer)
	if err := eng1.Save(buf); err != nil {
		t.Fatalf("eng1.Save() failed: %v", err)
	}

	eng2 := New(nil, nil)
	if err := eng2.Load(buf); err != nil {
		t.Fatalf("eng2.Load() failed: %v", err)
	}

	if eng2.Count() != eng1.Count() {
		t.Errorf("eng2.Count() = %d, want %d", eng2.Count(), eng1.Count())
	}
}
```

---

### 🔍 Результаты замеров производительности и прохождения тестов:

#### 1. Сравнительная таблица бенчмарков `sort.Ints` vs `sort.Float64s`:

| Бенчмарк | Размер $N$ | Выполнено итераций | Время (`ns/op`) | Память (`B/op`) | Аллокации (`allocs/op`) |
| :--- | :---: | :---: | :---: | :---: | :---: |
| `BenchmarkSortInts` | **10** | 3 868 743 | **52.84 ns/op** | 0 B/op | 0 allocs/op |
| `BenchmarkSortFloat64s` | **10** | 3 252 450 | **42.12 ns/op** | 0 B/op | 0 allocs/op |
| `BenchmarkSortInts` | **100** | 184 539 | **626.5 ns/op** | 0 B/op | 0 allocs/op |
| `BenchmarkSortFloat64s` | **100** | 147 826 | **787.4 ns/op** | 0 B/op | 0 allocs/op |
| `BenchmarkSortInts` | **1 000** | 12 744 | **9 645 ns/op** | 0 B/op | 0 allocs/op |
| `BenchmarkSortFloat64s` | **1 000** | 7 758 | **16 249 ns/op** | 0 B/op | 0 allocs/op |
| `BenchmarkSortInts` | **10 000** | 286 | **420 378 ns/op** | 0 B/op | 0 allocs/op |
| `BenchmarkSortFloat64s` | **10 000** | 213 | **572 177 ns/op** | 0 B/op | 0 allocs/op |

#### 2. Лог прогона тестов и покрытия кода:

```bash
$ cd go_course_4 && go test -v -cover -bench=. -benchmem -benchtime=100ms ./lessons/07-testing/homework-07/...
=== RUN   TestExecuteSearch
--- PASS: TestExecuteSearch (0.00s)
PASS
ok  	go-core-4/homework-07/cmd/gosearch	0.234s
=== RUN   TestEngine_ScanAndSearch
--- PASS: TestEngine_ScanAndSearch (0.00s)
=== RUN   TestEngine_SaveAndLoad
--- PASS: TestEngine_SaveAndLoad (0.00s)
=== RUN   TestEngine_Encapsulation
--- PASS: TestEngine_Encapsulation (0.00s)
=== RUN   TestEngine_NilCrawlerError
--- PASS: TestEngine_NilCrawlerError (0.00s)
=== RUN   TestEngine_SaveLoadErrors
--- PASS: TestEngine_SaveLoadErrors (0.00s)
PASS
coverage: 96.3% of statements
ok  	go-core-4/homework-07/pkg/engine	0.200s
=== RUN   TestService_AddAndSearch
--- PASS: TestService_AddAndSearch (0.00s)
=== RUN   TestBinarySearch_Table
--- PASS: TestBinarySearch_Table (0.00s)
PASS
coverage: 100.0% of statements
ok  	go-core-4/homework-07/pkg/index	0.196s
=== RUN   TestSortInts
--- PASS: TestSortInts (0.00s)
=== RUN   TestSortStrings
--- PASS: TestSortStrings (0.00s)
PASS
coverage: 100.0% of statements
ok  	go-core-4/homework-07/pkg/sorts	1.712s
=== RUN   TestJSONStorageRoundtrip
--- PASS: TestJSONStorageRoundtrip (0.00s)
=== RUN   TestGOBStorageRoundtrip
--- PASS: TestGOBStorageRoundtrip (0.00s)
=== RUN   TestStorage_Errors
--- PASS: TestStorage_Errors (0.00s)
PASS
coverage: 92.6% of statements
ok  	go-core-4/homework-07/pkg/storage	0.202s
```

---

## Домашнее задание 8: Профилирование, Отладка и Трассировка в Go

> 📖 **Теория:** [Том 5: Разделы 17.6–17.9 (pprof, трассировка)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/05-go-core-testing-arch.md#17-тестирование-бенчмарки-и-профилирование-testing--pprof)  
> 💻 **Код лекции:** [08-prof_debug](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/08-prof_debug)  
> 📁 **Каталог решения:** [homework-08](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/08-prof_debug/homework-08)  
> 🧪 **Тесты:** `go test -v -bench=. ./lessons/08-prof_debug/homework-08/...`  
> 🗺️ **Путеводитель:** [course-codebase-guide.md (Урок 08)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/course-codebase-guide.md#урок-08-профилирование-бенчмарки-и-трейсинг)

Решение полностью реализовано в директории [`go_course_4/08-prof_debug/homework-08`](../lessons/08-prof_debug/homework-08).

#### 1. Отладка приложения в Delve ([`pkg/twosum/twosum.go`](../lessons/08-prof_debug/homework-08/pkg/twosum/twosum.go)):
- **Локализация бага**: В исходном варианте `TwoSumBuggy` вложенный цикл `for j := 0; j < len(nums); j++` при `nums = [2, 10, 6, 12, 14], target = 20` брал один и тот же элемент `nums[1] = 10` дважды и возвращал `[1, 1]` вместо `[2, 4]`.
- **Сессия Delve CLI**:
  ```text
  (dlv) break twosum.go:10
  (dlv) condition 1 i == j && nums[i] + nums[j] == target
  (dlv) continue
  (dlv) locals
  (dlv) print nums[i]
  ```
- **Исправление**:
  1. `TwoSumBruteForce`: внутренний цикл начинается с `j := i + 1` ($O(N^2)$ время, $O(1)$ память).
  2. `TwoSumOptimal`: сохранение дополнения `target - num` в хеш-таблице `map[int]int` ($O(N)$ время, $O(N)$ память).
- **Бенчмарки**:
  ```text
  BenchmarkTwoSumBruteForce_1000-8    5846    205892 ns/op
  BenchmarkTwoSumOptimal_1000-8      81206     14813 ns/op
  ```
  Оптимальный алгоритм на хеш-таблице работает в **14 раз быстрее**.

#### 2. Профилирование процессора и памяти ([`pkg/search`](../lessons/08-prof_debug/homework-08/pkg/search), [`cmd/app_profile`](../lessons/08-prof_debug/homework-08/cmd/app_profile)):
- **Снятие профилей бенчмарка**:
  ```bash
  go test -bench=BenchmarkBinarySearch_100k -benchmem -cpuprofile=cpu.out -memprofile=mem.out ./lessons/08-prof_debug/homework-08/pkg/search
  ```
  Замер показал: `48.54 ns/op, 0 B/op, 0 allocs/op`.
- **Инспекция `go tool pprof -text cpu.out`**:
  Функция `Binary` занимает **86.59%** чистого процессорного времени (`flat`), отсутствуют паразитные аллокации.
- **Профилирование сервиса через `net/http/pprof`**:
  Сервер на `:6060` предоставляет эндпоинты `/debug/pprof/profile` и `/debug/pprof/heap`. Анализ разделяет утечки (`inuse_space`) и нагрузку на GC (`alloc_space`).

#### 3. Трассировка выполнения (`runtime/trace`) ([`cmd/app_trace`](../lessons/08-prof_debug/homework-08/cmd/app_trace)):
- Запись журнала событий через `trace.Start(f)` и `trace.Stop()`.
- Просмотр через `go tool trace trace.out`.

---

## Домашнее задание 9: Интерфейсы в Go

> 📖 **Теория:** [Том 2: Раздел 9 (Методы и интерфейсы)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/02-go-core-oop-methods.md#9-методы-и-интерфейсы-methods--interfaces)  
> 💻 **Код лекции:** [09-interfaces](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/09-interfaces)  
> 📁 **Каталог решения:** [homework-09](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/09-interfaces/homework-09)  
> 🧪 **Тесты:** `go test -v ./lessons/09-interfaces/homework-09/...`  
> 🗺️ **Путеводитель:** [course-codebase-guide.md (Урок 09)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/course-codebase-guide.md#урок-09-интерфейсы-и-полиморфизм)

Решение полностью реализовано и протестировано в директории [`go_course_4/09-interfaces/homework-09`](../lessons/09-interfaces/homework-09).

#### 1. Полиморфизм через общий контракт интерфейса ([`pkg/users/users.go`](../lessons/09-interfaces/homework-09/pkg/users/users.go)):
- Интерфейс `Ager` с методом `Age() int` и утиная типизация для `Employee` и `Customer`.
- Функция `MaxAge(people ...Ager) int` находит максимальный возраст среди произвольных типов.

#### 2. Обобщенная обработка без методов через `any` и `Type Switch`:
- Структуры без методов `EmployeeSimple` и `CustomerSimple`.
- Извлечение данных и возврат самого объекта через функцию `OldestPerson(people ...any) any`.

#### 3. Фильтрация типов и потоковый вывод в `io.Writer`:
- Функция `WriteOnlyStrings(w io.Writer, args ...any) (int, error)` проверяет типы через `switch v := arg.(type)` и пишет строковые аргументы в поток.

#### 4. Тестирование ([`pkg/users/users_test.go`](../lessons/09-interfaces/homework-09/pkg/users/users_test.go)):
- Покрытие: 100% statements, 0 data races.

---

## Домашнее задание 10: Конкурентное программирование (Игра в Пинг-Понг)

> 📖 **Теория:** [Том 3: Разделы 11, 15 (Конкурентность, горутины, каналы, sync)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/03-go-core-concurrency.md#11-конкурентность-горутины-и-каналы-concurrency)  
> 💻 **Код лекции:** [10-concurrency](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/10-concurrency)  
> 📁 **Каталог решения:** [homework-10](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/10-concurrency/homework-10)  
> 🧪 **Тесты:** `go test -v -race ./lessons/10-concurrency/homework-10/...`  
> 🗺️ **Путеводитель:** [course-codebase-guide.md (Урок 10)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/course-codebase-guide.md#урок-10-конкурентность-горутины-и-каналы)

Решение полностью реализовано и протестировано в директории [`go_course_4/10-concurrency/homework-10`](../lessons/10-concurrency/homework-10).

#### 1. Игровой движок на горутинах и каналах ([`pkg/pingpong/pingpong.go`](../lessons/10-concurrency/homework-10/pkg/pingpong/pingpong.go)):
- Два независимых игрока в отдельных горутинах обмениваются строковыми сообщениями (`"ping"`, `"pong"`) через небуферизованный канал `table`.
- Подача начинается с команды `"begin"`.
- 20% вероятность критического удара в угол (`"stop"`), приносящего победное очко.
- **Синхронизация через `sync.WaitGroup`**: судья ведет счет до `TargetScore` (по умолчанию 11 очков), закрывает канал `close(table)` и через `wg.Wait()` гарантирует корректное завершение обеих горутин без утечек памяти (`Goroutine Leaks`).

#### 2. Консольное приложение матча ([`cmd/pingpong/main.go`](../lessons/10-concurrency/homework-10/cmd/pingpong/main.go)):
- Поддерживает флаги командной строки: `-p1` (имя первого игрока), `-p2` (имя второго игрока), `-score` (очки для победы), `-delay` (задержка между ударами для зрелищности).

#### 3. Тестирование ([`pkg/pingpong/pingpong_test.go`](../lessons/10-concurrency/homework-10/pkg/pingpong/pingpong_test.go)):
- Покрытие пакета: **98.8%** statements.
- Полная безопасность конкурентности: **0 race conditions** при запуске `go test -race`.

---

## Домашнее задание 11: Сетевое программирование (Сетевая служба GoSearch)

> 📖 **Теория:** [Том 4: Раздел 28 (Сетевое программирование, TCP, UDP, сокеты)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/04-go-core-backend-web.md#28-сетевое-программирование-network-programming)  
> 💻 **Код лекции:** [11-network](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/11-network)  
> 📁 **Каталог решения:** [homework-11](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/11-network/homework-11)  
> 🧪 **Тесты:** `go test -v -race ./lessons/11-network/homework-11/...`  
> 🗺️ **Путеводитель:** [course-codebase-guide.md (Урок 11)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/course-codebase-guide.md#урок-11-сетевое-программирование-tcpudp)

Решение полностью реализовано и протестировано в директории [`go_course_4/11-network/homework-11`](../lessons/11-network/homework-11).

#### 1. Очистка поисковика от устаревших артефактов (Задача №1, [`cmd/netsrv/main.go`](../lessons/11-network/homework-11/cmd/netsrv/main.go)):
- Полностью удалены флаги командной строки (`-s`), интерактивный ввод `os.Stdin` в сервере и рудиментарные файлы.
- Сервер фокусируется исключительно на сетевом взаимодействии: сканирует/индексирует документы, инициализирует движок и запускает службу `netsrv`.
- Реализован Graceful Shutdown по сигналам `SIGINT` и `SIGTERM`.

#### 2. Сетевой пакет `netsrv` (Задача №2, [`pkg/netsrv/server.go`](../lessons/11-network/homework-11/pkg/netsrv/server.go)):
- Служба принимает подключения на порту 8000 на всех сетевых интерфейсах (`:8000` / `0.0.0.0:8000`).
- **Конкурентная обработка:** каждое входящее клиентское соединение передается в отдельную горутину `go s.HandleConn(conn)`.
- **Потоковый протокол:** читает текстовые запросы от клиентов, ищет документы через абстракцию `Searcher`, форматирует результаты и отправляет обратно в сокет.
- **Безопасность:** выставляется таймаут на бездействие (`conn.SetDeadline`), корректно обрабатываются команды завершения (`exit`/`quit`) и обрыв связи.

#### 3. Интерактивный консольный клиент (Задача №3, [`cmd/netclient/main.go`](../lessons/11-network/homework-11/cmd/netclient/main.go)):
- Подключается по TCP (`net.Dial`) к `localhost:8000`.
- В интерактивном цикле передает поисковые запросы из `os.Stdin` на сервер и асинхронно транслирует ответы сервера пользователю.
- Поддерживает подключение как специализированным клиентом, так и стандартной утилитой `telnet localhost 8000`.

#### 4. Тестирование сетевой службы в памяти ([`pkg/netsrv/server_test.go`](../lessons/11-network/homework-11/pkg/netsrv/server_test.go)):
- Тестирование реализовано с применением **`net.Pipe()`** в полном соответствии со слайдом 22.
- Обеспечивает мгновенные изолированные тесты без конфликта занятых портов ОС (`bind: address already in use`).
- Покрытие пакета: **88.6%** statements, **0 data races** под `-race`.

---

## Домашнее задание 12: Веб-приложения на Go (Веб-служба GoSearch)

> 📖 **Теория:** [Том 4: Раздел 12 (Backend Web-разработка, HTTP, routing, templates)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/04-go-core-backend-web.md#12-backend-web-разработка)  
> 💻 **Код лекции:** [12-web-apps](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/12-web-apps)  
> 📁 **Каталог решения:** [homework-12](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/12-web-apps/homework-12)  
> 🧪 **Тесты:** `go test -v -race ./lessons/12-web-apps/homework-12/...`  
> 🗺️ **Путеводитель:** [course-codebase-guide.md (Урок 12)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/course-codebase-guide.md#урок-12-веб-приложения-и-шаблонизация)

Решение полностью реализовано и протестировано в отдельной директории [`go_course_4/12-web-apps/homework-12`](../lessons/12-web-apps/homework-12) без изменения или затирания предыдущих домашних работ.

#### 1. Веб-пакет `webapp` (Задача №1, [`pkg/webapp/webapp.go`](../lessons/12-web-apps/homework-12/pkg/webapp/webapp.go)):
- **Маршрутизация на базе `github.com/gorilla/mux`:**
  - Привязаны эндпоинты `GET /docs` и `GET /index`.
  - Строгая валидация методов: для любых других HTTP-методов (`POST`, `PUT`, `DELETE`) автоматически возвращается статус `405 Method Not Allowed`.
- **Шаблонизация `html/template` (Server-Side Rendering):**
  - Шаблоны компилируются единожды при старте через `template.Must(template.New().Parse())`.
  - Автоматическое экранирование предотвращает атаки XSS (Cross-Site Scripting).
  - Стилизованные таблицы с бейджами, подсветкой строк и адаптивной версткой.
- **Поддержка согласования контента (Content Negotiation):**
  - При открытии в обычном веб-браузере отдается полноценная HTML-страница (`Content-Type: text/html; charset=utf-8`).
  - При наличии query-параметра `?format=json` или HTTP-заголовка `Accept: application/json` отдается структурированный JSON (`Content-Type: application/json; charset=utf-8`), что превращает веб-пакет также в полноценный REST API.
- **Потокобезопасное хранилище `SearchDB` ([`pkg/webapp/db.go`](../lessons/12-web-apps/homework-12/pkg/webapp/db.go)):**
  - Реализует контракт `DataProvider` с инвертированным поисковым индексом и хранилищем документов.
  - Защищено `sync.RWMutex` для безопасного конкурентного чтения и обновления данных.

#### 2. Модульные тесты обработчиков (Задача №2, [`pkg/webapp/webapp_test.go`](../lessons/12-web-apps/homework-12/pkg/webapp/webapp_test.go)):
- Реализованы через стандартный пакет `net/http/httptest` (`httptest.NewRequest`, `httptest.NewRecorder`).
- Запросы отправляются напрямую в метод `router.ServeHTTP(rec, req)` без открытия реальных TCP-сокетов ОС, что исключает флаки-тесты и конфликты портов.
- **Протестированные сценарии:**
  - `TestHandleDocs_HTML` и `TestHandleDocs_JSON` (проверка статус-кода 200, Content-Type, структуры JSON, названий и URL).
  - `TestHandleIndex_HTML` и `TestHandleIndex_JSON` (проверка рендеринга токенов, количества и ID документов).
  - `TestMethodNotAllowed` (проверка отклонения недопустимых методов со статусом 405).
  - `TestEmptyData` (корректная обработка пустого индекса и пустой базы документов).
  - `TestNilProvider` (обработка системных ошибок со статусом 500).
  - `TestNotFoundRoute` (обработка несуществующих путей со статусом 404).
- **Результаты тестирования:** Покрытие пакета составляет **98.6%** statements, **0 data races** при проверке с флагом `-race`.

#### 3. Серверное приложение ([`cmd/webapp/main.go`](../lessons/12-web-apps/homework-12/cmd/webapp/main.go)):
- Инициализирует поисковую базу с демонстрационными данными.
- Настраивает промышленный `http.Server` с таймаутами (`ReadTimeout`, `WriteTimeout`, `IdleTimeout`).
- Реализует корректное завершение работы (**Graceful Shutdown**) по сигналам `os.Interrupt` и `syscall.SIGTERM`.

---

## Домашнее задание 13: Разработка REST API (Поисковик GoSearch API и Модель памяти)

> 📖 **Теория:** [Том 4: Раздел 12 (Backend Web-разработка, REST API)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/04-go-core-backend-web.md#12-backend-web-разработка) | [Том 1: Раздел 1 (Модель памяти Go)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/01-go-core-basics.md#1-внутреннее-устройство-go-и-система-типов)  
> 💻 **Код лекции:** [13-api](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/13-api)  
> 📁 **Каталог решения:** [homework-13](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/13-api/homework-13)  
> 🧪 **Тесты:** `go test -v -race ./lessons/13-api/homework-13/...`  
> 🗺️ **Путеводитель:** [course-codebase-guide.md (Урок 13)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/course-codebase-guide.md#урок-13-rest-api-и-промежуточное-по-middleware)

Решение полностью реализовано и протестировано в отдельной изолированной директории [`go_course_4/13-api/homework-13`](../lessons/13-api/homework-13) без изменения или затирания предыдущих решений.

#### 1. REST API поисковика GoSearch (Задача №1):
- **Пакет хранения данных ([`pkg/storage/storage.go`](../lessons/13-api/homework-13/pkg/storage/storage.go)):**
  - Реализует потокобезопасную структуру `DB` с `sync.RWMutex`, инвертированным поисковым индексом `map[string][]int` и коллекцией документов `map[int]Document`.
  - Предоставляет методы интерфейса `storage.Interface`:
    - `Search(query string) []Document` — полнотекстовый поиск по токенизированным словам с дедупликацией и сортировкой по ID.
    - `Documents() []Document` — возврат среза всех документов.
    - `Document(id int) (Document, error)` — получение документа по его идентификатору (возвращает `ErrNotFound` при отсутствии).
    - `Add(doc Document) Document` — создание документа с автогенерацией уникального ID (`maxID++`) и добавлением в инвертированный индекс.
    - `Update(doc Document) error` — полное обновление документа (по условию перестраивать индекс не требуется).
    - `Delete(id int) error` — удаление документа из хранилища.
- **Пакет API ([`pkg/api/api.go`](../lessons/13-api/homework-13/pkg/api/api.go)):**
  - Спроектирован в ООП-стиле: структура `API` инкапсулирует `*mux.Router` и `storage.Interface`.
  - Конструктор `New(db storage.Interface) *API` конфигурирует цепочки middleware и регистрирует endpoints.
  - **Спецификация REST-маршрутов:**
    - `GET /api/v1/search?q={query}` — поиск документов (возвращает массив `[]Document`, статус `200 OK`; при пустом параметре — `400 Bad Request`).
    - `GET /api/v1/docs` — получение списка всех документов (`200 OK`).
    - `GET /api/v1/docs/{id:[0-9]+}` — получение документа по ID (`200 OK`, либо `404 Not Found`).
    - `POST /api/v1/docs` — создание нового документа из JSON тела (`201 Created`, валидация обязательных полей `title` и `url`).
    - `PUT /api/v1/docs/{id:[0-9]+}` — обновление документа (`200 OK`, либо `404 Not Found`).
    - `DELETE /api/v1/docs/{id:[0-9]+}` — удаление документа (`200 OK`, либо `404 Not Found`).
- **Сквозное промежуточное ПО ([`pkg/api/middleware.go`](../lessons/13-api/homework-13/pkg/api/middleware.go)):**
  - `headersMiddleware`: автоматически выставляет заголовок `Content-Type: application/json; charset=utf-8`.
  - `requestIDMiddleware`: читает заголовок `X-Request-ID` от клиента или генерирует криптографически случайный идентификатор, внедряет его в контекст запроса через `context.WithValue` и `r.WithContext(ctx)`, а также выставляет в заголовках ответа.
  - `loggingMiddleware`: логирует метод, URL, Request ID и продолжительность обработки каждого HTTP-запроса.

#### 2. Модульное тестирование методов API (Задача №2):
- **Тесты пакета API ([`pkg/api/api_test.go`](../lessons/13-api/homework-13/pkg/api/api_test.go)):**
  - Инициализация тестового окружения выполнена через стандартную функцию `TestMain(m *testing.M)` в соответствии со слайдом 13 презентации курса.
  - Тестирование обработчиков изолировано в памяти через `httptest.NewRequest` и `httptest.NewRecorder` без поднятия TCP-портов.
  - Проверены сценарии:
    - Полнотекстовый поиск (валидный запрос, альтернативный параметр `?query=`, пустой запрос, отсутствие совпадений).
    - CRUD-методы: создание с валидным/битым JSON и пустыми полями, чтение по существующему/несуществующему ID, обновление и удаление.
    - Проверка фильтрации методов: неподдерживаемые методы возвращают `405 Method Not Allowed`.
  - Покрытие тестами: **87.5%** statements для пакета `api`, **99.0%** statements для пакета `storage`, **0 data races** при проверке с флагом `-race`.

#### 3. Разбор Задачи №3 (Модель памяти в Go и `size.Of`):

- **Ссылка на задачу:** https://go.dev/play/p/_an3p5w9qB9
- **Исполняемая программа с подробным выводом:** [`cmd/task3/main.go`](../lessons/13-api/homework-13/cmd/task3/main.go).

##### Исходный код структуры:
```go
type example struct {
    a []int
    b bool
    c int32
    d string
}
```

##### 1. Расчет выравнивания и поверхностного размера (`unsafe.Sizeof`):

| Поле | Тип | Размер | Выравнивание | Смещение в памяти | Описание |
| :--- | :--- | :---: | :---: | :---: | :--- |
| `a` | `[]int` | 24 байта | 8 байт | **0 .. 23** | Заголовок среза: указатель `Data` (8B) + `Len` (8B) + `Cap` (8B) |
| `b` | `bool` | 1 байт | 1 байт | **24** | Логический флаг |
| *padding* | — | 3 байта | — | **25 .. 27** | Поле `c int32` требует выравнивания кратно 4 (`Alignof = 4`). Смещение 28 кратно 4. |
| `c` | `int32` | 4 байта | 4 байта | **28 .. 31** | 32-битное целое число |
| `d` | `string` | 16 байт | 8 байт | **32 .. 47** | Заголовок строки: указатель `Data` (8B) + `Len` (8B). Смещение 32 кратно 8. |

Общий поверхностный размер структуры (`unsafe.Sizeof(example)`): **48 байт** (кратно 8, tail padding = 0).

##### 2. Заполненный код с подстановкой значений вместо `?`:

```go
ex := example{
    a: []int{1, 2, 3}, // 24 байта (заголовок слайса) + 24 байта (массив из 3 int по 8 байт в куче) = 48 байт
    b: true,           // 1 байт (внутри структуры)
    d: "1234",         // 16 байт (заголовок строки) + 4 байта (символы UTF-8) = 20 байт
} // unsafe.Sizeof: 48 байт | size.Of: 76 байт

ex1 := example{
    a: []int{1, 2, 3}, // 24 байта (заголовок слайса) + 24 байта (массив из 3 int по 8 байт в куче) = 48 байт
    b: true,           // 1 байт (внутри структуры)
    d: "1234",         // 16 байт (заголовок строки) + 4 байта (символы UTF-8) = 20 байт
    c: 100,            // 4 байта (внутри структуры на смещении 28..31)
} // unsafe.Sizeof: 48 байт | size.Of: 76 байт
```

##### 3. Архитектурное объяснение равенства `size.Of(ex) == 76` и `size.Of(ex1) == 76`:
1. Библиотека `github.com/DmitriyVTitov/size` считает **глубокий размер (Deep Size)** в байтах:
   $$\text{Размер структуры (48)} + \text{Данные среза } a \text{ (} 3 \times 8 = 24 \text{)} + \text{Байты строки } d \text{ (} 4 \times 1 = 4 \text{)} = 48 + 24 + 4 = \mathbf{76 \text{ байт}}.$$
2. Поле `c int32` **уже зарезервировано компилятором Go в 48 байтах структуры** на этапе компиляции. В `ex` поле не инициализировано и заполнено нулем (`zero value: 0`), а в `ex1` заполнено значением `100`. Скалярный тип `int32` не создает никаких дополнительных аллокаций в куче, поэтому глубокий размер обеих переменных идентичен и равен **76 байт**.

---

#### 4. Лог выполнения тестов и замеров покрытия:

```bash
$ cd go_course_4 && go test -v -race -cover ./lessons/13-api/homework-13/...
=== RUN   TestTask3_Sizes
--- PASS: TestTask3_Sizes (0.00s)
PASS
ok  	go-core-4/homework-13/cmd/task3	2.135s	coverage: 0.0% of statements
=== RUN   TestAPI_Search
=== RUN   TestAPI_Search/Search_with_match
=== RUN   TestAPI_Search/Search_with_alternative_query_param
=== RUN   TestAPI_Search/Search_with_no_matches
=== RUN   TestAPI_Search/Search_with_empty_query
--- PASS: TestAPI_Search (0.00s)
=== RUN   TestAPI_Docs
--- PASS: TestAPI_Docs (0.00s)
=== RUN   TestAPI_CreateDoc
=== RUN   TestAPI_CreateDoc/Valid_Document_Creation
=== RUN   TestAPI_CreateDoc/Invalid_JSON_Payload
=== RUN   TestAPI_CreateDoc/Empty_Required_Fields
--- PASS: TestAPI_CreateDoc (0.00s)
=== RUN   TestAPI_DocByID
=== RUN   TestAPI_DocByID/Existing_Document
=== RUN   TestAPI_DocByID/Non-existing_Document
--- PASS: TestAPI_DocByID (0.00s)
=== RUN   TestAPI_UpdateDoc
=== RUN   TestAPI_UpdateDoc/Successful_Update
=== RUN   TestAPI_UpdateDoc/Update_Non-existing_Document
=== RUN   TestAPI_UpdateDoc/Update_with_Invalid_JSON
--- PASS: TestAPI_UpdateDoc (0.00s)
=== RUN   TestAPI_DeleteDoc
=== RUN   TestAPI_DeleteDoc/Delete_Existing_Document
=== RUN   TestAPI_DeleteDoc/Delete_Already_Deleted_Document
--- PASS: TestAPI_DeleteDoc (0.00s)
=== RUN   TestAPI_MethodNotAllowed
--- PASS: TestAPI_MethodNotAllowed (0.00s)
=== RUN   TestMiddlewares
--- PASS: TestMiddlewares (0.00s)
PASS
coverage: 87.5% of statements
ok  	go-core-4/homework-13/pkg/api	2.279s	coverage: 87.5% of statements
=== RUN   TestStorage_CRUD
--- PASS: TestStorage_CRUD (0.00s)
=== RUN   TestStorage_DocumentsAndSearch
--- PASS: TestStorage_DocumentsAndSearch (0.00s)
=== RUN   TestStorage_Concurrency
--- PASS: TestStorage_Concurrency (0.00s)
PASS
coverage: 99.0% of statements
ok  	go-core-4/homework-13/pkg/storage	2.460s	coverage: 99.0% of statements
```

---

## Домашнее задание 14: Удалённый вызов процедур (RPC-служба сообщений)

> 📖 **Теория:** [Том 5: Раздел 20 (Удаленный вызов процедур RPC / gRPC)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/05-go-core-testing-arch.md#20-удаленный-вызов-процедур-rpc--grpc)  
> 💻 **Код лекции:** [14-RPC](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/14-RPC)  
> 📁 **Каталог решения:** [homework-14](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/14-RPC/homework-14)  
> 🧪 **Тесты:** `go test -v -race ./lessons/14-RPC/homework-14/...`  
> 🗺️ **Путеводитель:** [course-codebase-guide.md (Урок 14)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/course-codebase-guide.md#урок-14-удаленный-вызов-процедур-rpc)

Решение полностью реализовано и протестировано в отдельной изолированной директории [`go_course_4/14-RPC/homework-14`](../lessons/14-RPC/homework-14) без изменения или затирания предыдущих решений.

#### 1. Архитектура решения и модель данных:
- **Структура сообщения ([`pkg/messages/messages.go`](../lessons/14-RPC/homework-14/pkg/messages/messages.go)):**
  - Содержит обязательные поля по условию задачи:
    - `ID int` — числовой идентификатор сообщения;
    - `Time time.Time` — точная метка времени отправки;
    - `Text string` — строковое содержимое сообщения.
- **RPC-сервис сообщений ([`pkg/rpcservice/service.go`](../lessons/14-RPC/homework-14/pkg/rpcservice/service.go)):**
  - Реализован строго по контракту стандартного пакета Go `net/rpc`:
    1. `Service` — экспортируемый тип данных структуры.
    2. Потокобезопасность обеспечена через `sync.RWMutex`.
    3. Метод `Send(req []messages.Message, resp *int) error` — принимает массив сообщений, конкатенирует их с внутренним слайсом сервиса и записывает в `*resp` количество добавленных элементов.
    4. Метод `Messages(req struct{}, resp *[]messages.Message) error` — возвращает клиенту полную копию среза всех накопленных сообщений.

#### 2. Серверное приложение ([`cmd/server/main.go`](../lessons/14-RPC/homework-14/cmd/server/main.go)):
- Регистрирует сервис через `rpc.Register(svc)`.
- Открывает TCP-слушатель на порту `:8080` (`net.Listen("tcp", ":8080")`).
- Каждое входящее подключение обслуживается конкурентно в отдельной горутине: `go rpc.ServeConn(conn)`.
- Реализует Graceful Shutdown при перехвате сигналов `os.Interrupt` и `syscall.SIGTERM`.

#### 3. Клиентское приложение ([`cmd/client/main.go`](../lessons/14-RPC/homework-14/cmd/client/main.go)):
- Подключается к RPC-серверу через `rpc.Dial("tcp", "localhost:8080")`.
- Создает пакет из нескольких тестовых сообщений и отправляет их методом `client.Call("Service.Send", newBatch, &ackCount)`.
- Запрашивает всю накопленную историю сообщений с сервера: `client.Call("Service.Messages", struct{}{}, &allMessages)`.
- Форматирует и выводит таблицу сообщений в терминал с таймстемпами.

#### 4. Изолированное тестирование в памяти ([`pkg/rpcservice/service_test.go`](../lessons/14-RPC/homework-14/pkg/rpcservice/service_test.go)):
- **Сетевое тестирование без реальных сокетов (`net.Pipe()`):**
  - Используется виртуальная двунаправленная сетевая трубка `net.Pipe()`, позволяющая тестировать полный стек `net/rpc` (бинарная GOB-сериализация, передача байт, десериализация и вызов обработчика) мгновенно и без конфликтов сетевых портов.
  - Протестирован как синхронный `client.Call`, так и асинхронный вызов через канал `client.Go(...).Done`.
- **Нагрузочный стресс-тест:** 20 одновременных горутин-клиентов посылают и читают сообщения параллельно.
- **Покрытие тестами:** **100.0%** statements, **0 data races** под `-race`.

#### 5. Лог выполнения тестов:

```bash
$ cd go_course_4 && go test -v -race -cover ./lessons/14-RPC/homework-14/...
=== RUN   TestService_DirectMethods
--- PASS: TestService_DirectMethods (0.00s)
=== RUN   TestService_RPCOverNetPipe
--- PASS: TestService_RPCOverNetPipe (0.00s)
=== RUN   TestService_Concurrency
--- PASS: TestService_Concurrency (0.01s)
PASS
coverage: 100.0% of statements
ok  	go-core-4/homework-14/pkg/rpcservice	1.495s	coverage: 100.0% of statements
```

---

## Домашнее задание 15: Реляционные базы данных (Схема БД онлайн-кинотеатра)

> 📖 **Теория:** [Том 4: Раздел 13 (Работа с базами данных: SQL, PostgreSQL, pgx)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/04-go-core-backend-web.md#13-работа-с-базами-данных-sql-и-реляционные-субд)  
> 📚 **Руководства:** [go-database-interview-guide.md](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/go-database-interview-guide.md) | [sql-livecoding-guide.md](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/sql-livecoding-guide.md)  
> 💻 **Код лекции:** [15-sql](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/15-sql)  
> 📁 **Каталог решения:** [homework-15](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/15-sql/homework-15)  
> 🧪 **Тесты:** `go test -v ./lessons/15-sql/homework-15/...`  
> 🗺️ **Путеводитель:** [course-codebase-guide.md (Урок 15)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/course-codebase-guide.md#урок-15-проектирование-реляционных-бд-и-sql)

Решение полностью реализовано в отдельной изолированной директории [`go_course_4/15-sql/homework-15`](../lessons/15-sql/homework-15) без изменения или затирания предыдущих решений.

#### 1. Схема базы данных ([`schema.sql`](../lessons/15-sql/homework-15/schema.sql)):
- **Нормализованная модель данных:**
  - `studios` — киностудии (`id SERIAL PRIMARY KEY`, `name TEXT NOT NULL UNIQUE`).
  - `actors` — актёры (`id SERIAL PRIMARY KEY`, `first_name`, `last_name`, `birth_date DATE`).
  - `directors` — режиссёры (аналогичная структура).
  - `movies` — фильмы (`id BIGSERIAL`, `title`, `release_year`, `box_office BIGINT`, `rating` с ограничением `CHECK (rating IN ('PG-10', 'PG-13', 'PG-18'))`, `studio_id` с `REFERENCES studios(id) ON DELETE RESTRICT`).
  - `movies_actors` и `movies_directors` — связи N:M с составными первичными ключами.
- **Ограничения целостности:**
  - `UNIQUE (title, release_year)` — запрет дублей по названию+год.
  - `CHECK (release_year >= 1800)` и `CHECK (box_office >= 0)`.
  - PL/pgSQL триггер `check_movie_year()` — проверка, что год не превышает текущий +10.
- **Индексы:** B-Tree на внешние ключи, даты и `lower(title)`.

#### 2. Демонстрационные данные ([`data.sql`](../lessons/15-sql/homework-15/data.sql)):
- 5 киностудий, 6 режиссёров, 10 актёров, 12 культовых фильмов.
- Полные ассоциативные связи `movies_actors` и `movies_directors`.

#### 3. Аналитические SQL-запросы ([`queries.sql`](../lessons/15-sql/homework-15/queries.sql)):
- 10 типовых аналитических задач: `JOIN`, `GROUP BY`, `HAVING`, `UNION ALL`, `COUNT(DISTINCT ...)`, подзапросы `WHERE id IN (SELECT ...)`, поиск дубликатов по `title`.

#### 4. Лог выполнения тестов:

```bash
$ cd go_course_4 && go test -v ./lessons/15-sql/homework-15/...
--- PASS: TestSchemaSQL_ValidSyntax
--- PASS: TestSchemaSQL_AllConstraints
PASS
ok  	go-core-4/homework-15
```

---

## Домашнее задание 16: Приложения с базами данных (Пакет для БД фильмов с паттерном Repository)

> 📖 **Теория:** [Том 4: Раздел 13 (Работа с базами данных: SQL, транзакции, Repository)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/04-go-core-backend-web.md#13-работа-с-базами-данных-sql-и-реляционные-субд)  
> 📚 **Руководство:** [go-database-interview-guide.md](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/go-database-interview-guide.md)  
> 💻 **Код лекции:** [16-db-apps](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/16-db-apps)  
> 📁 **Каталог решения:** [homework-16](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/16-db-apps/homework-16)  
> 🧪 **Тесты:** `go test -v -race ./lessons/16-db-apps/homework-16/...`  
> 🗺️ **Путеводитель:** [course-codebase-guide.md (Урок 16)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/course-codebase-guide.md#урок-16-работа-с-бд-в-go-интерфейсы-и-драйверы)

Решение полностью реализовано и протестировано в отдельной изолированной директории [`go_course_4/16-db-apps/homework-16`](../lessons/16-db-apps/homework-16) на основе схемы БД информационного сайта о фильмах из Занятия №15.

#### 1. Интерфейс взаимодействия с БД — Задача №1 ([`pkg/db/db.go`](../lessons/16-db-apps/homework-16/pkg/db/db.go)):
- **Модель данных `Movie`** с JSON-тегами (поля: `id`, `title`, `release_year`, `box_office`, `rating`, `studio_id`, `studio_name`).
- **Интерфейс `db.Interface`** по контракту Задачи №1:
  - `AddMovies(ctx, movies)` — добавление массива фильмов;
  - `DeleteMovie(ctx, id)` — удаление фильма по ID;
  - `UpdateMovie(ctx, m)` — обновление данных фильма;
  - `Movies(ctx, studioID)` — получение массива фильмов (если `studioID == 0` — возвращаются все фильмы, если `studioID > 0` — фильтрация по ID киностудии);
  - `MovieByID(ctx, id)` — получение фильма по ID.
- **Sentinel Errors** ([`pkg/db/errors.go`](../lessons/16-db-apps/homework-16/pkg/db/errors.go)): `db.ErrNotFound`.

#### 2. Пакет для PostgreSQL — Задача №2 ([`pkg/db/pgsql/pgsql.go`](../lessons/16-db-apps/homework-16/pkg/db/pgsql/pgsql.go)):
- Пул соединений через `github.com/jackc/pgx/v5/pgxpool`.
- **Пакетная вставка (`AddMovies`)** в рамках одной транзакции через `pgx.Batch` (слайд 21 лекции №16).
- **Фильтрация по студии (`Movies`)** через `LEFT JOIN studios` с условием `WHERE ($1 = 0 OR m.studio_id = $1)`.
- **Маппинг ошибок**: `pgx.ErrNoRows` → `db.ErrNotFound`, проверка `RowsAffected() == 0` для `UPDATE`/`DELETE`.

#### 3. Тестирование функций — Задача №3:
- **`pkg/db/memsql/memsql_test.go`**:
  - `TestAddMovies_And_Movies_All`: добавление массива фильмов и получение всех при `studioID == 0`.
  - `TestMovies_FilterByStudio`: фильтрация по ID студии (`studioID > 0`) и по несуществующей студии.
  - `TestUpdateMovie`: обновление и возврат `db.ErrNotFound`.
  - `TestDeleteMovie`: удаление, повторное удаление и проверка отсутствия.
  - `TestConcurrency`: конкурентная запись и чтение под `-race`.
- **`pkg/db/pgsql/pgsql_test.go`**: статическая проверка соответствия интерфейсу и обработка ошибок.
- **`pkg/api/api_test.go`**: тесты HTTP REST API для фильмов (`GET /api/v1/movies?studio_id=X`, `POST`, `PUT`, `DELETE`).

#### 4. REST API и веб-сервер:
- Маршрутизация на базе `gorilla/mux` ([`pkg/api/api.go`](../lessons/16-db-apps/homework-16/pkg/api/api.go)).
- Демонстрационный сервер с таймаутами и Graceful Shutdown ([`cmd/server/main.go`](../lessons/16-db-apps/homework-16/cmd/server/main.go)).

#### 5. Лог выполнения тестов:

```bash
$ cd go_course_4 && go test -count=1 -v -race -cover ./lessons/16-db-apps/homework-16/...
=== RUN   TestHandleMovies_Empty
--- PASS: TestHandleMovies_Empty (0.00s)
=== RUN   TestHandleMovies_WithFilterByStudio
--- PASS: TestHandleMovies_WithFilterByStudio (0.00s)
=== RUN   TestHandleMovieByID
--- PASS: TestHandleMovieByID (0.00s)
=== RUN   TestHandleCreateMovies_BatchAndSingle
--- PASS: TestHandleCreateMovies_BatchAndSingle (0.00s)
=== RUN   TestHandleUpdateMovie
--- PASS: TestHandleUpdateMovie (0.00s)
=== RUN   TestHandleDeleteMovie
--- PASS: TestHandleDeleteMovie (0.00s)
PASS
ok  	go-core-4/homework-16/pkg/api	1.241s
=== RUN   TestAddMovies_And_Movies_All
--- PASS: TestAddMovies_And_Movies_All (0.00s)
=== RUN   TestMovies_FilterByStudio
--- PASS: TestMovies_FilterByStudio (0.00s)
=== RUN   TestUpdateMovie
--- PASS: TestUpdateMovie (0.00s)
=== RUN   TestDeleteMovie
--- PASS: TestDeleteMovie (0.00s)
=== RUN   TestConcurrency
--- PASS: TestConcurrency (0.00s)
PASS
ok  	go-core-4/homework-16/pkg/db/memsql	1.448s
=== RUN   TestInterfaceCompliance
--- PASS: TestInterfaceCompliance (0.00s)
=== RUN   TestNew_ConnectionError
--- PASS: TestNew_ConnectionError (0.00s)
=== RUN   TestAddMovies_Empty
--- PASS: TestAddMovies_Empty (0.00s)
PASS
ok  	go-core-4/homework-16/pkg/db/pgsql	1.497s
```

---

## Домашнее задание 17: Архитектура Go-приложения (Clean Architecture и SOLID)

> 📖 **Теория:** [Том 5: Разделы 18, 19 (Архитектура сервисов, Clean Architecture, SOLID)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/05-go-core-testing-arch.md#18-архитектура-сервисов-и-паттерны-проектирования)  
> 📚 **Руководство:** [go-system-design-guide.md](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/go-system-design-guide.md)  
> 💻 **Код лекции:** [17-system-design](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/17-system-design)  
> 📁 **Каталог решения:** [homework-17](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/17-system-design/homework-17)  
> 🧪 **Тесты:** `go test -v -race ./lessons/17-system-design/homework-17/...`  
> 🗺️ **Путеводитель:** [course-codebase-guide.md (Урок 17)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/course-codebase-guide.md#урок-17-системный-дизайн-и-чистая-архитектура)

Решение полностью реализовано и протестировано в отдельной изолированной директории [`go_course_4/17-system-design/homework-17`](../lessons/17-system-design/homework-17) без изменения или затирания предыдущих решений.

#### 1. Анализ книги Роберта Мартина «Чистая архитектура» (Clean Architecture):
- **Цель архитектуры**: минимизировать человеческие трудозатраты на создание и сопровождение системы, откладывая решения по деталям реализации.
- **Правило зависимостей (The Dependency Rule)**: зависимости исходного кода направлены строго внутрь — от деталей (Web, БД, UI) к высокоуровневым сценариям (Use Cases) и сущностям (Entities).
- **Кричащая архитектура (Screaming Architecture)**: структура пакетов отражает предметную область (Cinema/Movies), а не обезличенные технические слои MVC.
- **База данных и Веб — это детали**: бизнес-правила не зависят от протокола HTTP, драйверов СУБД или форматов сериализации.

#### 2. Практическая реализация принципов SOLID в Go ([`solid/`](../lessons/17-system-design/homework-17/solid)):
- **SRP** ([`solid/srp/srp.go`](../lessons/17-system-design/homework-17/solid/srp/srp.go)): Разделение бизнес-логики управления фильмами (`MovieService`) и форматирования отчетов (`ReportFormatter`) по независимым акторам.
- **OCP** ([`solid/ocp/ocp.go`](../lessons/17-system-design/homework-17/solid/ocp/ocp.go)): Фильтрация фильмов по спецификациям (`YearFilter`, `RatingFilter`, `AndSpecification`) — добавление новых условий без модификации алгоритма фильтрации.
- **LSP** ([`solid/lsp/lsp.go`](../lessons/17-system-design/homework-17/solid/lsp/lsp.go)): Взаимозаменяемость реализаций (`EmailNotifier`, `SMSNotifier`) через единый контракт `Notifier` с сохранением семантики.
- **ISP** ([`solid/isp/isp.go`](../lessons/17-system-design/homework-17/solid/isp/isp.go)): Разделение толстого интерфейса на узкие: `Reader`, `Writer`, `Deleter`. Клиент (`ReadOnlyService`) зависит исключительно от `Reader`.
- **DIP** ([`solid/dip/dip.go`](../lessons/17-system-design/homework-17/solid/dip/dip.go)): Сервис бронирования билетов зависит от интерфейсов `PaymentGateway` и `TicketRepository`, адаптеры инжектируются через конструктор.

#### 3. Эталонная реализация Clean Architecture ([`cleanarch/`](../lessons/17-system-design/homework-17/cleanarch)):
- **Entities** ([`cleanarch/entity/movie.go`](../lessons/17-system-design/homework-17/cleanarch/entity/movie.go)): чистая сущность `Movie` с доменной валидацией, 0 внешних зависимостей.
- **Use Cases** ([`cleanarch/usecase/movie_usecase.go`](../lessons/17-system-design/homework-17/cleanarch/usecase/movie_usecase.go)): бизнес-сценарии приложения, Input Port (`MovieUseCase`), Output Port (`MovieRepository`), изоляция DTO.
- **Adapters**:
  - `repository.MemoryMovieRepository` ([`cleanarch/adapter/repository/memory_repo.go`](../lessons/17-system-design/homework-17/cleanarch/adapter/repository/memory_repo.go)): реализация хранилища.
  - `controller.MovieHTTPController` ([`cleanarch/adapter/controller/http_controller.go`](../lessons/17-system-design/homework-17/cleanarch/adapter/controller/http_controller.go)): HTTP-адаптер маршрутов.
- **Composition Root** ([`cmd/app/main.go`](../lessons/17-system-design/homework-17/cmd/app/main.go)): инициализация и связывание слоев.

#### 4. Анализ архитектуры: Монолит vs Микросервисы:
- Согласно рекомендациям занятия и Мартина Фаулера, начинать проект практически всегда следует с модульного монолита с чистой архитектурой.
- Микросервисы вводятся только при необходимости независимого масштабирования команд и имеют высокую стоимость инфраструктурной сложности (распределенные транзакции, сетевые задержки, саги).

#### 5. Лог выполнения тестов:

```bash
$ cd go_course_4 && go test -v -race ./lessons/17-system-design/homework-17/...
=== RUN   TestCleanArchitecture_FullFlow
--- PASS: TestCleanArchitecture_FullFlow (0.00s)
PASS
ok  	go-core-4/homework-17/cleanarch	3.126s
=== RUN   TestTicketBookingService_DIP
--- PASS: TestTicketBookingService_DIP (0.00s)
PASS
ok  	go-core-4/homework-17/solid/dip	2.924s
=== RUN   TestReadOnlyService_ISP
--- PASS: TestReadOnlyService_ISP (0.00s)
PASS
ok  	go-core-4/homework-17/solid/isp	2.388s
=== RUN   TestNotificationService_LSP
=== RUN   TestNotificationService_LSP/Email
=== RUN   TestNotificationService_LSP/SMS
--- PASS: TestNotificationService_LSP (0.00s)
    --- PASS: TestNotificationService_LSP/Email (0.00s)
    --- PASS: TestNotificationService_LSP/SMS (0.00s)
PASS
ok  	go-core-4/homework-17/solid/lsp	2.150s
=== RUN   TestOCP_Filters
--- PASS: TestOCP_Filters (0.00s)
PASS
ok  	go-core-4/homework-17/solid/ocp	2.743s
=== RUN   TestMovieService_SRP
--- PASS: TestMovieService_SRP (0.00s)
PASS
ok  	go-core-4/homework-17/solid/srp	2.561s
```

---

## Домашнее задание 18: Разработка микросервисов (URL Shortener Service)

> 📖 **Теория:** [Том 5: Раздел 23 (Cloud-Ready и микросервисная архитектура, Docker, Prometheus)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/05-go-core-testing-arch.md#23-cloud-ready-и-микросервисная-архитектура)  
> 📚 **Руководство:** [go-system-design-guide.md](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/go-system-design-guide.md)  
> 💻 **Код лекции:** [18-microservices](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/18-microservices)  
> 📁 **Каталог решения:** [homework-18](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/18-microservices/homework-18)  
> 🧪 **Тесты:** `go test -v -race ./lessons/18-microservices/homework-18/...`  
> 🗺️ **Путеводитель:** [course-codebase-guide.md (Урок 18)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/course-codebase-guide.md#урок-18-микросервисы-docker-и-контейнеризация)

Решение полностью реализовано и протестировано в отдельной изолированной директории [`go_course_4/18-microservices/homework-18`](../lessons/18-microservices/homework-18) со всеми требованиями к микросервисной архитектуре и контейнеризации.

#### 1. HTTP REST API и хранилище в памяти — Задача №1:
- **Хранилище на базе слайса** ([`pkg/storage/slicestore.go`](../lessons/18-microservices/homework-18/pkg/storage/slicestore.go)):
  - Строгое выполнение условия ТЗ: хранилище реализовано в виде среза в памяти (`[]Link`).
  - Потокобезопасность (`sync.RWMutex`), валидация URL и исключение дубликатов.
  - Генерация уникальных 8-символьных идентификаторов ссылок.
- **REST API** ([`pkg/api/api.go`](../lessons/18-microservices/homework-18/pkg/api/api.go)):
  - `POST /api/v1/links`: добавление ссылки в хранилище и генерация короткой.
  - `GET /api/v1/links/{id}`: получение оригинальной ссылки по идентификатору.
  - `POST /api/v1/links/resolve`: получение оригинальной ссылки по переданному короткому URL.
  - `GET /r/{id}`: перенаправление на оригинальный ресурс (HTTP 302 Redirect).
  - `GET /api/v1/links`: просмотр всех ссылок.
- **Микросервисные функции**:
  - Метрики Prometheus (`GET /metrics`): счетчик запросов `http_requests_total` и количество ссылок в памяти `shortener_links_stored_total`.
  - Healthcheck probes (`GET /healthz`, `GET /ready`).
  - Graceful Shutdown по сигналам `SIGINT`/`SIGTERM` с таймаутом контекста.

#### 2. Контейнеризация — Задача №2:
- **Multi-stage `Dockerfile`** ([`Dockerfile`](../lessons/18-microservices/homework-18/Dockerfile)):
  - `builder` на базе `golang:1.22-alpine`: сборка статического бинарника (`CGO_ENABLED=0`, `-ldflags="-s -w"`).
  - `runtime` на базе `alpine:3.19` (размер образа < 25 МБ).
  - Запуск от имени непривилегированного пользователя `appuser` (UID 1000).
- **Декларативная оркестрация**: [`docker-compose.yml`](../lessons/18-microservices/homework-18/docker-compose.yml) со встроенным Healthcheck.

#### 3. Лог выполнения тестов:

```bash
$ cd go_course_4 && go test -v -race -cover ./lessons/18-microservices/homework-18/...
=== RUN   TestAPI_CreateAndGetLink
--- PASS: TestAPI_CreateAndGetLink (0.00s)
=== RUN   TestAPI_ErrorCases
--- PASS: TestAPI_ErrorCases (0.00s)
=== RUN   TestAPI_HealthzAndMetrics
--- PASS: TestAPI_HealthzAndMetrics (0.00s)
PASS
coverage: 89.1% of statements
ok  	go-core-4/homework-18/pkg/api	1.555s
=== RUN   TestSliceStore_SaveAndGet
--- PASS: TestSliceStore_SaveAndGet (0.00s)
=== RUN   TestSliceStore_ValidationAndNotFound
--- PASS: TestSliceStore_ValidationAndNotFound (0.00s)
=== RUN   TestSliceStore_Concurrency
--- PASS: TestSliceStore_Concurrency (0.00s)
PASS
coverage: 94.9% of statements
ok  	go-core-4/homework-18/pkg/storage	1.797s
```

---

## Домашнее задание 19: Очереди сообщений и Асинхронная аналитика (Kafka / Event Sourcing)

> 📖 **Теория:** [Том 5: Раздел 21 (Очереди сообщений: Kafka, RabbitMQ, EDA)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/05-go-core-testing-arch.md#21-очереди-сообщений-и-брокеры-kafka-rabbitmq)  
> 📚 **Руководство:** [go-system-design-guide.md](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/go-system-design-guide.md)  
> 💻 **Код лекции:** [19-queue](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/19-queue)  
> 📁 **Каталог решения:** [homework-19](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/19-queue/homework-19)  
> 🧪 **Тесты:** `go test -v -race ./lessons/19-queue/homework-19/...`  
> 🗺️ **Путеводитель:** [course-codebase-guide.md (Урок 19)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/course-codebase-guide.md#урок-19-очереди-сообщений-и-брокеры-kafka)

Решение полностью реализовано и протестировано в отдельной изолированной директории [`go_course_4/19-queue/homework-19`](../lessons/19-queue/homework-19) без затирания решений предыдущих уроков.

#### 1. Архитектура распределенной событийно-ориентированной системы:
- **Event-Driven Architecture (Событийная модель)**:
  - Микросервисы общаются асинхронно через брокер сообщений (**Apache Kafka**).
  - Событие `LinkCreatedEvent` ([`pkg/event/event.go`](../lessons/19-queue/homework-19/pkg/event/event.go)) содержит `event_id`, `link_id`, `original_url`, `short_url`, `created_at`.
- **Гарантии доставки и фиксация смещений (слайд 14 лекции)**:
  - Реализован паттерн At-Least-Once Delivery: **Вычитка сообщения (`FetchMessage`) ➔ Обработка в сервисе аналитики (`RecordLink`) ➔ Ручная фиксация смещения (`CommitMessages`)**.
  - Защита от Poison Pill: при поврежденном JSON смещение фиксируется с выводом предупреждения в лог для недопущения вечного цикла.

#### 2. Микросервисы системы:
1. **URL Shortener Microservice (`shortener`)**:
   - In-memory хранилище на базе среза (`SliceStore`) с защитой `sync.RWMutex`.
   - При успешном вызове `POST /api/v1/links` ссылка сохраняется, и событие асинхронно публикуется через Kafka Producer в топик `links-events`.
   - Реализованы эндпоинты `GET /api/v1/links/{id}`, `POST /api/v1/links/resolve`, `GET /r/{id}` (HTTP 302 Redirect), `GET /healthz`, `GET /metrics`.
2. **Analytics Microservice (`analytics`)**:
   - Фоновый воркер `consumer.Worker` подписан на топик `links-events` в рамках Consumer Group `analytics-group`.
   - Бизнес-логика `service.Service` вычисляет:
     - `total_links` (общее число ссылок);
     - `total_original_length` и `avg_original_length` (средняя длина исходных ссылок);
     - `total_short_length` и `avg_short_length` (средняя длина коротких ссылок);
     - `domain_counts` (частотное распределение по доменным именам);
     - `min_original_length` и `max_original_length`.
   - Предоставляет REST API: `GET /api/v1/stats`, `GET /healthz`, `GET /ready`, `GET /metrics`.

#### 3. Тестируемость и In-Memory брокер:
- Для обеспечения мгновенного прохождения тестов без обязательного наличия запущенного кластера Kafka спроектирован [`ChannelBroker`](../lessons/19-queue/homework-19/pkg/event/channel.go) на каналах Go.
- Написан сквозной интеграционный тест [`e2e_test.go`](../lessons/19-queue/homework-19/e2e_test.go), проверяющий весь цикл: `Shortener POST -> Broker -> Analytics Consumer -> Analytics GET /api/v1/stats`.

#### 4. Контейнеризация и Docker Compose:
- **Multi-Stage Dockerfile** ([`Dockerfile`](../lessons/19-queue/homework-19/Dockerfile)): компилирует оба статических бинарника (`shortener` и `analytics`) на этапе `golang:1.22-alpine` и упаковывает в защищенный образ `alpine:3.19` с непривилегированным пользователем `appuser`.
- **Оркестрация** ([`docker-compose.yml`](../lessons/19-queue/homework-19/docker-compose.yml)): одновременный запуск `zookeeper`, `kafka`, `shortener` и `analytics` с Healthchecks.

#### 5. Лог выполнения тестов:

```bash
$ cd go_course_4 && go test -v -race ./lessons/19-queue/homework-19/...
=== RUN   TestE2E_ShortenerToAnalytics_AsyncFlow
--- PASS: TestE2E_ShortenerToAnalytics_AsyncFlow (0.02s)
PASS
ok  	go-core-4/homework-19	1.288s
=== RUN   TestAnalyticsAPI_GetStats
--- PASS: TestAnalyticsAPI_GetStats (0.00s)
=== RUN   TestAnalyticsAPI_HealthAndMetrics
--- PASS: TestAnalyticsAPI_HealthAndMetrics (0.00s)
PASS
ok  	go-core-4/homework-19/pkg/analytics/api	0.512s
=== RUN   TestWorker_Processing
--- PASS: TestWorker_Processing (0.02s)
PASS
ok  	go-core-4/homework-19/pkg/analytics/consumer	0.485s
=== RUN   TestService_RecordLink_CalculatesStats
--- PASS: TestService_RecordLink_CalculatesStats (0.00s)
=== RUN   TestService_Concurrency
--- PASS: TestService_Concurrency (0.00s)
PASS
ok  	go-core-4/homework-19/pkg/analytics/service	0.450s
=== RUN   TestChannelBroker_PublishAndSubscribe
--- PASS: TestChannelBroker_PublishAndSubscribe (0.00s)
=== RUN   TestChannelBroker_Closed
--- PASS: TestChannelBroker_Closed (0.00s)
=== RUN   TestKafkaConfig_Validation
--- PASS: TestKafkaConfig_Validation (0.00s)
PASS
ok  	go-core-4/homework-19/pkg/event	0.520s
=== RUN   TestAPI_CreateLink_PublishesEvent
--- PASS: TestAPI_CreateLink_PublishesEvent (0.00s)
=== RUN   TestAPI_GetAndResolve
--- PASS: TestAPI_GetAndResolve (0.00s)
=== RUN   TestAPI_ErrorsAndHealth
--- PASS: TestAPI_ErrorsAndHealth (0.00s)
PASS
ok  	go-core-4/homework-19/pkg/shortener/api	0.490s
=== RUN   TestSliceStore_SaveAndGet
--- PASS: TestSliceStore_SaveAndGet (0.00s)
=== RUN   TestSliceStore_Validation
--- PASS: TestSliceStore_Validation (0.00s)
=== RUN   TestSliceStore_Concurrency
--- PASS: TestSliceStore_Concurrency (0.00s)
PASS
ok  	go-core-4/homework-19/pkg/shortener/storage	0.510s
```

---

## Домашнее задание 20: Итоговый проект Lynks (NoSQL, Redis, PostgreSQL)

> 📖 **Теория:** [Том 4: Раздел 13 (Базы данных)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/04-go-core-backend-web.md#13-работа-с-базами-данных-sql-и-реляционные-субд) | [Том 5: Разделы 18, 23 (Архитектура, микросервисы, NoSQL/Redis)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/05-go-core-testing-arch.md#21-очереди-сообщений-и-брокеры-kafka-rabbitmq)  
> 📚 **Руководства:** [go-system-design-guide.md](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/go-system-design-guide.md) | [go-database-interview-guide.md](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-database-interview-guide.md)  
> 💻 **Код лекции:** [20-NoSQL](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/20-NoSQL)  
> 📁 **Каталог решения:** [homework-20](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/20-NoSQL/homework-20) | Исходная спецификация: [`lynks/final.md`](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lynks/final.md)  
> 🧪 **Тесты:** `go test -v -race ./lessons/20-NoSQL/homework-20/...`  
> 🗺️ **Путеводитель:** [course-codebase-guide.md (Урок 20)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/course-codebase-guide.md#урок-20-nosql-базы-данных-и-кэширование-redis)

Решение полностью реализовано и протестировано в отдельной изолированной директории [`go_course_4/20-NoSQL/homework-20`](../lessons/20-NoSQL/homework-20) без затирания решений предыдущих уроков в строгом соответствии со спецификацией [`final.md`](../lynks/final.md).

#### 1. Архитектура сервиса Lynks и Cache-Aside Pattern:
- **Микросервис коротких ссылок (`shortener`)**:
  - `POST /` и `POST /api/v1/links`: валидирует входящий `destination` URL, генерирует уникальный 5-значный идентификатор (`qz6d7`), сохраняет пару в PostgreSQL и асинхронно прогревает кэш в Memcache.
  - `GET /{short}`: алгоритм редиректа по паттерну **Cache-Aside**:
    1. Проверка в микросервисе `memcache` (`GET /api/v1/cache/{short}`).
    2. **Cache Hit**: мгновенный редирект (HTTP 302 Found) без обращения к реляционной СУБД.
    3. **Cache Miss**: запрос в PostgreSQL по индексированному полю `short_url`.
    4. При нахождении — асинхронная запись в `memcache` (`POST /api/v1/cache`) и редирект.
- **Микросервис кэширования (`memcache`)**:
  - `POST /api/v1/cache`: принимает `{ "shortUrl": "...", "destination": "..." }` и сохраняет в Redis со временем жизни (TTL) **строго не более 24 часов**.
  - `GET /api/v1/cache/{short}`: возвращает пару по ключу короткой ссылки.

#### 2. Базы данных:
- **PostgreSQL**:
  - Таблица `links` с первичным ключом `id` и полями `short_url`, `destination`, `created_at`.
  - Индекс `idx_links_short_url` на поле `short_url` для быстрого $O(\log N)$ поиска длинной ссылки по короткой.
- **Redis (NoSQL Key-Value)**:
  - Хранение пар с ключом `lynks:<short_url>` и TTL `24 * time.Hour`.

#### 3. Наблюдаемость (Observability):
- **Prometheus метрики**:
  - `http_requests_total` (`CounterVec` с метками `method`, `path`, `status`).
  - `http_request_duration_seconds` (`HistogramVec` с метками `method`, `path`).
  - Счетчики `cache_hits_total` и `cache_misses_total` для отслеживания эффективности кэширования.
  - Экспорт в стандартном формате по пути `GET /metrics`.
- **Структурированное логирование (Zerolog)**:
  - Пакет `github.com/rs/zerolog` выводит в `stdout` структурированные JSON-логи каждого запроса (метод, путь, статус, длительность, ошибки).

#### 4. Контейнеризация и Docker Compose:
- Отдельные Dockerfile: [`Dockerfile.shortener`](../lessons/20-NoSQL/homework-20/Dockerfile.shortener) и [`Dockerfile.memcache`](../lessons/20-NoSQL/homework-20/Dockerfile.memcache) на базе `alpine:3.19` с непривилегированным пользователем `appuser`.
- Итоговый [`docker-compose.yml`](../lessons/20-NoSQL/homework-20/docker-compose.yml), поднимающий весь комплекс: `postgres`, `redis`, `memcache` и `shortener` с автоматическими проверками работоспособности (`healthcheck`).

#### 5. Лог выполнения тестов:

```bash
$ cd go_course_4 && go test -v -race ./lessons/20-NoSQL/homework-20/...
=== RUN   TestE2E_Lynks_FullFlow
--- PASS: TestE2E_Lynks_FullFlow (0.02s)
PASS
ok  	go-core-4/homework-20	1.698s
=== RUN   TestMemcacheAPI_SetAndGet
--- PASS: TestMemcacheAPI_SetAndGet (0.00s)
PASS
ok  	go-core-4/homework-20/memcache/pkg/api	1.850s
=== RUN   TestMemoryStore_BasicOperations
--- PASS: TestMemoryStore_BasicOperations (0.00s)
=== RUN   TestMemoryStore_TTLExpiration
--- PASS: TestMemoryStore_TTLExpiration (0.07s)
=== RUN   TestMemoryStore_Concurrency
--- PASS: TestMemoryStore_Concurrency (0.00s)
PASS
ok  	go-core-4/homework-20/memcache/pkg/store	1.710s
=== RUN   TestShortenerAPI_CreateAndRedirect_CacheAside
--- PASS: TestShortenerAPI_CreateAndRedirect_CacheAside (0.00s)
PASS
ok  	go-core-4/homework-20/shortener/pkg/api	1.800s
=== RUN   TestMemoryStorage_CRUD
--- PASS: TestMemoryStorage_CRUD (0.00s)
=== RUN   TestMemoryStorage_Concurrency
--- PASS: TestMemoryStorage_Concurrency (0.00s)
PASS
ok  	go-core-4/homework-20/shortener/pkg/storage	1.527s
```


