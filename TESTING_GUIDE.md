# Testing and Benchmarking Guide

This guide covers comprehensive testing and performance analysis for the enhanced zod-go API, including benchmarks, edge cases, and real-world scenario testing.

## 🧪 Testing Overview

### Test Structure

The enhanced zod-go API includes multiple levels of testing:

- **Unit Tests**: Basic functionality validation
- **Integration Tests**: Complex schema interactions  
- **Edge Case Tests**: Boundary conditions and error scenarios
- **Performance Tests**: Benchmarks and memory analysis
- **Concurrency Tests**: Thread safety verification
- **Real-world Tests**: Practical usage scenarios

### Test Files Overview

```
zod/validators/
├── enhanced_dx_test.go          # Basic functionality tests
├── comprehensive_test.go        # Edge cases and integration tests
├── benchmarks_test.go          # Performance benchmarks
└── ...

examples/
├── enhanced_dx_demo.go         # Usage demonstrations
├── real_world_examples.go      # Practical scenarios
└── framework_integration.go   # Framework patterns

tools/
└── test_runner.go             # Performance analysis tool
```

## 🚀 Running Tests

### Basic Test Execution

```bash
# Run all tests
go test ./zod/validators/...

# Run tests with verbose output
go test -v ./zod/validators/...

# Run specific test functions
go test -v -run TestTypeSafeStateManagement ./zod/validators/

# Run tests with coverage
go test -cover ./zod/validators/...

# Generate detailed coverage report
go test -coverprofile=coverage.out ./zod/validators/...
go tool cover -html=coverage.out -o coverage.html
```

### Running Benchmarks

```bash
# Run all benchmarks
go test -bench=. ./zod/validators/

# Run specific benchmark patterns
go test -bench=BenchmarkNewAPI ./zod/validators/

# Run benchmarks with memory allocation stats
go test -bench=. -benchmem ./zod/validators/

# Run benchmarks multiple times for stability
go test -bench=. -count=5 ./zod/validators/

# Run benchmarks for a specific duration
go test -bench=. -benchtime=10s ./zod/validators/

# Generate CPU profile during benchmarks
go test -bench=. -cpuprofile=cpu.prof ./zod/validators/

# Generate memory profile
go test -bench=. -memprofile=mem.prof ./zod/validators/
```

### Advanced Testing

```bash
# Test with race detection
go test -race ./zod/validators/...

# Test with different build tags
go test -tags=integration ./zod/validators/...

# Parallel test execution
go test -parallel=4 ./zod/validators/...

# Test with timeout
go test -timeout=30s ./zod/validators/...
```

## 📊 Performance Analysis

### Using the Test Runner Tool

```bash
# Run comprehensive performance analysis
go run tools/test_runner.go

# This will output:
# - Performance metrics for different validation types
# - Memory usage analysis
# - Concurrency testing results
# - Real-world scenario benchmarks
# - Comparison between API approaches
# - Generated performance report (JSON file)
```

### Key Performance Metrics

The test runner measures:

1. **Validation Speed**: Operations per second for different schema types
2. **Memory Usage**: Allocation patterns and garbage collection impact
3. **Schema Creation Cost**: Time and memory for building schemas
4. **Concurrency Performance**: Thread safety and parallel execution
5. **Error Handling Overhead**: Cost of validation failures

### Benchmark Categories

#### 1. Basic Operations
- String validation (simple and complex)
- Number validation with constraints
- Email and URL validation
- Array validation with elements
- Object validation with nested structures

#### 2. API Comparison
- Method-based vs constant error keys
- Different state transition patterns
- Pre-configured vs custom validators

#### 3. Real-world Scenarios
- User registration forms
- API request/response validation
- Configuration file validation
- E-commerce data structures

#### 4. Edge Cases
- Large data structures
- Deep nesting levels
- Unicode and special characters
- Boundary value testing

## 🎯 Test Categories

### 1. Functionality Tests

Located in `enhanced_dx_test.go`:

```go
// Test clean package interface
func TestCleanPackageInterface(t *testing.T)

// Test type-safe state management  
func TestTypeSafeStateManagement(t *testing.T)

// Test different validation types
func TestStringValidation(t *testing.T)
func TestNumberValidation(t *testing.T)
func TestArrayValidation(t *testing.T)
func TestObjectValidation(t *testing.T)
func TestBoolValidation(t *testing.T)

// Test error key autocompletion
func TestErrorKeyAutocompletion(t *testing.T)

// Test complex validation patterns
func TestComplexValidationPatterns(t *testing.T)
```

### 2. Edge Case Tests

Located in `comprehensive_test.go`:

```go
// Boundary conditions and extreme values
func TestEdgeCases(t *testing.T)

// Thread safety and concurrent access
func TestConcurrencyAndThreadSafety(t *testing.T)

// Custom error message handling
func TestErrorMessageCustomization(t *testing.T)

// Complex nested data structures
func TestComplexDataTypes(t *testing.T)

// Type conversion and validation
func TestTypeConversions(t *testing.T)

// Schema composition patterns
func TestSchemaComposition(t *testing.T)

// Memory efficiency and performance
func TestMemoryAndPerformance(t *testing.T)

// Comprehensive error scenarios
func TestErrorHandling(t *testing.T)

// Integration with Go patterns
func TestCompatibility(t *testing.T)

// Custom validation functions
func TestCustomValidators(t *testing.T)
```

### 3. Performance Tests

Located in `benchmarks_test.go`:

```go
// Compare new vs old API performance
func BenchmarkNewVsOldAPI(b *testing.B)

// Test complex schema performance
func BenchmarkComplexSchemas(b *testing.B)

// Array validation with different sizes
func BenchmarkArrayValidation(b *testing.B)

// State transition performance
func BenchmarkStateTransitions(b *testing.B)

// Error key approach comparison
func BenchmarkErrorKeyPerformance(b *testing.B)

// Memory allocation testing
func BenchmarkMemoryAllocations(b *testing.B)

// Pre-configured validator performance
func BenchmarkPreConfiguredValidators(b *testing.B)

// Validation failure scenarios
func BenchmarkValidationFailures(b *testing.B)

// Concurrent validation testing
func BenchmarkConcurrentValidation(b *testing.B)

// Type conversion performance
func BenchmarkTypeConversions(b *testing.B)
```

## 📈 Performance Expectations

### Baseline Performance Targets

Based on benchmarking, the enhanced API should achieve:

- **String Validation**: >100,000 ops/sec
- **Email Validation**: >50,000 ops/sec  
- **Number Validation**: >200,000 ops/sec
- **Simple Object Validation**: >25,000 ops/sec
- **Complex Object Validation**: >5,000 ops/sec
- **Array Validation (10 items)**: >10,000 ops/sec

### Memory Efficiency

- **Schema Creation**: <1KB per simple schema
- **Validation Memory**: <100 bytes per validation
- **Concurrent Safety**: No memory leaks under load

### Scaling Characteristics

The API should maintain performance with:
- Objects with 100+ fields
- Arrays with 1000+ elements
- Nested structures 10+ levels deep
- 100+ concurrent goroutines

## 🧩 Test Data Generation

### Helper Functions for Testing

```go
// Generate test data for schemas
func generateTestData(schema interface{}) map[string]interface{} {
    // Implementation would analyze schema and generate valid test data
    return map[string]interface{}{
        "generated": "data",
    }
}

// Create stress test data
func generateLargeDataset(size int) []map[string]interface{} {
    data := make([]map[string]interface{}, size)
    for i := 0; i < size; i++ {
        data[i] = map[string]interface{}{
            "id":   i,
            "name": fmt.Sprintf("item_%d", i),
            "data": generateRandomString(100),
        }
    }
    return data
}

// Generate edge case inputs
func generateEdgeCaseInputs() []interface{} {
    return []interface{}{
        nil,
        "",
        0,
        math.Inf(1),
        math.NaN(),
        strings.Repeat("x", 10000),
        // ... more edge cases
    }
}
```

### Property-Based Testing

For more comprehensive testing, consider property-based testing:

```go
func TestStringValidationProperties(t *testing.T) {
    schema := validators.String().Min(5).Max(20).Required()
    
    // Property: all strings within bounds should pass
    for i := 0; i < 1000; i++ {
        length := 5 + rand.Intn(16) // 5-20 characters
        testString := generateRandomString(length)
        
        if err := schema.Validate(testString); err != nil {
            t.Errorf("Valid string of length %d should pass: %v", length, err)
        }
    }
    
    // Property: all strings outside bounds should fail
    for i := 0; i < 1000; i++ {
        var length int
        if rand.Intn(2) == 0 {
            length = rand.Intn(5) // 0-4 characters (too short)
        } else {
            length = 21 + rand.Intn(100) // 21+ characters (too long)
        }
        
        testString := generateRandomString(length)
        if err := schema.Validate(testString); err == nil {
            t.Errorf("Invalid string of length %d should fail", length)
        }
    }
}
```

## 🔍 Debugging and Profiling

### Memory Profiling

```bash
# Generate memory profile
go test -bench=BenchmarkComplexSchemas -memprofile=mem.prof ./zod/validators/

# Analyze with pprof
go tool pprof mem.prof

# Commands in pprof:
# (pprof) top10          # Show top memory allocators
# (pprof) list funcname  # Show source code for function
# (pprof) web            # Generate web interface
```

### CPU Profiling

```bash
# Generate CPU profile
go test -bench=BenchmarkComplexSchemas -cpuprofile=cpu.prof ./zod/validators/

# Analyze with pprof
go tool pprof cpu.prof

# Generate flame graph
go tool pprof -http=:8080 cpu.prof
```

### Race Detection

```bash
# Test for race conditions
go test -race ./zod/validators/...

# Run specific concurrent tests
go test -race -run TestConcurrency ./zod/validators/
```

## 📋 Test Checklist

### Before Release

- [ ] All unit tests pass
- [ ] All integration tests pass
- [ ] All edge case tests pass
- [ ] Benchmarks meet performance targets
- [ ] No race conditions detected
- [ ] Memory usage is within limits
- [ ] Error handling works correctly
- [ ] Custom validators work as expected
- [ ] Framework integration examples work
- [ ] Migration guide tested with real code

### Performance Verification

- [ ] String validation >100k ops/sec
- [ ] Complex object validation >5k ops/sec
- [ ] Memory usage <100 bytes per validation
- [ ] No memory leaks in concurrent tests
- [ ] Performance scales with data size
- [ ] Error scenarios don't degrade performance significantly

### API Safety Verification

- [ ] Invalid method chains don't compile
- [ ] Type safety prevents logical errors
- [ ] Error key autocompletion works
- [ ] State transitions are enforced
- [ ] Default values only work on optional fields

## 🛠️ Continuous Integration

### GitHub Actions Example

```yaml
name: Test and Benchmark

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    - uses: actions/setup-go@v3
      with:
        go-version: 1.19
    
    - name: Run tests
      run: go test -v ./zod/validators/...
    
    - name: Run tests with race detection
      run: go test -race ./zod/validators/...
    
    - name: Run benchmarks
      run: go test -bench=. -benchmem ./zod/validators/
    
    - name: Generate coverage
      run: |
        go test -coverprofile=coverage.out ./zod/validators/...
        go tool cover -html=coverage.out -o coverage.html
    
    - name: Run performance analysis
      run: go run tools/test_runner.go
```

## 📊 Performance Monitoring

### Tracking Performance Over Time

Create a performance tracking system:

```go
// tools/performance_tracker.go
func TrackPerformance() {
    results := RunBenchmarks()
    
    // Store results with timestamp
    record := PerformanceRecord{
        Timestamp: time.Now(),
        Results:   results,
        Version:   getGitVersion(),
    }
    
    // Save to database or file
    SavePerformanceRecord(record)
    
    // Check for regressions
    if HasPerformanceRegression(record) {
        log.Fatalf("Performance regression detected!")
    }
}
```

### Alerting on Regressions

Set up alerts when performance degrades:

- Response time increases by >10%
- Memory usage increases by >20%
- Throughput decreases by >15%

This comprehensive testing and benchmarking approach ensures the enhanced DX API maintains high performance while providing better developer experience and type safety.
