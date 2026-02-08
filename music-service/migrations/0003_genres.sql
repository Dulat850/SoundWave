-- 0003_genres.sql
-- Жанры

CREATE TABLE IF NOT EXISTS genres (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL UNIQUE
);

CREATE INDEX IF NOT EXISTS idx_genres_name ON genres (name);