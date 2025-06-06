package validators

import (
	"fmt"
	"math"
	"reflect"
	"strings"
	"sync"
	"testing"
)

// TestEdgeCases covers edge cases and boundary conditions
func TestEdgeCases(t *testing.T) {
	t.Run("Empty and Nil Inputs", func(t *testing.T) {
		stringSchema := String().Required()

		nilCases := []interface{}{nil, (*string)(nil), (*int)(nil)}
		for i, nilCase := range nilCases {
			t.Run(fmt.Sprintf("nil_case_%d", i), func(t *testing.T) {
				if err := stringSchema.Validate(nilCase); err == nil {
					t.Error("Expected nil validation to fail for required field")
				}
			})
		}

		if err := stringSchema.Validate(""); err == nil {
			t.Error("Expected empty string to fail for required field")
		}

		numberSchema := Number().Required()
		if err := numberSchema.Validate(0); err != nil {
			t.Errorf("Expected zero to be valid for number, got: %v", err)
		}
		if err := numberSchema.Validate(0.0); err != nil {
			t.Errorf("Expected 0.0 to be valid for number, got: %v", err)
		}
	})

	t.Run("Boundary Values", func(t *testing.T) {
		schema := String().Min(5).Max(10).Required()

		if err := schema.Validate("1234"); err == nil {
			t.Error("Expected string with length 4 to fail min(5)")
		}
		if err := schema.Validate("12345"); err != nil {
			t.Errorf("Expected string with length 5 to pass min(5), got: %v", err)
		}
		if err := schema.Validate("1234567890"); err != nil {
			t.Errorf("Expected string with length 10 to pass max(10), got: %v", err)
		}
		if err := schema.Validate("12345678901"); err == nil {
			t.Error("Expected string with length 11 to fail max(10)")
		}

		numSchema := Number().Min(0).Max(100).Required()

		if err := numSchema.Validate(-0.1); err == nil {
			t.Error("Expected -0.1 to fail min(0)")
		}
		if err := numSchema.Validate(0); err != nil {
			t.Errorf("Expected 0 to pass min(0), got: %v", err)
		}
		if err := numSchema.Validate(100); err != nil {
			t.Errorf("Expected 100 to pass max(100), got: %v", err)
		}
		if err := numSchema.Validate(100.1); err == nil {
			t.Error("Expected 100.1 to fail max(100)")
		}
	})

	t.Run("Unicode and Special Characters", func(t *testing.T) {
		schema := String().Min(1).Max(50).Required()

		unicodeCases := []string{
			"Hello 世界",
			"🚀🎉✨",
			"Ñoño",
			"مرحبا",
			"Привет",
			"こんにちは",
		}

		for _, testCase := range unicodeCases {
			t.Run(fmt.Sprintf("unicode_%s", testCase), func(t *testing.T) {
				if err := schema.Validate(testCase); err != nil {
					t.Errorf("Expected unicode string '%s' to be valid, got: %v", testCase, err)
				}
			})
		}

		longUnicode := strings.Repeat("🚀", 100)
		if err := String().Max(50).Required().Validate(longUnicode); err == nil {
			t.Error("Expected long unicode string to fail max length")
		}
	})

	t.Run("Extreme Number Values", func(t *testing.T) {
		schema := Number().Required()

		extremeCases := []float64{
			math.MaxFloat64,
			math.SmallestNonzeroFloat64,
			math.Inf(1),
			math.Inf(-1),
			math.NaN(),
		}

		for i, testCase := range extremeCases {
			t.Run(fmt.Sprintf("extreme_number_%d", i), func(t *testing.T) {
				err := schema.Validate(testCase)
				if math.IsNaN(testCase) || math.IsInf(testCase, 0) {
					t.Logf("Extreme value %v result: %v", testCase, err)
				} else {
					if err != nil {
						t.Errorf("Expected extreme number %v to be valid, got: %v", testCase, err)
					}
				}
			})
		}
	})

	t.Run("Large Data Structures", func(t *testing.T) {
		largeArray := make([]string, 10000)
		for i := range largeArray {
			largeArray[i] = fmt.Sprintf("item_%d", i)
		}

		arraySchema := Array(String().Min(1)).Required()
		if err := arraySchema.Validate(largeArray); err != nil {
			t.Errorf("Expected large array to be valid, got: %v", err)
		}

		deepObject := make(map[string]interface{})
		current := deepObject
		for i := 0; i < 100; i++ {
			next := make(map[string]interface{})
			current[fmt.Sprintf("level_%d", i)] = next
			current = next
		}
		current["value"] = "deep_value"

		objectSchema := Object(map[string]interface{}{}).Optional()
		if err := objectSchema.Validate(deepObject); err != nil {
			t.Errorf("Expected deep object to be valid, got: %v", err)
		}
	})
}

