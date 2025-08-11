package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"exit-debt/internal/domain/entities"
	"exit-debt/internal/domain/interfaces"
	"exit-debt/internal/models"
)

// debtListRepositoryGORM implements the DebtListRepository interface using GORM
type debtListRepositoryGORM struct {
	db *gorm.DB
}

// NewDebtListRepositoryGORM creates a new debt list repository with GORM
func NewDebtListRepositoryGORM(db *gorm.DB) interfaces.DebtListRepository {
	return &debtListRepositoryGORM{
		db: db,
	}
}

func (r *debtListRepositoryGORM) Create(ctx context.Context, debtList *entities.DebtList) error {
	gormDebtList := r.entityToGORM(debtList)
	if err := r.db.WithContext(ctx).Create(gormDebtList).Error; err != nil {
		return fmt.Errorf("failed to create debt list: %w", err)
	}
	// Update the entity with the created ID and timestamps
	debtList.ID = gormDebtList.ID
	debtList.CreatedAt = gormDebtList.CreatedAt
	debtList.UpdatedAt = gormDebtList.UpdatedAt
	return nil
}

func (r *debtListRepositoryGORM) GetByID(ctx context.Context, id uuid.UUID) (*entities.DebtList, error) {
	var gormDebtList models.DebtList
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&gormDebtList).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrDebtListNotFound
		}
		return nil, fmt.Errorf("failed to get debt list by ID: %w", err)
	}
	return r.gormToEntity(&gormDebtList), nil
}

func (r *debtListRepositoryGORM) GetByIDWithRelations(ctx context.Context, id uuid.UUID) (*entities.DebtListResponse, error) {
	var gormDebtList models.DebtList
	if err := r.db.WithContext(ctx).
		Preload("Contact").
		Preload("User").
		Preload("Payments").
		Where("id = ?", id).
		First(&gormDebtList).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrDebtListNotFound
		}
		return nil, fmt.Errorf("failed to get debt list with relations: %w", err)
	}
	return r.gormToResponseEntity(&gormDebtList), nil
}

func (r *debtListRepositoryGORM) GetUserDebtLists(ctx context.Context, userID uuid.UUID) ([]entities.DebtListResponse, error) {
	var gormDebtLists []models.DebtList
	if err := r.db.WithContext(ctx).
		Preload("Contact").
		Preload("User").
		Preload("Payments").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&gormDebtLists).Error; err != nil {
		return nil, fmt.Errorf("failed to get user debt lists: %w", err)
	}

	debtLists := make([]entities.DebtListResponse, len(gormDebtLists))
	for i, gormDebtList := range gormDebtLists {
		debtLists[i] = *r.gormToResponseEntity(&gormDebtList)
	}

	return debtLists, nil
}

func (r *debtListRepositoryGORM) GetDebtListsWhereUserIsContact(ctx context.Context, userID uuid.UUID) ([]entities.DebtListResponse, error) {
	// Find debt lists where the current user is referenced as a contact
	var gormDebtLists []models.DebtList
	if err := r.db.WithContext(ctx).
		Preload("Contact").
		Preload("User").
		Preload("Payments").
		Joins("JOIN contacts ON debt_lists.contact_id = contacts.id").
		Where("contacts.user_id_ref = ?", userID).
		Order("created_at DESC").
		Find(&gormDebtLists).Error; err != nil {
		return nil, fmt.Errorf("failed to get debt lists where user is contact: %w", err)
	}

	debtLists := make([]entities.DebtListResponse, len(gormDebtLists))
	for i, gormDebtList := range gormDebtLists {
		debtLists[i] = *r.gormToResponseEntity(&gormDebtList)
	}

	return debtLists, nil
}

func (r *debtListRepositoryGORM) Update(ctx context.Context, debtList *entities.DebtList) error {
	gormDebtList := r.entityToGORM(debtList)
	if err := r.db.WithContext(ctx).Save(gormDebtList).Error; err != nil {
		return fmt.Errorf("failed to update debt list: %w", err)
	}
	// Update the entity with the updated timestamp
	debtList.UpdatedAt = gormDebtList.UpdatedAt
	return nil
}

