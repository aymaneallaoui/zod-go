module github.com/aymaneallaoui/zod-go

go 1.21

// Production-ready validation library inspired by TypeScript's Zod
// 
// This module provides comprehensive data validation with:
// - Fluent, chainable API
// - Rich type support (string, number, boolean, array, object)
// - Custom validation functions
// - Detailed error reporting
// - High performance with concurrent validation
// - Zero external dependencies for core functionality

require (
	// No external dependencies required for core functionality
	// This ensures minimal attack surface and easy integration
)

// Development and testing dependencies (not included in production builds)
require (
	// These are only needed for development and will not be included
	// in applications that import this library
)

// Go version compatibility:
// - Minimum: Go 1.21 (for improved type inference and performance)
// - Tested: Go 1.21, 1.22
// - Recommended: Go 1.22+ for best performance
//
// The library uses only standard library features to ensure
// maximum compatibility and minimal dependencies.
