package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// DebtList represents the core debt list entity
type DebtList struct {
	ID                  uuid.UUID
	UserID              uuid.UUID
	ContactID           uuid.UUID
	DebtType            string
	TotalAmount         decimal.Decimal
	InstallmentAmount   decimal.Decimal
	TotalPaymentsMade   decimal.Decimal
	TotalRemainingDebt  decimal.Decimal
	Currency            string
	Status              string
	DueDate             time.Time
	NextPaymentDate     time.Time
	InstallmentPlan     string
	NumberOfPayments    *int
	Description         *string
	Notes               *string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// DebtItem represents the core debt item (payment) entity
type DebtItem struct {
	ID            uuid.UUID
	DebtListID    uuid.UUID
	Amount        decimal.Decimal
	Currency      string
	PaymentDate   time.Time
	PaymentMethod string
	Description   *string
	Status        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// CreateDebtListRequest represents a request to create a new debt list
type CreateDebtListRequest struct {
	ContactID        uuid.UUID  `json:"contact_id" validate:"required"`
	DebtType         string     `json:"debt_type" validate:"required,oneof=owed_to_me i_owe"`
	TotalAmount      string     `json:"total_amount" validate:"required"`
	Currency         string     `json:"currency"`
	DueDate          *time.Time `json:"due_date"`
	InstallmentPlan  string     `json:"installment_plan" validate:"omitempty,oneof=onetime weekly biweekly monthly quarterly yearly"`
	NumberOfPayments *int       `json:"number_of_payments"`
	Description      *string    `json:"description"`
	Notes            *string    `json:"notes"`
}

// UpdateDebtListRequest represents a request to update a debt list
type UpdateDebtListRequest struct {
	TotalAmount      *string    `json:"total_amount"`
	Currency         *string    `json:"currency"`
	Status           *string    `json:"status" validate:"omitempty,oneof=active settled archived overdue"`
	DueDate          *time.Time `json:"due_date"`
	InstallmentPlan  *string    `json:"installment_plan" validate:"omitempty,oneof=onetime weekly biweekly monthly quarterly yearly"`
	NumberOfPayments *int       `json:"number_of_payments"`
	Description      *string    `json:"description"`
	Notes            *string    `json:"notes"`
}

// CreateDebtItemRequest represents a request to create a new debt item (payment)
type CreateDebtItemRequest struct {
	DebtListID    uuid.UUID `json:"debt_list_id" validate:"required"`
	Amount        string    `json:"amount" validate:"required"`
	Currency      string    `json:"currency"`
	PaymentDate   time.Time `json:"payment_date" validate:"required"`
	PaymentMethod string    `json:"payment_method" validate:"required,oneof=cash bank_transfer check digital_wallet other"`
	Description   *string   `json:"description"`
}

// UpdateDebtItemRequest represents a request to update a debt item
type UpdateDebtItemRequest struct {
	Amount        *string    `json:"amount"`
	Currency      *string    `json:"currency"`
	PaymentDate   *time.Time `json:"payment_date"`
	PaymentMethod *string    `json:"payment_method" validate:"omitempty,oneof=cash bank_transfer check digital_wallet other"`
	Description   *string    `json:"description"`
	Status        *string    `json:"status" validate:"omitempty,oneof=completed pending failed refunded"`
}

// PaymentScheduleItem represents a scheduled payment
type PaymentScheduleItem struct {
	PaymentNumber int             `json:"payment_number"`
	DueDate       time.Time       `json:"due_date"`
	Amount        decimal.Decimal `json:"amount"`
	Status        string          `json:"status"` // pending, paid, overdue, missed
}

// PaymentSummary represents a summary of payments for a debt list
type PaymentSummary struct {
	DebtListID       uuid.UUID       `json:"debt_list_id"`
	TotalAmount      decimal.Decimal `json:"total_amount"`
	TotalPaid        decimal.Decimal `json:"total_paid"`
	RemainingDebt    decimal.Decimal `json:"remaining_debt"`
	PercentagePaid   decimal.Decimal `json:"percentage_paid"`
	NumberOfPayments int             `json:"number_of_payments"`
	Payments         []DebtItem      `json:"payments"`
}

// UpcomingPayment represents an upcoming payment
type UpcomingPayment struct {
	DebtListID      uuid.UUID       `json:"debt_list_id"`
	ContactName     string          `json:"contact_name"`
	DebtType        string          `json:"debt_type"`
	NextPaymentDate time.Time       `json:"next_payment_date"`
	Amount          decimal.Decimal `json:"amount"`
	Currency        string          `json:"currency"`
	Description     *string         `json:"description"`
}

// DebtListResponse represents a debt list response with related data
type DebtListResponse struct {
	ID                  uuid.UUID    `json:"id"`
	UserID              uuid.UUID    `json:"user_id"`
	ContactID           uuid.UUID    `json:"contact_id"`
	DebtType            string       `json:"debt_type"`
	TotalAmount         decimal.Decimal `json:"total_amount"`
	InstallmentAmount   decimal.Decimal `json:"installment_amount"`
	TotalPaymentsMade   decimal.Decimal `json:"total_payments_made"`
	TotalRemainingDebt  decimal.Decimal `json:"total_remaining_debt"`
	Currency            string       `json:"currency"`
	Status              string       `json:"status"`
	DueDate             time.Time    `json:"due_date"`
	NextPaymentDate     time.Time    `json:"next_payment_date"`
	InstallmentPlan     string       `json:"installment_plan"`
	NumberOfPayments    *int         `json:"number_of_payments"`
	Description         *string      `json:"description"`
	Notes               *string      `json:"notes"`
	CreatedAt           time.Time    `json:"created_at"`
	UpdatedAt           time.Time    `json:"updated_at"`
	Contact             Contact      `json:"contact,omitempty"`
	Payments            []DebtItem   `json:"payments,omitempty"`
}

// IsValid validates the debt list entity
func (d *DebtList) IsValid() error {
	if d.UserID == uuid.Nil {
		return ErrInvalidInput
	}
	if d.ContactID == uuid.Nil {
		return ErrInvalidInput
	}
	if d.DebtType != "owed_to_me" && d.DebtType != "i_owe" {
		return ErrInvalidDebtType
	}
	if d.TotalAmount.LessThanOrEqual(decimal.Zero) {
		return ErrInvalidAmount
	}
	if d.Currency == "" {
		return ErrInvalidCurrency
	}
	return nil
}

// IsValid validates the debt item entity
func (d *DebtItem) IsValid() error {
	if d.DebtListID == uuid.Nil {
		return ErrInvalidInput
	}
	if d.Amount.LessThanOrEqual(decimal.Zero) {
		return ErrInvalidAmount
	}
	if d.Currency == "" {
		return ErrInvalidCurrency
	}
	if d.PaymentMethod == "" {
		return ErrInvalidPaymentMethod
	}
	return nil
}

// IsSettled checks if the debt is fully settled
func (d *DebtList) IsSettled() bool {
	return d.TotalRemainingDebt.LessThanOrEqual(decimal.Zero) || d.Status == "settled"
}

// IsOverdue checks if the debt is overdue
func (d *DebtList) IsOverdue() bool {
	return time.Now().After(d.NextPaymentDate) && !d.IsSettled()
}

// CalculateProgress returns the payment progress as a percentage
func (d *DebtList) CalculateProgress() decimal.Decimal {
	if d.TotalAmount.LessThanOrEqual(decimal.Zero) {
		return decimal.Zero
	}
	return d.TotalPaymentsMade.Div(d.TotalAmount).Mul(decimal.NewFromInt(100))
}






