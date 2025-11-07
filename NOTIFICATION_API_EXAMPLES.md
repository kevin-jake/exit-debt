# Notification API Examples & Testing

This document provides complete API request/response examples for testing the notification system.

---

## 1. User Settings - Get Notification Settings

### Request

```http
GET /api/users/settings/notifications
Authorization: Bearer {jwt_token}
```

### Response

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "notification_email": true,
  "notification_sms": false,
  "notification_reminder_days": [7, 3, 1],
  "notification_time": "09:00:00",
  "overdue_reminder_frequency": "daily",
  "custom_email_message": null,
  "custom_sms_message": null,
  "default_currency": "Php",
  "timezone": "UTC",
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-01T00:00:00Z"
}
```

---

## 2. User Settings - Update Notification Settings

### Request

```http
PUT /api/users/settings/notifications
Authorization: Bearer {jwt_token}
Content-Type: application/json

{
  "notification_email": true,
  "notification_sms": true,
  "notification_reminder_days": [10, 5, 2, 1],
  "notification_time": "14:00:00",
  "overdue_reminder_frequency": "daily",
  "custom_email_message": "Hi {{user_first_name}}, you have a payment of {{currency}}{{amount}} due to {{debtor_name}} on {{due_date}}. Please don't forget!",
  "custom_sms_message": "Payment reminder: {{currency}}{{amount}} due to {{debtor_name}} on {{due_date}}"
}
```

### Response

```json
{
  "success": true,
  "message": "Notification settings updated successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "notification_email": true,
    "notification_sms": true,
    "notification_reminder_days": [10, 5, 2, 1],
    "notification_time": "14:00:00",
    "overdue_reminder_frequency": "daily",
    "custom_email_message": "Hi {{user_first_name}}, you have a payment...",
    "custom_sms_message": "Payment reminder: {{currency}}{{amount}}...",
    "updated_at": "2025-11-06T10:30:00Z"
  }
}
```

---

## 3. Create Notification Schedule for Debt

### Request

```http
POST /api/debt-lists/123e4567-e89b-12d3-a456-426614174001/schedule-reminders
Authorization: Bearer {jwt_token}
Content-Type: application/json

{
  "use_custom_schedule": false
}
```

### Response

```json
{
  "success": true,
  "message": "Reminder notifications scheduled successfully",
  "data": {
    "debt_list_id": "123e4567-e89b-12d3-a456-426614174001",
    "notifications_created": 3,
    "notifications": [
      {
        "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
        "debt_item_id": "987e6543-e21b-12d3-a456-426614174000",
        "schedule_type": "reminder",
        "scheduled_for": "2025-01-13T09:00:00Z",
        "reminder_days_before": 7,
        "status": "pending",
        "enabled": true
      },
      {
        "id": "b2c3d4e5-f6a7-8901-bcde-f12345678901",
        "debt_item_id": "987e6543-e21b-12d3-a456-426614174000",
        "schedule_type": "reminder",
        "scheduled_for": "2025-01-17T09:00:00Z",
        "reminder_days_before": 3,
        "status": "pending",
        "enabled": true
      },
      {
        "id": "c3d4e5f6-a7b8-9012-cdef-123456789012",
        "debt_item_id": "987e6543-e21b-12d3-a456-426614174000",
        "schedule_type": "reminder",
        "scheduled_for": "2025-01-19T09:00:00Z",
        "reminder_days_before": 1,
        "status": "pending",
        "enabled": true
      }
    ]
  }
}
```

---

## 4. Create Notification Schedule with Custom Settings

### Request

```http
POST /api/debt-lists/123e4567-e89b-12d3-a456-426614174001/schedule-reminders
Authorization: Bearer {jwt_token}
Content-Type: application/json

