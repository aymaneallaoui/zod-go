#!/bin/bash

# Make benchmark scripts executable
chmod +x run_benchmarks.sh
chmod +x compare_benchmarks.sh

echo "✅ Benchmark scripts are now executable"
echo ""
echo "🚀 Quick start:"
echo "  ./run_benchmarks.sh -q     # Quick benchmarks"
echo "  ./run_benchmarks.sh -f     # Full suite with profiling"
echo "  ../Makefile targets:"
echo "    make bench               # Run all benchmarks"
echo "    make bench-api           # API benchmarks only"
echo "    make status              # Quick performance check"
