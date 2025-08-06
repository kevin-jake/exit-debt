package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"exit-debt/internal/models"
)

type ContactServiceGORM struct {
	db *gorm.DB
}

func NewContactServiceGORM(db *gorm.DB) *ContactServiceGORM {
	return &ContactServiceGORM{db: db}
}

func (s *ContactServiceGORM) CreateContact(userID uuid.UUID, req *models.CreateContactRequest) (*models.Contact, error) {
	// Check if contact with same name already exists for this user
	var existingUserContact models.UserContact
	if err := s.db.Joins("Contact").Where("user_contacts.user_id = ? AND \"Contact\".\"name\" = ?", userID, req.Name).First(&existingUserContact).Error; err == nil {
		return nil, errors.New("contact with this name already exists")
	}

	// Check if this contact is also a user (by email)
	var isUser bool
	var userIDRef *uuid.UUID
	if req.Email != nil {
		var user models.User
		if err := s.db.Where("email = ?", *req.Email).First(&user).Error; err == nil {
			isUser = true
			userIDRef = &user.ID
		}
	}

	// Create or find existing contact
	var contact *models.Contact
	if req.Email != nil {
		// Try to find existing contact with same email
		var existingContact models.Contact
		if err := s.db.Where("email = ?", *req.Email).First(&existingContact).Error; err == nil {
			contact = &existingContact
		}
	}

	if contact == nil {
		// Create new contact
		contact = &models.Contact{
			ID:         uuid.New(),
			Name:       req.Name,
			Email:      req.Email,
			Phone:      req.Phone,
			FacebookID: req.FacebookID,
			Notes:      req.Notes,
			IsUser:     isUser,
			UserIDRef:  userIDRef,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		if err := s.db.Create(contact).Error; err != nil {
			return nil, err
		}
	}

	// Create user-contact relationship
	userContact := &models.UserContact{
		ID:        uuid.New(),
		UserID:    userID,
		ContactID: contact.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.db.Create(userContact).Error; err != nil {
		return nil, err
	}

	// Create reciprocal contact if this contact is also a user
	if req.Email != nil && isUser {
		if err := s.CreateReciprocalContact(*req.Email, userID); err != nil {
			// Log the error but don't fail contact creation
			// TODO: Add proper logging here
		}
	}

	return contact, nil
}

func (s *ContactServiceGORM) GetContact(id uuid.UUID, userID uuid.UUID) (*models.Contact, error) {
	var userContact models.UserContact
	if err := s.db.Joins("Contact").Where("user_contacts.contact_id = ? AND user_contacts.user_id = ?", id, userID).First(&userContact).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("contact not found")
		}
		return nil, err
	}
	return &userContact.Contact, nil
}

func (s *ContactServiceGORM) GetUserContacts(userID uuid.UUID) ([]models.Contact, error) {
	var userContacts []models.UserContact
	if err := s.db.Joins("Contact").Where("user_contacts.user_id = ?", userID).Order("\"Contact\".\"name\" ASC").Find(&userContacts).Error; err != nil {
		return nil, err
	}
	
	var contacts []models.Contact
	for _, uc := range userContacts {
		contacts = append(contacts, uc.Contact)
	}
	
	return contacts, nil
}

func (s *ContactServiceGORM) UpdateContact(id uuid.UUID, userID uuid.UUID, req *models.UpdateContactRequest) (*models.Contact, error) {
	var userContact models.UserContact
	if err := s.db.Joins("Contact").Where("user_contacts.contact_id = ? AND user_contacts.user_id = ?", id, userID).First(&userContact).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("contact not found")
		}
		return nil, err
	}

	contact := &userContact.Contact

	// Update fields if provided
	if req.Name != nil {
		contact.Name = *req.Name
	}
	if req.Email != nil {
		contact.Email = req.Email
		// Check if this contact is also a user (by email)
		if req.Email != nil {
			var user models.User
			if err := s.db.Where("email = ?", *req.Email).First(&user).Error; err == nil {
				contact.IsUser = true
				contact.UserIDRef = &user.ID
			} else {
				contact.IsUser = false
				contact.UserIDRef = nil
			}
		}
	}
	if req.Phone != nil {
		contact.Phone = req.Phone
	}
	if req.FacebookID != nil {
		contact.FacebookID = req.FacebookID
	}
	if req.Notes != nil {
		contact.Notes = req.Notes
	}

	contact.UpdatedAt = time.Now()

	if err := s.db.Save(contact).Error; err != nil {
		return nil, err
	}

	return contact, nil
}

