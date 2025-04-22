package http

import (
	"encoding/json"
	"net/http"
	"time"
	"user-service/internal/repository"

	"github.com/google/uuid"
)

type Handler struct {
	repo repository.IRepository
}

func NewHandler(repo repository.IRepository) *Handler {
	return &Handler{repo: repo}
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	userID, err := h.repo.Create(req.Username, req.Password)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int64{"id": userID})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.repo.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token := uuid.New().String()
	if err := h.repo.CreateToken(int64(user.Id), token); err != nil {
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}

	expiresAt := time.Now().Add(time.Hour)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(TokenResponse{
		Token:     token,
		ExpiresAt: expiresAt,
	})
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Authorization token is required", http.StatusUnauthorized)
		return
	}

	if err := h.repo.DeleteToken(token); err != nil {
		http.Error(w, "Failed to delete token", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
