package tealog

import (
	"io"
)

type WriterCloser[R Record] interface {
	Writer(r R) io.Writer
	Close()
}

func MultiWriterCloser[R Record](writerClosers ...WriterCloser[R]) WriterCloser[R] {
	allWriterClosers := make([]WriterCloser[R], 0, len(writerClosers))
	for _, wc := range writerClosers {
		if mwc, ok := wc.(*multiWriterCloser[R]); ok {
			allWriterClosers = append(allWriterClosers, mwc.writerClosers...)
		} else {
			allWriterClosers = append(allWriterClosers, wc)
		}
	}
	return &multiWriterCloser[R]{allWriterClosers}
}

type multiWriterCloser[R Record] struct {
	writerClosers []WriterCloser[R]
}

func (mwc *multiWriterCloser[R]) Writer(r R) io.Writer {
	writers := make([]io.Writer, 0, len(mwc.writerClosers))
	for _, wc := range mwc.writerClosers {
		writers = append(writers, wc.Writer(r))
	}
	return io.MultiWriter(writers...)
}

func (mwc *multiWriterCloser[R]) Close() {
	for _, wc := range mwc.writerClosers {
		wc.Close()
	}
}

func NewErrorWriter(err error) io.Writer {
	return &errorWriter{err}
}

type errorWriter struct {
	err error
}

func (w *errorWriter) Write(p []byte) (n int, err error) {
	return 0, w.err
}
