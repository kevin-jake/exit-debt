# Notification System - Quick Start Implementation Guide

This guide provides step-by-step instructions with code snippets to implement the notification system.

---

## Step 1: Install Dependencies

```bash
# Core dependencies
go get github.com/go-co-op/gocron
go get gopkg.in/gomail.v2
go get github.com/twilio/twilio-go

# Optional: For distributed locking (multi-instance deployments)
go get github.com/go-co-op/gocron-redis-lock
```

---

## Step 2: Update Environment Variables

Add to your `.env` file:

```env
# Email Configuration (SMTP)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM_EMAIL=noreply@exitdebt.com
SMTP_FROM_NAME=Exit Debt

# SMS Configuration (Twilio)
TWILIO_ACCOUNT_SID=ACxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
TWILIO_AUTH_TOKEN=your_auth_token_here
TWILIO_PHONE_NUMBER=+1234567890

# Notification Service
NOTIFICATION_WORKER_INTERVAL=1m
NOTIFICATION_MAX_RETRY=3
NOTIFICATION_BATCH_SIZE=100

# Go-Cron UI
CRON_UI_ENABLED=true
CRON_UI_PORT=8081
```

---

## Step 3: Database Migrations

### 3.1 Update `user_settings` table

```sql
-- Migration: Add notification schedule fields to user_settings
ALTER TABLE user_settings
ADD COLUMN notification_reminder_days INTEGER[] DEFAULT '{7,3,1}',
ADD COLUMN notification_time TIME DEFAULT '09:00:00',
ADD COLUMN overdue_reminder_frequency VARCHAR(50) DEFAULT 'daily',
ADD COLUMN custom_email_message TEXT,
ADD COLUMN custom_sms_message TEXT;
```

### 3.2 Update `notifications` table

**Important Update:** Notification is now linked to `debt_list_id` (not `debt_item_id`). Contact email/phone comes from `user_contacts` table.

```sql
-- Migration: Update foreign key and add scheduling fields to notifications

-- Step 1: Update the foreign key relationship
ALTER TABLE notifications DROP COLUMN IF EXISTS debt_item_id;
ALTER TABLE notifications ADD COLUMN debt_list_id UUID REFERENCES debt_lists(id) ON DELETE CASCADE;

-- Step 2: Add scheduling fields
ALTER TABLE notifications
ADD COLUMN schedule_type VARCHAR(50),
ADD COLUMN scheduled_for TIMESTAMP,
ADD COLUMN cron_job_id VARCHAR(255),
ADD COLUMN reminder_days_before INTEGER,
ADD COLUMN use_custom_schedule BOOLEAN DEFAULT false,
ADD COLUMN custom_reminder_days INTEGER[],
ADD COLUMN custom_notification_time TIME,
ADD COLUMN custom_message TEXT,
ADD COLUMN last_sent_at TIMESTAMP,
ADD COLUMN next_run_at TIMESTAMP,
ADD COLUMN is_recurring BOOLEAN DEFAULT false,
ADD COLUMN enabled BOOLEAN DEFAULT true;

-- Step 3: Add indexes for performance
CREATE INDEX idx_notifications_debt_list ON notifications(debt_list_id);
CREATE INDEX idx_notifications_scheduled_for ON notifications(scheduled_for) WHERE status = 'pending';
CREATE INDEX idx_notifications_next_run ON notifications(next_run_at) WHERE enabled = true;
CREATE INDEX idx_notifications_schedule_type ON notifications(schedule_type);

-- Note: recipient_email and recipient_phone are optional override fields
-- Primary contact info is fetched from user_contacts table via:
--   notifications.debt_list_id -> debt_lists.contact_id + debt_lists.user_id
--   -> user_contacts (WHERE contact_id AND user_id) -> email, phone
```

### 3.3 Create `notification_templates` table

```sql
-- Migration: Create notification_templates table
CREATE TABLE notification_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    template_name VARCHAR(255) NOT NULL,
    template_type VARCHAR(50) NOT NULL,
    subject VARCHAR(500),
    body TEXT NOT NULL,
    is_default BOOLEAN DEFAULT false,
    variables TEXT[],
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_notification_templates_user ON notification_templates(user_id);
CREATE INDEX idx_notification_templates_type ON notification_templates(template_type);
```

---

## Step 4: Update Go Models

### 4.1 Update `internal/models/user_settings.go`