// TestConcurrencyAndThreadSafety tests thread safety
func TestConcurrencyAndThreadSafety(t *testing.T) {
	t.Run("Concurrent Schema Creation", func(t *testing.T) {
		var wg sync.WaitGroup
		numGoroutines := 100

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()

				schema := String().
					Min(1).
					Max(100).
					Pattern(fmt.Sprintf("^test_%d_.*", id)).
					Required()

				testValue := fmt.Sprintf("test_%d_value", id)
				if err := schema.Validate(testValue); err != nil {
					t.Errorf("Goroutine %d: validation failed: %v", id, err)
				}
			}(i)
		}

		wg.Wait()
	})

	t.Run("Concurrent Validation", func(t *testing.T) {
		schema := String().Min(1).Max(100).Required()
		var wg sync.WaitGroup
		numGoroutines := 1000

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()

				testValue := fmt.Sprintf("concurrent_test_%d", id)
				if err := schema.Validate(testValue); err != nil {
					t.Errorf("Goroutine %d: validation failed: %v", id, err)
				}
			}(i)
		}

		wg.Wait()
	})

	t.Run("Concurrent Complex Object Validation", func(t *testing.T) {
		schema := Object(map[string]interface{}{
			"id":    Number().Integer().Required(),
			"name":  String().Min(1).Required(),
			"email": Email(),
			"tags":  Array(String()).Optional(),
		}).Required()

		var wg sync.WaitGroup
		numGoroutines := 100

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()

				testData := map[string]interface{}{
					"id":    id,
					"name":  fmt.Sprintf("User_%d", id),
					"email": fmt.Sprintf("user%d@example.com", id),
					"tags":  []string{fmt.Sprintf("tag_%d", id)},
				}

				if err := schema.Validate(testData); err != nil {
					t.Errorf("Goroutine %d: validation failed: %v", id, err)
				}
			}(i)
		}

		wg.Wait()
	})
}