func (r *debtListRepositoryGORM) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.DebtList{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete debt list: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return entities.ErrDebtListNotFound
	}
	return nil
}

func (r *debtListRepositoryGORM) GetOverdueForUser(ctx context.Context, userID uuid.UUID) ([]entities.DebtList, error) {
	var gormDebtLists []models.DebtList
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND next_payment_date < ? AND status = ?", userID, time.Now(), "active").
		Order("next_payment_date ASC").
		Find(&gormDebtLists).Error; err != nil {
		return nil, fmt.Errorf("failed to get overdue debt lists: %w", err)
	}

	debtLists := make([]entities.DebtList, len(gormDebtLists))
	for i, gormDebtList := range gormDebtLists {
		debtLists[i] = *r.gormToEntity(&gormDebtList)
	}

	return debtLists, nil
}

func (r *debtListRepositoryGORM) GetDueSoonForUser(ctx context.Context, userID uuid.UUID, dueDate time.Time) ([]entities.DebtList, error) {
	var gormDebtLists []models.DebtList
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND next_payment_date BETWEEN ? AND ? AND status = ?", userID, time.Now(), dueDate, "active").
		Order("next_payment_date ASC").
		Find(&gormDebtLists).Error; err != nil {
		return nil, fmt.Errorf("failed to get due soon debt lists: %w", err)
	}

	debtLists := make([]entities.DebtList, len(gormDebtLists))
	for i, gormDebtList := range gormDebtLists {
		debtLists[i] = *r.gormToEntity(&gormDebtList)
	}

	return debtLists, nil
}

func (r *debtListRepositoryGORM) BelongsToUser(ctx context.Context, debtListID, userID uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.DebtList{}).
		Where("id = ? AND user_id = ?", debtListID, userID).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check debt list ownership: %w", err)
	}
	return count > 0, nil
}

func (r *debtListRepositoryGORM) UpdatePaymentTotals(ctx context.Context, debtListID uuid.UUID, totalPaid, remaining decimal.Decimal) error {
	if err := r.db.WithContext(ctx).Model(&models.DebtList{}).
		Where("id = ?", debtListID).
		Updates(map[string]interface{}{
			"total_payments_made":   totalPaid,
			"total_remaining_debt":  remaining,
			"updated_at":           time.Now(),
		}).Error; err != nil {
		return fmt.Errorf("failed to update payment totals: %w", err)
	}
	return nil
}

func (r *debtListRepositoryGORM) UpdateStatus(ctx context.Context, debtListID uuid.UUID, status string) error {
	if err := r.db.WithContext(ctx).Model(&models.DebtList{}).
		Where("id = ?", debtListID).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error; err != nil {
		return fmt.Errorf("failed to update debt list status: %w", err)
	}
	return nil
}

func (r *debtListRepositoryGORM) UpdateNextPaymentDate(ctx context.Context, debtListID uuid.UUID, nextPaymentDate time.Time) error {
	if err := r.db.WithContext(ctx).Model(&models.DebtList{}).
		Where("id = ?", debtListID).
		Updates(map[string]interface{}{
			"next_payment_date": nextPaymentDate,
			"updated_at":        time.Now(),
		}).Error; err != nil {
		return fmt.Errorf("failed to update next payment date: %w", err)
	}
	return nil
}

// entityToGORM converts a domain entity to GORM model
func (r *debtListRepositoryGORM) entityToGORM(debtList *entities.DebtList) *models.DebtList {
	return &models.DebtList{
		ID:                  debtList.ID,
		UserID:              debtList.UserID,
		ContactID:           debtList.ContactID,
		DebtType:            debtList.DebtType,
		TotalAmount:         debtList.TotalAmount,
		InstallmentAmount:   debtList.InstallmentAmount,
		TotalPaymentsMade:   debtList.TotalPaymentsMade,
		TotalRemainingDebt:  debtList.TotalRemainingDebt,
		Currency:            debtList.Currency,
		Status:              debtList.Status,
		DueDate:             debtList.DueDate,
		NextPaymentDate:     debtList.NextPaymentDate,
		InstallmentPlan:     debtList.InstallmentPlan,
		NumberOfPayments:    debtList.NumberOfPayments,
		Description:         debtList.Description,
		Notes:               debtList.Notes,
		CreatedAt:           debtList.CreatedAt,
		UpdatedAt:           debtList.UpdatedAt,
	}
}

