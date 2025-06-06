package validators

// StringBuilder represents the initial string builder state.
// From this state, you can configure validation rules and then transition to
// either a required or optional state. This prevents invalid method chaining.
type StringBuilder interface {
	// Configuration methods - these return StringBuilder to allow chaining
	Min(length int) StringBuilder
	Max(length int) StringBuilder
	Pattern(pattern string) StringBuilder
	Email() StringBuilder
	URL() StringBuilder
	Custom(fn func(string) error) StringBuilder

	// State transition methods - these change the type to prevent invalid chaining
	Required() RequiredStringBuilder // Transitions to required state
	Optional() OptionalStringBuilder // Transitions to optional state

	// Error message configuration methods
	WithMessage(validationType, message string) StringBuilder
	WithMinLengthMessage(message string) StringBuilder
	WithMaxLengthMessage(message string) StringBuilder
	WithPatternMessage(message string) StringBuilder
	WithEmailMessage(message string) StringBuilder
	WithURLMessage(message string) StringBuilder
}

// RequiredStringBuilder represents a string builder in the required state.
// Once in this state, you cannot:
// - Call Required() again (prevents .Required().Required())
// - Set a Default() value (required fields cannot have defaults)
// This enforces logical validation rules at compile time.
type RequiredStringBuilder interface {
	// Configuration methods - these return RequiredStringBuilder to maintain state
	Min(length int) RequiredStringBuilder
	Max(length int) RequiredStringBuilder
	Pattern(pattern string) RequiredStringBuilder
	Email() RequiredStringBuilder
	URL() RequiredStringBuilder
	Custom(fn func(string) error) RequiredStringBuilder

	// Error message configuration methods
	WithMessage(validationType, message string) RequiredStringBuilder
	WithMinLengthMessage(message string) RequiredStringBuilder
	WithMaxLengthMessage(message string) RequiredStringBuilder
	WithPatternMessage(message string) RequiredStringBuilder
	WithEmailMessage(message string) RequiredStringBuilder
	WithURLMessage(message string) RequiredStringBuilder
	WithRequiredMessage(message string) RequiredStringBuilder

	// Validation method - final step in the builder chain
	Validate(data interface{}) error
}

// OptionalStringBuilder represents a string builder in the optional state.
// Once in this state, you cannot:
// - Call Optional() again (prevents .Optional().Optional())
// But you can:
// - Set a Default() value (only optional fields can have defaults)
// This enforces logical validation rules at compile time.
type OptionalStringBuilder interface {
	// Configuration methods - these return OptionalStringBuilder to maintain state
	Min(length int) OptionalStringBuilder
	Max(length int) OptionalStringBuilder
	Pattern(pattern string) OptionalStringBuilder
	Email() OptionalStringBuilder
	URL() OptionalStringBuilder
	Custom(fn func(string) error) OptionalStringBuilder
	Default(value string) OptionalStringBuilder // Only available on optional builders!

	// Error message configuration methods
	WithMessage(validationType, message string) OptionalStringBuilder
	WithMinLengthMessage(message string) OptionalStringBuilder
	WithMaxLengthMessage(message string) OptionalStringBuilder
	WithPatternMessage(message string) OptionalStringBuilder
	WithEmailMessage(message string) OptionalStringBuilder
	WithURLMessage(message string) OptionalStringBuilder

	// Validation method - final step in the builder chain
	Validate(data interface{}) error
}
