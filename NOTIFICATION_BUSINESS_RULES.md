# Notification System - Business Rules

## Critical Rule: Only Send Notifications for Unpaid Debts

---

## The Rule

**💡 Before sending ANY notification, the system MUST check if the debt is still unpaid.**

```
IF debt_status == "settled" → SKIP notification + DISABLE future notifications
IF total_remaining_debt <= 0 → SKIP notification + DISABLE future notifications
ELSE → SEND notification
```

---

## Decision Flow

```
┌─────────────────────────────────────┐
│ Notification Scheduled to Send     │
└────────────────┬────────────────────┘
                 ▼
┌─────────────────────────────────────┐
│ Fetch Debt Information              │
│ - status                            │
│ - total_remaining_debt              │
│ - next_payment_date                 │
└────────────────┬────────────────────┘
                 ▼
         ┌───────┴───────┐
         │ Check Status  │
         └───────┬───────┘
                 │
    ┌────────────┼────────────┐
    ▼            ▼            ▼
┌─────────┐ ┌─────────┐ ┌──────────┐
│Status = │ │Remaining│ │All Good  │
│"settled"│ │Debt = 0 │ │Continue  │
└────┬────┘ └────┬────┘ └────┬─────┘
     │           │            │
     │    ┌──────┴──────┐     │
     │    ▼             ▼     │
     │  ┌─────────────────┐   │
     └─>│ SKIP & DISABLE  │   │
        │ Future Notifs   │   │
        └─────────────────┘   │
                              ▼
                    ┌──────────────────┐
                    │ Fetch Contact    │
                    │ Send Notification│
                    └──────────────────┘
```

---

## When to Check

### 1. Reminder Notifications (Before Due Date)

```go
// Check BEFORE sending reminder
if debtStatus.Status == "settled" || debtStatus.TotalRemainingDebt <= 0 {
    skipAndDisable(notification)
    return
}

// Additional check for reminders
if time.Now().After(debtStatus.NextPaymentDate) {
    skip(notification)  // Let overdue notification handle this
    return
}
```

### 2. Overdue Notifications (After Due Date)

```go
// Check BEFORE sending overdue notification
if debtStatus.Status == "settled" || debtStatus.TotalRemainingDebt <= 0 {
    skipAndDisable(notification)
    return
}

// Check 24-hour limit
if time.Since(notification.LastSentAt) < 24*time.Hour {
    skip(notification)
    return
}
```

### 3. Manual Notifications (API Triggered)

```go
// Even manual notifications respect this rule
if debtStatus.Status == "settled" || debtStatus.TotalRemainingDebt <= 0 {
    return error("Cannot send notification - debt already paid/settled")
}
```

---

## Auto-Disable Feature

When a debt is paid or settled, the system automatically disables all future notifications for that debt.

```go
func processNotification(notif *models.Notification) error {
    debtStatus, _ := getDebtStatus(notif.DebtListID)

    shouldSend, reason := ShouldSendNotification(debtStatus)

    if !shouldSend {
        log.Infof("Skipping notification: %s", reason)

        // Auto-disable for paid/settled debts
        if reason == "debt already settled" || reason == "debt fully paid" {
            db.Model(&models.Notification{}).
                Where("debt_list_id = ?", notif.DebtListID).
                Update("enabled", false)
        }

        return nil
    }

    // Continue with sending...
}
```

---

## Benefits

| Benefit             | Description                                          |
| ------------------- | ---------------------------------------------------- |
| **No Spam**         | Users don't receive reminders for already paid debts |
| **Professional**    | Maintains system credibility and trust               |
| **Cost Savings**    | Reduces unnecessary SMS/email costs                  |
| **Automatic**       | No manual intervention needed                        |
| **Clean Database**  | Auto-disables obsolete notifications                 |
| **User Experience** | Only relevant notifications are sent                 |

---

## Example Scenarios

### Scenario 1: Debt Paid During Reminder Period

