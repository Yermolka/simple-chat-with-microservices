package http

import (
	"net/http"
	"user-service/internal/repository"
)

func NewServer(repo repository.IRepository) *http.Server {
	handler := NewHandler(repo)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/users", handler.CreateUser)
	mux.HandleFunc("POST /api/login", handler.Login)
	mux.HandleFunc("POST /api/logout", handler.Logout)

	return &http.Server{
		Addr:    ":3001",
		Handler: mux,
	}
}
