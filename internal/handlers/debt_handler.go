package handlers

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"exit-debt/internal/domain/entities"
	"exit-debt/internal/domain/interfaces"
)

// DebtHandler handles debt-related HTTP requests
type DebtHandler struct {
	debtService        interfaces.DebtService
	fileStorageService interfaces.FileStorageService
	logger             zerolog.Logger
}

// NewDebtHandler creates a new debt handler
func NewDebtHandler(debtService interfaces.DebtService, fileStorageService interfaces.FileStorageService, logger zerolog.Logger) *DebtHandler {
	return &DebtHandler{
		debtService:        debtService,
		fileStorageService: fileStorageService,
		logger:             logger.With().Str("handler", "debt").Logger(),
	}
}

// CreateDebtList handles debt list creation
func (h *DebtHandler) CreateDebtList(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	// Extract request ID and user ID for logging
	requestID := getRequestID(c)
	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Error().Str("request_id", requestID).Msg("User ID not found in context")
		c.JSON(http.StatusUnauthorized, NewErrorResponse("Unauthorized", "", requestID))
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		h.logger.Error().Str("request_id", requestID).Msg("Invalid user ID type in context")
		c.JSON(http.StatusUnauthorized, NewErrorResponse("Unauthorized", "", requestID))
		return
	}

	logger := h.logger.With().Str("request_id", requestID).Str("user_id", userUUID.String()).Str("method", "CreateDebtList").Logger()

	var req entities.CreateDebtListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn().Err(err).Msg("Invalid request body")
		c.JSON(http.StatusBadRequest, NewErrorResponse("Invalid request body", err.Error(), requestID))
		return
	}

	// Sanitize input
	req.TotalAmount = sanitizeString(req.TotalAmount)
	req.Currency = sanitizeString(req.Currency)
	req.InstallmentPlan = sanitizeString(req.InstallmentPlan)
	req.DebtType = sanitizeString(req.DebtType)
	if req.Description != nil {
		sanitized := sanitizeString(*req.Description)
		req.Description = &sanitized
	}
	if req.Notes != nil {
		sanitized := sanitizeString(*req.Notes)
		req.Notes = &sanitized
	}

	logger.Info().Str("debt_type", req.DebtType).Str("contact_id", req.ContactID.String()).Msg("Debt list creation attempt")

	debtList, err := h.debtService.CreateDebtList(ctx, userUUID, &req)
	if err != nil {
		logger.Error().Err(err).Str("debt_type", req.DebtType).Msg("Debt list creation failed")

		// Handle specific error types
		switch err {
		case entities.ErrInvalidDebtType, entities.ErrInvalidAmount, entities.ErrInvalidCurrency, entities.ErrInvalidDueDate:
			c.JSON(http.StatusBadRequest, NewErrorResponse("Invalid input", err.Error(), requestID))
		case entities.ErrContactNotFound:
			c.JSON(http.StatusNotFound, NewErrorResponse("Contact not found", "", requestID))
		default:
			c.JSON(http.StatusInternalServerError, NewErrorResponse("Internal server error", "", requestID))
		}
		return
	}

	logger.Info().Str("debt_list_id", debtList.ID.String()).Str("debt_type", req.DebtType).Msg("Debt list created successfully")

	c.JSON(http.StatusCreated, NewSuccessResponse("Debt list created successfully", debtList, requestID))
}

// GetUserDebtLists handles retrieving all debt lists for a user
func (h *DebtHandler) GetUserDebtLists(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	// Extract request ID and user ID for logging
	requestID := getRequestID(c)
	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Error().Str("request_id", requestID).Msg("User ID not found in context")
		c.JSON(http.StatusUnauthorized, NewErrorResponse("Unauthorized", "", requestID))
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		h.logger.Error().Str("request_id", requestID).Msg("Invalid user ID type in context")
		c.JSON(http.StatusUnauthorized, NewErrorResponse("Unauthorized", "", requestID))
		return
	}

	logger := h.logger.With().Str("request_id", requestID).Str("user_id", userUUID.String()).Str("method", "GetUserDebtLists").Logger()

	logger.Info().Msg("Retrieving user debt lists")

	debtLists, err := h.debtService.GetUserDebtLists(ctx, userUUID)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to retrieve user debt lists")
		c.JSON(http.StatusInternalServerError, NewErrorResponse("Internal server error", "", requestID))
		return
	}

	logger.Info().Int("count", len(debtLists)).Msg("User debt lists retrieved successfully")

	c.JSON(http.StatusOK, NewSuccessResponse("Debt lists retrieved successfully", debtLists, requestID))
}

