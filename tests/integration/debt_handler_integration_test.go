package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"pay-your-dues/internal/domain/entities"
	"pay-your-dues/internal/handlers"
	"pay-your-dues/internal/mocks"
)

func TestDebtHandler_CreateDebtList(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	contactID := uuid.New()
	futureDate := time.Now().AddDate(0, 0, 30)

	tests := []struct {
		name           string
		requestBody    interface{}
		setupMock      func(*mocks.MockDebtService)
		setupContext   func(*gin.Context)
		expectedStatus int
		validateBody   func(*testing.T, map[string]interface{})
	}{
		{
			name: "successful debt list creation",
			requestBody: map[string]interface{}{
				"contact_id":     contactID.String(),
				"debt_type":      "to_pay",
				"total_amount":   "1000.00",
				"currency":       "USD",
				"due_date":       futureDate.Format(time.RFC3339),
				"installment_plan": "onetime",
				"description":    "Loan from friend",
			},
			setupMock: func(mockDebtService *mocks.MockDebtService) {
				expectedDebtList := &entities.DebtList{
					ID:                uuid.New(),
					UserID:            userID,
					ContactID:         contactID,
					DebtType:          "to_pay",
					TotalAmount:       decimal.RequireFromString("1000.00"),
					InstallmentAmount: decimal.RequireFromString("1000.00"),
					Currency:          "USD",
					Status:            "active",
					DueDate:           futureDate,
					InstallmentPlan:   "onetime",
					Description:       stringPtr("Loan from friend"),
				}
				mockDebtService.On("CreateDebtList", mock.Anything, userID, mock.AnythingOfType("*entities.CreateDebtListRequest")).Return(expectedDebtList, nil)
			},
			setupContext: func(c *gin.Context) {
				c.Set("user_id", userID)
			},
			expectedStatus: http.StatusCreated,
			validateBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Debt list created successfully", body["message"])
				assert.NotNil(t, body["data"])
				debtData := body["data"].(map[string]interface{})
				assert.Equal(t, "to_pay", debtData["DebtType"])
				assert.Equal(t, "1000", debtData["TotalAmount"])
			},
		},
		{
			name: "unauthorized user",
			requestBody: map[string]interface{}{
				"contact_id":   contactID.String(),
				"debt_type":    "to_pay",
				"total_amount": "500.00",
			},
			setupMock: func(mockDebtService *mocks.MockDebtService) {
				// No mock setup needed
			},
			setupContext: func(c *gin.Context) {
				// Don't set user_id to simulate unauthorized access
			},
			expectedStatus: http.StatusUnauthorized,
			validateBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Unauthorized", body["error"])
			},
		},
		{
			name: "invalid debt type",
			requestBody: map[string]interface{}{
				"contact_id":   contactID.String(),
				"debt_type":    "invalid_type",
				"total_amount": "500.00",
			},
			setupMock: func(mockDebtService *mocks.MockDebtService) {
				mockDebtService.On("CreateDebtList", mock.Anything, userID, mock.AnythingOfType("*entities.CreateDebtListRequest")).Return(nil, entities.ErrInvalidDebtType)
			},
			setupContext: func(c *gin.Context) {
				c.Set("user_id", userID)
			},
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Invalid input", body["error"])
			},
		},
		{
			name: "contact not found",
			requestBody: map[string]interface{}{
				"contact_id":   uuid.New().String(),
				"debt_type":    "to_pay",
				"total_amount": "500.00",
			},
			setupMock: func(mockDebtService *mocks.MockDebtService) {
				mockDebtService.On("CreateDebtList", mock.Anything, userID, mock.AnythingOfType("*entities.CreateDebtListRequest")).Return(nil, entities.ErrContactNotFound)
			},
			setupContext: func(c *gin.Context) {
				c.Set("user_id", userID)
			},
			expectedStatus: http.StatusNotFound,
			validateBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Contact not found", body["error"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockDebtService := &mocks.MockDebtService{}
			mockFileStorageService := &mocks.MockFileStorageService{}
			tt.setupMock(mockDebtService)

			logger := zerolog.New(nil)
			debtHandler := handlers.NewDebtHandler(mockDebtService, mockFileStorageService, logger)

			// Prepare request
			requestBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/debt-lists", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Setup Gin context
			router := gin.New()
			router.POST("/api/debt-lists", func(c *gin.Context) {
				tt.setupContext(c)
				debtHandler.CreateDebtList(c)
			})
			
			// Execute
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			var responseBody map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &responseBody)
			assert.NoError(t, err)

			tt.validateBody(t, responseBody)

			// Verify mock expectations
			mockDebtService.AssertExpectations(t)
		})
	}
}

