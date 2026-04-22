package httpserver

import (
	"RESTAPI/internal/errs"
	"encoding/json"
	"net/http"
	"strconv"

	"RESTAPI/internal/dto"
	"RESTAPI/internal/mapper"
	"RESTAPI/internal/model"
	"RESTAPI/internal/service"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// 1. Получаем ID из URL (chi.URLParam)
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr) // "1" → 1
	if err != nil {
		handleError(w, errs.ErrInvalidInput) // ← 400 Bad Request
		return
	}

	// 2. Вызываем сервис
	user, err := h.userService.GetByID(r.Context(), id)
	if err != nil {
		handleError(w, err) // ← 404 или 500
		return
	}

	// 3. Маппим модель → DTO для ответа
	response := mapper.ToUserResponse(user)

	// 4. Отправляем JSON клиенту
	writeJSON(w, http.StatusOK, response) // ← 200 + {"id":1,"name":"Alice",...}
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	limit := 20
	offset := 0

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil {
			handleError(w, errs.ErrInvalidInput)
			return
		}
		limit = parsedLimit
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err != nil {
			handleError(w, errs.ErrInvalidInput)
			return
		}
		offset = parsedOffset
	}

	users, err := h.userService.List(r.Context(), limit, offset)
	if err != nil {
		handleError(w, err)
		return
	}

	resp := make([]dto.UserResponse, 0, len(users))
	for _, user := range users {
		userCopy := user
		resp = append(resp, mapper.ToUserResponse(&userCopy))
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	// --- request ---
	var req dto.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	// --- model ---
	user := &model.User{
		Name:        req.Name,
		Surname:     req.Surname,
		Email:       req.Email,
		TelegramTag: req.TelegramTag,
	}

	// --- service ---
	if err := h.userService.Create(r.Context(), user); err != nil {
		handleError(w, err)
		return
	}

	// --- response ---
	resp := mapper.ToUserResponse(user)
	writeJSON(w, http.StatusCreated, resp)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		handleError(w, errs.ErrInvalidInput)
		return
	}

	if err := h.userService.Delete(r.Context(), id); err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
