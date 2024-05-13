package validators

import "github.com/aymaneallaoui/zod-Go/zod"

type StringSchema struct {
    minLength int
    maxLength int
    required  bool
}

func String() *StringSchema {
    return &StringSchema{}
}

func (s *StringSchema) Min(length int) *StringSchema {
    s.minLength = length
    return s
}

func (s *StringSchema) Max(length int) *StringSchema {
    s.maxLength = length
    return s
}

func (s *StringSchema) Required() *StringSchema {
    s.required = true
    return s
}

func (s *StringSchema) Validate(data interface{}) error {
    str, ok := data.(string)
    if !ok {
        return zod.NewValidationError("string", data, "invalid type, expected string")
    }

    if s.required && str == "" {
        return zod.NewValidationError("string", data, "string is required")
    }

    if len(str) < s.minLength {
        return zod.NewValidationError("string", str, "string is too short")
    }

    if len(str) > s.maxLength {
        return zod.NewValidationError("string", str, "string is too long")
    }

    return nil
}
