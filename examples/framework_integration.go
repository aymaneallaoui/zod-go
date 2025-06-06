// Package main demonstrates integration of zod-go enhanced DX API with popular Go web frameworks
// This shows real-world usage patterns for building robust, validated APIs
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/aymaneallaoui/zod-go/zod/validators"
)

func main() {
	fmt.Println("🌐 zod-go Framework Integration Examples")
	fmt.Println("=======================================")

	// Since we can't import external frameworks in this example,
	// we'll demonstrate the patterns using standard library HTTP
	// and show how they would work with popular frameworks

	runStandardLibraryExample()
	runMiddlewareExamples()
	runValidationPatterns()
	runErrorHandlingExamples()
}

// Example using Go standard library HTTP
func runStandardLibraryExample() {
	fmt.Println("\n🔧 Standard Library HTTP Integration")
	fmt.Println("===================================")

	// Define validation schemas for different endpoints
	userCreateSchema := validators.Object(map[string]interface{}{
		"username": validators.String().
			Min(3).WithMinLengthMessage("Username must be at least 3 characters").
			Max(20).WithMaxLengthMessage("Username cannot exceed 20 characters").
			Pattern(`^[a-zA-Z0-9_]+$`).WithPatternMessage("Username can only contain letters, numbers, and underscores").
			Required().WithRequiredMessage("Username is required"),
		"email": validators.Email().
			WithEmailMessage("Please provide a valid email address"),
		"password": validators.String().
			Min(8).WithMinLengthMessage("Password must be at least 8 characters").
			Required().WithRequiredMessage("Password is required"),
		"age": validators.Number().
			Min(13).WithMinMessage("Must be at least 13 years old").
			Max(120).WithMaxMessage("Please provide a valid age").
			Integer().WithIntegerMessage("Age must be a whole number").
			Optional(),
	}).Required().WithRequiredMessage("User data is required")

	userUpdateSchema := validators.Object(map[string]interface{}{
		"username": validators.String().
			Min(3).Max(20).
			Pattern(`^[a-zA-Z0-9_]+$`).
			Optional(), // Updates are optional
		"email": validators.Email().Optional(),
		"age": validators.Number().
			Min(13).Max(120).Integer().
			Optional(),
	}).Required()

	// HTTP handlers using the schemas
	createUserHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var userData map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
			writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON", err.Error())
			return
		}

		if err := userCreateSchema.Validate(userData); err != nil {
			writeErrorResponse(w, http.StatusBadRequest, "Validation failed", err.Error())
			return
		}

		// Simulate user creation
		userID := generateID()
		response := map[string]interface{}{
			"success": true,
			"message": "User created successfully",
			"data": map[string]interface{}{
				"id":       userID,
				"username": userData["username"],
				"email":    userData["email"],
				"created_at": time.Now().Format(time.RFC3339),
			},
		}

		writeJSONResponse(w, http.StatusCreated, response)
	}

	updateUserHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var updateData map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
			writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON", err.Error())
			return
		}

		if err := userUpdateSchema.Validate(updateData); err != nil {
			writeErrorResponse(w, http.StatusBadRequest, "Validation failed", err.Error())
			return
		}

		response := map[string]interface{}{
			"success": true,
			"message": "User updated successfully",
			"data": map[string]interface{}{
				"updated_fields": getUpdatedFields(updateData),
				"updated_at":     time.Now().Format(time.RFC3339),
			},
		}

		writeJSONResponse(w, http.StatusOK, response)
	}

	fmt.Println("✅ Standard library handlers configured with zod-go validation")
	fmt.Println("   Example usage patterns:")
	fmt.Println("   - POST /users - Create user with full validation")
	fmt.Println("   - PUT /users/:id - Update user with partial validation")

	// Test the handlers with sample data
	testCreateUser(createUserHandler)
	testUpdateUser(updateUserHandler)

	_ = createUserHandler
	_ = updateUserHandler
}

