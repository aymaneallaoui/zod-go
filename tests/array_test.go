package tests

import (
	"strings"
	"testing"

	"github.com/aymaneallaoui/zod-go/zod"
	"github.com/aymaneallaoui/zod-go/zod/validators"
)

func TestArrayValidator_BasicValidation(t *testing.T) {
	elementSchema := validators.String().Required()
	schema := validators.Array(elementSchema)

	tests := []struct {
		name      string
		input     interface{}
		wantError bool
		errorMsg  string
	}{
		{
			name:      "valid string array",
			input:     []interface{}{"hello", "world"},
			wantError: false,
		},
		{
			name:      "valid string slice",
			input:     []string{"hello", "world"},
			wantError: false,
		},
		{
			name:      "empty array",
			input:     []interface{}{},
			wantError: false,
		},
		{
			name:      "nil array optional",
			input:     nil,
			wantError: false, // Arrays are optional by default
		},
		{
			name:      "invalid type",
			input:     "not an array",
			wantError: true,
			errorMsg:  "expected array",
		},
		{
			name:      "invalid element",
			input:     []interface{}{"hello", 123},
			wantError: true,
			errorMsg:  "invalid type", // Updated to match actual error format
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

func TestArrayValidator_RequiredValidation(t *testing.T) {
	elementSchema := validators.String().Required()
	schema := validators.Array(elementSchema).Required()

	tests := []struct {
		name      string
		input     interface{}
		wantError bool
	}{
		{"valid required array", []string{"hello"}, false},
		{"nil required array", nil, true},
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

func TestArrayValidator_LengthValidation(t *testing.T) {
	elementSchema := validators.String().Required()
	
	tests := []struct {
		name      string
		schema    *validators.ArraySchema
		input     []string
		wantError bool
		errorMsg  string
	}{
		{
			name:      "valid length range",
			schema:    validators.Array(elementSchema).Min(1).Max(3),
			input:     []string{"hello", "world"},
			wantError: false,
		},
		{
			name:      "too short",
			schema:    validators.Array(elementSchema).Min(3),
			input:     []string{"hello"},
			wantError: true,
			errorMsg:  "too short",
		},
		{
			name:      "too long",
			schema:    validators.Array(elementSchema).Max(2),
			input:     []string{"a", "b", "c"},
			wantError: true,
			errorMsg:  "too long",
		},
		{
			name:      "exact min length",
			schema:    validators.Array(elementSchema).Min(2),
			input:     []string{"a", "b"},
			wantError: false,
		},
		{
			name:      "exact max length",
			schema:    validators.Array(elementSchema).Max(2),
			input:     []string{"a", "b"},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.schema.Validate(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("Length validation error = %v, wantError %v", err, tt.wantError)
				return
			}
			if tt.wantError && err != nil && !strings.Contains(strings.ToLower(err.Error()), strings.ToLower(tt.errorMsg)) {
				t.Errorf("Expected error message to contain %q, got %q", tt.errorMsg, err.Error())
			}
		})
	}
}

func TestArrayValidator_NonEmptyValidation(t *testing.T) {
	elementSchema := validators.String().Required()
	schema := validators.Array(elementSchema).NonEmpty()

	tests := []struct {
		name      string
		input     []string
		wantError bool
	}{
		{"non-empty array", []string{"hello"}, false},
		{"empty array", []string{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := schema.Validate(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("NonEmpty validation error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestArrayValidator_UniqueValidation(t *testing.T) {
	elementSchema := validators.String().Required()
	schema := validators.Array(elementSchema).Unique()

	tests := []struct {
		name      string
		input     []interface{}
		wantError bool
	}{
		{"unique elements", []interface{}{"a", "b", "c"}, false},
		{"duplicate elements", []interface{}{"a", "b", "a"}, true},
		{"empty array", []interface{}{}, false},
		{"single element", []interface{}{"a"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := schema.Validate(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("Unique validation error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestArrayValidator_UniqueComplexTypes(t *testing.T) {
	elementSchema := validators.Object(map[string]zod.Schema{
		"id": validators.Number().Required(),
	})
	schema := validators.Array(elementSchema).Unique()

	tests := []struct {
		name      string
		input     []interface{}
		wantError bool
	}{
		{
			name: "unique objects",
			input: []interface{}{
				map[string]interface{}{"id": 1},
				map[string]interface{}{"id": 2},
			},
			wantError: false,
		},
		{
			name: "duplicate objects",
			input: []interface{}{
				map[string]interface{}{"id": 1},
				map[string]interface{}{"id": 1},
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := schema.Validate(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("Unique complex validation error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestArrayValidator_ElementValidation(t *testing.T) {
	numberSchema := validators.Number().Min(0).Max(100)
	schema := validators.Array(numberSchema)

	tests := []struct {
		name      string
		input     []interface{}
		wantError bool
	}{
		{"valid elements", []interface{}{10, 20, 30}, false},
		{"invalid element type", []interface{}{10, "invalid", 30}, true},
		{"invalid element range", []interface{}{10, 150, 30}, true},
		{"mixed valid/invalid", []interface{}{10, -5, 30}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := schema.Validate(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("Element validation error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestArrayValidator_NestedArrays(t *testing.T) {
	innerArraySchema := validators.Array(validators.String().Required()).Min(1)
	outerArraySchema := validators.Array(innerArraySchema).Min(1)

	tests := []struct {
		name      string
		input     []interface{}
		wantError bool
	}{
		{
			name: "valid nested arrays",
			input: []interface{}{
				[]interface{}{"a", "b"},
				[]interface{}{"c", "d"},
			},
			wantError: false,
		},
		{
			name: "invalid nested array",
			input: []interface{}{
				[]interface{}{"a", "b"},
				[]interface{}{}, // empty inner array
			},
			wantError: true,
		},
		{
			name: "invalid nested element",
			input: []interface{}{
				[]interface{}{"a", "b"},
				[]interface{}{"c", 123}, // invalid element type
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := outerArraySchema.Validate(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("Nested array validation error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestArrayValidator_CustomValidation(t *testing.T) {
	// Custom validator that checks if array has even number of elements
	evenLengthValidator := func(arr []interface{}) error {
		if len(arr)%2 != 0 {
			return zod.NewValidationError("array", arr, "array must have even number of elements")
		}
		return nil
	}

	schema := validators.Array(validators.String().Required()).Custom(evenLengthValidator)

	tests := []struct {
		name      string
		input     []string
		wantError bool
	}{
		{"even length", []string{"a", "b", "c", "d"}, false},
		{"odd length", []string{"a", "b", "c"}, true},
		{"empty array", []string{}, false}, // 0 is even
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

func TestArrayValidator_CustomMessages(t *testing.T) {
	elementSchema := validators.String().Required()
	schema := validators.Array(elementSchema).
		Min(2).WithMessage("minLength", "Custom min length error").
		Max(5).WithMessage("maxLength", "Custom max length error").
		NonEmpty().WithMessage("nonEmpty", "Custom non-empty error").
		Unique().WithMessage("unique", "Custom unique error")

	tests := []struct {
		name        string
		input       interface{}
		expectedMsg string
	}{
		{"custom min length message", []string{"a"}, "Custom min length error"},
		{"custom max length message", []string{"a", "b", "c", "d", "e", "f"}, "Custom max length error"},
		{"custom non-empty message", []string{}, "Custom non-empty error"},
		{"custom unique message", []string{"a", "a"}, "Custom unique error"},
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

func TestArrayValidator_TypeConversion(t *testing.T) {
	elementSchema := validators.Number().Required()
	schema := validators.Array(elementSchema)

	tests := []struct {
		name      string
		input     interface{}
		wantError bool
	}{
		{"[]interface{}", []interface{}{1, 2, 3}, false},
		{"[]int", []int{1, 2, 3}, false},
		{"[]float64", []float64{1.0, 2.0, 3.0}, false},
		{"[3]int array", [3]int{1, 2, 3}, false},
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

func TestArrayValidator_HelperMethods(t *testing.T) {
	t.Run("StringArray", func(t *testing.T) {
		schema := validators.StringArray()
		if err := schema.Validate([]string{"hello", "world"}); err != nil {
			t.Errorf("StringArray validation failed: %v", err)
		}
		if err := schema.Validate([]interface{}{"hello", 123}); err == nil {
			t.Error("StringArray should reject non-string elements")
		}
	})

	t.Run("NumberArray", func(t *testing.T) {
		schema := validators.NumberArray()
		if err := schema.Validate([]interface{}{1, 2, 3}); err != nil {
			t.Errorf("NumberArray validation failed: %v", err)
		}
		if err := schema.Validate([]interface{}{1, "2", 3}); err == nil {
			t.Error("NumberArray should reject non-number elements")
		}
	})

	t.Run("IntegerArray", func(t *testing.T) {
		schema := validators.IntegerArray()
		if err := schema.Validate([]interface{}{1, 2, 3}); err != nil {
			t.Errorf("IntegerArray validation failed: %v", err)
		}
		if err := schema.Validate([]interface{}{1, 2.5, 3}); err == nil {
			t.Error("IntegerArray should reject non-integer elements")
		}
	})

	t.Run("UniqueStringArray", func(t *testing.T) {
		schema := validators.UniqueStringArray()
		if err := schema.Validate([]string{"a", "b", "c"}); err != nil {
			t.Errorf("UniqueStringArray validation failed: %v", err)
		}
		if err := schema.Validate([]string{"a", "b", "a"}); err == nil {
			t.Error("UniqueStringArray should reject duplicate elements")
		}
	})

	t.Run("NonEmptyArray", func(t *testing.T) {
		schema := validators.NonEmptyArray(validators.String().Required())
		if err := schema.Validate([]string{"hello"}); err != nil {
			t.Errorf("NonEmptyArray validation failed: %v", err)
		}
		if err := schema.Validate([]string{}); err == nil {
			t.Error("NonEmptyArray should reject empty arrays")
		}
	})

	t.Run("FixedLengthArray", func(t *testing.T) {
		schema := validators.FixedLengthArray(3, validators.String().Required())
		if err := schema.Validate([]string{"a", "b", "c"}); err != nil {
			t.Errorf("FixedLengthArray validation failed: %v", err)
		}
		if err := schema.Validate([]string{"a", "b"}); err == nil {
			t.Error("FixedLengthArray should reject wrong length")
		}
	})

	t.Run("BoundedArray", func(t *testing.T) {
		schema := validators.BoundedArray(1, 3, validators.String().Required())
		if err := schema.Validate([]string{"a", "b"}); err != nil {
			t.Errorf("BoundedArray validation failed: %v", err)
		}
		if err := schema.Validate([]string{}); err == nil {
			t.Error("BoundedArray should reject too short arrays")
		}
		if err := schema.Validate([]string{"a", "b", "c", "d"}); err == nil {
			t.Error("BoundedArray should reject too long arrays")
		}
	})
}

func TestArrayValidator_DefaultValues(t *testing.T) {
	defaultArray := []interface{}{"default1", "default2"}
	schema := validators.Array(validators.String().Required()).Default(defaultArray)
	
	// Test with nil - should use default
	err := schema.Validate(nil)
	if err != nil {
		t.Errorf("Expected no error for nil with default, got %v", err)
	}
}

func TestArrayValidator_ComplexChaining(t *testing.T) {
	// Complex validation: array of unique integers between 1-100, length 3-5
	schema := validators.Array(validators.Number().Integer().Min(1).Max(100).Required()).
		Min(3).
		Max(5).
		Unique().
		Required()

	tests := []struct {
		name      string
		input     interface{}
		wantError bool
	}{
		{"valid complex array", []interface{}{10, 25, 50, 75}, false},
		{"duplicate elements", []interface{}{10, 25, 10}, true},
		{"out of range element", []interface{}{10, 25, 150}, true},
		{"too short", []interface{}{10, 25}, true},
		{"too long", []interface{}{10, 25, 50, 75, 90, 95}, true},
		{"non-integer element", []interface{}{10, 25.5, 50}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := schema.Validate(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("Complex chaining validation error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}
