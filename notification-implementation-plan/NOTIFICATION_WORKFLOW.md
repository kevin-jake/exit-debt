# Notification System Workflow Diagrams

## 1. Overall System Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                         Exit Debt Application                        │
└─────────────────────────────────────────────────────────────────────┘
                                    │
        ┌───────────────────────────┼───────────────────────────┐
        │                           │                           │
        ▼                           ▼                           ▼
┌───────────────┐          ┌─────────────────┐        ┌────────────────┐
│  API Endpoint │          │  Scheduler      │        │  Manual Trigger│
│  (Create Debt)│          │  Service        │        │  API           │
└───────┬───────┘          └────────┬────────┘        └────────┬───────┘
        │                           │                          │
        └───────────────┬───────────┴──────────────────────────┘
                        ▼
        ┌──────────────────────────────────────┐
        │   Notification Scheduler Service      │
        │   - Reads User Settings               │
        │   - Reads Debt Records                │
        │   - Creates Notification Records      │
        │   - Creates Cron Jobs                 │
        └──────────────┬───────────────────────┘
                       │
        ┌──────────────┴──────────────┐
        │     PostgreSQL Database      │
        │  - notifications             │
        │  - user_settings             │
        │  - debt_lists                │
        │  - notification_templates    │
        └──────────────┬───────────────┘
                       │
        ┌──────────────┴──────────────┐
        ▼                             ▼
┌───────────────────┐        ┌────────────────────┐
│   Go-Cron Engine  │        │ Notification Worker│
│   - Job Queue     │───────>│ - Fetches Pending  │
│   - Scheduler     │        │ - Renders Template │
│   - Job Manager   │        │ - Sends Email/SMS  │
└───────────────────┘        └─────────┬──────────┘
        │                              │
        │                    ┌─────────┴─────────┐
        │                    ▼                   ▼
        │            ┌───────────────┐   ┌──────────────┐
        │            │ Email Service │   │  SMS Service │
        │            │ (SMTP)        │   │  (Twilio)    │
        │            └───────────────┘   └──────────────┘
        │
        ▼
┌───────────────────┐
│   Go-Cron UI      │
│   - View Jobs     │
│   - Monitor Status│
│   - Trigger Jobs  │
└───────────────────┘
```

---

## 2. Notification Creation Flow

```
User Creates Debt
       │
       ▼
┌─────────────────────────────────────────────────┐
│ 1. Check User Settings                          │
│    - Get notification_reminder_days [7,3,1]     │
│    - Get notification_time "09:00"              │
│    - Get custom messages (if any)               │
└────────────────────┬────────────────────────────┘
                     ▼
┌─────────────────────────────────────────────────┐
│ 2. Calculate Reminder Dates                     │
│    due_date = 2025-01-20                        │
│    reminders:                                   │
│      - 2025-01-13 09:00 (7 days before)         │
│      - 2025-01-17 09:00 (3 days before)         │
│      - 2025-01-19 09:00 (1 day before)          │
└────────────────────┬────────────────────────────┘
                     ▼
┌─────────────────────────────────────────────────┐
│ 3. Create Notification Records                  │
│    FOR EACH reminder_date:                      │
│      INSERT INTO notifications (                │
│        debt_item_id,                            │
│        schedule_type = 'reminder',              │
│        scheduled_for = reminder_date,           │
│        status = 'pending',                      │
│        enabled = true                           │
│      )                                          │
└────────────────────┬────────────────────────────┘
                     ▼
┌─────────────────────────────────────────────────┐
│ 4. Create Cron Jobs                             │
│    FOR EACH notification:                       │
│      cron_expr = toCronExpr(scheduled_for)      │
│      job_id = scheduler.AddFunc(cron_expr, ...)  │
│      UPDATE notifications SET cron_job_id = ...  │
└────────────────────┬────────────────────────────┘
                     ▼
                  SUCCESS
