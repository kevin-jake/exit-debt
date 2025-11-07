# Notification Architecture Update - Contact Information

## Summary of Changes

The notification system has been updated with two major improvements:

1. **Contact information** is now fetched from the `user_contacts` table instead of storing them directly in the `notifications` table
2. **Payment status checking** ensures notifications are only sent for unpaid debts

---

## Why This Change?

### Before (Old Approach)

```
Notification table stored:
- debt_item_id
- recipient_email
- recipient_phone

Problem: Duplicated contact data, hard to update if contact info changes
```

### After (New Approach)

```
Notification table stores:
- debt_list_id (link to debt)
- recipient_email (optional override)
- recipient_phone (optional override)

Contact info fetched from user_contacts table via relationship chain
```

---

## Relationship Chain

```
┌─────────────────┐
│  notifications  │
│  (pending send) │
└────────┬────────┘
         │
         │ debt_list_id
         ▼
┌─────────────────┐
│   debt_lists    │
│  (debt record)  │
└────────┬────────┘
         │
         │ contact_id + user_id
         ▼
┌─────────────────┐
│ user_contacts   │◄──── Email & Phone stored here!
│ (contact info)  │
└─────────────────┘
```

---

## How It Works

### 1. When Creating a Notification

```go
// Create notification linked to debt_list
notification := models.Notification{
    DebtListID:       debtListID,
    NotificationType: "email",
    Message:          "Payment reminder",
    Status:           "pending",
}

// NO need to specify recipient_email/phone
// System will fetch from user_contacts automatically
```

### 2. When Sending a Notification

```go
// Step 1: Fetch contact info from user_contacts
contactInfo := fetchContactInfo(notification.DebtListID)
// Returns: email, phone from user_contacts table

// Step 2: Check for override
email := notification.RecipientEmail  // Check override first
if email == nil {
    email = contactInfo.Email  // Use from user_contacts
}

// Step 3: Send notification
sendEmail(email, message)
```

### 3. SQL Query to Fetch Contact Info

```sql
SELECT
    uc.email,
    uc.phone,
    uc.name as contact_name,
    dl.user_id
FROM notifications n
JOIN debt_lists dl ON n.debt_list_id = dl.id
JOIN user_contacts uc ON uc.contact_id = dl.contact_id
                     AND uc.user_id = dl.user_id
WHERE n.id = $1;
```

---

## Benefits

### ✅ Single Source of Truth

- Contact information is stored in ONE place: `user_contacts` table
- No data duplication
- Easy to update contact info globally

### ✅ Automatic Updates

- If user updates contact email/phone in `user_contacts`
- Future notifications automatically use new contact info
- No need to update old notification records

### ✅ User-Specific Contact Info

- Same contact can have different email/phone for different users
- Respects the many-to-many relationship: users ↔ contacts

### ✅ Optional Override

- Can still override contact info per notification if needed
- Useful for testing or temporary changes
- Override is stored in notification record

---

## Override Behavior

### When to Use Overrides

1. **Testing/Debugging**

   ```go
   notification.RecipientEmail = stringPtr("test@example.com")
   // Temporarily send to test email
   ```

2. **Manual Notifications**

   ```go
   // User manually enters alternative contact
   notification.RecipientPhone = stringPtr("+1234567890")
   ```

3. **Temporary Contact Changes**
   ```go
   // Don't want to update user_contacts, just this notification
   notification.RecipientEmail = stringPtr("alternative@example.com")
   ```

### Override Logic

```
IF notification.recipient_email IS SET:
    Use notification.recipient_email
ELSE:
    Fetch email from user_contacts table

IF notification.recipient_phone IS SET:
    Use notification.recipient_phone
ELSE:
    Fetch phone from user_contacts table
```

---

## Code Implementation

### Helper Function

```go
// GetEffectiveRecipient determines which email/phone to use
func GetEffectiveRecipient(notification *models.Notification, contactInfo *ContactInfo) (email *string, phone *string) {
    // Email: override takes precedence
    if notification.RecipientEmail != nil && *notification.RecipientEmail != "" {
        email = notification.RecipientEmail  // Use override
    } else {
        email = contactInfo.Email  // Use from user_contacts
    }

    // Phone: override takes precedence
    if notification.RecipientPhone != nil && *notification.RecipientPhone != "" {
        phone = notification.RecipientPhone  // Use override
    } else {
        phone = contactInfo.Phone  // Use from user_contacts
    }

    return email, phone
}
```

### Notification Worker

