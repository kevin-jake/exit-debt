package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"exit-debt/internal/domain/entities"
)

// MockContactService is a mock implementation of ContactService
type MockContactService struct {
	mock.Mock
}

func (m *MockContactService) CreateContact(ctx context.Context, userID uuid.UUID, req *entities.CreateContactRequest) (*entities.Contact, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Contact), args.Error(1)
}

func (m *MockContactService) GetContact(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*entities.Contact, error) {
	args := m.Called(ctx, id, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Contact), args.Error(1)
}

func (m *MockContactService) GetUserContacts(ctx context.Context, userID uuid.UUID) ([]entities.Contact, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.Contact), args.Error(1)
}

func (m *MockContactService) UpdateContact(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *entities.UpdateContactRequest) (*entities.Contact, error) {
	args := m.Called(ctx, id, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Contact), args.Error(1)
}

func (m *MockContactService) DeleteContact(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

func (m *MockContactService) CreateContactsForNewUser(ctx context.Context, userID uuid.UUID, userEmail string) error {
	args := m.Called(ctx, userID, userEmail)
	return args.Error(0)
}

func (m *MockContactService) CreateReciprocalContact(ctx context.Context, contactEmail string, contactOwnerID uuid.UUID) error {
	args := m.Called(ctx, contactEmail, contactOwnerID)
	return args.Error(0)
}
