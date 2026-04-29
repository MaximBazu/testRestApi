package mapper

import (
	"RESTAPI/internal/dto"
	"RESTAPI/internal/model"
)

func ToProductSizeResponse(ps *model.ProductSize) dto.ProductSizeResponse {
	return dto.ProductSizeResponse{
		ID:        ps.ID,
		ProductID: ps.ProductID,
		Size:      ps.Size,
		Stock:     ps.Stock,
	}
}
