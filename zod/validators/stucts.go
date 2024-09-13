package validators

import (
	"fmt"
	"sync"

	"github.com/aymaneallaoui/zod-Go/zod"
)

type ObjectSchema struct {
	fields      map[string]zod.Schema
	required    bool
	defaults    map[string]interface{}
	customError map[string]string
}

func Object(fields map[string]zod.Schema) *ObjectSchema {
	return &ObjectSchema{
		fields:      fields,
		defaults:    make(map[string]interface{}),
		customError: make(map[string]string),
	}
}

func (o *ObjectSchema) Required() *ObjectSchema {
	o.required = true
	return o
}

func (o *ObjectSchema) Default(field string, value interface{}) *ObjectSchema {
	o.defaults[field] = value
	return o
}

func (o *ObjectSchema) WithMessage(validationType, message string) *ObjectSchema {
	o.customError[validationType] = message
	return o
}

func (o *ObjectSchema) getErrorMessage(validationType, defaultMessage string) string {
	if msg, exists := o.customError[validationType]; exists {
		return msg
	}
	return defaultMessage
}

func (o *ObjectSchema) Validate(data interface{}) error {
	obj, ok := data.(map[string]interface{})
	if !ok {
		return zod.NewValidationError("object", data, o.getErrorMessage("type", "invalid type, expected object"))
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(o.fields))

	for key, schema := range o.fields {
		value, exists := obj[key]

		wg.Add(1)
		go func(k string, v interface{}, s zod.Schema, fieldExists bool) {
			defer wg.Done()

			if !fieldExists {
				if defaultValue, hasDefault := o.defaults[k]; hasDefault {
					obj[k] = defaultValue
					return
				} else if o.required {
					errChan <- zod.NewValidationError(k, v, o.getErrorMessage("required", fmt.Sprintf("missing required field: %s", k)))
					return
				}
			}

			// If the value is an object, validate recursively and collect errors
			if subObj, isMap := v.(map[string]interface{}); isMap {
				if subErr := s.Validate(subObj); subErr != nil {
					if nestedValidationErr, ok := subErr.(*zod.ValidationError); ok {
						// Collect nested errors and pass them into NewNestedValidationError
						errChan <- zod.NewNestedValidationError(k, v, "Validation failed", nestedValidationErr.Details)
					}
					return
				}
			}

			// Handle regular validation errors
			if err := s.Validate(v); err != nil {
				errChan <- zod.NewValidationError(k, v, err.Error())
			}
		}(key, value, schema, exists)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	// Collect errors from all fields
	var combinedErrors []zod.ValidationError
	for err := range errChan {
		if validationErr, ok := err.(*zod.ValidationError); ok {
			combinedErrors = append(combinedErrors, *validationErr)
		}
	}

	if len(combinedErrors) > 0 {
		// Return structured error with nested details
		return zod.NewNestedValidationError("object", data, "Validation failed", combinedErrors)
	}

	return nil
}
