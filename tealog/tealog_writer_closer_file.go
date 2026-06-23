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
	fileOld        *os.File
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
		fileOld:        nil,
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
	fileRotateDate.year, fileRotateDate.month, fileRotateDate.day = r.Time().Date()
	switch f.fileRotateDate.compare(fileRotateDate) {
	case -1:
		f.closeOld()
		f.fileOld = f.file
		f.file = nil
		filePath := filepath.Join(f.fileDir, r.Time().Format(time.DateOnly), f.fileName+".log")
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
		return f.file
	case +1:
		return f.fileOld
	default:
		return f.file
	}
}

func (f *writerCloserFile) Close() {
	f.closeOld()
	if f.file != nil {
		f.close(f.file)
		f.file = nil
	}
}

func (f *writerCloserFile) closeOld() {
	if f.fileOld != nil {
		f.close(f.fileOld)
		f.fileOld = nil
	}
}

func (f *writerCloserFile) close(file *os.File) {
	err := file.Sync()
	if err != nil {
		fmt.Fprintf(os.Stderr, "tealog: file.Sync error: %v", err)
	}
	file.Close()
}

func (fileRotateDate fileRotateDate) compare(otherFileRotateDate fileRotateDate) int {
	if fileRotateDate.year < otherFileRotateDate.year {
		return -1
	}
	if fileRotateDate.year > otherFileRotateDate.year {
		return +1
	}
	if fileRotateDate.month < otherFileRotateDate.month {
		return -1
	}
	if fileRotateDate.month > otherFileRotateDate.month {
		return +1
	}
	if fileRotateDate.day < otherFileRotateDate.day {
		return -1
	}
	if fileRotateDate.day > otherFileRotateDate.day {
		return +1
	}
	return 0
}

const FileDir = "log"
