package main

import (
	"fmt"
	"log"

	"github.com/aymaneallaoui/zod-go/zod"
	"github.com/aymaneallaoui/zod-go/zod/validators"
)

func main() {
	fmt.Println("=== ZOD-GO IMPROVED DEVELOPER EXPERIENCE DEMO ===")
	fmt.Println()

	// Demonstration of improved autocompletion for WithMessage
	demonstrateImprovedWithMessage()
	fmt.Println()

	// Demonstration of improved array method chaining
	demonstrateImprovedArrayChaining()
	fmt.Println()

	// Practical examples showing the improvements
	practicalExamples()
}

func demonstrateImprovedWithMessage() {
	fmt.Println("🚀 IMPROVED WithMessage() Autocompletion:")
	fmt.Println("=" + "=40")

	// Now you get autocompletion for validation types!
	schema := validators.String().
		Min(3).
		Max(50).
		Required().
		WithMessage(zod.ValidationTypeMinLength, "Username must be at least 3 characters").
		WithMessage(zod.ValidationTypeMaxLength, "Username cannot exceed 50 characters").
		WithMessage(zod.ValidationTypeRequired, "Username is required")

	// Test with valid input
	if err := schema.Validate("john_doe"); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Valid username: john_doe")
	}

	// Test with invalid input to see custom messages
	if err := schema.Validate("jo"); err != nil {
		fmt.Printf("❌ Invalid username 'jo': %v\n", err)
	}

	if err := schema.Validate(""); err != nil {
		fmt.Printf("❌ Empty username: %v\n", err)
	}
}

func demonstrateImprovedArrayChaining() {
	fmt.Println("🔗 IMPROVED Array Method Chaining:")
	fmt.Println("=" + "=35")

	// Example 1: Regular array - can call Unique()
	elementSchema := validators.String().Min(1)
	regularArray := validators.Array(elementSchema).
		Min(1).
		Max(10)

	// This returns *UniqueArraySchema, preventing double .Unique() calls
	uniqueArray := regularArray.Unique()

	// After calling Unique(), you can still chain other methods
	finalSchema := uniqueArray.
		Required().
		WithMessage(zod.ValidationTypeUnique, "All tags must be unique").
		WithMessage(zod.ValidationTypeMinLength, "At least one tag is required")

	// Test valid unique array
	validTags := []interface{"golang", "typescript", "python"}
	if err := finalSchema.Validate(validTags); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Valid unique tags:", validTags)
	}

	// Test invalid array with duplicates
	invalidTags := []interface{"golang", "typescript", "golang"}
	if err := finalSchema.Validate(invalidTags); err != nil {
		fmt.Printf("❌ Invalid tags with duplicates: %v\n", err)
	}

	// Note: After calling .Unique(), you cannot call .Unique() again
	// This would cause a compile-time error:
	// uniqueArray.Unique() // ❌ This method doesn't exist on UniqueArraySchema
}

func practicalExamples() {
	fmt.Println("💼 PRACTICAL EXAMPLES:")
	fmt.Println("=" + "=20")

	// User registration schema with improved DX
	userRegistrationSchema := validators.Object(map[string]zod.Schema{
		"username": validators.String().
			Min(3).Max(30).
			Pattern(`^[a-zA-Z0-9_]+$`).
			Required().
			WithMessage(zod.ValidationTypeMinLength, "Username must be at least 3 characters").
			WithMessage(zod.ValidationTypeMaxLength, "Username cannot exceed 30 characters").
			WithMessage(zod.ValidationTypePattern, "Username can only contain letters, numbers, and underscores"),

		"email": validators.String().
			Email().
			Required().
			WithMessage(zod.ValidationTypeEmail, "Please provide a valid email address").
			WithMessage(zod.ValidationTypeRequired, "Email is required"),

		"tags": validators.Array(validators.String().Min(1)).
			Min(1).Max(5).
			Unique(). // Returns UniqueArraySchema, no more double .Unique() calls!
			Required().
			WithMessage(zod.ValidationTypeUnique, "All tags must be unique").
			WithMessage(zod.ValidationTypeMinLength, "At least one tag is required").
			WithMessage(zod.ValidationTypeMaxLength, "Maximum 5 tags allowed"),
	})

	// Test data
	userData := map[string]interface{}{
		"username": "john_doe_123",
		"email":    "john@example.com",
		"tags":     []interface{}{"developer", "golang", "typescript"},
	}

	if err := userRegistrationSchema.Validate(userData); err != nil {
		log.Printf("❌ User registration validation failed: %v", err)
	} else {
		fmt.Println("✅ User registration data is valid!")
		fmt.Printf("   Username: %s\n", userData["username"])
		fmt.Printf("   Email: %s\n", userData["email"])
		fmt.Printf("   Tags: %v\n", userData["tags"])
	}

	fmt.Println()

	// Example with validation errors to show improved error messages
	invalidUserData := map[string]interface{}{
		"username": "jo", // Too short
		"email":    "invalid-email",
		"tags":     []interface{}{"dev", "dev"}, // Duplicate tags
	}

	if err := userRegistrationSchema.Validate(invalidUserData); err != nil {
		fmt.Println("❌ Validation errors for invalid data:")
		if validationErr, ok := err.(*zod.ValidationError); ok {
			fmt.Printf("   %s\n", validationErr.ErrorJSON())
		} else {
			fmt.Printf("   %v\n", err)
		}
	}
}

func demonstrateValidationTypes() {
	fmt.Println("📝 AVAILABLE VALIDATION TYPE CONSTANTS:")
	fmt.Println("=" + "=40")

	fmt.Println("String validation types:")
	fmt.Printf("  - %s\n", zod.ValidationTypeRequired)
	fmt.Printf("  - %s\n", zod.ValidationTypeMinLength)
	fmt.Printf("  - %s\n", zod.ValidationTypeMaxLength)
	fmt.Printf("  - %s\n", zod.ValidationTypePattern)
	fmt.Printf("  - %s\n", zod.ValidationTypeEmail)
	fmt.Printf("  - %s\n", zod.ValidationTypeURL)

	fmt.Println("\nArray validation types:")
	fmt.Printf("  - %s\n", zod.ValidationTypeUnique)
	fmt.Printf("  - %s\n", zod.ValidationTypeNonEmpty)
	fmt.Printf("  - %s\n", zod.ValidationTypeElements)

	fmt.Println("\nNumber validation types:")
	fmt.Printf("  - %s\n", zod.ValidationTypeMin)
	fmt.Printf("  - %s\n", zod.ValidationTypeMax)
	fmt.Printf("  - %s\n", zod.ValidationTypeInteger)
	fmt.Printf("  - %s\n", zod.ValidationTypePositive)
}
