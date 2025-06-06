# Migration Guide: Upgrading to Enhanced DX API

This guide helps you migrate from the legacy zod-go API to the new enhanced developer experience API with type-safe state management and smart autocompletion.

## 🎯 Migration Overview

### What's Changed

| Aspect | Legacy API | Enhanced API |
|--------|------------|--------------|
| **Package Interface** | 50+ exports visible | Only 10 main functions visible |
| **Method Chaining** | Can call `.Required().Required()` | Compile-time prevention of invalid chains |
| **Error Keys** | Manual strings like `"minLength"` | Smart autocompletion with `validators.Errors.MinLength()` |
| **State Management** | No compile-time safety | Type-safe state transitions |
| **Default Values** | Can set on required fields | Only available on optional fields |

### Benefits of Migrating

- ✅ **Compile-time Safety** - Invalid method chains won't compile
- ✅ **Better IDE Support** - Smart autocompletion for error keys
- ✅ **Cleaner API** - Only essential functions visible
- ✅ **Type Safety** - Prevents logical errors at compile time
- ✅ **Better Performance** - Optimized implementation

## 📋 Step-by-Step Migration

### Step 1: Update Import Patterns

**Before:**
```go
import "github.com/aymaneallaoui/zod-go/zod/validators"

// Creating schemas with struct literals (deprecated)
schema := &validators.StringSchema{}
schema.Min(3).Max(50).Required()
```

**After:**
```go
import "github.com/aymaneallaoui/zod-go/zod/validators"

// Creating schemas with builder functions
schema := validators.String().Min(3).Max(50).Required()
```

### Step 2: Replace Struct Literals with Builder Functions

**Before:**
```go
// Legacy struct literal approach
stringSchema := &validators.StringSchema{}
stringSchema.Min(3).Max(20).Required().WithMessage("minLength", "Too short")

numberSchema := &validators.NumberSchema{}
numberSchema.Min(0).Max(100).Integer().Required()

arraySchema := &validators.ArraySchema{ElementSchema: stringSchema}
arraySchema.MinItems(1).MaxItems(10).Required()
```

**After:**
```go
// Enhanced builder approach
stringSchema := validators.String().
    Min(3).Max(20).
    Required().
    WithMessage(validators.Errors.MinLength(), "Too short")

numberSchema := validators.Number().
    Min(0).Max(100).
    Integer().
    Required()

arraySchema := validators.Array(stringSchema).
    MinItems(1).MaxItems(10).
    Required()
```

### Step 3: Update Error Key Usage

**Before:**
```go
schema := validators.String().
    Min(5).WithMessage("minLength", "Too short").
    Max(50).WithMessage("maxLength", "Too long").
    Required().WithMessage("required", "Field required")
```

**After (Multiple Options):**

**Option 1: Method-based error keys (Recommended)**
```go
schema := validators.String().
    Min(5).WithMessage(validators.Errors.MinLength(), "Too short").
    Max(50).WithMessage(validators.Errors.MaxLength(), "Too long").
    Required().WithMessage(validators.Errors.Required(), "Field required")
```

**Option 2: Constant error keys (Shorter)**
```go
schema := validators.String().
    Min(5).WithMessage(validators.ErrMinLength, "Too short").
    Max(50).WithMessage(validators.ErrMaxLength, "Too long").
    Required().WithMessage(validators.ErrRequired, "Field required")
```

**Option 3: Convenience methods (Most fluent)**
```go
schema := validators.String().
    Min(5).WithMinLengthMessage("Too short").
    Max(50).WithMaxLengthMessage("Too long").
    Required().WithRequiredMessage("Field required")
```

### Step 4: Leverage Type-Safe State Management

**Before (Could cause issues):**
```go
// These patterns were possible but problematic
schema := validators.String().Required().Required() // Redundant
schema := validators.String().Required().Default("value") // Illogical
```

**After (Type-safe):**
```go
// Required fields - cannot have defaults
requiredSchema := validators.String().
    Min(3).
    Required().
    WithRequiredMessage("Username is required")
    // .Default("value") // ❌ Won't compile!

// Optional fields - can have defaults
optionalSchema := validators.String().
    Min(3).
    Optional().
    Default("default_value") // ✅ Only works on optional!
    // .Required() // ❌ Won't compile!
```

### Step 5: Use Pre-configured Validators

