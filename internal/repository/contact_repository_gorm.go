package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"exit-debt/internal/domain/entities"
	"exit-debt/internal/domain/interfaces"
	"exit-debt/internal/models"
)

// contactRepositoryGORM implements the ContactRepository interface using GORM
type contactRepositoryGORM struct {
	db *gorm.DB
}

// NewContactRepositoryGORM creates a new contact repository with GORM
func NewContactRepositoryGORM(db *gorm.DB) interfaces.ContactRepository {
	return &contactRepositoryGORM{
		db: db,
	}
}

func (r *contactRepositoryGORM) Create(ctx context.Context, contact *entities.Contact) error {
	gormContact := r.entityToGORM(contact)
	if err := r.db.WithContext(ctx).Create(gormContact).Error; err != nil {
		return fmt.Errorf("failed to create contact: %w", err)
	}
	// Update the entity with the created ID and timestamps
	contact.ID = gormContact.ID
	contact.CreatedAt = gormContact.CreatedAt
	contact.UpdatedAt = gormContact.UpdatedAt
	return nil
}

func (r *contactRepositoryGORM) GetByID(ctx context.Context, id uuid.UUID) (*entities.Contact, error) {
	var gormContact models.Contact
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&gormContact).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrContactNotFound
		}
		return nil, fmt.Errorf("failed to get contact by ID: %w", err)
	}
	return r.gormToEntity(&gormContact), nil
}

func (r *contactRepositoryGORM) GetByEmail(ctx context.Context, email string) (*entities.Contact, error) {
	var gormContact models.Contact
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&gormContact).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrContactNotFound
		}
		return nil, fmt.Errorf("failed to get contact by email: %w", err)
	}
	return r.gormToEntity(&gormContact), nil
}

func (r *contactRepositoryGORM) GetUserContacts(ctx context.Context, userID uuid.UUID) ([]entities.Contact, error) {
	var userContacts []models.UserContact
	if err := r.db.WithContext(ctx).Joins("Contact").Where("user_contacts.user_id = ?", userID).Order("\"Contact\".\"name\" ASC").Find(&userContacts).Error; err != nil {
		return nil, fmt.Errorf("failed to get user contacts: %w", err)
	}

	contacts := make([]entities.Contact, len(userContacts))
	for i, uc := range userContacts {
		contacts[i] = *r.gormToEntity(&uc.Contact)
	}

	return contacts, nil
}

func (r *contactRepositoryGORM) Update(ctx context.Context, contact *entities.Contact) error {
	gormContact := r.entityToGORM(contact)
	if err := r.db.WithContext(ctx).Save(gormContact).Error; err != nil {
		return fmt.Errorf("failed to update contact: %w", err)
	}
	// Update the entity with the updated timestamp
	contact.UpdatedAt = gormContact.UpdatedAt
	return nil
}

func (r *contactRepositoryGORM) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.Contact{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete contact: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return entities.ErrContactNotFound
	}
	return nil
}

func (r *contactRepositoryGORM) GetUserContactRelation(ctx context.Context, userID, contactID uuid.UUID) (*entities.UserContact, error) {
	var gormUserContact models.UserContact
	if err := r.db.WithContext(ctx).Where("user_id = ? AND contact_id = ?", userID, contactID).First(&gormUserContact).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrContactNotFound
		}
		return nil, fmt.Errorf("failed to get user contact relation: %w", err)
	}
	return r.userContactGormToEntity(&gormUserContact), nil
}

func (r *contactRepositoryGORM) CreateUserContactRelation(ctx context.Context, userContact *entities.UserContact) error {
	gormUserContact := r.userContactEntityToGORM(userContact)
	if err := r.db.WithContext(ctx).Create(gormUserContact).Error; err != nil {
		return fmt.Errorf("failed to create user contact relation: %w", err)
	}
	// Update the entity with the created ID and timestamps
	userContact.ID = gormUserContact.ID
	userContact.CreatedAt = gormUserContact.CreatedAt
	userContact.UpdatedAt = gormUserContact.UpdatedAt
	return nil
}

