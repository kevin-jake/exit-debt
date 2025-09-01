package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"exit-debt/internal/domain/entities"
	"exit-debt/internal/domain/interfaces"
)

// ContactHandler handles contact-related HTTP requests
type ContactHandler struct {
	contactService interfaces.ContactService
	logger         zerolog.Logger
}

// NewContactHandler creates a new contact handler
func NewContactHandler(contactService interfaces.ContactService, logger zerolog.Logger) *ContactHandler {
	return &ContactHandler{
		contactService: contactService,
		logger:         logger.With().Str("handler", "contact").Logger(),
	}
}

// CreateContact handles contact creation
func (h *ContactHandler) CreateContact(c *gin.Context) {
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

	logger := h.logger.With().Str("request_id", requestID).Str("user_id", userUUID.String()).Str("method", "CreateContact").Logger()

	var req entities.CreateContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn().Err(err).Msg("Invalid request body")
		c.JSON(http.StatusBadRequest, NewErrorResponse("Invalid request body", err.Error(), requestID))
		return
	}

	// Sanitize input
	req.Name = sanitizeString(req.Name)
	if req.Email != nil {
		sanitized := sanitizeEmail(*req.Email)
		req.Email = &sanitized
	}
	if req.Phone != nil {
		sanitized := sanitizeString(*req.Phone)
		req.Phone = &sanitized
	}
	if req.Notes != nil {
		sanitized := sanitizeString(*req.Notes)
		req.Notes = &sanitized
	}

	logger.Info().Str("contact_name", req.Name).Msg("Contact creation attempt")

	contact, err := h.contactService.CreateContact(ctx, userUUID, &req)
	if err != nil {
		logger.Error().Err(err).Str("contact_name", req.Name).Msg("Contact creation failed")
		
		// Handle specific error types
		switch err {
		case entities.ErrContactAlreadyExists:
			c.JSON(http.StatusConflict, NewErrorResponse("Contact already exists", "", requestID))
		case entities.ErrContactPhoneExists:
			c.JSON(http.StatusConflict, NewErrorResponse("Contact with this phone number already exists", "", requestID))
		case entities.ErrInvalidContactName:
			c.JSON(http.StatusBadRequest, NewErrorResponse("Invalid input", err.Error(), requestID))
		default:
			c.JSON(http.StatusInternalServerError, NewErrorResponse("Internal server error", "", requestID))
		}
		return
	}

	logger.Info().Str("contact_name", req.Name).Str("contact_id", contact.ID.String()).Msg("Contact created successfully")

	c.JSON(http.StatusCreated, NewSuccessResponse("Contact created successfully", contact, requestID))
}

// GetUserContacts handles retrieving all contacts for a user
func (h *ContactHandler) GetUserContacts(c *gin.Context) {
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

	logger := h.logger.With().Str("request_id", requestID).Str("user_id", userUUID.String()).Str("method", "GetUserContacts").Logger()

	logger.Info().Msg("Retrieving user contacts")

	contacts, err := h.contactService.GetUserContacts(ctx, userUUID)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to retrieve user contacts")
		c.JSON(http.StatusInternalServerError, NewErrorResponse("Internal server error", "", requestID))
		return
	}

	logger.Info().Int("count", len(contacts)).Msg("User contacts retrieved successfully")

	c.JSON(http.StatusOK, NewSuccessResponse("Contacts retrieved successfully", contacts, requestID))
}

// GetContact handles retrieving a specific contact
func (h *ContactHandler) GetContact(c *gin.Context) {
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

	// Parse contact ID from URL parameter
	contactIDStr := c.Param("id")
	contactID, err := uuid.Parse(contactIDStr)
	if err != nil {
		h.logger.Warn().Str("request_id", requestID).Str("contact_id", contactIDStr).Msg("Invalid contact ID format")
		c.JSON(http.StatusBadRequest, NewErrorResponse("Invalid contact ID", "", requestID))
		return
	}

	logger := h.logger.With().Str("request_id", requestID).Str("user_id", userUUID.String()).Str("contact_id", contactID.String()).Str("method", "GetContact").Logger()

	logger.Info().Msg("Retrieving contact")

	contact, err := h.contactService.GetContact(ctx, contactID, userUUID)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to retrieve contact")
		
		// Handle specific error types
		switch err {
		case entities.ErrContactNotFound:
			c.JSON(http.StatusNotFound, NewErrorResponse("Contact not found", "", requestID))
		default:
			c.JSON(http.StatusInternalServerError, NewErrorResponse("Internal server error", "", requestID))
		}
		return
	}

	logger.Info().Msg("Contact retrieved successfully")

	c.JSON(http.StatusOK, NewSuccessResponse("Contact retrieved successfully", contact, requestID))
}

