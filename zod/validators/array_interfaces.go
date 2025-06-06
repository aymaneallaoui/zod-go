package validators

// ArrayBuilder represents the initial array builder state.
// From this state, you can configure validation rules and then transition to
// either a required or optional state. This prevents invalid method chaining.
type ArrayBuilder interface {
	// Configuration methods - these return ArrayBuilder to allow chaining
	MinItems(count int) ArrayBuilder
	MaxItems(count int) ArrayBuilder
	Contains(value interface{}) ArrayBuilder
	Custom(fn func([]interface{}) error) ArrayBuilder

	// State transition methods - these change the type to prevent invalid chaining
	Required() RequiredArrayBuilder // Transitions to required state
	Optional() OptionalArrayBuilder // Transitions to optional state

	// Error message configuration methods
	WithMessage(validationType, message string) ArrayBuilder
	WithMinItemsMessage(message string) ArrayBuilder
	WithMaxItemsMessage(message string) ArrayBuilder
	WithContainsMessage(message string) ArrayBuilder

	// Validation method - final step in the builder chain
	//Validate(data interface{}) error
}

// RequiredArrayBuilder represents an array builder in the required state.
// Once in this state, you cannot:
// - Call Required() again (prevents .Required().Required())
// - Set a Default() value (required fields cannot have defaults)
// This enforces logical validation rules at compile time.
type RequiredArrayBuilder interface {
	// Configuration methods - these return RequiredArrayBuilder to maintain state
	MinItems(count int) RequiredArrayBuilder
	MaxItems(count int) RequiredArrayBuilder
	Contains(value interface{}) RequiredArrayBuilder
	Custom(fn func([]interface{}) error) RequiredArrayBuilder

	// Error message configuration methods
	WithMessage(validationType, message string) RequiredArrayBuilder
	WithMinItemsMessage(message string) RequiredArrayBuilder
	WithMaxItemsMessage(message string) RequiredArrayBuilder
	WithContainsMessage(message string) RequiredArrayBuilder
	WithRequiredMessage(message string) RequiredArrayBuilder

	// Validation method - final step in the builder chain
	Validate(data interface{}) error
}

// OptionalArrayBuilder represents an array builder in the optional state.
// Once in this state, you cannot:
// - Call Optional() again (prevents .Optional().Optional())
// But you can:
// - Set a Default() value (only optional fields can have defaults)
// This enforces logical validation rules at compile time.
type OptionalArrayBuilder interface {
	// Configuration methods - these return OptionalArrayBuilder to maintain state
	MinItems(count int) OptionalArrayBuilder
	MaxItems(count int) OptionalArrayBuilder
	Contains(value interface{}) OptionalArrayBuilder
	Custom(fn func([]interface{}) error) OptionalArrayBuilder
	Default(value []interface{}) OptionalArrayBuilder // Only available on optional builders!

	// Error message configuration methods
	WithMessage(validationType, message string) OptionalArrayBuilder
	WithMinItemsMessage(message string) OptionalArrayBuilder
	WithMaxItemsMessage(message string) OptionalArrayBuilder
	WithContainsMessage(message string) OptionalArrayBuilder

	// Validation method - final step in the builder chain
	Validate(data interface{}) error
}
