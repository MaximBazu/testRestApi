package repository

import (
	"RESTAPI/internal/model"
	"context"
)

type ProductRepository interface {
	GetByID(ctx context.Context, id int) (*model.Product, error)
	List(ctx context.Context, limit, offset int) ([]model.Product, error)
	Create(ctx context.Context, Product *model.Product) error
	Delete(ctx context.Context, id int) error
}
