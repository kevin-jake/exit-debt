# Notification System - Phase Implementation Guide

**Quick reference for implementing each phase of the notification system.**

---

## 📋 Phase Overview

```
┌─────────────────────────────────────────────────────────────┐
│ PLANNING PHASE ✅ COMPLETE                                   │
│ • 12 documentation files created                            │
│ • 7,146 lines of comprehensive documentation                │
│ • All features designed and documented                      │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│ PHASE 1: Database & Models (Days 1-2)                       │
│ Documents: IMPLEMENTATION_PLAN, ARCHITECTURE_UPDATE          │
│ Goal: Set up database schema and Go models                  │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│ PHASE 2: Core Services (Days 2-3)                           │
│ Documents: QUICK_START, WEBHOOK_GUIDE                       │
│ Goal: Implement email, SMS, webhook senders                 │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│ PHASE 3: Business Logic (Days 3-4)                          │
│ Documents: INSTALLMENT_SYSTEM, BUSINESS_RULES, EVENT_BASED  │
│ Goal: Notification scheduling and payment checking          │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│ PHASE 4: API Layer (Days 4-5)                               │
│ Documents: API_EXAMPLES, QUICK_START                        │
│ Goal: Create API endpoints for management                   │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│ PHASE 5: Integration & Testing (Days 5-6)                   │
│ Documents: SYSTEM_SUMMARY, WEBHOOK_GUIDE                    │
│ Goal: Configure, integrate UI, and test                     │
└─────────────────────────────────────────────────────────────┘
```

---

## 🎯 Phase 1: Database & Models

### What You'll Build

- User_settings table with notification preferences
- Notifications table with installment tracking
- Notification_templates table
- Complete GORM models

### Documents to Read

1. **Primary:** `NOTIFICATION_DATABASE_SETUP.md` ⭐ **START HERE** - Complete copy-paste ready code
2. **Reference:** `NOTIFICATION_IMPLEMENTATION_PLAN.md` (Lines 86-240) - Field descriptions
3. **Reference:** `NOTIFICATION_ARCHITECTURE_UPDATE.md` - Data flow diagrams
4. **Reference:** `NOTIFICATION_INSTALLMENT_SYSTEM.md` (Lines 1-150) - Installment fields

### Step-by-Step

```bash
# 1. Read the database setup guide (has complete code ready to copy)
code NOTIFICATION_DATABASE_SETUP.md

# 2. Update user_settings.go model
code internal/models/user_settings.go
# Copy the UserSettings struct from DATABASE_SETUP.md
# Copy the UpdateUserSettingsRequest struct

# 3. Update notification.go model
code internal/models/notification.go
# Copy the Notification struct from DATABASE_SETUP.md
# Copy the CreateNotificationRequest struct

# 4. Create notification_template.go model
touch internal/models/notification_template.go
code internal/models/notification_template.go
# Copy the complete file content from DATABASE_SETUP.md

# 5. Add AutoMigrate to main.go (or database init function)
code main.go
# Add the initDatabase function from DATABASE_SETUP.md

# 6. Run your application
go run main.go
# GORM AutoMigrate will create tables automatically

# 7. Verify tables created
psql your_database -c "\dt"
# Should show user_settings, notifications, notification_templates
```

### Key Files to Create/Update

- `internal/models/user_settings.go` (update with notification fields)
- `internal/models/notification.go` (update with all new fields)
- `internal/models/notification_template.go` (create new)
- Update main.go or init function to run AutoMigrate

### Success Criteria

- [ ] All models compile successfully
- [ ] GORM AutoMigrate creates tables without errors
- [ ] Can query new fields from database
- [ ] Foreign keys and indexes created automatically
- [ ] No migration files needed

### Estimated Time: 4-6 hours

---

## 🎯 Phase 2: Core Services

### What You'll Build

- Email sender (SMTP)
- SMS sender (Twilio)
- Slack webhook sender
- Telegram bot sender
- Discord webhook sender
- Template engine
- Cron scheduler

