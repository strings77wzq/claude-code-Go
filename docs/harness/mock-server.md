# Mock Server Guide

The mock server simulates the Anthropic API for testing and development purposes.

## Overview

The mock server:
- Accepts the same API requests as the real Anthropic API
- Returns realistic responses based on configured scenarios
- Eliminates API costs during development
- Enables deterministic testing

## Starting the Mock Server

```bash
# From the harness directory
cd harness

# Run the mock server
python -m mock_server
```

The mock server runs on `http://localhost:8080` by default.

## Configuring go-code to Use the Mock Server

Set the API base URL in your config or environment:

```bash
# Environment variable
export ANTHROPIC_API_BASE=http://localhost:8080
export ANTHROPIC_API_KEY=mock-key
```

Or in `~/.config/go-code/config.yaml`:

```yaml
api_key: "mock-key"
api_base: "http://localhost:8080"
```

## Response Scenarios

The mock server can be configured to return different responses:

### Simple Response

Returns a basic text response without tool calls.

### Tool Call Response

Configures the model to request specific tool execution.

### Error Response

Simulates API errors for testing error handling.

## Custom Responses

To customize mock responses, edit the mock server configuration:

```python
# harness/mock_server/__init__.py

# Define custom response scenarios
RESPONSES = {
    "simple": {
        "type": "text",
        "content": "This is a mock response."
    },
    "tool_call": {
        "type": "tool_use",
        "name": "Read",
        "input": {"file_path": "test.txt"}
    }
}
```

## Testing with the Mock Server

1. Start the mock server
2. Configure go-code to use it
3. Run your tests or development sessions

```bash
# Terminal 1: Start mock server
cd harness && python -m mock_server

# Terminal 2: Run go-code
cd go-code
export ANTHROPIC_API_BASE=http://localhost:8080
export ANTHROPIC_API_KEY=mock-key
./bin/go-code "Read test.txt"
```

## Limitations

The mock server:
- Does not implement the full Anthropic API
- Has simplified tool execution (reads actual files but doesn't run commands)
- Cannot replicate complex model behaviors

Use for basic testing and development only. Always test against the real API before release.

## Next Steps

- [Harness Overview](overview.md) - Back to harness documentation
- [Configuration](../guide/config.md) - Configure go-code