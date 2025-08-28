package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"exit-debt/internal/domain/entities"
	"exit-debt/internal/domain/interfaces"
)

// debtService implements the DebtService interface
type debtService struct {
	debtListRepo           interfaces.DebtListRepository
	debtItemRepo           interfaces.DebtItemRepository
	contactRepo            interfaces.ContactRepository
	paymentScheduleService interfaces.PaymentScheduleService
	fileStorageService     interfaces.FileStorageService
}

// NewDebtService creates a new debt service
func NewDebtService(
	debtListRepo interfaces.DebtListRepository,
	debtItemRepo interfaces.DebtItemRepository,
	contactRepo interfaces.ContactRepository,
	paymentScheduleService interfaces.PaymentScheduleService,
	fileStorageService interfaces.FileStorageService,
) interfaces.DebtService {
	return &debtService{
		debtListRepo:           debtListRepo,
		debtItemRepo:           debtItemRepo,
		contactRepo:            contactRepo,
		paymentScheduleService: paymentScheduleService,
		fileStorageService:     fileStorageService,
	}
}

// Debt List operations

func (s *debtService) CreateDebtList(ctx context.Context, userID uuid.UUID, req *entities.CreateDebtListRequest) (*entities.DebtList, error) {
	// Validate input
	if err := s.validateCreateDebtListRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Verify contact exists and belongs to user
	_, err := s.contactRepo.GetUserContactRelation(ctx, userID, req.ContactID)
	if err != nil {
		return nil, fmt.Errorf("contact verification failed: %w", err)
	}

	// Parse total amount
	totalAmount, err := decimal.NewFromString(req.TotalAmount)
	if err != nil {
		return nil, entities.ErrInvalidAmount
	}
	
	// Validate that amount is positive
	if totalAmount.LessThanOrEqual(decimal.Zero) {
		return nil, entities.ErrInvalidAmount
	}

	// Set default currency if not provided
	currency := req.Currency
	if currency == "" {
		currency = "Php"
	}

	// Validation: If number_of_payments is provided, installment_plan is required
	if req.NumberOfPayments != nil && *req.NumberOfPayments > 0 && req.InstallmentPlan == "" {
		return nil, fmt.Errorf("installment_plan is required when number_of_payments is provided")
	}

	// Validation: If due_date is provided but installment_plan is not, default to 1-time payment
	installmentPlan := req.InstallmentPlan
	if req.DueDate != nil && installmentPlan == "" {
		installmentPlan = "onetime" // Default to onetime for 1-time payment calculation
	}

	// Determine due date and installment amount based on input
	var dueDate time.Time
	var installmentAmount decimal.Decimal
	var numberOfPayments *int

	createdAt := time.Now()

	if req.NumberOfPayments != nil && *req.NumberOfPayments > 0 {
		// Use number of payments to calculate due date and installment amount
		numberOfPayments = req.NumberOfPayments
		dueDate = s.paymentScheduleService.CalculateDueDateFromNumberOfPayments(createdAt, *req.NumberOfPayments, installmentPlan)
		installmentAmount = s.paymentScheduleService.CalculateInstallmentAmountFromNumberOfPayments(totalAmount, *req.NumberOfPayments)
	} else if req.DueDate != nil {
		// Use provided due date (existing behavior)
		dueDate = *req.DueDate
		installmentAmount = s.paymentScheduleService.CalculateInstallmentAmount(totalAmount, installmentPlan, createdAt, dueDate)
	} else {
		// Default to 1 payment if neither is provided
		defaultPayments := 1
		numberOfPayments = &defaultPayments
		dueDate = s.paymentScheduleService.CalculateDueDateFromNumberOfPayments(createdAt, defaultPayments, installmentPlan)
		installmentAmount = s.paymentScheduleService.CalculateInstallmentAmountFromNumberOfPayments(totalAmount, defaultPayments)
	}

	// Validate that due date is in the future
	now := time.Now()
	if dueDate.Before(now) || dueDate.Equal(now) {
		return nil, entities.ErrInvalidDueDate
	}

	// Calculate next payment date
	var nextPaymentDate time.Time
	if installmentPlan == "onetime" {
		nextPaymentDate = dueDate
	} else {
		nextPaymentDate = s.paymentScheduleService.CalculateNextPaymentDate(&entities.DebtList{
			InstallmentPlan: installmentPlan,
			CreatedAt:       createdAt,
			DueDate:         dueDate,
		}, nil)
	}

	debtList := &entities.DebtList{
		ID:                  uuid.New(),
		UserID:              userID,
		ContactID:           req.ContactID,
		DebtType:            req.DebtType,
		TotalAmount:         totalAmount,
		InstallmentAmount:   installmentAmount,
		TotalPaymentsMade:   decimal.Zero,
		TotalRemainingDebt:  totalAmount,
		Currency:            currency,
		Status:              "active",
		DueDate:             dueDate,
		NextPaymentDate:     nextPaymentDate,
		InstallmentPlan:     installmentPlan,
		NumberOfPayments:    numberOfPayments,
		Description:         req.Description,
		Notes:               req.Notes,
		CreatedAt:           createdAt,
		UpdatedAt:           createdAt,
	}

	// Validate debt list entity
	if err := debtList.IsValid(); err != nil {
		return nil, fmt.Errorf("invalid debt list entity: %w", err)
	}

	if err := s.debtListRepo.Create(ctx, debtList); err != nil {
		return nil, fmt.Errorf("failed to create debt list: %w", err)
	}

	return debtList, nil
}

