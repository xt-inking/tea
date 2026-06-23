package tealog

import (
	"context"
	"time"

	"github.com/tea-frame-go/tea/internal/bufferpool"
)

type recordText struct {
	time  time.Time
	level string
	msg   string
}

func NewRecordText(ctx context.Context, level string, msg string) Record {
	return &recordText{
		time:  time.Now(),
		level: level,
		msg:   msg,
	}
}

func (r *recordText) WriteBuffer(buf *bufferpool.Buffer) {
	const layout = "2006-01-02T15:04:05.000Z07:00"
	buf.WriteString(r.time.Format(layout))
	buf.WriteByte(' ')
	buf.WriteString(r.level)
	buf.WriteByte(' ')
	buf.WriteString(r.msg)
	buf.WriteByte('\n')
}

func (r *recordText) Time() time.Time {
	return r.time
}
