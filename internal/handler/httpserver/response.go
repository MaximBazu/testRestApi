package httpserver

import (
	"RESTAPI/internal/errs"
	"encoding/json"
	"errors"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, ErrorResponse{
		Error: msg,
	})
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
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
