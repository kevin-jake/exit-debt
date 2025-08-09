package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"exit-debt/internal/domain/interfaces"
	"exit-debt/internal/handlers"
	"exit-debt/internal/middleware"
	"exit-debt/internal/models"
	"exit-debt/internal/repository"
	"exit-debt/internal/services"
)

type AuthAPITestSuite struct {
	suite.Suite
	db             *gorm.DB
	router         *gin.Engine
	authService    interfaces.AuthService
	contactService interfaces.ContactService
	authHandler    *handlers.AuthHandler
}

func (suite *AuthAPITestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)

	// Setup in-memory SQLite database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	suite.Require().NoError(err)

	// Auto-migrate models
	err = db.AutoMigrate(
		&models.User{},
		&models.Contact{},
		&models.UserContact{},
	)
	suite.Require().NoError(err)

	suite.db = db

	// Initialize repositories and services
	userRepo := repository.NewUserRepositoryGORM(db)
	contactRepo := repository.NewContactRepositoryGORM(db)
	
	suite.contactService = services.NewContactService(contactRepo, userRepo)
	
	authService, err := services.NewAuthService(userRepo, suite.contactService, "test-secret-key", "24h")
	suite.Require().NoError(err)
	suite.authService = authService

	// Initialize handlers
	logger := zerolog.New(nil)
	suite.authHandler = handlers.NewAuthHandler(suite.authService, logger)

	// Setup router
	suite.router = gin.New()
	suite.setupRoutes()
}

func (suite *AuthAPITestSuite) setupRoutes() {
	api := suite.router.Group("/api")
	auth := api.Group("/auth")
	{
		auth.POST("/register", suite.authHandler.Register)
		auth.POST("/login", suite.authHandler.Login)
	}

	// Protected routes for testing authorization
	protected := api.Group("/protected")
	protected.Use(middleware.AuthMiddleware(suite.authService))
	{
		protected.GET("/profile", func(c *gin.Context) {
			userID, _ := c.Get("user_id")
			c.JSON(http.StatusOK, gin.H{
				"message": "Protected endpoint accessed",
				"user_id": userID,
			})
		})
	}
}

func (suite *AuthAPITestSuite) TearDownTest() {
	// Clean up database after each test
	suite.db.Exec("DELETE FROM user_contacts")
	suite.db.Exec("DELETE FROM contacts")
	suite.db.Exec("DELETE FROM users")
}

func (suite *AuthAPITestSuite) TestUserRegistration() {
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
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
			expectedStatus: http.StatusCreated,
			validateBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "User registered successfully", body["message"])
				assert.True(t, body["success"].(bool))
				assert.NotNil(t, body["data"])
				
				userData := body["data"].(map[string]interface{})
				assert.Equal(t, "test@example.com", userData["email"])
				assert.Equal(t, "John", userData["first_name"])
				assert.Equal(t, "Doe", userData["last_name"])
				assert.NotEmpty(t, userData["id"])
			},
		},
		{
			name: "duplicate email registration",
			requestBody: map[string]interface{}{
				"email":      "test@example.com", // Same email as above
				"password":   "password456",
				"first_name": "Jane",
				"last_name":  "Smith",
			},
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Registration failed", body["message"])
				assert.False(t, body["success"].(bool))
			},
		},
		{
			name: "invalid email format",
			requestBody: map[string]interface{}{
				"email":      "invalid-email",
				"password":   "password123",
				"first_name": "John",
				"last_name":  "Doe",
			},
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Invalid request body", body["message"])
				assert.False(t, body["success"].(bool))
			},
		},
		{
			name: "password too short",
			requestBody: map[string]interface{}{
				"email":      "short@example.com",
				"password":   "123",
				"first_name": "John",
				"last_name":  "Doe",
			},
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Invalid request body", body["message"])
				assert.False(t, body["success"].(bool))
			},
		},
		{
			name: "missing required fields",
			requestBody: map[string]interface{}{
				"email":    "missing@example.com",
				"password": "password123",
				// Missing first_name and last_name
			},
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Invalid request body", body["message"])
				assert.False(t, body["success"].(bool))
			},
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			// Prepare request
			requestBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Execute
			suite.router.ServeHTTP(w, req)

			// Assert
			assert.Equal(suite.T(), tt.expectedStatus, w.Code)

			var responseBody map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &responseBody)
			assert.NoError(suite.T(), err)

			tt.validateBody(suite.T(), responseBody)
		})
	}
}

