package tealog

import (
	"context"
	"time"
)

type RecordText struct {
	time  time.Time
	level string
	msg   string
}

func NewRecordText(ctx context.Context, level string, msg string) RecordText {
	return RecordText{
		time:  time.Now(),
		level: level,
		msg:   msg,
	}
}

func (r RecordText) WriteInto(b *Buffer) {
	const layout = "2006-01-02T15:04:05.000Z07:00"
	*b = r.time.AppendFormat(*b, layout)
	b.WriteByte(' ')
	b.WriteString(r.level)
	b.WriteByte(' ')
	b.WriteString(r.msg)
	b.WriteByte('\n')
}

func (r RecordText) Time() time.Time {
	return r.time
}
