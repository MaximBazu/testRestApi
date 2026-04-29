package httpserver

import (
	"net/http"
	"strconv"
)

func getLimitOffset(r *http.Request) (int, int, error) {
	limit := 20
	offset := 0

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil {
			return 0, 0, err
		}
		limit = parsedLimit
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err != nil {
			return 0, 0, err
		}
		offset = parsedOffset
	}

	return limit, offset, nil
}
