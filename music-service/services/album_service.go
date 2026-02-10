package services

import (
	"context"
	"database/sql"
	"errors"
	"music-service/models"
	"music-service/repositories"
)

type AlbumService interface {
	Create(ctx context.Context, artistID int, title string, coverPath *string, releasedAt *string) (*models.Album, error)
	GetByIDAlbum(ctx context.Context, id int) (*models.Album, error)
	ListByArtistID(ctx context.Context, artistID int, limit, offset int) ([]models.Album, error)
}

type albumService struct {
	repo repositories.AlbumRepository
}

func NewAlbumService(repo repositories.AlbumRepository) AlbumService {
	return &albumService{repo: repo}
}

func (s *albumService) Create(ctx context.Context, artistID int, title string, coverPath *string, releasedAt *string) (*models.Album, error) {
	if artistID <= 0 {
		return nil, errors.New("invalid artist_id")
	}
	if title == "" {
		return nil, errors.New("title is required")
	}

	a := &models.Album{
		ArtistID:   artistID,
		Title:      title,
		CoverPath:  coverPath,
		ReleasedAt: releasedAt,
	}

	if err := s.repo.Create(ctx, a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *albumService) GetByIDAlbum(ctx context.Context, id int) (*models.Album, error) {
	a, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return a, nil
}

func (s *albumService) ListByArtistID(ctx context.Context, artistID int, limit, offset int) ([]models.Album, error) {
	if artistID <= 0 {
		return nil, errors.New("invalid artist_id")
	}
	return s.repo.ListByArtistID(ctx, artistID, limit, offset)
}
