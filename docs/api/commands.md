---
title: REPL Commands Reference
description: Complete reference for all REPL commands in go-code
---

# REPL Commands Reference

go-code provides a shared set of slash commands for controlling interactive sessions. The default TUI and legacy REPL use the same command layer for the commands listed here.

## Command List

| Command | Description | Example |
|---------|-------------|---------|
| [`/help`](#help) | Show help information | `/help` |
| [`/clear`](#clear) | Clear conversation history | `/clear` |
| [`/model`](#model) | Show or switch model | `/model claude-opus-4-20250514` |
| [`/models`](#models) | List available models | `/models` |
| [`/sessions`](#sessions) | List saved sessions | `/sessions` |
| [`/resume`](#resume) | Resume a session | `/resume session-id` |
| [`/compact`](#compact) | Compress conversation context | `/compact` |
| [`/permissions`](#permissions) | Show permission status | `/permissions` |
| [`/update`](#update) | Check for updates | `/update` |
| [`/exit`](#exit) | Exit the application | `/exit` |
| [`/skills`](#skills) | List available skills in legacy REPL | `/skills` |

---

## /help

Displays help information with all available commands.

### Usage

```
/help
```

### Description

Shows a summary of all available REPL commands with brief descriptions. This is the quickest way to learn about go-code's capabilities.

### Example Output

```
Available commands:
  /help        - Show this help
  /clear       - Clear conversation history
  /model       - Show current model
  /model <n>   - Switch model
  /models      - List available models
  /sessions    - List sessions
  /resume <id> - Resume session
  /compact     - Compact context
  /permissions - Show permission status
  /update      - Check for updates
  /exit        - Exit

Type /<command> to use a command.
```

---

## /clear

Clears the conversation history.

### Usage

```
/clear
```

### Description

Clears all conversation history from the current session. The agent will lose context of previous messages but the session remains active. This is useful for starting fresh without ending the session.

### Behavior

- Clears the agent's message history
- Does not affect saved sessions
- Preserves current model and settings

### Example

```
go-code> /clear
Conversation history cleared
```

---

## /model

Display or change the current model.

### Usage

```
/model              # Show current model
/model <model-name> # Switch to a different model
```

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `model-name` | string | No | The model to switch to |

### Description

Without arguments, displays the currently active model. With a model name, switches to that model for subsequent requests.

### Available Models

The supported model list is defined by the provider registry. Use `/models` for the current list.

### Examples

**Show current model:**
```
go-code> /model
Current model: claude-sonnet-4-20250514
```

**Switch to a different model:**
```
go-code> /model claude-opus-4-20250514
Model switched to: claude-opus-4-20250514
```

---

## /models

Lists all available models with descriptions.

### Usage

```
/models
```

### Description

Displays a comprehensive list of all models that can be used with go-code, organized by provider. Shows both Anthropic models and Tencent Coding Plan models.

### Example Output

```
Available models:

  Anthropic:
    claude-opus-4-6-20251001 - Most powerful model for complex reasoning
    claude-sonnet-4-6-20251001 - Balanced model for everyday tasks
    claude-haiku-4-20250514 - Fast and efficient model

  Openai:
    gpt-4o - OpenAI's most capable model
    gpt-4o-mini - Fast and affordable model

Switch model: /model <model-name>
```

---

## /sessions

Lists all saved sessions.

### Usage

```
/sessions
```

### Description

Displays all saved sessions from the sessions directory. Each session shows its ID, model, turn count, and start time. Sessions are automatically saved and can be resumed later.

### Example Output

```
Available sessions:
  abc123  model=claude-sonnet-4-20250514 turns=5 started=2026-04-05 10:30:00
  def456  model=claude-opus-4-20250514 turns=12 started=2026-04-04 15:45:30
  ghi789  model=claude-sonnet-4-20250514 turns=3 started=2026-04-03 09:20:15
```

### See Also

- [`/resume`](#resume) — Resume a specific session

---

## /resume

Resumes a previous session.

### Usage

```
/resume <session-id>
```

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `session-id` | string | Yes | The ID of the session to resume |

### Description

Loads a previous session's conversation history and continues from where it left off. The session must exist in the sessions directory.

### Examples

**Resume a session:**
```
go-code> /resume abc123
Resumed session abc123 with 10 messages
Session model: claude-sonnet-4-20250514
```

**Invalid session:**
```
go-code> /resume nonexistent
Session not found: nonexistent
```

---

## /compact

Compresses the conversation context.

### Usage

```
/compact
```

### Description

Triggers context compression for long conversations. This reduces the memory footprint of the session by summarizing older messages while preserving the key information. Useful for very long sessions to maintain performance.

### Behavior

- Summarizes older conversation messages
- Reduces token count for API calls
- Preserves important context

### Example

```
go-code> /compact
Conversation compacted
```

---

## /permissions

Shows permission status.

### Usage

```
/permissions
```

### Description

Displays the current permission status when available. The current implementation exposes a placeholder status while the safe-default approval flow is being completed.

### Example

```
go-code> /permissions
Permission mode details are not exposed yet. Safe-default approval flow is tracked in PARITY.md.
```

---

## /update

Checks for updates.

### Usage

```
/update
```

### Description

Connects to the release server to check if a newer version is available. The shared command layer reports the available version and download URL; automatic replacement is not performed from the shared TUI/REPL command.

### Behavior

1. Checks latest version from GitHub releases
2. Compares with current version
3. If update available, prints the download URL

### Examples

**No update available:**
```
go-code> /update
Already up to date (v0.1.0)
```

**Update available:**
```
go-code> /update
Update available: v0.1.0 -> v0.1.1
Download: https://github.com/strings77wzq/claude-code-Go/releases/...
```

---

## /exit

Exits the application.

### Usage

```
/exit
# or
/quit
```

### Description

Ends the current REPL session and exits the application. Sessions are automatically saved before exit.

### Example

```
go-code> /exit

Goodbye!
```

---

## /skills

Lists all available skills in the legacy REPL.

### Usage

```
/skills
```

### Description

Displays all custom skills that have been configured in the legacy REPL. Skills are custom commands that can be invoked with `/skillname`. Each skill shows its name and description.

### Example Output

```
Available skills:
  /brainstorming - Use before any creative work
  /debugging    - Systematic debugging workflow
  /refactor     - Intelligent refactoring workflow
  /tdd          - Test-driven development workflow
```

---

## Tips

### Command History

Use the up/down arrow keys to navigate through previous commands in the current session.

### Partial Input

For `/model`, you can type just the model name after the command:
```
go-code> /model claude-opus-4-20250514
```

### Tab Completion

Not supported in the basic REPL. Use the TUI mode for enhanced features.

---

## Related Documentation

- [Configuration Reference](./config.md) — Model configuration
- [Session Management](../guide/session-management.md) — Session persistence
- [Skills System](../extension/skills.md) — Custom commands
