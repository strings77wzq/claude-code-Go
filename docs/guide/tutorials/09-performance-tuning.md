# Tutorial 9: Performance Tuning

Optimize claude-code-Go for speed and efficiency.

## Startup Performance

### Measure startup time

```bash
time go-code -p "Hello"
```

Typical results:
- Cold start: ~50ms
- Warm start: ~10ms

### Reduce startup time

1. **Use native binary** (not `go run`)
2. **Smaller context** - Clear old sessions
3. **Disable unused features** in settings

## Context Management

### Token Estimation

```
> /tokens

Current context: 15,234 tokens
Budget: 100,000 tokens
Usage: 15%
```

### Compaction Timing

```
# Manual compaction
> /compact

# Auto-compact threshold in settings.json
{
  "autoCompactThreshold": 0.8  // Compact at 80% capacity
}
```

## Network Optimization

### Connection Pooling

Settings automatically reuse connections.

### Timeout Configuration

```json
{
  "timeout": "30s",
  "maxRetries": 3
}
```

## Memory Usage

### Check memory

```bash
ps aux | grep go-code
```

Typical: 10-50 MB

### Reduce memory

1. Compact sessions
2. Clear history
3. Disable unused MCP servers

## Next Steps

- [Tutorial 10: Contributing](10-contributing.md)
- [Benchmarks](../../benchmark.md)
