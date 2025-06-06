// Package main demonstrates the enhanced developer experience (DX) of zod-go
// with clean package interface, type-safe state management, and smart autocompletion.
package main

import (
	"fmt"
	"log"

	"github.com/aymaneallaoui/zod-go/zod/validators"
)

func main() {
	fmt.Println("🎉 Enhanced DX Demo for zod-go")
	fmt.Println("===============================")

	demonstrateCleanInterface()
	demonstrateTypeSafetyFeatures()
	demonstrateErrorKeyAutocompletion()
	demonstrateComplexValidations()
	demonstrateCommonPatterns()
}

// demonstrateCleanInterface shows how the package interface is now clean
func demonstrateCleanInterface() {
	fmt.Println("\n✨ Clean Package Interface")
	fmt.Println("When you type 'validators.' you only see:")
	fmt.Println("- String(), Number(), Array(), Object(), Bool()")
	fmt.Println("- Email(), URL(), OptionalString(), RequiredString()")
	fmt.Println("- Errors (for error keys), Err* constants")
	fmt.Println("- NO internal functions or utilities!")

	// Clean, discoverable API
	_ = validators.String()        // ✅ Visible
	_ = validators.Number()        // ✅ Visible
	_ = validators.Array(nil)      // ✅ Visible
	_ = validators.Object(nil)     // ✅ Visible
	_ = validators.Bool()          // ✅ Visible
	_ = validators.Email()         // ✅ Visible
	_ = validators.URL()           // ✅ Visible
	_ = validators.Errors          // ✅ Visible
	_ = validators.ErrRequired     // ✅ Visible
}

// demonstrateTypeSafetyFeatures shows compile-time safety improvements
func demonstrateTypeSafetyFeatures() {
	fmt.Println("\n🔒 Type-Safe State Management")
	fmt.Println("The API prevents invalid method chaining at compile time:")

	// ✅ Valid progressions
	requiredUsername := validators.String().
		Min(3).
		Max(50).
		Required().  // Can only call this once!
		WithRequiredMessage("Username is required")

	optionalUsername := validators.String().
		Min(3).
		Max(50).
		Optional().                      // Can only call this once!
		Default("anonymous").            // Only available on optional!
		WithMinLengthMessage("Too short")

	// ❌ These would NOT compile (demonstrating type safety):
	// requiredUsername.Required()      // Cannot call Required() twice!
	// requiredUsername.Default("x")    // Required fields cannot have defaults!
	// optionalUsername.Optional()      // Cannot call Optional() twice!

	fmt.Println("✅ Valid: String -> Min -> Max -> Required -> Validation")
	fmt.Println("✅ Valid: String -> Min -> Max -> Optional -> Default -> Validation")
	fmt.Println("❌ Invalid: .Required().Required() (won't compile)")
	fmt.Println("❌ Invalid: .Required().Default() (won't compile)")

	// Test the validators
	testValidation("Required Username", requiredUsername, "john")
	testValidation("Optional Username with nil", optionalUsername, nil)
}

// demonstrateErrorKeyAutocompletion shows smart error key features
func demonstrateErrorKeyAutocompletion() {
	fmt.Println("\n🎯 Smart Error Key Autocompletion")
	fmt.Println("Multiple approaches for error key autocompletion:")

	// Approach 1: Method-based (cleanest autocompletion)
	fmt.Println("\n1. Method-based (validators.Errors.MinLength()):")
	emailSchema1 := validators.String().
		Email().
		Required().
		WithMessage(validators.Errors.Email(), "Please enter a valid email address").
		WithMessage(validators.Errors.Required(), "Email is required")

	// Approach 2: Constants (shorter syntax)
	fmt.Println("2. Constants (validators.ErrMinLength):")
	emailSchema2 := validators.String().
		Email().
		Required().
		WithMessage(validators.ErrEmail, "Please enter a valid email address").
		WithMessage(validators.ErrRequired, "Email is required")

	// Approach 3: Convenience methods (most fluent)
	fmt.Println("3. Convenience methods (.WithEmailMessage()):")
	emailSchema3 := validators.String().
		Email().WithEmailMessage("Please enter a valid email address").
		Required().WithRequiredMessage("Email is required")

	// Approach 4: Mixed (flexibility)
	fmt.Println("4. Mixed approaches:")
	phoneSchema := validators.String().
		Pattern(`^\+?[1-9]\d{1,14}$`).
		Required().
		WithMessage(validators.Errors.Pattern(), "Please enter a valid phone number").
		WithMessage(validators.ErrRequired, "Phone number is required")

	// Test all approaches
	testValidation("Email Schema 1", emailSchema1, "invalid-email")
	testValidation("Email Schema 2", emailSchema2, "test@example.com")
	testValidation("Email Schema 3", emailSchema3, "")
	testValidation("Phone Schema", phoneSchema, "+1234567890")
}

