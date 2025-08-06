package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"exit-debt/internal/models"
)

type DebtServiceGORM struct {
	db *gorm.DB
	paymentScheduleService *PaymentScheduleService
}

func NewDebtServiceGORM(db *gorm.DB) *DebtServiceGORM {
	return &DebtServiceGORM{
		db: db,
		paymentScheduleService: NewPaymentScheduleService(),
	}
}

// Debt List Methods
func (s *DebtServiceGORM) CreateDebtList(userID uuid.UUID, req *models.CreateDebtListRequest) (*models.DebtList, error) {
	// Verify contact exists and belongs to user using user_contacts table
	var userContact models.UserContact
	if err := s.db.Joins("Contact").Where("user_contacts.contact_id = ? AND user_contacts.user_id = ?", req.ContactID, userID).First(&userContact).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("contact not found")
		}
		return nil, err
	}

	// Parse total amount
	totalAmount, err := decimal.NewFromString(req.TotalAmount)
	if err != nil {
		return nil, errors.New("invalid total amount format")
	}

	// Set default currency if not provided
	currency := req.Currency
	if currency == "" {
		currency = "Php"
	}

	// Determine due date and installment amount based on input
	var dueDate time.Time
	var installmentAmount decimal.Decimal
	var numberOfPayments *int

	createdAt := time.Now()

	if req.NumberOfPayments != nil && *req.NumberOfPayments > 0 {
		// Use number of payments to calculate due date and installment amount
		numberOfPayments = req.NumberOfPayments
		dueDate = s.paymentScheduleService.CalculateDueDateFromNumberOfPayments(createdAt, *req.NumberOfPayments, req.InstallmentPlan)
		installmentAmount = s.paymentScheduleService.CalculateInstallmentAmountFromNumberOfPayments(totalAmount, *req.NumberOfPayments)
	} else if req.DueDate != nil {
		// Use provided due date (existing behavior)
		dueDate = *req.DueDate
		installmentAmount = s.paymentScheduleService.CalculateInstallmentAmount(totalAmount, req.InstallmentPlan, createdAt, dueDate)
	} else {
		// Default to 1 payment if neither is provided
		defaultPayments := 1
		numberOfPayments = &defaultPayments
		dueDate = s.paymentScheduleService.CalculateDueDateFromNumberOfPayments(createdAt, defaultPayments, req.InstallmentPlan)
		installmentAmount = s.paymentScheduleService.CalculateInstallmentAmountFromNumberOfPayments(totalAmount, defaultPayments)
	}

	// Validate that due date is in the future
	now := time.Now()
	if dueDate.Before(now) || dueDate.Equal(now) {
		return nil, errors.New("due date must be in the future")
	}

	debtList := &models.DebtList{
		ID:                uuid.New(),
		UserID:            userID,
		ContactID:         req.ContactID,
		DebtType:          req.DebtType,
		TotalAmount:       totalAmount,
		InstallmentAmount: installmentAmount,
		TotalPaymentsMade: decimal.Zero,
		TotalRemainingDebt: totalAmount,
		Currency:          currency,
		Status:            "active",
		DueDate:           dueDate,
		NextPaymentDate:   s.paymentScheduleService.CalculateNextPaymentDate(&models.DebtList{
			InstallmentPlan: req.InstallmentPlan,
			CreatedAt:       createdAt,
		}, nil),
		InstallmentPlan:   req.InstallmentPlan,
		NumberOfPayments:  numberOfPayments,
		Description:       req.Description,
		Notes:             req.Notes,
		CreatedAt:         createdAt,
		UpdatedAt:         createdAt,
	}

	if err := s.db.Create(debtList).Error; err != nil {
		return nil, err
	}

	return debtList, nil
}

func (s *DebtServiceGORM) GetDebtList(id uuid.UUID, userID uuid.UUID) (*models.DebtListResponse, error) {
	var debtList models.DebtList
	if err := s.db.Preload("Contact").Preload("User").Preload("Payments").Where("id = ? AND user_id = ?", id, userID).First(&debtList).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("debt list not found")
		}
		return nil, err
	}
	
	response := debtList.ToDebtListResponse()
	return &response, nil
}

