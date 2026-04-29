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

type productSizeRepository struct {
	db DBTX
}

func NewProductSizeRepository(db *pgxpool.Pool) *productSizeRepository {
	return &productSizeRepository{db: db}
}

func NewProductSizeRepositoryTx(tx pgx.Tx) *productSizeRepository {
	return &productSizeRepository{db: tx}
}

func (r *productSizeRepository) GetByID(ctx context.Context, id int) (*model.ProductSize, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := `SELECT product_size_id, product_id, size, stock FROM product_sizes WHERE product_size_id = $1`

	var ps model.ProductSize
	if err := r.db.QueryRow(ctx, query, id).Scan(&ps.ID, &ps.ProductID, &ps.Size, &ps.Stock); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: %v", errs.ErrProductSizeNotFound, err)
		}
		return nil, MapPGError(err)
	}
	return &ps, nil
}

func (r *productSizeRepository) List(ctx context.Context, limit, offset int) ([]model.ProductSize, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := `
		SELECT product_size_id, product_id, size, stock
		FROM product_sizes
		ORDER BY product_size_id DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, MapPGError(err)
	}
	defer rows.Close()

	sizes := make([]model.ProductSize, 0, limit)
	for rows.Next() {
		var ps model.ProductSize
		if err := rows.Scan(&ps.ID, &ps.ProductID, &ps.Size, &ps.Stock); err != nil {
			return nil, MapPGError(err)
		}
		sizes = append(sizes, ps)
	}

	if err := rows.Err(); err != nil {
		return nil, MapPGError(err)
	}
	return sizes, nil
}

func (r *productSizeRepository) Create(ctx context.Context, productSize *model.ProductSize) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := `
		INSERT INTO product_sizes (product_id, size, stock)
		VALUES ($1, $2, $3)
		RETURNING product_size_id
	`

	err := r.db.QueryRow(ctx, query, productSize.ProductID, productSize.Size, productSize.Stock).Scan(&productSize.ID)

	if err != nil {
		return MapPGError(err)
	}

	return nil
}

func (r *productSizeRepository) Delete(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	tag, err := r.db.Exec(ctx, `DELETE FROM product_sizes WHERE product_size_id = $1`, id)
	if err != nil {
		return MapPGError(err)
	}
	if tag.RowsAffected() == 0 {
		return errs.ErrProductSizeNotFound
	}
	return nil
}
