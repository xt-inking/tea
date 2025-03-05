package teaerrors

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

type errorStack struct {
	err error
}

func New(err error) *errorStack {
	e := &errorStack{
		err: err,
	}
	return e
}

func NewAny(e any) *errorStack {
	err, ok := e.(error)
	if !ok {
		err = fmt.Errorf("%s", e)
	}
	return New(err)
}

func (e *errorStack) ErrorStack(skip int) string {
	return e.Error() + "\n" + e.Stack(skip+1)
}

func (e *errorStack) Error() string {
	return e.err.Error()
}

func (e *errorStack) Stack(skip int) string {
	pc := make([]uintptr, 64)
	for {
		n := runtime.Callers(skip+2, pc)
		if n < len(pc) {
			pc = pc[:n]
			break
		}
		pc = make([]uintptr, len(pc)+64)
	}
	frames := runtime.CallersFrames(pc)
	var sb strings.Builder
	sb.WriteString("Stack:\n")
	for {
		frame, more := frames.Next()
		sb.WriteString(frame.Function)
		sb.WriteString("\n\t")
		sb.WriteString(frame.File)
		sb.WriteByte(':')
		sb.WriteString(strconv.Itoa(frame.Line))
		sb.WriteByte('\n')
		if !more {
			break
		}
	}
	return sb.String()
}