func (s *DebtServiceGORM) GetUserDebtLists(userID uuid.UUID) ([]models.DebtListResponse, error) {
	// Get debt lists that belong to the current user
	var userDebtLists []models.DebtList
	if err := s.db.Preload("Contact").Preload("User").Preload("Payments").Where("user_id = ?", userID).Find(&userDebtLists).Error; err != nil {
		return nil, err
	}

	// Convert to response models
	responses := make([]models.DebtListResponse, len(userDebtLists))
	for i, debtList := range userDebtLists {
		responses[i] = debtList.ToDebtListResponse()
	}

	return responses, nil
}

func (s *DebtServiceGORM) UpdateDebtList(id uuid.UUID, userID uuid.UUID, req *models.UpdateDebtListRequest) (*models.DebtList, error) {
	var debtList models.DebtList
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&debtList).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("debt list not found")
		}
		return nil, err
	}

	// Update fields if provided
	if req.TotalAmount != nil {
		totalAmount, err := decimal.NewFromString(*req.TotalAmount)
		if err != nil {
			return nil, errors.New("invalid total amount format")
		}
		debtList.TotalAmount = totalAmount
		// Recalculate installment amount based on current settings
		if debtList.NumberOfPayments != nil {
			debtList.InstallmentAmount = s.paymentScheduleService.CalculateInstallmentAmountFromNumberOfPayments(totalAmount, *debtList.NumberOfPayments)
		} else {
			debtList.InstallmentAmount = s.paymentScheduleService.CalculateInstallmentAmount(totalAmount, debtList.InstallmentPlan, debtList.CreatedAt, debtList.DueDate)
		}
	}
	
	if req.InstallmentPlan != nil {
		debtList.InstallmentPlan = *req.InstallmentPlan
		// Recalculate installment amount and due date if number of payments is set
		if debtList.NumberOfPayments != nil {
			debtList.InstallmentAmount = s.paymentScheduleService.CalculateInstallmentAmountFromNumberOfPayments(debtList.TotalAmount, *debtList.NumberOfPayments)
			debtList.DueDate = s.paymentScheduleService.CalculateDueDateFromNumberOfPayments(debtList.CreatedAt, *debtList.NumberOfPayments, *req.InstallmentPlan)
		} else {
			debtList.InstallmentAmount = s.paymentScheduleService.CalculateInstallmentAmount(debtList.TotalAmount, *req.InstallmentPlan, debtList.CreatedAt, debtList.DueDate)
		}
	}

	if req.NumberOfPayments != nil {
		debtList.NumberOfPayments = req.NumberOfPayments
		// Recalculate installment amount and due date
		debtList.InstallmentAmount = s.paymentScheduleService.CalculateInstallmentAmountFromNumberOfPayments(debtList.TotalAmount, *req.NumberOfPayments)
		debtList.DueDate = s.paymentScheduleService.CalculateDueDateFromNumberOfPayments(debtList.CreatedAt, *req.NumberOfPayments, debtList.InstallmentPlan)
	}
	
	if req.Currency != nil {
		debtList.Currency = *req.Currency
	}
	if req.Status != nil {
		debtList.Status = *req.Status
	}
	if req.DueDate != nil {
		// Validate that due date is in the future
		now := time.Now()
		if req.DueDate.Before(now) || req.DueDate.Equal(now) {
			return nil, errors.New("due date must be in the future")
		}
		
		debtList.DueDate = *req.DueDate
		// If due date is manually set, clear number of payments as they're now independent
		debtList.NumberOfPayments = nil
		// Recalculate installment amount based on new due date
		debtList.InstallmentAmount = s.paymentScheduleService.CalculateInstallmentAmount(debtList.TotalAmount, debtList.InstallmentPlan, debtList.CreatedAt, *req.DueDate)
	}
	if req.Description != nil {
		debtList.Description = req.Description
	}
	if req.Notes != nil {
		debtList.Notes = req.Notes
	}

	debtList.UpdatedAt = time.Now()

	if err := s.db.Save(&debtList).Error; err != nil {
		return nil, err
	}

	return &debtList, nil
}

