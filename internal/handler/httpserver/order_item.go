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

type OrderItemHandler struct {
	orderItemService service.OrderItemService
}

func NewOrderItemHandler(orderItemService service.OrderItemService) *OrderItemHandler {
	return &OrderItemHandler{orderItemService: orderItemService}
}

func (h *OrderItemHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		handleError(w, errs.ErrInvalidInput)
		return
	}

	orderItem, err := h.orderItemService.GetByID(r.Context(), id)
	if err != nil {
		handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, mapper.ToOrderItemResponse(orderItem))
}

func (h *OrderItemHandler) List(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := getLimitOffset(r)
	if err != nil {
		handleError(w, errs.ErrInvalidInput)
		return
	}

	items, err := h.orderItemService.List(r.Context(), limit, offset)
	if err != nil {
		handleError(w, err)
		return
	}

	resp := make([]dto.OrderItemResponse, 0, len(items))
	for _, item := range items {
		itemCopy := item
		resp = append(resp, mapper.ToOrderItemResponse(&itemCopy))
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *OrderItemHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateOrderItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	orderItem := &model.OrderItem{
		OrderID:         req.OrderID,
		ProductID:       req.ProductID,
		ProductSizeID:   req.ProductSizeID,
		Quantity:        req.Quantity,
		PriceAtPurchase: req.PriceAtPurchase,
	}
	if err := h.orderItemService.Create(r.Context(), orderItem); err != nil {
		handleError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, mapper.ToOrderItemResponse(orderItem))
}

func (h *OrderItemHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		handleError(w, errs.ErrInvalidInput)
		return
	}

	if err := h.orderItemService.Delete(r.Context(), id); err != nil {
		handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
