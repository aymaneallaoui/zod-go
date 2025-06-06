package validators

import (
	"fmt"
	"testing"

	"github.com/aymaneallaoui/zod-go/zod/validators"
)

// BenchmarkNewVsOldAPI compares performance between new and legacy APIs
func BenchmarkNewVsOldAPI(b *testing.B) {
	testValue := "test_user_123"

	b.Run("NewAPI_String_Simple", func(b *testing.B) {
		schema := validators.String().Min(3).Max(50).Required()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = schema.Validate(testValue)
		}
	})

	b.Run("NewAPI_String_Complex", func(b *testing.B) {
		schema := validators.String().
			Min(3).WithMinLengthMessage("Too short").
			Max(50).WithMaxLengthMessage("Too long").
			Pattern(`^[a-zA-Z0-9_]+$`).WithPatternMessage("Invalid format").
			Required().WithRequiredMessage("Username required")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = schema.Validate(testValue)
		}
	})

	b.Run("NewAPI_Email_Validation", func(b *testing.B) {
		schema := validators.Email().WithEmailMessage("Invalid email")
		testEmail := "user@example.com"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = schema.Validate(testEmail)
		}
	})

	b.Run("NewAPI_Number_Validation", func(b *testing.B) {
		schema := validators.Number().Min(0).Max(100).Integer().Required()
		testNum := 42
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = schema.Validate(testNum)
		}
	})
}

