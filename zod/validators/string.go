package validators

import (
	"fmt"

	"github.com/aymaneallaoui/zod-Go/zod"
)

type StringSchema struct {
	minLength   int
	maxLength   int
	required    bool
	customError map[string]string
}

func String() *StringSchema {
	return &StringSchema{
		customError: make(map[string]string),
	}
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

func (s *StringSchema) WithMessage(validationType, message string) *StringSchema {
	s.customError[validationType] = message
	return s
}

func (s *StringSchema) getErrorMessage(validationType, defaultMessage string) string {
	if msg, exists := s.customError[validationType]; exists {
		return msg
	}
	return defaultMessage
}

func (s *StringSchema) Validate(data interface{}) error {
	str, ok := data.(string)
	if !ok {

		return zod.NewValidationError(fmt.Sprintf("%v", data), data, s.getErrorMessage("type", "invalid type, expected string"))
	}

	if s.required && str == "" {

		return zod.NewValidationError(fmt.Sprintf("%v", data), data, s.getErrorMessage("required", "string is required"))
	}

	if len(str) < s.minLength {

		return zod.NewValidationError(fmt.Sprintf("%v", str), str, s.getErrorMessage("minLength", fmt.Sprintf("The string is too short! Minimum length is %d", s.minLength)))
	}

	if s.maxLength > 0 && len(str) > s.maxLength {

		return zod.NewValidationError(fmt.Sprintf("%v", str), str, s.getErrorMessage("maxLength", fmt.Sprintf("The string is too long! Maximum length is %d", s.maxLength)))
	}

	return nil
}
