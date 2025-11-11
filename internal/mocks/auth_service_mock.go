package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"pay-your-dues/internal/domain/entities"
)

// MockAuthService is a mock implementation of AuthService
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(ctx context.Context, req *entities.CreateUserRequest) (*entities.RegisterResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.RegisterResponse), args.Error(1)
}

func (m *MockAuthService) Login(ctx context.Context, req *entities.LoginRequest) (*entities.LoginResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.LoginResponse), args.Error(1)
}

func (m *MockAuthService) ValidateToken(ctx context.Context, token string) (uuid.UUID, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockAuthService) GenerateJWT(ctx context.Context, userID uuid.UUID) (string, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Error(1)
}