// BenchmarkComplexSchemas tests performance with realistic complex schemas
func BenchmarkComplexSchemas(b *testing.B) {
	userSchema := validators.Object(map[string]interface{}{
		"username": validators.String().
			Min(3).Max(20).
			Pattern(`^[a-zA-Z0-9_]+$`).
			Required(),
		"email": validators.Email(),
		"password": validators.String().
			Min(8).
			Pattern(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]`).
			Required(),
		"age":       validators.Number().Min(13).Max(120).Integer().Optional(),
		"interests": validators.Array(validators.String().Min(1)).MinItems(0).MaxItems(10).Optional(),
		"settings": validators.Object(map[string]interface{}{
			"theme":         validators.String().Optional().Default("light"),
			"notifications": validators.Bool().Optional().Default(true),
			"language":      validators.String().Optional().Default("en"),
		}).Optional(),
	}).Required()

	validUser := map[string]interface{}{
		"username":  "john_doe",
		"email":     "john@example.com",
		"password":  "SecurePass123!",
		"age":       25,
		"interests": []string{"programming", "music", "travel"},
		"settings": map[string]interface{}{
			"theme":         "dark",
			"notifications": true,
			"language":      "en",
		},
	}

	b.Run("ComplexUserSchema_Valid", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = userSchema.Validate(validUser)
		}
	})

	invalidUser := map[string]interface{}{
		"username": "jo", // Too short
		"email":    "invalid-email",
		"password": "weak",
		"age":      5, // Too young
	}

	b.Run("ComplexUserSchema_Invalid", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = userSchema.Validate(invalidUser)
		}
	})
}

// BenchmarkArrayValidation tests array validation performance
func BenchmarkArrayValidation(b *testing.B) {
	schema := validators.Array(validators.String().Min(1).Max(100)).
		MinItems(1).MaxItems(1000).
		Required()

	b.Run("SmallArray_10_Items", func(b *testing.B) {
		testArray := make([]string, 10)
		for i := range testArray {
			testArray[i] = fmt.Sprintf("item_%d", i)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = schema.Validate(testArray)
		}
	})

	b.Run("MediumArray_100_Items", func(b *testing.B) {
		testArray := make([]string, 100)
		for i := range testArray {
			testArray[i] = fmt.Sprintf("item_%d", i)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = schema.Validate(testArray)
		}
	})

	b.Run("LargeArray_1000_Items", func(b *testing.B) {
		testArray := make([]string, 1000)
		for i := range testArray {
			testArray[i] = fmt.Sprintf("item_%d", i)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = schema.Validate(testArray)
		}
	})
}

// BenchmarkStateTransitions tests the performance impact of type-safe state management
func BenchmarkStateTransitions(b *testing.B) {
	b.Run("StringBuilder_To_Required", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			schema := validators.String().Min(3).Max(50).Required()
			_ = schema.Validate("test")
		}
	})

	b.Run("StringBuilder_To_Optional", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			schema := validators.String().Min(3).Max(50).Optional().Default("default")
			_ = schema.Validate(nil)
		}
	})

	b.Run("NumberBuilder_To_Required", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			schema := validators.Number().Min(0).Max(100).Integer().Required()
			_ = schema.Validate(42)
		}
	})

	b.Run("NumberBuilder_To_Optional", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			schema := validators.Number().Min(0).Max(100).Optional().Default(0.0)
			_ = schema.Validate(nil)
		}
	})
}

// BenchmarkErrorKeyPerformance tests different error key approaches
func BenchmarkErrorKeyPerformance(b *testing.B) {
	testValue := "hi"

	b.Run("ErrorKeys_Method_Based", func(b *testing.B) {
		schema := validators.String().Min(5).Required().
			WithMessage(validators.Errors.MinLength(), "Too short").
			WithMessage(validators.Errors.Required(), "Required")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = schema.Validate(testValue)
		}
	})

	b.Run("ErrorKeys_Constants", func(b *testing.B) {
		schema := validators.String().Min(5).Required().
			WithMessage(validators.ErrMinLength, "Too short").
			WithMessage(validators.ErrRequired, "Required")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = schema.Validate(testValue)
		}
	})

	b.Run("ErrorKeys_Convenience_Methods", func(b *testing.B) {
		schema := validators.String().Min(5).WithMinLengthMessage("Too short").
			Required().WithRequiredMessage("Required")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = schema.Validate(testValue)
		}
	})

	b.Run("ErrorKeys_Legacy_Strings", func(b *testing.B) {
		schema := validators.String().Min(5).Required().
			WithMessage("minLength", "Too short").
			WithMessage("required", "Required")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = schema.Validate(testValue)
		}
	})
}

// BenchmarkMemoryAllocations tests memory efficiency
func BenchmarkMemoryAllocations(b *testing.B) {
	b.Run("String_Schema_Creation", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = validators.String().Min(3).Max(50).Required()
		}
	})

	b.Run("Complex_Object_Schema_Creation", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = validators.Object(map[string]interface{}{
				"name":  validators.String().Min(1).Required(),
				"email": validators.Email(),
				"age":   validators.Number().Min(0).Optional(),
			}).Required()
		}
	})

	b.Run("Array_Schema_Creation", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = validators.Array(validators.String().Min(1)).MinItems(1).MaxItems(10).Required()
		}
	})
}

// BenchmarkPreConfiguredValidators tests convenience validators
func BenchmarkPreConfiguredValidators(b *testing.B) {
	emailValue := "test@example.com"
	urlValue := "https://example.com"
	stringValue := "test"

	b.Run("PreConfigured_Email", func(b *testing.B) {
		schema := validators.Email()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = schema.Validate(emailValue)
		}
	})

	b.Run("PreConfigured_URL", func(b *testing.B) {
		schema := validators.URL()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = schema.Validate(urlValue)
		}
	})

	b.Run("PreConfigured_RequiredString", func(b *testing.B) {
		schema := validators.RequiredString()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = schema.Validate(stringValue)
		}
	})

	b.Run("PreConfigured_OptionalString", func(b *testing.B) {
		schema := validators.OptionalString()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = schema.Validate(nil)
		}
	})

	b.Run("PreConfigured_PositiveNumber", func(b *testing.B) {
		schema := validators.PositiveNumber().Required()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = schema.Validate(42)
		}
	})

	b.Run("PreConfigured_IntegerNumber", func(b *testing.B) {
		schema := validators.IntegerNumber().Required()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = schema.Validate(42)
		}
	})
}

// BenchmarkValidationFailures tests performance when validation fails
func BenchmarkValidationFailures(b *testing.B) {
	b.Run("String_MinLength_Failure", func(b *testing.B) {
		schema := validators.String().Min(10).Required()
		testValue := "short"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = schema.Validate(testValue)
		}
	})

	b.Run("String_MaxLength_Failure", func(b *testing.B) {
		schema := validators.String().Max(5).Required()
		testValue := "this_is_too_long"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = schema.Validate(testValue)
		}
	})

	b.Run("String_Pattern_Failure", func(b *testing.B) {
		schema := validators.String().Pattern(`^[a-zA-Z]+$`).Required()
		testValue := "invalid123"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = schema.Validate(testValue)
		}
	})

	b.Run("Email_Format_Failure", func(b *testing.B) {
		schema := validators.Email()
		testValue := "invalid-email"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = schema.Validate(testValue)
		}
	})

	b.Run("Number_Range_Failure", func(b *testing.B) {
		schema := validators.Number().Min(0).Max(100).Required()
		testValue := -5
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = schema.Validate(testValue)
		}
	})

	b.Run("Array_MinItems_Failure", func(b *testing.B) {
		schema := validators.Array(validators.String()).MinItems(5).Required()
		testValue := []string{"item1", "item2"}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = schema.Validate(testValue)
		}
	})
}

// BenchmarkConcurrentValidation tests thread safety and concurrent performance
func BenchmarkConcurrentValidation(b *testing.B) {
	schema := validators.String().Min(3).Max(50).Required()
	testValue := "concurrent_test"

	b.Run("Sequential_Validation", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = schema.Validate(testValue)
		}
	})

	b.Run("Concurrent_Validation", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = schema.Validate(testValue)
			}
		})
	})
}

// BenchmarkTypeConversions tests performance of type checking and conversion
func BenchmarkTypeConversions(b *testing.B) {
	numberSchema := validators.Number().Required()

	b.Run("Int_To_Float64", func(b *testing.B) {
		testValue := 42
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = numberSchema.Validate(testValue)
		}
	})

	b.Run("Float32_To_Float64", func(b *testing.B) {
		testValue := float32(42.5)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = numberSchema.Validate(testValue)
		}
	})

	b.Run("String_Type_Check", func(b *testing.B) {
		stringSchema := validators.String().Required()
		testValue := "test_string"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = stringSchema.Validate(testValue)
		}
	})

	b.Run("Array_Type_Check", func(b *testing.B) {
		arraySchema := validators.Array(validators.String()).Required()
		testValue := []string{"item1", "item2", "item3"}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = arraySchema.Validate(testValue)
		}
	})

	b.Run("Map_Type_Check", func(b *testing.B) {
		objectSchema := validators.Object(map[string]interface{}{
			"name": validators.String().Required(),
			"age":  validators.Number().Required(),
		}).Required()
		testValue := map[string]interface{}{
			"name": "John",
			"age":  30,
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = objectSchema.Validate(testValue)
		}
	})
}
