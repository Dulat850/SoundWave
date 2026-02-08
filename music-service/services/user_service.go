package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"music-service/models"
	"music-service/repositories"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrConflict = errors.New("conflict")
)

type UserService interface {
	Me(ctx context.Context, userID int) (*models.User, error)
	UpdateMe(ctx context.Context, userID int, username string, email string) (*models.User, error)
	ChangePassword(ctx context.Context, userID int, oldPassword, newPassword string) error
	List(ctx context.Context, limit, offset int) ([]models.User, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Me(ctx context.Context, userID int) (*models.User, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user_id")
	}

	u, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	// пароль и так json:"-" но на всякий случай
	u.Password = ""
	return u, nil
}

func (s *userService) UpdateMe(ctx context.Context, userID int, username string, email string) (*models.User, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user_id")
	}
	if username == "" {
		return nil, errors.New("username is required")
	}
	if email == "" {
		return nil, errors.New("email is required")
	}

	u, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	// Проверяем уникальность username (если меняется)
	if u.Username != username {
		found, err := s.repo.GetByUsername(ctx, username)
		if err == nil && found != nil && found.ID != userID {
			return nil, fmt.Errorf("%w: username already used", ErrConflict)
		}
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	// Проверяем уникальность email (если меняется)
	if u.Email != email {
		found, err := s.repo.GetByEmail(ctx, email)
		if err == nil && found != nil && found.ID != userID {
			return nil, fmt.Errorf("%w: email already used", ErrConflict)
		}
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	u.Username = username
	u.Email = email

	if err := s.repo.Update(ctx, u); err != nil {
		return nil, err
	}

	u.Password = ""
	return u, nil
}

func (s *userService) ChangePassword(ctx context.Context, userID int, oldPassword, newPassword string) error {
	if userID <= 0 {
		return errors.New("invalid user_id")
	}
	if oldPassword == "" {
		return errors.New("old_password is required")
	}
	if newPassword == "" {
		return errors.New("new_password is required")
	}

	u, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(oldPassword)); err != nil {
		return errors.New("invalid old password")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hash)
	return s.repo.Update(ctx, u)
}

func (s *userService) List(ctx context.Context, limit, offset int) ([]models.User, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.List(ctx, limit, offset)
}
