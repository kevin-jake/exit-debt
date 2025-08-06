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
	var existingContact models.Contact
	if err := s.db.Where("user_id = ? AND name = ?", userID, req.Name).First(&existingContact).Error; err == nil {
		return nil, errors.New("contact with this name already exists")
	}

	contact := &models.Contact{
		ID:         uuid.New(),
		UserID:     userID,
		Name:       req.Name,
		Email:      req.Email,
		Phone:      req.Phone,
		FacebookID: req.FacebookID,
		Notes:      req.Notes,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.db.Create(contact).Error; err != nil {
		return nil, err
	}

	return contact, nil
}

func (s *ContactServiceGORM) GetContact(id uuid.UUID, userID uuid.UUID) (*models.Contact, error) {
	var contact models.Contact
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&contact).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("contact not found")
		}
		return nil, err
	}
	return &contact, nil
}

func (s *ContactServiceGORM) GetUserContacts(userID uuid.UUID) ([]models.Contact, error) {
	var contacts []models.Contact
	if err := s.db.Where("user_id = ?", userID).Order("name ASC").Find(&contacts).Error; err != nil {
		return nil, err
	}
	return contacts, nil
}

func (s *ContactServiceGORM) UpdateContact(id uuid.UUID, userID uuid.UUID, req *models.UpdateContactRequest) (*models.Contact, error) {
	var contact models.Contact
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&contact).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("contact not found")
		}
		return nil, err
	}

	// Update fields if provided
	if req.Name != nil {
		contact.Name = *req.Name
	}
	if req.Email != nil {
		contact.Email = req.Email
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

	if err := s.db.Save(&contact).Error; err != nil {
		return nil, err
	}

	return &contact, nil
}

func (s *ContactServiceGORM) DeleteContact(id uuid.UUID, userID uuid.UUID) error {
	result := s.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Contact{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("contact not found")
	}
	return nil
} 