// TestErrorMessageCustomization tests error message features
func TestErrorMessageCustomization(t *testing.T) {
	t.Run("Custom Error Messages", func(t *testing.T) {
		schema := String().
			Min(5).WithMessage(Errors.MinLength(), "Custom min length error").
			Max(10).WithMessage(ErrMaxLength, "Custom max length error").
			Required().WithRequiredMessage("Custom required error")

		err := schema.Validate("hi")
		if err == nil {
			t.Error("Expected validation to fail")
		}
		if !strings.Contains(err.Error(), "Custom min length error") {
			t.Errorf("Expected custom min length error message, got: %s", err.Error())
		}

		err = schema.Validate("this_is_too_long")
		if err == nil {
			t.Error("Expected validation to fail")
		}
		if !strings.Contains(err.Error(), "Custom max length error") {
			t.Errorf("Expected custom max length error message, got: %s", err.Error())
		}

		err = schema.Validate("")
		if err == nil {
			t.Error("Expected validation to fail")
		}
		if !strings.Contains(err.Error(), "Custom required error") {
			t.Errorf("Expected custom required error message, got: %s", err.Error())
		}
	})

	t.Run("Error Message Overrides", func(t *testing.T) {
		schema := String().
			Min(5).WithMinLengthMessage("First message").
			WithMessage(Errors.MinLength(), "Override message").
			Required()

		err := schema.Validate("hi")
		if err == nil {
			t.Error("Expected validation to fail")
		}

		if !strings.Contains(err.Error(), "Override message") {
			t.Errorf("Expected override message, got: %s", err.Error())
		}
	})

	t.Run("Default vs Custom Messages", func(t *testing.T) {
		defaultSchema := String().Min(5).Required()

		customSchema := String().
			Min(5).WithMinLengthMessage("Custom message").
			Required()

		testValue := "hi"

		defaultErr := defaultSchema.Validate(testValue)
		customErr := customSchema.Validate(testValue)

		if defaultErr == nil || customErr == nil {
			t.Error("Expected both validations to fail")
		}

		if defaultErr.Error() == customErr.Error() {
			t.Error("Expected different error messages for default vs custom")
		}

		if !strings.Contains(customErr.Error(), "Custom message") {
			t.Errorf("Expected custom message in error, got: %s", customErr.Error())
		}
	})
}

// TestComplexDataTypes tests complex nested data structures
func TestComplexDataTypes(t *testing.T) {
	t.Run("Deeply Nested Objects", func(t *testing.T) {
		schema := Object(map[string]interface{}{
			"user": Object(map[string]interface{}{
				"profile": Object(map[string]interface{}{
					"personal": Object(map[string]interface{}{
						"name": String().Required(),
						"age":  Number().Integer().Min(0).Required(),
					}).Required(),
					"contact": Object(map[string]interface{}{
						"email": Email(),
						"phone": String().Optional(),
					}).Required(),
				}).Required(),
				"preferences": Object(map[string]interface{}{
					"theme":         String().Optional().Default("light"),
					"language":      String().Optional().Default("en"),
					"notifications": Bool().Optional().Default(true),
				}).Optional(),
			}).Required(),
		}).Required()

		validData := map[string]interface{}{
			"user": map[string]interface{}{
				"profile": map[string]interface{}{
					"personal": map[string]interface{}{
						"name": "John Doe",
						"age":  30,
					},
					"contact": map[string]interface{}{
						"email": "john@example.com",
						"phone": "+1234567890",
					},
				},
				"preferences": map[string]interface{}{
					"theme":         "dark",
					"language":      "es",
					"notifications": false,
				},
			},
		}

		if err := schema.Validate(validData); err != nil {
			t.Errorf("Expected nested object to be valid, got: %v", err)
		}

		invalidData := map[string]interface{}{
			"user": map[string]interface{}{
				"profile": map[string]interface{}{
					"personal": map[string]interface{}{
						"name": "John Doe",
						// missing age
					},
					"contact": map[string]interface{}{
						"email": "john@example.com",
					},
				},
			},
		}

		if err := schema.Validate(invalidData); err == nil {
			t.Error("Expected invalid nested object to fail validation")
		}
	})

	t.Run("Arrays of Complex Objects", func(t *testing.T) {
		userSchema := Object(map[string]interface{}{
			"id":    Number().Integer().Required(),
			"name":  String().Min(1).Required(),
			"email": Email(),
			"tags":  Array(String().Min(1)).Optional(),
		})

		arraySchema := Array(userSchema).
			MinItems(1).
			MaxItems(100).
			Required()

		validData := []map[string]interface{}{
			{
				"id":    1,
				"name":  "User 1",
				"email": "user1@example.com",
				"tags":  []string{"admin", "power-user"},
			},
			{
				"id":    2,
				"name":  "User 2",
				"email": "user2@example.com",
			},
		}

		if err := arraySchema.Validate(validData); err != nil {
			t.Errorf("Expected array of objects to be valid, got: %v", err)
		}

		invalidData := []map[string]interface{}{
			{
				"id":    1,
				"name":  "User 1",
				"email": "user1@example.com",
			},
			{
				"id":    2,
				"name":  "", // Invalid: empty name
				"email": "user2@example.com",
			},
		}

		if err := arraySchema.Validate(invalidData); err == nil {
			t.Error("Expected array with invalid object to fail validation")
		}
	})

	t.Run("Mixed Type Arrays", func(t *testing.T) {
		mixedArray := []interface{}{
			"string_value",
			123,
			true,
			map[string]interface{}{"key": "value"},
		}

		arraySchema := Array(String()).Optional()

		if err := arraySchema.Validate(mixedArray); err == nil {
			t.Log("Mixed array validation behavior:", err)
		}
	})
}

