package tests

import (
	"strings"
	"testing"

	"github.com/aymaneallaoui/zod-go/zod"
	"github.com/aymaneallaoui/zod-go/zod/validators"
)

func TestBoolValidator_BasicValidation(t *testing.T) {
	schema := validators.Bool()

	tests := []struct {
		name      string
		input     interface{}
		wantError bool
		errorMsg  string
	}{
		{
			name:      "valid true",
			input:     true,
			wantError: false,
		},
		{
			name:      "valid false",
			input:     false,
			wantError: false,
		},
		{
			name:      "nil input optional",
			input:     nil,
			wantError: false, // Bool is optional by default
		},
		{
			name:      "invalid type",
			input:     "not a boolean",
			wantError: true,
			errorMsg:  "expected boolean",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := schema.Validate(tt.input)
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

func TestBoolValidator_RequiredValidation(t *testing.T) {
	schema := validators.Bool().Required()

	tests := []struct {
		name      string
		input     interface{}
		wantError bool
	}{
		{"valid required true", true, false},
		{"valid required false", false, false},
		{"nil required bool", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := schema.Validate(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("Required validation error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestBoolValidator_TypeConversion(t *testing.T) {
	schema := validators.Bool()

	tests := []struct {
		name      string
		input     interface{}
		expected  bool
		wantError bool
	}{
		// Boolean values
		{"true boolean", true, true, false},
		{"false boolean", false, false, false},

		// String representations
		{"string true", "true", true, false},
		{"string True", "True", true, false},
		{"string TRUE", "TRUE", true, false},
		{"string 1", "1", true, false},
		{"string yes", "yes", true, false},
		{"string Yes", "Yes", true, false},
		{"string YES", "YES", true, false},
		{"string on", "on", true, false},
		{"string On", "On", true, false},
		{"string ON", "ON", true, false},

		{"string false", "false", false, false},
		{"string False", "False", false, false},
		{"string FALSE", "FALSE", false, false},
		{"string 0", "0", false, false},
		{"string no", "no", false, false},
		{"string No", "No", false, false},
		{"string NO", "NO", false, false},
		{"string off", "off", false, false},
		{"string Off", "Off", false, false},
		{"string OFF", "OFF", false, false},

		// Numeric representations
		{"int 1", 1, true, false},
		{"int 0", 0, false, false},
		{"int 42", 42, true, false},
		{"int -1", -1, true, false},
		{"int8 1", int8(1), true, false},
		{"int16 0", int16(0), false, false},
		{"int32 1", int32(1), true, false},
		{"int64 0", int64(0), false, false},
		{"uint 1", uint(1), true, false},
		{"uint8 0", uint8(0), false, false},
		{"uint16 1", uint16(1), true, false},
		{"uint32 0", uint32(0), false, false},
		{"uint64 1", uint64(1), true, false},
		{"float32 1.0", float32(1.0), true, false},
		{"float32 0.0", float32(0.0), false, false},
		{"float64 1.0", float64(1.0), true, false},
		{"float64 0.0", float64(0.0), false, false},

		// Invalid conversions
		{"invalid string", "invalid", false, true},
		{"empty string", "", false, true},
		{"complex type", []int{1, 2, 3}, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := schema.Validate(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("Type conversion error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestBoolValidator_TrueValidation(t *testing.T) {
	schema := validators.Bool().True()

	tests := []struct {
		name      string
		input     interface{}
		wantError bool
	}{
		{"true value", true, false},
		{"string true", "true", false},
		{"int 1", 1, false},
		{"false value", false, true},
		{"string false", "false", true},
		{"int 0", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := schema.Validate(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("True validation error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestBoolValidator_FalseValidation(t *testing.T) {
	schema := validators.Bool().False()

	tests := []struct {
		name      string
		input     interface{}
		wantError bool
	}{
		{"false value", false, false},
		{"string false", "false", false},
		{"int 0", 0, false},
		{"true value", true, true},
		{"string true", "true", true},
		{"int 1", 1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := schema.Validate(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("False validation error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestBoolValidator_CustomValidation(t *testing.T) {
	// Custom validator that only allows true during business hours
	businessHoursValidator := func(b bool) error {
		if b {
			// In a real scenario, you'd check the actual time
			// For testing, let's assume it's business hours
			return nil
		}
		return zod.NewValidationError("bool", b, "feature only available during business hours")
	}

	schema := validators.Bool().Custom(businessHoursValidator)

	tests := []struct {
		name      string
		input     bool
		wantError bool
	}{
		{"true during business hours", true, false},
		{"false during business hours", false, true},
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

func TestBoolValidator_CustomMessages(t *testing.T) {
	schema := validators.Bool().
		Required().WithMessage("required", "Custom required error").
		True().WithMessage("true", "Custom true error")

	tests := []struct {
		name        string
		input       interface{}
		expectedMsg string
	}{
		{"custom required message", nil, "Custom required error"},
		{"custom true message", false, "Custom true error"},
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

func TestBoolValidator_DefaultValues(t *testing.T) {
	tests := []struct {
		name         string
		schema       *validators.BoolSchema
		input        interface{}
		expectError  bool
	}{
		{
			name:        "default true",
			schema:      validators.Bool().Default(true),
			input:       nil,
			expectError: false,
		},
		{
			name:        "default false",
			schema:      validators.Bool().Default(false),
			input:       nil,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.schema.Validate(tt.input)
			if (err != nil) != tt.expectError {
				t.Errorf("Default value validation error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

func TestBoolValidator_HelperMethods(t *testing.T) {
	t.Run("TrueBoolean", func(t *testing.T) {
		schema := validators.TrueBoolean()
		if err := schema.Validate(true); err != nil {
			t.Errorf("TrueBoolean validation failed: %v", err)
		}
		if err := schema.Validate(false); err == nil {
			t.Error("TrueBoolean should reject false")
		}
	})

	t.Run("FalseBoolean", func(t *testing.T) {
		schema := validators.FalseBoolean()
		if err := schema.Validate(false); err != nil {
			t.Errorf("FalseBoolean validation failed: %v", err)
		}
		if err := schema.Validate(true); err == nil {
			t.Error("FalseBoolean should reject true")
		}
	})

	t.Run("OptionalBoolean", func(t *testing.T) {
		schema := validators.OptionalBoolean()
		if err := schema.Validate(nil); err != nil {
			t.Errorf("OptionalBoolean should allow nil: %v", err)
		}
		if err := schema.Validate(true); err != nil {
			t.Errorf("OptionalBoolean validation failed: %v", err)
		}
	})

	t.Run("RequiredBoolean", func(t *testing.T) {
		schema := validators.RequiredBoolean()
		if err := schema.Validate(true); err != nil {
			t.Errorf("RequiredBoolean validation failed: %v", err)
		}
		if err := schema.Validate(false); err != nil {
			t.Errorf("RequiredBoolean validation failed: %v", err)
		}
		if err := schema.Validate(nil); err == nil {
			t.Error("RequiredBoolean should reject nil")
		}
	})

	t.Run("ConsentBoolean", func(t *testing.T) {
		schema := validators.ConsentBoolean()
		if err := schema.Validate(true); err != nil {
			t.Errorf("ConsentBoolean validation failed: %v", err)
		}
		if err := schema.Validate(false); err == nil {
			t.Error("ConsentBoolean should reject false")
		}
		if err := schema.Validate("yes"); err != nil {
			t.Errorf("ConsentBoolean should accept 'yes': %v", err)
		}
	})

	t.Run("ToggleBoolean", func(t *testing.T) {
		schema := validators.ToggleBoolean(true)
		if err := schema.Validate(nil); err != nil {
			t.Errorf("ToggleBoolean should use default: %v", err)
		}
		if err := schema.Validate(false); err != nil {
			t.Errorf("ToggleBoolean should accept any value: %v", err)
		}
	})
}

func TestBoolValidator_ChainedValidation(t *testing.T) {
	// Test chaining multiple validations
	schema := validators.Bool().
		Required().
		True()

	tests := []struct {
		name      string
		input     interface{}
		wantError bool
	}{
		{"valid chained", true, false},
		{"valid chained string", "true", false},
		{"invalid nil", nil, true},
		{"invalid false", false, true},
		{"invalid type", "invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := schema.Validate(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("Chained validation error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestBoolValidator_EdgeCases(t *testing.T) {
	schema := validators.Bool()

	tests := []struct {
		name      string
		input     interface{}
		wantError bool
	}{
		// Edge cases for string conversion
		{"empty space", " ", true},
		{"tab character", "\t", true},
		{"newline", "\n", true},
		{"mixed case", "TrUe", true}, // Only exact case matches work
		{"partial match", "tr", true},
		{"extra characters", "true ", true},
		
		// Edge cases for numeric conversion
		{"float 0.1", 0.1, true}, // Non-zero, non-one float
		{"negative zero", -0.0, false}, // Should be false
		{"very large number", 1000000, true}, // Non-zero
		{"very small positive", 0.0001, true}, // Non-zero
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := schema.Validate(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("Edge case validation error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestBoolValidator_ConflictingRequirements(t *testing.T) {
	// Test conflicting requirements (both True and False)
	// This should be handled gracefully - the last one wins or it's an error
	
	t.Run("True then False", func(t *testing.T) {
		schema := validators.Bool().True().False()
		
		// Should require false (last constraint wins)
		err := schema.Validate(false)
		if err != nil {
			t.Errorf("Conflicting requirements should use last constraint: %v", err)
		}
		
		err = schema.Validate(true)
		if err == nil {
			t.Error("Should reject true when false is required")
		}
	})

	t.Run("False then True", func(t *testing.T) {
		schema := validators.Bool().False().True()
		
		// Should require true (last constraint wins)
		err := schema.Validate(true)
		if err != nil {
			t.Errorf("Conflicting requirements should use last constraint: %v", err)
		}
		
		err = schema.Validate(false)
		if err == nil {
			t.Error("Should reject false when true is required")
		}
	})
}
