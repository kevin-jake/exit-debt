package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type DebtItem struct {
	ID                uuid.UUID     `json:"id" gorm:"type:uuid;primary_key"`
	UserID            uuid.UUID     `json:"user_id" gorm:"type:uuid;not null;index"`
	DebtListID        uuid.UUID     `json:"debt_list_id" gorm:"type:uuid;not null;index"`
	Amount            decimal.Decimal `json:"amount" gorm:"type:decimal(15,2);not null"`
	Currency          string        `json:"currency" gorm:"default:'Php'"`
	PaymentDate       time.Time     `json:"payment_date" gorm:"not null"`
	PaymentMethod     string        `json:"payment_method" gorm:"default:'cash';check:payment_method IN ('cash', 'bank_transfer', 'check', 'digital_wallet', 'other')"`
	Description       *string       `json:"description"`
	Status            string        `json:"status" gorm:"default:'pending';index;check:status IN ('completed', 'pending', 'failed', 'refunded', 'rejected')"`
	ReceiptPhotoURL   *string       `json:"receipt_photo_url"`
	VerifiedBy        *uuid.UUID    `json:"verified_by" gorm:"type:uuid"`
	VerifiedAt        *time.Time    `json:"verified_at"`
	VerificationNotes *string       `json:"verification_notes"`
	CreatedAt         time.Time     `json:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at"`
	
	// Relationships
	DebtList DebtList `json:"debt_list,omitempty" gorm:"foreignKey:DebtListID"`
}

type CreateDebtItemRequest struct {
	DebtListID        uuid.UUID `json:"debt_list_id" binding:"required"`
	Amount            string    `json:"amount" binding:"required"`
	Currency          string    `json:"currency"`
	PaymentDate       time.Time `json:"payment_date" binding:"required"`
	PaymentMethod     string    `json:"payment_method" binding:"required,oneof=cash bank_transfer check digital_wallet other"`
	Description       *string   `json:"description"`
	ReceiptPhotoURL   *string   `json:"receipt_photo_url"`
	VerificationNotes *string   `json:"verification_notes"`
}

type UpdateDebtItemRequest struct {
	Amount            *string    `json:"amount"`
	Currency          *string    `json:"currency"`
	PaymentDate       *time.Time `json:"payment_date"`
	PaymentMethod     *string    `json:"payment_method"`
	Description       *string    `json:"description"`
	Status            *string    `json:"status"`
	ReceiptPhotoURL   *string    `json:"receipt_photo_url"`
	VerificationNotes *string    `json:"verification_notes"`
}

// VerifyDebtItemRequest represents a request to verify a debt item
type VerifyDebtItemRequest struct {
	Status            string  `json:"status" binding:"required,oneof=completed rejected"`
	VerificationNotes *string `json:"verification_notes"`
}

// DebtItemVerificationResponse represents verification details for a debt item
type DebtItemVerificationResponse struct {
	ID                uuid.UUID  `json:"id"`
	Status            string     `json:"status"`
	VerifiedBy        *uuid.UUID `json:"verified_by"`
	VerifiedAt        *time.Time `json:"verified_at"`
	VerificationNotes *string    `json:"verification_notes"`
	ReceiptPhotoURL   *string    `json:"receipt_photo_url"`
} 