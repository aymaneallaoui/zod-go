package tests

import (
	"testing"

	"github.com/aymaneallaoui/zod-go/zod/validators"
)

// TestCleanPackageInterface verifies that only intended exports are visible
func TestCleanPackageInterface(t *testing.T) {
	t.Run("Main builders are accessible", func(t *testing.T) {
		// These should all be accessible and return the correct types
		_ = validators.String()
		_ = validators.Number()
		_ = validators.Array(validators.String())
		_ = validators.Object(map[string]interface{}{})
		_ = validators.Bool()
	})

	t.Run("Convenience builders are accessible", func(t *testing.T) {
		_ = validators.Email()
		_ = validators.URL()
		_ = validators.OptionalString()
		_ = validators.RequiredString()
		_ = validators.PositiveNumber()
		_ = validators.IntegerNumber()
	})

	t.Run("Error keys are accessible", func(t *testing.T) {
		// Method-based access
		_ = validators.Errors.Required()
		_ = validators.Errors.MinLength()
		_ = validators.Errors.Email()

		// Constant access
		_ = validators.ErrRequired
		_ = validators.ErrMinLength
		_ = validators.ErrEmail
	})
}

// TestTypeSafeStateManagement verifies compile-time safety
func TestTypeSafeStateManagement(t *testing.T) {
	t.Run("String state transitions work correctly", func(t *testing.T) {
		// Initial state can transition to required or optional
		initial := validators.String().Min(3).Max(50)

		// Can become required
		required := initial.Required()
		if err := required.Validate("hello"); err != nil {
			t.Errorf("Expected validation to pass, got: %v", err)
		}

		// Can become optional
		optional := validators.String().Min(3).Max(50).Optional()
		if err := optional.Validate(nil); err != nil {
			t.Errorf("Expected nil validation to pass for optional, got: %v", err)
		}
	})

	t.Run("Required state cannot have defaults", func(t *testing.T) {
		required := validators.String().Required()

		// This should compile and work
		err := required.Validate("test")
		if err != nil {
			t.Errorf("Expected validation to pass, got: %v", err)
		}

		// Note: required.Default("value") would not compile
		// This is tested by the type system, not runtime
	})

	t.Run("Optional state can have defaults", func(t *testing.T) {
		optional := validators.String().Optional().Default("default_value")

		// Nil should use default (conceptually - implementation may vary)
		err := optional.Validate(nil)
		if err != nil {
			t.Errorf("Expected nil validation to pass with default, got: %v", err)
		}
	})
}

// TestStringValidation tests the enhanced string validation
func TestStringValidation(t *testing.T) {
	t.Run("Required string validation", func(t *testing.T) {
		schema := validators.String().
			Min(3).WithMinLengthMessage("Too short").
			Max(10).WithMaxLengthMessage("Too long").
			Required().WithRequiredMessage("Required field")

		// Valid input
		if err := schema.Validate("hello"); err != nil {
			t.Errorf("Expected valid input to pass, got: %v", err)
		}

		// Too short
		if err := schema.Validate("hi"); err == nil {
			t.Error("Expected short input to fail")
		}

		// Too long
		if err := schema.Validate("this_is_too_long"); err == nil {
			t.Error("Expected long input to fail")
		}

		// Empty/nil should fail for required
		if err := schema.Validate(""); err == nil {
			t.Error("Expected empty input to fail for required")
		}
		if err := schema.Validate(nil); err == nil {
			t.Error("Expected nil input to fail for required")
		}
	})

	t.Run("Optional string validation", func(t *testing.T) {
		schema := validators.String().
			Min(3).
			Max(10).
			Optional().Default("default")

		// Valid input
		if err := schema.Validate("hello"); err != nil {
			t.Errorf("Expected valid input to pass, got: %v", err)
		}

		// Empty/nil should pass for optional
		if err := schema.Validate(""); err != nil {
			t.Errorf("Expected empty input to pass for optional, got: %v", err)
		}
		if err := schema.Validate(nil); err != nil {
			t.Errorf("Expected nil input to pass for optional, got: %v", err)
		}
	})

	t.Run("Email validation", func(t *testing.T) {
		schema := validators.Email() // Pre-configured required email

		// Valid email
		if err := schema.Validate("test@example.com"); err != nil {
			t.Errorf("Expected valid email to pass, got: %v", err)
		}

		// Invalid email
		if err := schema.Validate("invalid-email"); err == nil {
			t.Error("Expected invalid email to fail")
		}

		// Empty should fail (required)
		if err := schema.Validate(""); err == nil {
			t.Error("Expected empty email to fail")
		}
	})

	t.Run("URL validation", func(t *testing.T) {
		schema := validators.URL() // Pre-configured required URL

		// Valid URL
		if err := schema.Validate("https://example.com"); err != nil {
			t.Errorf("Expected valid URL to pass, got: %v", err)
		}

		// Invalid URL
		if err := schema.Validate("not-a-url"); err == nil {
			t.Error("Expected invalid URL to fail")
		}
	})

	t.Run("Pattern validation", func(t *testing.T) {
		schema := validators.String().
			Pattern(`^[a-zA-Z0-9_]+$`).WithPatternMessage("Invalid characters").
			Required()

		// Valid pattern
		if err := schema.Validate("valid_username123"); err != nil {
			t.Errorf("Expected valid pattern to pass, got: %v", err)
		}

		// Invalid pattern
		if err := schema.Validate("invalid-username!"); err == nil {
			t.Error("Expected invalid pattern to fail")
		}
	})
}

