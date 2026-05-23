package model

import "time"

type Order struct {
	ID              int
	UserID          int
	ShippingAddress string
	TotalAmount     int64
	IdempotencyKey  string
	CreatedAt       time.Time
}
