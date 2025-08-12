package unit

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"exit-debt/internal/domain/entities"
	"exit-debt/internal/mocks"
	"exit-debt/internal/services"
)

func TestContactService_CreateContact(t *testing.T) {
	userID := uuid.New()
	contactID := uuid.New()
	existingUserID := uuid.New()

	tests := []struct {
		name           string
		userID         uuid.UUID
		request        *entities.CreateContactRequest
		setupMocks     func(*mocks.MockContactRepository, *mocks.MockUserRepository)
		expectedError  error
		expectSuccess  bool
		validateResult func(*testing.T, *entities.Contact)
	}{
		{
			name:   "create regular contact (non-user)",
			userID: userID,
			request: &entities.CreateContactRequest{
				Name:  "Alice Johnson",
				Email: stringPtr("alice@example.com"),
				Phone: stringPtr("+1987654321"),
				Notes: stringPtr("Friend from college"),
			},
			setupMocks: func(contactRepo *mocks.MockContactRepository, userRepo *mocks.MockUserRepository) {
				contactRepo.On("ExistsByEmailForUser", mock.Anything, userID, "alice@example.com").Return(false, nil)
				userRepo.On("GetByEmail", mock.Anything, "alice@example.com").Return(nil, entities.ErrUserNotFound)
				contactRepo.On("GetByEmail", mock.Anything, "alice@example.com").Return(nil, entities.ErrContactNotFound)
				contactRepo.On("Create", mock.Anything, mock.AnythingOfType("*entities.Contact")).Return(nil)
				contactRepo.On("CreateUserContactRelation", mock.Anything, mock.AnythingOfType("*entities.UserContact")).Return(nil)
			},
			expectedError: nil,
			expectSuccess: true,
			validateResult: func(t *testing.T, contact *entities.Contact) {
				assert.Equal(t, "Alice Johnson", contact.Name)
				assert.Equal(t, "alice@example.com", *contact.Email)
				assert.False(t, contact.IsUser)
				assert.Nil(t, contact.UserIDRef)
			},
		},
		{
			name:   "create contact for existing user",
			userID: userID,
			request: &entities.CreateContactRequest{
				Name:  "Bob Smith",
				Email: stringPtr("bob@example.com"),
			},
			setupMocks: func(contactRepo *mocks.MockContactRepository, userRepo *mocks.MockUserRepository) {
				// Mock the initial check for existing contact with same email
				contactRepo.On("ExistsByEmailForUser", mock.Anything, userID, "bob@example.com").Return(false, nil)
				
				existingUser := &entities.User{
					ID:        existingUserID,
					Email:     "bob@example.com",
					FirstName: "Bob",
					LastName:  "Smith",
				}
				userRepo.On("GetByEmail", mock.Anything, "bob@example.com").Return(existingUser, nil)
				contactRepo.On("GetByEmail", mock.Anything, "bob@example.com").Return(nil, entities.ErrContactNotFound)
				contactRepo.On("Create", mock.Anything, mock.AnythingOfType("*entities.Contact")).Return(nil)
				contactRepo.On("CreateUserContactRelation", mock.Anything, mock.AnythingOfType("*entities.UserContact")).Return(nil)
				
				// Mock calls for CreateReciprocalContact
				contactOwner := &entities.User{
					ID:        userID,
					Email:     "user@example.com",
					FirstName: "Test",
					LastName:  "User",
				}
				userRepo.On("GetByID", mock.Anything, userID).Return(contactOwner, nil)
				contactRepo.On("ExistsByEmailForUser", mock.Anything, userID, "user@example.com").Return(false, nil)
				contactRepo.On("GetByEmail", mock.Anything, "user@example.com").Return(nil, entities.ErrContactNotFound)
				contactRepo.On("Create", mock.Anything, mock.AnythingOfType("*entities.Contact")).Return(nil)
				contactRepo.On("CreateUserContactRelation", mock.Anything, mock.AnythingOfType("*entities.UserContact")).Return(nil)
			},
			expectedError: nil,
			expectSuccess: true,
			validateResult: func(t *testing.T, contact *entities.Contact) {
				assert.Equal(t, "Bob Smith", contact.Name)
				assert.True(t, contact.IsUser)
				assert.Equal(t, existingUserID, *contact.UserIDRef)
			},
		},
		{
			name:   "contact already exists for user",
			userID: userID,
			request: &entities.CreateContactRequest{
				Name:  "Existing Contact",
				Email: stringPtr("existing@example.com"),
			},
			setupMocks: func(contactRepo *mocks.MockContactRepository, userRepo *mocks.MockUserRepository) {
				contactRepo.On("ExistsByEmailForUser", mock.Anything, userID, "existing@example.com").Return(true, nil)
			},
			expectedError: entities.ErrContactAlreadyExists,
			expectSuccess: false,
		},
		{
			name:   "missing contact name",
			userID: userID,
			request: &entities.CreateContactRequest{
				Name:  "",
				Email: stringPtr("test@example.com"),
			},
			setupMocks:    func(contactRepo *mocks.MockContactRepository, userRepo *mocks.MockUserRepository) {},
			expectedError: entities.ErrInvalidContactName,
			expectSuccess: false,
		},
		{
			name:   "reuse existing global contact",
			userID: userID,
			request: &entities.CreateContactRequest{
				Name:  "Shared Contact",
				Email: stringPtr("shared@example.com"),
			},
			setupMocks: func(contactRepo *mocks.MockContactRepository, userRepo *mocks.MockUserRepository) {
				contactRepo.On("ExistsByEmailForUser", mock.Anything, userID, "shared@example.com").Return(false, nil)
				userRepo.On("GetByEmail", mock.Anything, "shared@example.com").Return(nil, entities.ErrUserNotFound)
				existingContact := &entities.Contact{
					ID:        contactID,
					Name:      "Shared Contact",
					Email:     stringPtr("shared@example.com"),
					IsUser:    false,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				contactRepo.On("GetByEmail", mock.Anything, "shared@example.com").Return(existingContact, nil)
				contactRepo.On("CreateUserContactRelation", mock.Anything, mock.AnythingOfType("*entities.UserContact")).Return(nil)
			},
			expectedError: nil,
			expectSuccess: true,
			validateResult: func(t *testing.T, contact *entities.Contact) {
				assert.Equal(t, contactID, contact.ID)
				assert.Equal(t, "Shared Contact", contact.Name)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			contactRepo := &mocks.MockContactRepository{}
			userRepo := &mocks.MockUserRepository{}
			tt.setupMocks(contactRepo, userRepo)

			// Create service
			contactService := services.NewContactService(contactRepo, userRepo)

			// Execute
			ctx := context.Background()
			result, err := contactService.CreateContact(ctx, tt.userID, tt.request)

			// Assert
			if tt.expectSuccess {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if tt.validateResult != nil {
					tt.validateResult(t, result)
				}
			} else {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedError)
				assert.Nil(t, result)
			}

			// Verify mock expectations
			contactRepo.AssertExpectations(t)
			userRepo.AssertExpectations(t)
		})
	}
}

