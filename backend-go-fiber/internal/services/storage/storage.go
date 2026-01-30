package storage

import (
	"context"
	"io"
	"time"
)

// FileInfo represents metadata about an uploaded file
type FileInfo struct {
	// Key is the unique identifier/path of the file in storage
	Key string `json:"key"`

	// OriginalName is the original filename from upload
	OriginalName string `json:"originalName"`

	// Size in bytes
	Size int64 `json:"size"`

	// ContentType is the MIME type
	ContentType string `json:"contentType"`

	// URL is the public/signed URL to access the file
	URL string `json:"url"`

	// CreatedAt is when the file was uploaded
	CreatedAt time.Time `json:"createdAt"`
}

// Storage defines the interface for file storage backends
// Implement this interface to add support for different storage providers
type Storage interface {
	// Upload stores a file and returns its info
	Upload(ctx context.Context, key string, reader io.Reader, size int64, contentType string) (*FileInfo, error)

	// Download retrieves a file by its key
	Download(ctx context.Context, key string) (io.ReadCloser, error)

	// Delete removes a file by its key
	Delete(ctx context.Context, key string) error

	// GetURL returns the URL to access the file
	// For local storage, this returns a relative path
	// For S3, this can return a signed URL
	GetURL(ctx context.Context, key string) (string, error)

	// Exists checks if a file exists
	Exists(ctx context.Context, key string) (bool, error)
}

// Config holds storage configuration
type Config struct {
	// Type is the storage backend: "local" or "s3"
	Type string

	// Local storage settings
	LocalPath string // Path to store files (e.g., "./data/uploads")
	BaseURL   string // Base URL for serving files (e.g., "/uploads")

	// S3 storage settings
	S3Bucket    string
	S3Region    string
	S3Endpoint  string // For S3-compatible services (MinIO, DigitalOcean Spaces)
	S3AccessKey string
	S3SecretKey string
}
