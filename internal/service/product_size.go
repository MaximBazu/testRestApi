package service

import (
	"RESTAPI/internal/errs"
	"RESTAPI/internal/model"
	"RESTAPI/internal/repository"
	"context"
	"strings"
)

type ProductSizeService interface {
	GetByID(ctx context.Context, id int) (*model.ProductSize, error)
	List(ctx context.Context, limit, offset int) ([]model.ProductSize, error)
	Create(ctx context.Context, productSize *model.ProductSize) error
	Delete(ctx context.Context, id int) error
}

type productSizeService struct {
	repo repository.ProductSizeRepository
}

func NewProductSizeService(repo repository.ProductSizeRepository) ProductSizeService {
	return &productSizeService{repo: repo}
}

func (s *productSizeService) GetByID(ctx context.Context, id int) (*model.ProductSize, error) {
	if id <= 0 {
		return nil, errs.ErrInvalidInput
	}
	return s.repo.GetByID(ctx, id)
}

func (s *productSizeService) List(ctx context.Context, limit, offset int) ([]model.ProductSize, error) {
	if limit <= 0 || limit > 100 || offset < 0 {
		return nil, errs.ErrInvalidInput
	}
	return s.repo.List(ctx, limit, offset)
}

func (s *productSizeService) Create(ctx context.Context, productSize *model.ProductSize) error {
	if productSize.ProductID <= 0 || strings.TrimSpace(productSize.Size) == "" || productSize.Stock < 0 {
		return errs.ErrInvalidInput
	}
	return s.repo.Create(ctx, productSize)
}

func (s *productSizeService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return errs.ErrInvalidInput
	}
	return s.repo.Delete(ctx, id)
}
