#!/bin/bash

# Benchmark Comparison Script for zod-go
# Compares two benchmark result files and shows performance differences

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
BOLD='\033[1m'
NC='\033[0m' # No Color

show_help() {
    cat << EOF
Usage: $0 [OPTIONS] OLD_FILE NEW_FILE

Compare two benchmark result files and show performance differences.

OPTIONS:
    -h, --help          Show this help message
    -t, --threshold N   Only show changes above N% (default: 5)
    -f, --format FORMAT Output format: table, json, csv (default: table)
    -o, --output FILE   Save comparison to file
    -s, --summary       Show summary statistics only
    -v, --verbose       Verbose output with raw numbers

ARGUMENTS:
    OLD_FILE           Baseline benchmark results file
    NEW_FILE           Current benchmark results file

EXAMPLES:
    $0 baseline.txt current.txt
    $0 -t 10 baseline.txt current.txt      # Only show >10% changes
    $0 -f json baseline.txt current.txt    # JSON output
    $0 -s baseline.txt current.txt         # Summary only
EOF
}

# Default values
THRESHOLD=5
FORMAT="table"
OUTPUT_FILE=""
SUMMARY_ONLY=false
VERBOSE=false

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        -t|--threshold)
            THRESHOLD="$2"
            shift 2
            ;;
        -f|--format)
            FORMAT="$2"
            shift 2
            ;;
        -o|--output)
            OUTPUT_FILE="$2"
            shift 2
            ;;
        -s|--summary)
            SUMMARY_ONLY=true
            shift
            ;;
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        *)
            break
            ;;
    esac
done

