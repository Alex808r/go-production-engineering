/*
    =============================================================================
    Схема учебной реляционной базы данных "Книги" (books_db) для PostgreSQL.
    Демонстрирует:
    1. Третью нормальную форму (3NF).
    2. Первичные (PRIMARY KEY) и внешние ключи (FOREIGN KEY с ON DELETE/UPDATE CASCADE).
    3. Ограничения целостности (CHECK, UNIQUE, NOT NULL).
    4. Функциональные индексы (B-Tree по lower(title)).
    5. Процедурный язык PL/pgSQL и функции-триггеры (BEFORE INSERT OR UPDATE).
    =============================================================================
*/

/*
    Удаление таблиц в обратном порядке создания
    (сначала зависимые таблицы с внешними ключами, затем родительские).
*/
DROP TABLE IF EXISTS books_authors CASCADE;
DROP TABLE IF EXISTS books CASCADE;
DROP TABLE IF EXISTS authors CASCADE;
DROP TABLE IF EXISTS publishers CASCADE;

/*
    1. Таблица авторов (authors)
    Родительская таблица: на неё ссылается таблица связей books_authors.
*/
CREATE TABLE authors (
    id SERIAL PRIMARY KEY,                    -- автоинкрементный 32-битный суррогатный ключ
    first_name TEXT NOT NULL DEFAULT '',      -- имя автора
    last_name TEXT NOT NULL DEFAULT '',       -- фамилия автора
    year_of_birth INTEGER NOT NULL DEFAULT 0  -- год рождения
);

/*
    2. Таблица издательств (publishers)
    Родительская таблица: отношение 1:N с таблицей books (у книги одно издательство).
*/
CREATE TABLE publishers (
    id SERIAL PRIMARY KEY,                    -- автоинкрементный первичный ключ
    name TEXT NOT NULL,                       -- наименование издательства
    website TEXT NOT NULL DEFAULT ''          -- веб-сайт
);

/*
    3. Таблица книг (books)
    Основная сущность каталога.
*/
CREATE TABLE books (
    id BIGSERIAL PRIMARY KEY,                 -- 64-битный автоинкрементный первичный ключ
    isbn BIGINT UNIQUE,                       -- уникальный номер ISBN (международный стандарт)
    title TEXT NOT NULL,                      -- название произведения
    year INTEGER DEFAULT 0,                   -- год выпуска (валидируется триггером)
    public_domain BOOLEAN DEFAULT FALSE,      -- флаг общественного достояния (свободное распространение)
    -- Внешний ключ на издательство:
    -- ON DELETE CASCADE: при удалении издательства удаляются связанные книги
    -- ON UPDATE CASCADE: при обновлении id издательства обновляется publisher_id в книгах
    publisher_id INTEGER REFERENCES publishers(id) ON DELETE CASCADE ON UPDATE CASCADE DEFAULT 0,
    price INTEGER DEFAULT 0 CHECK (price >= 0), -- цена в копейках/центах; CHECK гарантирует неотрицательность
    genres TEXT[] DEFAULT '{"не указано"}',   -- массив текстовых меток жанров (поддержка типов массивов в Postgres)
    info JSONB DEFAULT '{}'                   -- полуструктурированные метаданные (JSONB с бинарным индексированием)
);

-- Индекс B-Tree по выражению lower(title) для быстрого регистронезависимого поиска книг
CREATE INDEX IF NOT EXISTS books_title_idx ON books USING btree (lower(title));

/*
    4. Связующая таблица связи многие-ко-многим (books_authors)
    Реализует отношение N:M: у одной книги может быть несколько соавторов,
    и один автор может написать множество книг.
*/
CREATE TABLE books_authors (
    id BIGSERIAL PRIMARY KEY,
    book_id BIGINT NOT NULL REFERENCES books(id) ON DELETE CASCADE ON UPDATE CASCADE,
    author_id INTEGER NOT NULL REFERENCES authors(id) ON DELETE CASCADE ON UPDATE CASCADE,
    -- Уникальный составной ключ предотвращает дублирование одного и того же автора у книги
    UNIQUE(book_id, author_id)
);

-- Индекс на внешний ключ author_id для быстрого поиска всех книг заданного автора
CREATE INDEX IF NOT EXISTS books_authors_author_idx ON books_authors (author_id);

/*
    5. Функция-триггер на PL/pgSQL для бизнес-валидации года издания книги.
    Проверяет, что год книги не опережает текущую дату более чем на 10 лет
    и не меньше текущего года минус 10 лет (учебное правило).
*/
CREATE OR REPLACE FUNCTION check_book_year()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.year > (SELECT (EXTRACT(YEAR FROM CURRENT_DATE) - 10)) AND
       NEW.year < (SELECT (EXTRACT(YEAR FROM CURRENT_DATE) + 10))
    THEN
        RETURN NEW; -- Проверка пройдена, вставка/обновление разрешены
    ELSE
        -- Исключение откатывает транзакцию и возвращает ошибку клиенту
        RAISE EXCEPTION 'Некорректный год издания книги %: допустим диапазон от % до %',
            NEW.year,
            (EXTRACT(YEAR FROM CURRENT_DATE) - 10),
            (EXTRACT(YEAR FROM CURRENT_DATE) + 10);
    END IF;
END;
$$ LANGUAGE plpgsql;

-- Регистрация триггера: срабатывает BEFORE INSERT OR UPDATE для каждой строки
CREATE OR REPLACE TRIGGER trg_check_book_year
BEFORE INSERT OR UPDATE ON books
FOR EACH ROW
EXECUTE FUNCTION check_book_year();