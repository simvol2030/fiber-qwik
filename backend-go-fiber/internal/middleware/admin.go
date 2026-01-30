package middleware

import (
	"backend-go-fiber/internal/models"
	"backend-go-fiber/internal/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// AdminOnly middleware ensures that only users with admin role can access the route
func AdminOnly(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user payload from context (set by AuthMiddleware)
		payload, ok := c.Locals("user").(*utils.JWTPayload)
		if !ok || payload == nil {
			return utils.SendError(c, "UNAUTHORIZED", "Authentication required", fiber.StatusUnauthorized)
		}

		// Fetch user from database to check role
		var user models.User
		if err := db.First(&user, "id = ?", payload.UserID).Error; err != nil {
			return utils.SendError(c, "UNAUTHORIZED", "User not found", fiber.StatusUnauthorized)
		}

		// Check if user is admin
		if !user.IsAdmin() {
			return utils.SendError(c, "FORBIDDEN", "Admin access required", fiber.StatusForbidden)
		}

		// Store full user in context for handlers
		c.Locals("adminUser", &user)

		return c.Next()
	}
}
