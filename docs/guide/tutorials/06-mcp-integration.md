# Tutorial 6: MCP Integration

Learn to integrate the Model Context Protocol.

## What is MCP?

MCP (Model Context Protocol) is a standard for connecting AI tools to external services. It allows claude-code-Go to:
- Connect to databases
- Call external APIs
- Use cloud services
- Integrate with development tools

## MCP Architecture

```
claude-code-Go ←→ MCP Client ←→ MCP Server ←→ External Service
```

## Configuring MCP Servers

Add MCP servers to `~/.go-code/settings.json`:

```json
{
  "mcpServers": {
    "filesystem": {
      "command": "mcp-server-filesystem",
      "args": ["/home/user/projects"],
      "env": {}
    },
    "github": {
      "command": "mcp-server-github",
      "args": [],
      "env": {
        "GITHUB_TOKEN": "your-token"
      }
    },
    "postgres": {
      "command": "mcp-server-postgres",
      "args": ["postgresql://localhost/mydb"],
      "env": {}
    }
  }
}
```

## Built-in MCP Support

claude-code-Go supports MCP tools automatically. Once configured, MCP tools appear alongside built-in tools:

```
> List all MCP tools

Available MCP tools:
- filesystem_read_file
- filesystem_write_file
- github_create_issue
- postgres_query
```

## Using MCP Tools

MCP tools work like built-in tools:

```
> Use the filesystem MCP to read /home/user/projects/README.md

🛠️  Using MCP tool: filesystem_read_file
   path: /home/user/projects/README.md

📄 Content:
# My Project
...
```

## Creating an MCP Server

Here's a simple MCP server in Python:

```python
#!/usr/bin/env python3
import json
import sys

def main():
    while True:
        line = sys.stdin.readline()
        if not line:
            break
        
        request = json.loads(line)
        
        if request["method"] == "tools/list":
            response = {
                "jsonrpc": "2.0",
                "id": request["id"],
                "result": {
                    "tools": [
                        {
                            "name": "hello",
                            "description": "Say hello",
                            "inputSchema": {
                                "type": "object",
                                "properties": {
                                    "name": {
                                        "type": "string",
                                        "description": "Name to greet"
                                    }
                                },
                                "required": ["name"]
                            }
                        }
                    ]
                }
            }
        elif request["method"] == "tools/call":
            name = request["params"]["arguments"]["name"]
            response = {
                "jsonrpc": "2.0",
                "id": request["id"],
                "result": {
                    "content": [
                        {"type": "text", "text": f"Hello, {name}!"}
                    ]
                }
            }
        
        print(json.dumps(response), flush=True)

if __name__ == "__main__":
    main()
```

## MCP Protocol

MCP uses JSON-RPC 2.0:

### List Tools
```json
{"jsonrpc": "2.0", "id": 1, "method": "tools/list"}
```

### Call Tool
```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/call",
  "params": {
    "name": "hello",
    "arguments": {"name": "World"}
  }
}
```

## Best Practices

1. **Security**: Store credentials in environment variables
2. **Error Handling**: Always return meaningful error messages
3. **Timeouts**: Set reasonable timeouts for MCP calls
4. **Logging**: Enable MCP logging for debugging

## Next Steps

- [Tutorial 7: Session Management](07-session-management.md)
- [MCP Specification](https://modelcontextprotocol.io)
