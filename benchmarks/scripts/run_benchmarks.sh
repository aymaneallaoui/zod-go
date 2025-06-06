#!/bin/bash

# Benchmark Runner Script for zod-go
# This script provides easy ways to run different benchmark scenarios

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
BENCHMARKS_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
RESULTS_DIR="${BENCHMARKS_DIR}/results"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")

# Create results directory if it doesn't exist
mkdir -p "${RESULTS_DIR}"

echo -e "${BLUE}🚀 zod-go Benchmark Runner${NC}"
echo -e "${BLUE}================================${NC}"

show_help() {
    cat << EOF
Usage: $0 [OPTIONS] [BENCHMARK_PATTERN]

OPTIONS:
    -h, --help          Show this help message
    -a, --all           Run all benchmarks (default)
    -q, --quick         Run quick benchmarks only
    -f, --full          Run full benchmark suite with profiling
    -c, --compare FILE  Compare with previous benchmark results
    -s, --save          Save results to timestamped file
    -v, --verbose       Verbose output
    --count N           Run benchmarks N times (default: 1)
    --cpu-profile       Generate CPU profile
    --mem-profile       Generate memory profile

BENCHMARK_PATTERNS:
    api                 API performance benchmarks
    complex             Complex schema benchmarks
    array               Array validation benchmarks
    concurrent          Concurrency benchmarks
    memory              Memory allocation benchmarks
    errors              Error handling benchmarks
    types               Type conversion benchmarks

EXAMPLES:
    $0                              # Run all benchmarks
    $0 -s                           # Run all and save results
    $0 -f                           # Full suite with profiling
    $0 api                          # Run only API benchmarks
    $0 --count=5 complex            # Run complex benchmarks 5 times
    $0 -c results/baseline.txt      # Compare with baseline
EOF
}

# Parse command line arguments
BENCHMARK_PATTERN="."
SAVE_RESULTS=false
VERBOSE=false
COUNT=1
CPU_PROFILE=false
MEM_PROFILE=false
COMPARE_FILE=""
QUICK_MODE=false
FULL_MODE=false

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        -a|--all)
            BENCHMARK_PATTERN="."
            shift
            ;;
        -q|--quick)
            QUICK_MODE=true
            shift
            ;;
        -f|--full)
            FULL_MODE=true
            CPU_PROFILE=true
            MEM_PROFILE=true
            COUNT=5
            shift
            ;;
        -s|--save)
            SAVE_RESULTS=true
            shift
            ;;
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -c|--compare)
            COMPARE_FILE="$2"
            shift 2
            ;;
        --count)
            COUNT="$2"
            shift 2
            ;;
        --cpu-profile)
            CPU_PROFILE=true
            shift
            ;;
        --mem-profile)
            MEM_PROFILE=true
            shift
            ;;
        api)
            BENCHMARK_PATTERN="BenchmarkNewVsOldAPI"
            shift
            ;;
        complex)
            BENCHMARK_PATTERN="BenchmarkComplexSchemas"
            shift
            ;;
        array)
            BENCHMARK_PATTERN="BenchmarkArrayValidation"
            shift
            ;;
        concurrent)
            BENCHMARK_PATTERN="BenchmarkConcurrentValidation"
            shift
            ;;
        memory)
            BENCHMARK_PATTERN="BenchmarkMemoryAllocations"
            shift
            ;;
        errors)
            BENCHMARK_PATTERN="BenchmarkErrorKeyPerformance|BenchmarkValidationFailures"
            shift
            ;;
        types)
            BENCHMARK_PATTERN="BenchmarkTypeConversions"
            shift
            ;;
        *)
            BENCHMARK_PATTERN="$1"
            shift
            ;;
    esac
done

# Quick mode patterns
if [[ "$QUICK_MODE" == true ]]; then
    BENCHMARK_PATTERN="BenchmarkNewVsOldAPI/NewAPI_String_Simple|BenchmarkNewVsOldAPI/NewAPI_Number_Validation"
    COUNT=1
    echo -e "${YELLOW}📋 Running quick benchmarks...${NC}"
fi

# Full mode setup
if [[ "$FULL_MODE" == true ]]; then
    echo -e "${BLUE}🔬 Running full benchmark suite with profiling...${NC}"
