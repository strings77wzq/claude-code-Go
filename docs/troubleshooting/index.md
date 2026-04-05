---
title: Troubleshooting Guide
description: Common issues and solutions for go-code
---

# Troubleshooting Guide

This guide covers common issues you might encounter when using go-code and provides solutions to resolve them.

---

## Common Issues

### API Key Errors

#### "API key is required"

**Problem:** The application cannot find a valid API key.

**Solutions:**

1. **Set environment variable:**
   ```bash
   export ANTHROPIC_API_KEY=sk-ant-your-key-here
   ```

2. **Create config file:**
   ```bash
   mkdir -p ~/.go-code
   echo '{"apiKey": "sk-ant-your-key-here"}' > ~/.go-code/settings.json
   ```

3. **Verify the key is set:**
   ```bash
   echo $ANTHROPIC_API_KEY
   ```

#### "Invalid API key"

**Problem:** The API key format is incorrect.

**Solutions:**

1. **Check key format:** Anthropic keys start with `sk-ant-`
2. **Regenerate key:** Visit [Anthropic Console](https://console.anthropic.com/)
3. **Check for extra characters:** Ensure no quotes or spaces in the key

---

### Network Errors

#### "Connection timeout"

**Problem:** The request to the API timed out.

**Solutions:**

1. **Check internet connection:**
   ```bash
   ping api.anthropic.com
   ```

2. **Increase timeout** (via config):
   ```json
   { "timeout": 300 }
   ```

3. **Check firewall/proxy** settings
4. **Try a different network** to diagnose

#### "Network error: dial tcp"

**Problem:** Cannot establish TCP connection to the API.

**Solutions:**

1. **Verify base URL** is correct
2. **Check DNS resolution:**
   ```bash
   nslookup api.anthropic.com
   ```
3. **Disable VPN/proxy** if causing issues
4. **Check corporate firewall** restrictions

---

### Permission Denied

#### "Permission denied" for file operations

**Problem:** The tool cannot access or modify a file.

**Solutions:**

1. **Check file permissions:**
   ```bash
   ls -la /path/to/file
   ```

2. **Verify path is within working directory:**
   - go-code restricts file access to the working directory tree
   - Ensure the file path is relative or within the project

3. **Grant execute permission for Bash:**
   ```bash
   chmod +x /path/to/script.sh
   ```

---

### Model Not Found

#### "Model not found: xxx"

**Problem:** The specified model is not available.

**Solutions:**

1. **List available models:**
   ```
   /models
   ```

2. **Use a different model:**
   ```
   /model claude-sonnet-4-20250514
   ```

3. **Check model name spelling**
4. **Verify API subscription** includes the model

---

### Session Errors

#### "Failed to save session"

**Problem:** Cannot save the current session.

**Solutions:**

1. **Check sessions directory exists:**
   ```bash
   mkdir -p ~/.go-code/sessions
   ```

2. **Check write permissions:**
   ```bash
   ls -la ~/.go-code/
   ```

3. **Check disk space**

#### "Session not found"

**Problem:** Cannot resume the specified session.

**Solutions:**

1. **List available sessions:**
   ```
   /sessions
   ```

2. **Check session file exists:**
   ```bash
   ls ~/.go-code/sessions/
   ```

3. **Session may have been deleted**

---

## Error Code Reference

| Error Code | Description | Solution |
|------------|-------------|----------|
| `E001` | API key required | Set `ANTHROPIC_API_KEY` environment variable |
| `E002` | Invalid API key | Verify key format and regenerate if needed |
| `E003` | Network timeout | Check connection and increase timeout |
| `E004` | Permission denied | Check file permissions and working directory |
| `E005` | Model not found | Use `/models` to list available models |
| `E006` | Session not found | Use `/sessions` to list valid session IDs |
| `E007` | File too large | Maximum file size is 200KB |
| `E008` | Invalid JSON | Check configuration file syntax |
| `E009` | MCP server error | Check MCP server configuration |
| `E010` | Context exceeded | Use `/compact` to compress context |

---

## FAQ

### General

**Q: How do I get an API key?**
> Visit [Anthropic Console](https://console.anthropic.com/) and create an API key.

**Q: Why is my response slow?**
> Check your network connection. Large responses or complex tasks take longer.

**Q: Can I use go-code offline?**
> No, go-code requires an API connection for the model to generate responses.

**Q: Where are sessions saved?**
> Sessions are saved in `~/.go-code/sessions/` by default.

### Configuration

**Q: How do I switch models?**
> Use the `/model` command in the REPL:
> ```
> /model claude-opus-4-20250514
> ```

**Q: How do I set a different base URL?**
> Set `ANTHROPIC_BASE_URL` environment variable or `baseUrl` in config.

**Q: Can I use multiple API keys?**
> Not simultaneously. You can switch keys by updating the config.

### Tools

**Q: Why can't I read a file?**
> Check that the file is within the working directory and under 200KB.

**Q: Why was my write operation blocked?**
> The tool requires permission approval. Confirm when prompted.

**Q: How do I use MCP servers?**
> Configure MCP servers in `~/.go-code/mcp.json`. See [MCP Integration](../extension/mcp.md).

### Troubleshooting

**Q: How do I debug issues?**
> Set `GO_CODE_TRACE=true` environment variable for detailed logging.

**Q: How do I report a bug?**
> Open an issue on [GitHub](https://github.com/strings77wzq/claude-code-Go/issues) with details.

**Q: Where are logs stored?**
> Logs are written to stdout in the REPL. Use trace mode for more detail.

---

## Getting Help

If you encounter an issue not covered here:

1. **Check the documentation** — See related docs at the bottom of this page
2. **Search existing issues** — [GitHub Issues](https://github.com/strings77wzq/claude-code-Go/issues)
3. **Ask on discussions** — [GitHub Discussions](https://github.com/strings77wzq/claude-code-Go/discussions)
4. **Report a bug** — Include: Go version, OS, error message, and reproduction steps

---

## Related Documentation

- [Configuration Guide](../guide/configuration.md) — Configuration options
- [Tool System](../tools/overview.md) — Built-in tools
- [MCP Integration](../extension/mcp.md) — MCP servers
- [Session Management](../guide/session-management.md) — Session persistence