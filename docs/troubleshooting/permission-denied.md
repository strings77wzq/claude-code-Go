# Permission Denied

Understanding and fixing permission errors.

## Default Safety Model

`go-code` starts in `WorkspaceWrite` by default. This mode allows workspace reads without prompting and asks before writes, edits, bash execution, network access, or other side effects.

When prompted:

- Choose `y` for a one-time approval.
- Choose `a` to remember the same operation for the current session.
- Choose `n` to deny it. The tool is not executed, and the agent receives an error result.

Permission decisions are written to the session trace as tool name, decision, summary, and timestamp. Secrets and full tool payloads are not stored in the permission audit entry.

## Common Scenarios

### Reading Sensitive Files

```
> Read .env

❌ Permission Denied: .env matches blocked pattern
```

**Solutions**:
1. One-time access: `/allow read .env`
2. Pattern access: `/allow read *.env`
3. Switch mode: `/mode ReadOnly` (can read anything)

### Writing to System Directories

```
> Write /etc/config "data"

❌ Permission Denied: /etc/* is blocked
```

**Solutions**:
1. Use user directory instead
2. Switch to DangerFullAccess (not recommended)
3. Use sudo manually outside claude-code-Go

### Executing Dangerous Commands

```
> Bash rm -rf /

⚠️  Permission Required: DangerFullAccess needed
```

**Solutions**:
1. Confirm you want to do this
2. Use `/allow` for one-time
3. Switch mode if automation needed

## Permission Levels

| Level | Can Read | Can Write | Can Execute |
|-------|----------|-----------|-------------|
| ReadOnly | ✅ All | ❌ None | ❌ None |
| WorkspaceWrite | ✅ All | ✅ Workspace | ✅ Safe |
| DangerFullAccess | ✅ All | ✅ All | ✅ All |

## Custom Rules

Add to `~/.go-code/settings.json`:

```json
{
  "rules": [
    {"pattern": "*.secret", "allowed": false},
    {"pattern": "docs/*", "allowed": true},
    {"pattern": "*.tmp", "allowed": true}
  ]
}
```

## Session Memory

Remember permissions:

```
> /remember allow read *.log
> /remember allow bash go test
```

## Debugging

Check current permissions:

```
> /mode
Current: WorkspaceWrite

> /rules
Active rules:
- *.env → DENY
- *.go → ALLOW
- * → ASK
```