{
  "use_custom_schedule": true,
  "custom_reminder_days": [14, 7, 3],
  "custom_notification_time": "16:00:00",
  "custom_message": "Important: Your payment of {{currency}}{{amount}} to {{debtor_name}} is due on {{due_date}}. Please prepare your payment."
}
```

### Response

```json
{
  "success": true,
  "message": "Custom reminder notifications scheduled successfully",
  "data": {
    "debt_list_id": "123e4567-e89b-12d3-a456-426614174001",
    "notifications_created": 3,
    "use_custom_schedule": true,
    "notifications": [
      {
        "id": "d4e5f6a7-b8c9-0123-def1-234567890123",
        "debt_item_id": "987e6543-e21b-12d3-a456-426614174000",
        "schedule_type": "reminder",
        "scheduled_for": "2025-01-06T16:00:00Z",
        "reminder_days_before": 14,
        "use_custom_schedule": true,
        "custom_reminder_days": [14, 7, 3],
        "custom_notification_time": "16:00:00",
        "custom_message": "Important: Your payment...",
        "status": "pending",
        "enabled": true
      }
    ]
  }
}
```

---

## 5. Send Manual Notification

### Request

```http
POST /api/notifications/send
Authorization: Bearer {jwt_token}
Content-Type: application/json

{
  "debt_item_id": "987e6543-e21b-12d3-a456-426614174000",
  "notification_type": "email",
  "recipient_email": "user@example.com",
  "message": "This is a manual notification about your payment."
}
```

### Response

```json
{
  "success": true,
  "message": "Notification sent successfully",
  "data": {
    "id": "e5f6a7b8-c9d0-1234-ef12-345678901234",
    "debt_item_id": "987e6543-e21b-12d3-a456-426614174000",
    "notification_type": "email",
    "schedule_type": "manual",
    "status": "sent",
    "sent_at": "2025-11-06T10:35:00Z",
    "recipient_email": "user@example.com",
    "message": "This is a manual notification about your payment."
  }
}
```

---

## 6. List All Notifications

### Request

```http
GET /api/notifications?status=pending&schedule_type=reminder&limit=20&offset=0
Authorization: Bearer {jwt_token}
```

### Query Parameters

- `status` (optional): pending, sent, failed
- `schedule_type` (optional): reminder, overdue, manual
- `enabled` (optional): true, false
- `limit` (optional): default 20
- `offset` (optional): default 0

### Response

```json
{
  "success": true,
  "data": {
    "total": 5,
    "limit": 20,
    "offset": 0,
    "notifications": [
      {
        "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
        "debt_item_id": "987e6543-e21b-12d3-a456-426614174000",
        "notification_type": "email",
        "schedule_type": "reminder",
        "scheduled_for": "2025-01-13T09:00:00Z",
        "reminder_days_before": 7,
        "status": "pending",
        "enabled": true,
        "use_custom_schedule": false,
        "is_recurring": false,
        "created_at": "2025-01-01T10:00:00Z"
      }
    ]
  }
}
```

---

## 7. Get Notification Details

### Request

```http
GET /api/notifications/a1b2c3d4-e5f6-7890-abcd-ef1234567890
Authorization: Bearer {jwt_token}
```

### Response

```json
{
  "success": true,
  "data": {
    "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "debt_item_id": "987e6543-e21b-12d3-a456-426614174000",
    "notification_type": "email",
    "recipient_email": "user@example.com",
    "recipient_phone": null,
    "message": "Payment reminder message...",
    "status": "pending",
    "schedule_type": "reminder",
    "scheduled_for": "2025-01-13T09:00:00Z",
    "cron_job_id": "cron_a1b2c3d4",
    "reminder_days_before": 7,
    "use_custom_schedule": false,
    "custom_reminder_days": null,
    "custom_notification_time": null,
    "custom_message": null,
    "last_sent_at": null,
    "next_run_at": "2025-01-13T09:00:00Z",
    "is_recurring": false,
    "enabled": true,
    "sent_at": null,
    "created_at": "2025-01-01T10:00:00Z"
  }
}
```

---

## 8. Update Notification Schedule

### Request

```http
PUT /api/notifications/a1b2c3d4-e5f6-7890-abcd-ef1234567890
Authorization: Bearer {jwt_token}
Content-Type: application/json

