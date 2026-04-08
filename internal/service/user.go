package service

import (
	"context"
	"strings"

	"RESTAPI/internal/errs"
	"RESTAPI/internal/model"
	"RESTAPI/internal/repository"
)

type UserService interface {
	GetByID(ctx context.Context, id int) (*model.User, error)
	Create(ctx context.Context, user *model.User) error
}
type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetByID(ctx context.Context, id int) (*model.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) Create(ctx context.Context, user *model.User) error {
	if err := validateUser(user); err != nil {
		return err
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return err
	}

	return nil
}

func validateUser(u *model.User) error {
	if strings.TrimSpace(u.Name) == "" {
		return errs.ErrInvalidInput
	}

	if strings.TrimSpace(u.Email) == "" {
		return errs.ErrInvalidInput
	}

	if !strings.Contains(u.Email, "@") {
		return errs.ErrInvalidInput
	}

	return nil
}
