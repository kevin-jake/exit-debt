package interfaces

import (
	"context"

	"github.com/google/uuid"

	"exit-debt/internal/domain/entities"
)

// ContactRepository defines the interface for contact data access operations
type ContactRepository interface {
	Create(ctx context.Context, contact *entities.Contact) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Contact, error)
	GetByEmail(ctx context.Context, email string) (*entities.Contact, error)
	GetUserContacts(ctx context.Context, userID uuid.UUID) ([]entities.Contact, error)
	Update(ctx context.Context, contact *entities.Contact) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetUserContactRelation(ctx context.Context, userID, contactID uuid.UUID) (*entities.UserContact, error)
	CreateUserContactRelation(ctx context.Context, userContact *entities.UserContact) error
	DeleteUserContactRelation(ctx context.Context, userID, contactID uuid.UUID) error
	ExistsByEmailForUser(ctx context.Context, userID uuid.UUID, email string) (bool, error)
	GetContactsWithEmail(ctx context.Context, email string) ([]entities.Contact, error)
}