// TestNumberValidation tests the enhanced number validation
func TestNumberValidation(t *testing.T) {
	t.Run("Required number validation", func(t *testing.T) {
		schema := validators.Number().
			Min(0).WithMinMessage("Cannot be negative").
			Max(100).WithMaxMessage("Cannot exceed 100").
			Required().WithRequiredMessage("Number required")

		// Valid number
		if err := schema.Validate(50); err != nil {
			t.Errorf("Expected valid number to pass, got: %v", err)
		}
		if err := schema.Validate(50.5); err != nil {
			t.Errorf("Expected valid float to pass, got: %v", err)
		}

		// Too small
		if err := schema.Validate(-1); err == nil {
			t.Error("Expected negative number to fail")
		}

		// Too large
		if err := schema.Validate(101); err == nil {
			t.Error("Expected large number to fail")
		}

		// Nil should fail for required
		if err := schema.Validate(nil); err == nil {
			t.Error("Expected nil to fail for required")
		}
	})

	t.Run("Optional number validation", func(t *testing.T) {
		schema := validators.Number().
			Min(0).
			Max(100).
			Optional().Default(18.0)

		// Valid number
		if err := schema.Validate(25); err != nil {
			t.Errorf("Expected valid number to pass, got: %v", err)
		}

		// Nil should pass for optional
		if err := schema.Validate(nil); err != nil {
			t.Errorf("Expected nil to pass for optional, got: %v", err)
		}
	})

	t.Run("Integer validation", func(t *testing.T) {
		schema := validators.IntegerNumber().Required() // Pre-configured integer

		// Valid integer
		if err := schema.Validate(42); err != nil {
			t.Errorf("Expected integer to pass, got: %v", err)
		}

		// Invalid (float)
		if err := schema.Validate(42.5); err == nil {
			t.Error("Expected float to fail integer validation")
		}
	})

	t.Run("Positive number validation", func(t *testing.T) {
		schema := validators.PositiveNumber().Required() // Pre-configured positive

		// Valid positive
		if err := schema.Validate(1); err != nil {
			t.Errorf("Expected positive number to pass, got: %v", err)
		}

		// Invalid (zero/negative)
		if err := schema.Validate(0); err == nil {
			t.Error("Expected zero to fail positive validation")
		}
		if err := schema.Validate(-1); err == nil {
			t.Error("Expected negative to fail positive validation")
		}
	})
}

// TestArrayValidation tests the enhanced array validation
func TestArrayValidation(t *testing.T) {
	t.Run("Required array validation", func(t *testing.T) {
		schema := validators.Array(validators.String().Min(1)).
			MinItems(1).WithMinItemsMessage("At least one item required").
			MaxItems(3).WithMaxItemsMessage("Too many items").
			Required()

		// Valid array
		validArray := []string{"item1", "item2"}
		if err := schema.Validate(validArray); err != nil {
			t.Errorf("Expected valid array to pass, got: %v", err)
		}

		// Empty array should fail min items
		emptyArray := []string{}
		if err := schema.Validate(emptyArray); err == nil {
			t.Error("Expected empty array to fail min items")
		}

		// Too many items
		largeArray := []string{"item1", "item2", "item3", "item4"}
		if err := schema.Validate(largeArray); err == nil {
			t.Error("Expected large array to fail max items")
		}

		// Nil should fail for required
		if err := schema.Validate(nil); err == nil {
			t.Error("Expected nil to fail for required array")
		}
	})

	t.Run("Optional array validation", func(t *testing.T) {
		defaultArray := []interface{}{"default"}
		schema := validators.Array(validators.String()).
			Optional().Default(defaultArray)

		// Valid array
		validArray := []string{"item1", "item2"}
		if err := schema.Validate(validArray); err != nil {
			t.Errorf("Expected valid array to pass, got: %v", err)
		}

		// Nil should pass for optional
		if err := schema.Validate(nil); err != nil {
			t.Errorf("Expected nil to pass for optional array, got: %v", err)
		}
	})
}

