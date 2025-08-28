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

// debtItemRepositoryGORM implements the DebtItemRepository interface using GORM
type debtItemRepositoryGORM struct {
	db *gorm.DB
}

// NewDebtItemRepositoryGORM creates a new debt item repository with GORM
func NewDebtItemRepositoryGORM(db *gorm.DB) interfaces.DebtItemRepository {
	return &debtItemRepositoryGORM{
		db: db,
	}
}

func (r *debtItemRepositoryGORM) Create(ctx context.Context, debtItem *entities.DebtItem) error {
	gormDebtItem := r.entityToGORM(debtItem)
	if err := r.db.WithContext(ctx).Create(gormDebtItem).Error; err != nil {
		return fmt.Errorf("failed to create debt item: %w", err)
	}
	// Update the entity with the created ID and timestamps
	debtItem.ID = gormDebtItem.ID
	debtItem.CreatedAt = gormDebtItem.CreatedAt
	debtItem.UpdatedAt = gormDebtItem.UpdatedAt
	return nil
}

func (r *debtItemRepositoryGORM) GetByID(ctx context.Context, id uuid.UUID) (*entities.DebtItem, error) {
	var gormDebtItem models.DebtItem
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&gormDebtItem).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrDebtItemNotFound
		}
		return nil, fmt.Errorf("failed to get debt item by ID: %w", err)
	}
	return r.gormToEntity(&gormDebtItem), nil
}

func (r *debtItemRepositoryGORM) GetByDebtListID(ctx context.Context, debtListID uuid.UUID) ([]entities.DebtItem, error) {
	var gormDebtItems []models.DebtItem
	if err := r.db.WithContext(ctx).
		Where("debt_list_id = ?", debtListID).
		Order("payment_date DESC").
		Find(&gormDebtItems).Error; err != nil {
		return nil, fmt.Errorf("failed to get debt items by debt list ID: %w", err)
	}

	debtItems := make([]entities.DebtItem, len(gormDebtItems))
	for i, gormDebtItem := range gormDebtItems {
		debtItems[i] = *r.gormToEntity(&gormDebtItem)
	}

	return debtItems, nil
}

func (r *debtItemRepositoryGORM) Update(ctx context.Context, debtItem *entities.DebtItem) error {
	gormDebtItem := r.entityToGORM(debtItem)
	if err := r.db.WithContext(ctx).Save(gormDebtItem).Error; err != nil {
		return fmt.Errorf("failed to update debt item: %w", err)
	}
	// Update the entity with the updated timestamp
	debtItem.UpdatedAt = gormDebtItem.UpdatedAt
	return nil
}

func (r *debtItemRepositoryGORM) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.DebtItem{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete debt item: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return entities.ErrDebtItemNotFound
	}
	return nil
}

func (r *debtItemRepositoryGORM) GetTotalPaidForDebtList(ctx context.Context, debtListID uuid.UUID) (decimal.Decimal, error) {
	var totalPaid decimal.Decimal
	if err := r.db.WithContext(ctx).Model(&models.DebtItem{}).
		Where("debt_list_id = ? AND status = ?", debtListID, "completed").
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalPaid).Error; err != nil {
		return decimal.Zero, fmt.Errorf("failed to get total paid for debt list: %w", err)
	}
	return totalPaid, nil
}

func (r *debtItemRepositoryGORM) GetCompletedPaymentsForDebtList(ctx context.Context, debtListID uuid.UUID) ([]entities.DebtItem, error) {
	var gormDebtItems []models.DebtItem
	if err := r.db.WithContext(ctx).
		Where("debt_list_id = ? AND status = ?", debtListID, "completed").
		Order("payment_date ASC").
		Find(&gormDebtItems).Error; err != nil {
		return nil, fmt.Errorf("failed to get completed payments for debt list: %w", err)
	}

	debtItems := make([]entities.DebtItem, len(gormDebtItems))
	for i, gormDebtItem := range gormDebtItems {
		debtItems[i] = *r.gormToEntity(&gormDebtItem)
	}

	return debtItems, nil
}

func (r *debtItemRepositoryGORM) GetLastPaymentDate(ctx context.Context, debtListID uuid.UUID) (*time.Time, error) {
	var lastPayment models.DebtItem
	if err := r.db.WithContext(ctx).
		Where("debt_list_id = ? AND status = ?", debtListID, "completed").
		Order("payment_date DESC").
		First(&lastPayment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // No payments found
		}
		return nil, fmt.Errorf("failed to get last payment date: %w", err)
	}
	return &lastPayment.PaymentDate, nil
}

