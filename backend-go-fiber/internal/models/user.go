package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Role represents user role for basic RBAC
type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

// User represents the user model with soft delete support
type User struct {
	ID           string `gorm:"primaryKey;type:text"`
	Email        string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	Name         *string
	Role         Role           `gorm:"type:text;default:user;not null"`
	IsActive     bool           `gorm:"default:true;not null"`  // Account active status
	LastLoginAt  *time.Time     `json:"lastLoginAt"`            // Last login timestamp
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"` // Soft delete support

	RefreshTokens []RefreshToken `gorm:"foreignKey:UserID"`
}

// IsAdmin checks if user has admin role
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

type RefreshToken struct {
	ID        string `gorm:"primaryKey;type:text"`
	Token     string `gorm:"uniqueIndex;not null"`
	UserID    string `gorm:"not null"`
	User      User   `gorm:"constraint:OnDelete:CASCADE"`
	ExpiresAt time.Time
	CreatedAt time.Time
}

func (r *RefreshToken) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	return nil
}

// UserResponse is the response format for user data
type UserResponse struct {
	ID          string     `json:"id"`
	Email       string     `json:"email"`
	Name        *string    `json:"name"`
	Role        Role       `json:"role"`
	IsActive    bool       `json:"isActive"`
	LastLoginAt *time.Time `json:"lastLoginAt"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:          u.ID,
		Email:       u.Email,
		Name:        u.Name,
		Role:        u.Role,
		IsActive:    u.IsActive,
		LastLoginAt: u.LastLoginAt,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}

// AdminUserResponse is the response format for admin user management
type AdminUserResponse struct {
	ID          string     `json:"id"`
	Email       string     `json:"email"`
	Name        *string    `json:"name"`
	Role        Role       `json:"role"`
	IsActive    bool       `json:"isActive"`
	LastLoginAt *time.Time `json:"lastLoginAt"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

func (u *User) ToAdminResponse() AdminUserResponse {
	return AdminUserResponse{
		ID:          u.ID,
		Email:       u.Email,
		Name:        u.Name,
		Role:        u.Role,
		IsActive:    u.IsActive,
		LastLoginAt: u.LastLoginAt,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}

// PasswordResetToken stores password reset tokens
type PasswordResetToken struct {
	ID        string `gorm:"primaryKey;type:text"`
	Token     string `gorm:"uniqueIndex;not null"`
	UserID    string `gorm:"not null"`
	User      User   `gorm:"constraint:OnDelete:CASCADE"`
	ExpiresAt time.Time
	UsedAt    *time.Time // Null if not used yet
	CreatedAt time.Time
}

func (p *PasswordResetToken) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return nil
}

// IsValid checks if the token is still valid (not expired and not used)
func (p *PasswordResetToken) IsValid() bool {
	return p.UsedAt == nil && time.Now().Before(p.ExpiresAt)
}