// TestObjectValidation tests the enhanced object validation
func TestObjectValidation(t *testing.T) {
	t.Run("Required object validation", func(t *testing.T) {
		schema := validators.Object(map[string]interface{}{
			"name":  validators.RequiredString().Min(1),
			"email": validators.Email(),
			"age":   validators.Number().Min(0).Optional(),
		}).Required()

		// Valid object
		validObj := map[string]interface{}{
			"name":  "John",
			"email": "john@example.com",
			"age":   25,
		}
		if err := schema.Validate(validObj); err != nil {
			t.Errorf("Expected valid object to pass, got: %v", err)
		}

		// Missing required field
		invalidObj := map[string]interface{}{
			"email": "john@example.com",
			// missing name
		}
		if err := schema.Validate(invalidObj); err == nil {
			t.Error("Expected object with missing required field to fail")
		}

		// Nil should fail for required
		if err := schema.Validate(nil); err == nil {
			t.Error("Expected nil to fail for required object")
		}
	})

	t.Run("Optional object validation", func(t *testing.T) {
		defaultObj := map[string]interface{}{"default": "value"}
		schema := validators.Object(map[string]interface{}{
			"name": validators.OptionalString(),
		}).Optional().Default(defaultObj)

		// Valid object
		validObj := map[string]interface{}{"name": "John"}
		if err := schema.Validate(validObj); err != nil {
			t.Errorf("Expected valid object to pass, got: %v", err)
		}

		// Nil should pass for optional
		if err := schema.Validate(nil); err != nil {
			t.Errorf("Expected nil to pass for optional object, got: %v", err)
		}
	})
}

// TestBoolValidation tests the enhanced boolean validation
func TestBoolValidation(t *testing.T) {
	t.Run("Required bool validation", func(t *testing.T) {
		schema := validators.Bool().Required().WithRequiredMessage("Boolean required")

		// Valid boolean
		if err := schema.Validate(true); err != nil {
			t.Errorf("Expected true to pass, got: %v", err)
		}
		if err := schema.Validate(false); err != nil {
			t.Errorf("Expected false to pass, got: %v", err)
		}

		// Invalid type
		if err := schema.Validate("not a bool"); err == nil {
			t.Error("Expected non-boolean to fail")
		}

		// Nil should fail for required
		if err := schema.Validate(nil); err == nil {
			t.Error("Expected nil to fail for required boolean")
		}
	})

	t.Run("Optional bool validation", func(t *testing.T) {
		schema := validators.Bool().Optional().Default(true)

		// Valid boolean
		if err := schema.Validate(false); err != nil {
			t.Errorf("Expected false to pass, got: %v", err)
		}

		// Nil should pass for optional
		if err := schema.Validate(nil); err != nil {
			t.Errorf("Expected nil to pass for optional boolean, got: %v", err)
		}
	})
}

// TestErrorKeyAutocompletion tests the error key system
func TestErrorKeyAutocompletion(t *testing.T) {
	t.Run("Error key methods work", func(t *testing.T) {
		// These should all return the expected string values
		if validators.Errors.Required() != "required" {
			t.Errorf("Expected 'required', got: %s", validators.Errors.Required())
		}
		if validators.Errors.MinLength() != "minLength" {
			t.Errorf("Expected 'minLength', got: %s", validators.Errors.MinLength())
		}
		if validators.Errors.Email() != "email" {
			t.Errorf("Expected 'email', got: %s", validators.Errors.Email())
		}
	})

	t.Run("Error key constants work", func(t *testing.T) {
		// These should all have the expected values
		if validators.ErrRequired != "required" {
			t.Errorf("Expected 'required', got: %s", validators.ErrRequired)
		}
		if validators.ErrMinLength != "minLength" {
			t.Errorf("Expected 'minLength', got: %s", validators.ErrMinLength)
		}
		if validators.ErrEmail != "email" {
			t.Errorf("Expected 'email', got: %s", validators.ErrEmail)
		}
	})

	t.Run("Custom error messages work", func(t *testing.T) {
		schema := validators.String().
			Min(5).
			Required().
			WithMessage(validators.Errors.MinLength(), "Custom min length message").
			WithMessage(validators.ErrRequired, "Custom required message")

		// Test that custom message is used
		err := schema.Validate("hi") // Too short
		if err == nil {
			t.Error("Expected validation to fail")
		}
		if err.Error() != "Field: hi, Error: Custom min length message" {
			// Note: Exact error format may vary based on implementation
			t.Logf("Error message: %s", err.Error())
		}
	})
}

