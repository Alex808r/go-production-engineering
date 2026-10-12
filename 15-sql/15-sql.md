# 15-sql — Реляционные базы данных и язык SQL

В данной директории изучаются проектирование реляционных схем данных для PostgreSQL, язык запросов SQL (DDL, DML), нормализация данных (3NF), внешние ключи (Foreign Keys), многотабличные соединения (`JOIN`), агрегатные функции, индексы B-Tree и функции-триггеры на PL/pgSQL.

---

## 📋 Структура каталога

```text
15-sql/
├── books_db/        # Учебная база данных книжного каталога
│   ├── 1-schema.sql # DDL: создание таблиц (authors, publishers, books, books_authors), B-Tree индексов и триггера
│   └── 2-data.sql   # DML: тестовые данные, последовательности (SERIAL), примеры JOIN и разбор исключения года издания
└── homework-15/     # Практическое ДЗ: База данных онлайн-кинотеатра
    ├── schema.sql   # DDL: сущности (studios, actors, directors, movies), M:N таблицы, CHECK, UNIQUE, триггер года
    ├── data.sql     # DML: киностудии, культовые режиссеры, актеры, фильмы и связи
    ├── queries.sql  # Аналитические SQL-запросы: JOIN, GROUP BY, HAVING, UNION ALL, подзапросы, агрегаты
    ├── db_test.go   # Модульные тесты схемы, ограничений, наличия всех требуемых запросов и бенчмарки
    └── README.md    # Подробное описание решений всех задач домашнего задания
```

---

## 🔍 Ключевые концепции

1. **Реляционная модель и нормализация:**
   - Моделирование отношений:
     - «Один-ко-многим» (1:N): фильмы и киностудии (`movies.studio_id -> studios.id`), книги и издательства (`books.publisher_id -> publishers.id`);
     - «Многие-ко-многим» (N:M): фильмы и актеры (`movies_actors`), фильмы и режиссеры (`movies_directors`), книги и соавторы (`books_authors`).
   - Защита целостности: составные уникальные ключи `UNIQUE(book_id, author_id)` и `PRIMARY KEY (movie_id, actor_id)` исключают дублирование связей.

2. **DDL (Data Definition Language) и ограничения:**
   - Первичные ключи: `SERIAL` и `BIGSERIAL`.
   - Ограничения целостности (Constraints):
     - `NOT NULL` и `DEFAULT`;
     - `UNIQUE (title, release_year)` — запрет одинаковых названий в один год;
     - `CHECK (release_year >= 1800)` и `CHECK (price >= 0)`;
     - `CHECK (rating IN ('PG-10', 'PG-13', 'PG-18'))` — валидация возрастных рейтингов.
   - Каскадные операции: `ON DELETE CASCADE / RESTRICT`, `ON UPDATE CASCADE`.

3. **Индексы и оптимизация (B-Tree):**
   - Индексирование внешних ключей для ускорения соединений (`JOIN`).
   - Функциональные индексы по вычисляемым выражениям: `CREATE INDEX ... ON books USING btree (lower(title))`.

4. **Процедурная логика PL/pgSQL (Триггеры):**
   - Функции-триггеры (`RETURNS TRIGGER`) с доступом к записи `NEW`.
   - Проверка диапазонов дат относительно `CURRENT_DATE`.
   - Прерывание недопустимых операций через `RAISE EXCEPTION`.

5. **Аналитический SQL (DML):**
   - Внутренние и внешние соединения: `INNER JOIN`, `LEFT JOIN`.
   - Группировка и фильтрация агрегатов: `GROUP BY ... HAVING COUNT(*) > 1`.
   - Объединение результатов независимых выборок: `UNION ALL`.
   - Подзапросы с предикатом `IN (SELECT ...)`.

---

## 🚀 Запуск и проверка

```bash
# Запуск валидационных unit-тестов и бенчмарков синтаксиса
go test -v -race -bench=. ./15-sql/homework-15/...

# Развертывание учебной базы books_db в PostgreSQL
psql -h localhost -U postgres -d postgres -f ./15-sql/books_db/1-schema.sql
psql -h localhost -U postgres -d postgres -f ./15-sql/books_db/2-data.sql

# Развертывание базы данных кинотеатра (ДЗ-15) и запуск аналитики
psql -h localhost -U postgres -d movies_db -f ./15-sql/homework-15/schema.sql
psql -h localhost -U postgres -d movies_db -f ./15-sql/homework-15/data.sql
psql -h localhost -U postgres -d movies_db -f ./15-sql/homework-15/queries.sql
```

---

## 🔗 Связанные материалы
- **Теория:** [Том 4: Backend и Web (Раздел 13: Реляционные БД)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/04-go-core-backend-web.md)
- **Практическое руководство:** [sql-livecoding-guide.md](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/sql-livecoding-guide.md)
- **Практическое ДЗ:** [homework-15](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/15-sql/homework-15) (База данных онлайн-кинотеатра: студии, актеры, режиссеры, фильмы и 10 аналитических запросов)

