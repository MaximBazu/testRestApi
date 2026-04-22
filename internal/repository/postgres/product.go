package postgres

import (
	"RESTAPI/internal/errs"
	"context"
	"errors"
	"fmt"
	"time"

	"RESTAPI/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository struct {
	db DBTX
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{db: db}
}

func NewProductRepositoryTx(tx pgx.Tx) *ProductRepository {
	return &ProductRepository{db: tx}
}

func (r *ProductRepository) GetByID(ctx context.Context, id int) (*model.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := `SELECT id, name, surname, email, telegram_tag, created_at FROM Products WHERE id=$1`

	var Product model.Product

	err := r.db.QueryRow(ctx, query, id).Scan(
		&Product.ID,
		&Product.Name,
		&Product.Surname,
		&Product.Email,
		&Product.TelegramTag,
		&Product.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: %v", errs.ErrProductNotFound, err)
		}
		return nil, err
	}

	return &Product, nil
}

func (r *ProductRepository) List(ctx context.Context, limit, offset int) ([]model.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := `
		SELECT id, name, surname, email, telegram_tag, created_at
		FROM Products
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	Products := make([]model.Product, 0, limit)
	for rows.Next() {
		var Product model.Product
		if err := rows.Scan(
			&Product.ID,
			&Product.Name,
			&Product.Surname,
			&Product.Email,
			&Product.TelegramTag,
			&Product.CreatedAt,
		); err != nil {
			return nil, err
		}
		Products = append(Products, Product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return Products, nil
}

func (r *ProductRepository) Create(ctx context.Context, Product *model.Product) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := `
		INSERT INTO Products (name, surname, email, telegram_tag)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	return r.db.QueryRow(
		ctx,
		query,
		Product.Name,
		Product.Surname,
		Product.Email,
		Product.TelegramTag,
	).Scan(&Product.ID, &Product.CreatedAt)
}

func (r *ProductRepository) Delete(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := `DELETE FROM Products WHERE id=$1`
	tag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return errs.ErrProductNotFound
	}

	return nil
}
