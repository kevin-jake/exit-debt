package entities

import (
	"time"

	"github.com/google/uuid"
)

// Contact represents the core contact entity
type Contact struct {
	ID         uuid.UUID
	Name       string
	Email      *string
	Phone      *string
	FacebookID *string
	Notes      *string
	IsUser     bool
	UserIDRef  *uuid.UUID
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// UserContact represents the many-to-many relationship between users and contacts
type UserContact struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	ContactID uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CreateContactRequest represents a request to create a new contact
type CreateContactRequest struct {
	Name       string  `json:"name" validate:"required"`
	Email      *string `json:"email" validate:"omitempty,email"`
	Phone      *string `json:"phone"`
	FacebookID *string `json:"facebook_id"`
	Notes      *string `json:"notes"`
}

// UpdateContactRequest represents a request to update a contact
type UpdateContactRequest struct {
	Name       *string `json:"name"`
	Email      *string `json:"email" validate:"omitempty,email"`
	Phone      *string `json:"phone"`
	FacebookID *string `json:"facebook_id"`
	Notes      *string `json:"notes"`
}

// IsValid validates the contact entity
func (c *Contact) IsValid() error {
	if c.Name == "" {
		return ErrInvalidContactName
	}
	return nil
}

// HasEmail returns true if the contact has an email address
func (c *Contact) HasEmail() bool {
	return c.Email != nil && *c.Email != ""
}

// HasPhone returns true if the contact has a phone number
func (c *Contact) HasPhone() bool {
	return c.Phone != nil && *c.Phone != ""
}
