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
	customError   map[string]string
}

func Array(elementSchema zod.Schema) *ArraySchema {
	return &ArraySchema{
		elementSchema: elementSchema,
		customError:   make(map[string]string),
	}
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

func (a *ArraySchema) WithMessage(validationType, message string) *ArraySchema {
	a.customError[validationType] = message
	return a
}

func (a *ArraySchema) getErrorMessage(validationType, defaultMessage string) string {
	if msg, exists := a.customError[validationType]; exists {
		return msg
	}
	return defaultMessage
}

func (a *ArraySchema) Validate(data interface{}) error {
	array, ok := data.([]interface{})
	if !ok {
		return zod.NewValidationError("array", data, a.getErrorMessage("type", "invalid type, expected array"))
	}

	if a.required && len(array) == 0 {
		return zod.NewValidationError("array", data, a.getErrorMessage("required", "array is required"))
	}

	if len(array) < a.minLength {
		return zod.NewValidationError("array", data, a.getErrorMessage("min", "array is too short"))
	}

	if len(array) > a.maxLength {
		return zod.NewValidationError("array", data, a.getErrorMessage("max", "array is too long"))
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
