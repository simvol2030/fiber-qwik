package utils

import (
	"math"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// PaginationQuery represents pagination parameters from query string
type PaginationQuery struct {
	Page  int `query:"page" validate:"omitempty,gte=1"`
	Limit int `query:"limit" validate:"omitempty,gte=1,lte=100"`
}

// PaginationMeta represents pagination metadata in API response
type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"totalPages"`
	HasNext    bool  `json:"hasNext"`
	HasPrev    bool  `json:"hasPrev"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Items      interface{}    `json:"items"`
	Pagination PaginationMeta `json:"pagination"`
}

// DefaultPage is the default page number
const DefaultPage = 1

// DefaultLimit is the default items per page
const DefaultLimit = 20

// MaxLimit is the maximum items per page
const MaxLimit = 100

// NewPaginationQuery creates a PaginationQuery with defaults
func NewPaginationQuery() PaginationQuery {
	return PaginationQuery{
		Page:  DefaultPage,
		Limit: DefaultLimit,
	}
}

// ParsePagination extracts pagination from Fiber context
func ParsePagination(c *fiber.Ctx) PaginationQuery {
	page := c.QueryInt("page", DefaultPage)
	limit := c.QueryInt("limit", DefaultLimit)

	// Validate bounds
	if page < 1 {
		page = DefaultPage
	}
	if limit < 1 {
		limit = DefaultLimit
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}

	return PaginationQuery{
		Page:  page,
		Limit: limit,
	}
}

// Offset calculates the database offset for pagination
func (p PaginationQuery) Offset() int {
	return (p.Page - 1) * p.Limit
}

// Paginate applies pagination to a GORM query and returns metadata
// Usage:
//
//	var users []User
//	pagination := utils.ParsePagination(c)
//	meta, err := pagination.Paginate(db.Model(&User{}), &users)
func (p PaginationQuery) Paginate(query *gorm.DB, dest interface{}) (*PaginationMeta, error) {
	var total int64

	// Count total records (without limit/offset)
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(p.Limit)))

	// Fetch paginated results
	if err := query.Offset(p.Offset()).Limit(p.Limit).Find(dest).Error; err != nil {
		return nil, err
	}

	return &PaginationMeta{
		Page:       p.Page,
		Limit:      p.Limit,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    p.Page < totalPages,
		HasPrev:    p.Page > 1,
	}, nil
}

// SendPaginated sends a paginated success response
func SendPaginated(c *fiber.Ctx, items interface{}, meta *PaginationMeta, statusCode ...int) error {
	code := fiber.StatusOK
	if len(statusCode) > 0 {
		code = statusCode[0]
	}

	response := PaginatedResponse{
		Items:      items,
		Pagination: *meta,
	}

	return SendSuccess(c, response, code)
}
