package zod

import (
	"encoding/json"
	"fmt"
)

type ValidationError struct {
	ErrorType string            `json:"errorType"`
	Message   string            `json:"message"`
	Field     string            `json:"field"`
	Value     interface{}       `json:"value"`
	Details   []ValidationError `json:"details,omitempty"`
}

func NewValidationError(field string, value interface{}, message string) *ValidationError {
	return &ValidationError{
		ErrorType: "Validation Error",
		Message:   message,
		Field:     field,
		Value:     value,
	}
}

func NewNestedValidationError(field string, value interface{}, message string, details []ValidationError) *ValidationError {
	return &ValidationError{
		ErrorType: "Nested Validation Error",
		Message:   message,
		Field:     field,
		Value:     value,
		Details:   details,
	}
}

func (v *ValidationError) Error() string {
	if len(v.Details) > 0 {
		return v.formatNestedError()
	}
	return fmt.Sprintf("Field: %s, Error: %s", v.Field, v.Message)
}

func (v *ValidationError) formatNestedError() string {
	var result string
	for _, detail := range v.Details {
		result += fmt.Sprintf("Field: %s, Error: %s\n", detail.Field, detail.Message)
	}
	return result
}

func (v *ValidationError) ErrorJSON() string {

	if len(v.Details) > 0 {
		return v.flattenNestedErrors()
	}

	errJSON, _ := json.Marshal(map[string]string{
		"field":   v.Field,
		"message": v.Message,
	})
	return string(errJSON)
}

func (v *ValidationError) flattenNestedErrors() string {
	var flatErrors []map[string]string

	v.collectErrors(&flatErrors)

	flatErrorJSON, _ := json.Marshal(flatErrors)
	return string(flatErrorJSON)
}

func (v *ValidationError) collectErrors(flatErrors *[]map[string]string) {

	if v.Field != "" && v.Message != "" {
		*flatErrors = append(*flatErrors, map[string]string{
			"field":   v.Field,
			"message": v.Message,
		})
	}

	for _, detail := range v.Details {
		detail.collectErrors(flatErrors)
	}
}
