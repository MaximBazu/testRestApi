package service

import (
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
}
type ProductService struct {
	repo repository.ProductRepository
}

func NewProductrService(repo repository.ProductRepository) ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetByID(ctx context.Context, id int) (*model.Product, error) {
	Product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return Product, nil
}

func (s *ProductService) List(ctx context.Context, limit, offset int) ([]model.Product, error) {
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

func (s *ProductService) Create(ctx context.Context, Product *model.Product) error {
	if err := validateProduct(Product); err != nil {
		return err
	}

	if err := s.repo.Create(ctx, Product); err != nil {
		return err
	}

	return nil
}

func (s *ProductService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return errs.ErrInvalidInput
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

func validateProduct(u *model.Product) error {
	if strings.TrimSpace(u.Name) == "" {
		return errs.ErrInvalidInput
	}

	if strings.TrimSpace(u.Email) == "" {
		return errs.ErrInvalidInput
	}

	if !strings.Contains(u.Email, "@") {
		return errs.ErrInvalidInput
	}

	return nil
}
