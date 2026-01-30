package utils

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator is the global validator instance
var Validator *validator.Validate

func init() {
	Validator = validator.New()

	// Use JSON tag names in error messages
	Validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// ValidationError represents a single validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidateStruct validates a struct and returns a slice of validation errors
func ValidateStruct(s interface{}) []ValidationError {
	var errors []ValidationError

	err := Validator.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var message string

			switch err.Tag() {
			case "required":
				message = err.Field() + " is required"
			case "email":
				message = err.Field() + " must be a valid email address"
			case "min":
				message = err.Field() + " must be at least " + err.Param() + " characters"
			case "max":
				message = err.Field() + " must be at most " + err.Param() + " characters"
			case "oneof":
				message = err.Field() + " must be one of: " + err.Param()
			case "url":
				message = err.Field() + " must be a valid URL"
			case "uuid":
				message = err.Field() + " must be a valid UUID"
			case "alphanum":
				message = err.Field() + " must contain only alphanumeric characters"
			case "numeric":
				message = err.Field() + " must be numeric"
			case "gte":
				message = err.Field() + " must be greater than or equal to " + err.Param()
			case "lte":
				message = err.Field() + " must be less than or equal to " + err.Param()
			case "eqfield":
				message = err.Field() + " must be equal to " + err.Param()
			default:
				message = err.Field() + " is invalid"
			}

			errors = append(errors, ValidationError{
				Field:   err.Field(),
				Message: message,
			})
		}
	}

	return errors
}

// HasValidationErrors returns true if there are validation errors
func HasValidationErrors(errors []ValidationError) bool {
	return len(errors) > 0
}
