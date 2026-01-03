package logger

import (
	"log/slog"
	"os"
)

const (
	LevelDebug = slog.LevelDebug
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
)

var globalLogger = New(LevelDebug)

func SetGlobalLogger(logger Logger) {
	globalLogger = logger
}

func New(level slog.Level) Logger {
	return &slogAdapter{slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}))}
}

type Entry struct {
	Key   string
	Value any
}

type Logger interface {
	Debug(string, ...Entry)
	Info(string, ...Entry)
	Warn(string, ...Entry)
	Error(string, ...Entry)
}

func Err(err error) Entry {
	return Entry{
		Key:   "error",
		Value: err,
	}
}

func Any(key string, value any) Entry {
	return Entry{
		Key:   key,
		Value: value,
	}
}

func Debug(s string, fields ...Entry) {
	globalLogger.Debug(s, fields...)
}

func Info(s string, fields ...Entry) {
	globalLogger.Info(s, fields...)
}

func Warn(s string, fields ...Entry) {
	globalLogger.Warn(s, fields...)
}

func Error(s string, fields ...Entry) {
	globalLogger.Error(s, fields...)
}
