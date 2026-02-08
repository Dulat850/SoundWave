package repositories

import (
	"context"
	"database/sql"
	"music-service/models"
)

type playListRepository struct {
	db *sql.DB
}

func (r *playListRepository) Create(ctx context.Context, playlist *models.Playlist) error {
	query := `INSERT INTO playlists (name, user_id) VALUES ($1, $2) RETURNING id`
	return r.db.QueryRowContext(ctx, query, playlist.Name, playlist.UserID).Scan(&playlist.ID)
}

func (r *playListRepository) GetByID(ctx context.Context, id int) (*models.Playlist, error) {
	p := &models.Playlist{}
	query := `SELECT id, name, user_id FROM playlists WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&p.ID, &p.Name, &p.UserID)
	return p, err
}

func (r *playListRepository) GetByUserID(ctx context.Context, userID int) ([]models.Playlist, error) {
	query := `SELECT id, name, user_id
			  FROM playlists
			  WHERE user_id = $1
			  ORDER BY id DESC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []models.Playlist
	for rows.Next() {
		var p models.Playlist
		if err := rows.Scan(&p.ID, &p.Name, &p.UserID); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, rows.Err()
}

func (r *playListRepository) Update(ctx context.Context, playlist *models.Playlist) error {
	query := `UPDATE playlists SET name = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, playlist.Name, playlist.ID)
	return err
}

func (r *playListRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM playlists WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *playListRepository) AddTrack(ctx context.Context, playlistID, trackID int) error {
	query := `INSERT INTO playlist_tracks (playlist_id, track_id)
			  VALUES ($1, $2)
			  ON CONFLICT DO NOTHING`
	_, err := r.db.ExecContext(ctx, query, playlistID, trackID)
	return err
}

func (r *playListRepository) RemoveTrack(ctx context.Context, playlistID, trackID int) error {
	query := `DELETE FROM playlist_tracks WHERE playlist_id = $1 AND track_id = $2`
	_, err := r.db.ExecContext(ctx, query, playlistID, trackID)
	return err
}
