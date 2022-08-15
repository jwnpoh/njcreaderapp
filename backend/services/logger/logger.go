package logger

import (
	"log"
	"net/http"
	"os"
	"time"

	fire "github.com/jwnpoh/njcreaderapp/backend/external/firestore"
	"github.com/jwnpoh/njcreaderapp/backend/services/serializer"
)

// Logger provides methods to interface with the app logger.
type Logger interface {
	Info(s serializer.Serializer, r *http.Request)
	Error(s serializer.Serializer, r *http.Request)
	Success(s serializer.Serializer, r *http.Request)
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
func (l *appLogger) Info(s serializer.Serializer, r *http.Request) {
	entry := logEntry{
		Date:     time.Now().Format("Jan 2, 2006 15:04:05"),
		LogLevel: "INFO",
		Method:   r.Method,
		URL:      r.URL.Path,
		Query:    r.URL.Query().Encode(),
		Message:  s,
	}

	l.info.Println(entry)
	l.db.Log(entry)
}

// Error logs an error level log entry.
func (l *appLogger) Error(s serializer.Serializer, r *http.Request) {
	entry := logEntry{
		Date:     time.Now().Format("Jan 2, 2006 15:04:05"),
		LogLevel: "ERROR",
		Method:   r.Method,
		URL:      r.URL.Path,
		Query:    r.URL.Query().Encode(),
		Message:  s,
	}

	l.err.Println(entry)
	l.db.Log(entry)
}

// Success logs a success level log entry.
func (l *appLogger) Success(s serializer.Serializer, r *http.Request) {
	entry := logEntry{
		Date:     time.Now().Format("Jan 2, 2006 15:04:05"),
		LogLevel: "SUCCESS",
		Method:   r.Method,
		URL:      r.URL.Path,
		Query:    r.URL.Query().Encode(),
		Message:  "",
	}

	l.success.Println(entry)
	l.db.Log(entry)
}
