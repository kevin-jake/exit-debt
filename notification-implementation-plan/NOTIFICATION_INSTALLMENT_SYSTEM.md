# Installment Notification System - Complete Guide

## Overview

For installment debts (weekly, monthly, quarterly, etc.), **each individual installment payment** has its own notification schedule. This ensures users receive timely reminders for EVERY payment in their payment plan.

---

## Key Concept

### One Debt = Multiple Installments = Multiple Notification Schedules

**Example:**
```
Debt: $12,000 paid monthly over 12 months
├── Installment #1: $1,000 due 2025-01-15
│   ├── Reminder: 7 days before (2025-01-08)
│   ├── Reminder: 3 days before (2025-01-12)
│   └── Reminder: 1 day before (2025-01-14)
├── Installment #2: $1,000 due 2025-02-15
│   ├── Reminder: 7 days before (2025-02-08)
│   ├── Reminder: 3 days before (2025-02-12)
│   └── Reminder: 1 day before (2025-02-14)
├── ... (continues for all 12 installments)
└── Installment #12: $1,000 due 2025-12-15
    ├── Reminder: 7 days before (2025-12-08)
    ├── Reminder: 3 days before (2025-12-12)
    └── Reminder: 1 day before (2025-12-14)

TOTAL NOTIFICATIONS: 12 installments × 3 reminders = 36 scheduled notifications
```

---

## Database Schema

### Notification Table Fields

```sql
CREATE TABLE notifications (
    id UUID PRIMARY KEY,
    debt_list_id UUID NOT NULL REFERENCES debt_lists(id),
    
    -- NEW: Installment tracking fields
    installment_number INT,              -- Which payment: 1, 2, 3, etc.
    installment_due_date TIMESTAMP,      -- Due date of THIS specific installment
    
    notification_type VARCHAR(50),       -- email, sms
    message TEXT,
    status VARCHAR(50),                  -- pending, sent, failed
    sent_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    
    -- Contact info (optional override)
    recipient_email VARCHAR(255),
    recipient_phone VARCHAR(50)
);

CREATE INDEX idx_notifications_debt_installment 
ON notifications(debt_list_id, installment_number);

CREATE INDEX idx_notifications_installment_due 
ON notifications(installment_due_date) WHERE status = 'pending';
```

---

## How It Works

### Step 1: User Creates Installment Debt

```http
POST /api/debt-lists
{
  "contact_id": "...",
  "debt_type": "i_owe",
  "total_amount": "12000.00",
  "installment_plan": "monthly",     // Key: installment frequency
  "number_of_payments": 12,          // Key: total installments
  "due_date": "2025-01-15"           // First payment date
}
```

### Step 2: System Calculates Payment Schedule

System automatically generates `PaymentScheduleItem` for each installment:

```go
paymentSchedule := []PaymentScheduleItem{
    {
        PaymentNumber: 1,
        DueDate: time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
        ScheduledAmount: decimal.NewFromInt(1000),
        Status: "pending"
    },
    {
        PaymentNumber: 2,
        DueDate: time.Date(2025, 2, 15, 0, 0, 0, 0, time.UTC),
        ScheduledAmount: decimal.NewFromInt(1000),
        Status: "pending"
    },
    // ... continues for all 12 installments
}
```

### Step 3: System Creates Notifications for Each Installment

For **EACH** installment in the payment schedule:

```go
func CreateNotificationsForInstallmentDebt(debtListID uuid.UUID, paymentSchedule []PaymentScheduleItem, userSettings UserSettings) {
    for _, installment := range paymentSchedule {
        // Get reminder days from user settings (default: [7, 3, 1])
        reminderDays := userSettings.NotificationReminderDays
        
        for _, daysBefore := range reminderDays {
            // Calculate when to send reminder
            scheduledFor := installment.DueDate.AddDate(0, 0, -daysBefore).
                Add(parseTime(userSettings.NotificationTime))
            
            // Create notification
            notification := Notification{
                ID:                 uuid.New(),
                DebtListID:         debtListID,
                InstallmentNumber:  &installment.PaymentNumber,      // Track which installment
                InstallmentDueDate: &installment.DueDate,             // Track installment due date
                NotificationType:   "email",
                ScheduleType:       "reminder",
                ScheduledFor:       scheduledFor,
                Status:             "pending",
                Enabled:            true,
            }
            
            // Save to database
            db.Create(&notification)
            
            // Schedule cron job
            cronService.ScheduleAt(
                notification.ID.String(),
                scheduledFor,
                func() { processNotification(notification.ID) }
            )
        }
    }
}
```

