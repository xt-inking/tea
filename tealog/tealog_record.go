package tealog

import (
	"context"
	"time"
)

type Record struct {
	Ctx   context.Context
	Time  time.Time
	Level string
	Msg   string
}

func newRecord(ctx context.Context, level string, msg string) Record {
	return Record{
		Ctx:   ctx,
		Time:  time.Now(),
		Level: level,
		Msg:   msg,
	}
}
