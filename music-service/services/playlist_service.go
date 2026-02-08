package services

import (
	"context"
	"database/sql"
	"errors"
	"music-service/models"
	"music-service/repositories"
)

var (
	ErrForbidden = errors.New("forbidden")
)

type PlaylistService interface {
	Create(ctx context.Context, userID int, name string) (*models.Playlist, error)
	GetByID(ctx context.Context, id int) (*models.Playlist, error)
	GetByUserID(ctx context.Context, userID int) ([]models.Playlist, error)
	Rename(ctx context.Context, userID, playlistID int, name string) error
	Delete(ctx context.Context, userID, playlistID int) error
	AddTrack(ctx context.Context, userID, playlistID, trackID int) error
	RemoveTrack(ctx context.Context, userID, playlistID, trackID int) error
}

type playlistService struct {
	repo repositories.PlaylistRepository
}

func NewPlaylistService(repo repositories.PlaylistRepository) PlaylistService {
	return &playlistService{repo: repo}
}

func (s *playlistService) Create(ctx context.Context, userID int, name string) (*models.Playlist, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user_id")
	}
	if name == "" {
		return nil, errors.New("name is required")
	}

	p := &models.Playlist{
		Name:   name,
		UserID: userID,
	}

	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *playlistService) GetByID(ctx context.Context, id int) (*models.Playlist, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return p, nil
}

func (s *playlistService) GetByUserID(ctx context.Context, userID int) ([]models.Playlist, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user_id")
	}
	return s.repo.GetByUserID(ctx, userID)
}

func (s *playlistService) Rename(ctx context.Context, userID, playlistID int, name string) error {
	if name == "" {
		return errors.New("name is required")
	}

	p, err := s.GetByID(ctx, playlistID)
	if err != nil {
		return err
	}
	if p.UserID != userID {
		return ErrForbidden
	}

	p.Name = name
	return s.repo.Update(ctx, p)
}

func (s *playlistService) Delete(ctx context.Context, userID, playlistID int) error {
	p, err := s.GetByID(ctx, playlistID)
	if err != nil {
		return err
	}
	if p.UserID != userID {
		return ErrForbidden
	}

	return s.repo.Delete(ctx, playlistID)
}

func (s *playlistService) AddTrack(ctx context.Context, userID, playlistID, trackID int) error {
	if trackID <= 0 {
		return errors.New("invalid track_id")
	}

	p, err := s.GetByID(ctx, playlistID)
	if err != nil {
		return err
	}
	if p.UserID != userID {
		return ErrForbidden
	}

	return s.repo.AddTrack(ctx, playlistID, trackID)
}

func (s *playlistService) RemoveTrack(ctx context.Context, userID, playlistID, trackID int) error {
	if trackID <= 0 {
		return errors.New("invalid track_id")
	}

	p, err := s.GetByID(ctx, playlistID)
	if err != nil {
		return err
	}
	if p.UserID != userID {
		return ErrForbidden
	}

	return s.repo.RemoveTrack(ctx, playlistID, trackID)
}
