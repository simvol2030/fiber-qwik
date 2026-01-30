package handlers

import (
	"backend-go-fiber/internal/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db *gorm.DB
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) Health(c *fiber.Ctx) error {
	return utils.SendSuccess(c, fiber.Map{
		"status":    "ok",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *HealthHandler) Ready(c *fiber.Ctx) error {
	// Check database connection
	sqlDB, err := h.db.DB()
	if err != nil {
		return utils.SendError(c, "NOT_READY", "Service is not ready", fiber.StatusServiceUnavailable)
	}

	if err := sqlDB.Ping(); err != nil {
		return utils.SendError(c, "NOT_READY", "Service is not ready", fiber.StatusServiceUnavailable)
	}

	return utils.SendSuccess(c, fiber.Map{
		"status": "ready",
		"checks": fiber.Map{
			"database": "ok",
		},
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}
