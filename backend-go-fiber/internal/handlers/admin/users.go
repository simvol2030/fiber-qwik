package admin

import (
	"strconv"

	"backend-go-fiber/internal/services/admin"
	"backend-go-fiber/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type UsersHandler struct {
	service *admin.UsersService
}

func NewUsersHandler(service *admin.UsersService) *UsersHandler {
	return &UsersHandler{service: service}
}

// List returns paginated list of users
// GET /api/admin/users
func (h *UsersHandler) List(c *fiber.Ctx) error {
	// Parse query parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "10"))
	search := c.Query("search", "")
	sortBy := c.Query("sortBy", "created_at")
	sortDir := c.Query("sortDir", "desc")
	role := c.Query("role", "")

	var isActive *bool
	if isActiveStr := c.Query("isActive", ""); isActiveStr != "" {
		val := isActiveStr == "true"
		isActive = &val
	}

	params := admin.ListParams{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
		SortBy:   sortBy,
		SortDir:  sortDir,
		Role:     role,
		IsActive: isActive,
	}

	result, err := h.service.List(params)
	if err != nil {
		return utils.SendError(c, "INTERNAL_ERROR", "Failed to list users", fiber.StatusInternalServerError)
	}

	return utils.SendSuccess(c, result, fiber.StatusOK)
}

// Get returns a single user by ID
// GET /api/admin/users/:id
func (h *UsersHandler) Get(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return utils.SendError(c, "VALIDATION_ERROR", "User ID is required", fiber.StatusBadRequest)
	}

	user, err := h.service.GetByID(id)
	if err != nil {
		if err.Error() == "user not found" {
			return utils.SendError(c, "NOT_FOUND", "User not found", fiber.StatusNotFound)
		}
		return utils.SendError(c, "INTERNAL_ERROR", "Failed to get user", fiber.StatusInternalServerError)
	}

	return utils.SendSuccess(c, user.ToAdminResponse(), fiber.StatusOK)
}

// Create creates a new user
// POST /api/admin/users
func (h *UsersHandler) Create(c *fiber.Ctx) error {
	var input admin.CreateUserInput
	if err := c.BodyParser(&input); err != nil {
		return utils.SendError(c, "VALIDATION_ERROR", "Invalid request body", fiber.StatusBadRequest)
	}

	// Validate input
	if errors := utils.ValidateStruct(input); errors != nil {
		return utils.SendValidationError(c, errors)
	}

	user, err := h.service.Create(input)
	if err != nil {
		if err.Error() == "email already exists" {
			return utils.SendError(c, "CONFLICT", "Email already exists", fiber.StatusConflict)
		}
		return utils.SendError(c, "INTERNAL_ERROR", "Failed to create user", fiber.StatusInternalServerError)
	}

	return utils.SendSuccess(c, user.ToAdminResponse(), fiber.StatusCreated)
}

// Update updates an existing user
// PUT /api/admin/users/:id
func (h *UsersHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return utils.SendError(c, "VALIDATION_ERROR", "User ID is required", fiber.StatusBadRequest)
	}

	var input admin.UpdateUserInput
	if err := c.BodyParser(&input); err != nil {
		return utils.SendError(c, "VALIDATION_ERROR", "Invalid request body", fiber.StatusBadRequest)
	}

	// Validate input
	if errors := utils.ValidateStruct(input); errors != nil {
		return utils.SendValidationError(c, errors)
	}

	user, err := h.service.Update(id, input)
	if err != nil {
		if err.Error() == "user not found" {
			return utils.SendError(c, "NOT_FOUND", "User not found", fiber.StatusNotFound)
		}
		if err.Error() == "email already exists" {
			return utils.SendError(c, "CONFLICT", "Email already exists", fiber.StatusConflict)
		}
		return utils.SendError(c, "INTERNAL_ERROR", "Failed to update user", fiber.StatusInternalServerError)
	}

	return utils.SendSuccess(c, user.ToAdminResponse(), fiber.StatusOK)
}

// Delete soft deletes a user
// DELETE /api/admin/users/:id
func (h *UsersHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return utils.SendError(c, "VALIDATION_ERROR", "User ID is required", fiber.StatusBadRequest)
	}

	if err := h.service.Delete(id); err != nil {
		if err.Error() == "user not found" {
			return utils.SendError(c, "NOT_FOUND", "User not found", fiber.StatusNotFound)
		}
		return utils.SendError(c, "INTERNAL_ERROR", "Failed to delete user", fiber.StatusInternalServerError)
	}

	return utils.SendSuccess(c, fiber.Map{"message": "User deleted successfully"}, fiber.StatusOK)
}
