package upload

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"backend-go-fiber/internal/services/storage"

	"github.com/google/uuid"
)

// Service handles file uploads with validation
type Service struct {
	storage storage.Storage
	config  Config
}

// Config defines upload service configuration
type Config struct {
	// MaxFileSize in bytes (default: 10MB)
	MaxFileSize int64

	// AllowedMimeTypes is a list of allowed MIME types
	// Empty means all types are allowed
	AllowedMimeTypes []string

	// AllowedExtensions is a list of allowed file extensions (with dot, e.g., ".jpg")
	// Empty means all extensions are allowed
	AllowedExtensions []string

	// PathPrefix is prepended to generated file keys (e.g., "avatars/")
	PathPrefix string
}

// DefaultConfig returns sensible defaults for file uploads
func DefaultConfig() Config {
	return Config{
		MaxFileSize: 10 * 1024 * 1024, // 10MB
		AllowedMimeTypes: []string{
			"image/jpeg",
			"image/png",
			"image/gif",
			"image/webp",
			"application/pdf",
		},
		AllowedExtensions: []string{
			".jpg", ".jpeg", ".png", ".gif", ".webp", ".pdf",
		},
	}
}

// ImageOnlyConfig returns config for image uploads only
func ImageOnlyConfig() Config {
	return Config{
		MaxFileSize: 5 * 1024 * 1024, // 5MB
		AllowedMimeTypes: []string{
			"image/jpeg",
			"image/png",
			"image/gif",
			"image/webp",
		},
		AllowedExtensions: []string{
			".jpg", ".jpeg", ".png", ".gif", ".webp",
		},
	}
}

// NewService creates a new upload service
func NewService(storage storage.Storage, config Config) *Service {
	return &Service{
		storage: storage,
		config:  config,
	}
}

// UploadFile handles a single file upload from multipart form
func (s *Service) UploadFile(ctx context.Context, fileHeader *multipart.FileHeader) (*storage.FileInfo, error) {
	// Validate file size
	if s.config.MaxFileSize > 0 && fileHeader.Size > s.config.MaxFileSize {
		return nil, fmt.Errorf("file too large: %d bytes (max: %d)", fileHeader.Size, s.config.MaxFileSize)
	}

	// Get content type
	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Validate MIME type
	if len(s.config.AllowedMimeTypes) > 0 {
		allowed := false
		for _, mime := range s.config.AllowedMimeTypes {
			if strings.EqualFold(contentType, mime) {
				allowed = true
				break
			}
		}
		if !allowed {
			return nil, fmt.Errorf("file type not allowed: %s", contentType)
		}
	}

	// Validate extension
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if len(s.config.AllowedExtensions) > 0 {
		allowed := false
		for _, allowedExt := range s.config.AllowedExtensions {
			if strings.EqualFold(ext, allowedExt) {
				allowed = true
				break
			}
		}
		if !allowed {
			return nil, fmt.Errorf("file extension not allowed: %s", ext)
		}
	}

	// Generate unique key
	key := s.generateKey(fileHeader.Filename)

	// Open uploaded file
	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer file.Close()

	// Upload to storage
	info, err := s.storage.Upload(ctx, key, file, fileHeader.Size, contentType)
	if err != nil {
		return nil, err
	}

	// Preserve original filename
	info.OriginalName = fileHeader.Filename

	return info, nil
}

// DeleteFile removes a file from storage
func (s *Service) DeleteFile(ctx context.Context, key string) error {
	return s.storage.Delete(ctx, key)
}

// GetFileURL returns the URL for accessing a file
func (s *Service) GetFileURL(ctx context.Context, key string) (string, error) {
	return s.storage.GetURL(ctx, key)
}

// generateKey creates a unique file key with timestamp and UUID
func (s *Service) generateKey(originalName string) string {
	ext := filepath.Ext(originalName)
	timestamp := time.Now().Format("2006/01/02")
	uniqueID := uuid.New().String()

	key := fmt.Sprintf("%s/%s%s", timestamp, uniqueID, ext)

	if s.config.PathPrefix != "" {
		key = s.config.PathPrefix + key
	}

	return key
}
