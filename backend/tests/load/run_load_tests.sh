#!/bin/bash
set -e

# WanderSphere Load Testing Suite
# Usage: ./run_load_tests.sh [smoke|load|stress|spike] [test_name]
# Examples:
#   ./run_load_tests.sh smoke          # Run all tests with smoke profile
#   ./run_load_tests.sh load auth      # Run auth test with load profile
#   ./run_load_tests.sh stress posts   # Run posts test with stress profile

PROFILE=${1:-smoke}
TEST=${2:-all}
BASE_URL=${BASE_URL:-http://localhost:19003/api/v1}
RESULTS_DIR="./results/$(date +%Y%m%d_%H%M%S)"

echo "=========================================="
echo "  WanderSphere Load Testing Suite"
echo "=========================================="
echo "Profile:  $PROFILE"
echo "Test:     $TEST"
echo "Base URL: $BASE_URL"
echo "Results:  $RESULTS_DIR"
echo "=========================================="

# Check k6 installation
if ! command -v k6 &> /dev/null; then
    echo "ERROR: k6 is not installed. Install with: brew install k6"
    exit 1
fi

# Create results directory
mkdir -p "$RESULTS_DIR"

# Function to run a single test
run_test() {
    local test_name=$1
    local test_file=$2
    echo ""
    echo "--- Running: $test_name ($PROFILE profile) ---"
    k6 run \
        --env PROFILE="$PROFILE" \
        --env BASE_URL="$BASE_URL" \
        --summary-export="$RESULTS_DIR/${test_name}_summary.json" \
        --out json="$RESULTS_DIR/${test_name}_results.json" \
        "$test_file" 2>&1 | tee "$RESULTS_DIR/${test_name}.log"
    echo "--- $test_name complete ---"
}

# Run selected tests
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

if [ "$TEST" = "all" ] || [ "$TEST" = "auth" ]; then
    run_test "auth" "$SCRIPT_DIR/auth_test.js"
fi

if [ "$TEST" = "all" ] || [ "$TEST" = "newsfeed" ]; then
    run_test "newsfeed" "$SCRIPT_DIR/newsfeed_test.js"
fi

if [ "$TEST" = "all" ] || [ "$TEST" = "posts" ]; then
    run_test "posts" "$SCRIPT_DIR/posts_test.js"
fi

if [ "$TEST" = "all" ] || [ "$TEST" = "social" ]; then
    run_test "social" "$SCRIPT_DIR/social_test.js"
fi

echo ""
echo "=========================================="
echo "  Load Test Results Summary"
echo "=========================================="
echo "Results saved to: $RESULTS_DIR"
echo ""
echo "Key metrics to check:"
echo "  - http_req_duration p95 < 500ms"
echo "  - http_req_duration p99 < 1000ms"
echo "  - http_req_failed rate < 1%"
echo "  - http_reqs rate > 100 RPS"
echo "=========================================="

# Generate summary report
echo ""
echo "Generating summary report..."
cat > "$RESULTS_DIR/REPORT.md" << 'REPORT_EOF'
# WanderSphere Load Test Report

## Test Configuration
- **Profile**: $PROFILE
- **Base URL**: $BASE_URL
- **Date**: $(date)

## Performance Targets
| Metric | Target | Status |
|--------|--------|--------|
| HTTP p95 latency | < 500ms | Check results |
| HTTP p99 latency | < 1000ms | Check results |
| Error rate | < 1% | Check results |
| Throughput | > 100 RPS | Check results |

## Test Suites
1. **Auth** - Signup, login, profile retrieval
2. **Newsfeed** - Feed retrieval with cursor pagination
3. **Posts** - Create, read, like, comment
4. **Social** - Follow, unfollow, get followers

## CV-Ready Metrics (fill after test run)
- Sustained X RPS with p99 latency under Y ms across Z concurrent users
- Feed retrieval p95 < Xms with cursor-based pagination under N concurrent readers
- Zero-downtime under spike load of N→500 concurrent users
REPORT_EOF

echo "Report template saved to: $RESULTS_DIR/REPORT.md"
echo "Done!"
