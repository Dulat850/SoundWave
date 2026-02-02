package repositories

import (
	"context"
	"database/sql"
	"music-service/models"
)

type musicRepository struct {
	db *sql.DB
}

// ===== ARTISTS =====
func (r *musicRepository) CreateArtist(ctx context.Context, artist *models.Artist) error {
	query := `INSERT INTO artists (name, description) VALUES ($1, $2) RETURNING id`
	return r.db.QueryRowContext(ctx, query, artist.Name, artist.Description).Scan(&artist.ID)
}

func (r *musicRepository) GetArtistByID(ctx context.Context, id int) (*models.Artist, error) {
	artist := &models.Artist{}
	query := `SELECT id, name, description FROM artists WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&artist.ID, &artist.Name, &artist.Description)
	return artist, err
}

func (r *musicRepository) UpdateArtist(ctx context.Context, artist *models.Artist) error {
	query := `UPDATE artists SET name = $1, description = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, artist.Name, artist.Description, artist.ID)
	return err
}

func (r *musicRepository) DeleteArtist(ctx context.Context, id int) error {
	query := `DELETE FROM artists WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *musicRepository) ListArtists(ctx context.Context, limit, offset int) ([]models.Artist, error) {
	query := `SELECT id, name, description FROM artists ORDER BY name LIMIT $1 OFFSET $2`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var artists []models.Artist
	for rows.Next() {
		var a models.Artist
		rows.Scan(&a.ID, &a.Name, &a.Description)
		artists = append(artists, a)
	}
	return artists, nil
}

// ===== ALBUMS =====
func (r *musicRepository) CreateAlbum(ctx context.Context, album *models.Album) error {
	query := `INSERT INTO albums (title, artist_id, cover_url) VALUES ($1, $2, $3) RETURNING id`
	return r.db.QueryRowContext(ctx, query, album.Title, album.ArtistID, album.CoverURL).Scan(&album.ID)
}

func (r *musicRepository) GetAlbumByID(ctx context.Context, id int) (*models.Album, error) {
	album := &models.Album{}
	query := `SELECT id, title, artist_id, cover_url FROM albums WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&album.ID, &album.Title, &album.ArtistID, &album.CoverURL)
	return album, err
}

func (r *musicRepository) ListAlbumsByArtist(ctx context.Context, artistID int) ([]models.Album, error) {
	query := `SELECT id, title, artist_id, cover_url FROM albums WHERE artist_id = $1 ORDER BY title`
	rows, err := r.db.QueryContext(ctx, query, artistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var albums []models.Album
	for rows.Next() {
		var a models.Album
		rows.Scan(&a.ID, &a.Title, &a.ArtistID, &a.CoverURL)
		albums = append(albums, a)
	}
	return albums, nil
}

// ===== TRACKS =====
func (r *musicRepository) CreateTrack(ctx context.Context, track *models.Track) error {
	query := `INSERT INTO tracks (title, artist_id, album_id, duration, audio_url) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return r.db.QueryRowContext(ctx, query, track.Title, track.ArtistID, track.AlbumID, track.Duration, track.AudioURL).Scan(&track.ID)
}

func (r *musicRepository) GetTrackByID(ctx context.Context, id int) (*models.Track, error) {
	track := &models.Track{}
	query := `SELECT id, title, artist_id, album_id, duration, audio_url FROM tracks WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&track.ID, &track.Title, &track.ArtistID, &track.AlbumID, &track.Duration, &track.AudioURL)
	return track, err
}

func (r *musicRepository) ListTracksByArtist(ctx context.Context, artistID int) ([]models.Track, error) {
	query := `SELECT id, title, artist_id, album_id, duration, audio_url FROM tracks WHERE artist_id = $1 ORDER BY title`
	rows, err := r.db.QueryContext(ctx, query, artistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tracks []models.Track
	for rows.Next() {
		var t models.Track
		rows.Scan(&t.ID, &t.Title, &t.ArtistID, &t.AlbumID, &t.Duration, &t.AudioURL)
		tracks = append(tracks, t)
	}
	return tracks, nil
}

func (r *musicRepository) ListTracksByAlbum(ctx context.Context, albumID int) ([]models.Track, error) {
	query := `SELECT id, title, artist_id, album_id, duration, audio_url FROM tracks WHERE album_id = $1 ORDER BY title`
	rows, err := r.db.QueryContext(ctx, query, albumID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tracks []models.Track
	for rows.Next() {
		var t models.Track
		rows.Scan(&t.ID, &t.Title, &t.ArtistID, &t.AlbumID, &t.Duration, &t.AudioURL)
		tracks = append(tracks, t)
	}
	return tracks, nil
}

func (r *musicRepository) SearchTracks(ctx context.Context, query string) ([]models.Track, error) {
	q := `%` + query + `%`
	sqlQuery := `SELECT id, title, artist_id, album_id, duration, audio_url 
                 FROM tracks WHERE title ILIKE $1 OR audio_url ILIKE $1 ORDER BY title LIMIT 50`
	rows, err := r.db.QueryContext(ctx, sqlQuery, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tracks []models.Track
	for rows.Next() {
		var t models.Track
		rows.Scan(&t.ID, &t.Title, &t.ArtistID, &t.AlbumID, &t.Duration, &t.AudioURL)
		tracks = append(tracks, t)
	}
	return tracks, nil
}
