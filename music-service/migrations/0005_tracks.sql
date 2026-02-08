-- 0005_tracks.sql
-- Треки + локальные пути к аудио/обложке + счетчик прослушиваний для "популярного"

CREATE TABLE IF NOT EXISTS tracks (
  id BIGSERIAL PRIMARY KEY,
  artist_id BIGINT NOT NULL REFERENCES artists(id) ON DELETE CASCADE,
  album_id BIGINT REFERENCES albums(id) ON DELETE SET NULL,
  genre_id BIGINT REFERENCES genres(id) ON DELETE SET NULL,

  title TEXT NOT NULL,
  duration_seconds INT NOT NULL DEFAULT 0,

  audio_path TEXT NOT NULL,      -- например: storage/audio/<file>
  cover_path TEXT,               -- например: storage/covers/<file>

  play_count BIGINT NOT NULL DEFAULT 0,  -- можно увеличивать при прослушивании
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_tracks_artist_id ON tracks (artist_id);
CREATE INDEX IF NOT EXISTS idx_tracks_album_id ON tracks (album_id);
CREATE INDEX IF NOT EXISTS idx_tracks_genre_id ON tracks (genre_id);
CREATE INDEX IF NOT EXISTS idx_tracks_title ON tracks (title);
CREATE INDEX IF NOT EXISTS idx_tracks_play_count ON tracks (play_count DESC);