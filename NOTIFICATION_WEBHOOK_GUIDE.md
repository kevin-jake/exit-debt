# Webhook Notifications Guide - Slack, Telegram & Discord

## Overview

The notification system supports sending payment reminders via webhooks to popular messaging platforms:
- **Slack** - Using Incoming Webhooks
- **Telegram** - Using Bot API
- **Discord** - Using Webhooks

This allows users to receive debt payment notifications in their preferred messaging platform.

---

## Supported Platforms

### 1. Slack
- **Method**: Incoming Webhooks
- **Configuration**: Webhook URL
- **Message Format**: Rich formatted messages with blocks
- **Documentation**: https://api.slack.com/messaging/webhooks

### 2. Telegram
- **Method**: Bot API
- **Configuration**: Bot Token + Chat ID
- **Message Format**: Markdown or HTML formatted messages
- **Documentation**: https://core.telegram.org/bots/api

### 3. Discord
- **Method**: Webhooks
- **Configuration**: Webhook URL
- **Message Format**: Embeds (rich cards)
- **Documentation**: https://discord.com/developers/docs/resources/webhook

---

## Database Schema Updates

### User Settings Table

```sql
-- Add webhook notification settings
ALTER TABLE user_settings ADD COLUMN notification_webhook BOOLEAN DEFAULT false;
ALTER TABLE user_settings ADD COLUMN slack_webhook_url VARCHAR(500);
ALTER TABLE user_settings ADD COLUMN telegram_bot_token VARCHAR(255);
ALTER TABLE user_settings ADD COLUMN telegram_chat_id VARCHAR(255);
ALTER TABLE user_settings ADD COLUMN discord_webhook_url VARCHAR(500);
```

### Notification Table

```sql
-- Add webhook type field
ALTER TABLE notifications ADD COLUMN webhook_type VARCHAR(50);  -- 'slack', 'telegram', 'discord'

-- Add constraint
ALTER TABLE notifications ADD CONSTRAINT check_webhook_type 
    CHECK (
        (notification_type != 'webhook') OR 
        (notification_type = 'webhook' AND webhook_type IN ('slack', 'telegram', 'discord'))
    );
```

---

## Go Models

### User Settings (Updated)

```go
type UserSettings struct {
    ID                  uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
    UserID              uuid.UUID `json:"user_id" gorm:"type:uuid;not null;uniqueIndex"`
    NotificationEmail   bool      `json:"notification_email" gorm:"default:true"`
    NotificationSMS     bool      `json:"notification_sms" gorm:"default:false"`
    NotificationWebhook bool      `json:"notification_webhook" gorm:"default:false"`
    
    // Webhook configurations
    SlackWebhookURL     *string   `json:"slack_webhook_url"`
    TelegramBotToken    *string   `json:"telegram_bot_token"`
    TelegramChatID      *string   `json:"telegram_chat_id"`
    DiscordWebhookURL   *string   `json:"discord_webhook_url"`
    
    // ... other fields
}
```

### Notification (Updated)

```go
type Notification struct {
    ID                   uuid.UUID  `json:"id"`
    DebtListID           uuid.UUID  `json:"debt_list_id"`
    InstallmentNumber    *int       `json:"installment_number"`
    InstallmentDueDate   *time.Time `json:"installment_due_date"`
    NotificationType     string     `json:"notification_type"`  // 'email', 'sms', 'webhook'
    WebhookType          *string    `json:"webhook_type"`       // 'slack', 'telegram', 'discord'
    Message              string     `json:"message"`
    Status               string     `json:"status"`
    // ... other fields
}
```

---

## Environment Variables

Add to `.env`:

```bash
# Telegram Bot Configuration (optional, can be per-user)
TELEGRAM_DEFAULT_BOT_TOKEN=your_bot_token_here

# Webhook Settings
WEBHOOK_TIMEOUT=10s
WEBHOOK_MAX_RETRY=3
WEBHOOK_RATE_LIMIT=10  # requests per second per platform
```

---

## Dependencies

```bash
# Slack SDK (optional, or use standard HTTP)
go get github.com/slack-go/slack

# Telegram Bot API
go get github.com/go-telegram-bot-api/telegram-bot-api/v5

# Discord Webhooks (standard HTTP is sufficient, but SDK available)
go get github.com/bwmarrin/discordgo

# HTTP client with retry
go get github.com/hashicorp/go-retryablehttp
```

---

## Implementation

