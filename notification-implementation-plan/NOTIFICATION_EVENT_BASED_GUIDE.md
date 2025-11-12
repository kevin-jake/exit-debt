# Event-Based Notifications Guide

## Overview

In addition to **scheduled reminders**, the notification system sends **instant notifications** when certain events occur:

### Event Types

1. **Payment Made** - When a debt item (payment) is recorded
2. **Debt Item Added** - When a new payment transaction is added to a debt

### Recipients

For each event, **TWO notifications** are sent:

1. ✅ **To the User** - The person who owns/created the debt
2. ✅ **To the Contact** - The other party involved in the debt

---

## Architecture

### Notification Types

```
Schedule Type:
  - 'reminder'      → Time-based, scheduled for future
  - 'overdue'       → Time-based, for overdue payments
  - 'manual'        → User-triggered
  - 'event'         → NEW! Event-triggered (payment made, debt item added)

Recipient Type:
  - 'user'          → Sent to the debt owner
  - 'contact'       → Sent to the contact involved
```

### Database Schema Updates

```sql
-- Add event-based notification fields
ALTER TABLE notifications ADD COLUMN debt_item_id UUID REFERENCES debt_items(id) ON DELETE CASCADE;
ALTER TABLE notifications ADD COLUMN recipient_type VARCHAR(50);  -- 'user' or 'contact'

-- Add index for debt items
CREATE INDEX idx_notifications_debt_item ON notifications(debt_item_id);
CREATE INDEX idx_notifications_recipient_type ON notifications(recipient_type);

-- Add constraint
ALTER TABLE notifications ADD CONSTRAINT check_recipient_type
    CHECK (recipient_type IN ('user', 'contact'));
```

### Go Model Updates

```go
type Notification struct {
    ID                   uuid.UUID  `json:"id"`
    DebtListID           uuid.UUID  `json:"debt_list_id"`
    DebtItemID           *uuid.UUID `json:"debt_item_id"`  // NEW: For payment confirmations
    InstallmentNumber    *int       `json:"installment_number"`
    InstallmentDueDate   *time.Time `json:"installment_due_date"`
    NotificationType     string     `json:"notification_type"`
    WebhookType          *string    `json:"webhook_type"`
    RecipientType        string     `json:"recipient_type"`  // NEW: 'user' or 'contact'
    Message              string     `json:"message"`
    Status               string     `json:"status"`
    ScheduleType         string     `json:"schedule_type"`  // 'reminder', 'overdue', 'manual', 'event'
    // ... other fields
}
```

---

## Event Flow

### Scenario 1: Payment Made on "I Owe" Debt

```
User Kevin owes $1,000 to John Doe (contact)
Kevin makes a $500 payment
```

**System Actions:**

1. Create debt item record (status: "pending")
2. Update debt list (remaining balance, etc.)
3. **Create 2 event notifications:**
   - Notification #1 (to User - Kevin):
     - `recipient_type: 'user'`
     - `schedule_type: 'event'`
     - Message: "Payment recorded: You paid $500 to John Doe (pending verification)"
   - Notification #2 (to Contact - John):
     - `recipient_type: 'contact'`
     - `schedule_type: 'event'`
     - Message: "Payment received: Kevin paid you $500 and still under verification"
4. Send both notifications immediately
5. Check if installment paid → disable remaining reminders

### Scenario 2: Payment Made on "Owed to Me" Debt

```
John Doe owes User Kevin $1,000
John makes a $500 payment
```

**System Actions:**

1. Create debt item record (status: "pending")
2. Update debt list
3. **Create 2 event notifications:**
   - Notification #1 (to User - Kevin):
     - `recipient_type: 'user'`
     - `schedule_type: 'event'`
     - Message: "Payment received: John Doe paid you $500 and still under verification"
   - Notification #2 (to Contact - John):
     - `recipient_type: 'contact'`
     - `schedule_type: 'event'`
     - Message: "Payment recorded: You paid $500 to Kevin (pending verification)"
4. Send both notifications immediately
5. Check if installment paid → disable remaining reminders

### Scenario 3: Payment Verified

```
User Kevin verifies John Doe's $500 payment
Payment status changes from "pending" to "completed"
```

