package models

import (
	"time"

	"github.com/google/uuid"
)

type UserSettings struct {
	ID                  uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	UserID              uuid.UUID `json:"user_id" gorm:"type:uuid;not null;uniqueIndex"`
	NotificationEmail   bool      `json:"notification_email" gorm:"default:true"`
	NotificationSMS     bool      `json:"notification_sms" gorm:"default:false"`
	NotificationWebhook bool      `json:"notification_webhook" gorm:"default:false"`
	
	// Webhook configurations
	SlackWebhookURL     *string   `json:"slack_webhook_url"`
	TelegramBotToken    *string   `json:"telegram_bot_token"`
	TelegramChatID      *string   `json:"telegram_chat_id"`
	DiscordWebhookURL   *string   `json:"discord_webhook_url"`
	
	// Event notification settings
	EventNotificationsEnabled bool  `json:"event_notifications_enabled" gorm:"default:true"`
	NotifyContactOnPayment    bool  `json:"notify_contact_on_payment" gorm:"default:true"`
	
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
	NotificationWebhook *bool   `json:"notification_webhook"`
	SlackWebhookURL     *string `json:"slack_webhook_url"`
	TelegramBotToken         *string `json:"telegram_bot_token"`
	TelegramChatID           *string `json:"telegram_chat_id"`
	DiscordWebhookURL        *string `json:"discord_webhook_url"`
	EventNotificationsEnabled *bool  `json:"event_notifications_enabled"`
	NotifyContactOnPayment    *bool  `json:"notify_contact_on_payment"`
	DefaultCurrency          *string `json:"default_currency"`
	Timezone                 *string `json:"timezone"`
} 