# Check arguments
if [[ $# -ne 2 ]]; then
    echo -e "${RED}Error: Two benchmark files required${NC}"
    show_help
    exit 1
fi

OLD_FILE="$1"
NEW_FILE="$2"

# Check if files exist
if [[ ! -f "$OLD_FILE" ]]; then
    echo -e "${RED}Error: Old file '$OLD_FILE' not found${NC}"
    exit 1
fi

if [[ ! -f "$NEW_FILE" ]]; then
    echo -e "${RED}Error: New file '$NEW_FILE' not found${NC}"
    exit 1
fi

echo -e "${BLUE}📊 Benchmark Comparison${NC}"
echo -e "${BLUE}=======================${NC}"
echo -e "Old: ${OLD_FILE}"
echo -e "New: ${NEW_FILE}"
echo -e "Threshold: ${THRESHOLD}%"
echo ""

# Function to extract benchmark data
extract_benchmarks() {
    local file="$1"
    grep -E "^Benchmark.*-[0-9]+.*ns/op" "$file" | \
    awk '{
        name = $1
        gsub(/-[0-9]+$/, "", name)  # Remove CPU count suffix
        time = $3
        gsub(/ns\/op/, "", time)
        mem = $4
        gsub(/B\/op/, "", mem)
        allocs = $5
        gsub(/allocs\/op/, "", allocs)
        print name "," time "," mem "," allocs
    }'
}

# Function to calculate percentage change
calc_change() {
    local old_val="$1"
    local new_val="$2"
    
    if [[ "$old_val" == "0" ]]; then
        echo "inf"
    else
        echo "scale=2; (($new_val - $old_val) / $old_val) * 100" | bc -l
    fi
}

# Function to format change with color
format_change() {
    local change="$1"
    local metric="$2"  # time, mem, or allocs
    
    if [[ "$change" == "inf" ]]; then
        echo -e "${YELLOW}+∞%${NC}"
        return
    fi
    
    local abs_change=$(echo "$change" | sed 's/-//')
    
    # For time, mem, and allocs: lower is better (green), higher is worse (red)
    if (( $(echo "$change < -$THRESHOLD" | bc -l) )); then
        echo -e "${GREEN}${change}%${NC}"
    elif (( $(echo "$change > $THRESHOLD" | bc -l) )); then
        echo -e "${RED}+${change}%${NC}"
    else
        echo "${change}%"
    fi
}

# Extract benchmark data
OLD_DATA=$(extract_benchmarks "$OLD_FILE")
NEW_DATA=$(extract_benchmarks "$NEW_FILE")

# Create temporary files for processing
OLD_TEMP=$(mktemp)
NEW_TEMP=$(mktemp)
echo "$OLD_DATA" > "$OLD_TEMP"
echo "$NEW_DATA" > "$NEW_TEMP"

# Function to generate comparison
generate_comparison() {
    local output_func="$1"
    
    # Join the data on benchmark name
    join -t',' -j1 -o 1.1,1.2,1.3,1.4,2.2,2.3,2.4 \
        <(sort "$OLD_TEMP") <(sort "$NEW_TEMP") | \
    while IFS=',' read -r name old_time old_mem old_allocs new_time new_mem new_allocs; do
        # Calculate changes
        time_change=$(calc_change "$old_time" "$new_time")
        mem_change=$(calc_change "$old_mem" "$new_mem")
        allocs_change=$(calc_change "$old_allocs" "$new_allocs")
        
        # Check if any change exceeds threshold
        time_abs=$(echo "$time_change" | sed 's/-//')
        mem_abs=$(echo "$mem_change" | sed 's/-//')
        allocs_abs=$(echo "$allocs_change" | sed 's/-//')
        
        if [[ "$time_change" == "inf" ]] || \
           (( $(echo "$time_abs > $THRESHOLD" | bc -l) )) || \
           (( $(echo "$mem_abs > $THRESHOLD" | bc -l) )) || \
           (( $(echo "$allocs_abs > $THRESHOLD" | bc -l) )); then
            
            $output_func "$name" "$old_time" "$new_time" "$time_change" \
                        "$old_mem" "$new_mem" "$mem_change" \
                        "$old_allocs" "$new_allocs" "$allocs_change"
        fi
    done
}

# Output format functions
output_table() {
    local name="$1" old_time="$2" new_time="$3" time_change="$4"
    local old_mem="$5" new_mem="$6" mem_change="$7"
    local old_allocs="$8" new_allocs="$9" allocs_change="${10}"
    
    if [[ "$TABLE_HEADER_PRINTED" != "true" ]]; then
        printf "%-50s %15s %15s %15s %15s %15s\n" \
            "Benchmark" "Time Change" "Memory Change" "Allocs Change" "New Time (ns)" "New Memory (B)"
        printf "%-50s %15s %15s %15s %15s %15s\n" \
            "$(printf '%*s' 50 '' | tr ' ' '-')" \
            "$(printf '%*s' 15 '' | tr ' ' '-')" \
            "$(printf '%*s' 15 '' | tr ' ' '-')" \
            "$(printf '%*s' 15 '' | tr ' ' '-')" \
            "$(printf '%*s' 15 '' | tr ' ' '-')" \
            "$(printf '%*s' 15 '' | tr ' ' '-')"
        TABLE_HEADER_PRINTED=true
    fi
    
    printf "%-50s %15s %15s %15s %15s %15s\n" \
        "${name:0:47}..." \
        "$(format_change "$time_change" "time")" \
        "$(format_change "$mem_change" "mem")" \
        "$(format_change "$allocs_change" "allocs")" \
        "$new_time" \
        "$new_mem"
}

output_json() {
    local name="$1" old_time="$2" new_time="$3" time_change="$4"
    local old_mem="$5" new_mem="$6" mem_change="$7"
    local old_allocs="$8" new_allocs="$9" allocs_change="${10}"
    
    if [[ "$JSON_FIRST" != "false" ]]; then
        echo "["
        JSON_FIRST=false
    else
        echo ","
    fi
    
    cat << EOF
  {
    "benchmark": "$name",
    "time": {
      "old": $old_time,
      "new": $new_time,
      "change_percent": $time_change
    },
    "memory": {
      "old": $old_mem,
      "new": $new_mem,
      "change_percent": $mem_change
    },
    "allocations": {
      "old": $old_allocs,
      "new": $new_allocs,
      "change_percent": $allocs_change
    }
  }
EOF
}

output_csv() {
    local name="$1" old_time="$2" new_time="$3" time_change="$4"
    local old_mem="$5" new_mem="$6" mem_change="$7"
    local old_allocs="$8" new_allocs="$9" allocs_change="${10}"
    
    if [[ "$CSV_HEADER_PRINTED" != "true" ]]; then
        echo "benchmark,old_time_ns,new_time_ns,time_change_percent,old_mem_b,new_mem_b,mem_change_percent,old_allocs,new_allocs,allocs_change_percent"
        CSV_HEADER_PRINTED=true
    fi
    
    echo "$name,$old_time,$new_time,$time_change,$old_mem,$new_mem,$mem_change,$old_allocs,$new_allocs,$allocs_change"
}

# Generate output
if [[ -n "$OUTPUT_FILE" ]]; then
    exec > "$OUTPUT_FILE"
fi

if [[ "$SUMMARY_ONLY" == "true" ]]; then
    # Generate summary statistics
    echo -e "${BOLD}Summary Statistics${NC}"
    echo -e "${BOLD}==================${NC}"
    
    total_benchmarks=$(wc -l < "$NEW_TEMP")
    echo "Total benchmarks: $total_benchmarks"
    
    # Count improvements/regressions
    improvements=0
    regressions=0
    
    join -t',' -j1 -o 1.1,1.2,2.2 <(sort "$OLD_TEMP") <(sort "$NEW_TEMP") | \
    while IFS=',' read -r name old_time new_time; do
        if (( $(echo "$new_time < $old_time" | bc -l) )); then
            ((improvements++))
        elif (( $(echo "$new_time > $old_time" | bc -l) )); then
            ((regressions++))
        fi
    done
    
    echo "Performance improvements: $improvements"
    echo "Performance regressions: $regressions"
else
    case "$FORMAT" in
        table)
            TABLE_HEADER_PRINTED=false
            generate_comparison output_table
            ;;
        json)
            JSON_FIRST=true
            generate_comparison output_json
            echo "]"
            ;;
        csv)
            CSV_HEADER_PRINTED=false
            generate_comparison output_csv
            ;;
        *)
            echo -e "${RED}Error: Unknown format '$FORMAT'${NC}"
            exit 1
            ;;
    esac
fi

# Cleanup
rm -f "$OLD_TEMP" "$NEW_TEMP"

if [[ -n "$OUTPUT_FILE" ]]; then
    echo -e "${GREEN}✅ Comparison saved to: $OUTPUT_FILE${NC}" >&2
fi
