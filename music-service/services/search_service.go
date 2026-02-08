package services

import (
	"context"
	"errors"
	"music-service/models"
	"music-service/repositories"
)

type SearchService interface {
	All(ctx context.Context, q string, limit, offset int) (tracks []models.Track, artists []models.Artist, genres []models.Genre, err error)
	Tracks(ctx context.Context, q string, limit, offset int) ([]models.Track, error)
	Artists(ctx context.Context, q string, limit, offset int) ([]models.Artist, error)
	Genres(ctx context.Context, q string, limit, offset int) ([]models.Genre, error)
}

type searchService struct {
	repo repositories.SearchRepository
}

func NewSearchService(repo repositories.SearchRepository) SearchService {
	return &searchService{repo: repo}
}

func (s *searchService) All(ctx context.Context, q string, limit, offset int) ([]models.Track, []models.Artist, []models.Genre, error) {
	if q == "" {
		return nil, nil, nil, errors.New("q is required")
	}

	tracks, err := s.repo.SearchTracks(ctx, q, limit, offset)
	if err != nil {
		return nil, nil, nil, err
	}
	artists, err := s.repo.SearchArtists(ctx, q, limit, offset)
	if err != nil {
		return nil, nil, nil, err
	}
	genres, err := s.repo.SearchGenres(ctx, q, limit, offset)
	if err != nil {
		return nil, nil, nil, err
	}

	return tracks, artists, genres, nil
}

func (s *searchService) Tracks(ctx context.Context, q string, limit, offset int) ([]models.Track, error) {
	if q == "" {
		return nil, errors.New("q is required")
	}
	return s.repo.SearchTracks(ctx, q, limit, offset)
}

func (s *searchService) Artists(ctx context.Context, q string, limit, offset int) ([]models.Artist, error) {
	if q == "" {
		return nil, errors.New("q is required")
	}
	return s.repo.SearchArtists(ctx, q, limit, offset)
}

func (s *searchService) Genres(ctx context.Context, q string, limit, offset int) ([]models.Genre, error) {
	if q == "" {
		return nil, errors.New("q is required")
	}
	return s.repo.SearchGenres(ctx, q, limit, offset)
}
