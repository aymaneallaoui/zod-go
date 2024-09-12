package validators

import (
	"fmt"
	"sync"

	"github.com/aymaneallaoui/zod-Go/zod"
)

type ObjectSchema struct {
	fields   map[string]zod.Schema
	required bool
	defaults map[string]interface{} // Default values
}

// Object returns a schema for an object (similar to struct in Go)
func Object(fields map[string]zod.Schema) *ObjectSchema {
	return &ObjectSchema{
		fields:   fields,
		defaults: make(map[string]interface{}),
	}
}

func (o *ObjectSchema) Required() *ObjectSchema {
	o.required = true
	return o
}

// Default sets a default value for a field.
func (o *ObjectSchema) Default(field string, value interface{}) *ObjectSchema {
	o.defaults[field] = value
	return o
}

// Validate uses goroutines to parallelize the validation of fields in nested structures.
func (o *ObjectSchema) Validate(data interface{}) error {
	obj, ok := data.(map[string]interface{})
	if !ok {
		return zod.NewValidationError("object", data, "invalid type, expected object")
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(o.fields))

	for key, schema := range o.fields {
		value, exists := obj[key]

		wg.Add(1)
		go func(k string, v interface{}, s zod.Schema) {
			defer wg.Done()

			// If the field doesn't exist and a default value is defined, assign the default.
			if !exists {
				if defaultValue, hasDefault := o.defaults[k]; hasDefault {
					obj[k] = defaultValue
					return
				} else if o.required {
					errChan <- zod.NewValidationError(k, v, "missing required field")
					return
				}
			}

			// Validate the nested schema.
			if err := s.Validate(v); err != nil {
				errChan <- zod.NewValidationError(k, v, err.Error())
			}
		}(key, value, schema)
	}

	// Close the error channel after all goroutines are done.
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// Collect errors.
	var combinedErrors []string
	for err := range errChan {
		if err != nil {
			// If it's a validation error, convert to JSON or string.
			combinedErrors = append(combinedErrors, err.(*zod.ValidationError).ErrorJSON())
		}
	}

	if len(combinedErrors) > 0 {
		return zod.NewValidationError("object", data, fmt.Sprintf("Validation failed: %s", combinedErrors))
	}

	return nil
}
