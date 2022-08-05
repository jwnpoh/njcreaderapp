package logger

import (
	"log"
	"os"
)

const (
	Success = "SUCCESS"
)

type Logger interface {
	Info(v ...interface{})
	Error(v ...interface{})
	Success(v ...interface{})
}

type appLogger struct {
	info    *log.Logger
	err     *log.Logger
	success *log.Logger
}

func NewAppLogger() Logger {
	flags := log.LstdFlags
	infoLogger := log.New(os.Stdout, "INFO: ", flags)
	errorLogger := log.New(os.Stdout, "ERROR: ", flags)
	successLogger := log.New(os.Stdout, "SUCCESS: ", flags)
	return &appLogger{
		info:    infoLogger,
		err:     errorLogger,
		success: successLogger,
	}
}

func (l *appLogger) Info(v ...interface{}) {
	l.info.Println(v...)
}

func (l *appLogger) Error(v ...interface{}) {
	l.err.Println(v...)
}

func (l *appLogger) Success(v ...interface{}) {
	l.success.Println(v...)
}
