package interfaces

import (
	"context"

	"github.com/google/uuid"

	"exit-debt/internal/domain/entities"
)

// AuthService defines the interface for authentication operations
type AuthService interface {
	Register(ctx context.Context, req *entities.CreateUserRequest) (*entities.RegisterResponse, error)
	Login(ctx context.Context, req *entities.LoginRequest) (*entities.LoginResponse, error)
	ValidateToken(ctx context.Context, tokenString string) (uuid.UUID, error)
	GenerateJWT(ctx context.Context, userID uuid.UUID) (string, error)
}
