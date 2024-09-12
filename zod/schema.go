package zod

// Schema defines the interface for all schema types.
type Schema interface {
	Validate(data interface{}) error
}

// Validator is an alias for validation functions.
type Validator func(data interface{}) error

// ValidateSchema validates the provided data against the schema.
func ValidateSchema(schema Schema, data interface{}) error {
	return schema.Validate(data)
}
