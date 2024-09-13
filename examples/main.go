package main

import (
	"fmt"

	"github.com/aymaneallaoui/zod-Go/zod"
	"github.com/aymaneallaoui/zod-Go/zod/validators"
)

func main() {

	stringSchema := validators.String().Min(3).Max(5).Required().
		WithMessage("minLength", "The string is too short! like ur ...").
		WithMessage("maxLength", "The string is too long! Maximum length is 5.").
		WithMessage("required", "A string value is required.")

	err := stringSchema.Validate("ab")
	if err != nil {
		fmt.Println("String validation failed:", err.(*zod.ValidationError).Error())
		fmt.Println("String validation failed (JSON):", err.(*zod.ValidationError).ErrorJSON())
	} else {
		fmt.Println("String validation succeeded")
	}

	numberSchema := validators.Number().Min(18).Max(65).Required().
		WithMessage("min", "Age must be at least 18.").
		WithMessage("max", "Age must be at most 65 ur too old.").
		WithMessage("required", "An age value is required.")

	err = numberSchema.Validate(17)
	if err != nil {
		fmt.Println("Number validation failed:", err.(*zod.ValidationError).Error())
		fmt.Println("Number validation failed (JSON):", err.(*zod.ValidationError).ErrorJSON())
	} else {
		fmt.Println("Number validation succeeded")
	}

	boolSchema := validators.Bool().Required().
		WithMessage("required", "A boolean value is required.")

	err = boolSchema.Validate(false)
	if err != nil {
		fmt.Println("Boolean validation failed:", err.(*zod.ValidationError).Error())
		fmt.Println("Boolean validation failed (JSON):", err.(*zod.ValidationError).ErrorJSON())
	} else {
		fmt.Println("Boolean validation succeeded")
	}

	arraySchema := validators.Array(validators.String().Min(3).Max(5)).
		Min(1).Max(3).Required().
		WithMessage("min", "The array must have at least 1 element.").
		WithMessage("max", "The array must have no more than 3 elements.").
		WithMessage("required", "An array is required.")

	err = arraySchema.Validate([]interface{}{"short", "this is too long like my D", "valid"})
	if err != nil {
		fmt.Println("Array validation failed:", err.(*zod.ValidationError).Error())
		fmt.Println("Array validation failed (JSON):", err.(*zod.ValidationError).ErrorJSON())
	} else {
		fmt.Println("Array validation succeeded")
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
		}).Required().WithMessage("required", "Address is required."),
	})

	userData := map[string]interface{}{
		"name":    "aymane",
		"age":     21,
		"address": map[string]interface{}{"street": "marjan 2", "city": "meknes"},
	}

	err = userSchema.Validate(userData)
	if err != nil {
		fmt.Println("User validation failed:", err.(*zod.ValidationError).Error())
		fmt.Println("User validation failed (JSON):", err.(*zod.ValidationError).ErrorJSON())
	} else {
		fmt.Println("User validation succeeded")
	}
}
