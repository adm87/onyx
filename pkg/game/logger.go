package game

import (
	"fmt"
	"io"
)

type Logger interface {
	With(args ...any) Logger
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type loggerImpl struct {
	out io.Writer
}

func (l *loggerImpl) log(level, msg string, args ...any) {
	fmt.Fprintf(l.out, "[%s] %s\n", level, fmt.Sprintf(msg, args...))
}

func (l *loggerImpl) With(args ...any) Logger {
	return l
}

func (l *loggerImpl) Debug(msg string, args ...any) {
	l.log("DEBUG", msg, args...)
}

func (l *loggerImpl) Info(msg string, args ...any) {
	l.log("INFO", msg, args...)
}

func (l *loggerImpl) Warn(msg string, args ...any) {
	l.log("WARN", msg, args...)
}

func (l *loggerImpl) Error(msg string, args ...any) {
	l.log("ERROR", msg, args...)
}
