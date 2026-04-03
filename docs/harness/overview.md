# Python Harness Overview

The Python harness is an optional component that provides testing and development utilities for go-code.

## Purpose

The harness serves several purposes:

1. **Mock Server**: Provides a mock Anthropic API server for testing without API costs
2. **Replay Tests**: Record and replay agent sessions for debugging
3. **Integration Tests**: Test the full agent loop with realistic scenarios

## Installation

The harness requires Python 3.x:

```bash
# Navigate to harness directory
cd harness

# Install dependencies (if requirements.txt exists)
pip install -r requirements.txt
```

## Components

```
harness/
├── mock_server/       # Mock Anthropic API server
│   └── __init__.py    # Server implementation
├── tests/             # Test suite (if present)
└── pytest.ini         # Pytest configuration (if present)
```

## Running Tests

```bash
# Run Python tests only
cd harness && python -m pytest -v

# Or use the Makefile (runs both Go and Python tests)
make test
```

The Makefile will automatically detect if the harness exists and has tests.

## When to Use

- **Development**: Use the mock server during development to avoid API costs
- **CI/CD**: Use for continuous integration testing
- **Debugging**: Replay tests help reproduce and debug issues

## When Not to Use

The harness is optional. For production use, you only need the Go binary and API key.

## Related Documentation

- [Mock Server](mock-server.md) - Detailed mock server guide
- [Architecture Overview](../architecture/overview.md) - Main architecture