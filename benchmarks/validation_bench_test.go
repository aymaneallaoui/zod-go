package benchmarks

import (
	"testing"

	"github.com/aymaneallaoui/zod-go/zod"
	"github.com/aymaneallaoui/zod-go/zod/validators"
)

// String validation benchmarks
func BenchmarkStringValidation_Simple(b *testing.B) {
	schema := validators.String().Required()
	testString := "hello world"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = schema.Validate(testString)
	}
}

func BenchmarkStringValidation_Complex(b *testing.B) {
	schema := validators.String().
		Min(5).
		Max(100).
		Pattern(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).
		Required()
	testString := "user@example.com"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = schema.Validate(testString)
	}
}

func BenchmarkStringValidation_Email(b *testing.B) {
	schema := validators.String().Email().Required()
	testString := "user@example.com"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = schema.Validate(testString)
	}
}

func BenchmarkStringValidation_URL(b *testing.B) {
	schema := validators.String().URL().Required()
	testString := "https://example.com/path?query=value"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = schema.Validate(testString)
	}
}

// Number validation benchmarks
func BenchmarkNumberValidation_Simple(b *testing.B) {
	schema := validators.Number().Required()
	testNumber := 42.5
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = schema.Validate(testNumber)
	}
}

func BenchmarkNumberValidation_Complex(b *testing.B) {
	schema := validators.Number().
		Min(0).
		Max(100).
		Integer().
		MultipleOf(5).
		Required()
	testNumber := 25
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = schema.Validate(testNumber)
	}
}

func BenchmarkNumberValidation_TypeConversion(b *testing.B) {
	schema := validators.Number().Required()
	
	b.Run("int", func(b *testing.B) {
		testValue := 42
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = schema.Validate(testValue)
		}
	})
	
	b.Run("int64", func(b *testing.B) {
		testValue := int64(42)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = schema.Validate(testValue)
		}
	})
	
	b.Run("float32", func(b *testing.B) {
		testValue := float32(42.5)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = schema.Validate(testValue)
		}
	})
	
	b.Run("float64", func(b *testing.B) {
		testValue := float64(42.5)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = schema.Validate(testValue)
		}
	})
}

// Boolean validation benchmarks
func BenchmarkBoolValidation_Simple(b *testing.B) {
	schema := validators.Bool().Required()
	testBool := true
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = schema.Validate(testBool)
	}
}

func BenchmarkBoolValidation_TypeConversion(b *testing.B) {
	schema := validators.Bool()
	
	b.Run("bool", func(b *testing.B) {
		testValue := true
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = schema.Validate(testValue)
		}
	})
	
	b.Run("string", func(b *testing.B) {
		testValue := "true"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = schema.Validate(testValue)
		}
	})
	
	b.Run("int", func(b *testing.B) {
		testValue := 1
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = schema.Validate(testValue)
		}
	})
}

// Array validation benchmarks
func BenchmarkArrayValidation_SmallArray(b *testing.B) {
	elementSchema := validators.String().Required()
	schema := validators.Array(elementSchema).Required()
	testArray := []interface{}{"hello", "world", "test"}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = schema.Validate(testArray)
	}
}

func BenchmarkArrayValidation_LargeArray(b *testing.B) {
	elementSchema := validators.Number().Min(0).Max(1000).Required()
	schema := validators.Array(elementSchema).Required()
	
	// Create large test array
	testArray := make([]interface{}, 1000)
	for i := range testArray {
		testArray[i] = i
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = schema.Validate(testArray)
	}
}

func BenchmarkArrayValidation_UniqueElements(b *testing.B) {
	elementSchema := validators.String().Required()
	schema := validators.Array(elementSchema).Unique().Required()
	
	// Create test array with unique elements
	testArray := make([]interface{}, 100)
	for i := range testArray {
		testArray[i] = string(rune('a' + i%26)) + string(rune('0' + i/26))
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = schema.Validate(testArray)
	}
}

func BenchmarkArrayValidation_NestedArrays(b *testing.B) {
	innerSchema := validators.Array(validators.String().Required()).Min(1)
	outerSchema := validators.Array(innerSchema).Required()
	
	testArray := []interface{}{
		[]interface{}{"a", "b", "c"},
		[]interface{}{"d", "e", "f"},
		[]interface{}{"g", "h", "i"},
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = schema.Validate(testArray)
	}
}

// Object validation benchmarks
func BenchmarkObjectValidation_SimpleObject(b *testing.B) {
	schema := validators.Object(map[string]zod.Schema{
		"name": validators.String().Required(),
		"age":  validators.Number().Required(),
	})
	
	testObject := map[string]interface{}{
		"name": "John Doe",
		"age":  30,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = schema.Validate(testObject)
	}
}

func BenchmarkObjectValidation_ComplexObject(b *testing.B) {
	schema := validators.Object(map[string]zod.Schema{
		"id":       validators.Number().Integer().Positive().Required(),
		"name":     validators.String().Min(1).Max(100).Required(),
		"email":    validators.String().Email().Required(),
		"age":      validators.Number().Min(0).Max(150).Optional(),
		"tags":     validators.Array(validators.String().Required()).Unique().Optional(),
		"active":   validators.Bool().Required(),
	})
	
	testObject := map[string]interface{}{
		"id":     1,
		"name":   "John Doe",
		"email":  "john@example.com",
		"age":    30,
		"tags":   []interface{}{"developer", "golang"},
		"active": true,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = schema.Validate(testObject)
	}
}

