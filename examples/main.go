package main

import (
	"fmt"

	"github.com/aymaneallaoui/zod-Go/zod"
	"github.com/aymaneallaoui/zod-Go/zod/validators"
)

func main() {

	// userSchema := validators.Object(map[string]zod.Schema{
	// 	"name":     validators.String().Min(3).Max(30).Required(),
	// 	"age":      validators.Number().Min(18).Max(65),
	// 	"email":    validators.String().Min(5).Max(50),
	// 	"isActive": validators.Bool().Required(),
	// 	"address": validators.Object(map[string]zod.Schema{
	// 		"street": validators.String().Min(5).Max(50).Required(),
	// 		"city":   validators.String().Min(3).Max(30).Required(),
	// 	}).Required(),
	// 	"preferences": validators.Array(validators.String()).Min(0).Max(5),
	// }).Default("age", 30)

	// data := map[string]interface{}{
	// 	"name":     "aymane allaoui",
	// 	"email":    "aymane@aallaoui.com",
	// 	"isActive": true,
	// 	"address": map[string]interface{}{
	// 		"street": "marjane 2",
	// 		"city":   "meknes",
	// 	},
	// 	"preferences": []interface{}{"emil", "sms"},
	// }

	// err := userSchema.Validate(data)
	// if err != nil {
	// 	fmt.Println("Validation failed:", err.(*zod.ValidationError).ErrorJSON())
	// } else {
	// 	fmt.Println("Validation succeeded")
	// }

	// fmt.Println("Validated data with defaults:", data)

	keySchema := validators.String().Min(1).Max(10)
	valueSchema := validators.Number().Min(10).Max(100)
	mapSchema := validators.Map(keySchema, valueSchema)

	zbData := map[interface{}]interface{}{
		"key1": 50.0,
		"key2": 70.00,
	}

	err := mapSchema.Validate(zbData)
	if err != nil {
		fmt.Println("Map validation failed:", err.(*zod.ValidationError).ErrorJSON())
	} else {
		fmt.Println("3tini visa w passport")
	}
}
