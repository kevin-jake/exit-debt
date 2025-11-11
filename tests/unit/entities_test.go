package unit

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"pay-your-dues/internal/domain/entities"
)

func TestUser_IsValid(t *testing.T) {
	tests := []struct {
		name          string
		user          *entities.User
		expectedError error
		expectValid   bool
	}{
		{
			name: "valid user",
			user: &entities.User{
				ID:           uuid.New(),
				Email:        "test@example.com",
				PasswordHash: "hashed_password",
				FirstName:    "John",
				LastName:     "Doe",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			},
			expectedError: nil,
			expectValid:   true,
		},
		{
			name: "missing email",
			user: &entities.User{
				ID:           uuid.New(),
				Email:        "",
				PasswordHash: "hashed_password",
				FirstName:    "John",
				LastName:     "Doe",
			},
			expectedError: entities.ErrInvalidEmail,
			expectValid:   false,
		},
		{
			name: "missing first name",
			user: &entities.User{
				ID:           uuid.New(),
				Email:        "test@example.com",
				PasswordHash: "hashed_password",
				FirstName:    "",
				LastName:     "Doe",
			},
			expectedError: entities.ErrInvalidFirstName,
			expectValid:   false,
		},
		{
			name: "missing last name",
			user: &entities.User{
				ID:           uuid.New(),
				Email:        "test@example.com",
				PasswordHash: "hashed_password",
				FirstName:    "John",
				LastName:     "",
			},
			expectedError: entities.ErrInvalidLastName,
			expectValid:   false,
		},
		{
			name: "missing password hash",
			user: &entities.User{
				ID:           uuid.New(),
				Email:        "test@example.com",
				PasswordHash: "",
				FirstName:    "John",
				LastName:     "Doe",
			},
			expectedError: entities.ErrInvalidPassword,
			expectValid:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.IsValid()

			if tt.expectValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedError)
			}
		})
	}
}

func TestUser_FullName(t *testing.T) {
	user := &entities.User{
		FirstName: "John",
		LastName:  "Doe",
	}

	fullName := user.FullName()
	assert.Equal(t, "John Doe", fullName)
}

func TestUserContact_IsValid(t *testing.T) {
	tests := []struct {
		name          string
		userContact   *entities.UserContact
		expectedError error
		expectValid   bool
	}{
		{
			name: "valid user contact",
			userContact: &entities.UserContact{
				ID:        uuid.New(),
				UserID:    uuid.New(),
				ContactID: uuid.New(),
				Name:      "Test Contact",
				Email:     stringPtr("test@example.com"),
				Phone:     stringPtr("+1234567890"),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			expectedError: nil,
			expectValid:   true,
		},
		{
			name: "missing name",
			userContact: &entities.UserContact{
				ID:        uuid.New(),
				UserID:    uuid.New(),
				ContactID: uuid.New(),
				Name:      "",
				Email:     stringPtr("test@example.com"),
			},
			expectedError: entities.ErrInvalidContactName,
			expectValid:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.userContact.IsValid()

			if tt.expectValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedError)
			}
		})
	}
}

