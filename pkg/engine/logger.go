package engine

import (
	"fmt"
	"io"
)

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type logger struct {
	writer io.Writer
}

func newLogger(writer io.Writer) *logger {
	return &logger{
		writer: writer,
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

func (l *logger) log(level, msg string, args ...any) {
	formattedMsg := fmt.Sprintf(msg, args...)
	fmt.Fprintf(l.writer, "[%s] %s\n", level, formattedMsg)
}
