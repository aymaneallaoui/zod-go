package main

import (
	"fmt"
	"log"

	"github.com/aymaneallaoui/zod-go/zod"
	"github.com/aymaneallaoui/zod-go/zod/validators"
)

// Example 1: Basic string validation
func basicStringValidation() {
	fmt.Println("=== Basic String Validation ===")

	// Create a string schema with various constraints
	schema := validators.String().
		Min(3).
		Max(50).
		Required().
		WithMessage("minLength", "Username must be at least 3 characters").
		WithMessage("maxLength", "Username cannot exceed 50 characters")

	// Test valid input
	if err := schema.Validate("john_doe"); err != nil {
		fmt.Printf("Validation failed: %v\n", err)
	} else {
		fmt.Println("✓ Valid username: john_doe")
	}

	// Test invalid input
	if err := schema.Validate("jo"); err != nil {
		fmt.Printf("✗ Invalid username: %v\n", err)
	}

	fmt.Println()
}

// Example 2: Email validation
func emailValidation() {
	fmt.Println("=== Email Validation ===")

	emailSchema := validators.String().Email().Required()

	emails := []string{
		"user@example.com",
		"john.doe+test@company.org",
		"invalid-email",
		"",
	}

	for _, email := range emails {
		if err := emailSchema.Validate(email); err != nil {
			fmt.Printf("✗ Invalid email '%s': %v\n", email, err)
		} else {
			fmt.Printf("✓ Valid email: %s\n", email)
		}
	}

	fmt.Println()
}

// Example 3: Number validation with ranges
func numberValidation() {
	fmt.Println("=== Number Validation ===")

	ageSchema := validators.Number().
		Integer().
		Min(0).
		Max(150).
		Required().
		WithMessage("min", "Age must be at least 0").
		WithMessage("max", "Age cannot exceed 150").
		WithMessage("integer", "Age must be a whole number")

	ages := []interface{}{25, 150, -5, 200, 25.5, "30"}

	for _, age := range ages {
		if err := ageSchema.Validate(age); err != nil {
			fmt.Printf("✗ Invalid age %v: %v\n", age, err)
		} else {
			fmt.Printf("✓ Valid age: %v\n", age)
		}
	}

	fmt.Println()
}

// Example 4: Array validation
func arrayValidation() {
	fmt.Println("=== Array Validation ===")

	// Array of unique strings with length constraints
	tagsSchema := validators.Array(validators.String().Min(1).Required()).
		Min(1).
		Max(5).
		Unique().
		Required()

	testArrays := []interface{}{
		[]string{"golang", "typescript", "python"},
		[]string{"a", "b", "a"}, // duplicate
		[]string{},              // empty
		[]string{"", "valid"},   // empty string element
	}

	for i, arr := range testArrays {
		if err := tagsSchema.Validate(arr); err != nil {
			fmt.Printf("✗ Invalid array %d: %v\n", i+1, err)
		} else {
			fmt.Printf("✓ Valid array %d: %v\n", i+1, arr)
		}
	}

	fmt.Println()
}

// Example 5: Object validation
func objectValidation() {
	fmt.Println("=== Object Validation ===")

	userSchema := validators.Object(map[string]zod.Schema{
		"id":    validators.Number().Integer().Positive().Required(),
		"name":  validators.String().Min(1).Max(100).Required(),
		"email": validators.String().Email().Required(),
		"age":   validators.Number().Integer().Min(13).Max(120).Optional(),
		"tags":  validators.Array(validators.String().Required()).Unique().Optional(),
	})

	// Valid user
	validUser := map[string]interface{}{
		"id":    1,
		"name":  "John Doe",
		"email": "john@example.com",
		"age":   30,
		"tags":  []interface{}{"developer", "golang"},
	}

	if err := userSchema.Validate(validUser); err != nil {
		fmt.Printf("✗ Invalid user: %v\n", err)
	} else {
		fmt.Println("✓ Valid user object")
	}

	// Invalid user (missing required field)
	invalidUser := map[string]interface{}{
		"id":   1,
		"name": "John Doe",
		// email is missing
	}

	if err := userSchema.Validate(invalidUser); err != nil {
		fmt.Printf("✗ Invalid user: %v\n", err)
	}

	fmt.Println()
}

