package httpserver

import (
	"RESTAPI/internal/dto"
	"RESTAPI/internal/errs"
	"RESTAPI/internal/mapper"
	"RESTAPI/internal/model"
	"RESTAPI/internal/service"
	"encoding/json"
	"net/http"
	"strconv"
)

type ProductImageHandler struct {
	productImageService service.ProductImageService
}

func NewProductImageHandler(productImageService service.ProductImageService) *ProductImageHandler {
	return &ProductImageHandler{productImageService: productImageService}
}

func (h *ProductImageHandler) GetByKey(w http.ResponseWriter, r *http.Request) {
	productID, imageURL, err := readProductImageKey(r)
	if err != nil {
		handleError(w, errs.ErrInvalidInput)
		return
	}

	productImage, err := h.productImageService.GetByKey(r.Context(), productID, imageURL)
	if err != nil {
		handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, mapper.ToProductImageResponse(productImage))
}

func (h *ProductImageHandler) List(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := getLimitOffset(r)
	if err != nil {
		handleError(w, errs.ErrInvalidInput)
		return
	}

	images, err := h.productImageService.List(r.Context(), limit, offset)
	if err != nil {
		handleError(w, err)
		return
	}

	resp := make([]dto.ProductImageResponse, 0, len(images))
	for _, image := range images {
		imageCopy := image
		resp = append(resp, mapper.ToProductImageResponse(&imageCopy))
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *ProductImageHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateProductImageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	productImage := &model.ProductImage{ProductID: req.ProductID, ImageURL: req.ImageURL}
	if err := h.productImageService.Create(r.Context(), productImage); err != nil {
		handleError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, mapper.ToProductImageResponse(productImage))
}

func (h *ProductImageHandler) Delete(w http.ResponseWriter, r *http.Request) {
	productID, imageURL, err := readProductImageKey(r)
	if err != nil {
		handleError(w, errs.ErrInvalidInput)
		return
	}

	if err := h.productImageService.Delete(r.Context(), productID, imageURL); err != nil {
		handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func readProductImageKey(r *http.Request) (int, string, error) {
	productIDStr := r.URL.Query().Get("product_id")
	imageURL := r.URL.Query().Get("image_url")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		return 0, "", err
	}
	return productID, imageURL, nil
}
