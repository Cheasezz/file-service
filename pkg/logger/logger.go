package logger

import (
	"context"
	"log/slog"
	"os"
)

type Logger interface {
	Debug(message string, args ...any)
	Info(message string, args ...any)
	Warn(message string, args ...any)
	Error(message string, args ...any)
	With(fields ...any) Logger
	Log(ctx context.Context, level int, msg string, args ...any)
}

type slogLG struct {
	l *slog.Logger
}

const (
	envLocal = "local"
	envProd  = "prod"
)

func New(env string) Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return &slogLG{log}
}

func (s *slogLG) Debug(message string, args ...any) {
	s.l.Debug(message, args...)
}

func (s *slogLG) Info(message string, args ...any) {
	s.l.Info(message, args...)
}

func (s *slogLG) Warn(message string, args ...any) {
	s.l.Warn(message, args...)
}

func (s *slogLG) Error(message string, args ...any) {
	s.l.Error(message, args...)
}

func (s *slogLG) With(fields ...any) Logger {
	return &slogLG{l: s.l.With(fields...)}
}

func (s *slogLG) Log(ctx context.Context, level int, msg string, args ...any) {
	s.l.Log(ctx, slog.Level(level), msg, args...)
}
