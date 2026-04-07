# Tutorial 2: Your First Tool Call

Learn how claude-code-Go uses tools to interact with your system.

## What are Tools?

Tools are functions the AI can call to:
- Read and write files
- Execute shell commands
- Search code
- Browse the web
- And more...

Think of tools as the AI's hands — they let it actually *do* things rather than just talk about them.

## Built-in Tools

claude-code-Go includes 10 built-in tools:

| Tool | Purpose | Example |
|------|---------|---------|
| `Read` | Read file contents | `Read main.go` |
| `Write` | Create new files | `Write config.json {"key": "value"}` |
| `Edit` | Modify existing files | `Edit main.go:10 change "foo" to "bar"` |
| `Glob` | Find files by pattern | `Glob *.go` |
| `Grep` | Search file contents | `Grep "func main"` |
| `Bash` | Execute shell commands | `Bash ls -la` |
| `Diff` | Show file differences | `Diff file.go` |
| `Tree` | Show directory structure | `Tree src/` |
| `WebFetch` | Fetch web content | `WebFetch https://example.com` |
| `TodoWrite` | Manage todo lists | `TodoWrite "Fix bug"` |

## How Tool Calls Work

When you ask the AI to do something, it follows this process:

```
1. You: "Read the README"
   ↓
2. AI thinks: "I need to use the Read tool"
   ↓
3. AI calls tool: Read {"file_path": "README.md"}
   ↓
4. Tool executes and returns content
   ↓
5. AI receives result and summarizes it
   ↓
6. You see the summary
```

## Hands-On Example

### Step 1: Create a test file

```bash
echo "Hello, World!" > test.txt
```

### Step 2: Ask the AI to read it

In the claude-code-Go REPL:

```
> Read the file test.txt
```

You'll see output like:

```
🛠️  Using tool: Read
   file_path: test.txt

📄 Content:
Hello, World!

💬 The file contains a simple greeting message.
```

### Step 3: Edit the file

```
> Edit test.txt to say "Hello, claude-code-Go!"
```

The AI will:
1. Read the current content
2. Determine what needs to change
3. Use the Edit tool
4. Confirm the change

```
🛠️  Using tool: Edit
   file_path: test.txt
   old_string: Hello, World!
   new_string: Hello, claude-code-Go!

✅ File updated successfully.
```

### Step 4: Verify the change

```
> Read test.txt
```

## Permission System

By default, claude-code-Go runs in "WorkspaceWrite" mode:
- ✅ Can read any file
- ✅ Can write to files in your workspace
- ⚠️  Will prompt before dangerous operations
- ❌ Cannot access sensitive files (`.env`, `.ssh`, etc.)

### Granting Permission

If the AI needs permission:

```
> Delete the test.txt file

⚠️  Permission required
The tool Bash with command "rm test.txt" requires DangerFullAccess mode.

Options:
1. Type "/allow" to grant one-time permission
2. Type "/mode DangerFullAccess" to switch modes
3. Cancel and try a safer approach
```

## Exercise: File Operations

Try these exercises:

1. **Create a file**: Ask the AI to create a new file
2. **Edit it**: Make changes to the file
3. **Search**: Use Grep to find content
4. **Directory listing**: Use Tree to see folder structure

## Understanding Tool Selection

The AI decides which tool to use based on your request:

| Your Request | AI Chooses |
|--------------|------------|
| "Show me main.go" | Read tool |
| "Find all Go files" | Glob tool |
| "Search for TODO comments" | Grep tool |
| "What files are in this directory?" | Tree or Bash tool |
| "Run the tests" | Bash tool |

The AI is smart about tool selection, but you can be explicit:

```
> Use the Glob tool to find all test files
```

## Next Steps

- [Tutorial 3: Understanding the Agent Loop](03-agent-loop.md) - Learn how the AI makes decisions
- [Tutorial 4: Permission System](04-permission-system.md) - Deep dive into security
- [API Reference](../../api/tools.md) - Complete tool documentation
