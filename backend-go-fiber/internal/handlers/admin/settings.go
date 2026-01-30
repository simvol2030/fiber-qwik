package admin

import (
	"backend-go-fiber/internal/models"
	"backend-go-fiber/internal/services/admin"
	"backend-go-fiber/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type SettingsHandler struct {
	service *admin.SettingsService
}

func NewSettingsHandler(service *admin.SettingsService) *SettingsHandler {
	return &SettingsHandler{service: service}
}

// GetAll returns all settings
// GET /api/admin/settings
func (h *SettingsHandler) GetAll(c *fiber.Ctx) error {
	// Check if group filter is provided
	group := c.Query("group", "")

	var settings []models.AppSettingsResponse
	var err error

	if group != "" {
		settings, err = h.service.GetByGroup(group)
	} else {
		settings, err = h.service.GetAll()
	}

	if err != nil {
		return utils.SendError(c, "INTERNAL_ERROR", "Failed to get settings", fiber.StatusInternalServerError)
	}

	return utils.SendSuccess(c, settings, fiber.StatusOK)
}

// Get returns a single setting by key
// GET /api/admin/settings/:key
func (h *SettingsHandler) Get(c *fiber.Ctx) error {
	key := c.Params("key")
	if key == "" {
		return utils.SendError(c, "VALIDATION_ERROR", "Setting key is required", fiber.StatusBadRequest)
	}

	setting, err := h.service.GetByKey(key)
	if err != nil {
		if err.Error() == "setting not found" {
			return utils.SendError(c, "NOT_FOUND", "Setting not found", fiber.StatusNotFound)
		}
		return utils.SendError(c, "INTERNAL_ERROR", "Failed to get setting", fiber.StatusInternalServerError)
	}

	return utils.SendSuccess(c, setting.ToResponse(), fiber.StatusOK)
}

// UpdateInput represents the input for updating a single setting
type UpdateInput struct {
	Value string `json:"value"`
}

// Update updates a single setting
// PUT /api/admin/settings/:key
func (h *SettingsHandler) Update(c *fiber.Ctx) error {
	key := c.Params("key")
	if key == "" {
		return utils.SendError(c, "VALIDATION_ERROR", "Setting key is required", fiber.StatusBadRequest)
	}

	var input UpdateInput
	if err := c.BodyParser(&input); err != nil {
		return utils.SendError(c, "VALIDATION_ERROR", "Invalid request body", fiber.StatusBadRequest)
	}

	setting, err := h.service.Update(key, input.Value)
	if err != nil {
		if err.Error() == "setting not found" {
			return utils.SendError(c, "NOT_FOUND", "Setting not found", fiber.StatusNotFound)
		}
		return utils.SendError(c, "INTERNAL_ERROR", "Failed to update setting", fiber.StatusInternalServerError)
	}

	return utils.SendSuccess(c, setting.ToResponse(), fiber.StatusOK)
}

// UpdateBatch updates multiple settings at once
// PUT /api/admin/settings
func (h *SettingsHandler) UpdateBatch(c *fiber.Ctx) error {
	var input admin.UpdateBatchInput
	if err := c.BodyParser(&input); err != nil {
		return utils.SendError(c, "VALIDATION_ERROR", "Invalid request body", fiber.StatusBadRequest)
	}

	// Validate input
	if errors := utils.ValidateStruct(input); errors != nil {
		return utils.SendValidationError(c, errors)
	}

	settings, err := h.service.UpdateBatch(input.Settings)
	if err != nil {
		return utils.SendError(c, "INTERNAL_ERROR", "Failed to update settings", fiber.StatusInternalServerError)
	}

	return utils.SendSuccess(c, settings, fiber.StatusOK)
}