// gormToEntity converts a GORM model to domain entity
func (r *debtListRepositoryGORM) gormToEntity(gormDebtList *models.DebtList) *entities.DebtList {
	return &entities.DebtList{
		ID:                  gormDebtList.ID,
		UserID:              gormDebtList.UserID,
		ContactID:           gormDebtList.ContactID,
		DebtType:            gormDebtList.DebtType,
		TotalAmount:         gormDebtList.TotalAmount,
		InstallmentAmount:   gormDebtList.InstallmentAmount,
		TotalPaymentsMade:   gormDebtList.TotalPaymentsMade,
		TotalRemainingDebt:  gormDebtList.TotalRemainingDebt,
		Currency:            gormDebtList.Currency,
		Status:              gormDebtList.Status,
		DueDate:             gormDebtList.DueDate,
		NextPaymentDate:     gormDebtList.NextPaymentDate,
		InstallmentPlan:     gormDebtList.InstallmentPlan,
		NumberOfPayments:    gormDebtList.NumberOfPayments,
		Description:         gormDebtList.Description,
		Notes:               gormDebtList.Notes,
		CreatedAt:           gormDebtList.CreatedAt,
		UpdatedAt:           gormDebtList.UpdatedAt,
	}
}

// gormToResponseEntity converts a GORM model to response entity with relations
func (r *debtListRepositoryGORM) gormToResponseEntity(gormDebtList *models.DebtList) *entities.DebtListResponse {
	// Convert contact
	contact := entities.Contact{
		ID:         gormDebtList.Contact.ID,
		Name:       gormDebtList.Contact.Name,
		Email:      gormDebtList.Contact.Email,
		Phone:      gormDebtList.Contact.Phone,
		FacebookID: gormDebtList.Contact.FacebookID,
		Notes:      gormDebtList.Contact.Notes,
		IsUser:     gormDebtList.Contact.IsUser,
		UserIDRef:  gormDebtList.Contact.UserIDRef,
		CreatedAt:  gormDebtList.Contact.CreatedAt,
		UpdatedAt:  gormDebtList.Contact.UpdatedAt,
	}

	// Convert payments
	payments := make([]entities.DebtItem, len(gormDebtList.Payments))
	for i, payment := range gormDebtList.Payments {
		payments[i] = entities.DebtItem{
			ID:            payment.ID,
			DebtListID:    payment.DebtListID,
			Amount:        payment.Amount,
			Currency:      payment.Currency,
			PaymentDate:   payment.PaymentDate,
			PaymentMethod: payment.PaymentMethod,
			Description:   payment.Description,
			Status:        payment.Status,
			CreatedAt:     payment.CreatedAt,
			UpdatedAt:     payment.UpdatedAt,
		}
	}

	return &entities.DebtListResponse{
		ID:                  gormDebtList.ID,
		UserID:              gormDebtList.UserID,
		ContactID:           gormDebtList.ContactID,
		DebtType:            gormDebtList.DebtType,
		TotalAmount:         gormDebtList.TotalAmount,
		InstallmentAmount:   gormDebtList.InstallmentAmount,
		TotalPaymentsMade:   gormDebtList.TotalPaymentsMade,
		TotalRemainingDebt:  gormDebtList.TotalRemainingDebt,
		Currency:            gormDebtList.Currency,
		Status:              gormDebtList.Status,
		DueDate:             gormDebtList.DueDate,
		NextPaymentDate:     gormDebtList.NextPaymentDate,
		InstallmentPlan:     gormDebtList.InstallmentPlan,
		NumberOfPayments:    gormDebtList.NumberOfPayments,
		Description:         gormDebtList.Description,
		Notes:               gormDebtList.Notes,
		CreatedAt:           gormDebtList.CreatedAt,
		UpdatedAt:           gormDebtList.UpdatedAt,
		Contact:             contact,
		Payments:            payments,
	}
}






