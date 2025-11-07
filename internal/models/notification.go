package models

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID                   uuid.UUID  `json:"id" db:"id"`
	DebtListID           uuid.UUID  `json:"debt_list_id" db:"debt_list_id"`
	DebtItemID           *uuid.UUID `json:"debt_item_id" db:"debt_item_id"` // For payment confirmation notifications
	InstallmentNumber    *int       `json:"installment_number" db:"installment_number"` // Which payment in the schedule (1, 2, 3, etc.)
	InstallmentDueDate   *time.Time `json:"installment_due_date" db:"installment_due_date"` // Due date of this specific installment
	NotificationType     string     `json:"notification_type" db:"notification_type"` // email, sms, webhook
	WebhookType          *string    `json:"webhook_type" db:"webhook_type"` // slack, telegram, discord (if notification_type is webhook)
	RecipientType        string     `json:"recipient_type" db:"recipient_type"` // 'user' or 'contact' - who receives this notification
	Message              string     `json:"message" db:"message"`
	Status               string     `json:"status" db:"status"`
	SentAt               *time.Time `json:"sent_at" db:"sent_at"`
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`
	
	// Contact info is fetched from user_contacts table via debt_list -> contact_id
	// These fields are for caching/override purposes only (optional)
	RecipientEmail       *string    `json:"recipient_email,omitempty" db:"recipient_email"`
	RecipientPhone       *string    `json:"recipient_phone,omitempty" db:"recipient_phone"`
}

type CreateNotificationRequest struct {
	DebtListID           uuid.UUID  `json:"debt_list_id" binding:"required"`
	InstallmentNumber    *int       `json:"installment_number"` // For installment debts: which payment (1, 2, 3, etc.)
	InstallmentDueDate   *time.Time `json:"installment_due_date"` // For installment debts: due date of this payment
	NotificationType     string     `json:"notification_type" binding:"required,oneof=email sms webhook"`
	WebhookType          *string    `json:"webhook_type" binding:"omitempty,oneof=slack telegram discord"` // Required if notification_type is webhook
	Message              string     `json:"message" binding:"required"`
	// Optional: Override contact info from user_contacts table
	RecipientEmail       *string    `json:"recipient_email,omitempty"`
	RecipientPhone       *string    `json:"recipient_phone,omitempty"`
} 