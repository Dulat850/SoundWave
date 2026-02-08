package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"music-service/models"
)

type trackRepository struct {
	db *sql.DB
}

func (r *trackRepository) Create(ctx context.Context, t *models.Track) error {
	query := `
INSERT INTO tracks (artist_id, album_id, genre_id, title, duration_seconds, audio_path, cover_path)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, play_count, created_at
`
	return r.db.QueryRowContext(
		ctx,
		query,
		t.ArtistID,
		t.AlbumID,
		t.GenreID,
		t.Title,
		t.DurationSeconds,
		t.AudioPath,
		t.CoverPath,
	).Scan(&t.ID, &t.PlayCount, &t.CreatedAt)
}

func (r *trackRepository) GetByID(ctx context.Context, id int) (*models.Track, error) {
	t := &models.Track{}
	query := `
SELECT id, artist_id, album_id, genre_id, title, duration_seconds, audio_path, cover_path, play_count, created_at
FROM tracks
WHERE id = $1
`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
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
	)
	return t, err
}

func (r *trackRepository) Top(ctx context.Context, limit int) ([]models.Track, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	query := `
SELECT id, artist_id, album_id, genre_id, title, duration_seconds, audio_path, cover_path, play_count, created_at
FROM tracks
ORDER BY play_count DESC, id DESC
LIMIT $1
`
	rows, err := r.db.QueryContext(ctx, query, limit)
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

func (r *trackRepository) List(ctx context.Context, limit, offset int, sort string) ([]models.Track, error) {
	return r.list(ctx, 0, "all", limit, offset, sort)
}

func (r *trackRepository) ListByArtistID(ctx context.Context, artistID int, limit, offset int, sort string) ([]models.Track, error) {
	return r.list(ctx, artistID, "artist", limit, offset, sort)
}

func (r *trackRepository) ListByAlbumID(ctx context.Context, albumID int, limit, offset int, sort string) ([]models.Track, error) {
	return r.list(ctx, albumID, "album", limit, offset, sort)
}

func (r *trackRepository) list(ctx context.Context, id int, mode string, limit, offset int, sort string) ([]models.Track, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	orderBy := orderByForTracks(sort)

	where := ""
	args := []any{}
	switch mode {
	case "artist":
		where = "WHERE artist_id = $1"
		args = append(args, id, limit, offset)
	case "album":
		where = "WHERE album_id = $1"
		args = append(args, id, limit, offset)
	default:
		args = append(args, limit, offset)
	}

	query := fmt.Sprintf(`
SELECT id, artist_id, album_id, genre_id, title, duration_seconds, audio_path, cover_path, play_count, created_at
FROM tracks
%s
%s
LIMIT $%d OFFSET $%d
`, where, orderBy, len(args)-1, len(args))

	rows, err := r.db.QueryContext(ctx, query, args...)
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

func orderByForTracks(sort string) string {
	switch sort {
	case "popular":
		return "ORDER BY play_count DESC, id DESC"
	case "old":
		return "ORDER BY created_at ASC, id ASC"
	default: // "new"
		return "ORDER BY created_at DESC, id DESC"
	}
}