// demonstrateComplexValidations shows advanced validation patterns
func demonstrateComplexValidations() {
	fmt.Println("\n🔧 Complex Validation Patterns")

	// Number validation with various constraints
	ageSchema := validators.Number().
		Min(0).WithMinMessage("Age cannot be negative").
		Max(150).WithMaxMessage("Age cannot exceed 150 years").
		Integer().WithIntegerMessage("Age must be a whole number").
		Required().WithRequiredMessage("Age is required")

	// Array validation with element schemas
	tagsSchema := validators.Array(validators.String().Min(1)).
		MinItems(1).WithMinItemsMessage("At least one tag is required").
		MaxItems(10).WithMaxItemsMessage("Cannot have more than 10 tags").
		Required().WithRequiredMessage("Tags are required")

	// Complex object validation
	userSchema := validators.Object(map[string]interface{}{
		"name":     validators.RequiredString().Min(1),
		"email":    validators.Email(),
		"age":      validators.Number().Min(0).Optional().Default(18.0),
		"active":   validators.Bool().Optional().Default(true),
		"tags":     validators.Array(validators.String()).Optional(),
	}).Required().WithRequiredMessage("User object is required")

	// Test complex validations
	testValidation("Age Validation", ageSchema, 25)
	testValidation("Tags Validation", tagsSchema, []string{"golang", "validation"})
	
	userData := map[string]interface{}{
		"name":   "John Doe",
		"email":  "john@example.com",
		"age":    30,
		"active": true,
		"tags":   []string{"developer", "go"},
	}
	testValidation("User Validation", userSchema, userData)
}

