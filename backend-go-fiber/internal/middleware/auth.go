package middleware

import (
	"strings"

	"backend-go-fiber/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return utils.SendError(c, "UNAUTHORIZED", "Missing or invalid authorization header", fiber.StatusUnauthorized)
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		payload, err := utils.VerifyAccessToken(token)

		if err != nil {
			return utils.SendError(c, "UNAUTHORIZED", "Invalid or expired token", fiber.StatusUnauthorized)
		}

		c.Locals("user", payload)
		return c.Next()
	}
}
