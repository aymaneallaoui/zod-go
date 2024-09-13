package tests

import (
	"errors"
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
		"name": "John l7way",
		"address": map[string]interface{}{
			"street": "123 Elm St",
			"city":   "zbcity",
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

type MockSchema struct {
	isValid bool
}

type DynamicMockSchema struct {
	validateFunc func(data interface{}) error
}

func (ms DynamicMockSchema) Validate(data interface{}) error {
	return ms.validateFunc(data)
}

func (ms MockSchema) Validate(data interface{}) error {
	if ms.isValid {
		return nil
	}
	return errors.New("invalid data")
}

func TestValidateConcurrently_Success(t *testing.T) {
	schema := MockSchema{isValid: true}
	dataList := []interface{}{"data1", "data2", "data3"}

	results := zod.ValidateConcurrently(schema, dataList, 3)

	for _, result := range results {
		if !result.IsValid || result.Error != nil {
			t.Errorf("Expected valid data, got invalid result: %+v", result)
		}
	}
}

func TestValidateConcurrently_Failure(t *testing.T) {
	schema := MockSchema{isValid: false}
	dataList := []interface{}{"data1", "data2", "data3"}

	results := zod.ValidateConcurrently(schema, dataList, 3)

	for _, result := range results {
		if result.IsValid || result.Error == nil {
			t.Errorf("Expected invalid data, got valid result: %+v", result)
		}
	}
}

func TestValidateConcurrently_EmptyData(t *testing.T) {
	schema := MockSchema{isValid: true}
	dataList := []interface{}{}

	results := zod.ValidateConcurrently(schema, dataList, 3)

	if len(results) != 0 {
		t.Errorf("Expected no results, got %d results", len(results))
	}
}

func TestValidateConcurrently_LargeDataset(t *testing.T) {
	schema := MockSchema{isValid: true}
	largeDataList := make([]interface{}, 10000)
	for i := range largeDataList {
		largeDataList[i] = i
	}

	results := zod.ValidateConcurrently(schema, largeDataList, 100)

	if len(results) != len(largeDataList) {
		t.Errorf("Expected %d results, got %d", len(largeDataList), len(results))
	}

	for _, result := range results {
		if !result.IsValid {
			t.Errorf("Expected all valid data, got invalid result: %+v", result)
		}
	}
}
