---
title: Skills System
description: How to create and use custom skills
---

# Skills

Skills are custom commands that enhance the agent's capabilities. When you invoke a skill in the REPL using `/<skill-name>`, the skill's prompt is injected into the agent's system prompt, providing specialized instructions for the current task.

## What Are Skills

Skills are essentially named prompts stored as JSON files. They allow you to:

- Define reusable prompts for common tasks
- Customize the agent's behavior for specific workflows
- Create domain-specific instructions

## How to Create Custom Skills

Create a `.json` file in the `.go-code/skills/` directory with the following format:

```json
{
  "name": "skill-name",
  "description": "What this skill does",
  "prompt": "The instruction prompt to inject",
  "examples": ["/skill-name"]
}
```

### Fields

| Field | Required | Description |
|-------|----------|-------------|
| `name` | Yes | Unique identifier for the skill |
| `description` | Yes | Brief description shown in help |
| `prompt` | Yes | The prompt injected when skill is invoked |
| `examples` | No | Usage examples |

### Example

Create `.go-code/skills/refactor.json`:

```json
{
  "name": "refactor",
  "description": "Refactor code for better quality",
  "prompt": "Refactor the following code to improve readability, performance, and maintainability. Ensure all existing tests continue to pass.",
  "examples": ["/refactor"]
}
```

## Built-in Skills

The following skills are available by default:

| Skill | Description |
|-------|-------------|
| `review-pr` | Review a pull request |
| `explain-code` | Explain how code works |
| `write-tests` | Write tests for code |

## Usage

In the REPL, type `/<skill-name>` to invoke a skill:

```
> /review-pr
> /explain-code
> /write-tests
```

When a skill is invoked, its prompt is prepended to your message and sent to the agent, providing specialized context for the task.