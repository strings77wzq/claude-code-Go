# Benchmarks

> **Status: Methodology defined — measurements in progress.** Reproducible performance benchmarks will be published here once the v0.2 measurement harness is finalized.

## Benchmark Methodology

All benchmarks follow this protocol:

### Environment

- **Hardware:** Recorded per run (CPU model, cores, RAM)
- **OS:** Recorded per run (Linux kernel version / macOS version)
- **Go version:** `go version` output recorded
- **Binary:** `go-code` built with `CGO_ENABLED=0` and `-ldflags="-s -w"`

### Commands

```bash
# Build the release binary
make build

# Run benchmarks (placeholder — harness under development)
# go-code bench --scenario standard --runs 5 --output results.json
```

### Metrics Tracked

| Metric | Description |
|--------|-------------|
| Startup latency | Time from binary invocation to ready-for-input |
| First-token latency | Time from prompt to first response token |
| Tool execution latency | Round-trip time for Read/Write/Edit/Bash tools |
| Memory RSS | Peak resident set size during a standard session |
| Binary size | Stripped binary size per platform |

### Reproduction

All benchmark results include:
- Exact command with arguments
- Environment snapshot (OS, Go version, CPU, RAM)
- Date and commit hash
- Raw output file (linked when published)

## Why Go?

This project uses Go for single-binary distribution, straightforward cross-compilation, mature HTTP/SSE primitives, and a small runtime dependency surface. Quantitative comparisons with Python or Rust will be published only after reproducible benchmark commands and raw results exist.

*Last updated: April 2026. Measurements pending.*
