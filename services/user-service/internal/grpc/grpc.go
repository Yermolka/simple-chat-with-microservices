package grpc

import (
	"user-service/internal/repository"
	"user-service/proto"

	"google.golang.org/grpc"
)

func NewServer(repo *repository.Repository) *grpc.Server {
	srv := grpc.NewServer()
	authService := NewAuthService(repo)
	proto.RegisterAuthServiceServer(srv, authService)
	return srv
}
