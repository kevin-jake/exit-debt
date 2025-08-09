# Exit-Debt - a Go Debt Tracker Application

A comprehensive debt tracking application built with Go, PostgreSQL, and Gin framework. Track your debts, money lent to others, and send notifications via email, SMS, or Facebook Messenger.

## Features

- **User Authentication**: Secure JWT-based authentication
- **Debt Management**: Track money you owe and money owed to you
- **Multiple Debt Lists**: Organize debts into different categories
- **Notifications**: Send reminders via email, SMS, or Facebook Messenger
- **User Settings**: Customize notification preferences and default currency
- **RESTful API**: Clean and well-documented API endpoints

## Tech Stack

- **Backend**: Go 1.21+
- **Framework**: Gin (HTTP web framework)
- **Database**: PostgreSQL
- **Authentication**: JWT tokens
- **Password Hashing**: bcrypt
- **Logging**: Zerolog
- **Hot Reloading**: Air

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Air (for development hot reloading)

## Installation

### 1. Clone the repository

```bash
git clone <repository-url>
cd exit-debt
```

### 2. Install dependencies

```bash
go mod download
```

### 3. Install Air (for development)

```bash
go install github.com/cosmtrek/air@latest
```

### 4. Set up PostgreSQL

Create a database for the application:

```sql
CREATE DATABASE debt_tracker;
```

### 5. Configure environment variables

Copy the example environment file and update it with your settings:

```bash
cp env.example .env
```

Edit `.env` with your database credentials:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=debt_tracker
DB_SSL_MODE=disable

SERVER_PORT=8080
SERVER_HOST=localhost

JWT_SECRET=your-secret-key-here
JWT_EXPIRY=24h

LOG_LEVEL=debug
```

## Running the Application

### Development (with hot reloading)

```bash
air
```

### Production

```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`

## API Endpoints

### Authentication

- `POST /api/auth/register` - Register a new user
- `POST /api/auth/login` - Login user

### Protected Endpoints (require JWT token)

- `GET /api/health` - Health check

## Database Schema

The application uses the following main tables:

- **users**: User accounts and authentication
- **user_settings**: User preferences and notification settings
- **debt_lists**: Collections of debt items
- **debt_items**: Individual debt records
- **notifications**: Notification history and tracking

## Development

### Project Structure

```
├── cmd/
│   └── server/          # Main application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── handler/         # HTTP handlers
│   ├── middleware/      # HTTP middleware
│   ├── models/          # Data models
│   ├── database/        # Database connection and GORM setup
│   └── service/         # Business logic
├── pkg/                 # Public packages
└── scripts/             # Utility scripts
```

### Running Tests

```bash
go test ./...
```

### Database Schema

The application uses GORM's auto-migration feature to automatically create and update the database schema based on the model definitions. No manual migrations are required.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License.