```

---

## 3. Notification Processing Flow

```
┌──────────────────────────────────────────────────┐
│ Cron Job Triggers                                │
│ (Every minute, checks for pending notifications) │
└────────────────────┬─────────────────────────────┘
                     ▼
┌─────────────────────────────────────────────────┐
│ Fetch Pending Notifications                     │
│ WHERE scheduled_for <= NOW()                    │
│   AND status = 'pending'                        │
│   AND enabled = true                            │
└────────────────────┬────────────────────────────┘
                     ▼
              ┌──────┴──────┐
              │ Has Results?│
              └──────┬──────┘
                NO   │   YES
                     ▼
        ┌────────────────────────────┐
        │ FOR EACH notification:     │
        └────────┬───────────────────┘
                 ▼
        ┌─────────────────────────────────────┐
        │ Check Debt Payment Status           │
        │ Get debt_list info:                 │
        │ - status                            │
        │ - total_remaining_debt              │
        │ - next_payment_date                 │
        └────────┬────────────────────────────┘
                 ▼
        ┌──────────────────┐
        │ Is debt settled? │
        │ OR fully paid?   │
        └─────┬────────────┘
         YES  │   NO
              │
     ┌────────┴──────────┐
     ▼                   ▼
┌──────────┐      ┌───────────────────────┐
│ DISABLE  │      │ Continue Processing   │
│ Future   │      └───────────┬───────────┘
│ Notifs   │                  ▼
│ & SKIP   │      ┌─────────────────────────────────────┐
└──────────┘      │ Check if Overdue with 24hr Limit    │
                  │ IF schedule_type = 'overdue':       │
                  │   IF (NOW - last_sent_at) < 24hrs:  │
                  │     SKIP                            │
                  └────────┬────────────────────────────┘
                           ▼
                  ┌─────────────────────────────────────┐
                  │ Get User Settings                    │
                  │ - Check notification_email enabled   │
                  │ - Check notification_sms enabled     │
                  │ - Get custom messages                │
                  └────────┬────────────────────────────┘
                 ▼
        ┌─────────────────────────────────────┐
        │ Load & Render Template              │
        │ - Get template from DB              │
        │ - Replace variables                 │
        │   {{debtor_name}} → "John Doe"      │
        │   {{amount}} → "1000.00"            │
        │   {{due_date}} → "2025-01-20"       │
        └────────┬────────────────────────────┘
                 ▼
        ┌──────────────────┐
        │ Email Enabled?   │
        └─────┬────────────┘
         YES  │   NO
              ▼
    ┌──────────────────────┐
    │ Send Email (SMTP)    │
    └──────────┬───────────┘
               │
               ▼
        ┌──────────────────┐
        │ SMS Enabled?     │
        └─────┬────────────┘
         YES  │   NO
              ▼
    ┌──────────────────────┐
    │ Send SMS (Twilio)    │
    └──────────┬───────────┘
               │
               ▼
        ┌─────────────────────────────────────┐
        │ Update Notification Record          │
        │ - status = 'sent'                   │
        │ - sent_at = NOW()                   │
        │ - last_sent_at = NOW()              │
        │ - next_run_at = calculateNext()     │
        └─────────────────────────────────────┘
```

---

## 4. Custom Schedule Override Flow

```
User Edits Debt Notification Settings
       │
       ▼
┌─────────────────────────────────────────────────┐
│ API: PUT /api/notifications/:id                 │
│ Body: {                                         │
│   "use_custom_schedule": true,                  │
│   "custom_reminder_days": [10, 5, 2],           │
│   "custom_notification_time": "14:00",          │
│   "custom_message": "Custom reminder..."        │
│ }                                               │
└────────────────────┬────────────────────────────┘
                     ▼
┌─────────────────────────────────────────────────┐
│ 1. Update Notification Record                   │
│    UPDATE notifications SET                     │
│      use_custom_schedule = true,                │
│      custom_reminder_days = [10,5,2],           │
│      ...                                        │
└────────────────────┬────────────────────────────┘
                     ▼