```go
func processNotification(notif *models.Notification) error {
    // 1. Check debt payment status - CRITICAL: Only send if NOT yet paid
    debtStatus, err := contactFetcher.GetDebtStatus(notif.DebtListID)
    if err != nil {
        return fmt.Errorf("failed to get debt status: %w", err)
    }

    // 2. Check if notification should be sent
    shouldSend, reason := ShouldSendNotification(debtStatus)
    if !shouldSend {
        log.Infof("Skipping notification: %s", reason)
        // Auto-disable future notifications for settled/paid debts
        if reason == "debt already settled" || reason == "debt fully paid" {
            disableNotification(notif.ID)
        }
        return nil
    }

    // 3. Fetch contact info from user_contacts
    contactInfo, err := contactFetcher.GetContactInfoForNotification(notif.ID)
    if err != nil {
        return fmt.Errorf("contact info not found: %w", err)
    }

    // 4. Determine effective recipient (applies override logic)
    email, phone := GetEffectiveRecipient(notif, contactInfo)

    // 5. Validate recipient exists
    if notif.NotificationType == "email" && email == nil {
        return fmt.Errorf("no email address available")
    }
    if notif.NotificationType == "sms" && phone == nil {
        return fmt.Errorf("no phone number available")
    }

    // 6. Send notification
    if notif.NotificationType == "email" {
        return emailSender.SendEmail(*email, subject, body)
    } else {
        return smsSender.SendSMS(*phone, message)
    }
}

// ShouldSendNotification checks if notification should be sent
func ShouldSendNotification(debtStatus *DebtStatus) (bool, string) {
    // Don't send if debt is settled
    if debtStatus.Status == "settled" {
        return false, "debt already settled"
    }

    // Don't send if fully paid
    if debtStatus.TotalRemainingDebt.LessThanOrEqual(decimal.Zero) {
        return false, "debt fully paid"
    }

    return true, ""
}
```

---

## Database Schema Changes

### Notifications Table

**Old Schema:**

```sql
CREATE TABLE notifications (
    id UUID PRIMARY KEY,
    debt_item_id UUID REFERENCES debt_items(id),  -- Old: linked to debt_item
    recipient_email VARCHAR,  -- Old: required field
    recipient_phone VARCHAR,  -- Old: required field
    ...
);
```

**New Schema:**

```sql
CREATE TABLE notifications (
    id UUID PRIMARY KEY,
    debt_list_id UUID REFERENCES debt_lists(id),  -- New: linked to debt_list
    recipient_email VARCHAR,  -- New: optional override
    recipient_phone VARCHAR,  -- New: optional override
    ...
);
```

### Migration Script

```sql
-- Step 1: Change foreign key
ALTER TABLE notifications DROP COLUMN debt_item_id;
ALTER TABLE notifications ADD COLUMN debt_list_id UUID REFERENCES debt_lists(id);

-- Step 2: Make email/phone optional (already nullable in most DBs)
ALTER TABLE notifications ALTER COLUMN recipient_email DROP NOT NULL;
ALTER TABLE notifications ALTER COLUMN recipient_phone DROP NOT NULL;

-- Step 3: Add index for performance
CREATE INDEX idx_notifications_debt_list ON notifications(debt_list_id);
```

---

## Example Scenarios

### Scenario 1: Normal Notification (No Override)

```go
// Create notification
notification := &models.Notification{
    DebtListID:       debtListUUID,
    NotificationType: "email",
    Message:          "Payment reminder",
}

// At send time:
// 1. System queries user_contacts via debt_list relationship
// 2. Finds email: "john@example.com"
// 3. Sends to: "john@example.com"
```

### Scenario 2: Testing with Override

```go
// Create notification with override
notification := &models.Notification{
    DebtListID:       debtListUUID,
    NotificationType: "email",
    Message:          "Payment reminder",
    RecipientEmail:   stringPtr("test@internal.com"),  // Override
}

// At send time:
// 1. System checks notification.RecipientEmail
// 2. Finds override: "test@internal.com"
// 3. Sends to: "test@internal.com" (ignores user_contacts)
```

### Scenario 3: Contact Updated

```sql
-- User updates their contact's email
UPDATE user_contacts
SET email = 'newemail@example.com'
WHERE user_id = $1 AND contact_id = $2;

-- Future notifications automatically use new email
-- No need to update notification records
```

---

## API Changes

### Create Notification

**Old API:**

```json
POST /api/notifications/send
{
  "debt_item_id": "uuid...",
  "recipient_email": "required@example.com",  // Required
  "recipient_phone": "+1234567890"            // Required
}
```

**New API:**

```json
POST /api/notifications/send
{
  "debt_list_id": "uuid...",
  "recipient_email": "optional@example.com",  // Optional override
  "recipient_phone": "+1234567890"            // Optional override
}

// If not provided, system fetches from user_contacts automatically
```

---

## Error Handling

### Missing Contact Info

