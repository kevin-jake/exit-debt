package entities

import "errors"

// Domain errors
var (
	// User errors
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrPhoneNumberExists = errors.New("phone number already exists")
	ErrInvalidEmail      = errors.New("invalid email address")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrInvalidFirstName  = errors.New("first name is required")
	ErrInvalidLastName   = errors.New("last name is required")
	ErrInvalidCredentials = errors.New("invalid credentials")

	// Contact errors
	ErrContactNotFound     = errors.New("contact not found")
	ErrContactAlreadyExists = errors.New("contact already exists")
	ErrContactPhoneExists   = errors.New("contact with this phone number already exists")
	ErrInvalidContactName  = errors.New("contact name is required")

	// Debt errors
	ErrDebtListNotFound     = errors.New("debt list not found")
	ErrDebtItemNotFound     = errors.New("debt item not found")
	ErrInvalidDebtType      = errors.New("invalid debt type")
	ErrInvalidAmount        = errors.New("invalid amount")
	ErrInvalidCurrency      = errors.New("invalid currency")
	ErrInvalidPaymentMethod = errors.New("invalid payment method")
	ErrInvalidDueDate       = errors.New("due date must be in the future")
	ErrInvalidPaymentStatus = errors.New("invalid payment status")

	// Generic errors
	ErrInvalidInput       = errors.New("invalid input")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrInternalServer     = errors.New("internal server error")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token expired")
)
