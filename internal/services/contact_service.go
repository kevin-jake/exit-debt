package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"exit-debt/internal/domain/entities"
	"exit-debt/internal/domain/interfaces"
)

// contactService implements the ContactService interface
type contactService struct {
	contactRepo interfaces.ContactRepository
	userRepo    interfaces.UserRepository
}

// NewContactService creates a new contact service
func NewContactService(contactRepo interfaces.ContactRepository, userRepo interfaces.UserRepository) interfaces.ContactService {
	return &contactService{
		contactRepo: contactRepo,
		userRepo:    userRepo,
	}
}

func (s *contactService) CreateContact(ctx context.Context, userID uuid.UUID, req *entities.CreateContactRequest) (*entities.Contact, error) {
	// Validate input
	if err := s.validateCreateContactRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Check if contact with same email already exists for this user
	if req.Email != nil {
		exists, err := s.contactRepo.ExistsByEmailForUser(ctx, userID, *req.Email)
		if err != nil {
			return nil, fmt.Errorf("failed to check if contact exists: %w", err)
		}
		if exists {
			return nil, entities.ErrContactAlreadyExists
		}
	}

	// Check if contact with same phone number already exists for this user
	if req.Phone != nil {
		exists, err := s.contactRepo.ExistsByPhoneForUser(ctx, userID, *req.Phone)
		if err != nil {
			return nil, fmt.Errorf("failed to check if contact exists: %w", err)
		}
		if exists {
			return nil, entities.ErrContactPhoneExists
		}
	}

	// Check if this contact is also a user (by email)
	var isUser bool
	var userIDRef *uuid.UUID
	if req.Email != nil {
		user, err := s.userRepo.GetByEmail(ctx, *req.Email)
		if err == nil {
			isUser = true
			userIDRef = &user.ID
		} else if err != entities.ErrUserNotFound {
			return nil, fmt.Errorf("failed to check if email belongs to user: %w", err)
		}
	}

	// Try to find existing contact by email first
	var contact *entities.Contact
	if req.Email != nil {
		existingContact, err := s.contactRepo.GetByEmail(ctx, *req.Email)
		if err == nil {
			contact = existingContact
		} else if err != entities.ErrContactNotFound {
			return nil, fmt.Errorf("failed to get existing contact: %w", err)
		}
	}

	// If no contact found by email, try to find by phone
	if contact == nil && req.Phone != nil {
		existingContact, err := s.contactRepo.GetByPhone(ctx, *req.Phone)
		if err == nil {
			contact = existingContact
		} else if err != entities.ErrContactNotFound {
			return nil, fmt.Errorf("failed to get existing contact: %w", err)
		}
	}

	if contact == nil {
		// Create new contact
		contact = &entities.Contact{
			ID:         uuid.New(),
			Name:       req.Name,
			Email:      req.Email,
			Phone:      req.Phone,
			Notes:      req.Notes,
			IsUser:     isUser,
			UserIDRef:  userIDRef,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		// Validate contact entity
		if err := contact.IsValid(); err != nil {
			return nil, fmt.Errorf("invalid contact entity: %w", err)
		}

		if err := s.contactRepo.Create(ctx, contact); err != nil {
			return nil, fmt.Errorf("failed to create contact: %w", err)
		}
	} else {
		// Update existing contact with new information if provided
		updated := false
		if req.Name != "" && contact.Name != req.Name {
			contact.Name = req.Name
			updated = true
		}
		if req.Phone != nil && (contact.Phone == nil || *contact.Phone != *req.Phone) {
			contact.Phone = req.Phone
			updated = true
		}
		if req.Notes != nil && (contact.Notes == nil || *contact.Notes != *req.Notes) {
			contact.Notes = req.Notes
			updated = true
		}
		if isUser && !contact.IsUser {
			contact.IsUser = true
			contact.UserIDRef = userIDRef
			updated = true
		}
		
		if updated {
			contact.UpdatedAt = time.Now()
			if err := s.contactRepo.Update(ctx, contact); err != nil {
				return nil, fmt.Errorf("failed to update existing contact: %w", err)
			}
		}
	}

	// Create user-contact relationship
	userContact := &entities.UserContact{
		ID:        uuid.New(),
		UserID:    userID,
		ContactID: contact.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.contactRepo.CreateUserContactRelation(ctx, userContact); err != nil {
		return nil, fmt.Errorf("failed to create user contact relation: %w", err)
	}

	// Create reciprocal contact if this contact is also a user
	if req.Email != nil && isUser {
		if err := s.CreateReciprocalContact(ctx, *req.Email, userID); err != nil {
			// Log the error but don't fail contact creation
			logger := zerolog.Ctx(ctx)
			logger.Warn().
				Err(err).
				Str("user_id", userID.String()).
				Str("contact_email", *req.Email).
				Str("contact_id", contact.ID.String()).
				Msg("Failed to create reciprocal contact")
		}
	}

	return contact, nil
}

func (s *contactService) GetContact(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*entities.Contact, error) {
	// Check if user has access to this contact
	_, err := s.contactRepo.GetUserContactRelation(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("failed to verify contact access: %w", err)
	}

	contact, err := s.contactRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get contact: %w", err)
	}

	return contact, nil
}

func (s *contactService) GetUserContacts(ctx context.Context, userID uuid.UUID) ([]entities.Contact, error) {
	contacts, err := s.contactRepo.GetUserContacts(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user contacts: %w", err)
	}

	return contacts, nil
}

func (s *contactService) UpdateContact(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *entities.UpdateContactRequest) (*entities.Contact, error) {
	// Validate input
	if err := s.validateUpdateContactRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Check if user has access to this contact
	_, err := s.contactRepo.GetUserContactRelation(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("failed to verify contact access: %w", err)
	}

	// Get the existing contact
	contact, err := s.contactRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get contact: %w", err)
	}

	// Update fields if provided
	if req.Name != nil {
		contact.Name = *req.Name
	}
	if req.Email != nil {
		contact.Email = req.Email
		// Check if this contact is also a user (by email)
		if req.Email != nil {
			user, err := s.userRepo.GetByEmail(ctx, *req.Email)
			if err == nil {
				contact.IsUser = true
				contact.UserIDRef = &user.ID
			} else if err == entities.ErrUserNotFound {
				contact.IsUser = false
				contact.UserIDRef = nil
			} else {
				return nil, fmt.Errorf("failed to check if email belongs to user: %w", err)
			}
		}
	}
	if req.Phone != nil {
		contact.Phone = req.Phone
	}
	if req.Notes != nil {
		contact.Notes = req.Notes
	}

	contact.UpdatedAt = time.Now()

	// Validate updated contact entity
	if err := contact.IsValid(); err != nil {
		return nil, fmt.Errorf("invalid updated contact entity: %w", err)
	}

	if err := s.contactRepo.Update(ctx, contact); err != nil {
		return nil, fmt.Errorf("failed to update contact: %w", err)
	}

	return contact, nil
}

