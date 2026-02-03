package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"music-service/models"
)

type playlistRepository struct {
	db *sql.DB
}

func (r *playlistRepository) Create(ctx context.Context, playlist *models.Playlist) error {
	query := `INSERT INTO playlists (name, user_id, tracks) VALUES ($1, $2, $3) RETURNING id`
	return r.db.QueryRowContext(ctx, query, playlist.Name, playlist.UserID, "[]").Scan(&playlist.ID)
}

func (r *playlistRepository) GetByID(ctx context.Context, id int) (*models.Playlist, error) {
	playlist := &models.Playlist{}
	query := `SELECT id, name, user_id, tracks FROM playlists WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&playlist.ID, &playlist.Name, &playlist.UserID, &playlist.Tracks)
	if err != nil {
		return nil, err
	}
	// Парсим JSON tracks
	jsonTracks, _ := json.Marshal(playlist.Tracks)
	err = json.Unmarshal(jsonTracks, &playlist.Tracks)
	return playlist, err
}

func (r *playlistRepository) GetByUserID(ctx context.Context, userID int) ([]models.Playlist, error) {
	query := `SELECT id, name, user_id, tracks FROM playlists WHERE user_id = $1`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var playlists []models.Playlist
	for rows.Next() {
		var p models.Playlist
		var tracksJSON []byte
		if err := rows.Scan(&p.ID, &p.Name, &p.UserID, &tracksJSON); err != nil {
			return nil, err
		}
		json.Unmarshal(tracksJSON, &p.Tracks)
		playlists = append(playlists, p)
	}
	return playlists, nil
}

func (r *playlistRepository) Update(ctx context.Context, playlist *models.Playlist) error {
	tracksJSON, _ := json.Marshal(playlist.Tracks)
	query := `UPDATE playlists SET name = $1, tracks = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, playlist.Name, tracksJSON, playlist.ID)
	return err
}

func (r *playlistRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM playlists WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *playlistRepository) AddTrack(ctx context.Context, playlistID, trackID int) error {
	playlist, err := r.GetByID(ctx, playlistID)
	if err != nil {
		return err
	}
	playlist.Tracks = append(playlist.Tracks, trackID)
	return r.Update(ctx, playlist)
}

func (r *playlistRepository) RemoveTrack(ctx context.Context, playlistID, trackID int) error {
	playlist, err := r.GetByID(ctx, playlistID)
	if err != nil {
		return err
	}
	for i, tid := range playlist.Tracks {
		if tid == trackID {
			playlist.Tracks = append(playlist.Tracks[:i], playlist.Tracks[i+1:]...)
			break
		}
	}
	return r.Update(ctx, playlist)
}

//
