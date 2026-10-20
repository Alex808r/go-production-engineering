# GoSearch — Поисковый движок на Go

**GoSearch** — курсовой проект, разрабатываемый на протяжении практических модулей курса Go. Представляет собой поисковый движок, включающий сетевой робот-краулер, механизм обратного (инвертированного) индексирования и быстрый бинарный поиск документов.

---

## 📋 Архитектура и структура проекта

```text
GoSearch/
├── cmd/
│   └── gosearch/
│       └── gosearh.go         # CLI и интерактивный REPL поиска
└── pkg/
    ├── crawler/               # Сетевой поисковый робот
    │   ├── crawler.go         # Интерфейс Interface и структура Document
    │   ├── membot/            # In-memory эмулятор для локального тестирования
    │   └── spider/            # Реальный HTTP-краулер с разбором HTML
    └── index/                 # Поисковый инвертированный индекс
        ├── index.go           # Inverted Index, токенизация, BinarySearch
        └── index_test.go      # Unit-тесты и бенчмарки
```

---

## 🔍 Ключевые компоненты

### 1. Краулер ([`pkg/crawler`](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/GoSearch/pkg/crawler/crawler.go))
Интерфейс `crawler.Interface`:
```go
type Interface interface {
    Scan(url string, depth int) ([]Document, error)
}
```
- [`spider.Service`](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/GoSearch/pkg/crawler/spider): рекурсивно обходит веб-страницы по ссылкам `href`, извлекая URL и Title с заданной глубиной `depth`.
- [`membot.Service`](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/GoSearch/pkg/crawler/membot/membot.go): возвращает статичный набор документов для офлайн-разработки и юнит-тестов.

### 2. Инвертированный индекс ([`pkg/index`](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/GoSearch/pkg/index/index.go))
- Индекс хранит отображение `map[string][]int`: **слово $\to$ срез идентификаторов документов (ID)**.
- При индексации заголовки документов токенизируются, приводятся к нижнему регистру и очищаются от знаков препинания (`cleanWord`).
- Дубликаты слов внутри одного документа устраняются через `seen[word]`.

### 3. Бинарный поиск документов `O(log N)`
- Документы предварительно сортируются по возрастанию ID:
  ```go
  sort.Slice(docs, func(i, j int) bool {
      return docs[i].ID < docs[j].ID
  })
  ```
- Для быстрого извлечения полного документа по ID, возвращенному из индекса, используется функция `BinarySearch`:
  ```go
  func BinarySearch(docs []crawler.Document, targetID int) int {
      idx := sort.Search(len(docs), func(i int) bool {
          return docs[i].ID >= targetID
      })
      if idx < len(docs) && docs[idx].ID == targetID {
          return idx
      }
      return -1
  }
  ```

---

## 📊 Результаты бенчмарков

```text
BenchmarkIndex_Add-8         187532      6150 ns/op
BenchmarkIndex_Search-8     7128392       168 ns/op
BenchmarkBinarySearch-8    41294821      28.9 ns/op
```

- Поиск по хэш-таблице индекса выполняется за **~168 нс**.
- Извлечение документа через бинарный поиск `O(log N)` среди 1000 элементов занимает всего **~29 нс**.

---

## 🚀 Запуск и использование

### 1. Тестирование с детектором гонок
```bash
go test -v -race -bench=. ./GoSearch/...
```

### 2. Быстрый поиск через membot
```bash
# Поиск конкретного слова
go run ./GoSearch/cmd/gosearch/gosearh.go -mem -s яндекс

# Интерактивный режим поиска
go run ./GoSearch/cmd/gosearch/gosearh.go -mem
```

### 3. Полноценное сканирование сети
```bash
go run ./GoSearch/cmd/gosearch/gosearh.go -s Go
```

---

## 🔗 Связанные материалы курса
- [03-algorithms — Поиск и сортировка](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/03-algorithms/03-algorithms.md)
- [04-datastructs — Структуры данных и деревья поиска](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/04-datastructs/04-datastructs.md)
- [05-io — Потоковый ввод-вывод и чтение](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/05-io/05-io.md)
