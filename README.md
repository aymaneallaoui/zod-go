# Zod-Go: Schema Validation Library for Go

Zod-Go is a Go-based validation library inspired by the popular Zod library in TypeScript(before you say anything yes i'm a TS soy dev, and i don't use go std lib cry about it + i only use go cuz it's blue like TS). It allows developers to easily define schemas to validate complex data structures, including strings, numbers, arrays(we know it's sLiCeS), maps, and nested objects.

## Features

- **Schema Definitions**: Validate strings, numbers, booleans, arrays(sLiCeS), and nested objects.
- **Custom Error Handling**: Get detailed validation errors with custom messages.
- **Concurrent Validation**: "Improve" performance for large shit through concurrent validation.
- **Optional Fields and Default Values**: Handle optional fields gracefully and set defaults where necessary.

## Installation

Install the package using `go get`:

```bash
go get github.com/aymaneallaoui/zod-Go/zod
```

## Usage

here's a example of how to use it (it's shit i know but i'm lazy and dumb):

```go
package main

import (
	"fmt"
	"github.com/aymaneallaoui/zod-Go/zod/validators"
)

func main() {
	stringSchema := validators.String().
		Min(3).Max(5).Required().
		WithMessage("minLength", "This string is too short like arch users!").
		WithMessage("maxLength", "This string is too long!")

	err := stringSchema.Validate("ab")
	if err != nil {
		fmt.Println("Validation failed:", err.(*zod.ValidationError).ErrorJSON())
	}
}

```

as you can see (i guess) we validate a string with a minimum length of 3 and a maximum length of 5.
Custom error messages (thanks for a friend for suggesting that) are used for both validation rules.
