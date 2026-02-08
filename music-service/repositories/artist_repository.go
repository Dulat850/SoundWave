package repositories

import (
	"context"
	"database/sql"
	"music-service/models"
)

type artistRepository struct {
	db *sql.DB
}

func (r *artistRepository) Create(ctx context.Context, a *models.Artist) error {
	query := `INSERT INTO artists (user_id, name, bio, avatar_path) VALUES ($1, $2, $3, $4) RETURNING id`
	return r.db.QueryRowContext(ctx, query, a.UserID, a.Name, a.Bio, a.AvatarPath).Scan(&a.ID)
}

func (r *artistRepository) GetByID(ctx context.Context, id int) (*models.Artist, error) {
	a := &models.Artist{}
	query := `SELECT id, user_id, name, bio, avatar_path FROM artists WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&a.ID, &a.UserID, &a.Name, &a.Bio, &a.AvatarPath)
	return a, err
}

func (r *artistRepository) GetByUserID(ctx context.Context, userID int) (*models.Artist, error) {
	a := &models.Artist{}
	query := `SELECT id, user_id, name, bio, avatar_path FROM artists WHERE user_id = $1`
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&a.ID, &a.UserID, &a.Name, &a.Bio, &a.AvatarPath)
	return a, err
}
