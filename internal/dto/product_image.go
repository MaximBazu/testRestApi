package dto

type CreateProductImageRequest struct {
	ProductID int    `json:"product_id"`
	ImageURL  string `json:"image_url"`
}

type ProductImageResponse struct {
	ProductID int    `json:"product_id"`
	ImageURL  string `json:"image_url"`
}