**System Actions:**

1. Update debt item status to "completed"
2. Update debt list calculations (if needed)
3. **Create 1 event notification:**
   - Notification (to Payer - John):
     - `recipient_type: 'contact'` or `'user'` (depending on who made payment)
     - `schedule_type: 'event'`
     - Message: "Payment verified: Your payment of $500 to Kevin has been verified and accepted"
4. Send notification immediately

### Scenario 4: Payment Rejected

```
User Kevin rejects John Doe's $500 payment
Payment status changes from "pending" to "rejected"
```

**System Actions:**

1. Update debt item status to "rejected"
2. Reverse debt list calculations
3. **Create 1 event notification:**
   - Notification (to Payer - John):
     - `recipient_type: 'contact'` or `'user'`
     - `schedule_type: 'event'`
     - Message: "Payment rejected: Your payment of $500 to Kevin has been rejected. Reason: [rejection reason]"
4. Send notification immediately

---

## Implementation

### 1. Event Hook on Debt Item Creation

```go
// internal/handlers/debt_item_handler.go
func (h *DebtItemHandler) CreateDebtItem(c *gin.Context) {
    var req CreateDebtItemRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // Step 1: Create debt item (status: "pending")
    debtItem, err := h.debtItemService.Create(req)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    // Step 2: Update debt list calculations
    h.debtListService.RecalculateDebt(req.DebtListID)

    // Step 3: Send event notifications to both user and contact (with "pending verification" message)
    go h.notificationService.SendPaymentNotifications(debtItem.ID, req.DebtListID)

    // Step 4: Check if installment is now paid, disable scheduled reminders
    go h.notificationService.CheckAndDisableReminders(req.DebtListID, debtItem)

    c.JSON(201, debtItem)
}

// VerifyPayment handles payment verification
func (h *DebtItemHandler) VerifyPayment(c *gin.Context) {
    debtItemID := c.Param("id")

    var req VerifyPaymentRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // Step 1: Update debt item status to "completed"
    debtItem, err := h.debtItemService.VerifyPayment(debtItemID)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    // Step 2: Send verification notification to the payer
    go h.notificationService.SendPaymentVerificationNotification(debtItem.ID, debtItem.DebtListID, true, nil)

    c.JSON(200, gin.H{"message": "Payment verified successfully", "debt_item": debtItem})
}

// RejectPayment handles payment rejection
func (h *DebtItemHandler) RejectPayment(c *gin.Context) {
    debtItemID := c.Param("id")

    var req RejectPaymentRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // Step 1: Update debt item status to "rejected"
    debtItem, err := h.debtItemService.RejectPayment(debtItemID, req.Reason)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    // Step 2: Send rejection notification to the payer
    go h.notificationService.SendPaymentVerificationNotification(debtItem.ID, debtItem.DebtListID, false, &req.Reason)

    c.JSON(200, gin.H{"message": "Payment rejected", "debt_item": debtItem})
}

type VerifyPaymentRequest struct {
    // Empty for now, can add notes later
}

type RejectPaymentRequest struct {
    Reason string `json:"reason" binding:"required"`
}
```

### 2. Notification Service - Send Payment Notifications

