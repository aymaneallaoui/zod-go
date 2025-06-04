package validators

import (
	"fmt"
	"reflect"

	"github.com/aymaneallaoui/zod-go/zod"
)

// ObjectSchema represents an object validation schema
type ObjectSchema struct {
	fields         map[string]zod.Schema
	required       bool
	strict         bool
	allowUnknown   bool
	customFunc     func(map[string]interface{}) error
	optional       bool
	defaultValue   map[string]interface{}
	customError    map[string]string
}

// Object creates a new object schema with field definitions (optional by default)
func Object(fields map[string]zod.Schema) *ObjectSchema {
	return &ObjectSchema{
		fields:      fields,
		customError: make(map[string]string),
		optional:    true, // Default to optional
	}
}

// Required marks the object as required
func (o *ObjectSchema) Required() *ObjectSchema {
	o.required = true
	o.optional = false
	return o
}

// Optional marks the object as optional
func (o *ObjectSchema) Optional() *ObjectSchema {
	o.optional = true
	o.required = false
	return o
}

// Default sets a default value for optional objects
func (o *ObjectSchema) Default(value map[string]interface{}) *ObjectSchema {
	o.defaultValue = value
	o.optional = true
	return o
}

// Strict enables strict mode - rejects unknown properties
func (o *ObjectSchema) Strict() *ObjectSchema {
	o.strict = true
	o.allowUnknown = false
	return o
}

// AllowUnknown allows unknown properties (opposite of strict)
func (o *ObjectSchema) AllowUnknown() *ObjectSchema {
	o.allowUnknown = true
	o.strict = false
	return o
}

// Custom adds a custom validation function
func (o *ObjectSchema) Custom(fn func(map[string]interface{}) error) *ObjectSchema {
	o.customFunc = fn
	return o
}

// WithMessage sets a custom error message for a validation type
func (o *ObjectSchema) WithMessage(validationType, message string) *ObjectSchema {
	o.customError[validationType] = message
	return o
}

// getErrorMessage returns custom or default error message
func (o *ObjectSchema) getErrorMessage(validationType, defaultMessage string) string {
	if msg, exists := o.customError[validationType]; exists {
		return msg
	}
	return defaultMessage
}

// convertToMap converts various map types to map[string]interface{}
func convertToMap(data interface{}) (map[string]interface{}, bool) {
	value := reflect.ValueOf(data)
	if value.Kind() != reflect.Map {
		return nil, false
	}

	result := make(map[string]interface{})
	for _, key := range value.MapKeys() {
		keyStr, ok := key.Interface().(string)
		if !ok {
			return nil, false
		}
		result[keyStr] = value.MapIndex(key).Interface()
	}
	return result, true
}

// isFieldRequired checks if a field is required based on its schema
func isFieldRequired(schema zod.Schema) bool {
	// This is a simplified check - in a real implementation, you might want
	// to add a method to the Schema interface to check if it's required
	if stringSchema, ok := schema.(*StringSchema); ok {
		return stringSchema.required
	}
	if numberSchema, ok := schema.(*NumberSchema); ok {
		return numberSchema.required
	}
	if boolSchema, ok := schema.(*BoolSchema); ok {
		return boolSchema.required
	}
	if arraySchema, ok := schema.(*ArraySchema); ok {
		return arraySchema.required
	}
	if objectSchema, ok := schema.(*ObjectSchema); ok {
		return objectSchema.required
	}
	// Default to optional (not required) - changed from true to false
	return false
}

// validateFields validates each field in the object
func (o *ObjectSchema) validateFields(obj map[string]interface{}) error {
	var errors []zod.ValidationError

	// Validate defined fields
	for fieldName, fieldSchema := range o.fields {
		value, exists := obj[fieldName]
		
		if !exists {
			// Check if field is required
			if isFieldRequired(fieldSchema) {
				errors = append(errors, *zod.NewValidationError(
					fieldName,
					nil,
					fmt.Sprintf("field '%s' is required", fieldName),
				))
			}
			continue
		}

		if err := fieldSchema.Validate(value); err != nil {
			if validationErr, ok := err.(*zod.ValidationError); ok {
				// Preserve field name in nested error
				fieldErr := zod.NewValidationError(
					fieldName,
					value,
					validationErr.Message,
				)
				if len(validationErr.Details) > 0 {
					fieldErr.Details = validationErr.Details
				}
				errors = append(errors, *fieldErr)
			} else {
				errors = append(errors, *zod.NewValidationError(
					fieldName,
					value,
					err.Error(),
				))
			}
		}
	}

	// Check for unknown fields in strict mode
	if o.strict {
		for fieldName := range obj {
			if _, exists := o.fields[fieldName]; !exists {
				errors = append(errors, *zod.NewValidationError(
					fieldName,
					obj[fieldName],
					fmt.Sprintf("unknown field '%s' is not allowed", fieldName),
				))
			}
		}
	}

	if len(errors) > 0 {
		return zod.NewNestedValidationError(
			"object",
			obj,
			o.getErrorMessage("fields", "object contains invalid fields"),
			errors,
		)
	}

	return nil
}

