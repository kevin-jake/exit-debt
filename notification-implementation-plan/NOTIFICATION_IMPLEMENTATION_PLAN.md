# Email & SMS Notification System Implementation Plan

## Overview

This document outlines the complete implementation plan for a scheduled notification system using go-cron with go-cron UI. The system sends email and SMS notifications for debt payment reminders with:

- **User-customizable schedules**
- **Per-installment notifications** for installment debts
- **Smart payment status checking** (only send for unpaid installments)
- **Contact info from `user_contacts` table**

---

## Critical Architecture Decisions

### 1. ✅ Installment-Based Notifications

**For installment debts (weekly, monthly, etc.), EACH installment payment gets its own notification schedule.**

**Example:**

```
Debt: $12,000, Monthly plan, 12 installments
├── Installment #1 (due 2025-01-15)
│   ├── Notification: 7 days before
│   ├── Notification: 3 days before
│   └── Notification: 1 day before
├── Installment #2 (due 2025-02-15)
│   ├── Notification: 7 days before
│   ├── Notification: 3 days before
│   └── Notification: 1 day before
... (continues for all installments)

TOTAL: 12 installments × 3 reminders = 36 notifications
```

### 2. ✅ Contact Info from user_contacts

**Recipient email/phone is fetched from `user_contacts` table, not stored in notifications.**

**Relationship Chain:**

```
Notification (debt_list_id + installment_number)
    → DebtList (contact_id + user_id)
        → UserContact (email, phone, name)
```

### 3. ✅ Webhook Support for Slack, Telegram & Discord

**Send notifications via webhooks to messaging platforms.**

**Supported Platforms:**

- **Slack** - Incoming Webhooks with rich blocks
- **Telegram** - Bot API with markdown formatting
- **Discord** - Webhooks with embed cards

### 4. ✅ Event-Based Notifications with Payment Verification

**Instant notifications when payments are made or debt items are added, with verification workflow.**

**Two notification types:**

1. **Time-based** (scheduled reminders) - Sent at scheduled times before due dates
2. **Event-based** (payment confirmations + verifications) - Sent immediately when events occur

**For each payment event:**

**Stage 1: Payment Made (Pending)**

- Notify BOTH **user** (debt owner) and **contact** (other party)
- Messages include "pending verification" status

**Stage 2: Payment Verified/Rejected**

- Notify the **payer** only (whoever made the payment)
- Verification message: "Payment verified and accepted"
- Rejection message: "Payment rejected. Reason: [reason]"

### 5. ✅ Payment Status Checking

**Only send notifications if the specific installment is NOT yet paid.**

```go
if installment.Status == "paid" || installment.PaidAmount >= installment.ScheduledAmount {
    skipNotification()  // Don't send
    disableFutureNotifications()  // Auto-disable remaining
}
```

---

## Database Schema Updates

> **Note:** This project uses a **clean database** approach. All schema changes will be implemented via GORM models and AutoMigrate, not SQL migrations. The SQL shown below is for reference - implement these as GORM struct fields.

### 1.1 User Settings Table Enhancement

Add default notification schedule fields to `user_settings`:

```sql
-- New fields to add to user_settings table
ALTER TABLE user_settings ADD COLUMN notification_reminder_days INTEGER[] DEFAULT '{7,3,1}';
ALTER TABLE user_settings ADD COLUMN notification_time TIME DEFAULT '09:00:00';
ALTER TABLE user_settings ADD COLUMN overdue_reminder_frequency VARCHAR(50) DEFAULT 'daily';
ALTER TABLE user_settings ADD COLUMN custom_email_message TEXT;
ALTER TABLE user_settings ADD COLUMN custom_sms_message TEXT;

-- Webhook notification settings
ALTER TABLE user_settings ADD COLUMN notification_webhook BOOLEAN DEFAULT false;
ALTER TABLE user_settings ADD COLUMN slack_webhook_url VARCHAR(500);
ALTER TABLE user_settings ADD COLUMN telegram_bot_token VARCHAR(255);
ALTER TABLE user_settings ADD COLUMN telegram_chat_id VARCHAR(255);
ALTER TABLE user_settings ADD COLUMN discord_webhook_url VARCHAR(500);

-- Event notification settings
ALTER TABLE user_settings ADD COLUMN event_notifications_enabled BOOLEAN DEFAULT true;
ALTER TABLE user_settings ADD COLUMN notify_contact_on_payment BOOLEAN DEFAULT true;
```

