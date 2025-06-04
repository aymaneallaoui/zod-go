# Makefile for zod-go development and maintenance

.PHONY: help build test test-verbose test-race test-coverage benchmark lint fmt clean install-tools dev-setup release-check

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build and test targets
build: ## Build the library (verify compilation)
	@echo "Building zod-go..."
	@go build ./...
	@echo "✓ Build successful"

test: ## Run all tests
	@echo "Running tests..."
	@go test ./...
	@echo "✓ Tests passed"

test-verbose: ## Run tests with verbose output
	@echo "Running tests with verbose output..."
	@go test -v ./...

test-race: ## Run tests with race detector
	@echo "Running tests with race detector..."
	@go test -race ./...
	@echo "✓ Race tests passed"

test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report generated: coverage.html"
	@go tool cover -func=coverage.out | grep total:

test-all: test test-race test-coverage ## Run all test suites

# Benchmarking targets
benchmark: ## Run performance benchmarks
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem ./benchmarks

benchmark-cpu: ## Run CPU profiling benchmarks
	@echo "Running CPU profiling benchmarks..."
	@go test -bench=. -benchmem -cpuprofile=cpu.prof ./benchmarks
	@echo "✓ CPU profile saved to cpu.prof"

benchmark-mem: ## Run memory profiling benchmarks
	@echo "Running memory profiling benchmarks..."
	@go test -bench=. -benchmem -memprofile=mem.prof ./benchmarks
	@echo "✓ Memory profile saved to mem.prof"

benchmark-compare: ## Run benchmarks and save results for comparison
	@echo "Running benchmarks for comparison..."
	@go test -bench=. -benchmem ./benchmarks > benchmark_results.txt
	@echo "✓ Benchmark results saved to benchmark_results.txt"

# Code quality targets
lint: ## Run golangci-lint
	@echo "Running linter..."
	@golangci-lint run
	@echo "✓ Linting passed"

fmt: ## Format code using gofumpt
	@echo "Formatting code..."
	@gofumpt -w .
	@echo "✓ Code formatted"

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...
	@echo "✓ Vet checks passed"

security: ## Run security checks
	@echo "Running security checks..."
	@gosec ./...
	@nancy sleuth --input go.list
	@echo "✓ Security checks passed"

# Code quality combined
quality-check: fmt vet lint security ## Run all code quality checks

# Development setup
install-tools: ## Install development tools
	@echo "Installing development tools..."
	@go install mvdan.cc/gofumpt@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	@go install github.com/sonatypecommunity/nancy@latest
	@echo "✓ Development tools installed"

dev-setup: install-tools ## Set up development environment
	@echo "Setting up development environment..."
	@go mod download
	@go mod tidy
	@echo "✓ Development environment ready"

# Dependencies management
deps-update: ## Update dependencies
	@echo "Updating dependencies..."
	@go get -u ./...
	@go mod tidy
	@echo "✓ Dependencies updated"

deps-verify: ## Verify dependencies
	@echo "Verifying dependencies..."
	@go mod verify
	@echo "✓ Dependencies verified"

deps-clean: ## Clean module cache
	@echo "Cleaning module cache..."
	@go clean -modcache
	@echo "✓ Module cache cleaned"

# Documentation
docs-generate: ## Generate documentation
	@echo "Generating documentation..."
	@mkdir -p docs
	@go doc -all ./zod > docs/api.md
	@echo "✓ Documentation generated in docs/api.md"

docs-serve: ## Serve documentation locally (requires godoc)
	@echo "Serving documentation at http://localhost:6060"
	@godoc -http=:6060

# Examples
examples-run: ## Run all examples
	@echo "Running examples..."
	@cd examples && go run basic_usage.go
	@echo "✓ Examples completed"

examples-test: ## Test that examples compile
	@echo "Testing examples compilation..."
	@go build ./examples/...
	@echo "✓ Examples compile successfully"

# Maintenance and cleanup
clean: ## Clean build artifacts and temporary files
	@echo "Cleaning up..."
	@go clean ./...
	@rm -f coverage.out coverage.html
	@rm -f cpu.prof mem.prof
	@rm -f benchmark_results.txt
	@echo "✓ Cleanup completed"

tidy: ## Tidy up go.mod and format code
	@echo "Tidying up..."
	@go mod tidy
	@$(MAKE) fmt
	@echo "✓ Tidying completed"