func (s *debtService) GetDebtList(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*entities.DebtListResponse, error) {
	// First check if debt list belongs to user (user is the owner)
	belongs, err := s.debtListRepo.BelongsToUser(ctx, id, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify ownership: %w", err)
	}
	
	if belongs {
		// User owns the debt list, get it with relations
		debtListResponse, err := s.debtListRepo.GetByIDWithRelations(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("failed to get debt list: %w", err)
		}
		return debtListResponse, nil
	}
	
	// If user doesn't own it, check if they are a contact in the debt list
	// This allows users to view debt lists where they owe money or are owed money
	contactDebtLists, err := s.debtListRepo.GetDebtListsWhereUserIsContact(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get debt lists where user is contact: %w", err)
	}
	
	// Find the specific debt list where user is a contact
	for _, contactDebtList := range contactDebtLists {
		if contactDebtList.ID == id {
			// Flip the debt type for the user's perspective
			// When a user views a debt list where they are the contact,
			// the debt type should be from their perspective
			if contactDebtList.DebtType == "owed_to_me" {
				contactDebtList.DebtType = "i_owe"
			} else if contactDebtList.DebtType == "i_owe" {
				contactDebtList.DebtType = "owed_to_me"
			}
			return &contactDebtList, nil
		}
	}
	
	// User neither owns the debt list nor is a contact in it
	return nil, entities.ErrDebtListNotFound
}

func (s *debtService) GetUserDebtLists(ctx context.Context, userID uuid.UUID) ([]entities.DebtListResponse, error) {
	// Get debt lists owned by the user
	ownedDebtLists, err := s.debtListRepo.GetUserDebtLists(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user debt lists: %w", err)
	}

	// Get debt lists where the user is referenced as a contact (other users owe them or they owe other users)
	contactDebtLists, err := s.debtListRepo.GetDebtListsWhereUserIsContact(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get debt lists where user is contact: %w", err)
	}

	// Flip the debt type for debt lists where the user is the contact
	for i := range contactDebtLists {
		if contactDebtLists[i].DebtType == "owed_to_me" {
			contactDebtLists[i].DebtType = "i_owe"
		} else if contactDebtLists[i].DebtType == "i_owe" {
			contactDebtLists[i].DebtType = "owed_to_me"
		}
	}

	// Combine both lists
	allDebtLists := append(ownedDebtLists, contactDebtLists...)

	return allDebtLists, nil
}

