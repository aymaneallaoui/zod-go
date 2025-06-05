package zod

// ValidationTypeConstant represents a validation type for better autocompletion and type safety
type ValidationTypeConstant string

// String validation type constants
const (
	// StringValidationTypes - validation types for string schemas
	ValidationTypeRequired   ValidationTypeConstant = "required"
	ValidationTypeMinLength  ValidationTypeConstant = "minLength"
	ValidationTypeMaxLength  ValidationTypeConstant = "maxLength"
	ValidationTypePattern    ValidationTypeConstant = "pattern"
	ValidationTypeEmail      ValidationTypeConstant = "email"
	ValidationTypeURL        ValidationTypeConstant = "url"
	ValidationTypeType       ValidationTypeConstant = "type"
	ValidationTypeCustom     ValidationTypeConstant = "custom"
)

// Number validation type constants
const (
	ValidationTypeMin      ValidationTypeConstant = "min"
	ValidationTypeMax      ValidationTypeConstant = "max"
	ValidationTypeInteger  ValidationTypeConstant = "integer"
	ValidationTypePositive ValidationTypeConstant = "positive"
	ValidationTypeNegative ValidationTypeConstant = "negative"
)

// Array validation type constants
const (
	ValidationTypeUnique     ValidationTypeConstant = "unique"
	ValidationTypeNonEmpty   ValidationTypeConstant = "nonEmpty"
	ValidationTypeElements   ValidationTypeConstant = "elements"
)

// Boolean validation type constants
const (
	ValidationTypeTrue  ValidationTypeConstant = "true"
	ValidationTypeFalse ValidationTypeConstant = "false"
)

// Object validation type constants
const (
	ValidationTypeStrict    ValidationTypeConstant = "strict"
	ValidationTypeKeys      ValidationTypeConstant = "keys"
	ValidationTypeShape     ValidationTypeConstant = "shape"
	ValidationTypeUnknown   ValidationTypeConstant = "unknown"
)

// String returns the string representation of the validation type
func (v ValidationTypeConstant) String() string {
	return string(v)
}