# Release preparation
pre-commit: quality-check test-all ## Run all checks before committing
	@echo "✓ Pre-commit checks passed"

release-check: ## Verify project is ready for release
	@echo "Checking release readiness..."
	@$(MAKE) quality-check
	@$(MAKE) test-all
	@$(MAKE) benchmark
	@$(MAKE) examples-test
	@echo "✓ Project is ready for release"

# Git hooks
git-hooks-install: ## Install git hooks
	@echo "Installing git hooks..."
	@echo '#!/bin/sh\nmake pre-commit' > .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit
	@echo "✓ Git hooks installed"

# Performance monitoring
perf-baseline: ## Create performance baseline
	@echo "Creating performance baseline..."
	@go test -bench=. -benchmem ./benchmarks > perf_baseline.txt
	@echo "✓ Performance baseline saved to perf_baseline.txt"

perf-compare: ## Compare current performance with baseline
	@echo "Comparing performance with baseline..."
	@go test -bench=. -benchmem ./benchmarks > perf_current.txt
	@echo "Baseline vs Current:"
	@echo "==================="
	@benchcmp perf_baseline.txt perf_current.txt || echo "Note: benchcmp not installed, showing raw results"
	@echo ""
	@echo "Current results:"
	@cat perf_current.txt

# CI/CD simulation
ci-test: ## Simulate CI environment testing
	@echo "Simulating CI environment..."
	@$(MAKE) clean
	@$(MAKE) deps-verify
	@$(MAKE) build
	@$(MAKE) quality-check
	@$(MAKE) test-all
	@$(MAKE) benchmark
	@echo "✓ CI simulation completed successfully"

# Multi-platform build verification
build-all-platforms: ## Build for multiple platforms
	@echo "Building for multiple platforms..."
	@GOOS=linux GOARCH=amd64 go build ./...
	@echo "✓ Linux/amd64 build successful"
	@GOOS=linux GOARCH=arm64 go build ./...
	@echo "✓ Linux/arm64 build successful"
	@GOOS=windows GOARCH=amd64 go build ./...
	@echo "✓ Windows/amd64 build successful"
	@GOOS=darwin GOARCH=amd64 go build ./...
	@echo "✓ macOS/amd64 build successful"
	@GOOS=darwin GOARCH=arm64 go build ./...
	@echo "✓ macOS/arm64 build successful"
	@echo "✓ All platform builds successful"

# Development workflow helpers
quick-check: fmt vet test ## Quick development check
	@echo "✓ Quick check completed"

full-check: clean quality-check test-all benchmark ## Comprehensive check
	@echo "✓ Full check completed"

# Help target (detailed)
help-detailed: ## Show detailed help with examples
	@echo 'zod-go Development Makefile'
	@echo '=========================='
	@echo ''
	@echo 'Common workflows:'
	@echo '  make dev-setup     - Set up development environment'
	@echo '  make quick-check   - Fast development feedback loop'
	@echo '  make pre-commit    - Run before committing code'
	@echo '  make release-check - Verify project is release-ready'
	@echo ''
	@echo 'Testing:'
	@echo '  make test          - Run basic tests'
	@echo '  make test-all      - Run all test suites'
	@echo '  make benchmark     - Run performance benchmarks'
	@echo ''
	@echo 'Code Quality:'
	@echo '  make fmt           - Format code'
	@echo '  make lint          - Run linter'
	@echo '  make security      - Run security checks'
	@echo '  make quality-check - Run all quality checks'
	@echo ''
	@echo 'Maintenance:'
	@echo '  make clean         - Clean up build artifacts'
	@echo '  make deps-update   - Update dependencies'
	@echo '  make tidy          - Tidy code and dependencies'
	@echo ''
	@echo 'For a full list of targets, run: make help'

# Version and info
version: ## Show project version info
	@echo "zod-go Development Environment"
	@echo "=============================="
	@echo "Go version: $$(go version)"
	@echo "Module: $$(go list -m)"
	@echo "Git branch: $$(git branch --show-current 2>/dev/null || echo 'unknown')"
	@echo "Git commit: $$(git rev-parse --short HEAD 2>/dev/null || echo 'unknown')"
	@echo "Last modified: $$(git log -1 --format=%cd --date=short 2>/dev/null || echo 'unknown')"