func (r *contactRepositoryGORM) DeleteUserContactRelation(ctx context.Context, userID, contactID uuid.UUID) error {
	result := r.db.WithContext(ctx).Where("user_id = ? AND contact_id = ?", userID, contactID).Delete(&models.UserContact{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete user contact relation: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return entities.ErrContactNotFound
	}
	return nil
}

func (r *contactRepositoryGORM) ExistsByEmailForUser(ctx context.Context, userID uuid.UUID, email string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Joins("Contact").Where("user_contacts.user_id = ? AND \"Contact\".\"email\" = ?", userID, email).Model(&models.UserContact{}).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check if contact exists by email for user: %w", err)
	}
	return count > 0, nil
}

func (r *contactRepositoryGORM) GetContactsWithEmail(ctx context.Context, email string) ([]entities.Contact, error) {
	var gormContacts []models.Contact
	if err := r.db.WithContext(ctx).Where("email = ?", email).Find(&gormContacts).Error; err != nil {
		return nil, fmt.Errorf("failed to get contacts with email: %w", err)
	}

	contacts := make([]entities.Contact, len(gormContacts))
	for i, gc := range gormContacts {
		contacts[i] = *r.gormToEntity(&gc)
	}

	return contacts, nil
}

func (r *contactRepositoryGORM) GetUserContactRelationsByContactID(ctx context.Context, contactID uuid.UUID) ([]entities.UserContact, error) {
	var gormUserContacts []models.UserContact
	if err := r.db.WithContext(ctx).Where("contact_id = ?", contactID).Find(&gormUserContacts).Error; err != nil {
		return nil, fmt.Errorf("failed to get user contact relations by contact ID: %w", err)
	}

	userContacts := make([]entities.UserContact, len(gormUserContacts))
	for i, gormUserContact := range gormUserContacts {
		userContacts[i] = *r.userContactGormToEntity(&gormUserContact)
	}

	return userContacts, nil
}

// entityToGORM converts a domain entity to GORM model
func (r *contactRepositoryGORM) entityToGORM(contact *entities.Contact) *models.Contact {
	return &models.Contact{
		ID:         contact.ID,
		Name:       contact.Name,
		Email:      contact.Email,
		Phone:      contact.Phone,
		FacebookID: contact.FacebookID,
		Notes:      contact.Notes,
		IsUser:     contact.IsUser,
		UserIDRef:  contact.UserIDRef,
		CreatedAt:  contact.CreatedAt,
		UpdatedAt:  contact.UpdatedAt,
	}
}

// gormToEntity converts a GORM model to domain entity
func (r *contactRepositoryGORM) gormToEntity(gormContact *models.Contact) *entities.Contact {
	return &entities.Contact{
		ID:         gormContact.ID,
		Name:       gormContact.Name,
		Email:      gormContact.Email,
		Phone:      gormContact.Phone,
		FacebookID: gormContact.FacebookID,
		Notes:      gormContact.Notes,
		IsUser:     gormContact.IsUser,
		UserIDRef:  gormContact.UserIDRef,
		CreatedAt:  gormContact.CreatedAt,
		UpdatedAt:  gormContact.UpdatedAt,
	}
}

// userContactEntityToGORM converts a domain entity to GORM model
func (r *contactRepositoryGORM) userContactEntityToGORM(userContact *entities.UserContact) *models.UserContact {
	return &models.UserContact{
		ID:        userContact.ID,
		UserID:    userContact.UserID,
		ContactID: userContact.ContactID,
		CreatedAt: userContact.CreatedAt,
		UpdatedAt: userContact.UpdatedAt,
	}
}

// userContactGormToEntity converts a GORM model to domain entity
func (r *contactRepositoryGORM) userContactGormToEntity(gormUserContact *models.UserContact) *entities.UserContact {
	return &entities.UserContact{
		ID:        gormUserContact.ID,
		UserID:    gormUserContact.UserID,
		ContactID: gormUserContact.ContactID,
		CreatedAt: gormUserContact.CreatedAt,
		UpdatedAt: gormUserContact.UpdatedAt,
	}
}
