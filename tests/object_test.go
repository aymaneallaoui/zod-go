package tests

import (
	"strings"
	"testing"

	"github.com/aymaneallaoui/zod-go/zod"
	"github.com/aymaneallaoui/zod-go/zod/validators"
)

func TestObjectValidator_BasicValidation(t *testing.T) {
	schema := validators.Object(map[string]zod.Schema{
		"name": validators.String().Required(),
		"age":  validators.Number().Required(),
	})

	tests := []struct {
		name      string
		input     interface{}
		wantError bool
		errorMsg  string
	}{
		{
			name: "valid object",
			input: map[string]interface{}{
				"name": "John",
				"age":  30,
			},
			wantError: false,
		},
		{
			name:      "nil input optional",
			input:     nil,
			wantError: false, // Objects are optional by default
		},
		{
			name:      "invalid type",
			input:     "not an object",
			wantError: true,
			errorMsg:  "expected object",
		},
		{
			name: "missing required field",
			input: map[string]interface{}{
				"name": "John",
				// age is missing
			},
			wantError: true,
			errorMsg:  "required",
		},
		{
			name: "invalid field type",
			input: map[string]interface{}{
				"name": 123, // should be string
				"age":  30,
			},
			wantError: true,
			errorMsg:  "invalid",
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

func TestObjectValidator_RequiredValidation(t *testing.T) {
	schema := validators.Object(map[string]zod.Schema{
		"name": validators.String().Required(),
	}).Required()

	tests := []struct {
		name      string
		input     interface{}
		wantError bool
	}{
		{"valid required object", map[string]interface{}{"name": "John"}, false},
		{"nil required object", nil, true},
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

func TestObjectValidator_StrictMode(t *testing.T) {
	schema := validators.Object(map[string]zod.Schema{
		"name": validators.String().Required(),
		"age":  validators.Number().Required(),
	}).Strict()

	tests := []struct {
		name      string
		input     map[string]interface{}
		wantError bool
	}{
		{
			name: "valid strict object",
			input: map[string]interface{}{
				"name": "John",
				"age":  30,
			},
			wantError: false,
		},
		{
			name: "unknown field in strict mode",
			input: map[string]interface{}{
				"name":    "John",
				"age":     30,
				"unknown": "field",
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := schema.Validate(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("Strict mode validation error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestObjectValidator_AllowUnknown(t *testing.T) {
	schema := validators.Object(map[string]zod.Schema{
		"name": validators.String().Required(),
	}).AllowUnknown()

	tests := []struct {
		name      string
		input     map[string]interface{}
		wantError bool
	}{
		{
			name: "unknown field allowed",
			input: map[string]interface{}{
				"name":    "John",
				"unknown": "field",
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := schema.Validate(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("AllowUnknown validation error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestObjectValidator_NestedObjects(t *testing.T) {
	addressSchema := validators.Object(map[string]zod.Schema{
		"street": validators.String().Required(),
		"city":   validators.String().Required(),
	})

	userSchema := validators.Object(map[string]zod.Schema{
		"name":    validators.String().Required(),
		"age":     validators.Number().Required(),
		"address": addressSchema.Required(),
	})

	tests := []struct {
		name      string
		input     map[string]interface{}
		wantError bool
	}{
		{
			name: "valid nested object",
			input: map[string]interface{}{
				"name": "John",
				"age":  30,
				"address": map[string]interface{}{
					"street": "123 Main St",
					"city":   "New York",
				},
			},
			wantError: false,
		},
		{
			name: "invalid nested object",
			input: map[string]interface{}{
				"name": "John",
				"age":  30,
				"address": map[string]interface{}{
					"street": "123 Main St",
					// city is missing
				},
			},
			wantError: true,
		},
		{
			name: "invalid nested field type",
			input: map[string]interface{}{
				"name": "John",
				"age":  30,
				"address": map[string]interface{}{
					"street": 123, // should be string
					"city":   "New York",
				},
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := userSchema.Validate(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("Nested object validation error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestObjectValidator_OptionalFields(t *testing.T) {
	schema := validators.Object(map[string]zod.Schema{
		"name": validators.String().Required(),
		"age":  validators.Number().Optional(),
		"bio":  validators.String().Optional(),
	})

	tests := []struct {
		name      string
		input     map[string]interface{}
		wantError bool
	}{
		{
			name: "all fields present",
			input: map[string]interface{}{
				"name": "John",
				"age":  30,
				"bio":  "A developer",
			},
			wantError: false,
		},
		{
			name: "optional fields missing",
			input: map[string]interface{}{
				"name": "John",
			},
			wantError: false,
		},
		{
			name: "some optional fields present",
			input: map[string]interface{}{
				"name": "John",
				"age":  30,
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := schema.Validate(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("Optional fields validation error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestObjectValidator_CustomValidation(t *testing.T) {
	// Custom validator that checks if age and birth year are consistent
	ageConsistencyValidator := func(obj map[string]interface{}) error {
		age, hasAge := obj["age"]
		birthYear, hasBirthYear := obj["birthYear"]

		if hasAge && hasBirthYear {
			ageVal, ok1 := age.(float64)
			birthYearVal, ok2 := birthYear.(float64)

			if ok1 && ok2 {
				currentYear := 2024.0
				expectedAge := currentYear - birthYearVal
				if ageVal != expectedAge {
					return zod.NewValidationError("age", obj, "age and birth year are inconsistent")
				}
			}
		}
		return nil
	}

	schema := validators.Object(map[string]zod.Schema{
		"name":      validators.String().Required(),
		"age":       validators.Number().Required(),
		"birthYear": validators.Number().Required(),
	}).Custom(ageConsistencyValidator)

	tests := []struct {
		name      string
		input     map[string]interface{}
		wantError bool
	}{
		{
			name: "consistent age and birth year",
			input: map[string]interface{}{
				"name":      "John",
				"age":       30.0,
				"birthYear": 1994.0,
			},
			wantError: false,
		},
		{
			name: "inconsistent age and birth year",
			input: map[string]interface{}{
				"name":      "John",
				"age":       25.0,
				"birthYear": 1994.0,
			},
			wantError: true,
		},
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

func TestObjectValidator_FieldManagement(t *testing.T) {
	baseSchema := validators.Object(map[string]zod.Schema{
		"name": validators.String().Required(),
		"age":  validators.Number().Required(),
	})

	t.Run("AddField", func(t *testing.T) {
		schema := baseSchema.AddField("email", validators.String().Email().Required())

		err := schema.Validate(map[string]interface{}{
			"name":  "John",
			"age":   30,
			"email": "john@example.com",
		})
		if err != nil {
			t.Errorf("AddField validation failed: %v", err)
		}

		// Should fail without email
		err = schema.Validate(map[string]interface{}{
			"name": "John",
			"age":  30,
		})
		if err == nil {
			t.Error("AddField should require the new field")
		}
	})

	t.Run("RemoveField", func(t *testing.T) {
		schema := baseSchema.RemoveField("age")

		err := schema.Validate(map[string]interface{}{
			"name": "John",
			// age is removed, so it's not required
		})
		if err != nil {
			t.Errorf("RemoveField validation failed: %v", err)
		}
	})
}

func TestObjectValidator_SchemaComposition(t *testing.T) {
	baseSchema := validators.Object(map[string]zod.Schema{
		"name": validators.String().Required(),
		"age":  validators.Number().Required(),
	})

	t.Run("Extend", func(t *testing.T) {
		extendedSchema := baseSchema.Extend(map[string]zod.Schema{
			"email": validators.String().Email().Required(),
			"phone": validators.String().Optional(),
		})

		err := extendedSchema.Validate(map[string]interface{}{
			"name":  "John",
			"age":   30,
			"email": "john@example.com",
		})
		if err != nil {
			t.Errorf("Extend validation failed: %v", err)
		}

		// Should fail without email (new required field)
		err = extendedSchema.Validate(map[string]interface{}{
			"name": "John",
			"age":  30,
		})
		if err == nil {
			t.Error("Extend should require new required fields")
		}
	})

	t.Run("Pick", func(t *testing.T) {
		pickedSchema := baseSchema.Pick("name")

		err := pickedSchema.Validate(map[string]interface{}{
			"name": "John",
			// age is not required in picked schema
		})
		if err != nil {
			t.Errorf("Pick validation failed: %v", err)
		}
	})

	t.Run("Omit", func(t *testing.T) {
		omittedSchema := baseSchema.Omit("age")

		err := omittedSchema.Validate(map[string]interface{}{
			"name": "John",
			// age is omitted, so it's not required
		})
		if err != nil {
			t.Errorf("Omit validation failed: %v", err)
		}
	})
}

func TestObjectValidator_CustomMessages(t *testing.T) {
	schema := validators.Object(map[string]zod.Schema{
		"name": validators.String().Required(),
	}).Required().WithMessage("required", "Custom object required error").
		WithMessage("fields", "Custom fields error")

	tests := []struct {
		name        string
		input       interface{}
		expectedMsg string
	}{
		{"custom required message", nil, "Custom object required error"},
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

func TestObjectValidator_TypeConversion(t *testing.T) {
	schema := validators.Object(map[string]zod.Schema{
		"name": validators.String().Required(),
	})

	// Test different map types
	tests := []struct {
		name      string
		input     interface{}
		wantError bool
	}{
		{
			"map[string]interface{}",
			map[string]interface{}{"name": "John"},
			false,
		},
		{
			"map[string]string",
			map[string]string{"name": "John"},
			false,
		},
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

func TestObjectValidator_HelperMethods(t *testing.T) {
	t.Run("StrictObject", func(t *testing.T) {
		schema := validators.StrictObject(map[string]zod.Schema{
			"name": validators.String().Required(),
		})

		err := schema.Validate(map[string]interface{}{
			"name":    "John",
			"unknown": "field",
		})
		if err == nil {
			t.Error("StrictObject should reject unknown fields")
		}
	})

	t.Run("OptionalObject", func(t *testing.T) {
		schema := validators.OptionalObject(map[string]zod.Schema{
			"name": validators.String().Required(),
		})

		err := schema.Validate(nil)
		if err != nil {
			t.Errorf("OptionalObject should allow nil: %v", err)
		}
	})

	t.Run("UserSchema", func(t *testing.T) {
		schema := validators.UserSchema()

		user := map[string]interface{}{
			"id":    1,
			"name":  "John Doe",
			"email": "john@example.com",
			"age":   30,
		}

		err := schema.Validate(user)
		if err != nil {
			t.Errorf("UserSchema validation failed: %v", err)
		}

		// Test invalid email
		invalidUser := map[string]interface{}{
			"id":    1,
			"name":  "John Doe",
			"email": "invalid-email",
		}

		err = schema.Validate(invalidUser)
		if err == nil {
			t.Error("UserSchema should reject invalid email")
		}
	})

	t.Run("AddressSchema", func(t *testing.T) {
		schema := validators.AddressSchema()

		address := map[string]interface{}{
			"street":  "123 Main St",
			"city":    "New York",
			"state":   "NY",
			"zipCode": "10001",
		}

		err := schema.Validate(address)
		if err != nil {
			t.Errorf("AddressSchema validation failed: %v", err)
		}

		// Test with default country
		minimalAddress := map[string]interface{}{
			"street": "123 Main St",
			"city":   "New York",
		}

		err = schema.Validate(minimalAddress)
		if err != nil {
			t.Errorf("AddressSchema with defaults failed: %v", err)
		}
	})
}

func TestObjectValidator_DefaultValues(t *testing.T) {
	defaultObj := map[string]interface{}{
		"name": "Default Name",
		"age":  0,
	}

	schema := validators.Object(map[string]zod.Schema{
		"name": validators.String().Required(),
		"age":  validators.Number().Required(),
	}).Default(defaultObj)

	// Test with nil - should use default
	err := schema.Validate(nil)
	if err != nil {
		t.Errorf("Expected no error for nil with default, got %v", err)
	}
}

func TestObjectValidator_ComplexNesting(t *testing.T) {
	// Complex nested structure: User with multiple addresses and preferences
	addressSchema := validators.Object(map[string]zod.Schema{
		"type":   validators.String().Required(),
		"street": validators.String().Required(),
		"city":   validators.String().Required(),
	})

	preferencesSchema := validators.Object(map[string]zod.Schema{
		"theme":         validators.String().Required(),
		"notifications": validators.Bool().Required(),
		"language":      validators.String().Default("en"),
	})

	userSchema := validators.Object(map[string]zod.Schema{
		"id":          validators.Number().Integer().Positive().Required(),
		"name":        validators.String().Min(1).Required(),
		"email":       validators.String().Email().Required(),
		"addresses":   validators.Array(addressSchema).Min(1).Required(),
		"preferences": preferencesSchema.Required(),
		"tags":        validators.Array(validators.String().Required()).Unique().Optional(),
	}).Strict()

	validUser := map[string]interface{}{
		"id":    1,
		"name":  "John Doe",
		"email": "john@example.com",
		"addresses": []interface{}{
			map[string]interface{}{
				"type":   "home",
				"street": "123 Main St",
				"city":   "New York",
			},
			map[string]interface{}{
				"type":   "work",
				"street": "456 Business Ave",
				"city":   "New York",
			},
		},
		"preferences": map[string]interface{}{
			"theme":         "dark",
			"notifications": true,
		},
		"tags": []interface{}{"developer", "golang", "typescript"},
	}

	err := userSchema.Validate(validUser)
	if err != nil {
		t.Errorf("Complex nested validation failed: %v", err)
	}

	// Test with invalid nested data
	invalidUser := map[string]interface{}{
		"id":    1,
		"name":  "John Doe",
		"email": "invalid-email", // Invalid email
		"addresses": []interface{}{
			map[string]interface{}{
				"type": "home",
				// missing street and city
			},
		},
		"preferences": map[string]interface{}{
			"theme": "dark",
			// missing notifications
		},
		"tags": []interface{}{"developer", "developer"}, // Duplicate tags
	}

	err = userSchema.Validate(invalidUser)
	if err == nil {
		t.Error("Complex nested validation should have failed")
	}
}