// Middleware examples for different frameworks
func runMiddlewareExamples() {
	fmt.Println("\n🔌 Middleware Examples")
	fmt.Println("======================")

	// Generic validation middleware pattern
	fmt.Println("\n1. Generic Validation Middleware:")
	validationMiddleware := createValidationMiddleware(validators.Object(map[string]interface{}{
		"name": validators.String().Min(1).Required(),
		"age":  validators.Number().Integer().Min(0).Required(),
	}).Required())

	// Gin-style middleware pattern
	fmt.Println("\n2. Gin-style Middleware Pattern:")
	fmt.Println("```go")
	fmt.Println("func ValidateJSON(schema interface{ Validate(interface{}) error }) gin.HandlerFunc {")
	fmt.Println("    return func(c *gin.Context) {")
	fmt.Println("        var data map[string]interface{}")
	fmt.Println("        if err := c.ShouldBindJSON(&data); err != nil {")
	fmt.Println("            c.JSON(400, gin.H{\"error\": \"Invalid JSON\"})")
	fmt.Println("            c.Abort()")
	fmt.Println("            return")
	fmt.Println("        }")
	fmt.Println("        if err := schema.Validate(data); err != nil {")
	fmt.Println("            c.JSON(400, gin.H{\"error\": err.Error()})")
	fmt.Println("            c.Abort()")
	fmt.Println("            return")
	fmt.Println("        }")
	fmt.Println("        c.Set(\"validatedData\", data)")
	fmt.Println("        c.Next()")
	fmt.Println("    }")
	fmt.Println("}")
	fmt.Println("```")

	// Echo-style middleware pattern
	fmt.Println("\n3. Echo-style Middleware Pattern:")
	fmt.Println("```go")
	fmt.Println("func ValidateRequest(schema interface{ Validate(interface{}) error }) echo.MiddlewareFunc {")
	fmt.Println("    return func(next echo.HandlerFunc) echo.HandlerFunc {")
	fmt.Println("        return func(c echo.Context) error {")
	fmt.Println("            var data map[string]interface{}")
	fmt.Println("            if err := c.Bind(&data); err != nil {")
	fmt.Println("                return echo.NewHTTPError(400, \"Invalid JSON\")")
	fmt.Println("            }")
	fmt.Println("            if err := schema.Validate(data); err != nil {")
	fmt.Println("                return echo.NewHTTPError(400, err.Error())")
	fmt.Println("            }")
	fmt.Println("            c.Set(\"validatedData\", data)")
	fmt.Println("            return next(c)")
	fmt.Println("        }")
	fmt.Println("    }")
	fmt.Println("}")
	fmt.Println("```")

	// Fiber-style middleware pattern
	fmt.Println("\n4. Fiber-style Middleware Pattern:")
	fmt.Println("```go")
	fmt.Println("func ValidateBody(schema interface{ Validate(interface{}) error }) fiber.Handler {")
	fmt.Println("    return func(c *fiber.Ctx) error {")
	fmt.Println("        var data map[string]interface{}")
	fmt.Println("        if err := c.BodyParser(&data); err != nil {")
	fmt.Println("            return c.Status(400).JSON(fiber.Map{\"error\": \"Invalid JSON\"})")
	fmt.Println("        }")
	fmt.Println("        if err := schema.Validate(data); err != nil {")
	fmt.Println("            return c.Status(400).JSON(fiber.Map{\"error\": err.Error()})")
	fmt.Println("        }")
	fmt.Println("        c.Locals(\"validatedData\", data)")
	fmt.Println("        return c.Next()")
	fmt.Println("    }")
	fmt.Println("}")
	fmt.Println("```")

	fmt.Printf("✅ Middleware patterns demonstrate integration with %s\n", 
		"Gin, Echo, Fiber, and other frameworks")

	_ = validationMiddleware
}

