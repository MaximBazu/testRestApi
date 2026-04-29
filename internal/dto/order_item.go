package dto

type CreateOrderItemRequest struct {
	OrderID         int     `json:"order_id"`
	ProductID       int     `json:"product_id"`
	ProductSizeID   int     `json:"product_size_id"`
	Quantity        int     `json:"quantity"`
	PriceAtPurchase float64 `json:"price_at_purchase"`
}

type OrderItemResponse struct {
	ID              int     `json:"id"`
	OrderID         int     `json:"order_id"`
	ProductID       int     `json:"product_id"`
	ProductSizeID   int     `json:"product_size_id"`
	Quantity        int     `json:"quantity"`
	PriceAtPurchase float64 `json:"price_at_purchase"`
}
