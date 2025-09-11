package unit

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"exit-debt/internal/config"
	"exit-debt/internal/services"
)

// MockS3Client is a mock implementation of the S3 client
type MockS3Client struct {
	mock.Mock
}

func (m *MockS3Client) PutObject(ctx context.Context, params interface{}, optFns ...func(*interface{})) (interface{}, error) {
	args := m.Called(ctx, params)
	return args.Get(0), args.Error(1)
}

func (m *MockS3Client) DeleteObject(ctx context.Context, params interface{}, optFns ...func(*interface{})) (interface{}, error) {
	args := m.Called(ctx, params)
	return args.Get(0), args.Error(1)
}

func (m *MockS3Client) HeadBucket(ctx context.Context, params interface{}, optFns ...func(*interface{})) (interface{}, error) {
	args := m.Called(ctx, params)
	return args.Get(0), args.Error(1)
}

func TestS3Service_ValidateFile(t *testing.T) {
	// Create a minimal S3Service for testing validation methods
	s3Service := &services.S3Service{}

	tests := []struct {
		name        string
		filename    string
		contentType string
		size        int64
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid JPEG file",
			filename:    "receipt.jpg",
			contentType: "image/jpeg",
			size:        1024 * 1024, // 1MB
			expectError: false,
		},
		{
			name:        "valid PNG file",
			filename:    "receipt.png",
			contentType: "image/png",
			size:        2 * 1024 * 1024, // 2MB
			expectError: false,
		},
		{
			name:        "file too large",
			filename:    "receipt.jpg",
			contentType: "image/jpeg",
			size:        15 * 1024 * 1024, // 15MB
			expectError: true,
			errorMsg:    "exceeds maximum allowed size",
		},
		{
			name:        "invalid content type",
			filename:    "receipt.pdf",
			contentType: "application/pdf",
			size:        1024 * 1024, // 1MB
			expectError: true,
			errorMsg:    "invalid file type",
		},
		{
			name:        "invalid file extension",
			filename:    "receipt.txt",
			contentType: "image/jpeg",
			size:        1024 * 1024, // 1MB
			expectError: true,
			errorMsg:    "invalid file extension",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s3Service.ValidateFile(tt.filename, tt.contentType, tt.size)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestS3Service_ExtractKeyFromURL(t *testing.T) {
	// Create a minimal S3Service for testing URL parsing methods
	s3Service := &services.S3Service{}

	tests := []struct {
		name        string
		fileURL     string
		expectKey   string
		expectError bool
	}{
		{
			name:        "valid API path format",
			fileURL:     "/api/v1/debts/123e4567-e89b-12d3-a456-426614174000/receipts/20240115-143022-abc123.jpg",
			expectKey:   "123e4567-e89b-12d3-a456-426614174000/receipts/20240115-143022-abc123.jpg",
			expectError: false,
		},
		{
			name:        "valid S3 URL",
			fileURL:     "s3://my-bucket/receipts/2024/01/15/abc123.jpg",
			expectKey:   "receipts/2024/01/15/abc123.jpg",
			expectError: false,
		},
		{
			name:        "valid HTTPS S3 URL",
			fileURL:     "https://my-bucket.s3.us-east-1.amazonaws.com/receipts/2024/01/15/abc123.jpg",
			expectKey:   "receipts/2024/01/15/abc123.jpg",
			expectError: false,
		},
		{
			name:        "invalid API path format - missing filename",
			fileURL:     "/api/v1/debts/123e4567-e89b-12d3-a456-426614174000/receipts/",
			expectError: true,
		},
		{
			name:        "invalid API path format - invalid debt ID",
			fileURL:     "/api/v1/debts/invalid-uuid/receipts/20240115-143022-abc123.jpg",
			expectError: true,
		},
		{
			name:        "invalid S3 URL format",
			fileURL:     "s3://my-bucket",
			expectError: true,
		},
		{
			name:        "unsupported URL format",
			fileURL:     "ftp://example.com/file.jpg",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := s3Service.ExtractKeyFromURL(tt.fileURL)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, key)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectKey, key)
			}
		})
	}
}

func TestS3Service_IsValidImageType(t *testing.T) {
	// Create a minimal S3Service for testing image type validation
	s3Service := &services.S3Service{}

	tests := []struct {
		name        string
		contentType string
		expected    bool
	}{
		{"JPEG", "image/jpeg", true},
		{"JPG", "image/jpg", true},
		{"PNG", "image/png", true},
		{"GIF", "image/gif", true},
		{"WebP", "image/webp", true},
		{"PDF", "application/pdf", false},
		{"Text", "text/plain", false},
		{"Empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := s3Service.IsValidImageType(tt.contentType)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestS3Service_NewS3Service(t *testing.T) {
	tests := []struct {
		name        string
		config      *config.Config
		expectError bool
		errorMsg    string
	}{
		{
			name: "missing bucket name",
			config: &config.Config{
				S3Region:          "us-east-1",
				S3AccessKeyID:     "test-key",
				S3SecretAccessKey: "test-secret",
			},
			expectError: true,
			errorMsg:    "S3_BUCKET_NAME is required",
		},
		{
			name: "missing access key",
			config: &config.Config{
				S3Region:          "us-east-1",
				S3BucketName:      "test-bucket",
				S3SecretAccessKey: "test-secret",
			},
			expectError: true,
			errorMsg:    "S3_ACCESS_KEY_ID and S3_SECRET_ACCESS_KEY are required",
		},
		{
			name: "missing secret key",
			config: &config.Config{
				S3Region:      "us-east-1",
				S3BucketName:  "test-bucket",
				S3AccessKeyID: "test-key",
			},
			expectError: true,
			errorMsg:    "S3_ACCESS_KEY_ID and S3_SECRET_ACCESS_KEY are required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := zerolog.New(io.Discard)
			_, err := services.NewS3Service(tt.config, logger)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			}
		})
	}
}

// Test helper function to create a test config
func createTestConfig() *config.Config {
	return &config.Config{
		S3Region:          "us-east-1",
		S3BucketName:      "test-bucket",
		S3AccessKeyID:     "test-key",
		S3SecretAccessKey: "test-secret",
	}
}

// Test helper function to create a test file reader
func createTestFileReader(content string) io.Reader {
	return strings.NewReader(content)
}
