package services

import (
	"context"
	"database/sql"
	"errors"
	"music-service/models"
	"music-service/repositories"
)

type GenreService interface {
	List(ctx context.Context, limit, offset int) ([]models.Genre, error)
	Create(ctx context.Context, name string) (*models.Genre, error)
}

type genreService struct {
	repo repositories.GenreRepository
}

func NewGenreService(repo repositories.GenreRepository) GenreService {
	return &genreService{repo: repo}
}

func (s *genreService) List(ctx context.Context, limit, offset int) ([]models.Genre, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *genreService) Create(ctx context.Context, name string) (*models.Genre, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}

	// уникальность обеспечит БД (UNIQUE), но можно проверить заранее
	if _, err := s.repo.GetByName(ctx, name); err == nil {
		return nil, ErrConflict
	} else if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	g := &models.Genre{Name: name}
	if err := s.repo.Create(ctx, g); err != nil {
		return nil, err
	}
	return g, nil
}