func TestContactService_GetContact(t *testing.T) {
	userID := uuid.New()
	contactID := uuid.New()
	otherUserID := uuid.New()

	tests := []struct {
		name          string
		contactID     uuid.UUID
		userID        uuid.UUID
		setupMocks    func(*mocks.MockContactRepository, *mocks.MockUserRepository)
		expectedError error
		expectSuccess bool
	}{
		{
			name:      "successful contact retrieval",
			contactID: contactID,
			userID:    userID,
			setupMocks: func(contactRepo *mocks.MockContactRepository, userRepo *mocks.MockUserRepository) {
				userContact := &entities.UserContact{
					ID:        uuid.New(),
					UserID:    userID,
					ContactID: contactID,
				}
				contactRepo.On("GetUserContactRelation", mock.Anything, userID, contactID).Return(userContact, nil)
				contact := &entities.Contact{
					ID:   contactID,
					Name: "Test Contact",
				}
				contactRepo.On("GetByID", mock.Anything, contactID).Return(contact, nil)
			},
			expectedError: nil,
			expectSuccess: true,
		},
		{
			name:      "contact access denied",
			contactID: contactID,
			userID:    otherUserID,
			setupMocks: func(contactRepo *mocks.MockContactRepository, userRepo *mocks.MockUserRepository) {
				contactRepo.On("GetUserContactRelation", mock.Anything, otherUserID, contactID).Return(nil, entities.ErrContactNotFound)
			},
			expectedError: entities.ErrContactNotFound,
			expectSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			contactRepo := &mocks.MockContactRepository{}
			userRepo := &mocks.MockUserRepository{}
			tt.setupMocks(contactRepo, userRepo)

			// Create service
			contactService := services.NewContactService(contactRepo, userRepo)

			// Execute
			ctx := context.Background()
			result, err := contactService.GetContact(ctx, tt.contactID, tt.userID)

			// Assert
			if tt.expectSuccess {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.contactID, result.ID)
			} else {
				assert.Error(t, err)
				assert.Nil(t, result)
			}

			// Verify mock expectations
			contactRepo.AssertExpectations(t)
			userRepo.AssertExpectations(t)
		})
	}
}

