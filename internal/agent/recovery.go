// Package agent provides the core agent loop for the Claude Code clone.
package agent

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/strings77wzq/claude-code-Go/internal/api"
)

// Error types for recovery
const (
	ErrorTypeAPITimeout   = "api_timeout"
	ErrorTypeRateLimit    = "rate_limit"
	ErrorTypeToolError    = "tool_error"
	ErrorTypeContextFull  = "context_full"
	ErrorTypeServerError  = "server_error"
	ErrorTypeNetworkError = "network_error"
)

// RecoveryRecipe defines the recovery strategy for a specific error type
type RecoveryRecipe struct {
	ErrorType         string
	MaxRetries        int
	BackoffMultiplier time.Duration
	Description       string
	ShouldRetry       bool
	RecoveryAction    func(ctx context.Context, attempt int, err error) error
}

// RecoveryManager manages recovery recipes for different error types
type RecoveryManager struct {
	recipes map[string]*RecoveryRecipe
}

// NewRecoveryManager creates a new RecoveryManager with default recipes
func NewRecoveryManager() *RecoveryManager {
	rm := &RecoveryManager{
		recipes: make(map[string]*RecoveryRecipe),
	}

	// Register default recipes
	rm.registerDefaultRecipes()

	return rm
}

// registerDefaultRecipes registers the default error-specific recovery recipes
func (rm *RecoveryManager) registerDefaultRecipes() {
	// API Timeout: retry with exponential backoff (max 3 retries, base 1s)
	rm.RegisterRecipe(ErrorTypeAPITimeout, RecoveryRecipe{
		ErrorType:         ErrorTypeAPITimeout,
		MaxRetries:        3,
		BackoffMultiplier: time.Second,
		Description:       "API timeout - retry with exponential backoff",
		ShouldRetry:       true,
		RecoveryAction: func(ctx context.Context, attempt int, err error) error {
			log.Printf("[Recovery] Attempting to recover from timeout (attempt %d/%d)", attempt+1, 3)
			// Exponential backoff: 1s, 2s, 4s
			backoff := time.Second * (1 << attempt)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
			}
			return nil
		},
	})

	// Rate Limit: wait and retry (extract retry-after header, fallback to 60s)
	rm.RegisterRecipe(ErrorTypeRateLimit, RecoveryRecipe{
		ErrorType:         ErrorTypeRateLimit,
		MaxRetries:        3,
		BackoffMultiplier: 60 * time.Second,
		Description:       "Rate limited - wait and retry with retry-after header",
		ShouldRetry:       true,
		RecoveryAction: func(ctx context.Context, attempt int, err error) error {
			log.Printf("[Recovery] Rate limited, waiting before retry (attempt %d/%d)", attempt+1, 3)

			// Try to extract retry-after from error message
			retryAfter := 60 * time.Second
			var apiErr *api.APIError
			if errors.As(err, &apiErr) {
				// Check if there's a retry-after in the message
				if strings.Contains(apiErr.Message, "retry-after") {
					// Parse retry-after from message if present
					// The actual header extraction happens at the API level
					retryAfter = 60 * time.Second
				}
			}

			// Fallback to exponential backoff starting at 60s
			backoff := retryAfter * time.Duration(1<<attempt)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
			}
			return nil
		},
	})

	// Tool Error: report and continue (don't retry)
	rm.RegisterRecipe(ErrorTypeToolError, RecoveryRecipe{
		ErrorType:         ErrorTypeToolError,
		MaxRetries:        0,
		BackoffMultiplier: 0,
		Description:       "Tool error - report and continue without retry",
		ShouldRetry:       false,
		RecoveryAction: func(ctx context.Context, attempt int, err error) error {
			log.Printf("[Recovery] Tool error occurred: %v", err)
			// Don't retry tool errors - they should be reported back to the agent
			return nil
		},
	})

	// Context Full: compact and retry
	rm.RegisterRecipe(ErrorTypeContextFull, RecoveryRecipe{
		ErrorType:         ErrorTypeContextFull,
		MaxRetries:        2,
		BackoffMultiplier: 500 * time.Millisecond,
		Description:       "Context full - compact context and retry",
		ShouldRetry:       true,
		RecoveryAction: func(ctx context.Context, attempt int, err error) error {
			log.Printf("[Recovery] Context full, compacting context (attempt %d/%d)", attempt+1, 2)
			// The actual compaction will be handled by the caller
			// This is a placeholder for the recovery action
			return nil
		},
	})

	// Server Error: retry with exponential backoff
	rm.RegisterRecipe(ErrorTypeServerError, RecoveryRecipe{
		ErrorType:         ErrorTypeServerError,
		MaxRetries:        3,
		BackoffMultiplier: time.Second,
		Description:       "Server error - retry with exponential backoff",
		ShouldRetry:       true,
		RecoveryAction: func(ctx context.Context, attempt int, err error) error {
			log.Printf("[Recovery] Server error, retrying (attempt %d/%d)", attempt+1, 3)
			backoff := time.Second * (1 << attempt)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
			}
			return nil
		},
	})

	// Network Error: retry with exponential backoff
	rm.RegisterRecipe(ErrorTypeNetworkError, RecoveryRecipe{
		ErrorType:         ErrorTypeNetworkError,
		MaxRetries:        3,
		BackoffMultiplier: time.Second,
		Description:       "Network error - retry with exponential backoff",
		ShouldRetry:       true,
		RecoveryAction: func(ctx context.Context, attempt int, err error) error {
			log.Printf("[Recovery] Network error, retrying (attempt %d/%d)", attempt+1, 3)
			backoff := time.Second * (1 << attempt)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
			}
			return nil
		},
	})
}

