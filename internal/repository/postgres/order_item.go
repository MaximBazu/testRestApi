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

type orderItemRepository struct {
	db DBTX
}

func NewOrderItemRepository(db *pgxpool.Pool) *orderItemRepository {
	return &orderItemRepository{db: db}
}

func NewOrderItemRepositoryTx(tx pgx.Tx) *orderItemRepository {
	return &orderItemRepository{db: tx}
}

func (r *orderItemRepository) GetByID(ctx context.Context, id int) (*model.OrderItem, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := `
		SELECT order_item_id, order_id, product_id, product_size_id, quantity, price_at_purchase
		FROM order_items
		WHERE order_item_id = $1
	`

	var oi model.OrderItem
	if err := r.db.QueryRow(ctx, query, id).Scan(&oi.ID, &oi.OrderID, &oi.ProductID, &oi.ProductSizeID, &oi.Quantity, &oi.PriceAtPurchase); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: %v", errs.ErrOrderItemNotFound, err)
		}
		return nil, MapPGError(err)
	}
	return &oi, nil
}

func (r *orderItemRepository) List(ctx context.Context, limit, offset int) ([]model.OrderItem, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := `
		SELECT order_item_id, order_id, product_id, product_size_id, quantity, price_at_purchase
		FROM order_items
		ORDER BY order_item_id DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, MapPGError(err)
	}
	defer rows.Close()

	items := make([]model.OrderItem, 0, limit)
	for rows.Next() {
		var oi model.OrderItem
		if err := rows.Scan(&oi.ID, &oi.OrderID, &oi.ProductID, &oi.ProductSizeID, &oi.Quantity, &oi.PriceAtPurchase); err != nil {
			return nil, MapPGError(err)
		}
		items = append(items, oi)
	}

	if err := rows.Err(); err != nil {
		return nil, MapPGError(err)
	}
	return items, nil
}

func (r *orderItemRepository) Create(ctx context.Context, orderItem *model.OrderItem) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := `
		INSERT INTO order_items (order_id, product_id, product_size_id, quantity, price_at_purchase)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING order_item_id
	`

	err := r.db.QueryRow(ctx, query, orderItem.OrderID, orderItem.ProductID, orderItem.ProductSizeID, orderItem.Quantity, orderItem.PriceAtPurchase).Scan(&orderItem.ID)

	if err != nil {
		return MapPGError(err)
	}

	return nil
}

func (r *orderItemRepository) Delete(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	tag, err := r.db.Exec(ctx, `DELETE FROM order_items WHERE order_item_id = $1`, id)
	if err != nil {
		return MapPGError(err)
	}
	if tag.RowsAffected() == 0 {
		return errs.ErrOrderItemNotFound
	}
	return nil
}
