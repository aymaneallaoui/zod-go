// Package validators - Legacy Compatibility Layer
// This file provides backwards compatibility for existing code while encouraging migration to the new enhanced API.
//
// DEPRECATED: These functions are provided for backwards compatibility only.

package validators

import (
	"regexp"
)

// Legacy StringSchema - DEPRECATED
// Use validators.String() instead for the enhanced API with type safety
type StringSchema struct {
	minLength    int
	maxLength    int
	required     bool
	pattern      *regexp.Regexp
	emailFormat  bool
	urlFormat    bool
	customFunc   func(string) error
	optional     bool
	defaultValue *string
	customError  map[string]string
}

// Legacy NumberSchema - DEPRECATED
// Use validators.Number() instead for the enhanced API with type safety
type NumberSchema struct {
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

// Legacy ArraySchema - DEPRECATED
// Use validators.Array() instead for the enhanced API with type safety
type ArraySchema struct {
	elementSchema interface{}
	minItems      int
	maxItems      int
	contains      interface{}
	customFunc    func([]interface{}) error
	required      bool
	optional      bool
	defaultValue  []interface{}
	customError   map[string]string
}

// Legacy BoolSchema - DEPRECATED
// Use validators.Bool() instead for the enhanced API with type safety
type BoolSchema struct {
	customFunc   func(bool) error
	required     bool
	optional     bool
	defaultValue *bool
	customError  map[string]string
}

// Legacy constructor functions - DEPRECATED
// These are provided for backwards compatibility only

// LegacyString creates a legacy string schema - DEPRECATED
// Use validators.String() instead for enhanced DX
func LegacyString() *StringSchema {
	return &StringSchema{
		customError: make(map[string]string),
	}
}

// LegacyNumber creates a legacy number schema - DEPRECATED
// Use validators.Number() instead for enhanced DX
func LegacyNumber() *NumberSchema {
	return &NumberSchema{
		customError: make(map[string]string),
	}
}

// LegacyArray creates a legacy array schema - DEPRECATED
// Use validators.Array() instead for enhanced DX
func LegacyArray(elementSchema interface{}) *ArraySchema {
	return &ArraySchema{
		elementSchema: elementSchema,
		customError:   make(map[string]string),
	}
}

// LegacyBool creates a legacy bool schema - DEPRECATED
// Use validators.Bool() instead for enhanced DX
func LegacyBool() *BoolSchema {
	return &BoolSchema{
		customError: make(map[string]string),
	}
}

// Legacy StringSchema methods - DEPRECATED
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
	s.optional = false
	return s
}

func (s *StringSchema) Optional() *StringSchema {
	s.optional = true
	s.required = false
	return s
}

func (s *StringSchema) Default(value string) *StringSchema {
	s.defaultValue = &value
	s.optional = true
	return s
}

func (s *StringSchema) Pattern(pattern string) *StringSchema {
	s.pattern = regexp.MustCompile(pattern)
	return s
}

func (s *StringSchema) Email() *StringSchema {
	s.emailFormat = true
	return s
}

func (s *StringSchema) URL() *StringSchema {
	s.urlFormat = true
	return s
}

func (s *StringSchema) Custom(fn func(string) error) *StringSchema {
	s.customFunc = fn
	return s
}

func (s *StringSchema) WithMessage(validationType, message string) *StringSchema {
	s.customError[validationType] = message
	return s
}

func (s *StringSchema) Validate(data interface{}) error {
	// Use the new implementation internally for consistency
	newSchema := String()

	// Copy settings to new schema
	if s.minLength > 0 {
		newSchema = newSchema.Min(s.minLength)
	}
	if s.maxLength > 0 {
		newSchema = newSchema.Max(s.maxLength)
	}
	if s.pattern != nil {
		newSchema = newSchema.Pattern(s.pattern.String())
	}
	if s.emailFormat {
		newSchema = newSchema.Email()
	}
	if s.urlFormat {
		newSchema = newSchema.URL()
	}
	if s.customFunc != nil {
		newSchema = newSchema.Custom(s.customFunc)
	}

	// Apply error messages
	for key, msg := range s.customError {
		newSchema = newSchema.WithMessage(key, msg)
	}

	// Apply state
	if s.required {
		return newSchema.Required().Validate(data)
	} else if s.optional {
		optSchema := newSchema.Optional()
		if s.defaultValue != nil {
			optSchema = optSchema.Default(*s.defaultValue)
		}
		return optSchema.Validate(data)
	}

	// Default to required if neither is explicitly set
	return newSchema.Required().Validate(data)
}

// Similar legacy implementations for other types...
// (Abbreviated for brevity, but would follow the same pattern)

// Compatibility warning function
func init() {
	// This could log a warning about using legacy API, but we'll keep it quiet for now
	// to avoid breaking existing code
}
