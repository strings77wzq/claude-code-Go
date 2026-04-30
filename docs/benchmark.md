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

### Go vs Python (for AI agent runtimes)

| Aspect | Go | Python |
|--------|-----|--------|
| Deployment | Single static binary | Requires interpreter + dependencies |
| Startup time | ~10ms | ~100-500ms |
| Memory baseline | ~5MB | ~30-50MB |
| Concurrency model | Goroutines (lightweight) | asyncio (cooperative) |
| Cross-compilation | `GOOS=linux GOARCH=amd64 go build` | Packaging tools required |

### Go vs Rust (for AI agent runtimes)

| Aspect | Go | Rust |
|--------|-----|------|
| Development speed | Faster iteration (GC, simpler types) | Slower (borrow checker, longer compiles) |
| Compilation time | <1s for incremental builds | 5-30s+ |
| Ecosystem (AI/HTTP) | Mature `net/http`, JSON, SSE | Growing but less turnkey |
| Binary size | ~10-15MB (stripped) | ~5-10MB (stripped) |

*Last updated: April 2026. Measurements pending.*
