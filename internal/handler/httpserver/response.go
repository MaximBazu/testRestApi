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
	case errors.Is(err, errs.ErrUserNotFound):
		writeError(w, http.StatusNotFound, "user not found")

	case errors.Is(err, errs.ErrInvalidInput):
		writeError(w, http.StatusBadRequest, "invalid input")

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