// GetDebtList handles retrieving a specific debt list
func (h *DebtHandler) GetDebtList(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	// Extract request ID and user ID for logging
	requestID := getRequestID(c)
	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Error().Str("request_id", requestID).Msg("User ID not found in context")
		c.JSON(http.StatusUnauthorized, NewErrorResponse("Unauthorized", "", requestID))
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		h.logger.Error().Str("request_id", requestID).Msg("Invalid user ID type in context")
		c.JSON(http.StatusUnauthorized, NewErrorResponse("Unauthorized", "", requestID))
		return
	}

	// Parse debt list ID from URL parameter
	debtListIDStr := c.Param("id")
	debtListID, err := uuid.Parse(debtListIDStr)
	if err != nil {
		h.logger.Warn().Str("request_id", requestID).Str("debt_list_id", debtListIDStr).Msg("Invalid debt list ID format")
		c.JSON(http.StatusBadRequest, NewErrorResponse("Invalid debt list ID", "", requestID))
		return
	}

	logger := h.logger.With().Str("request_id", requestID).Str("user_id", userUUID.String()).Str("debt_list_id", debtListID.String()).Str("method", "GetDebtList").Logger()

	logger.Info().Msg("Retrieving debt list")

	debtList, err := h.debtService.GetDebtList(ctx, debtListID, userUUID)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to retrieve debt list")

		// Handle specific error types
		switch err {
		case entities.ErrDebtListNotFound:
			c.JSON(http.StatusNotFound, NewErrorResponse("Debt list not found", "", requestID))
		default:
			c.JSON(http.StatusInternalServerError, NewErrorResponse("Internal server error", "", requestID))
		}
		return
	}

	logger.Info().Msg("Debt list retrieved successfully")

	c.JSON(http.StatusOK, NewSuccessResponse("Debt list retrieved successfully", debtList, requestID))
}

// UpdateDebtList handles debt list updates
func (h *DebtHandler) UpdateDebtList(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	// Extract request ID and user ID for logging
	requestID := getRequestID(c)
	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Error().Str("request_id", requestID).Msg("User ID not found in context")
		c.JSON(http.StatusUnauthorized, NewErrorResponse("Unauthorized", "", requestID))
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		h.logger.Error().Str("request_id", requestID).Msg("Invalid user ID type in context")
		c.JSON(http.StatusUnauthorized, NewErrorResponse("Unauthorized", "", requestID))
		return
	}

	// Parse debt list ID from URL parameter
	debtListIDStr := c.Param("id")
	debtListID, err := uuid.Parse(debtListIDStr)
	if err != nil {
		h.logger.Warn().Str("request_id", requestID).Str("debt_list_id", debtListIDStr).Msg("Invalid debt list ID format")
		c.JSON(http.StatusBadRequest, NewErrorResponse("Invalid debt list ID", "", requestID))
		return
	}

	logger := h.logger.With().Str("request_id", requestID).Str("user_id", userUUID.String()).Str("debt_list_id", debtListID.String()).Str("method", "UpdateDebtList").Logger()

	var req entities.UpdateDebtListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn().Err(err).Msg("Invalid request body")
		c.JSON(http.StatusBadRequest, NewErrorResponse("Invalid request body", err.Error(), requestID))
		return
	}

	// Sanitize input
	if req.TotalAmount != nil {
		sanitized := sanitizeString(*req.TotalAmount)
		req.TotalAmount = &sanitized
	}
	if req.Currency != nil {
		sanitized := sanitizeString(*req.Currency)
		req.Currency = &sanitized
	}
	if req.Status != nil {
		sanitized := sanitizeString(*req.Status)
		req.Status = &sanitized
	}
	if req.InstallmentPlan != nil {
		sanitized := sanitizeString(*req.InstallmentPlan)
		req.InstallmentPlan = &sanitized
	}
	if req.Description != nil {
		sanitized := sanitizeString(*req.Description)
		req.Description = &sanitized
	}
	if req.Notes != nil {
		sanitized := sanitizeString(*req.Notes)
		req.Notes = &sanitized
	}

	logger.Info().Msg("Debt list update attempt")

	debtList, err := h.debtService.UpdateDebtList(ctx, debtListID, userUUID, &req)
	if err != nil {
		logger.Error().Err(err).Msg("Debt list update failed")

		// Handle specific error types
		switch err {
		case entities.ErrDebtListNotFound:
			c.JSON(http.StatusNotFound, NewErrorResponse("Debt list not found", "", requestID))
		case entities.ErrInvalidAmount, entities.ErrInvalidCurrency, entities.ErrInvalidDueDate:
			c.JSON(http.StatusBadRequest, NewErrorResponse("Invalid input", err.Error(), requestID))
		default:
			c.JSON(http.StatusInternalServerError, NewErrorResponse("Internal server error", "", requestID))
		}
		return
	}

	logger.Info().Msg("Debt list updated successfully")

	c.JSON(http.StatusOK, NewSuccessResponse("Debt list updated successfully", debtList, requestID))
}