**Go Model Fields:**

```go
NotificationReminderDays    pq.Int64Array  `json:"notification_reminder_days" gorm:"type:integer[];default:'{7,3,1}'"`
NotificationTime            string         `json:"notification_time" gorm:"default:'09:00:00'"`
OverdueReminderFrequency    string         `json:"overdue_reminder_frequency" gorm:"default:'daily'"`
CustomEmailMessage          *string        `json:"custom_email_message"`
CustomSMSMessage            *string        `json:"custom_sms_message"`

// Webhook settings
NotificationWebhook         bool           `json:"notification_webhook" gorm:"default:false"`
SlackWebhookURL             *string        `json:"slack_webhook_url"`
TelegramBotToken            *string        `json:"telegram_bot_token"`
TelegramChatID              *string        `json:"telegram_chat_id"`
DiscordWebhookURL           *string        `json:"discord_webhook_url"`

// Event notification settings
EventNotificationsEnabled   bool           `json:"event_notifications_enabled" gorm:"default:true"`
NotifyContactOnPayment      bool           `json:"notify_contact_on_payment" gorm:"default:true"`
```

### 1.2 Notification Table Enhancement

**Key Changes:**

1. Linked to `debt_list_id` (not `debt_item_id`)
2. Added `installment_number` to track which payment
3. Added `installment_due_date` for the specific installment

```sql
-- Update foreign key relationship
ALTER TABLE notifications DROP COLUMN IF EXISTS debt_item_id;
ALTER TABLE notifications ADD COLUMN debt_list_id UUID REFERENCES debt_lists(id) ON DELETE CASCADE;

-- Add installment tracking fields
ALTER TABLE notifications ADD COLUMN installment_number INT;
ALTER TABLE notifications ADD COLUMN installment_due_date TIMESTAMP;

-- Add event-based notification fields
ALTER TABLE notifications ADD COLUMN debt_item_id UUID REFERENCES debt_items(id) ON DELETE CASCADE;
ALTER TABLE notifications ADD COLUMN recipient_type VARCHAR(50);  -- 'user' or 'contact'

-- Add webhook type field
ALTER TABLE notifications ADD COLUMN webhook_type VARCHAR(50);  -- 'slack', 'telegram', 'discord'

-- Add scheduling fields
ALTER TABLE notifications ADD COLUMN schedule_type VARCHAR(50);  -- 'reminder', 'overdue', 'manual', 'event'
ALTER TABLE notifications ADD COLUMN scheduled_for TIMESTAMP;
ALTER TABLE notifications ADD COLUMN cron_job_id VARCHAR(255);
ALTER TABLE notifications ADD COLUMN reminder_days_before INT;
ALTER TABLE notifications ADD COLUMN use_custom_schedule BOOLEAN DEFAULT false;
ALTER TABLE notifications ADD COLUMN custom_reminder_days INTEGER[];
ALTER TABLE notifications ADD COLUMN custom_notification_time TIME;
ALTER TABLE notifications ADD COLUMN custom_message TEXT;
ALTER TABLE notifications ADD COLUMN last_sent_at TIMESTAMP;
ALTER TABLE notifications ADD COLUMN next_run_at TIMESTAMP;
ALTER TABLE notifications ADD COLUMN is_recurring BOOLEAN DEFAULT false;
ALTER TABLE notifications ADD COLUMN enabled BOOLEAN DEFAULT true;

-- Add indexes for performance
CREATE INDEX idx_notifications_debt_list ON notifications(debt_list_id);
CREATE INDEX idx_notifications_debt_item ON notifications(debt_item_id);
CREATE INDEX idx_notifications_debt_installment ON notifications(debt_list_id, installment_number);
CREATE INDEX idx_notifications_scheduled_for ON notifications(scheduled_for) WHERE status = 'pending';
CREATE INDEX idx_notifications_installment_due ON notifications(installment_due_date) WHERE status = 'pending';
CREATE INDEX idx_notifications_next_run ON notifications(next_run_at) WHERE enabled = true;
CREATE INDEX idx_notifications_schedule_type ON notifications(schedule_type);
CREATE INDEX idx_notifications_recipient_type ON notifications(recipient_type);
```

