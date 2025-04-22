package http

import (
	"net/http"
	"user-service/internal/repository"
)

func NewServer(repo repository.IRepository) *http.Server {
	handler := NewHandler(repo)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /users", handler.CreateUser)
	mux.HandleFunc("POST /login", handler.Login)
	mux.HandleFunc("POST /logout", handler.Logout)

	return &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
}