```go
package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type UserSettings struct {
	ID                       uuid.UUID      `json:"id" gorm:"type:uuid;primary_key"`
	UserID                   uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;uniqueIndex"`
	NotificationEmail        bool           `json:"notification_email" gorm:"default:true"`
	NotificationSMS          bool           `json:"notification_sms" gorm:"default:false"`
	DefaultCurrency          string         `json:"default_currency" gorm:"default:'Php'"`
	Timezone                 string         `json:"timezone" gorm:"default:'UTC'"`

	// New notification schedule fields
	NotificationReminderDays pq.Int64Array  `json:"notification_reminder_days" gorm:"type:integer[];default:'{7,3,1}'"`
	NotificationTime         string         `json:"notification_time" gorm:"default:'09:00:00'"`
	OverdueReminderFrequency string         `json:"overdue_reminder_frequency" gorm:"default:'daily'"`
	CustomEmailMessage       *string        `json:"custom_email_message"`
	CustomSMSMessage         *string        `json:"custom_sms_message"`

	CreatedAt                time.Time      `json:"created_at"`
	UpdatedAt                time.Time      `json:"updated_at"`
}

type UpdateUserSettingsRequest struct {
	NotificationEmail        *bool          `json:"notification_email"`
	NotificationSMS          *bool          `json:"notification_sms"`
	DefaultCurrency          *string        `json:"default_currency"`
	Timezone                 *string        `json:"timezone"`
	NotificationReminderDays *[]int         `json:"notification_reminder_days"`
	NotificationTime         *string        `json:"notification_time"`
	OverdueReminderFrequency *string        `json:"overdue_reminder_frequency"`
	CustomEmailMessage       *string        `json:"custom_email_message"`
	CustomSMSMessage         *string        `json:"custom_sms_message"`
}
```

### 4.2 Update `internal/models/notification.go`

**Key Changes:**

- Changed from `debt_item_id` to `debt_list_id`
- Contact email/phone fetched from `user_contacts` table
- `recipient_email` and `recipient_phone` are now optional override fields

```go
package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Notification struct {
	ID                      uuid.UUID      `json:"id" gorm:"type:uuid;primary_key"`
	DebtListID              uuid.UUID      `json:"debt_list_id" gorm:"type:uuid;not null;index"`
	NotificationType        string         `json:"notification_type" gorm:"not null"`
	Message                 string         `json:"message" gorm:"type:text;not null"`
	Status                  string         `json:"status" gorm:"default:'pending';index"`
	SentAt                  *time.Time     `json:"sent_at"`
	CreatedAt               time.Time      `json:"created_at"`

	// Contact info (optional override - fetched from user_contacts by default)
	// Primary source: user_contacts table via debt_list -> contact_id
	RecipientEmail          *string        `json:"recipient_email,omitempty"`
	RecipientPhone          *string        `json:"recipient_phone,omitempty"`

	// Scheduling fields
	ScheduleType            string         `json:"schedule_type" gorm:"index"`
	ScheduledFor            *time.Time     `json:"scheduled_for" gorm:"index"`
	CronJobID               *string        `json:"cron_job_id"`
	ReminderDaysBefore      *int           `json:"reminder_days_before"`
	UseCustomSchedule       bool           `json:"use_custom_schedule" gorm:"default:false"`
	CustomReminderDays      pq.Int64Array  `json:"custom_reminder_days" gorm:"type:integer[]"`
	CustomNotificationTime  *string        `json:"custom_notification_time"`
	CustomMessage           *string        `json:"custom_message" gorm:"type:text"`
	LastSentAt              *time.Time     `json:"last_sent_at"`
	NextRunAt               *time.Time     `json:"next_run_at" gorm:"index"`
	IsRecurring             bool           `json:"is_recurring" gorm:"default:false"`
	Enabled                 bool           `json:"enabled" gorm:"default:true;index"`

	// Relationships
	DebtList                DebtList       `json:"debt_list,omitempty" gorm:"foreignKey:DebtListID"`
}

