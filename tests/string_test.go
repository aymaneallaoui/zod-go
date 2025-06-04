package tests

import (
	"strings"
	"testing"

	"github.com/aymaneallaoui/zod-go/zod"
	"github.com/aymaneallaoui/zod-go/zod/validators"
)

func TestStringValidator_BasicValidation(t *testing.T) {
	tests := []struct {
		name      string
		schema    *validators.StringSchema
		input     interface{}
		wantError bool
		errorMsg  string
	}{
		{
			name:      "valid string",
			schema:    validators.String().Required(),
			input:     "hello",
			wantError: false,
		},
		{
			name:      "empty string required",
			schema:    validators.String().Required(),
			input:     "",
			wantError: true,
			errorMsg:  "required",
		},
		{
			name:      "nil input required",
			schema:    validators.String().Required(),
			input:     nil,
			wantError: true,
			errorMsg:  "required",
		},
		{
			name:      "invalid type",
			schema:    validators.String().Required(),
			input:     123,
			wantError: true,
			errorMsg:  "expected string",
		},
		{
			name:      "optional empty string",
			schema:    validators.String().Optional(),
			input:     "",
			wantError: false,
		},
		{
			name:      "optional nil",
			schema:    validators.String().Optional(),
			input:     nil,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.schema.Validate(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("Validate() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if tt.wantError && err != nil && !strings.Contains(err.Error(), tt.errorMsg) {
				t.Errorf("Expected error message to contain %q, got %q", tt.errorMsg, err.Error())
			}
		})
	}
}

func TestStringValidator_LengthValidation(t *testing.T) {
	tests := []struct {
		name      string
		schema    *validators.StringSchema
		input     string
		wantError bool
		errorMsg  string
	}{
		{
			name:      "valid length",
			schema:    validators.String().Min(3).Max(10),
			input:     "hello",
			wantError: false,
		},
		{
			name:      "too short",
			schema:    validators.String().Min(5),
			input:     "hi",
			wantError: true,
			errorMsg:  "too short",
		},
		{
			name:      "too long",
			schema:    validators.String().Max(5),
			input:     "this is too long",
			wantError: true,
			errorMsg:  "too long",
		},
		{
			name:      "exact min length",
			schema:    validators.String().Min(5),
			input:     "hello",
			wantError: false,
		},
		{
			name:      "exact max length",
			schema:    validators.String().Max(5),
			input:     "hello",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.schema.Validate(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("Validate() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if tt.wantError && err != nil && !strings.Contains(strings.ToLower(err.Error()), strings.ToLower(tt.errorMsg)) {
				t.Errorf("Expected error message to contain %q, got %q", tt.errorMsg, err.Error())
			}
		})
	}
}

