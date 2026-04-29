package provider

import (
	"context"
	"fmt"
	"net"
	"net/http"
)

type ErrorKind string

const (
	ErrorAuth           ErrorKind = "auth"
	ErrorRateLimit      ErrorKind = "rate_limit"
	ErrorTimeout        ErrorKind = "timeout"
	ErrorServer         ErrorKind = "server"
	ErrorNetwork        ErrorKind = "network"
	ErrorInvalidRequest ErrorKind = "invalid_request"
	ErrorUnexpected     ErrorKind = "unexpected"
)

type ClassifiedError struct {
	Kind       ErrorKind
	StatusCode int
	Message    string
	Err        error
}

func (e *ClassifiedError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	if e.Err != nil {
		return e.Err.Error()
	}
	return string(e.Kind)
}

func (e *ClassifiedError) Unwrap() error {
	return e.Err
}

func ClassifyHTTPStatus(statusCode int, body string) *ClassifiedError {
	switch statusCode {
	case http.StatusUnauthorized, http.StatusForbidden:
		return &ClassifiedError{Kind: ErrorAuth, StatusCode: statusCode, Message: fmt.Sprintf("provider authentication failed (%d)", statusCode)}
	case http.StatusTooManyRequests:
		return &ClassifiedError{Kind: ErrorRateLimit, StatusCode: statusCode, Message: "provider rate limit reached"}
	case http.StatusBadRequest, http.StatusUnprocessableEntity:
		return &ClassifiedError{Kind: ErrorInvalidRequest, StatusCode: statusCode, Message: fmt.Sprintf("provider rejected request (%d): %s", statusCode, body)}
	}
	if statusCode >= 500 {
		return &ClassifiedError{Kind: ErrorServer, StatusCode: statusCode, Message: fmt.Sprintf("provider server error (%d)", statusCode)}
	}
	return &ClassifiedError{Kind: ErrorUnexpected, StatusCode: statusCode, Message: fmt.Sprintf("unexpected provider status (%d): %s", statusCode, body)}
}

func ClassifyError(err error) *ClassifiedError {
	if err == nil {
		return nil
	}
	if err == context.DeadlineExceeded {
		return &ClassifiedError{Kind: ErrorTimeout, Message: "provider request timed out", Err: err}
	}
	if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		return &ClassifiedError{Kind: ErrorTimeout, Message: "provider request timed out", Err: err}
	}
	return &ClassifiedError{Kind: ErrorNetwork, Message: "provider network error", Err: err}
}
