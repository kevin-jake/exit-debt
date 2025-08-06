package models

import (
	"time"

	"github.com/google/uuid"
)

type Contact struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID      uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	Name        string    `json:"name" gorm:"not null;index"`
	Email       *string   `json:"email" gorm:"uniqueIndex"`
	Phone       *string   `json:"phone"`
	FacebookID  *string   `json:"facebook_id"`
	Notes       *string   `json:"notes"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	// Relationships
	User      User        `json:"user,omitempty" gorm:"foreignKey:UserID"`
	DebtLists []DebtList  `json:"debt_lists,omitempty" gorm:"foreignKey:ContactID"`
}

type CreateContactRequest struct {
	Name       string  `json:"name" binding:"required"`
	Email      *string `json:"email"`
	Phone      *string `json:"phone"`
	FacebookID *string `json:"facebook_id"`
	Notes      *string `json:"notes"`
}

type UpdateContactRequest struct {
	Name       *string `json:"name"`
	Email      *string `json:"email"`
	Phone      *string `json:"phone"`
	FacebookID *string `json:"facebook_id"`
	Notes      *string `json:"notes"`
} 