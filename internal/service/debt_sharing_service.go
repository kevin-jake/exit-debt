package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"exit-debt/internal/models"
)

type DebtSharingService struct {
	db *gorm.DB
}

func NewDebtSharingService(db *gorm.DB) *DebtSharingService {
	return &DebtSharingService{db: db}
}

// ShareDebtListsWithUser shares debt lists with a newly registered user based on email matching
func (s *DebtSharingService) ShareDebtListsWithUser(userID uuid.UUID, userEmail string) error {
	// Find all contacts with matching email
	var contacts []models.Contact
	if err := s.db.Where("email = ?", userEmail).Find(&contacts).Error; err != nil {
		return err
	}

	if len(contacts) == 0 {
		return nil // No matching contacts found
	}

	// For each contact, share their debt lists with the new user
	for _, contact := range contacts {
		if err := s.shareContactDebtLists(userID, contact); err != nil {
			return err
		}
	}

	return nil
}

// shareContactDebtLists shares all debt lists from a specific contact with the new user
func (s *DebtSharingService) shareContactDebtLists(userID uuid.UUID, contact models.Contact) error {
	// Get all debt lists for this contact
	var debtLists []models.DebtList
	if err := s.db.Where("contact_id = ?", contact.ID).Find(&debtLists).Error; err != nil {
		return err
	}

	// Share each debt list by creating a new debt list for the new user
	for _, debtList := range debtLists {
		if err := s.shareDebtList(userID, contact, debtList); err != nil {
			return err
		}
	}

	return nil
}

// shareDebtList creates a shared debt list for the new user with inverted debt type
func (s *DebtSharingService) shareDebtList(userID uuid.UUID, contact models.Contact, originalDebtList models.DebtList) error {
	// Determine the new debt type based on the original type
	var newDebtType string
	switch originalDebtList.DebtType {
	case "owed_to_me":
		// If the original user was owed money, the new user owes money
		newDebtType = "i_owe"
	case "i_owe":
		// If the original user owed money, the new user is owed money
		newDebtType = "owed_to_me"
	default:
		return errors.New("invalid debt type")
	}

	// Create a new debt list for the new user
	sharedDebtList := &models.DebtList{
		ID:                uuid.New(),
		UserID:            userID,
		ContactID:         contact.ID,
		DebtType:          newDebtType,
		TotalAmount:       originalDebtList.TotalAmount,
		InstallmentAmount: originalDebtList.InstallmentAmount,
		Currency:          originalDebtList.Currency,
		Status:            originalDebtList.Status,
		DueDate:           originalDebtList.DueDate,
		NextPaymentDate:   originalDebtList.NextPaymentDate,
		InstallmentPlan:   originalDebtList.InstallmentPlan,
		Description:       originalDebtList.Description,
		Notes:             originalDebtList.Notes,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if err := s.db.Create(sharedDebtList).Error; err != nil {
		return err
	}

	// Copy all payments from the original debt list to the new shared debt list
	if err := s.copyPayments(originalDebtList.ID, sharedDebtList.ID); err != nil {
		return err
	}

	return nil
}

// copyPayments copies all payments from the original debt list to the new shared debt list
func (s *DebtSharingService) copyPayments(originalDebtListID uuid.UUID, newDebtListID uuid.UUID) error {
	var payments []models.DebtItem
	if err := s.db.Where("debt_list_id = ?", originalDebtListID).Find(&payments).Error; err != nil {
		return err
	}

	for _, payment := range payments {
		// Create a new payment record for the shared debt list
		sharedPayment := &models.DebtItem{
			ID:            uuid.New(),
			DebtListID:    newDebtListID,
			Amount:        payment.Amount,
			Currency:      payment.Currency,
			PaymentDate:   payment.PaymentDate,
			PaymentMethod: payment.PaymentMethod,
			Description:   payment.Description,
			Status:        payment.Status,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		if err := s.db.Create(sharedPayment).Error; err != nil {
			return err
		}
	}

	return nil
}

// GetSharingSummary returns a summary of shared debt lists
func (s *DebtSharingService) GetSharingSummary(userID uuid.UUID, userEmail string) (*SharingSummary, error) {
	// Find all contacts with matching email
	var contacts []models.Contact
	if err := s.db.Where("email = ?", userEmail).Find(&contacts).Error; err != nil {
		return nil, err
	}

	if len(contacts) == 0 {
		return &SharingSummary{
			ContactsFound: 0,
			DebtListsShared: 0,
			TotalAmountShared: "0.00",
		}, nil
	}

	var totalDebtLists int
	var totalAmount float64

	// Count debt lists and calculate total amount
	for _, contact := range contacts {
		var debtLists []models.DebtList
		if err := s.db.Where("contact_id = ?", contact.ID).Find(&debtLists).Error; err != nil {
			return nil, err
		}

		totalDebtLists += len(debtLists)
		for _, debtList := range debtLists {
			amount, _ := debtList.TotalAmount.Float64()
			totalAmount += amount
		}
	}

	return &SharingSummary{
		ContactsFound: len(contacts),
		DebtListsShared: totalDebtLists,
		TotalAmountShared: fmt.Sprintf("%.2f", totalAmount),
	}, nil
}

// SharingSummary represents a summary of shared debt lists
type SharingSummary struct {
	ContactsFound      int    `json:"contacts_found"`
	DebtListsShared    int    `json:"debt_lists_shared"`
	TotalAmountShared  string `json:"total_amount_shared"`
} 