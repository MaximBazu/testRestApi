package dto

import "time"

type CreateOrderRequest struct {
	UserID          int    `json:"user_id"`
	ShippingAddress string `json:"shipping_address"`
}

type OrderResponse struct {
	ID              int       `json:"id"`
	UserID          int       `json:"user_id"`
	ShippingAddress string    `json:"shipping_address"`
	CreatedAt       time.Time `json:"created_at"`
}
