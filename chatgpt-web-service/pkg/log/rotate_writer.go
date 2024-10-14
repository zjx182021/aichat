package log

import (
	"io"
	"os"
	"sync"

	"gopkg.in/natefinch/lumberjack.v2"
)

type fileRotateWriter struct {
	data map[string]io.Writer
	sync.RWMutex
}

func (frw *fileRotateWriter) Getwriter(logpath string) io.Writer {
	frw.RLock()
	defer frw.RUnlock()
	w, ok := frw.data[logpath]
	if !ok {
		return nil
	}
	return w
}

func (frw *fileRotateWriter) SetWriter(logpath string, w io.Writer) error {
	frw.Lock()
	defer frw.Unlock()
	frw.data[logpath] = w
	_, err := os.Stat(logpath)
	if os.IsNotExist(err) {
		os.Create(logpath)
	}
	return nil
}

var _fileRotateWriter *fileRotateWriter

func init() {
	_fileRotateWriter = &fileRotateWriter{
		data: map[string]io.Writer{},
	}
}

func GetRotateWriter(logpath string) io.Writer {
	if logpath == "" {
		My_log.Error("logpath 为空")
		return nil
	}
	writer := _fileRotateWriter.Getwriter(logpath)
	if writer != nil {
		return writer
	}
	writer = &lumberjack.Logger{
		Filename:   logpath,
		MaxSize:    10, // MB
		MaxBackups: 15,
		MaxAge:     7, // days
		Compress:   false,
		LocalTime:  true,
	}
	_fileRotateWriter.SetWriter(logpath, writer)
	return writer
}
