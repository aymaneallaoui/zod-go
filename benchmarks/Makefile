# Makefile for zod-go benchmarks
# Provides convenient targets for running and analyzing benchmarks

.PHONY: help bench bench-all bench-quick bench-full bench-api bench-complex bench-array bench-memory bench-concurrent bench-errors bench-types
.PHONY: profile-cpu profile-mem profile-all compare baseline clean setup install-tools
.PHONY: ci-bench bench-save bench-compare

# Default target
help: ## Show this help message
	@echo "🚀 zod-go Benchmark Makefile"
	@echo "=============================="
	@echo ""
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Configuration
BENCHMARK_DIR := .
RESULTS_DIR := $(BENCHMARK_DIR)/results
SCRIPTS_DIR := $(BENCHMARK_DIR)/scripts
TIMESTAMP := $(shell date +"%Y%m%d_%H%M%S")

# Ensure results directory exists
$(RESULTS_DIR):
	@mkdir -p $(RESULTS_DIR)

# Basic benchmark targets
bench: ## Run all benchmarks
	@echo "🚀 Running all benchmarks..."
	@go test -bench=. -benchmem

bench-all: bench ## Alias for bench

bench-quick: ## Run quick benchmarks (fast operations only)
	@echo "⚡ Running quick benchmarks..."
	@$(SCRIPTS_DIR)/run_benchmarks.sh -q

bench-full: $(RESULTS_DIR) ## Run full benchmark suite with profiling
	@echo "🔬 Running full benchmark suite..."
	@$(SCRIPTS_DIR)/run_benchmarks.sh -f -s

# Category-specific benchmarks
bench-api: ## Run API performance benchmarks
	@echo "📊 Running API benchmarks..."
	@go test -bench=BenchmarkNewVsOldAPI -benchmem

bench-complex: ## Run complex schema benchmarks
	@echo "🏗️ Running complex schema benchmarks..."
	@go test -bench=BenchmarkComplexSchemas -benchmem

bench-array: ## Run array validation benchmarks
	@echo "📋 Running array validation benchmarks..."
	@go test -bench=BenchmarkArrayValidation -benchmem

bench-memory: ## Run memory allocation benchmarks
	@echo "🧠 Running memory benchmarks..."
	@go test -bench=BenchmarkMemoryAllocations -benchmem

bench-concurrent: ## Run concurrency benchmarks
	@echo "🔄 Running concurrency benchmarks..."
	@go test -bench=BenchmarkConcurrentValidation -benchmem

bench-errors: ## Run error handling benchmarks
	@echo "❌ Running error handling benchmarks..."
	@go test -bench="BenchmarkErrorKeyPerformance|BenchmarkValidationFailures" -benchmem

bench-types: ## Run type conversion benchmarks
	@echo "🔄 Running type conversion benchmarks..."
	@go test -bench=BenchmarkTypeConversions -benchmem

# Profiling targets
profile-cpu: $(RESULTS_DIR) ## Generate CPU profile
	@echo "🖥️ Generating CPU profile..."
	@go test -bench=BenchmarkComplexSchemas -cpuprofile=$(RESULTS_DIR)/cpu_$(TIMESTAMP).prof -benchmem
	@echo "📁 CPU profile saved to: $(RESULTS_DIR)/cpu_$(TIMESTAMP).prof"
	@echo "🔍 Analyze with: go tool pprof $(RESULTS_DIR)/cpu_$(TIMESTAMP).prof"

profile-mem: $(RESULTS_DIR) ## Generate memory profile
	@echo "🧠 Generating memory profile..."
	@go test -bench=BenchmarkMemoryAllocations -memprofile=$(RESULTS_DIR)/mem_$(TIMESTAMP).prof -benchmem
	@echo "📁 Memory profile saved to: $(RESULTS_DIR)/mem_$(TIMESTAMP).prof"
	@echo "🔍 Analyze with: go tool pprof $(RESULTS_DIR)/mem_$(TIMESTAMP).prof"

profile-all: profile-cpu profile-mem ## Generate both CPU and memory profiles

# Benchmark comparison and storage
bench-save: $(RESULTS_DIR) ## Run benchmarks and save results
	@echo "💾 Running benchmarks and saving results..."
	@$(SCRIPTS_DIR)/run_benchmarks.sh -s

baseline: $(RESULTS_DIR) ## Create baseline benchmark results
	@echo "📊 Creating baseline benchmark results..."
	@go test -bench=. -benchmem -count=3 | tee $(RESULTS_DIR)/baseline.txt
	@echo "✅ Baseline saved to: $(RESULTS_DIR)/baseline.txt"

bench-compare: ## Compare with baseline (requires existing baseline.txt)
	@if [ ! -f "$(RESULTS_DIR)/baseline.txt" ]; then \
		echo "❌ No baseline found. Run 'make baseline' first."; \
		exit 1; \
	fi
	@echo "📈 Comparing with baseline..."
	@$(SCRIPTS_DIR)/run_benchmarks.sh -c $(RESULTS_DIR)/baseline.txt -s