**Go Model:**

```go
type Notification struct {
	ID                      uuid.UUID      `json:"id" db:"id"`
	DebtListID              uuid.UUID      `json:"debt_list_id" db:"debt_list_id"`
	DebtItemID              *uuid.UUID     `json:"debt_item_id" db:"debt_item_id"`  // For event-based notifications
	InstallmentNumber       *int           `json:"installment_number" db:"installment_number"`  // Which payment (1, 2, 3...)
	InstallmentDueDate      *time.Time     `json:"installment_due_date" db:"installment_due_date"`  // Due date of this installment
	NotificationType        string         `json:"notification_type" db:"notification_type"`  // 'email', 'sms', 'webhook'
	WebhookType             *string        `json:"webhook_type" db:"webhook_type"`  // 'slack', 'telegram', 'discord'
	RecipientType           string         `json:"recipient_type" db:"recipient_type"`  // 'user' or 'contact'
	Message                 string         `json:"message" db:"message"`
	Status                  string         `json:"status" db:"status"`
	SentAt                  *time.Time     `json:"sent_at" db:"sent_at"`
	CreatedAt               time.Time      `json:"created_at" db:"created_at"`

	// Scheduling fields
	ScheduleType            string         `json:"schedule_type" db:"schedule_type"`
	ScheduledFor            *time.Time     `json:"scheduled_for" db:"scheduled_for"`
	CronJobID               *string        `json:"cron_job_id" db:"cron_job_id"`
	ReminderDaysBefore      *int           `json:"reminder_days_before" db:"reminder_days_before"`
	UseCustomSchedule       bool           `json:"use_custom_schedule" db:"use_custom_schedule"`
	CustomReminderDays      pq.Int64Array  `json:"custom_reminder_days" db:"custom_reminder_days"`
	CustomNotificationTime  *string        `json:"custom_notification_time" db:"custom_notification_time"`
	CustomMessage           *string        `json:"custom_message" db:"custom_message"`
	LastSentAt              *time.Time     `json:"last_sent_at" db:"last_sent_at"`
	NextRunAt               *time.Time     `json:"next_run_at" db:"next_run_at"`
	IsRecurring             bool           `json:"is_recurring" db:"is_recurring"`
	Enabled                 bool           `json:"enabled" db:"enabled"`

	// Contact info (optional override - normally fetched from user_contacts)
	RecipientEmail          *string        `json:"recipient_email,omitempty" db:"recipient_email"`
	RecipientPhone          *string        `json:"recipient_phone,omitempty" db:"recipient_phone"`
}
```

### 1.3 Notification Templates Table (New)

Store reusable notification templates:

```sql
CREATE TABLE notification_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    template_name VARCHAR(255) NOT NULL,
    template_type VARCHAR(50) NOT NULL,  -- 'email' or 'sms'
    subject VARCHAR(500),  -- For email only
    body TEXT NOT NULL,
    is_default BOOLEAN DEFAULT false,
    variables TEXT[],  -- Available template variables
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_notification_templates_user ON notification_templates(user_id);
CREATE INDEX idx_notification_templates_type ON notification_templates(template_type);
```

---

## Installment Notification Creation Logic

### When Debt is Created

