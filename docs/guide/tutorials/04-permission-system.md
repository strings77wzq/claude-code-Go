# Tutorial 4: Permission System Deep Dive

Learn how claude-code-Go keeps you safe.

## The Three Tiers

claude-code-Go uses a 3-tier permission model:

```
ReadOnly < WorkspaceWrite < DangerFullAccess
  Safe       Default          Power User
```

### Tier 1: ReadOnly
**What it can do:**
- ✅ Read any file
- ❌ Cannot write files
- ❌ Cannot execute commands

**Best for:** Exploring codebases, learning, documentation

### Tier 2: WorkspaceWrite (Default)
**What it can do:**
- ✅ Read files
- ✅ Write files in workspace
- ✅ Execute safe commands (ls, cat, grep)
- ⚠️ Prompts before dangerous operations

**Best for:** Daily development, coding tasks

### Tier 3: DangerFullAccess
**What it can do:**
- ✅ Everything
- ⚠️ Minimal prompting
- ⚠️ Can delete files, run any command

**Best for:** Trusted automation, CI/CD, expert users

## Permission Modes in Action

### Switching Modes

```bash
# Check current mode
> /mode
Current mode: WorkspaceWrite

# Switch to ReadOnly
> /mode ReadOnly
Mode changed to: ReadOnly

# Switch to DangerFullAccess
> /mode DangerFullAccess
⚠️  WARNING: You are entering DangerFullAccess mode.
The AI will be able to execute any command without confirmation.

Type "I understand" to continue: I understand
Mode changed to: DangerFullAccess
```

## Glob Rules

Glob rules provide fine-grained control:

```json
{
  "mode": "WorkspaceWrite",
  "rules": [
    {"pattern": "*.go", "allowed": true},
    {"pattern": "*.md", "allowed": true},
    {"pattern": "*.env", "allowed": false},
    {"pattern": ".ssh/*", "allowed": false},
    {"pattern": "/etc/*", "allowed": false}
  ]
}
```

### Rule Priority

Rules are evaluated in order. First match wins:

```
"*.env" → DENY    # Sensitive files blocked
"*.go"  → ALLOW   # Go files allowed
"*"     → ASK     # Everything else requires permission
```

## Session Memory

Grant permissions for the entire session:

```bash
# One-time grant
> /allow read secret.txt
✅ Granted read permission for secret.txt

# Pattern-based grant
> /allow write *.json
✅ Granted write permission for *.json files

# Persist across sessions
> /remember allow bash
✅ Permission remembered in settings
```

## Permission Prompts

When the AI needs permission, you'll see:

```
> Delete the build directory

⚠️  Permission Required
The tool "Bash" wants to execute:
  rm -rf build/

This is a destructive operation that cannot be undone.

Options:
1. [Y] Yes - Allow this operation
2. [N] No  - Cancel this operation  
3. [A] Always - Allow all similar operations this session
4. [M] Mode - Switch to DangerFullAccess

Choice: 
```

## Common Scenarios

### Scenario 1: Reading Sensitive Files

```
> Read .env

❌ Permission Denied
Reading .env files is blocked by default.

To allow this:
1. Use "/allow read .env" for one-time access
2. Or add to settings.json rules
3. Or switch to ReadOnly mode (can read anything)
```

### Scenario 2: Running Tests

```
> Run go test ./...

🛠️  Using tool: Bash
   command: go test ./...

✅ Tests passed (safe command, auto-allowed)
```

### Scenario 3: Installing Dependencies

```
> Install the package

⚠️  Permission Required
The tool "Bash" wants to execute:
  go get github.com/example/package

This will modify go.mod and download code.

[Y/n/a/m] y
✅ Package installed
```

## Best Practices

### 1. Start Conservative

Begin with `ReadOnly` mode when exploring new codebases.

### 2. Use WorkspaceWrite for Development

The default mode is safe for most coding tasks.

### 3. Reserve DangerFullAccess for Automation

Only use `DangerFullAccess` in:
- CI/CD pipelines
- Scripts you fully control
- When you're an expert user

### 4. Review Before Allowing

Always read what the AI wants to do before granting permission.

### 5. Use Session Memory

If you trust a pattern, use `/remember` to avoid repeated prompts.

## Configuration

Edit `~/.go-code/settings.json`:

```json
{
  "mode": "WorkspaceWrite",
  "rules": [
    {"pattern": "*.go", "allowed": true},
    {"pattern": "*.md", "allowed": true},
    {"pattern": "*.env", "allowed": false},
    {"pattern": ".git/*", "allowed": false},
    {"pattern": "node_modules/*", "allowed": false}
  ],
  "sessionMemory": [
    "allow bash go test",
    "allow bash npm install"
  ]
}
```

## Troubleshooting

### "Permission denied" errors

```
> Check if the file exists

❌ Permission Denied for tool: Bash
```

**Solution**: The AI tried to run a command that requires permission. Grant it with `/allow` or switch modes.

### Too many prompts

**Solution**: Use `/mode DangerFullAccess` temporarily, or add common patterns to session memory.

### Can't read a file you need

**Solution**: Use `/allow read filename` or switch to `ReadOnly` mode (can read anything but not write).

## Security Checklist

- [ ] Review your glob rules regularly
- [ ] Don't commit `.go-code/settings.json` with DangerFullAccess
- [ ] Be cautious with `/remember` commands
- [ ] Review what the AI wants to do before allowing
- [ ] Use ReadOnly mode when exploring untrusted code

## Next Steps

- [Tutorial 5: Building Custom Tools](05-custom-tools.md)
- [Architecture: Permission Model](../../architecture/permission-model-deep-dive.md)
- [ADR-002: Permission Model](../../adr/002-permission-model.md)
