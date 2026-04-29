package model

import "time"

type Order struct {
	ID              int
	UserID          int
	ShippingAddress string
	CreatedAt       time.Time
}
