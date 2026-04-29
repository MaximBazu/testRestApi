package repository

import (
	"RESTAPI/internal/model"
	"context"
)

type ProductSizeRepository interface {
	GetByID(ctx context.Context, id int) (*model.ProductSize, error)
	List(ctx context.Context, limit, offset int) ([]model.ProductSize, error)
	Create(ctx context.Context, productSize *model.ProductSize) error
	Delete(ctx context.Context, id int) error
}
