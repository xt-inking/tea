package teaerrors

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

type ErrorStack interface {
	error
	ErrorStack() string
}

type errorStack struct {
	err   error
	stack []uintptr
}

func New(err error, skip int) ErrorStack {
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

func NewAny(e any, skip int) ErrorStack {
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
	var sb strings.Builder
	sb.WriteString("Stack:")
	for {
		frame, more := frames.Next()
		sb.WriteString("\n\t")
		sb.WriteString(frame.Function)
		sb.WriteString("\n\t\t")
		sb.WriteString(frame.File)
		sb.WriteByte(':')
		sb.WriteString(strconv.Itoa(frame.Line))
		if !more {
			break
		}
	}
	return sb.String()
}

func (e *errorStack) Wrap(msg string) error {
	e.err = fmt.Errorf(msg+": %w", e.err)
	return e
}

func (e *errorStack) Wrapf(format string, a ...any) error {
	e.err = fmt.Errorf(format+": %w", append(a, e.err)...)
	return e
}
