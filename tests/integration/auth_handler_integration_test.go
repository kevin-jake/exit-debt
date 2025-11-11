package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"pay-your-dues/internal/domain/entities"
	"pay-your-dues/internal/handlers"
	"pay-your-dues/internal/mocks"
)

func TestAuthHandler_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)
	

	tests := []struct {
		name           string
		requestBody    interface{}
		setupMock      func(*mocks.MockAuthService)
		expectedStatus int
		validateBody   func(*testing.T, map[string]interface{})
	}{
		{
			name: "successful registration",
			requestBody: map[string]interface{}{
				"email":      "test@example.com",
				"password":   "password123",
				"first_name": "John",
				"last_name":  "Doe",
				"phone":      "+1234567890",
			},
			setupMock: func(mockAuthService *mocks.MockAuthService) {
				expectedUser := &entities.User{
					ID:        uuid.New(),
					Email:     "test@example.com",
					FirstName: "John",
					LastName:  "Doe",
					Phone:     stringPtr("+1234567890"),
				}
				mockAuthService.On("Register", mock.Anything, mock.AnythingOfType("*entities.CreateUserRequest")).Return(&entities.RegisterResponse{
					User: *expectedUser,
				}, nil)
			},
			expectedStatus: http.StatusCreated,
			validateBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "User registered successfully", body["message"])
				assert.NotNil(t, body["data"])
				userData := body["data"].(map[string]interface{})
				assert.NotNil(t, userData["user"])
				user := userData["user"].(map[string]interface{})
				assert.Equal(t, "test@example.com", user["Email"])
				assert.Equal(t, "John", user["FirstName"])
				assert.Equal(t, "Doe", user["LastName"])
				assert.Equal(t, "+1234567890", user["Phone"])
			},
		},
		{
			name: "user already exists",
			requestBody: map[string]interface{}{
				"email":      "existing@example.com",
				"password":   "password123",
				"first_name": "Jane",
				"last_name":  "Smith",
			},
			setupMock: func(mockAuthService *mocks.MockAuthService) {
				mockAuthService.On("Register", mock.Anything, mock.AnythingOfType("*entities.CreateUserRequest")).Return(nil, entities.ErrUserAlreadyExists)
			},
			expectedStatus: http.StatusConflict,
			validateBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "User already exists", body["error"])
			},
		},
		{
			name: "invalid request body",
			requestBody: map[string]interface{}{
				"email":    "invalid-email",
				"password": "123",
			},
			setupMock: func(mockAuthService *mocks.MockAuthService) {
				mockAuthService.On("Register", mock.Anything, mock.AnythingOfType("*entities.CreateUserRequest")).Return(nil, entities.ErrInvalidEmail)
			},
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Invalid input", body["error"])
			},
		},
		{
			name:        "malformed JSON",
			requestBody: "invalid json",
			setupMock: func(mockAuthService *mocks.MockAuthService) {
				// No mock setup needed
			},
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Invalid request body", body["error"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockAuthService := &mocks.MockAuthService{}
			tt.setupMock(mockAuthService)

			logger := zerolog.New(nil)
			authHandler := handlers.NewAuthHandler(mockAuthService, logger)

			// Prepare request
			var requestBody []byte
			if str, ok := tt.requestBody.(string); ok {
				requestBody = []byte(str)
			} else {
				requestBody, _ = json.Marshal(tt.requestBody)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Setup Gin context
			router := gin.New()
			router.POST("/api/auth/register", authHandler.Register)
			
			// Execute
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			var responseBody map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &responseBody)
			assert.NoError(t, err)

			tt.validateBody(t, responseBody)

			// Verify mock expectations
			mockAuthService.AssertExpectations(t)
		})
	}
	
}

func TestAuthHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    interface{}
		setupMock      func(*mocks.MockAuthService)
		expectedStatus int
		validateBody   func(*testing.T, map[string]interface{})
	}{
		{
			name: "successful login",
			requestBody: map[string]interface{}{
				"email":    "test@example.com",
				"password": "password123",
			},
			setupMock: func(mockAuthService *mocks.MockAuthService) {
				expectedUser := &entities.User{
					ID:        uuid.New(),
					Email:     "test@example.com",
					FirstName: "John",
					LastName:  "Doe",
				}
				mockAuthService.On("Login", mock.Anything, mock.AnythingOfType("*entities.LoginRequest")).Return(&entities.LoginResponse{
					Token: "jwt.token.here",
					User:  *expectedUser,
				}, nil)
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Login successful", body["message"])
				assert.NotNil(t, body["data"])
				loginData := body["data"].(map[string]interface{})
				assert.Equal(t, "jwt.token.here", loginData["token"])
				assert.NotNil(t, loginData["user"])
			},
		},
		{
			name: "invalid credentials",
			requestBody: map[string]interface{}{
				"email":    "test@example.com",
				"password": "wrongpassword",
			},
			setupMock: func(mockAuthService *mocks.MockAuthService) {
				mockAuthService.On("Login", mock.Anything, mock.AnythingOfType("*entities.LoginRequest")).Return(nil, entities.ErrInvalidCredentials)
			},
			expectedStatus: http.StatusUnauthorized,
			validateBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Invalid credentials", body["error"])
			},
		},
		{
			name: "missing email",
			requestBody: map[string]interface{}{
				"password": "password123",
			},
			setupMock: func(mockAuthService *mocks.MockAuthService) {
				// The handler will call the service with empty email, so we need to mock it
				mockAuthService.On("Login", mock.Anything, mock.AnythingOfType("*entities.LoginRequest")).Return(nil, entities.ErrInvalidEmail)
			},
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Invalid input", body["error"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockAuthService := &mocks.MockAuthService{}
			tt.setupMock(mockAuthService)

			logger := zerolog.New(nil)
			authHandler := handlers.NewAuthHandler(mockAuthService, logger)

			// Prepare request
			requestBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Setup Gin context
			router := gin.New()
			router.POST("/api/auth/login", authHandler.Login)
			
			// Execute
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			var responseBody map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &responseBody)
			assert.NoError(t, err)

			tt.validateBody(t, responseBody)

			// Verify mock expectations
			mockAuthService.AssertExpectations(t)
		})
	}
}
