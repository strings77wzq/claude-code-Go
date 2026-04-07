# Performance Issues

Optimizing claude-code-Go performance.

## Slow Startup

### Symptoms
- Takes > 5 seconds to start
- Shows spinner for long time

### Solutions

**Check binary location**:
```bash
# Should be native binary, not 'go run'
file $(which go-code)
# Expected: ELF 64-bit executable
```

**Clear old sessions**:
```bash
rm ~/.go-code/sessions/*.jsonl
```

**Disable unused features**:
```json
{
  "enableMCP": false,
  "enableLSP": false
}
```

## High Memory Usage

### Symptoms
- Memory usage > 500MB
- System slowing down
- OOM kills

### Solutions

**Check memory usage**:
```bash
ps aux | grep go-code
```

**Compact context**:
```
> /compact
```

**Start fresh**:
```
> /clear
```

**Reduce history size**:
```json
{
  "maxHistoryMessages": 20
}
```

## Slow Responses

### Symptoms
- Long delays between messages
- Timeouts

### Solutions

**Check network**:
```bash
ping api.anthropic.com
```

**Increase timeout**:
```json
{
  "timeout": "60s"
}
```

**Switch model**:
```
> /model claude-haiku-4-6-20251001  # Faster model
```

**Simplify requests**:
- Break complex tasks into smaller ones
- Clear context with `/clear`
- Use specific file paths

## Token Usage

### Monitor usage

```
> /tokens

Context: 45,234 / 100,000 tokens (45%)
```

### Reduce usage

1. **Compact regularly**: `/compact`
2. **Start fresh**: `/clear`
3. **Be specific**: "Read main.go" not "Read all files"

### Auto-compact

```json
{
  "autoCompactThreshold": 0.7  // Compact at 70%
}
```

## Benchmarking

Measure performance:

```bash
# Startup time
time go-code -p "Hello"

# Memory usage
/usr/bin/time -v go-code -p "Hello"
```
