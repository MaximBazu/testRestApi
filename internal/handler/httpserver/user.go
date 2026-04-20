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
