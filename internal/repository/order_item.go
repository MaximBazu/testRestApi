package repository

import (
	"RESTAPI/internal/model"
	"context"
)

type OrderItemRepository interface {
	GetByID(ctx context.Context, id int) (*model.OrderItem, error)
	List(ctx context.Context, limit, offset int) ([]model.OrderItem, error)
	Create(ctx context.Context, orderItem *model.OrderItem) error
	Delete(ctx context.Context, id int) error
}
