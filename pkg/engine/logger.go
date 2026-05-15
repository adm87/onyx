package engine

import (
	"fmt"
	"io"
	"time"
)

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type logger struct {
	out io.Writer
}

func NewLogger(stdout io.Writer) Logger {
	if stdout == nil {
		stdout = io.Discard
	}
	return &logger{
		out: stdout,
	}
}

func (l *logger) Debug(msg string, args ...any) {
	l.log("DEBUG", msg, args...)
}

func (l *logger) Info(msg string, args ...any) {
	l.log("INFO", msg, args...)
}

func (l *logger) Warn(msg string, args ...any) {
	l.log("WARN", msg, args...)
}

func (l *logger) Error(msg string, args ...any) {
	l.log("ERROR", msg, args...)
}

func (l *logger) log(level string, msg string, args ...any) {
	formattedMsg := msg
	if len(args) > 0 {
		formattedMsg = fmt.Sprintf(msg, args...)
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	logLine := fmt.Sprintf("%s [%s] %s\n", now, level, formattedMsg)

	l.out.Write([]byte(logLine))
}
