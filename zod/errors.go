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

// NewValidationError creates a new simple validation error
func NewValidationError(field string, value interface{}, message string) *ValidationError {
	return &ValidationError{
		ErrorType: "Validation Error",
		Message:   message, // Directly store the relevant message without extra formatting
		Field:     field,
		Value:     value,
	}
}

// NewNestedValidationError creates a new nested validation error
func NewNestedValidationError(field string, value interface{}, message string, details []ValidationError) *ValidationError {
	return &ValidationError{
		ErrorType: "Nested Validation Error",
		Message:   message,
		Field:     field,
		Value:     value,
		Details:   details,
	}
}

// Error returns a human-readable string for the top-level error and its nested details
func (v *ValidationError) Error() string {
	if len(v.Details) > 0 {
		return v.formatNestedError()
	}
	return fmt.Sprintf("Field: %s, Error: %s", v.Field, v.Message)
}

// formatNestedError recursively formats the nested validation errors
func (v *ValidationError) formatNestedError() string {
	var result string
	for _, detail := range v.Details {
		result += fmt.Sprintf("Field: %s, Error: %s\n", detail.Field, detail.Message)
	}
	return result
}

// ErrorJSON returns the JSON representation of the error with only the message
func (v *ValidationError) ErrorJSON() string {
	// If there are nested details, we need to flatten them
	if len(v.Details) > 0 {
		return v.flattenNestedErrors()
	}

	// Return simple error if there are no nested details
	errJSON, _ := json.Marshal(map[string]string{
		"field":   v.Field,
		"message": v.Message, // Use only the message itself, no "Field: string, Error: ..." formatting
	})
	return string(errJSON)
}

// flattenNestedErrors returns a simplified JSON with just the most relevant error messages
func (v *ValidationError) flattenNestedErrors() string {
	var flatErrors []map[string]string

	// Collect errors from both top-level and nested errors
	v.collectErrors(&flatErrors)

	flatErrorJSON, _ := json.Marshal(flatErrors)
	return string(flatErrorJSON)
}

// collectErrors recursively collects both top-level and nested validation errors
func (v *ValidationError) collectErrors(flatErrors *[]map[string]string) {
	// Add the current error to the list
	if v.Field != "" && v.Message != "" {
		*flatErrors = append(*flatErrors, map[string]string{
			"field":   v.Field,
			"message": v.Message, // Only use the clean message without extra formatting
		})
	}

	// Recursively collect any nested errors
	for _, detail := range v.Details {
		detail.collectErrors(flatErrors)
	}
}
