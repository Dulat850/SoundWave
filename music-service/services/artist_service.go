package services

import (
	"context"
	"database/sql"
	"errors"
	"music-service/models"
	"music-service/repositories"
)

// Ошибки
var (
	ErrNotFoundArt = errors.New("artist not found")
	ErrConflictArt = errors.New("artist already exists")
)

// Интерфейс сервиса артистов
type ArtistService interface {
	Me(ctx context.Context, userID int) (*models.Artist, error)
	UpsertMe(ctx context.Context, userID int, name string, bio string, avatarPath *string) (*models.Artist, error)
	GetAll(ctx context.Context) ([]models.Artist, error) // <- исправлено
}

type artistService struct {
	repo repositories.ArtistRepository
}

func NewArtistService(repo repositories.ArtistRepository) ArtistService {
	return &artistService{repo: repo}
}

func (s *artistService) Me(ctx context.Context, userID int) (*models.Artist, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user_id")
	}

	a, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFoundArt
		}
		return nil, err
	}
	return a, nil
}

func (s *artistService) UpsertMe(ctx context.Context, userID int, name string, bio string, avatarPath *string) (*models.Artist, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user_id")
	}
	if name == "" {
		return nil, errors.New("name is required")
	}

	_, err := s.repo.GetByUserID(ctx, userID)
	if err == nil {
		return nil, ErrConflictArt
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	a := &models.Artist{
		UserID:     userID,
		Name:       name,
		Bio:        bio,
		AvatarPath: avatarPath,
	}

	if err := s.repo.Create(ctx, a); err != nil {
		return nil, err
	}
	return a, nil
}

// Реализация метода GetAll
func (s *artistService) GetAll(ctx context.Context) ([]models.Artist, error) {
	return s.repo.GetAll(ctx)
}
