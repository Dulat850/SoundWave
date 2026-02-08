package repositories

import (
	"context"
	"database/sql"
	"music-service/models"
)

type genreRepository struct {
	db *sql.DB
}

func (r *genreRepository) Create(ctx context.Context, g *models.Genre) error {
	query := `INSERT INTO genres (name) VALUES ($1) RETURNING id`
	return r.db.QueryRowContext(ctx, query, g.Name).Scan(&g.ID)
}

func (r *genreRepository) GetByID(ctx context.Context, id int) (*models.Genre, error) {
	g := &models.Genre{}
	query := `SELECT id, name FROM genres WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&g.ID, &g.Name)
	return g, err
}

func (r *genreRepository) GetByName(ctx context.Context, name string) (*models.Genre, error) {
	g := &models.Genre{}
	query := `SELECT id, name FROM genres WHERE name = $1`
	err := r.db.QueryRowContext(ctx, query, name).Scan(&g.ID, &g.Name)
	return g, err
}

func (r *genreRepository) List(ctx context.Context, limit, offset int) ([]models.Genre, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	query := `SELECT id, name FROM genres ORDER BY name LIMIT $1 OFFSET $2`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
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
