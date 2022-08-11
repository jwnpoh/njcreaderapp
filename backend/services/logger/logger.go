package logger

import (
	"log"
	"os"
	"time"

	fire "github.com/jwnpoh/njcreaderapp/backend/external/firestore"
)

// Logger provides methods to interface with the app logger.
type Logger interface {
	Info(method, urlPath, query string, message interface{})
	Error(method, urlPath, query string, message interface{})
	Success(method, urlPath, query string, message interface{})
}

type appLogger struct {
	info    *log.Logger
	err     *log.Logger
	success *log.Logger
	db      *fire.FireStoreRepo
}

type logEntry struct {
	Date     string      `firestore:"date"`
	LogLevel string      `firestore:"logLevel"`
	Method   string      `firestore:"method"`
	URL      string      `firestore:"url"`
	Query    string      `firestore:"query"`
	Message  interface{} `firestore:"message"`
}

// NewAppLogger returns an appLogger struct of aggregated loggers for logging.
func NewAppLogger() Logger {
	flags := log.LstdFlags

	infoLogger := log.New(os.Stdout, "INFO: ", flags)
	errorLogger := log.New(os.Stdout, "ERROR: ", flags)
	successLogger := log.New(os.Stdout, "SUCCESS: ", flags)

	return &appLogger{
		info:    infoLogger,
		err:     errorLogger,
		success: successLogger,
		db:      fire.NewFireStoreRepo(),
	}
}

// Info logs an info level log entry.
func (l *appLogger) Info(method, urlPath, query string, message interface{}) {
	entry := logEntry{
		Date:     time.Now().Format("Jan 2, 2006 15:04:05"),
		LogLevel: "SUCCESS",
		Method:   method,
		URL:      urlPath,
		Query:    query,
		Message:  message,
	}

	l.info.Println(entry)
	l.db.Log(entry)
}

// Error logs an error level log entry.
func (l *appLogger) Error(method, urlPath, query string, message interface{}) {
	entry := logEntry{
		Date:     time.Now().Format("Jan 2, 2006 15:04:05"),
		LogLevel: "SUCCESS",
		Method:   method,
		URL:      urlPath,
		Query:    query,
		Message:  message,
	}

	l.err.Println(entry)
	l.db.Log(entry)
}

// Success logs a success level log entry.
func (l *appLogger) Success(method, urlPath, query string, message interface{}) {
	entry := logEntry{
		Date:     time.Now().Format("Jan 2, 2006 15:04:05"),
		LogLevel: "SUCCESS",
		Method:   method,
		URL:      urlPath,
		Query:    query,
		Message:  message,
	}

	l.success.Println(entry)
	l.db.Log(entry)
}
