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

type OrderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		handleError(w, errs.ErrInvalidInput)
		return
	}

	order, err := h.orderService.GetByID(r.Context(), id)
	if err != nil {
		handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, mapper.ToOrderResponse(order))
}

func (h *OrderHandler) List(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := getLimitOffset(r)
	if err != nil {
		handleError(w, errs.ErrInvalidInput)
		return
	}

	orders, err := h.orderService.List(r.Context(), limit, offset)
	if err != nil {
		handleError(w, err)
		return
	}

	resp := make([]dto.OrderResponse, 0, len(orders))
	for _, order := range orders {
		orderCopy := order
		resp = append(resp, mapper.ToOrderResponse(&orderCopy))
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	order := &model.Order{UserID: req.UserID, ShippingAddress: req.ShippingAddress}
	if err := h.orderService.Create(r.Context(), order); err != nil {
		handleError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, mapper.ToOrderResponse(order))
}

func (h *OrderHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		handleError(w, errs.ErrInvalidInput)
		return
	}

	if err := h.orderService.Delete(r.Context(), id); err != nil {
		handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
