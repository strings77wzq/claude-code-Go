package builtin

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/strings77wzq/claude-code-Go/internal/permission"
	"github.com/strings77wzq/claude-code-Go/internal/tool"
)

const (
	maxFetchSize = 50 * 1024
	fetchTimeout = 30 * time.Second
)

type WebFetchTool struct{}

func NewWebFetchTool() tool.Tool {
	return &WebFetchTool{}
}

func (w *WebFetchTool) Name() string {
	return "WebFetch"
}

func (w *WebFetchTool) Description() string {
	return "Fetches a URL and returns readable text content with HTML tags stripped."
}

func (w *WebFetchTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"url": map[string]any{
				"type":        "string",
				"description": "The URL to fetch",
			},
		},
		"required": []string{"url"},
	}
}

func (w *WebFetchTool) RequiresPermission() bool {
	return true
}

func (w *WebFetchTool) RequiredPermissionLevel() permission.PermissionLevel {
	return permission.LevelReadOnly
}

func (w *WebFetchTool) Execute(ctx context.Context, input map[string]any) tool.Result {
	url, ok := input["url"].(string)
	if !ok || url == "" {
		return tool.Error("url is required")
	}

	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return tool.Error(fmt.Sprintf("failed to create request: %v", err))
	}

	req.Header.Set("User-Agent", "go-code/1.0")

	client := &http.Client{
		Timeout: fetchTimeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return tool.Error(fmt.Sprintf("failed to fetch URL: %v", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return tool.Error(fmt.Sprintf("HTTP error: status %d", resp.StatusCode))
	}

	reader := io.LimitReader(resp.Body, maxFetchSize)
	body, err := io.ReadAll(reader)
	if err != nil {
		return tool.Error(fmt.Sprintf("failed to read response: %v", err))
	}

	content := stripHTML(string(body))

	content = strings.TrimSpace(content)

	if len(content) > maxFetchSize {
		content = content[:maxFetchSize] + "\n... (output truncated)"
	}

	if content == "" {
		return tool.Error("no readable content found")
	}

	return tool.Success(content)
}

func stripHTML(html string) string {
	var result strings.Builder
	inTag := false
	inScript := false
	inStyle := false

	for i := 0; i < len(html); i++ {
		c := html[i]

		if i+7 < len(html) && strings.ToLower(html[i:i+8]) == "<script" {
			inScript = true
		}
		if i+6 < len(html) && strings.ToLower(html[i:i+7]) == "</script" {
			inScript = false
			i += 8
			continue
		}
		if i+6 < len(html) && strings.ToLower(html[i:i+7]) == "<style" {
			inStyle = true
		}
		if i+7 < len(html) && strings.ToLower(html[i:i+8]) == "</style" {
			inStyle = false
			i += 7
			continue
		}

		if inScript || inStyle {
			continue
		}

		if c == '<' {
			inTag = true
			continue
		}
		if c == '>' {
			inTag = false
			continue
		}

		if !inTag {
			if c == '\n' || c == '\r' || c == '\t' || c == ' ' {
				if result.Len() > 0 && result.String()[result.Len()-1] != ' ' && result.String()[result.Len()-1] != '\n' {
					result.WriteByte(' ')
				}
			} else {
				result.WriteByte(c)
			}
		}
	}

	return strings.Join(strings.Fields(result.String()), "\n")
}
