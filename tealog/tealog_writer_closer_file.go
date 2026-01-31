package tealog

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type writerCloserFile struct {
	fileDir        string
	fileName       string
	file           *os.File
	fileRotateDate fileRotateDate
}

type fileRotateDate struct {
	year  int
	month time.Month
	day   int
}

func NewWriterCloserFile(fileDir string, fileName string) *writerCloserFile {
	f := &writerCloserFile{
		fileDir:        fileDir,
		fileName:       fileName,
		file:           nil,
		fileRotateDate: newFileRotateDate(),
	}
	return f
}

func newFileRotateDate() fileRotateDate {
	return fileRotateDate{0, 0, 0}
}

var _ writerCloser = (*writerCloserFile)(nil)

func (f *writerCloserFile) Writer(r Record) io.Writer {
	fileRotateDate := newFileRotateDate()
	fileRotateDate.year, fileRotateDate.month, fileRotateDate.day = r.Time.Date()
	if fileRotateDate != f.fileRotateDate {
		f.Close()
		filePath := filepath.Join(f.fileDir, r.Time.Format(time.DateOnly), f.fileName+".log")
		err := os.MkdirAll(filepath.Dir(filePath), 0o755)
		if err != nil {
			return newErrorWriter(err)
		}
		file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o644)
		if err != nil {
			return newErrorWriter(err)
		}
		f.file = file
		f.fileRotateDate = fileRotateDate
	}
	return f.file
}

func (f *writerCloserFile) Close() {
	if f.file != nil {
		err := f.file.Sync()
		if err != nil {
			fmt.Fprintf(os.Stderr, "tealog: file.Sync error: %v", err)
		}
		f.file.Close()
		f.file = nil
	}
}

const FileDir = "log"
