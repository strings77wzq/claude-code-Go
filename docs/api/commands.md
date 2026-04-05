---
title: REPL Commands Reference
description: Complete reference for all REPL commands in go-code
---

# REPL Commands Reference

go-code provides a set of slash commands for controlling the REPL session. These commands allow you to get help, manage sessions, switch models, and control the application behavior.

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
| [`/update`](#update) | Check and apply updates | `/update` |
| [`/exit`](#exit) | Exit the application | `/exit` |
| [`/skills`](#skills) | List available skills | `/skills` |

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
  /help       - Show this help message
  /clear      - Clear conversation history
  /model      - Show/switch current model
  /models     - List available models
  /sessions   - List saved sessions
  /resume     - Resume a session
  /compact    - Compress conversation context
  /update     - Check for updates
  /exit       - Exit the application
  /skills     - List available skills

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

**Anthropic:**
- `claude-sonnet-4-20250514` (default)
- `claude-opus-4-20250514`
- `claude-haiku-4-20250514`

**Tencent Coding Plan:**
- `tc-code-latest` (Auto)
- `hunyuan-2.0-instruct`
- `hunyuan-2.0-thinking`
- `minimax-m2.5`
- `kimi-k2.5`
- `glm-5`
- `hunyuan-t1`
- `hunyuan-turbos`

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

## /update

Checks for and applies updates.

### Usage

```
/update
```

### Description

Connects to the release server to check if a newer version is available. If an update is available, prompts to download and replace the current binary.

### Behavior

1. Checks latest version from GitHub releases
2. Compares with current version
3. If update available, prompts for confirmation
4. Downloads and replaces the binary
5. Requires restart to apply

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
Download and replace binary? [y/N]: y
Update successful. Please restart go-code.
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

Lists all available skills.

### Usage

```
/skills
```

### Description

Displays all custom skills that have been configured. Skills are custom commands that can be invoked with `/skillname`. Each skill shows its name and description.

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

- [Configuration Guide](../guide/configuration.md) — Model configuration
- [Session Management](../guide/session-management.md) — Session persistence
- [Skills System](../extension/skills.md) — Custom commands