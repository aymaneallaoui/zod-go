package zod

import (
	"encoding/json"
	"fmt"
)

// ValidationError holds detailed information about validation errors
type ValidationError struct {
	Field   string            `json:"field"`
	Value   interface{}       `json:"value"`
	Message string            `json:"message"`
	Details []ValidationError `json:"details,omitempty"` // Nested errors go here
}

// Error formats the error as a human-readable message
func (e *ValidationError) Error() string {
	return fmt.Sprintf("Validation error on field '%s': %s (Value: %v)", e.Field, e.Message, e.Value)
}

// ErrorJSON formats the error as a JSON string
func (e *ValidationError) ErrorJSON() string {
	jsonData, _ := json.MarshalIndent(e, "", "  ") // Pretty print JSON for readability
	return string(jsonData)
}

// NewValidationError creates a new validation error
func NewValidationError(field string, value interface{}, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Value:   value,
		Message: message,
	}
}

// NewNestedValidationError creates a new validation error with nested details
func NewNestedValidationError(field string, value interface{}, message string, details []ValidationError) *ValidationError {
	return &ValidationError{
		Field:   field,
		Value:   value,
		Message: message,
		Details: details, // Add nested errors here
	}
}
