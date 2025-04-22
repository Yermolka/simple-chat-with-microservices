package repository

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(username, password string) (int64, error) {
	args := m.Called(username, password)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRepository) GetAll() ([]User, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]User), args.Error(1)
}

func (m *MockRepository) GetById(id int64) (User, error) {
	args := m.Called(id)
	return args.Get(0).(User), args.Error(1)
}

func (m *MockRepository) AuthenticateUser(username, password string) (*User, error) {
	args := m.Called(username, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockRepository) CreateToken(userID int64, token string) error {
	args := m.Called(userID, token)
	return args.Error(0)
}

func (m *MockRepository) DeleteToken(token string) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockRepository) VerifyToken(ctx context.Context, userID string, token string) (bool, error) {
	args := m.Called(ctx, userID, token)
	return args.Bool(0), args.Error(1)
}
