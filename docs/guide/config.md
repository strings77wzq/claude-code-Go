# Configuration

go-code can be configured via config files or environment variables.

## Config File Location

- Linux/macOS: `~/.config/go-code/config.yaml`
- Windows: `%APPDATA%\go-code\config.yaml`

## Config File Format

```yaml
# API Configuration
api_key: "sk-ant-your-api-key-here"

# Model settings
model: "claude-sonnet-4-20250514"
max_tokens: 4096

# Permission settings
auto_approve_read: true
auto_approve_write: false
auto_approve_bash: false

# MCP servers
mcp_servers:
  filesystem:
    command: "npx"
    args: ["-y", "@modelcontextprotocol/server-filesystem", "/path/to/dir"]
  github:
    command: "python"
    args: ["-m", "mcp.server.github", "--token", "your-token"]

# Logging
log_level: "info"
log_file: "go-code.log"
```

## Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `ANTHROPIC_API_KEY` | Your Anthropic API key | Yes |
| `ANTHROPIC_MODEL` | Model to use | No (default: claude-sonnet-4-20250514) |
| `ANTHROPIC_MAX_TOKENS` | Max response tokens | No (default: 4096) |
| `GO_CODE_LOG_LEVEL` | Log level (debug, info, warn, error) | No (default: info) |
| `GO_CODE_CONFIG` | Path to config file | No |

## Precedence

Configuration is loaded in this order (later overrides earlier):

1. Default values
2. Config file
3. Environment variables
4. Command-line flags

## Command-Line Flags

```
--help      Show help message
--version   Show version
--model     Specify model (overrides config)
--config    Specify config file path
```

## Model Options

Available models (subject to API availability):

- `claude-sonnet-4-20250514` (default)
- `claude-opus-4-20250514`
- `claude-3-5-sonnet-20241022`

## MCP Configuration

Model Context Protocol servers can be configured in the config file:

```yaml
mcp_servers:
  server_name:
    command: "path/to/executable"
    args: ["arg1", "arg2"]
    env:
      KEY: "value"
```

See [MCP Integration](../architecture/tools.md) for details.