package builtin

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/strings77wzq/claude-code-Go/internal/permission"
	"github.com/strings77wzq/claude-code-Go/internal/tool"
)

type TodoItem struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
	Status  string `json:"status"`
}

type TodoList struct {
	items  []TodoItem
	mu     sync.Mutex
	nextID int
}

var globalTodoList = &TodoList{
	items:  []TodoItem{},
	nextID: 1,
}

type TodoTool struct{}

func NewTodoTool() tool.Tool {
	return &TodoTool{}
}

func (t *TodoTool) Name() string {
	return "TodoWrite"
}

func (t *TodoTool) Description() string {
	return "Create, update, and manage todo items for tracking task progress"
}

func (t *TodoTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"todos": map[string]any{
				"type":        "array",
				"description": "Array of todo items. Each item can have: content (string), status (string: pending/in_progress/completed), id (optional int). If id is provided, update existing todo; if not, create new.",
				"items": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"id": map[string]any{
							"type":        "integer",
							"description": "Optional. If provided, update existing todo; if not, create new.",
						},
						"content": map[string]any{
							"type":        "string",
							"description": "The content of the todo item",
						},
						"status": map[string]any{
							"type":        "string",
							"description": "Status: pending, in_progress, or completed",
							"enum":        []string{"pending", "in_progress", "completed"},
						},
					},
				},
			},
		},
		"required": []string{"todos"},
	}
}

func (t *TodoTool) RequiresPermission() bool {
	return false
}

func (t *TodoTool) RequiredPermissionLevel() permission.PermissionLevel {
	return permission.LevelWorkspaceWrite
}

func (t *TodoTool) Execute(ctx context.Context, input map[string]any) tool.Result {
	todosRaw, ok := input["todos"].([]any)
	if !ok {
		return tool.Error("todos must be an array")
	}

	globalTodoList.mu.Lock()
	defer globalTodoList.mu.Unlock()

	for _, todoRaw := range todosRaw {
		todoMap, ok := todoRaw.(map[string]any)
		if !ok {
			return tool.Error("each todo must be an object")
		}

		var id int
		if idVal, hasID := todoMap["id"]; hasID {
			idFloat, ok := idVal.(float64)
			if !ok {
				return tool.Error("id must be an integer")
			}
			id = int(idFloat)
		}

		content, _ := todoMap["content"].(string)

		status := "pending"
		if statusVal, hasStatus := todoMap["status"]; hasStatus {
			status, ok = statusVal.(string)
			if !ok {
				return tool.Error("status must be a string")
			}
			if status != "pending" && status != "in_progress" && status != "completed" {
				return tool.Error("status must be pending, in_progress, or completed")
			}
		}

		if id > 0 {
			found := false
			for i, item := range globalTodoList.items {
				if item.ID == id {
					if content != "" {
						globalTodoList.items[i].Content = content
					}
					if status != "" {
						globalTodoList.items[i].Status = status
					}
					found = true
					break
				}
			}
			if !found {
				return tool.Error(fmt.Sprintf("todo with id %d not found", id))
			}
		} else {
			newTodo := TodoItem{
				ID:      globalTodoList.nextID,
				Content: content,
				Status:  status,
			}
			globalTodoList.items = append(globalTodoList.items, newTodo)
			globalTodoList.nextID++
		}
	}

	return tool.Success(formatTodoList(globalTodoList.items))
}

func formatTodoList(items []TodoItem) string {
	if len(items) == 0 {
		return "No todos"
	}

	var sb strings.Builder
	sb.WriteString("Todo List:\n")

	for _, item := range items {
		var statusIcon string
		switch item.Status {
		case "in_progress":
			statusIcon = "🔄"
		case "completed":
			statusIcon = "✅"
		default:
			statusIcon = "📋"
		}
		fmt.Fprintf(&sb, "%s [%d] %s - %s\n", statusIcon, item.ID, item.Status, item.Content)
	}

	return sb.String()
}
