package mapper

import (
	"RESTAPI/internal/dto"
	"RESTAPI/internal/model"
)

func ToOrderItemResponse(oi *model.OrderItem) dto.OrderItemResponse {
	return dto.OrderItemResponse{
		ID:              oi.ID,
		OrderID:         oi.OrderID,
		ProductID:       oi.ProductID,
		ProductSizeID:   oi.ProductSizeID,
		Quantity:        oi.Quantity,
		PriceAtPurchase: oi.PriceAtPurchase,
	}
}
