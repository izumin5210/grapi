package grapiserver

import (
	"log"
	"os"
)

// LogFields represents logged fields.
type LogFields map[string]interface{}

// Logger is an interface for logging in a grapi server.
type Logger interface {
	Info(msg string, fields LogFields)
	Error(msg string, fields LogFields)
}

var (
	// DefaultLogger is a default logger implementation.
	DefaultLogger Logger = &defaultLogger{
		Logger: log.New(os.Stdout, "[grapi] ", log.LstdFlags),
	}
)

type defaultLogger struct {
	*log.Logger
}

func (l *defaultLogger) Info(msg string, fields LogFields) {
	l.Logger.Printf("%s: %#v", msg, fields)
}

func (l *defaultLogger) Error(msg string, fields LogFields) {
	l.Logger.Printf("%s: %#v", msg, fields)
}
