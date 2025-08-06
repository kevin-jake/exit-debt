# Debt Tracker API - Postman Collection

This Postman collection provides comprehensive testing for all endpoints of the Debt Tracker API. The collection includes authentication, debt management, and special endpoint tests with automated validation.

## 📋 Table of Contents

- [Setup Instructions](#setup-instructions)
- [Collection Structure](#collection-structure)
- [Test Coverage](#test-coverage)
- [Running Tests](#running-tests)
- [Variables](#variables)
- [Test Scenarios](#test-scenarios)

## 🚀 Setup Instructions

### Prerequisites

1. **Install Postman**: Download and install [Postman](https://www.postman.com/downloads/)
2. **Start the API Server**: Ensure your Debt Tracker API is running on `http://localhost:8080`
3. **Database Setup**: Make sure your PostgreSQL database is running and migrations are applied

### Import Collection

1. Open Postman
2. Click "Import" button
3. Select the `Debt_Tracker_API.postman_collection.json` file
4. The collection will be imported with all endpoints and tests

## 📁 Collection Structure

The collection is organized into the following folders:

### 1. Authentication

- **Register User**: Creates a new user account
- **Login User**: Authenticates and retrieves JWT token

### 2. Health Check

- **Health Check (Public)**: Tests the public health endpoint
- **Health Check (Protected)**: Tests the protected health endpoint

### 3. Debt Lists

- **Create Debt List**: Creates a new debt list
- **Get User Debt Lists**: Retrieves all debt lists for the authenticated user
- **Get Specific Debt List**: Retrieves a specific debt list by ID
- **Update Debt List**: Updates an existing debt list
- **Delete Debt List**: Deletes a debt list

### 4. Debt Items

- **Create Debt Item**: Creates a new debt item within a debt list
- **Get Specific Debt Item**: Retrieves a specific debt item by ID
- **Get Debt List Items**: Retrieves all items in a specific debt list
- **Update Debt Item**: Updates an existing debt item
- **Delete Debt Item**: Deletes a debt item

### 5. Special Endpoints

- **Get Overdue Items**: Retrieves all overdue debt items
- **Get Due Soon Items**: Retrieves debt items due within specified days

## 🧪 Test Coverage

### Authentication Tests

- ✅ User registration validation
- ✅ Login token generation
- ✅ Response structure validation
- ✅ Variable storage for subsequent requests

### Debt List Tests

- ✅ CRUD operations validation
- ✅ Response structure validation
- ✅ Data integrity checks
- ✅ Authorization validation

### Debt Item Tests

- ✅ CRUD operations validation
- ✅ Relationship validation (debt list association)
- ✅ Amount and currency validation
- ✅ Status updates validation

### Special Endpoint Tests

- ✅ Overdue items filtering
- ✅ Due soon items with parameter validation
- ✅ Response array structure validation

## 🏃‍♂️ Running Tests

### Manual Testing

1. **Start with Authentication**:

   - Run "Register User" to create a test account
   - Run "Login User" to get authentication token

2. **Test Debt Lists**:

   - Run "Create Debt List" to create a test debt list
   - Run "Get User Debt Lists" to verify retrieval
   - Run "Get Specific Debt List" to test individual retrieval
   - Run "Update Debt List" to test modifications
   - Run "Delete Debt List" to test deletion

3. **Test Debt Items**:

   - Run "Create Debt Item" to add items to the debt list
   - Run "Get Debt List Items" to verify item retrieval
   - Run "Update Debt Item" to test item modifications
   - Run "Delete Debt Item" to test item deletion

4. **Test Special Endpoints**:
   - Run "Get Overdue Items" to test overdue filtering
   - Run "Get Due Soon Items" to test due soon filtering

### Automated Testing

1. **Run Entire Collection**:

   - Click the collection name
   - Click "Run collection"
   - Select all requests or specific folders
   - Click "Run Debt Tracker API"

2. **Run Specific Folders**:
   - Right-click on a folder
   - Select "Run folder"
   - Choose test options and run

## 🔧 Variables

The collection uses the following variables:

| Variable       | Description              | Auto-populated |
| -------------- | ------------------------ | -------------- |
| `base_url`     | API base URL             | ❌ Manual      |
| `auth_token`   | JWT authentication token | ✅ Auto        |
| `user_id`      | Current user ID          | ✅ Auto        |
| `debt_list_id` | Current debt list ID     | ✅ Auto        |
| `debt_item_id` | Current debt item ID     | ✅ Auto        |

### Variable Flow

1. **Registration**: Sets `user_id`
2. **Login**: Sets `auth_token` and `user_id`
3. **Create Debt List**: Sets `debt_list_id`
4. **Create Debt Item**: Sets `debt_item_id`

## 🎯 Test Scenarios

### Scenario 1: Complete User Workflow

1. Register a new user
2. Login to get authentication token
3. Create a debt list
4. Add multiple debt items
5. Update debt items
6. Check overdue and due soon items
7. Delete items and lists

### Scenario 2: Error Handling

1. Test invalid authentication
2. Test invalid request data
3. Test unauthorized access
4. Test resource not found scenarios

### Scenario 3: Data Validation

1. Test required field validation
2. Test data type validation
3. Test business rule validation
4. Test relationship validation

## 📊 Test Results

Each test includes validation for:

- **HTTP Status Codes**: Correct response codes (200, 201, 400, 401, 404, 500)
- **Response Structure**: Proper JSON structure and required fields
- **Data Integrity**: Valid data types and business rules
- **Authentication**: Proper token validation and user authorization
- **Error Handling**: Appropriate error messages and codes

## 🔍 Troubleshooting

### Common Issues

1. **Authentication Errors**:

   - Ensure you've run the Login request first
   - Check that the `auth_token` variable is set
   - Verify the token hasn't expired

2. **404 Errors**:

   - Ensure the API server is running on `http://localhost:8080`
   - Check that the base_url variable is correct
   - Verify the endpoint paths are correct

3. **400 Bad Request Errors**:

   - Check request body format (JSON)
   - Verify required fields are provided
   - Ensure data types are correct

4. **Database Errors**:
   - Ensure PostgreSQL is running
   - Verify database migrations are applied
   - Check database connection settings

### Debug Tips

1. **Check Variables**: Use the "Environment" tab to verify variable values
2. **Console Logs**: Check the Postman console for detailed error messages
3. **Response Headers**: Verify Content-Type and other headers
4. **Request Body**: Ensure JSON is properly formatted

## 📝 Customization

### Adding New Tests

1. **Create New Request**:

   - Right-click on appropriate folder
   - Select "Add request"
   - Configure method, URL, and headers

2. **Add Test Scripts**:

   - Use the "Tests" tab
   - Write JavaScript validation
   - Use `pm.test()` for assertions

3. **Set Variables**:
   - Use `pm.collectionVariables.set()` to store values
   - Use `{{variable_name}}` to reference variables

### Modifying Existing Tests

1. **Update Test Scripts**:

   - Edit the "Tests" tab
   - Modify validation logic
   - Add new assertions

2. **Update Request Data**:
   - Modify request body
   - Update URL parameters
   - Change headers as needed

## 🚀 Best Practices

1. **Run Tests in Order**: Follow the folder structure for proper test flow
2. **Check Variables**: Verify variables are set before running dependent tests
3. **Monitor Console**: Watch for errors and warnings in the console
4. **Validate Responses**: Ensure all response validations pass
5. **Clean Up**: Delete test data after testing if needed

## 📞 Support

For issues with the API or collection:

1. Check the API documentation in `/docs/API.md`
2. Verify server logs for detailed error messages
3. Ensure all prerequisites are met
4. Test with curl commands to isolate issues

---

**Happy Testing! 🎉**