func (s *DebtServiceGORM) DeleteDebtList(id uuid.UUID, userID uuid.UUID) error {
	result := s.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.DebtList{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("debt list not found")
	}
	return nil
}

// Payment Methods (Debt Items are now Payments)
func (s *DebtServiceGORM) CreateDebtItem(userID uuid.UUID, req *models.CreateDebtItemRequest) (*models.DebtItem, error) {
	// Verify debt list exists and belongs to user
	var debtList models.DebtList
	if err := s.db.Where("id = ? AND user_id = ?", req.DebtListID, userID).First(&debtList).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("debt list not found")
		}
		return nil, err
	}

	// Parse amount
	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		return nil, errors.New("invalid amount format")
	}

	// Set default currency if not provided
	currency := req.Currency
	if currency == "" {
		currency = debtList.Currency
	}

	payment := &models.DebtItem{
		ID:            uuid.New(),
		DebtListID:    req.DebtListID,
		Amount:        amount,
		Currency:      currency,
		PaymentDate:   req.PaymentDate,
		PaymentMethod: req.PaymentMethod,
		Description:   req.Description,
		Status:        "completed",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.db.Create(payment).Error; err != nil {
		return nil, err
	}

	// Update debt list status, next payment date, and payment totals
	if err := s.updateDebtListStatusAndPaymentTotals(debtList.ID); err != nil {
		return nil, err
	}

	return payment, nil
}

func (s *DebtServiceGORM) GetDebtItem(id uuid.UUID, userID uuid.UUID) (*models.DebtItem, error) {
	var debtItem models.DebtItem
	if err := s.db.Preload("DebtList").Preload("DebtList.Contact").Preload("DebtList.User").Joins("JOIN debt_lists ON debt_items.debt_list_id = debt_lists.id").
		Where("debt_items.id = ? AND debt_lists.user_id = ?", id, userID).
		First(&debtItem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("debt item not found")
		}
		return nil, err
	}
	return &debtItem, nil
}

func (s *DebtServiceGORM) GetDebtListItems(debtListID uuid.UUID, userID uuid.UUID) ([]models.DebtItem, error) {
	var debtItems []models.DebtItem
	if err := s.db.Preload("DebtList").Preload("DebtList.Contact").Preload("DebtList.User").Joins("JOIN debt_lists ON debt_items.debt_list_id = debt_lists.id").
		Where("debt_items.debt_list_id = ? AND debt_lists.user_id = ?", debtListID, userID).
		Order("debt_items.created_at DESC").
		Find(&debtItems).Error; err != nil {
		return nil, err
	}
	return debtItems, nil
}

func (s *DebtServiceGORM) UpdateDebtItem(id uuid.UUID, userID uuid.UUID, req *models.UpdateDebtItemRequest) (*models.DebtItem, error) {
	var debtItem models.DebtItem
	if err := s.db.Joins("JOIN debt_lists ON debt_items.debt_list_id = debt_lists.id").
		Where("debt_items.id = ? AND debt_lists.user_id = ?", id, userID).
		First(&debtItem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("debt item not found")
		}
		return nil, err
	}

	// Update fields if provided
	if req.Amount != nil {
		amount, err := decimal.NewFromString(*req.Amount)
		if err != nil {
			return nil, errors.New("invalid amount format")
		}
		debtItem.Amount = amount
	}
	if req.Currency != nil {
		debtItem.Currency = *req.Currency
	}
	if req.PaymentDate != nil {
		debtItem.PaymentDate = *req.PaymentDate
	}
	if req.PaymentMethod != nil {
		debtItem.PaymentMethod = *req.PaymentMethod
	}
	if req.Description != nil {
		debtItem.Description = req.Description
	}
	if req.Status != nil {
		debtItem.Status = *req.Status
	}

	debtItem.UpdatedAt = time.Now()

	if err := s.db.Save(&debtItem).Error; err != nil {
		return nil, err
	}

	// Update debt list status, next payment date, and payment totals
	if err := s.updateDebtListStatusAndPaymentTotals(debtItem.DebtListID); err != nil {
		return nil, err
	}

	return &debtItem, nil
}

