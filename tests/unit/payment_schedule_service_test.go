package unit

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"pay-your-dues/internal/domain/entities"
	"pay-your-dues/internal/services"
)

func TestCalculatePaymentSchedule(t *testing.T) {
	tests := []struct {
		name              string
		debtList          *entities.DebtList
		payments          []entities.DebtItem
		expectedSchedule  int // number of schedule items
		validateSchedule  func(t *testing.T, schedule []entities.PaymentScheduleItem, debtList *entities.DebtList)
	}{
		{
			name: "no payments made - all pending",
			debtList: &entities.DebtList{
				ID:                uuid.New(),
				TotalAmount:       decimal.RequireFromString("1000.00"),
				InstallmentAmount: decimal.RequireFromString("250.00"),
				InstallmentPlan:   "monthly",
				CreatedAt:         time.Now(),
				DueDate:           time.Now().AddDate(0, 4, 0),
			},
			payments:         []entities.DebtItem{},
			expectedSchedule: 4,
			validateSchedule: func(t *testing.T, schedule []entities.PaymentScheduleItem, debtList *entities.DebtList) {
				assert.Len(t, schedule, 4)
				
				// All payments should be pending with no paid amount
				for i, item := range schedule {
					assert.Equal(t, i+1, item.PaymentNumber, "Payment number should be sequential")
					assert.Equal(t, "pending", item.Status, "All payments should be pending")
					assert.True(t, item.PaidAmount.IsZero(), "Paid amount should be zero")
					assert.True(t, item.ScheduledAmount.Equal(decimal.RequireFromString("250.00")), "Scheduled amount should be 250")
					assert.True(t, item.Amount.Equal(decimal.RequireFromString("250.00")), "Amount (remaining) should be 250")
				}
			},
		},
		{
			name: "one full payment made",
			debtList: &entities.DebtList{
				ID:                uuid.New(),
				TotalAmount:       decimal.RequireFromString("1000.00"),
				InstallmentAmount: decimal.RequireFromString("250.00"),
				InstallmentPlan:   "monthly",
				CreatedAt:         time.Now(),
				DueDate:           time.Now().AddDate(0, 4, 0),
			},
			payments: []entities.DebtItem{
				{
					ID:          uuid.New(),
					Amount:      decimal.RequireFromString("250.00"),
					Status:      "completed",
					PaymentDate: time.Now(),
				},
			},
			expectedSchedule: 4,
			validateSchedule: func(t *testing.T, schedule []entities.PaymentScheduleItem, debtList *entities.DebtList) {
				assert.Len(t, schedule, 4)
				
				// First payment should be paid
				assert.Equal(t, "paid", schedule[0].Status)
				assert.True(t, schedule[0].PaidAmount.Equal(decimal.RequireFromString("250.00")))
				assert.True(t, schedule[0].ScheduledAmount.Equal(decimal.RequireFromString("250.00")))
				assert.True(t, schedule[0].Amount.IsZero(), "Remaining amount should be zero")
				
				// Rest should be pending
				for i := 1; i < 4; i++ {
					assert.Equal(t, "pending", schedule[i].Status)
					assert.True(t, schedule[i].PaidAmount.IsZero())
					assert.True(t, schedule[i].ScheduledAmount.Equal(decimal.RequireFromString("250.00")))
					assert.True(t, schedule[i].Amount.Equal(decimal.RequireFromString("250.00")))
				}
			},
		},
		{
			name: "partial payment made",
			debtList: &entities.DebtList{
				ID:                uuid.New(),
				TotalAmount:       decimal.RequireFromString("1000.00"),
				InstallmentAmount: decimal.RequireFromString("250.00"),
				InstallmentPlan:   "monthly",
				CreatedAt:         time.Now(),
				DueDate:           time.Now().AddDate(0, 4, 0),
			},
			payments: []entities.DebtItem{
				{
					ID:          uuid.New(),
					Amount:      decimal.RequireFromString("150.00"),
					Status:      "completed",
					PaymentDate: time.Now(),
				},
			},
			expectedSchedule: 4,
			validateSchedule: func(t *testing.T, schedule []entities.PaymentScheduleItem, debtList *entities.DebtList) {
				assert.Len(t, schedule, 4)
				
				// First payment should be partially paid
				assert.Equal(t, "pending", schedule[0].Status, "Status should be pending for partial payment")
				assert.True(t, schedule[0].PaidAmount.Equal(decimal.RequireFromString("150.00")), "Paid amount should be 150")
				assert.True(t, schedule[0].ScheduledAmount.Equal(decimal.RequireFromString("250.00")), "Scheduled amount should be 250")
				assert.True(t, schedule[0].Amount.Equal(decimal.RequireFromString("100.00")), "Remaining amount should be 100")
				
				// Rest should be pending with no payment
				for i := 1; i < 4; i++ {
					assert.Equal(t, "pending", schedule[i].Status)
					assert.True(t, schedule[i].PaidAmount.IsZero())
					assert.True(t, schedule[i].ScheduledAmount.Equal(decimal.RequireFromString("250.00")))
					assert.True(t, schedule[i].Amount.Equal(decimal.RequireFromString("250.00")))
				}
			},
		},
		{
			name: "multiple payments - some full, one partial",
			debtList: &entities.DebtList{
				ID:                uuid.New(),
				TotalAmount:       decimal.RequireFromString("1000.00"),
				InstallmentAmount: decimal.RequireFromString("250.00"),
				InstallmentPlan:   "monthly",
				CreatedAt:         time.Now(),
				DueDate:           time.Now().AddDate(0, 4, 0),
			},
			payments: []entities.DebtItem{
				{
					ID:          uuid.New(),
					Amount:      decimal.RequireFromString("250.00"),
					Status:      "completed",
					PaymentDate: time.Now().AddDate(0, 0, -60),
				},
				{
					ID:          uuid.New(),
					Amount:      decimal.RequireFromString("250.00"),
					Status:      "completed",
					PaymentDate: time.Now().AddDate(0, 0, -30),
				},
				{
					ID:          uuid.New(),
					Amount:      decimal.RequireFromString("100.00"),
					Status:      "completed",
					PaymentDate: time.Now(),
				},
			},
			expectedSchedule: 4,
			validateSchedule: func(t *testing.T, schedule []entities.PaymentScheduleItem, debtList *entities.DebtList) {
				assert.Len(t, schedule, 4)
				
				// First two payments should be fully paid
				for i := 0; i < 2; i++ {
					assert.Equal(t, "paid", schedule[i].Status)
					assert.True(t, schedule[i].PaidAmount.Equal(decimal.RequireFromString("250.00")))
					assert.True(t, schedule[i].ScheduledAmount.Equal(decimal.RequireFromString("250.00")))
					assert.True(t, schedule[i].Amount.IsZero())
				}
				
				// Third payment should be partially paid
				assert.Equal(t, "pending", schedule[2].Status)
				assert.True(t, schedule[2].PaidAmount.Equal(decimal.RequireFromString("100.00")))
				assert.True(t, schedule[2].ScheduledAmount.Equal(decimal.RequireFromString("250.00")))
				assert.True(t, schedule[2].Amount.Equal(decimal.RequireFromString("150.00")))
				
				// Fourth payment should be pending
				assert.Equal(t, "pending", schedule[3].Status)
				assert.True(t, schedule[3].PaidAmount.IsZero())
				assert.True(t, schedule[3].ScheduledAmount.Equal(decimal.RequireFromString("250.00")))
				assert.True(t, schedule[3].Amount.Equal(decimal.RequireFromString("250.00")))
			},
		},
		{
			name: "all payments completed",
			debtList: &entities.DebtList{
				ID:                uuid.New(),
				TotalAmount:       decimal.RequireFromString("1000.00"),
				InstallmentAmount: decimal.RequireFromString("250.00"),
				InstallmentPlan:   "monthly",
				CreatedAt:         time.Now(),
				DueDate:           time.Now().AddDate(0, 4, 0),
			},
			payments: []entities.DebtItem{
				{
					ID:          uuid.New(),
					Amount:      decimal.RequireFromString("1000.00"),
					Status:      "completed",
					PaymentDate: time.Now(),
				},
			},
			expectedSchedule: 4,
			validateSchedule: func(t *testing.T, schedule []entities.PaymentScheduleItem, debtList *entities.DebtList) {
				assert.Len(t, schedule, 4)
				
				// All payments should be marked as paid
				for _, item := range schedule {
					assert.Equal(t, "paid", item.Status)
					assert.True(t, item.PaidAmount.Equal(decimal.RequireFromString("250.00")))
					assert.True(t, item.ScheduledAmount.Equal(decimal.RequireFromString("250.00")))
					assert.True(t, item.Amount.IsZero())
				}
			},
		},
		{
			name: "overpayment scenario",
			debtList: &entities.DebtList{
				ID:                uuid.New(),
				TotalAmount:       decimal.RequireFromString("1000.00"),
				InstallmentAmount: decimal.RequireFromString("250.00"),
				InstallmentPlan:   "monthly",
				CreatedAt:         time.Now(),
				DueDate:           time.Now().AddDate(0, 4, 0),
			},
			payments: []entities.DebtItem{
				{
					ID:          uuid.New(),
					Amount:      decimal.RequireFromString("1200.00"),
					Status:      "completed",
					PaymentDate: time.Now(),
				},
			},
			expectedSchedule: 4,
			validateSchedule: func(t *testing.T, schedule []entities.PaymentScheduleItem, debtList *entities.DebtList) {
				assert.Len(t, schedule, 4)
				
				// All payments should be marked as paid (overpayment covers all)
				for _, item := range schedule {
					assert.Equal(t, "paid", item.Status)
					assert.True(t, item.Amount.IsZero())
				}
			},
		},
		{
			name: "weekly installment plan",
			debtList: &entities.DebtList{
				ID:                uuid.New(),
				TotalAmount:       decimal.RequireFromString("400.00"),
				InstallmentAmount: decimal.RequireFromString("100.00"),
				InstallmentPlan:   "weekly",
				CreatedAt:         time.Now(),
				DueDate:           time.Now().AddDate(0, 0, 28),
			},
			payments: []entities.DebtItem{
				{
					ID:          uuid.New(),
					Amount:      decimal.RequireFromString("100.00"),
					Status:      "completed",
					PaymentDate: time.Now(),
				},
			},
			expectedSchedule: 4,
			validateSchedule: func(t *testing.T, schedule []entities.PaymentScheduleItem, debtList *entities.DebtList) {
				assert.Len(t, schedule, 4)
				
				// First payment should be paid
				assert.Equal(t, "paid", schedule[0].Status)
				assert.True(t, schedule[0].PaidAmount.Equal(decimal.RequireFromString("100.00")))
				
				// Verify dates are weekly intervals
				for i := 1; i < len(schedule); i++ {
					diff := schedule[i].DueDate.Sub(schedule[i-1].DueDate)
					days := int(diff.Hours() / 24)
					assert.Equal(t, 7, days, "Weekly payments should be 7 days apart")
				}
			},
		},
		{
			name: "biweekly installment plan",
			debtList: &entities.DebtList{
				ID:                uuid.New(),
				TotalAmount:       decimal.RequireFromString("600.00"),
				InstallmentAmount: decimal.RequireFromString("200.00"),
				InstallmentPlan:   "biweekly",
				CreatedAt:         time.Now(),
				DueDate:           time.Now().AddDate(0, 0, 42),
			},
			payments:         []entities.DebtItem{},
			expectedSchedule: 3,
			validateSchedule: func(t *testing.T, schedule []entities.PaymentScheduleItem, debtList *entities.DebtList) {
				assert.Len(t, schedule, 3)
				
				// Verify dates are biweekly intervals
				for i := 1; i < len(schedule); i++ {
					diff := schedule[i].DueDate.Sub(schedule[i-1].DueDate)
					days := int(diff.Hours() / 24)
					assert.Equal(t, 14, days, "Biweekly payments should be 14 days apart")
				}
			},
		},
		{
			name: "irregular final payment (smaller)",
			debtList: &entities.DebtList{
				ID:                uuid.New(),
				TotalAmount:       decimal.RequireFromString("550.00"),
				InstallmentAmount: decimal.RequireFromString("200.00"),
				InstallmentPlan:   "monthly",
				CreatedAt:         time.Now(),
				DueDate:           time.Now().AddDate(0, 3, 0),
			},
			payments:         []entities.DebtItem{},
			expectedSchedule: 3,
			validateSchedule: func(t *testing.T, schedule []entities.PaymentScheduleItem, debtList *entities.DebtList) {
				assert.Len(t, schedule, 3)
				
				// First two payments should be 200
				assert.True(t, schedule[0].ScheduledAmount.Equal(decimal.RequireFromString("200.00")))
				assert.True(t, schedule[1].ScheduledAmount.Equal(decimal.RequireFromString("200.00")))
				
				// Final payment should be 150 (remaining balance)
				assert.True(t, schedule[2].ScheduledAmount.Equal(decimal.RequireFromString("150.00")))
			},
		},
		{
			name: "pending payments excluded",
			debtList: &entities.DebtList{
				ID:                uuid.New(),
				TotalAmount:       decimal.RequireFromString("1000.00"),
				InstallmentAmount: decimal.RequireFromString("250.00"),
				InstallmentPlan:   "monthly",
				CreatedAt:         time.Now(),
				DueDate:           time.Now().AddDate(0, 4, 0),
			},
			payments: []entities.DebtItem{
				{
					ID:          uuid.New(),
					Amount:      decimal.RequireFromString("250.00"),
					Status:      "completed",
					PaymentDate: time.Now(),
				},
				{
					ID:          uuid.New(),
					Amount:      decimal.RequireFromString("250.00"),
					Status:      "pending", // Should not be counted
					PaymentDate: time.Now(),
				},
			},
			expectedSchedule: 4,
			validateSchedule: func(t *testing.T, schedule []entities.PaymentScheduleItem, debtList *entities.DebtList) {
				assert.Len(t, schedule, 4)
				
				// Only first payment should be paid (pending payment not counted)
				assert.Equal(t, "paid", schedule[0].Status)
				assert.True(t, schedule[0].PaidAmount.Equal(decimal.RequireFromString("250.00")))
				
				// Rest should be pending
				for i := 1; i < 4; i++ {
					assert.Equal(t, "pending", schedule[i].Status)
					assert.True(t, schedule[i].PaidAmount.IsZero())
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create service
			service := services.NewPaymentScheduleService()

			// Execute
			schedule := service.CalculatePaymentSchedule(tt.debtList, tt.payments)

			// Assert
			require.NotNil(t, schedule)
			assert.Equal(t, tt.expectedSchedule, len(schedule))
			
			// Custom validation
			if tt.validateSchedule != nil {
				tt.validateSchedule(t, schedule, tt.debtList)
			}
		})
	}
}

func TestCalculateNextPaymentDate(t *testing.T) {
	baseTime := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	
	tests := []struct {
		name            string
		debtList        *entities.DebtList
		lastPaymentDate *time.Time
		expectedDay     int
		expectedMonth   time.Month
	}{
		{
			name: "weekly from creation date",
			debtList: &entities.DebtList{
				InstallmentPlan: "weekly",
				CreatedAt:       baseTime,
				DueDate:         baseTime.AddDate(0, 1, 0),
			},
			lastPaymentDate: nil,
			expectedDay:     22,
			expectedMonth:   time.January,
		},
		{
			name: "weekly from last payment",
			debtList: &entities.DebtList{
				InstallmentPlan: "weekly",
				CreatedAt:       baseTime,
				DueDate:         baseTime.AddDate(0, 1, 0),
			},
			lastPaymentDate: timePtr(baseTime.AddDate(0, 0, 7)),
			expectedDay:     29,
			expectedMonth:   time.January,
		},
		{
			name: "biweekly from creation date",
			debtList: &entities.DebtList{
				InstallmentPlan: "biweekly",
				CreatedAt:       baseTime,
				DueDate:         baseTime.AddDate(0, 1, 0),
			},
			lastPaymentDate: nil,
			expectedDay:     29,
			expectedMonth:   time.January,
		},
		{
			name: "monthly from creation date",
			debtList: &entities.DebtList{
				InstallmentPlan: "monthly",
				CreatedAt:       baseTime,
				DueDate:         baseTime.AddDate(0, 3, 0),
			},
			lastPaymentDate: nil,
			expectedDay:     15,
			expectedMonth:   time.February,
		},
		{
			name: "onetime uses due date",
			debtList: &entities.DebtList{
				InstallmentPlan: "onetime",
				CreatedAt:       baseTime,
				DueDate:         baseTime.AddDate(0, 2, 5),
			},
			lastPaymentDate: nil,
			expectedDay:     20,
			expectedMonth:   time.March,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := services.NewPaymentScheduleService()
			
			result := service.CalculateNextPaymentDate(tt.debtList, tt.lastPaymentDate)
			
			assert.Equal(t, tt.expectedDay, result.Day())
			assert.Equal(t, tt.expectedMonth, result.Month())
		})
	}
}
