package grpc

import (
	"context"
	"testing"
	"user-service/internal/repository"
	"user-service/proto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
	repository.Repository
}

func (m *MockRepository) VerifyToken(ctx context.Context, userID string, token string) (bool, error) {
	args := m.Called(ctx, userID, token)
	return args.Bool(0), args.Error(1)
}

func TestAuthService_VerifyToken(t *testing.T) {
	tests := []struct {
		name          string
		userID        string
		token         string
		repoResponse  bool
		repoError     error
		expectedValid bool
		expectedError string
	}{
		{
			name:          "Valid token",
			userID:        "1",
			token:         "valid-token",
			repoResponse:  true,
			repoError:     nil,
			expectedValid: true,
			expectedError: "",
		},
		{
			name:          "Invalid token",
			userID:        "1",
			token:         "invalid-token",
			repoResponse:  false,
			repoError:     nil,
			expectedValid: false,
			expectedError: "",
		},
		{
			name:          "Repository error",
			userID:        "1",
			token:         "any-token",
			repoResponse:  false,
			repoError:     assert.AnError,
			expectedValid: false,
			expectedError: assert.AnError.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			mockRepo.On("VerifyToken", mock.Anything, tt.userID, tt.token).
				Return(tt.repoResponse, tt.repoError)

			service := NewAuthService(mockRepo)
			resp, err := service.VerifyToken(context.Background(), &proto.VerifyTokenRequest{
				UserId: tt.userID,
				Token:  tt.token,
			})

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedValid, resp.Valid)
			assert.Equal(t, tt.expectedError, resp.Error)
			mockRepo.AssertExpectations(t)
		})
	}
}
