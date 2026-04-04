package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

const (
	anthropicVersion = "2023-06-01"
	defaultBaseURL   = "https://api.anthropic.com"
	maxRetries       = 3
	retryDelayBase   = time.Second
)

type ErrorType string

const (
	ErrorAuth       ErrorType = "auth"
	ErrorRateLimit  ErrorType = "rate_limit"
	ErrorServer     ErrorType = "server"
	ErrorTimeout    ErrorType = "timeout"
	ErrorNetwork    ErrorType = "network"
	ErrorUnexpected ErrorType = "unexpected"
)

type APIError struct {
	Type    ErrorType
	Code    int
	Message string
}

func (e *APIError) Error() string {
	return e.Message
}

func classifyError(statusCode int, body string, originalErr error) *APIError {
	if originalErr != nil {
		if netErr, ok := originalErr.(interface{ Timeout() bool }); ok && netErr.Timeout() {
			return &APIError{
				Type:    ErrorTimeout,
				Message: "Request timed out. Please check your network connection and API key.",
			}
		}
		return &APIError{
			Type:    ErrorNetwork,
			Message: "Network error. Please check your internet connection.",
		}
	}

	switch statusCode {
	case http.StatusUnauthorized:
		return &APIError{
			Type:    ErrorAuth,
			Code:    401,
			Message: "Invalid API key. Please check your ANTHROPIC_API_KEY.",
		}
	case http.StatusForbidden:
		return &APIError{
			Type:    ErrorAuth,
			Code:    403,
			Message: "API access denied. Check your API key permissions.",
		}
	case http.StatusTooManyRequests:
		return &APIError{
			Type:    ErrorRateLimit,
			Code:    429,
			Message: "Rate limited. Retrying automatically...",
		}
	}

	if statusCode >= 500 {
		return &APIError{
			Type:    ErrorServer,
			Code:    statusCode,
			Message: "Server error. Please try again later.",
		}
	}

	return &APIError{
		Type:    ErrorUnexpected,
		Code:    statusCode,
		Message: fmt.Sprintf("Unexpected error (%d): %s", statusCode, body),
	}
}

type ConnectionStatus int

const (
	ConnConnecting ConnectionStatus = iota
	ConnConnected
	ConnTimeout
)

type Client struct {
	apiKey     string
	baseURL    string
	model      string
	httpClient *http.Client
	mu         sync.RWMutex
}

func NewClient(apiKey, baseURL, model string) *Client {
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	return &Client{
		apiKey:  apiKey,
		baseURL: baseURL,
		model:   model,
		httpClient: &http.Client{
			Timeout: 5 * time.Minute,
		},
	}
}

func (c *Client) SendMessage(ctx context.Context, req *ApiRequest) (*ApiResponse, error) {
	if req.Model == "" {
		req.Model = c.Model()
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/v1/messages", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.setHeaders(httpReq)

	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			delay := retryDelayBase * time.Duration(1<<(attempt-1))
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
			}
		}

		resp, err := c.httpClient.Do(httpReq)
		if err != nil {
			return nil, classifyError(0, "", err)
		}

		defer resp.Body.Close()

		if resp.StatusCode == http.StatusUnauthorized {
			return nil, classifyError(http.StatusUnauthorized, "", nil)
		}

		if resp.StatusCode == http.StatusForbidden {
			return nil, classifyError(http.StatusForbidden, "", nil)
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			lastErr = classifyError(http.StatusTooManyRequests, "", nil)
			continue
		}

		if resp.StatusCode >= 500 {
			lastErr = classifyError(resp.StatusCode, "", nil)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
			return nil, classifyError(resp.StatusCode, string(body), nil)
		}

		var apiResp ApiResponse
		if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}
		return &apiResp, nil
	}

	return nil, fmt.Errorf("request failed after %d retries: %w", maxRetries, lastErr)
}

func (c *Client) SendMessageStream(ctx context.Context, req *ApiRequest, onTextDelta func(text string)) (*ApiResponse, error) {
	req.Stream = true
	if req.Model == "" {
		req.Model = c.Model()
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/v1/messages", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.setHeaders(httpReq)

	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			delay := retryDelayBase * time.Duration(1<<(attempt-1))
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
			}
		}

		resp, err := c.httpClient.Do(httpReq)
		if err != nil {
			return nil, classifyError(0, "", err)
		}

		defer resp.Body.Close()

		if resp.StatusCode == http.StatusUnauthorized {
			return nil, classifyError(http.StatusUnauthorized, "", nil)
		}

		if resp.StatusCode == http.StatusForbidden {
			return nil, classifyError(http.StatusForbidden, "", nil)
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			lastErr = classifyError(http.StatusTooManyRequests, "", nil)
			continue
		}

		if resp.StatusCode >= 500 {
			lastErr = classifyError(resp.StatusCode, "", nil)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
			return nil, classifyError(resp.StatusCode, string(body), nil)
		}

		return parseStreamResponse(resp.Body, onTextDelta)
	}

	return nil, fmt.Errorf("request failed after %d retries: %w", maxRetries, lastErr)
}

func (c *Client) setHeaders(req *http.Request) {
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", anthropicVersion)
	req.Header.Set("content-type", "application/json")
}

func (c *Client) SetModel(model string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.model = model
}

func (c *Client) Model() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.model
}
