package integration

import (
	"testing"
	"time"
)

// BenchmarkAgentProcessing measures agent processing performance
func BenchmarkAgentProcessing(b *testing.B) {
	ag := setupTestAgent()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ag.Process("Simple query")
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkToolExecution measures tool execution performance
func BenchmarkToolExecution(b *testing.B) {
	ag := setupTestAgent()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ag.ExecuteTool("Read", map[string]interface{}{
			"file_path": "test.txt",
		})
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkSessionSave measures session persistence performance
func BenchmarkSessionSave(b *testing.B) {
	ag := setupTestAgent()
	// Create large session
	for i := 0; i < 100; i++ {
		ag.Process(fmt.Sprintf("Message %d", i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := ag.Session.Save()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkContextCompaction measures context compaction performance
func BenchmarkContextCompaction(b *testing.B) {
	ag := setupTestAgent()
	// Create large context
	for i := 0; i < 50; i++ {
		ag.AddMessage(agent.Message{
			Role:    "user",
			Content: strings.Repeat("x", 1000),
		})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := ag.Context.Compact()
		if err != nil {
			b.Fatal(err)
		}
	}
}
