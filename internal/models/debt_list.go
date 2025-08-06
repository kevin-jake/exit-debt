package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Struct-level validation for CreateDebtListRequest
func (req CreateDebtListRequest) Validate() error {
	hasDueDate := req.DueDate != nil
	hasNumberOfPayments := req.NumberOfPayments != nil && *req.NumberOfPayments > 0
	
	if !hasDueDate && !hasNumberOfPayments {
		return fmt.Errorf("either due_date or number_of_payments must be provided")
	}
	
	return nil
}

type DebtList struct {
	ID              uuid.UUID     `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID          uuid.UUID     `json:"user_id" gorm:"type:uuid;not null;index"`
	ContactID       uuid.UUID     `json:"contact_id" gorm:"type:uuid;not null;index"`
	DebtType        string        `json:"debt_type" gorm:"not null;index;check:debt_type IN ('owed_to_me', 'i_owe')"`
	TotalAmount     decimal.Decimal `json:"total_amount" gorm:"type:decimal(15,2);not null"`
	InstallmentAmount decimal.Decimal `json:"installment_amount" gorm:"type:decimal(15,2);not null"`
	TotalPaymentsMade decimal.Decimal `json:"total_payments_made" gorm:"type:decimal(15,2);default:0"`
	TotalRemainingDebt decimal.Decimal `json:"total_remaining_debt" gorm:"type:decimal(15,2);not null"`
	Currency        string        `json:"currency" gorm:"default:'Php'"`
	Status          string        `json:"status" gorm:"default:'active';index;check:status IN ('active', 'settled', 'archived', 'overdue')"`
	DueDate         time.Time     `json:"due_date" gorm:"not null"`
	NextPaymentDate time.Time     `json:"next_payment_date" gorm:"not null"`
	InstallmentPlan string        `json:"installment_plan" gorm:"default:'monthly';check:installment_plan IN ('weekly', 'biweekly', 'monthly', 'quarterly', 'yearly')"`
	NumberOfPayments *int         `json:"number_of_payments" gorm:"default:null"`
	Description     *string       `json:"description"`
	Notes           *string       `json:"notes"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
	
	// Relationships
	User      User        `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Contact   Contact     `json:"contact,omitempty" gorm:"foreignKey:ContactID"`
	Payments  []DebtItem  `json:"payments,omitempty" gorm:"foreignKey:DebtListID"`
}

type CreateDebtListRequest struct {
	ContactID         uuid.UUID `json:"contact_id" binding:"required"`
	DebtType          string    `json:"debt_type" binding:"required,oneof=owed_to_me i_owe"`
	TotalAmount       string    `json:"total_amount" binding:"required"`
	Currency          string    `json:"currency"`
	DueDate           *time.Time `json:"due_date"`
	InstallmentPlan   string    `json:"installment_plan" binding:"required,oneof=weekly biweekly monthly quarterly yearly"`
	NumberOfPayments  *int      `json:"number_of_payments"`
	Description       *string   `json:"description"`
	Notes             *string   `json:"notes"`
}

type UpdateDebtListRequest struct {
	TotalAmount       *string    `json:"total_amount"`
	Currency          *string    `json:"currency"`
	Status            *string    `json:"status"`
	DueDate           *time.Time `json:"due_date"`
	InstallmentPlan   *string    `json:"installment_plan"`
	NumberOfPayments  *int       `json:"number_of_payments"`
	Description       *string    `json:"description"`
	Notes             *string    `json:"notes"`
} 