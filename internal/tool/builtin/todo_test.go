package builtin

import (
	"context"
	"sync"
	"testing"

	"github.com/strings77wzq/claude-code-Go/internal/tool"
)

func resetTodos() {
	globalTodoList.mu.Lock()
	defer globalTodoList.mu.Unlock()
	globalTodoList.items = []TodoItem{}
	globalTodoList.nextID = 1
}

func TestTodoToolName(t *testing.T) {
	td := NewTodoTool()
	if td.Name() != "TodoWrite" {
		t.Errorf("expected 'TodoWrite', got '%s'", td.Name())
	}
}

func TestTodoToolRequiresPermission(t *testing.T) {
	td := NewTodoTool()
	if td.RequiresPermission() {
		t.Errorf("TodoWrite should not require permission")
	}
}

func TestTodoToolInputSchema(t *testing.T) {
	td := NewTodoTool()
	schema := td.InputSchema()
	if schema["type"] != "object" {
		t.Errorf("expected type 'object'")
	}
}

func TestTodoToolExecuteHappyPath(t *testing.T) {
	resetTodos()
	td := NewTodoTool()
	result := td.Execute(context.Background(), map[string]any{
		"todos": []any{
			map[string]any{"content": "write tests", "status": "in_progress"},
			map[string]any{"content": "review code", "status": "pending"},
		},
	})

	if result.IsError {
		t.Fatalf("expected success, got error: %s", result.Content)
	}
}

func TestTodoToolExecuteEmptyTodos(t *testing.T) {
	resetTodos()
	td := NewTodoTool()
	result := td.Execute(context.Background(), map[string]any{
		"todos": []any{},
	})

	if result.IsError {
		t.Fatalf("expected success, got error: %s", result.Content)
	}
	if result.Content != "No todos" {
		t.Errorf("expected 'No todos', got '%s'", result.Content)
	}
}

func TestTodoToolExecuteVariousStatuses(t *testing.T) {
	resetTodos()
	td := NewTodoTool()
	result := td.Execute(context.Background(), map[string]any{
		"todos": []any{
			map[string]any{"content": "pending task", "status": "pending"},
			map[string]any{"content": "in progress task", "status": "in_progress"},
			map[string]any{"content": "done task", "status": "completed"},
		},
	})

	if result.IsError {
		t.Fatalf("expected success, got error: %s", result.Content)
	}
}

func TestTodoToolExecuteInvalidTodos(t *testing.T) {
	td := NewTodoTool()
	result := td.Execute(context.Background(), map[string]any{
		"todos": "not an array",
	})
	if !result.IsError {
		t.Errorf("expected error for invalid todos type")
	}
}

func TestTodoToolConcurrentSafety(t *testing.T) {
	resetTodos()
	td := NewTodoTool()

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			td.Execute(context.Background(), map[string]any{
				"todos": []any{
					map[string]any{"content": "task", "status": "pending"},
				},
			})
		}()
	}
	wg.Wait()

	globalTodoList.mu.Lock()
	count := len(globalTodoList.items)
	globalTodoList.mu.Unlock()

	if count != 10 {
		t.Errorf("expected 10 todos after concurrent writes, got %d", count)
	}
}

var _ tool.Tool = (*TodoTool)(nil)
