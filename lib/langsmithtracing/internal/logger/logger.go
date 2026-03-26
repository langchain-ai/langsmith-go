package logger

import "log"

// Logger is an interface that the tracing client uses for all diagnostic output.
// Users can supply their own implementation via [langsmithtracing.WithLogger].
type Logger interface {
	Debug(msg string, keyvals ...any)
	Info(msg string, keyvals ...any)
	Warn(msg string, keyvals ...any)
	Error(msg string, keyvals ...any)
}

// DefaultLogger writes to the standard log package with a [langsmith] prefix.
type DefaultLogger struct{}

func (DefaultLogger) Debug(msg string, keyvals ...any) {
	logWithKeyvals("DEBUG", msg, keyvals)
}

func (DefaultLogger) Info(msg string, keyvals ...any) {
	logWithKeyvals("INFO", msg, keyvals)
}

func (DefaultLogger) Warn(msg string, keyvals ...any) {
	logWithKeyvals("WARN", msg, keyvals)
}

func (DefaultLogger) Error(msg string, keyvals ...any) {
	logWithKeyvals("ERROR", msg, keyvals)
}

func logWithKeyvals(level, msg string, keyvals []any) {
	if len(keyvals) == 0 {
		log.Printf("[langsmith] %s: %s", level, msg)
		return
	}
	log.Printf("[langsmith] %s: %s %v", level, msg, keyvals)
}