**Before:**
```go
emailSchema := validators.String().Email().Required()
urlSchema := validators.String().URL().Required()
```

**After:**
```go
// Pre-configured for common patterns
emailSchema := validators.Email() // Already required
urlSchema := validators.URL()     // Already required

// Or with additional configuration
emailSchema := validators.Email().WithEmailMessage("Invalid email")
```

### Step 6: Update Object and Array Validation

**Before:**
```go
objectSchema := &validators.ObjectSchema{
    Schema: map[string]interface{}{
        "name": validators.String().Required(),
        "age":  validators.Number().Integer().Optional(),
    },
}
objectSchema.Required()
```

**After:**
```go
objectSchema := validators.Object(map[string]interface{}{
    "name": validators.String().Required(),
    "age":  validators.Number().Integer().Optional(),
}).Required()
```

## 🔄 Common Migration Patterns

### Pattern 1: Basic String Validation

**Before:**
```go
schema := &validators.StringSchema{}
schema.Min(3).Max(50).Required().WithMessage("minLength", "Too short")
```

**After:**
```go
schema := validators.String().
    Min(3).WithMinLengthMessage("Too short").
    Max(50).
    Required()
```

### Pattern 2: Complex Object Validation

**Before:**
```go
userSchema := &validators.ObjectSchema{
    Schema: map[string]interface{}{
        "username": validators.String().Min(3).Required(),
        "email":    validators.String().Email().Required(),
        "age":      validators.Number().Integer().Min(0).Optional(),
    },
}
userSchema.Required()
```

**After:**
```go
userSchema := validators.Object(map[string]interface{}{
    "username": validators.String().Min(3).Required(),
    "email":    validators.Email(), // Pre-configured
    "age":      validators.Number().Integer().Min(0).Optional(),
}).Required()
```

### Pattern 3: Array with Element Validation

**Before:**
```go
elementSchema := validators.String().Min(1).Required()
arraySchema := &validators.ArraySchema{ElementSchema: elementSchema}
arraySchema.MinItems(1).MaxItems(10).Required()
```

**After:**
```go
arraySchema := validators.Array(validators.String().Min(1).Required()).
    MinItems(1).WithMinItemsMessage("At least one item required").
    MaxItems(10).WithMaxItemsMessage("Too many items").
    Required()
```

### Pattern 4: Conditional Validation

**Before:**
```go
schema := validators.String().Optional()
if isRequired {
    schema = validators.String().Required()
}
```

**After:**
```go
// Use builder pattern for conditional logic
var schema interface{ Validate(interface{}) error }
if isRequired {
    schema = validators.String().Required()
} else {
    schema = validators.String().Optional().Default("default")
}
```

## 🧪 Testing Your Migration

### 1. Compile-time Verification

After migration, these should NOT compile:

```go
// ❌ These won't compile with the new API
validators.String().Required().Required()
validators.String().Required().Default("value")
validators.String().Optional().Optional()
```

### 2. Runtime Verification

Create a test to verify your schemas work as expected:

```go
func TestMigration(t *testing.T) {
    // Your migrated schema
    schema := validators.String().
        Min(3).WithMinLengthMessage("Too short").
        Required().WithRequiredMessage("Required field")

    // Test valid input
    if err := schema.Validate("valid"); err != nil {
        t.Errorf("Expected valid input to pass: %v", err)
    }

    // Test invalid input
    if err := schema.Validate("hi"); err == nil {
        t.Error("Expected short input to fail")
    }

    // Test custom error message
    err := schema.Validate("hi")
    if err != nil && !strings.Contains(err.Error(), "Too short") {
        t.Errorf("Expected custom error message, got: %s", err.Error())
    }
}
```

### 3. Performance Comparison

```go
func BenchmarkOldVsNew(b *testing.B) {
    // New API
    newSchema := validators.String().Min(3).Max(50).Required()
    
    b.Run("NewAPI", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = newSchema.Validate("test_string")
        }
    })
    
    // The new API should be at least as fast as the old one
}
```

## 📊 Migration Checklist

