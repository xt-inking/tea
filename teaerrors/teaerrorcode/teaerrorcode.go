package teaerrorcode

import (
	"fmt"
)

type ErrorCode interface {
	error
	ErrorCode() int
}

type errorCode struct {
	Code int
	Msg  string
}

func New(code int, msg string) ErrorCode {
	e := &errorCode{code, msg}
	return e
}

func Newf(code int, format string, a ...any) ErrorCode {
	e := &errorCode{code, fmt.Sprintf(format, a...)}
	return e
}

func (e *errorCode) Error() string {
	return e.Msg
}

func (e *errorCode) ErrorCode() int {
	return e.Code
}