// DeleteDebtList handles debt list deletion
func (h *DebtHandler) DeleteDebtList(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	// Extract request ID and user ID for logging
	requestID := getRequestID(c)
	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Error().Str("request_id", requestID).Msg("User ID not found in context")
		c.JSON(http.StatusUnauthorized, NewErrorResponse("Unauthorized", "", requestID))
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		h.logger.Error().Str("request_id", requestID).Msg("Invalid user ID type in context")
		c.JSON(http.StatusUnauthorized, NewErrorResponse("Unauthorized", "", requestID))
		return
	}

	// Parse debt list ID from URL parameter
	debtListIDStr := c.Param("id")
	debtListID, err := uuid.Parse(debtListIDStr)
	if err != nil {
		h.logger.Warn().Str("request_id", requestID).Str("debt_list_id", debtListIDStr).Msg("Invalid debt list ID format")
		c.JSON(http.StatusBadRequest, NewErrorResponse("Invalid debt list ID", "", requestID))
		return
	}

	logger := h.logger.With().Str("request_id", requestID).Str("user_id", userUUID.String()).Str("debt_list_id", debtListID.String()).Str("method", "DeleteDebtList").Logger()

	logger.Info().Msg("Debt list deletion attempt")

	err = h.debtService.DeleteDebtList(ctx, debtListID, userUUID)
	if err != nil {
		logger.Error().Err(err).Msg("Debt list deletion failed")

		// Handle specific error types
		switch err {
		case entities.ErrDebtListNotFound:
			c.JSON(http.StatusNotFound, NewErrorResponse("Debt list not found", "", requestID))
		default:
			c.JSON(http.StatusInternalServerError, NewErrorResponse("Internal server error", "", requestID))
		}
		return
	}

	logger.Info().Msg("Debt list deleted successfully")

	c.JSON(http.StatusOK, NewSuccessResponse("Debt list deleted successfully", nil, requestID))
}

// CreateDebtItem handles debt item (payment) creation
func (h *DebtHandler) CreateDebtItem(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	// Extract request ID and user ID for logging
	requestID := getRequestID(c)
	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Error().Str("request_id", requestID).Msg("User ID not found in context")
		c.JSON(http.StatusUnauthorized, NewErrorResponse("Unauthorized", "", requestID))
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		h.logger.Error().Str("request_id", requestID).Msg("Invalid user ID type in context")
		c.JSON(http.StatusUnauthorized, NewErrorResponse("Unauthorized", "", requestID))
		return
	}

	logger := h.logger.With().Str("request_id", requestID).Str("user_id", userUUID.String()).Str("method", "CreateDebtItem").Logger()

	var req entities.CreateDebtItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn().Err(err).Msg("Invalid request body")
		c.JSON(http.StatusBadRequest, NewErrorResponse("Invalid request body", err.Error(), requestID))
		return
	}

	// Sanitize input
	req.Amount = sanitizeString(req.Amount)
	req.Currency = sanitizeString(req.Currency)
	req.PaymentMethod = sanitizeString(req.PaymentMethod)
	if req.Description != nil {
		sanitized := sanitizeString(*req.Description)
		req.Description = &sanitized
	}

	logger.Info().Str("debt_list_id", req.DebtListID.String()).Str("amount", req.Amount).Msg("Debt item creation attempt")

	debtItem, err := h.debtService.CreateDebtItem(ctx, userUUID, &req)
	if err != nil {
		logger.Error().Err(err).Str("debt_list_id", req.DebtListID.String()).Msg("Debt item creation failed")

		// Handle specific error types
		switch err {
		case entities.ErrInvalidAmount, entities.ErrInvalidCurrency, entities.ErrInvalidPaymentMethod:
			c.JSON(http.StatusBadRequest, NewErrorResponse("Invalid input", err.Error(), requestID))
		case entities.ErrDebtListNotFound:
			c.JSON(http.StatusNotFound, NewErrorResponse("Debt list not found", "", requestID))
		default:
			c.JSON(http.StatusInternalServerError, NewErrorResponse("Internal server error", "", requestID))
		}
		return
	}

	logger.Info().Str("debt_item_id", debtItem.ID.String()).Str("amount", req.Amount).Msg("Debt item created successfully")

	c.JSON(http.StatusCreated, NewSuccessResponse("Payment recorded successfully", debtItem, requestID))
}