func (s *debtService) UpdateDebtList(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *entities.UpdateDebtListRequest) (*entities.DebtList, error) {
	// Validate input
	if err := s.validateUpdateDebtListRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Check if debt list belongs to user
	belongs, err := s.debtListRepo.BelongsToUser(ctx, id, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify ownership: %w", err)
	}
	if !belongs {
		return nil, entities.ErrDebtListNotFound
	}

	// Get existing debt list
	debtList, err := s.debtListRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get debt list: %w", err)
	}

	// Update fields if provided
	if req.TotalAmount != nil {
		totalAmount, err := decimal.NewFromString(*req.TotalAmount)
		if err != nil {
			return nil, entities.ErrInvalidAmount
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
			return nil, entities.ErrInvalidDueDate
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

	// Validate updated debt list entity
	if err := debtList.IsValid(); err != nil {
		return nil, fmt.Errorf("invalid updated debt list entity: %w", err)
	}

	if err := s.debtListRepo.Update(ctx, debtList); err != nil {
		return nil, fmt.Errorf("failed to update debt list: %w", err)
	}

	return debtList, nil
}

func (s *debtService) DeleteDebtList(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	// Check if debt list belongs to user
	belongs, err := s.debtListRepo.BelongsToUser(ctx, id, userID)
	if err != nil {
		return fmt.Errorf("failed to verify ownership: %w", err)
	}
	if !belongs {
		return entities.ErrDebtListNotFound
	}

	if err := s.debtListRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete debt list: %w", err)
	}

	return nil
}

// Debt Item (Payment) operations

func (s *debtService) CreateDebtItem(ctx context.Context, userID uuid.UUID, req *entities.CreateDebtItemRequest) (*entities.DebtItem, error) {
	// Validate input
	if err := s.validateCreateDebtItemRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// First check if debt list belongs to user (user is the owner)
	belongs, err := s.debtListRepo.BelongsToUser(ctx, req.DebtListID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify debt list ownership: %w", err)
	}
	
	if belongs {
		// User owns the debt list, proceed with creation
	} else {
		// If user doesn't own it, check if they are a contact in the debt list
		// This allows users to create debt items for debt lists where they owe money or are owed money
		contactDebtLists, err := s.debtListRepo.GetDebtListsWhereUserIsContact(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to get debt lists where user is contact: %w", err)
		}
		
		// Check if the user is a contact in the specific debt list
		userIsContact := false
		for _, contactDebtList := range contactDebtLists {
			if contactDebtList.ID == req.DebtListID {
				userIsContact = true
				break
			}
		}
		
		if !userIsContact {
			// User neither owns the debt list nor is a contact in it
			return nil, entities.ErrDebtListNotFound
		}
	}

	// Get debt list for currency default
	debtList, err := s.debtListRepo.GetByID(ctx, req.DebtListID)
	if err != nil {
		return nil, fmt.Errorf("failed to get debt list: %w", err)
	}

	// Parse amount
	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		return nil, entities.ErrInvalidAmount
	}
	
	// Validate that amount is positive
	if amount.LessThanOrEqual(decimal.Zero) {
		return nil, entities.ErrInvalidAmount
	}

	// Set default currency if not provided
	currency := req.Currency
	if currency == "" {
		currency = debtList.Currency
	}

	// Determine initial status based on the user's perspective
	// When a user creates a payment, we need to consider their perspective:
	// - If they owe money (i_owe from their perspective), payments are pending until verified
	// - If they are owed money (owed_to_me from their perspective), payments are completed
	initialStatus := "completed"
	
	// Check if the user is the owner or a contact to determine their perspective
	if belongs {
		// User owns the debt list, use the debt list's debt type
		if debtList.DebtType == "i_owe" {
			initialStatus = "pending"
		}
	} else {
		// User is a contact, determine their perspective by flipping the debt type
		// If the debt list is "owed_to_me" (someone owes them), then from the contact's perspective it's "i_owe"
		if debtList.DebtType == "owed_to_me" {
			initialStatus = "pending"
		}
	}

	debtItem := &entities.DebtItem{
		ID:                uuid.New(),
		DebtListID:        req.DebtListID,
		Amount:            amount,
		Currency:          currency,
		PaymentDate:       req.PaymentDate,
		PaymentMethod:     req.PaymentMethod,
		Description:       req.Description,
		Status:            initialStatus,
		ReceiptPhotoURL:   req.ReceiptPhotoURL,
		VerificationNotes: req.VerificationNotes,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	// Validate debt item entity
	if err := debtItem.IsValid(); err != nil {
		return nil, fmt.Errorf("invalid debt item entity: %w", err)
	}

	if err := s.debtItemRepo.Create(ctx, debtItem); err != nil {
		return nil, fmt.Errorf("failed to create debt item: %w", err)
	}

	// Update debt list status, next payment date, and payment totals
	if err := s.updateDebtListStatusAndPaymentTotals(ctx, debtList.ID); err != nil {
		return nil, fmt.Errorf("failed to update debt list totals: %w", err)
	}

	return debtItem, nil
}

