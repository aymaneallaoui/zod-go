package zod

import (
	"encoding/json"
	"fmt"
)

type ValidationError struct {
	Field   string            `json:"field"`
	Value   interface{}       `json:"value"`
	Message string            `json:"message"`
	Details []ValidationError `json:"details,omitempty"`
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("Validation error on field '%s': %s (Value: %v)", e.Field, e.Message, e.Value)
}

func (e *ValidationError) ErrorJSON() string {
	jsonData, _ := json.MarshalIndent(e, "", "  ")
	return string(jsonData)
}

func NewValidationError(field string, value interface{}, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Value:   value,
		Message: message,
	}
}

func NewNestedValidationError(field string, value interface{}, message string, details []ValidationError) *ValidationError {
	return &ValidationError{
		Field:   field,
		Value:   value,
		Message: message,
		Details: details,
	}
}
