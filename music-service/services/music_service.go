package services

import (
	"context"
	"database/sql"
	"errors"

	"music-service/models"
	"music-service/repositories"
)

var ErrNotFound = errors.New("not found")

type MusicService interface {
	ListArtists(ctx context.Context, limit, offset int) ([]models.Artist, error)
	GetArtistByID(ctx context.Context, id int) (*models.Artist, error)
	SearchTracks(ctx context.Context, query string) ([]models.Track, error)
}

type musicService struct {
	repo repositories.MusicRepository
}

func NewMusicService(repo repositories.MusicRepository) MusicService {
	return &musicService{repo: repo}
}

func (s *musicService) ListArtists(ctx context.Context, limit, offset int) ([]models.Artist, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.ListArtists(ctx, limit, offset)
}

func (s *musicService) GetArtistByID(ctx context.Context, id int) (*models.Artist, error) {
	artist, err := s.repo.GetArtistByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return artist, nil
}

func (s *musicService) SearchTracks(ctx context.Context, query string) ([]models.Track, error) {
	return s.repo.SearchTracks(ctx, query)
}
