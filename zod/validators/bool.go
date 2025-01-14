package validators

import (
	"fmt"

	"github.com/aymaneallaoui/zod-go/zod"
)

// BoolSchema represents a boolean validation schema
type BoolSchema struct {
	required     bool
	mustBeTrue   bool
	mustBeFalse  bool
	customFunc   func(bool) error
	optional     bool
	defaultValue *bool
	customError  map[string]string
}

// Bool creates a new boolean schema
func Bool() *BoolSchema {
	return &BoolSchema{
		customError: make(map[string]string),
	}
}

// Required marks the boolean field as required
func (b *BoolSchema) Required() *BoolSchema {
	b.required = true
	b.optional = false
	return b
}

// Optional marks the boolean field as optional
func (b *BoolSchema) Optional() *BoolSchema {
	b.optional = true
	b.required = false
	return b
}

// Default sets a default value for optional boolean fields
func (b *BoolSchema) Default(value bool) *BoolSchema {
	b.defaultValue = &value
	b.optional = true
	return b
}

// True requires the boolean to be true
func (b *BoolSchema) True() *BoolSchema {
	b.mustBeTrue = true
	b.mustBeFalse = false
	return b
}

// False requires the boolean to be false
func (b *BoolSchema) False() *BoolSchema {
	b.mustBeFalse = true
	b.mustBeTrue = false
	return b
}

// Custom adds a custom validation function
func (b *BoolSchema) Custom(fn func(bool) error) *BoolSchema {
	b.customFunc = fn
	return b
}

// WithMessage sets a custom error message for a validation type
func (b *BoolSchema) WithMessage(validationType, message string) *BoolSchema {
	b.customError[validationType] = message
	return b
}

// getErrorMessage returns custom or default error message
func (b *BoolSchema) getErrorMessage(validationType, defaultMessage string) string {
	if msg, exists := b.customError[validationType]; exists {
		return msg
	}
	return defaultMessage
}

// convertToBool converts various types to boolean
func convertToBool(data interface{}) (bool, bool) {
	switch v := data.(type) {
	case bool:
		return v, true
	case string:
		// Accept common string representations
		switch v {
		case "true", "True", "TRUE", "1", "yes", "Yes", "YES", "on", "On", "ON":
			return true, true
		case "false", "False", "FALSE", "0", "no", "No", "NO", "off", "Off", "OFF":
			return false, true
		default:
			return false, false
		}
	case int, int8, int16, int32, int64:
		// Non-zero integers are true
		return fmt.Sprintf("%v", v) != "0", true
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%v", v) != "0", true
	case float32, float64:
		return fmt.Sprintf("%v", v) != "0", true
	default:
		return false, false
	}
}

// Validate performs validation on the provided data
func (b *BoolSchema) Validate(data interface{}) error {
	// Handle nil values
	if data == nil {
		if b.required {
			return zod.NewValidationError("", nil, b.getErrorMessage("required", "boolean field is required"))
		}
		if b.defaultValue != nil {
			return b.Validate(*b.defaultValue)
		}
		if b.optional {
			return nil
		}
		return zod.NewValidationError("", nil, b.getErrorMessage("required", "boolean field is required"))
	}

	// Convert to boolean
	value, ok := convertToBool(data)
	if !ok {
		return zod.NewValidationError(fmt.Sprintf("%v", data), data,
			b.getErrorMessage("type", "invalid type, expected boolean"))
	}

	// Value-specific validations
	if b.mustBeTrue && !value {
		return zod.NewValidationError(fmt.Sprintf("%v", value), value,
			b.getErrorMessage("true", "value must be true"))
	}

	if b.mustBeFalse && value {
		return zod.NewValidationError(fmt.Sprintf("%v", value), value,
			b.getErrorMessage("false", "value must be false"))
	}

	// Custom validation
	if b.customFunc != nil {
		if err := b.customFunc(value); err != nil {
			return err
		}
	}

	return nil
}

// Helper functions for common boolean schemas

// TrueBoolean creates a boolean schema that must be true
func TrueBoolean() *BoolSchema {
	return Bool().True().Required()
}

// FalseBoolean creates a boolean schema that must be false
func FalseBoolean() *BoolSchema {
	return Bool().False().Required()
}

// OptionalBoolean creates an optional boolean schema
func OptionalBoolean() *BoolSchema {
	return Bool().Optional()
}

// RequiredBoolean creates a required boolean schema (any value)
func RequiredBoolean() *BoolSchema {
	return Bool().Required()
}

// ConsentBoolean creates a boolean schema for consent (must be true)
// Commonly used for terms of service, privacy policy acceptance, etc.
func ConsentBoolean() *BoolSchema {
	return Bool().True().Required().WithMessage("true", "consent is required")
}

// ToggleBoolean creates a boolean schema with a default value
func ToggleBoolean(defaultValue bool) *BoolSchema {
	return Bool().Default(defaultValue)
}