compare: ## Compare two benchmark files (usage: make compare OLD=file1.txt NEW=file2.txt)
	@if [ -z "$(OLD)" ] || [ -z "$(NEW)" ]; then \
		echo "❌ Usage: make compare OLD=file1.txt NEW=file2.txt"; \
		exit 1; \
	fi
	@$(SCRIPTS_DIR)/compare_benchmarks.sh $(OLD) $(NEW)

# CI/CD targets
ci-bench: ## Run benchmarks for CI (with consistent settings)
	@echo "🤖 Running CI benchmarks..."
	@go test -bench=. -benchmem -count=1 -timeout=30m

# Maintenance targets
clean: ## Clean up benchmark artifacts
	@echo "🧹 Cleaning up benchmark artifacts..."
	@rm -rf $(RESULTS_DIR)/*.prof
	@rm -rf $(RESULTS_DIR)/tmp_*
	@echo "✅ Cleanup complete"

clean-all: ## Clean all benchmark results and profiles
	@echo "🧹 Cleaning all benchmark results..."
	@rm -rf $(RESULTS_DIR)
	@echo "✅ All benchmark data cleaned"

setup: $(RESULTS_DIR) ## Setup benchmark environment
	@echo "🛠️ Setting up benchmark environment..."
	@chmod +x $(SCRIPTS_DIR)/*.sh
	@echo "✅ Benchmark environment ready"

install-tools: ## Install benchmark analysis tools
	@echo "🔧 Installing benchmark tools..."
	@go install golang.org/x/tools/cmd/benchcmp@latest
	@go install golang.org/x/perf/cmd/benchstat@latest
	@echo "✅ Tools installed: benchcmp, benchstat"

# Advanced targets
bench-stability: $(RESULTS_DIR) ## Test benchmark stability (run multiple times)
	@echo "🎯 Testing benchmark stability..."
	@go test -bench=BenchmarkNewVsOldAPI -benchmem -count=10 | tee $(RESULTS_DIR)/stability_$(TIMESTAMP).txt
	@echo "📊 Results saved to: $(RESULTS_DIR)/stability_$(TIMESTAMP).txt"

bench-regression: ## Run regression test suite
	@echo "🔍 Running regression test suite..."
	@$(SCRIPTS_DIR)/run_benchmarks.sh --count=5 -s
	@echo "✅ Regression tests complete"

# Performance targets for different scenarios
bench-small: ## Benchmark with small datasets
	@echo "📏 Running small dataset benchmarks..."
	@go test -bench="SmallArray|NewAPI_String_Simple" -benchmem

bench-large: ## Benchmark with large datasets
	@echo "📏 Running large dataset benchmarks..."
	@go test -bench="LargeArray|ComplexUserSchema" -benchmem

# Development helpers
watch: ## Watch for changes and re-run benchmarks
	@echo "👀 Watching for changes..."
	@while true; do \
		inotifywait -e modify,create,delete -r . --exclude="\.git|results|\.prof" 2>/dev/null; \
		echo "🔄 Files changed, running quick benchmarks..."; \
		make bench-quick; \
		echo "✅ Benchmarks complete. Watching for changes..."; \
	done

# Documentation
bench-report: $(RESULTS_DIR) ## Generate comprehensive benchmark report
	@echo "📋 Generating benchmark report..."
	@go test -bench=. -benchmem -count=3 > $(RESULTS_DIR)/full_report_$(TIMESTAMP).txt
	@echo "# zod-go Benchmark Report" > $(RESULTS_DIR)/report_$(TIMESTAMP).md
	@echo "Generated: $(shell date)" >> $(RESULTS_DIR)/report_$(TIMESTAMP).md
	@echo "" >> $(RESULTS_DIR)/report_$(TIMESTAMP).md
	@echo "## Results" >> $(RESULTS_DIR)/report_$(TIMESTAMP).md
	@echo '```' >> $(RESULTS_DIR)/report_$(TIMESTAMP).md
	@cat $(RESULTS_DIR)/full_report_$(TIMESTAMP).txt >> $(RESULTS_DIR)/report_$(TIMESTAMP).md
	@echo '```' >> $(RESULTS_DIR)/report_$(TIMESTAMP).md
	@echo "📄 Report saved to: $(RESULTS_DIR)/report_$(TIMESTAMP).md"

# Show current performance status
status: ## Show quick performance status
	@echo "⚡ Quick performance check..."
	@go test -bench="NewAPI_String_Simple|NewAPI_Number_Validation" -benchmem -count=1 | \
		awk '/^Benchmark/ { printf "%-50s %10s %10s %10s\n", $$1, $$3, $$4, $$5 }'
