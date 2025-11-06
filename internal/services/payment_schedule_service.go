package services

import (
	"time"

	"github.com/shopspring/decimal"

	"exit-debt/internal/domain/entities"
	"exit-debt/internal/domain/interfaces"
)

// paymentScheduleService implements the PaymentScheduleService interface
type paymentScheduleService struct{}

// NewPaymentScheduleService creates a new payment schedule service
func NewPaymentScheduleService() interfaces.PaymentScheduleService {
	return &paymentScheduleService{}
}

func (s *paymentScheduleService) CalculateNextPaymentDate(debtList *entities.DebtList, lastPaymentDate *time.Time) time.Time {
	var startDate time.Time
	
	if lastPaymentDate != nil {
		// Use the last payment date as reference
		startDate = *lastPaymentDate
	} else {
		// Use the debt creation date as reference
		startDate = debtList.CreatedAt
	}

	switch debtList.InstallmentPlan {
	case "onetime":
		// For 1-time payments, return the due date itself
		return debtList.DueDate
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
		return debtList.DueDate // Default to onetime
	}
}
func (s *paymentScheduleService) CalculatePaymentSchedule(debtList *entities.DebtList, payments []entities.DebtItem) []entities.PaymentScheduleItem {
	var schedule []entities.PaymentScheduleItem
	currentDate := debtList.CreatedAt
	paymentNumber := 1

	// Calculate total payments made
	totalPaymentsMade := decimal.Zero
	for _, payment := range payments {
		if payment.Status == "completed" {
			totalPaymentsMade = totalPaymentsMade.Add(payment.Amount)
		}
	}

	// Use the original installment amount from the debt list
	originalInstallmentAmount := debtList.InstallmentAmount

	// Variable 1: Track how much of the USER'S PAYMENTS remains to be allocated to installments
	// This is used to determine which installments should be marked as "paid" or "partially paid"
	remainingPaymentAmount := totalPaymentsMade

	// Variable 2: Track how much of the TOTAL DEBT remains to be scheduled as installments
	// This starts at the full debt amount and decrements by each installment amount scheduled
	remainingDebtToSchedule := debtList.TotalAmount

	// Continue generating payment schedule items until all debt is accounted for
	for remainingDebtToSchedule.GreaterThan(decimal.Zero) {
		// Use the original installment amount, but cap it at remaining debt
		paymentAmount := originalInstallmentAmount
		if remainingDebtToSchedule.LessThan(paymentAmount) {
			paymentAmount = remainingDebtToSchedule
		}

		// Determine if this installment is already paid
		status := "pending"
		amountNeeded := paymentAmount
		paidAmount := decimal.Zero

		if remainingPaymentAmount.GreaterThanOrEqual(paymentAmount) {
			status = "paid"
			paidAmount = paymentAmount
			amountNeeded = decimal.Zero
			remainingPaymentAmount = remainingPaymentAmount.Sub(paymentAmount)
		} else if remainingPaymentAmount.GreaterThan(decimal.Zero) {
			paidAmount = remainingPaymentAmount
			amountNeeded = paymentAmount.Sub(remainingPaymentAmount)
			remainingPaymentAmount = decimal.Zero
		}

		// Calculate next payment date
		nextDate := s.CalculateNextPaymentDate(debtList, &currentDate)

		scheduleItem := entities.PaymentScheduleItem{
			PaymentNumber:   paymentNumber,
			DueDate:         nextDate,
			Amount:          amountNeeded,
			ScheduledAmount: paymentAmount,
			PaidAmount:      paidAmount,
			Status:          status,
		}

		schedule = append(schedule, scheduleItem)

		// Update for next iteration
		remainingDebtToSchedule = remainingDebtToSchedule.Sub(paymentAmount)
		currentDate = nextDate
		paymentNumber++
	}

	return schedule
}

func (s *paymentScheduleService) CalculateDueDateFromNumberOfPayments(createdAt time.Time, numberOfPayments int, installmentPlan string) time.Time {
	if numberOfPayments <= 0 {
		numberOfPayments = 1
	}

	var dueDate time.Time
	switch installmentPlan {
	case "onetime":
		// For 1-time payments, due date is the creation date (single payment)
		dueDate = createdAt
	case "weekly":
		// Each payment is 7 days apart, final payment is N weeks from creation
		dueDate = createdAt.AddDate(0, 0, numberOfPayments*7)
	case "biweekly":
		// Each payment is 14 days apart, final payment is N biweekly periods from creation
		dueDate = createdAt.AddDate(0, 0, numberOfPayments*14)
	case "monthly":
		// Each payment is 1 month apart, final payment is N months from creation
		dueDate = createdAt.AddDate(0, numberOfPayments, 0)
	case "quarterly":
		// Each payment is 3 months apart, final payment is N quarters from creation
		dueDate = createdAt.AddDate(0, numberOfPayments*3, 0)
	case "yearly":
		// Each payment is 1 year apart, final payment is N years from creation
		dueDate = createdAt.AddDate(numberOfPayments, 0, 0)
	default:
		// Default to onetime (single payment)
		dueDate = createdAt
	}

	return dueDate
}

func (s *paymentScheduleService) CalculateInstallmentAmountFromNumberOfPayments(totalAmount decimal.Decimal, numberOfPayments int) decimal.Decimal {
	if numberOfPayments <= 0 {
		numberOfPayments = 1
	}
	return totalAmount.Div(decimal.NewFromInt(int64(numberOfPayments)))
}

func (s *paymentScheduleService) CalculateInstallmentAmount(totalAmount decimal.Decimal, installmentPlan string, createdAt time.Time, dueDate time.Time) decimal.Decimal {
	numberOfPayments := s.CalculateNumberOfPayments(installmentPlan, createdAt, dueDate)
	if numberOfPayments <= 0 {
		numberOfPayments = 1 // At least 1 payment
	}
	return totalAmount.Div(decimal.NewFromInt(int64(numberOfPayments)))
}

// Helper methods

func (s *paymentScheduleService) CalculateNumberOfPayments(installmentPlan string, createdAt time.Time, dueDate time.Time) int {
	duration := dueDate.Sub(createdAt)
	days := int(duration.Hours() / 24)
	
	switch installmentPlan {
	case "onetime":
		// For 1-time payments, always return 1
		return 1
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

func (s *paymentScheduleService) calculateMonthsBetween(start, end time.Time) int {
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

func (s *paymentScheduleService) calculateQuartersBetween(start, end time.Time) int {
	months := s.calculateMonthsBetween(start, end)
	quarters := months / 3
	if months%3 > 0 {
		quarters++
	}
	return quarters
}

func (s *paymentScheduleService) calculateYearsBetween(start, end time.Time) int {
	years := end.Year() - start.Year()
	if end.Month() < start.Month() || (end.Month() == start.Month() && end.Day() < start.Day()) {
		years--
	}
	if years < 0 {
		years = 0
	}
	return years + 1 // Add 1 to include both start and end years
}






