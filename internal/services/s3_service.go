package services

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"exit-debt/internal/config"
)

// S3Service implements the FileStorageService interface
type S3Service struct {
	s3Client   *s3.Client
	bucketName string
	logger     zerolog.Logger
}

// NewS3Service creates a new S3 service instance
func NewS3Service(cfg *config.Config, logger zerolog.Logger) (*S3Service, error) {
	// Validate required S3 configuration
	if cfg.S3BucketName == "" {
		return nil, fmt.Errorf("S3_BUCKET_NAME is required")
	}
	if cfg.S3AccessKeyID == "" || cfg.S3SecretAccessKey == "" {
		return nil, fmt.Errorf("S3_ACCESS_KEY_ID and S3_SECRET_ACCESS_KEY are required")
	}

	// Create AWS config
	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(cfg.S3Region),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.S3AccessKeyID,
			cfg.S3SecretAccessKey,
			"",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Set custom endpoint if provided (useful for local development with MinIO)
	if cfg.S3Endpoint != "" {
		awsCfg.EndpointResolverWithOptions = aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:               cfg.S3Endpoint,
				SigningRegion:     cfg.S3Region,
				HostnameImmutable: true,
			}, nil
		})
	}

	// Create S3 client
	s3Client := s3.NewFromConfig(awsCfg)

	// Test bucket access
	_, err = s3Client.HeadBucket(context.Background(), &s3.HeadBucketInput{
		Bucket: aws.String(cfg.S3BucketName),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to access S3 bucket %s: %w", cfg.S3BucketName, err)
	}

	return &S3Service{
		s3Client:   s3Client,
		bucketName: cfg.S3BucketName,
		logger:     logger,
	}, nil
}

// UploadReceipt uploads a receipt photo and returns the relative path
func (s *S3Service) UploadReceipt(ctx context.Context, file io.Reader, filename string, contentType string) (string, error) {
	// Validate file type
	if !s.IsValidImageType(contentType) {
		return "", fmt.Errorf("invalid file type: %s. Only images are allowed", contentType)
	}

	// Generate unique filename
	ext := filepath.Ext(filename)
	uniqueFilename := fmt.Sprintf("receipts/%s/%s%s", 
		time.Now().Format("2006/01/02"), 
		uuid.New().String(), 
		ext)

	// Upload file to S3
	_, err := s.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(uniqueFilename),
		Body:        file,
		ContentType: aws.String(contentType),
		Metadata: map[string]string{
			"original-filename": filename,
			"uploaded-at":       time.Now().Format(time.RFC3339),
		},
	})
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to upload receipt to S3")
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	// Return relative path instead of S3 URL
	relativePath := fmt.Sprintf("/api/v1/receipts/%s", uniqueFilename)
	s.logger.Info().Str("filename", uniqueFilename).Str("relative_path", relativePath).Msg("Receipt uploaded successfully to S3")

	return relativePath, nil
}

// DeleteReceipt deletes a receipt photo from S3
func (s *S3Service) DeleteReceipt(ctx context.Context, fileURL string) error {
	// Extract key from S3 URL
	key, err := s.ExtractKeyFromURL(fileURL)
	if err != nil {
		return fmt.Errorf("invalid S3 URL: %w", err)
	}

	// Delete object from S3
	_, err = s.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		s.logger.Error().Err(err).Str("key", key).Msg("Failed to delete receipt from S3")
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}

	s.logger.Info().Str("key", key).Msg("Receipt deleted successfully from S3")
	return nil
}

// GetReceiptURL generates a secure, time-limited URL for receipt access
func (s *S3Service) GetReceiptURL(ctx context.Context, fileURL string) (string, error) {
	// Extract key from S3 URL
	key, err := s.ExtractKeyFromURL(fileURL)
	if err != nil {
		return "", fmt.Errorf("invalid S3 URL: %w", err)
	}

	// Generate presigned URL (valid for 1 hour)
	presignClient := s3.NewPresignClient(s.s3Client)
	req, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(1*time.Hour))
	if err != nil {
		s.logger.Error().Err(err).Str("key", key).Msg("Failed to generate presigned URL")
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	s.logger.Info().Str("key", key).Time("expires", time.Now().Add(1*time.Hour)).Msg("Generated presigned URL for receipt")
	return req.URL, nil
}

