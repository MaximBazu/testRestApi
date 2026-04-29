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

	"github.com/go-chi/chi/v5"
)

type ProductSizeHandler struct {
	productSizeService service.ProductSizeService
}

func NewProductSizeHandler(productSizeService service.ProductSizeService) *ProductSizeHandler {
	return &ProductSizeHandler{productSizeService: productSizeService}
}

func (h *ProductSizeHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		handleError(w, errs.ErrInvalidInput)
		return
	}

	productSize, err := h.productSizeService.GetByID(r.Context(), id)
	if err != nil {
		handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, mapper.ToProductSizeResponse(productSize))
}

func (h *ProductSizeHandler) List(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := getLimitOffset(r)
	if err != nil {
		handleError(w, errs.ErrInvalidInput)
		return
	}

	sizes, err := h.productSizeService.List(r.Context(), limit, offset)
	if err != nil {
		handleError(w, err)
		return
	}

	resp := make([]dto.ProductSizeResponse, 0, len(sizes))
	for _, size := range sizes {
		sizeCopy := size
		resp = append(resp, mapper.ToProductSizeResponse(&sizeCopy))
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *ProductSizeHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateProductSizeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	productSize := &model.ProductSize{ProductID: req.ProductID, Size: req.Size, Stock: req.Stock}
	if err := h.productSizeService.Create(r.Context(), productSize); err != nil {
		handleError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, mapper.ToProductSizeResponse(productSize))
}

func (h *ProductSizeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		handleError(w, errs.ErrInvalidInput)
		return
	}

	if err := h.productSizeService.Delete(r.Context(), id); err != nil {
		handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
