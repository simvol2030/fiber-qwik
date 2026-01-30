package admin

import (
	"errors"
	"strings"
	"time"

	"backend-go-fiber/internal/models"
	"backend-go-fiber/internal/utils"

	"gorm.io/gorm"
)

// escapeLikeWildcards escapes SQL LIKE wildcards to prevent injection
func escapeLikeWildcards(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "%", "\\%")
	s = strings.ReplaceAll(s, "_", "\\_")
	return s
}

type UsersService struct {
	db *gorm.DB
}

func NewUsersService(db *gorm.DB) *UsersService {
	return &UsersService{db: db}
}

// ListParams contains pagination and filtering parameters
type ListParams struct {
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
	Search   string `json:"search"`
	SortBy   string `json:"sortBy"`
	SortDir  string `json:"sortDir"`
	Role     string `json:"role"`
	IsActive *bool  `json:"isActive"`
}

// ListResult contains paginated list result
type ListResult struct {
	Items      []models.AdminUserResponse `json:"items"`
	Total      int64                      `json:"total"`
	Page       int                        `json:"page"`
	PageSize   int                        `json:"pageSize"`
	TotalPages int                        `json:"totalPages"`
}

// CreateUserInput contains input for creating a user
type CreateUserInput struct {
	Email    string       `json:"email" validate:"required,email"`
	Password string       `json:"password" validate:"required,min=8"`
	Name     *string      `json:"name"`
	Role     models.Role  `json:"role" validate:"omitempty,oneof=user admin"`
	IsActive *bool        `json:"isActive"`
}

// UpdateUserInput contains input for updating a user
type UpdateUserInput struct {
	Email    *string      `json:"email" validate:"omitempty,email"`
	Password *string      `json:"password" validate:"omitempty,min=8"`
	Name     *string      `json:"name"`
	Role     *models.Role `json:"role" validate:"omitempty,oneof=user admin"`
	IsActive *bool        `json:"isActive"`
}

// List returns paginated list of users
func (s *UsersService) List(params ListParams) (*ListResult, error) {
	// Defaults
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 10
	}
	if params.SortBy == "" {
		params.SortBy = "created_at"
	}
	if params.SortDir == "" {
		params.SortDir = "desc"
	}

	// Build query
	query := s.db.Model(&models.User{})

	// Search filter (escape wildcards to prevent SQL injection)
	if params.Search != "" {
		escapedSearch := escapeLikeWildcards(params.Search)
		searchPattern := "%" + escapedSearch + "%"
		query = query.Where("email LIKE ? ESCAPE '\\' OR name LIKE ? ESCAPE '\\'", searchPattern, searchPattern)
	}

	// Role filter
	if params.Role != "" {
		query = query.Where("role = ?", params.Role)
	}

	// Active filter
	if params.IsActive != nil {
		query = query.Where("is_active = ?", *params.IsActive)
	}

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Sorting
	orderClause := params.SortBy
	if params.SortDir == "desc" {
		orderClause += " DESC"
	} else {
		orderClause += " ASC"
	}

	// Fetch users
	var users []models.User
	offset := (params.Page - 1) * params.PageSize
	if err := query.Order(orderClause).Offset(offset).Limit(params.PageSize).Find(&users).Error; err != nil {
		return nil, err
	}

	// Convert to response
	items := make([]models.AdminUserResponse, len(users))
	for i, user := range users {
		items[i] = user.ToAdminResponse()
	}

	// Calculate total pages
	totalPages := int(total) / params.PageSize
	if int(total)%params.PageSize > 0 {
		totalPages++
	}

	return &ListResult{
		Items:      items,
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: totalPages,
	}, nil
}

// GetByID returns a user by ID
func (s *UsersService) GetByID(id string) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// Create creates a new user
func (s *UsersService) Create(input CreateUserInput) (*models.User, error) {
	// Check if email already exists
	var existingUser models.User
	if err := s.db.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	// Set defaults
	role := models.RoleUser
	if input.Role != "" {
		role = input.Role
	}
	isActive := true
	if input.IsActive != nil {
		isActive = *input.IsActive
	}

	user := &models.User{
		Email:        input.Email,
		PasswordHash: hashedPassword,
		Name:         input.Name,
		Role:         role,
		IsActive:     isActive,
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// Update updates an existing user
func (s *UsersService) Update(id string, input UpdateUserInput) (*models.User, error) {
	user, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Check if new email already exists
	if input.Email != nil && *input.Email != user.Email {
		var existingUser models.User
		if err := s.db.Where("email = ? AND id != ?", *input.Email, id).First(&existingUser).Error; err == nil {
			return nil, errors.New("email already exists")
		}
		user.Email = *input.Email
	}

	// Update password if provided
	if input.Password != nil && *input.Password != "" {
		hashedPassword, err := utils.HashPassword(*input.Password)
		if err != nil {
			return nil, err
		}
		user.PasswordHash = hashedPassword
	}

	// Update other fields
	if input.Name != nil {
		user.Name = input.Name
	}
	if input.Role != nil {
		user.Role = *input.Role
	}
	if input.IsActive != nil {
		user.IsActive = *input.IsActive
	}

	if err := s.db.Save(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// Delete soft deletes a user
func (s *UsersService) Delete(id string) error {
	user, err := s.GetByID(id)
	if err != nil {
		return err
	}

	// Soft delete
	if err := s.db.Delete(user).Error; err != nil {
		return err
	}

	return nil
}

// UpdateLastLogin updates the last login timestamp
func (s *UsersService) UpdateLastLogin(id string) error {
	now := time.Now()
	return s.db.Model(&models.User{}).Where("id = ?", id).Update("last_login_at", now).Error
}
