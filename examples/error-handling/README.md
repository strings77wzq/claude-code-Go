# Error Handling Best Practices

This guide covers best practices for handling errors in claude-code-Go.

## Error Types

### 1. API Errors

API errors occur when communicating with the LLM provider.

**Common API Errors:**

| Error Code | Cause | Solution |
|------------|-------|----------|
| `rate_limit_exceeded` | Too many requests | Wait and retry |
| `invalid_api_key` | Wrong API key | Check configuration |
| `timeout` | Request took too long | Increase timeout or simplify |
| `context_length_exceeded` | Too much context | Compact session |
| `model_unavailable` | Model temporarily down | Try different model |

**Best Practice:**
```go
if apiErr, ok := err.(*api.Error); ok {
    switch apiErr.Code {
    case "rate_limit_exceeded":
        time.Sleep(60 * time.Second)
        return retry()
    case "timeout":
        return retryWithLongerTimeout()
    default:
        return fmt.Errorf("API error: %w", err)
    }
}
```

### 2. Permission Errors

Permission errors occur when a tool operation violates security rules.

**Handling Permission Errors:**

1. **Inform the user clearly** what was blocked and why
2. **Provide alternatives:**
   - Use `/allow` to grant one-time permission
   - Switch mode with `/mode`
   - Update glob rules in settings

**Best Practice:**
```go
if permErr, ok := err.(*permission.Error); ok {
    fmt.Println("Permission denied:", permErr.Resource)
    fmt.Println("Run '/allow read", permErr.Resource, "' to grant access")
}
```

### 3. Tool Execution Errors

Tool execution errors occur when a tool fails to complete its operation.

**Common Tool Errors:**

- `ErrCommandFailed`: Command returned non-zero exit
- `ErrTimeout`: Tool execution timed out
- `ErrInvalidInput`: Missing or invalid parameters
- `ErrNotFound`: Target file/resource not found

**Best Practice:**
```go
if toolErr, ok := err.(*tool.ExecutionError); ok {
    switch toolErr.Type {
    case tool.ErrTimeout:
        // Retry with longer timeout
    case tool.ErrNotFound:
        // Check if file exists first
    }
}
```

### 4. Timeout Errors

Timeout errors occur when operations exceed their time limit.

**Strategies:**

1. **Retry with exponential backoff**
2. **Break complex tasks into smaller steps**
3. **Use async execution for long-running tasks**
4. **Increase timeout in settings**

**Example:**
```go
for i := 0; i < 3; i++ {
    result, err := operation()
    if err == nil {
        return result, nil
    }
    if !isTimeout(err) {
        return nil, err // Don't retry non-timeout errors
    }
    time.Sleep(time.Duration(i+1) * time.Second)
}
```

## Error Recovery Patterns

### 1. Graceful Degradation

When a feature fails, continue with reduced functionality:

```go
if err := loadSkills(); err != nil {
    log.Warn("Skills not loaded, continuing without them")
    // Continue with basic functionality
}
```

### 2. User Confirmation

For recoverable errors, ask the user what to do:

```go
if err := riskyOperation(); err != nil {
    fmt.Printf("Operation failed: %v\n", err)
    fmt.Print("Retry? [y/N]: ")
    if confirm() {
        return retryOperation()
    }
}
```

### 3. Automatic Retry

For transient errors, retry automatically:

```go
retry := retry.New(
    retry.MaxAttempts(3),
    retry.Backoff(retry.ExponentialBackoff),
    retry.RetryableErrors(api.IsTemporaryError),
)

result, err := retry.Do(func() error {
    return api.Call()
})
```

## Logging Errors

Always log errors with context:

```go
log.Errorw("API call failed",
    "error", err,
    "model", config.Model,
    "attempt", attempt,
    "duration", time.Since(start),
)
```

## User-Friendly Error Messages

Convert technical errors to user-friendly messages:

```go
func UserFriendlyError(err error) string {
    if apiErr, ok := err.(*api.Error); ok {
        switch apiErr.Code {
        case "rate_limit_exceeded":
            return "You've hit the rate limit. Please wait a minute before trying again."
        case "timeout":
            return "The request took too long. Try simplifying your question or check your connection."
        }
    }
    return fmt.Sprintf("An error occurred: %v", err)
}
```

## Examples

See the [examples/error-handling/](.) directory for complete examples:

- `api_errors.go` - Handling API errors
- `permission_errors.go` - Handling permission errors
- `tool_errors.go` - Handling tool execution errors
- `timeout_errors.go` - Handling timeout errors