```go
// internal/services/notification/event_notification_service.go
package notification

import (
    "fmt"
    "log"
)

type EventNotificationService struct {
    emailSender    *EmailSender
    smsSender      *SMSSender
    webhookService *WebhookService
    db             *gorm.DB
}

func (s *EventNotificationService) SendPaymentNotifications(debtItemID, debtListID uuid.UUID) error {
    // Step 1: Fetch debt item, debt list, user, and contact information
    info, err := s.fetchPaymentInfo(debtItemID, debtListID)
    if err != nil {
        return fmt.Errorf("failed to fetch payment info: %w", err)
    }

    // Step 2: Create notification for USER
    userNotification := s.createUserPaymentNotification(info)
    s.db.Create(&userNotification)

    // Step 3: Create notification for CONTACT
    contactNotification := s.createContactPaymentNotification(info)
    s.db.Create(&contactNotification)

    // Step 4: Send both notifications immediately
    go s.sendNotification(userNotification)
    go s.sendNotification(contactNotification)

    return nil
}

func (s *EventNotificationService) fetchPaymentInfo(debtItemID, debtListID uuid.UUID) (*PaymentInfo, error) {
    var info PaymentInfo

    err := s.db.Raw(`
        SELECT
            di.id as debt_item_id,
            di.amount as payment_amount,
            di.payment_date,
            di.payment_method,
            dl.id as debt_list_id,
            dl.debt_type,
            dl.currency,
            dl.total_remaining_debt,
            u.id as user_id,
            u.first_name as user_first_name,
            u.last_name as user_last_name,
            u.email as user_email,
            u.phone as user_phone,
            uc.id as contact_id,
            uc.name as contact_name,
            uc.email as contact_email,
            uc.phone as contact_phone,
            us.notification_email as user_notification_email,
            us.notification_sms as user_notification_sms,
            us.notification_webhook as user_notification_webhook,
            us.slack_webhook_url,
            us.telegram_bot_token,
            us.telegram_chat_id,
            us.discord_webhook_url
        FROM debt_items di
        JOIN debt_lists dl ON di.debt_list_id = dl.id
        JOIN users u ON dl.user_id = u.id
        JOIN user_contacts uc ON uc.contact_id = dl.contact_id AND uc.user_id = dl.user_id
        JOIN user_settings us ON us.user_id = u.id
        WHERE di.id = ? AND dl.id = ?
    `, debtItemID, debtListID).Scan(&info).Error

    if err != nil {
        return nil, err
    }

    return &info, nil
}

func (s *EventNotificationService) createUserPaymentNotification(info *PaymentInfo) Notification {
    // Determine message based on debt type
    var message string
    if info.DebtType == "i_owe" {
        message = fmt.Sprintf("Payment recorded: You paid %s %s to %s (pending verification)",
            info.Currency, info.PaymentAmount.String(), info.ContactName)
    } else {
        message = fmt.Sprintf("Payment received: %s paid you %s %s and still under verification",
            info.ContactName, info.Currency, info.PaymentAmount.String())
    }

    return Notification{
        ID:               uuid.New(),
        DebtListID:       info.DebtListID,
        DebtItemID:       &info.DebtItemID,
        NotificationType: "email", // Or based on user settings
        RecipientType:    "user",
        ScheduleType:     "event",
        Message:          message,
        Status:           "pending",
        RecipientEmail:   &info.UserEmail,
        RecipientPhone:   &info.UserPhone,
        CreatedAt:        time.Now(),
    }
}

func (s *EventNotificationService) createContactPaymentNotification(info *PaymentInfo) Notification {
    // Determine message based on debt type (opposite of user)
    var message string
    if info.DebtType == "i_owe" {
        message = fmt.Sprintf("Payment received: %s paid you %s %s and still under verification",
            info.UserFirstName, info.Currency, info.PaymentAmount.String())
    } else {
        message = fmt.Sprintf("Payment recorded: You paid %s %s to %s (pending verification)",
            info.Currency, info.PaymentAmount.String(),
            fmt.Sprintf("%s %s", info.UserFirstName, info.UserLastName))
    }

    return Notification{
        ID:               uuid.New(),
        DebtListID:       info.DebtListID,
        DebtItemID:       &info.DebtItemID,
        NotificationType: "email", // Or based on contact preferences
        RecipientType:    "contact",
        ScheduleType:     "event",
        Message:          message,
        Status:           "pending",
        RecipientEmail:   &info.ContactEmail,
        RecipientPhone:   &info.ContactPhone,
        CreatedAt:        time.Now(),
    }
}

// SendPaymentVerificationNotification sends a notification when payment status changes (verified/rejected)
func (s *EventNotificationService) SendPaymentVerificationNotification(debtItemID, debtListID uuid.UUID, verified bool, rejectionReason *string) error {
    // Step 1: Fetch payment info
    info, err := s.fetchPaymentInfo(debtItemID, debtListID)
    if err != nil {
        return fmt.Errorf("failed to fetch payment info: %w", err)
    }

    // Step 2: Determine who made the payment (they get the verification notification)
    // If debt_type is "i_owe", the USER made the payment
    // If debt_type is "owed_to_me", the CONTACT made the payment
    var notification Notification
    var recipientType string
    var recipientEmail, recipientPhone *string

    if info.DebtType == "i_owe" {
        // User made payment, notify user about verification
        recipientType = "user"
        recipientEmail = &info.UserEmail
        recipientPhone = &info.UserPhone
    } else {
        // Contact made payment, notify contact about verification
        recipientType = "contact"
        recipientEmail = &info.ContactEmail
        recipientPhone = &info.ContactPhone
    }

    // Step 3: Create verification message
    var message string
    if verified {
        if info.DebtType == "i_owe" {
            message = fmt.Sprintf("Payment verified: Your payment of %s %s to %s has been verified and accepted",
                info.Currency, info.PaymentAmount.String(), info.ContactName)
        } else {
            message = fmt.Sprintf("Payment verified: Your payment of %s %s to %s %s has been verified and accepted",
                info.Currency, info.PaymentAmount.String(), info.UserFirstName, info.UserLastName)
        }
    } else {
        reason := "No reason provided"
        if rejectionReason != nil {
            reason = *rejectionReason
        }
        if info.DebtType == "i_owe" {
            message = fmt.Sprintf("Payment rejected: Your payment of %s %s to %s has been rejected. Reason: %s",
                info.Currency, info.PaymentAmount.String(), info.ContactName, reason)
        } else {
            message = fmt.Sprintf("Payment rejected: Your payment of %s %s to %s %s has been rejected. Reason: %s",
                info.Currency, info.PaymentAmount.String(), info.UserFirstName, info.UserLastName, reason)
        }
    }

    notification = Notification{
        ID:               uuid.New(),
        DebtListID:       info.DebtListID,
        DebtItemID:       &info.DebtItemID,
        NotificationType: "email",
        RecipientType:    recipientType,
        ScheduleType:     "event",
        Message:          message,
        Status:           "pending",
        RecipientEmail:   recipientEmail,
        RecipientPhone:   recipientPhone,
        CreatedAt:        time.Now(),
    }

    s.db.Create(&notification)

    // Step 4: Send notification immediately
    go s.sendNotification(notification)

    return nil
}

func (s *EventNotificationService) sendNotification(notification Notification) {
    // Fetch user settings for channel preferences
    settings := s.getUserSettings(notification.DebtListID)

    // Prepare notification data
    data := s.prepareNotificationData(notification)

    var err error

    // Send based on recipient type and preferences
    if notification.RecipientType == "user" {
        // Send to user using their channel preferences
        if settings.NotificationEmail && notification.RecipientEmail != nil {
            err = s.emailSender.SendEmail(*notification.RecipientEmail, "Payment Notification", data.Message)
        }
        if settings.NotificationSMS && notification.RecipientPhone != nil {
            err = s.smsSender.SendSMS(*notification.RecipientPhone, data.Message)
        }
        if settings.NotificationWebhook {
            // Send to all configured webhooks
            if settings.SlackWebhookURL != nil {
                s.webhookService.SendNotification("slack", settings, data)
            }
            // ... telegram, discord
        }
    } else {
        // Send to contact (default to email if available)
        if notification.RecipientEmail != nil {
            err = s.emailSender.SendEmail(*notification.RecipientEmail, "Payment Notification", data.Message)
        }
        // Optionally send SMS if phone available
        if notification.RecipientPhone != nil {
            err = s.smsSender.SendSMS(*notification.RecipientPhone, data.Message)
        }
    }

    // Update notification status
    if err == nil {
        s.db.Model(&notification).Updates(map[string]interface{}{
            "status":  "sent",
            "sent_at": time.Now(),
        })
    } else {
        log.Printf("[ERROR] Failed to send notification: %v", err)
        s.db.Model(&notification).Update("status", "failed")
    }
}

type PaymentInfo struct {
    DebtItemID             uuid.UUID
    PaymentAmount          decimal.Decimal
    PaymentDate            time.Time
    PaymentMethod          string
    DebtListID             uuid.UUID
    DebtType               string
    Currency               string
    TotalRemainingDebt     decimal.Decimal
    UserID                 uuid.UUID
    UserFirstName          string
    UserLastName           string
    UserEmail              string
    UserPhone              *string
    ContactID              uuid.UUID
    ContactName            string
    ContactEmail           string
    ContactPhone           *string
    UserNotificationEmail  bool
    UserNotificationSMS    bool
    UserNotificationWebhook bool
    SlackWebhookURL        *string
    TelegramBotToken       *string
    TelegramChatID         *string
    DiscordWebhookURL      *string
}
```

