# Tutorial 5: Building Custom Tools

Extend claude-code-Go with your own tools.

## Why Custom Tools?

Built-in tools cover common tasks, but you might need:
- Company-specific commands
- Integration with internal systems
- Specialized operations

## Tool Structure

A tool consists of:
1. **Name**: Unique identifier
2. **Description**: What the tool does
3. **Parameters**: Input the tool accepts
4. **Handler**: Go function that executes the tool

## Simple Example

Let's create a tool that greets users:

```go
package main

import (
    "fmt"
    "github.com/strings77wzq/claude-code-Go/internal/tool"
)

func init() {
    tool.Register(tool.Definition{
        Name:        "Greet",
        Description: "Greet a user by name",
        Parameters: map[string]tool.Parameter{
            "name": {
                Type:        "string",
                Description: "Name of the person to greet",
                Required:    true,
            },
            "language": {
                Type:        "string",
                Description: "Language for greeting (en, es, fr)",
                Required:    false,
                Default:     "en",
            },
        },
        Handler: func(args map[string]interface{}) (string, error) {
            name := args["name"].(string)
            lang := args["language"].(string)
            
            greetings := map[string]string{
                "en": "Hello",
                "es": "Hola",
                "fr": "Bonjour",
            }
            
            greeting := greetings[lang]
            if greeting == "" {
                greeting = "Hello"
            }
            
            return fmt.Sprintf("%s, %s!", greeting, name), nil
        },
    })
}
```

## Using the Tool

Once registered, the AI can use it:

```
> Greet Alice in Spanish

🛠️  Using tool: Greet
   name: Alice
   language: es

💬 Response:
Hola, Alice!
```

## Real-World Example: Database Query

Let's create a tool that queries a database:

```go
package tools

import (
    "database/sql"
    "encoding/json"
    "fmt"
    
    "github.com/strings77wzq/claude-code-Go/internal/tool"
    _ "github.com/mattn/go-sqlite3"
)

type DBQueryTool struct {
    db *sql.DB
}

func NewDBQueryTool(dbPath string) (*DBQueryTool, error) {
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        return nil, err
    }
    return &DBQueryTool{db: db}, nil
}

func (t *DBQueryTool) Register() {
    tool.Register(tool.Definition{
        Name:        "DBQuery",
        Description: "Execute a read-only SQL query",
        Parameters: map[string]tool.Parameter{
            "query": {
                Type:        "string",
                Description: "SQL SELECT query to execute",
                Required:    true,
            },
            "limit": {
                Type:        "integer",
                Description: "Maximum rows to return",
                Required:    false,
                Default:     100,
            },
        },
        Handler: t.handleQuery,
    })
}

func (t *DBQueryTool) handleQuery(args map[string]interface{}) (string, error) {
    query := args["query"].(string)
    limit := args["limit"].(int)
    
    // Safety: Only allow SELECT queries
    if !strings.HasPrefix(strings.ToUpper(strings.TrimSpace(query)), "SELECT") {
        return "", fmt.Errorf("only SELECT queries are allowed")
    }
    
    // Add LIMIT if not present
    if !strings.Contains(strings.ToUpper(query), "LIMIT") {
        query = fmt.Sprintf("%s LIMIT %d", query, limit)
    }
    
    rows, err := t.db.Query(query)
    if err != nil {
        return "", err
    }
    defer rows.Close()
    
    // Convert to JSON
    results := []map[string]interface{}{}
    columns, _ := rows.Columns()
    
    for rows.Next() {
        values := make([]interface{}, len(columns))
        valuePtrs := make([]interface{}, len(columns))
        for i := range values {
            valuePtrs[i] = &values[i]
        }
        
        rows.Scan(valuePtrs...)
        
        row := map[string]interface{}{}
        for i, col := range columns {
            row[col] = values[i]
        }
        results = append(results, row)
    }
    
    jsonData, _ := json.MarshalIndent(results, "", "  ")
    return string(jsonData), nil
}
```

## Parameter Types

Supported parameter types:

| Type | Description | Example |
|------|-------------|---------|
| `string` | Text value | `"hello"` |
| `integer` | Whole number | `42` |
| `number` | Decimal number | `3.14` |
| `boolean` | true/false | `true` |
| `array` | List of values | `["a", "b"]` |
| `object` | Key-value pairs | `{"key": "value"}` |

## Advanced Features

### Validation

Add validation to parameters:

```go
Parameters: map[string]tool.Parameter{
    "email": {
        Type:        "string",
        Description: "Email address",
        Required:    true,
        Validate: func(value interface{}) error {
            email := value.(string)
            if !strings.Contains(email, "@") {
                return fmt.Errorf("invalid email format")
            }
            return nil
        },
    },
},
```

### Async Execution

For long-running operations:

```go
Handler: func(args map[string]interface{}) (string, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
    defer cancel()
    
    resultChan := make(chan string)
    go func() {
        // Long operation
        result := performLongTask()
        resultChan <- result
    }()
    
    select {
    case result := <-resultChan:
        return result, nil
    case <-ctx.Done():
        return "", fmt.Errorf("operation timed out")
    }
},
```

### Progress Updates

Report progress for long operations:

```go
Handler: func(args map[string]interface{}) (string, error) {
    progress := tool.GetProgressChannel()
    
    for i := 0; i < 10; i++ {
        progress <- fmt.Sprintf("Processing %d/10...", i+1)
        time.Sleep(time.Second)
    }
    
    return "Done!", nil
},
```

## Tool Categories

Organize tools by category:

```
tools/
├── database/       # Database tools
├── filesystem/     # File operations
├── network/        # HTTP, API calls
├── internal/       # Company-specific
└── utils/          # Utility functions
```

## Best Practices

### 1. Validate Inputs

Always validate and sanitize inputs:

```go
Handler: func(args map[string]interface{}) (string, error) {
    path := args["path"].(string)
    
    // Prevent directory traversal
    if strings.Contains(path, "..") {
        return "", fmt.Errorf("invalid path")
    }
    
    // Continue with operation
},
```

### 2. Handle Errors Gracefully

Provide helpful error messages:

```go
if err != nil {
    return "", fmt.Errorf("failed to connect to database: %w\n"+
        "Please check:\n"+
        "1. Database is running\n"+
        "2. Connection string is correct\n"+
        "3. Network connectivity", err)
}
```

### 3. Document Thoroughly

Write clear descriptions:

```go
Description: "Send a message to Slack channel. " +
    "Requires SLACK_TOKEN environment variable. " +
    "The bot must be invited to the channel first.",
```

### 4. Test Your Tools

Create unit tests:

```go
func TestGreetTool(t *testing.T) {
    result, err := greetHandler(map[string]interface{}{
        "name":     "Alice",
        "language": "es",
    })
    
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    
    if result != "Hola, Alice!" {
        t.Errorf("expected 'Hola, Alice!', got '%s'", result)
    }
}
```

## Example: GitHub Integration

```go
package tools

import (
    "context"
    "fmt"
    
    "github.com/google/go-github/v50/github"
    "github.com/strings77wzq/claude-code-Go/internal/tool"
    "golang.org/x/oauth2"
)

func init() {
    tool.Register(tool.Definition{
        Name:        "GitHubIssue",
        Description: "Create a GitHub issue",
        Parameters: map[string]tool.Parameter{
            "repo": {
                Type:        "string",
                Description: "Repository (owner/repo)",
                Required:    true,
            },
            "title": {
                Type:        "string",
                Description: "Issue title",
                Required:    true,
            },
            "body": {
                Type:        "string",
                Description: "Issue body (markdown)",
                Required:    true,
            },
            "labels": {
                Type:        "array",
                Description: "Labels to add",
                Required:    false,
            },
        },
        Handler: createGitHubIssue,
    })
}

func createGitHubIssue(args map[string]interface{}) (string, error) {
    token := os.Getenv("GITHUB_TOKEN")
    if token == "" {
        return "", fmt.Errorf("GITHUB_TOKEN not set")
    }
    
    ctx := context.Background()
    ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
    tc := oauth2.NewClient(ctx, ts)
    client := github.NewClient(tc)
    
    // Parse repo
    parts := strings.Split(args["repo"].(string), "/")
    if len(parts) != 2 {
        return "", fmt.Errorf("repo must be format 'owner/repo'")
    }
    
    issue := &github.IssueRequest{
        Title: github.String(args["title"].(string)),
        Body:  github.String(args["body"].(string)),
    }
    
    if labels, ok := args["labels"].([]string); ok {
        issue.Labels = labels
    }
    
    result, _, err := client.Issues.Create(ctx, parts[0], parts[1], issue)
    if err != nil {
        return "", err
    }
    
    return fmt.Sprintf("Created issue #%d: %s", result.GetNumber(), result.GetHTMLURL()), nil
}
```

## Next Steps

- [Tutorial 6: MCP Integration](06-mcp-integration.md)
- [API Reference: Tools](../../api/tools.md)
- [Examples Directory](../../../examples/)
