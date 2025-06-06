package validators

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

// TestRegexPatternFixes demonstrates how to replace unsupported regex patterns
func TestRegexPatternFixes(t *testing.T) {
	t.Run("Password Validation Without Lookaheads", func(t *testing.T) {
		// ❌ This pattern uses lookaheads which Go doesn't support:
		// "^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]"

		// ✅ Instead, use custom validation function:
		passwordValidator := func(password string) error {
			if len(password) < 8 {
				return fmt.Errorf("password must be at least 8 characters long")
			}
			
			hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
			hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
			hasDigit := regexp.MustCompile(`\d`).MatchString(password)
			hasSpecial := regexp.MustCompile(`[@$!%*?&]`).MatchString(password)
			
			if !hasLower {
				return fmt.Errorf("password must contain at least one lowercase letter")
			}
			if !hasUpper {
				return fmt.Errorf("password must contain at least one uppercase letter")
			}
			if !hasDigit {
				return fmt.Errorf("password must contain at least one digit")
			}
			if !hasSpecial {
				return fmt.Errorf("password must contain at least one special character (@$!%%*?&)")
			}
			
			return nil
		}

		// Create schema using custom validation instead of regex pattern
		schema := String().
			Min(8).
			Custom(passwordValidator).
			Required()

		// Test valid passwords
		validPasswords := []string{
			"Password1!",
			"MyStr0ng@Pass",
			"Secure123$",
		}

		for _, password := range validPasswords {
			if err := schema.Validate(password); err != nil {
				t.Errorf("Expected valid password '%s' to pass, got: %v", password, err)
			}
		}

		// Test invalid passwords
		invalidPasswords := []string{
			"password",      // No uppercase, digits, or special chars
			"PASSWORD1!",    // No lowercase
			"Password!",     // No digits
			"Password1",     // No special chars
			"Pass1!",        // Too short
		}

		for _, password := range invalidPasswords {
			if err := schema.Validate(password); err == nil {
				t.Errorf("Expected invalid password '%s' to fail", password)
			}
		}
	})

	t.Run("Email Domain Validation", func(t *testing.T) {
		// ❌ Lookahead pattern: "^(?=.*@.*\\.com$).*@.*\\.com$"
		
		// ✅ Use simple pattern + custom validation:
		domainValidator := func(email string) error {
			if !strings.HasSuffix(email, ".com") {
				return fmt.Errorf("email must have .com domain")
			}
			return nil
		}

		schema := String().
			Email().                    // Basic email validation
			Custom(domainValidator).    // Additional domain check
			Required()

		// Test valid emails
		if err := schema.Validate("user@example.com"); err != nil {
			t.Errorf("Expected valid .com email to pass, got: %v", err)
		}

		// Test invalid emails
		if err := schema.Validate("user@example.org"); err == nil {
			t.Error("Expected non-.com email to fail")
		}
	})

	t.Run("Multiple Condition Validation", func(t *testing.T) {
		// ❌ Complex lookahead: "^(?=.*foo)(?=.*bar)(?!.*baz).*$"
		
		// ✅ Use custom validation:
		multiValidator := func(text string) error {
			if !strings.Contains(text, "foo") {
				return fmt.Errorf("text must contain 'foo'")
			}
			if !strings.Contains(text, "bar") {
				return fmt.Errorf("text must contain 'bar'")
			}
			if strings.Contains(text, "baz") {
				return fmt.Errorf("text must not contain 'baz'")
			}
			return nil
		}

		schema := String().
			Custom(multiValidator).
			Required()

		// Test valid text
		if err := schema.Validate("foo and bar are here"); err != nil {
			t.Errorf("Expected valid text to pass, got: %v", err)
		}

		// Test invalid text
		if err := schema.Validate("foo and bar and baz"); err == nil {
			t.Error("Expected text with 'baz' to fail")
		}
	})

	t.Run("Alternatives to Common Lookahead Patterns", func(t *testing.T) {
		tests := []struct {
			name        string
			description string
			validator   func(string) error
			valid       []string
			invalid     []string
		}{
			{
				name:        "Username with requirements",
				description: "Must start with letter, contain letter and number",
				validator: func(s string) error {
					if len(s) < 3 {
						return fmt.Errorf("username too short")
					}
					if !regexp.MustCompile(`^[a-zA-Z]`).MatchString(s) {
						return fmt.Errorf("must start with letter")
					}
					if !regexp.MustCompile(`[a-zA-Z]`).MatchString(s) {
						return fmt.Errorf("must contain letter")
					}
					if !regexp.MustCompile(`[0-9]`).MatchString(s) {
						return fmt.Errorf("must contain number")
					}
					if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(s) {
						return fmt.Errorf("invalid characters")
					}
					return nil
				},
				valid:   []string{"user123", "test_user1", "admin42"},
				invalid: []string{"123user", "user", "user@name"},
			},
			{
				name:        "Phone number format",
				description: "Must be valid phone format",
				validator: func(s string) error {
					// Remove common separators for validation
					cleaned := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(s, "-", ""), " ", ""), "(", "")
					cleaned = strings.ReplaceAll(strings.ReplaceAll(cleaned, ")", ""), "+", "")
					
					if !regexp.MustCompile(`^\d{10,15}$`).MatchString(cleaned) {
						return fmt.Errorf("invalid phone number format")
					}
					return nil
				},
				valid:   []string{"+1234567890", "(123) 456-7890", "123-456-7890"},
				invalid: []string{"123", "abc-def-ghij", "123-45-6789"},
			},
			{
				name:        "Hex color validation",
				description: "Must be valid hex color",
				validator: func(s string) error {
					if !regexp.MustCompile(`^#[0-9a-fA-F]{6}$`).MatchString(s) {
						return fmt.Errorf("must be valid hex color (#RRGGBB)")
					}
					return nil
				},
				valid:   []string{"#FF0000", "#00ff00", "#0000FF", "#123ABC"},
				invalid: []string{"#FF", "#GGGGGG", "FF0000", "#FF00GG"},
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				schema := String().Custom(test.validator).Required()

				// Test valid cases
				for _, valid := range test.valid {
					if err := schema.Validate(valid); err != nil {
						t.Errorf("Expected '%s' to be valid for %s, got: %v", valid, test.description, err)
					}
				}

				// Test invalid cases
				for _, invalid := range test.invalid {
					if err := schema.Validate(invalid); err == nil {
						t.Errorf("Expected '%s' to be invalid for %s", invalid, test.description)
					}
				}
			})
		}
	})
}

// Example of a reusable password validation function
func StrongPasswordValidator(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}
	
	checks := []struct {
		pattern string
		message string
	}{
		{`[a-z]`, "password must contain at least one lowercase letter"},
		{`[A-Z]`, "password must contain at least one uppercase letter"},
		{`\d`, "password must contain at least one digit"},
		{`[@$!%*?&]`, "password must contain at least one special character (@$!%%*?&)"},
	}
	
	for _, check := range checks {
		if matched, _ := regexp.MatchString(check.pattern, password); !matched {
			return fmt.Errorf(check.message)
		}
	}
	
	return nil
}

// Usage example function (not a test)
func ExamplePasswordValidation() {
	// Create a strong password schema
	passwordSchema := String().
		Min(8).
		Max(128).
		Custom(StrongPasswordValidator).
		Required()

	// Test it
	passwords := []string{
		"Password123!",  // Valid
		"weakpass",      // Invalid - missing uppercase, digit, special char
	}

	for _, pwd := range passwords {
		if err := passwordSchema.Validate(pwd); err != nil {
			fmt.Printf("Password '%s' failed: %v\n", pwd, err)
		} else {
			fmt.Printf("Password '%s' is valid\n", pwd)
		}
	}
}
