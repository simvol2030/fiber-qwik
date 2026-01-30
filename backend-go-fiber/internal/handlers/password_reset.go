package handlers

import (
	"errors"

	"backend-go-fiber/internal/services"
	"backend-go-fiber/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// PasswordResetHandler handles password reset requests
type PasswordResetHandler struct {
	service *services.PasswordResetService
}

// NewPasswordResetHandler creates a new password reset handler
func NewPasswordResetHandler(service *services.PasswordResetService) *PasswordResetHandler {
	return &PasswordResetHandler{
		service: service,
	}
}

// ForgotPasswordRequest represents the forgot password request body
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ResetPasswordRequest represents the reset password request body
type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required,min=8"`
}

// ValidateTokenRequest represents the validate token request body
type ValidateTokenRequest struct {
	Token string `json:"token" validate:"required"`
}

// ForgotPassword handles POST /api/auth/forgot-password
// Initiates password reset process by sending email with reset link
func (h *PasswordResetHandler) ForgotPassword(c *fiber.Ctx) error {
	var req ForgotPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, "VALIDATION_ERROR", "Invalid request body", fiber.StatusBadRequest)
	}

	// Validate request
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return utils.SendValidationError(c, validationErrors)
	}

	// Request reset (always returns success for security)
	if err := h.service.RequestReset(c.Context(), req.Email); err != nil {
		// Log error but don't expose to user
		return utils.SendError(c, "INTERNAL_ERROR", "Failed to process request", fiber.StatusInternalServerError)
	}

	// Always return success (don't reveal if email exists)
	return utils.SendSuccess(c, fiber.Map{
		"message": "If an account with that email exists, a password reset link has been sent.",
	})
}

// ValidateToken handles POST /api/auth/validate-reset-token
// Checks if a reset token is valid (used by frontend before showing reset form)
func (h *PasswordResetHandler) ValidateToken(c *fiber.Ctx) error {
	var req ValidateTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, "VALIDATION_ERROR", "Invalid request body", fiber.StatusBadRequest)
	}

	// Validate request
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return utils.SendValidationError(c, validationErrors)
	}

	// Validate token
	user, err := h.service.ValidateToken(c.Context(), req.Token)
	if err != nil {
		if errors.Is(err, services.ErrInvalidResetToken) {
			return utils.SendError(c, "INVALID_TOKEN", "Invalid or expired reset token", fiber.StatusBadRequest)
		}
		return utils.SendError(c, "INTERNAL_ERROR", "Failed to validate token", fiber.StatusInternalServerError)
	}

	return utils.SendSuccess(c, fiber.Map{
		"valid": true,
		"email": user.Email, // Show which email the token belongs to
	})
}

// ResetPassword handles POST /api/auth/reset-password
// Resets password using valid token
func (h *PasswordResetHandler) ResetPassword(c *fiber.Ctx) error {
	var req ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, "VALIDATION_ERROR", "Invalid request body", fiber.StatusBadRequest)
	}

	// Validate request
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return utils.SendValidationError(c, validationErrors)
	}

	// Reset password
	if err := h.service.ResetPassword(c.Context(), req.Token, req.NewPassword); err != nil {
		if errors.Is(err, services.ErrInvalidResetToken) {
			return utils.SendError(c, "INVALID_TOKEN", "Invalid or expired reset token", fiber.StatusBadRequest)
		}
		return utils.SendError(c, "INTERNAL_ERROR", "Failed to reset password", fiber.StatusInternalServerError)
	}

	return utils.SendSuccess(c, fiber.Map{
		"message": "Password has been reset successfully. Please login with your new password.",
	})
}
