package models

import (
	"time"

	"github.com/google/uuid"
)

type UserSettings struct {
	ID                  uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID              uuid.UUID `json:"user_id" gorm:"type:uuid;not null;uniqueIndex"`
	NotificationEmail   bool      `json:"notification_email" gorm:"default:true"`
	NotificationSMS     bool      `json:"notification_sms" gorm:"default:false"`
	NotificationFacebook bool     `json:"notification_facebook" gorm:"default:false"`
	DefaultCurrency     string    `json:"default_currency" gorm:"default:'Php'"`
	Timezone           string    `json:"timezone" gorm:"default:'UTC'"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	
	// Relationships
	// User relationship is handled in User model to avoid circular reference
}

type UpdateUserSettingsRequest struct {
	NotificationEmail   *bool   `json:"notification_email"`
	NotificationSMS     *bool   `json:"notification_sms"`
	NotificationFacebook *bool  `json:"notification_facebook"`
	DefaultCurrency     *string `json:"default_currency"`
	Timezone           *string `json:"timezone"`
} 