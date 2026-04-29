package service

import (
	"RESTAPI/internal/errs"
	"RESTAPI/internal/model"
	"RESTAPI/internal/repository"
	"context"
)

type OrderItemService interface {
	GetByID(ctx context.Context, id int) (*model.OrderItem, error)
	List(ctx context.Context, limit, offset int) ([]model.OrderItem, error)
	Create(ctx context.Context, orderItem *model.OrderItem) error
	Delete(ctx context.Context, id int) error
}

type orderItemService struct {
	repo repository.OrderItemRepository
}

func NewOrderItemService(repo repository.OrderItemRepository) OrderItemService {
	return &orderItemService{repo: repo}
}

func (s *orderItemService) GetByID(ctx context.Context, id int) (*model.OrderItem, error) {
	if id <= 0 {
		return nil, errs.ErrInvalidInput
	}
	return s.repo.GetByID(ctx, id)
}

func (s *orderItemService) List(ctx context.Context, limit, offset int) ([]model.OrderItem, error) {
	if limit <= 0 || limit > 100 || offset < 0 {
		return nil, errs.ErrInvalidInput
	}
	return s.repo.List(ctx, limit, offset)
}

func (s *orderItemService) Create(ctx context.Context, orderItem *model.OrderItem) error {
	if orderItem.OrderID <= 0 || orderItem.ProductID <= 0 || orderItem.ProductSizeID <= 0 || orderItem.Quantity <= 0 || orderItem.PriceAtPurchase < 0 {
		return errs.ErrInvalidInput
	}
	return s.repo.Create(ctx, orderItem)
}

func (s *orderItemService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return errs.ErrInvalidInput
	}
	return s.repo.Delete(ctx, id)
}
