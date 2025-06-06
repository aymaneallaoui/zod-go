package validators

import (
	"fmt"
	"math"

	"github.com/aymaneallaoui/zod-go/zod"
)

// Core number schema struct (unexported)
// This contains all the validation configuration and is wrapped by state-specific types
type numberSchema struct {
	minValue     *float64
	maxValue     *float64
	integerOnly  bool
	positiveOnly bool
	negativeOnly bool
	customFunc   func(float64) error
	required     bool
	optional     bool
	defaultValue *float64
	customError  map[string]string
}

// State wrapper types for compile-time safety
type requiredNumberSchema struct {
	*numberSchema
}

type optionalNumberSchema struct {
	*numberSchema
}

// NumberBuilder implementation (initial state)
// These methods return NumberBuilder to allow continued configuration

func (n *numberSchema) Min(value float64) NumberBuilder {
	n.minValue = &value
	return n
}

func (n *numberSchema) Max(value float64) NumberBuilder {
	n.maxValue = &value
	return n
}

func (n *numberSchema) Integer() NumberBuilder {
	n.integerOnly = true
	return n
}

func (n *numberSchema) Positive() NumberBuilder {
	n.positiveOnly = true
	return n
}

func (n *numberSchema) Negative() NumberBuilder {
	n.negativeOnly = true
	return n
}

func (n *numberSchema) Custom(fn func(float64) error) NumberBuilder {
	n.customFunc = fn
	return n
}

// State transition methods - these change the return type to enforce compile-time safety
func (n *numberSchema) Required() RequiredNumberBuilder {
	n.required = true
	n.optional = false
	return &requiredNumberSchema{n}
}

func (n *numberSchema) Optional() OptionalNumberBuilder {
	n.optional = true
	n.required = false
	return &optionalNumberSchema{n}
}

// Error message methods for NumberBuilder
func (n *numberSchema) WithMessage(validationType, message string) NumberBuilder {
	n.customError[validationType] = message
	return n
}

func (n *numberSchema) WithMinMessage(message string) NumberBuilder {
	return n.WithMessage(errorKeys.Min, message)
}

func (n *numberSchema) WithMaxMessage(message string) NumberBuilder {
	return n.WithMessage(errorKeys.Max, message)
}

func (n *numberSchema) WithIntegerMessage(message string) NumberBuilder {
	return n.WithMessage(errorKeys.Integer, message)
}

func (n *numberSchema) WithPositiveMessage(message string) NumberBuilder {
	return n.WithMessage(errorKeys.Positive, message)
}

func (n *numberSchema) WithNegativeMessage(message string) NumberBuilder {
	return n.WithMessage(errorKeys.Negative, message)
}

// RequiredNumberBuilder implementation
// These methods return RequiredNumberBuilder to maintain the required state

func (r *requiredNumberSchema) Min(value float64) RequiredNumberBuilder {
	r.numberSchema.minValue = &value
	return r
}

func (r *requiredNumberSchema) Max(value float64) RequiredNumberBuilder {
	r.numberSchema.maxValue = &value
	return r
}

func (r *requiredNumberSchema) Integer() RequiredNumberBuilder {
	r.numberSchema.integerOnly = true
	return r
}

func (r *requiredNumberSchema) Positive() RequiredNumberBuilder {
	r.numberSchema.positiveOnly = true
	return r
}

func (r *requiredNumberSchema) Negative() RequiredNumberBuilder {
	r.numberSchema.negativeOnly = true
	return r
}

func (r *requiredNumberSchema) Custom(fn func(float64) error) RequiredNumberBuilder {
	r.numberSchema.customFunc = fn
	return r
}

// Error message methods for RequiredNumberBuilder
func (r *requiredNumberSchema) WithMessage(validationType, message string) RequiredNumberBuilder {
	r.numberSchema.customError[validationType] = message
	return r
}

func (r *requiredNumberSchema) WithMinMessage(message string) RequiredNumberBuilder {
	return r.WithMessage(errorKeys.Min, message)
}

func (r *requiredNumberSchema) WithMaxMessage(message string) RequiredNumberBuilder {
	return r.WithMessage(errorKeys.Max, message)
}

func (r *requiredNumberSchema) WithIntegerMessage(message string) RequiredNumberBuilder {
	return r.WithMessage(errorKeys.Integer, message)
}

func (r *requiredNumberSchema) WithPositiveMessage(message string) RequiredNumberBuilder {
	return r.WithMessage(errorKeys.Positive, message)
}

func (r *requiredNumberSchema) WithNegativeMessage(message string) RequiredNumberBuilder {
	return r.WithMessage(errorKeys.Negative, message)
}

func (r *requiredNumberSchema) WithRequiredMessage(message string) RequiredNumberBuilder {
	return r.WithMessage(errorKeys.Required, message)
}

// OptionalNumberBuilder implementation
// These methods return OptionalNumberBuilder to maintain the optional state

