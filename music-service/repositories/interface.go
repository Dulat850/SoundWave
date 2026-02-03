package repositories

import (
	"context"
	"database/sql"
	"music-service/models"
)

type MusicRepository interface {
	// Artists
	CreateArtist(ctx context.Context, artist *models.Artist) error
	GetArtistByID(ctx context.Context, id int) (*models.Artist, error)
	UpdateArtist(ctx context.Context, artist *models.Artist) error
	DeleteArtist(ctx context.Context, id int) error
	ListArtists(ctx context.Context, limit, offset int) ([]models.Artist, error)

	// Albums
	CreateAlbum(ctx context.Context, album *models.Album) error
	GetAlbumByID(ctx context.Context, id int) (*models.Album, error)
	ListAlbumsByArtist(ctx context.Context, artistID int) ([]models.Album, error)

	// Tracks
	CreateTrack(ctx context.Context, track *models.Track) error
	GetTrackByID(ctx context.Context, id int) (*models.Track, error)
	ListTracksByArtist(ctx context.Context, artistID int) ([]models.Track, error)
	ListTracksByAlbum(ctx context.Context, albumID int) ([]models.Track, error)
	SearchTracks(ctx context.Context, query string) ([]models.Track, error) // поиск по названию
}

func NewMusicRepository(db *sql.DB) MusicRepository {
	return &musicRepository{db: db}
}

type PlaylistRepository interface {
	Create(ctx context.Context, playlist *models.Playlist) error
	GetByID(ctx context.Context, id int) (*models.Playlist, error)
	GetByUserID(ctx context.Context, userID int) ([]models.Playlist, error)
	Update(ctx context.Context, playlist *models.Playlist) error
	Delete(ctx context.Context, id int) error
	AddTrack(ctx context.Context, playlistID, trackID int) error
	RemoveTrack(ctx context.Context, playlistID, trackID int) error
}

func NewPlaylistRepository(db *sql.DB) PlaylistRepository {
	return &playlistRepository{db: db}
}

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id int) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, limit, offset int) ([]models.User, error)
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

//
