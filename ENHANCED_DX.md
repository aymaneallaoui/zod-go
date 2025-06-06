# Enhanced Developer Experience (DX) for zod-go

This branch introduces a **complete overhaul** of the zod-go API to provide world-class developer experience with:

- 🎯 **Clean Package Interface** - Only essential functions visible when typing `validators.`
- 🔒 **Type-Safe State Management** - Prevents invalid method chaining at compile time
- ✨ **Smart Autocompletion** - Error keys available with IDE autocompletion
- 🚀 **Logical API Flow** - Clear progression from configuration → state → validation

## 🌟 Key Improvements

### Before vs After

| Problem (Before) | Solution (After) |
|------------------|------------------|
| 50+ exports when typing `validators.` | Only 10 main functions visible |
| Can call `.Required().Required()` | Compile-time error prevention |
| Manual error key strings | Smart autocompletion for error keys |
| Can set `.Default()` on required fields | Type system prevents this |
| Cluttered interface with internals | Clean, focused API surface |

### Type-Safe State Management

```go
// ✅ Valid progressions
username := validators.String().Min(3).Required()  // ✅ Can transition to required
optional := validators.String().Min(3).Optional().Default("test")  // ✅ Can set default on optional

// ❌ These WON'T compile (type safety!)
// username.Required()      // ❌ Can't call Required() twice!
// username.Default("x")    // ❌ Required fields can't have defaults!
// optional.Optional()      // ❌ Can't call Optional() twice!
```

### Smart Error Key Autocompletion

```go
// Method-based (cleanest autocompletion)
schema := validators.String().
    Min(3).
    Required().
    WithMessage(validators.Errors.MinLength(), "Too short").  // 🎯 IDE shows all options!
    WithMessage(validators.Errors.Required(), "Required field")

// Constants (shorter syntax)
schema := validators.String().
    Min(3).
    Required().
    WithMessage(validators.ErrMinLength, "Too short").  // 🎯 Type validators.Err and see all!
    WithMessage(validators.ErrRequired, "Required field")

// Convenience methods (most fluent)
schema := validators.String().
    Min(3).WithMinLengthMessage("Too short").
    Required().WithRequiredMessage("Required field")
```

### Clean Package Interface

When you type `validators.` you now see **ONLY**:

```go
validators.String()        // ✅ Main builders
validators.Number()        // ✅ 
validators.Array()         // ✅
validators.Object()        // ✅
validators.Bool()          // ✅

validators.Email()         // ✅ Convenience builders
validators.URL()           // ✅
validators.OptionalString() // ✅
validators.RequiredString() // ✅

validators.Errors          // ✅ Error keys
validators.ErrRequired     // ✅ Error constants
```

**No more internal functions cluttering the interface!**

## 🚀 Quick Start

```go
package main

import "github.com/aymaneallaoui/zod-go/zod/validators"

func main() {
    // Type-safe validation with smart autocompletion
    userSchema := validators.Object(map[string]interface{}{
        "name": validators.String().
            Min(1).WithMinLengthMessage("Name is required").
            Max(100).WithMaxLengthMessage("Name too long").
            Required(),
            
        "email": validators.Email().  // Pre-configured required email
            WithEmailMessage("Please enter a valid email"),
            
        "age": validators.Number().
            Min(0).WithMinMessage("Age must be positive").
            Max(150).WithMaxMessage("Age must be realistic").
            Integer().WithIntegerMessage("Age must be whole number").
            Optional().Default(18.0),  // Only optional can have defaults!
            
        "active": validators.Bool().Optional().Default(true),
    }).Required()

    // Validate data
    userData := map[string]interface{}{
        "name":   "John Doe",
        "email":  "john@example.com", 
        "age":    25,
        "active": true,
    }

    if err := userSchema.Validate(userData); err != nil {
        fmt.Printf("Validation failed: %s\n", err)
    } else {
        fmt.Println("Validation passed!")
    }
}
```

## 📚 Documentation

### String Validation

```go
// Basic string validation
username := validators.String().
    Min(3).
    Max(50).
    Pattern(`^[a-zA-Z0-9_]+$`).
    Required().
    WithMinLengthMessage("Username too short").
    WithMaxLengthMessage("Username too long").
    WithPatternMessage("Username contains invalid characters")

// Pre-configured validators
email := validators.Email()              // Required email
url := validators.URL()                  // Required URL
optional := validators.OptionalString()  // Optional string
required := validators.RequiredString()  // Required string
```

### Number Validation

