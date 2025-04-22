package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"user-service/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]string
		expectedStatus int
		mockSetup      func(*repository.MockRepository)
	}{
		{
			name: "Successful user creation",
			requestBody: map[string]string{
				"username": "testuser",
				"password": "testpass",
			},
			expectedStatus: http.StatusOK,
			mockSetup: func(m *repository.MockRepository) {
				m.On("Create", "testuser", "testpass").Return(int64(1), nil)
			},
		},
		{
			name: "Missing username",
			requestBody: map[string]string{
				"password": "testpass",
			},
			expectedStatus: http.StatusBadRequest,
			mockSetup:      func(m *repository.MockRepository) {},
		},
		{
			name: "Repository error",
			requestBody: map[string]string{
				"username": "testuser",
				"password": "testpass",
			},
			expectedStatus: http.StatusInternalServerError,
			mockSetup: func(m *repository.MockRepository) {
				m.On("Create", "testuser", "testpass").Return(int64(0), assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(repository.MockRepository)
			tt.mockSetup(mockRepo)

			handler := NewHandler(mockRepo)
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.CreateUser(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestLogin(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]string
		expectedStatus int
		mockSetup      func(*repository.MockRepository)
	}{
		{
			name: "Successful login",
			requestBody: map[string]string{
				"username": "testuser",
				"password": "testpass",
			},
			expectedStatus: http.StatusOK,
			mockSetup: func(m *repository.MockRepository) {
				user := &repository.User{Id: 1}
				m.On("AuthenticateUser", "testuser", "testpass").Return(user, nil)
				m.On("CreateToken", int64(1), mock.Anything).Return(nil)
			},
		},
		{
			name: "Invalid credentials",
			requestBody: map[string]string{
				"username": "testuser",
				"password": "wrongpass",
			},
			expectedStatus: http.StatusUnauthorized,
			mockSetup: func(m *repository.MockRepository) {
				m.On("AuthenticateUser", "testuser", "wrongpass").Return(nil, assert.AnError)
			},
		},
		{
			name: "Token creation error",
			requestBody: map[string]string{
				"username": "testuser",
				"password": "testpass",
			},
			expectedStatus: http.StatusInternalServerError,
			mockSetup: func(m *repository.MockRepository) {
				user := &repository.User{Id: 1}
				m.On("AuthenticateUser", "testuser", "testpass").Return(user, nil)
				m.On("CreateToken", int64(1), mock.Anything).Return(assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(repository.MockRepository)
			tt.mockSetup(mockRepo)

			handler := NewHandler(mockRepo)
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.Login(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedStatus == http.StatusOK {
				var response TokenResponse
				err := json.NewDecoder(w.Body).Decode(&response)
				assert.NoError(t, err)
				assert.NotEmpty(t, response.Token)
				assert.True(t, response.ExpiresAt.After(time.Now()))
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestLogout(t *testing.T) {
	tests := []struct {
		name           string
		token          string
		expectedStatus int
		mockSetup      func(*repository.MockRepository)
	}{
		{
			name:           "Successful logout",
			token:          "valid-token",
			expectedStatus: http.StatusOK,
			mockSetup: func(m *repository.MockRepository) {
				m.On("DeleteToken", "valid-token").Return(nil)
			},
		},
		{
			name:           "Missing token",
			token:          "",
			expectedStatus: http.StatusUnauthorized,
			mockSetup:      func(m *repository.MockRepository) {},
		},
		{
			name:           "Token deletion error",
			token:          "valid-token",
			expectedStatus: http.StatusInternalServerError,
			mockSetup: func(m *repository.MockRepository) {
				m.On("DeleteToken", "valid-token").Return(assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(repository.MockRepository)
			tt.mockSetup(mockRepo)

			handler := NewHandler(mockRepo)
			req := httptest.NewRequest("POST", "/logout", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", tt.token)
			}
			w := httptest.NewRecorder()

			handler.Logout(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockRepo.AssertExpectations(t)
		})
	}
}
