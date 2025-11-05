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

func (s *contactService) CreateContact(ctx context.Context, userID uuid.UUID, req *entities.CreateContactRequest) (*entities.ContactResponse, error) {
	// Validate input
	if err := s.validateCreateContactRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Check if contact with same email already exists for this user
	if req.Email != nil && *req.Email != "" {
		exists, err := s.contactRepo.ExistsByEmailForUser(ctx, userID, *req.Email)
		if err != nil {
			return nil, fmt.Errorf("failed to check if contact exists: %w", err)
		}
		if exists {
			return nil, entities.ErrContactAlreadyExists
		}
	}

	// Check if this contact is also a user (by email)
	var isUser bool
	var userIDRef *uuid.UUID
	if req.Email != nil && *req.Email != "" {
		user, err := s.userRepo.GetByEmail(ctx, *req.Email)
		if err == nil {
			isUser = true
			userIDRef = &user.ID
		} else if err != entities.ErrUserNotFound {
			return nil, fmt.Errorf("failed to check if email belongs to user: %w", err)
		}
	}

	// Create new contact (minimal identity)
	contact := &entities.Contact{
		ID:         uuid.New(),
		IsUser:     isUser,
		UserIDRef:  userIDRef,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.contactRepo.Create(ctx, contact); err != nil {
		return nil, fmt.Errorf("failed to create contact: %w", err)
	}

	// Create user-contact relationship with user-specific data
	userContact := &entities.UserContact{
		ID:        uuid.New(),
		UserID:    userID,
		ContactID: contact.ID,
		Name:      req.Name,
		Email:     req.Email,
		Phone:     req.Phone,
		Notes:     req.Notes,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Validate user contact entity
	if err := userContact.IsValid(); err != nil {
		return nil, fmt.Errorf("invalid contact entity: %w", err)
	}

	if err := s.contactRepo.CreateUserContactRelation(ctx, userContact); err != nil {
		return nil, fmt.Errorf("failed to create user contact relation: %w", err)
	}

	// Create reciprocal contact if this contact is also a user
	if req.Email != nil && *req.Email != "" && isUser {
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

	// Build and return ContactResponse
	return &entities.ContactResponse{
		ID:        contact.ID,
		Name:      userContact.Name,
		Email:     userContact.Email,
		Phone:     userContact.Phone,
		Notes:     userContact.Notes,
		IsUser:    contact.IsUser,
		UserIDRef: contact.UserIDRef,
		CreatedAt: contact.CreatedAt,
		UpdatedAt: contact.UpdatedAt,
	}, nil
}

func (s *contactService) GetContact(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*entities.ContactResponse, error) {
	// Get user contact relation which contains user-specific data
	userContact, err := s.contactRepo.GetUserContactRelation(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("failed to verify contact access: %w", err)
	}

	// Get the contact entity for IsUser and UserIDRef
	contact, err := s.contactRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get contact: %w", err)
	}

	// Build and return ContactResponse combining both
	return &entities.ContactResponse{
		ID:        contact.ID,
		Name:      userContact.Name,
		Email:     userContact.Email,
		Phone:     userContact.Phone,
		Notes:     userContact.Notes,
		IsUser:    contact.IsUser,
		UserIDRef: contact.UserIDRef,
		CreatedAt: userContact.CreatedAt,
		UpdatedAt: userContact.UpdatedAt,
	}, nil
}

func (s *contactService) GetUserContacts(ctx context.Context, userID uuid.UUID) ([]entities.ContactResponse, error) {
	// Get user contacts with user-specific data
	userContacts, err := s.contactRepo.GetUserContacts(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user contacts: %w", err)
	}

	// Build ContactResponse for each user contact
	responses := make([]entities.ContactResponse, 0, len(userContacts))
	for _, uc := range userContacts {
		// Get the contact to retrieve IsUser and UserIDRef
		contact, err := s.contactRepo.GetByID(ctx, uc.ContactID)
		if err != nil {
			// Log but continue
			continue
		}
		
		responses = append(responses, entities.ContactResponse{
			ID:        contact.ID,
			Name:      uc.Name,
			Email:     uc.Email,
			Phone:     uc.Phone,
			Notes:     uc.Notes,
			IsUser:    contact.IsUser,
			UserIDRef: contact.UserIDRef,
			CreatedAt: uc.CreatedAt,
			UpdatedAt: uc.UpdatedAt,
		})
	}

	return responses, nil
}

func (s *contactService) UpdateContact(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *entities.UpdateContactRequest) (*entities.ContactResponse, error) {
	// Validate input
	if err := s.validateUpdateContactRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Get the user contact relation (user-specific data)
	userContact, err := s.contactRepo.GetUserContactRelation(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("failed to verify contact access: %w", err)
	}

	// Get the contact entity
	contact, err := s.contactRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get contact: %w", err)
	}

	// Update user-specific fields in UserContact
	if req.Name != nil {
		userContact.Name = *req.Name
	}
	if req.Email != nil {
		userContact.Email = req.Email
	}
	if req.Phone != nil {
		userContact.Phone = req.Phone
	}
	if req.Notes != nil {
		userContact.Notes = req.Notes
	}

	userContact.UpdatedAt = time.Now()

	// Validate updated user contact entity
	if err := userContact.IsValid(); err != nil {
		return nil, fmt.Errorf("invalid updated contact entity: %w", err)
	}

	// Check if email changed and if new email belongs to a user
	if req.Email != nil && *req.Email != "" {
		user, err := s.userRepo.GetByEmail(ctx, *req.Email)
		if err == nil {
			// Email belongs to a user - update Contact entity
			contact.IsUser = true
			contact.UserIDRef = &user.ID
			contact.UpdatedAt = time.Now()
			if err := s.contactRepo.Update(ctx, contact); err != nil {
				return nil, fmt.Errorf("failed to update contact: %w", err)
			}
		} else if err == entities.ErrUserNotFound {
			// Email doesn't belong to a user
			if contact.IsUser {
				// Was a user before, now isn't
				contact.IsUser = false
				contact.UserIDRef = nil
				contact.UpdatedAt = time.Now()
				if err := s.contactRepo.Update(ctx, contact); err != nil {
					return nil, fmt.Errorf("failed to update contact: %w", err)
				}
			}
		} else {
			return nil, fmt.Errorf("failed to check if email belongs to user: %w", err)
		}
	}

	// Update the user contact relation
	if err := s.contactRepo.UpdateUserContactRelation(ctx, userContact); err != nil {
		return nil, fmt.Errorf("failed to update user contact relation: %w", err)
	}

	// Build and return ContactResponse
	return &entities.ContactResponse{
		ID:        contact.ID,
		Name:      userContact.Name,
		Email:     userContact.Email,
		Phone:     userContact.Phone,
		Notes:     userContact.Notes,
		IsUser:    contact.IsUser,
		UserIDRef: contact.UserIDRef,
		CreatedAt: userContact.CreatedAt,
		UpdatedAt: userContact.UpdatedAt,
	}, nil
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
	// Find all UserContact entries that have this user's email
	userContactsWithEmail, err := s.contactRepo.GetUserContactsByEmail(ctx, userEmail)
	if err != nil {
		return fmt.Errorf("failed to get user contacts with email: %w", err)
	}

	if len(userContactsWithEmail) == 0 {
		return nil // No existing user contacts found with this email
	}

	// For each existing UserContact with matching email, create reciprocal contact
	for _, uc := range userContactsWithEmail {
		// Skip if this is the new user's own contact (shouldn't happen but just in case)
		if uc.UserID == userID {
			continue
		}

		// Update the Contact to reference the new user
		contact, err := s.contactRepo.GetByID(ctx, uc.ContactID)
		if err != nil {
			// Log the error but continue
			continue
		}
		contact.IsUser = true
		contact.UserIDRef = &userID
		contact.UpdatedAt = time.Now()
		
		if err := s.contactRepo.Update(ctx, contact); err != nil {
			// Log the error but continue
			continue
		}

		// Get the contact owner (the user who added this contact)
		contactOwner, err := s.userRepo.GetByID(ctx, uc.UserID)
		if err != nil {
			// Log the error but continue
			continue
		}

		// Create a new Contact entity for the contact owner
		contactForOwner := &entities.Contact{
			ID:        uuid.New(),
			IsUser:    true,
			UserIDRef: &contactOwner.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := s.contactRepo.Create(ctx, contactForOwner); err != nil {
			// Log the error but continue
			continue
		}

		// Create UserContact for new user -> contact owner
		reciprocalUserContact := &entities.UserContact{
			ID:        uuid.New(),
			UserID:    userID,
			ContactID: contactForOwner.ID,
			Name:      contactOwner.FirstName + " " + contactOwner.LastName,
			Email:     &contactOwner.Email,
			Phone:     contactOwner.Phone,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := s.contactRepo.CreateUserContactRelation(ctx, reciprocalUserContact); err != nil {
			// Log the error but continue
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

	// Check if the reciprocal contact already exists (by checking UserContact with contact owner's email)
	if contactOwner.Email != "" {
		exists, err := s.contactRepo.ExistsByEmailForUser(ctx, existingUser.ID, contactOwner.Email)
		if err != nil {
			return fmt.Errorf("failed to check if reciprocal contact exists: %w", err)
		}
		if exists {
			return nil // Reciprocal contact already exists
		}
	}

	// Create new Contact entity for the contact owner
	contact := &entities.Contact{
		ID:        uuid.New(),
		IsUser:    true,
		UserIDRef: &contactOwner.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.contactRepo.Create(ctx, contact); err != nil {
		return fmt.Errorf("failed to create reciprocal contact: %w", err)
	}

	// Create user-contact relationship for the existing user with user-specific data
	userContact := &entities.UserContact{
		ID:        uuid.New(),
		UserID:    existingUser.ID,
		ContactID: contact.ID,
		Name:      contactOwner.FirstName + " " + contactOwner.LastName,
		Email:     &contactOwner.Email,
		Phone:     contactOwner.Phone,
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
