package logger

import (
	"time"
)

func TraceRequest(method, url string, body []byte) {
	logger := Log()
	if IsTraceHTTP() && len(body) > 0 {
		logger.Debug("HTTP request",
			"method", method,
			"url", url,
			"body_length", len(body),
		)
	} else {
		logger.Info("HTTP request",
			"method", method,
			"url", url,
		)
	}
}

func TraceResponse(statusCode int, duration time.Duration, body []byte) {
	logger := Log()
	if IsTraceHTTP() && len(body) > 0 {
		logger.Debug("HTTP response",
			"status", statusCode,
			"duration_ms", duration.Milliseconds(),
			"body_length", len(body),
		)
	} else {
		logger.Info("HTTP response",
			"status", statusCode,
			"duration_ms", duration.Milliseconds(),
		)
	}
}
