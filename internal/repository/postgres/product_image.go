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

type productImageRepository struct {
	db DBTX
}

func NewProductImageRepository(db *pgxpool.Pool) *productImageRepository {
	return &productImageRepository{db: db}
}

func NewProductImageRepositoryTx(tx pgx.Tx) *productImageRepository {
	return &productImageRepository{db: tx}
}

func (r *productImageRepository) GetByKey(ctx context.Context, productID int, imageURL string) (*model.ProductImage, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := `SELECT product_id, image_url FROM product_images WHERE product_id = $1 AND image_url = $2`

	var pi model.ProductImage
	if err := r.db.QueryRow(ctx, query, productID, imageURL).Scan(&pi.ProductID, &pi.ImageURL); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: %v", errs.ErrProductImageNotFound, err)
		}
		return nil, MapPGError(err)
	}
	return &pi, nil
}

func (r *productImageRepository) List(ctx context.Context, limit, offset int) ([]model.ProductImage, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := `
		SELECT product_id, image_url
		FROM product_images
		ORDER BY product_id DESC, image_url ASC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, MapPGError(err)
	}
	defer rows.Close()

	images := make([]model.ProductImage, 0, limit)
	for rows.Next() {
		var pi model.ProductImage
		if err := rows.Scan(&pi.ProductID, &pi.ImageURL); err != nil {
			return nil, MapPGError(err)
		}
		images = append(images, pi)
	}

	if err := rows.Err(); err != nil {
		return nil, MapPGError(err)
	}

	return images, nil
}

func (r *productImageRepository) Create(ctx context.Context, productImage *model.ProductImage) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	_, err := r.db.Exec(ctx, `INSERT INTO product_images (product_id, image_url) VALUES ($1, $2)`, productImage.ProductID, productImage.ImageURL)
	return MapPGError(err)
}

func (r *productImageRepository) Delete(ctx context.Context, productID int, imageURL string) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	tag, err := r.db.Exec(ctx, `DELETE FROM product_images WHERE product_id = $1 AND image_url = $2`, productID, imageURL)
	if err != nil {
		return MapPGError(err)
	}
	if tag.RowsAffected() == 0 {
		return errs.ErrProductImageNotFound
	}
	return nil
}
