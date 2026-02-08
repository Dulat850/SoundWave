package repositories

import (
	"context"
	"database/sql"
	"music-service/models"
)

type listenRepository struct {
	db *sql.DB
}

func (r *listenRepository) Record(ctx context.Context, userID *int, trackID int) error {
	query := `INSERT INTO listens (user_id, track_id) VALUES ($1, $2)`
	_, err := r.db.ExecContext(ctx, query, userID, trackID)
	return err
}

func (r *listenRepository) IncrementPlayCount(ctx context.Context, trackID int) error {
	query := `UPDATE tracks SET play_count = play_count + 1 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, trackID)
	return err
}

func (r *listenRepository) UserHistory(ctx context.Context, userID int, limit, offset int) ([]models.Listen, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	query := `
SELECT id, user_id, track_id, listened_at
FROM listens
WHERE user_id = $1
ORDER BY listened_at DESC
LIMIT $2 OFFSET $3
`
	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []models.Listen
	for rows.Next() {
		var l models.Listen
		if err := rows.Scan(&l.ID, &l.UserID, &l.TrackID, &l.ListenedAt); err != nil {
			return nil, err
		}
		res = append(res, l)
	}
	return res, rows.Err()
}