// Validate performs validation on the provided data
func (o *ObjectSchema) Validate(data interface{}) error {
	// Handle nil values
	if data == nil {
		if o.required {
			return zod.NewValidationError("", nil, o.getErrorMessage("required", "object is required"))
		}
		if o.defaultValue != nil {
			return o.Validate(o.defaultValue)
		}
		// If optional (default) or explicitly marked optional, allow nil
		return nil
	}

	// Convert to map
	obj, ok := convertToMap(data)
	if !ok {
		return zod.NewValidationError(fmt.Sprintf("%v", data), data,
			o.getErrorMessage("type", "invalid type, expected object"))
	}

	// Validate fields
	if err := o.validateFields(obj); err != nil {
		return err
	}

	// Custom validation
	if o.customFunc != nil {
		if err := o.customFunc(obj); err != nil {
			return err
		}
	}

	return nil
}

// AddField adds a new field to the object schema
func (o *ObjectSchema) AddField(name string, schema zod.Schema) *ObjectSchema {
	if o.fields == nil {
		o.fields = make(map[string]zod.Schema)
	}
	newFields := make(map[string]zod.Schema)
	for k, v := range o.fields {
		newFields[k] = v
	}
	newFields[name] = schema
	
	return &ObjectSchema{
		fields:       newFields,
		required:     o.required,
		strict:       o.strict,
		allowUnknown: o.allowUnknown,
		customFunc:   o.customFunc,
		optional:     o.optional,
		defaultValue: o.defaultValue,
		customError:  o.customError,
	}
}

// RemoveField removes a field from the object schema
func (o *ObjectSchema) RemoveField(name string) *ObjectSchema {
	if o.fields == nil {
		return o
	}
	
	newFields := make(map[string]zod.Schema)
	for k, v := range o.fields {
		if k != name {
			newFields[k] = v
		}
	}
	
	return &ObjectSchema{
		fields:       newFields,
		required:     o.required,
		strict:       o.strict,
		allowUnknown: o.allowUnknown,
		customFunc:   o.customFunc,
		optional:     o.optional,
		defaultValue: o.defaultValue,
		customError:  o.customError,
	}
}

// Extend creates a new object schema that extends this one with additional fields
func (o *ObjectSchema) Extend(additionalFields map[string]zod.Schema) *ObjectSchema {
	newFields := make(map[string]zod.Schema)
	
	// Copy existing fields
	for name, schema := range o.fields {
		newFields[name] = schema
	}
	
	// Add new fields
	for name, schema := range additionalFields {
		newFields[name] = schema
	}
	
	return &ObjectSchema{
		fields:       newFields,
		required:     o.required,
		strict:       o.strict,
		allowUnknown: o.allowUnknown,
		customFunc:   o.customFunc,
		optional:     o.optional,
		defaultValue: o.defaultValue,
		customError:  o.customError,
	}
}

// Pick creates a new object schema with only the specified fields
func (o *ObjectSchema) Pick(fieldNames ...string) *ObjectSchema {
	newFields := make(map[string]zod.Schema)
	
	for _, name := range fieldNames {
		if schema, exists := o.fields[name]; exists {
			newFields[name] = schema
		}
	}
	
	return Object(newFields)
}

// Omit creates a new object schema without the specified fields
func (o *ObjectSchema) Omit(fieldNames ...string) *ObjectSchema {
	omitSet := make(map[string]bool)
	for _, name := range fieldNames {
		omitSet[name] = true
	}
	
	newFields := make(map[string]zod.Schema)
	for name, schema := range o.fields {
		if !omitSet[name] {
			newFields[name] = schema
		}
	}
	
	return Object(newFields)
}

// Helper functions for common object schemas

// StrictObject creates a strict object schema that rejects unknown properties
func StrictObject(fields map[string]zod.Schema) *ObjectSchema {
	return Object(fields).Strict()
}

// OptionalObject creates an optional object schema
func OptionalObject(fields map[string]zod.Schema) *ObjectSchema {
	return Object(fields).Optional()
}

// UserSchema creates a common user object schema
func UserSchema() *ObjectSchema {
	return Object(map[string]zod.Schema{
		"id":    Number().Integer().Positive().Required(),
		"name":  String().Min(1).Max(100).Required(),
		"email": String().Email().Required(),
		"age":   Number().Integer().Min(0).Max(150).Optional(),
	})
}

// AddressSchema creates a common address object schema
func AddressSchema() *ObjectSchema {
	return Object(map[string]zod.Schema{
		"street":   String().Min(1).Required(),
		"city":     String().Min(1).Required(),
		"state":    String().Min(2).Max(2).Optional(),
		"zipCode":  String().Pattern(`^\d{5}(-\d{4})?$`).Optional(),
		"country":  String().Min(2).Default("US"),
	})
}
