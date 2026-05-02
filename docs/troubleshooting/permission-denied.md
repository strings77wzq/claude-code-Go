---
title: Permission Denied
description: Understanding permission prompts, modes, and how to configure access control in go-code.
---

# Permission Denied

The permission system in `internal/permission` protects your system from unintended operations.

## Permission Modes

Three modes from `internal/permission/policy.go`:

| Mode | Read | Write/Edit | Bash |
| --- | --- | --- | --- |
| `ReadOnly` | Allowed | Denied | Denied |
| `WorkspaceWrite` | Allowed | Allowed (workspace) | Prompts |
| `DangerFullAccess` | Allowed | Allowed | Allowed |

Default is `WorkspaceWrite`. Switch with `/mode`.

## Permission Prompts

When a tool needs approval, you see a bordered prompt (from `internal/permission/prompter.go`):

```
┌──────────────────────────────────────────────┐
│  Tool: Bash                                    │
│  Command: rm -rf node_modules                  │
│  Allow? (y)es / (n)o / (a)lways               │
└──────────────────────────────────────────────┘
```

- `y` / `yes` — AllowOnce: single operation only.
- `n` / `no` — Deny: blocks the operation (tool is not executed).
- `a` / `always` — AllowForSession: remembers for the session.

## Bash Validation

Before the prompt, the enforcer classifies the command (from `internal/permission/bash_validation.go`):

**Blocked dangerous patterns**: `rm -rf /`, `sudo`, `curl | bash`, `mkfs`, `fdisk`, `dd if=`, `chmod 777 /`, `reboot`, `shutdown`, `killall`, and others.

**Read-only commands skip the prompt**: `ls`, `cat`, `grep`, `find`, `wc`, `head`, `tail`, `echo`, `pwd`, `tree`, `stat`, `which`, and others.

## File Boundary Checks

Write/Edit operations validate the path stays within the working directory (from `internal/permission/file_boundary.go`):

```go
func ResolveAndValidatePath(filePath, workingDir string) (string, error)
```

This resolves symlinks and returns `ErrPathEscape` if the path points outside.

## Evaluation Priority

From `internal/permission/policy.go`, the `Evaluate` method checks in order:

1. **Deny rules** (highest priority)
2. **Allow rules**
3. **Session memory** (previous `AllowForSession` decisions)
4. **Tool requirement** (minimum mode for specific tools)
5. **Active mode** (`DangerFullAccess` allows all)
6. **Default** — prompt the user (`Ask`)

## Custom Rules

Configure rules in `~/.go-code/settings.json` under the `permissions` key:

```json
{
  "permissions": {
    "denyRules": [{"pattern": "Bash(*.secret)", "allowed": false}],
    "allowRules": [
      {"pattern": "Read(docs/*)", "allowed": true},
      {"pattern": "Bash(go test *)", "allowed": true}
    ]
  }
}
```

## Changing Behavior

- `/mode ReadOnly` — block all writes and bash.
- `/mode DangerFullAccess` — disable all prompts (not recommended).
- `/rules` — see active allow/deny patterns.
- Permission decisions are logged in the session trace (tool name, decision, summary; no secrets).

## Related

- [Common Issues](common-issues) — General permission scenarios
- [PARITY.md](https://github.com/strings77wzq/claude-code-Go/blob/main/PARITY.md) — Verified workflow status