// Example 6: Nested object validation
func nestedObjectValidation() {
	fmt.Println("=== Nested Object Validation ===")

	addressSchema := validators.Object(map[string]zod.Schema{
		"street":  validators.String().Min(1).Required(),
		"city":    validators.String().Min(1).Required(),
		"zipCode": validators.String().Pattern(`^\d{5}(-\d{4})?$`).Optional(),
	})

	userSchema := validators.Object(map[string]zod.Schema{
		"name":    validators.String().Required(),
		"email":   validators.String().Email().Required(),
		"address": addressSchema.Required(),
	})

	user := map[string]interface{}{
		"name":  "John Doe",
		"email": "john@example.com",
		"address": map[string]interface{}{
			"street":  "123 Main St",
			"city":    "New York",
			"zipCode": "10001",
		},
	}

	if err := userSchema.Validate(user); err != nil {
		fmt.Printf("✗ Invalid nested object: %v\n", err)
	} else {
		fmt.Println("✓ Valid nested object")
	}

	fmt.Println()
}

// Example 7: Custom validation
func customValidation() {
	fmt.Println("=== Custom Validation ===")

	// Custom validator for password strength
	passwordValidator := func(s string) error {
		hasUpper := false
		hasLower := false
		hasDigit := false

		for _, r := range s {
			switch {
			case r >= 'A' && r <= 'Z':
				hasUpper = true
			case r >= 'a' && r <= 'z':
				hasLower = true
			case r >= '0' && r <= '9':
				hasDigit = true
			}
		}

		if !hasUpper || !hasLower || !hasDigit {
			return zod.NewValidationError("password", s, "password must contain uppercase, lowercase, and digit")
		}
		return nil
	}

	passwordSchema := validators.String().
		Min(8).
		Custom(passwordValidator).
		Required()

	passwords := []string{
		"SecurePass123",
		"weakpass",
		"ALLUPPERCASE123",
		"alllowercase123",
		"NoDigitsHere",
	}

	for _, pwd := range passwords {
		if err := passwordSchema.Validate(pwd); err != nil {
			fmt.Printf("✗ Weak password '%s': %v\n", pwd, err)
		} else {
			fmt.Printf("✓ Strong password: %s\n", pwd)
		}
	}

	fmt.Println()
}

// Example 8: Concurrent validation
func concurrentValidation() {
	fmt.Println("=== Concurrent Validation ===")

	schema := validators.String().Email().Required()

	emails := []interface{}{
		"user1@example.com",
		"user2@example.com",
		"invalid-email",
		"user3@example.com",
		"another-invalid",
		"user4@example.com",
	}

	results := zod.ValidateConcurrently(schema, emails, 3)

	for i, result := range results {
		if result.IsValid {
			fmt.Printf("✓ Email %d is valid\n", i+1)
		} else {
			fmt.Printf("✗ Email %d is invalid: %v\n", i+1, result.Error)
		}
	}

	fmt.Println()
}

