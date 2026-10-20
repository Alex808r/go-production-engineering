# 01-intro — Введение в Go и структура первого приложения

В данной директории рассматриваются первые шаги в языке Go: инициализация приложения, точка входа `package main` и организация двухуровневой структуры пакетов (`cmd/` и `pkg/`).

---

## 📋 Структура каталога

```text
01-intro/
├── hello/
│   └── main.go                  # Минимальная программа "Hello, World!"
└── demoapp/
    ├── cmd/
    │   └── app/
    │       └── main.go          # Исполняемая точка входа приложения
    └── pkg/
        └── stringutils/
            ├── stringutils.go   # Библиотечный пакет манипуляций со строками
            └── stringutils_test.go # Модульный тест библиотеки
```

---

## 🔍 Разбираемые концепции

1. **Точка входа и пакет `main` (`hello/main.go`):**
   - Каждый исполняемый файл в Go должен принадлежать пакету `package main` и содержать функцию `func main()`.
   - Использование стандартного пакета форматированного вывода `fmt.Println()`.

2. **Архитектура приложения (`demoapp`):**
   - Разделение ответственности: код точки входа (`cmd/app/main.go`) импортирует библиотечную функциональность из собственного пакета (`pkg/stringutils`).
   - Экспортируемые идентификаторы: функции и типы с заглавной буквы (`stringutils.Concat`) доступны снаружи пакета.

3. **Базовое модульное тестирование (`pkg/stringutils/stringutils_test.go`):**
   - Написание первого unit-теста с помощью стандартного пакета `testing` и сигнатуры `func TestConcat(t *testing.T)`.

---

## 🚀 Запуск и проверка

```bash
# Запуск простейшей программы
go run ./lessons/01-intro/hello/main.go

# Запуск приложения demoapp
go run ./lessons/01-intro/demoapp/cmd/app/main.go

# Запуск тестов пакета demoapp
go test -v ./lessons/01-intro/demoapp/pkg/stringutils/...
```

---

## 🔗 Связанные материалы
- **Теория:** [Том 1: Основы языка Go (Разделы 0.1–0.7: Старт с нуля, компиляция, fmt)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/01-go-core-basics.md#0-основы-языка-go-старт-с-нуля)
- **Задачник:** [homework-tasks.md (Урок 1: Первая программа и базовый ввод-вывод)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/homework-tasks.md#урок-1-первая-программа-и-базовый-ввод-вывод-fmt)
