package repositories

import (
	"context"
	"database/sql"
	"music-service/models"
)

type playListRepository struct {
	db *sql.DB
}

func (r *playListRepository) CreateArtist(ctx context.Context, artist *models.Artist) error {
	query := `INSERT INTO artists (name, description) VALUES ($1, $2) RETURNING id`
	return r.db.QueryRowContext(ctx, query, artist.Name, artist.Description).Scan(&artist.ID)
}

func (r *playListRepository) GetArtistByID(ctx context.Context, id int) (*models.Artist, error) {
	artist := &models.Artist{}
	query := `SELECT id, name, description FROM artists WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, id).
		Scan(&artist.ID, &artist.Name, &artist.Description)

	return artist, err
}

func (r *playListRepository) UpdateArtist(ctx context.Context, artist *models.Artist) error {
	query := `UPDATE artists SET name=$1, description=$2 WHERE id=$3`
	_, err := r.db.ExecContext(ctx, query, artist.Name, artist.Description, artist.ID)
	return err
}

func (r *playListRepository) DeleteArtist(ctx context.Context, id int) error {
	query := `DELETE FROM artists WHERE id=$1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *playListRepository) ListArtists(ctx context.Context, limit, offset int) ([]models.Artist, error) {

	query := `SELECT id, name, description 
			  FROM artists 
			  ORDER BY name 
			  LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var artists []models.Artist

	for rows.Next() {
		var a models.Artist

		if err := rows.Scan(&a.ID, &a.Name, &a.Description); err != nil {
			return nil, err
		}

		artists = append(artists, a)
	}

	return artists, rows.Err()
}

// ... existing code ...

func (r *playListRepository) CreateAlbum(ctx context.Context, album *models.Album) error {
	query := `INSERT INTO albums (title, artist_id, cover_url) 
			  VALUES ($1, $2, $3) RETURNING id`

	return r.db.QueryRowContext(
		ctx,
		query,
		album.Title,
		album.ArtistID,
		album.CoverURL,
	).Scan(&album.ID)
}

// ... existing code ...

func (r *playListRepository) CreateTrack(ctx context.Context, track *models.Track) error {

	query := `INSERT INTO tracks 
			  (title, artist_id, album_id, duration, audio_url) 
			  VALUES ($1, $2, $3, $4, $5) 
			  RETURNING id`

	return r.db.QueryRowContext(
		ctx,
		query,
		track.Title,
		track.ArtistID,
		track.AlbumID,
		track.Duration,
		track.AudioURL,
	).Scan(&track.ID)
}

// ... existing code ...

func (r *playListRepository) ListAlbumsByArtist(ctx context.Context, artistID int) ([]models.Album, error) {
	query := `SELECT id, title, artist_id, cover_url 
			  FROM albums 
			  WHERE artist_id=$1 
			  ORDER BY title`

	rows, err := r.db.QueryContext(ctx, query, artistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var albums []models.Album
	for rows.Next() {
		var a models.Album
		if err := rows.Scan(&a.ID, &a.Title, &a.ArtistID, &a.CoverURL); err != nil {
			return nil, err
		}
		albums = append(albums, a)
	}
	return albums, rows.Err()
}

func (r *playListRepository) ListTracksByArtist(ctx context.Context, artistID int) ([]models.Track, error) {
	query := `SELECT id, title, artist_id, album_id, duration, audio_url 
			  FROM tracks 
			  WHERE artist_id=$1 
			  ORDER BY title`

	rows, err := r.db.QueryContext(ctx, query, artistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tracks []models.Track
	for rows.Next() {
		var t models.Track
		if err := rows.Scan(&t.ID, &t.Title, &t.ArtistID, &t.AlbumID, &t.Duration, &t.AudioURL); err != nil {
			return nil, err
		}
		tracks = append(tracks, t)
	}
	return tracks, rows.Err()
}

func (r *playListRepository) ListTracksByAlbum(ctx context.Context, albumID int) ([]models.Track, error) {
	query := `SELECT id, title, artist_id, album_id, duration, audio_url 
			  FROM tracks 
			  WHERE album_id=$1 
			  ORDER BY title`

	rows, err := r.db.QueryContext(ctx, query, albumID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tracks []models.Track
	for rows.Next() {
		var t models.Track
		if err := rows.Scan(&t.ID, &t.Title, &t.ArtistID, &t.AlbumID, &t.Duration, &t.AudioURL); err != nil {
			return nil, err
		}
		tracks = append(tracks, t)
	}
	return tracks, rows.Err()
}

func (r *playListRepository) SearchTracks(ctx context.Context, search string) ([]models.Track, error) {
	q := "%" + search + "%"

	query := `SELECT id, title, artist_id, album_id, duration, audio_url
			  FROM tracks
			  WHERE title ILIKE $1
			  ORDER BY title
			  LIMIT 50`

	rows, err := r.db.QueryContext(ctx, query, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tracks []models.Track
	for rows.Next() {
		var t models.Track
		if err := rows.Scan(&t.ID, &t.Title, &t.ArtistID, &t.AlbumID, &t.Duration, &t.AudioURL); err != nil {
			return nil, err
		}
		tracks = append(tracks, t)
	}
	return tracks, rows.Err()
}
