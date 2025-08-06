# Debt Tracker API Documentation

## Base URL

```
http://localhost:8080
```

## Authentication

All protected endpoints require a JWT token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## Endpoints

### Authentication

#### Register User

```http
POST /api/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123",
  "first_name": "John",
  "last_name": "Doe",
  "phone": "+1234567890"
}
```

#### Login User

```http
POST /api/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

### Debt Lists

#### Create Debt List

```http
POST /api/debt-lists
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Personal Loans",
  "description": "Money lent to friends and family",
  "currency": "Php"
}
```

#### Get User's Debt Lists

```http
GET /api/debt-lists
Authorization: Bearer <token>
```

#### Get Specific Debt List

```http
GET /api/debt-lists/{id}
Authorization: Bearer <token>
```

#### Update Debt List

```http
PUT /api/debt-lists/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Updated Title",
  "description": "Updated description",
  "currency": "EUR",
  "status": "archived"
}
```

#### Delete Debt List

```http
DELETE /api/debt-lists/{id}
Authorization: Bearer <token>
```

### Debt Items

#### Create Debt Item

```http
POST /api/debt-items
Authorization: Bearer <token>
Content-Type: application/json

{
  "debt_list_id": "uuid-of-debt-list",
  "debtor_name": "John Smith",
  "debtor_email": "john@example.com",
  "debtor_phone": "+1234567890",
  "debtor_facebook_id": "john.smith.123",
  "amount": "150.50",
  "currency": "Php",
  "description": "Lunch money",
  "due_date": "2024-02-15T00:00:00Z",
  "debt_type": "owed_to_me"
}
```

#### Get Specific Debt Item

```http
GET /api/debt-items/{id}
Authorization: Bearer <token>
```

#### Get Debt Items in a List

```http
GET /api/debt-lists/{debtListId}/items
Authorization: Bearer <token>
```

#### Update Debt Item

```http
PUT /api/debt-items/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "debtor_name": "Updated Name",
  "amount": "200.00",
  "status": "paid",
  "due_date": "2024-03-01T00:00:00Z"
}
```

#### Delete Debt Item

```http
DELETE /api/debt-items/{id}
Authorization: Bearer <token>
```

### Special Endpoints

#### Get Overdue Items

```http
GET /api/debt-items/overdue
Authorization: Bearer <token>
```

#### Get Due Soon Items

```http
GET /api/debt-items/due-soon?days=7
Authorization: Bearer <token>
```

## Data Models

### Debt List

```json
{
  "id": "uuid",
  "user_id": "uuid",
  "title": "string",
  "description": "string",
  "total_amount": "0.00",
  "currency": "Php",
  "status": "active",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### Debt Item

```json
{
  "id": "uuid",
  "debt_list_id": "uuid",
  "debtor_name": "string",
  "debtor_email": "string",
  "debtor_phone": "string",
  "debtor_facebook_id": "string",
  "amount": "150.50",
  "currency": "Php",
  "description": "string",
  "due_date": "2024-02-15T00:00:00Z",
  "status": "pending",
  "debt_type": "owed_to_me",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

## Status Values

### Debt List Status

- `active` - Currently active debt list
- `archived` - Archived debt list
- `settled` - All debts in list are settled

### Debt Item Status

- `pending` - Debt is pending payment
- `paid` - Debt has been paid
- `overdue` - Debt is past due date
- `cancelled` - Debt has been cancelled

### Debt Types

- `owed_to_me` - Money owed to the user
- `i_owe` - Money the user owes to others

## Error Responses

All endpoints return errors in the following format:

```json
{
  "error": "Error message description"
}
```

Common HTTP status codes:

- `200` - Success
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `404` - Not Found
- `500` - Internal Server Error

## Example Usage

### Complete Workflow

1. **Register a user:**

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe"
  }'
```

2. **Login to get token:**

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

3. **Create a debt list:**

```bash
curl -X POST http://localhost:8080/api/debt-lists \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Personal Loans",
    "description": "Money lent to friends",
    "currency": "Php"
  }'
```

4. **Add a debt item:**

```bash
curl -X POST http://localhost:8080/api/debt-items \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "debt_list_id": "DEBT_LIST_UUID",
    "debtor_name": "Alice Johnson",
    "debtor_email": "alice@example.com",
    "amount": "100.00",
    "debt_type": "owed_to_me",
    "description": "Lunch money",
    "due_date": "2024-02-15T00:00:00Z"
  }'
```

5. **Get all debt lists:**

```bash
curl -X GET http://localhost:8080/api/debt-lists \
  -H "Authorization: Bearer YOUR_TOKEN"
```

6. **Get overdue items:**

```bash
curl -X GET http://localhost:8080/api/debt-items/overdue \
  -H "Authorization: Bearer YOUR_TOKEN"
```