```go
func CreateNotificationsForDebt(debtList DebtList, userSettings UserSettings) error {
    // Step 1: Calculate payment schedule for all installments
    paymentSchedule := calculatePaymentSchedule(debtList)

    // Step 2: For EACH installment, create notification schedule
    for _, installment := range paymentSchedule {
        err := createNotificationsForInstallment(
            debtList.ID,
            installment,
            userSettings,
        )
        if err != nil {
            return err
        }
    }

    return nil
}

func createNotificationsForInstallment(
    debtListID uuid.UUID,
    installment PaymentScheduleItem,
    settings UserSettings,
) error {
    // Get reminder days from settings (or use default)
    reminderDays := settings.NotificationReminderDays  // [7, 3, 1]
    notificationTime := settings.NotificationTime       // "09:00:00"

    // Create notification for each reminder day
    for _, daysBefore := range reminderDays {
        // Calculate when to send this reminder
        scheduledFor := installment.DueDate.
            AddDate(0, 0, -daysBefore).  // Subtract days
            Add(parseTime(notificationTime))  // Add time

        notification := Notification{
            ID:                 uuid.New(),
            DebtListID:         debtListID,
            InstallmentNumber:  &installment.PaymentNumber,      // Track which installment
            InstallmentDueDate: &installment.DueDate,             // Track due date
            NotificationType:   "email",
            ScheduleType:       "reminder",
            ScheduledFor:       &scheduledFor,
            ReminderDaysBefore: &daysBefore,
            Status:             "pending",
            Enabled:            true,
        }

        // Save to database
        db.Create(&notification)

        // Schedule cron job
        cronJobID := scheduleCronJob(notification)

        // Update with cron job ID
        notification.CronJobID = &cronJobID
        db.Save(&notification)
    }

    return nil
}

func calculatePaymentSchedule(debtList DebtList) []PaymentScheduleItem {
    var schedule []PaymentScheduleItem

    numPayments := *debtList.NumberOfPayments
    currentDate := debtList.DueDate  // First payment date

    for i := 1; i <= numPayments; i++ {
        schedule = append(schedule, PaymentScheduleItem{
            PaymentNumber:   i,
            DueDate:         currentDate,
            ScheduledAmount: debtList.InstallmentAmount,
            Status:          "pending",
        })

        // Calculate next payment date based on installment plan
        currentDate = calculateNextPaymentDate(currentDate, debtList.InstallmentPlan)
    }

    return schedule
}
```

---

## Payment Status Checking Logic

### Before Sending Any Notification

```go
func processNotification(notificationID uuid.UUID) {
    // Step 1: Fetch notification
    notification := getNotification(notificationID)

    // Step 2: Check if specific installment is paid
    if !shouldSendNotificationForInstallment(notification) {
        log.Printf("Skipping notification - installment already paid")

        // Auto-disable remaining notifications for this installment
        disableNotificationsForInstallment(
            notification.DebtListID,
            *notification.InstallmentNumber,
        )

        return
    }

    // Step 3: Fetch contact info from user_contacts
    contactInfo, err := getContactInfoForNotification(notification.DebtListID)
    if err != nil {
        log.Printf("Failed to fetch contact info: %v", err)
        return
    }

    // Step 4: Get effective recipient (with override logic)
    email, phone := getEffectiveRecipient(notification, contactInfo)

    // Step 5: Render template with variables
    message := renderNotificationTemplate(notification, contactInfo)

    // Step 6: Send notification based on type
    switch notification.NotificationType {
    case "email":
        if email != nil {
            err = emailSender.SendEmail(*email, "Payment Reminder", message)
        }
    case "sms":
        if phone != nil {
            err = smsSender.SendSMS(*phone, message)
        }
    case "webhook":
        if notification.WebhookType != nil {
            userSettings := getUserSettings(notification.DebtListID)
            notifData := prepareNotificationData(notification, contactInfo)
            err = webhookService.SendNotification(*notification.WebhookType, userSettings, notifData)
        }
    }

    // Step 7: Update notification record
    if err == nil {
        updateNotificationSent(notificationID)
    } else {
        handleNotificationError(notificationID, err)
    }
}

func shouldSendNotificationForInstallment(notification Notification) bool {
    // Check if installment number is specified (nil for one-time debts)
    if notification.InstallmentNumber == nil {
        // For one-time debts, check overall debt status
        return shouldSendForOneTimeDebt(notification.DebtListID)
    }

    // Get installment info from payment schedule
    installment := getInstallmentInfo(
        notification.DebtListID,
        *notification.InstallmentNumber,
    )

    // Check 1: Is this installment already paid?
    if installment.Status == "paid" {
        return false
    }

    // Check 2: Is the paid amount >= scheduled amount?
    if installment.PaidAmount.GreaterThanOrEqual(installment.ScheduledAmount) {
        return false
    }

    // Check 3: Is the entire debt settled?
    debtList := getDebtList(notification.DebtListID)
    if debtList.Status == "settled" || debtList.TotalRemainingDebt.LessThanOrEqual(decimal.Zero) {
        return false
    }

    return true  // OK to send
}

func getInstallmentInfo(debtListID uuid.UUID, installmentNumber int) PaymentScheduleItem {
    // Query to get specific installment with payment status
    var installment PaymentScheduleItem

    db.Raw(`
        WITH payment_schedule AS (
            -- Calculate full payment schedule
            SELECT ... FROM debt_lists WHERE id = ?
        ),
        payments_made AS (
            -- Sum all payments for this installment
            SELECT
                SUM(amount) as paid_amount
            FROM debt_items
            WHERE debt_list_id = ?
              AND payment_date <= (
                  SELECT due_date FROM payment_schedule
                  WHERE payment_number = ?
              )
        )
        SELECT
            ps.*,
            COALESCE(pm.paid_amount, 0) as paid_amount,
            CASE
                WHEN COALESCE(pm.paid_amount, 0) >= ps.scheduled_amount THEN 'paid'
                ELSE 'pending'
            END as status
        FROM payment_schedule ps
        CROSS JOIN payments_made pm
        WHERE ps.payment_number = ?
    `, debtListID, debtListID, installmentNumber, installmentNumber).
    Scan(&installment)

    return installment
}
```

