/*
    Схема базы данных для информационного сайта о фильмах.
    Занятие №15. Реляционные БД. Домашнее задание (Задача №1).

    Правила предметной области:
    - Сущности: фильмы (movies), актеры (actors), режиссеры (directors), студии (studios).
    - У каждого фильма одна киностудия (связь 1:N).
    - У каждого фильма может быть много актеров и много режиссеров (связи N:M).
    - Рейтинг фильма: строго из набора ('PG-10', 'PG-13', 'PG-18').
    - Год выхода фильма: не может быть меньше 1800.
    - В один год не может быть двух фильмов с одинаковым названием: UNIQUE(title, release_year).
*/

-- Удаляем таблицы в порядке, обратном созданию (с учетом внешних ключей)
DROP TABLE IF EXISTS movies_actors;
DROP TABLE IF EXISTS movies_directors;
DROP TABLE IF EXISTS movies;
DROP TABLE IF EXISTS studios;
DROP TABLE IF EXISTS actors;
DROP TABLE IF EXISTS directors;

-- 1. Студии (studios)
CREATE TABLE studios (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

-- 2. Актеры (actors)
CREATE TABLE actors (
    id SERIAL PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    birth_date DATE NOT NULL
);

-- 3. Режиссеры (directors)
CREATE TABLE directors (
    id SERIAL PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    birth_date DATE NOT NULL
);

-- 4. Фильмы (movies)
CREATE TABLE movies (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    release_year INTEGER NOT NULL CHECK (release_year >= 1800),
    box_office BIGINT NOT NULL DEFAULT 0 CHECK (box_office >= 0),
    rating TEXT NOT NULL CHECK (rating IN ('PG-10', 'PG-13', 'PG-18')),
    studio_id INTEGER NOT NULL REFERENCES studios(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    CONSTRAINT unique_title_year UNIQUE (title, release_year)
);

-- Индексы для ускорения поиска и фильтрации
CREATE INDEX IF NOT EXISTS movies_title_idx ON movies USING btree (lower(title));
CREATE INDEX IF NOT EXISTS movies_release_year_idx ON movies (release_year);
CREATE INDEX IF NOT EXISTS movies_box_office_idx ON movies (box_office);
CREATE INDEX IF NOT EXISTS movies_studio_id_idx ON movies (studio_id);

-- 5. Связующая таблица Фильмы <-> Актеры (N:M)
CREATE TABLE movies_actors (
    movie_id BIGINT NOT NULL REFERENCES movies(id) ON DELETE CASCADE ON UPDATE CASCADE,
    actor_id INTEGER NOT NULL REFERENCES actors(id) ON DELETE CASCADE ON UPDATE CASCADE,
    PRIMARY KEY (movie_id, actor_id)
);
CREATE INDEX IF NOT EXISTS movies_actors_actor_idx ON movies_actors (actor_id);

-- 6. Связующая таблица Фильмы <-> Режиссеры (N:M)
CREATE TABLE movies_directors (
    movie_id BIGINT NOT NULL REFERENCES movies(id) ON DELETE CASCADE ON UPDATE CASCADE,
    director_id INTEGER NOT NULL REFERENCES directors(id) ON DELETE CASCADE ON UPDATE CASCADE,
    PRIMARY KEY (movie_id, director_id)
);
CREATE INDEX IF NOT EXISTS movies_directors_director_idx ON movies_directors (director_id);

-- 7. Функция и триггер бизнес-логики: проверка года выхода
-- Новая книга/фильм не может быть выпущена позже 10 лет с текущей даты (слайды 9 и 13)
CREATE OR REPLACE FUNCTION check_movie_year()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.release_year <= (EXTRACT(YEAR FROM CURRENT_DATE) + 10) THEN
        RETURN NEW;
    ELSE
        RAISE EXCEPTION 'Некорректный год фильма %: дата релиза не может превышать текущий год более чем на 10 лет', NEW.release_year;
    END IF;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER trg_check_movie_year
BEFORE INSERT OR UPDATE ON movies
FOR EACH ROW
EXECUTE FUNCTION check_movie_year();
