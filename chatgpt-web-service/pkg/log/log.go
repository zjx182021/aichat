package log

import (
	"errors"
	"io"
	"runtime"

	"github.com/sirupsen/logrus"
)

var (
	Info    = "info"
	Error   = "error"
	Fatal   = "fatal"
	Panic   = "panic"
	Trace   = "trace"
	Debug   = "debug"
	Warning = "warning"
)

type ILogger interface {
	SetLevel(lvl string)
	SetOutput(writer io.Writer)
	SetPrintCaller(bool)
	SetCaller(caller func() (file string, line int, funcName string, err error))
	Trace(args ...interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})
	TraceF(format string, args ...interface{})
	DebugF(format string, args ...interface{})
	InfoF(format string, args ...interface{})
	WarningF(format string, args ...interface{})
	ErrorF(format string, args ...interface{})
	FatalF(format string, args ...interface{})
	PanicF(format string, args ...interface{})
	WithFields(fields map[string]interface{}) ILogger
}

var My_log *Logger

type Logger struct {
	*logrus.Entry
	printCaller bool
	caller      func() (file string, line int, funcName string, err error)
}

func (l *Logger) SetLevel(lvl string) {
	if lvl == "" {
		return
	}
	level, err := logrus.ParseLevel(lvl)
	if err == nil {
		l.Entry.Logger.Level = level
	}
}

func NewLogger() *Logger {
	if My_log == nil {

		mylog := &Logger{}
		mylog.Entry = logrus.NewEntry(logrus.New())
		mylog.SetLevel(Info)
		mylog.Entry.Logger.AddHook(&errorhook{})
		mylog.SetCaller(Defaultcaller)
		My_log = mylog
	}
	return My_log

}

func Defaultcaller() (string, int, string, error) {
	pc, file, line, ok := runtime.Caller(4)
	if !ok {
		return "", 0, "", errors.New("set defaultcaller failure")
	}
	funcname := runtime.FuncForPC(pc).Name()
	return file, line, funcname, nil
}
func (l *Logger) SetOutput(writer io.Writer) {
	l.Logger.SetOutput(writer)
}
func (l *Logger) SetPrintCaller(printCaller bool) {
	l.printCaller = printCaller
}
func (l *Logger) SetCaller(caller func() (file string, line int, funcName string, err error)) {
	l.caller = caller
}
func (l *Logger) Trace(args ...interface{}) {
	l.Log(logrus.TraceLevel, args...)
}
func (l *Logger) Debug(args ...interface{}) {
	l.Log(logrus.DebugLevel, args...)
}
func (l *Logger) Info(args ...interface{}) {
	l.Log(logrus.InfoLevel, args...)
}
func (l *Logger) Warning(args ...interface{}) {
	l.Log(logrus.WarnLevel, args...)
}
func (l *Logger) Error(args ...interface{}) {
	l.Log(logrus.ErrorLevel, args...)
}
func (l *Logger) Fatal(args ...interface{}) {
	l.Log(logrus.FatalLevel, args...)
}
func (l *Logger) Panic(args ...interface{}) {
	l.Log(logrus.PanicLevel, args...)
}
func (l *Logger) TraceF(format string, args ...interface{}) {
	l.Logf(logrus.TraceLevel, format, args...)
}
func (l *Logger) DebugF(format string, args ...interface{}) {
	l.Logf(logrus.DebugLevel, format, args...)
}
func (l *Logger) InfoF(format string, args ...interface{}) {
	l.Logf(logrus.InfoLevel, format, args...)
}
func (l *Logger) WarningF(format string, args ...interface{}) {
	l.Logf(logrus.WarnLevel, format, args...)
}
func (l *Logger) ErrorF(format string, args ...interface{}) {
	l.Logf(logrus.ErrorLevel, format, args...)
}
func (l *Logger) FatalF(format string, args ...interface{}) {
	l.Logf(logrus.FatalLevel, format, args...)
}
func (l *Logger) PanicF(format string, args ...interface{}) {
	l.Logf(logrus.PanicLevel, format, args...)
}
func (l *Logger) WithFields(fields map[string]interface{}) ILogger {
	return &Logger{
		Entry:       l.Entry.WithFields(fields),
		printCaller: l.printCaller,
		caller:      l.caller,
	}
}
