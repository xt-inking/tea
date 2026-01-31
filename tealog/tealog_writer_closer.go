package tealog

import (
	"io"
)

type writerCloser interface {
	Writer(r Record) io.Writer
	Close()
}

type errorWriter struct {
	err error
}

func newErrorWriter(err error) errorWriter {
	return errorWriter{
		err: err,
	}
}

func (e errorWriter) Write(p []byte) (n int, err error) {
	return 0, e.err
}
