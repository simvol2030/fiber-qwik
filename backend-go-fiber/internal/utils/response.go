package utils

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
	Meta    *APIMeta    `json:"meta,omitempty"`
}

type APIError struct {
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Details []FieldError `json:"details,omitempty"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type APIMeta struct {
	Timestamp string `json:"timestamp"`
	RequestID string `json:"requestId,omitempty"`
}

func SendSuccess(c *fiber.Ctx, data interface{}, statusCode ...int) error {
	code := fiber.StatusOK
	if len(statusCode) > 0 {
		code = statusCode[0]
	}

	requestID := c.Locals("requestId")
	var reqIDStr string
	if requestID != nil {
		reqIDStr = requestID.(string)
	}

	return c.Status(code).JSON(APIResponse{
		Success: true,
		Data:    data,
		Meta: &APIMeta{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			RequestID: reqIDStr,
		},
	})
}

func SendError(c *fiber.Ctx, code string, message string, statusCode int, details ...[]FieldError) error {
	requestID := c.Locals("requestId")
	var reqIDStr string
	if requestID != nil {
		reqIDStr = requestID.(string)
	}

	apiError := &APIError{
		Code:    code,
		Message: message,
	}

	if len(details) > 0 {
		apiError.Details = details[0]
	}

	return c.Status(statusCode).JSON(APIResponse{
		Success: false,
		Error:   apiError,
		Meta: &APIMeta{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			RequestID: reqIDStr,
		},
	})
}

// SendValidationError sends a validation error response with field details
func SendValidationError(c *fiber.Ctx, validationErrors []ValidationError) error {
	fieldErrors := make([]FieldError, len(validationErrors))
	for i, ve := range validationErrors {
		fieldErrors[i] = FieldError{
			Field:   ve.Field,
			Message: ve.Message,
		}
	}
	return SendError(c, "VALIDATION_ERROR", "Validation failed", fiber.StatusBadRequest, fieldErrors)
}