// GetReceiptFile retrieves a receipt file from S3 and returns the file content and metadata
func (s *S3Service) GetReceiptFile(ctx context.Context, fileURL string) ([]byte, string, error) {
	// Extract key from relative path or S3 URL
	key, err := s.ExtractKeyFromURL(fileURL)
	if err != nil {
		return nil, "", fmt.Errorf("invalid file path: %w", err)
	}

	// Get object from S3
	result, err := s.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		s.logger.Error().Err(err).Str("key", key).Msg("Failed to retrieve receipt from S3")
		return nil, "", fmt.Errorf("failed to retrieve file from S3: %w", err)
	}
	defer result.Body.Close()

	// Read file content
	fileContent, err := io.ReadAll(result.Body)
	if err != nil {
		s.logger.Error().Err(err).Str("key", key).Msg("Failed to read receipt content from S3")
		return nil, "", fmt.Errorf("failed to read file content: %w", err)
	}

	// Get content type from S3 metadata
	contentType := "application/octet-stream" // default
	if result.ContentType != nil {
		contentType = *result.ContentType
	}

	s.logger.Info().Str("key", key).Str("content_type", contentType).Int("size", len(fileContent)).Msg("Retrieved receipt file from S3")
	return fileContent, contentType, nil
}

// IsValidImageType checks if the content type is a valid image type
func (s *S3Service) IsValidImageType(contentType string) bool {
	validTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}
	return validTypes[contentType]
}

// ExtractKeyFromURL extracts the S3 key from a relative path or S3 URL
func (s *S3Service) ExtractKeyFromURL(fileURL string) (string, error) {
	// Handle relative path format: /api/v1/receipts/receipts/2024/01/15/uuid.jpg
	if strings.HasPrefix(fileURL, "/api/v1/receipts/") {
		// Remove the /api/v1/receipts/ prefix to get the S3 key
		key := strings.TrimPrefix(fileURL, "/api/v1/receipts/")
		if key == "" {
			return "", fmt.Errorf("invalid relative path format: %s", fileURL)
		}
		return key, nil
	}
	
	// Handle S3 URL format: s3://bucket-name/key
	if strings.HasPrefix(fileURL, "s3://") {
		parts := strings.SplitN(fileURL, "/", 4)
		if len(parts) < 4 {
			return "", fmt.Errorf("invalid S3 URL format: %s", fileURL)
		}
		return parts[3], nil
	}
	
	// Handle HTTPS S3 URL format: https://bucket-name.s3.region.amazonaws.com/key
	if strings.HasPrefix(fileURL, "https://") {
		parts := strings.Split(fileURL, "/")
		if len(parts) < 4 {
			return "", fmt.Errorf("invalid S3 HTTPS URL format: %s", fileURL)
		}
		return strings.Join(parts[3:], "/"), nil
	}
	
	return "", fmt.Errorf("unsupported URL format: %s", fileURL)
}

// ValidateFile validates if the uploaded file is acceptable
func (s *S3Service) ValidateFile(filename string, contentType string, size int64) error {
	// Check file size (max 10MB)
	const maxSize = 10 * 1024 * 1024 // 10MB
	if size > maxSize {
		return fmt.Errorf("file size %d bytes exceeds maximum allowed size of %d bytes", size, maxSize)
	}

	// Check file type
	if !s.IsValidImageType(contentType) {
		return fmt.Errorf("invalid file type: %s. Only images (JPEG, PNG, GIF, WebP) are allowed", contentType)
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(filename))
	validExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}
	if !validExtensions[ext] {
		return fmt.Errorf("invalid file extension: %s. Only .jpg, .jpeg, .png, .gif, .webp are allowed", ext)
	}

	return nil
}
