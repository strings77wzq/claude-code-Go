# API Errors

Handling errors from AI providers.

## Rate Limiting

### "rate_limit_exceeded"

**Cause**: Too many requests.

**Solution**:
- Wait 60 seconds
- Upgrade your plan
- Reduce request frequency

### "rate_limit_request_enqueued"

**Cause**: Request queued due to high load.

**Solution**:
- Wait for automatic retry
- Reduce concurrent requests

## Authentication

### "invalid_api_key"

**Cause**: Wrong API key.

**Solution**:
- Check key format (sk-ant-...)
- Generate new key
- Verify key in settings

### "permission_denied"

**Cause**: Key lacks required permissions.

**Solution**:
- Check key scope
- Generate key with correct permissions
- Contact Anthropic support

## Content Errors

### "context_length_exceeded"

**Cause**: Too much context.

**Solution**:
```
> /compact
> Let's continue with the compacted context
```

### "invalid_request_error"

**Cause**: Malformed request.

**Solution**:
- Update to latest version
- Check settings.json format
- Report bug if persistent

## Server Errors

### "api_error" or "server_error"

**Cause**: Provider issue.

**Solution**:
- Wait and retry
- Check provider status page
- Try different model

## Retry Logic

claude-code-Go automatically retries:
- Rate limits (with exponential backoff)
- Timeouts (up to 3 attempts)
- Server errors (5xx)

Manual retry:
```
> That failed, please try again
```
