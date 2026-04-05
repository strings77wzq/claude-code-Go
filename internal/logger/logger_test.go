package logger

import (
	"os"
	"path/filepath"
	"testing"
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
