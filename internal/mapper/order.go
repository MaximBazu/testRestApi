package mapper

import (
	"RESTAPI/internal/dto"
	"RESTAPI/internal/model"
)

func ToOrderResponse(o *model.Order) dto.OrderResponse {
	return dto.OrderResponse{
		ID:              o.ID,
		UserID:          o.UserID,
		ShippingAddress: o.ShippingAddress,
		CreatedAt:       o.CreatedAt,
	}
}