func (s *ContactServiceGORM) DeleteContact(id uuid.UUID, userID uuid.UUID) error {
	// First check if the user-contact relationship exists
	var userContact models.UserContact
	if err := s.db.Where("user_contacts.contact_id = ? AND user_contacts.user_id = ?", id, userID).First(&userContact).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("contact not found")
		}
		return err
	}

	// Delete the user-contact relationship
	result := s.db.Where("user_contacts.contact_id = ? AND user_contacts.user_id = ?", id, userID).Delete(&models.UserContact{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("contact not found")
	}
	return nil
}

// CreateContactsForNewUser creates contacts for a new user based on existing contacts that have their email
func (s *ContactServiceGORM) CreateContactsForNewUser(userID uuid.UUID, userEmail string) error {
	// Find all contacts that have this user's email
	var existingContacts []models.Contact
	if err := s.db.Where("email = ?", userEmail).Find(&existingContacts).Error; err != nil {
		return err
	}

	if len(existingContacts) == 0 {
		return nil // No existing contacts found with this email
	}

	// For each existing contact, create a user-contact relationship for the new user
	for _, existingContact := range existingContacts {
		// Get all users who own this contact
		var userContacts []models.UserContact
		if err := s.db.Where("user_contacts.contact_id = ?", existingContact.ID).Find(&userContacts).Error; err != nil {
			continue // Skip if no user-contact relationships found
		}

		// Process each user who has this contact
		for _, userContact := range userContacts {
			var contactOwner models.User
			if err := s.db.Where("id = ?", userContact.UserID).First(&contactOwner).Error; err != nil {
				continue // Skip if user not found
			}

			// Check if the new user already has a relationship with the contact owner
			var existingUserContact models.UserContact
			if err := s.db.Where("user_contacts.user_id = ? AND user_contacts.contact_id = ?", userID, userContact.UserID).First(&existingUserContact).Error; err == nil {
				continue // Relationship already exists, skip
			}

			// Find or create a contact for the contact owner
			var contactForOwner *models.Contact
			var existingContactForOwner models.Contact
			if err := s.db.Where("email = ?", contactOwner.Email).First(&existingContactForOwner).Error; err == nil {
				contactForOwner = &existingContactForOwner
			} else {
				// Create new contact for the contact owner
				contactForOwner = &models.Contact{
					ID:         uuid.New(),
					Name:       contactOwner.FirstName + " " + contactOwner.LastName,
					Email:      &contactOwner.Email,
					Phone:      contactOwner.Phone,
					IsUser:     true,
					UserIDRef:  &contactOwner.ID,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				}

				if err := s.db.Create(contactForOwner).Error; err != nil {
					continue
				}
			}

			// Create user-contact relationship for the new user
			// user_id = newly registered user's ID
			// contact_id = the contact ID that represents the user who owns the original contact
			newUserContact := &models.UserContact{
				ID:        uuid.New(),
				UserID:    userID,
				ContactID: contactForOwner.ID,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			if err := s.db.Create(newUserContact).Error; err != nil {
				// Log error but continue with other contacts
				// TODO: Add proper logging here
				continue
			}
		}
	}

	return nil
}

// CreateReciprocalContact creates a contact for an existing user when someone adds them as a contact
func (s *ContactServiceGORM) CreateReciprocalContact(contactEmail string, contactOwnerID uuid.UUID) error {
	// Find the user who owns this email
	var existingUser models.User
	if err := s.db.Where("email = ?", contactEmail).First(&existingUser).Error; err != nil {
		return nil // User doesn't exist, no reciprocal contact needed
	}

	// Get the contact owner's details
	var contactOwner models.User
	if err := s.db.Where("id = ?", contactOwnerID).First(&contactOwner).Error; err != nil {
		return err
	}

	// Check if the reciprocal contact already exists
	var existingUserContact models.UserContact
	if err := s.db.Joins("Contact").Where("user_contacts.user_id = ? AND \"Contact\".\"email\" = ?", existingUser.ID, contactOwner.Email).First(&existingUserContact).Error; err == nil {
		return nil // Reciprocal contact already exists
	}

	// Find or create contact for the contact owner
	var contact *models.Contact
	var existingContact models.Contact
	if err := s.db.Where("email = ?", contactOwner.Email).First(&existingContact).Error; err == nil {
		contact = &existingContact
	} else {
		// Create new contact for the contact owner
		contact = &models.Contact{
			ID:         uuid.New(),
			Name:       contactOwner.FirstName + " " + contactOwner.LastName,
			Email:      &contactOwner.Email,
			Phone:      contactOwner.Phone,
			IsUser:     true,
			UserIDRef:  &contactOwner.ID,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		if err := s.db.Create(contact).Error; err != nil {
			return err
		}
	}

	// Create user-contact relationship for the existing user
	userContact := &models.UserContact{
		ID:        uuid.New(),
		UserID:    existingUser.ID,
		ContactID: contact.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.db.Create(userContact).Error; err != nil {
		return err
	}

	return nil
} 