package models

import (
	"time"

	"github.com/google/uuid"
)

type Contact struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	IsUser      bool      `json:"is_user" gorm:"default:false;index"`
	UserIDRef   *uuid.UUID `json:"user_id_ref" gorm:"type:uuid;index"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	// Relationships
	Users      []User      `json:"users,omitempty" gorm:"many2many:user_contacts;"`
	UserRef    *User       `json:"user_ref,omitempty" gorm:"foreignKey:UserIDRef"`
	DebtLists  []DebtList  `json:"debt_lists,omitempty" gorm:"foreignKey:ContactID"`
}

// UserContact represents the many-to-many relationship between users and contacts
// with user-specific contact information
type UserContact struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	ContactID uuid.UUID `json:"contact_id" gorm:"type:uuid;not null;index"`
	Name      string    `json:"name" gorm:"not null;index"`
	Email     *string   `json:"email" gorm:"index"`
	Phone     *string   `json:"phone"`
	Notes     *string   `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	// Relationships
	User    User    `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Contact Contact `json:"contact,omitempty" gorm:"foreignKey:ContactID"`
}

type CreateContactRequest struct {
	Name  string  `json:"name" binding:"required"`
	Email *string `json:"email"`
	Phone *string `json:"phone"`
	Notes *string `json:"notes"`
}

type UpdateContactRequest struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
	Phone *string `json:"phone"`
	Notes *string `json:"notes"`
} 