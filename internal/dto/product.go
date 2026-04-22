package dto

import "time"

type CreateProductRequest struct {
	Name        string
	Description string
	Price       float64
	Slug        string
}

type ProductResponse struct {
	ID          int
	Name        string
	Description string
	Price       float64
	Slug        string
	CreatedAt   time.Time
}

type UpdateProductRequest struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price"`
	Slug        *string  `json:"slug"`
}
