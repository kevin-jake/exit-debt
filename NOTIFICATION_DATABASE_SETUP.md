# Notification System - Database Setup Guide

**Clean Database Approach - No Migrations Required**

---

## 🎯 Overview

Since you're starting with a **clean database**, all tables will be created automatically using **GORM AutoMigrate**. No SQL migration files are needed.

---

## 📋 Setup Steps

### Step 1: Update GORM Models

Update your existing models and create new ones with all notification fields.

#### 1.1 Update `user_settings.go`

Add these fields to your existing `UserSettings` struct:

```go
// File: internal/models/user_settings.go
package models

import (
    "time"
    "github.com/google/uuid"
    "github.com/lib/pq"
)

type UserSettings struct {
    ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    UserID      uuid.UUID `json:"user_id" gorm:"type:uuid;not null;uniqueIndex"`

    // Existing fields
    NotificationEmail   bool   `json:"notification_email" gorm:"default:true"`
    NotificationSMS     bool   `json:"notification_sms" gorm:"default:false"`
    DefaultCurrency     string `json:"default_currency" gorm:"default:'Php'"`
    Timezone            string `json:"timezone" gorm:"default:'UTC'"`

    // NEW: Notification schedule settings
    NotificationReminderDays   pq.Int64Array `json:"notification_reminder_days" gorm:"type:integer[];default:'{7,3,1}'"`
    NotificationTime           string        `json:"notification_time" gorm:"default:'09:00:00'"`
    OverdueReminderFrequency   string        `json:"overdue_reminder_frequency" gorm:"default:'daily'"`
    CustomEmailMessage         *string       `json:"custom_email_message"`
    CustomSMSMessage           *string       `json:"custom_sms_message"`

    // NEW: Webhook notification settings
    NotificationWebhook        bool          `json:"notification_webhook" gorm:"default:false"`
    SlackWebhookURL            *string       `json:"slack_webhook_url"`
    TelegramBotToken           *string       `json:"telegram_bot_token"`
    TelegramChatID             *string       `json:"telegram_chat_id"`
    DiscordWebhookURL          *string       `json:"discord_webhook_url"`

    // NEW: Event notification settings
    EventNotificationsEnabled  bool          `json:"event_notifications_enabled" gorm:"default:true"`
    NotifyContactOnPayment     bool          `json:"notify_contact_on_payment" gorm:"default:true"`

    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type UpdateUserSettingsRequest struct {
    NotificationEmail            *bool         `json:"notification_email"`
    NotificationSMS              *bool         `json:"notification_sms"`
    NotificationWebhook          *bool         `json:"notification_webhook"`
    NotificationReminderDays     *pq.Int64Array `json:"notification_reminder_days"`
    NotificationTime             *string       `json:"notification_time"`
    OverdueReminderFrequency     *string       `json:"overdue_reminder_frequency"`
    CustomEmailMessage           *string       `json:"custom_email_message"`
    CustomSMSMessage             *string       `json:"custom_sms_message"`
    SlackWebhookURL              *string       `json:"slack_webhook_url"`
    TelegramBotToken             *string       `json:"telegram_bot_token"`
    TelegramChatID               *string       `json:"telegram_chat_id"`
    DiscordWebhookURL            *string       `json:"discord_webhook_url"`
    EventNotificationsEnabled    *bool         `json:"event_notifications_enabled"`
    NotifyContactOnPayment       *bool         `json:"notify_contact_on_payment"`
    DefaultCurrency              *string       `json:"default_currency"`
    Timezone                     *string       `json:"timezone"`
}
```

#### 1.2 Update `notification.go`

Add all new fields to your existing `Notification` struct:

```go
// File: internal/models/notification.go
package models

import (
    "time"
    "github.com/google/uuid"
    "github.com/lib/pq"
)

type Notification struct {
    ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`

    // Core relationships
    DebtListID  uuid.UUID  `json:"debt_list_id" gorm:"type:uuid;not null;index"`
    DebtItemID  *uuid.UUID `json:"debt_item_id" gorm:"type:uuid;index"` // For event-based notifications

    // NEW: Installment tracking
    InstallmentNumber    *int       `json:"installment_number" gorm:"index:idx_debt_installment"`
    InstallmentDueDate   *time.Time `json:"installment_due_date" gorm:"index:idx_installment_due"`

    // Notification settings
    NotificationType     string     `json:"notification_type" gorm:"not null"` // 'email', 'sms', 'webhook'
    WebhookType          *string    `json:"webhook_type"` // 'slack', 'telegram', 'discord'
    RecipientType        string     `json:"recipient_type" gorm:"not null"` // 'user' or 'contact'

    Message              string     `json:"message" gorm:"not null"`
    Status               string     `json:"status" gorm:"default:'pending';index"` // 'pending', 'sent', 'failed'
    SentAt               *time.Time `json:"sent_at"`

    // NEW: Scheduling fields
    ScheduleType         string     `json:"schedule_type"` // 'reminder', 'overdue', 'manual', 'event'
    ScheduledFor         *time.Time `json:"scheduled_for" gorm:"index:idx_scheduled_for"`
    CronJobID            *string    `json:"cron_job_id"`
    ReminderDaysBefore   *int       `json:"reminder_days_before"`

    // NEW: Custom schedule
    UseCustomSchedule    bool           `json:"use_custom_schedule" gorm:"default:false"`
    CustomReminderDays   pq.Int64Array  `json:"custom_reminder_days" gorm:"type:integer[]"`
    CustomNotificationTime *string      `json:"custom_notification_time"`
    CustomMessage        *string        `json:"custom_message"`

    // NEW: Recurring
    LastSentAt           *time.Time `json:"last_sent_at"`
    NextRunAt            *time.Time `json:"next_run_at" gorm:"index:idx_next_run"`
    IsRecurring          bool       `json:"is_recurring" gorm:"default:false"`
    Enabled              bool       `json:"enabled" gorm:"default:true"`

    // Optional recipient override
    RecipientEmail       *string    `json:"recipient_email,omitempty"`
    RecipientPhone       *string    `json:"recipient_phone,omitempty"`

    CreatedAt            time.Time  `json:"created_at"`
    UpdatedAt            time.Time  `json:"updated_at"`
}

type CreateNotificationRequest struct {
    DebtListID           uuid.UUID  `json:"debt_list_id" binding:"required"`
    InstallmentNumber    *int       `json:"installment_number"`
    InstallmentDueDate   *time.Time `json:"installment_due_date"`
    NotificationType     string     `json:"notification_type" binding:"required,oneof=email sms webhook"`
    WebhookType          *string    `json:"webhook_type" binding:"omitempty,oneof=slack telegram discord"`
    RecipientType        string     `json:"recipient_type" binding:"required,oneof=user contact"`
    Message              string     `json:"message" binding:"required"`
    RecipientEmail       *string    `json:"recipient_email,omitempty"`
    RecipientPhone       *string    `json:"recipient_phone,omitempty"`
}
```

#### 1.3 Create `notification_template.go`

Create a new model file:

```go
// File: internal/models/notification_template.go
package models

import (
    "time"
    "github.com/google/uuid"
    "github.com/lib/pq"
)