// Example 9: Schema composition and reuse
func schemaComposition() {
	fmt.Println("=== Schema Composition ===")

	// Base user schema
	baseUserSchema := validators.Object(map[string]zod.Schema{
		"id":   validators.Number().Integer().Positive().Required(),
		"name": validators.String().Min(1).Required(),
	})

	// Extend base schema for admin user
	adminUserSchema := baseUserSchema.Extend(map[string]zod.Schema{
		"role":        validators.String().Pattern(`^(admin|superadmin)$`).Required(),
		"permissions": validators.Array(validators.String().Required()).Min(1).Required(),
	})

	// Pick only specific fields
	publicUserSchema := baseUserSchema.Pick("name")

	adminUser := map[string]interface{}{
		"id":          1,
		"name":        "Admin User",
		"role":        "admin",
		"permissions": []interface{}{"read", "write", "delete"},
	}

	if err := adminUserSchema.Validate(adminUser); err != nil {
		fmt.Printf("✗ Invalid admin user: %v\n", err)
	} else {
		fmt.Println("✓ Valid admin user")
	}

	publicUser := map[string]interface{}{
		"name": "Public User",
	}

	if err := publicUserSchema.Validate(publicUser); err != nil {
		fmt.Printf("✗ Invalid public user: %v\n", err)
	} else {
		fmt.Println("✓ Valid public user")
	}

	fmt.Println()
}

// Example 10: Boolean validation with type conversion
func booleanValidation() {
	fmt.Println("=== Boolean Validation ===")

	consentSchema := validators.Bool().True().Required().
		WithMessage("true", "You must accept the terms and conditions")

	values := []interface{}{
		true,
		"true",
		1,
		false,
		"false",
		0,
		"yes",
		"no",
	}

	for _, val := range values {
		if err := consentSchema.Validate(val); err != nil {
			fmt.Printf("✗ Invalid consent %v: %v\n", val, err)
		} else {
			fmt.Printf("✓ Valid consent: %v\n", val)
		}
	}

	fmt.Println()
}

func basicUsage() {
	fmt.Println("Zod-Go Validation Examples")
	fmt.Println("==========================")
	fmt.Println()

	// Run all examples
	basicStringValidation()
	emailValidation()
	numberValidation()
	arrayValidation()
	objectValidation()
	nestedObjectValidation()
	customValidation()
	concurrentValidation()
	schemaComposition()
	booleanValidation()

	fmt.Println("All examples completed!")
}

// Helper function to demonstrate error JSON output
func demonstrateErrorJSON() {
	schema := validators.Object(map[string]zod.Schema{
		"user": validators.Object(map[string]zod.Schema{
			"name":  validators.String().Min(3).Required(),
			"email": validators.String().Email().Required(),
		}).Required(),
	})

	invalidData := map[string]interface{}{
		"user": map[string]interface{}{
			"name":  "Jo", // too short
			"email": "invalid-email",
		},
	}

	if err := schema.Validate(invalidData); err != nil {
		if validationErr, ok := err.(*zod.ValidationError); ok {
			fmt.Println("Error JSON:")
			fmt.Println(validationErr.ErrorJSON())
		}
	}
}

// Example of using the library in a web API context
func webAPIExample() {
	// This would typically be in a web handler
	userRegistrationSchema := validators.Object(map[string]zod.Schema{
		"username": validators.String().
			Min(3).Max(30).
			Pattern(`^[a-zA-Z0-9_]+$`).
			Required().
			WithMessage("pattern", "Username can only contain letters, numbers, and underscores"),
		"email": validators.String().Email().Required(),
		"password": validators.String().
			Min(8).
			Pattern(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]`).
			Required().
			WithMessage("pattern", "Password must contain uppercase, lowercase, digit, and special character"),
		"age":             validators.Number().Integer().Min(13).Max(120).Required(),
		"acceptedTerms":   validators.Bool().True().Required(),
		"marketingEmails": validators.Bool().Default(false),
	}).Strict() // Reject unknown fields

	// Simulate incoming request data
	requestData := map[string]interface{}{
		"username":        "john_doe123",
		"email":           "john@example.com",
		"password":        "SecurePass123!",
		"age":             25,
		"acceptedTerms":   true,
		"marketingEmails": false,
	}

	if err := userRegistrationSchema.Validate(requestData); err != nil {
		// In a real API, you'd return a 400 Bad Request with the error details
		log.Printf("Registration validation failed: %v", err)
		return
	}

	// Continue with user registration...
	log.Println("User registration data is valid")
}