fi

# Change to benchmarks directory
cd "${BENCHMARKS_DIR}"

# Prepare benchmark command
BENCH_CMD="go test -bench=${BENCHMARK_PATTERN} -benchmem"

if [[ "$COUNT" -gt 1 ]]; then
    BENCH_CMD="${BENCH_CMD} -count=${COUNT}"
fi

if [[ "$CPU_PROFILE" == true ]]; then
    BENCH_CMD="${BENCH_CMD} -cpuprofile=${RESULTS_DIR}/cpu_${TIMESTAMP}.prof"
fi

if [[ "$MEM_PROFILE" == true ]]; then
    BENCH_CMD="${BENCH_CMD} -memprofile=${RESULTS_DIR}/mem_${TIMESTAMP}.prof"
fi

if [[ "$VERBOSE" == true ]]; then
    BENCH_CMD="${BENCH_CMD} -v"
fi

echo -e "${GREEN}📊 Running benchmarks...${NC}"
echo -e "${BLUE}Command: ${BENCH_CMD}${NC}"
echo ""

# Run benchmarks
if [[ "$SAVE_RESULTS" == true ]]; then
    RESULTS_FILE="${RESULTS_DIR}/benchmark_${TIMESTAMP}.txt"
    echo -e "${YELLOW}💾 Saving results to: ${RESULTS_FILE}${NC}"
    
    # Add header to results file
    {
        echo "# zod-go Benchmark Results"
        echo "# Generated: $(date)"
        echo "# Command: ${BENCH_CMD}"
        echo "# Go Version: $(go version)"
        echo "# System: $(uname -a)"
        echo ""
    } > "${RESULTS_FILE}"
    
    # Run and save
    eval "${BENCH_CMD}" | tee -a "${RESULTS_FILE}"
    
    echo ""
    echo -e "${GREEN}✅ Results saved to: ${RESULTS_FILE}${NC}"
else
    eval "${BENCH_CMD}"
fi

# Compare with previous results if requested
if [[ -n "$COMPARE_FILE" && -f "$COMPARE_FILE" ]]; then
    echo ""
    echo -e "${BLUE}📈 Comparing with previous results...${NC}"
    
    if command -v benchcmp >/dev/null 2>&1; then
        if [[ "$SAVE_RESULTS" == true ]]; then
            benchcmp "${COMPARE_FILE}" "${RESULTS_FILE}"
        else
            echo -e "${YELLOW}⚠️  Cannot compare without saving current results. Use -s flag.${NC}"
        fi
    else
        echo -e "${YELLOW}⚠️  benchcmp not found. Install with: go install golang.org/x/tools/cmd/benchcmp@latest${NC}"
    fi
fi

# Show profile analysis commands if profiles were generated
if [[ "$CPU_PROFILE" == true || "$MEM_PROFILE" == true ]]; then
    echo ""
    echo -e "${BLUE}🔍 Profile Analysis Commands:${NC}"
    
    if [[ "$CPU_PROFILE" == true ]]; then
        echo -e "${GREEN}CPU Profile:${NC} go tool pprof ${RESULTS_DIR}/cpu_${TIMESTAMP}.prof"
    fi
    
    if [[ "$MEM_PROFILE" == true ]]; then
        echo -e "${GREEN}Memory Profile:${NC} go tool pprof ${RESULTS_DIR}/mem_${TIMESTAMP}.prof"
    fi
fi

echo ""
echo -e "${GREEN}🎉 Benchmark run completed!${NC}"

# Show results summary for full mode
if [[ "$FULL_MODE" == true ]]; then
    echo ""
    echo -e "${BLUE}📋 Results Summary:${NC}"
    echo -e "${BLUE}===================${NC}"
    
    if [[ "$SAVE_RESULTS" == true ]]; then
        # Extract key performance metrics
        echo -e "${GREEN}⚡ Fastest Operations:${NC}"
        grep -E "ns/op" "${RESULTS_FILE}" | sort -k3 -n | head -5
        
        echo ""
        echo -e "${GREEN}🧠 Most Memory Efficient:${NC}"
        grep -E "B/op.*allocs/op" "${RESULTS_FILE}" | sort -k4 -n | head -5
    fi
fi
