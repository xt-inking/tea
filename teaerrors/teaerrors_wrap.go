package teaerrors

import (
	"fmt"
)

func Warp(err error, msg string) *Error {
	e, ok := err.(*Error)
	if !ok {
		e = NewSkip(err, 1)
	}
	return e.Wrap(msg)
}

func Warpf(err error, msg string, a ...any) *Error {
	e, ok := err.(*Error)
	if !ok {
		e = NewSkip(err, 1)
	}
	return e.Warpf(msg, a...)
}

func (e *Error) Wrap(msg string) *Error {
	e.err = fmt.Errorf(msg+": %w", e.err)
	return e
}

func (e *Error) Warpf(msg string, a ...any) *Error {
	e.err = fmt.Errorf(msg+": %w", append(a, e.err)...)
	return e
}
