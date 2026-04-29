package service

import (
	"RESTAPI/internal/errs"
	"RESTAPI/internal/model"
	"RESTAPI/internal/repository"
	"context"
	"strings"
)

type ProductImageService interface {
	GetByKey(ctx context.Context, productID int, imageURL string) (*model.ProductImage, error)
	List(ctx context.Context, limit, offset int) ([]model.ProductImage, error)
	Create(ctx context.Context, productImage *model.ProductImage) error
	Delete(ctx context.Context, productID int, imageURL string) error
}

type productImageService struct {
	repo repository.ProductImageRepository
}

func NewProductImageService(repo repository.ProductImageRepository) ProductImageService {
	return &productImageService{repo: repo}
}

func (s *productImageService) GetByKey(ctx context.Context, productID int, imageURL string) (*model.ProductImage, error) {
	if productID <= 0 || strings.TrimSpace(imageURL) == "" {
		return nil, errs.ErrInvalidInput
	}
	return s.repo.GetByKey(ctx, productID, imageURL)
}

func (s *productImageService) List(ctx context.Context, limit, offset int) ([]model.ProductImage, error) {
	if limit <= 0 || limit > 100 || offset < 0 {
		return nil, errs.ErrInvalidInput
	}
	return s.repo.List(ctx, limit, offset)
}

func (s *productImageService) Create(ctx context.Context, productImage *model.ProductImage) error {
	if productImage.ProductID <= 0 || strings.TrimSpace(productImage.ImageURL) == "" {
		return errs.ErrInvalidInput
	}
	return s.repo.Create(ctx, productImage)
}

func (s *productImageService) Delete(ctx context.Context, productID int, imageURL string) error {
	if productID <= 0 || strings.TrimSpace(imageURL) == "" {
		return errs.ErrInvalidInput
	}
	return s.repo.Delete(ctx, productID, imageURL)
}
