package postgres

import (
	"RESTAPI/internal/dto"
	"RESTAPI/internal/errs"
	"RESTAPI/internal/model"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductDBTX interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

type ProductRepository struct {
	db ProductDBTX
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

	query := `
		SELECT product_id, name, description, price, slug, created_at
		FROM products
		WHERE product_id = $1
	`

	var p model.Product
	err := r.db.QueryRow(ctx, query, id).Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.Price,
		&p.Slug,
		&p.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: %v", errs.ErrProductNotFound, err)
		}
		return nil, err
	}

	return &p, nil
}

func (r *ProductRepository) List(ctx context.Context, limit, offset int) ([]model.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := `
		SELECT product_id, name, description, price, slug, created_at
		FROM products
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]model.Product, 0, limit)
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&p.Price,
			&p.Slug,
			&p.CreatedAt,
		); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductRepository) Create(ctx context.Context, p *model.Product) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := `
		INSERT INTO products (name, description, price, slug)
		VALUES ($1, $2, $3, $4)
		RETURNING product_id, created_at
	`

	return r.db.QueryRow(
		ctx,
		query,
		p.Name,
		p.Description,
		p.Price,
		p.Slug,
	).Scan(&p.ID, &p.CreatedAt)
}

func (r *ProductRepository) Delete(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := `DELETE FROM products WHERE product_id = $1`
	tag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return errs.ErrProductNotFound
	}

	return nil
}

func (r *ProductRepository) Update(ctx context.Context, id int, req dto.UpdateProductRequest) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := `
		UPDATE products
		SET
			name = COALESCE($1, name),
			description = COALESCE($2, description),
			price = COALESCE($3, price),
			slug = COALESCE($4, slug)
		WHERE product_id = $5
	`
	tag, err := r.db.Exec(ctx, query, req.Name, req.Description, req.Price, req.Slug, id)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return errs.ErrProductNotFound
	}

	return nil
}