func (s *contactService) DeleteContact(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	// Check if user has access to this contact
	_, err := s.contactRepo.GetUserContactRelation(ctx, userID, id)
	if err != nil {
		return fmt.Errorf("failed to verify contact access: %w", err)
	}

	// Delete the user-contact relationship
	if err := s.contactRepo.DeleteUserContactRelation(ctx, userID, id); err != nil {
		return fmt.Errorf("failed to delete user contact relation: %w", err)
	}

	return nil
}

func (s *contactService) CreateContactsForNewUser(ctx context.Context, userID uuid.UUID, userEmail string) error {
	// Find all contacts that have this user's email
	existingContacts, err := s.contactRepo.GetContactsWithEmail(ctx, userEmail)
	if err != nil {
		return fmt.Errorf("failed to get contacts with email: %w", err)
	}

	if len(existingContacts) == 0 {
		return nil // No existing contacts found with this email
	}

	// For each existing contact, find all users who have this contact and create reciprocal contacts
	for _, contact := range existingContacts {
		// Get all user-contact relationships for this contact
		userContactRelations, err := s.contactRepo.GetUserContactRelationsByContactID(ctx, contact.ID)
		if err != nil {
			// Log the error but continue with other contacts
			continue
		}

		// For each user who has this contact, create a reciprocal contact
		for _, userContactRel := range userContactRelations {
			// Skip if this is the new user (we don't want to create a contact for ourselves)
			if userContactRel.UserID == userID {
				continue
			}

			// Get the user details who owns this contact (from user_contacts.user_id)
			contactOwner, err := s.userRepo.GetByID(ctx, userContactRel.UserID)
			if err != nil {
				// Log the error but continue with other users
				continue
			}

			// Find or create a contact for the contact owner in the contacts table
			var contactForOwner *entities.Contact
			if contactOwner.Email != "" {
				existingContact, err := s.contactRepo.GetByEmail(ctx, contactOwner.Email)
				if err == nil {
					contactForOwner = existingContact
				} else if err == entities.ErrContactNotFound {
					// Create new contact for the contact owner
					contactForOwner = &entities.Contact{
						ID:         uuid.New(),
						Name:       contactOwner.FullName(),
						Email:      &contactOwner.Email,
						Phone:      contactOwner.Phone,
						IsUser:     true,
						UserIDRef:  &contactOwner.ID,
						CreatedAt:  time.Now(),
						UpdatedAt:  time.Now(),
					}

					if err := s.contactRepo.Create(ctx, contactForOwner); err != nil {
						// Log the error but continue with other users
						continue
					}
				} else {
					// Log the error but continue with other users
					continue
				}
			}

			// Create user-contact relationship for the new user with the contact owner's contact
			if contactForOwner != nil {
				newUserContact := &entities.UserContact{
					ID:        uuid.New(),
					UserID:    userID,
					ContactID: contactForOwner.ID, // Use the contact ID from the contacts table
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}

				if err := s.contactRepo.CreateUserContactRelation(ctx, newUserContact); err != nil {
					// Log the error but continue with other users
					continue
				}
			}
		}

		// Update the existing contact to reference the new user
		contact.IsUser = true
		contact.UserIDRef = &userID
		contact.UpdatedAt = time.Now()
		
		if err := s.contactRepo.Update(ctx, &contact); err != nil {
			// Log the error but continue with other contacts
			continue
		}
	}

	return nil
}

