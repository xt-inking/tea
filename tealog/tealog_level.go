package tealog

import (
	"context"
	"fmt"

	"github.com/tea-frame-go/tea/tealog/internal/mode"
)

const (
	LevelDebug = "[DEBUG]"
	LevelInfo  = "[INFO]"
	LevelError = "[ERROR]"
)

func (l *Logger) Debug(ctx context.Context, msg string) {
	if !mode.Debug {
		return
	}
	l.logWriter(ctx, LevelDebug, msg)
}

func (l *Logger) Info(ctx context.Context, msg string) {
	l.logWriter(ctx, LevelInfo, msg)
}

func (l *Logger) Error(ctx context.Context, msg string) {
	l.logWriter(ctx, LevelError, msg)
}

func (l *Logger) Debugf(ctx context.Context, format string, a ...any) {
	if !mode.Debug {
		return
	}
	l.Debug(ctx, fmt.Sprintf(format, a...))
}

func (l *Logger) Infof(ctx context.Context, format string, a ...any) {
	l.Info(ctx, fmt.Sprintf(format, a...))
}

func (l *Logger) Errorf(ctx context.Context, format string, a ...any) {
	l.Error(ctx, fmt.Sprintf(format, a...))
}

func (l *loggerStdout) Debug(ctx context.Context, msg string) {
	if !mode.Debug {
		return
	}
	l.l.logStdout(ctx, LevelDebug, msg)
}

func (l *loggerStdout) Info(ctx context.Context, msg string) {
	l.l.logStdout(ctx, LevelInfo, msg)
}

func (l *loggerStdout) Error(ctx context.Context, msg string) {
	l.l.logStdout(ctx, LevelError, msg)
}

func (l *loggerStdout) Debugf(ctx context.Context, format string, a ...any) {
	if !mode.Debug {
		return
	}
	l.Debug(ctx, fmt.Sprintf(format, a...))
}

func (l *loggerStdout) Infof(ctx context.Context, format string, a ...any) {
	l.Info(ctx, fmt.Sprintf(format, a...))
}

func (l *loggerStdout) Errorf(ctx context.Context, format string, a ...any) {
	l.Error(ctx, fmt.Sprintf(format, a...))
}
