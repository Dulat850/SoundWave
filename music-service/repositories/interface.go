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

type ArtistRepository interface {
	Create(ctx context.Context, a *models.Artist) error
	GetByID(ctx context.Context, id int) (*models.Artist, error)
	GetByUserID(ctx context.Context, userID int) (*models.Artist, error)
}

func NewArtistRepository(db *sql.DB) ArtistRepository {
	return &artistRepository{db: db}
}

type TrackRepository interface {
	Create(ctx context.Context, t *models.Track) error
	GetByID(ctx context.Context, id int) (*models.Track, error)
	Top(ctx context.Context, limit int) ([]models.Track, error)

	List(ctx context.Context, limit, offset int, sort string) ([]models.Track, error)
	ListByArtistID(ctx context.Context, artistID int, limit, offset int, sort string) ([]models.Track, error)
	ListByAlbumID(ctx context.Context, albumID int, limit, offset int, sort string) ([]models.Track, error)
}

func NewTrackRepository(db *sql.DB) TrackRepository {
	return &trackRepository{db: db}
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
	return &playListRepository{db: db}
}

// UserRepository interface
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id int) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, limit, offset int) ([]models.User, error)
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

type ListenRepository interface {
	Record(ctx context.Context, userID *int, trackID int) error
	UserHistory(ctx context.Context, userID int, limit, offset int) ([]models.Listen, error)
	IncrementPlayCount(ctx context.Context, trackID int) error
}

type SearchRepository interface {
	SearchTracks(ctx context.Context, q string, limit, offset int) ([]models.Track, error)
	SearchArtists(ctx context.Context, q string, limit, offset int) ([]models.Artist, error)
	SearchGenres(ctx context.Context, q string, limit, offset int) ([]models.Genre, error)
}

func NewSearchRepository(db *sql.DB) SearchRepository {
	return &searchRepository{db: db}
}

type LikeRepository interface {
	Like(ctx context.Context, userID, trackID int) error
	Unlike(ctx context.Context, userID, trackID int) error
	ListLikedTrackIDs(ctx context.Context, userID int, limit, offset int) ([]int, error)
}

type GenreRepository interface {
	Create(ctx context.Context, g *models.Genre) error
	GetByID(ctx context.Context, id int) (*models.Genre, error)
	GetByName(ctx context.Context, name string) (*models.Genre, error)
	List(ctx context.Context, limit, offset int) ([]models.Genre, error)
}

func NewGenreRepository(db *sql.DB) GenreRepository {
	return &genreRepository{db: db}
}

type AlbumRepository interface {
	Create(ctx context.Context, a *models.Album) error
	GetByID(ctx context.Context, id int) (*models.Album, error)
	ListByArtistID(ctx context.Context, artistID int, limit, offset int) ([]models.Album, error)
}

func NewAlbumRepository(db *sql.DB) AlbumRepository {
	return &albumRepository{db: db}
}

func NewLikeRepository(db *sql.DB) LikeRepository {
	return &likeRepository{db: db}
}

func NewListenRepository(db *sql.DB) ListenRepository {
	return &listenRepository{db: db}
}
