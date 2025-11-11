package interfaces

import (
	"context"

	"github.com/google/uuid"

	"pay-your-dues/internal/domain/entities"
)

// ContactRepository defines the interface for contact data access operations
type ContactRepository interface {
	Create(ctx context.Context, contact *entities.Contact) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Contact, error)
	Update(ctx context.Context, contact *entities.Contact) error
	Delete(ctx context.Context, id uuid.UUID) error
	
	// UserContact operations
	GetUserContactRelation(ctx context.Context, userID, contactID uuid.UUID) (*entities.UserContact, error)
	CreateUserContactRelation(ctx context.Context, userContact *entities.UserContact) error
	UpdateUserContactRelation(ctx context.Context, userContact *entities.UserContact) error
	DeleteUserContactRelation(ctx context.Context, userID, contactID uuid.UUID) error
	GetUserContacts(ctx context.Context, userID uuid.UUID) ([]entities.UserContact, error)
	GetUserContactRelationsByContactID(ctx context.Context, contactID uuid.UUID) ([]entities.UserContact, error)
	GetUserContactsByEmail(ctx context.Context, email string) ([]entities.UserContact, error)
	ExistsByEmailForUser(ctx context.Context, userID uuid.UUID, email string) (bool, error)
}
