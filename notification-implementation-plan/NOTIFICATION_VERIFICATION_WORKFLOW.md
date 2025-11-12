# Payment Verification Workflow - Summary

## Overview

The notification system now includes a **payment verification workflow** that sends different notifications based on payment status transitions.

---

## Payment Lifecycle

```
Payment Made (Pending) → Payment Verified/Rejected (Completed/Rejected)
```

### Stage 1: Payment Submitted (Status: "pending")

**Who gets notified:** BOTH user and contact

**Example: John Doe owes User Kevin $1,000, John pays $500**

1. **To User (Kevin):**
   - ✉️ "Payment received: John Doe paid you $500 and still under verification"

2. **To Contact (John):**
   - ✉️ "Payment recorded: You paid $500 to Kevin (pending verification)"

### Stage 2: Payment Verified (Status: "completed")

**Who gets notified:** The PAYER only (whoever made the payment)

**Continuing the example: Kevin verifies John's payment**

3. **To Contact (John) - the payer:**
   - ✉️ "Payment verified: Your payment of $500 to Kevin has been verified and accepted"

### Stage 3: Payment Rejected (Status: "rejected")

**Who gets notified:** The PAYER only

**Alternative: Kevin rejects John's payment**

3. **To Contact (John) - the payer:**
   - ✉️ "Payment rejected: Your payment of $500 to Kevin has been rejected. Reason: Payment proof not clear, please resubmit"

---

## API Endpoints

### Create Payment (Stage 1)

```http
POST /api/debt-items
Content-Type: application/json

{
  "debt_list_id": "debt-uuid",
  "amount": "500",
  "payment_method": "bank_transfer"
}

Response:
{
  "id": "debt-item-uuid",
  "status": "pending",
  "amount": "500.00",
  "debt_list_id": "debt-uuid"
}
```

**What happens:**
- 2 notifications sent immediately (to user and contact)
- Both messages include "pending verification" or "still under verification"

### Verify Payment (Stage 2)

```http
POST /api/debt-items/:debt_item_id/verify

Response:
{
  "message": "Payment verified successfully",
  "debt_item": {
    "id": "debt-item-uuid",
    "status": "completed",
    "amount": "500.00"
  }
}
```

**What happens:**
- 1 notification sent to payer
- Message: "Payment verified: Your payment of $XXX has been verified and accepted"

### Reject Payment (Stage 3)

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
    "id": "debt-item-uuid",
    "status": "rejected",
    "rejection_reason": "Payment proof not clear, please resubmit"
  }
}
```

**What happens:**
- 1 notification sent to payer
- Message: "Payment rejected: Your payment of $XXX has been rejected. Reason: [reason]"

---

## Implementation Details

### New Handler Functions

```go
// internal/handlers/debt_item_handler.go

// Verify payment
func (h *DebtItemHandler) VerifyPayment(c *gin.Context) {
    debtItemID := c.Param("id")
    
    // Update status to "completed"
    debtItem, err := h.debtItemService.VerifyPayment(debtItemID)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    // Send verification notification to payer
    go h.notificationService.SendPaymentVerificationNotification(
        debtItem.ID, 
        debtItem.DebtListID, 
        true,  // verified
        nil,   // no reason needed
    )

    c.JSON(200, gin.H{"message": "Payment verified successfully", "debt_item": debtItem})
}

// Reject payment
func (h *DebtItemHandler) RejectPayment(c *gin.Context) {
    debtItemID := c.Param("id")
    
    var req RejectPaymentRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // Update status to "rejected"
    debtItem, err := h.debtItemService.RejectPayment(debtItemID, req.Reason)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    // Send rejection notification to payer
    go h.notificationService.SendPaymentVerificationNotification(
        debtItem.ID,
        debtItem.DebtListID,
        false,       // rejected
        &req.Reason, // rejection reason
    )

    c.JSON(200, gin.H{"message": "Payment rejected", "debt_item": debtItem})
}
```

### New Notification Service Function

```go
// internal/services/notification/event_notification_service.go

