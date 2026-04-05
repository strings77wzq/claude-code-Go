---
title: Built-in Tools API Reference
description: Complete API reference for all 10 built-in tools in go-code
---

# Built-in Tools API Reference

go-code provides 10 built-in tools that enable the agent to interact with the filesystem, execute shell commands, and manage task tracking. This reference documents each tool's parameters, permissions, and usage examples.

## Tool List

| # | Tool Name | Description | Permission Required |
|---|-----------|-------------|---------------------|
| 1 | [Read](#read) | Read file contents with optional offset/limit | No (ReadOnly) |
| 2 | [Write](#write) | Create or overwrite files | Yes (WorkspaceWrite) |
| 3 | [Edit](#edit) | Make targeted code edits using exact string matching | Yes (WorkspaceWrite) |
| 4 | [Glob](#glob) | Find files by glob patterns | No (ReadOnly) |
| 5 | [Grep](#grep) | Search file contents using regular expressions | No (ReadOnly) |
| 6 | [Bash](#bash) | Execute shell commands | Yes (WorkspaceWrite) |
| 7 | [Diff](#diff) | Compare two content strings | No (ReadOnly) |
| 8 | [Tree](#tree) | Display directory tree structure | No (ReadOnly) |
| 9 | [WebFetch](#webfetch) | Fetch URL and return readable text | Yes (WorkspaceWrite) |
| 10 | [TodoWrite](#todowrite) | Create and manage todo items | No (WorkspaceWrite) |

---

## Read

Reads a file and returns its contents with line numbers.

### Tool Definition

```json
{
  "name": "Read",
  "description": "Reads a file and returns its contents with line numbers.",
  "input_schema": {
    "type": "object",
    "properties": {
      "file_path": {
        "type": "string",
        "description": "Path to the file to read"
      },
      "offset": {
        "type": "number",
        "description": "Line number to start reading from (0-based, default: 0)"
      },
      "limit": {
        "type": "number",
        "description": "Maximum number of lines to read (default: 2000)"
      }
    },
    "required": ["file_path"]
  }
}
```

### Parameters

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `file_path` | string | Yes | - | Path to the file to read |
| `offset` | number | No | 0 | Line number to start reading from (0-based) |
| `limit` | number | No | 2000 | Maximum number of lines to read |

### Permission Level

- **Requires Permission**: No
- **Permission Level**: `ReadOnly`

### Constraints

- Maximum file size: 200KB
- Directories cannot be read (will return an error)

### Examples

**Read entire file:**
```
Input: { "file_path": "/home/user/project/main.go" }
Output:
1: package main
2:
3: func main() {
4:     fmt.Println("Hello, World!")
5: }
```

**Read specific lines:**
```
Input: { "file_path": "/home/user/project/main.go", "offset": 10, "limit": 20 }
```

---

## Write

Creates a new file or overwrites an existing file.

### Tool Definition

```json
{
  "name": "Write",
  "description": "Creates a new file or overwrites an existing file.",
  "input_schema": {
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

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `file_path` | string | Yes | Path to the file to write |
| `content` | string | Yes | Content to write to the file |

### Permission Level

- **Requires Permission**: Yes
- **Permission Level**: `WorkspaceWrite`

### Behavior

- Creates parent directories if they don't exist
- Overwrites existing files without confirmation
- Returns success message with file path

### Examples

**Create a new file:**
```
Input: {
  "file_path": "/home/user/project/README.md",
  "content": "# My Project\n\nThis is my project."
}
Output: Successfully wrote to /home/user/project/README.md
```

**Create file with nested directories:**
```
Input: {
  "file_path": "/home/user/project/src/lib/utils.go",
  "content": "package utils\n\nfunc Hello() string { return \"Hello\" }"
}
Output: Successfully wrote to /home/user/project/src/lib/utils.go
```

---

## Edit

Performs exact string replacement in a file.

### Tool Definition

```json
{
  "name": "Edit",
  "description": "Performs exact string replacement in a file.",
  "input_schema": {
    "type": "object",
    "properties": {
      "file_path": {
        "type": "string",
        "description": "Path to the file to edit"
      },
      "old_string": {
        "type": "string",
        "description": "The exact string to replace"
      },
      "new_string": {
        "type": "string",
        "description": "The replacement string"
      },
      "replace_all": {
        "type": "boolean",
        "description": "Replace all occurrences (default: false)"
      }
    },
    "required": ["file_path", "old_string", "new_string"]
  }
}
```

### Parameters

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `file_path` | string | Yes | - | Path to the file to edit |
| `old_string` | string | Yes | - | The exact string to replace |
| `new_string` | string | Yes | - | The replacement string |
| `replace_all` | boolean | No | false | Replace all occurrences |

### Permission Level

- **Requires Permission**: Yes
- **Permission Level**: `WorkspaceWrite`

### Constraints

- Requires exact string match (whitespace-sensitive)
- If `old_string` appears multiple times, defaults to error unless `replace_all=true`
- Maximum file size: 200KB

### Examples

**Simple replacement:**
```
Input: {
  "file_path": "/home/user/project/main.go",
  "old_string": "fmt.Println(\"Hello\")",
  "new_string": "fmt.Println(\"Hello, World!\")"
}
Output: Successfully edited /home/user/project/main.go
```

**Replace all occurrences:**
```
Input: {
  "file_path": "/home/user/project/main.go",
  "old_string": "TODO",
  "new_string": "DONE",
  "replace_all": true
}
Output: Successfully edited /home/user/project/main.go
```

---

## Glob

Finds files matching glob patterns.

### Tool Definition

```json
{
  "name": "Glob",
  "description": "Finds files matching glob patterns.",
  "input_schema": {
    "type": "object",
    "properties": {
      "path": {
        "type": "string",
        "description": "Glob pattern to match files (e.g., *.go, **/*.ts)"
      }
    },
    "required": ["path"]
  }
}
```

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `path` | string | Yes | Glob pattern to match files |

### Permission Level

- **Requires Permission**: No
- **Permission Level**: `ReadOnly`

### Pattern Syntax

| Pattern | Description | Example |
|---------|-------------|---------|
| `*` | Match any characters in a single directory | `*.go` matches all `.go` files |
| `**` | Match across directories recursively | `**/*.ts` matches all `.ts` files |
| `?` | Match single character | `file?.txt` matches `file1.txt` |
| `[abc]` | Match character class | `[abc].txt` matches `a.txt`, `b.txt` |

### Examples

**Find all Go files:**
```
Input: { "path": "**/*.go" }
Output:
/home/user/project/main.go
/home/user/project/utils/helper.go
/home/user/project/cmd/app.go
```

**Find specific extension:**
```
Input: { "path": "*.md" }
Output:
README.md
CHANGELOG.md
```

---

## Grep

Searches file contents using regular expressions.

### Tool Definition

```json
{
  "name": "Grep",
  "description": "Searches file contents using regular expressions.",
  "input_schema": {
    "type": "object",
    "properties": {
      "pattern": {
        "type": "string",
        "description": "Regular expression pattern to search for"
      },
      "path": {
        "type": "string",
        "description": "File or directory path to search in (glob patterns supported)"
      },
      "include": {
        "type": "string",
        "description": "File pattern to include (e.g., *.go, *.ts)"
      },
      "output_mode": {
        "type": "string",
        "description": "Output format: files_with_matches, content, or count"
      }
    },
    "required": ["pattern", "path"]
  }
}
```

### Parameters

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `pattern` | string | Yes | - | Regular expression pattern |
| `path` | string | Yes | - | File or directory to search |
| `include` | string | No | - | File pattern to include |
| `output_mode` | string | No | content | Output format |

### Permission Level

- **Requires Permission**: No
- **Permission Level**: `ReadOnly`

### Output Modes

| Mode | Description |
|------|-------------|
| `content` | Shows matching lines with context (default) |
| `files_with_matches` | Shows only file paths containing matches |
| `count` | Shows match count per file |

### Examples

**Search for function definition:**
```
Input: { "pattern": "func.*main", "path": "/home/user/project", "include": "*.go" }
Output:
main.go:5: func main() {
main.go:10: func mainLogic() {
```

**Find all TODO comments:**
```
Input: { "pattern": "TODO", "path": "/home/user/project", "include": "*.go", "output_mode": "files_with_matches" }
Output:
/home/user/project/main.go
/home/user/project/utils/helper.go
```

---

## Bash

Executes shell commands.

### Tool Definition

```json
{
  "name": "Bash",
  "description": "Executes shell commands in the project directory.",
  "input_schema": {
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

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `command` | string | Yes | Shell command to execute |

### Permission Level

- **Requires Permission**: Yes
- **Permission Level**: `WorkspaceWrite`

### Behavior

- Executes in the project root directory
- Default timeout: 120 seconds
- Output truncated at 100KB
- Working directory: project root

### Examples

**Run a test:**
```
Input: { "command": "go test ./..." }
Output: PASS
ok  	github.com/user/project	0.523s
```

**List files:**
```
Input: { "command": "ls -la" }
Output: total 48
drwxr-xr-x  4 user  staff   128 Apr  5 10:00 .
drwxr-xr-x  2 user  staff   256 Apr  5 10:01 src
```

**Git operations:**
```
Input: { "command": "git status" }
Output: On branch main
Your branch is up to date with 'origin/main'.
```

---

## Diff

Compares two content strings and returns unified diff output.

### Tool Definition

```json
{
  "name": "Diff",
  "description": "Compares two content strings and returns unified diff output.",
  "input_schema": {
    "type": "object",
    "properties": {
      "old_string": {
        "type": "string",
        "description": "Original content"
      },
      "new_string": {
        "type": "string",
        "description": "Modified content"
      }
    },
    "required": ["old_string", "new_string"]
  }
}
```

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `old_string` | string | Yes | Original content |
| `new_string` | string | Yes | Modified content |

### Permission Level

- **Requires Permission**: No
- **Permission Level**: `ReadOnly`

### Behavior

- Uses system `diff` command if available
- Falls back to pure Go implementation
- Returns unified diff format

### Examples

**Compare two strings:**
```
Input: {
  "old_string": "func Hello() string {\n    return \"Hello\"\n}",
  "new_string": "func Hello() string {\n    return \"Hello, World!\"\n}"
}
Output:
--- 
+++ 
@@ -1,3 +1,3 @@
 func Hello() string {
-    return "Hello"
+    return "Hello, World!"
 }
```

---

## Tree

Displays directory tree structure as text.

### Tool Definition

```json
{
  "name": "Tree",
  "description": "Displays directory tree structure as text.",
  "input_schema": {
    "type": "object",
    "properties": {
      "path": {
        "type": "string",
        "description": "Root directory path (default: current directory)"
      },
      "limit": {
        "type": "number",
        "description": "Maximum directory depth to display (default: 3)"
      }
    },
    "required": []
  }
}
```

### Parameters

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `path` | string | No | current dir | Root directory path |
| `limit` | number | No | 3 | Maximum depth to display |

### Permission Level

- **Requires Permission**: No
- **Permission Level**: `ReadOnly`

### Examples

**Display project structure:**
```
Input: { "path": "/home/user/project", "limit": 2 }
Output:
project/
├── cmd/
│   └── app/
├── internal/
│   └── handler/
├── pkg/
│   └── utils/
└── go.mod
```

---

## WebFetch

Fetches a URL and returns readable text (HTML stripped).

### Tool Definition

```json
{
  "name": "WebFetch",
  "description": "Fetches a URL and returns readable text with HTML tags stripped.",
  "input_schema": {
    "type": "object",
    "properties": {
      "url": {
        "type": "string",
        "description": "URL to fetch"
      },
      "goal": {
        "type": "string",
        "description": "Specific information to extract from the page"
      }
    },
    "required": ["url"]
  }
}
```

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `url` | string | Yes | URL to fetch |
| `goal` | string | No | Specific information to extract |

### Permission Level

- **Requires Permission**: Yes
- **Permission Level**: `WorkspaceWrite`

### Constraints

- Output limited to 50KB
- HTML tags automatically stripped

### Examples

**Fetch a page:**
```
Input: { "url": "https://example.com" }
Output: Example Domain
This domain is for use in illustrative examples in documents...
```

**Extract specific info:**
```
Input: { "url": "https://api.github.com/repos/user/repo", "goal": "description" }
Output: A cool project that does stuff
```

---

## TodoWrite

Creates, updates, and manages todo items for tracking task progress.

### Tool Definition

```json
{
  "name": "TodoWrite",
  "description": "Create, update, and manage todo items for tracking task progress",
  "input_schema": {
    "type": "object",
    "properties": {
      "todos": {
        "type": "array",
        "description": "Array of todo items",
        "items": {
          "type": "object",
          "properties": {
            "id": {
              "type": "integer",
              "description": "Optional. If provided, update existing todo; if not, create new."
            },
            "content": {
              "type": "string",
              "description": "The content of the todo item"
            },
            "status": {
              "type": "string",
              "description": "Status: pending, in_progress, or completed",
              "enum": ["pending", "in_progress", "completed"]
            }
          }
        }
      }
    },
    "required": ["todos"]
  }
}
```

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `todos` | array | Yes | Array of todo items |

### Permission Level

- **Requires Permission**: No
- **Permission Level**: `WorkspaceWrite`

### Todo Item Properties

| Property | Type | Required | Description |
|----------|------|----------|-------------|
| `id` | integer | No | If provided, updates existing; if not, creates new |
| `content` | string | Yes | The todo content |
| `status` | string | No | One of: pending, in_progress, completed |

### Examples

**Create todos:**
```
Input: {
  "todos": [
    { "content": "Fix the login bug", "status": "in_progress" },
    { "content": "Add unit tests", "status": "pending" }
  ]
}
Output:
Todo List:
🔄 [1] in_progress - Fix the login bug
📋 [2] pending - Add unit tests
```

**Update todo status:**
```
Input: {
  "todos": [
    { "id": 1, "status": "completed" }
  ]
}
Output:
Todo List:
✅ [1] completed - Fix the login bug
📋 [2] pending - Add unit tests
```

---

## Related Documentation

- [Tool System Overview](../tools/overview.md) — Architecture and extension guide
- [Configuration Guide](../guide/configuration.md) — Tool configuration
- [MCP Integration](../extension/mcp.md) — External tools via MCP