// UpdateContact handles contact updates
func (h *ContactHandler) UpdateContact(c *gin.Context) {
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

	// Parse contact ID from URL parameter
	contactIDStr := c.Param("id")
	contactID, err := uuid.Parse(contactIDStr)
	if err != nil {
		h.logger.Warn().Str("request_id", requestID).Str("contact_id", contactIDStr).Msg("Invalid contact ID format")
		c.JSON(http.StatusBadRequest, NewErrorResponse("Invalid contact ID", "", requestID))
		return
	}

	logger := h.logger.With().Str("request_id", requestID).Str("user_id", userUUID.String()).Str("contact_id", contactID.String()).Str("method", "UpdateContact").Logger()

	var req entities.UpdateContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn().Err(err).Msg("Invalid request body")
		c.JSON(http.StatusBadRequest, NewErrorResponse("Invalid request body", err.Error(), requestID))
		return
	}

	// Sanitize input
	if req.Name != nil {
		sanitized := sanitizeString(*req.Name)
		req.Name = &sanitized
	}
	if req.Email != nil {
		sanitized := sanitizeEmail(*req.Email)
		req.Email = &sanitized
	}
	if req.Phone != nil {
		sanitized := sanitizeString(*req.Phone)
		req.Phone = &sanitized
	}
	if req.Notes != nil {
		sanitized := sanitizeString(*req.Notes)
		req.Notes = &sanitized
	}

	logger.Info().Msg("Contact update attempt")

	contact, err := h.contactService.UpdateContact(ctx, contactID, userUUID, &req)
	if err != nil {
		logger.Error().Err(err).Msg("Contact update failed")
		
		// Handle specific error types
		switch err {
		case entities.ErrContactNotFound:
			c.JSON(http.StatusNotFound, NewErrorResponse("Contact not found", "", requestID))
		case entities.ErrInvalidContactName:
			c.JSON(http.StatusBadRequest, NewErrorResponse("Invalid input", err.Error(), requestID))
		default:
			c.JSON(http.StatusInternalServerError, NewErrorResponse("Internal server error", "", requestID))
		}
		return
	}

	logger.Info().Msg("Contact updated successfully")

	c.JSON(http.StatusOK, NewSuccessResponse("Contact updated successfully", contact, requestID))
}

// DeleteContact handles contact deletion
func (h *ContactHandler) DeleteContact(c *gin.Context) {
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

	// Parse contact ID from URL parameter
	contactIDStr := c.Param("id")
	contactID, err := uuid.Parse(contactIDStr)
	if err != nil {
		h.logger.Warn().Str("request_id", requestID).Str("contact_id", contactIDStr).Msg("Invalid contact ID format")
		c.JSON(http.StatusBadRequest, NewErrorResponse("Invalid contact ID", "", requestID))
		return
	}

	logger := h.logger.With().Str("request_id", requestID).Str("user_id", userUUID.String()).Str("contact_id", contactID.String()).Str("method", "DeleteContact").Logger()

	logger.Info().Msg("Contact deletion attempt")

	err = h.contactService.DeleteContact(ctx, contactID, userUUID)
	if err != nil {
		logger.Error().Err(err).Msg("Contact deletion failed")
		
		// Handle specific error types
		switch err {
		case entities.ErrContactNotFound:
			c.JSON(http.StatusNotFound, NewErrorResponse("Contact not found", "", requestID))
		default:
			c.JSON(http.StatusInternalServerError, NewErrorResponse("Internal server error", "", requestID))
		}
		return
	}

	logger.Info().Msg("Contact deleted successfully")

	c.JSON(http.StatusOK, NewSuccessResponse("Contact deleted successfully", nil, requestID))
}