func TestContactService_UpdateContact(t *testing.T) {
	userID := uuid.New()
	contactID := uuid.New()
	existingUserID := uuid.New()

	tests := []struct {
		name          string
		contactID     uuid.UUID
		userID        uuid.UUID
		request       *entities.UpdateContactRequest
		setupMocks    func(*mocks.MockContactRepository, *mocks.MockUserRepository)
		expectedError error
		expectSuccess bool
	}{
		{
			name:      "successful contact update",
			contactID: contactID,
			userID:    userID,
			request: &entities.UpdateContactRequest{
				Name:  stringPtr("Updated Name"),
				Phone: stringPtr("+9876543210"),
			},
			setupMocks: func(contactRepo *mocks.MockContactRepository, userRepo *mocks.MockUserRepository) {
				userContact := &entities.UserContact{
					ID:        uuid.New(),
					UserID:    userID,
					ContactID: contactID,
				}
				contactRepo.On("GetUserContactRelation", mock.Anything, userID, contactID).Return(userContact, nil)
				contact := &entities.Contact{
					ID:        contactID,
					Name:      "Original Name",
					Phone:     stringPtr("+1234567890"),
					IsUser:    false,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				contactRepo.On("GetByID", mock.Anything, contactID).Return(contact, nil)
				contactRepo.On("Update", mock.Anything, mock.AnythingOfType("*entities.Contact")).Return(nil)
			},
			expectedError: nil,
			expectSuccess: true,
		},
		{
			name:      "update contact email to existing user",
			contactID: contactID,
			userID:    userID,
			request: &entities.UpdateContactRequest{
				Email: stringPtr("existing.user@example.com"),
			},
			setupMocks: func(contactRepo *mocks.MockContactRepository, userRepo *mocks.MockUserRepository) {
				userContact := &entities.UserContact{
					ID:        uuid.New(),
					UserID:    userID,
					ContactID: contactID,
				}
				contactRepo.On("GetUserContactRelation", mock.Anything, userID, contactID).Return(userContact, nil)
				contact := &entities.Contact{
					ID:        contactID,
					Name:      "Contact Name",
					IsUser:    false,
					UserIDRef: nil,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				contactRepo.On("GetByID", mock.Anything, contactID).Return(contact, nil)
				existingUser := &entities.User{
					ID:        existingUserID,
					Email:     "existing.user@example.com",
					FirstName: "Existing",
					LastName:  "User",
				}
				userRepo.On("GetByEmail", mock.Anything, "existing.user@example.com").Return(existingUser, nil)
				contactRepo.On("Update", mock.Anything, mock.AnythingOfType("*entities.Contact")).Return(nil)
			},
			expectedError: nil,
			expectSuccess: true,
		},
		{
			name:      "invalid contact name",
			contactID: contactID,
			userID:    userID,
			request: &entities.UpdateContactRequest{
				Name: stringPtr(""),
			},
			setupMocks:    func(contactRepo *mocks.MockContactRepository, userRepo *mocks.MockUserRepository) {},
			expectedError: entities.ErrInvalidContactName,
			expectSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			contactRepo := &mocks.MockContactRepository{}
			userRepo := &mocks.MockUserRepository{}
			tt.setupMocks(contactRepo, userRepo)

			// Create service
			contactService := services.NewContactService(contactRepo, userRepo)

			// Execute
			ctx := context.Background()
			result, err := contactService.UpdateContact(ctx, tt.contactID, tt.userID, tt.request)

			// Assert
			if tt.expectSuccess {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			} else {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedError)
				assert.Nil(t, result)
			}

			// Verify mock expectations
			contactRepo.AssertExpectations(t)
			userRepo.AssertExpectations(t)
		})
	}
}

