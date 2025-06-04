package main

import (
	"fmt"
	"log"

	"github.com/aymaneallaoui/zod-go/zod"
	"github.com/aymaneallaoui/zod-go/zod/validators"
)

// User represents a user in our system
type User struct {
	ID       int                    `json:"id"`
	Username string                 `json:"username"`
	Email    string                 `json:"email"`
	Profile  map[string]interface{} `json:"profile"`
	Tags     []string               `json:"tags"`
	Active   bool                   `json:"active"`
}

func main() {
	fmt.Println("=== Zod-Go User Registration Example ===")

	// Define comprehensive user validation schema
	userSchema := validators.Object(map[string]zod.Schema{
		"id": validators.Number().
			Integer().
			Positive().
			Required().
			WithMessage("positive", "User ID must be a positive integer"),

		"username": validators.String().
			Min(3).
			Max(30).
			Pattern(`^[a-zA-Z0-9_]+$`).
			Required().
			WithMessage("pattern", "Username can only contain letters, numbers, and underscores").
			WithMessage("minLength", "Username must be at least 3 characters").
			WithMessage("maxLength", "Username cannot exceed 30 characters"),

		"email": validators.String().
			Email().
			Required().
			WithMessage("email", "Please provide a valid email address"),

		"profile": validators.Object(map[string]zod.Schema{
			"firstName": validators.String().Min(1).Required(),
			"lastName":  validators.String().Min(1).Required(),
			"age":       validators.Number().Integer().Min(13).Max(120).Required(),
			"bio":       validators.String().Max(500).Optional(),
			"website":   validators.String().URL().Optional(),
		}).Required(),

		"tags": validators.Array(validators.String().Min(1).Required()).
			Unique().
			Max(10).
			Optional(),

		"active": validators.Bool().
			Required().
			WithMessage("required", "Active status is required"),
	}).Strict() // Reject unknown properties

	// Test valid user data
	fmt.Println("\n--- Testing Valid User Data ---")
	validUser := map[string]interface{}{
		"id":       1,
		"username": "john_doe123",
		"email":    "john.doe@example.com",
		"profile": map[string]interface{}{
			"firstName": "John",
			"lastName":  "Doe",
			"age":       28,
			"bio":       "Software developer passionate about Go",
			"website":   "https://johndoe.dev",
		},
		"tags":   []interface{}{"developer", "golang", "backend"},
		"active": true,
	}

	if err := userSchema.Validate(validUser); err != nil {
		log.Printf("❌ Valid user validation failed: %v", err)
	} else {
		fmt.Println("✅ Valid user passed validation!")
		printUserData(validUser)
	}

	// Test invalid user data
	fmt.Println("\n--- Testing Invalid User Data ---")
	invalidUsers := []map[string]interface{}{
		{
			"id":       0, // Invalid: must be positive
			"username": "john_doe123",
			"email":    "john.doe@example.com",
			"profile": map[string]interface{}{
				"firstName": "John",
				"lastName":  "Doe",
				"age":       28,
			},
			"active": true,
		},
		{
			"id":       1,
			"username": "jo", // Invalid: too short
			"email":    "john.doe@example.com",
			"profile": map[string]interface{}{
				"firstName": "John",
				"lastName":  "Doe",
				"age":       28,
			},
			"active": true,
		},
		{
			"id":       1,
			"username": "john_doe123",
			"email":    "invalid-email", // Invalid email format
			"profile": map[string]interface{}{
				"firstName": "John",
				"lastName":  "Doe",
				"age":       28,
			},
			"active": true,
		},
		{
			"id":       1,
			"username": "john_doe123",
			"email":    "john.doe@example.com",
			"profile": map[string]interface{}{
				"firstName": "John",
				"lastName":  "Doe",
				"age":       12, // Invalid: too young
			},
			"active": true,
		},
		{
			"id":       1,
			"username": "john_doe123",
			"email":    "john.doe@example.com",
			"profile": map[string]interface{}{
				"firstName": "John",
				"lastName":  "Doe",
				"age":       28,
			},
			"tags":   []interface{}{"developer", "developer"}, // Invalid: duplicate tags
			"active": true,
		},
		{
			"id":         1,
			"username":   "john_doe123",
			"email":      "john.doe@example.com",
			"profile": map[string]interface{}{
				"firstName": "John",
				"lastName":  "Doe",
				"age":       28,
			},
			"active":    true,
			"extraField": "not allowed", // Invalid: unknown field in strict mode
		},
	}

	for i, user := range invalidUsers {
		fmt.Printf("\nTesting invalid user #%d:\n", i+1)
		if err := userSchema.Validate(user); err != nil {
			if validationErr, ok := err.(*zod.ValidationError); ok {
				fmt.Printf("❌ Validation failed: %s\n", validationErr.ErrorJSON())
			} else {
				fmt.Printf("❌ Validation failed: %v\n", err)
			}
		} else {
			fmt.Println("⚠️  Unexpected: Invalid user passed validation!")
		}
	}

	// Demonstrate custom validation
	fmt.Println("\n--- Custom Validation Example ---")
	
	// Custom validator for username availability
	usernameAvailabilityValidator := func(s string) error {
		// Simulate checking against a database
		reservedUsernames := []string{"admin", "root", "system", "support"}
		for _, reserved := range reservedUsernames {
			if s == reserved {
				return zod.NewValidationError("username", s, 
					fmt.Sprintf("Username '%s' is reserved and cannot be used", s))
			}
		}
		return nil
	}

	customUserSchema := validators.Object(map[string]zod.Schema{
		"username": validators.String().
			Min(3).
			Max(30).
			Pattern(`^[a-zA-Z0-9_]+$`).
			Custom(usernameAvailabilityValidator).
			Required(),
		"email": validators.String().Email().Required(),
	})

	testUsers := []map[string]interface{}{
		{
			"username": "john_doe",
			"email":    "john@example.com",
		},
		{
			"username": "admin", // Reserved username
			"email":    "admin@example.com",
		},
	}

	for _, user := range testUsers {
		fmt.Printf("\nTesting username: %s\n", user["username"])
		if err := customUserSchema.Validate(user); err != nil {
			fmt.Printf("❌ %v\n", err)
		} else {
			fmt.Printf("✅ Username '%s' is available!\n", user["username"])
		}
	}
}

func printUserData(user map[string]interface{}) {
	fmt.Printf("User Details:\n")
	fmt.Printf("  ID: %v\n", user["id"])
	fmt.Printf("  Username: %v\n", user["username"])
	fmt.Printf("  Email: %v\n", user["email"])
	
	if profile, ok := user["profile"].(map[string]interface{}); ok {
		fmt.Printf("  Profile:\n")
		fmt.Printf("    Name: %v %v\n", profile["firstName"], profile["lastName"])
		fmt.Printf("    Age: %v\n", profile["age"])
		if bio, exists := profile["bio"]; exists {
			fmt.Printf("    Bio: %v\n", bio)
		}
		if website, exists := profile["website"]; exists {
			fmt.Printf("    Website: %v\n", website)
		}
	}
	
	if tags, ok := user["tags"].([]interface{}); ok && len(tags) > 0 {
		fmt.Printf("  Tags: %v\n", tags)
	}
	
	fmt.Printf("  Active: %v\n", user["active"])
}
