# Домашнее задание №15: Реляционные базы данных (Схема БД онлайн-кинотеатра)

Решение домашнего задания по курсу «Современная разработка на Go» (Занятие №15: Реляционные БД).  
Реализовано в отдельной директории `homework-15` без изменения или затирания предыдущих решений.

---

## 📑 Содержание

1. [Задача №1: Схема базы данных (schema.sql)](#задача-1-схема-базы-данных)
2. [Задача №2: Демонстрационные данные (data.sql)](#задача-2-демонстрационные-данные)
3. [Задача №3: Аналитические SQL-запросы (queries.sql)](#задача-3-аналитические-sql-запросы)
4. [Инструкция по развертыванию и проверке](#инструкция-по-развертыванию)

---

## Задача №1: Схема базы данных

Файл [`schema.sql`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/15-sql/homework-15/schema.sql) описывает DDL-структуру базы данных для кинопортала.

### Моделируемые сущности и связи:
- **`studios` (Киностудии):**
  - `id SERIAL PRIMARY KEY`, `name TEXT NOT NULL UNIQUE`.
- **`actors` (Актёры):**
  - `id SERIAL PRIMARY KEY`, `first_name TEXT NOT NULL`, `last_name TEXT NOT NULL`, `birth_date DATE NOT NULL`.
- **`directors` (Режиссёры):**
  - `id SERIAL PRIMARY KEY`, `first_name TEXT NOT NULL`, `last_name TEXT NOT NULL`, `birth_date DATE NOT NULL`.
- **`movies` (Фильмы):**
  - `id BIGSERIAL PRIMARY KEY`.
  - `title TEXT NOT NULL`.
  - `release_year INTEGER NOT NULL CHECK (release_year >= 1800)` — ограничение на минимальный год выхода.
  - `box_office BIGINT NOT NULL DEFAULT 0 CHECK (box_office >= 0)` — кассовые сборы.
  - `rating TEXT NOT NULL CHECK (rating IN ('PG-10', 'PG-13', 'PG-18'))` — ограничение на допустимый возрастной рейтинг.
  - `studio_id INTEGER NOT NULL REFERENCES studios(id) ON DELETE RESTRICT ON UPDATE CASCADE` — у каждого фильма строго одна студия.
  - `CONSTRAINT unique_title_year UNIQUE (title, release_year)` — запрет двух фильмов с одинаковым названием в один и тот же год.
- **Связующие таблицы многие-ко-многим (N:M):**
  - `movies_actors` (`movie_id`, `actor_id`, `PRIMARY KEY (movie_id, actor_id)`).
  - `movies_directors` (`movie_id`, `director_id`, `PRIMARY KEY (movie_id, director_id)`).
- **Индексы и триггеры:**
  - B-Tree индексы на внешние ключи, даты и названия (`lower(title)`).
  - PL/pgSQL функция-триггер `check_movie_year()`: проверяет, что год релиза не превышает текущую дату более чем на 10 лет.

---

## Задача №2: Демонстрационные данные

Файл [`data.sql`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/15-sql/homework-15/data.sql) наполняет БД тестовым набором:
- 5 мировых киностудий.
- 6 известных режиссёров (Christopher Nolan, Steven Spielberg, Quentin Tarantino, Denis Villeneuve, David Lynch, Martin Scorsese).
- 10 известных актёров.
- 12 культовых фильмов (включая блокбастеры со сборами > 1000 млн, такие как Titanic, Avatar, The Dark Knight, Jurassic Park, и дубликаты названий: фильм "Dune" 1984 и 2021 годов).
- Все ассоциативные связи `movies_actors` и `movies_directors`.

---

## Задача №3: Аналитические SQL-запросы

Файл [`queries.sql`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/15-sql/homework-15/queries.sql) содержит решение 10 типовых аналитических задач:
1. Выборка фильмов с названием студии (`JOIN studios`).
2. Выборка фильмов для некоторого актёра (`WHERE a.first_name = ... AND a.last_name = ...`).
3. Подсчёт фильмов для некоторого режиссёра (`COUNT(md.movie_id)`).
4. Выборка фильмов для нескольких режиссёров из списка (подзапрос `WHERE m.id IN (SELECT ...)`).
5. Подсчёт количества фильмов для каждого актёра (`GROUP BY`, `COUNT`).
6. Выборка актёров и режиссёров, участвовавших более чем в 2 фильмах (`UNION ALL`, `HAVING COUNT(*) > 2`).
7. Подсчёт количества фильмов со сборами больше 1000 ($ млн) (`WHERE box_office > 1000`).
8. Подсчёт уникальных режиссёров, чьи фильмы собрали свыше 1000 ($ млн) (`COUNT(DISTINCT md.director_id)`).
9. Выборка уникальных фамилий актёров (`SELECT DISTINCT last_name`).
10. Подсчёт фильмов, имеющих дубли по названию (`GROUP BY title HAVING COUNT(*) > 1`).

---

## Инструкция по развертыванию

### Прогон unit-тестов синтаксиса и ограничений схемы:

```bash
cd go_course_4
go test -v ./lessons/15-sql/homework-15/...
```

### Развертывание в PostgreSQL:

```bash
# Применение схемы и наполнение данными
psql -h localhost -U postgres -d movies_db -f ./lessons/15-sql/homework-15/schema.sql
psql -h localhost -U postgres -d movies_db -f ./lessons/15-sql/homework-15/data.sql

# Выполнение аналитических запросов
psql -h localhost -U postgres -d movies_db -f ./lessons/15-sql/homework-15/queries.sql
```