func BenchmarkObjectValidation_NestedObject(b *testing.B) {
	addressSchema := validators.Object(map[string]zod.Schema{
		"street": validators.String().Required(),
		"city":   validators.String().Required(),
		"state":  validators.String().Optional(),
	})
	
	userSchema := validators.Object(map[string]zod.Schema{
		"name":    validators.String().Required(),
		"email":   validators.String().Email().Required(),
		"address": addressSchema.Required(),
	})
	
	testObject := map[string]interface{}{
		"name":  "John Doe",
		"email": "john@example.com",
		"address": map[string]interface{}{
			"street": "123 Main St",
			"city":   "New York",
			"state":  "NY",
		},
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = schema.Validate(testObject)
	}
}

func BenchmarkObjectValidation_StrictMode(b *testing.B) {
	schema := validators.Object(map[string]zod.Schema{
		"name": validators.String().Required(),
		"age":  validators.Number().Required(),
	}).Strict()
	
	testObject := map[string]interface{}{
		"name": "John Doe",
		"age":  30,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = schema.Validate(testObject)
	}
}

// Concurrent validation benchmarks
func BenchmarkConcurrentValidation_SmallDataset(b *testing.B) {
	schema := validators.String().Min(5).Max(50).Required()
	
	dataList := []interface{}{
		"hello world",
		"test string",
		"another test",
		"validation test",
		"concurrent test",
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = zod.ValidateConcurrently(schema, dataList, 2)
	}
}

func BenchmarkConcurrentValidation_LargeDataset(b *testing.B) {
	schema := validators.Number().Min(0).Max(1000).Required()
	
	// Create large dataset
	dataList := make([]interface{}, 1000)
	for i := range dataList {
		dataList[i] = i
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = zod.ValidateConcurrently(schema, dataList, 4)
	}
}

func BenchmarkConcurrentValidation_WorkerCount(b *testing.B) {
	schema := validators.Object(map[string]zod.Schema{
		"id":   validators.Number().Required(),
		"name": validators.String().Required(),
	})
	
	// Create test dataset
	dataList := make([]interface{}, 100)
	for i := range dataList {
		dataList[i] = map[string]interface{}{
			"id":   i,
			"name": "test name",
		}
	}
	
	workerCounts := []int{1, 2, 4, 8, 16}
	
	for _, workers := range workerCounts {
		b.Run(string(rune('0'+workers)), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = zod.ValidateConcurrently(schema, dataList, workers)
			}
		})
	}
}

// Error path benchmarks (important for realistic performance)
func BenchmarkValidation_ErrorPath_String(b *testing.B) {
	schema := validators.String().Min(10).Required()
	invalidString := "short" // Will fail validation
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = schema.Validate(invalidString)
	}
}

func BenchmarkValidation_ErrorPath_Number(b *testing.B) {
	schema := validators.Number().Min(100).Required()
	invalidNumber := 50 // Will fail validation
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = schema.Validate(invalidNumber)
	}
}

func BenchmarkValidation_ErrorPath_Object(b *testing.B) {
	schema := validators.Object(map[string]zod.Schema{
		"name":  validators.String().Required(),
		"email": validators.String().Email().Required(),
	})
	
	invalidObject := map[string]interface{}{
		"name":  "John",
		"email": "invalid-email", // Will fail validation
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = schema.Validate(invalidObject)
	}
}

// Memory allocation benchmarks
func BenchmarkValidation_MemoryAllocation_String(b *testing.B) {
	schema := validators.String().Min(5).Max(50).Required()
	testString := "hello world"
	
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = schema.Validate(testString)
	}
}

func BenchmarkValidation_MemoryAllocation_Object(b *testing.B) {
	schema := validators.Object(map[string]zod.Schema{
		"name": validators.String().Required(),
		"age":  validators.Number().Required(),
	})
	
	testObject := map[string]interface{}{
		"name": "John Doe",
		"age":  30,
	}
	
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = schema.Validate(testObject)
	}
}

func BenchmarkValidation_MemoryAllocation_Array(b *testing.B) {
	elementSchema := validators.String().Required()
	schema := validators.Array(elementSchema).Required()
	testArray := []interface{}{"hello", "world", "test"}
	
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = schema.Validate(testArray)
	}
}

// Real-world scenario benchmarks
func BenchmarkRealWorld_UserRegistration(b *testing.B) {
	userSchema := validators.Object(map[string]zod.Schema{
		"username": validators.String().Min(3).Max(30).Pattern(`^[a-zA-Z0-9_]+$`).Required(),
		"email":    validators.String().Email().Required(),
		"password": validators.String().Min(8).Pattern(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)`).Required(),
		"age":      validators.Number().Integer().Min(13).Max(120).Required(),
		"terms":    validators.Bool().True().Required(),
	})
	
	userData := map[string]interface{}{
		"username": "john_doe123",
		"email":    "john.doe@example.com",
		"password": "SecurePass123",
		"age":      25,
		"terms":    true,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = userSchema.Validate(userData)
	}
}

func BenchmarkRealWorld_APIResponse(b *testing.B) {
	responseSchema := validators.Object(map[string]zod.Schema{
		"status":  validators.String().Pattern(`^(success|error)$`).Required(),
		"code":    validators.Number().Integer().Min(100).Max(599).Required(),
		"message": validators.String().Optional(),
		"data": validators.Object(map[string]zod.Schema{
			"id":         validators.Number().Integer().Positive().Required(),
			"name":       validators.String().Required(),
			"created_at": validators.String().Pattern(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`).Required(),
			"tags":       validators.Array(validators.String().Required()).Optional(),
		}).Optional(),
	})
	
	responseData := map[string]interface{}{
		"status": "success",
		"code":   200,
		"data": map[string]interface{}{
			"id":         1,
			"name":       "Test Item",
			"created_at": "2024-01-01T12:00:00Z",
			"tags":       []interface{}{"test", "benchmark"},
		},
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = responseSchema.Validate(responseData)
	}
}
