package httpserver

import (
	"RESTAPI/internal/errs"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func handleError(w http.ResponseWriter, err error) {
	switch {
	// 404
	case errors.Is(err, errs.ErrUserNotFound):
		writeError(w, http.StatusNotFound, "user not found")
	case errors.Is(err, errs.ErrProductNotFound):
		writeError(w, http.StatusNotFound, "product not found")
	case errors.Is(err, errs.ErrOrderNotFound):
		writeError(w, http.StatusNotFound, "order not found")
	case errors.Is(err, errs.ErrOrderItemNotFound):
		writeError(w, http.StatusNotFound, "order item not found")
	case errors.Is(err, errs.ErrProductSizeNotFound):
		writeError(w, http.StatusNotFound, "product size not found")
	case errors.Is(err, errs.ErrProductImageNotFound):
		writeError(w, http.StatusNotFound, "product image not found")

	// 400
	case errors.Is(err, errs.ErrInvalidInput):
		writeError(w, http.StatusBadRequest, "invalid input")
	case errors.Is(err, errs.ErrBadFormat):
		writeError(w, http.StatusBadRequest, "bad format")
	case errors.Is(err, errs.ErrValueTooLong):
		writeError(w, http.StatusBadRequest, "value too long")
	case errors.Is(err, errs.ErrNotNull):
		writeError(w, http.StatusBadRequest, "not null violation")

	// 409
	case errors.Is(err, errs.ErrConflict):
		writeError(w, http.StatusConflict, "conflict")

	// 422
	case errors.Is(err, errs.ErrForeignKey):
		writeError(w, http.StatusUnprocessableEntity, "foreign key violation")

	// 500
	default:
		writeError(w, http.StatusInternalServerError, "internal error")
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, ErrorResponse{
		Error: msg,
	})
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	// 1. Заголовок: говорим клиенту, что дальше будет JSON
	w.Header().Set("Content-Type", "application/json")

	// 2. Статус-код: сообщаем, успешно ли выполнился запрос
	w.WriteHeader(status) // ← После этой строки заголовки "улетели" клиенту!

	// 3. Тело: сериализуем данные в JSON и пишем в ответ
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}
