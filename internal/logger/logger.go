package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
)

var (
	globalLogger *slog.Logger
	mu           sync.RWMutex

	debug     bool
	traceHTTP bool
	logFile   *os.File
)

type multiHandler struct {
	handlers []slog.Handler
}

func (m *multiHandler) Handle(ctx context.Context, record slog.Record) error {
	for _, h := range m.handlers {
		if err := h.Handle(ctx, record); err != nil {
			return err
		}
	}
	return nil
}

func (m *multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range m.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (m *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		handlers[i] = h.WithAttrs(attrs)
	}
	return &multiHandler{handlers: handlers}
}

func (m *multiHandler) WithGroup(name string) slog.Handler {
	handlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		handlers[i] = h.WithGroup(name)
	}
	return &multiHandler{handlers: handlers}
}

func newMultiHandler(handlers ...slog.Handler) *multiHandler {
	return &multiHandler{handlers: handlers}
}

func Init(debugMode bool, httpTrace bool) error {
	mu.Lock()
	defer mu.Unlock()

	debug = debugMode
	traceHTTP = httpTrace

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	logDir := filepath.Join(homeDir, ".go-code")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	logPath := filepath.Join(logDir, "go-code.log")

	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	logFile = file

	fileHandler := slog.NewJSONHandler(file, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	var handler slog.Handler
	if debugMode {
		handler = newMultiHandler(
			fileHandler,
			slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
				Level: slog.LevelInfo,
			}),
		)
	} else {
		handler = fileHandler
	}

	globalLogger = slog.New(handler)
	slog.SetDefault(globalLogger)

	return nil
}

func Log() *slog.Logger {
	mu.RLock()
	defer mu.RUnlock()
	if globalLogger == nil {
		return slog.New(slog.NewJSONHandler(io.Discard, nil))
	}
	return globalLogger
}

func SetDebug(enabled bool) {
	mu.Lock()
	defer mu.Unlock()
	debug = enabled
}

func SetTraceHTTP(enabled bool) {
	mu.Lock()
	defer mu.Unlock()
	traceHTTP = enabled
}

func IsTraceHTTP() bool {
	mu.RLock()
	defer mu.RUnlock()
	return traceHTTP
}

func IsDebug() bool {
	mu.RLock()
	defer mu.RUnlock()
	return debug
}

func Cleanup() error {
	mu.Lock()
	defer mu.Unlock()

	if logFile != nil {
		err := logFile.Close()
		logFile = nil
		globalLogger = nil
		return err
	}
	return nil
}

func With(attrs ...any) *slog.Logger {
	return Log().With(attrs...)
}

func APIRequestStart(model string, messageCount int) {
	Log().Info("API request sent", "model", model, "messages", messageCount)
}

func APIResponseReceived(durationMs int64, tokens int) {
	Log().Info("API response received", "duration_ms", durationMs, "tokens", tokens)
}

func ToolExecuted(name string, durationMs int64) {
	Log().Info("Tool executed", "tool", name, "duration_ms", durationMs)
}

func APIError(err error, statusCode int) {
	Log().Error("API error", "error", err.Error(), "status", statusCode)
}

func SessionStarted() {
	Log().Info("Session started")
}

func SessionEnded() {
	Log().Info("Session ended")
}
