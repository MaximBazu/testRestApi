package service

import (
	"RESTAPI/internal/dto"
	"context"
	"strings"

	"RESTAPI/internal/errs"
	"RESTAPI/internal/model"
	"RESTAPI/internal/repository"
)

type ProductService interface {
	GetByID(ctx context.Context, id int) (*model.Product, error)
	List(ctx context.Context, limit, offset int) ([]model.Product, error)
	Create(ctx context.Context, Product *model.Product) error
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, id int, req dto.UpdateProductRequest) error
}
type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) GetByID(ctx context.Context, id int) (*model.Product, error) {
	Product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return Product, nil
}

func (s *productService) List(ctx context.Context, limit, offset int) ([]model.Product, error) {
	if limit <= 0 || limit > 100 {
		return nil, errs.ErrInvalidInput
	}
	if offset < 0 {
		return nil, errs.ErrInvalidInput
	}

	Products, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	return Products, nil
}

func (s *productService) Create(ctx context.Context, Product *model.Product) error {

	if err := s.repo.Create(ctx, Product); err != nil {
		return err
	}

	return nil
}

func (s *productService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return errs.ErrInvalidInput
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

func (s *productService) Update(ctx context.Context, id int, req dto.UpdateProductRequest) error {
	if id <= 0 {
		return errs.ErrInvalidInput
	}

	// хотя бы одно поле должно прийти
	if req.Name == nil && req.Description == nil && req.Price == nil && req.Slug == nil {
		return errs.ErrInvalidInput
	}

	if req.Name != nil && strings.TrimSpace(*req.Name) == "" {
		return errs.ErrInvalidInput
	}
	if req.Price != nil && *req.Price < 0 {
		return errs.ErrInvalidInput
	}
	if req.Slug != nil && strings.TrimSpace(*req.Slug) == "" {
		return errs.ErrInvalidInput
	}

	return s.repo.Update(ctx, id, req)
}