// GetDebtListItems handles retrieving all debt items for a debt list
func (h *DebtHandler) GetDebtListItems(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	// Extract request ID and user ID for logging
	requestID := getRequestID(c)
	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Error().Str("request_id", requestID).Msg("User ID not found in context")
		c.JSON(http.StatusUnauthorized, NewErrorResponse("Unauthorized", "", requestID))
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		h.logger.Error().Str("request_id", requestID).Msg("Invalid user ID type in context")
		c.JSON(http.StatusUnauthorized, NewErrorResponse("Unauthorized", "", requestID))
		return
	}

	// Parse debt list ID from URL parameter
	debtListIDStr := c.Param("id")
	debtListID, err := uuid.Parse(debtListIDStr)
	if err != nil {
		h.logger.Warn().Str("request_id", requestID).Str("debt_list_id", debtListIDStr).Msg("Invalid debt list ID format")
		c.JSON(http.StatusBadRequest, NewErrorResponse("Invalid debt list ID", "", requestID))
		return
	}

	logger := h.logger.With().Str("request_id", requestID).Str("user_id", userUUID.String()).Str("debt_list_id", debtListID.String()).Str("method", "GetDebtListItems").Logger()

	logger.Info().Msg("Retrieving debt list items")

	debtItems, err := h.debtService.GetDebtListItems(ctx, debtListID, userUUID)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to retrieve debt list items")

		// Handle specific error types
		switch err {
		case entities.ErrDebtListNotFound:
			c.JSON(http.StatusNotFound, NewErrorResponse("Debt list not found", "", requestID))
		default:
			c.JSON(http.StatusInternalServerError, NewErrorResponse("Internal server error", "", requestID))
		}
		return
	}

	logger.Info().Int("count", len(debtItems)).Msg("Debt list items retrieved successfully")

	c.JSON(http.StatusOK, NewSuccessResponse("Payments retrieved successfully", debtItems, requestID))
}

// GetOverdueItems handles retrieving overdue debt lists
func (h *DebtHandler) GetOverdueItems(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	// Extract request ID and user ID for logging
	requestID := getRequestID(c)
	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Error().Str("request_id", requestID).Msg("User ID not found in context")
		c.JSON(http.StatusUnauthorized, NewErrorResponse("Unauthorized", "", requestID))
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		h.logger.Error().Str("request_id", requestID).Msg("Invalid user ID type in context")
		c.JSON(http.StatusUnauthorized, NewErrorResponse("Unauthorized", "", requestID))
		return
	}

	logger := h.logger.With().Str("request_id", requestID).Str("user_id", userUUID.String()).Str("method", "GetOverdueItems").Logger()

	logger.Info().Msg("Retrieving overdue items")

	overdueItems, err := h.debtService.GetOverdueItems(ctx, userUUID)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to retrieve overdue items")
		c.JSON(http.StatusInternalServerError, NewErrorResponse("Internal server error", "", requestID))
		return
	}

	logger.Info().Int("count", len(overdueItems)).Msg("Overdue items retrieved successfully")

	c.JSON(http.StatusOK, NewSuccessResponse("Overdue items retrieved successfully", overdueItems, requestID))
}

// GetDueSoonItems handles retrieving debt lists due soon
func (h *DebtHandler) GetDueSoonItems(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	// Extract request ID and user ID for logging
	requestID := getRequestID(c)
	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Error().Str("request_id", requestID).Msg("User ID not found in context")
		c.JSON(http.StatusUnauthorized, NewErrorResponse("Unauthorized", "", requestID))
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		h.logger.Error().Str("request_id", requestID).Msg("Invalid user ID type in context")
		c.JSON(http.StatusUnauthorized, NewErrorResponse("Unauthorized", "", requestID))
		return
	}

	// Get days parameter from query (default to 7 days)
	daysStr := c.DefaultQuery("days", "7")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days < 1 {
		days = 7
	}

	logger := h.logger.With().Str("request_id", requestID).Str("user_id", userUUID.String()).Int("days", days).Str("method", "GetDueSoonItems").Logger()

	logger.Info().Msg("Retrieving due soon items")

	dueSoonItems, err := h.debtService.GetDueSoonItems(ctx, userUUID, days)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to retrieve due soon items")
		c.JSON(http.StatusInternalServerError, NewErrorResponse("Internal server error", "", requestID))
		return
	}

	logger.Info().Int("count", len(dueSoonItems)).Msg("Due soon items retrieved successfully")

	c.JSON(http.StatusOK, NewSuccessResponse("Due soon items retrieved successfully", dueSoonItems, requestID))
}

