package logger

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestInit(t *testing.T) {
	err := Init(false, false)
	if err != nil {
		t.Errorf("Init failed: %v", err)
	}
	defer Cleanup()

	homeDir, _ := os.UserHomeDir()
	logPath := filepath.Join(homeDir, ".go-code", "go-code.log")

	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		t.Error("Log file should be created")
	}
}

func TestInitDebugMode(t *testing.T) {
	err := Init(true, false)
	if err != nil {
		t.Errorf("Init with debug mode failed: %v", err)
	}
	defer Cleanup()
}

func TestLog(t *testing.T) {
	l := Log()
	if l == nil {
		t.Error("Log should not return nil")
	}
}

func TestLogBeforeInit(t *testing.T) {
	Cleanup()

	l := Log()
	if l == nil {
		t.Error("Log should return a valid logger even before Init")
	}
}

func TestSetDebug(t *testing.T) {
	Init(false, false)
	defer Cleanup()

	SetDebug(true)
	if !IsDebug() {
		t.Error("IsDebug should return true after SetDebug(true)")
	}

	SetDebug(false)
	if IsDebug() {
		t.Error("IsDebug should return false after SetDebug(false)")
	}
}

func TestSetTraceHTTP(t *testing.T) {
	Init(false, false)
	defer Cleanup()

	SetTraceHTTP(true)
	if !IsTraceHTTP() {
		t.Error("IsTraceHTTP should return true after SetTraceHTTP(true)")
	}

	SetTraceHTTP(false)
	if IsTraceHTTP() {
		t.Error("IsTraceHTTP should return false after SetTraceHTTP(false)")
	}
}

func TestAPIRequestStart(t *testing.T) {
	Init(false, false)
	defer Cleanup()

	APIRequestStart("claude-sonnet-4", 5)
}

func TestAPIResponseReceived(t *testing.T) {
	Init(false, false)
	defer Cleanup()

	APIResponseReceived(1500, 200)
}

func TestToolExecuted(t *testing.T) {
	Init(false, false)
	defer Cleanup()

	ToolExecuted("Read", 50)
}

func TestSessionStarted(t *testing.T) {
	Init(false, false)
	defer Cleanup()

	SessionStarted()
}

func TestSessionEnded(t *testing.T) {
	Init(false, false)
	defer Cleanup()

	SessionEnded()
}

func TestWith(t *testing.T) {
	Init(false, false)
	defer Cleanup()

	l := With("key", "value")
	if l == nil {
		t.Error("With should not return nil")
	}
}

func TestCleanup(t *testing.T) {
	Init(false, false)

	err := Cleanup()
	if err != nil {
		t.Errorf("Cleanup failed: %v", err)
	}
}

func TestMultipleInit(t *testing.T) {
	Init(false, false)
	Init(true, true)
	defer Cleanup()
}

func TestMultiHandlerDelegatesRecordsAttrsAndGroups(t *testing.T) {
	var first, second bytes.Buffer
	handler := newMultiHandler(
		slog.NewTextHandler(&first, &slog.HandlerOptions{Level: slog.LevelInfo}),
		slog.NewTextHandler(&second, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	if !handler.Enabled(context.Background(), slog.LevelInfo) {
		t.Fatal("expected info level to be enabled by at least one delegate")
	}
	if !handler.Enabled(context.Background(), slog.LevelDebug) {
		t.Fatal("expected debug level to be enabled by one delegate")
	}

	grouped := handler.WithAttrs([]slog.Attr{slog.String("request_id", "abc")}).WithGroup("api")
	if err := grouped.Handle(context.Background(), slog.NewRecord(timeNow(), slog.LevelInfo, "hello", 0)); err != nil {
		t.Fatalf("Handle returned error: %v", err)
	}

	for name, output := range map[string]string{"first": first.String(), "second": second.String()} {
		if !strings.Contains(output, "hello") || !strings.Contains(output, "request_id") {
			t.Fatalf("%s handler did not receive grouped record with attrs: %q", name, output)
		}
	}
}

func TestAPIError(t *testing.T) {
	if err := Init(false, false); err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer Cleanup()

	APIError(errors.New("provider failed"), 500)
}

func timeNow() time.Time {
	return time.Unix(1700000000, 0)
}
