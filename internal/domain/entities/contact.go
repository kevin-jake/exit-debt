package entities

import (
	"time"

	"github.com/google/uuid"
)

// Contact represents the core contact entity (minimal identity)
type Contact struct {
	ID         uuid.UUID
	IsUser     bool
	UserIDRef  *uuid.UUID
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// UserContact represents the many-to-many relationship between users and contacts
// with user-specific contact information
type UserContact struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	ContactID uuid.UUID
	Name      string
	Email     *string
	Phone     *string
	Notes     *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CreateContactRequest represents a request to create a new contact
type CreateContactRequest struct {
	Name  string  `json:"name" validate:"required"`
	Email *string `json:"email" validate:"omitempty,email"`
	Phone *string `json:"phone"`
	Notes *string `json:"notes"`
}

// UpdateContactRequest represents a request to update a contact
type UpdateContactRequest struct {
	Name  *string `json:"name"`
	Email *string `json:"email" validate:"omitempty,email"`
	Phone *string `json:"phone"`
	Notes *string `json:"notes"`
}

// ContactResponse represents a contact with user-specific information
// This combines Contact entity (identity, IsUser) with UserContact (user-specific data)
type ContactResponse struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Email     *string    `json:"email"`
	Phone     *string    `json:"phone"`
	Notes     *string    `json:"notes"`
	IsUser    bool       `json:"is_user"`
	UserIDRef *uuid.UUID `json:"user_id_ref,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// IsValid validates the user contact entity
func (uc *UserContact) IsValid() error {
	if uc.Name == "" {
		return ErrInvalidContactName
	}
	return nil
}

// HasEmail returns true if the user contact has an email address
func (uc *UserContact) HasEmail() bool {
	return uc.Email != nil && *uc.Email != ""
}

// HasPhone returns true if the user contact has a phone number
func (uc *UserContact) HasPhone() bool {
	return uc.Phone != nil && *uc.Phone != ""
}