---

## Notification Templates for Events

### Email Template - Payment Made (To User - "I Owe")

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
      .footer {
        text-align: center;
        padding: 20px;
        color: #666;
      }
    </style>
  </head>
  <body>
    <div class="container">
      <div class="header">
        <h2>✅ Payment Recorded</h2>
      </div>
      <div class="content">
        <p>Hi {{user_first_name}},</p>
        <p>Your payment has been successfully recorded.</p>

        <div class="details">
          <h3>Payment Details:</h3>
          <ul>
            <li>
              <strong>Amount Paid:</strong> {{currency}} {{payment_amount}}
            </li>
            <li><strong>To:</strong> {{contact_name}}</li>
            <li><strong>Payment Date:</strong> {{payment_date}}</li>
            <li><strong>Payment Method:</strong> {{payment_method}}</li>
            <li>
              <strong>Remaining Balance:</strong> {{currency}}
              {{remaining_debt}}
            </li>
          </ul>
        </div>

        <p>A notification has also been sent to {{contact_name}}.</p>
      </div>
      <div class="footer">
        <p>Exit Debt - Track your debts easily</p>
      </div>
    </div>
  </body>
</html>
```

### Email Template - Payment Received (To Contact)

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
        background-color: #2196f3;
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
        border-left: 4px solid #2196f3;
      }
      .footer {
        text-align: center;
        padding: 20px;
        color: #666;
      }
    </style>
  </head>
  <body>
    <div class="container">
      <div class="header">
        <h2>💰 Payment Received</h2>
      </div>
      <div class="content">
        <p>Hi {{contact_name}},</p>
        <p>
          You have received a payment from {{user_first_name}}
          {{user_last_name}}.
        </p>

        <div class="details">
          <h3>Payment Details:</h3>
          <ul>
            <li>
              <strong>Amount Received:</strong> {{currency}} {{payment_amount}}
            </li>
            <li>
              <strong>From:</strong> {{user_first_name}} {{user_last_name}}
            </li>
            <li><strong>Payment Date:</strong> {{payment_date}}</li>
            <li><strong>Payment Method:</strong> {{payment_method}}</li>
            <li>
              <strong>Remaining Balance:</strong> {{currency}}
              {{remaining_debt}}
            </li>
          </ul>
        </div>

        <p>This payment has been recorded in Exit Debt system.</p>
      </div>
      <div class="footer">
        <p>Exit Debt</p>
      </div>
    </div>
  </body>
</html>
```

