package validators

import "github.com/aymaneallaoui/zod-Go/zod"

type BoolSchema struct {
	required bool
}

func Bool() *BoolSchema {
	return &BoolSchema{}
}

func (b *BoolSchema) Required() *BoolSchema {
	b.required = true
	return b
}

func (b *BoolSchema) Validate(data interface{}) error {
	val, ok := data.(bool)
	if !ok {
		return zod.NewValidationError("bool", data, "invalid type, expected boolean")
	}

	if b.required && !val {
		return zod.NewValidationError("bool", data, "boolean is required")
	}

	return nil
}
