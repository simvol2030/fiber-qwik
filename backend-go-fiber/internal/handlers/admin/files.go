package admin

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"backend-go-fiber/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type FilesHandler struct {
	uploadDir string
}

func NewFilesHandler(uploadDir string) *FilesHandler {
	return &FilesHandler{uploadDir: uploadDir}
}

// FileInfo represents file information
type FileInfo struct {
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	Size      int64     `json:"size"`
	IsDir     bool      `json:"isDir"`
	ModTime   time.Time `json:"modTime"`
	Extension string    `json:"extension"`
	MimeType  string    `json:"mimeType"`
}

// ListFilesResult contains the list of files
type ListFilesResult struct {
	Files      []FileInfo `json:"files"`
	Total      int        `json:"total"`
	TotalSize  int64      `json:"totalSize"`
	CurrentDir string     `json:"currentDir"`
}

// List returns list of uploaded files
// GET /api/admin/files
func (h *FilesHandler) List(c *fiber.Ctx) error {
	// Get subdirectory from query (optional)
	subDir := c.Query("dir", "")

	// Sanitize path to prevent directory traversal
	subDir = filepath.Clean(subDir)
	if strings.Contains(subDir, "..") {
		return utils.SendError(c, "VALIDATION_ERROR", "Invalid directory path", fiber.StatusBadRequest)
	}

	targetDir := h.uploadDir
	if subDir != "" && subDir != "." {
		targetDir = filepath.Join(h.uploadDir, subDir)
	}

	// Check if directory exists
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		return utils.SendSuccess(c, ListFilesResult{
			Files:      []FileInfo{},
			Total:      0,
			TotalSize:  0,
			CurrentDir: subDir,
		}, fiber.StatusOK)
	}

	// Read directory
	entries, err := os.ReadDir(targetDir)
	if err != nil {
		return utils.SendError(c, "INTERNAL_ERROR", "Failed to read directory", fiber.StatusInternalServerError)
	}

	files := make([]FileInfo, 0, len(entries))
	var totalSize int64

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		relPath := entry.Name()
		if subDir != "" && subDir != "." {
			relPath = filepath.Join(subDir, entry.Name())
		}

		ext := ""
		mimeType := ""
		if !entry.IsDir() {
			ext = strings.TrimPrefix(filepath.Ext(entry.Name()), ".")
			mimeType = getMimeType(ext)
			totalSize += info.Size()
		}

		files = append(files, FileInfo{
			Name:      entry.Name(),
			Path:      relPath,
			Size:      info.Size(),
			IsDir:     entry.IsDir(),
			ModTime:   info.ModTime(),
			Extension: ext,
			MimeType:  mimeType,
		})
	}

	return utils.SendSuccess(c, ListFilesResult{
		Files:      files,
		Total:      len(files),
		TotalSize:  totalSize,
		CurrentDir: subDir,
	}, fiber.StatusOK)
}

// Delete deletes a file
// DELETE /api/admin/files/*
func (h *FilesHandler) Delete(c *fiber.Ctx) error {
	// Get file path from params
	filePath := c.Params("*")
	if filePath == "" {
		return utils.SendError(c, "VALIDATION_ERROR", "File path is required", fiber.StatusBadRequest)
	}

	// Sanitize path
	filePath = filepath.Clean(filePath)
	if strings.Contains(filePath, "..") {
		return utils.SendError(c, "VALIDATION_ERROR", "Invalid file path", fiber.StatusBadRequest)
	}

	fullPath := filepath.Join(h.uploadDir, filePath)

	// Check if file exists
	info, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		return utils.SendError(c, "NOT_FOUND", "File not found", fiber.StatusNotFound)
	}

	// Delete file or directory
	if info.IsDir() {
		err = os.RemoveAll(fullPath)
	} else {
		err = os.Remove(fullPath)
	}

	if err != nil {
		return utils.SendError(c, "INTERNAL_ERROR", "Failed to delete file", fiber.StatusInternalServerError)
	}

	return utils.SendSuccess(c, fiber.Map{"message": "File deleted successfully"}, fiber.StatusOK)
}

// getMimeType returns MIME type based on file extension
func getMimeType(ext string) string {
	mimeTypes := map[string]string{
		"jpg":  "image/jpeg",
		"jpeg": "image/jpeg",
		"png":  "image/png",
		"gif":  "image/gif",
		"webp": "image/webp",
		"svg":  "image/svg+xml",
		"pdf":  "application/pdf",
		"doc":  "application/msword",
		"docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"xls":  "application/vnd.ms-excel",
		"xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		"zip":  "application/zip",
		"txt":  "text/plain",
		"json": "application/json",
		"mp4":  "video/mp4",
		"webm": "video/webm",
		"mp3":  "audio/mpeg",
	}

	if mime, ok := mimeTypes[strings.ToLower(ext)]; ok {
		return mime
	}
	return "application/octet-stream"
}
