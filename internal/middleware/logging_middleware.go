package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// LoggingMiddleware provides structured logging for HTTP requests
type LoggingMiddleware struct {
	logger zerolog.Logger
}

// NewLoggingMiddleware creates a new logging middleware
func NewLoggingMiddleware(logger zerolog.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{
		logger: logger.With().Str("middleware", "logging").Logger(),
	}
}

// LogRequests returns a Gin middleware function for request logging
func (m *LoggingMiddleware) LogRequests() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Generate request ID if not present
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
			c.Header("X-Request-ID", requestID)
		}

		// Extract user ID if available
		userID := ""
		if uid, exists := c.Get("user_id"); exists {
			if uuidVal, ok := uid.(uuid.UUID); ok {
				userID = uuidVal.String()
			}
		}

		// Create logger with request context
		logger := m.logger.With().
			Str("request_id", requestID).
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Str("query", c.Request.URL.RawQuery).
			Str("client_ip", c.ClientIP()).
			Str("user_agent", c.Request.UserAgent()).
			Str("user_id", userID).
			Logger()

		// Log incoming request
		logger.Info().Msg("Request started")

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get response size
		responseSize := c.Writer.Size()

		// Log completed request
		logEvent := logger.With().
			Int("status", c.Writer.Status()).
			Dur("latency", latency).
			Int("response_size", responseSize).
			Logger()

		// Determine log level based on status code
		statusCode := c.Writer.Status()
		switch {
		case statusCode >= 500:
			logEvent.Error().Msg("Request completed with server error")
		case statusCode >= 400:
			logEvent.Warn().Msg("Request completed with client error")
		case statusCode >= 300:
			logEvent.Info().Msg("Request completed with redirect")
		default:
			logEvent.Info().Msg("Request completed successfully")
		}

		// Log any errors that occurred during request processing
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				logger.Error().Err(err.Err).Int("type", int(err.Type)).Msg("Request error")
			}
		}
	}
}

// CORS returns a Gin middleware function for handling CORS
func (m *LoggingMiddleware) CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// Set CORS headers
		c.Header("Access-Control-Allow-Origin", "*") // In production, should be more restrictive
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Request-ID")
		c.Header("Access-Control-Expose-Headers", "X-Request-ID")
		c.Header("Access-Control-Max-Age", "86400") // Cache preflight for 24 hours

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			m.logger.Debug().
				Str("origin", origin).
				Str("method", c.Request.Header.Get("Access-Control-Request-Method")).
				Str("headers", c.Request.Header.Get("Access-Control-Request-Headers")).
				Msg("CORS preflight request")
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// Recovery returns a Gin middleware function for panic recovery
func (m *LoggingMiddleware) Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		m.logger.Error().
			Interface("panic", recovered).
			Str("request_id", requestID).
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Str("client_ip", c.ClientIP()).
			Msg("Panic recovered")

		c.JSON(500, gin.H{
			"error":      "Internal server error",
			"request_id": requestID,
			"timestamp":  time.Now(),
		})
	})
}