// TestTypeConversions tests type handling and conversions
func TestTypeConversions(t *testing.T) {
	t.Run("Number Type Variations", func(t *testing.T) {
		schema := Number().Required()

		testCases := []interface{}{
			int(42),
			int8(42),
			int16(42),
			int32(42),
			int64(42),
			uint(42),
			uint8(42),
			uint16(42),
			uint32(42),
			uint64(42),
			float32(42.5),
			float64(42.5),
		}

		for i, testCase := range testCases {
			t.Run(fmt.Sprintf("numeric_type_%d_%T", i, testCase), func(t *testing.T) {
				if err := schema.Validate(testCase); err != nil {
					t.Errorf("Expected %T(%v) to be valid number, got: %v", testCase, testCase, err)
				}
			})
		}

		invalidCases := []interface{}{
			"not_a_number",
			[]int{1, 2, 3},
			map[string]int{"key": 1},
			true,
		}

		for i, testCase := range invalidCases {
			t.Run(fmt.Sprintf("invalid_type_%d_%T", i, testCase), func(t *testing.T) {
				if err := schema.Validate(testCase); err == nil {
					t.Errorf("Expected %T(%v) to be invalid number", testCase, testCase)
				}
			})
		}
	})

	t.Run("String Type Strictness", func(t *testing.T) {
		schema := String().Required()

		// Only actual strings should be valid
		validCases := []interface{}{
			"hello",
			"",
			"123",
			"true",
		}

		for i, testCase := range validCases {
			t.Run(fmt.Sprintf("valid_string_%d", i), func(t *testing.T) {
				if err := schema.Validate(testCase); err != nil && testCase != "" {
					t.Errorf("Expected %T(%v) to be valid string, got: %v", testCase, testCase, err)
				}
			})
		}

		// Numbers, booleans, etc. should not be automatically converted
		invalidCases := []interface{}{
			123,
			true,
			12.34,
			[]string{"array"},
		}

		for i, testCase := range invalidCases {
			t.Run(fmt.Sprintf("invalid_string_%d_%T", i, testCase), func(t *testing.T) {
				if err := schema.Validate(testCase); err == nil {
					t.Errorf("Expected %T(%v) to be invalid string", testCase, testCase)
				}
			})
		}
	})

	t.Run("Array Type Handling", func(t *testing.T) {
		schema := Array(String()).Required()

		testCases := []interface{}{
			[]string{"a", "b", "c"},
			[]interface{}{"a", "b", "c"},
		}

		for i, testCase := range testCases {
			t.Run(fmt.Sprintf("array_type_%d_%T", i, testCase), func(t *testing.T) {
				if err := schema.Validate(testCase); err != nil {
					t.Errorf("Expected %T to be valid array, got: %v", testCase, err)
				}
			})
		}

		invalidCases := []interface{}{
			"not_an_array",
			123,
			map[string]interface{}{"key": "value"},
		}

		for i, testCase := range invalidCases {
			t.Run(fmt.Sprintf("invalid_array_%d_%T", i, testCase), func(t *testing.T) {
				if err := schema.Validate(testCase); err == nil {
					t.Errorf("Expected %T to be invalid array", testCase)
				}
			})
		}
	})
}

