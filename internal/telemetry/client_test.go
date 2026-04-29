package telemetry

import (
	"os"
	"path/filepath"
	"testing"
)

func TestClient_DisabledByDefault(t *testing.T) {
	c := NewClient(t.TempDir())
	if c.IsEnabled() {
		t.Error("expected client to be disabled by default")
	}
}

func TestClient_Enable(t *testing.T) {
	dir := t.TempDir()
	c := NewClient(dir)

	if err := c.Enable(); err != nil {
		t.Fatal(err)
	}
	if !c.IsEnabled() {
		t.Error("expected client to be enabled after Enable()")
	}

	// Verify data directory was created
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Error("expected data directory to exist after Enable()")
	}
}

func TestClient_Disable(t *testing.T) {
	c := NewClient(t.TempDir())
	c.Enable()
	c.Disable()
	if c.IsEnabled() {
		t.Error("expected client to be disabled after Disable()")
	}
}

func TestClient_Track_WhenDisabled_ReturnsNil(t *testing.T) {
	c := NewClient(t.TempDir())

	err := c.Track("test", map[string]interface{}{"key": "value"})
	if err != nil {
		t.Errorf("expected nil error when tracking while disabled, got: %v", err)
	}

	// No telemetry file should exist
	events, err := c.GetEvents()
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 0 {
		t.Errorf("expected 0 events when disabled, got %d", len(events))
	}
}

func TestClient_Track_WhenEnabled(t *testing.T) {
	dir := t.TempDir()
	c := NewClient(dir)

	if err := c.Enable(); err != nil {
		t.Fatal(err)
	}

	err := c.Track("test_event", map[string]interface{}{"foo": "bar"})
	if err != nil {
		t.Fatal(err)
	}

	// Verify file was written
	filename := filepath.Join(dir, "telemetry.jsonl")
	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 {
		t.Fatal("expected non-empty telemetry file")
	}

	// Verify via GetEvents
	events, err := c.GetEvents()
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].EventType != "test_event" {
		t.Errorf("expected event type 'test_event', got %q", events[0].EventType)
	}
	if v, ok := events[0].Data["foo"]; !ok || v != "bar" {
		t.Errorf("expected data.foo = 'bar', got %v", v)
	}
}

func TestClient_Track_ErrorWithInvalidDataDir(t *testing.T) {
	dir := t.TempDir()
	c := NewClient(dir)

	if err := c.Enable(); err != nil {
		t.Fatal(err)
	}

	// Replace the data directory with a regular file so Track fails
	os.RemoveAll(dir)
	f, err := os.Create(dir)
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	err = c.Track("test", nil)
	if err == nil {
		t.Error("expected error when data directory is invalid, got nil")
	}
}

func TestClient_GetEvents_NoFile(t *testing.T) {
	c := NewClient(t.TempDir())
	events, err := c.GetEvents()
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 0 {
		t.Errorf("expected 0 events when no file exists, got %d", len(events))
	}
}
