package tealog

import (
	"sync"
)

var bufPool = sync.Pool{
	New: func() any {
		b := make([]byte, 0, 1024)
		return (*Buffer)(&b)
	},
}

type Buffer []byte

func newBuffer() *Buffer {
	return bufPool.Get().(*Buffer)
}

func (b *Buffer) free() {
	const maxBufferSize = 16 << 10
	if cap(*b) <= maxBufferSize {
		*b = (*b)[:0]
		bufPool.Put(b)
	}
}

func (b *Buffer) BufferWriteString(s string) {
	*b = append(*b, s...)
}

func (b *Buffer) BufferWriteByte(c byte) {
	*b = append(*b, c)
}
