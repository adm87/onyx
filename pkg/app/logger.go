package app

import (
	"fmt"
	"io"
	"runtime"
)

func debugStyle(msg string) string {
	if runtime.GOOS == "windows" {
		return msg // No styling on Windows
	}
	return fmt.Sprintf("\x1b[90m%s\x1b[0m", msg) // grey
}

func infoStyle(msg string) string {
	if runtime.GOOS == "windows" {
		return msg // No styling on Windows
	}
	return fmt.Sprintf("\x1b[37m%s\x1b[0m", msg) // white
}

func warnStyle(msg string) string {
	if runtime.GOOS == "windows" {
		return msg // No styling on Windows
	}
	return fmt.Sprintf("\x1b[33m%s\x1b[0m", msg) // yellow
}

func errorStyle(msg string) string {
	if runtime.GOOS == "windows" {
		return msg // No styling on Windows
	}
	return fmt.Sprintf("\x1b[31m%s\x1b[0m", msg) // red
}

type LogLevel uint8

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelSilent
)

type Logger struct {
	lvl    LogLevel
	writer io.Writer
}

func NewLogger(writer io.Writer, lvl LogLevel) *Logger {
	if writer == nil {
		writer = io.Discard
	}
	return &Logger{
		writer: writer,
		lvl:    lvl,
	}
}

func (l *Logger) log(level, msg string, style func(string) string, args ...any) {
	logMsg := fmt.Sprintf(msg, args...)
	fmt.Fprintln(l.writer, style(logMsg))
}

func (l *Logger) Debug(msg string, args ...any) {
	if l.lvl > LogLevelDebug {
		return
	}
	l.log("DEBUG", msg, debugStyle, args...)
}

func (l *Logger) Info(msg string, args ...any) {
	if l.lvl > LogLevelInfo {
		return
	}
	l.log("INFO", msg, infoStyle, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	if l.lvl > LogLevelWarn {
		return
	}
	l.log("WARN", msg, warnStyle, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	if l.lvl > LogLevelError {
		return
	}
	l.log("ERROR", msg, errorStyle, args...)
}
