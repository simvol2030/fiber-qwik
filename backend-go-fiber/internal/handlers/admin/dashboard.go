package admin

import (
	"backend-go-fiber/internal/services/admin"
	"backend-go-fiber/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type DashboardHandler struct {
	service *admin.DashboardService
}

func NewDashboardHandler(service *admin.DashboardService) *DashboardHandler {
	return &DashboardHandler{service: service}
}

// GetStats returns dashboard statistics
// GET /api/admin/dashboard
func (h *DashboardHandler) GetStats(c *fiber.Ctx) error {
	stats, err := h.service.GetStats()
	if err != nil {
		return utils.SendError(c, "INTERNAL_ERROR", "Failed to get dashboard stats", fiber.StatusInternalServerError)
	}

	return utils.SendSuccess(c, stats, fiber.StatusOK)
}
