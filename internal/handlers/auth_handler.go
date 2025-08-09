package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"exit-debt/internal/domain/entities"
	"exit-debt/internal/domain/interfaces"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	authService interfaces.AuthService
	logger      zerolog.Logger
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService interfaces.AuthService, logger zerolog.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger.With().Str("handler", "auth").Logger(),
	}
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	// Extract request ID for logging
	requestID := getRequestID(c)
	logger := h.logger.With().Str("request_id", requestID).Str("method", "Register").Logger()

	var req entities.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn().Err(err).Msg("Invalid request body")
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:     "Invalid request body",
			Details:   err.Error(),
			RequestID: requestID,
		})
		return
	}

	// Sanitize input
	req.Email = sanitizeEmail(req.Email)
	req.FirstName = sanitizeString(req.FirstName)
	req.LastName = sanitizeString(req.LastName)
	if req.Phone != nil {
		sanitized := sanitizeString(*req.Phone)
		req.Phone = &sanitized
	}

	logger.Info().Str("email", req.Email).Msg("User registration attempt")

	response, err := h.authService.Register(ctx, &req)
	if err != nil {
		logger.Error().Err(err).Str("email", req.Email).Msg("User registration failed")
		
		// Handle specific error types
		switch err {
		case entities.ErrUserAlreadyExists:
			c.JSON(http.StatusConflict, ErrorResponse{
				Error:     "User already exists",
				RequestID: requestID,
			})
		case entities.ErrInvalidEmail, entities.ErrInvalidPassword, entities.ErrInvalidFirstName, entities.ErrInvalidLastName:
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:     "Invalid input",
				Details:   err.Error(),
				RequestID: requestID,
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:     "Internal server error",
				RequestID: requestID,
			})
		}
		return
	}

	logger.Info().Str("email", req.Email).Str("user_id", response.User.ID.String()).Msg("User registered successfully")

	c.JSON(http.StatusCreated, SuccessResponse{
		Message:   "User registered successfully",
		Data:      response,
		RequestID: requestID,
	})
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	// Extract request ID for logging
	requestID := getRequestID(c)
	logger := h.logger.With().Str("request_id", requestID).Str("method", "Login").Logger()

	var req entities.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn().Err(err).Msg("Invalid request body")
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:     "Invalid request body",
			Details:   err.Error(),
			RequestID: requestID,
		})
		return
	}

	// Sanitize input
	req.Email = sanitizeEmail(req.Email)

	logger.Info().Str("email", req.Email).Msg("User login attempt")

	response, err := h.authService.Login(ctx, &req)
	if err != nil {
		logger.Warn().Err(err).Str("email", req.Email).Msg("User login failed")
		
		// Handle specific error types
		switch err {
		case entities.ErrInvalidCredentials:
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:     "Invalid credentials",
				RequestID: requestID,
			})
		case entities.ErrInvalidEmail, entities.ErrInvalidPassword:
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:     "Invalid input",
				Details:   err.Error(),
				RequestID: requestID,
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:     "Internal server error",
				RequestID: requestID,
			})
		}
		return
	}

	logger.Info().Str("email", req.Email).Str("user_id", response.User.ID.String()).Msg("User logged in successfully")

	c.JSON(http.StatusOK, SuccessResponse{
		Message:   "Login successful",
		Data:      response,
		RequestID: requestID,
	})
}

// ValidateToken validates a JWT token (used by middleware)
func (h *AuthHandler) ValidateToken(ctx context.Context, tokenString string) (string, error) {
	userID, err := h.authService.ValidateToken(ctx, tokenString)
	if err != nil {
		return "", err
	}
	return userID.String(), nil
}