func (s *DebtServiceGORM) DeleteDebtItem(id uuid.UUID, userID uuid.UUID) error {
	var debtItem models.DebtItem
	if err := s.db.Joins("JOIN debt_lists ON debt_items.debt_list_id = debt_lists.id").
		Where("debt_items.id = ? AND debt_lists.user_id = ?", id, userID).
		First(&debtItem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("debt item not found")
		}
		return err
	}

	debtListID := debtItem.DebtListID

	if err := s.db.Delete(&debtItem).Error; err != nil {
		return err
	}

	// Update debt list status, next payment date, and payment totals
	return s.updateDebtListStatusAndPaymentTotals(debtListID)
}

func (s *DebtServiceGORM) GetOverdueItems(userID uuid.UUID) ([]models.DebtList, error) {
	var debtLists []models.DebtList
	if err := s.db.Where("user_id = ? AND next_payment_date < ? AND status = ?", userID, time.Now(), "active").
		Order("next_payment_date ASC").
		Find(&debtLists).Error; err != nil {
		return nil, err
	}
	return debtLists, nil
}

func (s *DebtServiceGORM) GetDueSoonItems(userID uuid.UUID, days int) ([]models.DebtList, error) {
	var debtLists []models.DebtList
	dueDate := time.Now().AddDate(0, 0, days)
	if err := s.db.Where("user_id = ? AND next_payment_date BETWEEN ? AND ? AND status = ?", userID, time.Now(), dueDate, "active").
		Order("next_payment_date ASC").
		Find(&debtLists).Error; err != nil {
		return nil, err
	}
	return debtLists, nil
}

// GetPaymentSchedule returns the complete payment schedule for a debt
func (s *DebtServiceGORM) GetPaymentSchedule(debtListID uuid.UUID, userID uuid.UUID) ([]PaymentScheduleItem, error) {
	// Verify debt list exists and belongs to user
	var debtList models.DebtList
	if err := s.db.Where("id = ? AND user_id = ?", debtListID, userID).First(&debtList).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("debt list not found")
		}
		return nil, err
	}

	return s.paymentScheduleService.CalculatePaymentSchedule(&debtList, s.db), nil
}

// GetUpcomingPayments returns payments due in the next X days for a user
func (s *DebtServiceGORM) GetUpcomingPayments(userID uuid.UUID, days int) ([]UpcomingPayment, error) {
	var debtLists []models.DebtList
	if err := s.db.Where("user_id = ? AND status IN (?)", userID, []string{"active", "overdue"}).
		Preload("Contact").
		Find(&debtLists).Error; err != nil {
		return nil, err
	}

	var upcomingPayments []UpcomingPayment
	cutoffDate := time.Now().AddDate(0, 0, days)

	for _, debtList := range debtLists {
		if debtList.NextPaymentDate.After(time.Now()) && debtList.NextPaymentDate.Before(cutoffDate) {
			upcomingPayment := UpcomingPayment{
				DebtListID:      debtList.ID,
				ContactName:      debtList.Contact.Name,
				DebtType:        debtList.DebtType,
				NextPaymentDate: debtList.NextPaymentDate,
				Amount:          debtList.InstallmentAmount,
				Currency:        debtList.Currency,
				Description:     debtList.Description,
			}
			upcomingPayments = append(upcomingPayments, upcomingPayment)
		}
	}

	return upcomingPayments, nil
}

