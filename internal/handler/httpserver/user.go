package httpserver

import (
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
	// --- params ---
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	// --- service ---
	user, err := h.userService.GetByID(r.Context(), id)
	if err != nil {
		handleError(w, err)
		return
	}

	// --- response ---
	resp := mapper.ToUserResponse(user)
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
