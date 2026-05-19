package api

// Thinking configures extended thinking for API requests.
type Thinking struct {
	Type         string `json:"type"`                    // "enabled"
	BudgetTokens int    `json:"budget_tokens"`           // token budget for thinking
}

// ApiRequest represents a request to the Anthropic Messages API
type ApiRequest struct {
	Model     string           `json:"model"`
	MaxTokens int              `json:"max_tokens"`
	System    interface{}      `json:"system,omitempty"`
	Stream    bool             `json:"stream,omitempty"`
	Messages  []Message        `json:"messages"`
	Tools     []ToolDefinition `json:"tools,omitempty"`
	Thinking  *Thinking        `json:"thinking,omitempty"`
}

// CachedSystemPrompt wraps a system prompt with a cache_control breakpoint.
type CachedSystemPrompt struct {
	Type         string                `json:"type"`
	Text         string                `json:"text"`
	CacheControl *CacheControl         `json:"cache_control,omitempty"`
}

// CachedToolDefinition wraps a tool definition with a cache_control breakpoint.
type CachedToolDefinition struct {
	Name         string                `json:"name"`
	Description  string                `json:"description"`
	InputSchema  map[string]any        `json:"input_schema"`
	CacheControl *CacheControl         `json:"cache_control,omitempty"`
}

// CacheControl marks content for Anthropic prompt caching.
type CacheControl struct {
	Type string `json:"type"` // "ephemeral"
}

// Message represents a chat message
type Message struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"` // Can be string or []ContentBlock
}

// ToolDefinition defines a tool that can be called
type ToolDefinition struct {
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	InputSchema  map[string]any `json:"input_schema"`
	CacheControl *CacheControl  `json:"cache_control,omitempty"`
}

// ApiResponse represents a response from the Anthropic Messages API
type ApiResponse struct {
	ID         string         `json:"id"`
	Type       string         `json:"type"`
	Role       string         `json:"role"`
	Content    []ContentBlock `json:"content"`
	Model      string         `json:"model"`
	StopReason string         `json:"stop_reason,omitempty"`
	Usage      Usage          `json:"usage"`
}

// ContentBlock represents a content block in a message
type ContentBlock struct {
	Type      string         `json:"type"`
	Text      string         `json:"text,omitempty"`
	ID        string         `json:"id,omitempty"`
	Name      string         `json:"name,omitempty"`
	Input     map[string]any `json:"input,omitempty"`
	ToolUseID string         `json:"tool_use_id,omitempty"`
	IsError   bool           `json:"is_error,omitempty"`
}

// Usage represents token usage statistics
type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// ErrorResponse represents an error response from the API
type ErrorResponse struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

// ExtractSystemPrompt extracts the plain text from a system prompt that may be
// a plain string or a CachedSystemPrompt struct.
func ExtractSystemPrompt(system interface{}) string {
	switch s := system.(type) {
	case string:
		return s
	case CachedSystemPrompt:
		return s.Text
	default:
		return ""
	}
}