// RegisterRecipe registers a recovery recipe for a specific error type
func (rm *RecoveryManager) RegisterRecipe(errorType string, recipe RecoveryRecipe) {
	rm.recipes[errorType] = &recipe
}

// GetRecipe returns the recovery recipe for a specific error type
func (rm *RecoveryManager) GetRecipe(errorType string) *RecoveryRecipe {
	return rm.recipes[errorType]
}

// classifyError classifies an error into a recovery error type
func classifyError(err error) string {
	if err == nil {
		return ""
	}

	// Check for API errors
	var apiErr *api.APIError
	if errors.As(err, &apiErr) {
		switch apiErr.Type {
		case api.ErrorTimeout:
			return ErrorTypeAPITimeout
		case api.ErrorRateLimit:
			return ErrorTypeRateLimit
		case api.ErrorServer:
			return ErrorTypeServerError
		case api.ErrorNetwork:
			return ErrorTypeNetworkError
		case api.ErrorAuth:
			// Auth errors should not be retried
			return ""
		}
	}

	// Check for context-related errors
	if strings.Contains(err.Error(), "context") && strings.Contains(err.Error(), "full") {
		return ErrorTypeContextFull
	}

	// Check for tool errors
	if strings.Contains(err.Error(), "tool") {
		return ErrorTypeToolError
	}

	// Check for HTTP status code in error
	if strings.Contains(err.Error(), "429") {
		return ErrorTypeRateLimit
	}
	if strings.Contains(err.Error(), "5") && len(err.Error()) > 0 {
		// Check for 5xx errors
		for i := 500; i <= 599; i++ {
			if strings.Contains(err.Error(), fmt.Sprintf("%d", i)) {
				return ErrorTypeServerError
			}
		}
	}

	return ""
}

// RecoveryContext holds context for recovery operations
type RecoveryContext struct {
	Manager    *RecoveryManager
	Agent      *Agent
	RetryCount int
	LastError  error
}

// ExecuteWithRecovery wraps a function with recovery logic
func (rc *RecoveryContext) ExecuteWithRecovery(ctx context.Context, fn func() error) error {
	err := fn()
	if err == nil {
		if rc.RetryCount > 0 {
			log.Printf("[Recovery] Successfully recovered after %d retry(s)", rc.RetryCount)
		}
		return nil
	}

	// Classify the error
	errorType := classifyError(err)
	if errorType == "" {
		// Unrecoverable error
		return err
	}

	recipe := rc.Manager.GetRecipe(errorType)
	if recipe == nil || !recipe.ShouldRetry {
		// No recipe or retry not recommended
		log.Printf("[Recovery] No recovery recipe for error type: %s", errorType)
		return err
	}

	// Check if we can still retry
	if rc.RetryCount >= recipe.MaxRetries {
		log.Printf("[Recovery] Max retries (%d) exceeded for error type: %s", recipe.MaxRetries, errorType)
		return err
	}

	// Execute recovery action
	rc.RetryCount++
	rc.LastError = err

	log.Printf("[Recovery] Executing recovery for error type: %s (attempt %d/%d)", errorType, rc.RetryCount, recipe.MaxRetries)

	if recipe.RecoveryAction != nil {
		if err := recipe.RecoveryAction(ctx, rc.RetryCount, err); err != nil {
			// Recovery action failed, return original error
			log.Printf("[Recovery] Recovery action failed: %v", err)
			return err
		}
	}

	// Special handling for context full - compact before retry
	if errorType == ErrorTypeContextFull && rc.Agent != nil {
		log.Printf("[Recovery] Compacting context before retry")
		rc.Agent.Compact()
	}

	// Retry the function
	return rc.ExecuteWithRecovery(ctx, fn)
}

// ExtractRetryAfter extracts retry-after duration from HTTP response
func ExtractRetryAfter(resp *http.Response) time.Duration {
	if resp == nil {
		return 0
	}

	// Check for Retry-After header
	retryAfter := resp.Header.Get("Retry-After")
	if retryAfter == "" {
		return 0
	}

	// Try to parse as seconds
	var seconds int
	if _, err := fmt.Sscanf(retryAfter, "%d", &seconds); err == nil {
		return time.Duration(seconds) * time.Second
	}

	// Try to parse as HTTP date (not implemented for simplicity)
	return 0
}
