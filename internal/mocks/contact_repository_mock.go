package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"exit-debt/internal/domain/entities"
)

// MockContactRepository is a mock implementation of ContactRepository
type MockContactRepository struct {
	mock.Mock
}

func (m *MockContactRepository) Create(ctx context.Context, contact *entities.Contact) error {
	args := m.Called(ctx, contact)
	return args.Error(0)
}

func (m *MockContactRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Contact, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Contact), args.Error(1)
}

func (m *MockContactRepository) GetByEmail(ctx context.Context, email string) (*entities.Contact, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Contact), args.Error(1)
}

func (m *MockContactRepository) GetByPhone(ctx context.Context, phone string) (*entities.Contact, error) {
	args := m.Called(ctx, phone)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Contact), args.Error(1)
}

func (m *MockContactRepository) GetUserContacts(ctx context.Context, userID uuid.UUID) ([]entities.Contact, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.Contact), args.Error(1)
}

func (m *MockContactRepository) Update(ctx context.Context, contact *entities.Contact) error {
	args := m.Called(ctx, contact)
	return args.Error(0)
}

func (m *MockContactRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockContactRepository) ExistsByEmailForUser(ctx context.Context, userID uuid.UUID, email string) (bool, error) {
	args := m.Called(ctx, userID, email)
	return args.Bool(0), args.Error(1)
}

func (m *MockContactRepository) ExistsByPhoneForUser(ctx context.Context, userID uuid.UUID, phone string) (bool, error) {
	args := m.Called(ctx, userID, phone)
	return args.Bool(0), args.Error(1)
}

func (m *MockContactRepository) GetContactsWithEmail(ctx context.Context, email string) ([]entities.Contact, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.Contact), args.Error(1)
}

func (m *MockContactRepository) CreateUserContactRelation(ctx context.Context, userContact *entities.UserContact) error {
	args := m.Called(ctx, userContact)
	return args.Error(0)
}

func (m *MockContactRepository) GetUserContactRelation(ctx context.Context, userID uuid.UUID, contactID uuid.UUID) (*entities.UserContact, error) {
	args := m.Called(ctx, userID, contactID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.UserContact), args.Error(1)
}

func (m *MockContactRepository) DeleteUserContactRelation(ctx context.Context, userID uuid.UUID, contactID uuid.UUID) error {
	args := m.Called(ctx, userID, contactID)
	return args.Error(0)
}

func (m *MockContactRepository) GetUserContactRelationsByContactID(ctx context.Context, contactID uuid.UUID) ([]entities.UserContact, error) {
	args := m.Called(ctx, contactID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.UserContact), args.Error(1)
}
