package main

import (
	"fmt"

	"github.com/aymaneallaoui/zod-Go/zod"
	"github.com/aymaneallaoui/zod-Go/zod/validators"
)

func main() {

	stringSchema := validators.String().Min(3).Max(5).Required().
		WithMessage("minLength", "This string is too short!").
		WithMessage("maxLength", "This string is too long!")

	err := stringSchema.Validate("ab")
	if err != nil {
		fmt.Println("Validation failed:", err.(*zod.ValidationError).ErrorJSON())
	}

	userSchema := validators.Object(map[string]zod.Schema{
		"name": validators.String().Min(3).Required().
			WithMessage("required", "Name is a required field!").
			WithMessage("minLength", "Name must be at least 3 characters."),
		"age": validators.Number().Min(18).Max(65).
			WithMessage("min", "Age must be at least 18.").
			WithMessage("max", "Age must be no more than 65."),
		"address": validators.Object(map[string]zod.Schema{
			"street": validators.String().Min(5).Max(50).Required().
				WithMessage("required", "Street is required.").
				WithMessage("minLength", "Street must be at least 5 characters."),
			"city": validators.String().Min(3).Max(30).Required().
				WithMessage("required", "City is required."),
		}).Required(),
	})

	userData := map[string]interface{}{
		"name": "Jo",
		"age":  17,
		"address": map[string]interface{}{
			"street": "",
			"city":   "NY",
		},
	}

	err = userSchema.Validate(userData)
	if err != nil {
		fmt.Println("User validation failed:", err.(*zod.ValidationError).ErrorJSON())
	} else {
		fmt.Println("User validation succeeded")
	}
}
