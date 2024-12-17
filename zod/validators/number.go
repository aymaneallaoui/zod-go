package validators

import (
	"fmt"
	"math"

	"github.com/aymaneallaoui/zod-go/zod"
)

// NumberSchema represents a number validation schema
type NumberSchema struct {
	minValue       *float64
	maxValue       *float64
	required       bool
	integerOnly    bool
	positiveOnly   bool
	negativeOnly   bool
	nonZero        bool
	multipleOf     *float64
	customFunc     func(float64) error
	optional       bool
	defaultValue   *float64
	customError    map[string]string
}

// Number creates a new number schema
func Number() *NumberSchema {
	return &NumberSchema{
		customError: make(map[string]string),
	}
}

// Min sets the minimum value requirement
func (n *NumberSchema) Min(value float64) *NumberSchema {
	n.minValue = &value
	return n
}

// Max sets the maximum value requirement
func (n *NumberSchema) Max(value float64) *NumberSchema {
	n.maxValue = &value
	return n
}

// Required marks the field as required
func (n *NumberSchema) Required() *NumberSchema {
	n.required = true
	n.optional = false
	return n
}

// Optional marks the field as optional
func (n *NumberSchema) Optional() *NumberSchema {
	n.optional = true
	n.required = false
	return n
}

// Default sets a default value for optional fields
func (n *NumberSchema) Default(value float64) *NumberSchema {
	n.defaultValue = &value
	n.optional = true
	return n
}

// Integer requires the number to be an integer
func (n *NumberSchema) Integer() *NumberSchema {
	n.integerOnly = true
	return n
}

// Positive requires the number to be positive (> 0)
func (n *NumberSchema) Positive() *NumberSchema {
	n.positiveOnly = true
	return n
}

// Negative requires the number to be negative (< 0)
func (n *NumberSchema) Negative() *NumberSchema {
	n.negativeOnly = true
	return n
}

// NonZero requires the number to be non-zero
func (n *NumberSchema) NonZero() *NumberSchema {
	n.nonZero = true
	return n
}

// MultipleOf requires the number to be a multiple of the given value
func (n *NumberSchema) MultipleOf(value float64) *NumberSchema {
	n.multipleOf = &value
	return n
}

// Custom adds a custom validation function
func (n *NumberSchema) Custom(fn func(float64) error) *NumberSchema {
	n.customFunc = fn
	return n
}

// WithMessage sets a custom error message for a validation type
func (n *NumberSchema) WithMessage(validationType, message string) *NumberSchema {
	n.customError[validationType] = message
	return n
}

// getErrorMessage returns custom or default error message
func (n *NumberSchema) getErrorMessage(validationType, defaultMessage string) string {
	if msg, exists := n.customError[validationType]; exists {
		return msg
	}
	return defaultMessage
}

// convertToFloat64 converts various numeric types to float64
func convertToFloat64(data interface{}) (float64, bool) {
	switch v := data.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int8:
		return float64(v), true
	case int16:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case uint:
		return float64(v), true
	case uint8:
		return float64(v), true
	case uint16:
		return float64(v), true
	case uint32:
		return float64(v), true
	case uint64:
		return float64(v), true
	default:
		return 0, false
	}
}

// isInteger checks if a float64 value is an integer
func isInteger(value float64) bool {
	return value == math.Trunc(value) && !math.IsInf(value, 0) && !math.IsNaN(value)
}

// isMultipleOf checks if a number is a multiple of another number
func isMultipleOf(value, multiple float64) bool {
	if multiple == 0 {
		return false
	}
	remainder := math.Mod(value, multiple)
	return math.Abs(remainder) < 1e-10 || math.Abs(remainder-multiple) < 1e-10
}