func (r *debtItemRepositoryGORM) BelongsToUserDebtList(ctx context.Context, debtItemID, userID uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.DebtItem{}).
		Joins("JOIN debt_lists ON debt_items.debt_list_id = debt_lists.id").
		Where("debt_items.id = ? AND debt_lists.user_id = ?", debtItemID, userID).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check debt item ownership: %w", err)
	}
	return count > 0, nil
}

// GetPendingVerifications gets all pending debt items that need verification
func (r *debtItemRepositoryGORM) GetPendingVerifications(ctx context.Context, userID uuid.UUID) ([]entities.DebtItem, error) {
	var gormDebtItems []models.DebtItem
	if err := r.db.WithContext(ctx).
		Joins("JOIN debt_lists ON debt_items.debt_list_id = debt_lists.id").
		Where("debt_lists.user_id = ? AND debt_items.status = ?", userID, "pending").
		Order("debt_items.created_at DESC").
		Find(&gormDebtItems).Error; err != nil {
		return nil, fmt.Errorf("failed to get pending verifications: %w", err)
	}

	debtItems := make([]entities.DebtItem, len(gormDebtItems))
	for i, gormDebtItem := range gormDebtItems {
		debtItems[i] = *r.gormToEntity(&gormDebtItem)
	}

	return debtItems, nil
}

// UpdatePaymentStatus updates the payment status and verification details
func (r *debtItemRepositoryGORM) UpdatePaymentStatus(ctx context.Context, debtItemID uuid.UUID, status string, verifiedBy uuid.UUID, notes *string) error {
	updates := map[string]interface{}{
		"status":             status,
		"verified_by":        verifiedBy,
		"verified_at":        time.Now(),
		"verification_notes": notes,
		"updated_at":         time.Now(),
	}

	result := r.db.WithContext(ctx).Model(&models.DebtItem{}).
		Where("id = ?", debtItemID).
		Updates(updates)

	if result.Error != nil {
		return fmt.Errorf("failed to update payment status: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return entities.ErrDebtItemNotFound
	}
	return nil
}

// UpdateReceiptPhoto updates the receipt photo URL for a debt item
func (r *debtItemRepositoryGORM) UpdateReceiptPhoto(ctx context.Context, debtItemID uuid.UUID, photoURL *string) error {
	updates := map[string]interface{}{
		"receipt_photo_url": photoURL,
		"updated_at":        time.Now(),
	}

	result := r.db.WithContext(ctx).Model(&models.DebtItem{}).
		Where("id = ?", debtItemID).
		Updates(updates)

	if result.Error != nil {
		return fmt.Errorf("failed to update receipt photo: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return entities.ErrDebtItemNotFound
	}
	return nil
}

// entityToGORM converts a domain entity to GORM model
func (r *debtItemRepositoryGORM) entityToGORM(debtItem *entities.DebtItem) *models.DebtItem {
	return &models.DebtItem{
		ID:                debtItem.ID,
		DebtListID:        debtItem.DebtListID,
		Amount:            debtItem.Amount,
		Currency:          debtItem.Currency,
		PaymentDate:       debtItem.PaymentDate,
		PaymentMethod:     debtItem.PaymentMethod,
		Description:       debtItem.Description,
		Status:            debtItem.Status,
		ReceiptPhotoURL:   debtItem.ReceiptPhotoURL,
		VerifiedBy:        debtItem.VerifiedBy,
		VerifiedAt:        debtItem.VerifiedAt,
		VerificationNotes: debtItem.VerificationNotes,
		CreatedAt:         debtItem.CreatedAt,
		UpdatedAt:         debtItem.UpdatedAt,
	}
}

// gormToEntity converts a GORM model to domain entity
func (r *debtItemRepositoryGORM) gormToEntity(gormDebtItem *models.DebtItem) *entities.DebtItem {
	return &entities.DebtItem{
		ID:                gormDebtItem.ID,
		DebtListID:        gormDebtItem.DebtListID,
		Amount:            gormDebtItem.Amount,
		Currency:          gormDebtItem.Currency,
		PaymentDate:       gormDebtItem.PaymentDate,
		PaymentMethod:     gormDebtItem.PaymentMethod,
		Description:       gormDebtItem.Description,
		Status:            gormDebtItem.Status,
		ReceiptPhotoURL:   gormDebtItem.ReceiptPhotoURL,
		VerifiedBy:        gormDebtItem.VerifiedBy,
		VerifiedAt:        gormDebtItem.VerifiedAt,
		VerificationNotes: gormDebtItem.VerificationNotes,
		CreatedAt:         gormDebtItem.CreatedAt,
		UpdatedAt:         gormDebtItem.UpdatedAt,
	}
}






