package memory

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"
	"sync"
)

// MemStore is an in-memory implementation of Store.
// KV data is stored in a map. Semantic search uses brute-force cosine similarity.
type MemStore struct {
	mu       sync.RWMutex
	data     map[string][]byte
	index    map[string]string // key → content for semantic search
	vectors  map[string][]float64
	embedder Embedder
}

// Embedder converts text to a vector representation.
type Embedder interface {
	Embed(ctx context.Context, text string) ([]float64, error)
}

// NewMemStore creates a new in-memory store.
func NewMemStore(embedder Embedder) *MemStore {
	return &MemStore{
		data:     make(map[string][]byte),
		index:    make(map[string]string),
		vectors:  make(map[string][]float64),
		embedder: embedder,
	}
}

func (s *MemStore) Put(_ context.Context, key string, value []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
	return nil
}

func (s *MemStore) Get(_ context.Context, key string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.data[key]
	if !ok {
		return nil, fmt.Errorf("key not found: %s", key)
	}
	return v, nil
}

func (s *MemStore) Delete(_ context.Context, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
	delete(s.index, key)
	delete(s.vectors, key)
	return nil
}

func (s *MemStore) List(_ context.Context, prefix string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var keys []string
	for k := range s.data {
		if strings.HasPrefix(k, prefix) {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	return keys, nil
}

func (s *MemStore) Index(ctx context.Context, key string, content string) error {
	if s.embedder == nil {
		return fmt.Errorf("no embedder configured")
	}
	vec, err := s.embedder.Embed(ctx, content)
	if err != nil {
		return fmt.Errorf("embedding failed: %w", err)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.index[key] = content
	s.vectors[key] = vec
	return nil
}

func (s *MemStore) Search(ctx context.Context, query string, topK int) ([]Result, error) {
	if s.embedder == nil {
		return nil, fmt.Errorf("no embedder configured")
	}
	queryVec, err := s.embedder.Embed(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("embedding failed: %w", err)
	}
	s.mu.RLock()
	defer s.mu.RUnlock()

	type scored struct {
		key   string
		score float64
	}
	var scores []scored
	for key, vec := range s.vectors {
		sc := cosineSimilarity(queryVec, vec)
		scores = append(scores, scored{key: key, score: sc})
	}
	sort.Slice(scores, func(i, j int) bool { return scores[i].score > scores[j].score })

	if topK > len(scores) {
		topK = len(scores)
	}
	results := make([]Result, topK)
	for i := 0; i < topK; i++ {
		results[i] = Result{
			Key:     scores[i].key,
			Content: s.index[scores[i].key],
			Score:   scores[i].score,
		}
	}
	return results, nil
}

func (s *MemStore) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = nil
	s.index = nil
	s.vectors = nil
	return nil
}

func cosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}
	var dot, normA, normB float64
	for i := range a {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}