// Validate performs validation on the provided data
func (n *NumberSchema) Validate(data interface{}) error {
	// Handle nil values
	if data == nil {
		if n.required {
			return zod.NewValidationError("", nil, n.getErrorMessage("required", "field is required"))
		}
		if n.defaultValue != nil {
			return n.Validate(*n.defaultValue)
		}
		if n.optional {
			return nil
		}
		return zod.NewValidationError("", nil, n.getErrorMessage("required", "field is required"))
	}

	// Type conversion
	value, ok := convertToFloat64(data)
	if !ok {
		return zod.NewValidationError(fmt.Sprintf("%v", data), data,
			n.getErrorMessage("type", "invalid type, expected number"))
	}

	// Check for special float values
	if math.IsNaN(value) {
		return zod.NewValidationError(fmt.Sprintf("%v", value), value,
			n.getErrorMessage("nan", "value cannot be NaN"))
	}

	if math.IsInf(value, 0) {
		return zod.NewValidationError(fmt.Sprintf("%v", value), value,
			n.getErrorMessage("infinity", "value cannot be infinite"))
	}

	// Integer validation
	if n.integerOnly && !isInteger(value) {
		return zod.NewValidationError(fmt.Sprintf("%v", value), value,
			n.getErrorMessage("integer", "value must be an integer"))
	}

	// Range validations
	if n.minValue != nil && value < *n.minValue {
		return zod.NewValidationError(fmt.Sprintf("%v", value), value,
			n.getErrorMessage("min", fmt.Sprintf("value must be at least %g", *n.minValue)))
	}

	if n.maxValue != nil && value > *n.maxValue {
		return zod.NewValidationError(fmt.Sprintf("%v", value), value,
			n.getErrorMessage("max", fmt.Sprintf("value must be at most %g", *n.maxValue)))
	}

	// Sign validations
	if n.positiveOnly && value <= 0 {
		return zod.NewValidationError(fmt.Sprintf("%v", value), value,
			n.getErrorMessage("positive", "value must be positive"))
	}

	if n.negativeOnly && value >= 0 {
		return zod.NewValidationError(fmt.Sprintf("%v", value), value,
			n.getErrorMessage("negative", "value must be negative"))
	}

	if n.nonZero && value == 0 {
		return zod.NewValidationError(fmt.Sprintf("%v", value), value,
			n.getErrorMessage("nonzero", "value must be non-zero"))
	}

	// Multiple validation
	if n.multipleOf != nil && !isMultipleOf(value, *n.multipleOf) {
		return zod.NewValidationError(fmt.Sprintf("%v", value), value,
			n.getErrorMessage("multipleOf", fmt.Sprintf("value must be a multiple of %g", *n.multipleOf)))
	}

	// Custom validation
	if n.customFunc != nil {
		if err := n.customFunc(value); err != nil {
			return err
		}
	}

	return nil
}

// Helper functions for common number schemas

// IntegerNumber creates a schema for integer validation
func IntegerNumber() *NumberSchema {
	return Number().Integer().Required()
}

// PositiveNumber creates a schema for positive number validation
func PositiveNumber() *NumberSchema {
	return Number().Positive().Required()
}

// NegativeNumber creates a schema for negative number validation
func NegativeNumber() *NumberSchema {
	return Number().Negative().Required()
}

// NonZeroNumber creates a schema for non-zero number validation
func NonZeroNumber() *NumberSchema {
	return Number().NonZero().Required()
}

// NumberBetween creates a schema with min and max constraints
func NumberBetween(min, max float64) *NumberSchema {
	return Number().Min(min).Max(max).Required()
}

// OptionalNumber creates an optional number schema
func OptionalNumber() *NumberSchema {
	return Number().Optional()
}

// PercentageNumber creates a schema for percentage values (0-100)
func PercentageNumber() *NumberSchema {
	return Number().Min(0).Max(100).Required()
}

// ProbabilityNumber creates a schema for probability values (0-1)
func ProbabilityNumber() *NumberSchema {
	return Number().Min(0).Max(1).Required()
}