┌─────────────────────────────────────────────────┐
│ 2. Cancel Existing Cron Jobs                    │
│    scheduler.Remove(old_cron_job_id)            │
└────────────────────┬────────────────────────────┘
                     ▼
┌─────────────────────────────────────────────────┐
│ 3. Recalculate Schedule                         │
│    Use custom_reminder_days instead of default  │
│    Use custom_notification_time                 │
└────────────────────┬────────────────────────────┘
                     ▼
┌─────────────────────────────────────────────────┐
│ 4. Create New Cron Jobs                         │
│    FOR EACH new reminder date:                  │
│      Create new cron job                        │
│      Update cron_job_id                         │
└─────────────────────────────────────────────────┘
```

---

## 5. Overdue Notification Flow

```
Daily Check (Automated)
       │
       ▼
┌─────────────────────────────────────────────────┐
│ Find Overdue Debts                              │
│ SELECT * FROM debt_lists                        │
│ WHERE next_payment_date < NOW()                 │
│   AND status != 'settled'                       │
└────────────────────┬────────────────────────────┘
                     ▼
┌─────────────────────────────────────────────────┐
│ FOR EACH overdue debt:                          │
│   Check if overdue notification exists          │
└────────────────────┬────────────────────────────┘
                     ▼
              ┌──────┴───────┐
              │ Exists?      │
              └──────┬───────┘
                NO   │   YES
                     │
     ┌───────────────┴───────────────┐
     ▼                               ▼
CREATE NEW                    CHECK LAST SENT
┌──────────────────┐         ┌─────────────────┐
│ Create Overdue   │         │ IF last_sent_at │
│ Notification:    │         │ > 24 hours ago: │
│ - schedule_type  │         │   SEND AGAIN    │
│   = 'overdue'    │         │ ELSE:           │
│ - is_recurring   │         │   SKIP          │
│   = true         │         └─────────────────┘
│ - Create daily   │
│   cron job       │
└──────────────────┘

Key Rules:
1. **Only send if debt is NOT yet paid** (status != 'settled' AND remaining_debt > 0)
2. Maximum 1 notification per 24 hours (hard limit)
3. Continues daily until debt is paid or manually disabled
4. Uses overdue template (different from reminder template)
5. Auto-disable notification when debt is settled/fully paid
```

---

## 6. Manual Trigger Flow (API Call)

```
User/System Triggers Manual Send
       │
       ▼
┌─────────────────────────────────────────────────┐
│ API: POST /api/notifications/send               │
│ Body: {                                         │
│   "debt_item_id": "uuid...",                    │
│   "notification_type": "email",                 │
│   "message": "Custom message...",               │
│   "recipient_email": "user@example.com"         │
│ }                                               │
└────────────────────┬────────────────────────────┘
                     ▼
┌─────────────────────────────────────────────────┐
│ 1. Create Notification Record                   │
│    schedule_type = 'manual'                     │
│    status = 'pending'                           │
│    scheduled_for = NOW()                        │
└────────────────────┬────────────────────────────┘
                     ▼
┌─────────────────────────────────────────────────┐
│ 2. Send Immediately                             │
│    (No cron job needed)                         │
│    - Render message                             │
│    - Send via email/SMS                         │
│    - Update status to 'sent'                    │
└────────────────────┬────────────────────────────┘
                     ▼
┌─────────────────────────────────────────────────┐
│ 3. Return Response                              │
│    {                                            │
│      "success": true,                           │
│      "notification_id": "uuid...",              │
│      "sent_at": "2025-11-06T10:30:00Z"          │
│    }                                            │
└─────────────────────────────────────────────────┘
```

---

## 7. Template Variable Replacement Flow

```
Notification Ready to Send
       │
       ▼
┌─────────────────────────────────────────────────┐
│ Load Template                                   │
│ - Check if custom_message exists               │
│ - Else, load from notification_templates       │
│ - Else, use default hardcoded template         │
└────────────────────┬────────────────────────────┘
                     ▼