### SMS Template - Payment Made

```
Payment recorded: You paid {{currency}}{{payment_amount}} to {{contact_name}}. Remaining: {{currency}}{{remaining_debt}}.
```

### SMS Template - Payment Received

```
Payment received: {{user_first_name}} paid you {{currency}}{{payment_amount}}. Remaining: {{currency}}{{remaining_debt}}.
```

### Slack Message - Payment Made

```json
{
  "blocks": [
    {
      "type": "header",
      "text": {
        "type": "plain_text",
        "text": "✅ Payment Recorded"
      }
    },
    {
      "type": "section",
      "fields": [
        {
          "type": "mrkdwn",
          "text": "*Amount Paid:*\n{{currency}} {{payment_amount}}"
        },
        {
          "type": "mrkdwn",
          "text": "*To:*\n{{contact_name}}"
        },
        {
          "type": "mrkdwn",
          "text": "*Date:*\n{{payment_date}}"
        },
        {
          "type": "mrkdwn",
          "text": "*Remaining:*\n{{currency}} {{remaining_debt}}"
        }
      ]
    }
  ]
}
```

---

## API Examples

### Manual Trigger Payment Notification

```http
POST /api/debt-items/:debt_item_id/notify

Response:
{
  "success": true,
  "message": "Payment notifications sent",
  "notifications_sent": {
    "user": {
      "id": "notification_uuid_1",
      "recipient_type": "user",
      "recipient_email": "user@example.com",
      "status": "sent",
      "sent_at": "2025-11-07T10:30:00Z"
    },
    "contact": {
      "id": "notification_uuid_2",
      "recipient_type": "contact",
      "recipient_email": "contact@example.com",
      "status": "sent",
      "sent_at": "2025-11-07T10:30:01Z"
    }
  }
}
```

