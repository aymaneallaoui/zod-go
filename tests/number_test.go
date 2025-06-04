package tests

import (
	"math"
	"strings"
	"testing"

	"github.com/aymaneallaoui/zod-go/zod"
	"github.com/aymaneallaoui/zod-go/zod/validators"
)

func TestNumberValidator_BasicValidation(t *testing.T) {
	tests := []struct {
		name      string
		schema    *validators.NumberSchema
		input     interface{}
		wantError bool
		errorMsg  string
	}{
		{
			name:      "valid int",
			schema:    validators.Number().Required(),
			input:     42,
			wantError: false,
		},
		{
			name:      "valid float64",
			schema:    validators.Number().Required(),
			input:     42.5,
			wantError: false,
		},
		{
			name:      "valid float32",
			schema:    validators.Number().Required(),
			input:     float32(42.5),
			wantError: false,
		},
		{
			name:      "nil input required",
			schema:    validators.Number().Required(),
			input:     nil,
			wantError: true,
			errorMsg:  "required",
		},
		{
			name:      "invalid type",
			schema:    validators.Number().Required(),
			input:     "not a number",
			wantError: true,
			errorMsg:  "expected number",
		},
		{
			name:      "optional nil",
			schema:    validators.Number().Optional(),
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
			if tt.wantError && err != nil && !strings.Contains(strings.ToLower(err.Error()), strings.ToLower(tt.errorMsg)) {
				t.Errorf("Expected error message to contain %q, got %q", tt.errorMsg, err.Error())
			}
		})
	}
}