// Common validation patterns for web applications
func runValidationPatterns() {
	fmt.Println("\n🎯 Common Validation Patterns")
	fmt.Println("=============================")

	// 1. Authentication/Registration patterns
	fmt.Println("\n1. Authentication & Registration:")
	
	loginSchema := validators.Object(map[string]interface{}{
		"email": validators.Email().
			WithEmailMessage("Please provide a valid email address").
			WithRequiredMessage("Email is required for login"),
		"password": validators.String().
			Min(1).WithMinLengthMessage("Password cannot be empty").
			Required().WithRequiredMessage("Password is required"),
		"remember_me": validators.Bool().Optional().Default(false),
	}).Required()

	registrationSchema := validators.Object(map[string]interface{}{
		"username": validators.String().
			Min(3).Max(20).
			Pattern(`^[a-zA-Z0-9_]+$`).
			Required(),
		"email": validators.Email(),
		"password": validators.String().
			Min(8).
			Pattern(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]`).
			WithPatternMessage("Password must contain uppercase, lowercase, digit, and special character").
			Required(),
		"confirm_password": validators.String().Required(),
		"terms_accepted": validators.Bool().
			Required().WithRequiredMessage("You must accept the terms of service"),
	}).Custom(func(data map[string]interface{}) error {
		password, _ := data["password"].(string)
		confirmPassword, _ := data["confirm_password"].(string)
		if password != confirmPassword {
			return fmt.Errorf("passwords do not match")
		}
		return nil
	}).Required()

	// 2. E-commerce patterns
	fmt.Println("\n2. E-commerce Validation:")
	
	productSchema := validators.Object(map[string]interface{}{
		"name": validators.String().
			Min(1).Max(200).
			Required(),
		"description": validators.String().
			Min(10).Max(5000).
			Required(),
		"price": validators.Number().
			Min(0.01).WithMinMessage("Price must be greater than 0").
			Required(),
		"currency": validators.String().
			Pattern(`^[A-Z]{3}$`).WithPatternMessage("Currency must be a 3-letter code").
			Required(),
		"category_id": validators.Number().Integer().Min(1).Required(),
		"tags": validators.Array(validators.String().Min(1)).
			MaxItems(20).Optional(),
		"in_stock": validators.Bool().Optional().Default(true),
	}).Required()

	orderSchema := validators.Object(map[string]interface{}{
		"items": validators.Array(validators.Object(map[string]interface{}{
			"product_id": validators.Number().Integer().Required(),
			"quantity":   validators.Number().Integer().Min(1).Required(),
			"price":      validators.Number().Min(0).Required(),
		})).MinItems(1).WithMinItemsMessage("Order must contain at least one item").Required(),
		"shipping_address": validators.Object(map[string]interface{}{
			"street":   validators.String().Min(1).Required(),
			"city":     validators.String().Min(1).Required(),
			"state":    validators.String().Min(1).Required(),
			"zip_code": validators.String().Pattern(`^\d{5}(-\d{4})?$`).Required(),
			"country":  validators.String().Pattern(`^[A-Z]{2}$`).Required(),
		}).Required(),
		"payment_method": validators.String().
			Pattern(`^(credit_card|paypal|bank_transfer)$`).
			Required(),
	}).Required()

	// 3. Content management patterns
	fmt.Println("\n3. Content Management:")
	
	articleSchema := validators.Object(map[string]interface{}{
		"title": validators.String().
			Min(1).Max(200).
			Required(),
		"slug": validators.String().
			Pattern(`^[a-z0-9-]+$`).
			Required(),
		"content": validators.String().
			Min(100).
			Required(),
		"excerpt": validators.String().
			Min(10).Max(500).
			Optional(),
		"tags": validators.Array(validators.String().Min(1)).
			MaxItems(10).
			Optional(),
		"status": validators.String().
			Pattern(`^(draft|published|archived)$`).
			Optional().Default("draft"),
		"publish_at": validators.String().Optional(), // ISO date string
		"featured": validators.Bool().Optional().Default(false),
	}).Required()

	// 4. API pagination patterns
	fmt.Println("\n4. API Pagination & Filtering:")
	
	paginationSchema := validators.Object(map[string]interface{}{
		"page": validators.Number().
			Integer().Min(1).
			Optional().Default(1.0),
		"per_page": validators.Number().
			Integer().Min(1).Max(100).
			Optional().Default(20.0),
		"sort_by": validators.String().
			Pattern(`^[a-zA-Z_]+$`).
			Optional().Default("created_at"),
		"sort_order": validators.String().
			Pattern(`^(asc|desc)$`).
			Optional().Default("desc"),
		"filters": validators.Object(map[string]interface{}{}).
			Optional(),
	}).Required()

	// Test the patterns
	testValidationPattern("Login", loginSchema, map[string]interface{}{
		"email":       "user@example.com",
		"password":    "securepass",
		"remember_me": true,
	})

	testValidationPattern("Registration", registrationSchema, map[string]interface{}{
		"username":         "newuser",
		"email":           "newuser@example.com",
		"password":        "SecurePass123!",
		"confirm_password": "SecurePass123!",
		"terms_accepted":  true,
	})

	testValidationPattern("Product", productSchema, map[string]interface{}{
		"name":        "Premium Headphones",
		"description": "High-quality wireless headphones with noise cancellation",
		"price":       299.99,
		"currency":    "USD",
		"category_id": 1,
		"tags":        []string{"electronics", "audio", "wireless"},
		"in_stock":    true,
	})

	testValidationPattern("Pagination", paginationSchema, map[string]interface{}{
		"page":       2,
		"per_page":   25,
		"sort_by":    "name",
		"sort_order": "asc",
	})
}

// Error handling examples
func runErrorHandlingExamples() {
	fmt.Println("\n❌ Error Handling Examples")
	fmt.Println("==========================")

	// Custom error response structure
	type ValidationErrorResponse struct {
		Success bool                   `json:"success"`
		Error   string                 `json:"error"`
		Details map[string]interface{} `json:"details,omitempty"`
		Code    string                 `json:"code"`
	}

	// Error handling middleware
	errorHandler := func(err error, statusCode int) *ValidationErrorResponse {
		response := &ValidationErrorResponse{
			Success: false,
			Error:   err.Error(),
			Code:    getErrorCode(statusCode),
		}

		// Parse validation errors for more detailed responses
		if validationErr, ok := err.(interface{ Error() string }); ok {
			response.Details = parseValidationError(validationErr.Error())
		}

		return response
	}

	// Example schema that will produce various errors
	testSchema := validators.Object(map[string]interface{}{
		"username": validators.String().
			Min(3).WithMinLengthMessage("Username too short").
			Max(20).WithMaxLengthMessage("Username too long").
			Pattern(`^[a-zA-Z0-9_]+$`).WithPatternMessage("Invalid username format").
			Required().WithRequiredMessage("Username is required"),
		"email": validators.Email().
			WithEmailMessage("Invalid email format").
			WithRequiredMessage("Email is required"),
		"age": validators.Number().
			Min(13).WithMinMessage("Must be at least 13").
			Max(120).WithMaxMessage("Invalid age").
			Integer().WithIntegerMessage("Age must be a whole number").
			Required(),
	}).Required()

	// Test different error scenarios
	errorScenarios := []struct {
		name string
		data map[string]interface{}
	}{
		{
			name: "Missing required fields",
			data: map[string]interface{}{},
		},
		{
			name: "Invalid username format",
			data: map[string]interface{}{
				"username": "user@name",
				"email":    "user@example.com",
				"age":      25,
			},
		},
		{
			name: "Invalid email",
			data: map[string]interface{}{
				"username": "username",
				"email":    "invalid-email",
				"age":      25,
			},
		},
		{
			name: "Age too young",
			data: map[string]interface{}{
				"username": "username",
				"email":    "user@example.com",
				"age":      10,
			},
		},
		{
			name: "Multiple validation errors",
			data: map[string]interface{}{
				"username": "ab",  // Too short
				"email":    "bad", // Invalid format
				"age":      5,     // Too young
			},
		},
	}

	fmt.Println("\nTesting error scenarios:")
	for _, scenario := range errorScenarios {
		fmt.Printf("\n🔍 %s:\n", scenario.name)
		
		if err := testSchema.Validate(scenario.data); err != nil {
			errorResponse := errorHandler(err, http.StatusBadRequest)
			errorJSON, _ := json.MarshalIndent(errorResponse, "   ", "  ")
			fmt.Printf("   Response: %s\n", string(errorJSON))
		} else {
			fmt.Println("   ✅ Validation passed (unexpected)")
		}
	}

	// Field-level error mapping example
	fmt.Println("\n📋 Field-Level Error Mapping:")
	fieldErrorMapper := func(err error) map[string]string {
		fieldErrors := make(map[string]string)
		
		// This would parse the error and map to specific fields
		errorStr := err.Error()
		
		// Simple parsing example (in real implementation, you'd have more sophisticated parsing)
		if containsError(errorStr, "username") {
			fieldErrors["username"] = extractFieldError(errorStr, "username")
		}
		if containsError(errorStr, "email") {
			fieldErrors["email"] = extractFieldError(errorStr, "email")
		}
		if containsError(errorStr, "age") {
			fieldErrors["age"] = extractFieldError(errorStr, "age")
		}
		
		return fieldErrors
	}

	// Test field error mapping
	invalidData := map[string]interface{}{
		"username": "ab",
		"email":    "invalid",
		"age":      5,
	}

	if err := testSchema.Validate(invalidData); err != nil {
		fieldErrors := fieldErrorMapper(err)
		fmt.Printf("   Field errors: %+v\n", fieldErrors)
	}
}

// Helper functions
func createValidationMiddleware(schema interface{ Validate(interface{}) error }) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var data map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
				writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON", err.Error())
				return
			}

			if err := schema.Validate(data); err != nil {
				writeErrorResponse(w, http.StatusBadRequest, "Validation failed", err.Error())
				return
			}

			// In a real implementation, you'd add the validated data to the request context
			next.ServeHTTP(w, r)
		})
	}
}

func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func writeErrorResponse(w http.ResponseWriter, statusCode int, message, details string) {
	response := map[string]interface{}{
		"success": false,
		"error":   message,
		"details": details,
		"code":    getErrorCode(statusCode),
	}
	writeJSONResponse(w, statusCode, response)
}

func getErrorCode(statusCode int) string {
	switch statusCode {
	case http.StatusBadRequest:
		return "VALIDATION_ERROR"
	case http.StatusUnauthorized:
		return "UNAUTHORIZED"
	case http.StatusForbidden:
		return "FORBIDDEN"
	case http.StatusNotFound:
		return "NOT_FOUND"
	case http.StatusInternalServerError:
		return "INTERNAL_ERROR"
	default:
		return "UNKNOWN_ERROR"
	}
}

func parseValidationError(errorStr string) map[string]interface{} {
	// Simple parsing - in real implementation, you'd have more sophisticated error parsing
	return map[string]interface{}{
		"raw_error": errorStr,
		"type":      "validation_error",
	}
}

func containsError(errorStr, field string) bool {
	// Simple check - in real implementation, you'd have better error parsing
	return len(errorStr) > 0 && len(field) > 0
}

func extractFieldError(errorStr, field string) string {
	// Simple extraction - in real implementation, you'd parse the actual error structure
	return fmt.Sprintf("Validation error for field: %s", field)
}

func testValidationPattern(name string, schema interface{ Validate(interface{}) error }, data interface{}) {
	fmt.Printf("   Testing %s pattern: ", name)
	if err := schema.Validate(data); err != nil {
		fmt.Printf("❌ %s\n", err.Error())
	} else {
		fmt.Printf("✅ Valid\n")
	}
}

func testCreateUser(handler http.HandlerFunc) {
	fmt.Println("\n🧪 Testing create user handler:")
	
	validUser := map[string]interface{}{
		"username": "john_doe",
		"email":    "john@example.com",
		"password": "securepass123",
		"age":      25,
	}

	testData, _ := json.Marshal(validUser)
	fmt.Printf("   Valid user data: %s\n", string(testData))
	fmt.Println("   ✅ Would create user successfully")

	invalidUser := map[string]interface{}{
		"username": "jo", // Too short
		"email":    "invalid-email",
		"password": "weak",
	}

	invalidData, _ := json.Marshal(invalidUser)
	fmt.Printf("   Invalid user data: %s\n", string(invalidData))
	fmt.Println("   ❌ Would return validation errors")
}

func testUpdateUser(handler http.HandlerFunc) {
	fmt.Println("\n🧪 Testing update user handler:")
	
	updateData := map[string]interface{}{
		"username": "new_username",
		"age":      26,
	}

	testData, _ := json.Marshal(updateData)
	fmt.Printf("   Update data: %s\n", string(testData))
	fmt.Println("   ✅ Would update user successfully")
}

func getUpdatedFields(data map[string]interface{}) []string {
	fields := make([]string, 0, len(data))
	for key := range data {
		fields = append(fields, key)
	}
	return fields
}

func generateID() string {
	return fmt.Sprintf("user_%d", time.Now().Unix())
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	fmt.Println("🚀 Framework integration examples initialized")
	fmt.Println("   Demonstrating patterns for Gin, Echo, Fiber, and standard library")
}