func TestStringValidator_EmailValidation(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
	}{
		{"valid email", "user@example.com", false},
		{"valid email with subdomain", "user@mail.example.com", false},
		{"valid email with plus", "user+test@example.com", false},
		{"valid email with numbers", "user123@example123.com", false},
		{"invalid no @", "userexample.com", true},
		{"invalid no domain", "user@", true},
		{"invalid no user", "@example.com", true},
		{"invalid no tld", "user@example", true},
		{"invalid spaces", "user @example.com", true},
		{"invalid double @", "user@@example.com", true},
		{"empty string", "", true},
	}

	schema := validators.String().Email()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := schema.Validate(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("Email validation error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestStringValidator_URLValidation(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
	}{
		{"valid http URL", "http://example.com", false},
		{"valid https URL", "https://example.com", false},
		{"valid URL with path", "https://example.com/path", false},
		{"valid URL with query", "https://example.com?query=value", false},
		{"valid URL with port", "http://example.com:8080", false},
		{"invalid no scheme", "example.com", true},
		{"invalid no host", "http://", true},
		{"invalid scheme only", "http", true},
		{"empty string", "", true},
		{"invalid characters", "http://ex ample.com", true},
	}

	schema := validators.String().URL()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := schema.Validate(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("URL validation error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestStringValidator_PatternValidation(t *testing.T) {
	tests := []struct {
		name      string
		pattern   string
		input     string
		wantError bool
	}{
		{"alphanumeric pattern match", `^[a-zA-Z0-9]+$`, "hello123", false},
		{"alphanumeric pattern no match", `^[a-zA-Z0-9]+$`, "hello-123", true},
		{"email pattern match", `^[^@]+@[^@]+\.[^@]+$`, "user@example.com", false},
		{"email pattern no match", `^[^@]+@[^@]+\.[^@]+$`, "invalid-email", true},
		{"phone pattern match", `^\d{3}-\d{3}-\d{4}$`, "123-456-7890", false},
		{"phone pattern no match", `^\d{3}-\d{3}-\d{4}$`, "1234567890", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema := validators.String().Pattern(tt.pattern)
			err := schema.Validate(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("Pattern validation error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestStringValidator_CustomMessages(t *testing.T) {
	schema := validators.String().
		Min(5).
		WithMessage("minLength", "Custom min length error").
		Max(10).
		WithMessage("maxLength", "Custom max length error").
		Required().
		WithMessage("required", "Custom required error")

	tests := []struct {
		name        string
		input       interface{}
		expectedMsg string
	}{
		{"custom required message", nil, "Custom required error"},
		{"custom min length message", "hi", "Custom min length error"},
		{"custom max length message", "this is way too long", "Custom max length error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := schema.Validate(tt.input)
			if err == nil {
				t.Errorf("Expected validation error, got nil")
				return
			}
			if !strings.Contains(err.Error(), tt.expectedMsg) {
				t.Errorf("Expected error message to contain %q, got %q", tt.expectedMsg, err.Error())
			}
		})
	}
}

func TestStringValidator_DefaultValues(t *testing.T) {
	schema := validators.String().Default("default value")

	// This test assumes the validator would somehow apply defaults
	// In a real implementation, you might need a separate method to get the processed value
	err := schema.Validate(nil)
	if err != nil {
		t.Errorf("Expected no error for nil with default, got %v", err)
	}
}

func TestStringValidator_CustomValidation(t *testing.T) {
	// Custom validator that checks if string contains "admin"
	adminValidator := func(s string) error {
		if strings.Contains(strings.ToLower(s), "admin") {
			return zod.NewValidationError("string", s, "string cannot contain 'admin'")
		}
		return nil
	}

	schema := validators.String().Custom(adminValidator)

	tests := []struct {
		name      string
		input     string
		wantError bool
	}{
		{"no admin", "regular_user", false},
		{"contains admin", "admin_user", true},
		{"contains Admin", "SuperAdmin", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := schema.Validate(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("Custom validation error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestStringValidator_HelperMethods(t *testing.T) {
	t.Run("EmailString", func(t *testing.T) {
		schema := validators.EmailString()
		if err := schema.Validate("user@example.com"); err != nil {
			t.Errorf("EmailString validation failed: %v", err)
		}
		if err := schema.Validate("invalid"); err == nil {
			t.Error("EmailString should reject invalid email")
		}
	})

	t.Run("URLString", func(t *testing.T) {
		schema := validators.URLString()
		if err := schema.Validate("https://example.com"); err != nil {
			t.Errorf("URLString validation failed: %v", err)
		}
		if err := schema.Validate("invalid"); err == nil {
			t.Error("URLString should reject invalid URL")
		}
	})

	t.Run("NonEmptyString", func(t *testing.T) {
		schema := validators.NonEmptyString()
		if err := schema.Validate("hello"); err != nil {
			t.Errorf("NonEmptyString validation failed: %v", err)
		}
		if err := schema.Validate(""); err == nil {
			t.Error("NonEmptyString should reject empty string")
		}
	})

	t.Run("StringLengthBetween", func(t *testing.T) {
		schema := validators.StringLengthBetween(5, 10)
		if err := schema.Validate("hello"); err != nil {
			t.Errorf("StringLengthBetween validation failed: %v", err)
		}
		if err := schema.Validate("hi"); err == nil {
			t.Error("StringLengthBetween should reject short string")
		}
	})
}