- [ ] Replace struct literals with builder functions
- [ ] Update error key usage to use autocompletion
- [ ] Leverage type-safe state management
- [ ] Use pre-configured validators where appropriate
- [ ] Test compile-time safety (invalid chains shouldn't compile)
- [ ] Verify runtime behavior matches expectations
- [ ] Update tests to use new API patterns
- [ ] Check performance is maintained or improved

## 🔧 Automated Migration Script

You can create a simple script to help with basic migrations:

```bash
#!/bin/bash
# migrate_zod_go.sh - Basic migration helper

# Replace struct literals with builder functions
sed -i 's/&validators\.StringSchema{}/validators.String()/g' *.go
sed -i 's/&validators\.NumberSchema{}/validators.Number()/g' *.go
sed -i 's/&validators\.ArraySchema{}/validators.Array()/g' *.go
sed -i 's/&validators\.ObjectSchema{}/validators.Object()/g' *.go

# Replace common error key patterns
sed -i 's/"minLength"/validators.ErrMinLength/g' *.go
sed -i 's/"maxLength"/validators.ErrMaxLength/g' *.go
sed -i 's/"required"/validators.ErrRequired/g' *.go
sed -i 's/"email"/validators.ErrEmail/g' *.go

echo "Basic migration completed. Please review and test manually."
```

## 🎯 Best Practices After Migration

### 1. Use Type-Safe Patterns

```go
// ✅ Good: Clear state progression
schema := validators.String().
    Min(3).                    // Configuration
    Max(50).                   // Configuration  
    Required().                // State transition
    WithRequiredMessage("...")  // Final configuration

// ✅ Good: Optional with default
optionalSchema := validators.String().
    Min(3).
    Optional().
    Default("default_value")
```

### 2. Leverage Pre-configured Validators

```go
// ✅ Good: Use pre-configured when appropriate
email := validators.Email()
url := validators.URL()
positiveNumber := validators.PositiveNumber().Required()
```

### 3. Use Consistent Error Key Patterns

```go
// ✅ Good: Consistent use of error key approach
schema := validators.String().
    Min(3).WithMessage(validators.Errors.MinLength(), "Custom message").
    Max(50).WithMessage(validators.Errors.MaxLength(), "Custom message").
    Required().WithMessage(validators.Errors.Required(), "Custom message")
```

### 4. Compose Reusable Schemas

```go
// ✅ Good: Create reusable components
func createNameSchema() validators.RequiredStringBuilder {
    return validators.String().
        Min(1).WithMinLengthMessage("Name cannot be empty").
        Max(100).WithMaxLengthMessage("Name too long").
        Required()
}

func createEmailSchema() validators.RequiredStringBuilder {
    return validators.Email().
        WithEmailMessage("Please provide a valid email address")
}

// Use in larger schemas
userSchema := validators.Object(map[string]interface{}{
    "firstName": createNameSchema(),
    "lastName":  createNameSchema(),
    "email":     createEmailSchema(),
}).Required()
```

## 🔍 Troubleshooting Common Issues

### Issue 1: Compilation Errors

**Problem:** Code doesn't compile after migration
**Solution:** Check for invalid method chains:

```go
// ❌ This won't compile
validators.String().Required().Required()

// ✅ This will compile
validators.String().Required()
```

### Issue 2: Missing Error Messages

**Problem:** Custom error messages not working
**Solution:** Use proper error key format:

```go
// ❌ Old format might not work
.WithMessage("minLength", "message")

// ✅ New format
.WithMessage(validators.Errors.MinLength(), "message")
```

### Issue 3: Type Assertion Issues

**Problem:** Interface type assertions failing
**Solution:** Use the new interfaces:

```go
// ✅ Correct interface usage
var validator interface{ Validate(interface{}) error }
validator = validators.String().Required()
```

## 📚 Additional Resources

- **Enhanced DX Documentation**: `ENHANCED_DX.md`
- **Examples**: `examples/enhanced_dx_demo.go`
- **Real-world Patterns**: `examples/real_world_examples.go`
- **Framework Integration**: `examples/framework_integration.go`
- **Performance Analysis**: `tools/test_runner.go`

## 💬 Getting Help

If you encounter issues during migration:

1. Check the **Enhanced DX Documentation** for complete API reference
2. Run the **test runner** to verify performance characteristics
3. Review the **examples** for common patterns
4. Check that your code compiles (the new API prevents many runtime errors at compile time)

The enhanced API is designed to catch common mistakes at compile time, so compilation errors usually indicate an improvement in code safety rather than a problem with the migration.