**Result:**
- 12 installments × 3 reminders = **36 total notifications** created and scheduled

---

## Payment Tracking & Smart Notifications

### Only Send for Unpaid Installments

Before sending any notification, the system checks:

```go
func ShouldSendNotificationForInstallment(notification Notification) bool {
    // Step 1: Get the specific installment info
    installment := getInstallmentInfo(
        notification.DebtListID, 
        *notification.InstallmentNumber
    )
    
    // Step 2: Check if this installment is already paid
    if installment.Status == "paid" {
        log.Printf("Installment #%d already paid, skipping notification", 
            *notification.InstallmentNumber)
        return false
    }
    
    // Step 3: Check if amount is settled
    if installment.PaidAmount >= installment.ScheduledAmount {
        log.Printf("Installment #%d fully paid (amount: %v >= %v), skipping",
            *notification.InstallmentNumber, 
            installment.PaidAmount,
            installment.ScheduledAmount)
        return false
    }
    
    // Step 4: Check overall debt status
    debtList := getDebtList(notification.DebtListID)
    if debtList.Status == "settled" {
        log.Println("Entire debt settled, skipping all notifications")
        return false
    }
    
    return true // OK to send
}
```

### Auto-Disable When Installment Paid

When a user makes a payment for a specific installment:

```go
func OnInstallmentPaid(debtListID uuid.UUID, installmentNumber int) {
    // Disable ALL remaining notifications for this specific installment
    db.Exec(`
        UPDATE notifications
        SET enabled = false, status = 'cancelled'
        WHERE debt_list_id = $1
          AND installment_number = $2
          AND status = 'pending'
    `, debtListID, installmentNumber)
    
    // Remove associated cron jobs
    notifications := getNotificationsByInstallment(debtListID, installmentNumber)
    for _, notif := range notifications {
        if notif.CronJobID != nil {
            cronService.Remove(*notif.CronJobID)
        }
    }
    
    log.Printf("Disabled notifications for installment #%d", installmentNumber)
}
```

---

## Complete Example Scenario

### Initial Setup

```
User creates debt:
- Total: $12,000
- Plan: Monthly
- Installments: 12
- Start date: 2025-01-15
- Reminder settings: [7, 3, 1] days before at 09:00

System creates 36 notifications:
```

| Installment | Due Date   | Notification | Scheduled For     | Status  |
|-------------|------------|--------------|-------------------|---------|
| 1           | 2025-01-15 | 7d reminder  | 2025-01-08 09:00  | pending |
| 1           | 2025-01-15 | 3d reminder  | 2025-01-12 09:00  | pending |
| 1           | 2025-01-15 | 1d reminder  | 2025-01-14 09:00  | pending |
| 2           | 2025-02-15 | 7d reminder  | 2025-02-08 09:00  | pending |
| 2           | 2025-02-15 | 3d reminder  | 2025-02-12 09:00  | pending |
| 2           | 2025-02-15 | 1d reminder  | 2025-02-14 09:00  | pending |
| ...         | ...        | ...          | ...               | ...     |
| 12          | 2025-12-15 | 7d reminder  | 2025-12-08 09:00  | pending |
| 12          | 2025-12-15 | 3d reminder  | 2025-12-12 09:00  | pending |
| 12          | 2025-12-15 | 1d reminder  | 2025-12-14 09:00  | pending |

### Timeline

**2025-01-08 09:00**: First notification sent
```
✉️ EMAIL SENT: 
"Reminder: Payment #1 of $1,000 due to John Doe in 7 days (Jan 15, 2025)"
```

**2025-01-10**: User pays first installment early

```go
// System detects payment
OnInstallmentPaid(debtListID, 1)

// Cancels remaining notifications for installment #1
- ❌ 2025-01-12 09:00 (3d reminder) - CANCELLED
- ❌ 2025-01-14 09:00 (1d reminder) - CANCELLED

// Installment #2 notifications remain active
- ✅ 2025-02-08 09:00 (7d reminder) - ACTIVE
- ✅ 2025-02-12 09:00 (3d reminder) - ACTIVE
- ✅ 2025-02-14 09:00 (1d reminder) - ACTIVE
```

