package services

import (
	"errors"
	"time"

	"backend-go-fiber/internal/models"
	"backend-go-fiber/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db: db}
}

type RegisterInput struct {
	Email    string  `json:"email" validate:"required,email"`
	Password string  `json:"password" validate:"required,min=8,max=128"`
	Name     *string `json:"name" validate:"omitempty,max=100"`
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResult struct {
	User        models.UserResponse `json:"user"`
	AccessToken string              `json:"accessToken"`
	ExpiresIn   int                 `json:"expiresIn"`
}

func (s *AuthService) Register(input RegisterInput) (*AuthResult, error) {
	// Check if user exists
	var existing models.User
	if err := s.db.Where("email = ?", input.Email).First(&existing).Error; err == nil {
		return nil, errors.New("user already exists")
	}

	// Hash password
	passwordHash, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := models.User{
		Email:        input.Email,
		PasswordHash: passwordHash,
		Name:         input.Name,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}

	// Generate access token
	accessToken, err := utils.GenerateAccessToken(utils.JWTPayload{
		UserID: user.ID,
		Email:  user.Email,
	})
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		User:        user.ToResponse(),
		AccessToken: accessToken,
		ExpiresIn:   utils.GetExpiresInSeconds(),
	}, nil
}

func (s *AuthService) Login(input LoginInput) (*AuthResult, error) {
	var user models.User
	if err := s.db.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !utils.VerifyPassword(input.Password, user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	// Update last login timestamp
	now := time.Now()
	user.LastLoginAt = &now
	s.db.Model(&user).Update("last_login_at", now)

	accessToken, err := utils.GenerateAccessToken(utils.JWTPayload{
		UserID: user.ID,
		Email:  user.Email,
	})
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		User:        user.ToResponse(),
		AccessToken: accessToken,
		ExpiresIn:   utils.GetExpiresInSeconds(),
	}, nil
}

func (s *AuthService) CreateRefreshToken(userID string) (string, error) {
	token := uuid.New().String()
	expiresAt := time.Now().AddDate(0, 0, utils.GetRefreshTokenExpiresDays())

	refreshToken := models.RefreshToken{
		Token:     token,
		UserID:    userID,
		ExpiresAt: expiresAt,
	}

	if err := s.db.Create(&refreshToken).Error; err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) RefreshAccessToken(refreshToken string) (*struct {
	AccessToken string `json:"accessToken"`
	ExpiresIn   int    `json:"expiresIn"`
}, error) {
	var storedToken models.RefreshToken
	if err := s.db.Preload("User").Where("token = ?", refreshToken).First(&storedToken).Error; err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if storedToken.ExpiresAt.Before(time.Now()) {
		s.db.Delete(&storedToken)
		return nil, errors.New("refresh token expired")
	}

	accessToken, err := utils.GenerateAccessToken(utils.JWTPayload{
		UserID: storedToken.User.ID,
		Email:  storedToken.User.Email,
	})
	if err != nil {
		return nil, err
	}

	return &struct {
		AccessToken string `json:"accessToken"`
		ExpiresIn   int    `json:"expiresIn"`
	}{
		AccessToken: accessToken,
		ExpiresIn:   utils.GetExpiresInSeconds(),
	}, nil
}

func (s *AuthService) RevokeRefreshToken(token string) error {
	return s.db.Where("token = ?", token).Delete(&models.RefreshToken{}).Error
}

func (s *AuthService) GetUserByID(userID string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateProfileInput represents the profile update request
type UpdateProfileInput struct {
	Name *string `json:"name" validate:"omitempty,max=100"`
}

// UpdateProfile updates the user's name
func (s *AuthService) UpdateProfile(userID string, input UpdateProfileInput) (*models.User, error) {
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	user.Name = input.Name
	if err := s.db.Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// ChangePasswordInput represents the password change request
type ChangePasswordInput struct {
	CurrentPassword string `json:"currentPassword" validate:"required"`
	NewPassword     string `json:"newPassword" validate:"required,min=8,max=128"`
}

// ChangePassword changes the user's password after verifying the current one
func (s *AuthService) ChangePassword(userID string, input ChangePasswordInput) error {
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return errors.New("user not found")
	}

	if !utils.VerifyPassword(input.CurrentPassword, user.PasswordHash) {
		return errors.New("current password is incorrect")
	}

	newHash, err := utils.HashPassword(input.NewPassword)
	if err != nil {
		return err
	}

	user.PasswordHash = newHash
	if err := s.db.Save(&user).Error; err != nil {
		return err
	}

	return nil
}
