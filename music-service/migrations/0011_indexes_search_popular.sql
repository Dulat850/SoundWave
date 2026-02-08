-- 0011_indexes_search_popular.sql
-- Индексы для поиска (artist name, genre, track title) и "популярного"
-- Используем pg_trgm для быстрого ILIKE.

CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX IF NOT EXISTS gin_trgm_artists_name ON artists USING gin (name gin_trgm_ops);
CREATE INDEX IF NOT EXISTS gin_trgm_genres_name ON genres USING gin (name gin_trgm_ops);
CREATE INDEX IF NOT EXISTS gin_trgm_tracks_title ON tracks USING gin (title gin_trgm_ops);

-- Дополнительно: если часто ищешь по имени альбома
CREATE INDEX IF NOT EXISTS gin_trgm_albums_title ON albums USING gin (title gin_trgm_ops);