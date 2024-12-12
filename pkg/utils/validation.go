package utils

import (
	"fmt"
	"mpc/pkg/errors"
	"mpc/pkg/logger"
	"regexp"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	// Register custom password validation
	validate.RegisterValidation("password", validatePassword)
}

// Password validation rules
const (
	minPasswordLength = 8
	maxPasswordLength = 72 // bcrypt max length
)

var (
	// Password must contain at least one special character
	specialChars     = `[!@#$%^&*(),.?":{}|<>]`
	specialCharRegex = regexp.MustCompile(specialChars)
)

// validatePassword checks if password meets requirements:
// - At least 8 characters
// - At least 1 uppercase letter
// - At least 1 lowercase letter
// - At least 1 number
// - At least 1 special character
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < minPasswordLength || len(password) > maxPasswordLength {
		return false
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		}
	}

	hasSpecial = specialCharRegex.MatchString(password)

	return hasUpper && hasLower && hasNumber && hasSpecial
}

// ValidateBody validates request body against a model
func ValidateBody(c *gin.Context, model interface{}) error {
	// Bind JSON to model
	if err := c.ShouldBindJSON(model); err != nil {
		logger.Error("Validation: Failed to bind JSON", err)
		return errors.ErrInvalidRequest
	}

	// Validate model
	if err := validate.Struct(model); err != nil {
		validationErrors := make(map[string]string)

		// Extract validation errors
		if errs, ok := err.(validator.ValidationErrors); ok {
			for _, e := range errs {
				validationErrors[e.Field()] = getErrorMessage(e)
			}
		}

		logger.Error("Validation: Failed to validate struct", err,
			logger.Any("validation_errors", validationErrors))
		return errors.ErrInvalidRequest
	}

	return nil
}

// getErrorMessage returns human readable error message for validation tags
func getErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return fmt.Sprintf("Must be at least %s characters", e.Param())
	case "max":
		return fmt.Sprintf("Must be at most %s characters", e.Param())
	case "uuid":
		return "Invalid UUID format"
	case "url":
		return "Invalid URL format"
	case "password":
		return fmt.Sprintf("Password must be %d-%d characters long and contain at least: 1 uppercase letter, 1 lowercase letter, 1 number, and 1 special character (%s)",
			minPasswordLength,
			maxPasswordLength,
			specialChars,
		)
	default:
		return fmt.Sprintf("Failed validation on %s", e.Tag())
	}
}
