package repository

import (
	"RESTAPI/internal/model"
	"context"
)

type UserRepository interface {
	GetByID(ctx context.Context, id int) (*model.User, error)
	Create(ctx context.Context, user *model.User) error
}
