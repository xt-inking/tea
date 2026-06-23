package teaerrors

import (
	"fmt"
)

func Wrap(err error, msg string) error {
	if e, ok := err.(interface {
		Wrap(msg string) error
	}); ok {
		return e.Wrap(msg)
	}
	return fmt.Errorf(msg+": %w", err)
}

func Wrapf(err error, format string, a ...any) error {
	if e, ok := err.(interface {
		Wrapf(format string, a ...any) error
	}); ok {
		return e.Wrapf(format, a...)
	}
	return fmt.Errorf(format+": %w", append(a, err)...)
}

func (e *errorStack) Wrap(msg string) error {
	e.err = fmt.Errorf(msg+": %w", e.err)
	return e
}

func (e *errorStack) Wrapf(format string, a ...any) error {
	e.err = fmt.Errorf(format+": %w", append(a, e.err)...)
	return e
}