```
Timeline:
Day 1:  Create debt (due in 10 days)
Day 2:  Schedule reminders (7, 3, 1 days before)
Day 3:  7-day reminder scheduled ✓
Day 4:  USER PAYS DEBT IN FULL
        - debt.status = "settled"
        - debt.total_remaining_debt = 0
Day 8:  3-day reminder scheduled
        → Check: debt settled? YES
        → Action: SKIP + DISABLE all future notifications
        → Result: No reminder sent ✓
Day 10: 1-day reminder scheduled
        → Check: notification disabled? YES
        → Action: SKIP
        → Result: No reminder sent ✓
```

### Scenario 2: Partial Payment

```
Timeline:
Day 1:  Create debt (total: $1000, remaining: $1000)
Day 5:  USER PAYS $300
        - debt.total_remaining_debt = $700
        - debt.status = "active"
Day 7:  Reminder scheduled
        → Check: debt settled? NO
        → Check: remaining > 0? YES ($700)
        → Action: SEND notification ✓
        → Result: Reminder sent for $700 remaining
```

### Scenario 3: Overdue Payment, Then Paid

```
Timeline:
Day 1:  Debt overdue (past due date)
Day 2:  Overdue notification sent
Day 3:  USER PAYS DEBT
        - debt.status = "settled"
        - debt.total_remaining_debt = 0
Day 4:  Daily overdue notification scheduled
        → Check: debt settled? YES
        → Action: SKIP + DISABLE
        → Result: No more overdue notifications ✓
```

---

## Code Implementation

### Helper Function

```go
type DebtStatus struct {
    Status             string
    TotalRemainingDebt decimal.Decimal
    NextPaymentDate    time.Time
    DueDate            time.Time
}

func GetDebtStatus(db *gorm.DB, debtListID uuid.UUID) (*DebtStatus, error) {
    var status DebtStatus

    err := db.Table("debt_lists").
        Select("status, total_remaining_debt, next_payment_date, due_date").
        Where("id = ?", debtListID).
        Scan(&status).Error

    return &status, err
}

func ShouldSendNotification(debtStatus *DebtStatus) (bool, string) {
    // Rule 1: Check if settled
    if debtStatus.Status == "settled" {
        return false, "debt already settled"
    }

    // Rule 2: Check if fully paid
    if debtStatus.TotalRemainingDebt.LessThanOrEqual(decimal.Zero) {
        return false, "debt fully paid"
    }

    // All checks passed
    return true, ""
}
```

### Integration in Worker

```go
func NotificationWorker() {
    scheduler.Every(1).Minute().Do(func() {
        notifications := fetchPendingNotifications()

        for _, notif := range notifications {
            // STEP 1: Check payment status FIRST
            debtStatus, err := GetDebtStatus(db, notif.DebtListID)
            if err != nil {
                log.Error("Failed to get debt status:", err)
                continue
            }

            // STEP 2: Validate should send
            shouldSend, reason := ShouldSendNotification(debtStatus)
            if !shouldSend {
                log.Infof("Skipping notification %s: %s", notif.ID, reason)

                // Auto-disable for paid/settled debts
                if strings.Contains(reason, "settled") || strings.Contains(reason, "paid") {
                    disableNotification(notif.DebtListID)
                }

                continue
            }

            // STEP 3: Proceed with sending
            sendNotification(notif)
        }
    })
}
```

---

## Database Queries

### Check Debt Status

```sql
-- Check if debt is paid/settled
SELECT
    status,
    total_remaining_debt,
    next_payment_date,
    due_date
FROM debt_lists
WHERE id = $1;
```

### Disable Notifications for Paid Debt

```sql
-- Auto-disable all notifications for a debt
UPDATE notifications
SET enabled = false,
    updated_at = NOW()
WHERE debt_list_id = $1
  AND enabled = true;
```

### Find Obsolete Notifications