**2025-02-08 09:00**: Second installment reminders begin
```
✉️ EMAIL SENT:
"Reminder: Payment #2 of $1,000 due to John Doe in 7 days (Feb 15, 2025)"
```

**2025-02-12 09:00**: 
```
✉️ EMAIL SENT:
"Reminder: Payment #2 of $1,000 due to John Doe in 3 days (Feb 15, 2025)"
```

**Continues until all installments are paid...**

---

## API Endpoints

### Schedule Notifications for Installment Debt

```http
POST /api/debt-lists/:debt_id/schedule-installment-notifications

Response:
{
  "success": true,
  "data": {
    "debt_list_id": "...",
    "total_installments": 12,
    "notifications_created": 36,
    "breakdown": [
      {
        "installment_number": 1,
        "due_date": "2025-01-15",
        "notifications": 3
      },
      {
        "installment_number": 2,
        "due_date": "2025-02-15",
        "notifications": 3
      },
      // ... 12 installments
    ]
  }
}
```

### Get Notifications for Specific Installment

```http
GET /api/debt-lists/:debt_id/installments/:number/notifications

Response:
{
  "success": true,
  "data": {
    "installment_number": 1,
    "due_date": "2025-01-15",
    "installment_status": "pending",
    "notifications": [
      {
        "id": "...",
        "reminder_days_before": 7,
        "scheduled_for": "2025-01-08T09:00:00Z",
        "status": "sent",
        "sent_at": "2025-01-08T09:00:15Z"
      },
      {
        "id": "...",
        "reminder_days_before": 3,
        "scheduled_for": "2025-01-12T09:00:00Z",
        "status": "cancelled"
      },
      {
        "id": "...",
        "reminder_days_before": 1,
        "scheduled_for": "2025-01-14T09:00:00Z",
        "status": "cancelled"
      }
    ]
  }
}
```

### Get All Notifications for Debt (Grouped by Installment)

```http
GET /api/debt-lists/:debt_id/notifications/grouped

Response:
{
  "success": true,
  "data": {
    "debt_list_id": "...",
    "total_installments": 12,
    "total_notifications": 36,
    "installments": [
      {
        "installment_number": 1,
        "due_date": "2025-01-15",
        "status": "paid",
        "notifications": [
          { "reminder_days": 7, "status": "sent" },
          { "reminder_days": 3, "status": "cancelled" },
          { "reminder_days": 1, "status": "cancelled" }
        ]
      },
      {
        "installment_number": 2,
        "due_date": "2025-02-15",
        "status": "pending",
        "notifications": [
          { "reminder_days": 7, "status": "pending" },
          { "reminder_days": 3, "status": "pending" },
          { "reminder_days": 1, "status": "pending" }
        ]
      }
    ]
  }
}
```

---

## Notification Template Variables for Installments

### Additional Variables for Installment Debts

```
{{installment_number}}        - Which payment (1, 2, 3, etc.)
{{installment_total}}         - Total number of installments (e.g., 12)
{{installment_due_date}}      - Due date of this specific installment
{{installment_amount}}        - Amount for this installment
{{remaining_installments}}    - How many payments left
{{payments_made_count}}       - How many installments already paid
```

### Example Templates

**Email Template:**
```html
<h2>Payment Reminder - Installment #{{installment_number}}</h2>

<p>Hi {{user_first_name}},</p>

<p>This is a reminder that payment #{{installment_number}} of {{installment_total}} 
is due to {{debtor_name}} in {{days_until_due}} days.</p>

<div class="details">
  <h3>Payment Details:</h3>
  <ul>
    <li><strong>Installment:</strong> #{{installment_number}} of {{installment_total}}</li>
    <li><strong>Amount:</strong> {{currency}} {{installment_amount}}</li>
    <li><strong>Due Date:</strong> {{installment_due_date}}</li>
    <li><strong>Remaining Installments:</strong> {{remaining_installments}}</li>
  </ul>
</div>

<p>Total debt remaining: {{currency}} {{remaining_debt}}</p>
```