---

## Fetching Contact Information

```go
func getContactInfoForNotification(debtListID uuid.UUID) (*ContactInfo, error) {
    var contactInfo ContactInfo

    err := db.Raw(`
        SELECT
            uc.email,
            uc.phone,
            uc.name as contact_name,
            u.first_name as user_first_name,
            u.last_name as user_last_name,
            dl.installment_amount,
            dl.total_remaining_debt,
            dl.currency
        FROM notifications n
        JOIN debt_lists dl ON n.debt_list_id = dl.id
        JOIN user_contacts uc ON uc.contact_id = dl.contact_id
                              AND uc.user_id = dl.user_id
        JOIN users u ON u.id = dl.user_id
        WHERE dl.id = ?
    `, debtListID).Scan(&contactInfo).Error

    if err != nil {
        return nil, fmt.Errorf("failed to fetch contact info: %w", err)
    }

    return &contactInfo, nil
}

func getEffectiveRecipient(notification Notification, contactInfo *ContactInfo) (*string, *string) {
    // Use override if provided, otherwise use from contact info
    email := notification.RecipientEmail
    if email == nil {
        email = contactInfo.Email
    }

    phone := notification.RecipientPhone
    if phone == nil {
        phone = contactInfo.Phone
    }

    return email, phone
}
```

---

## Auto-Disable When Installment Paid

```go
func OnInstallmentPaid(debtListID uuid.UUID, installmentNumber int) {
    log.Printf("Installment #%d paid, disabling remaining notifications", installmentNumber)

    // Disable all pending notifications for this installment
    db.Exec(`
        UPDATE notifications
        SET enabled = false,
            status = 'cancelled',
            updated_at = NOW()
        WHERE debt_list_id = $1
          AND installment_number = $2
          AND status = 'pending'
          AND enabled = true
    `, debtListID, installmentNumber)

    // Remove cron jobs
    var notifications []Notification
    db.Where("debt_list_id = ? AND installment_number = ? AND status = 'cancelled'",
        debtListID, installmentNumber).
        Find(&notifications)

    for _, notif := range notifications {
        if notif.CronJobID != nil {
            cronService.Remove(*notif.CronJobID)
        }
    }

    log.Printf("Disabled %d notifications for installment #%d", len(notifications), installmentNumber)
}
```

---

## Template Variables for Installments

### Standard Variables

```
{{user_first_name}}         - User first name
{{user_last_name}}          - User last name
{{debtor_name}}             - Contact name
{{amount}}                  - Payment amount
{{currency}}                - Currency code
{{due_date}}                - Payment due date
{{total_debt}}              - Total debt amount
{{remaining_debt}}          - Remaining debt amount
```

### Installment-Specific Variables

