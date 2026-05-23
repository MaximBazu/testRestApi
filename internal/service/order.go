package service

import (
	"RESTAPI/internal/dto"
	"RESTAPI/internal/errs"
	"RESTAPI/internal/model"
	"RESTAPI/internal/repository"
	"RESTAPI/internal/repository/postgres"
	"context"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderService interface {
	GetByID(ctx context.Context, id int) (*model.Order, error)
	List(ctx context.Context, limit, offset int) ([]model.Order, error)
	Create(ctx context.Context, order *model.Order) error
	Delete(ctx context.Context, id int) error

	Checkout(ctx context.Context, req dto.CreateOrderRequest) (*model.Order, error) // NEW
}

type orderService struct {
	repo repository.OrderRepository

	userRepo      repository.UserRepository
	orderItemRepo repository.OrderItemRepository
	sizeRepo      repository.ProductSizeRepository
	db            *pgxpool.Pool
}

func NewOrderService(
	db *pgxpool.Pool,
	orderRepo repository.OrderRepository,
	userRepo repository.UserRepository,
	orderItemRepo repository.OrderItemRepository,
	sizeRepo repository.ProductSizeRepository,
) OrderService {
	return &orderService{
		db:            db,
		repo:          orderRepo,
		userRepo:      userRepo,
		orderItemRepo: orderItemRepo,
		sizeRepo:      sizeRepo,
	}
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
	if order.UserID <= 0 || strings.TrimSpace(order.ShippingAddress) == "" || order.TotalAmount <= 0 {
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

func (s *orderService) Checkout(ctx context.Context, req dto.CreateOrderRequest) (*model.Order, error) {
	// 1) базовая валидация
	if strings.TrimSpace(req.ShippingAddress) == "" ||
		strings.TrimSpace(req.IdempotencyKey) == "" ||
		len(req.Items) == 0 ||
		req.TotalAmount <= 0 {
		return nil, errs.ErrInvalidInput
	}

	// 2) идемпотентность (первичная проверка до транзакции)
	existing, err := s.repo.GetByIdempotencyKey(ctx, req.IdempotencyKey)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		// MVP-вариант: возвращаем уже созданный заказ
		// (проверку "payload same/different" можешь добавить позже)
		return existing, nil
	}

	var createdOrder *model.Order

	err = postgres.WithTx(ctx, s.db, func(tx pgx.Tx) error {
		userRepoTx := postgres.NewUserRepositoryTx(tx)
		orderRepoTx := postgres.NewOrderRepositoryTx(tx)
		orderItemRepoTx := postgres.NewOrderItemRepositoryTx(tx)
		productRepoTx := postgres.NewProductRepositoryTx(tx)
		sizeRepoTx := postgres.NewProductSizeRepositoryTx(tx)

		// 3) create guest user
		user := &model.User{
			Name:       strings.TrimSpace(req.Buyer.Name),
			Surname:    strings.TrimSpace(req.Buyer.Surname),
			Patronymic: strings.TrimSpace(req.Buyer.Patronymic),
			Phone:      strings.TrimSpace(req.Buyer.Phone),
			Email:      strings.TrimSpace(req.Buyer.Email),
			TGTag:      strings.TrimSpace(req.Buyer.TGTag),
		}
		if err := userRepoTx.Create(ctx, user); err != nil {
			return err
		}

		// 4) проверка позиций + расчет server total
		var serverTotal int64

		for _, it := range req.Items {
			if it.ProductID <= 0 || it.SizeID <= 0 || it.Quantity <= 0 {
				return errs.ErrInvalidInput
			}

			size, err := sizeRepoTx.GetByID(ctx, it.SizeID)
			if err != nil {
				return err
			}
			if size.ProductID != it.ProductID {
				return errs.ErrInvalidInput
			}
			if size.Stock < it.Quantity {
				return errs.ErrInvalidInput // можно заменить на отдельную ErrInsufficientStock
			}

			product, err := productRepoTx.GetByID(ctx, it.ProductID)
			if err != nil {
				return err
			}

			serverTotal += product.Price * int64(it.Quantity)
		}

		// 5) сверка client total vs server total
		// req.TotalAmount у тебя int64; сравниваем через int64(serverTotal) в текущем MVP
		// (лучше позже перевести деньги полностью в копейки)
		if int64(serverTotal) != req.TotalAmount {
			return errs.ErrConflict
		}

		// 6) create order
		order := &model.Order{
			UserID:          user.ID,
			ShippingAddress: req.ShippingAddress,
			TotalAmount:     req.TotalAmount,
			IdempotencyKey:  req.IdempotencyKey,
		}
		if err := orderRepoTx.Create(ctx, order); err != nil {
			return err
		}

		// 7) create order_items + decrement stock
		for _, it := range req.Items {
			product, err := productRepoTx.GetByID(ctx, it.ProductID)
			if err != nil {
				return err
			}

			oi := &model.OrderItem{
				OrderID:         order.ID,
				ProductID:       it.ProductID,
				ProductSizeID:   it.SizeID,
				Quantity:        it.Quantity,
				PriceAtPurchase: product.Price,
			}
			if err := orderItemRepoTx.Create(ctx, oi); err != nil {
				return err
			}

			if err := sizeRepoTx.DecreaseStock(ctx, it.SizeID, it.Quantity); err != nil {
				return err
			}
		}

		createdOrder = order
		return nil
	})
	if err != nil {
		return nil, err
	}

	return createdOrder, nil
}
