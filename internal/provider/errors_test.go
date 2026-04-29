package provider

import (
	"context"
	"errors"
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

	networkErr := ClassifyError(errors.New("dial failed"))
	if networkErr.Kind != ErrorNetwork {
		t.Fatalf("network error classified as %s, want network", networkErr.Kind)
	}
}
