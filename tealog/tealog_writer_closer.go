package tealog

import (
	"io"
)

type writerCloser interface {
	Writer(r Record) io.Writer
	Close()
}

func MultiWriterCloser(writerClosers ...writerCloser) writerCloser {
	allWriterClosers := make([]writerCloser, 0, len(writerClosers))
	for _, wc := range writerClosers {
		if mwc, ok := wc.(*multiWriterCloser); ok {
			allWriterClosers = append(allWriterClosers, mwc.writerClosers...)
		} else {
			allWriterClosers = append(allWriterClosers, wc)
		}
	}
	return &multiWriterCloser{allWriterClosers}
}

type multiWriterCloser struct {
	writerClosers []writerCloser
}

func (mwc *multiWriterCloser) Writer(r Record) io.Writer {
	writers := make([]io.Writer, 0, len(mwc.writerClosers))
	for _, wc := range mwc.writerClosers {
		writers = append(writers, wc.Writer(r))
	}
	return io.MultiWriter(writers...)
}

func (mwc *multiWriterCloser) Close() {
	for _, wc := range mwc.writerClosers {
		wc.Close()
	}
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
