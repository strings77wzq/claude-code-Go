# Common Issues

Solutions to frequently encountered problems.

## Installation Issues

### "command not found: go-code"

**Problem**: The binary is not in your PATH.

**Solution**:
```bash
# Check where the binary is
which go-code

# If not found, add to PATH
export PATH="$HOME/go/bin:$PATH"

# For permanent fix, add to ~/.bashrc or ~/.zshrc
echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.bashrc
```

### "go install" fails with permission denied

**Problem**: No write permission to GOPATH.

**Solution**:
```bash
# Check GOPATH
go env GOPATH

# Change ownership
sudo chown -R $(whoami) $(go env GOPATH)

# Or install to different location
go build -o ~/bin/go-code ./cmd/go-code
```

## Configuration Issues

### "ANTHROPIC_API_KEY not set"

**Problem**: API key not configured.

**Solution**:
```bash
# Quick fix
export ANTHROPIC_API_KEY="sk-ant-..."

# Permanent fix
echo 'export ANTHROPIC_API_KEY="sk-ant-..."' >> ~/.bashrc

# Or create settings file
mkdir -p ~/.go-code
cat > ~/.go-code/settings.json << EOF
{
  "apiKey": "sk-ant-..."
}
EOF
```

### Settings file not loading

**Problem**: JSON syntax error or wrong location.

**Solution**:
```bash
# Check file location
ls -la ~/.go-code/settings.json

# Validate JSON
cat ~/.go-code/settings.json | python3 -m json.tool

# Check permissions
chmod 600 ~/.go-code/settings.json
```

## Connection Issues

### "connection refused" or timeout

**Problem**: Cannot connect to API.

**Solution**:
```bash
# Check internet
curl https://api.anthropic.com/v1/health

# Check firewall
# Some corporate networks block API calls

# Try different network
# Use mobile hotspot to test
```

### "invalid api key"

**Problem**: API key is wrong or expired.

**Solution**:
```bash
# Verify key format (should start with sk-ant-)
echo $ANTHROPIC_API_KEY

# Test key
curl -H "x-api-key: $ANTHROPIC_API_KEY" \
  https://api.anthropic.com/v1/models

# Generate new key at:
# https://console.anthropic.com/settings/keys
```

## Runtime Issues

### Sessions not saving

**Problem**: Cannot persist conversations.

**Solution**:
```bash
# Check directory exists
mkdir -p ~/.go-code/sessions

# Check disk space
df -h ~/.go-code

# Check permissions
ls -la ~/.go-code/sessions
```

### Tool execution fails

**Problem**: Tools return errors.

**Solution**:
1. Check current mode: `/mode`
2. Verify permissions: `/rules`
3. Check if tool exists: `/tools`
4. Review error message for specific issue

### Out of memory

**Problem**: Process killed or hangs.

**Solution**:
```bash
# Check memory usage
ps aux | grep go-code

# Clear old sessions
rm ~/.go-code/sessions/*.jsonl

# Compact current session
> /compact

# Start fresh
> /clear
```

## Getting More Help

1. Check [API Errors](api-errors.md)
2. Check [Permission Denied](permission-denied.md)
3. Check [Performance Issues](performance-issues.md)
4. Open [GitHub Discussion](https://github.com/strings77wzq/claude-code-Go/discussions)
5. Search [GitHub Issues](https://github.com/strings77wzq/claude-code-Go/issues)
