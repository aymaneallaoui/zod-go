package zod

import (
	"encoding/json"
	"fmt"
)

type ValidationError struct {
	Field   string      `json:"field"`
	Value   interface{} `json:"value"`
	Message string      `json:"message"`
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("Validation error on field '%s': %s (Value: %v)", e.Field, e.Message, e.Value)
}

func (e *ValidationError) ErrorJSON() string {
	jsonData, _ := json.Marshal(e)
	return string(jsonData)
}

func NewValidationError(field string, value interface{}, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Value:   value,
		Message: message,
	}
}
