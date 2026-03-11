package teaerrors

import (
	"fmt"
	"runtime"
	"strconv"

	"github.com/tea-frame-go/tea/internal/bufferpool"
)

type Error interface {
	error
	ErrorStack() string
}

type errorStack struct {
	err   error
	stack []uintptr
}

func New(err error, skip int) Error {
	pc := make([]uintptr, 64)
	for {
		n := runtime.Callers(skip+2, pc)
		if n < len(pc) {
			pc = pc[:n]
			break
		}
		pc = make([]uintptr, len(pc)+64)
	}
	e := &errorStack{
		err:   err,
		stack: pc,
	}
	return e
}

func NewAny(e any, skip int) Error {
	err, ok := e.(error)
	if !ok {
		err = fmt.Errorf("%v", e)
	}
	return New(err, skip+1)
}

func (e *errorStack) Error() string {
	return e.err.Error()
}

func (e *errorStack) ErrorStack() string {
	return e.Error() + "\n" + e.Stack()
}

func (e *errorStack) Stack() string {
	frames := runtime.CallersFrames(e.stack)
	buf := bufferpool.NewBuffer(bufPool)
	defer buf.Free(bufPool)
	buf.WriteString("Stack:")
	for {
		frame, more := frames.Next()
		buf.WriteString("\n\t")
		buf.WriteString(frame.Function)
		buf.WriteString("\n\t\t")
		buf.WriteString(frame.File)
		buf.WriteByte(':')
		buf.WriteString(strconv.Itoa(frame.Line))
		if !more {
			break
		}
	}
	return buf.String()
}

var bufPool = bufferpool.New()