{
  "use_custom_schedule": true,
  "custom_reminder_days": [5],
  "custom_notification_time": "15:00:00",
  "custom_message": "Urgent: Payment due soon!",
  "enabled": true
}
```

### Response

```json
{
  "success": true,
  "message": "Notification updated successfully",
  "data": {
    "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "use_custom_schedule": true,
    "custom_reminder_days": [5],
    "custom_notification_time": "15:00:00",
    "custom_message": "Urgent: Payment due soon!",
    "scheduled_for": "2025-01-15T15:00:00Z",
    "enabled": true,
    "updated_at": "2025-11-06T10:40:00Z"
  }
}
```

---

## 9. Enable/Disable Notification

### Enable Notification

```http
PATCH /api/notifications/a1b2c3d4-e5f6-7890-abcd-ef1234567890/enable
Authorization: Bearer {jwt_token}
```

### Response

```json
{
  "success": true,
  "message": "Notification enabled successfully",
  "data": {
    "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "enabled": true,
    "updated_at": "2025-11-06T10:45:00Z"
  }
}
```

### Disable Notification

```http
PATCH /api/notifications/a1b2c3d4-e5f6-7890-abcd-ef1234567890/disable
Authorization: Bearer {jwt_token}
```

### Response

```json
{
  "success": true,
  "message": "Notification disabled successfully",
  "data": {
    "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "enabled": false,
    "updated_at": "2025-11-06T10:46:00Z"
  }
}
```

---

## 10. Delete Notification

### Request

```http
DELETE /api/notifications/a1b2c3d4-e5f6-7890-abcd-ef1234567890
Authorization: Bearer {jwt_token}
```

### Response

```json
{
  "success": true,
  "message": "Notification deleted successfully"
}
```

---

## 11. Create Notification Template

### Request

```http
POST /api/notification-templates
Authorization: Bearer {jwt_token}
Content-Type: application/json