### Documents to Read

1. **Primary:** `NOTIFICATION_QUICK_START.md` (Lines 1-400)
2. **Primary:** `NOTIFICATION_WEBHOOK_GUIDE.md` (All)
3. **Reference:** `NOTIFICATION_IMPLEMENTATION_PLAN.md` (Lines 550-656)

### Step-by-Step

```bash
# 1. Read QUICK_START.md for email & SMS
# 2. Read WEBHOOK_GUIDE.md for Slack, Telegram, Discord

# 3. Create service directory structure
mkdir -p internal/services/notification
mkdir -p internal/services/scheduler

# 4. Implement services in order:
#    a. Email (simplest)
#    b. SMS (Twilio API)
#    c. Template engine (needed for all)
#    d. Slack webhook
#    e. Telegram bot
#    f. Discord webhook
#    g. Unified webhook service
#    h. Cron scheduler

# 5. Create environment configuration
touch internal/config/notification_config.go

# 6. Test each service individually
```

### Key Files to Create

- `internal/services/notification/email_sender.go`
- `internal/services/notification/sms_sender.go`
- `internal/services/notification/slack_sender.go`
- `internal/services/notification/telegram_sender.go`
- `internal/services/notification/discord_sender.go`
- `internal/services/notification/webhook_service.go`
- `internal/services/notification/template_engine.go`
- `internal/services/scheduler/cron_service.go`
- `internal/config/notification_config.go`

### Dependencies to Install

```bash
go get github.com/go-co-op/gocron
go get gopkg.in/gomail.v2
go get github.com/twilio/twilio-go
go get github.com/go-telegram-bot-api/telegram-bot-api/v5
go get github.com/hashicorp/go-retryablehttp
go get github.com/lib/pq
```

### Success Criteria

- [ ] Can send test email via SMTP
- [ ] Can send test SMS via Twilio
- [ ] Can send test Slack message
- [ ] Can send test Telegram message
- [ ] Can send test Discord message
- [ ] Templates render variables correctly
- [ ] Cron jobs can be scheduled

### Estimated Time: 8-12 hours

---

## 🎯 Phase 3: Business Logic

### What You'll Build

- Payment status checker
- Installment-based notification scheduler
- Event-based notification service
- Payment verification service
- Auto-disable logic
- Contact info fetcher
- Notification worker

### Documents to Read

1. **Primary:** `NOTIFICATION_INSTALLMENT_SYSTEM.md` (All)
2. **Primary:** `NOTIFICATION_BUSINESS_RULES.md` (All)
3. **Primary:** `NOTIFICATION_EVENT_BASED_GUIDE.md` (All)
4. **Primary:** `NOTIFICATION_VERIFICATION_WORKFLOW.md` (All)
5. **Reference:** `NOTIFICATION_IMPLEMENTATION_PLAN.md` (Lines 243-547)

### Step-by-Step

```bash
# 1. Read all primary documents in order

# 2. Implement in this order:
#    a. Payment checker (foundational)
#    b. Contact info fetcher (needed by all)
#    c. Installment scheduler (creates notifications)
#    d. Event notification service (payment events)
#    e. Verification service (verify/reject)
#    f. Auto-disable logic (cleanup)
#    g. Notification worker (processes queue)

# 3. Add hooks to existing handlers
#    - debt_item_handler.go (payment events)

# 4. Test each component
```

### Key Files to Create

- `internal/services/notification/payment_checker.go`
- `internal/services/notification/installment_scheduler.go`
- `internal/services/notification/event_notification_service.go`
- `internal/services/notification/scheduler.go`
- `internal/services/notification/service.go` (main service)

### Key Files to Update

- `internal/handlers/debt_item_handler.go` (add event hooks)

### Success Criteria

