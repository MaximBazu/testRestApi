package model

type OrderItem struct {
	ID              int
	OrderID         int
	ProductID       int
	ProductSizeID   int
	Quantity        int
	PriceAtPurchase float64
}
