package engine

import "log/slog"

type LogLevel slog.Level

const (
	LogLevelDebug LogLevel = LogLevel(slog.LevelDebug)
	LogLevelInfo  LogLevel = LogLevel(slog.LevelInfo)
	LogLevelWarn  LogLevel = LogLevel(slog.LevelWarn)
	LogLevelError LogLevel = LogLevel(slog.LevelError)
)

type Logger struct {
	slog *slog.Logger
}

func NewLogger() *Logger {
	return &Logger{
		slog: slog.Default(),
	}
}

func (l *Logger) With(args ...any) *Logger {
	return &Logger{
		slog: l.slog.With(args...),
	}
}

func (l *Logger) SetLevel(level LogLevel) {
	slog.SetLogLoggerLevel(slog.Level(level))
}

func (l *Logger) Debug(msg string, args ...any) {
	l.slog.Debug(msg, args...)
}

func (l *Logger) Info(msg string, args ...any) {
	l.slog.Info(msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	l.slog.Warn(msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	l.slog.Error(msg, args...)
}
