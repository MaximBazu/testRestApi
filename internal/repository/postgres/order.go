package postgres

import (
	"RESTAPI/internal/errs"
	"RESTAPI/internal/model"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type orderRepository struct {
	db DBTX
}

func NewOrderRepository(db *pgxpool.Pool) *orderRepository {
	return &orderRepository{db: db}
}

func NewOrderRepositoryTx(tx pgx.Tx) *orderRepository {
	return &orderRepository{db: tx}
}

func (r *orderRepository) GetByID(ctx context.Context, id int) (*model.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := `SELECT order_id, user_id, shipping_address, created_at FROM orders WHERE order_id = $1`

	var o model.Order
	if err := r.db.QueryRow(ctx, query, id).Scan(&o.ID, &o.UserID, &o.ShippingAddress, &o.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: %v", errs.ErrOrderNotFound, err)
		}
		return nil, MapPGError(err)
	}

	return &o, nil
}

func (r *orderRepository) List(ctx context.Context, limit, offset int) ([]model.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := `
		SELECT order_id, user_id, shipping_address, created_at
		FROM orders
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, MapPGError(err)
	}
	defer rows.Close()

	orders := make([]model.Order, 0, limit)
	for rows.Next() {
		var o model.Order
		if err := rows.Scan(&o.ID, &o.UserID, &o.ShippingAddress, &o.CreatedAt); err != nil {
			return nil, MapPGError(err)
		}
		orders = append(orders, o)
	}

	if err := rows.Err(); err != nil {
		return nil, MapPGError(err)
	}

	return orders, nil
}

func (r *orderRepository) Create(ctx context.Context, order *model.Order) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := `
		INSERT INTO orders (user_id, shipping_address)
		VALUES ($1, $2)
		RETURNING order_id, created_at
	`

	err := r.db.QueryRow(ctx, query, order.UserID, order.ShippingAddress).Scan(&order.ID, &order.CreatedAt)

	if err != nil {
		return MapPGError(err)
	}

	return nil
}

func (r *orderRepository) Delete(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	tag, err := r.db.Exec(ctx, `DELETE FROM orders WHERE order_id = $1`, id)
	if err != nil {
		return MapPGError(err)
	}
	if tag.RowsAffected() == 0 {
		return errs.ErrOrderNotFound
	}
	return nil
}