### 1. Slack Webhook Sender

```go
// internal/services/notification/slack_sender.go
package notification

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

type SlackSender struct {
    client *http.Client
}

func NewSlackSender() *SlackSender {
    return &SlackSender{
        client: &http.Client{
            Timeout: 10 * time.Second,
        },
    }
}

// SlackMessage represents a Slack message with blocks
type SlackMessage struct {
    Text   string        `json:"text"`
    Blocks []SlackBlock  `json:"blocks"`
}

type SlackBlock struct {
    Type string                 `json:"type"`
    Text *SlackTextObject       `json:"text,omitempty"`
    Fields []SlackTextObject    `json:"fields,omitempty"`
}

type SlackTextObject struct {
    Type string `json:"type"`
    Text string `json:"text"`
}

func (s *SlackSender) SendNotification(webhookURL string, data NotificationData) error {
    message := s.formatMessage(data)
    
    payload, err := json.Marshal(message)
    if err != nil {
        return fmt.Errorf("failed to marshal payload: %w", err)
    }
    
    resp, err := s.client.Post(webhookURL, "application/json", bytes.NewBuffer(payload))
    if err != nil {
        return fmt.Errorf("failed to send webhook: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("slack webhook returned status %d", resp.StatusCode)
    }
    
    return nil
}

func (s *SlackSender) formatMessage(data NotificationData) SlackMessage {
    // Create rich formatted message
    return SlackMessage{
        Text: "💰 Payment Reminder",
        Blocks: []SlackBlock{
            {
                Type: "header",
                Text: &SlackTextObject{
                    Type: "plain_text",
                    Text: "💰 Payment Reminder",
                },
            },
            {
                Type: "section",
                Text: &SlackTextObject{
                    Type: "mrkdwn",
                    Text: fmt.Sprintf("Hi *%s*,\n\nYou have a payment due soon!", data.UserFirstName),
                },
            },
            {
                Type: "section",
                Fields: []SlackTextObject{
                    {
                        Type: "mrkdwn",
                        Text: fmt.Sprintf("*To:*\n%s", data.DebtorName),
                    },
                    {
                        Type: "mrkdwn",
                        Text: fmt.Sprintf("*Amount:*\n%s %s", data.Currency, data.Amount),
                    },
                    {
                        Type: "mrkdwn",
                        Text: fmt.Sprintf("*Due Date:*\n%s", data.DueDate.Format("Jan 02, 2006")),
                    },
                    {
                        Type: "mrkdwn",
                        Text: fmt.Sprintf("*Days Until Due:*\n%d days", data.DaysUntilDue),
                    },
                },
            },
            {
                Type: "section",
                Text: &SlackTextObject{
                    Type: "mrkdwn",
                    Text: fmt.Sprintf("*Installment:* #%d of %d", data.InstallmentNumber, data.InstallmentTotal),
                },
            },
            {
                Type: "divider",
            },
            {
                Type: "context",
                Text: &SlackTextObject{
                    Type: "mrkdwn",
                    Text: fmt.Sprintf("Remaining Balance: %s %s", data.Currency, data.RemainingDebt),
                },
            },
        },
    }
}
```

### 2. Telegram Bot Sender

```go
// internal/services/notification/telegram_sender.go
package notification

import (
    "fmt"
    
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramSender struct {
    defaultBotToken string
}

func NewTelegramSender(defaultBotToken string) *TelegramSender {
    return &TelegramSender{
        defaultBotToken: defaultBotToken,
    }
}

func (t *TelegramSender) SendNotification(botToken, chatID string, data NotificationData) error {
    // Use user's bot token or default
    token := botToken
    if token == "" {
        token = t.defaultBotToken
    }
    
    bot, err := tgbotapi.NewBotAPI(token)
    if err != nil {
        return fmt.Errorf("failed to create bot: %w", err)
    }
    
    message := t.formatMessage(data)
    
    msg := tgbotapi.NewMessage(chatID, message)
    msg.ParseMode = "Markdown"
    
    _, err = bot.Send(msg)
    if err != nil {
        return fmt.Errorf("failed to send message: %w", err)
    }
    
    return nil
}

func (t *TelegramSender) formatMessage(data NotificationData) string {
    return fmt.Sprintf(`
💰 *Payment Reminder*

Hi %s,

You have a payment due soon!

*Payment Details:*
• To: %s
• Amount: %s %s
• Due Date: %s
• Days Until Due: %d days
• Installment: #%d of %d