func TestNumberValidator_TypeConversion(t *testing.T) {
	schema := validators.Number()

	tests := []struct {
		name      string
		input     interface{}
		wantError bool
	}{
		{"int8", int8(42), false},
		{"int16", int16(42), false},
		{"int32", int32(42), false},
		{"int64", int64(42), false},
		{"uint", uint(42), false},
		{"uint8", uint8(42), false},
		{"uint16", uint16(42), false},
		{"uint32", uint32(42), false},
		{"uint64", uint64(42), false},
		{"float32", float32(42.5), false},
		{"float64", float64(42.5), false},
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

func TestNumberValidator_RangeValidation(t *testing.T) {
	tests := []struct {
		name      string
		schema    *validators.NumberSchema
		input     float64
		wantError bool
		errorMsg  string
	}{
		{
			name:      "valid range",
			schema:    validators.Number().Min(0).Max(100),
			input:     50,
			wantError: false,
		},
		{
			name:      "below minimum",
			schema:    validators.Number().Min(10),
			input:     5,
			wantError: true,
			errorMsg:  "at least",
		},
		{
			name:      "above maximum",
			schema:    validators.Number().Max(100),
			input:     150,
			wantError: true,
			errorMsg:  "at most",
		},
		{
			name:      "exact minimum",
			schema:    validators.Number().Min(10),
			input:     10,
			wantError: false,
		},
		{
			name:      "exact maximum",
			schema:    validators.Number().Max(100),
			input:     100,
			wantError: false,
		},
		{
			name:      "negative range",
			schema:    validators.Number().Min(-100).Max(-10),
			input:     -50,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.schema.Validate(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("Range validation error = %v, wantError %v", err, tt.wantError)
				return
			}
			if tt.wantError && err != nil && !strings.Contains(strings.ToLower(err.Error()), strings.ToLower(tt.errorMsg)) {
				t.Errorf("Expected error message to contain %q, got %q", tt.errorMsg, err.Error())
			}
		})
	}
}

func TestNumberValidator_IntegerValidation(t *testing.T) {
	schema := validators.Number().Integer()

	tests := []struct {
		name      string
		input     interface{}
		wantError bool
	}{
		{"valid integer", 42, false},
		{"valid integer float", 42.0, false},
		{"invalid float", 42.5, true},
		{"negative integer", -42, false},
		{"zero", 0, false},
		{"large integer", 1000000, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := schema.Validate(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("Integer validation error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestNumberValidator_SignValidation(t *testing.T) {
	tests := []struct {
		name      string
		schema    *validators.NumberSchema
		input     float64
		wantError bool
	}{
		{"positive number valid", validators.Number().Positive(), 5, false},
		{"positive number zero invalid", validators.Number().Positive(), 0, true},
		{"positive number negative invalid", validators.Number().Positive(), -5, true},
		{"negative number valid", validators.Number().Negative(), -5, false},
		{"negative number zero invalid", validators.Number().Negative(), 0, true},
		{"negative number positive invalid", validators.Number().Negative(), 5, true},
		{"non-zero positive valid", validators.Number().NonZero(), 5, false},
		{"non-zero negative valid", validators.Number().NonZero(), -5, false},
		{"non-zero zero invalid", validators.Number().NonZero(), 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.schema.Validate(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("Sign validation error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestNumberValidator_MultipleOfValidation(t *testing.T) {
	tests := []struct {
		name      string
		multiple  float64
		input     float64
		wantError bool
	}{
		{"multiple of 5", 5, 15, false},
		{"multiple of 5 invalid", 5, 16, true},
		{"multiple of 0.5", 0.5, 2.5, false},
		{"multiple of 0.5 invalid", 0.5, 2.6, true},
		{"multiple of 1", 1, 42, false},
		{"multiple of 2", 2, 8, false},
		{"multiple of 2 invalid", 2, 9, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema := validators.Number().MultipleOf(tt.multiple)
			err := schema.Validate(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("MultipleOf validation error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestNumberValidator_SpecialValues(t *testing.T) {
	schema := validators.Number()

	tests := []struct {
		name      string
		input     float64
		wantError bool
		errorMsg  string
	}{
		{"NaN", math.NaN(), true, "NaN"},
		{"positive infinity", math.Inf(1), true, "infinite"},
		{"negative infinity", math.Inf(-1), true, "infinite"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := schema.Validate(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("Special value validation error = %v, wantError %v", err, tt.wantError)
				return
			}
			if tt.wantError && err != nil && !strings.Contains(strings.ToLower(err.Error()), strings.ToLower(tt.errorMsg)) {
				t.Errorf("Expected error message to contain %q, got %q", tt.errorMsg, err.Error())
			}
		})
	}
}

func TestNumberValidator_CustomMessages(t *testing.T) {
	schema := validators.Number().
		Min(0).WithMessage("min", "Custom min error").
		Max(100).WithMessage("max", "Custom max error").
		Required().WithMessage("required", "Custom required error")

	tests := []struct {
		name        string
		input       interface{}
		expectedMsg string
	}{
		{"custom required message", nil, "Custom required error"},
		{"custom min message", -5, "Custom min error"},
		{"custom max message", 150, "Custom max error"},
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

func TestNumberValidator_CustomValidation(t *testing.T) {
	// Custom validator that checks if number is prime
	isPrime := func(n float64) error {
		if n != float64(int(n)) || n < 2 {
			return zod.NewValidationError("number", n, "must be a prime number")
		}
		num := int(n)
		for i := 2; i*i <= num; i++ {
			if num%i == 0 {
				return zod.NewValidationError("number", n, "must be a prime number")
			}
		}
		return nil
	}

	schema := validators.Number().Custom(isPrime)

	tests := []struct {
		name      string
		input     float64
		wantError bool
	}{
		{"prime number 2", 2, false},
		{"prime number 3", 3, false},
		{"prime number 17", 17, false},
		{"composite number 4", 4, true},
		{"composite number 9", 9, true},
		{"number 1", 1, true},
		{"float number", 2.5, true},
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

func TestNumberValidator_HelperMethods(t *testing.T) {
	t.Run("IntegerNumber", func(t *testing.T) {
		schema := validators.IntegerNumber()
		if err := schema.Validate(42); err != nil {
			t.Errorf("IntegerNumber validation failed: %v", err)
		}
		if err := schema.Validate(42.5); err == nil {
			t.Error("IntegerNumber should reject float")
		}
	})

	t.Run("PositiveNumber", func(t *testing.T) {
		schema := validators.PositiveNumber()
		if err := schema.Validate(5); err != nil {
			t.Errorf("PositiveNumber validation failed: %v", err)
		}
		if err := schema.Validate(-5); err == nil {
			t.Error("PositiveNumber should reject negative")
		}
	})

	t.Run("NegativeNumber", func(t *testing.T) {
		schema := validators.NegativeNumber()
		if err := schema.Validate(-5); err != nil {
			t.Errorf("NegativeNumber validation failed: %v", err)
		}
		if err := schema.Validate(5); err == nil {
			t.Error("NegativeNumber should reject positive")
		}
	})

	t.Run("NonZeroNumber", func(t *testing.T) {
		schema := validators.NonZeroNumber()
		if err := schema.Validate(5); err != nil {
			t.Errorf("NonZeroNumber validation failed: %v", err)
		}
		if err := schema.Validate(0); err == nil {
			t.Error("NonZeroNumber should reject zero")
		}
	})

	t.Run("NumberBetween", func(t *testing.T) {
		schema := validators.NumberBetween(0, 100)
		if err := schema.Validate(50); err != nil {
			t.Errorf("NumberBetween validation failed: %v", err)
		}
		if err := schema.Validate(150); err == nil {
			t.Error("NumberBetween should reject out of range")
		}
	})

	t.Run("PercentageNumber", func(t *testing.T) {
		schema := validators.PercentageNumber()
		if err := schema.Validate(50); err != nil {
			t.Errorf("PercentageNumber validation failed: %v", err)
		}
		if err := schema.Validate(150); err == nil {
			t.Error("PercentageNumber should reject > 100")
		}
	})

	t.Run("ProbabilityNumber", func(t *testing.T) {
		schema := validators.ProbabilityNumber()
		if err := schema.Validate(0.5); err != nil {
			t.Errorf("ProbabilityNumber validation failed: %v", err)
		}
		if err := schema.Validate(1.5); err == nil {
			t.Error("ProbabilityNumber should reject > 1")
		}
	})
}

func TestNumberValidator_DefaultValues(t *testing.T) {
	schema := validators.Number().Default(42.0)

	// Test with nil - should use default
	err := schema.Validate(nil)
	if err != nil {
		t.Errorf("Expected no error for nil with default, got %v", err)
	}
}

func TestNumberValidator_ChainedValidation(t *testing.T) {
	// Complex chained validation
	schema := validators.Number().
		Min(0).
		Max(100).
		Integer().
		MultipleOf(5).
		Required()

	tests := []struct {
		name      string
		input     interface{}
		wantError bool
	}{
		{"valid chained", 25, false},
		{"valid boundary", 100, false},
		{"invalid too small", -5, true},
		{"invalid too large", 105, true},
		{"invalid not integer", 25.5, true},
		{"invalid not multiple", 23, true},
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
