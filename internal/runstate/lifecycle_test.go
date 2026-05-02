package runstate

import "testing"

func TestRequestLifecycleTerminalStates(t *testing.T) {
	lifecycle := NewRequestLifecycle("req-1")

	if lifecycle.State() != RequestStarted {
		t.Fatalf("initial state = %s, want %s", lifecycle.State(), RequestStarted)
	}

	if !lifecycle.Complete() {
		t.Fatal("expected first terminal transition to succeed")
	}
	if lifecycle.State() != RequestCompleted {
		t.Fatalf("state = %s, want %s", lifecycle.State(), RequestCompleted)
	}
	if lifecycle.Cancel() {
		t.Fatal("expected terminal lifecycle to reject later cancellation")
	}
	if lifecycle.State() != RequestCompleted {
		t.Fatalf("state changed after terminal transition: %s", lifecycle.State())
	}
}

func TestRequestLifecycleFailureAndCancellation(t *testing.T) {
	cancelled := NewRequestLifecycle("req-cancel")
	if !cancelled.Cancel() {
		t.Fatal("expected cancellation transition to succeed")
	}
	if cancelled.State() != RequestCancelled {
		t.Fatalf("state = %s, want %s", cancelled.State(), RequestCancelled)
	}

	failed := NewRequestLifecycle("req-fail")
	if !failed.Fail() {
		t.Fatal("expected failure transition to succeed")
	}
	if failed.State() != RequestFailed {
		t.Fatalf("state = %s, want %s", failed.State(), RequestFailed)
	}
}