💵 Remaining Balance: %s %s

_Stay on top of your payments!_
`,
        data.UserFirstName,
        data.DebtorName,
        data.Currency,
        data.Amount,
        data.DueDate.Format("Jan 02, 2006"),
        data.DaysUntilDue,
        data.InstallmentNumber,
        data.InstallmentTotal,
        data.Currency,
        data.RemainingDebt,
    )
}
```

### 3. Discord Webhook Sender

```go
// internal/services/notification/discord_sender.go
package notification

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

type DiscordSender struct {
    client *http.Client
}

func NewDiscordSender() *DiscordSender {
    return &DiscordSender{
        client: &http.Client{
            Timeout: 10 * time.Second,
        },
    }
}

type DiscordWebhook struct {
    Content string         `json:"content,omitempty"`
    Embeds  []DiscordEmbed `json:"embeds,omitempty"`
}

type DiscordEmbed struct {
    Title       string              `json:"title"`
    Description string              `json:"description"`
    Color       int                 `json:"color"`
    Fields      []DiscordEmbedField `json:"fields"`
    Footer      *DiscordEmbedFooter `json:"footer,omitempty"`
    Timestamp   string              `json:"timestamp,omitempty"`
}

type DiscordEmbedField struct {
    Name   string `json:"name"`
    Value  string `json:"value"`
    Inline bool   `json:"inline"`
}

type DiscordEmbedFooter struct {
    Text string `json:"text"`
}

func (d *DiscordSender) SendNotification(webhookURL string, data NotificationData) error {
    webhook := d.formatMessage(data)
    
    payload, err := json.Marshal(webhook)
    if err != nil {
        return fmt.Errorf("failed to marshal payload: %w", err)
    }
    
    resp, err := d.client.Post(webhookURL, "application/json", bytes.NewBuffer(payload))
    if err != nil {
        return fmt.Errorf("failed to send webhook: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
        return fmt.Errorf("discord webhook returned status %d", resp.StatusCode)
    }
    
    return nil
}

func (d *DiscordSender) formatMessage(data NotificationData) DiscordWebhook {
    // Color: Blue for reminder, Red for overdue
    color := 3447003 // Blue
    if data.DaysUntilDue < 0 {
        color = 15158332 // Red
    }
    
    return DiscordWebhook{
        Embeds: []DiscordEmbed{
            {
                Title:       "💰 Payment Reminder",
                Description: fmt.Sprintf("Hi **%s**,\n\nYou have a payment due soon!", data.UserFirstName),
                Color:       color,
                Fields: []DiscordEmbedField{
                    {
                        Name:   "To",
                        Value:  data.DebtorName,
                        Inline: true,
                    },
                    {
                        Name:   "Amount",
                        Value:  fmt.Sprintf("%s %s", data.Currency, data.Amount),
                        Inline: true,
                    },
                    {
                        Name:   "Due Date",
                        Value:  data.DueDate.Format("Jan 02, 2006"),
                        Inline: true,
                    },
                    {
                        Name:   "Days Until Due",
                        Value:  fmt.Sprintf("%d days", data.DaysUntilDue),
                        Inline: true,
                    },
                    {
                        Name:   "Installment",
                        Value:  fmt.Sprintf("#%d of %d", data.InstallmentNumber, data.InstallmentTotal),
                        Inline: true,
                    },
                    {
                        Name:   "Remaining Balance",
                        Value:  fmt.Sprintf("%s %s", data.Currency, data.RemainingDebt),
                        Inline: true,
                    },
                },
                Footer: &DiscordEmbedFooter{
                    Text: "Exit Debt Payment Reminder",
                },
                Timestamp: time.Now().Format(time.RFC3339),
            },
        },
    }
}
```

### 4. Unified Webhook Service

