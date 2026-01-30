package admin

import (
	"errors"

	"backend-go-fiber/internal/models"

	"gorm.io/gorm"
)

type SettingsService struct {
	db *gorm.DB
}

func NewSettingsService(db *gorm.DB) *SettingsService {
	return &SettingsService{db: db}
}

// GetAll returns all settings
func (s *SettingsService) GetAll() ([]models.AppSettingsResponse, error) {
	var settings []models.AppSettings
	if err := s.db.Order("setting_group ASC, key ASC").Find(&settings).Error; err != nil {
		return nil, err
	}

	result := make([]models.AppSettingsResponse, len(settings))
	for i, setting := range settings {
		result[i] = setting.ToResponse()
	}

	return result, nil
}

// GetByKey returns a setting by key
func (s *SettingsService) GetByKey(key string) (*models.AppSettings, error) {
	var setting models.AppSettings
	if err := s.db.First(&setting, "key = ?", key).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("setting not found")
		}
		return nil, err
	}
	return &setting, nil
}

// GetByGroup returns settings by group
func (s *SettingsService) GetByGroup(group string) ([]models.AppSettingsResponse, error) {
	var settings []models.AppSettings
	if err := s.db.Where("setting_group = ?", group).Order("key ASC").Find(&settings).Error; err != nil {
		return nil, err
	}

	result := make([]models.AppSettingsResponse, len(settings))
	for i, setting := range settings {
		result[i] = setting.ToResponse()
	}

	return result, nil
}

// UpdateSettingInput contains input for updating a setting
type UpdateSettingInput struct {
	Key   string `json:"key" validate:"required"`
	Value string `json:"value"`
}

// UpdateBatchInput contains input for batch updating settings
type UpdateBatchInput struct {
	Settings []UpdateSettingInput `json:"settings" validate:"required,dive"`
}

// Update updates a setting by key
func (s *SettingsService) Update(key string, value string) (*models.AppSettings, error) {
	setting, err := s.GetByKey(key)
	if err != nil {
		return nil, err
	}

	setting.Value = value
	if err := s.db.Save(setting).Error; err != nil {
		return nil, err
	}

	return setting, nil
}

// UpdateBatch updates multiple settings at once
func (s *SettingsService) UpdateBatch(inputs []UpdateSettingInput) ([]models.AppSettingsResponse, error) {
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	for _, input := range inputs {
		if err := tx.Model(&models.AppSettings{}).Where("key = ?", input.Key).Update("value", input.Value).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return s.GetAll()
}

// CreateSetting creates a new setting
func (s *SettingsService) CreateSetting(setting *models.AppSettings) error {
	return s.db.Create(setting).Error
}

// DeleteSetting deletes a setting by key
func (s *SettingsService) DeleteSetting(key string) error {
	return s.db.Where("key = ?", key).Delete(&models.AppSettings{}).Error
}
