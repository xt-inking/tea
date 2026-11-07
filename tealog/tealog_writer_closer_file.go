package tealog

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

func WriterCloserFile[R Record](fileDir string, fileName string) *writerCloserFile[R] {
	return &writerCloserFile[R]{
		fileDir:        fileDir,
		fileName:       fileName,
		file:           nil,
		fileRotateDate: newFileRotateDate(),
	}
}

type writerCloserFile[R Record] struct {
	fileDir        string
	fileName       string
	file           *os.File
	fileRotateDate fileRotateDate
}

func newFileRotateDate() fileRotateDate {
	return fileRotateDate{0, 0, 0}
}

type fileRotateDate struct {
	year  int
	month time.Month
	day   int
}

func (wc *writerCloserFile[R]) Writer(r R) io.Writer {
	fileRotateDate := newFileRotateDate()
	fileRotateDate.year, fileRotateDate.month, fileRotateDate.day = r.Time().Date()
	if wc.fileRotateDate != fileRotateDate {
		path := filepath.Join(wc.fileDir, r.Time().Format(time.DateOnly), wc.fileName+".log")
		err := os.MkdirAll(filepath.Dir(path), 0o755)
		if err != nil {
			return NewErrorWriter(fmt.Errorf("tealog: MkdirAll: error: %w", err))
		}
		file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return NewErrorWriter(fmt.Errorf("tealog: OpenFile: error: %w", err))
		}
		wc.file = file
		wc.fileRotateDate = fileRotateDate
	}
	return wc.file
}

func (wc *writerCloserFile[R]) Close() {
	if wc.file != nil {
		err := wc.file.Sync()
		if err != nil {
			fmt.Fprintf(os.Stderr, "tealog: Sync: error: %v", err)
		}
		wc.file.Close()
		wc.file = nil
	}
}

const FileDir = "log"