func (s *EventNotificationService) SendPaymentVerificationNotification(
    debtItemID, debtListID uuid.UUID, 
    verified bool, 
    rejectionReason *string,
) error {
    // Fetch payment info
    info, err := s.fetchPaymentInfo(debtItemID, debtListID)
    if err != nil {
        return fmt.Errorf("failed to fetch payment info: %w", err)
    }

    // Determine who made the payment (they get the notification)
    var recipientType string
    var recipientEmail, recipientPhone *string

    if info.DebtType == "i_owe" {
        // User made payment, notify user
        recipientType = "user"
        recipientEmail = &info.UserEmail
        recipientPhone = &info.UserPhone
    } else {
        // Contact made payment, notify contact
        recipientType = "contact"
        recipientEmail = &info.ContactEmail
        recipientPhone = &info.ContactPhone
    }

    // Create verification or rejection message
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
        // Similar for rejection messages
    }

    // Create and send notification
    notification := Notification{
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
    go s.sendNotification(notification)

    return nil
}
```

---

## Complete Notification Flow

### For "Owed to Me" Debt (Contact owes User)

```
1. Contact makes payment
   → User receives: "Payment received: Contact paid you $500 and still under verification"
   → Contact receives: "Payment recorded: You paid $500 to User (pending verification)"

2. User verifies payment
   → Contact receives: "Payment verified: Your payment of $500 to User has been verified and accepted"
```

### For "I Owe" Debt (User owes Contact)

```
1. User makes payment
   → User receives: "Payment recorded: You paid $500 to Contact (pending verification)"
   → Contact receives: "Payment received: User paid you $500 and still under verification"

2. Contact verifies payment
   → User receives: "Payment verified: Your payment of $500 to Contact has been verified and accepted"
```

---

## Benefits

✅ **Trust & Transparency** - Both parties always know payment status
✅ **Clear Communication** - Distinct messages for pending vs verified
✅ **Accountability** - Rejection requires a reason
✅ **Automatic** - No manual intervention needed
✅ **Multi-Channel** - Works with email, SMS, and webhooks

---

## Files Updated

1. **NOTIFICATION_EVENT_BASED_GUIDE.md**
   - Added Scenarios 3 & 4 (Payment Verified/Rejected)
   - Updated notification messages to include "pending verification"
   - Added verification notification function
   - Added API examples for verify/reject
   - Updated complete flow examples
   - Enhanced testing checklist

2. **NOTIFICATION_COMPLETE_SYSTEM_FINAL.md**
   - Updated event-based notification section
   - Enhanced flow examples with verification steps
   - Added verification API endpoints
   - Updated example scenarios

3. **NOTIFICATION_SYSTEM_SUMMARY.md**
   - Updated event-based notifications feature description
   - Added verification workflow mention

4. **NOTIFICATION_IMPLEMENTATION_PLAN.md**
   - Updated event-based notifications architecture
   - Added verification API endpoints

---

## Next Steps for Implementation

When implementing this feature:

1. **Database Changes**
   - Ensure `debt_items` table has `status` field (pending/completed/rejected)
   - Add `rejection_reason` field to `debt_items`

2. **Service Layer**
   - Implement `VerifyPayment()` method in debt item service
   - Implement `RejectPayment()` method in debt item service
   - Implement `SendPaymentVerificationNotification()` in notification service

3. **Handler Layer**
   - Add `VerifyPayment()` handler
   - Add `RejectPayment()` handler
   - Add validation for rejection reason

4. **Routes**
   - `POST /api/debt-items/:id/verify`
   - `POST /api/debt-items/:id/reject`

5. **Testing**
   - Test verification notification to correct recipient
   - Test rejection notification with reason
   - Test for both "i_owe" and "owed_to_me" debt types
   - Test multi-channel delivery (email, SMS, webhooks)

---

**Status:** Planning Complete ✅  
**Ready for Implementation:** Yes  
**Integration:** Fully compatible with existing notification system

