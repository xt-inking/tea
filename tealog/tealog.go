package tealog

import (
	"context"
	"fmt"
	"os"

	"github.com/tea-frame-go/tea/internal/bufferpool"
)

type Logger struct {
	newRecord    func(ctx context.Context, level string, msg string) Record
	writerCloser writerCloser
	chanWriter   chan Record
	chanStdout   chan Record
	close        chan struct{}
	closed       chan struct{}
}

func New(newRecord func(ctx context.Context, level string, msg string) Record, writerCloser writerCloser) *Logger {
	l := &Logger{
		newRecord:    newRecord,
		writerCloser: writerCloser,
		chanWriter:   make(chan Record, 1024),
		chanStdout:   make(chan Record, 1024),
		close:        make(chan struct{}),
		closed:       make(chan struct{}),
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
	buf := bufferpool.NewBuffer(bufPool)
	r.WriteBuffer(buf)
	_, err := l.writerCloser.Writer(r).Write(*buf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "tealog: Write error: %v; record:\n%s", err, *buf)
	}
	buf.Free(bufPool)
}

func (l *Logger) handleChanStdout(r Record) {
	for {
		select {
		case r := <-l.chanWriter:
			l.handleChanWriter(r)
		default:
			buf := bufferpool.NewBuffer(bufPool)
			r.WriteBuffer(buf)
			os.Stdout.Write(*buf)
			buf.Free(bufPool)
			return
		}
	}
}

var bufPool = bufferpool.New()

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

func (l *Logger) WithStdout() *loggerStdout {
	return &loggerStdout{
		l: l,
	}
}

type loggerStdout struct {
	l *Logger
}

func (l *Logger) logWriter(ctx context.Context, level string, msg string) {
	r := l.newRecord(ctx, level, msg)
	l.chanWriter <- r
}

func (l *Logger) logStdout(ctx context.Context, level string, msg string) {
	r := l.newRecord(ctx, level, msg)
	l.chanWriter <- r
	l.chanStdout <- r
}
