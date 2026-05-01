package provider

import (
	"context"
	"errors"
	"net"
	"testing"
)

func TestClassifyHTTPStatus(t *testing.T) {
	tests := []struct {
		status int
		kind   ErrorKind
	}{
		{status: 401, kind: ErrorAuth},
		{status: 429, kind: ErrorRateLimit},
		{status: 400, kind: ErrorInvalidRequest},
		{status: 500, kind: ErrorServer},
		{status: 418, kind: ErrorUnexpected},
	}

	for _, tt := range tests {
		err := ClassifyHTTPStatus(tt.status, "body")
		if err.Kind != tt.kind {
			t.Fatalf("status %d classified as %s, want %s", tt.status, err.Kind, tt.kind)
		}
	}
}

func TestClassifyError(t *testing.T) {
	timeoutErr := ClassifyError(context.DeadlineExceeded)
	if timeoutErr.Kind != ErrorTimeout {
		t.Fatalf("deadline classified as %s, want timeout", timeoutErr.Kind)
	}

	netTimeout := ClassifyError(timeoutNetError{})
	if netTimeout.Kind != ErrorTimeout {
		t.Fatalf("net timeout classified as %s, want timeout", netTimeout.Kind)
	}

	networkErr := ClassifyError(errors.New("dial failed"))
	if networkErr.Kind != ErrorNetwork {
		t.Fatalf("network error classified as %s, want network", networkErr.Kind)
	}

	if ClassifyError(nil) != nil {
		t.Fatal("nil error should classify to nil")
	}
}

func TestClassifiedErrorStringAndUnwrap(t *testing.T) {
	cause := errors.New("root cause")
	err := &ClassifiedError{Kind: ErrorUnexpected, Err: cause}
	if err.Error() != "root cause" {
		t.Fatalf("Error() = %q, want root cause", err.Error())
	}
	if !errors.Is(err, cause) {
		t.Fatal("expected ClassifiedError to unwrap cause")
	}

	err = &ClassifiedError{Kind: ErrorAuth, Message: "auth failed", Err: cause}
	if err.Error() != "auth failed" {
		t.Fatalf("Error() = %q, want auth failed", err.Error())
	}

	err = &ClassifiedError{Kind: ErrorRateLimit}
	if err.Error() != string(ErrorRateLimit) {
		t.Fatalf("Error() = %q, want kind string", err.Error())
	}
}

type timeoutNetError struct{}

func (timeoutNetError) Error() string   { return "timeout" }
func (timeoutNetError) Timeout() bool   { return true }
func (timeoutNetError) Temporary() bool { return true }

var _ net.Error = timeoutNetError{}
