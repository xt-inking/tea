package tealog

import (
	"context"
	"time"

	"github.com/xt-inking/tea/internal/bufferpool"
)

type Record interface {
	WriteInto(b *Buffer)
	Time() time.Time
}

type Buffer = bufferpool.Buffer

type NewRecord[R Record] = func(ctx context.Context, level string, msg string) R