```go
// Comprehensive number validation
price := validators.Number().
    Min(0).WithMinMessage("Price cannot be negative").
    Max(99999.99).WithMaxMessage("Price too high").
    Positive().WithPositiveMessage("Price must be positive").
    Required().WithRequiredMessage("Price is required")

// Age with default
age := validators.Number().
    Min(0).
    Max(150).
    Integer().
    Optional().Default(18.0)  // Only available on optional!
```

### Array Validation

```go
// Array with element validation
tags := validators.Array(validators.String().Min(1)).
    MinItems(1).WithMinItemsMessage("At least one tag required").
    MaxItems(10).WithMaxItemsMessage("Too many tags").
    Required()

// Optional array with default
categories := validators.Array(validators.String()).
    Optional().Default([]interface{}{"general"})
```

### Object Validation

```go
// Complex object schema
userSchema := validators.Object(map[string]interface{}{
    "profile": validators.Object(map[string]interface{}{
        "firstName": validators.RequiredString().Min(1),
        "lastName":  validators.RequiredString().Min(1),
        "bio":       validators.OptionalString().Max(500),
    }),
    "settings": validators.Object(map[string]interface{}{
        "theme":         validators.String().Optional().Default("light"),
        "notifications": validators.Bool().Optional().Default(true),
    }).Optional(),
}).Required()
```

### Error Key Reference

```go
// All available error keys for autocompletion:

// Common
validators.Errors.Required()    // or validators.ErrRequired
validators.Errors.Type()        // or validators.ErrType
validators.Errors.Custom()      // or validators.ErrCustom

// String-specific
validators.Errors.MinLength()   // or validators.ErrMinLength
validators.Errors.MaxLength()   // or validators.ErrMaxLength
validators.Errors.Pattern()     // or validators.ErrPattern
validators.Errors.Email()       // or validators.ErrEmail
validators.Errors.URL()         // or validators.ErrURL

// Number-specific
validators.Errors.Min()         // or validators.ErrMin
validators.Errors.Max()         // or validators.ErrMax
validators.Errors.Integer()     // or validators.ErrInteger
validators.Errors.Positive()    // or validators.ErrPositive
validators.Errors.Negative()    // or validators.ErrNegative

// Array-specific
validators.Errors.MinItems()    // or validators.ErrMinItems
validators.Errors.MaxItems()    // or validators.ErrMaxItems
validators.Errors.Contains()    // or validators.ErrContains
```

## 🏗️ Architecture

The new architecture uses **interface-based state management**:

```
StringBuilder (initial state)
    ↓ .Required()          ↓ .Optional()
RequiredStringBuilder   OptionalStringBuilder
    ↓ .Validate()           ↓ .Validate()  
    Error                   Error
```

This ensures:
- **Compile-time safety** - Invalid method chains won't compile
- **Logical progression** - Clear path from configuration to validation
- **State isolation** - Required and optional have different capabilities

## 🔄 Migration Guide

### Backwards Compatibility

The new API is designed to be backwards compatible. Existing code will continue to work.

### Migrating to Enhanced DX

```go
// Before (still works)
schema := validators.String().Min(3).Max(50).Required().WithMessage("minLength", "Too short")

// After (enhanced DX) - multiple options:

// Option 1: Error key methods
schema := validators.String().Min(3).Max(50).Required().
    WithMessage(validators.Errors.MinLength(), "Too short")

// Option 2: Error key constants  
schema := validators.String().Min(3).Max(50).Required().
    WithMessage(validators.ErrMinLength, "Too short")

// Option 3: Convenience methods
schema := validators.String().Min(3).WithMinLengthMessage("Too short").
    Max(50).Required()
```

## 🧪 Testing

Run the enhanced DX demo:

```bash
cd examples
go run enhanced_dx_demo.go
```

This demonstrates:
- Clean package interface
- Type-safe state management  
- Smart error key autocompletion
- Complex validation patterns
- Compile-time error prevention

## 🎯 Benefits

1. **Better IDE Experience** - Autocompletion shows only relevant options
2. **Compile-Time Safety** - Invalid method chains are caught at compile time
3. **Cleaner Code** - Less cluttered, more intuitive API
4. **Faster Development** - Smart autocompletion speeds up coding
5. **Fewer Bugs** - Type system prevents common mistakes
6. **Better Maintainability** - Clear interface contracts

## 🔮 Future Enhancements

- [ ] Add more pre-configured validators for common patterns
- [ ] Extend error key system for custom validators
- [ ] Add performance optimizations for large schemas
- [ ] Create TypeScript-style union type validators
- [ ] Add schema composition utilities

This enhanced DX provides **TypeScript-level developer experience in Go** while maintaining the performance and type safety Go developers expect.