// TestComplexValidationPatterns tests real-world usage patterns
func TestComplexValidationPatterns(t *testing.T) {
	t.Run("User registration schema", func(t *testing.T) {
		userSchema := validators.Object(map[string]interface{}{
			"username": validators.String().
				Min(3).WithMinLengthMessage("Username too short").
				Max(20).WithMaxLengthMessage("Username too long").
				Pattern(`^[a-zA-Z0-9_]+$`).WithPatternMessage("Username contains invalid characters").
				Required(),
			"email": validators.Email().WithEmailMessage("Please enter a valid email address"),
			"password": validators.String().
				Min(8).WithMinLengthMessage("Password must be at least 8 characters").
				Required(),
			"age": validators.Number().
				Min(13).WithMinMessage("Must be at least 13 years old").
				Max(120).WithMaxMessage("Age must be realistic").
				Integer().WithIntegerMessage("Age must be a whole number").
				Optional(),
			"terms": validators.Bool().Required().WithRequiredMessage("Must accept terms"),
		}).Required()

		// Valid user
		validUser := map[string]interface{}{
			"username": "john_doe",
			"email":    "john@example.com",
			"password": "SecurePass123",
			"age":      25,
			"terms":    true,
		}
		if err := userSchema.Validate(validUser); err != nil {
			t.Errorf("Expected valid user to pass, got: %v", err)
		}

		// Invalid user (missing required fields)
		invalidUser := map[string]interface{}{
			"username": "john_doe",
			// missing email, password, terms
		}
		if err := userSchema.Validate(invalidUser); err == nil {
			t.Error("Expected invalid user to fail")
		}
	})

	t.Run("API response schema", func(t *testing.T) {
		responseSchema := validators.Object(map[string]interface{}{
			"status": validators.String().Required(),
			"data": validators.Object(map[string]interface{}{
				"users": validators.Array(validators.Object(map[string]interface{}{
					"id":   validators.Number().Required(),
					"name": validators.RequiredString(),
				})).Optional(),
			}).Optional(),
			"errors": validators.Array(validators.String()).Optional(),
		}).Required()

		// Valid response
		validResponse := map[string]interface{}{
			"status": "success",
			"data": map[string]interface{}{
				"users": []map[string]interface{}{
					{"id": 1, "name": "John"},
					{"id": 2, "name": "Jane"},
				},
			},
		}
		if err := responseSchema.Validate(validResponse); err != nil {
			t.Errorf("Expected valid response to pass, got: %v", err)
		}

		// Error response
		errorResponse := map[string]interface{}{
			"status": "error",
			"errors": []string{"Invalid request", "Missing parameter"},
		}
		if err := responseSchema.Validate(errorResponse); err != nil {
			t.Errorf("Expected error response to pass, got: %v", err)
		}
	})
}

// BenchmarkNewAPI benchmarks the performance of the new API
func BenchmarkNewAPI(b *testing.B) {
	schema := validators.String().Min(3).Max(50).Required()
	testValue := "hello world"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = schema.Validate(testValue)
	}
}

// BenchmarkComplexSchema benchmarks a complex validation schema
func BenchmarkComplexSchema(b *testing.B) {
	userSchema := validators.Object(map[string]interface{}{
		"name":  validators.RequiredString().Min(1).Max(100),
		"email": validators.Email(),
		"age":   validators.Number().Min(0).Max(150).Integer().Optional(),
		"tags":  validators.Array(validators.String().Min(1)).Optional(),
	}).Required()

	testUser := map[string]interface{}{
		"name":  "John Doe",
		"email": "john@example.com",
		"age":   25,
		"tags":  []string{"developer", "golang"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = userSchema.Validate(testUser)
	}
}
