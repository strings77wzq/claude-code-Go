package runstate

import "sync"

type RequestState string

const (
	RequestStarted   RequestState = "started"
	RequestCompleted RequestState = "completed"
	RequestCancelled RequestState = "cancelled"
	RequestFailed    RequestState = "failed"
)

type RequestLifecycle struct {
	mu    sync.Mutex
	id    string
	state RequestState
}

func NewRequestLifecycle(id string) *RequestLifecycle {
	return &RequestLifecycle{
		id:    id,
		state: RequestStarted,
	}
}

func (l *RequestLifecycle) ID() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.id
}

func (l *RequestLifecycle) State() RequestState {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.state
}

func (l *RequestLifecycle) Complete() bool {
	return l.transition(RequestCompleted)
}

func (l *RequestLifecycle) Cancel() bool {
	return l.transition(RequestCancelled)
}

func (l *RequestLifecycle) Fail() bool {
	return l.transition(RequestFailed)
}

func (l *RequestLifecycle) transition(next RequestState) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.state != RequestStarted {
		return false
	}
	l.state = next
	return true
}
