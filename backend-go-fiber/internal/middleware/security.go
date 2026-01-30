package middleware

import (
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/google/uuid"
)

func CORSMiddleware() fiber.Handler {
	originsStr := os.Getenv("CORS_ORIGINS")
	if originsStr == "" {
		originsStr = "http://localhost:3000"
	}

	return cors.New(cors.Config{
		AllowOrigins:     originsStr,
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
	})
}

func HelmetMiddleware() fiber.Handler {
	// Build CSP based on environment
	// Production: More relaxed to allow CDN, fonts, analytics
	// Development: Stricter for testing
	var csp string
	if os.Getenv("NODE_ENV") == "production" {
		// Production CSP: Allow common CDNs, Google Fonts, analytics
		csp = strings.Join([]string{
			"default-src 'self'",
			"script-src 'self' 'unsafe-inline' https://www.googletagmanager.com https://www.google-analytics.com",
			"style-src 'self' 'unsafe-inline' https://fonts.googleapis.com",
			"img-src 'self' data: https: blob:",
			"connect-src 'self' https://www.google-analytics.com https://analytics.google.com",
			"font-src 'self' https://fonts.gstatic.com",
			"object-src 'none'",
			"media-src 'self'",
			"frame-src 'none'",
			"base-uri 'self'",
			"form-action 'self'",
		}, "; ")
	} else {
		// Development CSP: Stricter
		csp = strings.Join([]string{
			"default-src 'self'",
			"script-src 'self' 'unsafe-inline'", // unsafe-inline for Vite HMR
			"style-src 'self' 'unsafe-inline'",
			"img-src 'self' data: https:",
			"connect-src 'self' ws: wss:", // WebSocket for Vite HMR
			"font-src 'self'",
			"object-src 'none'",
			"media-src 'self'",
			"frame-src 'none'",
		}, "; ")
	}

	return helmet.New(helmet.Config{
		XSSProtection:             "1; mode=block",
		ContentTypeNosniff:        "nosniff",
		XFrameOptions:             "DENY",
		ReferrerPolicy:            "strict-origin-when-cross-origin",
		CrossOriginEmbedderPolicy: "require-corp",
		CrossOriginOpenerPolicy:   "same-origin",
		CrossOriginResourcePolicy: "same-origin",
		ContentSecurityPolicy:     csp,
	})
}

func RateLimiterMiddleware() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "RATE_LIMIT_EXCEEDED",
					"message": "Too many requests, please try again later",
				},
			})
		},
	})
}

// LoginRateLimiter limits login attempts: 5 per 5 minutes per IP
func LoginRateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        5,
		Expiration: 5 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return "login:" + c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "RATE_LIMIT_EXCEEDED",
					"message": "Too many login attempts. Try again in 5 minutes.",
				},
			})
		},
	})
}

// RegisterRateLimiter limits registration attempts: 3 per hour per IP
func RegisterRateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        3,
		Expiration: 1 * time.Hour,
		KeyGenerator: func(c *fiber.Ctx) string {
			return "register:" + c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "RATE_LIMIT_EXCEEDED",
					"message": "Too many registration attempts. Try again later.",
				},
			})
		},
	})
}

func RequestIDMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestID := uuid.New().String()
		c.Locals("requestId", requestID)
		c.Set("X-Request-ID", requestID)
		return c.Next()
	}
}