```
{{installment_number}}      - Which payment (1, 2, 3, etc.)
{{installment_total}}       - Total number of installments
{{installment_due_date}}    - Due date of this specific installment
{{installment_amount}}      - Amount for this installment
{{remaining_installments}}  - How many payments left
{{payments_made_count}}     - How many installments already paid
{{days_until_due}}          - Days remaining until due
```

### Example Email Template for Installments

```html
<!DOCTYPE html>
<html>
  <head>
    <style>
      body {
        font-family: Arial, sans-serif;
      }
      .container {
        max-width: 600px;
        margin: 0 auto;
        padding: 20px;
      }
      .header {
        background-color: #4caf50;
        color: white;
        padding: 20px;
        text-align: center;
      }
      .content {
        padding: 20px;
        background-color: #f9f9f9;
      }
      .details {
        background-color: white;
        padding: 15px;
        margin: 10px 0;
        border-left: 4px solid #4caf50;
      }
    </style>
  </head>
  <body>
    <div class="container">
      <div class="header">
        <h2>Payment Reminder - Installment #{{installment_number}}</h2>
      </div>
      <div class="content">
        <p>Hi {{user_first_name}},</p>
        <p>
          This is a reminder that payment
          <strong>#{{installment_number}} of {{installment_total}}</strong> is
          due to <strong>{{debtor_name}}</strong> in
          <strong>{{days_until_due}} days</strong>.
        </p>

        <div class="details">
          <h3>Payment Details:</h3>
          <ul>
            <li>
              <strong>Installment:</strong> #{{installment_number}} of
              {{installment_total}}
            </li>
            <li>
              <strong>Amount:</strong> {{currency}} {{installment_amount}}
            </li>
            <li><strong>Due Date:</strong> {{installment_due_date}}</li>
            <li>
              <strong>Remaining Installments:</strong>
              {{remaining_installments}}
            </li>
          </ul>
        </div>

        <p>
          <strong>Total Remaining Debt:</strong> {{currency}} {{remaining_debt}}
        </p>
      </div>
    </div>
  </body>
</html>
```

---

## API Endpoints

### Debt & Notification Management

```
POST   /api/debt-lists                                  - Create debt (auto-schedules notifications)
POST   /api/debt-lists/:id/schedule-notifications       - Manually trigger notification scheduling
GET    /api/debt-lists/:id/notifications                - Get all notifications for debt
GET    /api/debt-lists/:id/installments/:num/notifications - Get notifications for specific installment

POST   /api/notifications/send                          - Send manual notification

POST   /api/debt-items/:id/verify                       - Verify payment (sends verification notification to payer)
POST   /api/debt-items/:id/reject                       - Reject payment (sends rejection notification to payer)

GET    /api/notifications                               - List all notifications
GET    /api/notifications/:id                           - Get notification details
PUT    /api/notifications/:id                           - Update notification
DELETE /api/notifications/:id                           - Delete notification
PATCH  /api/notifications/:id/enable                    - Enable notification
PATCH  /api/notifications/:id/disable                   - Disable notification

POST   /api/notification-templates                      - Create template
GET    /api/notification-templates                      - List templates
GET    /api/notification-templates/defaults             - Get default templates
PUT    /api/notification-templates/:id                  - Update template
DELETE /api/notification-templates/:id                  - Delete template

GET    /api/users/settings/notifications                - Get notification settings
PUT    /api/users/settings/notifications                - Update notification settings
```

---

## Dependencies to Install

```bash
# Core dependencies
go get github.com/go-co-op/gocron                    # Cron scheduler
go get gopkg.in/gomail.v2                            # SMTP email sending
go get github.com/twilio/twilio-go                   # Twilio SMS

# Webhook dependencies
go get github.com/go-telegram-bot-api/telegram-bot-api/v5  # Telegram Bot API
# Note: Slack and Discord use standard HTTP, no special SDK needed

# Utilities
go get github.com/hashicorp/go-retryablehttp         # HTTP client with retry
go get github.com/lib/pq                             # PostgreSQL array support
```

---

## Implementation Phases

### Phase 1: Database & Models ✓

1. Update `user_settings` table and model
2. Update `notification` table and model with installment fields
3. Create `notification_templates` table and model
4. Run database migrations

### Phase 2: Core Services

