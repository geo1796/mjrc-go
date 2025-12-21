package logger

import (
	"log/slog"
	"os"
)

var logger Logger

var _ = initLogger()

func initLogger() bool {
	var lvl slog.Level

	if os.Getenv("ACTIVE_PROFILE") == "prod" {
		lvl = slog.LevelInfo
	} else {
		lvl = slog.LevelDebug
	}

	logger = newSlogAdapter(slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl})))

	return true
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
	logger.Debug(s, fields...)
}

func Info(s string, fields ...Entry) {
	logger.Info(s, fields...)
}

func Warn(s string, fields ...Entry) {
	logger.Warn(s, fields...)
}

func Error(s string, fields ...Entry) {
	logger.Error(s, fields...)
}
