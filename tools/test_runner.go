// Package main provides a comprehensive test runner and performance analyzer for zod-go
// This tool helps developers understand the performance characteristics and capabilities
// of the enhanced DX API compared to the legacy implementation.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/aymaneallaoui/zod-go/zod/validators"
)

func main() {
	fmt.Println("🧪 zod-go Enhanced DX Test Runner & Performance Analyzer")
	fmt.Println("=====================================================")
	
	runPerformanceAnalysis()
	runMemoryAnalysis()
	runConcurrencyTests()
	runRealWorldScenarios()
	runComparisonTests()
	generatePerformanceReport()
}

// Performance analysis for different API patterns
func runPerformanceAnalysis() {
	fmt.Println("\n📊 Performance Analysis")
	fmt.Println("======================")

	scenarios := []struct {
		name        string
		description string
		setupFunc   func() interface{ Validate(interface{}) error }
		testData    interface{}
	}{
		{
			name:        "SimpleString",
			description: "Basic string validation",
			setupFunc: func() interface{ Validate(interface{}) error } {
				return validators.String().Min(3).Max(50).Required()
			},
			testData: "test_string",
		},
		{
			name:        "ComplexString",
			description: "String with multiple constraints and custom messages",
			setupFunc: func() interface{ Validate(interface{}) error } {
				return validators.String().
					Min(3).WithMinLengthMessage("Too short").
					Max(50).WithMaxLengthMessage("Too long").
					Pattern(`^[a-zA-Z0-9_]+$`).WithPatternMessage("Invalid format").
					Required().WithRequiredMessage("Required field")
			},
			testData: "valid_string_123",
		},
		{
			name:        "EmailValidation",
			description: "Pre-configured email validation",
			setupFunc: func() interface{ Validate(interface{}) error } {
				return validators.Email().WithEmailMessage("Invalid email format")
			},
			testData: "user@example.com",
		},
		{
			name:        "NumberValidation",
			description: "Number with range and type constraints",
			setupFunc: func() interface{ Validate(interface{}) error } {
				return validators.Number().
					Min(0).WithMinMessage("Must be positive").
					Max(100).WithMaxMessage("Too large").
					Integer().WithIntegerMessage("Must be integer").
					Required()
			},
			testData: 42,
		},
		{
			name:        "ArrayValidation",
			description: "Array with element validation",
			setupFunc: func() interface{ Validate(interface{}) error } {
				return validators.Array(validators.String().Min(1)).
					MinItems(1).WithMinItemsMessage("At least one item").
					MaxItems(10).WithMaxItemsMessage("Too many items").
					Required()
			},
			testData: []string{"item1", "item2", "item3"},
		},
		{
			name:        "ComplexObject",
			description: "Nested object with multiple field types",
			setupFunc: func() interface{ Validate(interface{}) error } {
				return validators.Object(map[string]interface{}{
					"name":  validators.String().Min(1).Required(),
					"email": validators.Email(),
					"age":   validators.Number().Min(0).Integer().Optional(),
					"tags":  validators.Array(validators.String()).Optional(),
				}).Required()
			},
			testData: map[string]interface{}{
				"name":  "John Doe",
				"email": "john@example.com",
				"age":   30,
				"tags":  []string{"user", "active"},
			},
		},
	}

	for _, scenario := range scenarios {
		fmt.Printf("\n🔍 Testing: %s - %s\n", scenario.name, scenario.description)
		
		// Warm-up
		schema := scenario.setupFunc()
		for i := 0; i < 100; i++ {
			_ = schema.Validate(scenario.testData)
		}

		// Benchmark
		iterations := 10000
		start := time.Now()
		
		for i := 0; i < iterations; i++ {
			_ = schema.Validate(scenario.testData)
		}
		
		duration := time.Since(start)
		avgTime := duration / time.Duration(iterations)
		opsPerSec := float64(iterations) / duration.Seconds()

		fmt.Printf("   ⏱️  Total time: %v\n", duration)
		fmt.Printf("   📈 Average per operation: %v\n", avgTime)
		fmt.Printf("   🚀 Operations per second: %.0f\n", opsPerSec)
		
		// Memory usage during validation
		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)
		
		for i := 0; i < 1000; i++ {
			_ = schema.Validate(scenario.testData)
		}
		
		runtime.ReadMemStats(&m2)
		memUsed := m2.TotalAlloc - m1.TotalAlloc
		fmt.Printf("   💾 Memory per 1K ops: %d bytes\n", memUsed)
	}
}