// GetPaymentSchedule handles retrieving the payment schedule for a debt list
func (h *DebtHandler) GetPaymentSchedule(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	// Extract request ID and user ID for logging
	requestID := getRequestID(c)
	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Error().Str("request_id", requestID).Msg("User ID not found in context")
		c.JSON(http.StatusUnauthorized, NewErrorResponse("Unauthorized", "", requestID))
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		h.logger.Error().Str("request_id", requestID).Msg("Invalid user ID type in context")
		c.JSON(http.StatusUnauthorized, NewErrorResponse("Unauthorized", "", requestID))
		return
	}

	// Parse debt list ID from URL parameter
	debtListIDStr := c.Param("id")
	debtListID, err := uuid.Parse(debtListIDStr)
	if err != nil {
		h.logger.Warn().Str("request_id", requestID).Str("debt_list_id", debtListIDStr).Msg("Invalid debt list ID format")
		c.JSON(http.StatusBadRequest, NewErrorResponse("Invalid debt list ID", "", requestID))
		return
	}

	logger := h.logger.With().Str("request_id", requestID).Str("user_id", userUUID.String()).Str("debt_list_id", debtListID.String()).Str("method", "GetPaymentSchedule").Logger()

	logger.Info().Msg("Retrieving payment schedule")

	schedule, err := h.debtService.GetPaymentSchedule(ctx, debtListID, userUUID)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to retrieve payment schedule")

		// Handle specific error types
		switch err {
		case entities.ErrDebtListNotFound:
			c.JSON(http.StatusNotFound, NewErrorResponse("Debt list not found", "", requestID))
		default:
			c.JSON(http.StatusInternalServerError, NewErrorResponse("Internal server error", "", requestID))
		}
		return
	}

	logger.Info().Int("schedule_items", len(schedule)).Msg("Payment schedule retrieved successfully")

	c.JSON(http.StatusOK, NewSuccessResponse("Payment schedule retrieved successfully", schedule, requestID))
}

// GetUpcomingPayments handles retrieving upcoming payments
func (h *DebtHandler) GetUpcomingPayments(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	// Extract request ID and user ID for logging
	requestID := getRequestID(c)
	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Error().Str("request_id", requestID).Msg("User ID not found in context")
		c.JSON(http.StatusUnauthorized, NewErrorResponse("Unauthorized", "", requestID))
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		h.logger.Error().Str("request_id", requestID).Msg("Invalid user ID type in context")
		c.JSON(http.StatusUnauthorized, NewErrorResponse("Unauthorized", "", requestID))
		return
	}

	// Get days parameter from query (default to 30 days)
	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days < 1 {
		days = 30
	}

	logger := h.logger.With().Str("request_id", requestID).Str("user_id", userUUID.String()).Int("days", days).Str("method", "GetUpcomingPayments").Logger()

	logger.Info().Msg("Retrieving upcoming payments")

	upcomingPayments, err := h.debtService.GetUpcomingPayments(ctx, userUUID, days)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to retrieve upcoming payments")
		c.JSON(http.StatusInternalServerError, NewErrorResponse("Internal server error", "", requestID))
		return
	}

	logger.Info().Int("count", len(upcomingPayments)).Msg("Upcoming payments retrieved successfully")

	c.JSON(http.StatusOK, NewSuccessResponse("Upcoming payments retrieved successfully", upcomingPayments, requestID))
}

// GetTotalPaymentsForDebtList handles retrieving payment summary for a debt list
func (h *DebtHandler) GetTotalPaymentsForDebtList(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	// Extract request ID and user ID for logging
	requestID := getRequestID(c)
	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Error().Str("request_id", requestID).Msg("User ID not found in context")
		c.JSON(http.StatusUnauthorized, NewErrorResponse("Unauthorized", "", requestID))
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		h.logger.Error().Str("request_id", requestID).Msg("Invalid user ID type in context")
		c.JSON(http.StatusUnauthorized, NewErrorResponse("Unauthorized", "", requestID))
		return
	}

	// Parse debt list ID from URL parameter
	debtListIDStr := c.Param("id")
	debtListID, err := uuid.Parse(debtListIDStr)
	if err != nil {
		h.logger.Warn().Str("request_id", requestID).Str("debt_list_id", debtListIDStr).Msg("Invalid debt list ID format")
		c.JSON(http.StatusBadRequest, NewErrorResponse("Invalid debt list ID", "", requestID))
		return
	}

	logger := h.logger.With().Str("request_id", requestID).Str("user_id", userUUID.String()).Str("debt_list_id", debtListID.String()).Str("method", "GetTotalPaymentsForDebtList").Logger()

	logger.Info().Msg("Retrieving payment summary")

	summary, err := h.debtService.GetTotalPaymentsForDebtList(ctx, debtListID, userUUID)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to retrieve payment summary")

		// Handle specific error types
		switch err {
		case entities.ErrDebtListNotFound:
			c.JSON(http.StatusNotFound, NewErrorResponse("Debt list not found", "", requestID))
		default:
			c.JSON(http.StatusInternalServerError, NewErrorResponse("Internal server error", "", requestID))
		}
		return
	}

	logger.Info().Msg("Payment summary retrieved successfully")

	c.JSON(http.StatusOK, NewSuccessResponse("Payment summary retrieved successfully", summary, requestID))
}