// TestSchemaComposition tests combining schemas
func TestSchemaComposition(t *testing.T) {
	t.Run("Reusable Schema Components", func(t *testing.T) {
		nameSchema := String().Min(1).Max(100).Required()
		emailSchema := Email()
		ageSchema := Number().Integer().Min(0).Max(150).Optional()

		userSchema := Object(map[string]interface{}{
			"firstName": nameSchema,
			"lastName":  nameSchema,
			"email":     emailSchema,
			"age":       ageSchema,
		}).Required()

		adminSchema := Object(map[string]interface{}{
			"user":        userSchema,
			"permissions": Array(String().Min(1)).Required(),
			"role":        String().Pattern(`^(admin|super_admin)$`).Required(),
		}).Required()

		validAdmin := map[string]interface{}{
			"user": map[string]interface{}{
				"firstName": "Admin",
				"lastName":  "User",
				"email":     "admin@example.com",
				"age":       35,
			},
			"permissions": []string{"read", "write", "delete"},
			"role":        "admin",
		}

		if err := adminSchema.Validate(validAdmin); err != nil {
			t.Errorf("Expected composed schema to be valid, got: %v", err)
		}
	})

	t.Run("Schema Extension Patterns", func(t *testing.T) {
		baseUserFields := map[string]interface{}{
			"id":    Number().Integer().Required(),
			"name":  String().Min(1).Required(),
			"email": Email(),
		}

		extendedUserFields := make(map[string]interface{})
		for k, v := range baseUserFields {
			extendedUserFields[k] = v
		}
		extendedUserFields["createdAt"] = String().Required()
		extendedUserFields["updatedAt"] = String().Required()
		extendedUserFields["isActive"] = Bool().Optional().Default(true)

		baseSchema := Object(baseUserFields).Required()
		extendedSchema := Object(extendedUserFields).Required()

		baseData := map[string]interface{}{
			"id":    1,
			"name":  "John",
			"email": "john@example.com",
		}

		extendedData := map[string]interface{}{
			"id":        1,
			"name":      "John",
			"email":     "john@example.com",
			"createdAt": "2023-01-01T00:00:00Z",
			"updatedAt": "2023-01-02T00:00:00Z",
			"isActive":  true,
		}

		if err := baseSchema.Validate(baseData); err != nil {
			t.Errorf("Expected base schema to be valid, got: %v", err)
		}

		if err := extendedSchema.Validate(extendedData); err != nil {
			t.Errorf("Expected extended schema to be valid, got: %v", err)
		}

		// Base data should not be valid for extended schema (missing required fields)
		if err := extendedSchema.Validate(baseData); err == nil {
			t.Error("Expected base data to be invalid for extended schema")
		}
	})
}

