-- Схема базы данных PostgreSQL для сервиса Lynks (Урок 20: Итоговый проект)
CREATE TABLE IF NOT EXISTS links (
    id SERIAL PRIMARY KEY,
    short_url VARCHAR(64) UNIQUE NOT NULL,
    destination TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Индекс для эффективного поиска длинной гиперссылки по короткой (требование ТЗ)
CREATE INDEX IF NOT EXISTS idx_links_short_url ON links (short_url);

-- Индекс для предотвращения повторной генерации дубликатов
CREATE INDEX IF NOT EXISTS idx_links_destination ON links (destination);
