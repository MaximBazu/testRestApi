package dto

import "time"

type CheckoutBuyer struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	TGTag      string `json:"tg_tag"`
}

type CheckoutOrderItem struct {
	ProductID int `json:"product_id"`
	SizeID    int `json:"size_id"`
	Quantity  int `json:"quantity"`
}

type CreateOrderRequest struct {
	Buyer           CheckoutBuyer       `json:"buyer"`
	Items           []CheckoutOrderItem `json:"items"`
	ShippingAddress string              `json:"shipping_address"`
	TotalAmount     int64               `json:"total_amount"`
	IdempotencyKey  string              `json:"idempotency_key"`
}

type OrderResponse struct {
	ID              int       `json:"id"`
	UserID          int       `json:"user_id"`
	ShippingAddress string    `json:"shipping_address"`
	TotalAmount     int64     `json:"total_amount"`
	CreatedAt       time.Time `json:"created_at"`
}