func (s *contactService) CreateReciprocalContact(ctx context.Context, contactEmail string, contactOwnerID uuid.UUID) error {
	// Find the user who owns this email
	existingUser, err := s.userRepo.GetByEmail(ctx, contactEmail)
	if err != nil {
		if err == entities.ErrUserNotFound {
			return nil // User doesn't exist, no reciprocal contact needed
		}
		return fmt.Errorf("failed to get user by email: %w", err)
	}

	// Get the contact owner's details
	contactOwner, err := s.userRepo.GetByID(ctx, contactOwnerID)
	if err != nil {
		return fmt.Errorf("failed to get contact owner: %w", err)
	}

	// Check if the reciprocal contact already exists
	if contactOwner.Email != "" {
		exists, err := s.contactRepo.ExistsByEmailForUser(ctx, contactOwner.ID, contactOwner.Email)
		if err != nil {
			return fmt.Errorf("failed to check if reciprocal contact exists: %w", err)
		}
		if exists {
			return nil // Reciprocal contact already exists
		}
	}

	// Find or create contact for the contact owner
	var contact *entities.Contact
	if contactOwner.Email != "" {
		existingContact, err := s.contactRepo.GetByEmail(ctx, contactOwner.Email)
		if err == nil {
			contact = existingContact
		} else if err != entities.ErrContactNotFound {
			return fmt.Errorf("failed to get existing contact for owner: %w", err)
		}
	}

	if contact == nil {
		// Create new contact for the contact owner
		contact = &entities.Contact{
			ID:         uuid.New(),
			Name:       contactOwner.FullName(),
			Email:      &contactOwner.Email,
			Phone:      contactOwner.Phone,
			IsUser:     true,
			UserIDRef:  &contactOwner.ID,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		if err := s.contactRepo.Create(ctx, contact); err != nil {
			return fmt.Errorf("failed to create reciprocal contact: %w", err)
		}
	}

	// Create user-contact relationship for the existing user
	userContact := &entities.UserContact{
		ID:        uuid.New(),
		UserID:    existingUser.ID,
		ContactID: contact.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.contactRepo.CreateUserContactRelation(ctx, userContact); err != nil {
		return fmt.Errorf("failed to create reciprocal user contact relation: %w", err)
	}

	return nil
}

func (s *contactService) validateCreateContactRequest(req *entities.CreateContactRequest) error {
	if req.Name == "" {
		return entities.ErrInvalidContactName
	}
	return nil
}

func (s *contactService) validateUpdateContactRequest(req *entities.UpdateContactRequest) error {
	if req.Name != nil && *req.Name == "" {
		return entities.ErrInvalidContactName
	}
	return nil
}
