-- 0002_artists.sql
-- Профиль артиста. Привязываем artist-профиль к users (1 к 1)

CREATE TABLE IF NOT EXISTS artists (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  bio TEXT NOT NULL DEFAULT '',
  avatar_path TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_artists_name ON artists (name);