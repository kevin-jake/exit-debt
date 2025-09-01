package models

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID                   uuid.UUID `json:"id" db:"id"`
	DebtItemID           uuid.UUID `json:"debt_item_id" db:"debt_item_id"`
	NotificationType     string    `json:"notification_type" db:"notification_type"`
	RecipientEmail       *string   `json:"recipient_email" db:"recipient_email"`
	RecipientPhone       *string   `json:"recipient_phone" db:"recipient_phone"`
	Message              string    `json:"message" db:"message"`
	Status               string    `json:"status" db:"status"`
	SentAt               *time.Time `json:"sent_at" db:"sent_at"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
}

type CreateNotificationRequest struct {
	DebtItemID           uuid.UUID `json:"debt_item_id" binding:"required"`
	NotificationType     string    `json:"notification_type" binding:"required,oneof=email sms"`
	RecipientEmail       *string   `json:"recipient_email"`
	RecipientPhone       *string   `json:"recipient_phone"`
	Message              string    `json:"message" binding:"required"`
} 