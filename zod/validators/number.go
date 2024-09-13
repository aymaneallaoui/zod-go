package validators

import (
	"fmt"

	"github.com/aymaneallaoui/zod-Go/zod"
)

type NumberSchema struct {
	min         *float64
	max         *float64
	required    bool
	customError map[string]string
}

func Number() *NumberSchema {
	return &NumberSchema{
		customError: make(map[string]string),
	}
}

func (n *NumberSchema) Min(value float64) *NumberSchema {
	n.min = &value
	return n
}

func (n *NumberSchema) Max(value float64) *NumberSchema {
	n.max = &value
	return n
}

func (n *NumberSchema) Required() *NumberSchema {
	n.required = true
	return n
}

func (n *NumberSchema) WithMessage(validationType, message string) *NumberSchema {
	n.customError[validationType] = message
	return n
}

func (n *NumberSchema) getErrorMessage(validationType, defaultMessage string) string {
	if msg, exists := n.customError[validationType]; exists {
		return msg
	}
	return defaultMessage
}

func (n *NumberSchema) Validate(data interface{}) error {
	var num float64
	switch v := data.(type) {
	case int:
		num = float64(v)
	case float64:
		num = v
	default:
		return zod.NewValidationError("number", data, n.getErrorMessage("type", "invalid type, expected number (int or float64)"))
	}

	if n.required && num == 0 {
		return zod.NewValidationError("number", data, n.getErrorMessage("required", "number is required"))
	}

	if n.min != nil && num < *n.min {
		return zod.NewValidationError("number", num, n.getErrorMessage("min", fmt.Sprintf("number is too small, minimum is %f", *n.min)))
	}

	if n.max != nil && num > *n.max {
		return zod.NewValidationError("number", num, n.getErrorMessage("max", fmt.Sprintf("number is too large, maximum is %f", *n.max)))
	}

	return nil
}
