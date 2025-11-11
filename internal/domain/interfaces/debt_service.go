package interfaces

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"pay-your-dues/internal/domain/entities"
)

// DebtService defines the interface for debt management operations
type DebtService interface {
	// Debt List operations
	CreateDebtList(ctx context.Context, userID uuid.UUID, req *entities.CreateDebtListRequest) (*entities.DebtList, error)
	GetDebtList(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*entities.DebtListResponse, error)
	GetUserDebtLists(ctx context.Context, userID uuid.UUID) ([]entities.DebtListResponse, error)
	UpdateDebtList(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *entities.UpdateDebtListRequest) (*entities.DebtList, error)
	DeleteDebtList(ctx context.Context, id uuid.UUID, userID uuid.UUID) error

	// Debt Item (Payment) operations
	CreateDebtItem(ctx context.Context, userID uuid.UUID, req *entities.CreateDebtItemRequest) (*entities.DebtItem, error)
	GetDebtItem(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*entities.DebtItem, error)
	GetDebtListItems(ctx context.Context, debtListID uuid.UUID, userID uuid.UUID) ([]entities.DebtItem, error)
	UpdateDebtItem(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *entities.UpdateDebtItemRequest) (*entities.DebtItem, error)
	DeleteDebtItem(ctx context.Context, id uuid.UUID, userID uuid.UUID) error

	// Payment verification operations
	VerifyDebtItem(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *entities.VerifyDebtItemRequest) (*entities.DebtItem, error)
	GetPendingVerifications(ctx context.Context, userID uuid.UUID) ([]entities.DebtItem, error)
	RejectDebtItem(ctx context.Context, id uuid.UUID, userID uuid.UUID, notes *string) (*entities.DebtItem, error)

	// Debt analytics and reporting
	GetOverdueItems(ctx context.Context, userID uuid.UUID) ([]entities.DebtList, error)
	GetDueSoonItems(ctx context.Context, userID uuid.UUID, days int) ([]entities.DebtList, error)
	GetPaymentSchedule(ctx context.Context, debtListID uuid.UUID, userID uuid.UUID) ([]entities.PaymentScheduleItem, error)
	GetUpcomingPayments(ctx context.Context, userID uuid.UUID, days int) ([]entities.UpcomingPayment, error)
	GetTotalPaymentsForDebtList(ctx context.Context, debtListID uuid.UUID, userID uuid.UUID) (*entities.PaymentSummary, error)
}

// PaymentScheduleService defines the interface for payment schedule calculations
type PaymentScheduleService interface {
	CalculateNextPaymentDate(debtList *entities.DebtList, lastPaymentDate *time.Time) time.Time
	CalculatePaymentSchedule(debtList *entities.DebtList, payments []entities.DebtItem) []entities.PaymentScheduleItem
	CalculateDueDateFromNumberOfPayments(createdAt time.Time, numberOfPayments int, installmentPlan string) time.Time
	CalculateInstallmentAmountFromNumberOfPayments(totalAmount decimal.Decimal, numberOfPayments int) decimal.Decimal
	CalculateInstallmentAmount(totalAmount decimal.Decimal, installmentPlan string, createdAt time.Time, dueDate time.Time) decimal.Decimal
}