// VerifyDebtItem handles debt item verification
func (h *DebtHandler) VerifyDebtItem(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	// Extract request ID and user ID for logging
	requestID := getRequestID(c)
	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Error().Str("request_id", requestID).Msg("User ID not found in context")
		c.JSON(http.StatusUnauthorized, NewErrorResponse("Unauthorized", "", requestID))
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		h.logger.Error().Str("request_id", requestID).Msg("Invalid user ID type in context")
		c.JSON(http.StatusUnauthorized, NewErrorResponse("Unauthorized", "", requestID))
		return
	}

	// Parse debt item ID from URL parameter
	debtItemIDStr := c.Param("id")
	debtItemID, err := uuid.Parse(debtItemIDStr)
	if err != nil {
		h.logger.Warn().Str("request_id", requestID).Str("debt_item_id", debtItemIDStr).Msg("Invalid debt item ID format")
		c.JSON(http.StatusBadRequest, NewErrorResponse("Invalid debt item ID", "", requestID))
		return
	}

	var req entities.VerifyDebtItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn().Str("request_id", requestID).Err(err).Msg("Invalid request body")
		c.JSON(http.StatusBadRequest, NewErrorResponse("Invalid request body", err.Error(), requestID))
		return
	}

	// Sanitize input
	if req.VerificationNotes != nil {
		sanitized := sanitizeString(*req.VerificationNotes)
		req.VerificationNotes = &sanitized
	}

	logger := h.logger.With().Str("request_id", requestID).Str("user_id", userUUID.String()).Str("debt_item_id", debtItemID.String()).Str("method", "VerifyDebtItem").Logger()

	logger.Info().Str("status", req.Status).Msg("Debt item verification attempt")

	debtItem, err := h.debtService.VerifyDebtItem(ctx, debtItemID, userUUID, &req)
	if err != nil {
		logger.Error().Err(err).Str("status", req.Status).Msg("Debt item verification failed")

		// Handle specific error types
		switch err {
		case entities.ErrDebtItemNotFound:
			c.JSON(http.StatusNotFound, NewErrorResponse("Debt item not found", "", requestID))
		case entities.ErrInvalidPaymentStatus:
			c.JSON(http.StatusBadRequest, NewErrorResponse("Invalid payment status", err.Error(), requestID))
		default:
			c.JSON(http.StatusInternalServerError, NewErrorResponse("Internal server error", "", requestID))
		}
		return
	}

	logger.Info().Str("status", req.Status).Msg("Debt item verified successfully")

	c.JSON(http.StatusOK, NewSuccessResponse("Debt item verified successfully", debtItem, requestID))
}

// GetPendingVerifications handles retrieving pending verifications for a user
func (h *DebtHandler) GetPendingVerifications(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	// Extract request ID and user ID for logging
	requestID := getRequestID(c)
	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Error().Str("request_id", requestID).Msg("User ID not found in context")
		c.JSON(http.StatusUnauthorized, NewErrorResponse("Unauthorized", "", requestID))
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		h.logger.Error().Str("request_id", requestID).Msg("Invalid user ID type in context")
		c.JSON(http.StatusUnauthorized, NewErrorResponse("Unauthorized", "", requestID))
		return
	}

	logger := h.logger.With().Str("request_id", requestID).Str("user_id", userUUID.String()).Str("method", "GetPendingVerifications").Logger()

	logger.Info().Msg("Retrieving pending verifications")

	pendingVerifications, err := h.debtService.GetPendingVerifications(ctx, userUUID)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to retrieve pending verifications")
		c.JSON(http.StatusInternalServerError, NewErrorResponse("Internal server error", "", requestID))
		return
	}

	logger.Info().Int("count", len(pendingVerifications)).Msg("Pending verifications retrieved successfully")

	c.JSON(http.StatusOK, NewSuccessResponse("Pending verifications retrieved successfully", pendingVerifications, requestID))
}

