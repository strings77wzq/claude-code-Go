// Package memory provides cross-session knowledge persistence for agents.
package memory

import "context"

// Store provides key-value storage and semantic search for agent memory.
type Store interface {
	// KV operations (no embedding needed)
	Put(ctx context.Context, key string, value []byte) error
	Get(ctx context.Context, key string) ([]byte, error)
	Delete(ctx context.Context, key string) error
	List(ctx context.Context, prefix string) ([]string, error)

	// Semantic search (requires embedding)
	Index(ctx context.Context, key string, content string) error
	Search(ctx context.Context, query string, topK int) ([]Result, error)

	// Lifecycle
	Close() error
}

// Result holds a single search result with similarity score.
type Result struct {
	Key      string
	Content  string
	Score    float64
	Metadata map[string]string
}
