package repositories

import (
	"context"
	"database/sql"
)

type likeRepository struct {
	db *sql.DB
}

func (r *likeRepository) Like(ctx context.Context, userID, trackID int) error {
	query := `INSERT INTO likes (user_id, track_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := r.db.ExecContext(ctx, query, userID, trackID)
	return err
}

func (r *likeRepository) Unlike(ctx context.Context, userID, trackID int) error {
	query := `DELETE FROM likes WHERE user_id = $1 AND track_id = $2`
	_, err := r.db.ExecContext(ctx, query, userID, trackID)
	return err
}

func (r *likeRepository) ListLikedTrackIDs(ctx context.Context, userID int, limit, offset int) ([]int, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	query := `
SELECT track_id
FROM likes
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3
`
	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}