```sql
-- Find notifications for paid/settled debts
SELECT n.*
FROM notifications n
JOIN debt_lists dl ON n.debt_list_id = dl.id
WHERE n.enabled = true
  AND (dl.status = 'settled' OR dl.total_remaining_debt <= 0);
```

---

## Testing Checklist

### Unit Tests

- [ ] Test `ShouldSendNotification()` with settled debt
- [ ] Test `ShouldSendNotification()` with zero remaining debt
- [ ] Test `ShouldSendNotification()` with active unpaid debt
- [ ] Test `ShouldSendNotification()` with negative remaining debt
- [ ] Test auto-disable functionality

### Integration Tests

- [ ] Create debt → Pay fully → Verify notification skipped
- [ ] Create debt → Partial payment → Verify notification sent
- [ ] Create debt → Settle → Verify all future notifications disabled
- [ ] Schedule 3 reminders → Pay before first → Verify all skipped
- [ ] Overdue notification → Pay → Verify next overdue skipped

### Manual Tests

- [ ] Create test debt
- [ ] Schedule reminder notification
- [ ] Mark debt as paid
- [ ] Verify notification not sent
- [ ] Verify notification disabled in database

---

## Error Handling

### Missing Debt Information

```go
debtStatus, err := GetDebtStatus(db, notif.DebtListID)
if err != nil {
    log.Errorf("Cannot verify debt status for notification %s: %v", notif.ID, err)
    // Don't send notification if we can't verify status
    return fmt.Errorf("debt status verification failed: %w", err)
}
```

### Database Connection Issues

```go
// Retry logic for debt status check
var debtStatus *DebtStatus
var err error

for attempt := 0; attempt < 3; attempt++ {
    debtStatus, err = GetDebtStatus(db, notif.DebtListID)
    if err == nil {
        break
    }
    time.Sleep(time.Second * time.Duration(attempt+1))
}

if err != nil {
    // After retries, fail safe: don't send notification
    return fmt.Errorf("failed to verify debt status after retries: %w", err)
}
```

---

## Logging

### Log Levels

```go
// INFO: Normal skip due to paid debt
log.Info("Skipping notification - debt already paid")

// DEBUG: Detailed debt status
log.Debugf("Debt status: %s, Remaining: %s", debtStatus.Status, debtStatus.TotalRemainingDebt)

// WARN: Unexpected state
log.Warn("Notification enabled for settled debt - auto-disabling")

// ERROR: Cannot verify status
log.Error("Failed to fetch debt status - notification not sent")
```

### Log Format

```
[2025-11-06 10:30:15] INFO  Skipping notification abc-123: debt already settled
[2025-11-06 10:30:15] INFO  Auto-disabled 3 notifications for debt xyz-789
[2025-11-06 10:30:16] DEBUG Debt check passed: Status=active, Remaining=$500.00
[2025-11-06 10:30:16] INFO  Sending notification def-456 via email
```

---

## Monitoring & Alerts

### Metrics to Track

1. **Skip Rate**: % of notifications skipped due to paid debts
2. **Auto-Disable Count**: Number of notifications auto-disabled per day
3. **False Positive Rate**: Notifications skipped incorrectly (monitor complaints)
4. **Cost Savings**: Estimated savings from skipped SMS/emails

### Alert Conditions

```yaml
alerts:
  - name: high_skip_rate
    condition: skip_rate > 50%
    action: notify_admin
    reason: Possible issue with debt status updates

  - name: debt_status_check_failures
    condition: failures > 10 per hour
    action: page_oncall
    reason: Database or connectivity issues

  - name: no_notifications_sent
    condition: sent_count == 0 for 1 hour
    action: notify_admin
    reason: Possible system failure
```

---

## Summary

✅ **Always check debt status before sending**  
✅ **Auto-disable notifications for paid/settled debts**  
✅ **Log all skipped notifications with reason**  
✅ **Fail safe: don't send if status cannot be verified**  
✅ **Monitor skip rates and cost savings**

This business rule ensures the notification system is **smart**, **cost-effective**, and **user-friendly**.