┌─────────────────────────────────────────────────┐
│ Gather Data for Variables                      │
│ - Get User data (first_name, last_name)        │
│ - Get Contact data (name, email, phone)        │
│ - Get Debt data (amount, due_date, etc.)       │
│ - Calculate days_until_due / days_overdue      │
└────────────────────┬────────────────────────────┘
                     ▼
┌─────────────────────────────────────────────────┐
│ Replace Variables                               │
│                                                 │
│ Template:                                       │
│ "Hi {{user_first_name}}, payment of            │
│  {{currency}}{{amount}} due to {{debtor_name}} │
│  in {{days_until_due}} days."                  │
│                                                 │
│ Result:                                         │
│ "Hi Kevin, payment of Php1000.00 due to        │
│  John Doe in 3 days."                          │
└────────────────────┬────────────────────────────┘
                     ▼
            Return Rendered Message
```

---

## 8. Database Schema Relationships

```
┌──────────────────┐
│     users        │
│ ─────────────── │
│ id (PK)          │
│ email            │
│ first_name       │
│ last_name        │
│ phone            │
└────────┬─────────┘
         │ 1
         │
    ┌────┴──────────────────────────┐
    │                               │
    │ 1                             │ *
┌────┴────────────┐     ┌───────────┴─────────┐
│  user_settings  │     │   user_contacts     │◄─── Contact info stored here!
│ ─────────────── │     │ ──────────────────  │
│ id (PK)         │     │ id (PK)             │
│ user_id (FK)    │     │ user_id (FK)        │
│ notification_*  │     │ contact_id (FK)     │
│ custom_*_message│     │ name                │
└─────────────────┘     │ email ◄────────────┼─── Email for notifications
                        │ phone ◄────────────┼─── Phone for notifications
                        │ notes               │
                        └──────────┬──────────┘
                                   │ *
                                   │
                                   │ 1
                        ┌──────────┴──────────┐
                        │   contacts          │
                        │ ──────────────────  │
                        │ id (PK)             │
                        │ is_user             │
                        │ user_id_ref         │
                        └──────────┬──────────┘
                                   │ 1
                                   │
                                   │ *
                        ┌──────────┴──────────┐
                        │   debt_lists        │
                        │ ──────────────────  │
                        │ id (PK)             │
                        │ user_id (FK)        │
                        │ contact_id (FK) ◄───┼─── Links to contact
                        │ due_date            │
                        │ next_payment_date   │
                        └──────────┬──────────┘
                                   │ 1
                                   │
                                   │ *
                        ┌──────────┴───────────────────┐
                        │   notifications              │
                        │ ────────────────────────────│
                        │ id (PK)                      │
                        │ debt_list_id (FK) ◄──────────┼─── Links to debt_list
                        │ notification_type            │
                        │ schedule_type                │
                        │ scheduled_for                │
                        │ cron_job_id                  │
                        │ recipient_email (optional)   │◄─ Override field
                        │ recipient_phone (optional)   │◄─ Override field
                        │ use_custom_schedule          │
                        │ custom_reminder_days         │
                        │ custom_notification_time     │
                        │ custom_message               │
                        │ last_sent_at                 │
                        │ next_run_at                  │
                        │ is_recurring                 │
                        │ enabled                      │
                        │ status                       │
                        └──────────────────────────────┘

┌──────────────────────────┐
│ notification_templates   │
│ ────────────────────────│
│ id (PK)                  │
│ user_id (FK) nullable    │◄──── User-specific or global
│ template_type            │◄──── 'email' or 'sms'
│ subject                  │
│ body                     │
│ is_default               │
│ variables                │
└──────────────────────────┘

