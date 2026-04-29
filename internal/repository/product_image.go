package repository

import (
	"RESTAPI/internal/model"
	"context"
)

type ProductImageRepository interface {
	GetByKey(ctx context.Context, productID int, imageURL string) (*model.ProductImage, error)
	List(ctx context.Context, limit, offset int) ([]model.ProductImage, error)
	Create(ctx context.Context, productImage *model.ProductImage) error
	Delete(ctx context.Context, productID int, imageURL string) error
}