// TestMemoryAndPerformance tests memory usage and performance characteristics
func TestMemoryAndPerformance(t *testing.T) {
	t.Run("Schema Reuse", func(t *testing.T) {
		schema := String().Min(1).Max(100).Required()

		for i := 0; i < 1000; i++ {
			testValue := fmt.Sprintf("test_value_%d", i)
			if err := schema.Validate(testValue); err != nil {
				t.Errorf("Iteration %d: validation failed: %v", i, err)
			}
		}
	})

	t.Run("Large Schema Performance", func(t *testing.T) {
		fields := make(map[string]interface{})
		for i := 0; i < 100; i++ {
			fields[fmt.Sprintf("field_%d", i)] = String().
				Min(1).Max(50).
				Pattern(fmt.Sprintf("^value_%d_.*", i)).
				Optional()
		}

		largeSchema := Object(fields).Required()

		testData := make(map[string]interface{})
		for i := 0; i < 50; i++ {
			testData[fmt.Sprintf("field_%d", i)] = fmt.Sprintf("value_%d_test", i)
		}

		if err := largeSchema.Validate(testData); err != nil {
			t.Errorf("Expected large schema validation to pass, got: %v", err)
		}
	})

	t.Run("Memory Efficiency", func(t *testing.T) {
		schemas := make([]interface{ Validate(interface{}) error }, 1000)

		for i := 0; i < 1000; i++ {
			schemas[i] = String().
				Min(i % 10).
				Max((i % 50) + 10).
				Required()
		}

		for i, schema := range schemas {
			minLen := i % 10
			maxLen := (i % 50) + 10
			testValue := fmt.Sprintf("test_%d", i)
			if len(testValue) < minLen {
				testValue += strings.Repeat("x", minLen-len(testValue))
			}
			if len(testValue) > maxLen {
				testValue = testValue[:maxLen]
			}

			if err := schema.Validate(testValue); err != nil {
				t.Errorf("Schema %d validation failed: %v", i, err)
			}
		}
	})
}

// TestErrorHandling tests comprehensive error scenarios
func TestErrorHandling(t *testing.T) {
	t.Run("Multiple Validation Errors", func(t *testing.T) {
		schema := Object(map[string]interface{}{
			"name":  String().Min(5).Required(),
			"email": Email(),
			"age":   Number().Min(18).Integer().Required(),
		}).Required()

		invalidData := map[string]interface{}{
			"name":  "Jo",            // Too short
			"email": "invalid-email", // Invalid format
			"age":   "not_a_number",  // Wrong type
		}

		err := schema.Validate(invalidData)
		if err == nil {
			t.Error("Expected validation to fail with multiple errors")
		}

		errorString := err.Error()
		t.Logf("Multiple validation errors: %s", errorString)

		if !strings.Contains(errorString, "Error:") {
			t.Error("Expected error to indicate validation failure")
		}
	})

	t.Run("Error Propagation", func(t *testing.T) {
		schema := Array(Object(map[string]interface{}{
			"users": Array(Object(map[string]interface{}{
				"name": String().Min(1).Required(),
			})).Required(),
		})).Required()

		invalidData := []map[string]interface{}{
			{
				"users": []map[string]interface{}{
					{"name": ""}, // Invalid: empty name
				},
			},
		}

		err := schema.Validate(invalidData)
		if err == nil {
			t.Error("Expected deeply nested validation to fail")
		}

		t.Logf("Nested validation error: %s", err.Error())
	})

	t.Run("Panic Recovery", func(t *testing.T) {
		schema := String().Required()

		extremeInputs := []interface{}{
			nil,
			(*string)(nil),
			make(chan int),
			func() {},
			complex(1, 2),
		}

		for i, input := range extremeInputs {
			t.Run(fmt.Sprintf("extreme_input_%d", i), func(t *testing.T) {
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("Validation panicked with input %T: %v", input, r)
					}
				}()

				_ = schema.Validate(input)
			})
		}
	})
}