// Memory usage analysis
func runMemoryAnalysis() {
	fmt.Println("\n💾 Memory Analysis")
	fmt.Println("==================")

	fmt.Println("\n📋 Schema Creation Memory Usage:")
	
	var m1, m2 runtime.MemStats
	
	// Test schema creation memory usage
	runtime.GC()
	runtime.ReadMemStats(&m1)
	
	schemas := make([]interface{ Validate(interface{}) error }, 1000)
	for i := 0; i < 1000; i++ {
		schemas[i] = validators.String().
			Min(1).
			Max(100).
			Pattern(`^[a-zA-Z0-9_]+$`).
			Required().
			WithMessage(validators.Errors.MinLength(), "Custom message")
	}
	
	runtime.ReadMemStats(&m2)
	schemaMemory := m2.TotalAlloc - m1.TotalAlloc
	fmt.Printf("   🏗️  1000 string schemas: %d bytes (avg: %d bytes/schema)\n", 
		schemaMemory, schemaMemory/1000)

	// Test complex object schema memory
	runtime.GC()
	runtime.ReadMemStats(&m1)
	
	complexSchemas := make([]interface{ Validate(interface{}) error }, 100)
	for i := 0; i < 100; i++ {
		complexSchemas[i] = validators.Object(map[string]interface{}{
			"user": validators.Object(map[string]interface{}{
				"name":    validators.String().Min(1).Required(),
				"email":   validators.Email(),
				"age":     validators.Number().Integer().Optional(),
				"profile": validators.Object(map[string]interface{}{
					"bio":     validators.String().Max(500).Optional(),
					"website": validators.URL().Optional(),
				}).Optional(),
			}).Required(),
			"settings": validators.Object(map[string]interface{}{
				"theme":    validators.String().Optional().Default("light"),
				"language": validators.String().Optional().Default("en"),
			}).Optional(),
		}).Required()
	}
	
	runtime.ReadMemStats(&m2)
	complexMemory := m2.TotalAlloc - m1.TotalAlloc
	fmt.Printf("   🏗️  100 complex schemas: %d bytes (avg: %d bytes/schema)\n", 
		complexMemory, complexMemory/100)

	// Test validation memory usage
	schema := validators.Object(map[string]interface{}{
		"name":  validators.String().Required(),
		"items": validators.Array(validators.String()).Required(),
	}).Required()

	runtime.GC()
	runtime.ReadMemStats(&m1)
	
	for i := 0; i < 1000; i++ {
		testData := map[string]interface{}{
			"name":  fmt.Sprintf("test_%d", i),
			"items": []string{fmt.Sprintf("item_%d", i)},
		}
		_ = schema.Validate(testData)
	}
	
	runtime.ReadMemStats(&m2)
	validationMemory := m2.TotalAlloc - m1.TotalAlloc
	fmt.Printf("   ✅ 1000 validations: %d bytes (avg: %d bytes/validation)\n", 
		validationMemory, validationMemory/1000)
}

