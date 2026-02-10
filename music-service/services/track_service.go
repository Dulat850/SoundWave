package services

import (
	"context"
	"database/sql"
	"errors"
	"music-service/models"
	"music-service/repositories"
)

type TrackService interface {
	Create(ctx context.Context, t *models.Track) (*models.Track, error)
	GetByID(ctx context.Context, id int) (*models.Track, error)
	Top(ctx context.Context, limit int) ([]models.Track, error)

	List(ctx context.Context, limit, offset int, sort string) ([]models.Track, error)
	ListByArtistID(ctx context.Context, artistID int, limit, offset int, sort string) ([]models.Track, error)
	ListByAlbumID(ctx context.Context, albumID int, limit, offset int, sort string) ([]models.Track, error)
}
type trackService struct {
	repo repositories.TrackRepository
}

func NewTrackService(repo repositories.TrackRepository) TrackService {
	return &trackService{repo: repo}
}

func (s *trackService) Create(ctx context.Context, t *models.Track) (*models.Track, error) {
	if t == nil {
		return nil, errors.New("track is nil")
	}
	if t.ArtistID <= 0 {
		return nil, errors.New("artist_id is required")
	}
	if t.Title == "" {
		return nil, errors.New("title is required")
	}
	if t.AudioPath == "" {
		return nil, errors.New("audio_path is required")
	}

	if err := s.repo.Create(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *trackService) GetByID(ctx context.Context, id int) (*models.Track, error) {
	t, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return t, nil
}

func (s *trackService) Top(ctx context.Context, limit int) ([]models.Track, error) {
	return s.repo.Top(ctx, limit)
}

func (s *trackService) List(ctx context.Context, limit, offset int, sort string) ([]models.Track, error) {
	return s.repo.List(ctx, limit, offset, normalizeTrackSort(sort))
}

func (s *trackService) ListByArtistID(ctx context.Context, artistID int, limit, offset int, sort string) ([]models.Track, error) {
	if artistID <= 0 {
		return nil, errors.New("invalid artist_id")
	}
	return s.repo.ListByArtistID(ctx, artistID, limit, offset, normalizeTrackSort(sort))
}

func (s *trackService) ListByAlbumID(ctx context.Context, albumID int, limit, offset int, sort string) ([]models.Track, error) {
	if albumID <= 0 {
		return nil, errors.New("invalid album_id")
	}
	return s.repo.ListByAlbumID(ctx, albumID, limit, offset, normalizeTrackSort(sort))
}

func normalizeTrackSort(sort string) string {
	switch sort {
	case "popular", "old", "new", "":
		if sort == "" {
			return "new"
		}
		return sort
	default:
		return "new"
	}
}