```go
contactInfo, err := fetchContactInfo(debtListID)
if err != nil {
    return fmt.Errorf("contact not found in user_contacts: %w", err)
}

if contactInfo.Email == nil && notification.RecipientEmail == nil {
    return fmt.Errorf("no email address available for notification")
}
```

### Invalid Debt List

```go
// Ensure debt_list exists and has valid contact relationship
var count int64
db.Model(&models.DebtList{}).
    Joins("JOIN user_contacts uc ON uc.contact_id = debt_lists.contact_id").
    Where("debt_lists.id = ?", debtListID).
    Count(&count)

if count == 0 {
    return fmt.Errorf("invalid debt_list or missing contact relationship")
}
```

---

## Testing

### Test Contact Fetching

```go
func TestFetchContactInfo(t *testing.T) {
    // Setup
    user := createTestUser()
    contact := createTestContact()
    userContact := createTestUserContact(user.ID, contact.ID, "test@example.com", "+1234567890")
    debtList := createTestDebtList(user.ID, contact.ID)
    notification := createTestNotification(debtList.ID)

    // Test
    contactInfo, err := contactFetcher.GetContactInfoForNotification(notification.ID)

    // Assert
    assert.NoError(t, err)
    assert.Equal(t, "test@example.com", *contactInfo.Email)
    assert.Equal(t, "+1234567890", *contactInfo.Phone)
}
```

### Test Override Logic

```go
func TestOverrideLogic(t *testing.T) {
    // Setup
    contactInfo := &ContactInfo{
        Email: stringPtr("contact@example.com"),
        Phone: stringPtr("+1111111111"),
    }

    // Test with override
    notification := &models.Notification{
        RecipientEmail: stringPtr("override@example.com"),
    }

    email, phone := GetEffectiveRecipient(notification, contactInfo)

    // Assert
    assert.Equal(t, "override@example.com", *email)  // Uses override
    assert.Equal(t, "+1111111111", *phone)           // Uses from contact
}
```

---

## Migration Guide for Existing Data

If you have existing notification records with `debt_item_id`, run this migration:

```sql
-- Step 1: Add new debt_list_id column
ALTER TABLE notifications ADD COLUMN debt_list_id UUID;

-- Step 2: Populate debt_list_id from debt_item_id
UPDATE notifications n
SET debt_list_id = di.debt_list_id
FROM debt_items di
WHERE n.debt_item_id = di.id;

-- Step 3: Make debt_list_id NOT NULL and add foreign key
ALTER TABLE notifications ALTER COLUMN debt_list_id SET NOT NULL;
ALTER TABLE notifications ADD CONSTRAINT fk_notifications_debt_list
    FOREIGN KEY (debt_list_id) REFERENCES debt_lists(id) ON DELETE CASCADE;

-- Step 4: Drop old debt_item_id column
ALTER TABLE notifications DROP COLUMN debt_item_id;

-- Step 5: Add index
CREATE INDEX idx_notifications_debt_list ON notifications(debt_list_id);
```

---

## Business Rule: Only Send for Unpaid Debts

### The Problem

We shouldn't send payment reminders for debts that are already paid or settled. This would confuse users and damage credibility.

### The Solution

Before sending any notification, check the debt's payment status:

```go
// Check 1: Is debt settled?
if debtStatus.Status == "settled" {
    return false, "debt already settled"
}

// Check 2: Is remaining debt zero or negative?
if debtStatus.TotalRemainingDebt <= 0 {
    return false, "debt fully paid"
}
```

### When to Check

```
For EVERY notification before sending:
1. Reminder notifications (N days before due)
2. Overdue notifications (daily after due)
3. Manual notifications (API triggered)
```

### Auto-Disable Feature

When a debt is paid/settled, automatically disable future notifications:

```go
if !shouldSend {
    if reason == "debt already settled" || reason == "debt fully paid" {
        disableNotification(notif.ID)  // Prevent future sends
    }
}
```

### Benefits

✅ **Prevents spam**: No reminders for paid debts  
✅ **Professional**: Maintains system credibility  
✅ **Automatic**: No manual intervention needed  
✅ **Efficient**: Saves SMS/email costs  
✅ **Smart**: Auto-disables when debt is settled

---

## Summary

✅ **Simplified**: Contact info fetched from `user_contacts` (single source of truth)  
✅ **Flexible**: Can override per notification if needed  
✅ **Maintainable**: Updates to contacts automatically apply to future notifications  
✅ **Scalable**: Proper database relationships, no data duplication  
✅ **Smart**: Only sends notifications for unpaid debts  
✅ **Automatic**: Disables notifications when debts are settled

This architecture ensures that contact information is managed centrally while still providing flexibility for special cases through optional overrides, and guarantees that notifications are only sent when they make sense (unpaid debts).
