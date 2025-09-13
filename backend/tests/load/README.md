# WanderSphere Load Testing Suite

## Prerequisites

- [k6](https://k6.io/docs/get-started/installation/) installed (`brew install k6`)
- WanderSphere backend running locally on port 19003

## Quick Start

```bash
# Smoke test (5 VUs, 2 minutes)
./run_load_tests.sh smoke

# Load test (50-100 VUs, 10 minutes)
./run_load_tests.sh load

# Stress test (100-500 VUs, 14 minutes)
./run_load_tests.sh stress

# Spike test (10→500 VUs sudden spike)
./run_load_tests.sh spike
```

## Individual Tests

```bash
./run_load_tests.sh load auth      # Auth endpoints only
./run_load_tests.sh load newsfeed  # Newsfeed with pagination
./run_load_tests.sh load posts     # Post CRUD + engagement
./run_load_tests.sh load social    # Follow/unfollow
```

## Performance Targets

| Metric | Target |
|--------|--------|
| HTTP request p95 latency | < 500ms |
| HTTP request p99 latency | < 1000ms |
| Error rate | < 1% |
| Throughput | > 100 RPS |
| Login p95 | < 300ms |
| Newsfeed p95 | < 400ms |

## Results

Results are saved to `./results/<timestamp>/` with:
- `*_summary.json` - k6 summary export
- `*_results.json` - Detailed metrics
- `*.log` - Console output
- `REPORT.md` - Report template
