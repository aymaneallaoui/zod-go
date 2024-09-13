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

			if err := s.Validate(v); err != nil {

				if nestedValidationErr, ok := err.(*zod.ValidationError); ok && len(nestedValidationErr.Details) > 0 {
					for _, detail := range nestedValidationErr.Details {
						errChan <- &detail
					}
				} else {

					if validationErr, ok := err.(*zod.ValidationError); ok {
						errChan <- zod.NewValidationError(k, v, validationErr.Message)
					} else {
						errChan <- zod.NewValidationError(k, v, err.Error())
					}
				}
			}
		}(key, value, schema, exists)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	var combinedErrors []zod.ValidationError
	for err := range errChan {
		if validationErr, ok := err.(*zod.ValidationError); ok {
			combinedErrors = append(combinedErrors, *validationErr)
		}
	}

	if len(combinedErrors) > 0 {
		return zod.NewNestedValidationError("", data, "", combinedErrors)
	}

	return nil
}