// demonstrateCommonPatterns shows pre-configured validators
func demonstrateCommonPatterns() {
	fmt.Println("\n🚀 Common Validation Patterns")

	// Pre-configured validators for common use cases
	email := validators.Email()              // Already required email
	url := validators.URL()                  // Already required URL
	optionalStr := validators.OptionalString() // Already optional string
	requiredStr := validators.RequiredString() // Already required string
	positiveNum := validators.PositiveNumber() // Already positive constraint
	integerNum := validators.IntegerNumber()   // Already integer constraint

	// Demonstrate usage
	testValidation("Pre-configured Email", email, "test@example.com")
	testValidation("Pre-configured URL", url, "https://example.com")
	testValidation("Pre-configured Optional String", optionalStr, nil)
	testValidation("Pre-configured Required String", requiredStr, "hello")
	testValidation("Pre-configured Positive Number", positiveNum.Required(), 42)
	testValidation("Pre-configured Integer", integerNum.Required(), 42.5) // Should fail

	// Custom validation with fluent API
	customValidator := validators.String().
		Min(8).WithMinLengthMessage("Password must be at least 8 characters").
		Pattern(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d).*$`).WithPatternMessage("Password must contain uppercase, lowercase, and digit").
		Required().WithRequiredMessage("Password is required")

	testValidation("Custom Password Validator", customValidator, "weakpass")    // Should fail
	testValidation("Custom Password Validator", customValidator, "StrongPass1") // Should pass
}

// Helper function to test validations and display results
func testValidation(name string, validator interface{ Validate(interface{}) error }, value interface{}) {
	err := validator.Validate(value)
	if err != nil {
		fmt.Printf("❌ %s: %s\n", name, err.Error())
	} else {
		fmt.Printf("✅ %s: Valid\n", name)
	}
}

// Example of what won't compile (commented out to avoid compilation errors)
func demonstrateCompileTimeErrors() {
	// These examples show what the new API prevents at compile time:
	
	fmt.Println("\n🚫 Compile-Time Error Prevention")
	fmt.Println("The following would NOT compile:")

	/*
	// ❌ Cannot call Required() twice
	validators.String().Required().Required()

	// ❌ Cannot set Default() on required fields
	validators.String().Required().Default("value")

	// ❌ Cannot call Optional() twice  
	validators.String().Optional().Optional()

	// ❌ These methods don't exist in the wrong state
	requiredBuilder := validators.String().Required()
	requiredBuilder.Default("value")  // Method doesn't exist!

	optionalBuilder := validators.String().Optional()
	optionalBuilder.Required()        // Method doesn't exist!
	*/

	fmt.Println("- .Required().Required()")
	fmt.Println("- .Required().Default()")
	fmt.Println("- .Optional().Optional()")
	fmt.Println("- Invalid state transitions")
}

// Benchmark comparison (conceptual)
func demonstrateBenchmarkImprovement() {
	fmt.Println("\n📊 Developer Experience Improvements")
	
	fmt.Println("\nBefore (problems):")
	fmt.Println("- 50+ exports when typing 'validators.'")
	fmt.Println("- No compile-time safety for method chaining")
	fmt.Println("- Manual string literals for error keys")
	fmt.Println("- Possible to call .Required().Required()")
	fmt.Println("- Possible to set .Default() on required fields")

	fmt.Println("\nAfter (solutions):")
	fmt.Println("- Only 10 main exports when typing 'validators.'")
	fmt.Println("- Compile-time safety prevents invalid chaining")
	fmt.Println("- Smart autocompletion for error keys")
	fmt.Println("- Impossible to call .Required().Required()")
	fmt.Println("- Impossible to set .Default() on required fields")
	fmt.Println("- TypeScript-level DX in Go!")
}

// Advanced usage patterns
func demonstrateAdvancedPatterns() {
	fmt.Println("\n⚡ Advanced Usage Patterns")

	// Conditional validation based on other fields
	userSchema := validators.Object(map[string]interface{}{
		"type":     validators.String().Required(),
		"email":    validators.String().Optional(), // Will be conditionally required
		"phone":    validators.String().Optional(), // Will be conditionally required
	}).Custom(func(data map[string]interface{}) error {
		userType, _ := data["type"].(string)
		email, hasEmail := data["email"].(string)
		phone, hasPhone := data["phone"].(string)

		// Business rule: admins must have email, users must have email or phone
		if userType == "admin" && (!hasEmail || email == "") {
			return fmt.Errorf("admin users must have an email address")
		}
		
		if userType == "user" && (!hasEmail || email == "") && (!hasPhone || phone == "") {
			return fmt.Errorf("users must have either an email or phone number")
		}

		return nil
	}).Required()

	// Reusable validation components
	passwordValidator := validators.String().
		Min(8).WithMinLengthMessage("Password too short").
		Pattern(`^(?=.*[A-Z])(?=.*[a-z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]`).
		WithPatternMessage("Password must contain uppercase, lowercase, digit, and special character").
		Required()

	registrationSchema := validators.Object(map[string]interface{}{
		"username": validators.String().Min(3).Max(20).Required(),
		"email":    validators.Email(),
		"password": passwordValidator, // Reuse the validator
		"confirm":  validators.String().Required(),
	}).Custom(func(data map[string]interface{}) error {
		password, _ := data["password"].(string)
		confirm, _ := data["confirm"].(string)
		if password != confirm {
			return fmt.Errorf("passwords do not match")
		}
		return nil
	}).Required()

	// Test advanced patterns
	adminUser := map[string]interface{}{
		"type":  "admin",
		"email": "", // Should fail - admin needs email
	}

	registration := map[string]interface{}{
		"username": "john_doe",
		"email":    "john@example.com",
		"password": "StrongPass1!",
		"confirm":  "DifferentPass", // Should fail - passwords don't match
	}

	testValidation("Advanced User Schema", userSchema, adminUser)
	testValidation("Registration Schema", registrationSchema, registration)
}

func init() {
	// This demonstrates that the package is properly organized
	log.Println("🎯 zod-go enhanced DX demo initialized")
}
