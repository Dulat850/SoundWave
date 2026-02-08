package services

import (
	"context"
	"errors"
	"music-service/models"
	"music-service/repositories"
)

type ListenService interface {
	RecordListen(ctx context.Context, userID *int, trackID int) error
	MyHistory(ctx context.Context, userID int, limit, offset int) ([]models.Listen, error)
}

type listenService struct {
	repo repositories.ListenRepository
}

func NewListenService(repo repositories.ListenRepository) ListenService {
	return &listenService{repo: repo}
}

func (s *listenService) RecordListen(ctx context.Context, userID *int, trackID int) error {
	if trackID <= 0 {
		return errors.New("invalid track_id")
	}

	// Событие + увеличение счётчика для "популярного"
	if err := s.repo.Record(ctx, userID, trackID); err != nil {
		return err
	}
	return s.repo.IncrementPlayCount(ctx, trackID)
}

func (s *listenService) MyHistory(ctx context.Context, userID int, limit, offset int) ([]models.Listen, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user_id")
	}
	return s.repo.UserHistory(ctx, userID, limit, offset)
}