- [ ] Payment status correctly identified
- [ ] Installment notifications created for all payments
- [ ] Event notifications sent on payment
- [ ] Verification notifications sent correctly
- [ ] Notifications auto-disabled when paid
- [ ] Contact info fetched correctly
- [ ] Worker processes notification queue

### Estimated Time: 10-14 hours

---

## 🎯 Phase 4: API Layer

### What You'll Build

- Notification repository
- Template repository
- Notification handlers
- Template handlers
- Verification handlers
- API routes

### Documents to Read

1. **Primary:** `NOTIFICATION_API_EXAMPLES.md` (All)
2. **Primary:** `NOTIFICATION_QUICK_START.md` (Lines 800-1084)
3. **Reference:** `NOTIFICATION_IMPLEMENTATION_PLAN.md` (Lines 660-691)

### Step-by-Step

```bash
# 1. Read API_EXAMPLES for endpoint specifications
# 2. Read QUICK_START for handler implementations

# 3. Create repositories first (data access layer)
mkdir -p internal/repositories

# 4. Create handlers (business logic)
#    - notification_handler.go
#    - template_handler.go
#    - Add verify/reject to debt_item_handler.go

# 5. Set up routes in routes.go

# 6. Test each endpoint with curl or Postman
```

### Key Files to Create

- `internal/repositories/notification_repository.go`
- `internal/repositories/template_repository.go`
- `internal/handlers/notification_handler.go`
- `internal/handlers/template_handler.go`

### Key Files to Update

- `internal/handlers/debt_item_handler.go` (add verify/reject)
- `internal/routes/routes.go` (add routes)

### API Endpoints to Implement

```
POST   /api/debt-lists
POST   /api/debt-lists/:id/schedule-notifications
GET    /api/debt-lists/:id/notifications
GET    /api/debt-lists/:id/installments/:num/notifications

POST   /api/notifications/send
GET    /api/notifications
GET    /api/notifications/:id
PUT    /api/notifications/:id
DELETE /api/notifications/:id
PATCH  /api/notifications/:id/enable
PATCH  /api/notifications/:id/disable

POST   /api/debt-items/:id/verify
POST   /api/debt-items/:id/reject

POST   /api/notification-templates
GET    /api/notification-templates
GET    /api/notification-templates/defaults
PUT    /api/notification-templates/:id
DELETE /api/notification-templates/:id

GET    /api/users/settings/notifications
PUT    /api/users/settings/notifications
```

### Success Criteria

- [ ] All endpoints respond correctly
- [ ] Request validation works
- [ ] Error handling is robust
- [ ] Filtering and pagination work
- [ ] Authentication/authorization enforced

### Estimated Time: 6-10 hours

---

## 🎯 Phase 5: Integration & Testing

### What You'll Build

- Environment configuration
- Go-cron UI integration
- Initialization service
- Comprehensive test suite

### Documents to Read

1. **Primary:** `NOTIFICATION_SYSTEM_SUMMARY.md` (Lines 173-221, 396-426)
2. **Primary:** `NOTIFICATION_WEBHOOK_GUIDE.md` (Lines 600-841)
3. **Primary:** `NOTIFICATION_INSTALLMENT_SYSTEM.md` (Lines 450-579)
4. **Primary:** `NOTIFICATION_EVENT_BASED_GUIDE.md` (Lines 977-1022)
5. **Primary:** `NOTIFICATION_BUSINESS_RULES.md` (Lines 350-448)

### Step-by-Step

```bash
# 1. Set up environment variables
cp .env.example .env
# Add SMTP, Twilio, Telegram credentials

# 2. Configure test accounts
#    - Gmail SMTP
#    - Twilio (SMS)
#    - Slack webhook
#    - Telegram bot
#    - Discord webhook

# 3. Integrate go-cron UI
# Follow go-cron documentation

# 4. Create initialization service
# Loads default templates, seeds data

# 5. Run test suite
#    - Unit tests for each service
#    - Integration tests for flows
#    - E2E tests for complete scenarios

# 6. Performance testing
# 7. Error handling validation
```

