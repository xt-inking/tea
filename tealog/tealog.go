package tealog

import (
	"context"
	"fmt"
	"os"
)

type Logger struct {
	recordHandler recordHandler
	writerCloser  writerCloser
	chanWriter    chan Record
	chanStdout    chan Record
	close         chan struct{}
	closed        chan struct{}
}

func New(recordHandler recordHandler, writerCloser writerCloser) *Logger {
	l := &Logger{
		recordHandler: recordHandler,
		writerCloser:  writerCloser,
		chanWriter:    make(chan Record, 1024),
		chanStdout:    make(chan Record, 1024),
		close:         make(chan struct{}),
		closed:        make(chan struct{}),
	}
	go func() {
		for {
			select {
			case r := <-l.chanWriter:
				l.handleChanWriter(r)
			case r := <-l.chanStdout:
				l.handleChanStdout(r)
			case <-l.close:
				l.handleClose()
				return
			}
		}
	}()
	return l
}

func (l *Logger) handleChanWriter(r Record) {
	buf := newBuffer()
	l.recordHandler.HandleRecord(buf, r)
	_, err := l.writerCloser.Writer(r).Write(*buf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "tealog: Write error: %v; record:\n%s", err, *buf)
	}
	buf.free()
}

func (l *Logger) handleChanStdout(r Record) {
	for {
		select {
		case r := <-l.chanWriter:
			l.handleChanWriter(r)
		default:
			buf := newBuffer()
			l.recordHandler.HandleRecord(buf, r)
			os.Stdout.Write(*buf)
			buf.free()
			return
		}
	}
}

func (l *Logger) handleClose() {
	for {
		select {
		case r := <-l.chanWriter:
			l.handleChanWriter(r)
		case r := <-l.chanStdout:
			l.handleChanStdout(r)
		default:
			close(l.closed)
			return
		}
	}
}

func (l *Logger) Close() {
	close(l.close)
	<-l.closed
	l.writerCloser.Close()
}

func (l *Logger) Info(ctx context.Context, msg string) {
	l.logWriter(ctx, "[INFO]", msg)
}

func (l *Logger) Error(ctx context.Context, msg string) {
	l.logWriter(ctx, "[ERROR]", msg)
}

func (l *Logger) Infof(ctx context.Context, format string, a ...any) {
	l.Info(ctx, fmt.Sprintf(format, a...))
}

func (l *Logger) Errorf(ctx context.Context, format string, a ...any) {
	l.Error(ctx, fmt.Sprintf(format, a...))
}

func (l *Logger) WithStdout() *loggerStdout {
	return &loggerStdout{
		l: l,
	}
}

type loggerStdout struct {
	l *Logger
}

func (l *loggerStdout) Info(ctx context.Context, msg string) {
	l.l.logStdout(ctx, "[INFO]", msg)
}

func (l *loggerStdout) Error(ctx context.Context, msg string) {
	l.l.logStdout(ctx, "[ERROR]", msg)
}

func (l *loggerStdout) Infof(ctx context.Context, format string, a ...any) {
	l.Info(ctx, fmt.Sprintf(format, a...))
}

func (l *loggerStdout) Errorf(ctx context.Context, format string, a ...any) {
	l.Error(ctx, fmt.Sprintf(format, a...))
}

func (l *Logger) logWriter(ctx context.Context, level string, msg string) {
	r := newRecord(ctx, level, msg)
	l.chanWriter <- r
}

func (l *Logger) logStdout(ctx context.Context, level string, msg string) {
	r := newRecord(ctx, level, msg)
	l.chanWriter <- r
	l.chanStdout <- r
}