func (o *optionalNumberSchema) Min(value float64) OptionalNumberBuilder {
	o.numberSchema.minValue = &value
	return o
}

func (o *optionalNumberSchema) Max(value float64) OptionalNumberBuilder {
	o.numberSchema.maxValue = &value
	return o
}

func (o *optionalNumberSchema) Integer() OptionalNumberBuilder {
	o.numberSchema.integerOnly = true
	return o
}

func (o *optionalNumberSchema) Positive() OptionalNumberBuilder {
	o.numberSchema.positiveOnly = true
	return o
}

func (o *optionalNumberSchema) Negative() OptionalNumberBuilder {
	o.numberSchema.negativeOnly = true
	return o
}

func (o *optionalNumberSchema) Custom(fn func(float64) error) OptionalNumberBuilder {
	o.numberSchema.customFunc = fn
	return o
}

// Default is only available on optional builders - this is the key DX improvement!
func (o *optionalNumberSchema) Default(value float64) OptionalNumberBuilder {
	o.numberSchema.defaultValue = &value
	return o
}

// Error message methods for OptionalNumberBuilder
func (o *optionalNumberSchema) WithMessage(validationType, message string) OptionalNumberBuilder {
	o.numberSchema.customError[validationType] = message
	return o
}

func (o *optionalNumberSchema) WithMinMessage(message string) OptionalNumberBuilder {
	return o.WithMessage(errorKeys.Min, message)
}

func (o *optionalNumberSchema) WithMaxMessage(message string) OptionalNumberBuilder {
	return o.WithMessage(errorKeys.Max, message)
}

func (o *optionalNumberSchema) WithIntegerMessage(message string) OptionalNumberBuilder {
	return o.WithMessage(errorKeys.Integer, message)
}

func (o *optionalNumberSchema) WithPositiveMessage(message string) OptionalNumberBuilder {
	return o.WithMessage(errorKeys.Positive, message)
}

func (o *optionalNumberSchema) WithNegativeMessage(message string) OptionalNumberBuilder {
	return o.WithMessage(errorKeys.Negative, message)
}

// Validation methods - these are the final methods in the builder chain
func (r *requiredNumberSchema) Validate(data interface{}) error {
	return r.numberSchema.validate(data)
}

func (o *optionalNumberSchema) Validate(data interface{}) error {
	return o.numberSchema.validate(data)
}

// Core validation logic (shared between required and optional)
func (n *numberSchema) validate(data interface{}) error {
	// Handle nil values
	if data == nil {
		if n.required {
			return zod.NewValidationError("", nil, n.getErrorMessage(errorKeys.Required, "field is required"))
		}
		if n.defaultValue != nil {
			return n.validate(*n.defaultValue)
		}
		if n.optional {
			return nil
		}
		return zod.NewValidationError("", nil, n.getErrorMessage(errorKeys.Required, "field is required"))
	}

	// Type check and conversion - support multiple numeric types
	var num float64
	switch v := data.(type) {
	case int:
		num = float64(v)
	case int8:
		num = float64(v)
	case int16:
		num = float64(v)
	case int32:
		num = float64(v)
	case int64:
		num = float64(v)
	case uint:
		num = float64(v)
	case uint8:
		num = float64(v)
	case uint16:
		num = float64(v)
	case uint32:
		num = float64(v)
	case uint64:
		num = float64(v)
	case float32:
		num = float64(v)
	case float64:
		num = v
	default:
		return zod.NewValidationError(fmt.Sprintf("%v", data), data,
			n.getErrorMessage(errorKeys.Type, "invalid type, expected number"))
	}

	// Integer validation
	if n.integerOnly && num != math.Trunc(num) {
		return zod.NewValidationError(num, num,
			n.getErrorMessage(errorKeys.Integer, "value must be an integer"))
	}

	// Range validations
	if n.minValue != nil && num < *n.minValue {
		return zod.NewValidationError(num, num,
			n.getErrorMessage(errorKeys.Min,
				fmt.Sprintf("value is too small, minimum is %g", *n.minValue)))
	}

	if n.maxValue != nil && num > *n.maxValue {
		return zod.NewValidationError(num, num,
			n.getErrorMessage(errorKeys.Max,
				fmt.Sprintf("value is too large, maximum is %g", *n.maxValue)))
	}

	// Sign validations
	if n.positiveOnly && num <= 0 {
		return zod.NewValidationError(num, num,
			n.getErrorMessage(errorKeys.Positive, "value must be positive"))
	}

	if n.negativeOnly && num >= 0 {
		return zod.NewValidationError(num, num,
			n.getErrorMessage(errorKeys.Negative, "value must be negative"))
	}

	// Custom validation
	if n.customFunc != nil {
		if err := n.customFunc(num); err != nil {
			return err
		}
	}

	return nil
}

// Helper methods (unexported)
func (n *numberSchema) getErrorMessage(validationType, defaultMessage string) string {
	if msg, exists := n.customError[validationType]; exists {
		return msg
	}
	return defaultMessage
}
