package handlers

import (
	"backend-go-fiber/internal/services/upload"
	"backend-go-fiber/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// UploadHandler handles file upload requests
type UploadHandler struct {
	uploadService *upload.Service
}

// NewUploadHandler creates a new upload handler
func NewUploadHandler(uploadService *upload.Service) *UploadHandler {
	return &UploadHandler{
		uploadService: uploadService,
	}
}

// UploadSingle handles single file upload
// POST /api/upload
func (h *UploadHandler) UploadSingle(c *fiber.Ctx) error {
	// Get file from form
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return utils.SendError(c, "VALIDATION_ERROR", "No file provided", fiber.StatusBadRequest)
	}

	// Upload file
	info, err := h.uploadService.UploadFile(c.Context(), fileHeader)
	if err != nil {
		return utils.SendError(c, "UPLOAD_ERROR", err.Error(), fiber.StatusBadRequest)
	}

	return utils.SendSuccess(c, fiber.Map{
		"file": info,
	}, fiber.StatusCreated)
}

// UploadMultiple handles multiple file upload
// POST /api/upload/multiple
func (h *UploadHandler) UploadMultiple(c *fiber.Ctx) error {
	// Get multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return utils.SendError(c, "VALIDATION_ERROR", "Invalid form data", fiber.StatusBadRequest)
	}

	files := form.File["files"]
	if len(files) == 0 {
		return utils.SendError(c, "VALIDATION_ERROR", "No files provided", fiber.StatusBadRequest)
	}

	// Limit number of files
	const maxFiles = 10
	if len(files) > maxFiles {
		return utils.SendError(c, "VALIDATION_ERROR", "Too many files (max: 10)", fiber.StatusBadRequest)
	}

	// Upload all files
	var results []fiber.Map
	var errors []fiber.Map

	for _, fileHeader := range files {
		info, err := h.uploadService.UploadFile(c.Context(), fileHeader)
		if err != nil {
			errors = append(errors, fiber.Map{
				"filename": fileHeader.Filename,
				"error":    err.Error(),
			})
			continue
		}
		results = append(results, fiber.Map{
			"filename": fileHeader.Filename,
			"file":     info,
		})
	}

	return utils.SendSuccess(c, fiber.Map{
		"uploaded": results,
		"errors":   errors,
	}, fiber.StatusCreated)
}

// Delete handles file deletion
// DELETE /api/upload/:key
func (h *UploadHandler) Delete(c *fiber.Ctx) error {
	key := c.Params("*") // Use wildcard to capture full path

	if key == "" {
		return utils.SendError(c, "VALIDATION_ERROR", "File key is required", fiber.StatusBadRequest)
	}

	if err := h.uploadService.DeleteFile(c.Context(), key); err != nil {
		return utils.SendError(c, "DELETE_ERROR", err.Error(), fiber.StatusInternalServerError)
	}

	return utils.SendSuccess(c, fiber.Map{
		"message": "File deleted successfully",
	})
}
