# Benchmarks

This directory contains comprehensive performance benchmarks for the zod-go validation library.

## Quick Start

```bash
# Run all benchmarks
go test -bench=. -benchmem

# Run specific benchmark categories
go test -bench=BenchmarkNewVsOldAPI -benchmem
go test -bench=BenchmarkComplexSchemas -benchmem
go test -bench=BenchmarkArrayValidation -benchmem
```

## Benchmark Categories

### 🚀 **Core API Performance**
- `BenchmarkNewVsOldAPI` - String, Email, Number validation performance
- `BenchmarkStateTransitions` - Builder pattern efficiency
- `BenchmarkMemoryAllocations` - Memory usage analysis

### 🏗️ **Complex Scenarios**
- `BenchmarkComplexSchemas` - Real-world nested object validation
- `BenchmarkArrayValidation` - Array performance at different sizes
- `BenchmarkConcurrentValidation` - Thread safety and parallel performance

### 📊 **Specialized Tests**
- `BenchmarkErrorKeyPerformance` - Error handling efficiency
- `BenchmarkPreConfiguredValidators` - Pre-built validator performance
- `BenchmarkValidationFailures` - Error path performance
- `BenchmarkTypeConversions` - Type checking and conversion overhead

## Running Benchmarks

### Basic Benchmark Run
```bash
# All benchmarks with memory stats
go test -bench=. -benchmem

# Save results to file
go test -bench=. -benchmem | tee benchmark_results.txt
```

### Performance Comparison
```bash
# Run benchmarks multiple times for stability
go test -bench=. -benchmem -count=5

# Compare with CPU profiling
go test -bench=BenchmarkComplexSchemas -benchmem -cpuprofile=cpu.prof

# Memory profiling
go test -bench=BenchmarkMemoryAllocations -benchmem -memprofile=mem.prof
```

### Specific Benchmark Categories
```bash
# API Performance
go test -bench=BenchmarkNewVsOldAPI -benchmem

# Complex Object Validation
go test -bench=BenchmarkComplexSchemas -benchmem

# Array Scaling Performance
go test -bench=BenchmarkArrayValidation -benchmem

# Concurrent Performance
go test -bench=BenchmarkConcurrentValidation -benchmem

# Error Handling Performance
go test -bench=BenchmarkErrorKeyPerformance -benchmem
go test -bench=BenchmarkValidationFailures -benchmem

# Memory Efficiency
go test -bench=BenchmarkMemoryAllocations -benchmem

# Pre-configured Validators
go test -bench=BenchmarkPreConfiguredValidators -benchmem

# Type System Performance
go test -bench=BenchmarkTypeConversions -benchmem
```

## Understanding Results

### Benchmark Output Format
```
BenchmarkNewVsOldAPI/NewAPI_String_Simple-32    84704486    15.03 ns/op    16 B/op    1 allocs/op
```

- `84704486` - Number of iterations
- `15.03 ns/op` - Nanoseconds per operation
- `16 B/op` - Bytes allocated per operation
- `1 allocs/op` - Memory allocations per operation

### Performance Targets

#### ⚡ **Excellent Performance**
- Simple validation: `< 20 ns/op`
- Zero allocations for cached validators
- Memory usage scales linearly

#### ✅ **Good Performance**
- Complex validation: `< 1μs/op`
- Minimal allocations per validation
- Consistent performance across data sizes

#### 🎯 **Optimization Opportunities**
- Operations taking `> 10μs/op`
- High allocation counts
- Non-linear scaling

## Performance Analysis Scripts

### Automated Benchmark Runner
```bash
# Run the benchmark script
./scripts/run_benchmarks.sh

# Compare with previous results
./scripts/compare_benchmarks.sh baseline.txt current.txt
```

### Memory Profiling
```bash
# Generate memory profile
go test -bench=BenchmarkComplexSchemas -memprofile=mem.prof -memprofilerate=1

# Analyze with pprof
go tool pprof mem.prof
```

### CPU Profiling
```bash
# Generate CPU profile
go test -bench=BenchmarkArrayValidation -cpuprofile=cpu.prof

# Analyze with pprof
go tool pprof cpu.prof
```

## Benchmark Results History

### Latest Results (v2.0.0)
```
BenchmarkNewVsOldAPI/NewAPI_String_Simple-32         84704486    15.03 ns/op      16 B/op       1 allocs/op
BenchmarkNewVsOldAPI/NewAPI_Number_Validation-32     397649493   2.993 ns/op       0 B/op       0 allocs/op
BenchmarkComplexSchemas/ComplexUserSchema_Valid-32   198004      5726 ns/op      7588 B/op     107 allocs/op
BenchmarkArrayValidation/SmallArray_10_Items-32      3497109     295.5 ns/op      344 B/op      12 allocs/op
BenchmarkConcurrentValidation/Concurrent_Validation-32 273393046 4.720 ns/op      16 B/op       1 allocs/op
```

### Performance Highlights
- **Number validation**: 2.993 ns/op with 0 allocations ⚡
- **Concurrent speedup**: ~3.5x improvement over sequential
- **Linear array scaling**: Perfect O(n) performance
- **Memory efficient**: Zero-allocation fast paths

## Continuous Integration

Add to your CI pipeline:

```yaml
name: Benchmarks
on: [push, pull_request]

jobs:
  benchmark:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Run Benchmarks
        run: |
          cd benchmarks
          go test -bench=. -benchmem | tee benchmark_results.txt
      
      - name: Upload Results
        uses: actions/upload-artifact@v3
        with:
          name: benchmark-results
          path: benchmarks/benchmark_results.txt
```

## Contributing

When adding new benchmarks:

1. **Use descriptive names**: `BenchmarkFeature_Scenario_Condition`
2. **Include memory stats**: Always run with `-benchmem`
3. **Reset timer**: Use `b.ResetTimer()` after setup
4. **Test realistic data**: Use representative input sizes
5. **Document purpose**: Explain what the benchmark measures

### Example Benchmark
```go
func BenchmarkNewFeature_LargeInput(b *testing.B) {
    // Setup
    schema := validators.String().Min(1).Max(1000).Required()
    testData := strings.Repeat("x", 500)
    
    // Reset timer after setup
    b.ResetTimer()
    
    // Benchmark loop
    for i := 0; i < b.N; i++ {
        _ = schema.Validate(testData)
    }
}
```

## Troubleshooting

### Common Issues

**Unstable results**: Run with `-count=5` for multiple iterations
**High variance**: Check for background processes or thermal throttling
**Memory leaks**: Use `-memprofile` to analyze allocations
**Slow benchmarks**: Consider reducing test data size or iterations

### Platform Considerations

Results may vary between:
- Different CPU architectures (Intel vs ARM)
- Operating systems (Linux vs Windows vs macOS)
- Go versions (1.20 vs 1.21+)
- Hardware specifications (cache sizes, memory speed)

Always benchmark on your target platform for accurate results.