func TestUserContact_HasEmail(t *testing.T) {
	tests := []struct {
		name        string
		userContact *entities.UserContact
		expected    bool
	}{
		{
			name: "has email",
			userContact: &entities.UserContact{
				Email: stringPtr("test@example.com"),
			},
			expected: true,
		},
		{
			name: "nil email",
			userContact: &entities.UserContact{
				Email: nil,
			},
			expected: false,
		},
		{
			name: "empty email",
			userContact: &entities.UserContact{
				Email: stringPtr(""),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.userContact.HasEmail()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUserContact_HasPhone(t *testing.T) {
	tests := []struct {
		name        string
		userContact *entities.UserContact
		expected    bool
	}{
		{
			name: "has phone",
			userContact: &entities.UserContact{
				Phone: stringPtr("+1234567890"),
			},
			expected: true,
		},
		{
			name: "nil phone",
			userContact: &entities.UserContact{
				Phone: nil,
			},
			expected: false,
		},
		{
			name: "empty phone",
			userContact: &entities.UserContact{
				Phone: stringPtr(""),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.userContact.HasPhone()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDebtList_IsValid(t *testing.T) {
	tests := []struct {
		name          string
		debtList      *entities.DebtList
		expectedError error
		expectValid   bool
	}{
		{
			name: "valid debt list",
			debtList: &entities.DebtList{
				ID:            uuid.New(),
				UserID:        uuid.New(),
				ContactID:     uuid.New(),
				DebtType:      "to_pay",
				TotalAmount:   decimal.RequireFromString("1000.00"),
				Currency:      "USD",
				Status:        "active",
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			},
			expectedError: nil,
			expectValid:   true,
		},
		{
			name: "nil user ID",
			debtList: &entities.DebtList{
				UserID:      uuid.Nil,
				ContactID:   uuid.New(),
				DebtType:    "to_pay",
				TotalAmount: decimal.RequireFromString("1000.00"),
				Currency:    "USD",
			},
			expectedError: entities.ErrInvalidInput,
			expectValid:   false,
		},
		{
			name: "nil contact ID",
			debtList: &entities.DebtList{
				UserID:      uuid.New(),
				ContactID:   uuid.Nil,
				DebtType:    "to_pay",
				TotalAmount: decimal.RequireFromString("1000.00"),
				Currency:    "USD",
			},
			expectedError: entities.ErrInvalidInput,
			expectValid:   false,
		},
		{
			name: "invalid debt type",
			debtList: &entities.DebtList{
				UserID:      uuid.New(),
				ContactID:   uuid.New(),
				DebtType:    "invalid_type",
				TotalAmount: decimal.RequireFromString("1000.00"),
				Currency:    "USD",
			},
			expectedError: entities.ErrInvalidDebtType,
			expectValid:   false,
		},
		{
			name: "negative amount",
			debtList: &entities.DebtList{
				UserID:      uuid.New(),
				ContactID:   uuid.New(),
				DebtType:    "to_pay",
				TotalAmount: decimal.RequireFromString("-100.00"),
				Currency:    "USD",
			},
			expectedError: entities.ErrInvalidAmount,
			expectValid:   false,
		},
		{
			name: "empty currency",
			debtList: &entities.DebtList{
				UserID:      uuid.New(),
				ContactID:   uuid.New(),
				DebtType:    "to_pay",
				TotalAmount: decimal.RequireFromString("1000.00"),
				Currency:    "",
			},
			expectedError: entities.ErrInvalidCurrency,
			expectValid:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.debtList.IsValid()

			if tt.expectValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedError)
			}
		})
	}
}

func TestDebtList_IsSettled(t *testing.T) {
	tests := []struct {
		name     string
		debtList *entities.DebtList
		expected bool
	}{
		{
			name: "settled by status",
			debtList: &entities.DebtList{
				Status:              "settled",
				TotalRemainingDebt:  decimal.RequireFromString("100.00"),
			},
			expected: true,
		},
		{
			name: "settled by zero remaining debt",
			debtList: &entities.DebtList{
				Status:              "active",
				TotalRemainingDebt:  decimal.Zero,
			},
			expected: true,
		},
		{
			name: "not settled",
			debtList: &entities.DebtList{
				Status:              "active",
				TotalRemainingDebt:  decimal.RequireFromString("500.00"),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.debtList.IsSettled()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDebtList_IsOverdue(t *testing.T) {
	now := time.Now()
	pastDate := now.AddDate(0, 0, -5)
	futureDate := now.AddDate(0, 0, 5)

	tests := []struct {
		name     string
		debtList *entities.DebtList
		expected bool
	}{
		{
			name: "overdue and not settled",
			debtList: &entities.DebtList{
				NextPaymentDate:    pastDate,
				Status:             "active",
				TotalRemainingDebt: decimal.RequireFromString("500.00"),
			},
			expected: true,
		},
		{
			name: "overdue but settled",
			debtList: &entities.DebtList{
				NextPaymentDate:    pastDate,
				Status:             "settled",
				TotalRemainingDebt: decimal.Zero,
			},
			expected: false,
		},
		{
			name: "not overdue",
			debtList: &entities.DebtList{
				NextPaymentDate:    futureDate,
				Status:             "active",
				TotalRemainingDebt: decimal.RequireFromString("500.00"),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.debtList.IsOverdue()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDebtList_CalculateProgress(t *testing.T) {
	tests := []struct {
		name     string
		debtList *entities.DebtList
		expected string
	}{
		{
			name: "50% progress",
			debtList: &entities.DebtList{
				TotalAmount:       decimal.RequireFromString("1000.00"),
				TotalPaymentsMade: decimal.RequireFromString("500.00"),
			},
			expected: "50.00",
		},
		{
			name: "100% progress",
			debtList: &entities.DebtList{
				TotalAmount:       decimal.RequireFromString("1000.00"),
				TotalPaymentsMade: decimal.RequireFromString("1000.00"),
			},
			expected: "100.00",
		},
		{
			name: "zero total amount",
			debtList: &entities.DebtList{
				TotalAmount:       decimal.Zero,
				TotalPaymentsMade: decimal.RequireFromString("100.00"),
			},
			expected: "0.00",
		},
		{
			name: "partial progress",
			debtList: &entities.DebtList{
				TotalAmount:       decimal.RequireFromString("300.00"),
				TotalPaymentsMade: decimal.RequireFromString("100.00"),
			},
			expected: "33.33",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.debtList.CalculateProgress()
			assert.Equal(t, tt.expected, result.StringFixed(2))
		})
	}
}

func TestDebtItem_IsValid(t *testing.T) {
	tests := []struct {
		name          string
		debtItem      *entities.DebtItem
		expectedError error
		expectValid   bool
	}{
		{
			name: "valid debt item",
			debtItem: &entities.DebtItem{
				ID:            uuid.New(),
				DebtListID:    uuid.New(),
				Amount:        decimal.RequireFromString("200.00"),
				Currency:      "USD",
				PaymentMethod: "bank_transfer",
				PaymentDate:   time.Now(),
				Status:        "completed",
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			},
			expectedError: nil,
			expectValid:   true,
		},
		{
			name: "nil debt list ID",
			debtItem: &entities.DebtItem{
				DebtListID:    uuid.Nil,
				Amount:        decimal.RequireFromString("200.00"),
				Currency:      "USD",
				PaymentMethod: "bank_transfer",
			},
			expectedError: entities.ErrInvalidInput,
			expectValid:   false,
		},
		{
			name: "negative amount",
			debtItem: &entities.DebtItem{
				DebtListID:    uuid.New(),
				Amount:        decimal.RequireFromString("-50.00"),
				Currency:      "USD",
				PaymentMethod: "bank_transfer",
			},
			expectedError: entities.ErrInvalidAmount,
			expectValid:   false,
		},
		{
			name: "empty currency",
			debtItem: &entities.DebtItem{
				DebtListID:    uuid.New(),
				Amount:        decimal.RequireFromString("200.00"),
				Currency:      "",
				PaymentMethod: "bank_transfer",
			},
			expectedError: entities.ErrInvalidCurrency,
			expectValid:   false,
		},
		{
			name: "empty payment method",
			debtItem: &entities.DebtItem{
				DebtListID:    uuid.New(),
				Amount:        decimal.RequireFromString("200.00"),
				Currency:      "USD",
				PaymentMethod: "",
			},
			expectedError: entities.ErrInvalidPaymentMethod,
			expectValid:   false,
		},
		{
			name: "valid debt item with verification fields",
			debtItem: &entities.DebtItem{
				ID:                uuid.New(),
				DebtListID:        uuid.New(),
				Amount:            decimal.RequireFromString("200.00"),
				Currency:          "USD",
				PaymentMethod:     "bank_transfer",
				PaymentDate:       time.Now(),
				Status:            entities.PaymentStatusCompleted,
				ReceiptPhotoURL:   func() *string { s := "https://example.com/receipt.jpg"; return &s }(),
				VerifiedBy:        func() *uuid.UUID { id := uuid.New(); return &id }(),
				VerifiedAt:        func() *time.Time { t := time.Now(); return &t }(),
				VerificationNotes: func() *string { s := "Payment verified"; return &s }(),
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			},
			expectedError: nil,
			expectValid:   true,
		},
		{
			name: "invalid payment status",
			debtItem: &entities.DebtItem{
				DebtListID: uuid.New(),
				Amount:     decimal.RequireFromString("200.00"),
				Currency:   "USD",
				PaymentMethod: "bank_transfer",
				Status:     "invalid_status",
			},
			expectedError: entities.ErrInvalidPaymentStatus,
			expectValid:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.debtItem.IsValid()

			if tt.expectValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedError)
			}
		})
	}
}




