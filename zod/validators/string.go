package validators

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/aymaneallaoui/zod-go/zod"
)

// StringSchema represents a string validation schema
type StringSchema struct {
	minLength      int
	maxLength      int
	required       bool
	pattern        *regexp.Regexp
	emailFormat    bool
	urlFormat      bool
	customFunc     func(string) error
	optional       bool
	defaultValue   *string
	customError    map[string]string
}

// String creates a new string schema
func String() *StringSchema {
	return &StringSchema{
		customError: make(map[string]string),
	}
}

// Min sets the minimum length requirement
func (s *StringSchema) Min(length int) *StringSchema {
	s.minLength = length
	return s
}

// Max sets the maximum length requirement
func (s *StringSchema) Max(length int) *StringSchema {
	s.maxLength = length
	return s
}

// Required marks the field as required (non-empty)
func (s *StringSchema) Required() *StringSchema {
	s.required = true
	s.optional = false
	return s
}

// Optional marks the field as optional
func (s *StringSchema) Optional() *StringSchema {
	s.optional = true
	s.required = false
	return s
}

// Default sets a default value for optional fields
func (s *StringSchema) Default(value string) *StringSchema {
	s.defaultValue = &value
	s.optional = true
	return s
}

// Pattern adds regex pattern validation
func (s *StringSchema) Pattern(pattern string) *StringSchema {
	s.pattern = regexp.MustCompile(pattern)
	return s
}

// Email validates email format
func (s *StringSchema) Email() *StringSchema {
	s.emailFormat = true
	return s
}

// URL validates URL format
func (s *StringSchema) URL() *StringSchema {
	s.urlFormat = true
	return s
}

// Custom adds a custom validation function
func (s *StringSchema) Custom(fn func(string) error) *StringSchema {
	s.customFunc = fn
	return s
}

// WithMessage sets a custom error message for a validation type
func (s *StringSchema) WithMessage(validationType, message string) *StringSchema {
	s.customError[validationType] = message
	return s
}

// getErrorMessage returns custom or default error message
func (s *StringSchema) getErrorMessage(validationType, defaultMessage string) string {
	if msg, exists := s.customError[validationType]; exists {
		return msg
	}
	return defaultMessage
}

// isValidEmail validates email format using a more comprehensive regex
func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email) && len(email) <= 254
}

// isValidURL validates URL format
func isValidURL(urlStr string) bool {
	u, err := url.Parse(urlStr)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// Validate performs validation on the provided data
func (s *StringSchema) Validate(data interface{}) error {
	// Handle nil values
	if data == nil {
		if s.required {
			return zod.NewValidationError("", nil, s.getErrorMessage("required", "field is required"))
		}
		if s.defaultValue != nil {
			// In a real implementation, you might want to modify the data
			// For now, we just validate the default
			return s.Validate(*s.defaultValue)
		}
		if s.optional {
			return nil
		}
		return zod.NewValidationError("", nil, s.getErrorMessage("required", "field is required"))
	}

	// Type check
	str, ok := data.(string)
	if !ok {
		return zod.NewValidationError(fmt.Sprintf("%v", data), data, 
			s.getErrorMessage("type", "invalid type, expected string"))
	}

	// Handle empty strings
	if str == "" {
		if s.required {
			return zod.NewValidationError("", str, 
				s.getErrorMessage("required", "string is required"))
		}
		if s.defaultValue != nil {
			return s.Validate(*s.defaultValue)
		}
		if s.optional {
			return nil
		}
	}

	// Length validations
	if s.minLength > 0 && len(str) < s.minLength {
		return zod.NewValidationError(str, str, 
			s.getErrorMessage("minLength", 
				fmt.Sprintf("string is too short, minimum length is %d", s.minLength)))
	}

	if s.maxLength > 0 && len(str) > s.maxLength {
		return zod.NewValidationError(str, str, 
			s.getErrorMessage("maxLength", 
				fmt.Sprintf("string is too long, maximum length is %d", s.maxLength)))
	}

	// Pattern validation
	if s.pattern != nil && !s.pattern.MatchString(str) {
		return zod.NewValidationError(str, str, 
			s.getErrorMessage("pattern", "string does not match required pattern"))
	}

	// Email validation
	if s.emailFormat && !isValidEmail(str) {
		return zod.NewValidationError(str, str, 
			s.getErrorMessage("email", "invalid email format"))
	}

	// URL validation
	if s.urlFormat && !isValidURL(str) {
		return zod.NewValidationError(str, str, 
			s.getErrorMessage("url", "invalid URL format"))
	}

	// Custom validation
	if s.customFunc != nil {
		if err := s.customFunc(str); err != nil {
			return err
		}
	}

	return nil
}

// StringLengthBetween creates a string schema with length constraints
func StringLengthBetween(min, max int) *StringSchema {
	return String().Min(min).Max(max)
}

// EmailString creates a string schema for email validation
func EmailString() *StringSchema {
	return String().Email().Required()
}

// URLString creates a string schema for URL validation
func URLString() *StringSchema {
	return String().URL().Required()
}

// NonEmptyString creates a required non-empty string schema
func NonEmptyString() *StringSchema {
	return String().Min(1).Required()
}

// OptionalString creates an optional string schema
func OptionalString() *StringSchema {
	return String().Optional()
}
