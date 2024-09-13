package validators

import (
	"sync"

	"github.com/aymaneallaoui/zod-Go/zod"
)

type ArraySchema struct {
	elementSchema zod.Schema
	minLength     int
	maxLength     int
	required      bool
}

func Array(elementSchema zod.Schema) *ArraySchema {
	return &ArraySchema{elementSchema: elementSchema}
}

func (a *ArraySchema) Min(length int) *ArraySchema {
	a.minLength = length
	return a
}

func (a *ArraySchema) Max(length int) *ArraySchema {
	a.maxLength = length
	return a
}

func (a *ArraySchema) Required() *ArraySchema {
	a.required = true
	return a
}

func (a *ArraySchema) Validate(data interface{}) error {

	array, ok := data.([]interface{})
	if !ok {

		switch v := data.(type) {
		case []string:
			array = make([]interface{}, len(v))
			for i := range v {
				array[i] = v[i]
			}
		case []int:
			array = make([]interface{}, len(v))
			for i := range v {
				array[i] = v[i]
			}
		default:
			return zod.NewValidationError("array", data, "invalid type, expected array")
		}
	}

	if a.required && len(array) == 0 {
		return zod.NewValidationError("array", data, "array is required")
	}

	if len(array) < a.minLength {
		return zod.NewValidationError("array", data, "array is too short")
	}

	if len(array) > a.maxLength {
		return zod.NewValidationError("array", data, "array is too long")
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(array))

	for _, item := range array {
		wg.Add(1)
		go func(it interface{}) {
			defer wg.Done()
			if err := a.elementSchema.Validate(it); err != nil {
				errChan <- err
			}
		}(item)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}