func (s *debtService) GetDebtItem(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*entities.DebtItem, error) {
	// Check if debt item belongs to user's debt list
	belongs, err := s.debtItemRepo.BelongsToUserDebtList(ctx, id, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify ownership: %w", err)
	}
	if !belongs {
		return nil, entities.ErrDebtItemNotFound
	}

	debtItem, err := s.debtItemRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get debt item: %w", err)
	}

	return debtItem, nil
}

func (s *debtService) GetDebtListItems(ctx context.Context, debtListID uuid.UUID, userID uuid.UUID) ([]entities.DebtItem, error) {
	// First check if debt list belongs to user (user is the owner)
	belongs, err := s.debtListRepo.BelongsToUser(ctx, debtListID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify ownership: %w", err)
	}
	
	if belongs {
		// User owns the debt list, get the debt items
		debtItems, err := s.debtItemRepo.GetByDebtListID(ctx, debtListID)
		if err != nil {
			return nil, fmt.Errorf("failed to get debt list items: %w", err)
		}
		return debtItems, nil
	}
	
	// If user doesn't own it, check if they are a contact in the debt list
	// This allows users to view debt items for debt lists where they owe money or are owed money
	contactDebtLists, err := s.debtListRepo.GetDebtListsWhereUserIsContact(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get debt lists where user is contact: %w", err)
	}
	
	// Check if the user is a contact in the specific debt list
	for _, contactDebtList := range contactDebtLists {
		if contactDebtList.ID == debtListID {
			// User is a contact in this debt list, get the debt items
			debtItems, err := s.debtItemRepo.GetByDebtListID(ctx, debtListID)
			if err != nil {
				return nil, fmt.Errorf("failed to get debt list items: %w", err)
			}
			return debtItems, nil
		}
	}
	
	// User neither owns the debt list nor is a contact in it
	return nil, entities.ErrDebtListNotFound
}

