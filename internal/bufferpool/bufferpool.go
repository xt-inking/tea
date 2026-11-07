package bufferpool

import (
	"sync"
	"unsafe"
)

type BufferPool sync.Pool

func New() BufferPool {
	return BufferPool(sync.Pool{
		New: func() any {
			b := make([]byte, 0, 1024)
			return (*Buffer)(&b)
		},
	})
}

func (p *BufferPool) Get() *Buffer {
	return (*sync.Pool)(p).Get().(*Buffer)
}

func (p *BufferPool) Put(b *Buffer) {
	const maxBufferSize = 16 << 10
	if cap(*b) <= maxBufferSize {
		*b = (*b)[:0]
		(*sync.Pool)(p).Put(b)
	}
}

type Buffer []byte

func (b *Buffer) Write(p []byte) (n int, err error) {
	*b = append(*b, p...)
	return len(p), nil
}

func (b *Buffer) WriteByte(c byte) error {
	*b = append(*b, c)
	return nil
}

func (b *Buffer) WriteString(s string) (n int, err error) {
	*b = append(*b, s...)
	return len(s), nil
}

func (b *Buffer) String() string {
	return string(*b)
}

func (b *Buffer) StringUnsafe() string {
	return unsafe.String(unsafe.SliceData(*b), len(*b))
}
