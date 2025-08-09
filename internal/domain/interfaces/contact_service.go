package interfaces

import (
	"context"

	"github.com/google/uuid"

	"exit-debt/internal/domain/entities"
)

// ContactService defines the interface for contact operations
type ContactService interface {
	CreateContact(ctx context.Context, userID uuid.UUID, req *entities.CreateContactRequest) (*entities.Contact, error)
	GetContact(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*entities.Contact, error)
	GetUserContacts(ctx context.Context, userID uuid.UUID) ([]entities.Contact, error)
	UpdateContact(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *entities.UpdateContactRequest) (*entities.Contact, error)
	DeleteContact(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	CreateContactsForNewUser(ctx context.Context, userID uuid.UUID, userEmail string) error
	CreateReciprocalContact(ctx context.Context, contactEmail string, contactOwnerID uuid.UUID) error
}
