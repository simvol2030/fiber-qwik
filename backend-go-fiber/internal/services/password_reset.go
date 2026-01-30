package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os"
	"time"

	"backend-go-fiber/internal/models"
	"backend-go-fiber/internal/services/email"
	"backend-go-fiber/internal/utils"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

var (
	ErrInvalidResetToken = errors.New("invalid or expired reset token")
	ErrUserNotFound      = errors.New("user not found")
)

// PasswordResetService handles password reset logic
type PasswordResetService struct {
	db          *gorm.DB
	emailSender email.Sender
}

// NewPasswordResetService creates a new password reset service
func NewPasswordResetService(db *gorm.DB, emailSender email.Sender) *PasswordResetService {
	return &PasswordResetService{
		db:          db,
		emailSender: emailSender,
	}
}

// RequestReset creates a reset token and sends email
// Returns nil even if user doesn't exist (security: don't reveal if email exists)
func (s *PasswordResetService) RequestReset(ctx context.Context, emailAddr string) error {
	// Find user by email
	var user models.User
	if err := s.db.Where("email = ?", emailAddr).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Don't reveal if user exists - just log and return success
			log.Debug().Str("email", emailAddr).Msg("Password reset requested for non-existent user")
			return nil
		}
		return err
	}

	// Generate secure random token
	token, err := generateSecureToken(32)
	if err != nil {
		return err
	}

	// Create reset token (valid for 1 hour)
	resetToken := &models.PasswordResetToken{
		Token:     token,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	// Invalidate any existing tokens for this user
	s.db.Where("user_id = ? AND used_at IS NULL", user.ID).Delete(&models.PasswordResetToken{})

	// Save new token
	if err := s.db.Create(resetToken).Error; err != nil {
		return err
	}

	// Build reset URL
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}
	resetURL := frontendURL + "/reset-password?token=" + token

	// Send email
	if err := s.emailSender.SendTemplate(ctx, []string{user.Email}, email.TemplatePasswordReset, map[string]interface{}{
		"ResetURL":  resetURL,
		"ExpiresIn": "1 hour",
		"Name":      user.Name,
	}); err != nil {
		log.Error().Err(err).Str("email", emailAddr).Msg("Failed to send password reset email")
		// Don't return error to user - token was created successfully
	}

	log.Info().Str("email", emailAddr).Msg("Password reset token created")
	return nil
}

// ValidateToken checks if a reset token is valid
func (s *PasswordResetService) ValidateToken(ctx context.Context, token string) (*models.User, error) {
	var resetToken models.PasswordResetToken
	err := s.db.Preload("User").Where("token = ?", token).First(&resetToken).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidResetToken
		}
		return nil, err
	}

	if !resetToken.IsValid() {
		return nil, ErrInvalidResetToken
	}

	return &resetToken.User, nil
}

// ResetPassword validates token and updates password
func (s *PasswordResetService) ResetPassword(ctx context.Context, token, newPassword string) error {
	// Find and validate token
	var resetToken models.PasswordResetToken
	err := s.db.Preload("User").Where("token = ?", token).First(&resetToken).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrInvalidResetToken
		}
		return err
	}

	if !resetToken.IsValid() {
		return ErrInvalidResetToken
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Start transaction
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Update password
		if err := tx.Model(&models.User{}).Where("id = ?", resetToken.UserID).Update("password_hash", hashedPassword).Error; err != nil {
			return err
		}

		// Mark token as used
		now := time.Now()
		if err := tx.Model(&resetToken).Update("used_at", &now).Error; err != nil {
			return err
		}

		// Invalidate all refresh tokens (force re-login)
		if err := tx.Where("user_id = ?", resetToken.UserID).Delete(&models.RefreshToken{}).Error; err != nil {
			return err
		}

		return nil
	})
}

// CleanupExpiredTokens removes expired tokens (call periodically)
func (s *PasswordResetService) CleanupExpiredTokens(ctx context.Context) error {
	result := s.db.Where("expires_at < ?", time.Now()).Delete(&models.PasswordResetToken{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected > 0 {
		log.Debug().Int64("count", result.RowsAffected).Msg("Cleaned up expired password reset tokens")
	}

	return nil
}

// generateSecureToken generates a cryptographically secure random token
func generateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