func (s *debtService) UpdateDebtItem(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *entities.UpdateDebtItemRequest) (*entities.DebtItem, error) {
	// Validate input
	if err := s.validateUpdateDebtItemRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// First check if debt item belongs to user's debt list (user is the owner)
	belongs, err := s.debtItemRepo.BelongsToUserDebtList(ctx, id, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify ownership: %w", err)
	}
	
	if belongs {
		// User owns the debt list, proceed with update
	} else {
		// If user doesn't own it, check if they are a contact in the debt list
		// This allows users to update debt items for debt lists where they owe money or are owed money
		debtItem, err := s.debtItemRepo.GetByID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("failed to get debt item: %w", err)
		}
		
		// Check if the user is a contact in this debt list
		contactDebtLists, err := s.debtListRepo.GetDebtListsWhereUserIsContact(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to get debt lists where user is contact: %w", err)
		}
		
		// Check if the user is a contact in the specific debt list
		userIsContact := false
		for _, contactDebtList := range contactDebtLists {
			if contactDebtList.ID == debtItem.DebtListID {
				userIsContact = true
				break
			}
		}
		
		if !userIsContact {
			// User neither owns the debt list nor is a contact in it
			return nil, entities.ErrDebtItemNotFound
		}
	}

	// Get existing debt item
	debtItem, err := s.debtItemRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get debt item: %w", err)
	}

	// Update fields if provided
	if req.Amount != nil {
		amount, err := decimal.NewFromString(*req.Amount)
		if err != nil {
			return nil, entities.ErrInvalidAmount
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
	if req.ReceiptPhotoURL != nil {
		// If there's an old receipt photo, delete it from S3
		if debtItem.ReceiptPhotoURL != nil && *debtItem.ReceiptPhotoURL != "" {
			if err := s.fileStorageService.DeleteReceipt(ctx, *debtItem.ReceiptPhotoURL); err != nil {
				// Log the error but don't fail the update
				// In production, you might want to handle this differently
				fmt.Printf("Warning: failed to delete old receipt photo: %v\n", err)
			}
		}
		debtItem.ReceiptPhotoURL = req.ReceiptPhotoURL
	}
	if req.VerificationNotes != nil {
		debtItem.VerificationNotes = req.VerificationNotes
	}

	debtItem.UpdatedAt = time.Now()

	// Validate updated debt item entity
	if err := debtItem.IsValid(); err != nil {
		return nil, fmt.Errorf("invalid updated debt item entity: %w", err)
	}

	if err := s.debtItemRepo.Update(ctx, debtItem); err != nil {
		return nil, fmt.Errorf("failed to update debt item: %w", err)
	}

	// Update debt list status, next payment date, and payment totals
	if err := s.updateDebtListStatusAndPaymentTotals(ctx, debtItem.DebtListID); err != nil {
		return nil, fmt.Errorf("failed to update debt list totals: %w", err)
	}

	return debtItem, nil
}

func (s *debtService) DeleteDebtItem(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	// Check if debt item belongs to user's debt list
	belongs, err := s.debtItemRepo.BelongsToUserDebtList(ctx, id, userID)
	if err != nil {
		return fmt.Errorf("failed to verify ownership: %w", err)
	}
	if !belongs {
		return entities.ErrDebtItemNotFound
	}

	// Get debt item to know which debt list to update
	debtItem, err := s.debtItemRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get debt item: %w", err)
	}

	debtListID := debtItem.DebtListID

	// If there's a receipt photo, delete it from S3
	if debtItem.ReceiptPhotoURL != nil && *debtItem.ReceiptPhotoURL != "" {
		if err := s.fileStorageService.DeleteReceipt(ctx, *debtItem.ReceiptPhotoURL); err != nil {
			// Log the error but don't fail the deletion
			// In production, you might want to handle this differently
			fmt.Printf("Warning: failed to delete receipt photo: %v\n", err)
		}
	}

	if err := s.debtItemRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete debt item: %w", err)
	}

	// Update debt list status, next payment date, and payment totals
	if err := s.updateDebtListStatusAndPaymentTotals(ctx, debtListID); err != nil {
		return fmt.Errorf("failed to update debt list totals: %w", err)
	}

	return nil
}

// Debt analytics and reporting

func (s *debtService) GetOverdueItems(ctx context.Context, userID uuid.UUID) ([]entities.DebtList, error) {
	debtLists, err := s.debtListRepo.GetOverdueForUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get overdue items: %w", err)
	}

	return debtLists, nil
}

func (s *debtService) GetDueSoonItems(ctx context.Context, userID uuid.UUID, days int) ([]entities.DebtList, error) {
	dueDate := time.Now().AddDate(0, 0, days)
	debtLists, err := s.debtListRepo.GetDueSoonForUser(ctx, userID, dueDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get due soon items: %w", err)
	}

	return debtLists, nil
}

func (s *debtService) GetPaymentSchedule(ctx context.Context, debtListID uuid.UUID, userID uuid.UUID) ([]entities.PaymentScheduleItem, error) {
	// Check if debt list belongs to user
	belongs, err := s.debtListRepo.BelongsToUser(ctx, debtListID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify ownership: %w", err)
	}
	if !belongs {
		return nil, entities.ErrDebtListNotFound
	}

	// Get debt list
	debtList, err := s.debtListRepo.GetByID(ctx, debtListID)
	if err != nil {
		return nil, fmt.Errorf("failed to get debt list: %w", err)
	}

	// Get payments
	payments, err := s.debtItemRepo.GetCompletedPaymentsForDebtList(ctx, debtListID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payments: %w", err)
	}

	schedule := s.paymentScheduleService.CalculatePaymentSchedule(debtList, payments)
	return schedule, nil
}

