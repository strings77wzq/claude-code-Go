---
title: Skills System
description: Technical deep dive into the Skills system — JSON format, loading mechanism, custom skill creation, REPL integration, and best practices
---

# Skills System

go-code implements a powerful Skills system that allows users to extend and customize the agent's behavior through named prompts. This document provides a comprehensive technical overview.

## What Are Skills?

**Skills** are named prompts that get injected into the agent's system prompt when invoked. They provide a mechanism for:

- Defining reusable prompts for common tasks
- Customizing agent behavior for specific workflows
- Creating domain-specific instructions
- Enhancing the agent with specialized knowledge

When you invoke a skill in the REPL using `/<skill-name>`, the skill's prompt is prepended to your message and sent to the agent, providing specialized context for the task.

## How Skills Work

### Loading Mechanism

Skills are loaded from the `.go-code/skills/` directory at startup:

```
.go-code/skills/
├── review-pr.json
├── explain-code.json
├── write-tests.json
└── custom-skill.json
```

The loader reads all `.json` files in this directory and parses them into `Skill` structures:

```go
// Skill represents a custom command skill
type Skill struct {
    Name        string   `json:"name" yaml:"name"`
    Description string   `json:"description" yaml:"description"`
    Prompt      string   `json:"prompt" yaml:"prompt"`
    Examples    []string `json:"examples" yaml:"examples"`
}
```

The loading process:
1. Reads all entries in the skills directory
2. Filters for `.json` files only
3. Parses each JSON file into a Skill struct
4. Validates that required fields are present
5. Registers valid skills in the registry

### JSON Format

Each skill is defined as a JSON file with the following structure:

```json
{
  "name": "skill-name",
  "description": "What this skill does",
  "prompt": "The instruction prompt to inject into the agent's context",
  "examples": ["/skill-name"]
}
```

| Field | Required | Description |
|-------|----------|-------------|
| `name` | Yes | Unique identifier for the skill (used in REPL commands) |
| `description` | Yes | Brief description shown when listing skills |
| `prompt` | Yes | The actual prompt content injected when skill is invoked |
| `examples` | No | Usage examples for documentation |

### REPL Integration

Skills integrate with the REPL through slash commands:

```
> /skills              # List all available skills
> /review-pr          # Invoke the review-pr skill
> /explain-code       # Invoke the explain-code skill
> /write-tests        # Invoke the write-tests skill
```

When a skill is invoked:
1. The skill's prompt is retrieved from the registry
2. The prompt is prepended to the user's message
3. The combined message is sent to the agent
4. The agent responds with specialized context from the skill

## Creating Custom Skills

### Step-by-Step Guide

1. **Create the skills directory** (if it doesn't exist):
   ```bash
   mkdir -p ~/.config/go-code/skills
   ```

2. **Create a JSON file** for your skill:
   ```bash
   touch ~/.config/go-code/skills/my-skill.json
   ```

3. **Define the skill** with the JSON format:
   ```json
   {
     "name": "my-skill",
     "description": "Description of what this skill does",
     "prompt": "Your custom prompt content here...",
     "examples": ["/my-skill"]
   }
   ```

4. **Restart go-code** to load the new skill

### Example: Code Review Skill

Create `~/.config/go-code/skills/review-pr.json`:

```json
{
  "name": "review-pr",
  "description": "Review a pull request for code quality and issues",
  "prompt": "You are performing a code review. Analyze the provided code changes carefully and provide constructive feedback on:\n\n1. Code quality and readability\n2. Potential bugs or edge cases\n3. Performance considerations\n4. Security vulnerabilities\n5. Test coverage\n\nBe specific with line numbers and suggest improvements where applicable.",
  "examples": ["/review-pr"]
}
```

### Example: Code Explanation Skill

Create `~/.config/go-code/skills/explain-code.json`:

```json
{
  "name": "explain-code",
  "description": "Explain how code works in detail",
  "prompt": "Explain the following code in detail. Cover:\n\n1. What the code does (overall purpose)\n2. How it works (step-by-step logic)\n3. Key functions and their roles\n4. Any interesting patterns or idioms used\n5. Potential improvements or alternatives\n\nUse clear language and provide examples where helpful.",
  "examples": ["/explain-code"]
}
```

### Example: Test Generation Skill

Create `~/.config/go-code/skills/write-tests.json`:

```json
{
  "name": "write-tests",
  "description": "Write comprehensive tests for the given code",
  "prompt": "Write comprehensive tests for the provided code. Cover:\n\n1. Unit tests for individual functions/methods\n2. Edge cases and error conditions\n3. Integration tests where applicable\n4. Use appropriate testing frameworks for the language\n5. Include clear test names and documentation\n\nEnsure tests are maintainable and follow best practices.",
  "examples": ["/write-tests"]
}
```

### Example: Refactoring Skill

Create `~/.config/go-code/skills/refactor.json`:

```json
{
  "name": "refactor",
  "description": "Refactor code for better quality",
  "prompt": "Refactor the following code to improve:\n\n1. Readability - clear variable names, good formatting\n2. Performance - optimize expensive operations\n3. Maintainability - clean structure, reduced complexity\n4. Testability - easier to unit test\n5. DRY principle - eliminate code duplication\n\nPreserve the original functionality and ensure all existing tests continue to pass.",
  "examples": ["/refactor"]
}
```

## Best Practices

### Writing Effective Skills

1. **Be Specific**: Define clear, focused prompts that target specific tasks
2. **Use Context**: Include relevant context about the domain or task type
3. **Provide Structure**: Use numbered lists or sections to organize expectations
4. **Set Expectations**: Clearly define what output format or quality is expected
5. **Keep Updates**: Version your skills if you make breaking changes

### Skill Organization

- **Group related skills** in the same directory
- **Use consistent naming** conventions (e.g., `verb-noun` pattern)
- **Document complex skills** with detailed descriptions
- **Test skills** with actual use cases

### Common Patterns

```json
{
  "name": "security-audit",
  "description": "Perform security audit on code",
  "prompt": "Conduct a thorough security audit focusing on:\n- Input validation\n- Authentication/authorization\n- Data protection\n- Common vulnerabilities (OWASP Top 10)\n\nProvide specific findings with severity levels.",
  "examples": ["/security-audit"]
}
```

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                        Skills Architecture                           │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│   ┌─────────────┐          ┌─────────────────┐                     │
│   │   REPL      │─────────▶│  Skill Registry │                     │
│   │  (/skills)  │          │                 │                     │
│   └─────────────┘          └────────┬────────┘                     │
│                                      │                               │
│                                      ▼                               │
│                              ┌─────────────────┐                     │
│                              │  Skill Loader   │                     │
│                              │  (JSON files)   │                     │
│                              └────────┬────────┘                     │
│                                       │                              │
│                                       ▼                              │
│   ┌─────────────┐          ┌─────────────────┐                     │
│   │   Agent     │◄─────────│  Skills Dir     │                     │
│   │ (injected)  │          │ ~/.config/...   │                     │
│   └─────────────┘          └─────────────────┘                     │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

## Related Documentation

- [Configuration Guide](../guide/configuration.md) — Configuration file locations
- [Tool System Overview](../tools/overview.md) — Tool interface and registry
- [Agent Loop Implementation](../core-code/agent-loop-impl.md) — Tool execution flow

---

<div class="nav-prev-next">

- [Extension Overview](./overview.md) ←
- → [Hooks System](./hooks.md)

</div>
