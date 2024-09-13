package zod

import (
	"sync"
)

type Schema interface {
	Validate(data interface{}) error
}

func ValidateSchema(schema Schema, data interface{}) error {
	return schema.Validate(data)
}

type Result struct {
	IsValid bool
	Error   error
}

func Validator(schema Schema, data interface{}, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	err := schema.Validate(data)

	var isValid bool
	if err == nil {
		isValid = true
	} else {
		isValid = false
	}

	results <- Result{IsValid: isValid, Error: err}
}

func ValidateConcurrently(schema Schema, dataList []interface{}, workerCount int) []Result {

	results := make(chan Result, len(dataList))
	var wg sync.WaitGroup

	sem := make(chan struct{}, workerCount)

	for _, data := range dataList {
		wg.Add(1)
		go func(data interface{}) {
			sem <- struct{}{}
			defer func() { <-sem }()
			Validator(schema, data, results, &wg)
		}(data)
	}

	wg.Wait()
	close(results)

	var validationResults []Result
	for result := range results {
		validationResults = append(validationResults, result)
	}

	return validationResults
}
