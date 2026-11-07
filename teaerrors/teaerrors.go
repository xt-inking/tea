package teaerrors

import (
	"fmt"
	"runtime"
	"strconv"

	"github.com/xt-inking/tea/internal/bufferpool"
)

type Error struct {
	err   error
	stack []uintptr
}

func New(err error) *Error {
	return NewSkip(err, 1)
}

func NewSkip(err error, skip int) *Error {
	skip = skip + 2
	const size = 64
	stack := make([]uintptr, size)
	for {
		n := runtime.Callers(skip, stack)
		if n < len(stack) {
			stack = stack[:n]
			break
		}
		stack = make([]uintptr, len(stack)*2)
	}
	return &Error{err, stack}
}

func NewAny(e any) *Error {
	return NewAnySkip(e, 1)
}

func NewAnySkip(e any, skip int) *Error {
	err, ok := e.(error)
	if !ok {
		err = fmt.Errorf("%v", e)
	}
	return NewSkip(err, skip+1)
}

func (e *Error) Error() string {
	b := bufPool.Get()
	defer bufPool.Put(b)
	b.WriteString(e.err.Error())
	b.WriteString("\nStack:")
	frames := runtime.CallersFrames(e.stack)
	for i := 1; ; i++ {
		frame, more := frames.Next()
		b.WriteString("\n\t")
		b.WriteString(strconv.Itoa(i))
		b.WriteString(". ")
		b.WriteString(frame.Function)
		b.WriteString("\n\t\t")
		b.WriteString(frame.File)
		b.WriteByte(':')
		b.WriteString(strconv.Itoa(frame.Line))
		if !more {
			break
		}
	}
	return b.String()
}

var bufPool = bufferpool.New()
