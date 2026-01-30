package handlers

import (
	"backend-go-fiber/internal/services"
	"backend-go-fiber/internal/utils"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

const refreshTokenCookie = "refresh_token"

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var input services.RegisterInput
	if err := c.BodyParser(&input); err != nil {
		return utils.SendError(c, "VALIDATION_ERROR", "Invalid request body", fiber.StatusBadRequest)
	}

	// Validate input using go-playground/validator
	validationErrors := utils.ValidateStruct(input)
	if utils.HasValidationErrors(validationErrors) {
		fieldErrors := make([]utils.FieldError, len(validationErrors))
		for i, ve := range validationErrors {
			fieldErrors[i] = utils.FieldError{
				Field:   ve.Field,
				Message: ve.Message,
			}
		}
		return utils.SendError(c, "VALIDATION_ERROR", "Validation failed", fiber.StatusBadRequest, fieldErrors)
	}

	result, err := h.authService.Register(input)
	if err != nil {
		if err.Error() == "user already exists" {
			return utils.SendError(c, "USER_EXISTS", err.Error(), fiber.StatusConflict)
		}
		return utils.SendError(c, "INTERNAL_ERROR", "Registration failed", fiber.StatusInternalServerError)
	}

	// Create refresh token
	refreshToken, err := h.authService.CreateRefreshToken(result.User.ID)
	if err != nil {
		return utils.SendError(c, "INTERNAL_ERROR", "Failed to create refresh token", fiber.StatusInternalServerError)
	}

	h.setRefreshTokenCookie(c, refreshToken)
	return utils.SendSuccess(c, result, fiber.StatusCreated)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var input services.LoginInput
	if err := c.BodyParser(&input); err != nil {
		return utils.SendError(c, "VALIDATION_ERROR", "Invalid request body", fiber.StatusBadRequest)
	}

	// Validate input using go-playground/validator
	validationErrors := utils.ValidateStruct(input)
	if utils.HasValidationErrors(validationErrors) {
		fieldErrors := make([]utils.FieldError, len(validationErrors))
		for i, ve := range validationErrors {
			fieldErrors[i] = utils.FieldError{
				Field:   ve.Field,
				Message: ve.Message,
			}
		}
		return utils.SendError(c, "VALIDATION_ERROR", "Validation failed", fiber.StatusBadRequest, fieldErrors)
	}

	result, err := h.authService.Login(input)
	if err != nil {
		if err.Error() == "invalid credentials" {
			return utils.SendError(c, "INVALID_CREDENTIALS", err.Error(), fiber.StatusUnauthorized)
		}
		return utils.SendError(c, "INTERNAL_ERROR", "Login failed", fiber.StatusInternalServerError)
	}

	// Create refresh token
	refreshToken, err := h.authService.CreateRefreshToken(result.User.ID)
	if err != nil {
		return utils.SendError(c, "INTERNAL_ERROR", "Failed to create refresh token", fiber.StatusInternalServerError)
	}

	h.setRefreshTokenCookie(c, refreshToken)
	return utils.SendSuccess(c, result)
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	refreshToken := c.Cookies(refreshTokenCookie)
	if refreshToken == "" {
		return utils.SendError(c, "NO_REFRESH_TOKEN", "No refresh token provided", fiber.StatusUnauthorized)
	}

	result, err := h.authService.RefreshAccessToken(refreshToken)
	if err != nil {
		h.clearRefreshTokenCookie(c)
		return utils.SendError(c, "INVALID_REFRESH_TOKEN", "Invalid or expired refresh token", fiber.StatusUnauthorized)
	}

	return utils.SendSuccess(c, result)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	refreshToken := c.Cookies(refreshTokenCookie)
	if refreshToken != "" {
		h.authService.RevokeRefreshToken(refreshToken)
	}

	h.clearRefreshTokenCookie(c)
	return utils.SendSuccess(c, fiber.Map{"message": "Logged out successfully"})
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	userPayload := c.Locals("user").(*utils.JWTPayload)

	user, err := h.authService.GetUserByID(userPayload.UserID)
	if err != nil {
		return utils.SendError(c, "USER_NOT_FOUND", "User not found", fiber.StatusNotFound)
	}

	return utils.SendSuccess(c, user.ToResponse())
}

func (h *AuthHandler) UpdateProfile(c *fiber.Ctx) error {
	userPayload := c.Locals("user").(*utils.JWTPayload)

	var input services.UpdateProfileInput
	if err := c.BodyParser(&input); err != nil {
		return utils.SendError(c, "VALIDATION_ERROR", "Invalid request body", fiber.StatusBadRequest)
	}

	validationErrors := utils.ValidateStruct(input)
	if utils.HasValidationErrors(validationErrors) {
		return utils.SendValidationError(c, validationErrors)
	}

	user, err := h.authService.UpdateProfile(userPayload.UserID, input)
	if err != nil {
		return utils.SendError(c, "INTERNAL_ERROR", "Failed to update profile", fiber.StatusInternalServerError)
	}

	return utils.SendSuccess(c, user.ToResponse())
}

func (h *AuthHandler) ChangePassword(c *fiber.Ctx) error {
	userPayload := c.Locals("user").(*utils.JWTPayload)

	var input services.ChangePasswordInput
	if err := c.BodyParser(&input); err != nil {
		return utils.SendError(c, "VALIDATION_ERROR", "Invalid request body", fiber.StatusBadRequest)
	}

	validationErrors := utils.ValidateStruct(input)
	if utils.HasValidationErrors(validationErrors) {
		return utils.SendValidationError(c, validationErrors)
	}

	err := h.authService.ChangePassword(userPayload.UserID, input)
	if err != nil {
		if err.Error() == "current password is incorrect" {
			return utils.SendError(c, "INVALID_PASSWORD", err.Error(), fiber.StatusBadRequest)
		}
		return utils.SendError(c, "INTERNAL_ERROR", "Failed to change password", fiber.StatusInternalServerError)
	}

	return utils.SendSuccess(c, fiber.Map{"message": "Password changed successfully"})
}

func (h *AuthHandler) setRefreshTokenCookie(c *fiber.Ctx, token string) {
	secure := os.Getenv("NODE_ENV") == "production"
	maxAge := utils.GetRefreshTokenExpiresDays() * 24 * 60 * 60

	// COOKIE_SAMESITE: Lax (default), None (Telegram WebApp), Strict
	// Use "None" for embedded contexts (Telegram WebApp, iframes)
	sameSite := os.Getenv("COOKIE_SAMESITE")
	if sameSite != "None" && sameSite != "Strict" {
		sameSite = "Lax"
	}
	// SameSite=None requires Secure=true
	if sameSite == "None" {
		secure = true
	}

	c.Cookie(&fiber.Cookie{
		Name:     refreshTokenCookie,
		Value:    token,
		HTTPOnly: true,
		Secure:   secure,
		SameSite: sameSite,
		MaxAge:   maxAge,
		Path:     "/",
		Expires:  time.Now().Add(time.Duration(maxAge) * time.Second),
	})
}

func (h *AuthHandler) clearRefreshTokenCookie(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     refreshTokenCookie,
		Value:    "",
		HTTPOnly: true,
		MaxAge:   -1,
		Path:     "/",
		Expires:  time.Now().Add(-time.Hour),
	})
}