┌─────────────────────────────────────────────────────────┐
│ Contact Info Fetching Flow:                             │
│                                                          │
│ notification.debt_list_id                                │
│     → debt_list.contact_id + debt_list.user_id           │
│         → user_contact WHERE contact_id AND user_id      │
│             → user_contact.email, user_contact.phone     │
│                                                          │
│ Falls back to notification.recipient_email/phone if set  │
└─────────────────────────────────────────────────────────┘
```

---

## 9. Cron Job Lifecycle

```
┌─────────────────────────────────────────────────┐
│ CREATED                                         │
│ - Notification record created                   │
│ - Cron job scheduled                            │
│ - Status: 'pending'                             │
└────────────────────┬────────────────────────────┘
                     ▼
┌─────────────────────────────────────────────────┐
│ WAITING                                         │
│ - Job in scheduler queue                        │
│ - Waiting for scheduled_for time                │
└────────────────────┬────────────────────────────┘
                     ▼
┌─────────────────────────────────────────────────┐
│ TRIGGERED                                       │
│ - scheduled_for time reached                    │
│ - Cron job executes                             │
│ - Calls processNotification(id)                 │
└────────────────────┬────────────────────────────┘
                     ▼
              ┌──────┴──────┐
              │ Send Success?│
              └──────┬──────┘
                YES  │   NO
     ┌───────────────┴───────────────┐
     ▼                               ▼
┌────────────┐              ┌─────────────────┐
│ COMPLETED  │              │ RETRY           │
│ - Status:  │              │ - Attempt retry │
│   'sent'   │              │ - Max 3 attempts│
│ - Job done │              └────────┬────────┘
│ (one-time) │                       ▼
│ OR         │              ┌─────────────────┐
│ - Schedule │              │ FAILED          │
│   next run │              │ - Status:       │
│ (recurring)│              │   'failed'      │
└────────────┘              │ - Log error     │
                            │ - Alert admin   │
                            └─────────────────┘
```

---

## 10. Error Handling & Recovery

```
┌─────────────────────────────────────────────────┐
│ Notification Send Attempt                       │
└────────────────────┬────────────────────────────┘
                     ▼
              ┌──────┴──────┐
              │ Send Failed? │
              └──────┬──────┘
                NO   │   YES
                     ▼
        ┌─────────────────────────┐
        │ Log Error               │
        │ - Error message         │
        │ - Timestamp             │
        │ - Attempt number        │
        └─────────┬───────────────┘
                  ▼
        ┌─────────────────────────┐
        │ Check Retry Count       │
        │ IF attempts < max_retry │
        └─────────┬───────────────┘
                  ▼
           ┌──────┴──────┐
           │ Can Retry?   │
           └──────┬──────┘
             YES  │   NO
     ┌────────────┴────────────┐
     ▼                         ▼
┌──────────────┐      ┌────────────────┐
│ RETRY        │      │ MARK AS FAILED │
│ - Wait 2^N   │      │ - Update status│
│   seconds    │      │ - Send alert   │
│ - Try again  │      │ - Notify admin │
└──────────────┘      └────────────────┘

Common Error Scenarios:
1. SMTP connection failed → Retry with backoff
2. Twilio API error → Check credentials, retry
3. Invalid phone number → Mark failed, don't retry
4. Invalid email → Mark failed, don't retry
5. Template error → Log error, notify admin
6. Database connection lost → Retry up to max
```

---

## Summary

This notification system provides:

✅ **Flexibility**: Custom schedules per debt or use defaults  
✅ **Reliability**: Retry logic and error handling  
✅ **Scalability**: Cron-based scheduling handles any volume  
✅ **Control**: User preferences respected at all times  
✅ **Monitoring**: Go-Cron UI for real-time visibility  
✅ **Safety**: 24-hour limit on overdue notifications  
✅ **Customization**: Template system with variables

**Key Flows:**

1. Debt creation → Auto-schedule reminders
2. Due date approaching → Send reminders
3. Payment overdue → Send daily reminders (24hr limit)
4. Manual trigger → Immediate send
5. Custom schedule → Override defaults per debt