func (s *debtService) GetUpcomingPayments(ctx context.Context, userID uuid.UUID, days int) ([]entities.UpcomingPayment, error) {
	debtLists, err := s.debtListRepo.GetUserDebtLists(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user debt lists: %w", err)
	}

	var upcomingPayments []entities.UpcomingPayment
	cutoffDate := time.Now().AddDate(0, 0, days)

	for _, debtList := range debtLists {
		if debtList.Status == "active" || debtList.Status == "overdue" {
			if debtList.NextPaymentDate.After(time.Now()) && debtList.NextPaymentDate.Before(cutoffDate) {
				upcomingPayment := entities.UpcomingPayment{
					DebtListID:      debtList.ID,
					ContactName:     debtList.Contact.Name,
					DebtType:        debtList.DebtType,
					NextPaymentDate: debtList.NextPaymentDate,
					Amount:          debtList.InstallmentAmount,
					Currency:        debtList.Currency,
					Description:     debtList.Description,
				}
				upcomingPayments = append(upcomingPayments, upcomingPayment)
			}
		}
	}

	return upcomingPayments, nil
}

func (s *debtService) GetTotalPaymentsForDebtList(ctx context.Context, debtListID uuid.UUID, userID uuid.UUID) (*entities.PaymentSummary, error) {
	// Check if debt list belongs to user
	belongs, err := s.debtListRepo.BelongsToUser(ctx, debtListID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify ownership: %w", err)
	}
	if !belongs {
		return nil, entities.ErrDebtListNotFound
	}

	// Get debt list
	debtList, err := s.debtListRepo.GetByID(ctx, debtListID)
	if err != nil {
		return nil, fmt.Errorf("failed to get debt list: %w", err)
	}

	// Get all completed payments for this debt list
	payments, err := s.debtItemRepo.GetCompletedPaymentsForDebtList(ctx, debtListID)
	if err != nil {
		return nil, fmt.Errorf("failed to get completed payments: %w", err)
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

	return &entities.PaymentSummary{
		DebtListID:       debtListID,
		TotalAmount:      debtList.TotalAmount,
		TotalPaid:        totalPaid,
		RemainingDebt:    remainingDebt,
		PercentagePaid:   percentagePaid,
		NumberOfPayments: len(payments),
		Payments:         payments,
	}, nil
}

// Payment verification operations

func (s *debtService) VerifyDebtItem(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *entities.VerifyDebtItemRequest) (*entities.DebtItem, error) {
	// Get the debt item to check if it exists and get the debt list ID
	debtItem, err := s.GetDebtItem(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	// Update the payment status and verification details
	if err := s.debtItemRepo.UpdatePaymentStatus(ctx, id, req.Status, userID, req.VerificationNotes); err != nil {
		return nil, fmt.Errorf("failed to update payment status: %w", err)
	}

	// Get the updated debt item
	updatedDebtItem, err := s.GetDebtItem(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	// Update debt list totals if payment is completed
	if req.Status == entities.PaymentStatusCompleted {
		if err := s.updateDebtListStatusAndPaymentTotals(ctx, debtItem.DebtListID); err != nil {
			return nil, fmt.Errorf("failed to update debt list totals: %w", err)
		}
	}

	return updatedDebtItem, nil
}

func (s *debtService) GetPendingVerifications(ctx context.Context, userID uuid.UUID) ([]entities.DebtItem, error) {
	return s.debtItemRepo.GetPendingVerifications(ctx, userID)
}

func (s *debtService) RejectDebtItem(ctx context.Context, id uuid.UUID, userID uuid.UUID, notes *string) (*entities.DebtItem, error) {
	// Check if the debt item exists and belongs to the user
	_, err := s.GetDebtItem(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	// Update the payment status to rejected
	if err := s.debtItemRepo.UpdatePaymentStatus(ctx, id, entities.PaymentStatusRejected, userID, notes); err != nil {
		return nil, fmt.Errorf("failed to reject payment: %w", err)
	}

	// Get the updated debt item
	updatedDebtItem, err := s.GetDebtItem(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	return updatedDebtItem, nil
}

// Helper methods

func (s *debtService) updateDebtListStatusAndPaymentTotals(ctx context.Context, debtListID uuid.UUID) error {
	// Get the debt list
	debtList, err := s.debtListRepo.GetByID(ctx, debtListID)
	if err != nil {
		return fmt.Errorf("failed to get debt list: %w", err)
	}

	// Calculate total payments made
	totalPaid, err := s.debtItemRepo.GetTotalPaidForDebtList(ctx, debtListID)
	if err != nil {
		return fmt.Errorf("failed to get total paid: %w", err)
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
	} else if time.Now().After(debtList.NextPaymentDate) {
		newStatus = "overdue"
	} else {
		newStatus = "active"
	}

	// Get the last payment date to calculate next payment
	lastPaymentDate, err := s.debtItemRepo.GetLastPaymentDate(ctx, debtListID)
	if err != nil {
		return fmt.Errorf("failed to get last payment date: %w", err)
	}

	// Calculate next payment date
	nextPaymentDate := s.paymentScheduleService.CalculateNextPaymentDate(debtList, lastPaymentDate)

	// Update debt list totals
	if err := s.debtListRepo.UpdatePaymentTotals(ctx, debtListID, totalPaid, remainingAmount); err != nil {
		return fmt.Errorf("failed to update payment totals: %w", err)
	}

	// Update status
	if err := s.debtListRepo.UpdateStatus(ctx, debtListID, newStatus); err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	// Update next payment date
	if err := s.debtListRepo.UpdateNextPaymentDate(ctx, debtListID, nextPaymentDate); err != nil {
		return fmt.Errorf("failed to update next payment date: %w", err)
	}

	return nil
}

// Validation methods

func (s *debtService) validateCreateDebtListRequest(req *entities.CreateDebtListRequest) error {
	if req.ContactID == uuid.Nil {
		return entities.ErrInvalidInput
	}
	if req.DebtType != "owed_to_me" && req.DebtType != "i_owe" {
		return entities.ErrInvalidDebtType
	}
	if req.TotalAmount == "" {
		return entities.ErrInvalidAmount
	}
	// Validate that either due_date or number_of_payments is provided
	hasDueDate := req.DueDate != nil
	hasNumberOfPayments := req.NumberOfPayments != nil && *req.NumberOfPayments > 0
	if !hasDueDate && !hasNumberOfPayments {
		return fmt.Errorf("either due_date or number_of_payments must be provided")
	}
	return nil
}

func (s *debtService) validateUpdateDebtListRequest(req *entities.UpdateDebtListRequest) error {
	if req.TotalAmount != nil && *req.TotalAmount == "" {
		return entities.ErrInvalidAmount
	}
	if req.Currency != nil && *req.Currency == "" {
		return entities.ErrInvalidCurrency
	}
	return nil
}

func (s *debtService) validateCreateDebtItemRequest(req *entities.CreateDebtItemRequest) error {
	if req.DebtListID == uuid.Nil {
		return entities.ErrInvalidInput
	}
	if req.Amount == "" {
		return entities.ErrInvalidAmount
	}
	if req.PaymentMethod == "" {
		return entities.ErrInvalidPaymentMethod
	}
	
	// Validate payment method is one of the allowed values
	validPaymentMethods := map[string]bool{
		"cash":           true,
		"bank_transfer":  true,
		"check":          true,
		"digital_wallet": true,
		"other":          true,
	}
	if !validPaymentMethods[req.PaymentMethod] {
		return entities.ErrInvalidPaymentMethod
	}
	
	return nil
}

func (s *debtService) validateUpdateDebtItemRequest(req *entities.UpdateDebtItemRequest) error {
	if req.Amount != nil && *req.Amount == "" {
		return entities.ErrInvalidAmount
	}
	if req.Currency != nil && *req.Currency == "" {
		return entities.ErrInvalidCurrency
	}
	if req.PaymentMethod != nil && *req.PaymentMethod == "" {
		return entities.ErrInvalidPaymentMethod
	}
	return nil
}






