package logger

import (
	"testing"
	"time"
)

func TestTraceRequestWithTraceHTTPEnabled(t *testing.T) {
	Init(false, true)
	defer Cleanup()

	TraceRequest("POST", "https://api.example.com/v1/chat", []byte(`{"model": "claude"}`))
}

func TestTraceRequestWithTraceHTTPDisabled(t *testing.T) {
	Init(false, false)
	defer Cleanup()

	TraceRequest("POST", "https://api.example.com/v1/chat", []byte(`{"model": "claude"}`))
}

func TestTraceRequestWithEmptyBody(t *testing.T) {
	Init(false, true)
	defer Cleanup()

	TraceRequest("GET", "https://api.example.com/models", nil)
}

func TestTraceResponseWithTraceHTTPEnabled(t *testing.T) {
	Init(false, true)
	defer Cleanup()

	body := []byte(`{"content": "response"}`)
	duration := 1500 * time.Millisecond

	TraceResponse(200, duration, body)
}

func TestTraceResponseWithTraceHTTPDisabled(t *testing.T) {
	Init(false, false)
	defer Cleanup()

	body := []byte(`{"content": "response"}`)
	duration := 1500 * time.Millisecond

	TraceResponse(200, duration, body)
}

func TestTraceResponseWithEmptyBody(t *testing.T) {
	Init(false, true)
	defer Cleanup()

	duration := 500 * time.Millisecond

	TraceResponse(204, duration, nil)
}

func TestTraceResponseErrorStatus(t *testing.T) {
	Init(false, true)
	defer Cleanup()

	body := []byte(`{"error": "invalid request"}`)
	duration := 100 * time.Millisecond

	TraceResponse(400, duration, body)
}

func TestTraceHTTPDisabledDoesNotLogBody(t *testing.T) {
	Init(false, false)
	defer Cleanup()

	body := []byte("sensitive data that should not be logged")

	TraceRequest("POST", "https://api.example.com", body)
	TraceResponse(200, time.Second, body)
}
