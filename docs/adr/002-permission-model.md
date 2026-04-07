# ADR-002: Permission Model

## Status
Accepted

## Context
AI coding assistants can potentially:
- Read sensitive files (.env, .ssh, etc.)
- Execute arbitrary commands
- Delete or modify critical files
- Access private data

We need a permission system that:
1. Protects users from accidental damage
2. Allows powerful operations when needed
3. Isn't annoying with constant prompts
4. Supports automation in CI/CD

## Decision
We implemented a **3-Tier Permission Model**:

### Tier 1: ReadOnly
- Can only read files
- Cannot write or execute
- Safe for exploring codebases

### Tier 2: WorkspaceWrite (Default)
- Can read files
- Can write within workspace
- Can execute safe commands (ls, cat, grep)
- Prompts for dangerous operations

### Tier 3: DangerFullAccess
- Can do anything
- Minimal prompting
- For trusted automation

### Glob Rules

Each tier supports glob rules for fine-grained control:

```json
{
  "mode": "WorkspaceWrite",
  "rules": [
    {"pattern": "*.go", "allowed": true},
    {"pattern": "*.env", "allowed": false},
    {"pattern": "/etc/*", "allowed": false}
  ]
}
```

### Session Memory

Permissions can be granted for the session:
- `/allow read secret.txt` - One-time grant
- `/allow write *.json` - Pattern-based grant
- `/remember allow bash` - Persist across sessions

### Consequences

**Positive:**
- Users feel safe
- Supports both interactive and automated use
- Flexible rule system
- Clear audit trail

**Negative:**
- Complex rule evaluation
- Need to balance security vs usability
- Rules can be complex to configure

## Alternatives Considered

1. **Binary Allow/Deny**: Too coarse
2. **Per-Operation Prompts**: Too annoying
3. **Sandbox/Container**: Too heavy for local use

## Related Decisions

- ADR-001: Agent Loop (permission checks happen during Execute state)
