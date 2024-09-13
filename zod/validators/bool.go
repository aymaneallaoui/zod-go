package validators

import "github.com/aymaneallaoui/zod-Go/zod"

type BoolSchema struct {
	required    bool
	customError map[string]string
}

func Bool() *BoolSchema {
	return &BoolSchema{
		customError: make(map[string]string),
	}
}

func (b *BoolSchema) Required() *BoolSchema {
	b.required = true
	return b
}

// Add WithMessage method for BoolSchema
func (b *BoolSchema) WithMessage(validationType, message string) *BoolSchema {
	b.customError[validationType] = message
	return b
}

func (b *BoolSchema) getErrorMessage(validationType, defaultMessage string) string {
	if msg, exists := b.customError[validationType]; exists {
		return msg
	}
	return defaultMessage
}

func (b *BoolSchema) Validate(data interface{}) error {
	val, ok := data.(bool)
	if !ok {
		return zod.NewValidationError("bool", data, b.getErrorMessage("type", "invalid type, expected boolean"))
	}

	if b.required && !val {
		return zod.NewValidationError("bool", data, b.getErrorMessage("required", "boolean is required"))
	}

	return nil
}