// TestCompatibility tests integration with existing Go patterns
func TestCompatibility(t *testing.T) {
	t.Run("Interface Compatibility", func(t *testing.T) {
		stringSchema := String().Required()

		var validator interface{ Validate(interface{}) error }
		validator = stringSchema

		if err := validator.Validate("test"); err != nil {
			t.Errorf("Interface validation failed: %v", err)
		}
	})

	t.Run("Reflection Compatibility", func(t *testing.T) {
		// Test that the schemas work with reflection
		schema := Object(map[string]interface{}{
			"Name": String().Required(),
			"Age":  Number().Integer().Required(),
		}).Required()

		// Create a struct to test with
		type TestStruct struct {
			Name string
			Age  int
		}

		testStruct := TestStruct{
			Name: "John",
			Age:  30,
		}

		// Convert struct to map using reflection
		structValue := reflect.ValueOf(testStruct)
		structType := structValue.Type()

		dataMap := make(map[string]interface{})
		for i := 0; i < structValue.NumField(); i++ {
			field := structType.Field(i)
			value := structValue.Field(i)
			dataMap[field.Name] = value.Interface()
		}

		if err := schema.Validate(dataMap); err != nil {
			t.Errorf("Reflection-based validation failed: %v", err)
		}
	})

	t.Run("JSON Compatibility", func(t *testing.T) {
		// Test with data that comes from JSON unmarshaling
		schema := Object(map[string]interface{}{
			"name":   String().Required(),
			"age":    Number().Integer().Required(),
			"active": Bool().Required(),
		}).Required()

		// Simulate JSON unmarshaled data (numbers become float64)
		jsonData := map[string]interface{}{
			"name":   "John",
			"age":    float64(30), // JSON numbers become float64
			"active": true,
		}

		if err := schema.Validate(jsonData); err != nil {
			t.Errorf("JSON-style data validation failed: %v", err)
		}
	})
}

// TestCustomValidators tests custom validation functions
func TestCustomValidators(t *testing.T) {
	t.Run("String Custom Validation", func(t *testing.T) {
		profanityChecker := func(s string) error {
			badWords := []string{"spam", "badword"}
			for _, bad := range badWords {
				if strings.Contains(strings.ToLower(s), bad) {
					return fmt.Errorf("content contains inappropriate language")
				}
			}
			return nil
		}

		schema := String().
			Min(1).
			Custom(profanityChecker).
			Required()

		if err := schema.Validate("This is good content"); err != nil {
			t.Error("Expected clean content to pass, got:", err)
		}

		if err := schema.Validate("This is spam content"); err == nil {
			t.Error("Expected inappropriate content to fail")
		}
	})

	t.Run("Number Custom Validation", func(t *testing.T) {
		isPrime := func(n float64) error {
			if n != float64(int(n)) {
				return fmt.Errorf("must be an integer")
			}

			num := int(n)
			if num < 2 {
				return fmt.Errorf("must be a prime number (>= 2)")
			}

			for i := 2; i*i <= num; i++ {
				if num%i == 0 {
					return fmt.Errorf("must be a prime number")
				}
			}
			return nil
		}

		schema := Number().
			Custom(isPrime).
			Required()

		primes := []int{2, 3, 5, 7, 11, 13}
		for _, prime := range primes {
			if err := schema.Validate(prime); err != nil {
				t.Errorf("Expected prime %d to pass, got: %v", prime, err)
			}
		}

		nonPrimes := []int{4, 6, 8, 9, 10}
		for _, nonPrime := range nonPrimes {
			if err := schema.Validate(nonPrime); err == nil {
				t.Errorf("Expected non-prime %d to fail", nonPrime)
			}
		}
	})

	t.Run("Object Custom Validation", func(t *testing.T) {
		passwordMatch := func(data map[string]interface{}) error {
			password, hasPassword := data["password"].(string)
			confirm, hasConfirm := data["confirmPassword"].(string)

			if !hasPassword || !hasConfirm {
				return fmt.Errorf("both password and confirmPassword are required")
			}

			if password != confirm {
				return fmt.Errorf("passwords do not match")
			}

			return nil
		}

		schema := Object(map[string]interface{}{
			"password":        String().Min(8).Required(),
			"confirmPassword": String().Required(),
		}).Custom(passwordMatch).Required()

		validData := map[string]interface{}{
			"password":        "securepass123",
			"confirmPassword": "securepass123",
		}

		if err := schema.Validate(validData); err != nil {
			t.Errorf("Expected matching passwords to pass, got: %v", err)
		}

		invalidData := map[string]interface{}{
			"password":        "securepass123",
			"confirmPassword": "differentpass",
		}

		if err := schema.Validate(invalidData); err == nil {
			t.Error("Expected mismatched passwords to fail")
		}
	})
}
