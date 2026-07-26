package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"user-service/internal/entity"
)

type UserService interface {
	GetUser(ctx context.Context, id int) (entity.User, error)
	GetAllUsers(ctx context.Context) ([]entity.User, error)
	CreateUser(ctx context.Context, name string) (entity.User, error)
}

type Handler struct {
	userService UserService
}

func NewHandler(userService UserService) *Handler {
	return &Handler{userService: userService}
}

type userResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func toResponse(u entity.User) userResponse {
	return userResponse{ID: u.ID, Name: u.Name}
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "missing id parameter", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id parameter", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUser(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(toResponse(user))
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.GetAllUsers(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := make([]userResponse, 0, len(users))
	for _, u := range users {
		resp = append(resp, toResponse(u))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.userService.CreateUser(r.Context(), req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(toResponse(user))
}