func TestDebtHandler_GetUserDebtLists(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()

	tests := []struct {
		name           string
		setupMock      func(*mocks.MockDebtService)
		setupContext   func(*gin.Context)
		expectedStatus int
		validateBody   func(*testing.T, map[string]interface{})
	}{
		{
			name: "successful debt lists retrieval",
			setupMock: func(mockDebtService *mocks.MockDebtService) {
				debtLists := []entities.DebtListResponse{
					{
						ID:            uuid.New(),
						UserID:        userID,
						DebtType:      "to_pay",
						TotalAmount:   decimal.RequireFromString("1000.00"),
						Currency:      "USD",
						Status:        "active",
					},
					{
						ID:            uuid.New(),
						UserID:        userID,
						DebtType:      "to_receive",
						TotalAmount:   decimal.RequireFromString("500.00"),
						Currency:      "USD",
						Status:        "active",
					},
				}
				mockDebtService.On("GetUserDebtLists", mock.Anything, userID).Return(debtLists, nil)
			},
			setupContext: func(c *gin.Context) {
				c.Set("user_id", userID)
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Debt lists retrieved successfully", body["message"])
				assert.NotNil(t, body["data"])
				debtLists := body["data"].([]interface{})
				assert.Len(t, debtLists, 2)
			},
		},
		{
			name: "empty debt lists",
			setupMock: func(mockDebtService *mocks.MockDebtService) {
				mockDebtService.On("GetUserDebtLists", mock.Anything, userID).Return([]entities.DebtListResponse{}, nil)
			},
			setupContext: func(c *gin.Context) {
				c.Set("user_id", userID)
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Debt lists retrieved successfully", body["message"])
				assert.NotNil(t, body["data"])
				debtLists := body["data"].([]interface{})
				assert.Len(t, debtLists, 0)
			},
		},
		{
			name: "unauthorized access",
			setupMock: func(mockDebtService *mocks.MockDebtService) {
				// No mock setup needed
			},
			setupContext: func(c *gin.Context) {
				// Don't set user_id
			},
			expectedStatus: http.StatusUnauthorized,
			validateBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Unauthorized", body["error"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockDebtService := &mocks.MockDebtService{}
			tt.setupMock(mockDebtService)

			logger := zerolog.New(nil)
			mockFileStorageService := &mocks.MockFileStorageService{}
			debtHandler := handlers.NewDebtHandler(mockDebtService, mockFileStorageService, logger)

			// Prepare request
			req := httptest.NewRequest(http.MethodGet, "/api/debt-lists", nil)
			w := httptest.NewRecorder()

			// Setup Gin context
			router := gin.New()
			router.GET("/api/debt-lists", func(c *gin.Context) {
				tt.setupContext(c)
				debtHandler.GetUserDebtLists(c)
			})
			
			// Execute
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			var responseBody map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &responseBody)
			assert.NoError(t, err)

			tt.validateBody(t, responseBody)

			// Verify mock expectations
			mockDebtService.AssertExpectations(t)
		})
	}
}

