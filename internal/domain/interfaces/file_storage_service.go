package interfaces

import (
	"context"
	"io"
)

// FileStorageService defines the interface for file storage operations
type FileStorageService interface {
	// UploadReceipt uploads a receipt photo and returns the URL
	UploadReceipt(ctx context.Context, file io.Reader, filename string, contentType string) (string, error)
	
	// DeleteReceipt deletes a receipt photo from storage
	DeleteReceipt(ctx context.Context, fileURL string) error
	
	// GetReceiptURL generates a secure, time-limited URL for receipt access
	GetReceiptURL(ctx context.Context, fileURL string) (string, error)
}
