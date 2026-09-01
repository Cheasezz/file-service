package logger

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type LoggerMock struct {
	mock.Mock
}

func (m *LoggerMock) Debug(message string, args ...any) {
	m.Called(message, args)
}

func (m *LoggerMock) Info(message string, args ...any) {
	m.Called(message, args)
}

func (m *LoggerMock) Warn(message string, args ...any) {
	m.Called(message, args)
}

func (m *LoggerMock) Error(message string, args ...any) {
	m.Called(message, args)
}

func (m *LoggerMock) With(fields ...any) Logger {
	args := m.Called(fields...)
	return args.Get(0).(Logger)
}

func (m *LoggerMock) Log(ctx context.Context, level int, msg string, args ...any) {
	m.Called(ctx, level, msg, args)
}
