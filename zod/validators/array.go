package validators

import (
	"fmt"
	"reflect"

	"github.com/aymaneallaoui/zod-go/zod"
)

// ArraySchema represents an array validation schema
type ArraySchema struct {
	elementSchema   zod.Schema
	minLength       int
	maxLength       int
	required        bool
	uniqueElements  bool
	nonEmpty        bool
	customFunc      func([]interface{}) error
	optional        bool
	defaultValue    []interface{}
	customError     map[string]string
}

// UniqueArraySchema represents an array schema that has unique constraint applied
// This prevents calling Unique() multiple times and provides better method chaining
type UniqueArraySchema struct {
	*ArraySchema
}

// Array creates a new array schema with element validation (optional by default)
func Array(elementSchema zod.Schema) *ArraySchema {
	return &ArraySchema{
		elementSchema: elementSchema,
		customError:   make(map[string]string),
		optional:      true, // Default to optional
	}
}

// Min sets the minimum array length requirement
func (a *ArraySchema) Min(length int) *ArraySchema {
	a.minLength = length
	return a
}

// Max sets the maximum array length requirement
func (a *ArraySchema) Max(length int) *ArraySchema {
	a.maxLength = length
	return a
}

// Required marks the array as required
func (a *ArraySchema) Required() *ArraySchema {
	a.required = true
	a.optional = false
	return a
}

// Optional marks the array as optional
func (a *ArraySchema) Optional() *ArraySchema {
	a.optional = true
	a.required = false
	return a
}

// Default sets a default value for optional arrays
func (a *ArraySchema) Default(value []interface{}) *ArraySchema {
	a.defaultValue = value
	a.optional = true
	return a
}

// Unique requires all elements to be unique - returns UniqueArraySchema to prevent double calling
func (a *ArraySchema) Unique() *UniqueArraySchema {
	a.uniqueElements = true
	return &UniqueArraySchema{ArraySchema: a}
}

// NonEmpty requires the array to have at least one element
func (a *ArraySchema) NonEmpty() *ArraySchema {
	a.nonEmpty = true
	return a
}

// Custom adds a custom validation function
func (a *ArraySchema) Custom(fn func([]interface{}) error) *ArraySchema {
	a.customFunc = fn
	return a
}

// WithMessage sets a custom error message for a validation type
// Now uses ValidationTypeConstant for better autocompletion and type safety
func (a *ArraySchema) WithMessage(validationType zod.ValidationTypeConstant, message string) *ArraySchema {
	a.customError[validationType.String()] = message
	return a
}

// Methods for UniqueArraySchema - provide the same interface but prevent calling Unique() again

// Min sets the minimum array length requirement for unique arrays
func (u *UniqueArraySchema) Min(length int) *UniqueArraySchema {
	u.ArraySchema.minLength = length
	return u
}

// Max sets the maximum array length requirement for unique arrays
func (u *UniqueArraySchema) Max(length int) *UniqueArraySchema {
	u.ArraySchema.maxLength = length
	return u
}

// Required marks the unique array as required
func (u *UniqueArraySchema) Required() *UniqueArraySchema {
	u.ArraySchema.required = true
	u.ArraySchema.optional = false
	return u
}

// Optional marks the unique array as optional
func (u *UniqueArraySchema) Optional() *UniqueArraySchema {
	u.ArraySchema.optional = true
	u.ArraySchema.required = false
	return u
}

// NonEmpty requires the unique array to have at least one element
func (u *UniqueArraySchema) NonEmpty() *UniqueArraySchema {
	u.ArraySchema.nonEmpty = true
	return u
}

// Custom adds a custom validation function to unique arrays
func (u *UniqueArraySchema) Custom(fn func([]interface{}) error) *UniqueArraySchema {
	u.ArraySchema.customFunc = fn
	return u
}

// WithMessage sets a custom error message for unique arrays
func (u *UniqueArraySchema) WithMessage(validationType zod.ValidationTypeConstant, message string) *UniqueArraySchema {
	u.ArraySchema.customError[validationType.String()] = message
	return u
}

// Note: UniqueArraySchema does NOT have a Unique() method, preventing double application

// getErrorMessage returns custom or default error message
func (a *ArraySchema) getErrorMessage(validationType, defaultMessage string) string {
	if msg, exists := a.customError[validationType]; exists {
		return msg
	}
	return defaultMessage
}

// convertToSlice converts various slice types to []interface{}
func convertToSlice(data interface{}) ([]interface{}, bool) {
	value := reflect.ValueOf(data)
	if value.Kind() != reflect.Slice && value.Kind() != reflect.Array {
		return nil, false
	}

	length := value.Len()
	result := make([]interface{}, length)
	for i := 0; i < length; i++ {
		result[i] = value.Index(i).Interface()
	}
	return result, true
}

// areElementsUnique checks if all elements in the slice are unique
func areElementsUnique(slice []interface{}) bool {
	seen := make(map[interface{}]bool)
	for _, element := range slice {
		// For complex types, use reflect.DeepEqual approach
		elementKey := getComparableKey(element)
		if seen[elementKey] {
			return false
		}
		seen[elementKey] = true
	}
	return true
}