// GetTotalPaymentsForDebtList returns the total payments made for a specific debt list
func (s *DebtServiceGORM) GetTotalPaymentsForDebtList(debtListID uuid.UUID, userID uuid.UUID) (*PaymentSummary, error) {
	// Verify debt list exists and belongs to user
	var debtList models.DebtList
	if err := s.db.Where("id = ? AND user_id = ?", debtListID, userID).First(&debtList).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("debt list not found")
		}
		return nil, err
	}

	// Get all completed payments for this debt list
	var payments []models.DebtItem
	if err := s.db.Where("debt_list_id = ? AND status = ?", debtListID, "completed").
		Order("payment_date ASC").
		Find(&payments).Error; err != nil {
		return nil, err
	}

	// Calculate total payments made
	totalPaid := decimal.Zero
	for _, payment := range payments {
		totalPaid = totalPaid.Add(payment.Amount)
	}

	// Calculate remaining debt
	remainingDebt := debtList.TotalAmount.Sub(totalPaid)
	if remainingDebt.LessThan(decimal.Zero) {
		remainingDebt = decimal.Zero
	}

	// Calculate percentage paid
	percentagePaid := decimal.Zero
	if debtList.TotalAmount.GreaterThan(decimal.Zero) {
		percentagePaid = totalPaid.Div(debtList.TotalAmount).Mul(decimal.NewFromInt(100))
	}

	return &PaymentSummary{
		DebtListID:      debtListID,
		TotalAmount:      debtList.TotalAmount,
		TotalPaid:        totalPaid,
		RemainingDebt:    remainingDebt,
		PercentagePaid:   percentagePaid,
		NumberOfPayments: len(payments),
		Payments:         payments,
	}, nil
}

// PaymentSummary represents a summary of payments for a debt list
type PaymentSummary struct {
	DebtListID      uuid.UUID       `json:"debt_list_id"`
	TotalAmount     decimal.Decimal `json:"total_amount"`
	TotalPaid       decimal.Decimal `json:"total_paid"`
	RemainingDebt   decimal.Decimal `json:"remaining_debt"`
	PercentagePaid  decimal.Decimal `json:"percentage_paid"`
	NumberOfPayments int             `json:"number_of_payments"`
	Payments        []models.DebtItem `json:"payments"`
}

// UpcomingPayment represents an upcoming payment
type UpcomingPayment struct {
	DebtListID      uuid.UUID     `json:"debt_list_id"`
	ContactName      string        `json:"contact_name"`
	DebtType        string        `json:"debt_type"`
	NextPaymentDate time.Time     `json:"next_payment_date"`
	Amount          decimal.Decimal `json:"amount"`
	Currency        string        `json:"currency"`
	Description     *string       `json:"description"`
}

func (s *DebtServiceGORM) updateDebtListStatusAndPaymentTotals(debtListID uuid.UUID) error {
	// Get the debt list
	var debtList models.DebtList
	if err := s.db.Where("id = ?", debtListID).First(&debtList).Error; err != nil {
		return err
	}

	// Calculate total payments made
	var totalPaid decimal.Decimal
	if err := s.db.Model(&models.DebtItem{}).
		Where("debt_list_id = ? AND status = ?", debtListID, "completed").
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalPaid).Error; err != nil {
		return err
	}

	// Calculate remaining amount
	remainingAmount := debtList.TotalAmount.Sub(totalPaid)
	if remainingAmount.LessThan(decimal.Zero) {
		remainingAmount = decimal.Zero
	}

	// Determine new status
	var newStatus string
	if remainingAmount.LessThanOrEqual(decimal.Zero) {
		newStatus = "settled"
	} else if time.Now().After(debtList.DueDate) {
		newStatus = "overdue"
	} else {
		newStatus = "active"
	}

	// Get the last payment date to calculate next payment
	var lastPayment models.DebtItem
	var lastPaymentDate *time.Time
	if err := s.db.Where("debt_list_id = ? AND status = ?", debtListID, "completed").
		Order("payment_date DESC").
		First(&lastPayment).Error; err == nil {
		lastPaymentDate = &lastPayment.PaymentDate
	}

	// Calculate next payment date
	nextPaymentDate := s.paymentScheduleService.CalculateNextPaymentDate(&debtList, lastPaymentDate)

	// Update debt list status, next payment date, and payment totals
	return s.db.Model(&models.DebtList{}).
		Where("id = ?", debtListID).
		Updates(map[string]interface{}{
			"status":               newStatus,
			"next_payment_date":   nextPaymentDate,
			"total_payments_made": totalPaid,
			"total_remaining_debt": remainingAmount,
		}).Error
} 