func TestDebtHandler_CreateDebtItem(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	debtListID := uuid.New()
	paymentDate := time.Now()

	tests := []struct {
		name           string
		requestBody    interface{}
		setupMock      func(*mocks.MockDebtService)
		setupContext   func(*gin.Context)
		expectedStatus int
		validateBody   func(*testing.T, map[string]interface{})
	}{
		{
			name: "successful payment creation",
			requestBody: map[string]interface{}{
				"debt_list_id":    debtListID.String(),
				"amount":          "200.00",
				"currency":        "USD",
				"payment_date":    paymentDate.Format(time.RFC3339),
				"payment_method":  "bank_transfer",
				"description":     "First installment",
			},
			setupMock: func(mockDebtService *mocks.MockDebtService) {
				expectedDebtItem := &entities.DebtItem{
					ID:            uuid.New(),
					DebtListID:    debtListID,
					Amount:        decimal.RequireFromString("200.00"),
					Currency:      "USD",
					PaymentDate:   paymentDate,
					PaymentMethod: "bank_transfer",
					Description:   stringPtr("First installment"),
					Status:        "completed",
				}
				mockDebtService.On("CreateDebtItem", mock.Anything, userID, mock.AnythingOfType("*entities.CreateDebtItemRequest")).Return(expectedDebtItem, nil)
			},
			setupContext: func(c *gin.Context) {
				c.Set("user_id", userID)
			},
			expectedStatus: http.StatusCreated,
			validateBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Payment recorded successfully", body["message"])
				assert.NotNil(t, body["data"])
				paymentData := body["data"].(map[string]interface{})
				// Amount is returned as a string from decimal.Decimal
				assert.Equal(t, "200", paymentData["Amount"])
				assert.Equal(t, "bank_transfer", paymentData["PaymentMethod"])
			},
		},
		{
			name: "invalid amount",
			requestBody: map[string]interface{}{
				"debt_list_id":   debtListID.String(),
				"amount":         "0.00",
				"payment_date":   paymentDate.Format(time.RFC3339),
				"payment_method": "cash",
			},
			setupMock: func(mockDebtService *mocks.MockDebtService) {
				mockDebtService.On("CreateDebtItem", mock.Anything, userID, mock.AnythingOfType("*entities.CreateDebtItemRequest")).Return(nil, entities.ErrInvalidAmount)
			},
			setupContext: func(c *gin.Context) {
				c.Set("user_id", userID)
			},
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Invalid input", body["error"])
			},
		},
		{
			name: "debt list not found",
			requestBody: map[string]interface{}{
				"debt_list_id":   uuid.New().String(),
				"amount":         "100.00",
				"payment_date":   paymentDate.Format(time.RFC3339),
				"payment_method": "cash",
			},
			setupMock: func(mockDebtService *mocks.MockDebtService) {
				mockDebtService.On("CreateDebtItem", mock.Anything, userID, mock.AnythingOfType("*entities.CreateDebtItemRequest")).Return(nil, entities.ErrDebtListNotFound)
			},
			setupContext: func(c *gin.Context) {
				c.Set("user_id", userID)
			},
			expectedStatus: http.StatusNotFound,
			validateBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Debt list not found", body["error"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockDebtService := &mocks.MockDebtService{}
			tt.setupMock(mockDebtService)

			logger := zerolog.New(nil)
			mockFileStorageService := &mocks.MockFileStorageService{}
			debtHandler := handlers.NewDebtHandler(mockDebtService, mockFileStorageService, logger)

			// Prepare request
			requestBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/debt-items", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Setup Gin context
			router := gin.New()
			router.POST("/api/debt-items", func(c *gin.Context) {
				tt.setupContext(c)
				debtHandler.CreateDebtItem(c)
			})
			
			// Execute
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			var responseBody map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &responseBody)
			assert.NoError(t, err)



			tt.validateBody(t, responseBody)

			// Verify mock expectations
			mockDebtService.AssertExpectations(t)
		})
	}
}
