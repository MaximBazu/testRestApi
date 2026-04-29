package repository

import (
	"RESTAPI/internal/model"
	"context"
)

type OrderRepository interface {
	GetByID(ctx context.Context, id int) (*model.Order, error)
	List(ctx context.Context, limit, offset int) ([]model.Order, error)
	Create(ctx context.Context, order *model.Order) error
	Delete(ctx context.Context, id int) error
}
