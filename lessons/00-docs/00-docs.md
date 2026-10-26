# 00-docs — Стандарты разработки и модульная система Go

Данная директория содержит ключевые правила, конвенции оформления кода и руководство по работе со стандартной системой модулей Go (`Go Modules`), принятые в курсе и современной индустриальной разработке.

---

## 📋 Содержимое директории

- **[Modules.MD](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/00-docs/Modules.MD)** — Руководство по системе Go Modules:
  - Инициализация модуля (`go mod init <name>`);
  - Добавление и обновление внешних зависимостей (`go get`);
  - Очистка неиспользуемых пакетов (`go mod tidy`);
  - Работа с `go.mod` и `go.sum` (контрольные суммы и воспроизводимость сборки);
  - Локальная вендоризация (`go mod vendor`).
- **[Recommendations-and-rules.MD](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/00-docs/Recommendations-and-rules.MD)** — Стандарты оформления и инженерные правила:
  - Архитектура Standard Go Project Layout: разделение на `cmd/` (исполняемые файлы) и `pkg/` (библиотечный переиспользуемый код);
  - Отказ от аварийных завершений (`os.Exit`, `log.Fatal`) внутри пакетов библиотеки — возврат ошибок через `error` вызывающему коду;
  - Именование пакетов и идентификаторов по официальному гайду Go (Effective Go, Uber Go Style Guide);
  - Требования к Pull Request: наличие тестов, форматирование `gofmt` и отсутствие ошибок линтера.

---

## 🎯 Практическое применение

Приступая к изучению языка, важно сразу следовать профессиональным стандартам структуры Go-проектов:
1. Код, предназначенный для сборки исполняемых бинарных файлов, размещается в `cmd/<app-name>/main.go`.
2. Библиотечные пакеты размещаются в `pkg/<package-name>/` и никогда не вызывают `log.Fatal` или `os.Exit`.
3. Все зависимости фиксируются в корневом файле `go.mod`.

---

## 🔗 Связанные материалы курса
- **Теория:** [Том 1: Основы языка Go (Раздел 0.8: Структура проекта и тулчейн Go)](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/go-documentation/01-go-core-basics.md#08-структура-проекта-и-cli-команды-тулчейна-go)
- **Общий путеводитель:** [course-codebase-guide.md](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/go-documentation/course-codebase-guide.md)
