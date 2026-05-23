package model

import "time"

type Product struct {
	ID          int
	Name        string
	Description string
	Price       int64
	Slug        string
	CreatedAt   time.Time
}
