package model

import "time"

type Product struct {
	ID          int
	Name        string
	Description string
	Price       float64
	Slug        string
	CreatedAt   time.Time
}