// Concurrency and thread safety tests
func runConcurrencyTests() {
	fmt.Println("\n🔄 Concurrency Tests")
	fmt.Println("====================")

	schema := validators.Object(map[string]interface{}{
		"id":   validators.Number().Integer().Required(),
		"name": validators.String().Min(1).Required(),
	}).Required()

	fmt.Println("\n🧵 Testing concurrent validation...")
	
	start := time.Now()
	
	// Run concurrent validations
	done := make(chan bool, 100)
	for i := 0; i < 100; i++ {
		go func(id int) {
			defer func() { done <- true }()
			
			for j := 0; j < 100; j++ {
				testData := map[string]interface{}{
					"id":   id*100 + j,
					"name": fmt.Sprintf("user_%d_%d", id, j),
				}
				
				if err := schema.Validate(testData); err != nil {
					fmt.Printf("❌ Concurrent validation failed: %v\n", err)
					return
				}
			}
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 100; i++ {
		<-done
	}
	
	duration := time.Since(start)
	fmt.Printf("   ✅ 10,000 concurrent validations completed in %v\n", duration)
	fmt.Printf("   🚀 Concurrent throughput: %.0f ops/sec\n", 10000.0/duration.Seconds())
}

// Real-world scenario testing
func runRealWorldScenarios() {
	fmt.Println("\n🌍 Real-World Scenarios")
	fmt.Println("=======================")

	scenarios := []struct {
		name        string
		description string
		setup       func() (interface{ Validate(interface{}) error }, interface{})
	}{
		{
			name:        "UserRegistration",
			description: "Complete user registration form validation",
			setup: func() (interface{ Validate(interface{}) error }, interface{}) {
				schema := validators.Object(map[string]interface{}{
					"username": validators.String().
						Min(3).Max(20).
						Pattern(`^[a-zA-Z0-9_]+$`).
						Required(),
					"email":    validators.Email(),
					"password": validators.String().Min(8).Required(),
					"age":      validators.Number().Min(13).Integer().Required(),
					"terms":    validators.Bool().Required(),
				}).Required()

				data := map[string]interface{}{
					"username": "john_doe",
					"email":    "john@example.com",
					"password": "securepass123",
					"age":      25,
					"terms":    true,
				}

				return schema, data
			},
		},
		{
			name:        "APIResponse",
			description: "REST API response validation",
			setup: func() (interface{ Validate(interface{}) error }, interface{}) {
				schema := validators.Object(map[string]interface{}{
					"success": validators.Bool().Required(),
					"data": validators.Object(map[string]interface{}{
						"users": validators.Array(validators.Object(map[string]interface{}{
							"id":   validators.Number().Integer().Required(),
							"name": validators.String().Required(),
						})).Required(),
					}).Optional(),
					"pagination": validators.Object(map[string]interface{}{
						"page":  validators.Number().Integer().Min(1).Required(),
						"total": validators.Number().Integer().Min(0).Required(),
					}).Optional(),
				}).Required()

				data := map[string]interface{}{
					"success": true,
					"data": map[string]interface{}{
						"users": []map[string]interface{}{
							{"id": 1, "name": "User 1"},
							{"id": 2, "name": "User 2"},
						},
					},
					"pagination": map[string]interface{}{
						"page":  1,
						"total": 50,
					},
				}

				return schema, data
			},
		},
		{
			name:        "ConfigFile",
			description: "Application configuration validation",
			setup: func() (interface{ Validate(interface{}) error }, interface{}) {
				schema := validators.Object(map[string]interface{}{
					"database": validators.Object(map[string]interface{}{
						"host":     validators.String().Required(),
						"port":     validators.Number().Integer().Min(1).Max(65535).Required(),
						"username": validators.String().Required(),
						"password": validators.String().Required(),
					}).Required(),
					"server": validators.Object(map[string]interface{}{
						"port": validators.Number().Integer().Optional().Default(8080.0),
						"host": validators.String().Optional().Default("localhost"),
					}).Optional(),
				}).Required()

				data := map[string]interface{}{
					"database": map[string]interface{}{
						"host":     "localhost",
						"port":     5432,
						"username": "dbuser",
						"password": "dbpass",
					},
					"server": map[string]interface{}{
						"port": 3000,
						"host": "0.0.0.0",
					},
				}

				return schema, data
			},
		},
	}

	for _, scenario := range scenarios {
		fmt.Printf("\n🎯 Testing: %s - %s\n", scenario.name, scenario.description)
		
		schema, testData := scenario.setup()
		
		// Benchmark the scenario
		iterations := 1000
		start := time.Now()
		
		for i := 0; i < iterations; i++ {
			if err := schema.Validate(testData); err != nil {
				fmt.Printf("   ❌ Validation failed: %v\n", err)
				break
			}
		}
		
		duration := time.Since(start)
		fmt.Printf("   ⏱️  %d validations: %v (avg: %v)\n", 
			iterations, duration, duration/time.Duration(iterations))
		fmt.Printf("   🚀 Throughput: %.0f validations/sec\n", 
			float64(iterations)/duration.Seconds())
	}
}

// Comparison tests between different API approaches
func runComparisonTests() {
	fmt.Println("\n⚖️  API Comparison Tests")
	fmt.Println("========================")

	testData := "test_string_123"

	// Test different error key approaches
	errorKeyApproaches := []struct {
		name  string
		setup func() interface{ Validate(interface{}) error }
	}{
		{
			name: "Method-based error keys",
			setup: func() interface{ Validate(interface{}) error } {
				return validators.String().
					Min(5).WithMessage(validators.Errors.MinLength(), "Too short").
					Required().WithMessage(validators.Errors.Required(), "Required")
			},
		},
		{
			name: "Constant error keys",
			setup: func() interface{ Validate(interface{}) error } {
				return validators.String().
					Min(5).WithMessage(validators.ErrMinLength, "Too short").
					Required().WithMessage(validators.ErrRequired, "Required")
			},
		},
		{
			name: "Convenience methods",
			setup: func() interface{ Validate(interface{}) error } {
				return validators.String().
					Min(5).WithMinLengthMessage("Too short").
					Required().WithRequiredMessage("Required")
			},
		},
		{
			name: "Legacy string approach",
			setup: func() interface{ Validate(interface{}) error } {
				return validators.String().
					Min(5).WithMessage("minLength", "Too short").
					Required().WithMessage("required", "Required")
			},
		},
	}

	fmt.Println("\n📊 Error Key Approach Performance:")
	for _, approach := range errorKeyApproaches {
		schema := approach.setup()
		
		iterations := 10000
		start := time.Now()
		
		for i := 0; i < iterations; i++ {
			_ = schema.Validate(testData)
		}
		
		duration := time.Since(start)
		fmt.Printf("   %s: %v (%.0f ops/sec)\n", 
			approach.name, duration, float64(iterations)/duration.Seconds())
	}

	// Test state transition performance
	fmt.Println("\n🔄 State Transition Performance:")
	
	stateApproaches := []struct {
		name  string
		setup func() interface{ Validate(interface{}) error }
	}{
		{
			name: "Direct required",
			setup: func() interface{ Validate(interface{}) error } {
				return validators.String().Min(5).Required()
			},
		},
		{
			name: "Direct optional with default",
			setup: func() interface{ Validate(interface{}) error } {
				return validators.String().Min(5).Optional().Default("default")
			},
		},
		{
			name: "Pre-configured required string",
			setup: func() interface{ Validate(interface{}) error } {
				return validators.RequiredString().Min(5)
			},
		},
		{
			name: "Pre-configured optional string",
			setup: func() interface{ Validate(interface{}) error } {
				return validators.OptionalString().Min(5)
			},
		},
	}

	for _, approach := range stateApproaches {
		iterations := 10000
		start := time.Now()
		
		for i := 0; i < iterations; i++ {
			schema := approach.setup()
			_ = schema.Validate(testData)
		}
		
		duration := time.Since(start)
		fmt.Printf("   %s: %v (%.0f ops/sec)\n", 
			approach.name, duration, float64(iterations)/duration.Seconds())
	}
}

// Generate a comprehensive performance report
func generatePerformanceReport() {
	fmt.Println("\n📋 Performance Report Generation")
	fmt.Println("================================")

	report := struct {
		Timestamp    string                 `json:"timestamp"`
		GoVersion    string                 `json:"go_version"`
		OS           string                 `json:"os"`
		Arch         string                 `json:"arch"`
		NumCPU       int                    `json:"num_cpu"`
		Benchmarks   map[string]interface{} `json:"benchmarks"`
		Summary      map[string]interface{} `json:"summary"`
	}{
		Timestamp: time.Now().Format(time.RFC3339),
		GoVersion: runtime.Version(),
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
		NumCPU:    runtime.NumCPU(),
		Benchmarks: make(map[string]interface{}),
		Summary: map[string]interface{}{
			"api_version": "enhanced_dx_v1",
			"features": []string{
				"type_safe_state_management",
				"smart_error_key_autocompletion", 
				"clean_package_interface",
				"compile_time_safety",
			},
		},
	}

	// Run quick benchmarks for the report
	quickBenchmarks := []struct {
		name     string
		runFunc  func() (time.Duration, int)
	}{
		{
			name: "string_validation",
			runFunc: func() (time.Duration, int) {
				schema := validators.String().Min(1).Max(100).Required()
				iterations := 10000
				start := time.Now()
				for i := 0; i < iterations; i++ {
					_ = schema.Validate("test_string")
				}
				return time.Since(start), iterations
			},
		},
		{
			name: "email_validation",
			runFunc: func() (time.Duration, int) {
				schema := validators.Email()
				iterations := 10000
				start := time.Now()
				for i := 0; i < iterations; i++ {
					_ = schema.Validate("test@example.com")
				}
				return time.Since(start), iterations
			},
		},
		{
			name: "complex_object_validation",
			runFunc: func() (time.Duration, int) {
				schema := validators.Object(map[string]interface{}{
					"name":  validators.String().Required(),
					"email": validators.Email(),
					"age":   validators.Number().Integer().Optional(),
				}).Required()
				
				testData := map[string]interface{}{
					"name":  "John",
					"email": "john@example.com",
					"age":   30,
				}
				
				iterations := 5000
				start := time.Now()
				for i := 0; i < iterations; i++ {
					_ = schema.Validate(testData)
				}
				return time.Since(start), iterations
			},
		},
	}

	for _, benchmark := range quickBenchmarks {
		duration, iterations := benchmark.runFunc()
		opsPerSec := float64(iterations) / duration.Seconds()
		
		report.Benchmarks[benchmark.name] = map[string]interface{}{
			"duration_ns":     duration.Nanoseconds(),
			"iterations":      iterations,
			"ops_per_second":  opsPerSec,
			"avg_duration_ns": duration.Nanoseconds() / int64(iterations),
		}
		
		fmt.Printf("   ✅ %s: %.0f ops/sec\n", benchmark.name, opsPerSec)
	}

	// Save report to file
	reportJSON, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		fmt.Printf("❌ Failed to generate report: %v\n", err)
		return
	}

	filename := fmt.Sprintf("performance_report_%s.json", 
		time.Now().Format("2006-01-02_15-04-05"))
	
	if err := os.WriteFile(filename, reportJSON, 0644); err != nil {
		fmt.Printf("❌ Failed to save report: %v\n", err)
		return
	}

	fmt.Printf("   📄 Report saved to: %s\n", filename)
	fmt.Printf("   📊 Report size: %d bytes\n", len(reportJSON))
}

// Helper function for running Go benchmarks programmatically
func runGoBenchmark(benchFunc func(*testing.B)) testing.BenchmarkResult {
	return testing.Benchmark(benchFunc)
}

// Utility functions for testing
func warmup(fn func()) {
	for i := 0; i < 100; i++ {
		fn()
	}
}

func measureMemory(fn func()) uint64 {
	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)
	fn()
	runtime.ReadMemStats(&m2)
	return m2.TotalAlloc - m1.TotalAlloc
}

func init() {
	// Set GOMAXPROCS to use all available CPUs
	runtime.GOMAXPROCS(runtime.NumCPU())
	
	fmt.Printf("🖥️  System Info: %s %s, %d CPUs, Go %s\n", 
		runtime.GOOS, runtime.GOARCH, runtime.NumCPU(), runtime.Version())
}
