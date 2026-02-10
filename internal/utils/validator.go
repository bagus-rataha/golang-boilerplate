package utils

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

// ValidateStruct validates struct using validator tags
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

// ParseAndValidate combines body parsing and validation
// Returns error after writing detailed response
func ParseAndValidate(c *fiber.Ctx, input interface{}) error {
	// Parse JSON body
	if err := c.BodyParser(input); err != nil {
		// JSON syntax error
		if jsonErr, ok := err.(*json.SyntaxError); ok {
			c.Status(fiber.StatusBadRequest).JSON(Response{
				Success: false,
				Message: fmt.Sprintf("JSON syntax error at position %d", jsonErr.Offset),
				Data:    nil,
			})
			return fiber.ErrBadRequest
		}

		// Type mismatch error
		if jsonErr, ok := err.(*json.UnmarshalTypeError); ok {
			c.Status(fiber.StatusBadRequest).JSON(Response{
				Success: false,
				Message: fmt.Sprintf("Invalid type for field '%s': expected %s, got %s",
					jsonErr.Field, jsonErr.Type, jsonErr.Value),
				Data: nil,
			})
			return fiber.ErrBadRequest
		}

		// Generic parse error
		c.Status(fiber.StatusBadRequest).JSON(Response{
			Success: false,
			Message: "Invalid request body",
			Data:    nil,
		})
		return fiber.ErrBadRequest
	}

	// Validate struct
	if err := validate.Struct(input); err != nil {
		errors := make(map[string]string)

		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, fieldError := range validationErrors {
				field := strings.ToLower(fieldError.Field())

				switch fieldError.Tag() {
				case "required":
					errors[field] = fmt.Sprintf("%s is required", field)
				case "email":
					errors[field] = "Invalid email format"
				case "min":
					errors[field] = fmt.Sprintf("%s must be at least %s characters", field, fieldError.Param())
				case "max":
					errors[field] = fmt.Sprintf("%s must not exceed %s characters", field, fieldError.Param())
				case "oneof":
					errors[field] = fmt.Sprintf("%s must be one of: %s", field, fieldError.Param())
				case "gtfield":
					errors[field] = fmt.Sprintf("%s must be greater than %s", field, fieldError.Param())
				case "gt":
					errors[field] = fmt.Sprintf("%s must be greater than %s", field, fieldError.Param())
				case "gte":
					errors[field] = fmt.Sprintf("%s must be greater than or equal to %s", field, fieldError.Param())
				default:
					errors[field] = fmt.Sprintf("Validation failed on '%s'", fieldError.Tag())
				}
			}
		}

		c.Status(fiber.StatusBadRequest).JSON(Response{
			Success: false,
			Message: "Validation failed",
			Data:    errors,
		})
		return fiber.ErrBadRequest
	}

	return nil
}
