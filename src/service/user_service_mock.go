package service

import (
	"auth-service/src/domain"
	"context"

	"github.com/stretchr/testify/mock"
)

type UserServiceMock struct {
	mock.Mock
}

func (m *UserServiceMock) Register(ctx context.Context, name, email, password string) (*domain.User, error) {
	args := m.Called(ctx, name, email, password)
	if user, ok := args.Get(0).(*domain.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *UserServiceMock) Login(ctx context.Context, email, password string) (string, error) {
	args := m.Called(ctx, email, password)
	return args.String(0), args.Error(1)
}

func (m *UserServiceMock) GetProfile(ctx context.Context, userID string) (*domain.User, error) {
	args := m.Called(ctx, userID)
	if user, ok := args.Get(0).(*domain.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *UserServiceMock) ValidateToken(tokenString string) (map[string]interface{}, error) {
	args := m.Called(tokenString)
	if claims, ok := args.Get(0).(map[string]interface{}); ok {
		return claims, args.Error(1)
	}
	return nil, args.Error(1)
}
