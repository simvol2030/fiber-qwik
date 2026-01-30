package admin

import (
	"time"

	"backend-go-fiber/internal/models"

	"gorm.io/gorm"
)

type DashboardService struct {
	db *gorm.DB
}

func NewDashboardService(db *gorm.DB) *DashboardService {
	return &DashboardService{db: db}
}

// DashboardStats contains dashboard statistics
type DashboardStats struct {
	TotalUsers        int64              `json:"totalUsers"`
	ActiveUsers       int64              `json:"activeUsers"`
	AdminUsers        int64              `json:"adminUsers"`
	NewUsersToday     int64              `json:"newUsersToday"`
	NewUsersThisWeek  int64              `json:"newUsersThisWeek"`
	NewUsersThisMonth int64              `json:"newUsersThisMonth"`
	RecentUsers       []RecentUser       `json:"recentUsers"`
	RecentActivity    []ActivityLogEntry `json:"recentActivity"`
}

// RecentUser represents a recently registered user
type RecentUser struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      *string   `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

// ActivityLogEntry represents an activity log entry
type ActivityLogEntry struct {
	Type      string    `json:"type"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// GetStats returns dashboard statistics
func (s *DashboardService) GetStats() (*DashboardStats, error) {
	stats := &DashboardStats{}

	// Total users
	if err := s.db.Model(&models.User{}).Count(&stats.TotalUsers).Error; err != nil {
		return nil, err
	}

	// Active users
	if err := s.db.Model(&models.User{}).Where("is_active = ?", true).Count(&stats.ActiveUsers).Error; err != nil {
		return nil, err
	}

	// Admin users
	if err := s.db.Model(&models.User{}).Where("role = ?", models.RoleAdmin).Count(&stats.AdminUsers).Error; err != nil {
		return nil, err
	}

	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	weekStart := todayStart.AddDate(0, 0, -7)
	monthStart := todayStart.AddDate(0, -1, 0)

	// New users today
	if err := s.db.Model(&models.User{}).Where("created_at >= ?", todayStart).Count(&stats.NewUsersToday).Error; err != nil {
		return nil, err
	}

	// New users this week
	if err := s.db.Model(&models.User{}).Where("created_at >= ?", weekStart).Count(&stats.NewUsersThisWeek).Error; err != nil {
		return nil, err
	}

	// New users this month
	if err := s.db.Model(&models.User{}).Where("created_at >= ?", monthStart).Count(&stats.NewUsersThisMonth).Error; err != nil {
		return nil, err
	}

	// Recent users (last 5)
	var users []models.User
	if err := s.db.Model(&models.User{}).Order("created_at DESC").Limit(5).Find(&users).Error; err != nil {
		return nil, err
	}

	stats.RecentUsers = make([]RecentUser, len(users))
	for i, user := range users {
		stats.RecentUsers[i] = RecentUser{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			CreatedAt: user.CreatedAt,
		}
	}

	// Recent activity (mock for now - can be extended with real activity logging)
	stats.RecentActivity = []ActivityLogEntry{
		{Type: "info", Message: "Dashboard accessed", Timestamp: now},
	}

	return stats, nil
}
