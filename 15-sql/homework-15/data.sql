/*
    Наполнение базы данных демонстрационными данными.
    Занятие №15. Реляционные БД. Домашнее задание (Задача №2).
*/

-- 1. Студии
INSERT INTO studios (id, name) VALUES
(1, 'Warner Bros. Pictures'),
(2, 'Universal Pictures'),
(3, 'Paramount Pictures'),
(4, 'Columbia Pictures'),
(5, '20th Century Studios');

ALTER SEQUENCE studios_id_seq RESTART WITH 100;

-- 2. Режиссеры
INSERT INTO directors (id, first_name, last_name, birth_date) VALUES
(1, 'Christopher', 'Nolan', '1970-07-30'),
(2, 'Steven', 'Spielberg', '1946-12-18'),
(3, 'Quentin', 'Tarantino', '1963-03-27'),
(4, 'Denis', 'Villeneuve', '1967-10-03'),
(5, 'David', 'Lynch', '1946-01-20'),
(6, 'Martin', 'Scorsese', '1942-11-17');

ALTER SEQUENCE directors_id_seq RESTART WITH 100;

-- 3. Актеры
INSERT INTO actors (id, first_name, last_name, birth_date) VALUES
(1, 'Leonardo', 'DiCaprio', '1974-11-11'),
(2, 'Christian', 'Bale', '1974-01-30'),
(3, 'Tom', 'Hardy', '1977-09-15'),
(4, 'Cillian', 'Murphy', '1976-05-25'),
(5, 'Brad', 'Pitt', '1963-12-18'),
(6, 'Margot', 'Robbie', '1990-07-02'),
(7, 'Kyle', 'MacLachlan', '1959-02-22'),
(8, 'Timothee', 'Chalamet', '1995-12-27'),
(9, 'Michael', 'Caine', '1933-03-14'),
(10, 'Robert', 'De Niro', '1943-08-17');

ALTER SEQUENCE actors_id_seq RESTART WITH 100;

-- 4. Фильмы
-- Внимание: для демонстрации задачи "дубли по названию" включен фильм "Dune" (1984 и 2021 годов)
-- Сборы (box_office) указаны в миллионах долларов (для проверки фильтра > 1000)
INSERT INTO movies (id, title, release_year, box_office, rating, studio_id) VALUES
(1, 'Inception', 2010, 836, 'PG-13', 1),
(2, 'The Dark Knight', 2008, 1006, 'PG-13', 1),
(3, 'Oppenheimer', 2023, 957, 'PG-18', 2),
(4, 'Interstellar', 2014, 731, 'PG-13', 3),
(5, 'Once Upon a Time in Hollywood', 2019, 377, 'PG-18', 4),
(6, 'Dune', 1984, 30, 'PG-13', 2),
(7, 'Dune', 2021, 402, 'PG-13', 1),
(8, 'Jurassic Park', 1993, 1046, 'PG-13', 2),
(9, 'Schindler''s List', 1993, 322, 'PG-18', 2),
(10, 'The Wolf of Wall Street', 2013, 406, 'PG-18', 3),
(11, 'Titanic', 1997, 2257, 'PG-13', 3),
(12, 'Avatar', 2009, 2923, 'PG-13', 5);

ALTER SEQUENCE movies_id_seq RESTART WITH 100;

-- 5. Связи Фильмы <-> Актеры
INSERT INTO movies_actors (movie_id, actor_id) VALUES
(1, 1), -- Inception: Leonardo DiCaprio
(1, 3), -- Inception: Tom Hardy
(1, 4), -- Inception: Cillian Murphy
(1, 9), -- Inception: Michael Caine
(2, 2), -- The Dark Knight: Christian Bale
(2, 4), -- The Dark Knight: Cillian Murphy
(2, 9), -- The Dark Knight: Michael Caine
(3, 4), -- Oppenheimer: Cillian Murphy
(4, 9), -- Interstellar: Michael Caine
(5, 1), -- Once Upon a Time in Hollywood: Leonardo DiCaprio
(5, 5), -- Once Upon a Time in Hollywood: Brad Pitt
(5, 6), -- Once Upon a Time in Hollywood: Margot Robbie
(6, 7), -- Dune 1984: Kyle MacLachlan
(7, 8), -- Dune 2021: Timothee Chalamet
(10, 1), -- Wolf of Wall Street: Leonardo DiCaprio
(10, 6), -- Wolf of Wall Street: Margot Robbie
(11, 1); -- Titanic: Leonardo DiCaprio

-- 6. Связи Фильмы <-> Режиссеры
INSERT INTO movies_directors (movie_id, director_id) VALUES
(1, 1), -- Inception: Christopher Nolan
(2, 1), -- The Dark Knight: Christopher Nolan
(3, 1), -- Oppenheimer: Christopher Nolan
(4, 1), -- Interstellar: Christopher Nolan
(5, 3), -- Once Upon a Time in Hollywood: Quentin Tarantino
(6, 5), -- Dune 1984: David Lynch
(7, 4), -- Dune 2021: Denis Villeneuve
(8, 2), -- Jurassic Park: Steven Spielberg
(9, 2), -- Schindler's List: Steven Spielberg
(10, 6); -- The Wolf of Wall Street: Martin Scorsese
