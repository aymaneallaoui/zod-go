package validators

// String creates a new string validation builder.
// This is the primary entry point for string validation.
func String() StringBuilder {
	return &stringSchema{
		customError: make(map[string]string),
	}
}

// Number creates a new number validation builder.
// This is the primary entry point for number validation.
func Number() NumberBuilder {
	return &numberSchema{
		customError: make(map[string]string),
	}
}

// Array creates a new array validation builder.
// elementSchema defines the validation for each array element.
func Array(elementSchema interface{}) ArrayBuilder {
	return &arraySchema{
		elementSchema: elementSchema,
		customError:   make(map[string]string),
	}
}

// Object creates a new object validation builder.
// schema is a map defining validation rules for each object field.
func Object(schema map[string]interface{}) ObjectBuilder {
	return &objectSchema{
		schema:      schema,
		customError: make(map[string]string),
	}
}

// Bool creates a new boolean validation builder.
// This is the primary entry point for boolean validation.
func Bool() BoolBuilder {
	return &boolSchema{
		customError: make(map[string]string),
	}
}

// Convenience builders - these provide pre-configured common patterns
// These are the secondary entry points that make sense at package level

// Email creates a pre-configured required email string validator.
// Equivalent to String().Email().Required()
func Email() RequiredStringBuilder {
	return String().Email().Required()
}

// URL creates a pre-configured required URL string validator.
// Equivalent to String().URL().Required()
func URL() RequiredStringBuilder {
	return String().URL().Required()
}

// OptionalString creates a pre-configured optional string validator.
// Equivalent to String().Optional()
func OptionalString() OptionalStringBuilder {
	return String().Optional()
}

// RequiredString creates a pre-configured required string validator.
// Equivalent to String().Required()
func RequiredString() RequiredStringBuilder {
	return String().Required()
}

// PositiveNumber creates a pre-configured positive number validator.
// Equivalent to Number().Positive()
func PositiveNumber() NumberBuilder {
	return Number().Positive()
}

// IntegerNumber creates a pre-configured integer number validator.
// Equivalent to Number().Integer()
func IntegerNumber() NumberBuilder {
	return Number().Integer()
}
