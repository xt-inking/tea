package bufferpool

import (
	"sync"
	"unsafe"
)

func New() *sync.Pool {
	pool := &sync.Pool{
		New: func() any {
			b := make([]byte, 0, 1024)
			return (*Buffer)(&b)
		},
	}
	return pool
}

type Buffer []byte

func NewBuffer(pool *sync.Pool) *Buffer {
	return pool.Get().(*Buffer)
}

func (b *Buffer) Free(pool *sync.Pool) {
	const maxBufferSize = 16 << 10
	if cap(*b) <= maxBufferSize {
		*b = (*b)[:0]
		pool.Put(b)
	}
}

func (b *Buffer) WriteByte(c byte) error {
	*b = append(*b, c)
	return nil
}

func (b *Buffer) WriteString(s string) (int, error) {
	*b = append(*b, s...)
	return len(s), nil
}

func (b *Buffer) String() string {
	return string(*b)
}

func (b *Buffer) StringUnsafe() string {
	return unsafe.String(unsafe.SliceData(*b), len(*b))
}
