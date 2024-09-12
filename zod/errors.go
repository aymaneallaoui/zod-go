package zod

import (
	"encoding/json"
	"fmt"
)

// ValidationError stores details about validation failures.
type ValidationError struct {
	Field   string      `json:"field"`
	Value   interface{} `json:"value"`
	Message string      `json:"message"`
}

// Error formats the error in plain text.
func (e *ValidationError) Error() string {
	return fmt.Sprintf("Validation error on field '%s': %s (Value: %v)", e.Field, e.Message, e.Value)
}

// ErrorJSON formats the error as a JSON string.
func (e *ValidationError) ErrorJSON() string {
	jsonData, _ := json.Marshal(e)
	return string(jsonData)
}

// NewValidationError creates a new validation error for a specific field.
func NewValidationError(field string, value interface{}, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Value:   value,
		Message: message,
	}
}
