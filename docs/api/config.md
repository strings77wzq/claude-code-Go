---
title: Configuration Reference
description: Complete configuration reference for go-code including environment variables, settings.json schema, and priority chain
---

# Configuration Reference

This document provides a complete reference for all configuration options in go-code, including environment variables, configuration file settings, and the priority chain for configuration sources.

## Configuration Priority

go-code loads configuration from multiple sources with the following priority (highest to lowest):

```
1. CLI arguments (highest priority)
   ↓
2. Environment variables
   ↓
3. Project config file: ./.go-code/settings.json
   ↓
4. User config file: ~/.go-code/settings.json
   ↓
5. Built-in defaults (lowest priority)
```

This means you can set defaults in the user config and override them per-project, via environment variables, or via CLI arguments.

---

## Environment Variables

### API Configuration

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `ANTHROPIC_API_KEY` | Yes | - | API key for authentication |
| `ANTHROPIC_BASE_URL` | No | `https://api.anthropic.com` | Override the default API endpoint |
| `ANTHROPIC_MODEL` | No | `claude-sonnet-4-20250514` | Default model to use |

### MCP Configuration

| Variable | Required | Description |
|----------|----------|-------------|
| `MCP_SERVER_*` | No | Server-specific environment variables |

### Session Configuration

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `GO_CODE_SESSIONS_DIR` | No | `~/.go-code/sessions/` | Directory for session storage |

### Update Configuration

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `GO_CODE_UPDATE_URL` | No | GitHub releases | URL to check for updates |

### Debugging

| Variable | Required | Description |
|----------|----------|-------------|
| `GO_CODE_TRACE` | No | Enable trace logging |
| `GO_CODE_DEBUG` | No | Enable debug mode |

---

## settings.json Schema

The configuration file uses JSON format. Below is the complete schema:

### Root Object

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "properties": {
    "apiKey": {
      "type": "string",
      "description": "API key for authentication"
    },
    "baseUrl": {
      "type": "string",
      "description": "API endpoint URL",
      "default": "https://api.anthropic.com"
    },
    "model": {
      "type": "string",
      "description": "Default model to use",
      "default": "claude-sonnet-4-20250514"
    },
    "maxTokens": {
      "type": "integer",
      "description": "Maximum tokens per response",
      "default": 4096
    },
    "temperature": {
      "type": "number",
      "description": "Sampling temperature (0-1)",
      "default": 0.7
    },
    "timeout": {
      "type": "integer",
      "description": "Request timeout in seconds",
      "default": 120
    },
    "sessionsDir": {
      "type": "string",
      "description": "Directory for session storage",
      "default": "~/.go-code/sessions"
    },
    "autoSave": {
      "type": "boolean",
      "description": "Auto-save sessions",
      "default": true
    },
    "maxHistorySize": {
      "type": "integer",
      "description": "Maximum history messages to keep",
      "default": 100
    }
  }
}
```

### Example Configuration

```json
{
  "apiKey": "sk-ant-your-api-key-here",
  "baseUrl": "https://api.anthropic.com",
  "model": "claude-sonnet-4-20250514",
  "maxTokens": 8192,
  "temperature": 0.7,
  "timeout": 180,
  "sessionsDir": "~/.go-code/sessions",
  "autoSave": true,
  "maxHistorySize": 50
}
```

---

## Configuration Files

### User Configuration

Location: `~/.go-code/settings.json`

This is the user-level configuration that applies to all sessions for the current user.

```json
{
  "apiKey": "sk-ant-your-api-key-here",
  "model": "claude-opus-4-20250514"
}
```

### Project Configuration

Location: `./.go-code/settings.json` (in project root)

This is the project-level configuration that overrides user settings for the current project.

```json
{
  "model": "claude-haiku-4-20250514"
}
```

### Priority Example

Given:
- User config: `{ "apiKey": "user-key", "model": "sonnet" }`
- Project config: `{ "model": "opus" }`

Result:
- API key: `user-key` (from user config, not overridden)
- Model: `opus` (from project config, overrides user)

---

## MCP Configuration

### Location

`~/.go-code/mcp.json`

### Schema

```json
{
  "type": "object",
  "additionalProperties": {
    "type": "object",
    "properties": {
      "command": {
        "type": "string",
        "description": "The executable to run"
      },
      "args": {
        "type": "array",
        "items": {
          "type": "string"
        },
        "description": "Command-line arguments"
      },
      "env": {
        "type": "object",
        "additionalProperties": {
          "type": "string"
        },
        "description": "Environment variables (supports ${VAR} interpolation)"
      }
    },
    "required": ["command"]
  }
}
```

### Example

```json
{
  "filesystem": {
    "command": "npx",
    "args": ["-y", "@modelcontextprotocol/server-filesystem", "/path/to/directory"],
    "env": {
      "HOME": "${HOME}"
    }
  },
  "github": {
    "command": "uvx",
    "args": ["mcp-server-github"],
    "env": {
      "GITHUB_TOKEN": "${GITHUB_TOKEN}"
    }
  }
}
```

### Environment Variable Interpolation

MCP config supports `${VAR}` syntax to interpolate environment variables:

```json
{
  "server": {
    "command": "my-server",
    "env": {
      "API_KEY": "${ANTHROPIC_API_KEY}",
      "HOME": "${HOME}"
    }
  }
}
```

This allows sensitive credentials to be passed from the host environment without storing them in the config file.

---

## CLI Arguments

| Flag | Type | Description |
|------|------|-------------|
| `-p` | string | Single prompt to execute |
| `-f` | string | Output format: text or json |
| `-q` | bool | Quiet mode (no spinner) |
| `-m` | string | Model to use |
| `-c` | string | Config file path |

### Examples

```bash
# Single prompt
go-code -p "Explain the code"

# JSON output
go-code -p "List files" -f json

# Quiet mode
go-code -p "What is 2+2?" -q

# Specific model
go-code -m claude-opus-4-20250514
```

---

## Permission System Configuration

go-code uses a 3-tier permission system. Configuration is handled internally but can be controlled via:

1. **Session memory**: Permission decisions are remembered during a session
2. **Glob rules**: File path patterns that grant automatic permission
3. **Interactive prompts**: User is prompted for permission when needed

See [Permission System](../architecture/tools.md#permission-system) for details.

---

## Related Documentation

- [Configuration Guide](../guide/configuration.md) — User-friendly configuration guide
- [Tool System](../tools/overview.md) — Tool execution and permissions
- [MCP Integration](../extension/mcp.md) — MCP server configuration