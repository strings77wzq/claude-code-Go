---
title: API Errors
description: How to read and resolve errors from LLM providers (Anthropic, OpenAI, and compatible).
---

# API Errors

This page covers error types returned by LLM providers and what to do when something goes wrong.

## Error Classification

Errors are normalized into categories in `internal/provider/errors.go`:

| Kind | HTTP Status | Meaning |
| --- | --- | --- |
| `auth` | 401, 403 | Authentication or permission failure. |
| `rate_limit` | 429 | Too many requests. |
| `invalid_request` | 400, 422 | Malformed request payload. |
| `server` | 500+ | Temporary provider-side failure. |
| `timeout` | (none) | Request exceeded deadline. |
| `network` | (none) | DNS, dial, or connection failure. |

## Authentication Errors

**Message**: `Invalid API key. Please check your ANTHROPIC_API_KEY.`

Causes: missing key (`ErrAPIKeyRequired` from `internal/config`), wrong format (Anthropic keys start with `sk-ant-`), expired or revoked key.

```bash
echo $ANTHROPIC_API_KEY                      # verify set
echo ${ANTHROPIC_API_KEY:0:7}                # check prefix
curl -H "x-api-key: $ANTHROPIC_API_KEY" \
  https://api.anthropic.com/v1/models        # test directly
```

**Message**: `API access denied. Check your API key permissions.`

The key exists but lacks scope for the requested resource. Regenerate with correct permissions at the provider console.

## Rate Limiting (429)

**Message**: `Rate limited. Retrying automatically...`

The application retries automatically with exponential backoff (from `internal/api/client.go`):
- 3 retries with delays of 1s, 2s, 4s.
- After all retries fail: `request failed after 3 retries: rate_limit_exceeded`.

Reduce concurrent usage or upgrade your provider plan.

## Server Errors (5xx)

**Message**: `Server error. Please try again later.`

Retried automatically (same 3-attempt backoff). If persistent, check the provider status page or switch models with `/model`.

## Timeout Errors

**Message**: `provider request timed out`

The HTTP client timeout is 5 minutes (from `internal/api/client.go`). Try a faster model or check your network connection.

## Invalid Request (400, 422)

**Message format**: `Invalid provider request (status): body`

Causes: context length exceeded, invalid model name, malformed schema. If `context_length_exceeded`, use `/compact`.

## Reading Provider Errors

Errors from the Anthropic HTTP client are classified by `classifyError()` in `internal/api/client.go`. Provider adapters in `internal/provider/anthropic` and `internal/provider/openai` add a second classification layer:

```go
func ClassifyHTTPStatus(statusCode int, body string) *ClassifiedError
func ClassifyError(err error) *ClassifiedError
```

OpenAI-compatible providers transform their error payloads into the same categories before reaching the agent loop.

## Related

- [Common Issues](common-issues) — API key setup, connection issues
- [Provider Configuration](../architecture/providers) — Provider setup
