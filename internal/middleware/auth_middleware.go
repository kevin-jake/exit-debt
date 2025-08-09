package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"exit-debt/internal/domain/entities"
	"exit-debt/internal/domain/interfaces"
)

// AuthMiddleware provides JWT authentication middleware
type AuthMiddleware struct {
	authService interfaces.AuthService
	logger      zerolog.Logger
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(authService interfaces.AuthService, logger zerolog.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
		logger:      logger.With().Str("middleware", "auth").Logger(),
	}
}

// Authenticate returns a Gin middleware function for JWT authentication
func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()

		// Extract request ID for logging
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
			c.Header("X-Request-ID", requestID)
		}

		logger := m.logger.With().Str("request_id", requestID).Logger()

		// Get authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Warn().Str("path", c.Request.URL.Path).Msg("Missing authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":      "Authorization header required",
				"request_id": requestID,
				"timestamp":  time.Now(),
			})
			c.Abort()
			return
		}

		// Check if the header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			logger.Warn().Str("path", c.Request.URL.Path).Msg("Invalid authorization header format")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":      "Invalid authorization header format. Use: Bearer <token>",
				"request_id": requestID,
				"timestamp":  time.Now(),
			})
			c.Abort()
			return
		}

		// Extract the token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			logger.Warn().Str("path", c.Request.URL.Path).Msg("Empty token in authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":      "Empty token in authorization header",
				"request_id": requestID,
				"timestamp":  time.Now(),
			})
			c.Abort()
			return
		}

		// Validate the token
		userID, err := m.authService.ValidateToken(ctx, token)
		if err != nil {
			logger.Warn().Err(err).Str("path", c.Request.URL.Path).Msg("Token validation failed")
			
			// Handle specific error types
			var errorMessage string
			switch err {
			case entities.ErrInvalidToken:
				errorMessage = "Invalid token"
			case entities.ErrTokenExpired:
				errorMessage = "Token expired"
			default:
				errorMessage = "Token validation failed"
			}
			
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":      errorMessage,
				"request_id": requestID,
				"timestamp":  time.Now(),
			})
			c.Abort()
			return
		}

		// Set user ID in context for downstream handlers
		c.Set("user_id", userID)
		c.Set("request_id", requestID)

		logger.Debug().Str("user_id", userID.String()).Str("path", c.Request.URL.Path).Msg("Authentication successful")

		// Continue to the next handler
		c.Next()
	}
}

// OptionalAuthenticate returns a Gin middleware function for optional JWT authentication
// This middleware will set user_id in context if a valid token is provided, but won't block requests without tokens
func (m *AuthMiddleware) OptionalAuthenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()

		// Extract request ID for logging
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
			c.Header("X-Request-ID", requestID)
		}
		c.Set("request_id", requestID)

		logger := m.logger.With().Str("request_id", requestID).Logger()

		// Get authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// No token provided, but that's okay for optional auth
			logger.Debug().Str("path", c.Request.URL.Path).Msg("No authorization header provided for optional auth")
			c.Next()
			return
		}

		// Check if the header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			logger.Debug().Str("path", c.Request.URL.Path).Msg("Invalid authorization header format for optional auth")
			c.Next()
			return
		}

		// Extract the token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			logger.Debug().Str("path", c.Request.URL.Path).Msg("Empty token in authorization header for optional auth")
			c.Next()
			return
		}

		// Validate the token
		userID, err := m.authService.ValidateToken(ctx, token)
		if err != nil {
			logger.Debug().Err(err).Str("path", c.Request.URL.Path).Msg("Token validation failed for optional auth")
			// Don't abort, just continue without setting user_id
			c.Next()
			return
		}

		// Set user ID in context for downstream handlers
		c.Set("user_id", userID)
		logger.Debug().Str("user_id", userID.String()).Str("path", c.Request.URL.Path).Msg("Optional authentication successful")

		// Continue to the next handler
		c.Next()
	}
}