### Verify Payment

```http
POST /api/debt-items/:debt_item_id/verify

Response:
{
  "message": "Payment verified successfully",
  "debt_item": {
    "id": "debt_item_uuid",
    "status": "completed",
    "amount": "500.00",
    "payment_date": "2025-11-07T10:00:00Z"
  }
}
```

**What happens:**

1. Debt item status changes from "pending" to "completed"
2. Notification sent to payer: "Payment verified: Your payment of Php 500 to Kevin has been verified and accepted"

### Reject Payment

```http
POST /api/debt-items/:debt_item_id/reject
Content-Type: application/json

{
  "reason": "Payment proof not clear, please resubmit"
}

Response:
{
  "message": "Payment rejected",
  "debt_item": {
    "id": "debt_item_uuid",
    "status": "rejected",
    "amount": "500.00",
    "rejection_reason": "Payment proof not clear, please resubmit"
  }
}
```

**What happens:**

1. Debt item status changes from "pending" to "rejected"
2. Notification sent to payer: "Payment rejected: Your payment of Php 500 to Kevin has been rejected. Reason: Payment proof not clear, please resubmit"

### Get Event Notifications for Debt

```http
GET /api/debt-lists/:id/notifications?schedule_type=event

Response:
{
  "success": true,
  "data": {
    "total": 6,
    "notifications": [
      {
        "id": "uuid-1",
        "debt_item_id": "payment_uuid_1",
        "recipient_type": "user",
        "schedule_type": "event",
        "message": "Payment recorded: You paid Php 500 to John Doe (pending verification)",
        "status": "sent",
        "sent_at": "2025-11-07T10:30:00Z"
      },
      {
        "id": "uuid-2",
        "debt_item_id": "payment_uuid_1",
        "recipient_type": "contact",
        "schedule_type": "event",
        "message": "Payment received: Kevin paid you Php 500 and still under verification",
        "status": "sent",
        "sent_at": "2025-11-07T10:30:01Z"
      },
      {
        "id": "uuid-3",
        "debt_item_id": "payment_uuid_1",
        "recipient_type": "user",
        "schedule_type": "event",
        "message": "Payment verified: Your payment of Php 500 to John Doe has been verified and accepted",
        "status": "sent",
        "sent_at": "2025-11-07T12:00:00Z"
      }
    ]
  }
}
```

---

## Configuration

### Enable/Disable Event Notifications

```http
PUT /api/users/settings/notifications
{
  "event_notifications_enabled": true,
  "notify_contact_on_payment": true,
  "notification_email": true,
  "notification_sms": false,
  "notification_webhook": true
}
```

### User Settings Model Update

```go
type UserSettings struct {
    // ... existing fields

    // Event notification settings
    EventNotificationsEnabled bool `json:"event_notifications_enabled" gorm:"default:true"`
    NotifyContactOnPayment    bool `json:"notify_contact_on_payment" gorm:"default:true"`
}
```

---

## Benefits

| Benefit            | Description                          |
| ------------------ | ------------------------------------ |
| **Transparency**   | Both parties informed instantly      |
| **Trust Building** | Automatic receipts build trust       |
| **Record Keeping** | Notification history as proof        |
| **Real-Time**      | No delays, instant confirmation      |
| **Multi-Channel**  | Email, SMS, Slack, Telegram, Discord |
| **Professional**   | Automated, consistent communication  |

---

## Complete Flow Example

### Complete Workflow: User Kevin pays John $500 (with Verification)

