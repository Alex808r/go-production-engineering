/*
    =============================================================================
    Первичное наполнение учебной базы данных (DML) и демонстрация работы триггеров.
    =============================================================================
*/

-- 1. Служебные записи со специальным id=0 для значений по умолчанию
INSERT INTO authors (id, first_name, last_name) VALUES (0, 'автор не указан', '');
INSERT INTO authors (id, first_name, last_name) VALUES (1, 'автор неизвестен', '');

-- Внимание: так как id=0 и id=1 вставлены принудительно в обход генератора serial,
-- необходимо пересинхронизировать автоинкрементную последовательность:
ALTER SEQUENCE authors_id_seq RESTART WITH 100;

INSERT INTO publishers (id, name) VALUES (0, 'издательство не указано');
ALTER SEQUENCE publishers_id_seq RESTART WITH 100;

-- 2. Реальные издательства
INSERT INTO publishers (name, website) VALUES 
('O''Reilly Media', 'https://www.oreilly.com'),
('Manning Publications', 'https://www.manning.com'),
('Addison-Wesley', 'https://www.informit.com/imprint/addison-wesley');

-- 3. Реальные авторы
INSERT INTO authors (first_name, last_name, year_of_birth) VALUES 
('Alan', 'Donovan', 1976),
('Brian', 'Kernighan', 1942),
('Rob', 'Pike', 1956),
('Martin', 'Fowler', 1963);

-- 4. Вставка книг
-- Текущий год в базе данных динамически определяется через CURRENT_DATE.
-- Для успешной вставки год должен находиться в диапазоне [CURRENT_YEAR - 10, CURRENT_YEAR + 10].
INSERT INTO books (isbn, title, year, public_domain, publisher_id, price, genres, info) VALUES 
(9780134190440, 'The Go Programming Language', 2016, FALSE, 102, 3500, '{"Programming", "Go"}', '{"pages": 380, "rating": 4.9}'),
(9780201616224, 'The Pragmatic Programmer', 2019, FALSE, 102, 4200, '{"Software Engineering"}', '{"pages": 352, "rating": 4.8}'),
(9781617298691, '100 Go Mistakes and How to Avoid Them', 2022, FALSE, 101, 3900, '{"Programming", "Go"}', '{"pages": 384, "rating": 4.9}'),
(9780596007126, 'Head First Design Patterns', 2020, FALSE, 100, 2800, '{"Patterns", "Architecture"}', '{"pages": 650, "rating": 4.7}');

-- 5. Связи авторов и книг (M:N)
-- "The Go Programming Language" написана в соавторстве: Alan Donovan и Brian Kernighan
INSERT INTO books_authors (book_id, author_id) VALUES 
(1, 100), -- Donovan
(1, 101); -- Kernighan

-- =============================================================================
-- ВОПРОС: insert into books (title, year) VALUES ('Book', 2050); — что будет в результате выполнения?
--
-- ОТВЕТ:
-- Операция завершится с ОШИБКОЙ и откатом транзакции:
--   ERROR: Некорректный год издания книги 2050: дата релиза не может превышать текущий год более чем на 10 лет
--
-- Причина:
-- Таблица books защищена триггером trg_check_book_year, вызывающим функцию check_book_year().
-- Условие триггера требует, чтобы год был меньше текущего года + 10 (например, 2026 + 10 = 2036).
-- Для года 2050 проверка NEW.year < EXTRACT(YEAR FROM CURRENT_DATE) + 10 возвращает FALSE,
-- и триггер выполняет инструкцию RAISE EXCEPTION, прерывающую операцию INSERT.
-- =============================================================================

-- Пример аналитического запроса для проверки наполнения:
-- Выборка книг с авторами и издательствами через JOIN
SELECT 
    b.id,
    b.title,
    b.year,
    p.name AS publisher_name,
    STRING_AGG(a.first_name || ' ' || a.last_name, ', ') AS authors_list
FROM books b
JOIN publishers p ON b.publisher_id = p.id
LEFT JOIN books_authors ba ON b.id = ba.book_id
LEFT JOIN authors a ON ba.author_id = a.id
GROUP BY b.id, b.title, b.year, p.name
ORDER BY b.year DESC;