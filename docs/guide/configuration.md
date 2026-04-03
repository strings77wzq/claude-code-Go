---
title: Configuration Guide
description: Complete guide to configuring go-code — config files, environment variables, and MCP servers
---

# Configuration Guide

This guide covers all configuration options for go-code, including config files, environment variables, and MCP server setup.

## Config File Locations

go-code loads configuration from multiple locations with the following priority (highest to lowest):

1. **CLI arguments** (highest priority)
2. **Environment variables**
3. **Project config file**: `./.go-code/settings.json`
4. **User config file**: `~/.go-code/settings.json`
5. **Built-in defaults** (lowest priority)

This means you can set defaults in the user config and override them per-project, via environment variables, or via CLI arguments.

## Configuration File Format

Configuration files use JSON format:

```json
{
  "apiKey": "sk-ant-your-api-key-here",
  "baseUrl": "https://api.anthropic.com",
  "model": "claude-sonnet-4-20250514"
}
```

### Settings Options

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `apiKey` | string | (required) | API key |
| `baseUrl` | string | `https://api.anthropic.com` | API endpoint URL |
| `model` | string | `claude-sonnet-4-20250514` | Model to use |

## Environment Variables

| Variable | Description |
|----------|-------------|
| `ANTHROPIC_API_KEY` | Your API key |
| `ANTHROPIC_BASE_URL` | Override the default API endpoint |

### Example: Setting Environment Variables

```bash
# Set API key
export ANTHROPIC_API_KEY=sk-ant-your-api-key-here

# Optionally override base URL (for testing or proxy)
export ANTHROPIC_BASE_URL=https://custom-api.example.com
```

Add these to your shell profile for persistence:

```bash
# ~/.bashrc or ~/.zshrc
echo 'export ANTHROPIC_API_KEY=sk-ant-xxx' >> ~/.bashrc
source ~/.bashrc
```

## MCP Server Configuration

go-code supports the Model Context Protocol (MCP) for extending capabilities with external tools and services.

### MCP Config File Location

MCP server configurations are stored in:

```
~/.go-code/mcp.json
```

### MCP Configuration Format

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

### MCP Config Structure

Each MCP server entry has:

| Field | Type | Description |
|-------|------|-------------|
| `command` | string | The executable to run |
| `args` | array | Command-line arguments |
| `env` | object | Environment variables (supports `${VAR}` interpolation) |

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

## Full Configuration Example

Here's a complete example combining all configuration options:

### User Config (~/.go-code/settings.json)

```json
{
  "apiKey": "sk-ant-your-api-key-here",
  "baseUrl": "https://api.anthropic.com",
  "model": "claude-sonnet-4-20250514"
}
```

### Project Config (./.go-code/settings.json)

```json
{
  "model": "claude-opus-4-20250514"
}
```

This project-specific config overrides the model while using the API key from the user config.

### MCP Config (~/.go-code/mcp.json)

```json
{
  "filesystem": {
    "command": "npx",
    "args": ["-y", "@modelcontextprotocol/server-filesystem", "./workspace"]
  },
  "git": {
    "command": "uvx",
    "args": ["mcp-server-git"],
    "env": {
      "GIT_TOKEN": "${GIT_TOKEN}"
    }
  }
}
```

## Runtime Model Switching

You can switch models at runtime without restarting go-code.

### Using `/model` Command

```
go-code> /model
Current model: claude-sonnet-4-20250514

go-code> /model hunyuan-2.0-instruct
Model switched to: hunyuan-2.0-instruct
```

### Using `/models` Command

List all available models:

```
go-code> /models
Available models:

  Anthropic:
    claude-sonnet-4-20250514 (default)
    claude-opus-4-20250514
    claude-haiku-4-20250514

  Tencent Coding Plan:
    tc-code-latest (Auto)
    hunyuan-2.0-instruct
    hunyuan-2.0-thinking
    minimax-m2.5
    kimi-k2.5
    glm-5
    hunyuan-t1
    hunyuan-turbos

Switch model: /model <model-name>
```

### Tencent Coding Plan Setup

For Tencent Coding Plan, configure:

```bash
export ANTHROPIC_API_KEY="sk-sp-your-key"
export ANTHROPIC_BASE_URL="https://api.lkeap.cloud.tencent.com/coding/anthropic"
export ANTHROPIC_MODEL="tc-code-latest"
```

Or in `~/.go-code/settings.json`:

```json
{
  "apiKey": "sk-sp-your-key",
  "baseUrl": "https://api.lkeap.cloud.tencent.com/coding/anthropic",
  "model": "tc-code-latest"
}
```

## Troubleshooting

### "API key is required" Error

Ensure you've set either:
- `ANTHROPIC_API_KEY` environment variable, or
- `apiKey` in a config file

### Config File Not Found

Verify the config file exists and has valid JSON:

```bash
# Validate JSON syntax
cat ~/.go-code/settings.json | python -m json.tool
```

### MCP Server Not Loading

Check:
1. MCP config file exists at `~/.go-code/mcp.json`
2. The command executable is in your PATH
3. Required dependencies are installed

## Related Documentation

- [Quick Start Guide](quick-start.md) — Basic setup and first run
- [Architecture Overview](../architecture/overview.md) — System components
- [Tool System](../architecture/tools.md) — Built-in and MCP tools
