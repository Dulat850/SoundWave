package repositories

import (
	"context"
	"database/sql"
	"music-service/models"
)

type searchRepository struct {
	db *sql.DB
}

func (r *searchRepository) SearchTracks(ctx context.Context, q string, limit, offset int) ([]models.Track, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	pattern := "%" + q + "%"
	query := `
SELECT id, artist_id, album_id, genre_id, title, duration_seconds, audio_path, cover_path, play_count, created_at
FROM tracks
WHERE title ILIKE $1
ORDER BY play_count DESC, id DESC
LIMIT $2 OFFSET $3
`
	rows, err := r.db.QueryContext(ctx, query, pattern, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []models.Track
	for rows.Next() {
		var t models.Track
		if err := rows.Scan(
			&t.ID,
			&t.ArtistID,
			&t.AlbumID,
			&t.GenreID,
			&t.Title,
			&t.DurationSeconds,
			&t.AudioPath,
			&t.CoverPath,
			&t.PlayCount,
			&t.CreatedAt,
		); err != nil {
			return nil, err
		}
		res = append(res, t)
	}
	return res, rows.Err()
}

func (r *searchRepository) SearchArtists(ctx context.Context, q string, limit, offset int) ([]models.Artist, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	pattern := "%" + q + "%"
	query := `
SELECT id, user_id, name, bio, avatar_path, created_at
FROM artists
WHERE name ILIKE $1
ORDER BY name
LIMIT $2 OFFSET $3
`
	rows, err := r.db.QueryContext(ctx, query, pattern, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []models.Artist
	for rows.Next() {
		var a models.Artist
		if err := rows.Scan(&a.ID, &a.UserID, &a.Name, &a.Bio, &a.AvatarPath, &a.CreatedAt); err != nil {
			return nil, err
		}
		res = append(res, a)
	}
	return res, rows.Err()
}

func (r *searchRepository) SearchGenres(ctx context.Context, q string, limit, offset int) ([]models.Genre, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	pattern := "%" + q + "%"
	query := `
SELECT id, name
FROM genres
WHERE name ILIKE $1
ORDER BY name
LIMIT $2 OFFSET $3
`
	rows, err := r.db.QueryContext(ctx, query, pattern, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []models.Genre
	for rows.Next() {
		var g models.Genre
		if err := rows.Scan(&g.ID, &g.Name); err != nil {
			return nil, err
		}
		res = append(res, g)
	}
	return res, rows.Err()
}