type NotificationTemplate struct {
    ID           uuid.UUID     `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    UserID       *uuid.UUID    `json:"user_id" gorm:"type:uuid;index"` // nil for system defaults
    TemplateName string        `json:"template_name" gorm:"not null"`
    TemplateType string        `json:"template_type" gorm:"not null;index"` // 'email' or 'sms'
    Subject      *string       `json:"subject"` // For email only
    Body         string        `json:"body" gorm:"not null"`
    IsDefault    bool          `json:"is_default" gorm:"default:false"`
    Variables    pq.StringArray `json:"variables" gorm:"type:text[]"` // Available template variables
    CreatedAt    time.Time     `json:"created_at"`
    UpdatedAt    time.Time     `json:"updated_at"`
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
```

### Step 2: Run AutoMigrate

Add AutoMigrate to your application initialization (usually in `main.go` or a database init function):

```go
// File: main.go or internal/database/init.go
package main

import (
    "log"
    "your-project/internal/models"
    "gorm.io/gorm"
)

func initDatabase(db *gorm.DB) error {
    // AutoMigrate will create tables with all fields
    err := db.AutoMigrate(
        &models.UserSettings{},
        &models.Notification{},
        &models.NotificationTemplate{},
        // ... other models
    )

    if err != nil {
        log.Printf("Error running AutoMigrate: %v", err)
        return err
    }

    log.Println("✅ Database tables created successfully")
    return nil
}
```

### Step 3: Add Indexes Manually (Optional)

GORM will create basic indexes from the `gorm:"index"` tags, but for composite indexes, you can add them manually:

```go
// Add after AutoMigrate in your init function
func addCustomIndexes(db *gorm.DB) error {
    // Composite index for notifications
    if err := db.Exec(`
        CREATE INDEX IF NOT EXISTS idx_notifications_debt_installment
        ON notifications(debt_list_id, installment_number)
    `).Error; err != nil {
        return err
    }

    // Index for scheduled pending notifications
    if err := db.Exec(`
        CREATE INDEX IF NOT EXISTS idx_notifications_scheduled_pending
        ON notifications(scheduled_for)
        WHERE status = 'pending'
    `).Error; err != nil {
        return err
    }

    // Index for enabled next run
    if err := db.Exec(`
        CREATE INDEX IF NOT EXISTS idx_notifications_next_run_enabled
        ON notifications(next_run_at)
        WHERE enabled = true
    `).Error; err != nil {
        return err
    }

    log.Println("✅ Custom indexes created successfully")
    return nil
}
```

---

## 🧪 Verify Database Setup

After running your application with AutoMigrate:

```bash
# 1. Check if tables were created
psql your_database -c "\dt"

# Expected tables:
# - user_settings (with new notification fields)
# - notifications (with installment and webhook fields)
# - notification_templates (new table)

# 2. Check user_settings structure
psql your_database -c "\d user_settings"

# Should show:
# - notification_reminder_days (integer[])
# - notification_time (text/varchar)
# - slack_webhook_url (text/varchar)
# - telegram_bot_token (text/varchar)
# - etc.

# 3. Check notifications structure
psql your_database -c "\d notifications"

# Should show:
# - installment_number (integer)
# - installment_due_date (timestamp)
# - webhook_type (text/varchar)
# - recipient_type (text/varchar)
# - schedule_type (text/varchar)
# - etc.

# 4. Check notification_templates structure
psql your_database -c "\d notification_templates"

# Should show complete structure

# 5. Verify indexes
psql your_database -c "\d+ notifications"
```

---

## ✅ Success Checklist

Phase 1 Complete When:

- [ ] `user_settings.go` updated with all notification fields
- [ ] `notification.go` updated with installment, webhook, and event fields
- [ ] `notification_template.go` created
- [ ] AutoMigrate added to initialization code
- [ ] Application starts without errors
- [ ] Tables created in database
- [ ] Can query new fields successfully
- [ ] Indexes created (basic ones automatically, custom ones if added)

---

## 🚀 Next Steps

Once Phase 1 is complete:

1. **Phase 2:** Implement email, SMS, and webhook sender services

   - See: `NOTIFICATION_QUICK_START.md`
   - See: `NOTIFICATION_WEBHOOK_GUIDE.md`

2. **Phase 3:** Implement business logic (payment checking, scheduling)
   - See: `NOTIFICATION_INSTALLMENT_SYSTEM.md`
   - See: `NOTIFICATION_BUSINESS_RULES.md`
   - See: `NOTIFICATION_EVENT_BASED_GUIDE.md`

---

## 📚 Reference

### Complete Model References

For complete model definitions with all fields, see:

- `NOTIFICATION_IMPLEMENTATION_PLAN.md` (Lines 112-217)
- `internal/models/user_settings.go` (already updated)
- `internal/models/notification.go` (already updated)

### GORM Documentation

- [GORM AutoMigrate](https://gorm.io/docs/migration.html)
- [GORM Indexes](https://gorm.io/docs/indexes.html)
- [GORM Data Types](https://gorm.io/docs/data_types.html)

---

## 💡 Tips

1. **Import pq for PostgreSQL arrays:**

   ```go
   import "github.com/lib/pq"
   ```

2. **Test with sample data:**

   ```go
   settings := models.UserSettings{
       NotificationEmail: true,
       NotificationReminderDays: pq.Int64Array{7, 3, 1},
       NotificationTime: "09:00:00",
   }
   db.Create(&settings)
   ```

3. **Check for compilation errors:**

   ```bash
   go build ./...
   ```

4. **Run application to trigger AutoMigrate:**
   ```bash
   go run main.go
   ```

---

**Status:** Ready to implement Phase 1! 🎯

**No migrations needed** - GORM will handle everything automatically.
