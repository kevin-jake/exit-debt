package mocks

import (
	"context"
	"io"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockFileStorageService is a mock implementation of FileStorageService
type MockFileStorageService struct {
	mock.Mock
}

func (m *MockFileStorageService) UploadReceipt(ctx context.Context, file io.Reader, filename string, contentType string, debtID uuid.UUID) (string, error) {
	args := m.Called(ctx, file, filename, contentType, debtID)
	return args.String(0), args.Error(1)
}

func (m *MockFileStorageService) DeleteReceipt(ctx context.Context, fileURL string) error {
	args := m.Called(ctx, fileURL)
	return args.Error(0)
}

func (m *MockFileStorageService) GetReceiptURL(ctx context.Context, fileURL string) (string, error) {
	args := m.Called(ctx, fileURL)
	return args.String(0), args.Error(1)
}

func (m *MockFileStorageService) GetReceiptFile(ctx context.Context, fileURL string) ([]byte, string, error) {
	args := m.Called(ctx, fileURL)
	return args.Get(0).([]byte), args.String(1), args.Error(2)
}
