package tests

import (
	"testing"

	"github.com/aymaneallaoui/zod-Go/zod"
	"github.com/aymaneallaoui/zod-Go/zod/validators"
)

func TestStringSchema(t *testing.T) {
	schema := validators.String().Min(3).Max(5).Required()

	err := schema.Validate("hi")
	if err == nil {
		t.Error("Expected validation error for short string")
	}

	err = schema.Validate("hello")
	if err != nil {
		t.Error("Expected validation to pass")
	}

	err = schema.Validate("")
	if err == nil {
		t.Error("Expected validation error for empty string")
	}
}

func TestNumberSchema(t *testing.T) {
	schema := validators.Number().Min(10).Max(20).Required()

	err := schema.Validate(5.0)
	if err == nil {
		t.Error("Expected validation error for number too small")
	}

	err = schema.Validate(15.0)
	if err != nil {
		t.Error("Expected validation to pass")
	}

	err = schema.Validate(nil)
	if err == nil {
		t.Error("Expected validation error for nil input")
	}
}

func TestBoolSchema(t *testing.T) {
	schema := validators.Bool().Required()

	err := schema.Validate(false)
	if err == nil {
		t.Error("Expected validation error for false boolean")
	}

	err = schema.Validate(true)
	if err != nil {
		t.Error("Expected validation to pass")
	}
}

func TestArraySchema(t *testing.T) {
	elementSchema := validators.String().Min(3).Max(5)
	schema := validators.Array(elementSchema).Min(1).Max(3)

	err := schema.Validate([]interface{}{"abc", "def"})
	if err != nil {
		t.Errorf("Expected validation to pass, got error: %v", err)
	}

	err = schema.Validate([]interface{}{"a"})
	if err == nil {
		t.Error("Expected validation error for short string in array")
	}

	err = schema.Validate([]interface{}{"abc", "def", "ghi", "jkl"})
	if err == nil {
		t.Error("Expected validation error for array exceeding max length")
	}
}

func TestMapSchema(t *testing.T) {
	keySchema := validators.String().Min(1)
	valueSchema := validators.Number().Min(10)

	schema := validators.Map(keySchema, valueSchema)

	testMap := map[interface{}]interface{}{
		"a": 15.0,
		"b": 20.0,
	}

	err := schema.Validate(testMap)
	if err != nil {
		t.Error("Expected validation to pass")
	}

	testMap["a"] = 5.0
	err = schema.Validate(testMap)
	if err == nil {
		t.Error("Expected validation error for small number in map")
	}
}

func TestNestedObjectSchema(t *testing.T) {
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

	err := schema.Validate(data)
	if err != nil {
		t.Error("Expected validation to pass")
	}

	data["address"].(map[string]interface{})["street"] = "St"
	err = schema.Validate(data)
	if err == nil {
		t.Error("Expected validation error for short street name")
	}
}