{
  "template_name": "Friendly Payment Reminder",
  "template_type": "email",
  "subject": "Friendly Reminder: Payment Due to {{debtor_name}}",
  "body": "<h2>Hi {{user_first_name}}!</h2><p>Just a friendly reminder that you have a payment coming up:</p><ul><li>Amount: {{currency}}{{amount}}</li><li>Due Date: {{due_date}}</li><li>To: {{debtor_name}}</li></ul><p>Thanks for staying on top of your payments!</p>",
  "is_default": false,
  "variables": [
    "user_first_name",
    "debtor_name",
    "amount",
    "currency",
    "due_date"
  ]
}
```

### Response

```json
{
  "success": true,
  "message": "Template created successfully",
  "data": {
    "id": "f6a7b8c9-d0e1-2345-f123-456789012345",
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "template_name": "Friendly Payment Reminder",
    "template_type": "email",
    "subject": "Friendly Reminder: Payment Due to {{debtor_name}}",
    "body": "<h2>Hi {{user_first_name}}!</h2>...",
    "is_default": false,
    "variables": [
      "user_first_name",
      "debtor_name",
      "amount",
      "currency",
      "due_date"
    ],
    "created_at": "2025-11-06T10:50:00Z",
    "updated_at": "2025-11-06T10:50:00Z"
  }
}
```

---

## 12. List All Templates

### Request

```http
GET /api/notification-templates?template_type=email&is_default=false
Authorization: Bearer {jwt_token}
```

### Query Parameters

- `template_type` (optional): email, sms
- `is_default` (optional): true, false

### Response

```json
{
  "success": true,
  "data": {
    "total": 3,
    "templates": [
      {
        "id": "f6a7b8c9-d0e1-2345-f123-456789012345",
        "template_name": "Friendly Payment Reminder",
        "template_type": "email",
        "subject": "Friendly Reminder: Payment Due to {{debtor_name}}",
        "is_default": false,
        "created_at": "2025-11-06T10:50:00Z"
      },
      {
        "id": "a7b8c9d0-e1f2-3456-1234-567890123456",
        "template_name": "Urgent Overdue Notice",
        "template_type": "email",
        "subject": "URGENT: Overdue Payment to {{debtor_name}}",
        "is_default": false,
        "created_at": "2025-11-06T09:00:00Z"
      }
    ]
  }
}
```

---

## 13. Get Default Templates

### Request

```http
GET /api/notification-templates/defaults
Authorization: Bearer {jwt_token}
```

### Response

```json
{
  "success": true,
  "data": {
    "email_reminder": {
      "subject": "Payment Reminder - Due in {{days_until_due}} days",
      "body": "<!DOCTYPE html><html>...",
      "variables": [
        "user_first_name",
        "debtor_name",
        "amount",
        "currency",
        "due_date",
        "days_until_due",
        "remaining_debt"
      ]
    },
    "email_overdue": {
      "subject": "Overdue Payment Alert - Action Required",
      "body": "<!DOCTYPE html><html>...",
      "variables": [
        "user_first_name",
        "debtor_name",
        "amount",
        "currency",
        "due_date",
        "days_overdue",
        "remaining_debt"
      ]
    },
    "sms_reminder": {
      "body": "Reminder: Payment of {{currency}}{{amount}} due to {{debtor_name}} in {{days_until_due}} days ({{due_date}})",
      "variables": [
        "currency",
        "amount",
        "debtor_name",
        "days_until_due",
        "due_date"
      ]
    },
    "sms_overdue": {
      "body": "OVERDUE: Payment of {{currency}}{{amount}} to {{debtor_name}} is {{days_overdue}} days overdue (due: {{due_date}})",
      "variables": [
        "currency",
        "amount",
        "debtor_name",
        "days_overdue",
        "due_date"
      ]
    }
  }
}
```

---

## 14. Update Template

### Request

```http
PUT /api/notification-templates/f6a7b8c9-d0e1-2345-f123-456789012345
Authorization: Bearer {jwt_token}
Content-Type: application/json

{
  "template_name": "Updated Friendly Reminder",
  "subject": "Updated: Payment Due to {{debtor_name}}",
  "body": "<h2>Updated template body...</h2>"
}
```

### Response

```json
{
  "success": true,
  "message": "Template updated successfully",
  "data": {
    "id": "f6a7b8c9-d0e1-2345-f123-456789012345",
    "template_name": "Updated Friendly Reminder",
    "subject": "Updated: Payment Due to {{debtor_name}}",
    "body": "<h2>Updated template body...</h2>",
    "updated_at": "2025-11-06T11:00:00Z"
  }
}
```

---

## 15. Delete Template

### Request

```http
DELETE /api/notification-templates/f6a7b8c9-d0e1-2345-f123-456789012345
Authorization: Bearer {jwt_token}
```

### Response

```json
{
  "success": true,
  "message": "Template deleted successfully"
}
```

---

## 16. Sync All Debt Records with Notification System

This endpoint creates notification schedules for all existing debts that don't have notifications yet.

### Request

```http
POST /api/notifications/sync-all
Authorization: Bearer {jwt_token}
```

### Response

```json
{
  "success": true,
  "message": "Notification sync completed",
  "data": {
    "total_debts_processed": 50,
    "notifications_created": 150,
    "debts_with_notifications": [
      {
        "debt_list_id": "123e4567-e89b-12d3-a456-426614174001",
        "notifications_count": 3
      },
      {
        "debt_list_id": "234e5678-e89b-12d3-a456-426614174002",
        "notifications_count": 3
      }
    ]
  }
}
```

---

## Testing with cURL

### Example: Send Manual Notification

```bash
curl -X POST http://localhost:8080/api/notifications/send \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "debt_item_id": "987e6543-e21b-12d3-a456-426614174000",
    "notification_type": "email",
    "recipient_email": "user@example.com",
    "message": "Test notification"
  }'
