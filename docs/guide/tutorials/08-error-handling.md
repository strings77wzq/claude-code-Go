# Tutorial 8: Error Handling Best Practices

Learn to handle errors gracefully.

## Common Error Types

### API Errors
Errors from the LLM provider:

```
❌ API Error: rate_limit_exceeded
   You've exceeded the rate limit.
   Wait 60 seconds before retrying.
```

### Permission Errors
Security violations:

```
❌ Permission Denied
   Cannot write to /etc/passwd
   Use /mode DangerFullAccess or /allow write /etc/passwd
```

### Tool Errors
Tool execution failures:

```
❌ Tool Error: Bash
   Command 'npm install' failed with exit code 1
   Check the command and try again
```

## Recovery Strategies

### Retry with Backoff

```
> The request timed out

The AI will automatically retry with:
- Attempt 1: Immediate
- Attempt 2: After 2s
- Attempt 3: After 4s
```

### Context Compaction

```
> I'm getting context_length_exceeded

Running automatic compaction...
- Summarized 20 old messages
- Kept 5 recent messages
- Freed 40,000 tokens

Continue your conversation normally.
```

### Session Reset

When all else fails:

```
> /clear
> Let's start over
```

## Best Practices

1. **Read error messages carefully** - They often include solutions
2. **Don't retry immediately** - Wait for rate limits
3. **Use /compact for context issues** - Don't lose all context
4. **Check your API key** - Common cause of auth errors

## Next Steps

- [Troubleshooting Guide](../troubleshooting/common-issues.md)
- [Examples: Error Handling](../../../examples/error-handling/)
