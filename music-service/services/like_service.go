package services

import (
	"context"
	"errors"
	"music-service/repositories"
)

type LikeService interface {
	Like(ctx context.Context, userID, trackID int) error
	Unlike(ctx context.Context, userID, trackID int) error
	MyLikes(ctx context.Context, userID int, limit, offset int) ([]int, error)
}

type likeService struct {
	repo repositories.LikeRepository
}

func NewLikeService(repo repositories.LikeRepository) LikeService {
	return &likeService{repo: repo}
}

func (s *likeService) Like(ctx context.Context, userID, trackID int) error {
	if userID <= 0 {
		return errors.New("invalid user_id")
	}
	if trackID <= 0 {
		return errors.New("invalid track_id")
	}
	return s.repo.Like(ctx, userID, trackID)
}

func (s *likeService) Unlike(ctx context.Context, userID, trackID int) error {
	if userID <= 0 {
		return errors.New("invalid user_id")
	}
	if trackID <= 0 {
		return errors.New("invalid track_id")
	}
	return s.repo.Unlike(ctx, userID, trackID)
}

func (s *likeService) MyLikes(ctx context.Context, userID int, limit, offset int) ([]int, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user_id")
	}
	return s.repo.ListLikedTrackIDs(ctx, userID, limit, offset)
}
