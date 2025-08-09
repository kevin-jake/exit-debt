package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"exit-debt/internal/domain/entities"
	"exit-debt/internal/domain/interfaces"
)

// DebtHandler handles debt-related HTTP requests
type DebtHandler struct {
	debtService interfaces.DebtService
	logger      zerolog.Logger
}

// NewDebtHandler creates a new debt handler
func NewDebtHandler(debtService interfaces.DebtService, logger zerolog.Logger) *DebtHandler {
	return &DebtHandler{
		debtService: debtService,
		logger:      logger.With().Str("handler", "debt").Logger(),
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






