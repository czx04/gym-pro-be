package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator wraps go-playground validator
type Validator struct {
	validate *validator.Validate
}

// New creates a new validator instance
func New() *Validator {
	v := validator.New()

	// Register custom validators here if needed
	// v.RegisterValidation("custom", customValidator)

	return &Validator{validate: v}
}

// Validate validates a struct
func (v *Validator) Validate(i interface{}) error {
	if err := v.validate.Struct(i); err != nil {
		return v.formatValidationErrors(err)
	}
	return nil
}

// formatValidationErrors formats validation errors
func (v *Validator) formatValidationErrors(err error) error {
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		var messages []string
		for _, e := range validationErrs {
			messages = append(messages, v.formatFieldError(e))
		}
		return fmt.Errorf("%s", strings.Join(messages, "; "))
	}
	return err
}

// formatFieldError formats a single field error
func (v *Validator) formatFieldError(e validator.FieldError) string {
	field := strings.ToLower(e.Field())

	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", field, e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters long", field, e.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, e.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, e.Param())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, e.Param())
	case "lt":
		return fmt.Sprintf("%s must be less than %s", field, e.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of [%s]", field, e.Param())
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}

// GetValidationErrors returns validation errors as a map
func (v *Validator) GetValidationErrors(err error) map[string]interface{} {
	errors := make(map[string]interface{})

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrs {
			field := strings.ToLower(e.Field())
			errors[field] = v.formatFieldError(e)
		}
	}

	return errors
}
