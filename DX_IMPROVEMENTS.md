# 🚀 ZOD-GO Developer Experience Improvements

This document outlines the major improvements made to enhance the developer experience (DX) in zod-go.

## ✨ What's New

### 1. 🔍 **Better Autocompletion for `.WithMessage()`**

**Before:**
```go
schema := validators.String().
    Min(3).
    Max(50).
    WithMessage("minLength", "Username must be at least 3 characters"). // ❌ No autocompletion
    WithMessage("maxLenght", "typo here!") // ❌ Typos possible, no validation
```

**After:**
```go
schema := validators.String().
    Min(3).
    Max(50).
    WithMessage(zod.ValidationTypeMinLength, "Username must be at least 3 characters"). // ✅ Full autocompletion
    WithMessage(zod.ValidationTypeMaxLength, "Username cannot exceed 50 characters")     // ✅ Type-safe, no typos
```

### 2. 🔗 **Improved Method Chaining for Arrays**

**Before:**
```go
// This was possible and confusing:
schema := validators.Array(validators.String()).
    Unique().
    Unique() // ❌ Could call Unique() multiple times
```

**After:**
```go
// Now properly typed and prevents mistakes:
schema := validators.Array(validators.String()).
    Min(1).
    Max(10).
    Unique() // ✅ Returns UniqueArraySchema, preventing double .Unique() calls

// Can still chain other methods after Unique():
finalSchema := schema.
    Required().
    WithMessage(zod.ValidationTypeUnique, "All items must be unique")

// schema.Unique() // ❌ Compile-time error - method doesn't exist on UniqueArraySchema
```

## 📚 Available Validation Type Constants

### String Validation Types
```go
zod.ValidationTypeRequired      // "required"
zod.ValidationTypeMinLength     // "minLength" 
zod.ValidationTypeMaxLength     // "maxLength"
zod.ValidationTypePattern       // "pattern"
zod.ValidationTypeEmail         // "email"
zod.ValidationTypeURL           // "url"
zod.ValidationTypeType          // "type"
zod.ValidationTypeCustom        // "custom"
```

### Number Validation Types  
```go
zod.ValidationTypeMin           // "min"
zod.ValidationTypeMax           // "max"
zod.ValidationTypeInteger       // "integer"
zod.ValidationTypePositive      // "positive"
zod.ValidationTypeNegative      // "negative"
```

### Array Validation Types
```go
zod.ValidationTypeUnique        // "unique"
zod.ValidationTypeNonEmpty      // "nonEmpty"
zod.ValidationTypeElements      // "elements"
```

### Boolean Validation Types
```go
zod.ValidationTypeTrue          // "true"
zod.ValidationTypeFalse         // "false"
```

## 🎯 Practical Examples

### User Registration Schema
```go
userSchema := validators.Object(map[string]zod.Schema{
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
        WithMessage(zod.ValidationTypeEmail, "Please provide a valid email address"),

    "age": validators.Number().
        Integer().
        Min(13).Max(120).
        Required().
        WithMessage(zod.ValidationTypeMin, "Must be at least 13 years old").
        WithMessage(zod.ValidationTypeMax, "Age cannot exceed 120").
        WithMessage(zod.ValidationTypeInteger, "Age must be a whole number"),

    "tags": validators.Array(validators.String().Min(1)).
        Min(1).Max(5).
        Unique(). // ✅ Returns UniqueArraySchema, prevents double .Unique()
        Required().
        WithMessage(zod.ValidationTypeUnique, "All tags must be unique").
        WithMessage(zod.ValidationTypeMinLength, "At least one tag is required"),
})
```

### API Endpoint Validation
```go
func validateCreateUser(w http.ResponseWriter, r *http.Request) {
    schema := validators.Object(map[string]zod.Schema{
        "username": validators.String().
            Min(3).Max(30).
            Required().
            WithMessage(zod.ValidationTypeMinLength, "Username too short").
            WithMessage(zod.ValidationTypeMaxLength, "Username too long"),
            
        "password": validators.String().
            Min(8).
            Pattern(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]`).
            Required().
            WithMessage(zod.ValidationTypeMinLength, "Password must be at least 8 characters").
            WithMessage(zod.ValidationTypePattern, "Password must contain uppercase, lowercase, digit, and special character"),
            
        "interests": validators.Array(validators.String().Min(1)).
            Max(10).
            Unique().
            Optional().
            WithMessage(zod.ValidationTypeMaxLength, "Maximum 10 interests allowed").
            WithMessage(zod.ValidationTypeUnique, "Interests must be unique"),
    })

    var requestData map[string]interface{}
    if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    if err := schema.Validate(requestData); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Process valid data...
}
```

## 🛠️ Migration Guide

### Updating `.WithMessage()` Calls

**Old:**
```go
.WithMessage("minLength", "Too short")
.WithMessage("maxLength", "Too long")  
.WithMessage("required", "This field is required")
.WithMessage("email", "Invalid email")
.WithMessage("unique", "Must be unique")
```

**New:**
```go
.WithMessage(zod.ValidationTypeMinLength, "Too short")
.WithMessage(zod.ValidationTypeMaxLength, "Too long")
.WithMessage(zod.ValidationTypeRequired, "This field is required") 
.WithMessage(zod.ValidationTypeEmail, "Invalid email")
.WithMessage(zod.ValidationTypeUnique, "Must be unique")
```

### Array Schemas with Unique Constraint

**Old:**
```go
// This worked but had potential issues
schema := validators.Array(elementSchema).Unique().Min(1).Max(10)
```

**New:**
```go
// More explicit and prevents double .Unique() calls
schema := validators.Array(elementSchema).
    Min(1).Max(10).
    Unique(). // Returns UniqueArraySchema
    Required()
```

## 🎉 Benefits

1. **🔍 Better IDE Support**: Full autocompletion for validation types
2. **🛡️ Type Safety**: Prevents typos in validation type names
3. **🚀 Improved DX**: More intuitive and discoverable API
4. **🔗 Better Method Chaining**: Prevents invalid method combinations
5. **📝 Self-Documenting**: Validation types are clearly defined constants
6. **🐛 Fewer Bugs**: Compile-time errors instead of runtime issues

## 🔄 Backward Compatibility

The string-based `.WithMessage()` method is still supported internally, but using the new `ValidationTypeConstant` approach is strongly recommended for new code.

## 🚀 Try It Out

Run the improved DX demo:
```bash
go run examples/improved_dx_demo.go
```

This will show you all the new features in action!