```
1. Kevin submits payment:
   POST /api/debt-items
   {
     "debt_list_id": "xxx",
     "amount": "500",
     "payment_method": "bank_transfer"
   }

2. System creates debt item (status: "pending")
3. System updates debt calculations
4. System creates 2 notifications:

   Notification #1 (to Kevin):
   ✉️ "Payment recorded: You paid Php 500 to John Doe (pending verification)"
   Recipient: kevin@example.com
   Channels: Email + Slack

   Notification #2 (to John):
   ✉️ "Payment received: Kevin paid you Php 500 and still under verification"
   Recipient: john@example.com
   Channels: Email

5. Both notifications sent immediately
6. Returns success response to Kevin

--- Later, John verifies the payment ---

7. John verifies payment:
   POST /api/debt-items/:debt_item_id/verify

8. System updates debt item status to "completed"
9. System creates 1 verification notification:

   Notification #3 (to Kevin):
   ✉️ "Payment verified: Your payment of Php 500 to John Doe has been verified and accepted"
   Recipient: kevin@example.com
   Channels: Email + Slack

10. Verification notification sent immediately
11. System checks if installment paid → disables future reminders if fully paid
12. Returns success response to John
```

### Alternative Flow: Payment Rejected

```
--- If John rejects the payment instead ---

7. John rejects payment:
   POST /api/debt-items/:debt_item_id/reject
   {
     "reason": "Payment proof not clear, please resubmit"
   }

8. System updates debt item status to "rejected"
9. System reverses debt calculations
10. System creates 1 rejection notification:

   Notification #3 (to Kevin):
   ✉️ "Payment rejected: Your payment of Php 500 to John Doe has been rejected. Reason: Payment proof not clear, please resubmit"
   Recipient: kevin@example.com
   Channels: Email + Slack

11. Rejection notification sent immediately
12. Returns response to John
```

---

## Testing Checklist

### Payment Event Notifications

- [ ] Create debt item
- [ ] Verify 2 notifications created (user + contact)
- [ ] Check user notification includes "(pending verification)" message
- [ ] Check contact notification includes "still under verification" message
- [ ] Verify correct message for "i_owe" debt type
- [ ] Verify correct message for "owed_to_me" debt type
- [ ] Test with webhook enabled (Slack, Telegram, Discord)
- [ ] Test with SMS enabled
- [ ] Verify both notifications marked as "sent"

### Verification/Rejection Notifications

- [ ] Verify payment → debt item status changes to "completed"
- [ ] Verification notification sent to payer only (not receiver)
- [ ] Verification message includes "verified and accepted"
- [ ] Reject payment → debt item status changes to "rejected"
- [ ] Rejection notification sent to payer only
- [ ] Rejection message includes reason provided
- [ ] Test verification for "i_owe" debt (user is payer)
- [ ] Test verification for "owed_to_me" debt (contact is payer)
- [ ] Verify total of 3 notifications: 2 initial + 1 verification/rejection

### Edge Cases

- [ ] Contact has no email - should skip contact notification
- [ ] User has webhooks enabled - should send to all channels
- [ ] Payment notification when debt fully paid
- [ ] Multiple payments in quick succession
- [ ] Notification preferences disabled
- [ ] Verify payment that's already verified
- [ ] Reject payment that's already completed
- [ ] Empty rejection reason

### Integration

- [ ] Event notification + scheduled reminder both work
- [ ] Payment triggers reminder cancellation
- [ ] Event notifications appear in notification history
- [ ] Can filter by schedule_type="event"
- [ ] Can filter by recipient_type="user" or "contact"
- [ ] Verification after payment disables scheduled reminders if installment paid

---

## Summary

Event-based notifications provide:

✅ **Instant confirmation** when payments are made  
✅ **Two-way communication** between user and contact  
✅ **Verification workflow** with pending → verified/rejected status  
✅ **Trust building** through transparency  
✅ **Automatic receipts** for both parties  
✅ **Multi-channel delivery** (Email, SMS, webhooks)  
✅ **Integration** with scheduled reminders

### Notification Flow Summary

1. **Payment Made** → 2 notifications (user + contact) with "pending verification" message
2. **Payment Verified** → 1 notification to payer with "verified and accepted" message
3. **Payment Rejected** → 1 notification to payer with rejection reason

This completes the notification system with both **time-based** (scheduled reminders) and **event-based** (payment confirmations with verification) notifications!
