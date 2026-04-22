package mapper

import (
	"RESTAPI/internal/dto"
	"RESTAPI/internal/model"
)

func ToProductResponse(p *model.Product) dto.ProductResponse {
	return dto.ProductResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Slug:        p.Slug,
		CreatedAt:   p.CreatedAt,
	}
}
