package mapper

import (
	"RESTAPI/internal/dto"
	"RESTAPI/internal/model"
)

func ToProductImageResponse(pi *model.ProductImage) dto.ProductImageResponse {
	return dto.ProductImageResponse{
		ProductID: pi.ProductID,
		ImageURL:  pi.ImageURL,
	}
}