// RejectDebtItem handles debt item rejection
func (h *DebtHandler) RejectDebtItem(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	// Extract request ID and user ID for logging
	requestID := getRequestID(c)
	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Error().Str("request_id", requestID).Msg("User ID not found in context")
		c.JSON(http.StatusUnauthorized, NewErrorResponse("Unauthorized", "", requestID))
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		h.logger.Error().Str("request_id", requestID).Msg("Invalid user ID type in context")
		c.JSON(http.StatusUnauthorized, NewErrorResponse("Unauthorized", "", requestID))
		return
	}

	// Parse debt item ID from URL parameter
	debtItemIDStr := c.Param("id")
	debtItemID, err := uuid.Parse(debtItemIDStr)
	if err != nil {
		h.logger.Warn().Str("request_id", requestID).Str("debt_item_id", debtItemIDStr).Msg("Invalid debt item ID format")
		c.JSON(http.StatusBadRequest, NewErrorResponse("Invalid debt item ID", "", requestID))
		return
	}

	var req struct {
		Notes *string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn().Str("request_id", requestID).Err(err).Msg("Invalid request body")
		c.JSON(http.StatusBadRequest, NewErrorResponse("Invalid request body", err.Error(), requestID))
		return
	}

	// Sanitize input
	if req.Notes != nil {
		sanitized := sanitizeString(*req.Notes)
		req.Notes = &sanitized
	}

	logger := h.logger.With().Str("request_id", requestID).Str("user_id", userUUID.String()).Str("debt_item_id", debtItemID.String()).Str("method", "RejectDebtItem").Logger()

	logger.Info().Msg("Debt item rejection attempt")

	debtItem, err := h.debtService.RejectDebtItem(ctx, debtItemID, userUUID, req.Notes)
	if err != nil {
		logger.Error().Err(err).Msg("Debt item rejection failed")

		// Handle specific error types
		switch err {
		case entities.ErrDebtItemNotFound:
			c.JSON(http.StatusNotFound, NewErrorResponse("Debt item not found", "", requestID))
		default:
			c.JSON(http.StatusInternalServerError, NewErrorResponse("Internal server error", "", requestID))
		}
		return
	}

	logger.Info().Msg("Debt item rejected successfully")

	c.JSON(http.StatusOK, NewSuccessResponse("Debt item rejected successfully", debtItem, requestID))
}

// UploadReceipt handles receipt photo upload for a debt item
func (h *DebtHandler) UploadReceipt(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 60*time.Second) // Longer timeout for file uploads
	defer cancel()

	// Extract request ID and user ID for logging
	requestID := getRequestID(c)
	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Error().Str("request_id", requestID).Msg("User ID not found in context")
		c.JSON(http.StatusUnauthorized, NewErrorResponse("Unauthorized", "", requestID))
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		h.logger.Error().Str("request_id", requestID).Msg("Invalid user ID type in context")
		c.JSON(http.StatusUnauthorized, NewErrorResponse("Unauthorized", "", requestID))
		return
	}

	// Parse debt item ID from URL parameter
	debtItemIDStr := c.Param("id")
	debtItemID, err := uuid.Parse(debtItemIDStr)
	if err != nil {
		h.logger.Warn().Str("request_id", requestID).Str("debt_item_id", debtItemIDStr).Msg("Invalid debt item ID format")
		c.JSON(http.StatusBadRequest, NewErrorResponse("Invalid debt item ID", "", requestID))
		return
	}

	// Parse multipart form
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil { // 32MB max
		h.logger.Warn().Str("request_id", requestID).Err(err).Msg("Failed to parse multipart form")
		c.JSON(http.StatusBadRequest, NewErrorResponse("Failed to parse form data", "", requestID))
		return
	}

	// Get the uploaded file
	file, header, err := c.Request.FormFile("receipt")
	if err != nil {
		h.logger.Warn().Str("request_id", requestID).Err(err).Msg("No receipt file provided")
		c.JSON(http.StatusBadRequest, NewErrorResponse("Receipt file is required", "", requestID))
		return
	}
	defer file.Close()

	logger := h.logger.With().Str("request_id", requestID).Str("user_id", userUUID.String()).Str("debt_item_id", debtItemID.String()).Str("method", "UploadReceipt").Logger()

	// Validate file
	if err := h.validateReceiptFile(header); err != nil {
		logger.Warn().Err(err).Str("filename", header.Filename).Msg("Invalid receipt file")
		c.JSON(http.StatusBadRequest, NewErrorResponse("Invalid receipt file", err.Error(), requestID))
		return
	}

	// Upload file to S3
	photoURL, err := h.fileStorageService.UploadReceipt(ctx, file, header.Filename, header.Header.Get("Content-Type"), debtItemID)
	if err != nil {
		logger.Error().Err(err).Str("filename", header.Filename).Msg("Failed to upload receipt")
		c.JSON(http.StatusInternalServerError, NewErrorResponse("Failed to upload receipt", "", requestID))
		return
	}

	// Update debt item with receipt photo URL
	updateReq := &entities.UpdateDebtItemRequest{
		ReceiptPhotoURL: &photoURL,
	}

	debtItem, err := h.debtService.UpdateDebtItem(ctx, debtItemID, userUUID, updateReq)
	if err != nil {
		logger.Error().Err(err).Str("photo_url", photoURL).Msg("Failed to update debt item with receipt photo")
		c.JSON(http.StatusInternalServerError, NewErrorResponse("Failed to update debt item", "", requestID))
		return
	}

	logger.Info().Str("filename", header.Filename).Str("photo_url", photoURL).Msg("Receipt uploaded successfully")

	c.JSON(http.StatusOK, NewSuccessResponse("Receipt uploaded successfully", debtItem, requestID))
}

