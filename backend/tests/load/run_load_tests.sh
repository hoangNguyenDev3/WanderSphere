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
cat > "$RESULTS_DIR/REPORT.md" << REPORT_EOF
# WanderSphere Load Test Report

## Test Configuration
- **Profile**: ${PROFILE}
- **Base URL**: ${BASE_URL}
- **Date**: $(date)

## Performance Targets
| Metric | Target | Status |
|--------|--------|--------|
| HTTP p95 latency | < 500ms | Check results |
| HTTP p99 latency | < 1000ms | Check results |
| Error rate | < 1% | Check results |
| Throughput | > 100 RPS | Check results |

## Structured Metrics
See:
- \`METRICS_SUMMARY.md\` for per-suite throughput/latency/error rates and threshold status.
- \`*_summary.json\` for raw k6 summaries used to derive all metrics.
REPORT_EOF

echo "Report template saved to: $RESULTS_DIR/REPORT.md"

echo ""
echo "Generating structured metrics summary..."

python3 - "$RESULTS_DIR" << 'PY_EOF'
import json
import os
import sys
from datetime import datetime

results_dir = sys.argv[1]
summary_files = sorted(
    f for f in os.listdir(results_dir) if f.endswith("_summary.json")
)

if not summary_files:
    print("No *_summary.json files found, skipping METRICS_SUMMARY.md generation.")
    raise SystemExit(0)


def read_metric_value(metrics, metric_name, value_key):
    metric = metrics.get(metric_name, {})
    values = metric.get("values", {})
    return values.get(value_key)


def fmt_float(value, digits=2):
    if value is None:
        return "N/A"
    return f"{float(value):.{digits}f}"


def fmt_ms(seconds):
    if seconds is None:
        return "N/A"
    return f"{float(seconds) * 1000:.2f}ms"


def threshold_status(metric_obj):
    thresholds = metric_obj.get("thresholds", {})
    if not thresholds:
        return "N/A"
    statuses = []
    for name, detail in thresholds.items():
        ok = detail.get("ok")
        statuses.append(f"{name}:{'PASS' if ok else 'FAIL'}")
    return ", ".join(statuses)


lines = []
lines.append("# Load Test Metrics Summary")
lines.append("")
lines.append(f"- Generated: {datetime.now().isoformat(timespec='seconds')}")
lines.append(f"- Results directory: `{results_dir}`")
lines.append("")
lines.append("| Suite | Iter/s | HTTP req/s | HTTP p95 | HTTP p99 | HTTP error rate | Thresholds |")
lines.append("|---|---:|---:|---:|---:|---:|---|")

for filename in summary_files:
    path = os.path.join(results_dir, filename)
    with open(path, "r", encoding="utf-8") as f:
        data = json.load(f)

    metrics = data.get("metrics", {})
    suite = filename.replace("_summary.json", "")

    iter_rate = read_metric_value(metrics, "iterations", "rate")
    req_rate = read_metric_value(metrics, "http_reqs", "rate")
    p95 = read_metric_value(metrics, "http_req_duration", "p(95)")
    p99 = read_metric_value(metrics, "http_req_duration", "p(99)")
    failed_rate = read_metric_value(metrics, "http_req_failed", "rate")

    http_req_duration = metrics.get("http_req_duration", {})
    http_req_failed = metrics.get("http_req_failed", {})
    http_reqs = metrics.get("http_reqs", {})

    threshold_parts = []
    if http_req_duration.get("thresholds"):
        threshold_parts.append("http_req_duration[" + threshold_status(http_req_duration) + "]")
    if http_req_failed.get("thresholds"):
        threshold_parts.append("http_req_failed[" + threshold_status(http_req_failed) + "]")
    if http_reqs.get("thresholds"):
        threshold_parts.append("http_reqs[" + threshold_status(http_reqs) + "]")

    lines.append(
        "| "
        + suite
        + " | "
        + fmt_float(iter_rate)
        + " | "
        + fmt_float(req_rate)
        + " | "
        + fmt_ms(p95)
        + " | "
        + fmt_ms(p99)
        + " | "
        + (f"{float(failed_rate) * 100:.2f}%" if failed_rate is not None else "N/A")
        + " | "
        + ("; ".join(threshold_parts) if threshold_parts else "N/A")
        + " |"
    )

lines.append("")
lines.append("## Notes")
lines.append("- `http_req_duration` reflects k6 end-to-end request timing.")
lines.append("- Use this file for CV/README claims; pair each claim with the matching `*_summary.json` artifact.")

output_path = os.path.join(results_dir, "METRICS_SUMMARY.md")
with open(output_path, "w", encoding="utf-8") as f:
    f.write("\n".join(lines) + "\n")

print(f"Wrote {output_path}")
PY_EOF

echo "Done!"
