package tealog

import (
	"time"

	"github.com/tea-frame-go/tea/internal/bufferpool"
)

type Record interface {
	WriteBuffer(buf *bufferpool.Buffer)
	Time() time.Time
}