```go
// internal/services/notification/webhook_service.go
package notification

import (
    "fmt"
)

type WebhookService struct {
    slackSender    *SlackSender
    telegramSender *TelegramSender
    discordSender  *DiscordSender
}

func NewWebhookService(telegramDefaultToken string) *WebhookService {
    return &WebhookService{
        slackSender:    NewSlackSender(),
        telegramSender: NewTelegramSender(telegramDefaultToken),
        discordSender:  NewDiscordSender(),
    }
}

func (ws *WebhookService) SendNotification(
    webhookType string,
    settings UserSettings,
    data NotificationData,
) error {
    switch webhookType {
    case "slack":
        if settings.SlackWebhookURL == nil {
            return fmt.Errorf("slack webhook URL not configured")
        }
        return ws.slackSender.SendNotification(*settings.SlackWebhookURL, data)
        
    case "telegram":
        if settings.TelegramChatID == nil {
            return fmt.Errorf("telegram chat ID not configured")
        }
        botToken := ""
        if settings.TelegramBotToken != nil {
            botToken = *settings.TelegramBotToken
        }
        return ws.telegramSender.SendNotification(botToken, *settings.TelegramChatID, data)
        
    case "discord":
        if settings.DiscordWebhookURL == nil {
            return fmt.Errorf("discord webhook URL not configured")
        }
        return ws.discordSender.SendNotification(*settings.DiscordWebhookURL, data)
        
    default:
        return fmt.Errorf("unsupported webhook type: %s", webhookType)
    }
}

type NotificationData struct {
    UserFirstName      string
    UserLastName       string
    DebtorName         string
    Amount             string
    Currency           string
    DueDate            time.Time
    DaysUntilDue       int
    InstallmentNumber  int
    InstallmentTotal   int
    RemainingDebt      string
}
```

---

## Integration with Notification Processor

```go
func processNotification(notificationID uuid.UUID) {
    notification := getNotification(notificationID)
    
    // ... existing payment status checks ...
    
    // Fetch contact info and user settings
    contactInfo, _ := getContactInfoForNotification(notification.DebtListID)
    userSettings, _ := getUserSettings(notification.UserID)
    
    // Prepare notification data
    data := NotificationData{
        UserFirstName:     contactInfo.UserFirstName,
        DebtorName:        contactInfo.ContactName,
        Amount:            notification.Amount.String(),
        Currency:          notification.Currency,
        DueDate:           *notification.InstallmentDueDate,
        DaysUntilDue:      calculateDaysUntilDue(*notification.InstallmentDueDate),
        InstallmentNumber: *notification.InstallmentNumber,
        InstallmentTotal:  calculateTotalInstallments(notification.DebtListID),
        RemainingDebt:     contactInfo.TotalRemainingDebt.String(),
    }
    
    var err error
    
    // Send based on notification type
    switch notification.NotificationType {
    case "email":
        if userSettings.NotificationEmail {
            err = emailSender.SendEmail(contactInfo.Email, "Payment Reminder", renderTemplate(data))
        }
        
    case "sms":
        if userSettings.NotificationSMS {
            err = smsSender.SendSMS(contactInfo.Phone, renderTemplate(data))
        }
        
    case "webhook":
        if userSettings.NotificationWebhook && notification.WebhookType != nil {
            err = webhookService.SendNotification(*notification.WebhookType, userSettings, data)
        }
    }
    
    // Update notification record
    if err == nil {
        updateNotificationSent(notificationID)
    } else {
        handleNotificationError(notificationID, err)
    }
}
```

---

## Setup Guides for Users

### Setting Up Slack

1. Go to your Slack workspace
2. Navigate to: Apps → Incoming Webhooks
3. Click "Add to Slack"
4. Choose a channel to post to
5. Copy the Webhook URL
6. Add to Exit Debt settings

**API Endpoint:**
```http
PUT /api/users/settings/notifications
{
  "notification_webhook": true,
  "slack_webhook_url": "https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXX"
}
```

### Setting Up Telegram

1. Create a bot:
   - Open Telegram and search for @BotFather
   - Send `/newbot` command
   - Follow instructions to create your bot
   - Copy the Bot Token

2. Get your Chat ID:
   - Send a message to your bot
   - Visit: `https://api.telegram.org/bot<YOUR_BOT_TOKEN>/getUpdates`
   - Look for `"chat":{"id":123456789}`
   - Copy the Chat ID

3. Add to Exit Debt settings

**API Endpoint:**
```http
PUT /api/users/settings/notifications
{
  "notification_webhook": true,
  "telegram_bot_token": "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11",
  "telegram_chat_id": "123456789"
}
```

### Setting Up Discord

1. Go to your Discord server
2. Right-click on the channel → Edit Channel
3. Navigate to: Integrations → Webhooks
4. Click "New Webhook"
5. Customize name and avatar (optional)
6. Copy the Webhook URL
7. Add to Exit Debt settings

**API Endpoint:**
```http
PUT /api/users/settings/notifications
{
  "notification_webhook": true,
  "discord_webhook_url": "https://discord.com/api/webhooks/123456789/XXXXXXXXXXXXXXXXXXXX"
}
```

