package validators

import "github.com/aymaneallaoui/zod-Go/zod"

type NumberSchema struct {
	min      *float64
	max      *float64
	required bool
}

func Number() *NumberSchema {
	return &NumberSchema{}
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

func (n *NumberSchema) Validate(data interface{}) error {
	num, ok := data.(float64)
	if !ok {
		return zod.NewValidationError("number", data, "invalid type, expected number")
	}

	if n.required && num == 0 {
		return zod.NewValidationError("number", data, "number is required")
	}

	if n.min != nil && num < *n.min {
		return zod.NewValidationError("number", num, "number is too small")
	}

	if n.max != nil && num > *n.max {
		return zod.NewValidationError("number", num, "number is too large")
	}

	return nil
}
