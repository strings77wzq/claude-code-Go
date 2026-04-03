# Built-in Tools

go-code includes six built-in tools that provide the core file system and command execution capabilities.

## Tool Reference

| Tool | Description | Risk Level |
|------|-------------|------------|
| `Read` | Read file contents | Low |
| `Write` | Create or overwrite files | Medium |
| `Edit` | Make targeted code edits | Medium |
| `Glob` | Find files by pattern | Low |
| `Grep` | Search file contents | Low |
| `Bash` | Execute shell commands | High |

## Read

Reads the contents of a file from the file system.

**Schema:**
```json
{
  "name": "Read",
  "description": "Read the contents of a file from the filesystem. Use this when you need to examine the contents of a file.",
  "parameters": {
    "type": "object",
    "properties": {
      "file_path": {
        "type": "string",
        "description": "Path to the file to read"
      }
    },
    "required": ["file_path"]
  }
}
```

**Example:**
```
Tool: Read
Arguments: {"file_path": "src/main.go"}
```

## Write

Creates a new file or overwrites an existing file.

**Schema:**
```json
{
  "name": "Write",
  "description": "Create a new file or overwrite an existing file with new content.",
  "parameters": {
    "type": "object",
    "properties": {
      "file_path": {
        "type": "string",
        "description": "Path to the file to write"
      },
      "content": {
        "type": "string",
        "description": "Content to write to the file"
      }
    },
    "required": ["file_path", "content"]
  }
}
```

**Note:** Requires user approval unless auto-approve is enabled.

## Edit

Makes targeted modifications to an existing file. Uses a line-based approach.

**Schema:**
```json
{
  "name": "Edit",
  "description": "Make a targeted edit to an existing file. Replaces specific lines with new content.",
  "parameters": {
    "type": "object",
    "properties": {
      "file_path": {
        "type": "string",
        "description": "Path to the file to edit"
      },
      "old_string": {
        "type": "string",
        "description": "The exact text to replace (must match file content)"
      },
      "new_string": {
        "type": "string",
        "description": "The replacement text"
      }
    },
    "required": ["file_path", "old_string", "new_string"]
  }
}
```

**Note:** Requires user approval unless auto-approve is enabled.

## Glob

Finds files matching a glob pattern.

**Schema:**
```json
{
  "name": "Glob",
  "description": "Find files matching a glob pattern.",
  "parameters": {
    "type": "object",
    "properties": {
      "pattern": {
        "type": "string",
        "description": "Glob pattern (e.g., *.go, **/*.ts)"
      }
    },
    "required": ["pattern"]
  }
}
```

**Example patterns:**
- `*.go` - All Go files in current directory
- `**/*.js` - All JavaScript files recursively
- `src/**/*.ts` - TypeScript files in src/

## Grep

Searches for text within files.

**Schema:**
```json
{
  "name": "Grep",
  "description": "Search for text patterns in files.",
  "parameters": {
    "type": "object",
    "properties": {
      "pattern": {
        "type": "string",
        "description": "Regular expression to search for"
      },
      "path": {
        "type": "string",
        "description": "Path to search in (file or directory)"
      }
    },
    "required": ["pattern", "path"]
  }
}
```

## Bash

Executes shell commands.

**Schema:**
```json
{
  "name": "Bash",
  "description": "Execute a shell command.",
  "parameters": {
    "type": "object",
    "properties": {
      "command": {
        "type": "string",
        "description": "Shell command to execute"
      }
    },
    "required": ["command"]
  }
}
```

**Note:** High risk - always requires user approval.

## MCP Integration

In addition to built-in tools, go-code supports the Model Context Protocol (MCP) for extending capabilities.

### Configuring MCP Servers

In `config.yaml`:

```yaml
mcp_servers:
  filesystem:
    command: "npx"
    args: ["-y", "@modelcontextprotocol/server-filesystem", "/path"]
  github:
    command: "python"
    args: ["-m", "mcp.server.github", "--token", "your-token"]
```

### MCP Tool Discovery

When MCP servers are configured, the agent:
1. Connects to each MCP server on startup
2. Requests the list of available tools
3. Registers discovered tools in the tool registry
4. Routes MCP tool calls through the MCP client

## Tool Execution Flow

```
Model Response
      │
      ▼
┌─────────────────┐
│ Parse Tool Call │
└────────┬────────┘
         │
         ▼
┌─────────────────┐     ┌─────────────────┐
│ Check Registry  │────▶│ Permission      │
│ (tool exists?)  │     │ System          │
└────────┬────────┘     └────────┬────────┘
         │ YES                     │
         ▼                         ▼
┌─────────────────┐     ┌─────────────────┐
│ Execute Tool    │◀────│ User Approval  │
│                 │     │ (if required)   │
└────────┬────────┘     └─────────────────┘
         │
         ▼
┌─────────────────┐
│ Return Result   │
│ to Model        │
└─────────────────┘
```