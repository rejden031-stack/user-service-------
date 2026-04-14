package http

import (
	"encoding/json"
	"net/http"
	"user-service/internal/service"
)

type Handler struct {
	userService *service.UserService
}

func NewHandler() *Handler {
	return &Handler{
		userService: service.NewUserService(),
	}
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	user := h.userService.GetUser()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
