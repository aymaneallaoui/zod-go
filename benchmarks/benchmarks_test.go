package benchmarks

import (
	"testing"

	"github.com/aymaneallaoui/zod-Go/zod"
	"github.com/aymaneallaoui/zod-Go/zod/validators"
)

func BenchmarkStringSchema(b *testing.B) {
	schema := validators.String().Min(3).Max(30).Required()

	for i := 0; i < b.N; i++ {
		err := schema.Validate("benchmarktest")
		if err != nil {
			b.Error("Expected validation to pass")
		}
	}
}

func BenchmarkLargeArrayValidation(b *testing.B) {
	elementSchema := validators.String().Min(3).Max(10)
	schema := validators.Array(elementSchema).Min(1000).Max(10000)

	largeArray := make([]interface{}, 10000)
	for i := 0; i < 10000; i++ {
		largeArray[i] = "test"
	}

	for i := 0; i < b.N; i++ {
		err := schema.Validate(largeArray)
		if err != nil {
			b.Error("Validation failed unexpectedly")
		}
	}
}

func BenchmarkNestedObjectValidation(b *testing.B) {
	schema := validators.Object(map[string]zod.Schema{
		"name": validators.String().Min(3).Required(),
		"address": validators.Object(map[string]zod.Schema{
			"street": validators.String().Min(5).Required(),
			"city":   validators.String().Min(3).Required(),
		}).Required(),
	})

	data := map[string]interface{}{
		"name": "John",
		"address": map[string]interface{}{
			"street": "123 Elm St",
			"city":   "Somewhere",
		},
	}

	for i := 0; i < b.N; i++ {
		err := schema.Validate(data)
		if err != nil {
			b.Error("Validation failed unexpectedly")
		}
	}
}