---

## API Examples

### Create Webhook Notification

```http
POST /api/notifications/send
{
  "debt_list_id": "123e4567-e89b-12d3-a456-426614174000",
  "notification_type": "webhook",
  "webhook_type": "slack",
  "message": "Manual reminder"
}
```

### Schedule Webhook Notifications for Debt

```http
POST /api/debt-lists/:id/schedule-notifications
{
  "notification_channels": ["email", "webhook"],
  "webhook_types": ["slack", "telegram"]
}
```

This will create notifications for both email AND webhooks (Slack + Telegram).

### Get Webhook Configuration

```http
GET /api/users/settings/notifications

Response:
{
  "notification_webhook": true,
  "slack_webhook_url": "https://hooks.slack.com/services/...",
  "telegram_bot_token": "123456:ABC-DEF...",
  "telegram_chat_id": "123456789",
  "discord_webhook_url": "https://discord.com/api/webhooks/..."
}
```

---

## Testing

### Test Slack Webhook

```bash
curl -X POST https://hooks.slack.com/services/YOUR/WEBHOOK/URL \
  -H 'Content-Type: application/json' \
  -d '{
    "text": "Test notification from Exit Debt!",
    "blocks": [
      {
        "type": "section",
        "text": {
          "type": "mrkdwn",
          "text": "*Test Payment Reminder*\nAmount: $1,000\nDue: Jan 15, 2025"
        }
      }
    ]
  }'
```

### Test Telegram Bot

```bash
curl -X POST https://api.telegram.org/bot<BOT_TOKEN>/sendMessage \
  -H 'Content-Type: application/json' \
  -d '{
    "chat_id": "YOUR_CHAT_ID",
    "text": "Test notification from Exit Debt!\n\n*Payment Reminder*\nAmount: $1,000",
    "parse_mode": "Markdown"
  }'
```

### Test Discord Webhook

```bash
curl -X POST https://discord.com/api/webhooks/YOUR/WEBHOOK/URL \
  -H 'Content-Type: application/json' \
  -d '{
    "content": "Test notification from Exit Debt!",
    "embeds": [{
      "title": "Payment Reminder",
      "description": "Amount: $1,000",
      "color": 3447003
    }]
  }'
```

---

## Security Considerations

### 1. Webhook URL Storage
- Store webhook URLs encrypted in the database
- Never expose webhook URLs in API responses (mask them)
- Validate webhook URLs before saving

### 2. Rate Limiting
- Implement per-platform rate limits
- Slack: 1 message per second
- Telegram: 30 messages per second per chat
- Discord: 5 requests per second per webhook

### 3. Validation
- Validate webhook URLs format
- Test webhooks before saving
- Handle invalid configurations gracefully

### 4. Error Handling
- Retry failed webhook deliveries
- Log webhook errors separately
- Notify users of misconfigured webhooks

---

## Monitoring & Logs

```go
// Log webhook events
log.Printf("[WEBHOOK] Slack notification sent for debt %s", debtID)
log.Printf("[WEBHOOK] Telegram notification failed: %v", err)
log.Printf("[WEBHOOK] Discord webhook rate limit exceeded")
```

### Metrics to Track
- Webhook success rate per platform
- Average delivery time
- Failed webhook attempts
- Rate limit hits

---

## Benefits of Webhook Notifications

| Benefit | Description |
|---------|-------------|
| **Real-time** | Instant notifications in messaging apps |
| **Centralized** | All notifications in one place (Slack/Discord channel) |
| **Team Visibility** | Share debt reminders with team/family |
| **Rich Formatting** | Buttons, cards, and interactive elements |
| **Cross-platform** | Works on desktop and mobile |
| **No Email Fatigue** | Separate from email inbox |

---

## Troubleshooting

### Slack Webhook Not Working
- Check webhook URL is correct
- Ensure channel still exists
- Verify app is still installed
- Check payload format

### Telegram Bot Not Responding
- Verify bot token is valid
- Ensure chat ID is correct
- Check bot is not blocked
- Verify bot has permission to send messages

### Discord Webhook Failing
- Check webhook URL is valid
- Ensure webhook wasn't deleted
- Verify server permissions
- Check rate limits

---

This completes the webhook notification integration guide. Users can now receive payment reminders via Slack, Telegram, or Discord in addition to email and SMS!

