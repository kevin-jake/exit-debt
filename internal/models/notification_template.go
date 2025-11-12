package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type NotificationTemplate struct {
	ID           uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID       *uuid.UUID     `json:"user_id" gorm:"type:uuid;index"` // nil for system defaults
	TemplateName string         `json:"template_name" gorm:"not null"`
	TemplateType string         `json:"template_type" gorm:"not null;index"` // 'email' or 'sms'
	Subject      *string        `json:"subject"` // For email only
	Body         string         `json:"body" gorm:"not null"`
	IsDefault    bool           `json:"is_default" gorm:"default:false"`
	Variables    pq.StringArray `json:"variables" gorm:"type:text[]"` // Available template variables
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

type CreateTemplateRequest struct {
	TemplateName string   `json:"template_name" binding:"required"`
	TemplateType string   `json:"template_type" binding:"required,oneof=email sms"`
	Subject      *string  `json:"subject"`
	Body         string   `json:"body" binding:"required"`
	Variables    []string `json:"variables"`
}

type UpdateTemplateRequest struct {
	TemplateName *string  `json:"template_name"`
	Subject      *string  `json:"subject"`
	Body         *string  `json:"body"`
	Variables    []string `json:"variables"`
}

