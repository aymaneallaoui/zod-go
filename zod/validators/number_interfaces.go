package validators

// NumberBuilder represents the initial number builder state.
// From this state, you can configure validation rules and then transition to
// either a required or optional state. This prevents invalid method chaining.
type NumberBuilder interface {
	// Configuration methods - these return NumberBuilder to allow chaining
	Min(value float64) NumberBuilder
	Max(value float64) NumberBuilder
	Integer() NumberBuilder
	Positive() NumberBuilder
	Negative() NumberBuilder
	Custom(fn func(float64) error) NumberBuilder

	// State transition methods - these change the type to prevent invalid chaining
	Required() RequiredNumberBuilder // Transitions to required state
	Optional() OptionalNumberBuilder // Transitions to optional state

	// Error message configuration methods
	WithMessage(validationType, message string) NumberBuilder
	WithMinMessage(message string) NumberBuilder
	WithMaxMessage(message string) NumberBuilder
	WithIntegerMessage(message string) NumberBuilder
	WithPositiveMessage(message string) NumberBuilder
	WithNegativeMessage(message string) NumberBuilder
}

// RequiredNumberBuilder represents a number builder in the required state.
// Once in this state, you cannot:
// - Call Required() again (prevents .Required().Required())
// - Set a Default() value (required fields cannot have defaults)
// This enforces logical validation rules at compile time.
type RequiredNumberBuilder interface {
	// Configuration methods - these return RequiredNumberBuilder to maintain state
	Min(value float64) RequiredNumberBuilder
	Max(value float64) RequiredNumberBuilder
	Integer() RequiredNumberBuilder
	Positive() RequiredNumberBuilder
	Negative() RequiredNumberBuilder
	Custom(fn func(float64) error) RequiredNumberBuilder

	// Error message configuration methods
	WithMessage(validationType, message string) RequiredNumberBuilder
	WithMinMessage(message string) RequiredNumberBuilder
	WithMaxMessage(message string) RequiredNumberBuilder
	WithIntegerMessage(message string) RequiredNumberBuilder
	WithPositiveMessage(message string) RequiredNumberBuilder
	WithNegativeMessage(message string) RequiredNumberBuilder
	WithRequiredMessage(message string) RequiredNumberBuilder

	// Validation method - final step in the builder chain
	Validate(data interface{}) error
}

// OptionalNumberBuilder represents a number builder in the optional state.
// Once in this state, you cannot:
// - Call Optional() again (prevents .Optional().Optional())
// But you can:
// - Set a Default() value (only optional fields can have defaults)
// This enforces logical validation rules at compile time.
type OptionalNumberBuilder interface {
	// Configuration methods - these return OptionalNumberBuilder to maintain state
	Min(value float64) OptionalNumberBuilder
	Max(value float64) OptionalNumberBuilder
	Integer() OptionalNumberBuilder
	Positive() OptionalNumberBuilder
	Negative() OptionalNumberBuilder
	Custom(fn func(float64) error) OptionalNumberBuilder
	Default(value float64) OptionalNumberBuilder // Only available on optional builders!

	// Error message configuration methods
	WithMessage(validationType, message string) OptionalNumberBuilder
	WithMinMessage(message string) OptionalNumberBuilder
	WithMaxMessage(message string) OptionalNumberBuilder
	WithIntegerMessage(message string) OptionalNumberBuilder
	WithPositiveMessage(message string) OptionalNumberBuilder
	WithNegativeMessage(message string) OptionalNumberBuilder

	// Validation method - final step in the builder chain
	Validate(data interface{}) error
}