```

### Example: Update User Settings

```bash
curl -X PUT http://localhost:8080/api/users/settings/notifications \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "notification_email": true,
    "notification_sms": true,
    "notification_reminder_days": [7, 3, 1],
    "notification_time": "09:00:00"
  }'
```

### Example: Schedule Reminders

```bash
curl -X POST http://localhost:8080/api/debt-lists/123e4567-e89b-12d3-a456-426614174001/schedule-reminders \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "use_custom_schedule": false
  }'
```

---

## Testing with Postman

### Collection Setup

1. Create a new collection: "Exit Debt - Notifications"
2. Add environment variables:
   - `base_url`: http://localhost:8080
   - `jwt_token`: Your JWT token
3. Set Authorization header for all requests:
   - Type: Bearer Token
   - Token: `{{jwt_token}}`

### Pre-request Script (for all endpoints)

```javascript
// Generate timestamp
pm.environment.set("timestamp", new Date().toISOString());
```

### Test Script (for successful responses)

```javascript
pm.test("Status code is 200", function () {
  pm.response.to.have.status(200);
});

pm.test("Response has success field", function () {
  var jsonData = pm.response.json();
  pm.expect(jsonData).to.have.property("success");
  pm.expect(jsonData.success).to.be.true;
});
```

---

## Error Response Examples

### 400 Bad Request

```json
{
  "success": false,
  "error": "Invalid request body",
  "details": {
    "field": "notification_type",
    "message": "notification_type must be either 'email' or 'sms'"
  }
}
```

### 401 Unauthorized

```json
{
  "success": false,
  "error": "Unauthorized",
  "message": "Invalid or expired token"
}
```

### 404 Not Found

```json
{
  "success": false,
  "error": "Notification not found",
  "notification_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}
```

### 500 Internal Server Error

```json
{
  "success": false,
  "error": "Failed to send notification",
  "message": "SMTP connection failed",
  "details": "dial tcp: lookup smtp.gmail.com: no such host"
}
```

---

## Webhook Events (Future Enhancement)

For real-time notifications about notification status:

### Notification Sent Event

```json
{
  "event": "notification.sent",
  "timestamp": "2025-11-06T10:35:00Z",
  "data": {
    "notification_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "debt_item_id": "987e6543-e21b-12d3-a456-426614174000",
    "notification_type": "email",
    "status": "sent",
    "sent_at": "2025-11-06T10:35:00Z"
  }
}
```

### Notification Failed Event

```json
{
  "event": "notification.failed",
  "timestamp": "2025-11-06T10:35:00Z",
  "data": {
    "notification_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "debt_item_id": "987e6543-e21b-12d3-a456-426614174000",
    "notification_type": "sms",
    "status": "failed",
    "error": "Invalid phone number",
    "retry_count": 3
  }
}
```

---

## Performance Considerations

### Pagination

For large result sets, use pagination:

```http
GET /api/notifications?limit=50&offset=100
```

### Filtering

Combine multiple filters:

```http
GET /api/notifications?status=pending&schedule_type=reminder&enabled=true
```

### Bulk Operations

For creating multiple notifications at once:

```http
POST /api/notifications/bulk-create
```

---

## Security Best Practices

1. **Rate Limiting**: Limit API calls per user

   - Max 100 requests per minute per user
   - Max 10 manual notifications per hour

2. **Input Validation**: All inputs are validated

   - Email addresses validated
   - Phone numbers validated (E.164 format)
   - Template variables sanitized

3. **Authorization**: Users can only access their own data

   - Notifications linked to user's debts only
   - Templates are user-specific or global defaults

4. **HTTPS**: All API calls must use HTTPS in production

---

This completes the API examples and testing guide. Use these examples to test your notification system implementation.