5. Implement email sender service (SMTP)
6. Implement SMS sender service (Twilio)
7. Implement webhook senders (Slack, Telegram, Discord)
8. Create template engine with installment variables
9. Set up go-cron scheduler service
10. Implement installment payment schedule calculator
11. Create notification scheduler logic

### Phase 3: Business Logic

12. Implement payment status checking
13. Create installment-based notification creation
14. Implement event-based notification service (payment confirmations)
15. Implement auto-disable on payment
16. Create contact info fetcher (for both user and contact)
17. Build notification worker/processor with webhook support
18. Add payment event hooks in debt item handler

### Phase 4: API Layer

19. Create notification repository
20. Create template repository
21. Implement notification handlers (with webhook & event support)
22. Implement template handlers
23. Add manual trigger endpoint
24. Add webhook configuration endpoints
25. Add event notification trigger endpoints

### Phase 5: Integration & UI

26. Configure environment variables (SMTP, Twilio, Telegram)
27. Integrate go-cron UI
28. Create initialization service
29. Add webhook for payment events
30. Webhook testing and validation
31. Event notification testing
32. End-to-end testing & refinement

---

## Key Files to Create

```
exit-debt/
├── internal/
│   ├── models/
│   │   ├── notification.go              ✓ Updated
│   │   ├── user_settings.go             ✓ Updated
│   │   └── notification_template.go     ○ New
│   ├── services/
│   │   ├── notification/
│   │   │   ├── service.go               ○ New
│   │   │   ├── email_sender.go          ○ New
│   │   │   ├── sms_sender.go            ○ New
│   │   │   ├── slack_sender.go          ○ New (Slack webhooks)
│   │   │   ├── telegram_sender.go       ○ New (Telegram bot API)
│   │   │   ├── discord_sender.go        ○ New (Discord webhooks)
│   │   │   ├── webhook_service.go       ○ New (unified webhook service)
│   │   │   ├── template_engine.go       ○ New
│   │   │   ├── scheduler.go             ○ New
│   │   │   ├── installment_scheduler.go ○ New (handles per-installment logic)
│   │   │   ├── payment_checker.go       ○ New (checks installment status)
│   │   │   └── event_notification_service.go ○ New (handles payment event notifications)
│   │   └── scheduler/
│   │       ├── cron_service.go          ○ New
│   │       └── job_manager.go           ○ New
│   ├── handlers/
│   │   ├── notification_handler.go      ○ New
│   │   └── template_handler.go          ○ New
│   ├── repositories/
│   │   ├── notification_repository.go   ○ New
│   │   └── template_repository.go       ○ New
│   └── config/
│       └── notification_config.go       ○ New
```

---

## Testing Checklist

### Installment Debt Creation

- [ ] Create debt with 12 monthly installments
- [ ] Verify 36 notifications created (12 × 3 reminders)
- [ ] Check all installment numbers are correct (1-12)
- [ ] Verify all due dates calculated correctly
- [ ] Confirm all cron jobs scheduled

### Notification Sending

- [ ] First installment reminder sends on time
- [ ] Message includes correct installment number
- [ ] Contact info fetched from user_contacts
- [ ] Payment status checked before sending
- [ ] Skip if installment already paid

### Payment Processing

- [ ] Pay installment #1
- [ ] Verify remaining notifications for #1 cancelled
- [ ] Verify notifications for #2 still active
- [ ] Check cron jobs removed
- [ ] Verify database updated correctly

### Edge Cases

- [ ] Pay multiple installments at once
- [ ] Pay full debt early (all installments)
- [ ] Partial payment for installment
- [ ] Change installment due date
- [ ] Delete debt (cascade delete notifications)

---

## Monitoring & Logging

```go
// Log key events
- INFO: Installment debt created, X notifications scheduled
- INFO: Notification sent for installment #Y
- INFO: Installment #Y paid, Z notifications cancelled
- WARN: Installment payment status unclear
- ERROR: Failed to fetch contact info
- ERROR: Notification failed after max retries
```

---

This implementation plan provides a complete, production-ready notification system with per-installment scheduling, smart payment tracking, and contact integration.

**Estimated Implementation Time:** 4-6 days

**Next Steps:** Begin with Phase 1 (Database & Models already partially complete) and progress through each phase.
