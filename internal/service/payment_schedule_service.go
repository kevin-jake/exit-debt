package service

import (
	"time"

	"exit-debt/internal/models"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type PaymentScheduleService struct{}

func NewPaymentScheduleService() *PaymentScheduleService {
	return &PaymentScheduleService{}
}

// CalculateNextPaymentDate calculates when the next payment is due
func (s *PaymentScheduleService) CalculateNextPaymentDate(debtList *models.DebtList, lastPaymentDate *time.Time) time.Time {
	var startDate time.Time
	
	if lastPaymentDate != nil {
		// Use the last payment date as reference
		startDate = *lastPaymentDate
	} else {
		// Use the debt creation date as reference
		startDate = debtList.CreatedAt
	}

	switch debtList.InstallmentPlan {
	case "weekly":
		return startDate.AddDate(0, 0, 7)
	case "biweekly":
		return startDate.AddDate(0, 0, 14)
	case "monthly":
		return startDate.AddDate(0, 1, 0)
	case "quarterly":
		return startDate.AddDate(0, 3, 0)
	case "yearly":
		return startDate.AddDate(1, 0, 0)
	default:
		return startDate.AddDate(0, 1, 0) // Default to monthly
	}
}

// CalculateRemainingPayments calculates how many payments are left
func (s *PaymentScheduleService) CalculateRemainingPayments(debtList *models.DebtList, totalPaid decimal.Decimal) int {
	remainingAmount := debtList.TotalAmount.Sub(totalPaid)
	if remainingAmount.LessThanOrEqual(decimal.Zero) {
		return 0
	}

	// Calculate how many full installments are needed
	installmentsNeeded := remainingAmount.Div(debtList.InstallmentAmount).Ceil()
	return int(installmentsNeeded.IntPart())
}

// CalculatePaymentSchedule generates a complete payment schedule with payment tracking
func (s *PaymentScheduleService) CalculatePaymentSchedule(debtList *models.DebtList, db *gorm.DB) []PaymentScheduleItem {
	var schedule []PaymentScheduleItem
	currentDate := debtList.CreatedAt
	paymentNumber := 1

	// Calculate total number of payments based on time period
	totalPayments := s.CalculateNumberOfPayments(debtList.InstallmentPlan, debtList.CreatedAt, debtList.DueDate)
	
	// Get all payments made for this debt list
	var payments []models.DebtItem
	if err := db.Where("debt_list_id = ? AND status = ?", debtList.ID, "completed").Order("payment_date ASC").Find(&payments).Error; err != nil {
		// If error, continue with empty payments
		payments = []models.DebtItem{}
	}

	// Calculate total payments made
	totalPaymentsMade := decimal.Zero
	for _, payment := range payments {
		totalPaymentsMade = totalPaymentsMade.Add(payment.Amount)
	}

	// Calculate remaining debt after payments
	remainingAfterPayments := debtList.TotalAmount.Sub(totalPaymentsMade)
	if remainingAfterPayments.LessThan(decimal.Zero) {
		remainingAfterPayments = decimal.Zero
	}

	// Use the original installment amount from the debt list
	originalInstallmentAmount := debtList.InstallmentAmount

	// Track how much has been paid against each installment
	remainingPaymentAmount := totalPaymentsMade

	for paymentNumber <= totalPayments && remainingAfterPayments.GreaterThan(decimal.Zero) {
		// Use the original installment amount
		paymentAmount := originalInstallmentAmount
		if remainingAfterPayments.LessThan(paymentAmount) {
			paymentAmount = remainingAfterPayments
		}

		// Determine if this installment is already paid
		status := "pending"
		amountNeeded := paymentAmount

		if remainingPaymentAmount.GreaterThanOrEqual(paymentAmount) {
			// This installment is fully paid
			status = "paid"
			amountNeeded = decimal.Zero
			remainingPaymentAmount = remainingPaymentAmount.Sub(paymentAmount)
		} else if remainingPaymentAmount.GreaterThan(decimal.Zero) {
			// This installment is partially paid
			amountNeeded = paymentAmount.Sub(remainingPaymentAmount)
			remainingPaymentAmount = decimal.Zero
		}

		// Calculate next payment date
		nextDate := s.CalculateNextPaymentDate(debtList, &currentDate)

		scheduleItem := PaymentScheduleItem{
			PaymentNumber:  paymentNumber,
			DueDate:        nextDate,
			Amount:         amountNeeded,
			Status:         status,
		}

		schedule = append(schedule, scheduleItem)

		// Update for next iteration
		remainingAfterPayments = remainingAfterPayments.Sub(paymentAmount)
		currentDate = nextDate
		paymentNumber++
	}

	return schedule
}

// PaymentScheduleItem represents a scheduled payment
type PaymentScheduleItem struct {
	PaymentNumber int             `json:"payment_number"`
	DueDate       time.Time       `json:"due_date"`
	Amount        decimal.Decimal `json:"amount"`
	Status        string          `json:"status"` // pending, paid, overdue, missed
}

// GetUpcomingPayments returns payments due in the next X days
func (s *PaymentScheduleService) GetUpcomingPayments(debtList *models.DebtList, days int, db *gorm.DB) []PaymentScheduleItem {
	schedule := s.CalculatePaymentSchedule(debtList, db)
	var upcoming []PaymentScheduleItem
	
	cutoffDate := time.Now().AddDate(0, 0, days)
	
	for _, item := range schedule {
		if item.DueDate.After(time.Now()) && item.DueDate.Before(cutoffDate) {
			upcoming = append(upcoming, item)
		}
	}
	
	return upcoming
}

// GetOverduePayments returns payments that are overdue
func (s *PaymentScheduleService) GetOverduePayments(debtList *models.DebtList, db *gorm.DB) []PaymentScheduleItem {
	schedule := s.CalculatePaymentSchedule(debtList, db)
	var overdue []PaymentScheduleItem
	
	for _, item := range schedule {
		if item.DueDate.Before(time.Now()) {
			item.Status = "overdue"
			overdue = append(overdue, item)
		}
	}
	
	return overdue
}

// CalculateDueDateFromNumberOfPayments calculates the due date based on number of payments and installment plan
func (s *PaymentScheduleService) CalculateDueDateFromNumberOfPayments(createdAt time.Time, numberOfPayments int, installmentPlan string) time.Time {
	if numberOfPayments <= 0 {
		numberOfPayments = 1
	}

	var dueDate time.Time
	switch installmentPlan {
	case "weekly":
		// Each payment is 7 days apart
		dueDate = createdAt.AddDate(0, 0, (numberOfPayments-1)*7)
	case "biweekly":
		// Each payment is 14 days apart
		dueDate = createdAt.AddDate(0, 0, (numberOfPayments-1)*14)
	case "monthly":
		// Each payment is 1 month apart
		dueDate = createdAt.AddDate(0, numberOfPayments-1, 0)
	case "quarterly":
		// Each payment is 3 months apart
		dueDate = createdAt.AddDate(0, (numberOfPayments-1)*3, 0)
	case "yearly":
		// Each payment is 1 year apart
		dueDate = createdAt.AddDate(numberOfPayments-1, 0, 0)
	default:
		// Default to monthly
		dueDate = createdAt.AddDate(0, numberOfPayments-1, 0)
	}

	return dueDate
}

// CalculateInstallmentAmountFromNumberOfPayments calculates the installment amount based on total amount and number of payments
func (s *PaymentScheduleService) CalculateInstallmentAmountFromNumberOfPayments(totalAmount decimal.Decimal, numberOfPayments int) decimal.Decimal {
	if numberOfPayments <= 0 {
		numberOfPayments = 1
	}
	return totalAmount.Div(decimal.NewFromInt(int64(numberOfPayments)))
}

// CalculateInstallmentAmount calculates the installment amount based on total amount, installment plan, and time period
func (s *PaymentScheduleService) CalculateInstallmentAmount(totalAmount decimal.Decimal, installmentPlan string, createdAt time.Time, dueDate time.Time) decimal.Decimal {
	numberOfPayments := s.CalculateNumberOfPayments(installmentPlan, createdAt, dueDate)
	if numberOfPayments <= 0 {
		numberOfPayments = 1 // At least 1 payment
	}
	return totalAmount.Div(decimal.NewFromInt(int64(numberOfPayments)))
}

// CalculateNumberOfPayments returns the number of payments based on installment plan and time period
func (s *PaymentScheduleService) CalculateNumberOfPayments(installmentPlan string, createdAt time.Time, dueDate time.Time) int {
	duration := dueDate.Sub(createdAt)
	days := int(duration.Hours() / 24)
	
	switch installmentPlan {
	case "weekly":
		// Calculate weeks between creation and due date
		weeks := days / 7
		if weeks < 1 {
			weeks = 1
		}
		return weeks
	case "biweekly":
		// Calculate biweekly periods (every 2 weeks)
		biweeklyPeriods := days / 14
		if biweeklyPeriods < 1 {
			biweeklyPeriods = 1
		}
		return biweeklyPeriods
	case "monthly":
		// Calculate months between creation and due date
		months := s.calculateMonthsBetween(createdAt, dueDate)
		if months < 1 {
			months = 1
		}
		return months
	case "quarterly":
		// Calculate quarters between creation and due date
		quarters := s.calculateQuartersBetween(createdAt, dueDate)
		if quarters < 1 {
			quarters = 1
		}
		return quarters
	case "yearly":
		// Calculate years between creation and due date
		years := s.calculateYearsBetween(createdAt, dueDate)
		if years < 1 {
			years = 1
		}
		return years
	default:
		// Default to monthly
		months := s.calculateMonthsBetween(createdAt, dueDate)
		if months < 1 {
			months = 1
		}
		return months
	}
}

// Helper functions for calculating time periods
func (s *PaymentScheduleService) calculateMonthsBetween(start, end time.Time) int {
	years := end.Year() - start.Year()
	months := end.Month() - start.Month()
	
	if end.Day() < start.Day() {
		months--
	}
	
	totalMonths := years*12 + int(months)
	if totalMonths < 0 {
		totalMonths = 0
	}
	
	return totalMonths + 1 // Add 1 to include both start and end months
}

func (s *PaymentScheduleService) calculateQuartersBetween(start, end time.Time) int {
	months := s.calculateMonthsBetween(start, end)
	quarters := months / 3
	if months%3 > 0 {
		quarters++
	}
	return quarters
}

func (s *PaymentScheduleService) calculateYearsBetween(start, end time.Time) int {
	years := end.Year() - start.Year()
	if end.Month() < start.Month() || (end.Month() == start.Month() && end.Day() < start.Day()) {
		years--
	}
	if years < 0 {
		years = 0
	}
	return years + 1 // Add 1 to include both start and end years
} 