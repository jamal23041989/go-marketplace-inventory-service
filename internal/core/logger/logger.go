package logger

import (
	"io"
	"log"
)

type Logger interface {
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})
}

type logger struct {
	out      io.Writer
	infoLog  *log.Logger
	warnLog  *log.Logger
	errorLog *log.Logger
	fatalLog *log.Logger
}

func New(out io.Writer) Logger {
	flags := log.Ldate | log.Ltime | log.Lshortfile

	infoLog := log.New(out, "INFO\t", flags)
	warnLog := log.New(out, "WARN\t", flags)
	errorLog := log.New(out, "ERROR\t", flags)
	fatalLog := log.New(out, "FATAL\t", flags)

	return &logger{
		out:      out,
		infoLog:  infoLog,
		warnLog:  warnLog,
		errorLog: errorLog,
		fatalLog: fatalLog,
	}
}

func (l logger) Info(msg string, args ...interface{}) {
	l.infoLog.Printf(msg, args...)
}

func (l logger) Warn(msg string, args ...interface{}) {
	l.warnLog.Printf(msg, args...)
}

func (l logger) Error(msg string, args ...interface{}) {
	l.errorLog.Printf(msg, args...)
}

func (l logger) Fatal(msg string, args ...interface{}) {
	l.fatalLog.Fatalf(msg, args...)
}
