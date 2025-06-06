package validators

// ObjectBuilder represents the initial object builder state.
// From this state, you can configure validation rules and then transition to
// either a required or optional state. This prevents invalid method chaining.
type ObjectBuilder interface {
	// Configuration methods - these return ObjectBuilder to allow chaining
	Strict() ObjectBuilder                                   // Only allow defined keys
	Partial() ObjectBuilder                                  // All keys become optional
	Custom(fn func(map[string]interface{}) error) ObjectBuilder

	// State transition methods - these change the type to prevent invalid chaining
	Required() RequiredObjectBuilder // Transitions to required state
	Optional() OptionalObjectBuilder // Transitions to optional state

	// Error message configuration methods
	WithMessage(validationType, message string) ObjectBuilder
}

// RequiredObjectBuilder represents an object builder in the required state.
// Once in this state, you cannot:
// - Call Required() again (prevents .Required().Required())
// - Set a Default() value (required fields cannot have defaults)
// This enforces logical validation rules at compile time.
type RequiredObjectBuilder interface {
	// Configuration methods - these return RequiredObjectBuilder to maintain state
	Strict() RequiredObjectBuilder
	Partial() RequiredObjectBuilder
	Custom(fn func(map[string]interface{}) error) RequiredObjectBuilder

	// Error message configuration methods
	WithMessage(validationType, message string) RequiredObjectBuilder
	WithRequiredMessage(message string) RequiredObjectBuilder

	// Validation method - final step in the builder chain
	Validate(data interface{}) error
}

// OptionalObjectBuilder represents an object builder in the optional state.
// Once in this state, you cannot:
// - Call Optional() again (prevents .Optional().Optional())
// But you can:
// - Set a Default() value (only optional fields can have defaults)
// This enforces logical validation rules at compile time.
type OptionalObjectBuilder interface {
	// Configuration methods - these return OptionalObjectBuilder to maintain state
	Strict() OptionalObjectBuilder
	Partial() OptionalObjectBuilder
	Custom(fn func(map[string]interface{}) error) OptionalObjectBuilder
	Default(value map[string]interface{}) OptionalObjectBuilder // Only available on optional builders!

	// Error message configuration methods
	WithMessage(validationType, message string) OptionalObjectBuilder

	// Validation method - final step in the builder chain
	Validate(data interface{}) error
}

// BoolBuilder represents the initial bool builder state.
// From this state, you can configure validation rules and then transition to
// either a required or optional state. This prevents invalid method chaining.
type BoolBuilder interface {
	// Configuration methods - these return BoolBuilder to allow chaining
	Custom(fn func(bool) error) BoolBuilder

	// State transition methods - these change the type to prevent invalid chaining
	Required() RequiredBoolBuilder // Transitions to required state
	Optional() OptionalBoolBuilder // Transitions to optional state

	// Error message configuration methods
	WithMessage(validationType, message string) BoolBuilder
}

// RequiredBoolBuilder represents a bool builder in the required state.
// Once in this state, you cannot:
// - Call Required() again (prevents .Required().Required())
// - Set a Default() value (required fields cannot have defaults)
// This enforces logical validation rules at compile time.
type RequiredBoolBuilder interface {
	// Configuration methods - these return RequiredBoolBuilder to maintain state
	Custom(fn func(bool) error) RequiredBoolBuilder

	// Error message configuration methods
	WithMessage(validationType, message string) RequiredBoolBuilder
	WithRequiredMessage(message string) RequiredBoolBuilder

	// Validation method - final step in the builder chain
	Validate(data interface{}) error
}

// OptionalBoolBuilder represents a bool builder in the optional state.
// Once in this state, you cannot:
// - Call Optional() again (prevents .Optional().Optional())
// But you can:
// - Set a Default() value (only optional fields can have defaults)
// This enforces logical validation rules at compile time.
type OptionalBoolBuilder interface {
	// Configuration methods - these return OptionalBoolBuilder to maintain state
	Custom(fn func(bool) error) OptionalBoolBuilder
	Default(value bool) OptionalBoolBuilder // Only available on optional builders!

	// Error message configuration methods
	WithMessage(validationType, message string) OptionalBoolBuilder

	// Validation method - final step in the builder chain
	Validate(data interface{}) error
}
