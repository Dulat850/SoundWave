package repositories

import (
	"context"
	"database/sql"
	"music-service/models"
)

type albumRepository struct {
	db *sql.DB
}

func (r *albumRepository) Create(ctx context.Context, a *models.Album) error {
	query := `
INSERT INTO albums (artist_id, title, cover_path, released_at)
VALUES ($1, $2, $3, $4)
RETURNING id, created_at
`
	return r.db.QueryRowContext(ctx, query, a.ArtistID, a.Title, a.CoverPath, a.ReleasedAt).
		Scan(&a.ID, &a.CreatedAt)
}

func (r *albumRepository) GetByID(ctx context.Context, id int) (*models.Album, error) {
	a := &models.Album{}
	query := `
SELECT id, artist_id, title, cover_path, released_at, created_at
FROM albums
WHERE id = $1
`
	err := r.db.QueryRowContext(ctx, query, id).
		Scan(&a.ID, &a.ArtistID, &a.Title, &a.CoverPath, &a.ReleasedAt, &a.CreatedAt)
	return a, err
}

func (r *albumRepository) ListByArtistID(ctx context.Context, artistID int, limit, offset int) ([]models.Album, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	query := `
SELECT id, artist_id, title, cover_path, released_at, created_at
FROM albums
WHERE artist_id = $1
ORDER BY created_at DESC, id DESC
LIMIT $2 OFFSET $3
`
	rows, err := r.db.QueryContext(ctx, query, artistID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []models.Album
	for rows.Next() {
		var a models.Album
		if err := rows.Scan(&a.ID, &a.ArtistID, &a.Title, &a.CoverPath, &a.ReleasedAt, &a.CreatedAt); err != nil {
			return nil, err
		}
		res = append(res, a)
	}
	return res, rows.Err()
}
