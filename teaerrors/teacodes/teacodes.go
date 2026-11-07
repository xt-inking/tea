package teacodes

import (
	"fmt"
)

type Error struct {
	Code string
	Msg  string
}

func New(code string, msg string) *Error {
	return &Error{code, msg}
}

func Newf(code string, msg string, a ...any) *Error {
	return &Error{code, fmt.Sprintf(msg, a...)}
}

func (e *Error) Error() string {
	return e.Code + ": " + e.Msg
}
