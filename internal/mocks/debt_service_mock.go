package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"exit-debt/internal/domain/entities"
)

// MockDebtService is a mock implementation of DebtService
type MockDebtService struct {
	mock.Mock
}

func (m *MockDebtService) CreateDebtList(ctx context.Context, userID uuid.UUID, req *entities.CreateDebtListRequest) (*entities.DebtList, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.DebtList), args.Error(1)
}

func (m *MockDebtService) GetDebtList(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*entities.DebtListResponse, error) {
	args := m.Called(ctx, id, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.DebtListResponse), args.Error(1)
}

func (m *MockDebtService) GetUserDebtLists(ctx context.Context, userID uuid.UUID) ([]entities.DebtListResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.DebtListResponse), args.Error(1)
}

func (m *MockDebtService) UpdateDebtList(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *entities.UpdateDebtListRequest) (*entities.DebtList, error) {
	args := m.Called(ctx, id, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.DebtList), args.Error(1)
}

func (m *MockDebtService) DeleteDebtList(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

func (m *MockDebtService) CreateDebtItem(ctx context.Context, userID uuid.UUID, req *entities.CreateDebtItemRequest) (*entities.DebtItem, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.DebtItem), args.Error(1)
}

func (m *MockDebtService) GetDebtItem(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*entities.DebtItem, error) {
	args := m.Called(ctx, id, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.DebtItem), args.Error(1)
}

func (m *MockDebtService) GetDebtListItems(ctx context.Context, debtListID uuid.UUID, userID uuid.UUID) ([]entities.DebtItem, error) {
	args := m.Called(ctx, debtListID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.DebtItem), args.Error(1)
}

func (m *MockDebtService) UpdateDebtItem(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *entities.UpdateDebtItemRequest) (*entities.DebtItem, error) {
	args := m.Called(ctx, id, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.DebtItem), args.Error(1)
}

func (m *MockDebtService) DeleteDebtItem(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

func (m *MockDebtService) GetOverdueItems(ctx context.Context, userID uuid.UUID) ([]entities.DebtList, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.DebtList), args.Error(1)
}

func (m *MockDebtService) GetDueSoonItems(ctx context.Context, userID uuid.UUID, days int) ([]entities.DebtList, error) {
	args := m.Called(ctx, userID, days)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.DebtList), args.Error(1)
}

func (m *MockDebtService) GetPaymentSchedule(ctx context.Context, debtListID uuid.UUID, userID uuid.UUID) ([]entities.PaymentScheduleItem, error) {
	args := m.Called(ctx, debtListID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.PaymentScheduleItem), args.Error(1)
}

func (m *MockDebtService) GetUpcomingPayments(ctx context.Context, userID uuid.UUID, days int) ([]entities.UpcomingPayment, error) {
	args := m.Called(ctx, userID, days)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.UpcomingPayment), args.Error(1)
}

func (m *MockDebtService) GetTotalPaymentsForDebtList(ctx context.Context, debtListID uuid.UUID, userID uuid.UUID) (*entities.PaymentSummary, error) {
	args := m.Called(ctx, debtListID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.PaymentSummary), args.Error(1)
}
