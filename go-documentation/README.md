# База знаний и шпаргалка по Go (Go Core & Interview Prep)

> [!IMPORTANT]
> 📚 **Модульная библиотека Go Core (7 томов):**
> Для обеспечения максимальной производительности редакторов и 100% стабильной подсветки синтаксиса база знаний разделена на 7 специализированных томов (примерно по 1 600–4 200 строк каждый).
> 
> - 📘 **[Том 1.1: Основы языка, синтаксис и функции](01-go-core-basics.md)** (Разделы 0–3)
> - 📘 **[Том 1.2: Коллекции, структуры данных и память](01-go-core-data-structures.md)** (Разделы 4–8)
> - 📗 **[Том 2: ООП, Методы, Интерфейсы, I/O и Ошибки](02-go-core-oop-methods.md)** (Разделы 9–10)
> - 📙 **[Том 3: Конкурентность, Многопоточность и Go Runtime](03-go-core-concurrency.md)** (Разделы 11, 14, 15, 16)
> - 📕 **[Том 4: Backend, Web, Базы данных, Сокеты и Время](04-go-core-backend-web.md)** (Разделы 12, 13, 28, 29)
> - 📓 **[Том 5: Тестирование, Архитектура и Микросервисы](05-go-core-testing-arch.md)** (Разделы 17–25)
> - 📔 **[Том 6: Алгоритмы, Структуры данных и STAR-интервью](06-go-core-algorithms-interview.md)** (Разделы 26–27, 30–36)
> 
> 💡 *При составлении энциклопедии использованы материалы книги **Джона Боднера «Learning Go: An Idiomatic Approach to Real-World Go Programming» (2nd Edition, 2024)**, а канонические примеры кода и паттерны взяты из официального репозитория [github.com/learning-go-book-2e](https://github.com/learning-go-book-2e).*

> [!TIP]
> 🎯 **Практические задачники и гайды курса:**
> - 🚀 **[go-modern-features-guide.md (Гайд по Modern Go 1.21–1.24)](go-modern-features-guide.md)** — разбор `log/slog`, `iter.Seq` (Range over func), ServeMux `{id}`, `testing.F`, `slices`, `maps`, `cmp`.
> - ⚡ **[go-std-collections-guide.md (Справочник по strings, slices, maps и errors)](go-std-collections-guide.md)** — шпаргалка и cookbook по самым частым операциям, методам и ловушкам стандартных пакетов Go.
> - 💡 **[go-secrets-and-traps-guide.md (Секреты, подводные камни и Best Practices Go)](go-secrets-and-traps-guide.md)** — глубокие знания из Uber Go Style Guide и 100 Go Mistakes: struct padding, утечки памяти срезов и map, time.After leak, CFS CPU Throttling в k8s, typed nil.
> - 🏭 **[go-production-backend-guide.md (Гайд по Промышленной разработке на Go)](go-production-backend-guide.md)** — от учебного кода к боевому Enterprise-сервису: Cleanenv (12-Factor), Graceful Shutdown, `sqlc` + миграции, Prometheus + OpenTelemetry, Resilient HTTP-Client, валидация DTO, JWT, Multi-Stage Dockerfile, GitHub Actions CI/CD, Apache Kafka (Transactional Outbox & DLQ), gRPC Interceptors и AppSec (OWASP Top-10).
> - 🐳 **[docker-and-compose-starter-guide.md (Гайд по Docker и Docker Compose с нуля)](docker-and-compose-starter-guide.md)** — ментальная модель контейнеризации, шпаргалка по CLI, тома и сети, Multi-Stage Dockerfile (до 15 МБ), готовый compose-шаблон (Go + Postgres + Redis) и топ-10 ошибок новичков.
> - 🎯 **[homework-tasks.md (Практический задачник Goflex & GoSearch)](homework-tasks.md)**
> - 🗄️ **[go-database-interview-guide.md (Гайд по Базам Данных)](go-database-interview-guide.md)**
> - 💻 **[go-livecoding-guide.md (Live-Coding на Go: 47 задач + Гайд по интервью)](go-livecoding-guide.md)** — топ задач от Junior до Senior, 5-шаговый фреймворк решения, табличные тесты и топ-10 вопросов с подвохом.
> - 🗄️ **[sql-livecoding-guide.md (Live SQL Задачник: 29 задач)](sql-livecoding-guide.md)** — практические задачи с собеседований (Junior/Middle/Senior): оконные функции, JOIN, пагинация, блокировки, UPSERT, EXPLAIN ANALYZE.
> - 🏛️ **[go-system-design-guide.md (System Design для Go)](go-system-design-guide.md)**
> - 📊 **[mermaid-diagrams (Стандарты визуализации архитектуры Mermaid)](mermaid-diagrams/README.md)** — скилл и палитра для построения диаграмм C4, Sequence, State Machine и Flowchart.
> - 🗺️ **[course-codebase-guide.md (Путеводитель по кодовой базе и примерам курса)](course-codebase-guide.md)** — сквозная карта курса: 22 лекционных модуля (00–22), проекты GoSearch, Lynks и Wallet (1000+ RPS), навигатор по домашним заданиям.

---

## 🧭 С чего начать: путь изучения

> [!NOTE]
> **Для кого эта библиотека.** Тома 1–2 рассчитаны на человека, который начинает Go «с нуля» (базовое знание любого языка программирования желательно, но не обязательно). Тома 3–6 углубляются в конкурентность, runtime, веб-разработку, архитектуру и подготовку к собеседованиям.

| Этап | Что читать | Чему научитесь | Примеры кода лекций | Домашние задания и решения |
| :--- | :--- | :--- | :--- | :--- |
| **0. Подготовка** | [Том 1.1: разд. 0.1–0.8](01-go-core-basics.md#0-основы-языка-go-старт-с-нуля) | Установить Go, `go mod init`, тулчейн, ввод-вывод `fmt` | [00-docs](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/00-docs), [01-intro](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/01-intro) | [Уроки 1–3 (Goflex)](homework-tasks.md#урок-1-первая-программа-и-базовый-ввод-вывод-fmt) |
| **1. Основы языка** | [Том 1.1: разд. 1–3](01-go-core-basics.md#1-управление-потоком-control-flow),<br>[Том 1.2: разд. 4–8](01-go-core-data-structures.md#4-строки-байты-и-руны-strings-bytes--runes) | Управление потоком, функции, срезы, мапы, указатели, структуры | [02-syntax](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/02-syntax) | [Уроки 4–8 (Goflex)](homework-tasks.md#урок-4-условия-if-с-областью-видимости-и-switch), [homework-02 (GoSearch Crawler)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/02-syntax/homework-02) |
| **2. Абстракции и I/O** | [Том 2: разд. 9–10](02-go-core-oop-methods.md) | Методы, интерфейсы, полиморфизм, ошибки, `io.Reader/Writer` | [05-io](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/05-io), [06-oop](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/06-oop), [09-interfaces](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/09-interfaces) | [Урок 9](homework-tasks.md#урок-9-интерфейсы-полиморфизм-и-type-switch), [homework-05 (I/O & Cache)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/05-io/homework-05), [homework-06 (OOP Engine)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/06-oop/homework-06), [homework-09 (Ifaces)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/09-interfaces/homework-09) |
| **3. Конкурентность** | [Том 3: разд. 11, 14–16](03-go-core-concurrency.md) | Горутины, каналы, `sync`, `context`, планировщик GMP, GC | [10-concurrency](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/10-concurrency), [21-interview](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/21-interview) | [Урок 10](homework-tasks.md#урок-10-конкурентность-горутины-каналы-sync-worker-pool-context), [homework-10 (Ping-Pong)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/10-concurrency/homework-10), [Live-coding (47 задач)](go-livecoding-guide.md) |
| **4. Backend & Сеть** | [Том 4: разд. 12, 13, 28, 29](04-go-core-backend-web.md) | TCP-сокеты, HTTP, REST API, PostgreSQL, `pgxpool`, время | [11-network](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/11-network), [12-web-apps](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/12-web-apps), [13-api](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/13-api), [15-sql](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/15-sql), [16-db-apps](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/16-db-apps) | [Уроки 11–12](homework-tasks.md#урок-11-web-http-rest-api-и-middleware), [homework-11 (TCP)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/11-network/homework-11), [homework-12 (Web)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/12-web-apps/homework-12), [homework-13 (REST)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/13-api/homework-13), [homework-15 (SQL)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/15-sql/homework-15), [homework-16 (Repository)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/16-db-apps/homework-16) |
| **5. Инженерия & Архитектура** | [Том 5: разд. 17–25](05-go-core-testing-arch.md) | Тесты, бенчмарки, `pprof`, Clean Arch, gRPC, очереди, микросервисы | [07-testing](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/07-testing), [08-prof_debug](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/08-prof_debug), [14-RPC](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/14-RPC), [17-system-design](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/17-system-design), [18-microservices](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/18-microservices), [19-queue](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/19-queue), [20-NoSQL](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/20-NoSQL) | [homework-07 (TDD)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/07-testing/homework-07), [homework-08 (pprof)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/08-prof_debug/homework-08), [homework-14 (RPC)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/14-RPC/homework-14), [homework-17 (Clean Arch)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/17-system-design/homework-17), [homework-18 (Docker)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/18-microservices/homework-18), [homework-19 (Kafka)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/19-queue/homework-19), [homework-20 (Lynks)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/20-NoSQL/homework-20) |
| **6. Алгоритмы & Собеседования** | [Том 6: разд. 26–27, 30–36](06-go-core-algorithms-interview.md) | Сложность $O$, поиск, списки, деревья, heap, собеседование | [03-algorithms](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/03-algorithms), [04-datastructs](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/04-datastructs), [21-interview](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/21-interview) | [homework-03 (Index)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/03-algorithms/homework-03), [homework-04 (List)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/04-datastructs/homework-04), [Live-Coding (47 задач)](go-livecoding-guide.md), [SQL-Задачи (29 задач)](sql-livecoding-guide.md) |

**Как читать:** каждый раздел состоит из объяснения, примеров кода и «ловушек». Примеры лучше *набирать руками и запускать* (`go run .`), а не только читать. Блоки «🛡️ вопросы собеседований» и разборы внутренних структур (hchan, hmap, GMP) на первом проходе можно пропустить и вернуться к ним после практики. Краткая справка по синтаксису и командам — в [go-core-cheatsheet.md](go-core-cheatsheet.md).

> [!IMPORTANT]
> **Версии Go.** Материал ориентирован на актуальные версии языка (примеры проверялись с Go 1.25+; отдельные возможности помечены версией, например «Go 1.22+»). Если у вас старый Go, обновитесь: <https://go.dev/dl/>. Версию проверяет команда `go version`, минимальную версию модуля задаёт строка `go 1.xx` в `go.mod`.

---

## 📑 Полное интерактивное оглавление библиотеки


0. [**Основы языка Go: старт с нуля**](01-go-core-basics.md#0-основы-языка-go-старт-с-нуля)
   - [Структура программы: package main, func main() и первая программа](01-go-core-basics.md#01-структура-программы-package-main-func-main-и-первая-программа)
   - [Переменные в Go: явное var vs короткое объявление :=](01-go-core-basics.md#02-переменные-в-go-явное-var-vs-короткое-объявление)
   - [Базовые типы данных и концепция Zero Value](01-go-core-basics.md#03-базовые-типы-данных-и-концепция-zero-value)
   - [Строгая типизация и явное приведение типов (Type Conversion)](01-go-core-basics.md#04-строгая-типизация-и-явное-приведение-типов-type-conversion)
   - [Базовые арифметические и логические операторы](01-go-core-basics.md#05-базовые-арифметические-и-логические-операторы)
   - [Видимость идентификаторов: Экспортируемые (Public) vs Неэкспортируемые (Private)](01-go-core-basics.md#06-видимость-идентификаторов-экспортируемые-public-vs-неэкспортируемые-private)
   - [Базовый консольный ввод и вывод (fmt.Println, fmt.Printf, fmt.Scan)](01-go-core-basics.md#07-базовый-консольный-ввод-и-вывод-fmtprintln-fmtprintf-fmtscan)
   - [Структура проекта и CLI-команды тулчейна Go](01-go-core-basics.md#08-структура-проекта-и-cli-команды-тулчейна-go)
1. [**Управление потоком (Control Flow)**](01-go-core-basics.md#1-управление-потоком-control-flow)
   - [Формы цикла for в Go](01-go-core-basics.md#11-формы-цикла-for-в-go)
   - [Механика цикла for range со всеми типами данных](01-go-core-basics.md#механика-цикла-for-range-со-всеми-типами-данных)
   - [Конструкция if с коротким объявлением (If with Short Statement)](01-go-core-basics.md#12-конструкция-if-с-коротким-объявлением-if-with-short-statement)
   - [Особенности конструкции switch](01-go-core-basics.md#13-особенности-конструкции-switch)
   - [Прерывание вложенных циклов и switch/select: Именованные метки (Labels)](01-go-core-basics.md#14-прерывание-вложенных-циклов-и-switchselect-именованные-метки-labels)
   - [Оператор goto: строгие правила компилятора и легитимное применение](01-go-core-basics.md#141-оператор-goto-строгие-правила-компилятора-и-легитимное-применение)
   - [Механика конструкции defer](01-go-core-basics.md#15-механика-конструкции-defer)
   - [Особенности работы функции init()](01-go-core-basics.md#16-особенности-работы-функции-init)
2. [**Функции и модель передачи аргументов (Functions & Pass by Value)**](01-go-core-basics.md#2-функции-и-модель-передачи-аргументов-functions--pass-by-value)
   - [Объявление функций, параметры и множественный возврат](01-go-core-basics.md#21-объявление-функций-параметры-и-множественный-возврат)
   - [Именованные возвращаемые значения, голый возврат (Naked Return) и defer](01-go-core-basics.md#22-именованные-возвращаемые-значения-голый-возврат-naked-return-и-defer)
   - [Что фактически копируется при передаче в функцию (Модель памяти)](01-go-core-basics.md#23-что-фактически-копируется-при-передаче-в-функцию-модель-памяти)
   - [Сравнение с языками со ссылочной моделью (Ruby, Python, JavaScript)](01-go-core-basics.md#24-сравнение-с-языками-со-ссылочной-моделью-ruby-python-javascript-vs-системная-модель-go)
   - [Вариативные функции (Variadic Functions: ...Type)](01-go-core-basics.md#25-вариативные-функции-variadic-functions-type)
   - [Анонимные функции, Замыкания (Closures) и Callback](01-go-core-basics.md#26-анонимные-функции-замыкания-closures-и-callback)
   - [Переменная цикла for: изменение области видимости в Go 1.22+ и замыкания](01-go-core-basics.md#27-переменная-цикла-for-изменение-области-видимости-в-go-122-и-замыкания)
3. [**Числа, Константы и Enum-паттерн**](01-go-core-basics.md#3-числа-константы-и-enum-паттерн)
   - [Нетипизированные числовые константы (Numeric Constants)](01-go-core-basics.md#31-нетипизированные-числовые-константы-numeric-constants)
   - [Генератор констант iota и Enum-паттерн](01-go-core-basics.md#32-генератор-констант-iota-и-enum-паттерн)
   - [Форматирование строк: Спецификаторы fmt.Printf / fmt.Sprintf](01-go-core-basics.md#33-форматирование-строк-спецификаторы-fmtprintf--fmtsprintf)
   - [Битовые операции и пакет math/bits (интервью по алгоритмам)](01-go-core-basics.md#34-битовые-операции-и-пакет-mathbits-интервью-по-алгоритмам)
4. [**Строки, Байты и Руны (Strings, Bytes & Runes)**](01-go-core-data-structures.md#4-строки-байты-и-руны-strings-bytes--runes)
   - [Устройство строк и срезы UTF-8 (rune vs Байты)](01-go-core-data-structures.md#41-устройство-строк-и-срезы-utf-8-rune-vs-байты)
   - [Конкатенация строк в цикле: strings.Builder](01-go-core-data-structures.md#42-конкатенация-строк-в-цикле-stringsbuilder)
   - [Полезные функции пакета strings](01-go-core-data-structures.md#43-полезные-функции-пакета-strings)
5. [**Массивы и Слайсы (Arrays & Slices)**](01-go-core-data-structures.md#5-массивы-и-слайсы-arrays--slices)
   - [Массивы в Go: последовательная область памяти, фиксированная длина и тип [N]T](01-go-core-data-structures.md#51-массивы-в-go-последовательная-область-памяти-фиксированная-длина-и-тип-nt)
   - [Внутреннее устройство слайсов (Slice Header)](01-go-core-data-structures.md#52-внутреннее-устройство-слайсов-slice-header)
   - [Механика работы append, базового массива и алгоритм роста cap (Go 1.18+)](01-go-core-data-structures.md#53-механика-работы-append-базового-массива-и-алгоритм-роста-cap-go-118)
   - [Передача слайса в функцию: почему изменение элементов видно снаружи, а append — нет?](01-go-core-data-structures.md#54-передача-слайса-в-функцию-почему-изменение-элементов-видно-снаружи-а-append--нет)
   - [Полное выражение среза (Full Slice Expression: a[low:high:max])](01-go-core-data-structures.md#55-полное-выражение-среза-full-slice-expression-alowhighmax)
   - [Независимое копирование: Функция copy(dest, src) и slices.Clone (Go 1.21+)](01-go-core-data-structures.md#56-независимое-копирование-функция-copydest-src-и-slicesclone-go-121)
   - [Сравнение слайсов: Оператор == и пакет slices.Equal](01-go-core-data-structures.md#57-сравнение-слайсов-оператор--и-пакет-slicesequal)
   - [nil-слайс vs Пустой слайс (T{}): В чем разница?](01-go-core-data-structures.md#58-nil-слайс-vs-пустой-слайс-t-в-чем-разница)
   - [Слияние (конкатенация) слайсов через append](01-go-core-data-structures.md#59-слияние-конкатенация-слайсов-через-append)
   - [Удаление элементов из слайса (Slice Element Deletion)](01-go-core-data-structures.md#510-удаление-элементов-из-слайса-slice-element-deletion)
   - [Многомерные слайсы (Slices of Slices) и динамическая аллокация матриц](01-go-core-data-structures.md#511-многомерные-слайсы-slices-of-slices-и-динамическая-аллокация-матриц)
   - [Встроенная функция clear() для слайсов (Go 1.21+): clear(s) vs s[:0] vs s = nil](01-go-core-data-structures.md#512-встроенная-функция-clear-для-слайсов-go-121-clears-vs-s0-vs-s--nil)
   - [Современная функциональная обработка слайсов: библиотека samber/lo (Go Lodash)](01-go-core-data-structures.md#513-современная-функциональная-обработка-слайсов-библиотека-samberlo-go-lodash)
   - [Преобразование срезов в массивы и указатели на массив (Go 1.17+ и Go 1.20+)](01-go-core-data-structures.md#514-преобразование-срезов-в-массивы-и-указатели-на-массив-go-117-и-go-120)
6. [**Карты (Maps)**](01-go-core-data-structures.md#6-карты-maps)
   - [Способы создания и инициализации](01-go-core-data-structures.md#61-способы-создания-и-инициализации)
   - [Поведение nil-мапы, чтение, запись, удаление и очистка clear(m)](01-go-core-data-structures.md#62-поведение-nil-мапы-чтение-запись-и-удаление)
   - [Что может (и не может) быть ключом в Map? (Comparable Types)](01-go-core-data-structures.md#63-что-может-и-не-может-быть-ключом-в-map-comparable-types)
   - [Неупорядоченность итерации и сортировка ключей](01-go-core-data-structures.md#64-неупорядоченность-итерации-и-сортировка-ключей)
   - [Конкурентный доступ: Фатальная ошибка (fatal error: concurrent map writes)](01-go-core-data-structures.md#65-конкурентный-доступ-фатальная-ошибка-fatal-error-concurrent-map-writes)
   - [Почему в Go запрещено брать адрес элемента мапы (&m[k])?](01-go-core-data-structures.md#66-почему-в-go-запрещено-брать-адрес-элемента-мапы-mk)
   - [Внутреннее устройство Map: Старая реализация (hmap) vs Новая реализация (Swiss Tables в Go 1.24+)](01-go-core-data-structures.md#67-внутреннее-устройство-map-старая-реализация-hmap-vs-новая-реализация-swiss-tables-в-go-124)
   - [Практические паттерны использования карт (Map Idioms & Best Practices)](01-go-core-data-structures.md#68-практические-паттерны-использования-карт-map-idioms--best-practices)
7. [**Указатели в деталях (Pointers)**](01-go-core-data-structures.md#7-указатели-в-деталях-pointers)
   - [Указатели в Go: операторы адресации & и разыменования \*](01-go-core-data-structures.md#71-указатели-в-go-операторы-адресации-и-разыменования)
   - [Защита от паники: Проверка указателя на nil (Nil Pointer Check)](01-go-core-data-structures.md#72-защита-от-паники-проверка-указателя-на-nil-nil-pointer-check)
   - [Каверзные вопросы собеседований: Арифметика указателей и Escape-анализ](01-go-core-data-structures.md#73-каверзные-вопросы-собеседований-арифметика-указателей-и-escape-анализ)
   - [Функция new(T) vs make() vs &T{} (Топ-вопрос собеседований)](01-go-core-data-structures.md#74-функция-newt-vs-make-vs-t-топ-вопрос-собеседований)
   - [Адресуемость значений (Addressability) и понятия lvalue / rvalue в Go](01-go-core-data-structures.md#75-адресуемость-значений-addressability-и-понятия-lvalue--rvalue-в-go)
8. [**Пользовательские типы и Структуры (Structs)**](01-go-core-data-structures.md#8-пользовательские-типы-и-структуры-structs)
   - [Способы объявления и инициализации структур](01-go-core-data-structures.md#81-способы-объявления-и-инициализации-структур)
   - [Композиция вместо Наследования (Embedding)](01-go-core-data-structures.md#82-композиция-вместо-наследования-embedding)
   - [Паттерн Конструктора (New()) и отсутствие Деструкторов в Go](01-go-core-data-structures.md#паттерн-конструктора-new-и-отсутствие-деструкторов-в-go)
   - [Указатели и Структуры (Struct Pointers): авто-разыменование и Zero Copy](01-go-core-data-structures.md#83-указатели-и-структуры-struct-pointers-авто-разыменование-и-zero-copy)
   - [Анонимные структуры и Табличные тесты (Table-Driven Tests)](01-go-core-data-structures.md#84-анонимные-структуры-и-табличные-тесты-table-driven-tests)
   - [Пустая структура struct{} (0 байт) и паттерн Множество (Set)](01-go-core-data-structures.md#85-пустая-структура-struct-0-байт-и-паттерн-множество-set)
   - [Теги структур (Struct Tags)](01-go-core-data-structures.md#86-теги-структур-struct-tags)
   - [Паттерн Functional Options (Топ-вопрос Middle/Senior собеседований)](01-go-core-data-structures.md#87-паттерн-functional-options-топ-вопрос-middlesenior-собеседований)
   - [Продвинутая работа с JSON (encoding/json)](01-go-core-data-structures.md#88-продвинутая-работа-с-json-encodingjson)
   - [Работа с XML (encoding/xml): теги, атрибуты и сравнение с JSON](01-go-core-data-structures.md#89-работа-с-xml-encodingxml-теги-атрибуты-и-сравнение-с-json)
   - [Выравнивание данных в памяти и Memory Padding (Struct Alignment)](01-go-core-data-structures.md#810-выравнивание-данных-в-памяти-и-memory-padding-struct-alignment)
9. [**Методы и Интерфейсы (Methods & Interfaces)**](02-go-core-oop-methods.md#9-методы-и-интерфейсы-methods--interfaces)
   - [Парадигма ООП в Go: Сравнение ООП (Stateful), Процедурного (Stateless) и Функционального стилей](02-go-core-oop-methods.md#90-парадигма-ооп-в-go-методология-stateful-vs-stateless-и-dependency-injection)
   - [Указатели в ресиверах методов (Pointer vs Value Receivers)](02-go-core-oop-methods.md#91-указатели-в-ресиверах-методов-pointer-vs-value-receivers)
   - [Method Value vs Method Expression и безопасные методы с nil-ресиверами](02-go-core-oop-methods.md#-методы-с-nil-ресиверами-nil-receiver-methods)
   - [Что такое интерфейс: Контракт требований, Duck Typing и сравнение с Ruby](02-go-core-oop-methods.md#92-что-такое-интерфейс-контракт-требований-duck-typing-и-сравнение-с-ruby)
   - [Почему интерфейсы ДОЛЖНЫ быть аргументами функций (Accept interfaces, return structs)](02-go-core-oop-methods.md#почему-интерфейсы-должны-быть-аргументами-функций)
   - [Пустой интерфейс (any / interface{}), Type Assertion и Type Switch](02-go-core-oop-methods.md#93-пустой-интерфейс-any--interface-type-assertion-и-type-switch)
   - [Внутреннее устройство интерфейсов и коварная ловушка nil != nil](02-go-core-oop-methods.md#94-внутреннее-устройство-интерфейсов-и-коварная-ловушка-nil--nil)
   - [Встроенный интерфейс fmt.Stringer (Кастомный строковый вывод)](02-go-core-oop-methods.md#95-встроенный-интерфейс-fmtstringer-кастомный-строковый-вывод)
   - [Потоковый ввод-вывод: интерфейсы io.Reader и io.Writer](02-go-core-oop-methods.md#96-потоковый-ввод-вывод-интерфейсы-ioreader-и-iowriter)
   - [Архитектурный кейс GoSearch: Персистентность поисковых данных (io.Writer, io.Reader)](02-go-core-oop-methods.md#архитектурный-кейс-персистентность-поисковых-данных-урок-05-io)
   - [Работа с файлами и буферизованный ввод-вывод (os, bufio.Scanner, io.Copy)](02-go-core-oop-methods.md#97-работа-с-файлами-и-буферизованный-ввод-вывод-os-bufioscanner-iocopy)
   - [Эволюция io/ioutil: миграция на пакеты io и os (Go 1.16+)](02-go-core-oop-methods.md#98-эволюция-ioioutil-миграция-на-пакеты-io-и-os-go-116)
   - [Консольный ввод (os.Stdin) и флаги командной строки (пакет flag)](02-go-core-oop-methods.md#99-консольный-ввод-osstdin-и-флаги-командной-строки-пакет-flag)
   - [Бинарная сериализация GOB (encoding/gob) для Go-to-Go обмена](02-go-core-oop-methods.md#910-бинарная-сериализация-gob-encodinggob-и-форматы-данных)
   - [Интерфейсы как ограничения типов (Constraint Interfaces): Method Sets vs Type Sets](02-go-core-oop-methods.md#911-интерфейсы-как-ограничения-типов-constraint-interfaces-method-sets-vs-type-sets)
10. [**Обработка ошибок (Error Handling)**](02-go-core-oop-methods.md#10-обработка-ошибок-error-handling)

   - [Ошибки как обычные значения и встроенный интерфейс error](02-go-core-oop-methods.md#101-ошибки-как-обычные-значения-и-встроенный-интерфейс-error)
   - [Создание и Оборачивание ошибок (Error Wrapping с %w)](02-go-core-oop-methods.md#102-создание-и-оборачивание-ошибок-error-wrapping-с-w)
   - [Проверка ошибок: errors.Is и errors.As](02-go-core-oop-methods.md#103-проверка-ошибок-errorsis-и-errorsas)
   - [Объединение нескольких ошибок: errors.Join (Go 1.20+)](02-go-core-oop-methods.md#104-объединение-нескольких-ошибок-errorsjoin-go-120)
   - [Паника и Восстановление: panic, recover, defer](02-go-core-oop-methods.md#105-паника-и-восстановление-panic-recover-defer)
   - [Контракт возврата при ошибке: Идиома Zero Value и частичные результаты](02-go-core-oop-methods.md#106-контракт-возврата-при-ошибке-идиома-zero-value-и-частичные-результаты)

11. [**Конкурентность: Горутины и Каналы (Concurrency)**](03-go-core-concurrency.md#11-конкурентность-горутины-и-каналы-concurrency)

   - [Конкурентность vs Параллелизм (Concurrency is not Parallelism)](03-go-core-concurrency.md#110-конкурентность-vs-параллелизм-concurrency-is-not-parallelism)
   - [Горутины (Goroutines) vs Системные потоки (OS Threads)](03-go-core-concurrency.md#111-горутины-goroutines-vs-системные-потоки-os-threads)
   - [Каналы (Channels): Буферизованные и Небуферизованные](03-go-core-concurrency.md#112-каналы-channels-буферизованные-и-небуферизованные)
   - [Конструкция select и Таймауты](03-go-core-concurrency.md#113-конструкция-select-и-таймауты)
   - [Синхронизация: sync.WaitGroup](03-go-core-concurrency.md#114-синхронизация-syncwaitgroup)
   - [Паттерн Worker Pool (Пул воркеров)](03-go-core-concurrency.md#115-паттерн-worker-pool-пул-воркеров)
   - [Пакет context: Отмена операций, таймауты и метаданные (context.Context)](03-go-core-concurrency.md#116-пакет-context-отмена-операций-таймауты-и-метаданные-contextcontext)
   - [Анатомия каналов: структура hchan, очереди sudog, прямой стек-в-стек copy и матрица состояний](03-go-core-concurrency.md#117-анатомия-каналов-структура-hchan-очереди-sudog-прямой-стек-в-стек-copy-и-матрица-состояний)
   - [Типовые ошибки конкурентного программирования (Common Concurrency Pitfalls)](03-go-core-concurrency.md#118-типовые-ошибки-конкурентного-программирования-common-concurrency-pitfalls)

12. [**Backend: Web-разработка и HTTP (net/http)**](04-go-core-backend-web.md#12-backend-web-разработка-и-http-nethttp)

   - [Базовый HTTP-сервер и http.ServeMux](04-go-core-backend-web.md#121-базовый-http-сервер-и-httpservemux)
   - [Экосистема роутеров и фреймворков в Go (chi, gin, echo, fiber)](04-go-core-backend-web.md#122-экосистема-роутеров-и-фреймворков-в-go-chi-gin-echo-fiber)
   - [Middleware (Промежуточное ПО)](04-go-core-backend-web.md#123-middleware-промежуточное-по)
   - [Graceful Shutdown (Элегантное завершение сервера)](04-go-core-backend-web.md#124-graceful-shutdown-элегантное-завершение-сервера)
   - [Архитектура REST API и слои приложения (Standard Go Project Layout)](04-go-core-backend-web.md#125-архитектура-rest-api-и-слои-приложения-standard-go-project-layout)
   - [Идиоматичный REST API без фреймворков (CRUD, JSON Helpers, ловушки)](04-go-core-backend-web.md#126-идиоматичный-rest-api-без-фреймворков-crud-json-helpers-ловушки)
   - [Паттерн DTO и указатели для частичного обновления (Partial Update)](04-go-core-backend-web.md#127-паттерн-dto-и-указатели-для-частичного-обновления-partial-update)
   - [Эволюция маршрутизации: ручной диспатчинг (Pre-1.22) vs Native Routing (Go 1.22+)](04-go-core-backend-web.md#128-эволюция-маршрутизации-ручной-диспатчинг-pre-122-vs-native-routing-go-122)
   - [Docker Compose для PostgreSQL в локальной разработке](04-go-core-backend-web.md#129-docker-compose-для-postgresql-в-локальной-разработке)
   - [Подводные камни http.Client и http.Transport в Highload-продакшене](04-go-core-backend-web.md#1210-подводные-камни-httpclient-и-httptransport-в-highload-продакшене)
   - [Маршрутизатор Gorilla Mux (github.com/gorilla/mux)](04-go-core-backend-web.md#1211-маршрутизатор-gorilla-mux-githubcomgorillamux)
   - [Традиционные веб-приложения и SSR: Шаблонизация (html/template)](04-go-core-backend-web.md#1212-традиционные-веб-приложения-и-ssr-шаблонизация-htmltemplate)
   - [Файловый сервер, раздача статики и Single Page Applications (SPA)](04-go-core-backend-web.md#1213-файловый-сервер-раздача-статики-и-single-page-applications-spa)
   - [Запуск защищенного веб-сервера по TLS/HTTPS (ListenAndServeTLS)](04-go-core-backend-web.md#1214-запуск-защищенного-веб-сервера-по-tlshttps-listenandservetls)
   - [Архитектура пакета API и ООП-модель сервиса](04-go-core-backend-web.md#1215-архитектура-пакета-api-и-ооп-модель-сервиса)
   - [Контекст в HTTP-обработчиках: отмена цепочки вызовов и проброс метаданных](04-go-core-backend-web.md#1216-контекст-в-http-обработчиках-отмена-цепочки-вызовов-и-проброс-метаданных)
   - [Аутентификация и авторизация в API: Сессии vs JWT](04-go-core-backend-web.md#1217-аутентификация-и-авторизация-в-api-сессии-vs-jwt)
   - [Документирование и контрактное тестирование API (OpenAPI / Swagger и Postman)](04-go-core-backend-web.md#1218-документирование-и-контрактное-тестирование-api-openapi--swagger-и-postman)
   - [Web Security: CORS, Rate Limiting, Secure Cookies и CSRF](04-go-core-backend-web.md#1219-веб-безопасность-web-security-cors-rate-limiting-и-защита-cookiescsrf)
   - [Аутентификация: JWT и OAuth2 в Go (JWT, Access + Refresh токены, OAuth2)](04-go-core-backend-web.md#1220-аутентификация-jwt-и-oauth2-в-go)

13. [**Работа с Базами Данных (PostgreSQL & pgxpool)**](04-go-core-backend-web.md#13-работа-с-базами-данных-postgresql--pgxpool)

   - [Пул соединений (pgxpool.Pool)](04-go-core-backend-web.md#131-пул-соединений-pgxpoolpool)
   - [Запросы, Параметризация и обработка ErrNoRows](04-go-core-backend-web.md#132-запросы-параметризация-и-обработка-errnorows)
   - [Транзакции (ACID Transactions в Go)](04-go-core-backend-web.md#133-транзакции-acid-transactions-в-go)
   - [Продвинутый SQL для собеседований: Индексы, EXPLAIN ANALYZE, CTE и Оконные функции](04-go-core-backend-web.md#134-продвинутый-sql-для-собеседований-индексы-explain-analyze-cte-и-оконные-функции)
   - [Миграции Баз Данных: golang-migrate vs goose и генераторы кода](04-go-core-backend-web.md#135-миграции-баз-данных-golang-migrate-vs-goose-и-генераторы-кода)

14. [**Устройство Go Runtime: GMP, Память и GC**](03-go-core-concurrency.md#14-устройство-go-runtime-gmp-память-и-gc)

   - [Модель планировщика GMP (Goroutine, Machine, Processor)](03-go-core-concurrency.md#141-модель-планировщика-gmp-goroutine-machine-processor)
   - [Сборщик мусора: Tri-color Mark & Sweep, освобождение памяти и Scavenger](03-go-core-concurrency.md#142-сборщик-мусора-garbage-collector-tri-color-mark--sweep-освобождение-памяти-и-scavenger)
   - [Давление на сборщик мусора (GC Pressure) и деградация производительности](03-go-core-concurrency.md#6-давление-на-сборщик-мусора-gc-pressure-и-деградация-производительности)
   - [Escape-анализ и устройство стека горутины](03-go-core-concurrency.md#143-escape-анализ-и-устройство-стека-горутины)
   - [Граница cgo и FFI: Архитектура, накладные расходы и правила указателей](03-go-core-concurrency.md#144-граница-cgo-и-ffi-архитектура-накладные-расходы-и-правила-указателей)

15. [**Продвинутая Синхронизация и Паттерны Конкурентности**](03-go-core-concurrency.md#15-продвинутая-синхронизация-и-паттерны-конкурентности)

   - [sync.Mutex, sync.RWMutex, sync.Once и sync.Pool](03-go-core-concurrency.md#151-syncmutex-syncrwmutex-synconce-и-syncpool)
   - [sync.Map (Потокобезопасная карта: Two-Map Architecture, vs RWMutex)](03-go-core-concurrency.md#4-syncmap-потокобезопасная-карта-без-внешнего-мьютекса)
   - [Пакет sync/atomic и Lock-Free счетчики](03-go-core-concurrency.md#152-пакет-syncatomic-и-lock-free-счетчики)
   - [Паттерны: Fan-Out / Fan-In, Pipeline и ErrGroup](03-go-core-concurrency.md#153-паттерны-fan-out--fan-in-pipeline-и-errgroup)
   - [Детектор гонок памяти (Data Race Detector: -race)](03-go-core-concurrency.md#154-детектор-гонок-памяти-data-race-detector--race)
   - [Go Memory Model и гарантии Happens-Before](03-go-core-concurrency.md#155-go-memory-model-и-гарантии-happens-before)
   - [sync.Cond и Виды блокировок (Deadlock, Livelock, Starvation)](03-go-core-concurrency.md#156-synccond-и-виды-блокировок-deadlock-livelock-starvation)
   - [Пакет context.Context в деталях: анатомия, утечки горутин и новинки Go 1.21+](03-go-core-concurrency.md#157-пакет-contextcontext-в-деталях-анатомия-утечки-горутин-и-новинки-go-121)
   - [Hardware-Aware Go: Кэш-линии процессора, False Sharing и выравнивание структур](03-go-core-concurrency.md#158-hardware-aware-go-кэш-линии-cpu-64-байта-data-locality-и-false-sharing)

16. [**Обобщенное программирование (Generics в Go 1.18+)**](03-go-core-concurrency.md#16-обобщенное-программирование-generics-в-go-118)

   - [Философия Generics: Когда, Зачем и Какие проблемы решили](03-go-core-concurrency.md#160-философия-обобщённого-программирования-когда-зачем-и-какие-проблемы-решили-generics)
   - [Type Parameters, Any и Comparable](03-go-core-concurrency.md#161-type-parameters-any-и-comparable)
   - [Type Constraints, пакет cmp и оператор тильда ~](03-go-core-concurrency.md#162-type-constraints-пакет-cmp-и-оператор-тильда)
   - [Итераторы: range over func и пакет iter (Go 1.23+)](03-go-core-concurrency.md#163-итераторы-range-over-func-и-пакет-iter-go-123)

17. [**Тестирование, Бенчмарки и Профилирование (Testing & pprof)**](05-go-core-testing-arch.md#17-тестирование-бенчмарки-и-профилирование-testing--pprof)

   - [Философия и фундаментальные цели тестирования](05-go-core-testing-arch.md#171-философия-и-фундаментальные-цели-тестирования)
   - [Пирамида тестирования и виды тестов в Go](05-go-core-testing-arch.md#172-пирамида-тестирования-и-виды-тестов-в-go)
   - [Команды запуска, флаги go test и сброс кэша](05-go-core-testing-arch.md#173-команды-запуска-флаги-go-test-и-сброс-кэша)
   - [Управление жизненным циклом тестов пакета: TestMain(m *testing.M)](05-go-core-testing-arch.md#174-управление-жизненным-циклом-тестов-пакета-testmainm-testingm)
   - [Табличные тесты (Table-Driven Tests) и подтесты t.Run](05-go-core-testing-arch.md#175-табличные-тесты-table-driven-tests-и-подтесты-trun)
   - [Имитация зависимостей (Mocks, Stubs) и рефакторинг «Выделение функции»](05-go-core-testing-arch.md#176-имитация-зависимостей-mocks-stubs-и-рефакторинг-выделение-функции)
   - [Интеграционные тесты и идемпотентность работы с СУБД](05-go-core-testing-arch.md#177-интеграционные-тесты-и-идемпотентность-работы-с-субд)
   - [Тестирование HTTP-хендлеров (net/http/httptest)](05-go-core-testing-arch.md#178-тестирование-http-хендлеров-nethttphttptest)
   - [Test-Driven Development (TDD): Red, Green, Refactor](05-go-core-testing-arch.md#179-test-driven-development-tdd-red-green-refactor)
   - [Бенчмаркинг производительности (testing.B)](05-go-core-testing-arch.md#1710-бенчмаркинг-производительности-testingb)
   - [Fuzz-тестирование (Go 1.18+)](05-go-core-testing-arch.md#1711-fuzz-тестирование-go-118)
   - [Профилирование бенчмарков (-cpuprofile, -memprofile)](05-go-core-testing-arch.md#1712-профилирование-бенчмарков--cpuprofile--memprofile)
   - [Профилирование живых приложений в runtime (net/http/pprof)](05-go-core-testing-arch.md#1713-профилирование-живых-приложений-в-runtime-nethttppprof)
   - [Отладка программного обеспечения и отладчик Delve (dlv)](05-go-core-testing-arch.md#1714-отладка-программного-обеспечения-и-отладчик-delve-dlv)
   - [Трассировка выполнения (runtime/trace & go tool trace)](05-go-core-testing-arch.md#1715-трассировка-выполнения-runtimetrace--go-tool-trace)
   - [Сводная таблица методов пакета testing](05-go-core-testing-arch.md#1716-сводная-таблица-методов-пакета-testing)
   - [Детерминированное тестирование конкурентности: testing/synctest (Go 1.25+)](05-go-core-testing-arch.md#1717-детерминированное-тестирование-конкурентности-testingsynctest-go-125)


18. [**Архитектура сервисов и Структурированное логирование**](05-go-core-testing-arch.md#18-архитектура-сервисов-и-структурированное-логирование)

   - [Clean Architecture и Standard Go Project Layout](05-go-core-testing-arch.md#181-clean-architecture-и-standard-go-project-layout)
   - [Структурированное логирование: log/slog (в Go 1.21+)](05-go-core-testing-arch.md#182-структурированное-логирование-logslog-в-go-121)
   - [Принципы декомпозиции пакетов: Screaming Architecture, Сцепление и Связность](05-go-core-testing-arch.md#183-принципы-декомпозиции-пакетов-screaming-architecture-сцепление-и-связность)
   - [Архитектурный шаблон «Ядро и плагины» (Core & Plugins / Microkernel)](05-go-core-testing-arch.md#184-архитектурный-шаблон-ядро-и-плагины-core--plugins--microkernel)
   - [Циклические зависимости (import cycle not allowed) и стратегии их устранения](05-go-core-testing-arch.md#185-циклические-зависимости-import-cycle-not-allowed-и-стратегии-их-устранения)

19. [**Управление зависимостями: Go Modules и Vendoring**](05-go-core-testing-arch.md#19-управление-зависимостями-go-modules-и-vendoring)

   - [Ключевые команды go mod](05-go-core-testing-arch.md#191-ключевые-команды-go-mod)
   - [Назначение файлов go.mod и go.sum](05-go-core-testing-arch.md#192-назначение-файлов-gomod-и-gosum)
   - [Встраивание файлов в бинарник: //go:embed (Go 1.16+)](05-go-core-testing-arch.md#193-встраивание-файлов-в-бинарник-goembed-go-116)
   - [Кодогенерация: go generate](05-go-core-testing-arch.md#194-кодогенерация-go-generate)
   - [Чтение метаданных сборки в runtime: debug.ReadBuildInfo()](05-go-core-testing-arch.md#195-чтение-метаданных-сборки-в-runtime-debugreadbuildinfo)

20. [**Удаленный вызов процедур (RPC, net/rpc) и gRPC**](05-go-core-testing-arch.md#20-удаленный-вызов-процедур-rpc-netrpc-и-grpc)

   - [Концепция RPC (Remote Procedure Call) и сравнение с REST API](05-go-core-testing-arch.md#201-концепция-rpc-remote-procedure-call-и-сравнение-с-rest-api)
   - [Стандартный пакет Go net/rpc](05-go-core-testing-arch.md#202-стандартный-пакет-go-netrpc)
   - [Тестирование RPC-сервисов в памяти (net.Pipe())](05-go-core-testing-arch.md#203-тестирование-rpc-сервисов-в-памяти-netpipe)
   - [Почему gRPC быстрее REST JSON](05-go-core-testing-arch.md#204-почему-grpc-быстрее-rest-json)
   - [Protocol Buffers (.proto) и кодогенерация](05-go-core-testing-arch.md#205-protocol-buffers-proto-и-кодогенерация)
   - [Реализация gRPC сервера и клиента на Go](05-go-core-testing-arch.md#206-реализация-grpc-сервера-и-клиента-на-go)
   - [Enterprise gRPC: Интерцепторы, Метаданные и grpc-gateway](05-go-core-testing-arch.md#207-enterprise-grpc-интерцепторы-метаданные-и-grpc-gateway)

21. [**Очереди сообщений и Кэширование (RabbitMQ / Kafka & Redis)**](05-go-core-testing-arch.md#21-очереди-сообщений-и-кэширование-rabbitmq--kafka--redis)

   - [Асинхронные очереди (Producer-Consumer Pattern)](05-go-core-testing-arch.md#211-асинхронные-очереди-producer-consumer-pattern)
   - [In-Memory Кэширование с Redis (go-redis)](05-go-core-testing-arch.md#212-in-memory-кэширование-с-redis-go-redis)
   - [Transactional Outbox Pattern (Решение проблемы Dual Write)](05-go-core-testing-arch.md#213-transactional-outbox-pattern-решение-проблемы-dual-write)
   - [Распределенные блокировки (Distributed Locks)](05-go-core-testing-arch.md#214-распределенные-блокировки-distributed-locks)
   - [Практика работы с Apache Kafka в Go (segmentio/kafka-go)](05-go-core-testing-arch.md#215-практика-работы-с-apache-kafka-в-go-segmentiokafka-go)
   - [Отказоустойчивость Apache Kafka: Rebalance, Poison Pill, DLQ и Debezium CDC](05-go-core-testing-arch.md#216-отказоустойчивость-apache-kafka-rebalance-poison-pill-dlq-и-debezium-cdc)

22. [**Наблюдаемость систем (Observability: Prometheus & OpenTelemetry)**](05-go-core-testing-arch.md#22-наблюдаемость-систем-observability-prometheus--opentelemetry)

   - [Метрики Prometheus](05-go-core-testing-arch.md#221-метрики-prometheus)
   - [Распределенная трассировка (OpenTelemetry / Jaeger)](05-go-core-testing-arch.md#222-распределенная-трассировка-opentelemetry--jaeger)

23. [**Cloud-Ready приложения и 12-Factor App в Go**](05-go-core-testing-arch.md#23-cloud-ready-приложения-и-12-factor-app-в-go)

   - [Конфигурация через переменные окружения](05-go-core-testing-arch.md#231-конфигурация-через-переменные-окружения)
   - [Health checks: Liveness & Readiness Probes для Kubernetes](05-go-core-testing-arch.md#232-health-checks-liveness--readiness-probes-для-kubernetes)
   - [Docker: Multi-Stage Build для Go-приложений](05-go-core-testing-arch.md#233-docker-multi-stage-build-для-go-приложений)

24. [**Production-экосистема Go: Топ сторонних библиотек и инструментов**](05-go-core-testing-arch.md#24-production-экосистема-go-топ-сторонних-библиотек-и-инструментов)

   - [Сводная таблица золотого стандарта Go-библиотек](05-go-core-testing-arch.md#241-сводная-таблица-золотого-стандарта-go-библиотек)
   - [Линтеры и Статический анализ (golangci-lint)](05-go-core-testing-arch.md#242-линтеры-и-статический-анализ-golangci-lint)
   - [Валидация входных данных (go-playground/validator)](05-go-core-testing-arch.md#243-валидация-входных-данных-go-playgroundvalidator)
   - [Зеркалирование и тестирование трафика: GoReplay (утилита gor)](05-go-core-testing-arch.md#244-зеркалирование-и-тестирование-трафика-goreplay-утилита-gor)
   - [Золотой стандарт конфигурации .golangci.yml](05-go-core-testing-arch.md#245-золотой-стандарт-конфигурации-golangciyml)
   - [Продакшен Makefile для автоматизации разработки](05-go-core-testing-arch.md#246-продакшен-makefile-для-автоматизации-разработки)
   - [Паттерн Feature Flags (Управление фичами на лету)](05-go-core-testing-arch.md#247-паттерн-feature-flags-управление-фичами-на-лету)

25. [**Паттерны Отказоустойчивости в Микросервисах (Resilience & Stability)**](05-go-core-testing-arch.md#25-паттерны-отказоустойчивости-в-микросервисах-resilience--stability)

   - [Circuit Breaker (Предохранитель)](05-go-core-testing-arch.md#251-circuit-breaker-предохранитель)
   - [Rate Limiting и Throttling (Ограничение частоты запросов)](05-go-core-testing-arch.md#252-rate-limiting-и-throttling-ограничение-частоты-запросов)
   - [Retry + Exponential Backoff + Jitter](05-go-core-testing-arch.md#253-retry--exponential-backoff--jitter)

26. [**Базовые Алгоритмы и O-нотация в Go**](06-go-core-algorithms-interview.md#26-базовые-алгоритмы-и-o-нотация-в-go)

   - [Оценка сложности алгоритмов (Big-O Notation)](06-go-core-algorithms-interview.md#261-оценка-сложности-алгоритмов-big-o-notation)
   - [Бинарный поиск: ручной, sort.Search и slices.BinarySearch](06-go-core-algorithms-interview.md#262-бинарный-поиск-ручной-sortsearch-и-slicesbinarysearch)
   - [Алгоритмы сортировки и стандартный пакет sort (QuickSort, pdqsort, sort.Interface)](06-go-core-algorithms-interview.md#263-алгоритмы-сортировки-и-стандартный-пакет-sort-quicksort-pdqsort-sortinterface)
   - [Рекурсия и Стек вызовов (Call Stack & Stack Overflow)](06-go-core-algorithms-interview.md#264-рекурсия-и-стек-вызовов-call-stack--stack-overflow)
   - [Парадигмы алгоритмов: «Разделяй и властвуй», Динамическое программирование и Жадные алгоритмы](06-go-core-algorithms-interview.md#265-парадигмы-алгоритмов-разделяй-и-властвуй-динамическое-программирование-и-жадные-алгоритмы)
   - [Поисковый обратный индекс (Inverted Index) и быстрая выборка документов](06-go-core-algorithms-interview.md#266-поисковый-обратный-индекс-inverted-index-и-быстрая-выборка-документов)
   - [Топ-4 паттерна алгоритмических интервью (Шаблоны с наглядными аналогиями)](06-go-core-algorithms-interview.md#267-топ-4-паттерна-алгоритмических-интервью-шаблоны-с-наглядными-аналогиями)

27. [**Классические Структуры Данных (Data Structures в Go)**](06-go-core-algorithms-interview.md#27-классические-структуры-данных-data-structures-в-go)

   - [Связные списки (container/list vs Slice)](06-go-core-algorithms-interview.md#271-связные-списки-containerlist-vs-slice)
   - [Стек (LIFO) и Очередь (FIFO)](06-go-core-algorithms-interview.md#272-стек-lifo-и-очередь-fifo)
   - [Бинарные деревья поиска (BST: Binary Search Tree)](06-go-core-algorithms-interview.md#273-бинарные-деревья-поиска-bst-binary-search-tree)
   - [Графы и алгоритмы обхода (Представление в Go, BFS, DFS и остовные деревья)](06-go-core-algorithms-interview.md#274-графы-и-алгоритмы-обхода-представление-в-go-bfs-dfs-и-остовные-деревья)
   - [Кольцевой список (Пакет container/ring: Round-Robin и циклический буфер)](06-go-core-algorithms-interview.md#275-кольцевой-список-пакет-containerring-round-robin-и-циклический-буфер)
   - [Двоичная куча и Очередь с приоритетами (Пакет container/heap)](06-go-core-algorithms-interview.md#276-двоичная-куча-и-очередь-с-приоритетами-пакет-containerheap)


28. [**Сетевое Программирование: TCP и UDP сокеты (Пакет net)**](04-go-core-backend-web.md#28-сетевое-программирование-tcp-и-udp-сокеты-пакет-net)

   - [Сетевые модели (OSI vs TCP/IP), Сокеты и Принципы транспорта](04-go-core-backend-web.md#280-сетевые-модели-osi-vs-tcpip-сокеты-и-принципы-транспорта)
   - [TCP-сервер и TCP-клиент на сокетах](04-go-core-backend-web.md#281-tcp-сервер-и-tcp-клиент-на-сокетах)
   - [UDP: Дейтаграммы и отличие от TCP](04-go-core-backend-web.md#282-udp-дейтаграммы-и-отличие-от-tcp)
   - [Таймауты на сокетах (SetDeadline)](04-go-core-backend-web.md#283-таймауты-на-сокетах-setdeadline)
   - [WebSockets: Полнодуплексный обмен в реальном времени](04-go-core-backend-web.md#284-websockets-полнодуплексный-обмен-в-реальном-времени)
   - [Безопасность сетевых соединений: TLS и mTLS в Go](04-go-core-backend-web.md#285-безопасность-сетевых-соединений-tls-и-mtls-в-go)
   - [Тестирование сетевых служб в памяти: net.Pipe()](04-go-core-backend-web.md#286-тестирование-сетевых-служб-в-памяти-netpipe)

29. [**Работа со временем: пакет time**](04-go-core-backend-web.md#29-работа-со-временем-пакет-time)

   - [Форматирование и парсинг дат (Референсное время Go)](04-go-core-backend-web.md#291-форматирование-и-парсинг-дат-референсное-время-go)
   - [Измерение времени и таймеры](04-go-core-backend-web.md#292-измерение-времени-и-таймеры)
   - [Ловушки при сравнении time.Time](04-go-core-backend-web.md#293-ловушки-при-сравнении-timetime)

30. [**Рефлексия: пакет reflect**](06-go-core-algorithms-interview.md#30-рефлексия-пакет-reflect)

   - [reflect.TypeOf и reflect.ValueOf](06-go-core-algorithms-interview.md#301-reflecttypeof-и-reflectvalueof)
   - [Инспекция полей структур и тегов](06-go-core-algorithms-interview.md#302-инспекция-полей-структур-и-тегов)
   - [Проверка интерфейса на равенство nil: Value.IsNil()](06-go-core-algorithms-interview.md#303-проверка-интерфейса-на-равенство-nil-через-рефлексию-valueisnil)
   - [Почему рефлексия медленная?](06-go-core-algorithms-interview.md#304-почему-рефлексия-медленная)

31. [**Низкоуровневый доступ: unsafe.Pointer**](06-go-core-algorithms-interview.md#31-низкоуровневый-доступ-unsafepointer)

   - [Три правила преобразования unsafe.Pointer](06-go-core-algorithms-interview.md#311-три-правила-преобразования-unsafepointer)
   - [Практический пример: Zero-Copy конверсия byte ↔ string](06-go-core-algorithms-interview.md#312-практический-пример-zero-copy-конверсия-byte-и-string)
   - [Взлом инкапсуляции (приватные поля) и диагностический флаг -gcflags=-d=checkptr](06-go-core-algorithms-interview.md#313-взлом-инкапсуляции-чтение-и-изменение-неэкспортируемых-приватных-полей-структур)

32. [**Приоритетная очередь: container/heap**](06-go-core-algorithms-interview.md#32-приоритетная-очередь-containerheap)

   - [Реализация интерфейса heap.Interface](06-go-core-algorithms-interview.md#321-реализация-интерфейса-heapinterface)
   - [Использование приоритетной очереди (Min-Heap)](06-go-core-algorithms-interview.md#322-использование-приоритетной-очереди-min-heap)

33. [**Эрудиция Go-разработчика: Новинки Go 1.21, 1.22, 1.23, 1.24 и новее**](06-go-core-algorithms-interview.md#33-эрудиция-go-разработчика-новинки-go-121-122-123-124-и-новее)

   - [Хронология ключевых изменений в релизах Go](06-go-core-algorithms-interview.md#331-хронология-ключевых-изменений-в-релизах-go)
   - [Зачем разработчику следить за релизами и как отвечать на собеседовании](06-go-core-algorithms-interview.md#332-зачем-разработчику-следить-за-релизами-и-как-отвечать-на-собеседовании)
   - [Глубокий разбор ключевых новинок: math/rand/v2, unique, weak и директива tool](06-go-core-algorithms-interview.md#333-глубокий-разбор-ключевых-новинок-mathrandv2-unique-weak-и-директива-tool)

34. [**Стайлгайды, Чистый Код и Культура Code Review (Uber & Google Style Guides)**](06-go-core-algorithms-interview.md#34-стайлгайды-чистый-код-и-культура-code-review-uber--google-style-guides)

   - [Золотые правила идиоматичного Go (Uber, Google, Effective Go)](06-go-core-algorithms-interview.md#341-золотые-правила-идиоматичного-go-по-мотивам-uber-go-style-guide-google-style-guide-и-effective-go)
   - [Простота vs Избыточность (Пример рефакторинга кода)](06-go-core-algorithms-interview.md#342-простота-vs-избыточность-пример-рефакторинга-кода)
   - [Практики Code Review в распределенной команде](06-go-core-algorithms-interview.md#343-практики-code-review-в-распределенной-команде)

35. [**Production-Ready сборка, Dockerfile и CI/CD линтинг**](06-go-core-algorithms-interview.md#35-production-ready-сборка-dockerfile-и-cicd-линтинг)

   - [Идиоматичный Multi-Stage Dockerfile (Alpine -> Scratch, 15 MB)](06-go-core-algorithms-interview.md#351-идиоматичный-multi-stage-dockerfile-alpine---scratch-15-mb)
   - [Флаги компиляции и оптимизация бинарника](06-go-core-algorithms-interview.md#352-флаги-компиляции-и-оптимизация-бинарника)
   - [Настройка golangci-lint и топ-10 линтеров](06-go-core-algorithms-interview.md#353-настройка-golangci-lint-и-топ-10-линтеров)

36. [**Экспресс-подготовка к собеседованию и Поведенческая секция (STAR)**](06-go-core-algorithms-interview.md#36-экспресс-подготовка-к-собеседованию-и-поведенческая-секция-star)

   - [Топ-25 ключевых вопросов и ответов технического собеседования](06-go-core-algorithms-interview.md#361-топ-25-ключевых-вопросов-и-ответов-технического-собеседования)
   - [Каверзные задачи Live-кодинга и квизы (Code Snippet Puzzles)](06-go-core-algorithms-interview.md#362-каверзные-задачи-live-кодинга-и-квизы-code-snippet-puzzles)
   - [Топ-5 смертельных ловушек в Live-кодинге](06-go-core-algorithms-interview.md#363-топ-5-смертельных-ловушек-в-live-кодинге)
   - [Поведенческое интервью (Behavioral) и методология STAR](06-go-core-algorithms-interview.md#364-поведенческое-интервью-behavioral-и-методология-star)
   - [Структура собеседования Go-разработчика и методология прохождения (Вебинар Дмитрия Титова, ВК)](06-go-core-algorithms-interview.md#365-структура-собеседования-go-разработчика-и-методология-прохождения-вебинар-дмитрия-титова-вк)
   - [Золотая библиотека первоисточников и ресурсов Go-разработчика](06-go-core-algorithms-interview.md#366-золотая-библиотека-первоисточников-и-ресурсов-go-разработчика)

---