func TestContactService_CreateReciprocalContact(t *testing.T) {
	userAID := uuid.New()
	userBID := uuid.New()

	tests := []struct {
		name          string
		contactEmail  string
		contactOwnerID uuid.UUID
		setupMocks    func(*mocks.MockContactRepository, *mocks.MockUserRepository)
		expectedError error
		expectSuccess bool
	}{
		{
			name:           "create reciprocal contact successfully",
			contactEmail:   "userb@example.com",
			contactOwnerID: userAID,
			setupMocks: func(contactRepo *mocks.MockContactRepository, userRepo *mocks.MockUserRepository) {
				// User B exists
				userB := &entities.User{
					ID:        userBID,
					Email:     "userb@example.com",
					FirstName: "User",
					LastName:  "B",
				}
				userRepo.On("GetByEmail", mock.Anything, "userb@example.com").Return(userB, nil)
				
				// User A (contact owner) details
				userA := &entities.User{
					ID:        userAID,
					Email:     "usera@example.com",
					FirstName: "User",
					LastName:  "A",
				}
				userRepo.On("GetByID", mock.Anything, userAID).Return(userA, nil)
				
				// Check if reciprocal contact already exists
				contactRepo.On("ExistsByEmailForUser", mock.Anything, userAID, "usera@example.com").Return(false, nil)
				
				// No existing contact with User A's email
				contactRepo.On("GetByEmail", mock.Anything, "usera@example.com").Return(nil, entities.ErrContactNotFound)
				
				// Create contact and relation
				contactRepo.On("Create", mock.Anything, mock.AnythingOfType("*entities.Contact")).Return(nil)
				contactRepo.On("CreateUserContactRelation", mock.Anything, mock.AnythingOfType("*entities.UserContact")).Return(nil)
			},
			expectedError: nil,
			expectSuccess: true,
		},
		{
			name:           "reciprocal contact already exists",
			contactEmail:   "userb@example.com",
			contactOwnerID: userAID,
			setupMocks: func(contactRepo *mocks.MockContactRepository, userRepo *mocks.MockUserRepository) {
				// User B exists
				userB := &entities.User{
					ID:        userBID,
					Email:     "userb@example.com",
					FirstName: "User",
					LastName:  "B",
				}
				userRepo.On("GetByEmail", mock.Anything, "userb@example.com").Return(userB, nil)
				
				// User A (contact owner) details
				userA := &entities.User{
					ID:        userAID,
					Email:     "usera@example.com",
					FirstName: "User",
					LastName:  "A",
				}
				userRepo.On("GetByID", mock.Anything, userAID).Return(userA, nil)
				
				// Reciprocal contact already exists
				contactRepo.On("ExistsByEmailForUser", mock.Anything, userAID, "usera@example.com").Return(true, nil)
			},
			expectedError: nil,
			expectSuccess: true,
		},
		{
			name:           "contact email user does not exist",
			contactEmail:   "nonexistent@example.com",
			contactOwnerID: userAID,
			setupMocks: func(contactRepo *mocks.MockContactRepository, userRepo *mocks.MockUserRepository) {
				userRepo.On("GetByEmail", mock.Anything, "nonexistent@example.com").Return(nil, entities.ErrUserNotFound)
			},
			expectedError: nil,
			expectSuccess: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			contactRepo := &mocks.MockContactRepository{}
			userRepo := &mocks.MockUserRepository{}
			tt.setupMocks(contactRepo, userRepo)

			// Create service
			contactService := services.NewContactService(contactRepo, userRepo)

			// Execute
			ctx := context.Background()
			err := contactService.CreateReciprocalContact(ctx, tt.contactEmail, tt.contactOwnerID)

			// Assert
			if tt.expectSuccess {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedError)
			}

			// Verify mock expectations
			contactRepo.AssertExpectations(t)
			userRepo.AssertExpectations(t)
		})
	}
}