### Environment Variables Needed

```env
# SMTP
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password

# Twilio (SMS)
TWILIO_ACCOUNT_SID=ACxxxxxx
TWILIO_AUTH_TOKEN=your_token
TWILIO_PHONE_NUMBER=+1234567890

# Database
DATABASE_URL=postgresql://...
```

### Testing Checklist

#### Email & SMS

- [ ] Send test email
- [ ] Send test SMS
- [ ] Verify delivery
- [ ] Test with invalid credentials
- [ ] Test rate limiting

#### Webhooks

- [ ] Send Slack notification
- [ ] Send Telegram notification
- [ ] Send Discord notification
- [ ] Test with invalid webhook URLs
- [ ] Test error handling

#### Installment Flow

- [ ] Create installment debt
- [ ] Verify notifications created (12 × 3 = 36)
- [ ] Verify cron jobs scheduled
- [ ] Test reminder sending
- [ ] Pay installment #1
- [ ] Verify remaining reminders cancelled
- [ ] Verify installment #2 reminders still active

#### Event Notifications

- [ ] Make payment
- [ ] Verify 2 notifications sent (user + contact)
- [ ] Verify correct messages
- [ ] Verify payment verification
- [ ] Verify rejection notification

#### Edge Cases

- [ ] Contact has no email
- [ ] Multiple payments quickly
- [ ] Pay full debt early
- [ ] Change installment due date
- [ ] Delete debt (cascade)

### Success Criteria

- [ ] All channels working
- [ ] All test cases pass
- [ ] Performance acceptable
- [ ] Error handling robust
- [ ] Documentation complete
- [ ] Go-cron UI accessible

### Estimated Time: 8-12 hours

---

## 📊 Total Implementation Estimate

| Phase     | Time            | Difficulty   |
| --------- | --------------- | ------------ |
| Phase 1   | 4-6 hours       | Medium       |
| Phase 2   | 8-12 hours      | High         |
| Phase 3   | 10-14 hours     | High         |
| Phase 4   | 6-10 hours      | Medium       |
| Phase 5   | 8-12 hours      | Medium       |
| **Total** | **36-54 hours** | **4-6 days** |

---

## 🚀 Getting Started

### Before You Begin

1. Read `NOTIFICATION_SYSTEM_SUMMARY.md` for overview
2. Read `NOTIFICATION_IMPLEMENTATION_PLAN.md` for complete plan
3. Review `NOTIFICATION_WORKFLOW.md` for visual diagrams
4. Set up development environment
5. Install dependencies

### Start Implementation

1. Open `NOTIFICATION_DOCS_INDEX.md` for detailed phase guide
2. Begin with Phase 1
3. Follow step-by-step instructions
4. Use primary documents for each phase
5. Reference supplementary docs as needed
6. Test after each component
7. Mark off checklist items

### Track Progress

Use the TODO system:

```bash
# View todos
# Check off items as you complete them
```

---

## 💡 Pro Tips

1. **Don't skip testing** - Test each component before moving on
2. **Use QUICK_START.md** - Has the most detailed code examples
3. **Keep WORKFLOW.md open** - Visual reference helps
4. **Read BUSINESS_RULES.md** - Especially for payment logic
5. **Test with real accounts** - Set up actual SMTP, Twilio, etc.
6. **Follow the order** - Each phase builds on previous
7. **Use the docs index** - Quick navigation to specific topics

---

## 🎯 Next Steps

1. **Right now:** Read `NOTIFICATION_DOCS_INDEX.md`
2. **Today:** Set up development environment
3. **Tomorrow:** Start Phase 1 (Database & Models)
4. **This week:** Complete Phases 1-3
5. **Next week:** Complete Phases 4-5 and testing

---

**Ready to start?** Open `NOTIFICATION_DOCS_INDEX.md` and begin Phase 1!

**Questions?** Refer to the comprehensive documentation - everything is covered!

**Good luck!** 🚀
