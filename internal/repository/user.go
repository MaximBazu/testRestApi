package repository

import (
	"RESTAPI/internal/model"
	"context"
)

type UserRepository interface {
	GetByID(ctx context.Context, id int) (*model.User, error)
	List(ctx context.Context, limit, offset int) ([]model.User, error)
	Create(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id int) error
}