// getComparableKey creates a comparable key for any type
func getComparableKey(element interface{}) interface{} {
	value := reflect.ValueOf(element)
	switch value.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map, reflect.Struct:
		// For complex types, convert to string representation
		return fmt.Sprintf("%+v", element)
	default:
		return element
	}
}

// validateElements validates each element in the array
func (a *ArraySchema) validateElements(slice []interface{}) error {
	var errors []zod.ValidationError

	for i, element := range slice {
		if err := a.elementSchema.Validate(element); err != nil {
			if validationErr, ok := err.(*zod.ValidationError); ok {
				// Create a new error with array index information
				indexedErr := zod.NewValidationError(
					fmt.Sprintf("[%d]", i),
					element,
					validationErr.Message,
				)
				errors = append(errors, *indexedErr)
			} else {
				indexedErr := zod.NewValidationError(
					fmt.Sprintf("[%d]", i),
					element,
					err.Error(),
				)
				errors = append(errors, *indexedErr)
			}
		}
	}

	if len(errors) > 0 {
		return zod.NewNestedValidationError(
			"array",
			slice,
			a.getErrorMessage(zod.ValidationTypeElements.String(), "array contains invalid elements"),
			errors,
		)
	}

	return nil
}

// Validate performs validation on the provided data
func (a *ArraySchema) Validate(data interface{}) error {
	// Handle nil values
	if data == nil {
		if a.required {
			return zod.NewValidationError("", nil, a.getErrorMessage(zod.ValidationTypeRequired.String(), "array is required"))
		}
		if a.defaultValue != nil {
			return a.Validate(a.defaultValue)
		}
		// If optional (default) or explicitly marked optional, allow nil
		return nil
	}

	// Convert to slice
	slice, ok := convertToSlice(data)
	if !ok {
		return zod.NewValidationError(fmt.Sprintf("%v", data), data,
			a.getErrorMessage(zod.ValidationTypeType.String(), "invalid type, expected array or slice"))
	}

	// Length validations
	length := len(slice)

	if a.nonEmpty && length == 0 {
		return zod.NewValidationError("[]", slice,
			a.getErrorMessage(zod.ValidationTypeNonEmpty.String(), "array cannot be empty"))
	}

	if a.minLength > 0 && length < a.minLength {
		return zod.NewValidationError(fmt.Sprintf("length:%d", length), slice,
			a.getErrorMessage(zod.ValidationTypeMinLength.String(), fmt.Sprintf("array is too short, minimum length is %d", a.minLength)))
	}

	if a.maxLength > 0 && length > a.maxLength {
		return zod.NewValidationError(fmt.Sprintf("length:%d", length), slice,
			a.getErrorMessage(zod.ValidationTypeMaxLength.String(), fmt.Sprintf("array is too long, maximum length is %d", a.maxLength)))
	}

	// Unique elements validation
	if a.uniqueElements && !areElementsUnique(slice) {
		return zod.NewValidationError("array", slice,
			a.getErrorMessage(zod.ValidationTypeUnique.String(), "array elements must be unique"))
	}

	// Element validation
	if a.elementSchema != nil {
		if err := a.validateElements(slice); err != nil {
			return err
		}
	}

	// Custom validation
	if a.customFunc != nil {
		if err := a.customFunc(slice); err != nil {
			return err
		}
	}

	return nil
}

// Validate for UniqueArraySchema delegates to the underlying ArraySchema
func (u *UniqueArraySchema) Validate(data interface{}) error {
	return u.ArraySchema.Validate(data)
}

// Helper functions for common array schemas

// StringArray creates an array schema for string elements
func StringArray() *ArraySchema {
	return Array(String().Required())
}

// NumberArray creates an array schema for number elements
func NumberArray() *ArraySchema {
	return Array(Number().Required())
}

// IntegerArray creates an array schema for integer elements
func IntegerArray() *ArraySchema {
	return Array(Number().Integer().Required())
}

// UniqueStringArray creates an array schema for unique string elements
func UniqueStringArray() *UniqueArraySchema {
	return Array(String().Required()).Unique()
}

// NonEmptyArray creates an array schema that cannot be empty
func NonEmptyArray(elementSchema zod.Schema) *ArraySchema {
	return Array(elementSchema).NonEmpty()
}

// FixedLengthArray creates an array schema with exact length requirement
func FixedLengthArray(length int, elementSchema zod.Schema) *ArraySchema {
	return Array(elementSchema).Min(length).Max(length)
}

// OptionalArray creates an optional array schema
func OptionalArray(elementSchema zod.Schema) *ArraySchema {
	return Array(elementSchema).Optional()
}

// BoundedArray creates an array schema with min/max length constraints
func BoundedArray(min, max int, elementSchema zod.Schema) *ArraySchema {
	return Array(elementSchema).Min(min).Max(max)
}
