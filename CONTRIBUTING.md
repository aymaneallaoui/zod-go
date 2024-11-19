# Contributing to Zod-Go

Thank you for your interest in contributing to Zod-Go! This document provides guidelines and instructions for contributors.

## Code of Conduct

By participating in this project, you agree to abide by our code of conduct. We are committed to creating a welcoming and inclusive environment for all contributors.

## How to Contribute

### Reporting Issues

- Use the GitHub issue tracker to report bugs or request features
- Before creating a new issue, search existing issues to avoid duplicates
- Provide clear, detailed descriptions with steps to reproduce for bugs
- For feature requests, explain the use case and expected behavior

### Development Setup

1. **Fork and Clone**
   ```bash
   git clone https://github.com/YOUR_USERNAME/zod-go.git
   cd zod-go
   ```

2. **Install Dependencies**
   ```bash
   go mod tidy
   ```

3. **Install Development Tools**
   ```bash
   # Install golangci-lint for linting
   curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.54.2

   # Install gofumpt for formatting
   go install mvdan.cc/gofumpt@latest
   ```

### Making Changes

1. **Create a Branch**
   ```bash
   git checkout -b feature/your-feature-name
   # or
   git checkout -b fix/issue-description
   ```

2. **Write Code**
   - Follow Go conventions and best practices
   - Add tests for new functionality
   - Update documentation as needed
   - Ensure all tests pass

3. **Code Style**
   - Use `gofumpt` for formatting: `gofumpt -w .`
   - Run linter: `golangci-lint run`
   - Follow the existing code style
   - Use meaningful variable and function names

4. **Testing**
   ```bash
   # Run all tests
   go test ./...

   # Run tests with coverage
   go test -coverprofile=coverage.out ./...
   go tool cover -html=coverage.out

   # Run benchmarks
   go test -bench=. ./benchmarks

   # Test with race detector
   go test -race ./...
   ```

5. **Documentation**
   - Update README.md if adding new features
   - Add/update code comments for public APIs
   - Include examples for new validators

### Pull Request Process

1. **Before Submitting**
   - Ensure all tests pass
   - Run linter and fix any issues
   - Update documentation
   - Squash commits if necessary

2. **PR Description**
   - Provide a clear title and description
   - Reference related issues with "Fixes #123" or "Closes #123"
   - List breaking changes if any
   - Include examples of new functionality

3. **Review Process**
   - All PRs require review before merging
   - Address feedback promptly
   - Keep discussions constructive and professional

## Project Structure

```
├── zod/                    # Core library code
│   ├── errors.go          # Error types and handling
│   ├── schema.go          # Core schema interface
│   └── validators/        # Validator implementations
│       ├── string.go      # String validation
│       ├── number.go      # Number validation
│       ├── bool.go        # Boolean validation
│       ├── array.go       # Array validation
│       ├── object.go      # Object validation
│       └── map.go         # Map validation
├── tests/                 # Test files
├── benchmarks/            # Benchmark tests
├── examples/              # Usage examples
└── docs/                  # Additional documentation
```

## Writing Tests

### Test Structure
- Place tests in the `tests/` directory
- Use descriptive test names: `TestStringValidator_MinLength_ValidInput`
- Include both positive and negative test cases
- Test edge cases and error conditions

### Test Example
```go
func TestStringValidator_EmailValidation(t *testing.T) {
    tests := []struct {
        name    string
        input   interface{}
        wantErr bool
        errMsg  string
    }{
        {
            name:    "valid email",
            input:   "user@example.com",
            wantErr: false,
        },
        {
            name:    "invalid email format",
            input:   "invalid-email",
            wantErr: true,
            errMsg:  "invalid email format",
        },
        {
            name:    "non-string input",
            input:   123,
            wantErr: true,
            errMsg:  "expected string",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            schema := validators.String().Email()
            err := schema.Validate(tt.input)
            
            if (err != nil) != tt.wantErr {
                t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            
            if tt.wantErr && err != nil {
                if !strings.Contains(err.Error(), tt.errMsg) {
                    t.Errorf("Expected error message to contain %q, got %q", tt.errMsg, err.Error())
                }
            }
        })
    }
}
```

## Writing Benchmarks

- Place benchmarks in the `benchmarks/` directory
- Focus on performance-critical paths
- Include memory allocation benchmarks
- Compare against baseline implementations

### Benchmark Example
```go
func BenchmarkStringValidation(b *testing.B) {
    schema := validators.String().Min(5).Max(100)
    testString := "hello world"
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = schema.Validate(testString)
    }
}
```

## Adding New Validators

1. **Create Validator File**
   - Add new file in `zod/validators/`
   - Follow existing patterns and naming conventions

2. **Implement Interface**
   ```go
   type CustomSchema struct {
       // validator fields
   }

   func Custom() *CustomSchema {
       return &CustomSchema{}
   }

   func (c *CustomSchema) Validate(data interface{}) error {
       // validation logic
   }
   ```

3. **Add Tests**
   - Comprehensive test coverage
   - Test all validation rules
   - Test error cases and edge conditions

4. **Update Documentation**
   - Add examples to README.md
   - Document all public methods
   - Include usage examples

## Performance Guidelines

- Avoid unnecessary allocations in hot paths
- Use object pooling for frequently created objects
- Benchmark performance-critical changes
- Consider concurrent validation for large datasets
- Profile memory usage for complex validators

## Documentation Guidelines

- Use clear, concise language
- Provide practical examples
- Keep code comments up to date
- Use godoc conventions for public APIs
- Include error handling examples

## Release Process

1. Version bumps follow semantic versioning (SemVer)
2. Update CHANGELOG.md with release notes
3. Tag releases with `git tag v1.x.x`
4. Releases are automated through GitHub Actions

## Getting Help

- Open an issue for questions about contributing
- Join discussions in existing issues and PRs
- Reach out to maintainers for guidance on larger changes

## Recognition

Contributors will be recognized in:
- CHANGELOG.md for each release
- README.md contributors section
- Release notes for significant contributions

Thank you for contributing to Zod-Go! 🎉
