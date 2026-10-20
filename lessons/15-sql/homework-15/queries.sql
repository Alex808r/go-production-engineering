/*
    Аналитические SQL-запросы к базе данных онлайн-кинотеатра.
    Занятие №15. Реляционные БД. Домашнее задание (Задача №3).
*/

-- 1. Выборка фильмов с названием студии
SELECT 
    m.id,
    m.title,
    m.release_year,
    m.rating,
    m.box_office,
    s.name AS studio_name
FROM movies m
JOIN studios s ON m.studio_id = s.id
ORDER BY m.release_year DESC;

-- 2. Выборка фильмов для некоторого актёра (например, 'Leonardo DiCaprio')
SELECT 
    m.id,
    m.title,
    m.release_year,
    m.box_office,
    a.first_name || ' ' || a.last_name AS actor_name
FROM movies m
JOIN movies_actors ma ON m.id = ma.movie_id
JOIN actors a ON ma.actor_id = a.id
WHERE a.first_name = 'Leonardo' AND a.last_name = 'DiCaprio'
ORDER BY m.release_year;

-- 3. Подсчёт фильмов для некоторого режиссёра (например, 'Christopher Nolan')
SELECT 
    d.first_name || ' ' || d.last_name AS director_name,
    COUNT(md.movie_id) AS total_movies
FROM directors d
JOIN movies_directors md ON d.id = md.director_id
WHERE d.first_name = 'Christopher' AND d.last_name = 'Nolan'
GROUP BY d.id, d.first_name, d.last_name;

-- 4. Выборка фильмов для нескольких режиссёров из списка с использованием подзапроса
SELECT 
    m.id,
    m.title,
    m.release_year,
    m.box_office
FROM movies m
WHERE m.id IN (
    SELECT md.movie_id
    FROM movies_directors md
    JOIN directors d ON md.director_id = d.id
    WHERE d.last_name IN ('Nolan', 'Spielberg', 'Tarantino')
)
ORDER BY m.release_year;

-- 5. Подсчёт количества фильмов для каждого актёра
SELECT 
    a.id,
    a.first_name || ' ' || a.last_name AS actor_name,
    COUNT(ma.movie_id) AS movies_count
FROM actors a
LEFT JOIN movies_actors ma ON a.id = ma.actor_id
GROUP BY a.id, a.first_name, a.last_name
ORDER BY movies_count DESC, actor_name;

-- 6. Выборка актёров и режиссёров, участвовавших более чем в 2 фильмах
-- Используем UNION ALL и фильтрацию агрегатов через HAVING
SELECT 
    'Actor' AS role,
    a.first_name || ' ' || a.last_name AS full_name,
    COUNT(ma.movie_id) AS works_count
FROM actors a
JOIN movies_actors ma ON a.id = ma.actor_id
GROUP BY a.id, a.first_name, a.last_name
HAVING COUNT(ma.movie_id) > 2

UNION ALL

SELECT 
    'Director' AS role,
    d.first_name || ' ' || d.last_name AS full_name,
    COUNT(md.movie_id) AS works_count
FROM directors d
JOIN movies_directors md ON d.id = md.director_id
GROUP BY d.id, d.first_name, d.last_name
HAVING COUNT(md.movie_id) > 2

ORDER BY works_count DESC, full_name;

-- 7. Подсчёт количества фильмов со сборами больше 1000 ($ млн)
SELECT 
    COUNT(*) AS blockbuster_count
FROM movies
WHERE box_office > 1000;

-- 8. Подсчитать количество уникальных режиссёров, фильмы которых собрали больше 1000 ($ млн)
SELECT 
    COUNT(DISTINCT md.director_id) AS successful_directors_count
FROM movies m
JOIN movies_directors md ON m.id = md.movie_id
WHERE m.box_office > 1000;

-- 9. Выборка различных фамилий актёров
SELECT DISTINCT 
    last_name
FROM actors
ORDER BY last_name;

-- 10. Подсчёт количества фильмов, имеющих дубли по названию
SELECT 
    title,
    COUNT(*) AS duplicates_count
FROM movies
GROUP BY title
HAVING COUNT(*) > 1;