**SMS Template:**
```
Payment #{{installment_number}}/{{installment_total}}: {{currency}}{{installment_amount}} 
due to {{debtor_name}} in {{days_until_due}} days ({{installment_due_date}}). 
{{remaining_installments}} payments left.
```

---

## Database Queries

### Get All Pending Notifications for Specific Installment

```sql
SELECT n.*, dl.status as debt_status
FROM notifications n
JOIN debt_lists dl ON n.debt_list_id = dl.id
WHERE n.debt_list_id = $1
  AND n.installment_number = $2
  AND n.status = 'pending'
  AND n.enabled = true
ORDER BY n.scheduled_for;
```

### Get Upcoming Notifications Across All Installments

```sql
SELECT 
    n.*,
    dl.installment_plan,
    dl.status as debt_status
FROM notifications n
JOIN debt_lists dl ON n.debt_list_id = dl.id
WHERE n.scheduled_for BETWEEN NOW() AND NOW() + INTERVAL '7 days'
  AND n.status = 'pending'
  AND n.enabled = true
  AND dl.status != 'settled'
ORDER BY n.installment_number, n.scheduled_for;
```

### Count Notifications per Installment Status

```sql
SELECT 
    n.installment_number,
    n.installment_due_date,
    COUNT(*) as total_notifications,
    SUM(CASE WHEN n.status = 'sent' THEN 1 ELSE 0 END) as sent_count,
    SUM(CASE WHEN n.status = 'pending' THEN 1 ELSE 0 END) as pending_count,
    SUM(CASE WHEN n.status = 'cancelled' THEN 1 ELSE 0 END) as cancelled_count
FROM notifications n
WHERE n.debt_list_id = $1
GROUP BY n.installment_number, n.installment_due_date
ORDER BY n.installment_number;
```

---

## Business Rules Summary

### ✅ DO:

1. **Create notifications for EVERY installment** when debt is created
2. **Track installment number** and due date in each notification
3. **Check installment payment status** before sending each notification
4. **Auto-cancel remaining notifications** when installment is paid
5. **Include installment context** in notification messages

### ❌ DON'T:

1. **Don't send notifications** for already-paid installments
2. **Don't create notifications** for one-time (non-installment) debts using this system
3. **Don't send notifications** if entire debt is settled
4. **Don't forget to update** `debt_list.next_payment_date` after each payment

---

## Benefits of This Approach

| Benefit | Description |
|---------|-------------|
| **Precise Tracking** | Each installment is independently managed |
| **No Missed Payments** | Users get reminders for every payment |
| **Flexible** | Can customize reminders per installment if needed |
| **Smart Cancellation** | Auto-stops notifications when installment paid |
| **Clear Reporting** | Easy to see which installments have notifications |
| **User-Friendly** | Users know exactly which payment is due |

---

## Testing Checklist

### Create Installment Debt
- [ ] Create monthly debt with 12 installments
- [ ] Verify 36 notifications created (12 × 3 reminders)
- [ ] Check all installment numbers (1-12) are set correctly
- [ ] Verify all due dates are correct

### Notification Sending
- [ ] First installment reminder sends correctly
- [ ] Message includes installment number
- [ ] Check payment status before sending
- [ ] Skip if installment already paid

### Payment Processing
- [ ] Pay installment #1
- [ ] Verify remaining notifications for #1 are cancelled
- [ ] Verify notifications for #2 still active
- [ ] Check cron jobs removed for cancelled notifications

### Edge Cases
- [ ] What if user pays multiple installments at once?
- [ ] What if user pays partial amount for installment?
- [ ] What if user pays full debt early?
- [ ] What if installment date is changed?

---

## Migration from Single-Debt Notifications

If you have existing notifications without installment tracking:

```sql
-- For non-installment debts, set installment_number = 1
UPDATE notifications n
SET installment_number = 1,
    installment_due_date = (
        SELECT dl.next_payment_date 
        FROM debt_lists dl 
        WHERE dl.id = n.debt_list_id
    )
WHERE n.installment_number IS NULL
  AND EXISTS (
      SELECT 1 FROM debt_lists dl 
      WHERE dl.id = n.debt_list_id 
      AND dl.installment_plan = 'onetime'
  );
```

---

This completes the installment notification system design. Each installment is treated as an independent payment with its own notification schedule, ensuring users never miss a payment reminder.

