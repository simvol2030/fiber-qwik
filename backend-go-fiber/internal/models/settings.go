package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SettingType represents the type of setting value
type SettingType string

const (
	SettingTypeString  SettingType = "string"
	SettingTypeNumber  SettingType = "number"
	SettingTypeBoolean SettingType = "boolean"
	SettingTypeJSON    SettingType = "json"
)

// AppSettings stores application settings as key-value pairs
type AppSettings struct {
	ID        string      `gorm:"primaryKey;type:text" json:"id"`
	Key       string      `gorm:"uniqueIndex;not null" json:"key"`
	Value     string      `gorm:"type:text" json:"value"`
	Type      SettingType `gorm:"type:text;default:string" json:"type"`
	Label        string      `gorm:"type:text" json:"label"`                        // Human-readable label
	SettingGroup string      `gorm:"column:setting_group;type:text" json:"group"`   // Group for UI organization
	UpdatedAt time.Time   `json:"updatedAt"`
	CreatedAt time.Time   `json:"createdAt"`
}

func (s *AppSettings) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	if s.Type == "" {
		s.Type = SettingTypeString
	}
	return nil
}

// AppSettingsResponse is the response format for settings
type AppSettingsResponse struct {
	ID        string      `json:"id"`
	Key       string      `json:"key"`
	Value     string      `json:"value"`
	Type      SettingType `json:"type"`
	Label     string      `json:"label"`
	Group     string      `json:"group"`
	UpdatedAt time.Time   `json:"updatedAt"`
}

func (s *AppSettings) ToResponse() AppSettingsResponse {
	return AppSettingsResponse{
		ID:        s.ID,
		Key:       s.Key,
		Value:     s.Value,
		Type:      s.Type,
		Label:     s.Label,
		Group:     s.SettingGroup,
		UpdatedAt: s.UpdatedAt,
	}
}

// DefaultSettings returns default application settings
func DefaultSettings() []AppSettings {
	return []AppSettings{
		{Key: "app_name", Value: "My App", Type: SettingTypeString, Label: "Application Name", SettingGroup: "general"},
		{Key: "app_description", Value: "A Go Fiber + SvelteKit application", Type: SettingTypeString, Label: "Description", SettingGroup: "general"},
		{Key: "maintenance_mode", Value: "false", Type: SettingTypeBoolean, Label: "Maintenance Mode", SettingGroup: "general"},
		{Key: "allow_registration", Value: "true", Type: SettingTypeBoolean, Label: "Allow Registration", SettingGroup: "auth"},
		{Key: "max_login_attempts", Value: "5", Type: SettingTypeNumber, Label: "Max Login Attempts", SettingGroup: "auth"},
	}
}
