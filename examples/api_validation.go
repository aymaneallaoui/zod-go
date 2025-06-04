package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/aymaneallaoui/zod-go/zod"
	"github.com/aymaneallaoui/zod-go/zod/validators"
)

// API Request/Response validation schemas
var (
	// User creation request schema
	createUserSchema = validators.Object(map[string]zod.Schema{
		"username": validators.String().
			Min(3).
			Max(30).
			Pattern(`^[a-zA-Z0-9_]+$`).
			Required(),
		"email": validators.String().
			Email().
			Required(),
		"password": validators.String().
			Min(8).
			Pattern(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]`).
			Required().
			WithMessage("pattern", "Password must contain at least one uppercase, lowercase, digit, and special character"),
		"profile": validators.Object(map[string]zod.Schema{
			"firstName": validators.String().Min(1).Max(50).Required(),
			"lastName":  validators.String().Min(1).Max(50).Required(),
			"age":       validators.Number().Integer().Min(13).Max(120).Optional(),
		}).Required(),
		"preferences": validators.Object(map[string]zod.Schema{
			"newsletter": validators.Bool().Default(false),
			"theme":      validators.String().Pattern(`^(light|dark|auto)$`).Default("auto"),
			"language":   validators.String().Pattern(`^[a-z]{2}$`).Default("en"),
		}).Optional(),
	}).Strict()

	// User update request schema (all fields optional except ID)
	updateUserSchema = validators.Object(map[string]zod.Schema{
		"id": validators.Number().Integer().Positive().Required(),
		"username": validators.String().
			Min(3).
			Max(30).
			Pattern(`^[a-zA-Z0-9_]+$`).
			Optional(),
		"email": validators.String().Email().Optional(),
		"profile": validators.Object(map[string]zod.Schema{
			"firstName": validators.String().Min(1).Max(50).Optional(),
			"lastName":  validators.String().Min(1).Max(50).Optional(),
			"age":       validators.Number().Integer().Min(13).Max(120).Optional(),
		}).Optional(),
	})

	// Query parameters schema for user listing
	userListQuerySchema = validators.Object(map[string]zod.Schema{
		"page":   validators.Number().Integer().Min(1).Default(1),
		"limit":  validators.Number().Integer().Min(1).Max(100).Default(10),
		"sort":   validators.String().Pattern(`^(id|username|email|created_at)$`).Default("id"),
		"order":  validators.String().Pattern(`^(asc|desc)$`).Default("asc"),
		"search": validators.String().Max(100).Optional(),
		"active": validators.Bool().Optional(),
	})

	// API response schema
	apiResponseSchema = validators.Object(map[string]zod.Schema{
		"success": validators.Bool().Required(),
		"message": validators.String().Optional(),
		"data":    validators.Object(map[string]zod.Schema{}).Optional(), // Any object
		"errors":  validators.Array(validators.String().Required()).Optional(),
		"meta": validators.Object(map[string]zod.Schema{
			"timestamp": validators.String().Pattern(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`).Required(),
			"version":   validators.String().Pattern(`^v\d+\.\d+\.\d+$`).Required(),
		}).Optional(),
	})
)

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Errors  []string    `json:"errors,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// Meta represents response metadata
type Meta struct {
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
}

// ValidationMiddleware validates request body against a schema
func ValidationMiddleware(schema zod.Schema) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			var body map[string]interface{}

			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				respondWithError(w, http.StatusBadRequest, "Invalid JSON format", nil)
				return
			}

			if err := schema.Validate(body); err != nil {
				if validationErr, ok := err.(*zod.ValidationError); ok {
					respondWithValidationError(w, validationErr)
				} else {
					respondWithError(w, http.StatusBadRequest, "Validation failed", []string{err.Error()})
				}
				return
			}

			// Store validated data in request context for handlers to use
			r.Header.Set("X-Validated-Body", "true")
			next(w, r)
		}
	}
}

// QueryValidationMiddleware validates query parameters against a schema
func QueryValidationMiddleware(schema zod.Schema) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			query := make(map[string]interface{})

			for key, values := range r.URL.Query() {
				if len(values) > 0 {
					value := values[0]
					// Try to convert to number if possible
					if num, err := strconv.Atoi(value); err == nil {
						query[key] = num
					} else if b, err := strconv.ParseBool(value); err == nil {
						query[key] = b
					} else {
						query[key] = value
					}
				}
			}

			if err := schema.Validate(query); err != nil {
				if validationErr, ok := err.(*zod.ValidationError); ok {
					respondWithValidationError(w, validationErr)
				} else {
					respondWithError(w, http.StatusBadRequest, "Query validation failed", []string{err.Error()})
				}
				return
			}

			next(w, r)
		}
	}
}

// HTTP Handlers
func createUserHandler(w http.ResponseWriter, r *http.Request) {
	var userData map[string]interface{}
	json.NewDecoder(r.Body).Decode(&userData)

	// Simulate user creation
	user := map[string]interface{}{
		"id":         123,
		"username":   userData["username"],
		"email":      userData["email"],
		"profile":    userData["profile"],
		"active":     true,
		"created_at": "2024-01-01T12:00:00Z",
	}

	respondWithSuccess(w, http.StatusCreated, "User created successfully", user)
}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	var userData map[string]interface{}
	json.NewDecoder(r.Body).Decode(&userData)

	// Simulate user update
	updatedUser := map[string]interface{}{
		"id":         userData["id"],
		"username":   userData["username"],
		"email":      userData["email"],
		"profile":    userData["profile"],
		"updated_at": "2024-01-01T12:30:00Z",
	}

	respondWithSuccess(w, http.StatusOK, "User updated successfully", updatedUser)
}

func listUsersHandler(w http.ResponseWriter, r *http.Request) {
	// Simulate user listing with pagination
	users := []map[string]interface{}{
		{
			"id":       1,
			"username": "john_doe",
			"email":    "john@example.com",
			"active":   true,
		},
		{
			"id":       2,
			"username": "jane_smith",
			"email":    "jane@example.com",
			"active":   true,
		},
	}

	responseData := map[string]interface{}{
		"users": users,
		"pagination": map[string]interface{}{
			"page":       1,
			"limit":      10,
			"total":      2,
			"totalPages": 1,
		},
	}

	respondWithSuccess(w, http.StatusOK, "Users retrieved successfully", responseData)
}

// Response helpers
func respondWithSuccess(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	response := APIResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta: &Meta{
			Timestamp: "2024-01-01T12:00:00Z",
			Version:   "v1.0.0",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func respondWithError(w http.ResponseWriter, statusCode int, message string, errors []string) {
	response := APIResponse{
		Success: false,
		Message: message,
		Errors:  errors,
		Meta: &Meta{
			Timestamp: "2024-01-01T12:00:00Z",
			Version:   "v1.0.0",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func respondWithValidationError(w http.ResponseWriter, validationErr *zod.ValidationError) {
	var errors []string

	if len(validationErr.Details) > 0 {
		for _, detail := range validationErr.Details {
			errors = append(errors, fmt.Sprintf("%s: %s", detail.Field, detail.Message))
		}
	} else {
		errors = append(errors, fmt.Sprintf("%s: %s", validationErr.Field, validationErr.Message))
	}

	respondWithError(w, http.StatusBadRequest, "Validation failed", errors)
}

func apiMain() {
	fmt.Println("=== Zod-Go API Validation Example ===")

	// Setup routes with validation middleware
	http.HandleFunc("/users",
		QueryValidationMiddleware(userListQuerySchema)(listUsersHandler))

	http.HandleFunc("/users/create",
		ValidationMiddleware(createUserSchema)(createUserHandler))

	http.HandleFunc("/users/update",
		ValidationMiddleware(updateUserSchema)(updateUserHandler))

	// Demonstrate validation examples
	demonstrateAPIValidation()

	fmt.Println("\nStarting server on :8080...")
	fmt.Println("Try these endpoints:")
	fmt.Println("  GET  /users?page=1&limit=10&sort=username&order=desc")
	fmt.Println("  POST /users/create")
	fmt.Println("  PUT  /users/update")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func demonstrateAPIValidation() {
	fmt.Println("\n--- API Request Validation Examples ---")

	// Test valid create user request
	fmt.Println("\n1. Valid create user request:")
	validCreateRequest := map[string]interface{}{
		"username": "john_doe123",
		"email":    "john.doe@example.com",
		"password": "SecurePass123!",
		"profile": map[string]interface{}{
			"firstName": "John",
			"lastName":  "Doe",
			"age":       28,
		},
		"preferences": map[string]interface{}{
			"newsletter": true,
			"theme":      "dark",
			"language":   "en",
		},
	}

	if err := createUserSchema.Validate(validCreateRequest); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Valid create user request passed validation!")
	}

	// Test invalid create user requests
	fmt.Println("\n2. Invalid create user requests:")

	invalidRequests := []map[string]interface{}{
		{
			"username": "jo", // Too short
			"email":    "john.doe@example.com",
			"password": "SecurePass123!",
			"profile": map[string]interface{}{
				"firstName": "John",
				"lastName":  "Doe",
			},
		},
		{
			"username": "john_doe123",
			"email":    "invalid-email", // Invalid email
			"password": "SecurePass123!",
			"profile": map[string]interface{}{
				"firstName": "John",
				"lastName":  "Doe",
			},
		},
		{
			"username": "john_doe123",
			"email":    "john.doe@example.com",
			"password": "weak", // Weak password
			"profile": map[string]interface{}{
				"firstName": "John",
				"lastName":  "Doe",
			},
		},
	}

	for i, request := range invalidRequests {
		fmt.Printf("\nInvalid request #%d:\n", i+1)
		if err := createUserSchema.Validate(request); err != nil {
			if validationErr, ok := err.(*zod.ValidationError); ok {
				fmt.Printf("❌ %s\n", validationErr.ErrorJSON())
			} else {
				fmt.Printf("❌ %v\n", err)
			}
		}
	}

	// Test query parameter validation
	fmt.Println("\n3. Query parameter validation:")

	validQuery := map[string]interface{}{
		"page":   1,
		"limit":  20,
		"sort":   "username",
		"order":  "desc",
		"search": "john",
		"active": true,
	}

	if err := userListQuerySchema.Validate(validQuery); err != nil {
		fmt.Printf("❌ Query validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Valid query parameters passed validation!")
	}

	// Test invalid query
	invalidQuery := map[string]interface{}{
		"page":  0,         // Invalid: must be >= 1
		"limit": 200,       // Invalid: must be <= 100
		"sort":  "invalid", // Invalid sort field
		"order": "random",  // Invalid order
	}

	fmt.Println("\nInvalid query parameters:")
	if err := userListQuerySchema.Validate(invalidQuery); err != nil {
		if validationErr, ok := err.(*zod.ValidationError); ok {
			fmt.Printf("❌ %s\n", validationErr.ErrorJSON())
		} else {
			fmt.Printf("❌ %v\n", err)
		}
	}
}
