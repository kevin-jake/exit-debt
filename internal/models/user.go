package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email        string    `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash string    `json:"-" gorm:"column:password_hash;not null"`
	FirstName    string    `json:"first_name" gorm:"not null"`
	LastName     string    `json:"last_name" gorm:"not null"`
	Phone        *string   `json:"phone"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	
	// Relationships
	Contacts   []Contact   `json:"contacts,omitempty" gorm:"foreignKey:UserID"`
	DebtLists []DebtList  `json:"debt_lists,omitempty" gorm:"foreignKey:UserID"`
}

type CreateUserRequest struct {
	Email     string  `json:"email" binding:"required,email"`
	Password  string  `json:"password" binding:"required,min=6"`
	FirstName string  `json:"first_name" binding:"required"`
	LastName  string  `json:"last_name" binding:"required"`
	Phone     *string `json:"phone"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type RegisterResponse struct {
	User    User    `json:"user"`
	SharingSummary *SharingSummary `json:"sharing_summary,omitempty"`
}

type SharingSummary struct {
	ContactsFound     int    `json:"contacts_found"`
	DebtListsShared   int    `json:"debt_lists_shared"`
	TotalAmountShared string `json:"total_amount_shared"`
} 