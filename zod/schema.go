package zod

type Schema interface {
	Validate(data interface{}) error
}

type Validator func(data interface{}) error

func ValidateSchema(schema Schema, data interface{}) error {
	return schema.Validate(data)
}
