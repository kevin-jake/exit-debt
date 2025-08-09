package entities

import (
	"time"

	"github.com/google/uuid"
)

// User represents the core user entity
type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	FirstName    string
	LastName     string
	Phone        *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// CreateUserRequest represents a request to create a new user
type CreateUserRequest struct {
	Email     string  `json:"email" validate:"required,email"`
	Password  string  `json:"password" validate:"required,min=6"`
	FirstName string  `json:"first_name" validate:"required"`
	LastName  string  `json:"last_name" validate:"required"`
	Phone     *string `json:"phone"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// RegisterResponse represents a registration response
type RegisterResponse struct {
	User User `json:"user"`
}

// FullName returns the user's full name
func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

// IsValid validates the user entity
func (u *User) IsValid() error {
	if u.Email == "" {
		return ErrInvalidEmail
	}
	if u.FirstName == "" {
		return ErrInvalidFirstName
	}
	if u.LastName == "" {
		return ErrInvalidLastName
	}
	if u.PasswordHash == "" {
		return ErrInvalidPassword
	}
	return nil
}
