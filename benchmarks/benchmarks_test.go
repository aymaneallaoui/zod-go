package benchmarks

import (
	"testing"

	"github.com/aymaneallaoui/zod-Go/zod"
	"github.com/aymaneallaoui/zod-Go/zod/validators"
)

func BenchmarkStringSchema(b *testing.B) {
	schema := validators.String().Min(3).Max(5).Required()

	for i := 0; i < b.N; i++ {
		err := schema.Validate("test")
		if err != nil {
			b.Errorf("Expected validation to pass, got error: %v", err)
		}
	}
}

func BenchmarkLargeArrayValidation(b *testing.B) {
	elementSchema := validators.String().Min(3).Max(10)
	arraySchema := validators.Array(elementSchema).Min(1000).Max(10000)

	largeArray := make([]interface{}, 10000)
	for i := 0; i < 10000; i++ {
		largeArray[i] = "test"
	}

	for i := 0; i < b.N; i++ {
		err := arraySchema.Validate(largeArray)
		if err != nil {
			b.Errorf("Validation failed unexpectedly: %v", err)
		}
	}
}

func BenchmarkNestedObjectValidation(b *testing.B) {
	schema := validators.Object(map[string]zod.Schema{
		"name": validators.String().Min(3).Required(),
		"age":  validators.Number().Min(18).Max(65),
		"address": validators.Object(map[string]zod.Schema{
			"street": validators.String().Min(5).Max(50).Required(),
			"city":   validators.String().Min(3).Max(30).Required(),
		}).Required(),
	})

	userData := map[string]interface{}{
		"name": "John Doe",
		"age":  30,
		"address": map[string]interface{}{
			"street": "123 Elm St",
			"city":   "New York",
		},
	}

	for i := 0; i < b.N; i++ {
		err := schema.Validate(userData)
		if err != nil {
			b.Errorf("Validation failed unexpectedly: %v", err)
		}
	}
}
