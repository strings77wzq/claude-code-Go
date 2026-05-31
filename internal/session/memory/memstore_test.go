package memory

import (
	"context"
	"math"
	"testing"
)

// testEmbedder returns deterministic vectors for testing.
type testEmbedder struct{}

func (e *testEmbedder) Embed(_ context.Context, text string) ([]float64, error) {
	// Simple hash-based embedding: each char contributes to different dimensions
	vec := make([]float64, 8)
	for i, c := range text {
		vec[i%8] += float64(c) / 1000.0
	}
	// Normalize
	var norm float64
	for _, v := range vec {
		norm += v * v
	}
	if norm > 0 {
		for i := range vec {
			vec[i] /= math.Sqrt(norm)
		}
	}
	return vec, nil
}

func TestMemStore_PutGet(t *testing.T) {
	store := NewMemStore(nil)
	ctx := context.Background()

	if err := store.Put(ctx, "key1", []byte("value1")); err != nil {
		t.Fatalf("Put() error = %v", err)
	}
	got, err := store.Get(ctx, "key1")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if string(got) != "value1" {
		t.Errorf("Get() = %q, want %q", string(got), "value1")
	}
}

func TestMemStore_GetNotFound(t *testing.T) {
	store := NewMemStore(nil)
	ctx := context.Background()

	_, err := store.Get(ctx, "nonexistent")
	if err == nil {
		t.Error("Get() expected error for missing key")
	}
}

func TestMemStore_Delete(t *testing.T) {
	store := NewMemStore(nil)
	ctx := context.Background()

	store.Put(ctx, "key1", []byte("value1"))
	if err := store.Delete(ctx, "key1"); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	_, err := store.Get(ctx, "key1")
	if err == nil {
		t.Error("Get() expected error after Delete()")
	}
}

func TestMemStore_List(t *testing.T) {
	store := NewMemStore(nil)
	ctx := context.Background()

	store.Put(ctx, "user:1", []byte("a"))
	store.Put(ctx, "user:2", []byte("b"))
	store.Put(ctx, "session:1", []byte("c"))

	keys, err := store.List(ctx, "user:")
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(keys) != 2 {
		t.Errorf("List() returned %d keys, want 2", len(keys))
	}
	if keys[0] != "user:1" || keys[1] != "user:2" {
		t.Errorf("List() = %v, want [user:1 user:2]", keys)
	}
}

func TestMemStore_IndexAndSearch(t *testing.T) {
	store := NewMemStore(&testEmbedder{})
	ctx := context.Background()

	store.Index(ctx, "doc1", "golang testing patterns")
	store.Index(ctx, "doc2", "python web framework")
	store.Index(ctx, "doc3", "golang concurrency goroutines")

	results, err := store.Search(ctx, "golang", 2)
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}
	if len(results) == 0 {
		t.Fatal("Search() returned 0 results")
	}
	// First result should be most similar to "golang"
	if results[0].Key != "doc1" && results[0].Key != "doc3" {
		t.Errorf("Search() top result = %s, want doc1 or doc3", results[0].Key)
	}
}

func TestMemStore_SearchNoEmbedder(t *testing.T) {
	store := NewMemStore(nil)
	ctx := context.Background()

	_, err := store.Search(ctx, "query", 5)
	if err == nil {
		t.Error("Search() expected error when no embedder configured")
	}
}

func TestMemStore_IndexNoEmbedder(t *testing.T) {
	store := NewMemStore(nil)
	ctx := context.Background()

	err := store.Index(ctx, "key", "content")
	if err == nil {
		t.Error("Index() expected error when no embedder configured")
	}
}

func TestMemStore_Close(t *testing.T) {
	store := NewMemStore(nil)
	ctx := context.Background()

	store.Put(ctx, "key1", []byte("value1"))
	if err := store.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
	// After close, data should be nil
	if store.data != nil {
		t.Error("Close() did not clear data")
	}
}

func TestCosineSimilarity(t *testing.T) {
	tests := []struct {
		name string
		a, b []float64
		want float64
	}{
		{"identical", []float64{1, 0, 0}, []float64{1, 0, 0}, 1.0},
		{"orthogonal", []float64{1, 0, 0}, []float64{0, 1, 0}, 0.0},
		{"opposite", []float64{1, 0}, []float64{-1, 0}, -1.0},
		{"empty", nil, nil, 0.0},
		{"different lengths", []float64{1}, []float64{1, 2}, 0.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cosineSimilarity(tt.a, tt.b)
			if diff := got - tt.want; diff > 0.001 || diff < -0.001 {
				t.Errorf("cosineSimilarity() = %f, want %f", got, tt.want)
			}
		})
	}
}

// Verify MemStore implements Store
var _ Store = (*MemStore)(nil)
