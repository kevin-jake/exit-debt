package mocks

import (
	"context"
	"io"

	"github.com/stretchr/testify/mock"
)

// MockFileStorageService is a mock implementation of FileStorageService
type MockFileStorageService struct {
	mock.Mock
}

func (m *MockFileStorageService) UploadReceipt(ctx context.Context, file io.Reader, filename string, contentType string) (string, error) {
	args := m.Called(ctx, file, filename, contentType)
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
