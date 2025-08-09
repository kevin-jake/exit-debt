package handlers

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SuccessResponse represents a successful API response
type SuccessResponse struct {
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	RequestID string      `json:"request_id"`
	Timestamp time.Time   `json:"timestamp"`
}

// ErrorResponse represents an error API response
type ErrorResponse struct {
	Error     string `json:"error"`
	Details   string `json:"details,omitempty"`
	RequestID string `json:"request_id"`
	Timestamp time.Time `json:"timestamp"`
}

// NewSuccessResponse creates a new success response
func NewSuccessResponse(message string, data interface{}, requestID string) SuccessResponse {
	return SuccessResponse{
		Message:   message,
		Data:      data,
		RequestID: requestID,
		Timestamp: time.Now(),
	}
}

// NewErrorResponse creates a new error response
func NewErrorResponse(error, details, requestID string) ErrorResponse {
	return ErrorResponse{
		Error:     error,
		Details:   details,
		RequestID: requestID,
		Timestamp: time.Now(),
	}
}

// getRequestID extracts or generates a request ID for tracing
func getRequestID(c *gin.Context) string {
	requestID := c.GetHeader("X-Request-ID")
	if requestID == "" {
		// Generate a new request ID if not provided
		requestID = uuid.New().String()
		c.Header("X-Request-ID", requestID)
	}
	return requestID
}

// sanitizeString removes potentially dangerous characters from strings
func sanitizeString(input string) string {
	// Trim whitespace and remove null bytes
	cleaned := strings.TrimSpace(input)
	cleaned = strings.ReplaceAll(cleaned, "\x00", "")
	return cleaned
}

// sanitizeEmail sanitizes an email address
func sanitizeEmail(email string) string {
	return strings.ToLower(sanitizeString(email))
}
