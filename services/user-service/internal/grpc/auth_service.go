package grpc

import (
	"context"
	"user-service/internal/repository"
	"user-service/proto"
)

type authService struct {
	proto.UnimplementedAuthServiceServer
	repo repository.IRepository
}

func NewAuthService(repo repository.IRepository) *authService {
	return &authService{
		repo: repo,
	}
}

func (s *authService) VerifyToken(ctx context.Context, req *proto.VerifyTokenRequest) (*proto.VerifyTokenResponse, error) {
	valid, err := s.repo.VerifyToken(ctx, req.UserId, req.Token)
	if err != nil {
		return &proto.VerifyTokenResponse{
			Valid: false,
			Error: err.Error(),
		}, nil
	}

	return &proto.VerifyTokenResponse{
		Valid: valid,
	}, nil
}
