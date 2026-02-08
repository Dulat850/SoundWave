-- 0006_album_tracks.sql
-- Если треки могут быть в альбоме с порядком (track_number)

CREATE TABLE IF NOT EXISTS album_tracks (
  album_id BIGINT NOT NULL REFERENCES albums(id) ON DELETE CASCADE,
  track_id BIGINT NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,
  track_number INT NOT NULL DEFAULT 0,
  PRIMARY KEY (album_id, track_id)
);

CREATE INDEX IF NOT EXISTS idx_album_tracks_album_number ON album_tracks (album_id, track_number);