type CreateNotificationRequest struct {
	DebtListID              uuid.UUID  `json:"debt_list_id" binding:"required"`
	NotificationType        string     `json:"notification_type" binding:"required,oneof=email sms"`
	Message                 string     `json:"message" binding:"required"`

	// Optional: Override contact info from user_contacts
	RecipientEmail          *string    `json:"recipient_email,omitempty"`
	RecipientPhone          *string    `json:"recipient_phone,omitempty"`

	// Scheduling
	ScheduleType            string     `json:"schedule_type"`
	ScheduledFor            *time.Time `json:"scheduled_for"`
	UseCustomSchedule       *bool      `json:"use_custom_schedule"`
	CustomReminderDays      *[]int     `json:"custom_reminder_days"`
	CustomNotificationTime  *string    `json:"custom_notification_time"`
	CustomMessage           *string    `json:"custom_message"`
}

type UpdateNotificationRequest struct {
	UseCustomSchedule       *bool     `json:"use_custom_schedule"`
	CustomReminderDays      *[]int    `json:"custom_reminder_days"`
	CustomNotificationTime  *string   `json:"custom_notification_time"`
	CustomMessage           *string   `json:"custom_message"`
	Enabled                 *bool     `json:"enabled"`
}
```

### 4.3 Create `internal/models/notification_template.go`

```go
package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type NotificationTemplate struct {
	ID           uuid.UUID      `json:"id" gorm:"type:uuid;primary_key"`
	UserID       *uuid.UUID     `json:"user_id" gorm:"type:uuid"`
	TemplateName string         `json:"template_name" gorm:"not null"`
	TemplateType string         `json:"template_type" gorm:"not null"`
	Subject      *string        `json:"subject"`
	Body         string         `json:"body" gorm:"not null"`
	IsDefault    bool           `json:"is_default" gorm:"default:false"`
	Variables    pq.StringArray `json:"variables" gorm:"type:text[]"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

type CreateTemplateRequest struct {
	TemplateName string   `json:"template_name" binding:"required"`
	TemplateType string   `json:"template_type" binding:"required,oneof=email sms"`
	Subject      *string  `json:"subject"`
	Body         string   `json:"body" binding:"required"`
	IsDefault    *bool    `json:"is_default"`
	Variables    []string `json:"variables"`
}

type UpdateTemplateRequest struct {
	TemplateName *string  `json:"template_name"`
	Subject      *string  `json:"subject"`
	Body         *string  `json:"body"`
	IsDefault    *bool    `json:"is_default"`
	Variables    *[]string `json:"variables"`
}
```

---

## Step 5: Create Configuration

Create `internal/config/notification_config.go`:

```go
package config

import (
	"os"
	"strconv"
	"time"
)

type NotificationConfig struct {
	// SMTP Configuration
	SMTPHost      string
	SMTPPort      int
	SMTPUsername  string
	SMTPPassword  string
	SMTPFromEmail string
	SMTPFromName  string

	// Twilio Configuration
	TwilioAccountSID string
	TwilioAuthToken  string
	TwilioPhoneNumber string

	// Worker Configuration
	WorkerInterval time.Duration
	MaxRetry       int
	BatchSize      int

	// Cron UI
	CronUIEnabled bool
	CronUIPort    int
}

func LoadNotificationConfig() *NotificationConfig {
	smtpPort, _ := strconv.Atoi(getEnv("SMTP_PORT", "587"))
	maxRetry, _ := strconv.Atoi(getEnv("NOTIFICATION_MAX_RETRY", "3"))
	batchSize, _ := strconv.Atoi(getEnv("NOTIFICATION_BATCH_SIZE", "100"))
	cronUIPort, _ := strconv.Atoi(getEnv("CRON_UI_PORT", "8081"))

	workerInterval, _ := time.ParseDuration(getEnv("NOTIFICATION_WORKER_INTERVAL", "1m"))

	cronUIEnabled, _ := strconv.ParseBool(getEnv("CRON_UI_ENABLED", "true"))

	return &NotificationConfig{
		SMTPHost:          getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:          smtpPort,
		SMTPUsername:      getEnv("SMTP_USERNAME", ""),
		SMTPPassword:      getEnv("SMTP_PASSWORD", ""),
		SMTPFromEmail:     getEnv("SMTP_FROM_EMAIL", "noreply@exitdebt.com"),
		SMTPFromName:      getEnv("SMTP_FROM_NAME", "Exit Debt"),
		TwilioAccountSID:  getEnv("TWILIO_ACCOUNT_SID", ""),
		TwilioAuthToken:   getEnv("TWILIO_AUTH_TOKEN", ""),
		TwilioPhoneNumber: getEnv("TWILIO_PHONE_NUMBER", ""),
		WorkerInterval:    workerInterval,
		MaxRetry:          maxRetry,
		BatchSize:         batchSize,
		CronUIEnabled:     cronUIEnabled,
		CronUIPort:        cronUIPort,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
```

---

## Step 6: Create Helper Function to Fetch Contact Info

Create `internal/services/notification/contact_fetcher.go`:

```go
package notification

import (
	"fmt"

	"exit-debt/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ContactInfo struct {
	Email       *string
	Phone       *string
	Name        string
	UserID      uuid.UUID
	ContactID   uuid.UUID
}

type ContactFetcher struct {
	db *gorm.DB
}

func NewContactFetcher(db *gorm.DB) *ContactFetcher {
	return &ContactFetcher{db: db}
}

// GetContactInfoForNotification fetches contact information from user_contacts table
// via the relationship: notification -> debt_list -> contact_id + user_id -> user_contact
func (cf *ContactFetcher) GetContactInfoForNotification(notificationID uuid.UUID) (*ContactInfo, error) {
	var result ContactInfo

	err := cf.db.Table("notifications n").
		Select("uc.email, uc.phone, uc.name, dl.user_id, dl.contact_id").
		Joins("JOIN debt_lists dl ON n.debt_list_id = dl.id").
		Joins("JOIN user_contacts uc ON uc.contact_id = dl.contact_id AND uc.user_id = dl.user_id").
		Where("n.id = ?", notificationID).
		Scan(&result).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch contact info: %w", err)
	}

	return &result, nil
}

// GetContactInfoForDebtList fetches contact information directly from debt_list_id
func (cf *ContactFetcher) GetContactInfoForDebtList(debtListID uuid.UUID) (*ContactInfo, error) {
	var result ContactInfo

	err := cf.db.Table("debt_lists dl").
		Select("uc.email, uc.phone, uc.name, dl.user_id, dl.contact_id").
		Joins("JOIN user_contacts uc ON uc.contact_id = dl.contact_id AND uc.user_id = dl.user_id").
		Where("dl.id = ?", debtListID).
		Scan(&result).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch contact info: %w", err)
	}

	return &result, nil
}

// GetEffectiveRecipient determines the effective recipient email/phone
// Override logic: if notification has recipient_* field set, use it; else use from user_contacts
func GetEffectiveRecipient(notification *models.Notification, contactInfo *ContactInfo) (email *string, phone *string) {
	// Email: use override if present, else use from user_contacts
	if notification.RecipientEmail != nil && *notification.RecipientEmail != "" {
		email = notification.RecipientEmail
	} else {
		email = contactInfo.Email
	}

	// Phone: use override if present, else use from user_contacts
	if notification.RecipientPhone != nil && *notification.RecipientPhone != "" {
		phone = notification.RecipientPhone
	} else {
		phone = contactInfo.Phone
	}

	return email, phone
}

// DebtStatus represents the payment status of a debt
type DebtStatus struct {
	Status             string
	TotalRemainingDebt decimal.Decimal
	NextPaymentDate    time.Time
	DueDate            time.Time
}

// GetDebtStatus fetches the current payment status of a debt
func (cf *ContactFetcher) GetDebtStatus(debtListID uuid.UUID) (*DebtStatus, error) {
	var result DebtStatus

	err := cf.db.Table("debt_lists").
		Select("status, total_remaining_debt, next_payment_date, due_date").
		Where("id = ?", debtListID).
		Scan(&result).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch debt status: %w", err)
	}

	return &result, nil
}

// ShouldSendNotification checks if notification should be sent based on debt status
// Returns true if debt is NOT yet paid/settled
func ShouldSendNotification(debtStatus *DebtStatus) (bool, string) {
	// Check if debt is settled
	if debtStatus.Status == "settled" {
		return false, "debt already settled"
	}

	// Check if debt is fully paid
	if debtStatus.TotalRemainingDebt.LessThanOrEqual(decimal.Zero) {
		return false, "debt fully paid"
	}

	// All checks passed - should send notification
	return true, ""
}
```

**Usage Example:**

```go
// In your notification worker
func processNotification(notif *models.Notification) error {
	contactFetcher := NewContactFetcher(db)

	// 1. Check debt payment status - ONLY send if NOT yet paid
	debtStatus, err := contactFetcher.GetDebtStatus(notif.DebtListID)
	if err != nil {
		return fmt.Errorf("failed to get debt status: %w", err)
	}

	// 2. Check if notification should be sent
	shouldSend, reason := ShouldSendNotification(debtStatus)
	if !shouldSend {
		log.Infof("Skipping notification %s: %s", notif.ID, reason)
		// Disable future notifications for settled/paid debts
		if reason == "debt already settled" || reason == "debt fully paid" {
			disableNotification(notif.ID)
		}
		return nil
	}

	// 3. Fetch contact info from user_contacts
	contactInfo, err := contactFetcher.GetContactInfoForNotification(notif.ID)
	if err != nil {
		return fmt.Errorf("failed to get contact info: %w", err)
	}

	// 4. Get effective recipient (with override logic)
	email, phone := GetEffectiveRecipient(notif, contactInfo)

	// 5. Validate recipient exists
	if notif.NotificationType == "email" && email == nil {
		return fmt.Errorf("no email address available")
	}
	if notif.NotificationType == "sms" && phone == nil {
		return fmt.Errorf("no phone number available")
	}

	// 6. Send based on notification type
	if notif.NotificationType == "email" {
		return emailSender.SendEmail(*email, subject, body)
	}

	if notif.NotificationType == "sms" {
		return smsSender.SendSMS(*phone, message)
	}

	return fmt.Errorf("invalid notification type")
}
```

---

## Step 7: Create Email Service

Create `internal/services/notification/email_sender.go`:

```go
package notification

import (
	"fmt"

	"exit-debt/internal/config"
	"gopkg.in/gomail.v2"
)

type EmailSender struct {
	config *config.NotificationConfig
}

func NewEmailSender(cfg *config.NotificationConfig) *EmailSender {
	return &EmailSender{config: cfg}
}

func (s *EmailSender) SendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", s.config.SMTPFromName, s.config.SMTPFromEmail))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(
		s.config.SMTPHost,
		s.config.SMTPPort,
		s.config.SMTPUsername,
		s.config.SMTPPassword,
	)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
```

---

## Step 7: Create SMS Service

Create `internal/services/notification/sms_sender.go`:

```go
package notification

import (
	"fmt"

	"exit-debt/internal/config"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type SMSSender struct {
	config *config.NotificationConfig
	client *twilio.RestClient
}

func NewSMSSender(cfg *config.NotificationConfig) *SMSSender {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: cfg.TwilioAccountSID,
		Password: cfg.TwilioAuthToken,
	})

	return &SMSSender{
		config: cfg,
		client: client,
	}
}

func (s *SMSSender) SendSMS(to, message string) error {
	params := &twilioApi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(s.config.TwilioPhoneNumber)
	params.SetBody(message)

	_, err := s.client.Api.CreateMessage(params)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %w", err)
	}

	return nil
}
```

---

## Step 8: Create Template Engine

Create `internal/services/notification/template_engine.go`:

```go
package notification

import (
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

type TemplateData struct {
	UserFirstName    string
	UserLastName     string
	DebtorName       string
	Amount           decimal.Decimal
	Currency         string
	DueDate          time.Time
	NextPaymentDate  time.Time
	DaysUntilDue     int
	DaysOverdue      int
	TotalDebt        decimal.Decimal
	RemainingDebt    decimal.Decimal
	PaymentMethod    string
	ContactEmail     string
	ContactPhone     string
}

type TemplateEngine struct{}

func NewTemplateEngine() *TemplateEngine {
	return &TemplateEngine{}
}

func (e *TemplateEngine) Render(template string, data TemplateData) string {
	result := template

	replacements := map[string]string{
		"{{user_first_name}}":   data.UserFirstName,
		"{{user_last_name}}":    data.UserLastName,
		"{{debtor_name}}":       data.DebtorName,
		"{{amount}}":            data.Amount.String(),
		"{{currency}}":          data.Currency,
		"{{due_date}}":          data.DueDate.Format("2006-01-02"),
		"{{next_payment_date}}": data.NextPaymentDate.Format("2006-01-02"),
		"{{days_until_due}}":    fmt.Sprintf("%d", data.DaysUntilDue),
		"{{days_overdue}}":      fmt.Sprintf("%d", data.DaysOverdue),
		"{{total_debt}}":        data.TotalDebt.String(),
		"{{remaining_debt}}":    data.RemainingDebt.String(),
		"{{payment_method}}":    data.PaymentMethod,
		"{{contact_email}}":     data.ContactEmail,
		"{{contact_phone}}":     data.ContactPhone,
	}

	for placeholder, value := range replacements {
		result = strings.ReplaceAll(result, placeholder, value)
	}

	return result
}

// Default templates
func (e *TemplateEngine) GetDefaultEmailReminderTemplate() string {
	return `
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #4CAF50; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background-color: #f9f9f9; }
        .details { background-color: white; padding: 15px; margin: 10px 0; border-left: 4px solid #4CAF50; }
        .footer { text-align: center; padding: 20px; color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h2>Payment Reminder</h2>
        </div>
        <div class="content">
            <p>Hi {{user_first_name}},</p>
            <p>This is a friendly reminder that you have a payment due to <strong>{{debtor_name}}</strong> in <strong>{{days_until_due}} days</strong>.</p>

            <div class="details">
                <h3>Payment Details:</h3>
                <ul>
                    <li><strong>Amount:</strong> {{currency}} {{amount}}</li>
                    <li><strong>Due Date:</strong> {{due_date}}</li>
                    <li><strong>Remaining Balance:</strong> {{currency}} {{remaining_debt}}</li>
                </ul>
            </div>

            <p>Please ensure payment is made on time to maintain your payment schedule.</p>
        </div>
        <div class="footer">
            <p>Best regards,<br>Exit Debt Team</p>
        </div>
    </div>
</body>
</html>
`
}

func (e *TemplateEngine) GetDefaultSMSReminderTemplate() string {
	return "Reminder: Payment of {{currency}}{{amount}} due to {{debtor_name}} in {{days_until_due}} days ({{due_date}}). Remaining: {{currency}}{{remaining_debt}}."
}

func (e *TemplateEngine) GetDefaultEmailOverdueTemplate() string {
	return `
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #f44336; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background-color: #f9f9f9; }
        .details { background-color: white; padding: 15px; margin: 10px 0; border-left: 4px solid #f44336; }
        .footer { text-align: center; padding: 20px; color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h2>⚠️ Overdue Payment Alert</h2>
        </div>
        <div class="content">
            <p>Hi {{user_first_name}},</p>
            <p>Your payment to <strong>{{debtor_name}}</strong> is now <strong>{{days_overdue}} days overdue</strong>.</p>

            <div class="details">
                <h3>Payment Details:</h3>
                <ul>
                    <li><strong>Amount Due:</strong> {{currency}} {{amount}}</li>
                    <li><strong>Original Due Date:</strong> {{due_date}}</li>
                    <li><strong>Total Remaining:</strong> {{currency}} {{remaining_debt}}</li>
                </ul>
            </div>

            <p><strong>Please make this payment as soon as possible.</strong></p>
        </div>
        <div class="footer">
            <p>Best regards,<br>Exit Debt Team</p>
        </div>
    </div>
</body>
</html>
`
}

func (e *TemplateEngine) GetDefaultSMSOverdueTemplate() string {
	return "OVERDUE: Payment of {{currency}}{{amount}} to {{debtor_name}} is {{days_overdue}} days overdue (due: {{due_date}}). Please pay immediately."
}
```

---

## Step 9: Create Cron Service

Create `internal/services/scheduler/cron_service.go`:

```go
package scheduler

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
)

type CronService struct {
	scheduler *gocron.Scheduler
	jobs      map[string]*gocron.Job
	mu        sync.RWMutex
}

func NewCronService() *CronService {
	s := gocron.NewScheduler(time.UTC)
	s.StartAsync()

	return &CronService{
		scheduler: s,
		jobs:      make(map[string]*gocron.Job),
	}
}

// ScheduleAt schedules a job at a specific time
func (cs *CronService) ScheduleAt(jobID string, scheduledTime time.Time, task func()) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	// Remove existing job if any
	if existingJob, exists := cs.jobs[jobID]; exists {
		cs.scheduler.RemoveByReference(existingJob)
	}

	// Schedule new job
	job, err := cs.scheduler.Every(1).Day().At(scheduledTime.Format("15:04")).Do(task)
	if err != nil {
		return fmt.Errorf("failed to schedule job: %w", err)
	}

	cs.jobs[jobID] = job
	return nil
}

// ScheduleRecurring schedules a recurring job
func (cs *CronService) ScheduleRecurring(jobID string, interval time.Duration, task func()) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	// Remove existing job if any
	if existingJob, exists := cs.jobs[jobID]; exists {
		cs.scheduler.RemoveByReference(existingJob)
	}

	// Schedule new recurring job
	job, err := cs.scheduler.Every(interval).Do(task)
	if err != nil {
		return fmt.Errorf("failed to schedule recurring job: %w", err)
	}

	cs.jobs[jobID] = job
	return nil
}

// Remove removes a job by ID
func (cs *CronService) Remove(jobID string) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	job, exists := cs.jobs[jobID]
	if !exists {
		return fmt.Errorf("job not found: %s", jobID)
	}

	cs.scheduler.RemoveByReference(job)
	delete(cs.jobs, jobID)
	return nil
}

// Stop stops the scheduler
func (cs *CronService) Stop() {
	cs.scheduler.Stop()
}

// GetJobs returns all active jobs
func (cs *CronService) GetJobs() []*gocron.Job {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	return cs.scheduler.Jobs()
}
```

---

## Step 10: Update Database Migration

Update `internal/database/database.go` to include new models:

```go
// Auto migrate the schema
if err := db.AutoMigrate(
	&models.User{},
	&models.UserSettings{},
	&models.Contact{},
	&models.UserContact{},
	&models.DebtList{},
	&models.DebtItem{},
	&models.Notification{},
	&models.NotificationTemplate{},  // Add this
); err != nil {
	return nil, fmt.Errorf("failed to migrate database: %v", err)
}
```

---

## Step 11: Test Email Sending

Create a simple test file `test_email.go`:

```go
package main

import (
	"log"

	"exit-debt/internal/config"
	"exit-debt/internal/services/notification"
)

func main() {
	cfg := config.LoadNotificationConfig()
	emailSender := notification.NewEmailSender(cfg)

	err := emailSender.SendEmail(
		"test@example.com",
		"Test Email",
		"<h1>This is a test email</h1><p>From Exit Debt</p>",
	)

	if err != nil {
		log.Fatalf("Failed to send email: %v", err)
	}

	log.Println("Email sent successfully!")
}
```

Run: `go run test_email.go`

---

## Step 12: Test SMS Sending

Create a simple test file `test_sms.go`:

```go
package main

import (
	"log"

	"exit-debt/internal/config"
	"exit-debt/internal/services/notification"
)

func main() {
	cfg := config.LoadNotificationConfig()
	smsSender := notification.NewSMSSender(cfg)

	err := smsSender.SendSMS(
		"+1234567890",  // Replace with actual phone number
		"Test SMS from Exit Debt",
	)

	if err != nil {
		log.Fatalf("Failed to send SMS: %v", err)
	}

	log.Println("SMS sent successfully!")
}
```

Run: `go run test_sms.go`

---

## Next Steps

After completing these steps:

1. ✅ Dependencies installed
2. ✅ Environment configured
3. ✅ Database schema updated
4. ✅ Models created
5. ✅ Email service working
6. ✅ SMS service working
7. ✅ Template engine ready
8. ✅ Cron service initialized

**Now proceed to:**

- Implement notification repositories (database operations)
- Create notification service with business logic
- Build API handlers for notification management
- Integrate with existing debt creation flow
- Set up Go-Cron UI
- Add notification worker logic

Refer to `NOTIFICATION_IMPLEMENTATION_PLAN.md` for the complete implementation details.

---

## Quick Commands Reference

```bash
# Install dependencies
go mod tidy

# Run migrations (using your existing migration tool)
# Or apply SQL migrations manually

# Test email
go run test_email.go

# Test SMS
go run test_sms.go

# Run application
go run cmd/main.go

# Access Cron UI (once implemented)
# http://localhost:8081/cron-ui
```

---

## Troubleshooting

### Email not sending

- Check SMTP credentials
- For Gmail, use App Password (not regular password)
- Ensure "Less secure app access" is enabled or use OAuth2

### SMS not sending

- Verify Twilio credentials
- Check Twilio phone number is active
- Ensure target number is verified (for trial accounts)
- Check Twilio console for error messages

### Database migration errors

- Ensure PostgreSQL supports arrays (should be default)
- Check user permissions
- Verify column types match your PostgreSQL version

### Cron jobs not running

- Check system time/timezone
- Verify scheduler is started
- Check logs for errors
- Ensure database records have correct scheduled_for times