func (suite *AuthAPITestSuite) TestUserLogin() {
	// First, register a user to login with
	suite.registerTestUser("logintest@example.com", "password123", "Login", "Test")

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
		validateBody   func(*testing.T, map[string]interface{})
	}{
		{
			name: "successful login",
			requestBody: map[string]interface{}{
				"email":    "logintest@example.com",
				"password": "password123",
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Login successful", body["message"])
				assert.True(t, body["success"].(bool))
				assert.NotNil(t, body["data"])
				
				loginData := body["data"].(map[string]interface{})
				assert.NotEmpty(t, loginData["token"])
				assert.NotNil(t, loginData["user"])
				
				userData := loginData["user"].(map[string]interface{})
				assert.Equal(t, "logintest@example.com", userData["email"])
			},
		},
		{
			name: "wrong password",
			requestBody: map[string]interface{}{
				"email":    "logintest@example.com",
				"password": "wrongpassword",
			},
			expectedStatus: http.StatusUnauthorized,
			validateBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Invalid credentials", body["message"])
				assert.False(t, body["success"].(bool))
			},
		},
		{
			name: "non-existent email",
			requestBody: map[string]interface{}{
				"email":    "nonexistent@example.com",
				"password": "password123",
			},
			expectedStatus: http.StatusUnauthorized,
			validateBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Invalid credentials", body["message"])
				assert.False(t, body["success"].(bool))
			},
		},
		{
			name: "invalid request format",
			requestBody: map[string]interface{}{
				"email": "logintest@example.com",
				// Missing password
			},
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Invalid request body", body["message"])
				assert.False(t, body["success"].(bool))
			},
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			// Prepare request
			requestBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Execute
			suite.router.ServeHTTP(w, req)

			// Assert
			assert.Equal(suite.T(), tt.expectedStatus, w.Code)

			var responseBody map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &responseBody)
			assert.NoError(suite.T(), err)

			tt.validateBody(suite.T(), responseBody)
		})
	}
}

func (suite *AuthAPITestSuite) TestAuthenticationMiddleware() {
	// Register and login to get a valid token
	suite.registerTestUser("middleware@example.com", "password123", "Middleware", "Test")
	token := suite.loginAndGetToken("middleware@example.com", "password123")

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		validateBody   func(*testing.T, map[string]interface{})
	}{
		{
			name:           "valid token",
			authHeader:     "Bearer " + token,
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Protected endpoint accessed", body["message"])
				assert.NotNil(t, body["user_id"])
			},
		},
		{
			name:           "missing authorization header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			validateBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Authorization header required", body["message"])
				assert.False(t, body["success"].(bool))
			},
		},
		{
			name:           "invalid token format",
			authHeader:     "InvalidToken",
			expectedStatus: http.StatusUnauthorized,
			validateBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Invalid authorization header format", body["message"])
				assert.False(t, body["success"].(bool))
			},
		},
		{
			name:           "malformed token",
			authHeader:     "Bearer invalid.token.here",
			expectedStatus: http.StatusUnauthorized,
			validateBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Invalid token", body["message"])
				assert.False(t, body["success"].(bool))
			},
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			// Prepare request
			req := httptest.NewRequest(http.MethodGet, "/api/protected/profile", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			w := httptest.NewRecorder()

			// Execute
			suite.router.ServeHTTP(w, req)

			// Assert
			assert.Equal(suite.T(), tt.expectedStatus, w.Code)

			var responseBody map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &responseBody)
			assert.NoError(suite.T(), err)

			tt.validateBody(suite.T(), responseBody)
		})
	}
}

func (suite *AuthAPITestSuite) TestTokenExpiration() {
	// Create auth service with short expiry for testing
	userRepo := repository.NewUserRepositoryGORM(suite.db)
	shortExpiryAuthService, err := services.NewAuthService(userRepo, suite.contactService, "test-secret-key", "1ns") // 1 nanosecond
	suite.Require().NoError(err)

	// Register user with short expiry service
	suite.registerTestUser("expiry@example.com", "password123", "Expiry", "Test")

	// Login with short expiry service to get an expired token
	loginReq := map[string]interface{}{
		"email":    "expiry@example.com",
		"password": "password123",
	}

	requestBody, _ := json.Marshal(loginReq)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Temporarily replace auth service
	originalAuthService := suite.authService
	suite.authService = shortExpiryAuthService
	logger := zerolog.New(nil)
	suite.authHandler = handlers.NewAuthHandler(suite.authService, logger)
	suite.setupRoutes()

	suite.router.ServeHTTP(w, req)

	var loginResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &loginResponse)
	suite.NoError(err)

	loginData := loginResponse["data"].(map[string]interface{})
	expiredToken := loginData["token"].(string)

	// Wait a moment to ensure token is expired
	time.Sleep(10 * time.Millisecond)

	// Restore original auth service
	suite.authService = originalAuthService
	logger = zerolog.New(nil)
	suite.authHandler = handlers.NewAuthHandler(suite.authService, logger)
	suite.setupRoutes()

	// Try to access protected endpoint with expired token
	req = httptest.NewRequest(http.MethodGet, "/api/protected/profile", nil)
	req.Header.Set("Authorization", "Bearer "+expiredToken)
	w = httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)

	var responseBody map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Token expired", responseBody["message"])
}

// Helper methods
func (suite *AuthAPITestSuite) registerTestUser(email, password, firstName, lastName string) {
	registerReq := map[string]interface{}{
		"email":      email,
		"password":   password,
		"first_name": firstName,
		"last_name":  lastName,
	}

	requestBody, _ := json.Marshal(registerReq)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusCreated, w.Code)
}

func (suite *AuthAPITestSuite) loginAndGetToken(email, password string) string {
	loginReq := map[string]interface{}{
		"email":    email,
		"password": password,
	}

	requestBody, _ := json.Marshal(loginReq)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var responseBody map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.NoError(suite.T(), err)

	loginData := responseBody["data"].(map[string]interface{})
	return loginData["token"].(string)
}

func TestAuthAPITestSuite(t *testing.T) {
	suite.Run(t, new(AuthAPITestSuite))
}
