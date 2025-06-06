# Zod-Go

[![Go CI](https://github.com/aymaneallaoui/zod-go/actions/workflows/go.yml/badge.svg)](https://github.com/aymaneallaoui/zod-go/actions/workflows/go.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/aymaneallaoui/zod-go.svg)](https://pkg.go.dev/github.com/aymaneallaoui/zod-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/aymaneallaoui/zod-go)](https://goreportcard.com/report/github.com/aymaneallaoui/zod-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A TypeScript-inspired schema validation library for Go. Zod-Go provides a fluent, chainable API for validating complex data structures with detailed error reporting and excellent performance.

## Features

- **Fluent API**: Chain validation rules for readable and maintainable code
- **Rich Data Types**: Support for strings, numbers, booleans, arrays, objects, and maps
- **Detailed Errors**: Comprehensive error reporting with custom messages and nested validation details
- **High Performance**: Optimized validation with concurrent processing support
- **Extensible**: Easy to extend with custom validators and rules
- **Well Documented**: Comprehensive documentation with examples
- **Well Tested**: High test coverage with benchmarks

## Installation

```bash
go get github.com/aymaneallaoui/zod-go
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/aymaneallaoui/zod-go/zod/validators"
)

func main() {
    // String validation
    schema := validators.String().
        Min(3).
        Max(50).
        Required().
        WithMessage("minLength", "Username must be at least 3 characters").
        WithMessage("maxLength", "Username cannot exceed 50 characters")

    if err := schema.Validate("jo"); err != nil {
        fmt.Println("Validation failed:", err)
        // Output: Username must be at least 3 characters
    }

    // Nested object validation
    userSchema := validators.Object(map[string]zod.Schema{
        "name": validators.String().Min(2).Required(),
        "email": validators.String().Email().Required(),
        "age": validators.Number().Min(18).Max(120),
        "address": validators.Object(map[string]zod.Schema{
            "street": validators.String().Required(),
            "city": validators.String().Required(),
            "zipCode": validators.String().Pattern(`^\d{5}$`),
        }).Required(),
    })

    user := map[string]interface{}{
        "name": "John Doe",
        "email": "john@example.com",
        "age": 30,
        "address": map[string]interface{}{
            "street": "123 Main St",
            "city": "New York",
            "zipCode": "10001",
        },
    }

    if err := userSchema.Validate(user); err != nil {
        fmt.Println("User validation failed:", err)
    } else {
        fmt.Println("User data is valid!")
    }
}
```

## Validation Types

### String Validation

```go
schema := validators.String().
    Min(5).                          // Minimum length
    Max(100).                        // Maximum length
    Pattern(`^[a-zA-Z0-9]+$`).      // Regex pattern
    Email().                         // Email format
    URL().                           // URL format
    Required().                      // Non-empty required
    WithMessage("min", "Too short")  // Custom error message
```

### Number Validation

```go
schema := validators.Number().
    Min(0).                          // Minimum value
    Max(100).                        // Maximum value
    Integer().                       // Must be integer
    Positive().                      // Must be positive
    Required()                       // Required field
```

### Boolean Validation

```go
schema := validators.Bool().
    Required().                      // Required field
    True()                          // Must be true
```

### Array Validation

```go
elementSchema := validators.String().Min(1)
schema := validators.Array(elementSchema).
    Min(1).                         // Minimum array length
    Max(10).                        // Maximum array length
    Unique()                        // All elements must be unique
```

### Object Validation

```go
schema := validators.Object(map[string]zod.Schema{
    "name": validators.String().Required(),
    "age": validators.Number().Min(0),
    "tags": validators.Array(validators.String()),
}).Strict()  // Reject unknown properties
```

### Map Validation

```go
keySchema := validators.String().Min(1)
valueSchema := validators.Number().Min(0)
schema := validators.Map(keySchema, valueSchema)
```

## Advanced Features

### Custom Error Messages

```go
schema := validators.String().
    Min(8).
    WithMessage("minLength", "Password must be at least 8 characters").
    Pattern(`[A-Z]`).
    WithMessage("pattern", "Password must contain at least one uppercase letter")
```

### Concurrent Validation

For validating large datasets efficiently:

```go
schemas := []zod.Schema{userSchema, userSchema, userSchema}
data := []interface{}{user1, user2, user3}

results := zod.ValidateConcurrently(schemas, data, 4) // 4 workers
for _, result := range results {
    if !result.IsValid {
        fmt.Printf("Validation error: %v\n", result.Error)
    }
}
```

### Optional Fields and Defaults

```go
schema := validators.Object(map[string]zod.Schema{
    "name": validators.String().Required(),
    "role": validators.String().Default("user"),    // Default value
    "bio": validators.String().Optional(),          // Optional field
})
```

### Custom Validators

```go
// Custom validator function
emailDomainValidator := func(data interface{}) error {
    email, ok := data.(string)
    if !ok {
        return zod.NewValidationError("email", data, "must be a string")
    }
    if !strings.HasSuffix(email, "@company.com") {
        return zod.NewValidationError("email", email, "must be a company email")
    }
    return nil
}

schema := validators.String().
    Email().
    Custom(emailDomainValidator)
```

## Error Handling

Zod-Go provides detailed error information:

```go
err := schema.Validate(invalidData)
if err != nil {
    validationErr := err.(*zod.ValidationError)

    // Get JSON representation
    fmt.Println(validationErr.ErrorJSON())

    // Access error details
    fmt.Printf("Field: %s\n", validationErr.Field)
    fmt.Printf("Message: %s\n", validationErr.Message)
    fmt.Printf("Value: %v\n", validationErr.Value)

    // Handle nested errors
    for _, detail := range validationErr.Details {
        fmt.Printf("Nested error - Field: %s, Message: %s\n",
            detail.Field, detail.Message)
    }
}
```

## Performance

Zod-Go is optimized for performance with:

- Zero-allocation validation paths for simple types
- Concurrent validation for large datasets
- Efficient memory usage with object pooling
- Benchmark results show 10x+ performance improvement over reflection-based validators

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### Development Setup

```bash
# Clone the repository
git clone https://github.com/aymaneallaoui/zod-go.git
cd zod-go

# Install dependencies
go mod tidy

# Run tests
go test ./...

# Run benchmarks
go test -bench=. ./benchmarks

# Run linter
golangci-lint run
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Inspiration

This library is inspired by [Zod](https://github.com/colinhacks/zod), the popular TypeScript schema validation library.