// GetReceiptPhoto serves a receipt photo from S3
func (h *DebtHandler) GetReceiptPhoto(c *gin.Context) {
	requestID := c.GetString("request_id")
	if requestID == "" {
		requestID = uuid.New().String()
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Warn().Str("request_id", requestID).Msg("User ID not found in context")
		c.JSON(http.StatusUnauthorized, NewErrorResponse("Unauthorized", "", requestID))
		return
	}
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		h.logger.Warn().Str("request_id", requestID).Msg("Invalid user ID type in context")
		c.JSON(http.StatusUnauthorized, NewErrorResponse("Unauthorized", "", requestID))
		return
	}

	// Get debt ID and filename from URL parameters
	debtIDStr := c.Param("id")
	debtID, err := uuid.Parse(debtIDStr)
	if err != nil {
		h.logger.Warn().Str("request_id", requestID).Str("debt_id", debtIDStr).Msg("Invalid debt ID format")
		c.JSON(http.StatusBadRequest, NewErrorResponse("Invalid debt ID", "", requestID))
		return
	}

	filename := c.Param("filename")
	if filename == "" {
		h.logger.Warn().Str("request_id", requestID).Msg("Filename not provided")
		c.JSON(http.StatusBadRequest, NewErrorResponse("Filename is required", "", requestID))
		return
	}

	// Construct the full API path
	fullPath := fmt.Sprintf("/api/v1/debts/%s/receipts/%s", debtID.String(), filename)

	logger := h.logger.With().Str("request_id", requestID).Str("user_id", userUUID.String()).Str("debt_id", debtID.String()).Str("filename", filename).Str("method", "GetReceiptPhoto").Logger()

	// Get the file from S3
	fileContent, contentType, err := h.fileStorageService.GetReceiptFile(c.Request.Context(), fullPath)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to retrieve receipt photo")
		c.JSON(http.StatusNotFound, NewErrorResponse("Receipt photo not found", "", requestID))
		return
	}

	// Set appropriate headers for file serving
	c.Header("Content-Type", contentType)
	c.Header("Content-Length", fmt.Sprintf("%d", len(fileContent)))
	c.Header("Cache-Control", "public, max-age=3600") // Cache for 1 hour
	c.Header("Content-Disposition", "inline")

	// Serve the file
	c.Data(http.StatusOK, contentType, fileContent)

	logger.Info().Str("content_type", contentType).Int("size", len(fileContent)).Msg("Receipt photo served successfully")
}

// validateReceiptFile validates the uploaded receipt file
func (h *DebtHandler) validateReceiptFile(header *multipart.FileHeader) error {
	// Check file size (max 10MB)
	if header.Size > 10*1024*1024 {
		return fmt.Errorf("file size %d bytes exceeds maximum allowed size of 10MB", header.Size)
	}

	// Check file type
	contentType := header.Header.Get("Content-Type")
	if !h.isValidImageType(contentType) {
		return fmt.Errorf("invalid file type: %s. Only images (JPEG, PNG, GIF, WebP) are allowed", contentType)
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	validExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}
	if !validExtensions[ext] {
		return fmt.Errorf("invalid file extension: %s. Only .jpg, .jpeg, .png, .gif, .webp are allowed", ext)
	}

	return nil
}

// isValidImageType checks if the content type is a valid image type
func (h *DebtHandler) isValidImageType(contentType string) bool {
	validTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}
	return validTypes[contentType]
}






