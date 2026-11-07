package tealog

import (
	"context"
	"fmt"
	"os"

	"github.com/xt-inking/tea/internal/bufferpool"
	"github.com/xt-inking/tea/internal/profile"
)

type Logger[NR NewRecord[R], WC WriterCloser[R], R Record] struct {
	newRecord    NR
	writerCloser WC
	chanRecord   chan R
	close        chan struct{}
	closed       chan struct{}
}

func New[NR NewRecord[R], WC WriterCloser[R], R Record](newRecord NR, writerCloser WC) *Logger[NR, WC, R] {
	l := &Logger[NR, WC, R]{
		newRecord:    newRecord,
		writerCloser: writerCloser,
		chanRecord:   make(chan R, 1024),
		close:        make(chan struct{}),
		closed:       make(chan struct{}),
	}
	go func(l *Logger[NR, WC, R]) {
		for {
			select {
			case r := <-l.chanRecord:
				l.selectChanRecord(r)
			case <-l.close:
				l.selectClose()
				return
			}
		}
	}(l)
	return l
}

func (l *Logger[NR, WC, R]) selectChanRecord(r R) {
	b := bufPool.Get()
	r.WriteInto(b)
	_, err := l.writerCloser.Writer(r).Write(*b)
	if err != nil {
		fmt.Fprintf(os.Stderr, "tealog: Write: error: %v; record:\n%s", err, b.StringUnsafe())
	}
	bufPool.Put(b)
}

var bufPool = bufferpool.New()

func (l *Logger[NR, WC, R]) selectClose() {
	for {
		select {
		case r := <-l.chanRecord:
			l.selectChanRecord(r)
		default:
			close(l.closed)
			return
		}
	}
}

func (l *Logger[NR, WC, R]) Close() {
	close(l.close)
	<-l.closed
	l.writerCloser.Close()
}

func (l *Logger[NR, WC, R]) Debug(ctx context.Context, msg string) {
	if !profile.Debug {
		return
	}
	l.log(ctx, LevelDebug, msg)
}

func (l *Logger[NR, WC, R]) Info(ctx context.Context, msg string) {
	l.log(ctx, LevelInfo, msg)
}

func (l *Logger[NR, WC, R]) Error(ctx context.Context, msg string) {
	l.log(ctx, LevelError, msg)
}

func (l *Logger[NR, WC, R]) log(ctx context.Context, level string, msg string) {
	l.chanRecord <- l.newRecord(ctx, level, msg)
}

const (
	LevelDebug = "DBG"
	LevelInfo  = "INF"
	LevelError = "ERR"
)
