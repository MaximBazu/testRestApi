package service

import (
	"RESTAPI/internal/errs"
	"RESTAPI/internal/model"
	"RESTAPI/internal/repository"
	"context"
	"strings"
)

type OrderService interface {
	GetByID(ctx context.Context, id int) (*model.Order, error)
	List(ctx context.Context, limit, offset int) ([]model.Order, error)
	Create(ctx context.Context, order *model.Order) error
	Delete(ctx context.Context, id int) error
}

type orderService struct {
	repo repository.OrderRepository
}

func NewOrderService(repo repository.OrderRepository) OrderService {
	return &orderService{repo: repo}
}

func (s *orderService) GetByID(ctx context.Context, id int) (*model.Order, error) {
	if id <= 0 {
		return nil, errs.ErrInvalidInput
	}
	return s.repo.GetByID(ctx, id)
}

func (s *orderService) List(ctx context.Context, limit, offset int) ([]model.Order, error) {
	if limit <= 0 || limit > 100 || offset < 0 {
		return nil, errs.ErrInvalidInput
	}
	return s.repo.List(ctx, limit, offset)
}

func (s *orderService) Create(ctx context.Context, order *model.Order) error {
	if order.UserID <= 0 || strings.TrimSpace(order.ShippingAddress) == "" {
		return errs.ErrInvalidInput
	}
	return s.repo.Create(ctx, order)
}

func (s *orderService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return errs.ErrInvalidInput
	}
	return s.repo.Delete(ctx, id)
}
