-- 0010_listens_history.sql
-- История прослушиваний (события)
-- В коде можно: вставить событие + увеличить tracks.play_count

CREATE TABLE IF NOT EXISTS listens (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
  track_id BIGINT NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,
  listened_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_listens_track_id_time ON listens (track_id, listened_at DESC);
CREATE INDEX IF NOT EXISTS idx_listens_user_id_time ON listens (user_id, listened_at DESC);