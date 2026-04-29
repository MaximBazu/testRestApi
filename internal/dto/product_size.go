package dto

type CreateProductSizeRequest struct {
	ProductID int    `json:"product_id"`
	Size      string `json:"size"`
	Stock     int    `json:"stock"`
}

type ProductSizeResponse struct {
	ID        int    `json:"id"`
	ProductID int    `json:"product_id"`
	Size      string `json:"size"`
	Stock     int    `json:"stock"`
}
